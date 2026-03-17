package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Talent holds the schema definition for the Talent entity.
type Talent struct {
	ent.Schema
}

// Fields of the Talent.
func (Talent) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
		field.String("requirements").
			Optional().
			Comment("JSON string of prerequisites (skills, levels, etc.)"),
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
	}
}