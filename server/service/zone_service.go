package service

import (
	"context"
	"fmt"
	"regexp"

	"herbst-server/db"
	"herbst-server/repository"
)

var colorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{3}){1,2}$`)

type ZoneService struct {
	zoneRepo    repository.ZoneRepository
	npcTemplate repository.NPCTemplateRepo
}

func NewZoneService(zoneRepo repository.ZoneRepository, npcTemplate repository.NPCTemplateRepo) *ZoneService {
	return &ZoneService{
		zoneRepo:    zoneRepo,
		npcTemplate: npcTemplate,
	}
}

func (s *ZoneService) CreateZone(ctx context.Context, input repository.CreateZoneInput) (*db.Zone, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.WorldID == "" {
		return nil, fmt.Errorf("world_id is required")
	}
	existing, err := s.zoneRepo.GetByName(ctx, input.Name, input.WorldID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("zone with name %q already exists in world %s", input.Name, input.WorldID)
	}
	if input.ParentZoneID != "" {
		parent, err := s.zoneRepo.Get(ctx, input.ParentZoneID)
		if err != nil {
			return nil, fmt.Errorf("parent zone not found: %s", input.ParentZoneID)
		}
		if parent.WorldID != input.WorldID {
			return nil, fmt.Errorf("parent zone must be in the same world")
		}
	}
	if input.Color != "" && !colorRegex.MatchString(input.Color) {
		return nil, fmt.Errorf("color must be a valid hex color (e.g. #ff0000)")
	}
	if input.MinLevel < 0 {
		return nil, fmt.Errorf("min_level must be non-negative")
	}
	return s.zoneRepo.Create(ctx, input)
}

func (s *ZoneService) GetZone(ctx context.Context, id string) (*db.Zone, error) {
	return s.zoneRepo.Get(ctx, id)
}

func (s *ZoneService) ListZonesByWorld(ctx context.Context, worldID string) ([]*db.Zone, error) {
	return s.zoneRepo.ListByWorld(ctx, worldID)
}

func (s *ZoneService) UpdateZone(ctx context.Context, id string, updates repository.ZoneUpdates) (*db.Zone, error) {
	if updates.Name != nil && *updates.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if updates.ParentZoneID != nil && *updates.ParentZoneID != "" {
		parent, err := s.zoneRepo.Get(ctx, *updates.ParentZoneID)
		if err != nil {
			return nil, fmt.Errorf("parent zone not found: %s", *updates.ParentZoneID)
		}
		existing, err := s.zoneRepo.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if parent.WorldID != existing.WorldID {
			return nil, fmt.Errorf("parent zone must be in the same world")
		}
		if parent.ID == id {
			return nil, fmt.Errorf("zone cannot be its own parent")
		}
	}
	if updates.Color != nil && *updates.Color != "" && !colorRegex.MatchString(*updates.Color) {
		return nil, fmt.Errorf("color must be a valid hex color (e.g. #ff0000)")
	}
	if updates.MinLevel != nil && *updates.MinLevel < 0 {
		return nil, fmt.Errorf("min_level must be non-negative")
	}
	return s.zoneRepo.Update(ctx, id, updates)
}

func (s *ZoneService) DeleteZone(ctx context.Context, id string) error {
	return s.zoneRepo.Delete(ctx, id)
}
