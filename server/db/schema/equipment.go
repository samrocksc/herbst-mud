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
			Comment("weapon|armor|consumable|quest|misc|container|potion"),
		// Owner system - items can be owned by a character
		field.Int("ownerId").
			Optional().
			Nillable().
			Comment("Character ID that owns this item, nil if in a room"),
		// Consumable effects (unified effect system)
		field.String("effect_type").
			Default("").
			Comment("heal|damage|dot|buff_armor|buff_dodge|buff_crit|debuff"),
		field.Int("effect_value").
			Default(0).
			Comment("Effect magnitude"),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.Int("healing").
			Default(0).
			Comment("DEPRECATED: Use effect_type=heal and effect_value instead"),
		field.String("effect").
			Default("").
			Comment("DEPRECATED: Use effect_type instead"),
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
		// Corpse rotting (GitHub #22)
		field.Time("expiresAt").
			Optional().
			Nillable().
			Comment("When this item expires and is auto-deleted. nil = never rots."),
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