package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("world_id").
			Default("1").
			Comment("World this tag belongs to (for multi-world support)"),
		field.String("name").
			Comment("Display name for the tag, e.g. 'fire', 'magic', 'warrior'"),
		field.String("color").
			Optional().
			Comment("Hex color for UI display, e.g. '#ff6b6b'. Defaults handled in UI layer."),
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("world", World.Type).Ref("tags"),
		edge.To("races", Race.Type),
	}
}

// Indexes of the Tag.
func (Tag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "world_id").Unique(),
	}
}
