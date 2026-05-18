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
			Default("default").
			Comment("World this NPC template belongs to (for multi-world support)"),
		field.String("name"),
		field.Text("description"),
		field.String("race"),
		field.Enum("disposition").
			Values("hostile", "friendly", "neutral").
			Default("neutral"),
		field.Int("level").
			Default(1),
		field.Int("xp_value").
			Default(0).
			Comment("Base XP awarded when this NPC is killed by a player"),
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
	}
}
