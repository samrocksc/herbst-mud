package repository

import (
	"context"

	"herbst-server/db"
)

type entAchievementRepo struct {
	client *db.Client
}

func NewEntAchievementRepo(client *db.Client) AchievementRepo {
	return &entAchievementRepo{client: client}
}

func (r *entAchievementRepo) Get(ctx context.Context, id int) (*db.Achievement, error) {
	return r.client.Achievement.Get(ctx, id)
}

func (r *entAchievementRepo) List(ctx context.Context) ([]*db.Achievement, error) {
	return r.client.Achievement.Query().All(ctx)
}

func (r *entAchievementRepo) Create(ctx context.Context, input CreateAchievementInput) (*db.Achievement, error) {
	builder := r.client.Achievement.Create().
		SetName(input.Name).
		SetXpReward(input.XPReward).
		SetCriteria(input.Criteria)
	if input.Description != "" {
		builder = builder.SetDescription(input.Description)
	}
	if input.Icon != "" {
		builder = builder.SetIcon(input.Icon)
	}
	return builder.Save(ctx)
}

func (r *entAchievementRepo) Update(ctx context.Context, id int, updates AchievementUpdates) (*db.Achievement, error) {
	builder := r.client.Achievement.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.Icon != nil {
		builder = builder.SetIcon(*updates.Icon)
	}
	if updates.XPReward != nil {
		builder = builder.SetXpReward(*updates.XPReward)
	}
	if updates.Criteria != nil {
		builder = builder.SetCriteria(*updates.Criteria)
	}
	return builder.Save(ctx)
}

func (r *entAchievementRepo) Delete(ctx context.Context, id int) error {
	return r.client.Achievement.DeleteOneID(id).Exec(ctx)
}