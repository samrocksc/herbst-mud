package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
)

type entCharacterAbilityRepo struct {
	client *db.Client
}

func NewEntCharacterAbilityRepo(client *db.Client) CharacterAbilityRepo {
	return &entCharacterAbilityRepo{client: client}
}

func (r *entCharacterAbilityRepo) ListByCharacter(ctx context.Context, charID int) ([]*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Query().
		Where(characterability.HasCharacterWith(character.ID(charID))).
		All(ctx)
}

func (r *entCharacterAbilityRepo) ListByCharacterWithDetails(ctx context.Context, charID int) ([]*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Query().
		Where(characterability.HasCharacterWith(character.ID(charID))).
		WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
		All(ctx)
}

func (r *entCharacterAbilityRepo) ListByCharacterAndSlot(ctx context.Context, charID int, slot int) ([]*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Query().
		Where(
			characterability.HasCharacterWith(character.ID(charID)),
			characterability.SlotEQ(slot),
		).
		WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
		All(ctx)
}

func (r *entCharacterAbilityRepo) CountByCharacter(ctx context.Context, charID int) (int, error) {
	return r.client.CharacterAbility.Query().
		Where(characterability.HasCharacterWith(character.ID(charID))).
		Count(ctx)
}

func (r *entCharacterAbilityRepo) ExistsByCharacterAndAbility(ctx context.Context, charID int, abilityID int) (bool, error) {
	return r.client.CharacterAbility.Query().
		Where(
			characterability.HasCharacterWith(character.ID(charID)),
			characterability.HasAbilityWith(ability.ID(abilityID)),
		).
		Exist(ctx)
}

func (r *entCharacterAbilityRepo) Create(ctx context.Context, charID int, abilityID int, slot int) (*db.CharacterAbility, error) {
	return r.client.CharacterAbility.Create().
		SetCharacterID(charID).
		SetAbilityID(abilityID).
		SetSlot(slot).
		Save(ctx)
}

func (r *entCharacterAbilityRepo) Delete(ctx context.Context, id int) error {
	return r.client.CharacterAbility.DeleteOneID(id).Exec(ctx)
}

func (r *entCharacterAbilityRepo) DeleteByCharacterAndSlot(ctx context.Context, charID int, slot int) error {
	_, err := r.client.CharacterAbility.Delete().
		Where(
			characterability.HasCharacterWith(character.ID(charID)),
			characterability.Slot(slot),
		).Exec(ctx)
	return err
}

func (r *entCharacterAbilityRepo) DeleteByCharacterAndAbility(ctx context.Context, charID int, abilityID int) error {
	_, err := r.client.CharacterAbility.Delete().
		Where(
			characterability.HasCharacterWith(character.ID(charID)),
			characterability.HasAbilityWith(ability.ID(abilityID)),
		).Exec(ctx)
	return err
}