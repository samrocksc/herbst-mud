package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/zone"
)

type ZoneRepoImpl struct{}

func NewZoneRepoImpl() *ZoneRepoImpl {
	return &ZoneRepoImpl{}
}

type ZoneRepository interface {
	Create(ctx context.Context, input CreateZoneInput) (*db.Zone, error)
	Get(ctx context.Context, id string) (*db.Zone, error)
	GetByName(ctx context.Context, name, worldID string) (*db.Zone, error)
	ListByWorld(ctx context.Context, worldID string) ([]*db.Zone, error)
	Update(ctx context.Context, id string, updates ZoneUpdates) (*db.Zone, error)
	Delete(ctx context.Context, id string) error
	GetByParent(ctx context.Context, parentZoneID string) ([]*db.Zone, error)
}

type entZoneRepo struct {
	client *db.Client
}

func NewEntZoneRepo(client *db.Client) ZoneRepository {
	return &entZoneRepo{client: client}
}

func (r *entZoneRepo) Create(ctx context.Context, input CreateZoneInput) (*db.Zone, error) {
	id := input.ID
	if id == "" {
		panic("Zone.ID is required")
	}
	builder := r.client.Zone.Create().
		SetID(id).
		SetName(input.Name).
		SetWorldID(input.WorldID).
		SetMinLevel(input.MinLevel).
		SetColor(input.Color)
	if input.Description != "" {
		builder = builder.SetDescription(input.Description)
	}
	if input.ParentZoneID != "" {
		builder = builder.SetParentZoneID(input.ParentZoneID)
	}
	if len(input.RoomIDs) > 0 {
		builder = builder.SetRoomIds(input.RoomIDs)
	}
	return builder.Save(ctx)
}

func (r *entZoneRepo) Get(ctx context.Context, id string) (*db.Zone, error) {
	return r.client.Zone.Get(ctx, id)
}

func (r *entZoneRepo) GetByName(ctx context.Context, name, worldID string) (*db.Zone, error) {
	return r.client.Zone.Query().
		Where(zone.NameEQ(name), zone.WorldIDEQ(worldID)).
		Only(ctx)
}

func (r *entZoneRepo) ListByWorld(ctx context.Context, worldID string) ([]*db.Zone, error) {
	return r.client.Zone.Query().
		Where(zone.WorldIDEQ(worldID)).
		All(ctx)
}

func (r *entZoneRepo) Update(ctx context.Context, id string, updates ZoneUpdates) (*db.Zone, error) {
	builder := r.client.Zone.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.MinLevel != nil {
		builder = builder.SetMinLevel(*updates.MinLevel)
	}
	if updates.ParentZoneID != nil {
		builder = builder.SetParentZoneID(*updates.ParentZoneID)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.RoomIDs != nil {
		builder = builder.SetRoomIds(*updates.RoomIDs)
	}
	return builder.Save(ctx)
}

func (r *entZoneRepo) Delete(ctx context.Context, id string) error {
	// Detach all rooms that reference this zone.
	// ent's `field.Strings` does not generate a Contains predicate, so
	// we load all rooms in the world's zone and filter in Go. This is
	// not a hot path — Delete is called rarely and world-scoped.
	allRooms, err := r.client.Room.Query().All(ctx)
	if err != nil {
		return err
	}
	var rooms []*db.Room
	for _, rm := range allRooms {
		for _, zid := range rm.ZoneIds {
			if zid == id {
				rooms = append(rooms, rm)
				break
			}
		}
	}
	_ = err // re-uses err below
	for _, rm := range rooms {
		// Remove this zone from zone_ids slice.
		var newIDs []string
		for _, zid := range rm.ZoneIds {
			if zid != id {
				newIDs = append(newIDs, zid)
			}
		}
		if len(newIDs) == 0 {
			_, err = r.client.Room.UpdateOneID(rm.ID).ClearZoneIds().Save(ctx)
		} else {
			_, err = r.client.Room.UpdateOneID(rm.ID).SetZoneIds(newIDs).Save(ctx)
		}
		if err != nil {
			return err
		}
	}
	// Clear parent_zone_id on child zones.
	children, err := r.client.Zone.Query().Where(zone.ParentZoneIDEQ(id)).All(ctx)
	if err != nil {
		return err
	}
	for _, z := range children {
		_, err = r.client.Zone.UpdateOneID(z.ID).ClearParentZoneID().Save(ctx)
		if err != nil {
			return err
		}
	}
	return r.client.Zone.DeleteOneID(id).Exec(ctx)
}

func (r *entZoneRepo) GetByParent(ctx context.Context, parentZoneID string) ([]*db.Zone, error) {
	return r.client.Zone.Query().
		Where(zone.ParentZoneIDEQ(parentZoneID)).
		All(ctx)
}

type CreateZoneInput struct {
	ID           string
	WorldID      string
	Name         string
	Description  string
	MinLevel     int
	ParentZoneID string
	Color        string
	RoomIDs      []int
}

type ZoneUpdates struct {
	Name         *string
	Description  *string
	MinLevel     *int
	ParentZoneID *string
	Color        *string
	RoomIDs      *[]int
}
