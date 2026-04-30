package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CompetencyLevelThreshold holds the schema definition for the CompetencyLevelThreshold entity.
type CompetencyLevelThreshold struct {
	ent.Schema
}

// Fields of the CompetencyLevelThreshold.
func (CompetencyLevelThreshold) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(128)",
			}).
			Unique().
			Comment("Composite key, e.g., blades-1, blades-2"),
		field.Int("level").
			Comment("Competency level 1-10"),
		field.Int("xp_required").
			Comment("Cumulative XP needed to reach this level"),
		field.Float("damage_multiplier").
			Default(1.0).
			Comment("Damage multiplier at this competency level"),
		field.Float("defense_multiplier").
			Default(1.0).
			Comment("Defense multiplier at this competency level"),
	}
}

// Edges of the CompetencyLevelThreshold.
func (CompetencyLevelThreshold) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", CompetencyCategory.Type).
			Ref("thresholds").
			Unique().
			Required(),
	}
}