package service

import (
	"context"
	"errors"

	"herbst-server/db"
	"herbst-server/repository"
)

var (
	ErrSlotOutOfRange    = errors.New("slot must be between 1 and 5")
	ErrAbilityNotFound   = errors.New("ability not found")
	ErrNoAvailableSlots   = errors.New("no available slots")
	ErrAlreadyEquipped   = errors.New("ability already equipped for this character")
	ErrNotPassive        = errors.New("ability is not a passive ability")
	ErrNotEquipped       = errors.New("ability not equipped for this character")
	ErrMaxAbilities      = errors.New("cannot equip more than 5 abilities")
	ErrSkillRequirements  = errors.New("skill requirements not met")
)

type abilityService struct {
	charAbilityRepo repository.CharacterAbilityRepo
	abilityRepo     repository.AbilityRepo
	charRepo        repository.CharacterRepo
}

func NewAbilityService(
	charAbilityRepo repository.CharacterAbilityRepo,
	abilityRepo repository.AbilityRepo,
	charRepo repository.CharacterRepo,
) AbilityService {
	return &abilityService{
		charAbilityRepo: charAbilityRepo,
		abilityRepo:     abilityRepo,
		charRepo:        charRepo,
	}
}

func (s *abilityService) GetAbility(ctx context.Context, id int) (*db.Ability, error) {
	return s.abilityRepo.Get(ctx, id)
}

func (s *abilityService) ListAbilities(ctx context.Context) ([]*db.Ability, error) {
	return s.abilityRepo.List(ctx)
}

func (s *abilityService) ListClasslessAbilities(ctx context.Context) ([]*db.Ability, error) {
	return s.abilityRepo.ListClassless(ctx)
}

func (s *abilityService) ListPassiveAbilities(ctx context.Context) ([]*db.Ability, error) {
	return s.abilityRepo.ListByClass(ctx, "passive")
}

func (s *abilityService) CreateAbility(ctx context.Context, input CreateAbilityInput) (*db.Ability, error) {
	return s.abilityRepo.Create(ctx, repository.CreateAbilityInput{
		Name:            input.Name,
		Description:     input.Description,
		AbilityType:    input.AbilityType,
		AbilityClass:   input.AbilityClass,
		Cost:           input.Cost,
		Cooldown:       input.Cooldown,
		ManaCost:       input.ManaCost,
		StaminaCost:    input.StaminaCost,
		HPCost:         input.HPCost,
		Requirements:   input.Requirements,
		RequiredTag:    input.RequiredTag,
		ProcChance:     input.ProcChance,
		ProcEvent:      input.ProcEvent,
		CooldownSeconds: input.CooldownSeconds,
		Slug:           input.Slug,
		FactionID:      input.FactionID,
	})
}

func (s *abilityService) UpdateAbility(ctx context.Context, id int, input UpdateAbilityInput) (*db.Ability, error) {
	return s.abilityRepo.Update(ctx, id, repository.AbilityUpdates{
		Name:            input.Name,
		Description:     input.Description,
		AbilityType:    input.AbilityType,
		AbilityClass:   input.AbilityClass,
		Cost:           input.Cost,
		Cooldown:       input.Cooldown,
		ManaCost:       input.ManaCost,
		StaminaCost:    input.StaminaCost,
		HPCost:         input.HPCost,
		Requirements:   input.Requirements,
		RequiredTag:    input.RequiredTag,
		ProcChance:     input.ProcChance,
		ProcEvent:      input.ProcEvent,
		CooldownSeconds: input.CooldownSeconds,
		Slug:           input.Slug,
		FactionID:      input.FactionID,
	})
}

func (s *abilityService) DeleteAbility(ctx context.Context, id int) error {
	return s.abilityRepo.Delete(ctx, id)
}