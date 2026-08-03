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
		field.Bool("isRootRoom").
			Default(false),
		field.JSON("exits", map[string]int{}),
		field.Enum("atmosphere").
			Values("air", "water", "wind").
			Default("air"),
		field.JSON("tags", []string{}).
			Optional().
			Comment("Tags for station discovery (e.g. pizza_station, forge)"),
		field.Strings("zone_ids").
			Optional().
			Comment("Zone memberships for this room. First entry is the primary zone. Sub-zone membership is just appending the sub-zone ID."),
	}
}

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
		edge.To("equipment", Equipment.Type),
		edge.To("zones", Zone.Type),
	}
}