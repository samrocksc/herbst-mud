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
		field.Int("startingRoomId"),
		field.Bool("is_admin").
			Default(false),
		field.Int("class_id").
			Default(0),
		field.Int("race_id").
			Default(0),
		field.Int("gender_id").
			Default(0),
		field.Int("level").
			Default(1),
		field.Int("experience").
			Default(0),
		field.Int("skill_points").
			Default(0),
		field.Int("talent_points").
			Default(0),
		field.JSON("stats", map[string]int{
			"strength":     10,
			"dexterity":     10,
			"constitution": 10,
			"intelligence": 10,
			"wisdom":        10,
			"charisma":      10,
		}),
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
		edge.To("skills", Skill.Type),
		edge.To("talents", Talent.Type),
	}
}