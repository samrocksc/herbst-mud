package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DamageLog holds the schema definition for the DamageLog entity.
type DamageLog struct {
	ent.Schema
}

// Fields of the DamageLog.
func (DamageLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("attacker_id").
			Comment("ID of character who dealt damage"),
		field.Int("target_id").
			Comment("ID of character/NPC who took damage"),
		field.Int("damage").
			Comment("Damage amount"),
		field.Time("created_at").
			Default(time.Now).
			Comment("When the damage was dealt"),
	}
}

// Edges of the DamageLog.
func (DamageLog) Edges() []ent.Edge {
	return nil
}