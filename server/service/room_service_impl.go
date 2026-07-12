package service

import (
	"context"
	"fmt"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/repository"
)

type roomService struct {
	roomRepo  repository.RoomRepo
	charRepo  repository.CharacterRepo
	equipRepo repository.EquipmentRepo
	npcRepo   repository.NPCTemplateRepo
	tx        repository.TransactionRunner
	zoneRepo  repository.ZoneRepository
}

func NewRoomService(
	roomRepo repository.RoomRepo,
	charRepo repository.CharacterRepo,
	equipRepo repository.EquipmentRepo,
	npcRepo repository.NPCTemplateRepo,
	tx repository.TransactionRunner,
	zoneRepo repository.ZoneRepository,
) RoomService {
	return &roomService{
		roomRepo:  roomRepo,
		charRepo:  charRepo,
		equipRepo: equipRepo,
		npcRepo:   npcRepo,
		tx:        tx,
		zoneRepo:  zoneRepo,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, input CreateRoomInput) (*db.Room, error) {
	if input.IsRootRoom {
		rooms, err := s.roomRepo.GetRoot(ctx)
		if err == nil && len(rooms) > 0 {
			for _, r := range rooms {
				_, _ = s.roomRepo.Update(ctx, r.ID, repository.RoomUpdates{IsRootRoom: boolPtr(false)})
			}
		}
	}
	repoInput := repository.CreateRoomInput{
		Name:           input.Name,
		Description:    input.Description,
		IsStartingRoom: input.IsStartingRoom,
		IsRootRoom:     input.IsRootRoom,
		Exits:          input.Exits,
		Atmosphere:     input.Atmosphere,
		PosZ:           input.PosZ,
		WorldID:        input.WorldID,
		ZoneIDs:        input.ZoneIDs,
	}
	created, err := s.roomRepo.Create(ctx, repoInput)
	if err != nil {
		return nil, err
	}
	// Sync Zone.room_ids for each zone the room was added to.
	if len(input.ZoneIDs) > 0 {
		for _, zid := range input.ZoneIDs {
			zone, err := s.zoneRepo.Get(ctx, zid)
			if err != nil {
				continue
			}
			newRoomIDs := append([]int{}, zone.RoomIds...)
			already := false
			for _, rid := range newRoomIDs {
				if rid == created.ID {
					already = true
					break
				}
			}
			if !already {
				newRoomIDs = append(newRoomIDs, created.ID)
				_, _ = s.zoneRepo.Update(ctx, zid, repository.ZoneUpdates{RoomIDs: &newRoomIDs})
			}
		}
	}
	return created, nil
}

func (s *roomService) GetRoom(ctx context.Context, id int) (*db.Room, error) {
	return s.roomRepo.Get(ctx, id)
}

func (s *roomService) ListRooms(ctx context.Context, worldID string) ([]*db.Room, error) {
	return s.roomRepo.List(ctx, worldID)
}

func (s *roomService) UpdateRoom(ctx context.Context, id int, input UpdateRoomInput) (*db.Room, error) {
	existing, err := s.roomRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}
	if input.Version != nil && *input.Version > 0 && *input.Version != existing.Version {
		return nil, fmt.Errorf("version conflict: expected %d, got %d", existing.Version, *input.Version)
	}
	if input.IsRootRoom != nil && *input.IsRootRoom {
		rooms, err := s.roomRepo.GetRoot(ctx)
		if err == nil && len(rooms) > 0 {
			for _, r := range rooms {
				if r.ID != id {
					_, _ = s.roomRepo.Update(ctx, r.ID, repository.RoomUpdates{IsRootRoom: boolPtr(false)})
				}
			}
		}
	}
	updates := repository.RoomUpdates{
		Name:           input.Name,
		Description:    input.Description,
		IsStartingRoom: input.IsStartingRoom,
		IsRootRoom:     input.IsRootRoom,
		Exits:          input.Exits,
		Atmosphere:     input.Atmosphere,
		PosZ:           input.PosZ,
		ZoneIDs:        input.ZoneIDs,
	}
	updated, err := s.roomRepo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	// Keep Zone.room_ids in sync with the room's zone_ids. When a room is
	// added to or removed from a zone, append/remove the room id from the
	// zone's explicit room_ids list.
	if input.ZoneIDs != nil {
		oldZoneSet := make(map[string]bool, len(existing.ZoneIds))
		for _, zid := range existing.ZoneIds {
			oldZoneSet[zid] = true
		}
		newZoneSet := make(map[string]bool, len(*input.ZoneIDs))
		for _, zid := range *input.ZoneIDs {
			newZoneSet[zid] = true
		}
		// For each zone in old or new set, sync.
		allZones := make(map[string]bool, len(oldZoneSet)+len(newZoneSet))
		for z := range oldZoneSet {
			allZones[z] = true
		}
		for z := range newZoneSet {
			allZones[z] = true
		}
		for zid := range allZones {
			zone, err := s.zoneRepo.Get(ctx, zid)
			if err != nil {
				continue
			}
			var newRoomIDs []int
			hasRoom := false
			for _, rid := range zone.RoomIds {
				if rid == id {
					hasRoom = true
					continue
				}
				newRoomIDs = append(newRoomIDs, rid)
			}
			shouldBeInZone := newZoneSet[zid]
			if shouldBeInZone && !hasRoom {
				newRoomIDs = append(newRoomIDs, id)
			} else if !shouldBeInZone && hasRoom {
				// already removed above
			}
			_, _ = s.zoneRepo.Update(ctx, zid, repository.ZoneUpdates{RoomIDs: &newRoomIDs})
		}
	}
	return updated, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, id int) error {
	// Find the right "default" room for character relocation in this room's world.
	// Priority: the world's root room > any room in the same world > the original
	// (hardcoded) fallback if the world is somehow empty.
	defaultRoomID := 5
	if room, err := s.roomRepo.Get(ctx, id); err == nil && room != nil {
		if rootRooms, err := s.roomRepo.GetRoot(ctx); err == nil && len(rootRooms) > 0 {
			// Prefer a root room in the same world
			for _, r := range rootRooms {
				if r.WorldID == room.WorldID {
					defaultRoomID = r.ID
					break
				}
			}
			// Fall back to any root room
			if defaultRoomID == 5 {
				defaultRoomID = rootRooms[0].ID
			}
		} else {
			// No root at all — pick any room in the same world
			if sameWorld, err := s.roomRepo.List(ctx, room.WorldID); err == nil && len(sameWorld) > 0 {
				for _, r := range sameWorld {
					if r.ID != id {
						defaultRoomID = r.ID
						break
					}
				}
			}
		}
	}
	err := s.tx.WithTx(ctx, func(tx *db.Tx) error {
		_, err := tx.Room.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("room not found: %w", err)
		}
		_, err = tx.Character.Update().
			Where(character.CurrentRoomIdEQ(id)).
			SetCurrentRoomId(defaultRoomID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to relocate characters: %w", err)
		}
		rooms, err := tx.Room.Query().All(ctx)
		if err != nil {
			return fmt.Errorf("failed to load rooms for exit cleanup: %w", err)
		}
		for _, r := range rooms {
			if r.ID == id {
				continue
			}
			exits := r.Exits
			if exits == nil {
				continue
			}
			newExits := make(map[string]int)
			changed := false
			for dir, targetID := range exits {
				if targetID == id {
					changed = true
					continue
				}
				newExits[dir] = targetID
			}
			if changed {
				_, err := tx.Room.UpdateOneID(r.ID).
					SetExits(newExits).
					AddVersion(1).
					Save(ctx)
				if err != nil {
					continue
				}
			}
		}
		return tx.Room.DeleteOneID(id).Exec(ctx)
	})
	return err
}

// CleanupOrphanExits removes exits pointing to non-existent rooms.
//
// If worldID is empty (""), it operates on all rooms across all worlds.
// This works because roomRepo.List with "" returns all rooms regardless of world.
// If a worldID is provided, only exits within that world are cleaned.
func (s *roomService) CleanupOrphanExits(ctx context.Context, worldID string) (int, error) {
	rooms, err := s.roomRepo.List(ctx, worldID)
	if err != nil {
		return 0, err
	}
	validIDs := make(map[int]bool, len(rooms))
	for _, r := range rooms {
		validIDs[r.ID] = true
	}
	cleaned := 0
	for _, r := range rooms {
		exits := r.Exits
		if exits == nil || len(exits) == 0 {
			continue
		}
		newExits := make(map[string]int)
		changed := false
		for dir, targetID := range exits {
			if !validIDs[targetID] {
				changed = true
				continue
			}
			newExits[dir] = targetID
		}
		if changed {
			_, err := s.roomRepo.Update(ctx, r.ID, repository.RoomUpdates{Exits: &newExits})
			if err != nil {
				continue
			}
			cleaned += len(exits) - len(newExits)
		}
	}
	return cleaned, nil
}

var oppositeDir = map[string]string{
	"north":     "south",
	"south":     "north",
	"east":      "west",
	"west":      "east",
	"northeast": "southwest",
	"southwest": "northeast",
	"northwest": "southeast",
	"southeast": "northwest",
	"up":        "down",
	"down":      "up",
}

type BidirectionalExitResult struct {
	Source *db.Room `json:"source"`
	Target *db.Room `json:"target"`
}

func (s *roomService) CreateBidirectionalExit(ctx context.Context, sourceID int, direction string, targetID int) (*BidirectionalExitResult, error) {
	reverseDir, ok := oppositeDir[direction]
	if !ok {
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}
	source, err := s.roomRepo.Get(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("source room not found: %w", err)
	}
	target, err := s.roomRepo.Get(ctx, targetID)
	if err != nil {
		return nil, fmt.Errorf("target room not found: %w", err)
	}
	sourceExits := source.Exits
	if sourceExits == nil {
		sourceExits = map[string]int{}
	}
	sourceExits[direction] = targetID
	targetExits := target.Exits
	if targetExits == nil {
		targetExits = map[string]int{}
	}
	targetExits[reverseDir] = sourceID
	source, err = s.roomRepo.Update(ctx, sourceID, repository.RoomUpdates{Exits: &sourceExits})
	if err != nil {
		return nil, err
	}
	target, err = s.roomRepo.Update(ctx, targetID, repository.RoomUpdates{Exits: &targetExits})
	if err != nil {
		return nil, err
	}
	return &BidirectionalExitResult{Source: source, Target: target}, nil
}

func (s *roomService) DeleteBidirectionalExit(ctx context.Context, sourceID int, direction string) error {
	reverseDir, ok := oppositeDir[direction]
	if !ok {
		return fmt.Errorf("invalid direction: %s", direction)
	}
	source, err := s.roomRepo.Get(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("source room not found: %w", err)
	}
	sourceExits := source.Exits
	if sourceExits == nil {
		return nil
	}
	targetID, hasTarget := sourceExits[direction]
	delete(sourceExits, direction)
	_, err = s.roomRepo.Update(ctx, sourceID, repository.RoomUpdates{Exits: &sourceExits})
	if err != nil {
		return err
	}
	if hasTarget && targetID > 0 {
		target, err := s.roomRepo.Get(ctx, targetID)
		if err != nil {
			return nil
		}
		targetExits := target.Exits
		if targetExits != nil {
			delete(targetExits, reverseDir)
			_, _ = s.roomRepo.Update(ctx, targetID, repository.RoomUpdates{Exits: &targetExits})
		}
	}
	return nil
}

func boolPtr(b bool) *bool { return &b }
func intPtr(i int) *int     { return &i }
