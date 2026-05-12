package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// CharacterIgnore holds player ignore lists.
type CharacterIgnore struct {
	ent.Schema
}

func (CharacterIgnore) Fields() []ent.Field {
	return []ent.Field{
		field.Int("ignoredCharacterId").
			Comment("ID of the character being ignored"),
		field.Time("ignoredAt").
			Default(time.Now),
		field.String("reason").
			Optional().
			Comment("Optional reason for ignoring"),
	}
}

func (CharacterIgnore) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ignorer", Character.Type).
			Ref("ignoring").
			Unique(),
	}
}