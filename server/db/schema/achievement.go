package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Achievement holds the schema definition for the Achievement entity.
type Achievement struct {
	ent.Schema
}

// Fields of the Achievement.
func (Achievement) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description").
			Default(""),
		field.String("icon").
			Optional().
			Comment("Emoji or icon identifier for the achievement badge"),
		field.Int("xp_reward").
			Default(0).
			Comment("XP awarded when the achievement is unlocked"),
		field.String("criteria").
			Optional().
			Comment("JSON criteria describing how the achievement is earned"),
	}
}