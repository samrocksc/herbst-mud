package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
	"herbst-server/db/characterclasshistory"
	"herbst-server/db/characterfaction"
	"herbst-server/db/characterracehistory"
	"herbst-server/db/characterskill"
	"herbst-server/db/faction"
	"herbst-server/db/race"
	"herbst-server/db/world"
	"herbst-server/events"
)

// reclassReraceService implements ReclassReraceService.
type reclassReraceService struct {
	client *db.Client
	logger *slog.Logger
}

// NewReclassReraceService creates a new ReclassReraceService.
func NewReclassReraceService(client *db.Client, logger *slog.Logger) ReclassReraceService {
	if logger == nil {
		logger = slog.Default()
	}
	return &reclassReraceService{client: client, logger: logger}
}

// reclassConfig holds the reclass section of world config.
type reclassConfig struct {
	Allowed          bool    `json:"allowed"`
	Cost             int     `json:"cost"`
	MinLevel         int     `json:"min_level"`
	CooldownSeconds  int     `json:"cooldown_seconds"`
	SkillRetention   float64 `json:"skill_retention"`
}

// reraceConfig holds the rerace section of world config.
type reraceConfig struct {
	Allowed bool `json:"allowed"`
	Cost    int  `json:"cost"`
}

// loadReclassConfig loads the reclass config from the character's world.
func (s *reclassReraceService) loadReclassConfig(ctx context.Context, char *db.Character) (*reclassConfig, error) {
	if char.WorldID == 0 {
		return nil, fmt.Errorf("character has no world")
	}

	w, err := s.client.World.Query().Where(world.IDEQ(char.WorldID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("load world for reclass config: %w", err)
	}

	if w.Config == nil || len(w.Config) == 0 {
		return nil, fmt.Errorf("world has no config, reclass not configured")
	}

	raw, ok := w.Config["reclass"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("reclass not configured in world config")
	}

	cfg := &reclassConfig{}
	if v, ok := raw["allowed"].(bool); ok {
		cfg.Allowed = v
	}
	if v, ok := raw["cost"].(float64); ok {
		cfg.Cost = int(v)
	}
	if v, ok := raw["min_level"].(float64); ok {
		cfg.MinLevel = int(v)
	}
	if v, ok := raw["cooldown_seconds"].(float64); ok {
		cfg.CooldownSeconds = int(v)
	}
	if v, ok := raw["skill_retention"].(float64); ok {
		cfg.SkillRetention = v
	}

	return cfg, nil
}

// loadReraceConfig loads the rerace config from the character's world.
func (s *reclassReraceService) loadReraceConfig(ctx context.Context, char *db.Character) (*reraceConfig, error) {
	if char.WorldID == 0 {
		return nil, fmt.Errorf("character has no world")
	}

	w, err := s.client.World.Query().Where(world.IDEQ(char.WorldID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("load world for rerace config: %w", err)
	}

	if w.Config == nil || len(w.Config) == 0 {
		return nil, fmt.Errorf("world has no config, rerace not configured")
	}

	raw, ok := w.Config["rerace"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("rerace not configured in world config")
	}

	cfg := &reraceConfig{}
	if v, ok := raw["allowed"].(bool); ok {
		cfg.Allowed = v
	}
	if v, ok := raw["cost"].(float64); ok {
		cfg.Cost = int(v)
	}

	return cfg, nil
}

// Reclass switches a character's class (faction) with skill retention and history tracking.
func (s *reclassReraceService) Reclass(ctx context.Context, characterID int, newFactionID int) error {
	// 1. Load character
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	// 2. Load reclass config
	cfg, err := s.loadReclassConfig(ctx, char)
	if err != nil {
		return fmt.Errorf("reclass config: %w", err)
	}

	// 3. Check if reclass is allowed
	if !cfg.Allowed {
		return fmt.Errorf("reclass is not allowed in this world")
	}

	// 4. Check minimum level
	if cfg.MinLevel > 0 && char.Level < cfg.MinLevel {
		return fmt.Errorf("character level %d is below minimum level %d for reclass", char.Level, cfg.MinLevel)
	}

	// 5. Check cooldown — find the most recent class history record with a left_at
	if cfg.CooldownSeconds > 0 {
		history, err := s.client.CharacterClassHistory.Query().
			Where(characterclasshistory.HasCharacterWith(character.IDEQ(characterID))).
			All(ctx)
		if err != nil {
			return fmt.Errorf("query class history for cooldown: %w", err)
		}
		for _, h := range history {
			if h.LeftAt != nil {
				elapsed := time.Since(*h.LeftAt)
				if elapsed < time.Duration(cfg.CooldownSeconds)*time.Second {
					remaining := time.Duration(cfg.CooldownSeconds)*time.Second - elapsed
					return fmt.Errorf("reclass cooldown: %.0f seconds remaining", remaining.Seconds())
				}
			}
		}
	}

	// 6. Load current active faction
	activeFactions, err := s.client.CharacterFaction.Query().
		Where(characterfaction.HasCharacterWith(character.IDEQ(characterID))).
		Where(characterfaction.StatusEQ("active")).
		WithFaction().
		All(ctx)
	if err != nil {
		return fmt.Errorf("query active factions: %w", err)
	}

	if len(activeFactions) == 0 {
		return fmt.Errorf("character has no active class faction to reclass from")
	}

	oldCF := activeFactions[0]
	oldFaction := oldCF.Edges.Faction
	if oldFaction == nil {
		return fmt.Errorf("failed to load old faction details")
	}

	// Don't reclass to the same faction
	if oldFaction.ID == newFactionID {
		return fmt.Errorf("character is already in this class")
	}

	// 7. Load new faction
	newFaction, err := s.client.Faction.Get(ctx, newFactionID)
	if err != nil {
		return fmt.Errorf("get new faction: %w", err)
	}

	now := time.Now()

	// 8. Set old faction membership to inactive
	_, err = s.client.CharacterFaction.UpdateOne(oldCF).
		SetStatus("inactive").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("deactivate old faction membership: %w", err)
	}

	// 9. Create class history record
	err = s.client.CharacterClassHistory.Create().
		SetCharacterID(characterID).
		SetFactionID(oldFaction.ID).
		SetFactionName(oldFaction.Name).
		SetJoinedAt(oldCF.JoinedAt).
		SetLeftAt(now).
		SetReason("reclass").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("create class history: %w", err)
	}

	// 10. Create new active faction membership
	_, err = s.client.CharacterFaction.Create().
		SetCharacterID(characterID).
		SetFactionID(newFactionID).
		SetStatus("active").
		SetJoinedAt(now).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("create new faction membership: %w", err)
	}

	// 11. Update character's class field
	_, err = s.client.Character.UpdateOne(char).
		SetClass(newFaction.Name).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("update character class: %w", err)
	}

	// 12. Skill retention: multiply all character_skills levels by skill_retention
	if cfg.SkillRetention > 0 && cfg.SkillRetention < 1.0 {
		charSkills, err := s.client.CharacterSkill.Query().
			Where(characterskill.HasCharacterWith(character.IDEQ(characterID))).
			All(ctx)
		if err != nil {
			s.logger.Warn("failed to query character skills for retention", "error", err)
		} else {
			for _, cs := range charSkills {
				newLevel := int(float64(cs.Level) * cfg.SkillRetention)
				if newLevel < 1 {
					newLevel = 1
				}
				_, err := s.client.CharacterSkill.UpdateOne(cs).
					SetLevel(newLevel).
					Save(ctx)
				if err != nil {
					s.logger.Warn("failed to update skill level during reclass",
						"skill_id", cs.SkillID, "error", err)
				}
			}
		}
	}

	// 13. Remove old class-specific abilities (character_abilities where ability has faction edge to old faction)
	// First, find all abilities that belong to the old faction
	oldFactionAbilities, err := s.client.Ability.Query().
		Where(ability.HasFactionWith(faction.IDEQ(oldFaction.ID))).
		All(ctx)
	if err != nil {
		s.logger.Warn("failed to query old faction abilities", "error", err)
	} else {
		// Build a set of ability IDs to remove
		abilityIDsToRemove := make(map[int]bool)
		for _, a := range oldFactionAbilities {
			abilityIDsToRemove[a.ID] = true
		}

		// Find and delete character_abilities linked to those abilities
		if len(abilityIDsToRemove) > 0 {
			charAbilities, err := s.client.CharacterAbility.Query().
				Where(characterability.HasCharacterWith(character.IDEQ(characterID))).
				WithAbility().
				All(ctx)
			if err != nil {
				s.logger.Warn("failed to query character abilities for removal", "error", err)
			} else {
				for _, ca := range charAbilities {
					ab := ca.Edges.Ability
					if ab != nil && abilityIDsToRemove[ab.ID] {
						err := s.client.CharacterAbility.DeleteOne(ca).Exec(ctx)
						if err != nil {
							s.logger.Warn("failed to delete character ability during reclass",
								"ability_id", ab.ID, "error", err)
						}
					}
				}
			}
		}
	}

	// 14. Emit reclass event
	events.Publish(events.Event{
		Type: events.EventReclass,
		Payload: map[string]interface{}{
			"character_id": characterID,
			"old_faction":  oldFaction.Name,
			"new_faction":  newFaction.Name,
		},
	})

	s.logger.Info("reclass completed",
		"character_id", characterID,
		"old_faction", oldFaction.Name,
		"new_faction", newFaction.Name,
		"skill_retention", cfg.SkillRetention,
	)

	return nil
}

// Rerace changes a character's race with stat recalculation and history tracking.
func (s *reclassReraceService) Rerace(ctx context.Context, characterID int, newRaceName string) error {
	// 1. Load character
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return fmt.Errorf("get character: %w", err)
	}

	// 2. Load rerace config
	cfg, err := s.loadReraceConfig(ctx, char)
	if err != nil {
		return fmt.Errorf("rerace config: %w", err)
	}

	// 3. Check if rerace is allowed
	if !cfg.Allowed {
		return fmt.Errorf("rerace is not allowed in this world")
	}

	// 4. Don't rerace to the same race
	if char.Race == newRaceName {
		return fmt.Errorf("character is already race %s", newRaceName)
	}

	// 5. Load the new race by name (matching by name, first match)
	races, err := s.client.Race.Query().
		Where(race.NameEQ(newRaceName)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query race: %w", err)
	}
	if len(races) == 0 {
		return fmt.Errorf("race %s not found", newRaceName)
	}
	newRace := races[0]

	// 6. Find old race for history record (by name)
	var oldRaceID *int
	if char.Race != "" {
		oldRaces, err := s.client.Race.Query().
			Where(race.NameEQ(char.Race)).
			All(ctx)
		if err == nil && len(oldRaces) > 0 {
			oldRaceID = &oldRaces[0].ID
		}
	}

	now := time.Now()

	// 7. Create race history record
	historyCreate := s.client.CharacterRaceHistory.Create().
		SetCharacterID(characterID).
		SetRaceName(char.Race).
		SetChangedAt(now).
		SetReason("rerace")
	if oldRaceID != nil {
		historyCreate = historyCreate.SetNillableRaceID(oldRaceID)
	}
	err = historyCreate.Exec(ctx)
	if err != nil {
		return fmt.Errorf("create race history: %w", err)
	}

	// 8. Update character's race field
	_, err = s.client.Character.UpdateOne(char).
		SetRace(newRaceName).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("update character race: %w", err)
	}

	// 9. Recalculate base stats from new race's stat_modifiers
	// The stat_modifiers field is a JSON string like {"strength": 2, "dexterity": 0, ...}
	// We recalculate max HP, mana, stamina based on the new race's stat growth multipliers.
	if newRace.StatGrowthMultipliers != nil {
		hpMult := 1.0
		manaMult := 1.0
		staminaMult := 1.0
		if v, ok := newRace.StatGrowthMultipliers["hp"]; ok && v > 0 {
			hpMult = v
		}
		if v, ok := newRace.StatGrowthMultipliers["mana"]; ok && v > 0 {
			manaMult = v
		}
		if v, ok := newRace.StatGrowthMultipliers["stamina"]; ok && v > 0 {
			staminaMult = v
		}

		// Recalculate max stats based on level and race multipliers
		baseHP := DefaultHPPerLevel
		baseMana := DefaultManaPerLevel
		baseStamina := DefaultStaminaPerLevel

		newMaxHP := int(float64(baseHP*char.Level) * hpMult)
		newMaxMana := int(float64(baseMana*char.Level) * manaMult)
		newMaxStamina := int(float64(baseStamina*char.Level) * staminaMult)

		// Ensure minimums
		if newMaxHP < 1 {
			newMaxHP = 1
		}
		if newMaxMana < 1 {
			newMaxMana = 1
		}
		if newMaxStamina < 1 {
			newMaxStamina = 1
		}

		_, err = s.client.Character.UpdateOne(char).
			SetMaxHitpoints(newMaxHP).
			SetMaxMana(newMaxMana).
			SetMaxStamina(newMaxStamina).
			SetHitpoints(newMaxHP).
			SetMana(newMaxMana).
			SetStamina(newMaxStamina).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("recalculate stats after rerace: %w", err)
		}

		s.logger.Info("stats recalculated after rerace",
			"character_id", characterID,
			"new_max_hp", newMaxHP,
			"new_max_mana", newMaxMana,
			"new_max_stamina", newMaxStamina,
		)
	}

	// 10. Apply race stat modifiers if present (JSON string)
	if newRace.StatModifiers != "" {
		var modifiers map[string]int
		if err := json.Unmarshal([]byte(newRace.StatModifiers), &modifiers); err != nil {
			s.logger.Warn("failed to parse race stat_modifiers", "race", newRaceName, "error", err)
		}
		// Stat modifiers are informational for now — the ResistanceService uses the race
		// entity directly for resistance calculations, so they will apply automatically.
	}

	// 11. Emit rerace event
	events.Publish(events.Event{
		Type: events.EventRerace,
		Payload: map[string]interface{}{
			"character_id": characterID,
			"old_race":     char.Race,
			"new_race":     newRaceName,
		},
	})

	s.logger.Info("rerace completed",
		"character_id", characterID,
		"old_race", char.Race,
		"new_race", newRaceName,
	)

	return nil
}

// GetClassHistory returns the class (faction) change history for a character.
func (s *reclassReraceService) GetClassHistory(ctx context.Context, characterID int) ([]*db.CharacterClassHistory, error) {
	return s.client.CharacterClassHistory.Query().
		Where(characterclasshistory.CharacterIDEQ(characterID)).
		All(ctx)
}

// GetRaceHistory returns the race change history for a character.
func (s *reclassReraceService) GetRaceHistory(ctx context.Context, characterID int) ([]*db.CharacterRaceHistory, error) {
	return s.client.CharacterRaceHistory.Query().
		Where(characterracehistory.CharacterIDEQ(characterID)).
		All(ctx)
}