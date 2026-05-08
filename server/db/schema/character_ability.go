package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterAbility holds the schema definition for the CharacterAbility join entity.
// Links a Character to an Ability with a slot number.
type CharacterAbility struct {
	ent.Schema
}

// Fields of the CharacterAbility.
func (CharacterAbility) Fields() []ent.Field {
	return []ent.Field{
		field.Int("slot").
			Comment("Ability slot 1-5"),
	}
}

// Edges of the CharacterAbility.
func (CharacterAbility) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("abilities").
			Unique(),
		edge.From("ability", Ability.Type).
			Ref("characters").
			Unique(),
	}
}