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
		// Item system fields (GitHub #89)
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
			Comment("weapon|armor|consumable|quest|misc|container"),
		// Container system fields (GitHub #143)
		field.Bool("isContainer").
			Default(false).
			Comment("Can hold items if true"),
		field.Int("containerCapacity").
			Default(0).
			Comment("Max items container can hold"),
		field.Bool("isLocked").
			Default(false).
			Comment("Requires key to open"),
		field.String("keyItemID").
			Optional().
			Comment("ID of key item needed to unlock"),
		field.String("containedItems").
			Default("").
			Comment("JSON array of contained item IDs"),
		// Hidden items and reveal conditions (GitHub #12 - Look System)
		field.String("revealCondition").
			Default("").
			Comment("JSON: {type: examine|perception_check|use_item|event, target, minLevel}"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
	}
}