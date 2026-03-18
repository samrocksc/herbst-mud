package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Equipment holds the schema definition for the Equipment entity.
type Equipment struct {
	ent.Schema
}

// Fields of the Equipment.
func (Equipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		// Look/examine description fields (look-11)
		field.String("shortDesc").
			Default("").
			Comment("Short description for look command"),
		field.Text("examineDesc").
			Default("").
			Comment("Detailed description for examine command"),
		field.JSON("hiddenDetails", []map[string]interface{}{}).
			Optional().
			Comment("Hidden details revealed by examine"),
		field.JSON("onExamine", []map[string]interface{}{}).
			Optional().
			Comment("Event triggers on examine"),
		field.String("slot"). // e.g., "head", "chest", "weapon", "legs"
			Default(""),
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.Bool("isEquipped").
			Default(false),
		// New fields for item system (GitHub #89)
		field.Bool("isImmovable").
			Default(false).
			Comment("Cannot be picked up if true"),
		field.String("color").
			Default("").
			Comment("Custom display color (e.g., gold for immovable items)"),
		field.Bool("isVisible").
			Default(true).
			Comment("Shown in room list"),
		field.Bool("isContainer").
			Default(false).
			Comment("Can hold other items"),
		field.String("itemType").
			Default("misc").
			Comment("weapon|armor|consumable|quest|misc"),
		// Weapon-specific fields (GitHub #92)
		field.Int("minDamage").
			Default(1).
			Comment("Minimum damage for weapons"),
		field.Int("maxDamage").
			Default(1).
			Comment("Maximum damage for weapons"),
		field.String("weaponType").
			Default("").
			Comment("sword|dagger|staff|pipe|brawling - weapon style"),
		field.String("classRestriction").
			Default("").
			Comment("Class that can use this weapon (e.g., warrior, chef)"),
		field.Bool("isDroppable").
			Default(true).
			Comment("Can be dropped by NPCs"),
		field.Bool("guaranteedDrop").
			Default(false).
			Comment("Always drops on first NPC kill"),
		// Readable items (GitHub #look-07)
		field.Bool("isReadable").
			Default(false).
			Comment("Can be read if true"),
		field.Text("content").
			Default("").
			Comment("Text content for readable items"),
		field.String("readSkill").
			Default("").
			Comment("Skill required to read (e.g., tech, lore)"),
		field.Int("readSkillLevel").
			Default(0).
			Comment("Required skill level to read"),
		field.Text("decryptedContent").
			Default("").
			Comment("Content shown when skill check passes"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
		// Character inventory (GitHub #92)
		edge.To("character", Character.Type).
			Unique().
			Comment("Character carrying this item"),
	}
}