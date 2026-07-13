package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DefaultRespawnRooms is the default value for the respawn_rooms field.
var DefaultRespawnRooms = []string{}

// NPCTemplate holds the schema definition for the NPC Template entity.
type NPCTemplate struct {
	ent.Schema
}

// Fields of the NPCTemplate.
func (NPCTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique(),
		field.String("slug").
			Unique().
			Optional().
			Comment("Globally unique slug derived from name, e.g., goblin_scout"),
		field.String("world_id").
			Default("1").
			Comment("World this NPC template belongs to (for multi-world support)"),
		field.String("name"),
		field.Text("description"),
		field.Int("race_id").
			Optional().
			Comment("FK to Race.id — NPC's race"),
		field.Enum("disposition").
			Values("hostile", "friendly", "neutral").
			Default("neutral"),
		field.Int("level").
			Default(1),
		field.Int("xp_value").
			Default(0).
			Comment("Base XP awarded when this NPC is killed by a player"),
		field.Float("xp_multiplier").
			Default(1.0).
			Comment("Multiplier for XP gain scaling (1.0 = normal)"),
		field.JSON("skills", map[string]int{}),
		field.Strings("trades_with"),
		field.Text("greeting"),
		field.JSON("respawn_rooms", []string{}).
			Optional().
			Comment("Array of room IDs where this NPC can respawn (nil or empty = no respawn)"),
		field.Int("respawn_cooldown").
			Optional().
			Default(60).
			Comment("Seconds before this NPC respawns after death (0 = immediate, nil = default 60)"),
		field.Enum("roam_pattern").
			Values("static", "wander", "patrol", "return_home").
			Default("static").
			Comment("NPC roaming behavior: static (never moves), wander (random exits), patrol (round-robin), return_home (walk back to home)"),
		field.Strings("roam_zone_ids").
			Optional().
			Comment("Zones this NPC can roam inside. Empty = no zone restriction."),
		field.Int("roam_interval_seconds").
			Optional().
			Default(60).
			Comment("How often (in seconds) this NPC is eligible to move. Default 60."),
		field.Int("roam_pause_min_seconds").
			Optional().
			Default(15).
			Comment("Minimum seconds to add as jitter before the next move."),
		field.Int("roam_pause_max_seconds").
			Optional().
			Default(120).
			Comment("Maximum seconds to add as jitter before the next move."),
		field.Time("last_moved_at").
			Optional().
			Nillable().
			Comment("Last time this NPC was moved by the roaming service. Used to enforce roam_interval_seconds."),
		field.Bool("notify_on_enter").
			Optional().
			Default(false).
			Comment("If true, emit a chat/notification event when this NPC enters a room."),
	}
}

// Edges of the NPCTemplate.
func (NPCTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("npc_abilities", NPCAbility.Type),
		edge.To("hooks", EffectHook.Type),
		edge.To("dialog_nodes", DialogNode.Type),
		edge.From("characters", Character.Type).
			Ref("npcTemplate"),
		edge.From("race", Race.Type).Ref("npc_templates").Field("race_id").Unique(),
	}
}
