package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"herbst-server/db"
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
