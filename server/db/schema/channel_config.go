package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ChannelConfig holds global configuration for a chat channel.
type ChannelConfig struct {
	ent.Schema
}

func (ChannelConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Comment("Unique channel identifier, e.g., 'global', 'trade', 'newbie'"),
		field.String("description").
			Default("").
			Comment("Human-readable description of the channel"),
		field.String("color").
			Default("#FFFFFF").
			Comment("Hex color code for the channel messages"),
		field.Bool("default_enabled").
			Default(true).
			Comment("Whether characters are joined by default"),
		field.Int("cooldown_seconds").
			Default(0).
			Comment("Message cooldown in seconds"),
		field.Bool("admin_only").
			Default(false).
			Comment("Whether only admin characters can speak in this channel"),
	}
}

func (ChannelConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}
