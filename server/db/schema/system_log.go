package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SystemLog holds the schema definition for the SystemLog entity.
// SystemLog is used for tracking important events like shop transactions.
type SystemLog struct {
	ent.Schema
}

// Fields of the SystemLog.
func (SystemLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("action").
			Comment("Action type (e.g., shop_buy, shop_sell, npc_death, quest_complete)"),
		field.Int("character_id").
			Optional().
			Comment("Optional FK to character involved in the log event"),
		field.String("details").
			Optional().
			Comment("JSON details about the action"),
		field.Time("timestamp").
			Default(time.Now).
			Comment("When the log entry was created"),
	}
}

// Edges of the SystemLog.
func (SystemLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the SystemLog.
func (SystemLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("action"),
		index.Fields("timestamp"),
		index.Fields("action", "timestamp"),
	}
}
