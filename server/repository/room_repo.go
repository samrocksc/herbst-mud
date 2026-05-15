package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/room"
)

type entRoomRepo struct {
	client *db.Client
}

func NewEntRoomRepo(client *db.Client) RoomRepo {
	return &entRoomRepo{client: client}
}

func (r *entRoomRepo) Get(ctx context.Context, id int) (*db.Room, error) {
	return r.client.Room.Get(ctx, id)
}

func (r *entRoomRepo) List(ctx context.Context, worldID string) ([]*db.Room, error) {
	query := r.client.Room.Query()
	if worldID != "" {
		query = query.Where(room.WorldID(worldID))
	}
	return query.All(ctx)
}

func (r *entRoomRepo) GetRoot(ctx context.Context) ([]*db.Room, error) {
	return r.client.Room.Query().
		Where(room.IsRootRoom(true)).
		All(ctx)
}

func (r *entRoomRepo) Create(ctx context.Context, input CreateRoomInput) (*db.Room, error) {
	builder := r.client.Room.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetIsStartingRoom(input.IsStartingRoom).
		SetIsRootRoom(input.IsRootRoom).
		SetExits(input.Exits).
		SetPosX(input.PosX).
		SetPosY(input.PosY).
		SetPosZ(input.PosZ).
		SetWorldID(input.WorldID)
	if input.Atmosphere != "" {
		builder = builder.SetAtmosphere(room.Atmosphere(input.Atmosphere))
	}
	return builder.Save(ctx)
}

func (r *entRoomRepo) Update(ctx context.Context, id int, updates RoomUpdates) (*db.Room, error) {
	builder := r.client.Room.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.IsStartingRoom != nil {
		builder = builder.SetIsStartingRoom(*updates.IsStartingRoom)
	}
	if updates.IsRootRoom != nil {
		builder = builder.SetIsRootRoom(*updates.IsRootRoom)
	}
	if updates.Exits != nil {
		builder = builder.SetExits(*updates.Exits)
	}
	if updates.Atmosphere != nil {
		builder = builder.SetAtmosphere(room.Atmosphere(*updates.Atmosphere))
	}
	if updates.PosX != nil {
		builder = builder.SetPosX(*updates.PosX)
	}
	if updates.PosY != nil {
		builder = builder.SetPosY(*updates.PosY)
	}
	if updates.PosZ != nil {
		builder = builder.SetPosZ(*updates.PosZ)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entRoomRepo) Delete(ctx context.Context, id int) error {
	return r.client.Room.DeleteOneID(id).Exec(ctx)
}