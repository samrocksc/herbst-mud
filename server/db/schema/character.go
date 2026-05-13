package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Character holds the schema definition for the Character entity.
type Character struct {
	ent.Schema
}

// Fields of the Character.
func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("password").
			Optional(),
		field.Bool("isNPC").
			Default(false),
		field.Int("currentRoomId"),
		field.Int("startingRoomId"),
		field.Int("respawnRoomId").
			Default(5).
			Comment("Room ID where character respawns after death (default: The Hole)"),
		field.Bool("is_admin").
			Default(false),
		field.Bool("is_immortal").
			Default(false).
			Comment("Character cannot be killed - takes damage but never dies"),
		field.Bool("is_test").
			Default(false).
			Comment("Test player — can use /debug commands"),
		field.Bool("is_instance").
			Default(false).
			Comment("True if this NPC is an instance of a template"),
		field.Int("instance_number").
			Default(0).
			Comment("Auto-incremented instance number per template"),
		field.String("npc_template_id").
			Optional().
			Comment("Foreign key to npc_template ID"),
		field.String("npc_skill_id").
			Optional().
			Comment("NPC skill identifier (e.g., 'druid_heal')"),
		field.String("currentWorld").
			Default("default").
			Comment("World this character belongs to (for multi-world support)"),
		field.Int("npc_skill_cooldown").
			Default(0).
			Comment("Current cooldown ticks on NPC skill"),
		field.Int("hitpoints").
			Default(100),
		field.Int("max_hitpoints").
			Default(100),
		field.Int("stamina").
			Default(50),
		field.Int("max_stamina").
			Default(50),
		field.Int("mana").
			Default(25),
		field.Int("max_mana").
			Default(25),
		field.String("race").
			Default("human"),
		field.String("class").
			Default("adventurer"),
		field.String("specialty").
			Optional().
			Comment("Class specialty (e.g., fighter for warrior)"),
		field.Int("level").
			Default(1),
		field.Int("xp").
			Default(0).
			Comment("Current accumulated experience points"),
		field.Time("died_at").
			Optional().
			Nillable().
			Comment("When this NPC died (nil if alive or a player character)"),
			field.Time("lastSeenAt").
				Optional().
				Nillable().
				Comment("When the character was last online"),
		field.Int("constitution").
			Default(10),
		field.String("gender").
			Optional(),
		field.String("description").
			Optional(),
		field.Int("strength").
			Default(10),
		field.Int("dexterity").
			Default(10),
		field.Int("intelligence").
			Default(10),
		field.Int("wisdom").
			Default(10),
		field.Int("skill_blades").
			Default(0),
		field.Int("skill_staves").
			Default(0),
		field.Int("skill_knives").
			Default(0),
		field.Int("skill_martial").
			Default(0),
		field.Int("skill_brawling").
			Default(0),
		field.Int("skill_tech").
			Default(0),
		field.Int("skill_light_armor").
			Default(0),
		field.Int("skill_cloth_armor").
			Default(0),
		field.Int("skill_heavy_armor").
			Default(0),
	}
}

// Edges of the Character.
func (Character) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("characters").
			Unique(),
		edge.To("room", Room.Type).
			Field("currentRoomId").
			Required().
			Unique(),
		edge.To("npcTemplate", NPCTemplate.Type).
			Unique().
			Field("npc_template_id"),
		edge.To("abilities", CharacterAbility.Type),
		edge.To("tags", CharacterTag.Type),
		edge.To("faction_memberships", CharacterFaction.Type),
		edge.To("competencies", CharacterCompetency.Type),
		edge.To("active_effects", ActiveEffect.Type),
		edge.To("quest_progress", QuestProgress.Type),
		edge.To("channelSettings", CharacterChannel.Type),
		edge.To("ignoring", CharacterIgnore.Type),
		edge.To("tellQueue", TellQueue.Type),
	}
}