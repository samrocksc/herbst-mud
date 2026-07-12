package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/socialcommand"
)

type SocialCommandRepo interface {
	List(ctx context.Context, filter SocialCommandFilter) ([]*db.SocialCommand, error)
	Get(ctx context.Context, id int) (*db.SocialCommand, error)
	Create(ctx context.Context, input CreateSocialCommandInput) (*db.SocialCommand, error)
	Update(ctx context.Context, id int, updates SocialCommandUpdates) (*db.SocialCommand, error)
	Delete(ctx context.Context, id int) error
}

type SocialCommandFilter struct {
	Search string
}

type CreateSocialCommandInput struct {
	Name           string
	DisplayName    string
	SelfText       string
	RoomText       string
	TargetSelfText string
	TargetText     string
	TargetRoomText string
	RequiresTarget bool
	IsEmote        bool
	WorldID        string
}

type SocialCommandUpdates struct {
	Name           *string
	DisplayName    *string
	SelfText       *string
	RoomText       *string
	TargetSelfText *string
	TargetText     *string
	TargetRoomText *string
	RequiresTarget *bool
	IsEmote        *bool
}

type entSocialCommandRepo struct {
	client *db.Client
}

func NewEntSocialCommandRepo(client *db.Client) SocialCommandRepo {
	return &entSocialCommandRepo{client: client}
}

func (r *entSocialCommandRepo) List(ctx context.Context, filter SocialCommandFilter) ([]*db.SocialCommand, error) {
	query := r.client.SocialCommand.Query()
	if filter.Search != "" {
		query = query.Where(socialcommand.NameContains(filter.Search))
	}
	return query.All(ctx)
}

func (r *entSocialCommandRepo) Get(ctx context.Context, id int) (*db.SocialCommand, error) {
	return r.client.SocialCommand.Get(ctx, id)
}

func (r *entSocialCommandRepo) Create(ctx context.Context, input CreateSocialCommandInput) (*db.SocialCommand, error) {
	return r.client.SocialCommand.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName).
		SetSelfText(input.SelfText).
		SetRoomText(input.RoomText).
		SetTargetSelfText(input.TargetSelfText).
		SetTargetText(input.TargetText).
		SetTargetRoomText(input.TargetRoomText).
		SetRequiresTarget(input.RequiresTarget).
		SetIsEmote(input.IsEmote).
		SetWorldID(input.WorldID).
		Save(ctx)
}

func (r *entSocialCommandRepo) Update(ctx context.Context, id int, updates SocialCommandUpdates) (*db.SocialCommand, error) {
	builder := r.client.SocialCommand.UpdateOneID(id)
	if updates.Name != nil {
		builder.SetName(*updates.Name)
	}
	if updates.DisplayName != nil {
		builder.SetDisplayName(*updates.DisplayName)
	}
	if updates.SelfText != nil {
		builder.SetSelfText(*updates.SelfText)
	}
	if updates.RoomText != nil {
		builder.SetRoomText(*updates.RoomText)
	}
	if updates.TargetSelfText != nil {
		builder.SetTargetSelfText(*updates.TargetSelfText)
	}
	if updates.TargetText != nil {
		builder.SetTargetText(*updates.TargetText)
	}
	if updates.TargetRoomText != nil {
		builder.SetTargetRoomText(*updates.TargetRoomText)
	}
	if updates.RequiresTarget != nil {
		builder.SetRequiresTarget(*updates.RequiresTarget)
	}
	if updates.IsEmote != nil {
		builder.SetIsEmote(*updates.IsEmote)
	}
	return builder.Save(ctx)
}

func (r *entSocialCommandRepo) Delete(ctx context.Context, id int) error {
	return r.client.SocialCommand.DeleteOneID(id).Exec(ctx)
}
