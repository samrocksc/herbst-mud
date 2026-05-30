package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Trigger holds the schema definition for the Trigger entity.
// Triggers link objects (rooms, equipment) to actions (recipes, effects, dialog nodes).
// When a player interacts with an object via "use", "touch", or "press",
// the corresponding trigger fires and executes its target action.
type Trigger struct {
	ent.Schema
}

// Fields of the Trigger.
func (Trigger) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Display name for the trigger, e.g., 'Use Oven', 'Touch Statue'"),
		field.String("world_id").
			Default("1").
			Comment("World this trigger belongs to (for multi-world support)"),
		field.String("trigger_type").
			Comment("use|touch|press|enter|examine - what action triggers this"),
		field.String("target_type").
			Comment("recipe|effect|dialog_node - what gets executed"),
		field.Int("target_id").
			Comment("FK to target entity (recipe name, effect id, dialog_node id)"),
		field.Int("room_id").
			Optional().
			Nillable().
			Comment("FK to Room - trigger fires when interacting with room objects"),
		field.Int("equipment_id").
			Optional().
			Nillable().
			Comment("FK to Equipment - trigger fires when using this item"),
		field.String("condition").
			Optional().
			Comment("SPICE expression for conditional trigger firing"),
		field.Bool("enabled").
			Default(true).
			Comment("If false, trigger does not fire"),
	}
}

// Edges of the Trigger.
func (Trigger) Edges() []ent.Edge {
	return []ent.Edge{
		// Triggers reference effects
		edge.From("effect", Effect.Type).
			Ref("triggers").
			Unique(),
		// Triggers reference crafting recipes
		edge.From("recipe", CraftingRecipe.Type).
			Ref("triggers").
			Unique(),
		// Triggers reference dialog nodes
		edge.From("dialog_node", DialogNode.Type).
			Ref("triggers").
			Unique(),
	}
}
