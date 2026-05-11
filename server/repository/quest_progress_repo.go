package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
)

type entQuestProgressRepo struct {
	client *db.Client
}

func NewEntQuestProgressRepo(client *db.Client) QuestProgressRepo {
	return &entQuestProgressRepo{client: client}
}

func (r *entQuestProgressRepo) Get(ctx context.Context, id int) (*db.QuestProgress, error) {
	return r.client.QuestProgress.Get(ctx, id)
}

func (r *entQuestProgressRepo) GetWithRelations(ctx context.Context, id int) (*db.QuestProgress, error) {
	return r.client.QuestProgress.Query().
		Where(questprogress.ID(id)).
		WithCharacter().
		WithQuest().
		Only(ctx)
}

func (r *entQuestProgressRepo) ListByCharacter(ctx context.Context, charID int) ([]*db.QuestProgress, error) {
	return r.client.QuestProgress.Query().
		Where(questprogress.HasCharacterWith(character.ID(charID))).
		All(ctx)
}

func (r *entQuestProgressRepo) Create(ctx context.Context, input CreateQuestProgressInput) (*db.QuestProgress, error) {
	builder := r.client.QuestProgress.Create().
		SetCharacterID(input.CharacterID).
		SetQuestID(input.QuestID).
		SetStatus(input.Status).
		SetStartedAt(input.StartedAt).
		SetCurrentStep(input.CurrentStep).
		SetObjectiveCounts(input.ObjectiveCounts)
	return builder.Save(ctx)
}

func (r *entQuestProgressRepo) Update(ctx context.Context, id int, updates QuestProgressUpdates) (*db.QuestProgress, error) {
	builder := r.client.QuestProgress.UpdateOneID(id)
	if updates.Status != nil {
		builder = builder.SetStatus(*updates.Status)
	}
	if updates.CurrentStep != nil {
		builder = builder.SetCurrentStep(*updates.CurrentStep)
	}
	if updates.ObjectiveCounts != nil {
		builder = builder.SetObjectiveCounts(*updates.ObjectiveCounts)
	}
	if updates.CompletedAt != nil {
		builder = builder.SetCompletedAt(*updates.CompletedAt)
	}
	return builder.Save(ctx)
}

func (r *entQuestProgressRepo) Delete(ctx context.Context, id int) error {
	return r.client.QuestProgress.DeleteOneID(id).Exec(ctx)
}

func (r *entQuestProgressRepo) CountActiveByCharacter(ctx context.Context, charID int, questID int) (int, error) {
	return r.client.QuestProgress.Query().
		Where(
			questprogress.HasCharacterWith(character.ID(charID)),
			questprogress.HasQuestWith(quest.ID(questID)),
		).
		Count(ctx)
}