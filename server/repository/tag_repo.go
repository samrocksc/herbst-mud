package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/tag"
)

type entTagRepo struct {
	client *db.Client
}

func NewEntTagRepo(client *db.Client) TagRepo {
	return &entTagRepo{client: client}
}

func (r *entTagRepo) Get(ctx context.Context, id int) (*db.Tag, error) {
	return r.client.Tag.Get(ctx, id)
}

func (r *entTagRepo) GetByName(ctx context.Context, name, worldID string) (*db.Tag, error) {
	return r.client.Tag.Query().
		Where(tag.Name(name), tag.WorldID(worldID)).
		Only(ctx)
}

func (r *entTagRepo) List(ctx context.Context, worldID string) ([]*db.Tag, error) {
	return r.client.Tag.Query().Where(tag.WorldID(worldID)).All(ctx)
}

func (r *entTagRepo) Create(ctx context.Context, input CreateTagInput) (*db.Tag, error) {
	builder := r.client.Tag.Create().
		SetName(input.Name)
	if input.Color != "" {
		builder = builder.SetColor(input.Color)
	}
	if input.WorldID != "" {
		builder = builder.SetWorldID(input.WorldID)
	} else {
		builder = builder.SetWorldID("1")
	}
	return builder.Save(ctx)
}

func (r *entTagRepo) Update(ctx context.Context, id int, updates TagUpdates) (*db.Tag, error) {
	builder := r.client.Tag.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.WorldID != nil && *updates.WorldID != "" {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entTagRepo) Delete(ctx context.Context, id int) error {
	return r.client.Tag.DeleteOneID(id).Exec(ctx)
}
