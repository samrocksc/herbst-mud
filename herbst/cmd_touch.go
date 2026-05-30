package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ============================================================
// TOUCH COMMANDS — interact with objects via touch
// ============================================================

// handleTouchCommand handles the "touch" command for objects
func (m *model) handleTouchCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to touch objects.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: touch <object_name>\nType 'look' to see available objects.", "error")
		return
	}

	objectName := strings.Join(parts[1:], " ")

	// Get triggers for the current room
	roomTriggers, err := m.getTriggersForRoom(m.currentRoom)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching room triggers: %v", err), "error")
		return
	}

	// Check if there's a matching trigger for this object
	var matchedTrigger *triggerView
	for _, t := range roomTriggers {
		if fuzzyWordMatch(t.Name, objectName) {
			matchedTrigger = t
			break
		}
	}

	if matchedTrigger == nil {
		m.AppendMessage(fmt.Sprintf("No 'touch' interaction found for '%s'.", objectName), "error")
		return
	}

	// Process based on trigger type
	switch matchedTrigger.TriggerType {
	case "touch":
		m.handleTouchTrigger(matchedTrigger)
	default:
		m.AppendMessage(fmt.Sprintf("Trigger type '%s' not yet implemented.", matchedTrigger.TriggerType), "error")
	}
}

// handleTouchTrigger processes a touch trigger based on target type
func (m *model) handleTouchTrigger(trigger *triggerView) {
	switch trigger.TargetType {
	case "dialog_node":
		// Note: DialogNode uses string IDs, TargetID is int
		// For now, we just show a message - proper integration would require
		// updating Trigger schema to handle polymorphic references
		m.AppendMessage(fmt.Sprintf("Touch dialog node %d - full dialog support pending schema update", trigger.TargetID), "info")
	case "effect":
		m.AppendMessage("Effect trigger applied! (Effect execution requires ActiveEffect API)", "info")
	default:
		m.AppendMessage(fmt.Sprintf("Target type '%s' not yet implemented.", trigger.TargetType), "error")
	}
}

// showDialogNode displays a dialog node (for character creation, NPC interaction, etc.)
func (m *model) showDialogNode(dialogNodeID string) {
	url := fmt.Sprintf("%s/api/dialog-nodes/%s", RESTAPIBase, dialogNodeID)
	resp, err := httpGet(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching dialog node: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Could not load dialog node.", "error")
		return
	}

	var node struct {
		ID             string `json:"id"`
		NPCTemplateID  string `json:"npc_template_id"`
		NPCText        string `json:"npc_text"`
		Responses      []struct {
			Text      string `json:"text"`
			NextNode  string `json:"next_node"`
			Condition string `json:"condition"`
		} `json:"responses"`
		IsEntry        bool   `json:"is_entry"`
		EntryCondition string `json:"entry_condition"`
		WorldID        string `json:"world_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		m.AppendMessage("Error reading dialog node data", "error")
		return
	}

	// Display the dialog text
	m.AppendMessage(node.NPCText, "info")

	// Show available responses
	if len(node.Responses) > 0 {
		m.AppendMessage("Available responses:", "info")
		for i, resp := range node.Responses {
			m.AppendMessage(fmt.Sprintf("  %d. %s", i+1, resp.Text), "info")
		}
	}

	// For character creation, show the dialog node ID
	m.AppendMessage(fmt.Sprintf("Dialog node ID: %s", node.ID), "info")
}

// handleTouchWrapperCommand wraps the touch command for the command registry
func (m *model) handleTouchWrapperCommand(_ *model, args []string) {
	cmd := "touch"
	if len(args) > 0 {
		cmd = fmt.Sprintf("touch %s", strings.Join(args, " "))
	}
	m.handleTouchCommand(cmd)
}
