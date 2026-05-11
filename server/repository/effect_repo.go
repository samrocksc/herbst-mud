package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/effect"
)

type entEffectRepo struct {
	client *db.Client
}

func NewEntEffectRepo(client *db.Client) EffectRepo {
	return &entEffectRepo{client: client}
}

func (r *entEffectRepo) Get(ctx context.Context, id int) (*db.Effect, error) {
	return r.client.Effect.Get(ctx, id)
}

func (r *entEffectRepo) GetWithHooks(ctx context.Context, id int) (*db.Effect, error) {
	return r.client.Effect.Query().
		Where(effect.IDEQ(id)).
		WithHooks().
		Only(ctx)
}

func (r *entEffectRepo) List(ctx context.Context) ([]*db.Effect, error) {
	return r.client.Effect.Query().All(ctx)
}

func (r *entEffectRepo) ListWithHooks(ctx context.Context) ([]*db.Effect, error) {
	return r.client.Effect.Query().
		Order(db.Asc(effect.FieldName)).
		WithHooks().
		All(ctx)
}

func (r *entEffectRepo) Create(ctx context.Context, input CreateEffectInput) (*db.Effect, error) {
	builder := r.client.Effect.Create().
		SetName(input.Name).
		SetEffectType(input.EffectType).
		SetParameters(input.Parameters).
		SetStackLimit(input.StackLimit).
		SetIsPermanent(input.IsPermanent).
		SetDurationSecs(input.DurationSecs)
	if input.Description != "" {
		builder = builder.SetDescription(input.Description)
	}
	if input.StackMode != "" {
		builder = builder.SetStackMode(input.StackMode)
	}
	if input.Messages != nil {
		builder = builder.SetMessages(input.Messages)
	}
	return builder.Save(ctx)
}

func (r *entEffectRepo) Update(ctx context.Context, id int, updates EffectUpdates) (*db.Effect, error) {
	builder := r.client.Effect.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.EffectType != nil {
		builder = builder.SetEffectType(*updates.EffectType)
	}
	if updates.Parameters != nil {
		builder = builder.SetParameters(*updates.Parameters)
	}
	if updates.StackMode != nil {
		builder = builder.SetStackMode(*updates.StackMode)
	}
	if updates.StackLimit != nil {
		builder = builder.SetStackLimit(*updates.StackLimit)
	}
	if updates.IsPermanent != nil {
		builder = builder.SetIsPermanent(*updates.IsPermanent)
	}
	if updates.DurationSecs != nil {
		builder = builder.SetDurationSecs(*updates.DurationSecs)
	}
	if updates.Messages != nil {
		builder = builder.SetMessages(*updates.Messages)
	}
	return builder.Save(ctx)
}

func (r *entEffectRepo) Delete(ctx context.Context, id int) error {
	return r.client.Effect.DeleteOneID(id).Exec(ctx)
}