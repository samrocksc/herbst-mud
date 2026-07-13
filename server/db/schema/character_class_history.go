package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterClassHistory holds the schema definition for the CharacterClassHistory entity.
type CharacterClassHistory struct {
	ent.Schema
}

func (CharacterClassHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("character_id"),
		field.Int("faction_id"),
		field.String("faction_name"),
		field.Time("joined_at"),
		field.Time("left_at").Optional().Nillable(),
		field.String("reason").Optional().Default(""),
	}
}

func (CharacterClassHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("class_history").
			Field("character_id").
			Unique().
			Required(),
	}
}