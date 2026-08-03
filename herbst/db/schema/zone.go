package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Zone groups rooms into named geographic/quest areas.
type Zone struct {
	ent.Schema
}

// Fields of the Zone.
func (Zone) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique(),
		field.String("world_id").
			Default("1").
			Comment("World this zone belongs to (for multi-world support)"),
		field.String("name"),
		field.Text("description").
			Optional(),
		field.Int("min_level").
			Default(1),
		field.String("parent_zone_id").
			Optional().
			Comment("Parent zone for sub-zones (e.g. Old Car Garage inside Junkyard)"),
		field.String("color").
			Optional().
			Comment("Hex color for map UI, e.g. #8b4513"),
	}
}

// Edges of the Zone.
func (Zone) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent", Zone.Type).
			Ref("children").
			Field("parent_zone_id").
			Unique(),
		edge.To("children", Zone.Type),
		edge.From("rooms", Room.Type).Ref("zones"),
	}
}
