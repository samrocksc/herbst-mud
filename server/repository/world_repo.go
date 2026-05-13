package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/world"
)

type entWorldRepo struct {
	client *db.Client
}

func NewEntWorldRepo(client *db.Client) WorldRepo {
	return &entWorldRepo{client: client}
}

func (r *entWorldRepo) Get(ctx context.Context, id int) (*db.World, error) {
	return r.client.World.Get(ctx, id)
}

func (r *entWorldRepo) GetByName(ctx context.Context, name string) (*db.World, error) {
	return r.client.World.Query().Where(world.NameEQ(name)).Only(ctx)
}

func (r *entWorldRepo) List(ctx context.Context) ([]*db.World, error) {
	return r.client.World.Query().All(ctx)
}

func (r *entWorldRepo) GetActive(ctx context.Context) ([]*db.World, error) {
	return r.client.World.Query().Where(world.ActiveEQ(true)).All(ctx)
}

func (r *entWorldRepo) Create(ctx context.Context, input CreateWorldInput) (*db.World, error) {
	return r.client.World.Create().
		SetName(input.Name).
		SetTitle(input.Title).
		SetDescription(input.Description).
		SetActive(input.Active).
		Save(ctx)
}

func (r *entWorldRepo) Update(ctx context.Context, id int, updates WorldUpdates) (*db.World, error) {
	world, err := r.client.World.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	mutator := world.Update()
	if updates.Name != nil {
		mutator = mutator.SetName(*updates.Name)
	}
	if updates.Title != nil {
		mutator = mutator.SetTitle(*updates.Title)
	}
	if updates.Description != nil {
		mutator = mutator.SetDescription(*updates.Description)
	}
	if updates.Active != nil {
		mutator = mutator.SetActive(*updates.Active)
	}
	return mutator.Save(ctx)
}

func (r *entWorldRepo) Delete(ctx context.Context, id int) error {
	return r.client.World.DeleteOneID(id).Exec(ctx)
}
