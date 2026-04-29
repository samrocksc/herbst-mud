package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FactionCategory holds the schema definition for the FactionCategory entity.
type FactionCategory struct {
	ent.Schema
}

func (FactionCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Comment("e.g., class, alignment"),
		field.String("display_name").
			Comment("e.g., Class, Alignment"),
		field.String("description").
			Optional(),
		field.Int("max_memberships").
			Default(1).
			Comment("How many factions in this category a character can hold"),
		field.Bool("auto_join").
			Default(false).
			Comment("If true, earning required tag auto-joins faction"),
	}
}

func (FactionCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("factions", Faction.Type),
	}
}
