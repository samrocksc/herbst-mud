package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Quest holds the schema definition for the Quest entity.
type Quest struct {
	ent.Schema
}

func (Quest) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("description"),
		field.Strings("prerequisite_quest_ids").Optional(),
		field.JSON("objectives", []QuestObjective{}),
		field.JSON("rewards", QuestRewards{}),
		field.Enum("repeat_mode").Values("none", "cooldown", "always").Default("none"),
		field.Int("cooldown_hours").Default(0),
		field.Bool("is_active").Default(true),
	}
}

func (Quest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("progress", QuestProgress.Type),
	}
}