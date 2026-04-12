package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// GameConfig holds a single key/value config pair.
// Use as a simple key/value store for game settings (fountain_room_id, etc.)
type GameConfig struct {
	ent.Schema
}

// Fields of the GameConfig.
func (GameConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").
			Unique().
			Comment("Config key, e.g. fountain_room_id"),
		field.String("value").
			Comment("Config value"),
	}
}

// Edges of the GameConfig.
func (GameConfig) Edges() []ent.Edge {
	return nil
}