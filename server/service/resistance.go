package service

import (
	"context"
	"log/slog"
	"math"

	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/repository"
)

// resistanceService implements ResistanceService.
type resistanceService struct {
	charRepo repository.CharacterRepo
	raceRepo repository.RaceRepo
	client   *db.Client
	logger   *slog.Logger
}

// NewResistanceService creates a new ResistanceService.
func NewResistanceService(charRepo repository.CharacterRepo, raceRepo repository.RaceRepo, client *db.Client, logger *slog.Logger) ResistanceService {
	return &resistanceService{
		charRepo: charRepo,
		raceRepo: raceRepo,
		client:   client,
		logger:   logger,
	}
}

// GetCharacterResistances computes the total resistance percentages for a character.
// It combines race resistances, race vulnerabilities (negative values), and
// resistance modifiers from equipped equipment templates.
// Returns a map of {damage_type: total_resistance_percent}.
func (s *resistanceService) GetCharacterResistances(ctx context.Context, characterID int) (map[string]int, error) {
	result := make(map[string]int)

	// 1. Load the character
	char, err := s.charRepo.Get(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// 2. Load the character's race (by name + world_id)
	worldID := char.CurrentWorld
	if worldID == "" {
		worldID = "1"
	}

	race, err := s.raceRepo.GetByName(ctx, char.Race, worldID)
	if err != nil {
		// If race not found, continue with empty resistances (backward compatibility)
		s.logger.Debug("race not found for character, using no racial resistances",
			"character_id", characterID, "race", char.Race, "world_id", worldID, "error", err)
	} else {
		// 3. Start with race resistances (positive = damage reduction)
		for dmgType, pct := range race.Resistances {
			result[dmgType] += pct
		}
		// 4. Subtract race vulnerabilities (negative = damage increase)
		for dmgType, pct := range race.Vulnerabilities {
			result[dmgType] += pct // vulnerabilities are stored as negative values
		}
	}

	// 5. Load all equipped equipment for the character, eager-loading templates
	equipped, err := s.client.Equipment.Query().
		Where(equipment.OwnerId(characterID), equipment.IsEquippedEQ(true)).
		WithEquipmentTemplate().
		All(ctx)
	if err != nil {
		s.logger.Error("failed to query equipped items for resistances", "character_id", characterID, "error", err)
		return result, nil // return what we have so far (race resistances)
	}

	// 6. For each equipped item, add its template's resistance_modifiers
	for _, eq := range equipped {
		if eq.Edges.EquipmentTemplate != nil {
			tmpl := eq.Edges.EquipmentTemplate
			for dmgType, pct := range tmpl.ResistanceModifiers {
				result[dmgType] += pct
			}
		}
	}

	return result, nil
}

// ApplyResistances adjusts damage based on the target's resistance percentages.
// Positive resistance reduces damage; negative (vulnerability) increases it.
// Minimum final damage is 1 (never fully resist).
func (s *resistanceService) ApplyResistances(damage int, damageType string, resistances map[string]int) int {
	if damage <= 0 || len(resistances) == 0 {
		return damage
	}

	pct, ok := resistances[damageType]
	if !ok || pct == 0 {
		return damage
	}

	var finalDamage float64
	if pct > 0 {
		// Resistance: reduce damage
		finalDamage = float64(damage) * (1.0 - float64(pct)/100.0)
	} else {
		// Vulnerability (negative): increase damage
		finalDamage = float64(damage) * (1.0 + float64(-pct)/100.0)
	}

	// Round to nearest integer
	rounded := int(math.Round(finalDamage))

	// Minimum final damage = 1 (never fully resist)
	if rounded < 1 {
		rounded = 1
	}

	return rounded
}