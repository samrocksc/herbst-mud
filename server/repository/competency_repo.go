package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/competencycategory"
)

type entCompetencyRepo struct {
	client *db.Client
}

func NewEntCompetencyRepo(client *db.Client) CompetencyRepo {
	return &entCompetencyRepo{client: client}
}

func (r *entCompetencyRepo) GetCategory(ctx context.Context, id string) (*db.CompetencyCategory, error) {
	return r.client.CompetencyCategory.Get(ctx, id)
}

func (r *entCompetencyRepo) ListCategories(ctx context.Context) ([]*db.CompetencyCategory, error) {
	return r.client.CompetencyCategory.Query().All(ctx)
}

func (r *entCompetencyRepo) CreateCategory(ctx context.Context, input CreateCompetencyInput) (*db.CompetencyCategory, error) {
	builder := r.client.CompetencyCategory.Create().
		SetID(input.ID).
		SetName(input.Name)
	if input.XPMultiplier != 0 {
		builder = builder.SetXpMultiplier(input.XPMultiplier)
	}
	return builder.Save(ctx)
}

func (r *entCompetencyRepo) GetCharacterCompetency(ctx context.Context, charID int, categoryID string) (*db.CharacterCompetency, error) {
	return r.client.CharacterCompetency.Query().
		Where(
			charactercompetency.HasCategoryWith(competencycategory.ID(categoryID)),
		).
		Only(ctx)
}

func (r *entCompetencyRepo) UpsertCharacterCompetency(ctx context.Context, charID int, categoryID string, xp, level int) (*db.CharacterCompetency, error) {
	existing, err := r.GetCharacterCompetency(ctx, charID, categoryID)
	if err != nil {
		return r.client.CharacterCompetency.Create().
			SetCharacterID(charID).
			SetCategoryID(categoryID).
			SetXp(xp).
			SetLevel(level).
			Save(ctx)
	}
	return r.client.CharacterCompetency.UpdateOneID(existing.ID).
		SetXp(xp).
		SetLevel(level).
		Save(ctx)
}

func (r *entCompetencyRepo) CountCompetenciesByCategory(ctx context.Context, categoryID string) (int, error) {
	return r.client.CharacterCompetency.Query().
		Where(charactercompetency.HasCategoryWith(competencycategory.ID(categoryID))).
		Count(ctx)
}

func (r *entCompetencyRepo) GetCategoryWithThresholds(ctx context.Context, id string) (*db.CompetencyCategory, error) {
	return r.client.CompetencyCategory.Query().
		Where(competencycategory.ID(id)).
		WithThresholds().
		Only(ctx)
}

func (r *entCompetencyRepo) UpdateCategory(ctx context.Context, id string, updates CompetencyCategoryUpdates) (*db.CompetencyCategory, error) {
	builder := r.client.CompetencyCategory.UpdateOneID(id)
	if updates.Name != nil {
		builder = builder.SetName(*updates.Name)
	}
	if updates.XPMultiplier != nil {
		builder = builder.SetXpMultiplier(*updates.XPMultiplier)
	}
	return builder.Save(ctx)
}

func (r *entCompetencyRepo) DeleteCategory(ctx context.Context, id string) error {
	return r.client.CompetencyCategory.DeleteOneID(id).Exec(ctx)
}