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
		field.Int("slot").
			Comment("Skill slot 1-5, same as classless skill system"),
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