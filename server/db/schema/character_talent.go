package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterTalent holds the schema definition for the CharacterTalent entity.
type CharacterTalent struct {
	ent.Schema
}

// Fields of the CharacterTalent.
func (CharacterTalent) Fields() []ent.Field {
	return []ent.Field{
		field.Int("slot").
			Default(0).
			Comment("Equipment slot 0-3 for quick access"),
	}
}

// Edges of the CharacterTalent.
func (CharacterTalent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("talents").
			Unique(),
		edge.From("talent", Talent.Type).
			Ref("characters").
			Unique(),
	}
}