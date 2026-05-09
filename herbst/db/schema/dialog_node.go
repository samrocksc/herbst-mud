package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DialogResponse is a player response option within a dialog node.
type DialogResponse struct {
	Label         string `json:"label"`
	NextNodeID    string `json:"next_node_id"`
	Condition     string `json:"condition,omitempty"`
	QuestOfferID  string `json:"quest_offer_id,omitempty"`
	DeclineNodeID string `json:"decline_node_id,omitempty"`
	Effects       []int  `json:"effects,omitempty"`
}

// DialogNode holds the schema for NPC dialog tree nodes.
type DialogNode struct {
	ent.Schema
}

// Fields of the DialogNode.
func (DialogNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("npc_text"),
		field.JSON("responses", []DialogResponse{}).Optional(),
		field.Bool("is_entry").Default(false),
		field.String("entry_condition").Optional(),
		field.JSON("on_enter_effects", []int{}).Optional(),
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