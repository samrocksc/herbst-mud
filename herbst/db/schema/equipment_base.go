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
		field.String("revealCondition").
			Default("").
			Comment("JSON: {type: examine|perception_check|use_item|event, target, minLevel}"),
		field.String("equipment_template_id").
			Optional().
			Comment("FK to equipment_template"),
	}
}