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

func (r *entRaceRepo) GetByName(ctx context.Context, name, worldID string) (*db.Race, error) {
	return r.client.Race.Query().
		Where(race.Name(name), race.WorldID(worldID)).
		WithTags().
		Only(ctx)
}

func (r *entRaceRepo) List(ctx context.Context, worldID string) ([]*db.Race, error) {
	return r.client.Race.Query().Where(race.WorldID(worldID)).All(ctx)
}

func (r *entRaceRepo) Create(ctx context.Context, input CreateRaceInput) (*db.Race, error) {
	builder := r.client.Race.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName).
		SetDescription(input.Description).
		SetRequirementTags(input.RequirementTags)
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

func (r *entRaceRepo) CountCharactersByRaceName(ctx context.Context, raceName, worldID string) (int, error) {
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

func (r *entRaceRepo) ListWithTags(ctx context.Context, worldID string) ([]*db.Race, error) {
	return r.client.Race.Query().Where(race.WorldID(worldID)).WithTags().All(ctx)
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
	if updates.RequirementTags != nil {
		builder = builder.SetRequirementTags(updates.RequirementTags)
	}
	if updates.Color != nil {
		builder = builder.SetColor(*updates.Color)
	}
	if updates.EquipmentSlots != nil {
		builder = builder.SetEquipmentSlots(updates.EquipmentSlots)
	}
	if updates.Resistances != nil {
		builder = builder.SetResistances(updates.Resistances)
	}
	if updates.Vulnerabilities != nil {
		builder = builder.SetVulnerabilities(updates.Vulnerabilities)
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

func (r *entRaceRepo) CountByWorld(ctx context.Context, worldID string) (int, error) {
	return r.client.Race.Query().Where(race.WorldIDEQ(worldID)).Count(ctx)
}

// Race represents a playable race for character creation
type Race struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	WorldID     string `json:"world_id"`
}

func (r *entRaceRepo) ListPlayable(ctx context.Context, worldID string) ([]*Race, error) {
	all, err := r.client.Race.Query().
		Where(race.WorldID(worldID)).
		Order(race.ByName()).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*Race, 0, len(all))
	for _, item := range all {
		if len(item.RequirementTags) == 0 {
			result = append(result, &Race{
				ID:          item.ID,
				Name:        item.Name,
				DisplayName: item.DisplayName,
				WorldID:     item.WorldID,
			})
		}
	}
	return result, nil
}
