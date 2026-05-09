package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// QuestProgress holds the schema definition for the QuestProgress entity.
type QuestProgress struct {
	ent.Schema
}

func (QuestProgress) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("active", "completed", "failed", "abandoned").Default("active"),
		field.Time("started_at"),
		field.Time("completed_at").Optional().Nillable(),
		field.Int("current_step").Default(0),
		field.JSON("objective_counts", map[string]int{}),
	}
}

func (QuestProgress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).Ref("quest_progress").Unique().Required(),
		edge.From("quest", Quest.Type).Ref("progress").Unique().Required(),
	}
}