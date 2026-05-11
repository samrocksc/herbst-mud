package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/gender"
)

type entGenderRepo struct {
	client *db.Client
}

func NewEntGenderRepo(client *db.Client) GenderRepo {
	return &entGenderRepo{client: client}
}

func (r *entGenderRepo) Get(ctx context.Context, id int) (*db.Gender, error) {
	return r.client.Gender.Get(ctx, id)
}

func (r *entGenderRepo) GetByName(ctx context.Context, name string) (*db.Gender, error) {
	return r.client.Gender.Query().
		Where(gender.Name(name)).
		Only(ctx)
}

func (r *entGenderRepo) List(ctx context.Context) ([]*db.Gender, error) {
	return r.client.Gender.Query().All(ctx)
}