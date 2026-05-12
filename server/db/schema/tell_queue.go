package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// TellQueue holds offline tells with 7-day expiry.
type TellQueue struct {
	ent.Schema
}

func (TellQueue) Fields() []ent.Field {
	return []ent.Field{
		field.Int("senderId").
			Comment("ID of the character who sent the tell"),
		field.String("senderName").
			Comment("Name of sender (denormalized for when sender is deleted)"),
		field.String("message").
			Comment("The tell message content"),
		field.Time("sentAt").
			Default(time.Now),
		field.Time("expiresAt").
			Comment("7 days after sentAt, then auto-deleted"),
		field.Bool("isRead").
			Default(false),
	}
}

func (TellQueue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recipient", Character.Type).
			Ref("tellQueue").
			Unique(),
	}
}