package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Faction holds the schema definition for the Faction entity.
type Faction struct {
	ent.Schema
}

func (Faction) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Comment("e.g., ninja, foot_clan, surf_warden"),
		field.String("display_name").
			Comment("e.g., Ninja, Foot Clan, Surf Warden"),
		field.String("description").
			Optional(),
	}
}

func (Faction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", FactionCategory.Type).
			Ref("factions").
			Unique(),
		edge.To("required_tags", FactionRequiredTag.Type),
		edge.To("character_factions", CharacterFaction.Type),
		edge.To("skills", Skill.Type),
	}
}
