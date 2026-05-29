package service

import (
	"context"

	"herbst-server/db"
	"herbst-server/repository"
)

// npcService implements NPCService using repository interfaces.
type npcService struct {
	repo repository.NPCTemplateRepo
}

// NewNPCService creates a new NPCService.
func NewNPCService(repo repository.NPCTemplateRepo) NPCService {
	return &npcService{repo: repo}
}

func (s *npcService) GetTemplate(ctx context.Context, id string) (*db.NPCTemplate, error) {
	return s.repo.Get(ctx, id)
}

func (s *npcService) ListTemplates(ctx context.Context, worldID string) ([]*db.NPCTemplate, error) {
	return s.repo.List(ctx, worldID)
}

func (s *npcService) CreateTemplate(ctx context.Context, input CreateNPCTemplateInput) (*db.NPCTemplate, error) {
	raceID := 0
	if input.Race != "" {
		// TODO: resolve race name to ID
	}
	repoInput := repository.CreateNPCTemplateInput{
		ID:              input.ID,
		Name:            input.Name,
		Description:     input.Description,
		RaceID:          raceID,
		Disposition:     input.Disposition,
		Level:           input.Level,
		XPValue:         input.XPValue,
		Skills:          input.Skills,
		TradesWith:      input.TradesWith,
		Greeting:        input.Greeting,
		RespawnRooms:    input.RespawnRooms,
		RespawnCooldown: input.RespawnCooldown,
	}
	return s.repo.Create(ctx, repoInput)
}

func (s *npcService) UpdateTemplate(ctx context.Context, id string, input UpdateNPCTemplateInput) (*db.NPCTemplate, error) {
	updates := repository.NPCTemplateUpdates{
		Name:             input.Name,
		Description:      input.Description,
		Disposition:      input.Disposition,
		Level:            input.Level,
		XPValue:          input.XPValue,
		Skills:           input.Skills,
		TradesWith:       input.TradesWith,
		Greeting:         input.Greeting,
		RespawnRooms:     input.RespawnRooms,
		RespawnCooldown:  input.RespawnCooldown,
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *npcService) DeleteTemplate(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
