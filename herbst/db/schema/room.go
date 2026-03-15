package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Atmosphere type for room environment
type Atmosphere string

const (
	AtmosphereAir  Atmosphere = "air"
	AtmosphereWater Atmosphere = "water"
	AtmosphereWind  Atmosphere = "wind"
)

// Room holds the schema definition for the Room entity.
type Room struct {
	ent.Schema
}

// Fields of the Room.
func (Room) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Bool("isStartingRoom").
			Default(false),
		field.JSON("exits", map[string]int{}),
		field.Enum("atmosphere").
			Values("air", "water", "wind").
			Default("air"),
	}
}

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
		edge.To("equipment", Equipment.Type),
	}
}