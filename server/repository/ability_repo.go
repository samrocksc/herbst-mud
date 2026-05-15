package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/ability"
)

type entAbilityRepo struct {
	client *db.Client
}

func NewEntAbilityRepo(client *db.Client) AbilityRepo {
	return &entAbilityRepo{client: client}
}

func (r *entAbilityRepo) Get(ctx context.Context, id int) (*db.Ability, error) {
	return r.client.Ability.Get(ctx, id)
}

func (r *entAbilityRepo) List(ctx context.Context, worldID string) ([]*db.Ability, error) {
	query := r.client.Ability.Query()
	if worldID != "" {
		query = query.Where(ability.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entAbilityRepo) ListClassless(ctx context.Context, worldID string) ([]*db.Ability, error) {
	query := r.client.Ability.Query().Where(ability.AbilityClassEQ("classless"))
	if worldID != "" {
		query = query.Where(ability.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entAbilityRepo) ListByClass(ctx context.Context, worldID string, class string) ([]*db.Ability, error) {
	query := r.client.Ability.Query().Where(ability.AbilityClassEQ(class))
	if worldID != "" {
		query = query.Where(ability.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entAbilityRepo) Create(ctx context.Context, input CreateAbilityInput) (*db.Ability, error) {
	builder := r.client.Ability.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetAbilityType(input.AbilityType).
		SetCost(input.Cost).
		SetCooldown(input.Cooldown).
		SetManaCost(input.ManaCost).
		SetStaminaCost(input.StaminaCost).
		SetHpCost(input.HPCost).
		SetRequirements(input.Requirements).
		SetSlug(input.Slug).
		SetAbilityClass(input.AbilityClass).
		SetWorldID(input.WorldID)
	if input.RequiredTag != "" {
		builder = builder.SetRequiredTag(input.RequiredTag)
	}
	if input.ProcChance != 0 {
		builder = builder.SetProcChance(input.ProcChance)
	}
	if input.ProcEvent != "" {
		builder = builder.SetProcEvent(input.ProcEvent)
	}
	if input.CooldownSeconds != 0 {
		builder = builder.SetCooldownSeconds(input.CooldownSeconds)
	}
	if input.FactionID != nil {
		builder = builder.SetFactionID(*input.FactionID)
	}
	return builder.Save(ctx)
}

func (r *entAbilityRepo) Update(ctx context.Context, id int, updates AbilityUpdates) (*db.Ability, error) {
	builder := r.client.Ability.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.AbilityType != nil {
		builder = builder.SetAbilityType(*updates.AbilityType)
	}
	if updates.AbilityClass != nil {
		builder = builder.SetAbilityClass(*updates.AbilityClass)
	}
	if updates.Cost != nil {
		builder = builder.SetCost(*updates.Cost)
	}
	if updates.Cooldown != nil {
		builder = builder.SetCooldown(*updates.Cooldown)
	}
	if updates.ManaCost != nil {
		builder = builder.SetManaCost(*updates.ManaCost)
	}
	if updates.StaminaCost != nil {
		builder = builder.SetStaminaCost(*updates.StaminaCost)
	}
	if updates.HPCost != nil {
		builder = builder.SetHpCost(*updates.HPCost)
	}
	if updates.Requirements != nil {
		builder = builder.SetRequirements(*updates.Requirements)
	}
	if updates.RequiredTag != nil {
		builder = builder.SetRequiredTag(*updates.RequiredTag)
	}
	if updates.ProcChance != nil {
		builder = builder.SetProcChance(*updates.ProcChance)
	}
	if updates.ProcEvent != nil {
		builder = builder.SetProcEvent(*updates.ProcEvent)
	}
	if updates.CooldownSeconds != nil {
		builder = builder.SetCooldownSeconds(*updates.CooldownSeconds)
	}
	if updates.Slug != nil {
		builder = builder.SetSlug(*updates.Slug)
	}
	if updates.FactionID != nil {
		builder = builder.SetFactionID(*updates.FactionID)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entAbilityRepo) Delete(ctx context.Context, id int) error {
	return r.client.Ability.DeleteOneID(id).Exec(ctx)
}