package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NPCTemplate holds the schema definition for the NPC Template entity.
type NPCTemplate struct {
	ent.Schema
}

// Fields of the NPCTemplate.
func (NPCTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique(),
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
	}
}

// Edges of the NPCTemplate.
func (NPCTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("npc_skills", NPCSkill.Type),
	}
}