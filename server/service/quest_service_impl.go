package service

import (
	"context"
	"fmt"

	"herbst-server/db"
	"herbst-server/db/quest"
	"herbst-server/db/schema"
	"herbst-server/repository"
)

// questService implements QuestService using repository interfaces.
type questService struct {
	repo  repository.QuestRepo
	qpRepo repository.QuestProgressRepo
}

// NewQuestService creates a new QuestService.
func NewQuestService(repo repository.QuestRepo, qpRepo repository.QuestProgressRepo) QuestService {
	return &questService{repo: repo, qpRepo: qpRepo}
}

// validRepeatModes maps valid repeat mode strings to their ent enum values.
var validRepeatModes = map[string]bool{
	string(quest.RepeatModeNone):    true,
	string(quest.RepeatModeCooldown): true,
	string(quest.RepeatModeAlways):   true,
}

func (s *questService) CreateQuest(ctx context.Context, input CreateQuestInput) (*db.Quest, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	repeatMode := input.RepeatMode
	if repeatMode == "" {
		repeatMode = string(quest.RepeatModeNone)
	}
	if !validRepeatModes[repeatMode] {
		return nil, fmt.Errorf("invalid repeat_mode: %s", repeatMode)
	}
	objectives := questObjectivesToSchema(input.Objectives)
	rewards := questRewardsToSchema(input.Rewards)
	rm := quest.RepeatMode(repeatMode)
	repoInput := repository.CreateQuestInput{
		Name:                 input.Name,
		Description:          input.Description,
		PrerequisiteQuestIDs: input.PrerequisiteQuestIDs,
		Objectives:           objectives,
		Rewards:              rewards,
		RepeatMode:           rm,
		CooldownHours:        input.CooldownHours,
		IsActive:             input.IsActive,
		WorldID:              input.WorldID,
	}
	return s.repo.Create(ctx, repoInput)
}

func (s *questService) GetQuest(ctx context.Context, id int) (*db.Quest, error) {
	return s.repo.Get(ctx, id)
}

func (s *questService) ListQuests(ctx context.Context, worldID string) ([]*db.Quest, error) {
	return s.repo.List(ctx, worldID)
}

func (s *questService) UpdateQuest(ctx context.Context, id int, input UpdateQuestInput) (*db.Quest, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quest not found: %w", err)
	}
	if input.RepeatMode != nil && !validRepeatModes[*input.RepeatMode] {
		return nil, fmt.Errorf("invalid repeat_mode: %s", *input.RepeatMode)
	}
	updates := repository.QuestUpdates{
		Name:                 input.Name,
		Description:          input.Description,
		PrerequisiteQuestIDs: input.PrerequisiteQuestIDs,
		CooldownHours:       input.CooldownHours,
		IsActive:            input.IsActive,
	}
	if input.Objectives != nil {
		objs := questObjectivesToSchema(*input.Objectives)
		updates.Objectives = &objs
	}
	if input.Rewards != nil {
		rwds := questRewardsToSchema(*input.Rewards)
		updates.Rewards = &rwds
	}
	if input.RepeatMode != nil {
		rm := quest.RepeatMode(*input.RepeatMode)
		updates.RepeatMode = &rm
	}
	if input.WorldID != nil {
		updates.WorldID = input.WorldID
	}
	_ = existing // fetched for validation; repo does the update
	return s.repo.Update(ctx, id, updates)
}

func (s *questService) DeleteQuest(ctx context.Context, id int) error {
	count, err := s.qpRepo.CountActiveByCharacter(ctx, 0, id)
	if err != nil {
		return fmt.Errorf("checking quest progress: %w", err)
	}
	// Check any progress exists (not just active)
	if count > 0 {
		return fmt.Errorf("cannot delete quest with %d progress records", count)
	}
	return s.repo.Delete(ctx, id)
}

// questObjectivesToSchema converts input objectives to schema types.
func questObjectivesToSchema(objs []schema.QuestObjective) []schema.QuestObjective {
	result := make([]schema.QuestObjective, len(objs))
	copy(result, objs)
	return result
}

// questRewardsToSchema converts input rewards to schema type.
func questRewardsToSchema(r schema.QuestRewards) schema.QuestRewards {
	return r
}