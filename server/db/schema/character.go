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
		field.String("npc_skill_id").
			Optional().
			Comment("NPC skill identifier (e.g., 'druid_heal')"),
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
			Unique(),
		edge.To("available_talents", AvailableTalent.Type),
		edge.To("skills", CharacterSkill.Type),
		edge.To("talents", CharacterTalent.Type),
		edge.To("tags", CharacterTag.Type),
		edge.To("faction_memberships", CharacterFaction.Type),
	}
}