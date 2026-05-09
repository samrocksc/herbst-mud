package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Effect holds the schema definition for the Effect entity (herbst client).
// Mirrors server schema with game-logic fields only.
type Effect struct {
	ent.Schema
}

func (Effect) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("effect_type"),
		field.JSON("parameters", map[string]interface{}{}).
			Default(map[string]interface{}{}),
		field.String("stack_mode").Default("replace"),
		field.Int("stack_limit").Default(1),
		field.Bool("is_permanent").Default(false),
		field.Int("duration_secs").Default(0),
	}
}

func (Effect) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hooks", EffectHook.Type),
		edge.To("active_effect_instances", ActiveEffect.Type),
	}
}