package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Character holds the schema definition for the Character entity.
type Character struct {
	ent.Schema
}

// Fields of the Character.
func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("isNPC").
			Default(false),
		field.Int("currentRoomId"),
	}
}

// Edges of the Character.
func (Character) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("characters").
			Unique(),
		edge.To("room", Room.Type).
			Field("currentRoomId").
			Required().
			Unique(),
	}
}
