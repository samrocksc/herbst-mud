package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

/** Base fields for Equipment (non-combat). */
func equipmentBaseFields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.String("slot"),
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.Bool("isEquipped").
			Default(false),
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
		field.String("equipment_template_id").
			Optional().
			Comment("FK to equipment_template"),
		field.Int("ownerId").
			Optional().
			Nillable().
			Comment("Character ID that owns this item"),
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
		field.String("revealCondition").
			Default("").
			Comment("JSON: {type: examine|perception_check|use_item|event, target, minLevel}"),
		field.String("examineDesc").
			Default("").
			Comment("Detailed description shown with examine command"),
		field.JSON("hiddenDetails", []map[string]any{}).
			Default([]map[string]any{}).
			Comment("Details revealed based on examine skill"),
		field.Int("hiddenThreshold").
			Default(0).
			Comment("Examine skill required to reveal hidden details"),
		field.Time("expiresAt").
			Optional().
			Nillable().
			Comment("When this item expires and is auto-deleted. nil = never rots."),
		field.Int("quantity").
			Default(1).
			Comment("Stack size for consumable/ingredient items"),
	}
}