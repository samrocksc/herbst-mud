package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterCompetency holds the schema definition for the CharacterCompetency entity.
type CharacterCompetency struct {
	ent.Schema
}

// Fields of the CharacterCompetency.
func (CharacterCompetency) Fields() []ent.Field {
	return []ent.Field{
		field.Int("xp").
			Default(0).
			Comment("Accumulated XP in this competency (after multiplier applied)"),
		field.Int("level").
			Default(0).
			Comment("Cached competency level (recomputed on XP award)"),
	}
}

// Edges of the CharacterCompetency.
func (CharacterCompetency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("competencies").
			Unique().
			Required(),
		edge.From("category", CompetencyCategory.Type).
			Ref("character_competencies").
			Unique().
			Required(),
	}
}