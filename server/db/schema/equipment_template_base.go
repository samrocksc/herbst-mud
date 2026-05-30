package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

/** Base fields for EquipmentTemplate (non-combat). */
func equipmentTemplateBaseFields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Comment("Auto-increment primary key"),
		field.String("slug").
			Comment("URL-friendly unique identifier, e.g. 'pizza_dough' (unique per world_id)"),
		field.String("world_id").
			Default("1").
			Comment("World this item template belongs to (for multi-world support)"),
		field.String("name"),
		field.String("description"),
		field.String("slot"),
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.String("item_type").
			Default("misc").
			Comment("weapon|armor|consumable|quest|misc|container|potion"),
		field.JSON("stats", map[string]int{}).
			Optional().
			Comment("Stat bonuses e.g. {\"strength\": 5, \"dexterity\": 3}"),
		field.String("color").
			Default("").
			Comment("Custom display color"),
		field.Bool("is_visible").
			Default(true).
			Comment("Shown in room list"),
		field.Bool("is_immovable").
			Default(false).
			Comment("Cannot be picked up if true"),
		field.String("effect_type").
			Default("").
			Comment("heal|damage|dot|buff_armor|buff_dodge|buff_crit|debuff"),
		field.Int("effect_value").
			Default(0).
			Comment("Effect magnitude"),
		field.Int("effect_duration").
			Default(0).
			Comment("Duration in ticks (0 = instant)"),
		field.Bool("is_container").
			Default(false).
			Comment("Can hold items if true"),
		field.Int("container_capacity").
			Default(0).
			Comment("Max items container can hold"),
		field.Bool("is_locked").
			Default(false).
			Comment("Requires key to open"),
		field.String("key_item_id").
			Optional().
			Comment("ID of key item needed to unlock"),
		field.String("reveal_condition").
			Default("").
			Comment("JSON: {type: examine|perception_check|use_item|event, target, minLevel}"),
		field.Time("expires_at").
			Optional().
			Nillable().
			Comment("When this template's instances expire. nil = never rots."),
	}
}