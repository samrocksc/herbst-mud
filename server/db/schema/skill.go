package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

// Fields of the Skill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
		field.String("skill_type").
			Comment("e.g., combat, magic, utility"),
		field.Int("cost").
			Default(0).
			Comment("Points cost to learn/unlearn"),
		field.Int("cooldown").
			Default(0).
			Comment("Cooldown in seconds"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites"),
	}
}

// Edges of the Skill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", CharacterSkill.Type),
	}
}