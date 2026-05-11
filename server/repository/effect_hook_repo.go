package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/effect"
	"herbst-server/db/effecthook"
	"herbst-server/db/npctemplate"
)

type entEffectHookRepo struct {
	client *db.Client
}

func NewEntEffectHookRepo(client *db.Client) EffectHookRepo {
	return &entEffectHookRepo{client: client}
}

func (r *entEffectHookRepo) Get(ctx context.Context, id int) (*db.EffectHook, error) {
	return r.client.EffectHook.Get(ctx, id)
}

func (r *entEffectHookRepo) GetWithEdges(ctx context.Context, id int) (*db.EffectHook, error) {
	return r.client.EffectHook.Query().
		Where(effecthook.IDEQ(id)).
		WithEffect().
		WithNpcTemplate().
		Only(ctx)
}

func (r *entEffectHookRepo) List(ctx context.Context) ([]*db.EffectHook, error) {
	return r.client.EffectHook.Query().All(ctx)
}

func (r *entEffectHookRepo) ListWithEdges(ctx context.Context) ([]*db.EffectHook, error) {
	return r.client.EffectHook.Query().
		WithEffect().
		WithNpcTemplate().
		All(ctx)
}

func (r *entEffectHookRepo) ListByEvent(ctx context.Context, event string) ([]*db.EffectHook, error) {
	return r.client.EffectHook.Query().
		Where(effecthook.EventEQ(event)).
		All(ctx)
}

func (r *entEffectHookRepo) ListByTemplateWithEdges(ctx context.Context, templateID string) ([]*db.EffectHook, error) {
	return r.client.EffectHook.Query().
		Where(effecthook.HasNpcTemplateWith(npctemplate.IDEQ(templateID))).
		WithEffect().
		WithNpcTemplate().
		All(ctx)
}

func (r *entEffectHookRepo) CountByEffect(ctx context.Context, effectID int) (int, error) {
	return r.client.EffectHook.Query().
		Where(effecthook.HasEffectWith(effect.IDEQ(effectID))).
		Count(ctx)
}

func (r *entEffectHookRepo) Create(ctx context.Context, input CreateEffectHookInput) (*db.EffectHook, error) {
	builder := r.client.EffectHook.Create().
		SetName(input.Name).
		SetEvent(input.Event).
		SetEffectID(input.EffectID).
		SetEnabled(input.Enabled)
	if input.Target != "" {
		builder = builder.SetTarget(input.Target)
	}
	if input.Condition != "" {
		builder = builder.SetCondition(input.Condition)
	}
	if input.NPCTemplateID != nil {
		builder = builder.SetNpcTemplateID(*input.NPCTemplateID)
	}
	return builder.Save(ctx)
}

func (r *entEffectHookRepo) Update(ctx context.Context, id int, updates EffectHookUpdates) (*db.EffectHook, error) {
	builder := r.client.EffectHook.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Event != nil {
		builder = builder.SetEvent(*updates.Event)
	}
	if updates.Target != nil {
		builder = builder.SetTarget(*updates.Target)
	}
	if updates.Condition != nil {
		builder = builder.SetCondition(*updates.Condition)
	}
	if updates.Enabled != nil {
		builder = builder.SetEnabled(*updates.Enabled)
	}
	if updates.EffectID != nil {
		builder = builder.SetEffectID(*updates.EffectID)
	}
	return builder.Save(ctx)
}

func (r *entEffectHookRepo) Delete(ctx context.Context, id int) error {
	return r.client.EffectHook.DeleteOneID(id).Exec(ctx)
}