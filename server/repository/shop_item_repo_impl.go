package repository

import (
	"context"
	"log/slog"

	"herbst-server/db"
	"herbst-server/db/shopitem"
	"herbst-server/dblog"
)

// shopItemRepo implements ShopItemRepo.
type shopItemRepo struct {
	client *db.Client
}

// NewShopItemRepo creates a new ShopItemRepo.
func NewShopItemRepo(client *db.Client) ShopItemRepo {
	return &shopItemRepo{client: client}
}

// Get returns a ShopItem by shop and equipment template IDs.
func (r *shopItemRepo) Get(ctx context.Context, shopID, equipmentTemplateID int) (*db.ShopItem, error) {
	item, err := r.client.ShopItem.Query().
		Where(shopitem.ShopIDEQ(shopID), shopitem.EquipmentTemplateIDEQ(equipmentTemplateID)).
		Only(ctx)
	if err != nil {
		dblog.Error("failed to get shop item", err, slog.Int("shop_id", shopID), slog.Int("equipment_template_id", equipmentTemplateID))
		return nil, err
	}
	return item, nil
}

// ListByShop returns all ShopItems for a shop.
func (r *shopItemRepo) ListByShop(ctx context.Context, shopID int) ([]*db.ShopItem, error) {
	items, err := r.client.ShopItem.Query().
		Where(shopitem.ShopIDEQ(shopID)).
		All(ctx)
	if err != nil {
		dblog.Error("failed to list shop items", err, slog.Int("shop_id", shopID))
		return nil, err
	}
	return items, nil
}

// Create creates a new ShopItem.
func (r *shopItemRepo) Create(ctx context.Context, input CreateShopItemInput) (*db.ShopItem, error) {
	item, err := r.client.ShopItem.Create().
		SetShopID(input.ShopID).
		SetEquipmentTemplateID(input.EquipmentTemplateID).
		SetCategory(input.Category).
		SetPrice(input.Price).
		SetQuantity(input.Quantity).
		SetMaxStock(input.MaxStock).
		SetIsEnabled(input.IsEnabled).
		Save(ctx)
	if err != nil {
		dblog.Error("failed to create shop item", err, slog.Int("shop_id", input.ShopID))
		return nil, err
	}
	return item, nil
}

// Update updates an existing ShopItem.
func (r *shopItemRepo) Update(ctx context.Context, shopID, equipmentTemplateID int, updates ShopItemUpdates) (*db.ShopItem, error) {
	// First get the existing item
	item, err := r.Get(ctx, shopID, equipmentTemplateID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if updates.Category != nil {
		item.Category = *updates.Category
	}
	if updates.Price != nil {
		item.Price = *updates.Price
	}
	if updates.Quantity != nil {
		item.Quantity = *updates.Quantity
	}
	if updates.MaxStock != nil {
		item.MaxStock = *updates.MaxStock
	}
	if updates.IsEnabled != nil {
		item.IsEnabled = *updates.IsEnabled
	}

	// Save the updated item
	updatedItem, err := r.client.ShopItem.UpdateOne(item).Save(ctx)
	if err != nil {
		dblog.Error("failed to update shop item", err, slog.Int("shop_id", shopID), slog.Int("equipment_template_id", equipmentTemplateID))
		return nil, err
	}
	return updatedItem, nil
}

// Delete deletes a ShopItem by shop and equipment template IDs.
func (r *shopItemRepo) Delete(ctx context.Context, shopID, equipmentTemplateID int) error {
	_, err := r.client.ShopItem.Delete().
		Where(shopitem.ShopIDEQ(shopID), shopitem.EquipmentTemplateIDEQ(equipmentTemplateID)).
		Exec(ctx)
	if err != nil {
		dblog.Error("failed to delete shop item", err, slog.Int("shop_id", shopID), slog.Int("equipment_template_id", equipmentTemplateID))
		return err
	}
	return nil
}
