package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Nillable().
			Default(time.Now).
			Immutable().
			Comment("When the user account was created"),
		field.Time("updated_at").
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the user account was last updated"),
		field.String("email").
			Unique(),
		field.String("password"),
		field.Bool("is_admin").
			Default(false),
		field.Bool("god_mode").
			Default(false).
			Comment("Unkillable mode for the user"),
		field.String("allowed_worlds").
			Optional().
			Comment("Comma-separated list of world IDs this user can access (empty = all worlds for admins)"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("characters", Character.Type),
	}
}
