package main

import (
	"time"
)

// DialogResponse represents a player response option in a dialog node.
type DialogResponse struct {
	Label         string `json:"label"`
	NextNodeID    string `json:"next_node_id"`
	Condition     string `json:"condition,omitempty"`
	QuestOfferID  string `json:"quest_offer_id,omitempty"`
	DeclineNodeID string `json:"decline_node_id,omitempty"`
	Effects       []int  `json:"effects,omitempty"`
}

// DialogNode represents a single node in an NPC's dialog tree.
type DialogNode struct {
	ID             string           `json:"id"`
	NPCTemplateID  string           `json:"npc_template_id"`
	NPCText        string           `json:"npc_text"`
	Responses      []DialogResponse `json:"responses"`
	IsEntry        bool             `json:"is_entry"`
	EntryCondition string           `json:"entry_condition,omitempty"`
	OnEnterEffects []int            `json:"on_enter_effects,omitempty"`
}

// ConversationState tracks a player's active dialog with an NPC.
type ConversationState struct {
	NPCTemplateID string                 `json:"npc_template_id"`
	NPCName       string                 `json:"npc_name"`
	CurrentNodeID string                 `json:"current_node_id"`
	Nodes         map[string]*DialogNode `json:"nodes"`
	StartedAt     time.Time              `json:"started_at"`
}