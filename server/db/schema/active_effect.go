package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// ActiveEffect tracks a runtime effect applied to a character.
// Created when an Effect is applied; stores stack state and expiry.
type ActiveEffect struct {
	ent.Schema
}

func (ActiveEffect) Fields() []ent.Field {
	return []ent.Field{
		field.Int("character_id").
			Comment("FK to Character who has this effect"),
		field.Int("effect_id").
			Comment("FK to Effect definition"),
		field.Int("applied_by_id").
			Default(0).
			Comment("Character ID of who applied this effect"),
		field.Int("stack_count").
			Default(1),
		field.Time("started_at").
			Default(time.Now),
		field.Time("expires_at").
			Optional().
			Nillable().
			Comment("nil if permanent or instant"),
		field.Bool("is_active").
			Default(true),
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