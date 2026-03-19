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
		field.String("examineDesc").
			Default("").
			Comment("Detailed description shown with examine command"),
		field.JSON("hiddenDetails", []map[string]any{}).
			Default([]map[string]any{}).
			Comment("Details revealed based on examine skill"),
		field.Int("hiddenThreshold").
			Default(0).
			Comment("Examine skill required to reveal hidden details"),
		// Hidden items and reveal conditions (GitHub #12 - Look System)
		field.String("revealCondition").
			Default("").
			Comment("JSON: {type: examine|perception_check|use_item|event, target, minLevel}"),
		// Weapon fields (GitHub #92)
		field.Int("minDamage").
			Default(0).
			Comment("Minimum damage for weapons"),
		field.Int("maxDamage").
			Default(0).
			Comment("Maximum damage for weapons"),
		field.String("weaponType").
			Default("").
			Comment("Type of weapon: sword, dagger, staff, etc."),
		field.String("classRestriction").
			Default("").
			Comment("Class that can use this weapon"),
		field.Bool("guaranteedDrop").
			Default(false).
			Comment("Always drops from certain NPCs"),
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