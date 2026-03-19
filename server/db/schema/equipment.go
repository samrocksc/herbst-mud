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
		field.String("slot"), // e.g., "head", "chest", "weapon", "legs"
		field.Int("level").
			Default(1),
		field.Int("weight").
			Default(0),
		field.Bool("isEquipped").
			Default(false),
		// Item system fields (GitHub #89)
		field.Bool("isImmovable").
			Default(false).
			Comment("Cannot be picked up if true"),
		field.String("color").
			Default("").
			Comment("Custom display color (e.g., gold for immovable items)"),
		field.Bool("isVisible").
			Default(true).
			Comment("Shown in room list"),
		field.String("itemType").
			Default("misc").
			Comment("weapon|armor|consumable|quest|misc|container"),
		// Container system fields (GitHub #143)
		field.Bool("isContainer").
			Default(false).
			Comment("Whether this item can contain other items"),
		field.Int("capacity").
			Default(0).
			Comment("Maximum number of items this container can hold"),
		field.Bool("isLocked").
			Default(false).
			Comment("Whether the container is locked"),
		field.String("lockKey").
			Default("").
			Comment("Key ID required to unlock this container"),
		// Examine system fields
		field.String("examineDesc").
			Default("").
			Comment("Detailed description shown with examine command"),
		field.JSON("hiddenDetails", []map[string]any{}).
			Default([]map[string]any{}).
			Comment("Details revealed based on examine skill"),
		field.Int("hiddenThreshold").
			Default(0).
			Comment("Examine skill required to reveal hidden details"),
		// Readable item fields (GitHub #141)
		field.Bool("isReadable").
			Default(false).
			Comment("Item has readable content"),
		field.String("content").
			Default("").
			Comment("Text content for readable items"),
		field.String("readSkill").
			Optional().
			Comment("Skill required to read (e.g., 'tech')"),
		field.Int("readSkillLevel").
			Default(0).
			Comment("Minimum skill level required"),
	}
}

// Edges of the Equipment.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("equipment").
			Unique(),
		// Container system (GitHub #143) - self-referential edge for container contents
		edge.From("container", Equipment.Type).
			Ref("contents").
			Unique(),
		edge.To("contents", Equipment.Type).
			Comment("Items contained within this container"),
	}
}