package repository

import (
	"context"
	"time"

	"herbst-server/db"
	"herbst-server/db/activeeffect"
	"herbst-server/db/effect"
)

type entActiveEffectRepo struct {
	client *db.Client
}

func NewEntActiveEffectRepo(client *db.Client) ActiveEffectRepo {
	return &entActiveEffectRepo{client: client}
}

func (r *entActiveEffectRepo) ListByCharacter(ctx context.Context, charID int) ([]*db.ActiveEffect, error) {
	return r.client.ActiveEffect.Query().
		Where(activeeffect.CharacterID(charID)).
		All(ctx)
}

func (r *entActiveEffectRepo) ListActiveByCharacter(ctx context.Context, charID int) ([]*db.ActiveEffect, error) {
	return r.client.ActiveEffect.Query().
		Where(activeeffect.CharacterIDEQ(charID), activeeffect.IsActiveEQ(true)).
		WithEffect().
		All(ctx)
}

func (r *entActiveEffectRepo) GetActiveByCharacterAndEffect(ctx context.Context, charID, effectID int) (*db.ActiveEffect, error) {
	return r.client.ActiveEffect.Query().
		Where(
			activeeffect.CharacterIDEQ(charID),
			activeeffect.HasEffectWith(effect.IDEQ(effectID)),
			activeeffect.IsActiveEQ(true),
		).
		Only(ctx)
}

func (r *entActiveEffectRepo) GetWithEffect(ctx context.Context, id int) (*db.ActiveEffect, error) {
	return r.client.ActiveEffect.Query().
		Where(activeeffect.IDEQ(id)).
		WithEffect().
		Only(ctx)
}

func (r *entActiveEffectRepo) Create(ctx context.Context, input CreateActiveEffectInput) (*db.ActiveEffect, error) {
	builder := r.client.ActiveEffect.Create().
		SetCharacterID(input.CharacterID).
		SetEffectID(input.EffectID).
		SetStackCount(input.StackCount)
	if input.AppliedByID != 0 {
		builder = builder.SetAppliedByID(input.AppliedByID)
	}
	if input.ExpiresAt != nil {
		builder = builder.SetNillableExpiresAt(input.ExpiresAt)
	}
	return builder.Save(ctx)
}

func (r *entActiveEffectRepo) Update(ctx context.Context, id int, updates ActiveEffectUpdates) (*db.ActiveEffect, error) {
	builder := r.client.ActiveEffect.UpdateOneID(id)
	if updates.StackCount != nil {
		builder = builder.SetStackCount(*updates.StackCount)
	}
	if updates.IsActive != nil {
		builder = builder.SetIsActive(*updates.IsActive)
	}
	if updates.ExpiresAt != nil {
		builder = builder.SetExpiresAt(*updates.ExpiresAt)
	}
	if updates.StartedAt != nil {
		builder = builder.SetStartedAt(*updates.StartedAt)
	}
	return builder.Save(ctx)
}

func (r *entActiveEffectRepo) Delete(ctx context.Context, id int) error {
	return r.client.ActiveEffect.DeleteOneID(id).Exec(ctx)
}

func (r *entActiveEffectRepo) DeactivateExpired(ctx context.Context) ([]*db.ActiveEffect, error) {
	now := time.Now()
	expired, err := r.client.ActiveEffect.Query().
		Where(
			activeeffect.IsActiveEQ(true),
			activeeffect.ExpiresAtLTE(now),
		).
		WithEffect().
		All(ctx)
	if err != nil {
		return nil, err
	}
	var deactivated []*db.ActiveEffect
	for _, ae := range expired {
		updated, err := r.client.ActiveEffect.UpdateOneID(ae.ID).
			SetIsActive(false).
			Save(ctx)
		if err != nil {
			continue
		}
		// Preserve Effect edge from original query for response data
		updated.Edges.Effect = ae.Edges.Effect
		deactivated = append(deactivated, updated)
	}
	return deactivated, nil
}