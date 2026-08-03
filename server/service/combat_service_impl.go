package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
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
	deathSvc        DeathService
	logger          *slog.Logger
}

func NewCombatService(
	charRepo repository.CharacterRepo,
	damageRepo repository.DamageLogRepo,
	npcTmplRepo repository.NPCTemplateRepo,
	equipRepo repository.EquipmentRepo,
	resistanceSvc ResistanceService,
	deathSvc DeathService,
	logger *slog.Logger,
) CombatService {
	return &combatService{
		charRepo:      charRepo,
		damageRepo:    damageRepo,
		npcTmplRepo:   npcTmplRepo,
		equipRepo:     equipRepo,
		resistanceSvc: resistanceSvc,
		deathSvc:      deathSvc,
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
	var goldDropped int
	if defeated && updated.IsNPC {
		baseXP := updated.Level * 100

		// Enrich the event payload with NPC name, template ID, and attacker info
		// for biography tracking.
		npcName := updated.Name
		npcTemplateID := updated.NpcTemplateID
		attackerName := ""
		attackerLevel := 0
		if attackerID > 0 {
			if attacker, err := s.charRepo.Get(ctx, attackerID); err == nil {
				attackerName = attacker.Name
				attackerLevel = attacker.Level
			}
		}

		// Check if this NPC is a boss and roll gold drops
		isBoss := false
		if npcTemplateID != "" {
			if tmpl, err := s.npcTmplRepo.Get(ctx, npcTemplateID); err == nil {
				isBoss = tmpl.IsBoss
				// Roll gold: random amount between gold_min and gold_max
				if tmpl.GoldMax > 0 {
					min := tmpl.GoldMin
					max := tmpl.GoldMax
					if max < min {
						max = min
					}
					if min == max {
						goldDropped = min
					} else {
						goldDropped = min + rand.Intn(max-min+1)
					}
				}
			}
		}

		// Award gold to attacker
		if goldDropped > 0 && attackerID > 0 {
			if attacker, err := s.charRepo.Get(ctx, attackerID); err == nil {
				newGold := attacker.GoldCredits + goldDropped
				_, _ = s.charRepo.Update(ctx, attackerID, repository.CharacterUpdates{GoldCredits: &newGold})
			}
		}

		// Build the common payload
		payload := map[string]interface{}{
			"npc_id":          float64(updated.ID),
			"npc_name":        npcName,
			"npc_template_id": npcTemplateID,
			"npc_level":       float64(updated.Level),
			"base_xp":         float64(baseXP),
			"character_id":    float64(attackerID),
			"character_name":  attackerName,
			"character_level": float64(attackerLevel),
			"room_id":         float64(updated.CurrentRoomId),
			"world_id":        updated.CurrentWorld,
		}

		if isBoss {
			// Boss kill: publish boss_killed event for biography + XP
			payload["is_boss"] = true
			events.Publish(events.Event{
				Type:     events.EventBossKilled,
				Payload:  payload,
				Timestamp: time.Now().UnixMilli(),
			})
			// Also publish npc_defeated so the XP subscriber still awards XP
			events.Publish(events.Event{
				Type:     events.EventNPCDefeated,
				Payload:  payload,
				Timestamp: time.Now().UnixMilli(),
			})
		} else {
			// Normal NPC defeat
			events.Publish(events.Event{
				Type:     events.EventNPCDefeated,
				Payload:  payload,
				Timestamp: time.Now().UnixMilli(),
			})
		}
	}
	// Player death: emit EventCharacterDied and handle respawn.
	// The event triggers the DeathPenaltySubscriber (XP loss) and the
	// event log subscriber (biography tracking). The death service
	// handles respawn — moving to respawn room and restoring vitals.
	if defeated && !updated.IsNPC && s.deathSvc != nil {
		attackerName := ""
		if attackerID > 0 {
			if attacker, err := s.charRepo.Get(ctx, attackerID); err == nil {
				attackerName = attacker.Name
			}
		}
		deathResult, err := s.deathSvc.HandlePlayerDeath(ctx, updated.ID, attackerID, attackerName)
		if err != nil {
			s.logger.Error("failed to handle player death", "character_id", updated.ID, "error", err)
		} else if deathResult != nil {
			return &CombatResult{
				ID:       updated.ID,
				HP:       deathResult.NewHP,
				MaxHP:    updated.MaxHitpoints,
				Defeated: true,
				Message:  "You have been slain. Respawning at rally point...",
			}, nil
		}
	}
	// Persist damage to damage_logs table on every hit
	if damage > 0 {
		s.LogDamage(ctx, attackerID, targetID, damage)
	}
	return &CombatResult{ID: updated.ID, HP: updated.Hitpoints, MaxHP: updated.MaxHitpoints, Defeated: defeated, GoldDropped: goldDropped}, nil
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