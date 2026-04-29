package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// CharacterFaction holds the schema definition for the CharacterFaction entity.
type CharacterFaction struct {
	ent.Schema
}

func (CharacterFaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int("reputation").
			Default(0).
			Comment("Faction reputation 0-100"),
		field.String("status").
			Default("active").
			Comment("active, expelled, voluntarily_left"),
		field.Time("joined_at").
			Default(time.Now),
	}
}

func (CharacterFaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("faction_memberships").
			Unique(),
		edge.From("faction", Faction.Type).
			Ref("character_factions").
			Unique(),
	}
}
