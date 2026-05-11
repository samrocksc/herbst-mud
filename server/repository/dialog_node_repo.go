package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/db/npctemplate"
)

type entDialogNodeRepo struct {
	client *db.Client
}

func NewEntDialogNodeRepo(client *db.Client) DialogNodeRepo {
	return &entDialogNodeRepo{client: client}
}

func (r *entDialogNodeRepo) Get(ctx context.Context, id string) (*db.DialogNode, error) {
	return r.client.DialogNode.Get(ctx, id)
}

func (r *entDialogNodeRepo) List(ctx context.Context) ([]*db.DialogNode, error) {
	return r.client.DialogNode.Query().All(ctx)
}

func (r *entDialogNodeRepo) ListByTemplate(ctx context.Context, templateID string) ([]*db.DialogNode, error) {
	return r.client.DialogNode.Query().
		Where(dialognode.HasNpcTemplateWith(npctemplate.ID(templateID))).
		All(ctx)
}

func (r *entDialogNodeRepo) Create(ctx context.Context, input CreateDialogNodeInput) (*db.DialogNode, error) {
	builder := r.client.DialogNode.Create().
		SetID(input.ID).
		SetNpcTemplateID(input.NPCTemplateID).
		SetNpcText(input.NPCText).
		SetResponses(input.Responses).
		SetOnEnterEffects(input.OnEnterEffects)
	if input.IsEntry {
		builder = builder.SetIsEntry(input.IsEntry)
	}
	if input.EntryCondition != "" {
		builder = builder.SetEntryCondition(input.EntryCondition)
	}
	return builder.Save(ctx)
}

func (r *entDialogNodeRepo) Update(ctx context.Context, id string, updates DialogNodeUpdates) (*db.DialogNode, error) {
	builder := r.client.DialogNode.UpdateOneID(id)
	if updates.NPCText != nil {
		builder = builder.SetNpcText(*updates.NPCText)
	}
	if updates.Responses != nil {
		builder = builder.SetResponses(*updates.Responses)
	}
	if updates.IsEntry != nil {
		builder = builder.SetIsEntry(*updates.IsEntry)
	}
	if updates.EntryCondition != nil {
		builder = builder.SetEntryCondition(*updates.EntryCondition)
	}
	if updates.OnEnterEffects != nil {
		builder = builder.SetOnEnterEffects(*updates.OnEnterEffects)
	}
	if updates.NPCTemplateID != nil {
		builder = builder.SetNpcTemplateID(*updates.NPCTemplateID)
	}
	return builder.Save(ctx)
}

func (r *entDialogNodeRepo) Delete(ctx context.Context, id string) error {
	return r.client.DialogNode.DeleteOneID(id).Exec(ctx)
}