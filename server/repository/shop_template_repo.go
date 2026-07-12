package repository

import (
	"context"

	"herbst-server/db"
)

// ShopTemplateRepo defines the repository interface for ShopTemplate.
type ShopTemplateRepo interface {
	Get(ctx context.Context, id int) (*db.ShopTemplate, error)
	List(ctx context.Context, worldID string) ([]*db.ShopTemplate, error)
	Create(ctx context.Context, input CreateShopTemplateInput) (*db.ShopTemplate, error)
	Update(ctx context.Context, id int, updates ShopTemplateUpdates) (*db.ShopTemplate, error)
	Delete(ctx context.Context, id int) error
}

// CreateShopTemplateInput is the input for creating a ShopTemplate.
type CreateShopTemplateInput struct {
	Name             string `json:"name"`
	WorldID          string `json:"world_id"`
	NPCTemplateID    string `json:"npc_template_id,omitempty"`
	CurrencyItemType int    `json:"currency_item_type,omitempty"`
	MaxInventory     int    `json:"max_inventory"`
	GoldReserves     int    `json:"gold_reserves"`
	IsActive         bool   `json:"is_active"`
}

// ShopTemplateUpdates is the input for updating a ShopTemplate.
type ShopTemplateUpdates struct {
	Name             *string `json:"name,omitempty"`
	WorldID          *string `json:"world_id,omitempty"`
	NPCTemplateID    *string `json:"npc_template_id,omitempty"`
	CurrencyItemType *int    `json:"currency_item_type,omitempty"`
	MaxInventory     *int    `json:"max_inventory,omitempty"`
	GoldReserves     *int    `json:"gold_reserves,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
}
