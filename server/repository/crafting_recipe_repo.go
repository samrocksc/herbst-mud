package repository

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/craftingrecipe"
	"herbst-server/db/schema"
)

type CraftingRecipeRepo interface {
	Get(ctx context.Context, name string) (*db.CraftingRecipe, error)
	List(ctx context.Context, worldID, stationTag string) ([]*db.CraftingRecipe, error)
	Create(ctx context.Context, input CreateCraftingRecipeInput) (*db.CraftingRecipe, error)
	Update(ctx context.Context, name string, updates CraftingRecipeUpdates) (*db.CraftingRecipe, error)
	Delete(ctx context.Context, name string) error
}

type entCraftingRecipeRepo struct {
	client *db.Client
}

func NewEntCraftingRecipeRepo(client *db.Client) CraftingRecipeRepo {
	return &entCraftingRecipeRepo{client: client}
}

func (r *entCraftingRecipeRepo) Get(ctx context.Context, name string) (*db.CraftingRecipe, error) {
	return r.client.CraftingRecipe.Query().Where(craftingrecipe.Name(name)).Only(ctx)
}

func (r *entCraftingRecipeRepo) List(ctx context.Context, worldID, stationTag string) ([]*db.CraftingRecipe, error) {
	query := r.client.CraftingRecipe.Query()
	if worldID != "" {
		query = query.Where(craftingrecipe.WorldID(worldID))
	}
	if stationTag != "" {
		query = query.Where(craftingrecipe.RequiredStationTag(stationTag))
	}
	return query.All(ctx)
}

type CreateCraftingRecipeInput struct {
	Name                string
	DisplayName         string
	Description         string
	RequiredStationTag  string
	RequiredClass       string
	RequiredSkillLevel  int
	RequiredSkill       string
	Inputs              []schema.CraftingInput
	Outputs             []schema.CraftingOutput
	CraftTimeSecs       int
	WorldID             string
}

func (r *entCraftingRecipeRepo) Create(ctx context.Context, input CreateCraftingRecipeInput) (*db.CraftingRecipe, error) {
	builder := r.client.CraftingRecipe.Create().
		SetName(input.Name).
		SetDisplayName(input.DisplayName).
		SetRequiredStationTag(input.RequiredStationTag).
		SetInputs(input.Inputs).
		SetOutputs(input.Outputs).
		SetCraftTimeSecs(input.CraftTimeSecs).
		SetWorldID(input.WorldID)
	if input.Description != "" {
		builder = builder.SetDescription(input.Description)
	}
	if input.RequiredClass != "" {
		builder = builder.SetRequiredClass(input.RequiredClass)
	}
	if input.RequiredSkillLevel != 0 {
		builder = builder.SetRequiredSkillLevel(input.RequiredSkillLevel)
	}
	if input.RequiredSkill != "" {
		builder = builder.SetRequiredSkill(input.RequiredSkill)
	}
	return builder.Save(ctx)
}

type CraftingRecipeUpdates struct {
	DisplayName        *string
	Description        *string
	RequiredStationTag *string
	RequiredClass      *string
	RequiredSkillLevel *int
	RequiredSkill      *string
	Inputs             *[]schema.CraftingInput
	Outputs            *[]schema.CraftingOutput
	CraftTimeSecs      *int
	WorldID            *string
}

func (r *entCraftingRecipeRepo) Update(ctx context.Context, name string, updates CraftingRecipeUpdates) (*db.CraftingRecipe, error) {
	existing, err := r.client.CraftingRecipe.Query().Where(craftingrecipe.Name(name)).Only(ctx)
	if err != nil {
		return nil, err
	}
	builder := existing.Update()
	if updates.DisplayName != nil {
		builder = builder.SetDisplayName(*updates.DisplayName)
	}
	if updates.Description != nil {
		builder = builder.SetDescription(*updates.Description)
	}
	if updates.RequiredStationTag != nil {
		builder = builder.SetRequiredStationTag(*updates.RequiredStationTag)
	}
	if updates.RequiredClass != nil {
		builder = builder.SetRequiredClass(*updates.RequiredClass)
	}
	if updates.RequiredSkillLevel != nil {
		builder = builder.SetRequiredSkillLevel(*updates.RequiredSkillLevel)
	}
	if updates.RequiredSkill != nil {
		builder = builder.SetRequiredSkill(*updates.RequiredSkill)
	}
	if updates.Inputs != nil {
		builder = builder.SetInputs(*updates.Inputs)
	}
	if updates.Outputs != nil {
		builder = builder.SetOutputs(*updates.Outputs)
	}
	if updates.CraftTimeSecs != nil {
		builder = builder.SetCraftTimeSecs(*updates.CraftTimeSecs)
	}
	if updates.WorldID != nil {
		builder = builder.SetWorldID(*updates.WorldID)
	}
	return builder.Save(ctx)
}

func (r *entCraftingRecipeRepo) Delete(ctx context.Context, name string) error {
	existing, err := r.client.CraftingRecipe.Query().Where(craftingrecipe.Name(name)).Only(ctx)
	if err != nil {
		return err
	}
	return r.client.CraftingRecipe.DeleteOne(existing).Exec(ctx)
}