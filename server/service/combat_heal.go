package service

import (
	"context"
	"errors"

	"herbst-server/repository"
)

func (s *combatService) HealCharacter(ctx context.Context, charID int, amount int) (*HealResult, error) {
	if amount < 0 {
		return nil, errors.New("heal amount must be non-negative")
	}
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, ErrCharNotFound
	}
	newHP := char.Hitpoints + amount
	if newHP > char.MaxHitpoints {
		newHP = char.MaxHitpoints
	}
	updated, err := s.charRepo.Update(ctx, charID, repository.CharacterUpdates{Hitpoints: &newHP})
	if err != nil {
		return nil, err
	}
	return &HealResult{ID: updated.ID, HP: updated.Hitpoints, MaxHP: updated.MaxHitpoints}, nil
}

func (s *combatService) AdjustStamina(ctx context.Context, charID int, amount int) (*StatResult, error) {
	// Allow negative amounts for cost deduction (but not below 0)
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, ErrCharNotFound
	}
	newStamina := char.Stamina + amount
	if newStamina > char.MaxStamina {
		newStamina = char.MaxStamina
	}
	if newStamina < 0 {
		newStamina = 0
	}
	updated, err := s.charRepo.Update(ctx, charID, repository.CharacterUpdates{Stamina: &newStamina})
	if err != nil {
		return nil, err
	}
	return &StatResult{ID: updated.ID, Current: updated.Stamina, Max: updated.MaxStamina}, nil
}

func (s *combatService) AdjustMana(ctx context.Context, charID int, amount int) (*StatResult, error) {
	// Allow negative amounts for cost deduction (but not below 0)
	char, err := s.charRepo.Get(ctx, charID)
	if err != nil {
		return nil, ErrCharNotFound
	}
	newMana := char.Mana + amount
	if newMana > char.MaxMana {
		newMana = char.MaxMana
	}
	if newMana < 0 {
		newMana = 0
	}
	updated, err := s.charRepo.Update(ctx, charID, repository.CharacterUpdates{Mana: &newMana})
	if err != nil {
		return nil, err
	}
	return &StatResult{ID: updated.ID, Current: updated.Mana, Max: updated.MaxMana}, nil
}

func (s *combatService) HealNPCsInRoom(ctx context.Context, roomID int, amount int) (int, error) {
	if amount < 0 {
		return 0, errors.New("heal amount must be non-negative")
	}
	npcs, err := s.charRepo.ListNPCsByRoom(ctx, roomID)
	if err != nil {
		return 0, err
	}
	healedCount := 0
	for _, npc := range npcs {
		if npc.Hitpoints >= npc.MaxHitpoints || npc.Hitpoints <= 0 {
			continue
		}
		newHP := npc.Hitpoints + amount
		if newHP > npc.MaxHitpoints {
			newHP = npc.MaxHitpoints
		}
		_, err := s.charRepo.Update(ctx, npc.ID, repository.CharacterUpdates{Hitpoints: &newHP})
		if err == nil {
			healedCount++
		}
	}
	return healedCount, nil
}

func (s *combatService) PassiveHealNPCsInRoom(ctx context.Context, roomID int) (*PassiveHealResult, error) {
	npcs, err := s.charRepo.ListNPCsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	healedCount := 0
	fullyHealedCount := 0
	for _, npc := range npcs {
		if npc.Hitpoints <= 0 {
			newHP := npc.MaxHitpoints / 2
			_, err := s.charRepo.Update(ctx, npc.ID, repository.CharacterUpdates{Hitpoints: &newHP})
			if err == nil {
				healedCount++
			}
			continue
		}
		if npc.Hitpoints >= npc.MaxHitpoints {
			continue
		}
		healAmount := npc.MaxHitpoints / 4
		if healAmount < 1 {
			healAmount = 1
		}
		newHP := npc.Hitpoints + healAmount
		if newHP >= npc.MaxHitpoints {
			newHP = npc.MaxHitpoints
			fullyHealedCount++
		}
		_, err := s.charRepo.Update(ctx, npc.ID, repository.CharacterUpdates{Hitpoints: &newHP})
		if err == nil {
			healedCount++
		}
	}
	return &PassiveHealResult{Healed: healedCount}, nil
}
