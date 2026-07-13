package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"time"

	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
	"herbst-server/db/characterskill"
	"herbst-server/db/skill"
	"herbst-server/events"
	"herbst-server/repository"
)

// DefaultSkillXPPerLevel is the fallback XP required per skill level when no curve data exists.
const DefaultSkillXPPerLevel = 100

// DefaultSkillMaxLevel is the fallback max skill level when not set on the skill.
const DefaultSkillMaxLevel = 100

// skillXPService implements SkillXPService (interface defined in interface.go).
type skillXPService struct {
	client          *db.Client
	charRepo        repository.CharacterRepo
	charAbilityRepo repository.CharacterAbilityRepo
	abilitySvc      AbilityService
	logger          *slog.Logger
}

// NewSkillXPService creates a new SkillXPService.
func NewSkillXPService(client *db.Client, charRepo repository.CharacterRepo, charAbilityRepo repository.CharacterAbilityRepo, abilitySvc AbilityService, logger *slog.Logger) SkillXPService {
	if logger == nil {
		logger = slog.Default()
	}
	return &skillXPService{
		client:          client,
		charRepo:        charRepo,
		charAbilityRepo: charAbilityRepo,
		abilitySvc:      abilitySvc,
		logger:          logger,
	}
}

// AwardSkillXP adds XP to a character's skill and handles skill level-up.
// Returns the new XP, new level, whether a level-up occurred, and any error.
func (s *skillXPService) AwardSkillXP(ctx context.Context, characterID int, skillName string, xpGained int, source string) (newXP, newLevel int, leveledUp bool, err error) {
	if xpGained <= 0 {
		return 0, 0, false, nil
	}

	// 1. Get the character to find their world
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, false, fmt.Errorf("get character: %w", err)
	}

	// 2. Find the skill by name in the character's world
	skillObj, err := s.client.Skill.Query().
		Where(skill.NameEQ(skillName), skill.WorldIDEQ(char.WorldID)).
		Only(ctx)
	if err != nil {
		return 0, 0, false, fmt.Errorf("find skill %q in world %d: %w", skillName, char.WorldID, err)
	}

	// 3. Find or create the character_skills record
	charSkill, err := s.client.CharacterSkill.Query().
		Where(characterskill.HasCharacterWith(character.IDEQ(characterID))).
		Where(characterskill.HasSkillWith(skill.IDEQ(skillObj.ID))).
		Only(ctx)
	if err != nil {
		// Record doesn't exist — create it at level 1, 0 XP
		charSkill, err = s.client.CharacterSkill.Create().
			SetCharacterID(characterID).
			SetSkillID(skillObj.ID).
			SetLevel(1).
			SetXp(0).
			Save(ctx)
		if err != nil {
			return 0, 0, false, fmt.Errorf("create character_skill: %w", err)
		}
	}

	// 4. Add XP
	newXP = charSkill.Xp + xpGained
	newLevel = charSkill.Level

	// 5. Determine max level
	maxLevel := skillObj.MaxLevel
	if maxLevel == 0 {
		maxLevel = DefaultSkillMaxLevel
	}

	// 6. Check for skill level-up using the skill's xp_curve_mode
	for {
		if newLevel >= maxLevel {
			break
		}
		needed := s.xpForNextLevel(skillObj, newLevel)
		if needed == 0 || newXP < needed {
			break
		}
		newLevel++
		s.logger.Info("skill level up",
			"character_id", characterID,
			"skill_name", skillName,
			"new_level", newLevel,
			"xp_remaining", newXP,
		)
	}

	// 7. Update the character_skill record
	_, err = s.client.CharacterSkill.UpdateOne(charSkill).
		SetXp(newXP).
		SetLevel(newLevel).
		Save(ctx)
	if err != nil {
		return 0, 0, false, fmt.Errorf("update character_skill: %w", err)
	}

	leveledUp = newLevel > charSkill.Level

	// 8. Emit skill.leveled_up event and check ability unlocks
	if leveledUp {
		events.Publish(events.Event{
			Type: events.EventSkillLeveledUp,
			Payload: map[string]interface{}{
				"character_id": characterID,
				"skill_name":   skillName,
				"skill_id":     skillObj.ID,
				"new_level":    newLevel,
				"new_xp":      newXP,
				"source":       source,
			},
			Timestamp: time.Now().UnixMilli(),
		})

		// 9. After level-up, check for ability unlocks
		unlockedAbilityIDs, unlockErr := s.CheckAbilityUnlocks(ctx, characterID, skillName, newLevel)
		if unlockErr != nil {
			s.logger.Error("failed to check ability unlocks after skill level-up",
				"character_id", characterID,
				"skill_name", skillName,
				"new_level", newLevel,
				"error", unlockErr,
			)
		} else {
			for _, abilityID := range unlockedAbilityIDs {
				events.Publish(events.Event{
					Type: events.EventSkillAbilityUnlocked,
					Payload: map[string]interface{}{
						"character_id": characterID,
						"skill_name":   skillName,
						"ability_id":   abilityID,
						"skill_level":  newLevel,
					},
					Timestamp: time.Now().UnixMilli(),
				})
				s.logger.Info("ability unlocked via skill level-up",
					"character_id", characterID,
					"skill_name", skillName,
					"ability_id", abilityID,
					"skill_level", newLevel,
				)

				// Auto-equip passive abilities
				ab, abErr := s.client.Ability.Get(ctx, abilityID)
				if abErr != nil {
					s.logger.Error("failed to get ability for auto-equip check",
						"ability_id", abilityID,
						"error", abErr,
					)
					continue
				}
				if ab.AbilityClass == "passive" && s.abilitySvc != nil {
					if _, eqErr := s.abilitySvc.UnlockPassiveAbility(ctx, characterID, abilityID); eqErr != nil {
						s.logger.Error("failed to auto-equip passive ability after skill level-up",
							"character_id", characterID,
							"ability_id", abilityID,
							"error", eqErr,
						)
					} else {
						s.logger.Info("auto-equipped passive ability after skill level-up",
							"character_id", characterID,
							"ability_id", abilityID,
						)
					}
				}
			}
		}
	}

	// Emit skill_xp.gained event
	events.Publish(events.Event{
		Type: events.EventSkillXPGained,
		Payload: map[string]interface{}{
			"character_id": characterID,
			"skill_name":   skillName,
			"skill_id":     skillObj.ID,
			"amount":       xpGained,
			"new_xp":       newXP,
			"new_level":    newLevel,
			"source":       source,
		},
		Timestamp: time.Now().UnixMilli(),
	})

	return newXP, newLevel, leveledUp, nil
}

// GetSkillLevel returns the current level of a character's skill.
// Returns 0 if the character doesn't have the skill.
func (s *skillXPService) GetSkillLevel(ctx context.Context, characterID int, skillName string) (level int, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, fmt.Errorf("get character: %w", err)
	}

	skillObj, err := s.client.Skill.Query().
		Where(skill.NameEQ(skillName), skill.WorldIDEQ(char.WorldID)).
		Only(ctx)
	if err != nil {
		return 0, fmt.Errorf("find skill %q: %w", skillName, err)
	}

	charSkill, err := s.client.CharacterSkill.Query().
		Where(characterskill.HasCharacterWith(character.IDEQ(characterID))).
		Where(characterskill.HasSkillWith(skill.IDEQ(skillObj.ID))).
		Only(ctx)
	if err != nil {
		return 0, nil // No record = level 0
	}
	return charSkill.Level, nil
}

// CheckAbilityUnlocks queries abilities that require this skill and became available at newLevel.
// Returns a list of ability IDs that are newly unlockable (not already equipped).
func (s *skillXPService) CheckAbilityUnlocks(ctx context.Context, characterID int, skillName string, newLevel int) ([]int, error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}

	skillObj, err := s.client.Skill.Query().
		Where(skill.NameEQ(skillName), skill.WorldIDEQ(char.WorldID)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("find skill %q: %w", skillName, err)
	}

	return s.checkAbilityUnlocksForSkill(ctx, characterID, skillObj.ID, newLevel)
}

// checkAbilityUnlocksForSkill does the actual work of checking ability unlocks for a given skill ID and level.
func (s *skillXPService) checkAbilityUnlocksForSkill(ctx context.Context, characterID int, skillID int, newLevel int) ([]int, error) {
	// Query abilities where required_skill_id matches this skill
	// and required_skill_level <= newLevel
	abilities, err := s.client.Ability.Query().
		Where(ability.RequiredSkillIDEQ(skillID)).
		Where(ability.RequiredSkillLevelLTE(newLevel)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query abilities for skill_id %d: %w", skillID, err)
	}

	var unlockedIDs []int
	for _, ab := range abilities {
		// Check if the character already has this ability equipped
		exists, err := s.client.CharacterAbility.Query().
			Where(
				characterability.HasCharacterWith(character.ID(characterID)),
				characterability.HasAbilityWith(ability.ID(ab.ID)),
			).
			Exist(ctx)
		if err != nil {
			s.logger.Error("failed to check if ability is already equipped",
				"character_id", characterID,
				"ability_id", ab.ID,
				"error", err,
			)
			continue
		}
		if !exists {
			unlockedIDs = append(unlockedIDs, ab.ID)
		}
	}

	return unlockedIDs, nil
}

// xpForNextLevel computes the cumulative XP threshold required to advance from `currentLevel` to `currentLevel+1`.
// Uses the skill's xp_curve_mode and xp_curve_data:
//   - "percentage": base_xp * (1 + percentage/100)^(n-1)
//   - "hand_coded": explicit thresholds array from xp_curve_data
//   - Default (no curve data): linear DefaultSkillXPPerLevel per level
func (s *skillXPService) xpForNextLevel(skillObj *db.Skill, currentLevel int) int {
	mode := skillObj.XpCurveMode
	data := skillObj.XpCurveData

	switch mode {
	case "percentage":
		return s.xpForNextLevelPercentage(skillObj, currentLevel, data)
	case "hand_coded":
		return s.xpForNextLevelHandCoded(currentLevel, data)
	default:
		// No curve mode set — use linear default
		return DefaultSkillXPPerLevel * currentLevel
	}
}

// xpForNextLevelPercentage computes threshold using exponential growth:
// xp_for_level(n) = base_xp * (1 + percentage/100)^(n-1)
// This is the XP needed to go from level n to level n+1 (cumulative is tracked by the caller).
func (s *skillXPService) xpForNextLevelPercentage(skillObj *db.Skill, currentLevel int, data map[string]interface{}) int {
	baseXP := DefaultSkillXPPerLevel
	percentage := 50.0 // default 50% growth

	if data != nil {
		if v, ok := data["base_xp"].(float64); ok && v > 0 {
			baseXP = int(v)
		}
		if v, ok := data["percentage"].(float64); ok && v > 0 {
			percentage = v
		}
	}

	// XP needed to advance from currentLevel to currentLevel+1
	return int(math.Round(float64(baseXP) * math.Pow(1+percentage/100, float64(currentLevel-1))))
}

// xpForNextLevelHandCoded uses explicit thresholds from the curve data.
// The thresholds array is 0-indexed: thresholds[0] is XP for level 1→2, thresholds[1] for level 2→3, etc.
func (s *skillXPService) xpForNextLevelHandCoded(currentLevel int, data map[string]interface{}) int {
	if data == nil {
		return DefaultSkillXPPerLevel * currentLevel
	}

	thresholdsRaw, ok := data["thresholds"]
	if !ok {
		return DefaultSkillXPPerLevel * currentLevel
	}

	// thresholds can be []interface{} (from JSON) or []float64
	var thresholds []int
	switch tt := thresholdsRaw.(type) {
	case []interface{}:
		for _, t := range tt {
			if tf, ok := t.(float64); ok {
				thresholds = append(thresholds, int(tf))
			}
		}
	case []float64:
		for _, tf := range tt {
			thresholds = append(thresholds, int(tf))
		}
	case []int:
		thresholds = tt
	default:
		// Try JSON round-trip as fallback
		b, _ := json.Marshal(thresholdsRaw)
		_ = json.Unmarshal(b, &thresholds)
	}

	// thresholds[0] = XP for level 1→2, thresholds[1] = XP for level 2→3, etc.
	// So for currentLevel → currentLevel+1, we need thresholds[currentLevel-1]
	idx := currentLevel - 1
	if idx < 0 {
		return 0
	}
	if idx >= len(thresholds) {
		// Beyond defined thresholds — use the last threshold or fall back to linear
		if len(thresholds) > 0 {
			return thresholds[len(thresholds)-1]
		}
		return DefaultSkillXPPerLevel * currentLevel
	}
	return thresholds[idx]
}

// AwardSkillXPForAbility is a convenience method that awards skill XP to the ability's required skill.
// If the ability has a required_skill_id, it awards XP to that skill.
// If not, it tries to match the ability to a skill by name/category.
// Returns the skill name awarded to (empty string if no XP was awarded).
func (s *skillXPService) AwardSkillXPForAbility(ctx context.Context, characterID int, abilityID int, xpGained int, source string) (skillName string, err error) {
	ab, err := s.client.Ability.Get(ctx, abilityID)
	if err != nil {
		return "", fmt.Errorf("get ability %d: %w", abilityID, err)
	}

	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return "", fmt.Errorf("get character: %w", err)
	}

	// If the ability has a required_skill_id, award XP to that skill
	if ab.RequiredSkillID != nil && *ab.RequiredSkillID > 0 {
		skillObj, err := s.client.Skill.Get(ctx, *ab.RequiredSkillID)
		if err != nil {
			return "", fmt.Errorf("get skill %d: %w", *ab.RequiredSkillID, err)
		}
		_, _, _, err = s.AwardSkillXP(ctx, characterID, skillObj.Name, xpGained, source)
		if err != nil {
			return skillObj.Name, err
		}
		return skillObj.Name, nil
	}

	// Try to match by ability type → skill category
	if ab.AbilityType != "" {
		skillObj, err := s.client.Skill.Query().
			Where(skill.CategoryEQ(ab.AbilityType), skill.WorldIDEQ(char.WorldID)).
			First(ctx)
		if err == nil && skillObj != nil {
			_, _, _, err = s.AwardSkillXP(ctx, characterID, skillObj.Name, xpGained, source)
			if err != nil {
				return skillObj.Name, err
			}
			return skillObj.Name, nil
		}
	}

	return "", nil // No matching skill found
}

// ComputeDiminishingXP returns a diminished XP amount based on the character's current skill level.
// Higher levels receive less XP per use to prevent grinding.
// Formula: xp = baseXP * (1 - min(0.5, level*0.005))
func ComputeDiminishingSkillXP(baseXP int, currentLevel int) int {
	if currentLevel <= 1 {
		return baseXP
	}
	diminishFactor := 1.0 - math.Min(0.5, float64(currentLevel)*0.005)
	result := int(math.Round(float64(baseXP) * diminishFactor))
	if result < 1 {
		result = 1
	}
	return result
}