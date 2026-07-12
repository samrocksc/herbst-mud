package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/channelconfig"
)

type ChannelConfigRepo interface {
	List(ctx context.Context, filter ChannelConfigFilter) ([]*db.ChannelConfig, error)
	GetByName(ctx context.Context, name string) (*db.ChannelConfig, error)
	Create(ctx context.Context, input CreateChannelConfigInput) (*db.ChannelConfig, error)
	Update(ctx context.Context, name string, updates ChannelConfigUpdates) (*db.ChannelConfig, error)
	Delete(ctx context.Context, name string) error
}

type ChannelConfigFilter struct{ Search string }

type CreateChannelConfigInput struct {
	Name            string
	Description     string
	Color           string
	DefaultEnabled  bool
	CooldownSeconds int
	AdminOnly       bool
}

type ChannelConfigUpdates struct {
	Description     *string
	Color           *string
	DefaultEnabled  *bool
	CooldownSeconds *int
	AdminOnly       *bool
}

type entChannelConfigRepo struct{ client *db.Client }

func NewEntChannelConfigRepo(client *db.Client) ChannelConfigRepo {
	return &entChannelConfigRepo{client: client}
}

func (r *entChannelConfigRepo) List(ctx context.Context, filter ChannelConfigFilter) ([]*db.ChannelConfig, error) {
	query := r.client.ChannelConfig.Query()
	if filter.Search != "" {
		query = query.Where(channelconfig.NameContains(filter.Search))
	}
	return query.All(ctx)
}

func (r *entChannelConfigRepo) GetByName(ctx context.Context, name string) (*db.ChannelConfig, error) {
	return r.client.ChannelConfig.Query().Where(channelconfig.NameEQ(name)).Only(ctx)
}

func (r *entChannelConfigRepo) Create(ctx context.Context, input CreateChannelConfigInput) (*db.ChannelConfig, error) {
	return r.client.ChannelConfig.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetColor(input.Color).
		SetDefaultEnabled(input.DefaultEnabled).
		SetCooldownSeconds(input.CooldownSeconds).
		SetAdminOnly(input.AdminOnly).
		Save(ctx)
}

func (r *entChannelConfigRepo) Update(ctx context.Context, name string, updates ChannelConfigUpdates) (*db.ChannelConfig, error) {
	builder := r.client.ChannelConfig.Update().Where(channelconfig.NameEQ(name))
	if updates.Description != nil {
		builder.SetDescription(*updates.Description)
	}
	if updates.Color != nil {
		builder.SetColor(*updates.Color)
	}
	if updates.DefaultEnabled != nil {
		builder.SetDefaultEnabled(*updates.DefaultEnabled)
	}
	if updates.CooldownSeconds != nil {
		builder.SetCooldownSeconds(*updates.CooldownSeconds)
	}
	if updates.AdminOnly != nil {
		builder.SetAdminOnly(*updates.AdminOnly)
	}
	n, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	return r.GetByName(ctx, name)
}

func (r *entChannelConfigRepo) Delete(ctx context.Context, name string) error {
	_, err := r.client.ChannelConfig.Delete().Where(channelconfig.NameEQ(name)).Exec(ctx)
	return err
}
