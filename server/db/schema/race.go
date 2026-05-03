package schema

import (
	"entgo.io/ent"
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
		field.Bool("is_playable").
			Default(true).
			Comment("false = NPC-only race"),
		field.String("color").
			Optional().
			Comment("Hex color for UI display, e.g. '#8b5cf6'"),
	}
}

// Edges of the Race.
func (Race) Edges() []ent.Edge {
	return nil
}