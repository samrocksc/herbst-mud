package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ShopItem holds the schema definition for the ShopItem entity.
// A ShopItem represents an item available for purchase in a shop.
type ShopItem struct {
	ent.Schema
}

// Fields of the ShopItem.
func (ShopItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("shop_id").
			Comment("FK to shop_template"),
		field.Int("equipment_template_id").
			Comment("FK to equipment_template for the item"),
		field.String("category").
			Optional().
			Comment("Item category for filtering (weapons, armor, potions, etc.)"),
		field.Int("price").
			Default(0).
			Comment("Base price in currency"),
		field.Int("quantity").
			Default(0).
			Comment("Current stock in shop inventory"),
		field.Int("max_stock").
			Default(99).
			Comment("Maximum stack size for this item"),
		field.Bool("is_enabled").
			Default(true).
			Comment("If false, item cannot be purchased"),
		field.Time("last_restocked").
			Optional().
			Nillable().
			Comment("When this item was last restocked"),
	}
}

// Indexes of the ShopItem.
func (ShopItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("shop_id", "equipment_template_id").Unique(),
		index.Fields("shop_id", "category"),
		index.Fields("category"),
	}
}
