package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Effect holds the schema definition for the Effect entity.
// Effects are reusable state-change definitions: xp_drain, hp_change,
// teleport, tag_add, etc. Triggered by Hooks or applied directly.
type Effect struct {
	ent.Schema
}

func (Effect) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Display name, e.g., 'XP Drain', 'Shadow Bind'"),
		field.String("description").
			Default("").
			Comment("Human-readable description"),
		field.String("effect_type").
			Comment("xp_drain|xp_gain|xp_set|bind_point_set|hp_change|stamina_change|mana_change|message|send_message|teleport|apply_effect|tag_add|tag_remove"),
		field.JSON("parameters", map[string]interface{}{}).
			Default(map[string]interface{}{}).
			Comment(`Type-specific params: {amount:500}, {room_id:12}, etc. send_message: {message:"...",channel:"room"|"yell"|"shout"|"tell"|"whisper"|"chat"|"newbie"|"trade"|"ooc"|"admin"|"emote",target_type:"player"|"npc",target_id:123}`),
		field.String("stack_mode").
			Default("replace").
			Comment("replace|refresh|stack — how stacking works"),
		field.Int("stack_limit").
			Default(1).
			Comment("Max stacks when stack_mode=stack"),
		field.Bool("is_permanent").
			Default(false).
			Comment("If true, effect does not auto-expire"),
		field.Int("duration_secs").
			Default(0).
			Comment("0=instant/one-shot; >0=expires after N seconds"),
		field.JSON("messages", map[string]string{}).
			Default(map[string]string{}).
			Comment("Optional {on_start: '...', on_end: '...'}"),
	}
}

func (Effect) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hooks", EffectHook.Type),
		edge.To("active_effect_instances", ActiveEffect.Type),
		edge.To("triggers", Trigger.Type),
	}
}