package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/quest"
)

type entQuestRepo struct {
	client *db.Client
}

func NewEntQuestRepo(client *db.Client) QuestRepo {
	return &entQuestRepo{client: client}
}

func (r *entQuestRepo) Get(ctx context.Context, id int) (*db.Quest, error) {
	return r.client.Quest.Get(ctx, id)
}

func (r *entQuestRepo) List(ctx context.Context, worldID string) ([]*db.Quest, error) {
	query := r.client.Quest.Query()
	if worldID != "" {
		query = query.Where(quest.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entQuestRepo) Create(ctx context.Context, input CreateQuestInput) (*db.Quest, error) {
	builder := r.client.Quest.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetPrerequisiteQuestIds(input.PrerequisiteQuestIDs).
		SetObjectives(input.Objectives).
		SetRewards(input.Rewards).
		SetRepeatMode(input.RepeatMode).
		SetCooldownHours(input.CooldownHours).
		SetIsActive(input.IsActive).
		SetWorldID(input.WorldID)
	return builder.Save(ctx)
}

func (r *entQuestRepo) Update(ctx context.Context, id int, updates QuestUpdates) (*db.Quest, error) {
	builder := r.client.Quest.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.PrerequisiteQuestIDs != nil {
		builder = builder.SetPrerequisiteQuestIds(*updates.PrerequisiteQuestIDs)
	}
	if updates.Objectives != nil {
		builder = builder.SetObjectives(*updates.Objectives)
	}
	if updates.Rewards != nil {
		builder = builder.SetRewards(*updates.Rewards)
	}
	if updates.RepeatMode != nil {
		builder = builder.SetRepeatMode(*updates.RepeatMode)
	}
	if updates.CooldownHours != nil {
		builder = builder.SetCooldownHours(*updates.CooldownHours)
	}
	if updates.IsActive != nil {
		builder = builder.SetIsActive(*updates.IsActive)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entQuestRepo) Delete(ctx context.Context, id int) error {
	return r.client.Quest.DeleteOneID(id).Exec(ctx)
}