package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Quest holds the schema definition for the Quest entity.
// Quests are named, sequenced sets of objectives that characters progress through.
type Quest struct {
	ent.Schema
}

// Fields of the Quest.
func (Quest) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("description"),
		field.Strings("prerequisite_quest_ids").Optional().
			Comment("Quest IDs that must be completed before this quest can be accepted"),
		field.JSON("objectives", []QuestObjective{}).
			Comment("Ordered list of quest objectives"),
		field.JSON("rewards", QuestRewards{}).
			Comment("Rewards granted on quest completion"),
		field.Enum("repeat_mode").Values("none", "cooldown", "always").Default("none").
			Comment("Whether and how this quest can be repeated"),
		field.Enum("main_type").Values("hunter", "collector", "explorer", "general").Default("general").
			Comment("Primary quest type for categorization"),
		field.Int("cooldown_hours").Default(0).
			Comment("Hours before re-accept allowed (only if repeat_mode=cooldown)"),
		field.Bool("is_active").Default(true).
			Comment("Whether this quest can be accepted"),
	}
}

// Edges of the Quest.
func (Quest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("progress", QuestProgress.Type),
	}
}