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
		field.Bool("isImmovable").Default(false).Comment("Cannot be picked up"),
		field.Bool("isContainer").Default(false).Comment("Can hold other items"),
		field.Bool("isVisible").Default(true).Comment("Shown in room list"),
		field.Bool("isReadable").Default(false).Comment("Has text content"),
		field.Text("content").Default("").Comment("Text content if readable"),
		field.String("readSkill").Default("").Comment("Skill required to read (e.g., tech, lore)"),
		field.Int("readSkillLevel").Default(0).Comment("Required skill level"),
		field.String("itemType").Default("misc").Comment("weapon|armor|consumable|quest|misc"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
	}
}