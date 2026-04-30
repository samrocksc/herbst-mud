package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CompetencyCategory holds the schema definition for the CompetencyCategory entity.
type CompetencyCategory struct {
	ent.Schema
}

// Fields of the CompetencyCategory.
func (CompetencyCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(64)",
			}).
			Unique().
			Comment("e.g., blades, staves, knives, martial, brawling, tech, light_armor, cloth_armor, heavy_armor"),
		field.String("name").
			Unique().
			Comment("Display name, e.g., Blades, Staves"),
		field.Float("xp_multiplier").
			Default(0.20).
			Comment("Multiplier applied to raw XP before storing in character_competency"),
	}
}

// Edges of the CompetencyCategory.
func (CompetencyCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("thresholds", CompetencyLevelThreshold.Type),
		edge.To("character_competencies", CharacterCompetency.Type),
	}
}