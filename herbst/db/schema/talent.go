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
		field.String("name"),
		field.String("description"),
		field.JSON("requirements", map[string]int{}), // e.g., {"level": 5, "strength": 10}
	}
}

// Edges of the Talent.
func (Talent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
	}
}