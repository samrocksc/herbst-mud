package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NPCSkill holds the schema definition for the NPC Skill join table.
type NPCSkill struct {
	ent.Schema
}

// Fields of the NPCSkill.
func (NPCSkill) Fields() []ent.Field {
	return []ent.Field{
		field.Int("slot").
			Comment("Skill slot 1-5, same as player classless skills"),
	}
}

// Edges of the NPCSkill.
func (NPCSkill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("npc_template", NPCTemplate.Type).
			Ref("npc_skills").
			Required(),
		edge.From("skill", Skill.Type).
			Ref("npc_skills").
			Required(),
	}
}
