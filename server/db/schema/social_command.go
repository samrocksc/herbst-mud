package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SocialCommand holds social commands like smile, bow, wave with variants.
type SocialCommand struct {
	ent.Schema
}

func (SocialCommand) Fields() []ent.Field {
	return []ent.Field{
		field.String("world_id").
			Default("1").
			Comment("World this social command belongs to (for multi-world support)"),
		field.String("name").
			Comment("Social command name, e.g., smile, bow, wave"),
		field.String("displayName").
			Comment("Display name for the social"),
		field.String("selfText").
			Comment("Text when using social on self: {player} smiles."),
		field.String("roomText").
			Comment("Text shown to room: {player} smiles at {target}."),
		field.String("targetSelfText").
			Comment("Text shown to target when they are the target"),
		field.String("targetText").
			Comment("Text shown to target: {player} smiles at you."),
		field.String("targetRoomText").
			Comment("Text shown to others in room when target is involved"),
		field.Bool("requiresTarget").
			Default(false).
			Comment("True if social requires a target character"),
		field.Bool("isEmote").
			Default(true).
			Comment("True if this is an emote (displayed with *text*)"),
	}
}

func (SocialCommand) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("world", World.Type).Ref("social_commands"),
	}
}

// Indexes of the SocialCommand.
func (SocialCommand) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "world_id").Unique(),
	}
}
