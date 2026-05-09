package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EffectHook binds an event to an Effect on a character/NPC template.
// Named EffectHook to avoid conflict with ent's predeclared Hook identifier.
type EffectHook struct {
	ent.Schema
}

func (EffectHook) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("event"),
		field.String("target").Default("self"),
		field.String("condition").Optional(),
		field.Bool("enabled").Default(true),
	}
}

func (EffectHook) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("effect", Effect.Type).
			Ref("hooks").
			Unique().
			Required(),
		edge.From("npc_template", NPCTemplate.Type).
			Ref("hooks").
			Unique(),
	}
}