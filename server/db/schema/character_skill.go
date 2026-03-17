package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterSkill holds the schema definition for the CharacterSkill entity.
type CharacterSkill struct {
	ent.Schema
}

// Fields of the CharacterSkill.
func (CharacterSkill) Fields() []ent.Field {
	return []ent.Field{
		field.Int("level").
			Default(1).
			Comment("Current skill level"),
		field.Int("experience").
			Default(0).
			Comment("Experience points toward next level"),
	}
}

// Edges of the CharacterSkill.
func (CharacterSkill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("skills").
			Unique(),
		edge.From("skill", Skill.Type).
			Ref("characters").
			Unique(),
	}
}