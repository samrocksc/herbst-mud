package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
)

type entItemInstanceRepo struct {
	client *db.Client
}

func NewEntItemInstanceRepo(client *db.Client) ItemInstanceRepo {
	return &entItemInstanceRepo{client: client}
}

func (r *entItemInstanceRepo) ListNPCsByRoom(ctx context.Context, roomID int) ([]*db.Character, error) {
	return r.client.Character.Query().
		Where(
			character.CurrentRoomIdEQ(roomID),
			character.IsNPC(true),
		).
		All(ctx)
}