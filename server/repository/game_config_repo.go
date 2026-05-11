package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/gameconfig"
)

type entGameConfigRepo struct {
	client *db.Client
}

func NewEntGameConfigRepo(client *db.Client) GameConfigRepo {
	return &entGameConfigRepo{client: client}
}

func (r *entGameConfigRepo) Get(ctx context.Context, key string) (*db.GameConfig, error) {
	return r.client.GameConfig.Query().
		Where(gameconfig.Key(key)).
		Only(ctx)
}

func (r *entGameConfigRepo) GetOrCreate(ctx context.Context, key, defaultValue string) (*db.GameConfig, error) {
	cfg, err := r.Get(ctx, key)
	if err == nil {
		return cfg, nil
	}
	return r.client.GameConfig.Create().
		SetKey(key).
		SetValue(defaultValue).
		Save(ctx)
}

func (r *entGameConfigRepo) Set(ctx context.Context, key, value string) (*db.GameConfig, error) {
	cfg, err := r.Get(ctx, key)
	if err != nil {
		return r.client.GameConfig.Create().
			SetKey(key).
			SetValue(value).
			Save(ctx)
	}
	return r.client.GameConfig.UpdateOneID(cfg.ID).
		SetValue(value).
		Save(ctx)
}

func (r *entGameConfigRepo) List(ctx context.Context) ([]*db.GameConfig, error) {
	return r.client.GameConfig.Query().All(ctx)
}

func (r *entGameConfigRepo) Delete(ctx context.Context, key string) error {
	cfg, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return r.client.GameConfig.DeleteOne(cfg).Exec(ctx)
}