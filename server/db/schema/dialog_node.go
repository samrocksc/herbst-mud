package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DialogResponse is a player response option within a dialog node.
type DialogResponse struct {
	Label        string `json:"label"`
	NextNodeID   string `json:"next_node_id"`
	Condition    string `json:"condition,omitempty"`
	QuestOfferID string `json:"quest_offer_id,omitempty"`
	DeclineNodeID string `json:"decline_node_id,omitempty"`
	Effects      []int  `json:"effects,omitempty"`
}

// DialogNode holds the schema for NPC dialog tree nodes.
type DialogNode struct {
	ent.Schema
}

// Fields of the DialogNode.
func (DialogNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().
			Comment("Unique node identifier, e.g. node_001_greeting"),
		field.String("world_id").
			Default("default").
			Comment("World this dialog node belongs to (for multi-world support)"),
		field.String("npc_text").
			Comment("What the NPC says at this node"),
		field.JSON("responses", []DialogResponse{}).
			Optional().
			Comment("Player response options"),
		field.Bool("is_entry").Default(false).
			Comment("Whether this is the conversation start node"),
		field.String("entry_condition").Optional().
			Comment("SPICE expression gate for conditional entry"),
		field.JSON("on_enter_effects", []int{}).
			Optional().
			Comment("Effect IDs applied when this node is reached"),
	}
}

// Edges of the DialogNode.
func (DialogNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("npc_template", NPCTemplate.Type).
			Ref("dialog_nodes").
			Unique().
			Required(),
	}
}