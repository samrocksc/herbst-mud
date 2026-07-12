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

// StatBonuses represents per-class stat bonuses applied during character creation.
type StatBonuses struct {
	Strength     int `json:"strength,omitempty"`
	Dexterity    int `json:"dexterity,omitempty"`
	Constitution int `json:"constitution,omitempty"`
	Intelligence int `json:"intelligence,omitempty"`
	Wisdom       int `json:"wisdom,omitempty"`
	Charisma     int `json:"charisma,omitempty"`
}

// ClassSpecialty represents a specialty within a class faction.
type ClassSpecialty struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (Faction) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("world_id").
			Default("1").
			Comment("World this faction belongs to (for multi-world support)"),
		field.String("display_name").
			Comment("e.g., Ninja, Foot Clan, Surf Warden"),
		field.String("description").
			Optional(),
		field.JSON("member_tags", []string{}).
			Optional().
			Comment("Tags auto-applied to characters when they join this faction"),
		field.JSON("stat_bonuses", StatBonuses{}).
			Optional().
			Comment("Stat bonuses applied during character creation for class factions"),
		field.JSON("specialties", []ClassSpecialty{}).
			Optional().
			Comment("Specialties available for this class faction"),
	}
}

func (Faction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", FactionCategory.Type).
			Ref("factions").
			Unique(),
		edge.To("required_tags", FactionRequiredTag.Type),
		edge.To("character_factions", CharacterFaction.Type),
		edge.To("abilities", Ability.Type),
	}
}
