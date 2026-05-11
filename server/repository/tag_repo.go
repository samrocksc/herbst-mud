package repository

import (
	"context"

	"herbst-server/db"
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

func (r *entTagRepo) List(ctx context.Context) ([]*db.Tag, error) {
	return r.client.Tag.Query().All(ctx)
}

func (r *entTagRepo) Create(ctx context.Context, input CreateTagInput) (*db.Tag, error) {
	builder := r.client.Tag.Create().
		SetName(input.Name)
	if input.Color != "" {
		builder = builder.SetColor(input.Color)
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
	return builder.Save(ctx)
}

func (r *entTagRepo) Delete(ctx context.Context, id int) error {
	return r.client.Tag.DeleteOneID(id).Exec(ctx)
}