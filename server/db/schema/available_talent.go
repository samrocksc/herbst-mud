package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AvailableTalent holds the schema definition for the AvailableTalent entity.
// This tracks talents that a character has unlocked but not necessarily equipped.
type AvailableTalent struct {
	ent.Schema
}

// Fields of the AvailableTalent.
func (AvailableTalent) Fields() []ent.Field {
	return []ent.Field{
		field.String("unlock_reason").
			Optional().
			Default("level_up").
			Comment("Reason talent was unlocked: level_up, quest, skill_trainer, item"),
		field.Int("unlocked_at_level").
			Default(1).
			Comment("Character level when talent was unlocked"),
	}
}

// Edges of the AvailableTalent.
func (AvailableTalent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("available_talents").
			Unique(),
		edge.From("talent", Talent.Type).
			Ref("available_to_characters").
			Unique(),
	}
}