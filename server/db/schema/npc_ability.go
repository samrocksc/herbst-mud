package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NPCAbility holds the schema definition for the NPC Ability join table.
type NPCAbility struct {
	ent.Schema
}

// Fields of the NPCAbility.
func (NPCAbility) Fields() []ent.Field {
	return []ent.Field{
		field.Int("slot").
			Comment("Ability slot 1-5, same as player classless abilities"),
	}
}

// Edges of the NPCAbility.
func (NPCAbility) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("npc_template", NPCTemplate.Type).
			Ref("npc_abilities").
			Required(),
		edge.From("ability", Ability.Type).
			Ref("npc_abilities").
			Required(),
	}
}