package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FactionCategory holds the schema definition for the FactionCategory entity.
type FactionCategory struct {
	ent.Schema
}

func (FactionCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("world_id").
			Default("1").
			Comment("World this faction category belongs to (for multi-world support)"),
		field.String("name").
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
		field.Bool("initial_config").
			Default(false).
			Comment("If true, this category appears in the character creation wizard"),
	}
}

func (FactionCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("world", World.Type).Ref("faction_categories"),
		edge.To("factions", Faction.Type),
	}
}

// Indexes of the FactionCategory.
func (FactionCategory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "world_id").Unique(),
	}
}
