package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FactionRequiredTag holds the schema definition for the FactionRequiredTag entity.
type FactionRequiredTag struct {
	ent.Schema
}

func (FactionRequiredTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("required_tag").
			Comment("Tag the character must have to pledge to this faction"),
	}
}

func (FactionRequiredTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("faction", Faction.Type).
			Ref("required_tags").
			Unique(),
	}
}
