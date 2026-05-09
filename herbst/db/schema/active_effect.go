package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// ActiveEffect tracks a runtime effect applied to a character.
type ActiveEffect struct {
	ent.Schema
}

func (ActiveEffect) Fields() []ent.Field {
	return []ent.Field{
		field.Int("character_id"),
		field.Int("effect_id"),
		field.Int("applied_by_id").Default(0),
		field.Int("stack_count").Default(1),
		field.Time("started_at").Default(time.Now),
		field.Time("expires_at").Optional().Nillable(),
		field.Bool("is_active").Default(true),
	}
}

func (ActiveEffect) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("active_effects").
			Unique().
			Field("character_id").
			Required(),
		edge.From("effect", Effect.Type).
			Ref("active_effect_instances").
			Unique().
			Field("effect_id").
			Required(),
	}
}