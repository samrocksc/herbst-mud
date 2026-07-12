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
			Default("1").
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
			Optional().
			StructTag(`json:"posX"`),
		field.Int("posY").
			Default(0).
			Optional().
			StructTag(`json:"posY"`),
		field.Int("posZ").
			Default(0).
			Optional().
			StructTag(`json:"posZ"`).
			Comment("Z-level for map rendering; 0 = ground floor"),
		field.Int("version").
			Default(1),
		field.JSON("tags", []string{}).
			Optional().
			Comment("Tags for station discovery, room features (e.g. pizza_station, forge)"),
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