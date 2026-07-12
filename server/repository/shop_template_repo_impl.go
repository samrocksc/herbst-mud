package repository

import (
	"context"
	"log/slog"

	"herbst-server/db"
	"herbst-server/db/shoptemplate"
	"herbst-server/dblog"
)

// shopTemplateRepo implements ShopTemplateRepo.
type shopTemplateRepo struct {
	client *db.Client
}

// NewShopTemplateRepo creates a new ShopTemplateRepo.
func NewShopTemplateRepo(client *db.Client) ShopTemplateRepo {
	return &shopTemplateRepo{client: client}
}

// Get returns a ShopTemplate by ID.
func (r *shopTemplateRepo) Get(ctx context.Context, id int) (*db.ShopTemplate, error) {
	template, err := r.client.ShopTemplate.Query().
		Where(shoptemplate.ID(id)).
		Only(ctx)
	if err != nil {
		dblog.Error("failed to get shop template", err, slog.Int("shop_id", id))
		return nil, err
	}
	return template, nil
}

// List returns all ShopTemplates for a world.
func (r *shopTemplateRepo) List(ctx context.Context, worldID string) ([]*db.ShopTemplate, error) {
	templates, err := r.client.ShopTemplate.Query().
		Where(shoptemplate.WorldIDEQ(worldID)).
		All(ctx)
	if err != nil {
		dblog.Error("failed to list shop templates", err, slog.String("world_id", worldID))
		return nil, err
	}
	return templates, nil
}

// Create creates a new ShopTemplate.
func (r *shopTemplateRepo) Create(ctx context.Context, input CreateShopTemplateInput) (*db.ShopTemplate, error) {
	builder := r.client.ShopTemplate.Create().
		SetName(input.Name).
		SetWorldID(input.WorldID).
		SetMaxInventory(input.MaxInventory).
		SetGoldReserves(input.GoldReserves).
		SetIsActive(input.IsActive)

	// Set optional fields if provided
	if input.NPCTemplateID != "" {
		builder = builder.SetNpcTemplateID(input.NPCTemplateID)
	}
	if input.CurrencyItemType != 0 {
		builder = builder.SetCurrencyItemType(input.CurrencyItemType)
	}

	template, err := builder.Save(ctx)
	if err != nil {
		dblog.Error("failed to create shop template", err, slog.String("name", input.Name))
		return nil, err
	}
	return template, nil
}

// Update updates an existing ShopTemplate.
func (r *shopTemplateRepo) Update(ctx context.Context, id int, updates ShopTemplateUpdates) (*db.ShopTemplate, error) {
	query := r.client.ShopTemplate.UpdateOneID(id)

	if updates.Name != nil {
		query.SetName(*updates.Name)
	}
	if updates.WorldID != nil {
		query.SetWorldID(*updates.WorldID)
	}
	if updates.NPCTemplateID != nil {
		query.SetNpcTemplateID(*updates.NPCTemplateID)
	}
	if updates.CurrencyItemType != nil {
		query.SetCurrencyItemType(*updates.CurrencyItemType)
	}
	if updates.MaxInventory != nil {
		query.SetMaxInventory(*updates.MaxInventory)
	}
	if updates.GoldReserves != nil {
		query.SetGoldReserves(*updates.GoldReserves)
	}
	if updates.IsActive != nil {
		query.SetIsActive(*updates.IsActive)
	}

	template, err := query.Save(ctx)
	if err != nil {
		dblog.Error("failed to update shop template", err, slog.Int("shop_id", id))
		return nil, err
	}
	return template, nil
}

// Delete deletes a ShopTemplate by ID.
func (r *shopTemplateRepo) Delete(ctx context.Context, id int) error {
	err := r.client.ShopTemplate.DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		dblog.Error("failed to delete shop template", err, slog.Int("shop_id", id))
		return err
	}
	return nil
}
