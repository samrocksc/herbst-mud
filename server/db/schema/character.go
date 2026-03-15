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
		field.String("password").
			Optional(),
		field.Bool("isNPC").
			Default(false),
		field.Int("currentRoomId"),
		field.Int("startingRoomId"),
		field.Bool("is_admin").
			Default(false),
		field.Int("hitpoints").
			Default(100).
			Comment("Health points"),
		field.Int("max_hitpoints").
			Default(100).
			Comment("Maximum health points"),
		field.Int("stamina").
			Default(50).
			Comment("Stamina points"),
		field.Int("max_stamina").
			Default(50).
			Comment("Maximum stamina points"),
		field.Int("mana").
			Default(25).
			Comment("Magic/energy points"),
		field.Int("max_mana").
			Default(25).
			Comment("Maximum magic/energy points"),
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
