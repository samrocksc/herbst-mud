package repository

import (
	"context"

	"herbst-server/db"
)

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
	return builder.Save(ctx)
}