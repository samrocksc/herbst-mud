package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
)

// World holds the schema definition for the World entity.
type World struct {
	ent.Schema
}

// Fields of the World.
func (World) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("title"),
		field.String("description").
			Optional(),
		field.Bool("active").
			Default(false),
	}
}

// Edges of the World.
func (World) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
		edge.To("races", Race.Type),
		edge.To("genders", Gender.Type),
		edge.To("tags", Tag.Type),
		edge.To("social_commands", SocialCommand.Type),
		edge.To("faction_categories", FactionCategory.Type),
		edge.To("effect_hooks", EffectHook.Type),
	}
}
