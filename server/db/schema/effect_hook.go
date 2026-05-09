package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EffectHook binds an event to an Effect on a character/NPC template.
// When the event fires, the linked Effect is applied to the target.
// Named EffectHook to avoid conflict with ent's predeclared Hook identifier.
type EffectHook struct {
	ent.Schema
}

func (EffectHook) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Display name, e.g., 'Death Drain — XP from killer'"),
		field.String("event").
			Comment("on_death|on_hit_received|on_hit_dealt|on_kill|on_enter_room|on_leave_room|on_equip|on_unequip|on_login|on_effect_start|on_effect_end"),
		field.String("target").
			Default("self").
			Comment("self|attacker|killer|room|owner — who the effect targets"),
		field.String("condition").
			Optional().
			Comment("Optional SPICE condition expression (deferred)"),
		field.Bool("enabled").
			Default(true),
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