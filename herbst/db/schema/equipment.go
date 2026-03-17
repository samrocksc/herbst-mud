package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Equipment holds the schema definition for the Equipment entity.
type Equipment struct {
	ent.Schema
}

// Fields of the Equipment.
func (Equipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.String("slot"), // e.g., "head", "chest", "weapon", "legs"
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.Bool("isEquipped").
			Default(false),
		// New fields for item system (GitHub #89)
		field.Bool("isImmovable").
			Default(false).
			Comment("Cannot be picked up if true"),
		field.String("color").
			Default("").
			Comment("Custom display color (e.g., gold for immovable items)"),
		field.Bool("isVisible").
			Default(true).
			Comment("Shown in room list"),
		field.String("itemType").
			Default("misc").
			Comment("weapon|armor|consumable|quest|misc"),
		// Weapon-specific fields (GitHub #92)
		field.Int("minDamage").
			Default(1).
			Comment("Minimum damage for weapons"),
		field.Int("maxDamage").
			Default(1).
			Comment("Maximum damage for weapons"),
		field.String("weaponType").
			Default("").
			Comment("sword|dagger|staff|pipe|brawling - weapon style"),
		field.String("classRestriction").
			Default("").
			Comment("Class that can use this weapon (e.g., warrior, chef)"),
		field.Bool("isDroppable").
			Default(true).
			Comment("Can be dropped by NPCs"),
		field.Bool("guaranteedDrop").
			Default(false).
			Comment("Always drops on first NPC kill"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
		// Character inventory (GitHub #92)
		edge.To("character", Character.Type).
			Unique().
			Comment("Character carrying this item"),
	}
}