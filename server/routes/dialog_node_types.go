package routes

import (
	"herbst-server/db"
	"herbst-server/db/schema"
)

// dialogNodeView is the JSON response shape for a DialogNode entity.
type dialogNodeView struct {
	ID              string                    `json:"id"`
	NpcText         string                    `json:"npc_text"`
	Responses       []DialogResponseInput     `json:"responses"`
	IsEntry         bool                      `json:"is_entry"`
	EntryCondition  string                    `json:"entry_condition,omitempty"`
	OnEnterEffects  []int                     `json:"on_enter_effects"`
	NpcTemplateID   string                    `json:"npc_template_id"`
}

// DialogResponseInput mirrors schema.DialogResponse for JSON input/output.
type DialogResponseInput struct {
	Label         string `json:"label"`
	NextNodeID    string `json:"next_node_id"`
	Condition     string `json:"condition,omitempty"`
	QuestOfferID  string `json:"quest_offer_id,omitempty"`
	DeclineNodeID string `json:"decline_node_id,omitempty"`
	Effects       []int  `json:"effects,omitempty"`
}

// dialogNodeInput is the JSON request shape for creating/updating a DialogNode.
type dialogNodeInput struct {
	ID              *string                 `json:"id"`
	NpcText         *string                 `json:"npc_text"`
	Responses       *[]DialogResponseInput  `json:"responses"`
	IsEntry         *bool                   `json:"is_entry"`
	EntryCondition  *string                 `json:"entry_condition"`
	OnEnterEffects  *[]int                  `json:"on_enter_effects"`
	NpcTemplateID   *string                 `json:"npc_template_id"`
}

// dialogNodeToView converts a db.DialogNode (with edge loaded) to a dialogNodeView.
func dialogNodeToView(dn *db.DialogNode) dialogNodeView {
	npcTemplateID := ""
	if dn.Edges.NpcTemplate != nil {
		npcTemplateID = dn.Edges.NpcTemplate.ID
	}
	responses := responsesToView(dn.Responses)
	return dialogNodeView{
		ID:             dn.ID,
		NpcText:        dn.NpcText,
		Responses:      responses,
		IsEntry:        dn.IsEntry,
		EntryCondition: dn.EntryCondition,
		OnEnterEffects: dn.OnEnterEffects,
		NpcTemplateID:  npcTemplateID,
	}
}

// dialogNodeToViewSimple converts a db.DialogNode without loading edges.
func dialogNodeToViewSimple(dn *db.DialogNode) dialogNodeView {
	responses := responsesToView(dn.Responses)
	return dialogNodeView{
		ID:             dn.ID,
		NpcText:        dn.NpcText,
		Responses:      responses,
		IsEntry:        dn.IsEntry,
		EntryCondition: dn.EntryCondition,
		OnEnterEffects: dn.OnEnterEffects,
		NpcTemplateID:  "",
	}
}

// responsesToView converts schema.DialogResponse slice to DialogResponseInput slice.
func responsesToView(in []schema.DialogResponse) []DialogResponseInput {
	out := make([]DialogResponseInput, len(in))
	for i, r := range in {
		out[i] = DialogResponseInput{
			Label: r.Label, NextNodeID: r.NextNodeID,
			Condition: r.Condition, QuestOfferID: r.QuestOfferID,
			DeclineNodeID: r.DeclineNodeID, Effects: r.Effects,
		}
	}
	return out
}

// responsesToSchema converts DialogResponseInput slice to schema.DialogResponse slice.
func responsesToSchema(in []DialogResponseInput) []schema.DialogResponse {
	out := make([]schema.DialogResponse, len(in))
	for i, r := range in {
		out[i] = schema.DialogResponse{
			Label: r.Label, NextNodeID: r.NextNodeID,
			Condition: r.Condition, QuestOfferID: r.QuestOfferID,
			DeclineNodeID: r.DeclineNodeID, Effects: r.Effects,
		}
	}
	return out
}