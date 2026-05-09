package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AppLog holds the schema definition for the AppLog entity.
type AppLog struct {
	ent.Schema
}

// Fields of the AppLog.
func (AppLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Comment("Auto-incremented log entry ID"),
		field.String("level").
			Comment("Log level: DEBUG, INFO, WARN, ERROR"),
		field.String("message").
			Comment("Log message text"),
		field.String("service").
			Optional().
			Comment("Service or component that emitted the log"),
		field.Int("character_id").
			Optional().
			Nillable().
			Comment("Optional character ID context"),
		field.Int("room_id").
			Optional().
			Nillable().
			Comment("Optional room ID context"),
		field.String("template_id").
			Optional().
			Comment("Optional NPC template ID context"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Arbitrary key-value metadata"),
		field.Time("created_at").
			Default(time.Now).
			Comment("When the log entry was created"),
	}
}

// Edges of the AppLog.
func (AppLog) Edges() []ent.Edge {
	return nil
}

// Indexes of the AppLog.
func (AppLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("level"),
		index.Fields("service"),
		index.Fields("created_at"),
	}
}