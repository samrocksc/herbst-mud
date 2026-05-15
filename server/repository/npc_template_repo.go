package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/npctemplate"
)

type entNPCTemplateRepo struct {
	client *db.Client
}

func NewEntNPCTemplateRepo(client *db.Client) NPCTemplateRepo {
	return &entNPCTemplateRepo{client: client}
}

func (r *entNPCTemplateRepo) Get(ctx context.Context, id string) (*db.NPCTemplate, error) {
	return r.client.NPCTemplate.Get(ctx, id)
}

func (r *entNPCTemplateRepo) List(ctx context.Context, worldID string) ([]*db.NPCTemplate, error) {
	query := r.client.NPCTemplate.Query()
	if worldID != "" {
		query = query.Where(npctemplate.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entNPCTemplateRepo) Create(ctx context.Context, input CreateNPCTemplateInput) (*db.NPCTemplate, error) {
	builder := r.client.NPCTemplate.Create().
		SetID(input.ID).
		SetName(input.Name).
		SetDescription(input.Description).
		SetRace(input.Race).
		SetSkills(input.Skills).
		SetTradesWith(input.TradesWith).
		SetGreeting(input.Greeting).
		SetRespawnRooms(input.RespawnRooms).
		SetLevel(input.Level).
		SetXpValue(input.XPValue).
		SetDisposition(npctemplate.Disposition(input.Disposition)).
		SetWorldID(input.WorldID)
	if input.RespawnCooldown != nil {
		builder = builder.SetNillableRespawnCooldown(input.RespawnCooldown)
	}
	return builder.Save(ctx)
}

func (r *entNPCTemplateRepo) Update(ctx context.Context, id string, updates NPCTemplateUpdates) (*db.NPCTemplate, error) {
	builder := r.client.NPCTemplate.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.Race != nil {
		builder = builder.SetRace(*updates.Race)
	}
	if updates.Disposition != nil {
		builder = builder.SetDisposition(npctemplate.Disposition(*updates.Disposition))
	}
	if updates.Level != nil {
		builder = builder.SetLevel(*updates.Level)
	}
	if updates.XPValue != nil {
		builder = builder.SetXpValue(*updates.XPValue)
	}
	if updates.Skills != nil {
		builder = builder.SetSkills(*updates.Skills)
	}
	if updates.TradesWith != nil {
		builder = builder.SetTradesWith(*updates.TradesWith)
	}
	if updates.Greeting != nil {
		builder = builder.SetGreeting(*updates.Greeting)
	}
	if updates.RespawnRooms != nil {
		builder = builder.SetRespawnRooms(*updates.RespawnRooms)
	}
	if updates.RespawnCooldown != nil {
		builder = builder.SetNillableRespawnCooldown(updates.RespawnCooldown)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entNPCTemplateRepo) Delete(ctx context.Context, id string) error {
	return r.client.NPCTemplate.DeleteOneID(id).Exec(ctx)
}