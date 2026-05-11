package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/characterfaction"
	"herbst-server/db/faction"
)

type entCharacterFactionRepo struct {
	client *db.Client
}

func NewEntCharacterFactionRepo(client *db.Client) CharacterFactionRepo {
	return &entCharacterFactionRepo{client: client}
}

func (r *entCharacterFactionRepo) ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterFaction, error) {
	return r.client.CharacterFaction.Query().
		Where(characterfaction.HasCharacterWith(character.ID(charID))).
		All(ctx)
}

func (r *entCharacterFactionRepo) ListByFactionWithDetails(ctx context.Context, factionID int) ([]*db.CharacterFaction, error) {
	return r.client.CharacterFaction.Query().
		Where(characterfaction.HasFactionWith(faction.ID(factionID))).
		WithCharacter().
		WithFaction().
		All(ctx)
}

func (r *entCharacterFactionRepo) Create(ctx context.Context, charID int, factionID int, reputation int) (*db.CharacterFaction, error) {
	return r.client.CharacterFaction.Create().
		SetCharacterID(charID).
		SetFactionID(factionID).
		SetReputation(reputation).
		Save(ctx)
}

func (r *entCharacterFactionRepo) Delete(ctx context.Context, id int) error {
	return r.client.CharacterFaction.DeleteOneID(id).Exec(ctx)
}