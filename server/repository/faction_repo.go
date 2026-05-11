package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/faction"
)

type entFactionRepo struct {
	client *db.Client
}

func NewEntFactionRepo(client *db.Client) FactionRepo {
	return &entFactionRepo{client: client}
}

func (r *entFactionRepo) Get(ctx context.Context, id int) (*db.Faction, error) {
	return r.client.Faction.Get(ctx, id)
}

func (r *entFactionRepo) GetWithEdges(ctx context.Context, id int) (*db.Faction, error) {
	return r.client.Faction.Query().
		Where(faction.ID(id)).
		WithCategory().
		WithRequiredTags().
		WithAbilities().
		Only(ctx)
}

func (r *entFactionRepo) List(ctx context.Context) ([]*db.Faction, error) {
	return r.client.Faction.Query().All(ctx)
}

func (r *entFactionRepo) Create(ctx context.Context, input CreateFactionInput) (*db.Faction, error) {
	builder := r.client.Faction.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName)
	if input.Description != "" {
		builder = builder.SetDescription(input.Description)
	}
	if len(input.MemberTags) > 0 {
		builder = builder.SetMemberTags(input.MemberTags)
	}
	return builder.Save(ctx)
}

func (r *entFactionRepo) Update(ctx context.Context, id int, updates FactionUpdates) (*db.Faction, error) {
	builder := r.client.Faction.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.DisplayName != nil {
		builder = builder.SetDisplayName(*updates.DisplayName)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.MemberTags != nil {
		builder = builder.SetMemberTags(updates.MemberTags)
	}
	return builder.Save(ctx)
}

func (r *entFactionRepo) Delete(ctx context.Context, id int) error {
	return r.client.Faction.DeleteOneID(id).Exec(ctx)
}