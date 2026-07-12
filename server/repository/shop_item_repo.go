package repository

import (
	"context"

	"herbst-server/db"
)

// ShopItemRepo defines the repository interface for ShopItem.
type ShopItemRepo interface {
	Get(ctx context.Context, shopID, equipmentTemplateID int) (*db.ShopItem, error)
	ListByShop(ctx context.Context, shopID int) ([]*db.ShopItem, error)
	Create(ctx context.Context, input CreateShopItemInput) (*db.ShopItem, error)
	Update(ctx context.Context, shopID, equipmentTemplateID int, updates ShopItemUpdates) (*db.ShopItem, error)
	Delete(ctx context.Context, shopID, equipmentTemplateID int) error
}

// CreateShopItemInput is the input for creating a ShopItem.
type CreateShopItemInput struct {
	ShopID               int    `json:"shop_id"`
	EquipmentTemplateID  int    `json:"equipment_template_id"`
	Category             string `json:"category"`
	Price                int    `json:"price"`
	Quantity             int    `json:"quantity"`
	MaxStock             int    `json:"max_stock"`
	IsEnabled            bool   `json:"is_enabled"`
}

// ShopItemUpdates is the input for updating a ShopItem.
type ShopItemUpdates struct {
	Category    *string `json:"category,omitempty"`
	Price       *int    `json:"price,omitempty"`
	Quantity    *int    `json:"quantity,omitempty"`
	MaxStock    *int    `json:"max_stock,omitempty"`
	IsEnabled   *bool   `json:"is_enabled,omitempty"`
}
