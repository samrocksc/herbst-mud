package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterRaceHistory holds the schema definition for the CharacterRaceHistory entity.
type CharacterRaceHistory struct {
	ent.Schema
}

func (CharacterRaceHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("character_id"),
		field.Int("race_id").Optional().Nillable(),
		field.String("race_name"),
		field.Time("changed_at"),
		field.String("reason").Optional().Default(""),
	}
}

func (CharacterRaceHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("race_history").
			Field("character_id").
			Unique().
			Required(),
	}
}