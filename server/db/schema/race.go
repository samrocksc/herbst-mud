package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Race holds the schema definition for the Race entity.
type Race struct {
	ent.Schema
}

// Fields of the Race.
func (Race) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Comment("Internal ID: human, turtle, mutant"),
		field.String("display_name").
			Comment("Shown in UI: Human, Turtle, Mutant"),
		field.Text("description").
			Comment("Flavor text for character creation"),
		field.String("stat_modifiers").
			Optional().
			Comment(`JSON: {"strength": 2, "dexterity": 0, ...}`),
		field.String("skill_grants").
			Optional().
			Comment(`JSON array: ["swim", "shell_defense"]`),
		field.JSON("equipment_slots", []string{}).
			Default([]string{"head", "neck", "chest", "back", "hands", "legs", "feet", "finger_left", "finger_right", "main_hand", "off_hand"}).
			Comment(`Slots this race can equip: ["head","chest",...]`),
		field.Strings("requirement_tags").
			Default([]string{}).
			StorageKey("requirement_tags").
			Optional().
			Comment("Tags that must be satisfied for race to be selectable"),
		field.String("color").
			Optional().
			Comment("Hex color for UI display, e.g. '#8b5cf6'"),
	}
}

// Edges of the Race.
func (Race) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", Tag.Type).Ref("races"),
		edge.To("npc_templates", NPCTemplate.Type),
	}
}