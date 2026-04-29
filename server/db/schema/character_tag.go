package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// CharacterTag holds the schema definition for the CharacterTag entity.
type CharacterTag struct {
	ent.Schema
}

func (CharacterTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag").
			Comment("Tag identifier, e.g., first_class, wizard_complete"),
		field.String("source").
			Default("system").
			Comment("How the tag was earned: system, quest, achievement, admin"),
		field.Time("earned_at").
			Default(time.Now),
	}
}

func (CharacterTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("tags").
			Unique(),
	}
}
