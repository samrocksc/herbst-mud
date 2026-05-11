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

func (r *entAbilityRepo) List(ctx context.Context) ([]*db.Ability, error) {
	return r.client.Ability.Query().All(ctx)
}

func (r *entAbilityRepo) ListClassless(ctx context.Context) ([]*db.Ability, error) {
	return r.client.Ability.Query().
		Where(ability.AbilityClassEQ("classless")).
		All(ctx)
}

func (r *entAbilityRepo) ListByClass(ctx context.Context, class string) ([]*db.Ability, error) {
	return r.client.Ability.Query().
		Where(ability.AbilityClassEQ(class)).
		All(ctx)
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
		SetAbilityClass(input.AbilityClass)
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

func (r *entAbilityRepo) Delete(ctx context.Context, id int) error {
	return r.client.Ability.DeleteOneID(id).Exec(ctx)
}