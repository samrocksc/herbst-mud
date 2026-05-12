package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CharacterChannel holds per-character channel settings (chat, newbie, trade, etc.)
type CharacterChannel struct {
	ent.Schema
}

func (CharacterChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("chatEnabled").
			Default(true).
			Comment("General chat channel enabled"),
		field.Bool("newbieEnabled").
			Default(true).
			Comment("Newbie help channel enabled"),
		field.Bool("tradeEnabled").
			Default(false).
			Comment("Trade/sell channel enabled"),
		field.Bool("clanEnabled").
			Default(true).
			Comment("Clan/gossip channel enabled"),
		field.Bool("auctionEnabled").
			Default(false).
			Comment("Auction channel enabled"),
		field.String("chatColor").
			Default("#00FF00").
			Comment("Hex color for chat messages"),
		field.String("newbieColor").
			Default("#00FFFF").
			Comment("Hex color for newbie messages"),
		field.String("tradeColor").
			Default("#FFFF00").
			Comment("Hex color for trade messages"),
		field.String("clanColor").
			Default("#FF00FF").
			Comment("Hex color for clan messages"),
		field.Bool("timestamps").
			Default(true).
			Comment("Show timestamps in chat"),
		field.Bool("profanityFilter").
			Default(false).
			Comment("Filter profanity in incoming messages"),
	}
}

func (CharacterChannel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("character", Character.Type).
			Ref("channelSettings").
			Unique(),
	}
}