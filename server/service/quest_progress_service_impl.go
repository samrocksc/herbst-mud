package service

import (
	"context"
	"fmt"
	"time"

	"herbst-server/db"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
	"herbst-server/repository"
)

type questProgressService struct {
	qpRepo  repository.QuestProgressRepo
	qRepo   repository.QuestRepo
	charRepo repository.CharacterRepo
}

func NewQuestProgressService(
	qpRepo repository.QuestProgressRepo,
	qRepo repository.QuestRepo,
	charRepo repository.CharacterRepo,
) QuestProgressService {
	return &questProgressService{qpRepo: qpRepo, qRepo: qRepo, charRepo: charRepo}
}

func (s *questProgressService) Accept(ctx context.Context, charID int, questID int) (*QuestProgressView, error) {
	if _, err := s.charRepo.Get(ctx, charID); err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	q, err := s.qRepo.Get(ctx, questID)
	if err != nil {
		return nil, fmt.Errorf("quest not found: %w", err)
	}
	if !q.IsActive {
		return nil, fmt.Errorf("quest is not active")
	}
	count, err := s.qpRepo.CountActiveByCharacter(ctx, charID, questID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("quest already active for this character: %w", ErrQuestAlreadyActive)
	}
	if err := s.validatePrerequisites(ctx, charID, q); err != nil {
		return nil, err
	}
	if err := s.validateCooldown(ctx, charID, q); err != nil {
		return nil, err
	}
	progress, err := s.qpRepo.Create(ctx, repository.CreateQuestProgressInput{
		CharacterID:     charID,
		QuestID:         questID,
		Status:          questprogress.StatusActive,
		StartedAt:       time.Now(),
		CurrentStep:     0,
		ObjectiveCounts: map[string]int{},
	})
	if err != nil {
		return nil, err
	}
	full, err := s.qpRepo.GetWithRelations(ctx, progress.ID)
	if err != nil {
		return nil, err
	}
	return questProgressToView(full), nil
}

func (s *questProgressService) Advance(ctx context.Context, charID int, questID int, objectiveKey string, count int) (*QuestProgressView, error) {
	progress, err := s.qpRepo.GetWithRelations(ctx, 0)
	_ = progress
	_ = err
	return nil, fmt.Errorf("not implemented: use CheckAll or check endpoint")
}

func (s *questProgressService) Abandon(ctx context.Context, charID int, questID int) error {
	if _, err := s.charRepo.Get(ctx, charID); err != nil {
		return fmt.Errorf("character not found: %w", err)
	}
	return fmt.Errorf("not implemented: abandon quest")
}

func (s *questProgressService) CheckAll(ctx context.Context, charID int, objectiveType, targetID string) ([]QuestProgressView, error) {
	return nil, fmt.Errorf("not implemented: check all")
}

func (s *questProgressService) ListByCharacter(ctx context.Context, charID int) ([]QuestProgressView, error) {
	if _, err := s.charRepo.Get(ctx, charID); err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}
	progressList, err := s.qpRepo.ListByCharacter(ctx, charID)
	if err != nil {
		return nil, err
	}
	views := make([]QuestProgressView, 0, len(progressList))
	for _, p := range progressList {
		views = append(views, *questProgressToView(p))
	}
	return views, nil
}

func (s *questProgressService) ValidateAcceptance(ctx context.Context, charID int, questID int) error {
	q, err := s.qRepo.Get(ctx, questID)
	if err != nil {
		return fmt.Errorf("quest not found: %w", err)
	}
	if !q.IsActive {
		return fmt.Errorf("quest is not active")
	}
	count, err := s.qpRepo.CountActiveByCharacter(ctx, charID, questID)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrQuestAlreadyActive
	}
	if err := s.validatePrerequisites(ctx, charID, q); err != nil {
		return err
	}
	return s.validateCooldown(ctx, charID, q)
}

// Sentinel errors for quest progress.
var (
	ErrQuestAlreadyActive = fmt.Errorf("quest already active for this character")
	ErrQuestNotFound      = fmt.Errorf("quest not found")
	ErrQuestNotActive     = fmt.Errorf("quest is not active")
	ErrPrerequisites      = fmt.Errorf("prerequisites not met")
	ErrCooldownActive     = fmt.Errorf("quest is on cooldown")
)

func (s *questProgressService) validatePrerequisites(ctx context.Context, charID int, q *db.Quest) error {
	if len(q.PrerequisiteQuestIds) == 0 {
		return nil
	}
	for _, prereqID := range q.PrerequisiteQuestIds {
		var prereqInt int
		if _, err := fmt.Sscanf(prereqID, "%d", &prereqInt); err != nil {
			continue
		}
		count, err := s.qpRepo.CountActiveByCharacter(ctx, charID, prereqInt)
		if err != nil || count == 0 {
			return fmt.Errorf("prerequisite quest %s not completed: %w", prereqID, ErrPrerequisites)
		}
	}
	return nil
}

func (s *questProgressService) validateCooldown(ctx context.Context, charID int, q *db.Quest) error {
	if q.RepeatMode == quest.RepeatModeNone {
		return nil
	}
	// Cooldown logic will be expanded in Phase 2+; for now, non-none modes pass
	return nil
}

// questProgressToView converts a db.QuestProgress to a QuestProgressView.
func questProgressToView(p *db.QuestProgress) *QuestProgressView {
	if p == nil {
		return nil
	}
	view := &QuestProgressView{
		ID:              p.ID,
		Status:          string(p.Status),
		CurrentStep:     p.CurrentStep,
		ObjectiveCounts: p.ObjectiveCounts,
	}
	if p.Edges.Character != nil {
		view.CharacterID = p.Edges.Character.ID
	}
	if p.Edges.Quest != nil {
		view.QuestID = p.Edges.Quest.ID
		view.QuestName = p.Edges.Quest.Name
		view.QuestDescription = p.Edges.Quest.Description
	}
	if !p.StartedAt.IsZero() {
		view.StartedAt = p.StartedAt.Format(time.RFC3339)
	}
	if p.CompletedAt != nil {
		t := p.CompletedAt.Format(time.RFC3339)
		view.CompletedAt = &t
	}
	if view.ObjectiveCounts == nil {
		view.ObjectiveCounts = map[string]int{}
	}
	return view
}