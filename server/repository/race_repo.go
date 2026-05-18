package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/race"
)

type entRaceRepo struct {
	client *db.Client
}

func NewEntRaceRepo(client *db.Client) RaceRepo {
	return &entRaceRepo{client: client}
}

func (r *entRaceRepo) Get(ctx context.Context, id int) (*db.Race, error) {
	return r.client.Race.Get(ctx, id)
}

func (r *entRaceRepo) GetByName(ctx context.Context, name string) (*db.Race, error) {
	return r.client.Race.Query().
		Where(race.NameEQ(name)).
		WithTags().
		Only(ctx)
}

func (r *entRaceRepo) List(ctx context.Context) ([]*db.Race, error) {
	return r.client.Race.Query().All(ctx)
}

func (r *entRaceRepo) Create(ctx context.Context, input CreateRaceInput) (*db.Race, error) {
	builder := r.client.Race.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName).
		SetDescription(input.Description).
		SetIsPlayable(input.IsPlayable)
	if input.StatModifiers != nil {
		builder = builder.SetStatModifiers(*input.StatModifiers)
	}
	if input.Color != "" {
		builder = builder.SetColor(input.Color)
	}
	if len(input.EquipmentSlots) > 0 {
		builder = builder.SetEquipmentSlots(input.EquipmentSlots)
	}
	if len(input.TagIDs) > 0 {
		builder = builder.AddTagIDs(input.TagIDs...)
	}
	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.client.Race.Query().
		Where(race.IDEQ(created.ID)).
		WithTags().
		Only(ctx)
}

func (r *entRaceRepo) CountCharactersByRaceName(ctx context.Context, raceName string) (int, error) {
	return r.client.Character.Query().
		Where(character.RaceEQ(raceName)).
		Count(ctx)
}

func (r *entRaceRepo) GetWithTags(ctx context.Context, id int) (*db.Race, error) {
	return r.client.Race.Query().
		Where(race.ID(id)).
		WithTags().
		Only(ctx)
}

func (r *entRaceRepo) ListWithTags(ctx context.Context) ([]*db.Race, error) {
	return r.client.Race.Query().WithTags().All(ctx)
}

func (r *entRaceRepo) Update(ctx context.Context, id int, updates RaceUpdates) (*db.Race, error) {
	builder := r.client.Race.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.DisplayName != nil {
		builder = builder.SetDisplayName(*updates.DisplayName)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.StatModifiers != nil {
		builder = builder.SetStatModifiers(*updates.StatModifiers)
	}
	if updates.IsPlayable != nil {
		builder = builder.SetIsPlayable(*updates.IsPlayable)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.EquipmentSlots != nil {
		builder = builder.SetEquipmentSlots(updates.EquipmentSlots)
	}
	if updates.ClearTags {
		builder = builder.ClearTags()
	}
	if len(updates.AddTagIDs) > 0 {
		builder = builder.AddTagIDs(updates.AddTagIDs...)
	}
	_, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.client.Race.Query().
		Where(race.IDEQ(id)).
		WithTags().
		Only(ctx)
}

func (r *entRaceRepo) Delete(ctx context.Context, id int) error {
	return r.client.Race.DeleteOneID(id).Exec(ctx)
}

// Race represents a playable race for character creation
type Race struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func (r *entRaceRepo) ListPlayable(ctx context.Context) ([]*Race, error) {
	races, err := r.client.Race.Query().
		Where(race.IsPlayable(true)).
		Order(race.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*Race, len(races))
	for i, r := range races {
		result[i] = &Race{
			ID:          r.ID,
			Name:        r.Name,
			DisplayName: r.DisplayName,
		}
	}
	return result, nil
}