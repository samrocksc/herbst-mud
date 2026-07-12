package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
	"herbst-server/db/characterskill"
	"herbst-server/db/user"
)

type entCharacterRepo struct {
	client *db.Client
}

func NewEntCharacterRepo(client *db.Client) CharacterRepo {
	return &entCharacterRepo{client: client}
}

func (r *entCharacterRepo) Get(ctx context.Context, id int) (*db.Character, error) {
	return r.client.Character.Get(ctx, id)
}

func (r *entCharacterRepo) GetByName(ctx context.Context, name string) (*db.Character, error) {
	return r.client.Character.Query().
		Where(character.NameEQ(name)).
		Only(ctx)
}

func (r *entCharacterRepo) ListByUser(ctx context.Context, userID int) ([]*db.Character, error) {
	return r.client.Character.Query().
		Where(character.HasUserWith(user.ID(userID))).
		All(ctx)
}

func (r *entCharacterRepo) ListByRoom(ctx context.Context, roomID int) ([]*db.Character, error) {
	return r.client.Character.Query().
		Where(character.CurrentRoomIdEQ(roomID)).
		All(ctx)
}

func (r *entCharacterRepo) ListNPCsByRoom(ctx context.Context, roomID int) ([]*db.Character, error) {
	return r.client.Character.Query().
		Where(character.IsNPCEQ(true), character.CurrentRoomIdEQ(roomID)).
		All(ctx)
}

func (r *entCharacterRepo) ListAllNPCs(ctx context.Context) ([]*db.Character, error) {
	return r.client.Character.Query().
		Where(character.IsNPCEQ(true)).
		All(ctx)
}

func (r *entCharacterRepo) ListAll(ctx context.Context) ([]*db.Character, error) {
	return r.client.Character.Query().All(ctx)
}

func (r *entCharacterRepo) ListAllByWorld(ctx context.Context, worldID string) ([]*db.Character, error) {
	query := r.client.Character.Query().
		Where(character.IsNPCEQ(true))
	if worldID != "" {
		query = query.Where(character.CurrentWorldEQ(worldID))
	}
	return query.All(ctx)
}

func (r *entCharacterRepo) CountByUser(ctx context.Context, userID int) (int, error) {
	return r.client.Character.Query().
		Where(character.HasUserWith(user.ID(userID))).
		Count(ctx)
}

func (r *entCharacterRepo) Delete(ctx context.Context, id int) error {
	return r.client.Character.DeleteOneID(id).Exec(ctx)
}

func (r *entCharacterRepo) QueryAbilities(ctx context.Context, charID int) ([]*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Query().
		Where(characterability.HasCharacterWith(character.IDEQ(charID))).
		All(ctx)
}

func (r *entCharacterRepo) QueryAbilitiesWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Query().
		Where(characterability.HasCharacterWith(character.IDEQ(charID))).
		WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
		All(ctx)
}

func (r *entCharacterRepo) QueryTags(ctx context.Context, charID int) ([]*db.CharacterTag, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	return char.QueryTags().All(ctx)
}

func (r *entCharacterRepo) QueryFactions(ctx context.Context, charID int) ([]*db.CharacterFaction, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	return char.QueryFactionMemberships().All(ctx)
}

func (r *entCharacterRepo) QueryActiveEffects(ctx context.Context, charID int) ([]*db.ActiveEffect, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	return char.QueryActiveEffects().All(ctx)
}

func (r *entCharacterRepo) QueryQuestProgress(ctx context.Context, charID int) ([]*db.QuestProgress, error) {
	char, err := r.client.Character.Get(ctx, charID)
	if err != nil {
		return nil, err
	}
	return char.QueryQuestProgress().All(ctx)
}

func (r *entCharacterRepo) GetSkillLevels(ctx context.Context, charID int) (map[string]int, error) {
	charSkills, err := r.client.CharacterSkill.Query().
		Where(characterskill.HasCharacterWith(character.IDEQ(charID))).
		WithSkill().
		All(ctx)
	if err != nil {
		return nil, err
	}
	levels := make(map[string]int, len(charSkills))
	for _, cs := range charSkills {
		if cs.Edges.Skill != nil {
			levels[cs.Edges.Skill.Name] = cs.Level
		}
	}
	return levels, nil
}