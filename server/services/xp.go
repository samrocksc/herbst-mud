package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect/sql"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/competencycategory"
	"herbst-server/db/competencylevelthreshold"
	"herbst-server/db/gameconfig"
)

// DefaultXPPerLevel is used when no GameConfig override exists.
const DefaultXPPerLevel = 200

// XPThresholds holds the XP required to reach each level.
type XPThresholds map[int]int

// XPThresholdsFromConfig parses thresholds from a GameConfig value JSON string.
func XPThresholdsFromConfig(value string) (XPThresholds, error) {
	if value == "" {
		return nil, fmt.Errorf("empty config value")
	}
	var t XPThresholds
	if err := json.Unmarshal([]byte(value), &t); err != nil {
		return nil, err
	}
	return t, nil
}

// ThresholdForLevel returns the XP required to reach the given level.
// Returns 0 for level < 1.
func (t XPThresholds) ThresholdForLevel(level int) int {
	if level < 1 {
		return 0
	}
	if xp, ok := t[level]; ok {
		return xp
	}
	// Fallback: DefaultXPPerLevel * level
	return DefaultXPPerLevel * level
}

// XPAwardService handles XP awards and level-up logic.
type XPAwardService struct {
	client *db.Client
	logger *slog.Logger
}

// NewXPAwardService creates a new XP service.
func NewXPAwardService(client *db.Client, logger *slog.Logger) *XPAwardService {
	if logger == nil {
		logger = slog.Default()
	}
	return &XPAwardService{client: client, logger: logger}
}

// AwardXP adds XP to a character and handles level-up if thresholds are crossed.
// Returns new XP, new level, and whether a level-up occurred.
func (s *XPAwardService) AwardXP(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, false, fmt.Errorf("get character: %w", err)
	}

	charXP := char.Xp
	charLevel := char.Level

	// Apply XP
	charXP += xpGained

	// Get thresholds from config
	thresholds := s.getThresholds(ctx)

	// Check for level-up (can happen multiple times in one award)
	for {
		nextLevel := charLevel + 1
		needed := thresholds.ThresholdForLevel(nextLevel)
		if needed == 0 || charXP < needed {
			break
		}
		charLevel = nextLevel
		s.logger.Info("level up",
			"character_id", characterID,
			"new_level", charLevel,
			"xp_remaining", charXP,
		)
	}

	// Update character
	_, err = s.client.Character.UpdateOne(char).
		SetXp(charXP).
		SetLevel(charLevel).
		Save(ctx)
	if err != nil {
		return 0, 0, false, fmt.Errorf("update character xp/level: %w", err)
	}

	return charXP, charLevel, charLevel > char.Level, nil
}

// ApplyDeathPenalty reduces XP by penaltyPercent and returns XP lost and new XP.
func (s *XPAwardService) ApplyDeathPenalty(ctx context.Context, characterID, penaltyPercent int) (xpLost, newXP int, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, fmt.Errorf("get character: %w", err)
	}

	xpLost = (char.Xp * penaltyPercent) / 100
	newXP = char.Xp - xpLost

	_, err = s.client.Character.UpdateOne(char).
		SetXp(newXP).
		Save(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("update character xp after death: %w", err)
	}

	return xpLost, newXP, nil
}

// getThresholds loads XP thresholds from GameConfig or returns the default.
func (s *XPAwardService) getThresholds(ctx context.Context) XPThresholds {
		cfg, err := s.client.GameConfig.Query().Where(gameconfig.Key("xp_thresholds")).Only(ctx)
	if err != nil {
		// Fall back to default linear thresholds
		thresholds := make(XPThresholds)
		for i := 1; i <= 100; i++ {
			thresholds[i] = DefaultXPPerLevel * i
		}
		return thresholds
	}

	thresholds, err := XPThresholdsFromConfig(cfg.Value)
	if err != nil {
		s.logger.Warn("failed to parse xp_thresholds config, using defaults", "error", err)
		thresholds = make(XPThresholds)
		for i := 1; i <= 100; i++ {
			thresholds[i] = DefaultXPPerLevel * i
		}
	}
	return thresholds
}

// GetCharacterXP returns the current XP and level for a character.
func (s *XPAwardService) GetCharacterXP(ctx context.Context, characterID int) (xp, level int, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, err
	}
	return char.Xp, char.Level, nil
}

// QueryCharacter queries a character by ID (used by other services).
func (s *XPAwardService) QueryCharacter(ctx context.Context, id int) (*db.Character, error) {
	return s.client.Character.Get(ctx, id)
}

// AwardCompetencyXP awards XP to a character's competency in a category.
// It applies the category's xp_multiplier, updates the character's competency record,
// and recomputes the cached level based on thresholds.
func (s *XPAwardService) AwardCompetencyXP(ctx context.Context, characterID int, categoryID string, rawXP int) error {
	cat, err := s.client.CompetencyCategory.Get(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get competency category %s: %w", categoryID, err)
	}

	multiplied := int(float64(rawXP) * cat.XpMultiplier)

	cc, err := s.client.CharacterCompetency.Query().
		Where(charactercompetency.HasCharacterWith(character.ID(characterID))).
		Where(charactercompetency.HasCategoryWith(competencycategory.ID(categoryID))).
		Only(ctx)
	if err != nil {
		// Record doesn't exist — create it
		cc, err = s.client.CharacterCompetency.Create().
			SetXp(multiplied).
			SetLevel(1).
			SetCharacterID(characterID).
			SetCategoryID(categoryID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create character competency: %w", err)
		}
		s.logger.Info("competency started",
			"character_id", characterID, "category", categoryID, "xp", multiplied)
		return nil
	}

	// Update XP
	cc.Xp += multiplied

	// Recompute level from thresholds
	thresholds, err := s.client.CompetencyLevelThreshold.Query().
		Where(competencylevelthreshold.HasCategoryWith(competencycategory.ID(categoryID))).
		Order(competencylevelthreshold.ByLevel(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query thresholds: %w", err)
	}

	newLevel := cc.Level
	for _, t := range thresholds {
		if cc.Xp >= t.XpRequired {
			newLevel = t.Level
		}
	}

	cc.Level = newLevel

	_, err = s.client.CharacterCompetency.UpdateOne(cc).SetXp(cc.Xp).SetLevel(cc.Level).Save(ctx)
	if err != nil {
		return fmt.Errorf("update character competency: %w", err)
	}

	s.logger.Info("competency xp awarded",
		"character_id", characterID, "category", categoryID,
		"raw_xp", rawXP, "multiplied", multiplied, "total_xp", cc.Xp, "level", cc.Level)
	return nil
}

// SeedCompetencyCategories creates the default competency categories and level thresholds
// if they don't already exist. Safe to call on every startup.
func (s *XPAwardService) SeedCompetencyCategories(ctx context.Context) error {
	categories := []struct {
		id   string
		name string
		mult float64
	}{
		{"blades", "Blades", 0.20},
		{"staves", "Staves", 0.20},
		{"knives", "Knives", 0.20},
		{"martial", "Martial Arts", 0.20},
		{"brawling", "Brawling", 0.20},
		{"tech", "Tech", 0.20},
		{"light_armor", "Light Armor", 0.20},
		{"cloth_armor", "Cloth Armor", 0.20},
		{"heavy_armor", "Heavy Armor", 0.20},
	}

	// Level thresholds: levels 1-10 with escalating XP costs
	levelXP := []int{0, 100, 250, 500, 900, 1500, 2300, 3300, 4500, 6000}
	dmgMults := []float64{1.0, 1.05, 1.10, 1.15, 1.20, 1.30, 1.40, 1.55, 1.70, 1.90}
	defMults := []float64{1.0, 1.03, 1.06, 1.09, 1.12, 1.16, 1.22, 1.30, 1.40, 1.55}

	for _, cat := range categories {
		existing, _ := s.client.CompetencyCategory.Get(ctx, cat.id)
		if existing != nil {
			continue
		}

		created, err := s.client.CompetencyCategory.Create().
			SetID(cat.id).
			SetName(cat.name).
			SetXpMultiplier(cat.mult).
			Save(ctx)
		if err != nil {
			s.logger.Warn("failed to seed competency category", "id", cat.id, "error", err)
			continue
		}

		// Create thresholds for levels 1-10
		for lvl := 1; lvl <= 10; lvl++ {
			_, err = s.client.CompetencyLevelThreshold.Create().
				SetID(fmt.Sprintf("%s-%d", cat.id, lvl)).
				SetLevel(lvl).
				SetXpRequired(levelXP[lvl-1]).
				SetDamageMultiplier(dmgMults[lvl-1]).
				SetDefenseMultiplier(defMults[lvl-1]).
				SetCategory(created).
				Save(ctx)
			if err != nil {
				s.logger.Warn("failed to seed threshold", "id", cat.id, "level", lvl, "error", err)
			}
		}
		s.logger.Info("seeded competency category", "id", cat.id, "name", cat.name)
	}
	return nil
}
