package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"herbst-server/events"
	"herbst-server/repository"
)

type combatService struct {
	charRepo        repository.CharacterRepo
	damageRepo      repository.DamageLogRepo
	npcTmplRepo     repository.NPCTemplateRepo
	equipRepo       repository.EquipmentRepo
	resistanceSvc   ResistanceService
	logger          *slog.Logger
}

func NewCombatService(
	charRepo repository.CharacterRepo,
	damageRepo repository.DamageLogRepo,
	npcTmplRepo repository.NPCTemplateRepo,
	equipRepo repository.EquipmentRepo,
	resistanceSvc ResistanceService,
	logger *slog.Logger,
) CombatService {
	return &combatService{
		charRepo:      charRepo,
		damageRepo:    damageRepo,
		npcTmplRepo:   npcTmplRepo,
		equipRepo:     equipRepo,
		resistanceSvc: resistanceSvc,
		logger:        logger,
	}
}

var ErrCharNotFound = errors.New("character not found")

func (s *combatService) ApplyDamage(ctx context.Context, attackerID, targetID int, damage int) (*CombatResult, error) {
	if damage < 0 {
		return nil, errors.New("damage must be non-negative")
	}
	char, err := s.charRepo.Get(ctx, targetID)
	if err != nil {
		return nil, ErrCharNotFound
	}

	// Apply resistance-based damage adjustment
	if s.resistanceSvc != nil && damage > 0 {
		damageType := s.getDamageType(ctx, attackerID)
		resistances, err := s.resistanceSvc.GetCharacterResistances(ctx, targetID)
		if err != nil {
			s.logger.Error("failed to get character resistances", "target_id", targetID, "error", err)
		} else if len(resistances) > 0 {
			adjusted := s.resistanceSvc.ApplyResistances(damage, damageType, resistances)
			if adjusted != damage {
				s.logger.Debug("damage adjusted by resistances",
					"target_id", targetID, "damage_type", damageType,
					"original_damage", damage, "adjusted_damage", adjusted)
			}
			damage = adjusted
		}
	}

	if char.IsImmortal {
		newHP := char.Hitpoints - damage
		if newHP < 1 {
			newHP = 1
		}
		updated, err := s.charRepo.Update(ctx, targetID, repository.CharacterUpdates{Hitpoints: &newHP})
		if err != nil {
			return nil, err
		}
		return &CombatResult{ID: updated.ID, HP: updated.Hitpoints, MaxHP: updated.MaxHitpoints, Defeated: false, Immortal: true, Message: "Took damage but cannot be killed"}, nil
	}
	newHP := char.Hitpoints - damage
	if newHP < 0 {
		newHP = 0
	}
	updates := repository.CharacterUpdates{Hitpoints: &newHP}
	if newHP == 0 && char.IsNPC {
		now := time.Now()
		updates.DiedAt = &now
	}
	updated, err := s.charRepo.Update(ctx, targetID, updates)
	if err != nil {
		return nil, err
	}
	defeated := newHP == 0
	if defeated && updated.IsNPC {
		baseXP := updated.Level * 100
		events.Publish(events.Event{
			Type: events.EventNPCDefeated,
			Payload: map[string]interface{}{
				"npc_id":    updated.ID,
				"npc_level": updated.Level,
				"base_xp":   baseXP,
			},
			Timestamp: time.Now().UnixMilli(),
		})
	}
	// Persist damage to damage_logs table on every hit
	if damage > 0 {
		s.LogDamage(ctx, attackerID, targetID, damage)
	}
	return &CombatResult{ID: updated.ID, HP: updated.Hitpoints, MaxHP: updated.MaxHitpoints, Defeated: defeated}, nil
}

func (s *combatService) LogDamage(ctx context.Context, attackerID, targetID, damage int) {
	_, err := s.damageRepo.Create(ctx, attackerID, targetID, damage)
	if err != nil {
		s.logger.Error("failed to log damage", "error", err)
	}
}

func (s *combatService) GetCombatStatus(ctx context.Context, charID int) (*CombatStatusResult, error) {
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, ErrCharNotFound
	}
	return &CombatStatusResult{ID: char.ID, HP: char.Hitpoints, MaxHP: char.MaxHitpoints, IsNPC: char.IsNPC}, nil
}

// getDamageType determines the damage type from the attacker's equipped weapon.
// Returns "physical" as the default when no weapon is found or attackerID is 0.
func (s *combatService) getDamageType(ctx context.Context, attackerID int) string {
	if attackerID == 0 || s.equipRepo == nil {
		return "physical"
	}

	items, err := s.equipRepo.ListByOwner(ctx, attackerID)
	if err != nil {
		return "physical"
	}

	for _, item := range items {
		if item.IsEquipped && item.DamageType != "" {
			return item.DamageType
		}
	}

	return "physical"
}