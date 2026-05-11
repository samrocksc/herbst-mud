package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactertag"
)

type entCharacterTagRepo struct {
	client *db.Client
}

func NewEntCharacterTagRepo(client *db.Client) CharacterTagRepo {
	return &entCharacterTagRepo{client: client}
}

func (r *entCharacterTagRepo) ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterTag, error) {
	return r.client.CharacterTag.Query().
		Where(charactertag.HasCharacterWith(character.ID(charID))).
		All(ctx)
}

func (r *entCharacterTagRepo) Create(ctx context.Context, charID int, tag, source string) (*db.CharacterTag, error) {
	return r.client.CharacterTag.Create().
		SetCharacterID(charID).
		SetTag(tag).
		SetSource(source).
		Save(ctx)
}

func (r *entCharacterTagRepo) Delete(ctx context.Context, id int) error {
	return r.client.CharacterTag.DeleteOneID(id).Exec(ctx)
}