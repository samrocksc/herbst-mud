package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ShopTemplate holds the schema definition for the ShopTemplate entity.
// A ShopTemplate defines a shop configuration that can be attached to an NPC instance.
type ShopTemplate struct {
	ent.Schema
}

// Fields of the ShopTemplate.
func (ShopTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Shop name, e.g., 'Gems & Glitter'"),
		field.String("world_id").
			Default("1").
			Comment("World this shop belongs to (for multi-world support)"),
		field.String("npc_template_id").
			Optional().
			Comment("FK to npc_template for vendor NPCs"),
		field.Int("currency_item_type").
			Optional().
			Comment("FK to equipment_template for currency (gold, silver, etc.)"),
		field.Int("max_inventory").
			Default(50).
			Comment("Maximum number of items this shop can hold"),
		field.Int("gold_reserves").
			Default(1000).
			Comment("Starting gold reserves for the shop"),
		field.Bool("is_active").
			Default(true).
			Comment("If false, shop is closed but retains inventory"),
		field.Time("last_restocked").
			Optional().
			Nillable().
			Comment("When the shop was last restocked (nil = never)"),
	}
}

// Indexes of the ShopTemplate.
func (ShopTemplate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("world_id", "name").Unique(),
	}
}
