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
		field.String("world_id").
			Default("default").
			Comment("World this room belongs to (for multi-world support)"),
		field.String("description"),
		field.Bool("isStartingRoom").
			Default(false),
		field.Bool("isRootRoom").
			Default(false).
			Comment("Only one room can be root; new characters spawn here"),
		field.JSON("exits", map[string]int{}),
		field.Enum("atmosphere").
			Values("air", "water", "wind").
			Default("air"),
		field.Int("posX").
			Default(0).
			Optional(),
		field.Int("posY").
			Default(0).
			Optional(),
		field.Int("posZ").
			Default(0).
			Optional().
			Comment("Z-level for map rendering; 0 = ground floor"),
		field.Int("version").
			Default(1),
	}
}

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
		edge.To("equipment", Equipment.Type),
	}
}