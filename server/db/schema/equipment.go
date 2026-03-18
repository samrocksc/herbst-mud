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
		field.String("shortDesc"), // Short description for look command
		field.String("examineDesc"), // Detailed description for examine command
		field.JSON("hiddenDetails", []map[string]interface{}{}), // Hidden details revealed by examine
		field.JSON("onExamine", []map[string]interface{}{}), // Event triggers on examine
		field.String("slot"), // e.g., "head", "chest", "weapon", "legs"
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.Bool("isEquipped").
			Default(false),
		field.Bool("isImmovable").Default(false), // Cannot be picked up
		field.Bool("isContainer").Default(false), // Can hold other items
		field.Bool("isReadable").Default(false), // Has text content
		field.Text("content"), // Text content if readable
		field.String("readSkill"), // Skill required to read (e.g., "tech", "lore")
		field.Int("readSkillLevel").Default(0), // Required skill level
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