package main

import (
	"fmt"
	"strings"
)

// ============================================================
// PRESS COMMANDS — interact with objects by pressing them
// ============================================================

// handlePressCommand handles the "press" command for objects
func (m *model) handlePressCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to press objects.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: press <object_name>\nType 'look' to see available objects.", "error")
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
		m.AppendMessage(fmt.Sprintf("No 'press' interaction found for '%s'.", objectName), "error")
		return
	}

	// Process based on trigger type
	switch matchedTrigger.TriggerType {
	case "press":
		m.handlePressTrigger(matchedTrigger)
	default:
		m.AppendMessage(fmt.Sprintf("Trigger type '%s' not yet implemented.", matchedTrigger.TriggerType), "error")
	}
}

// handlePressTrigger processes a press trigger based on target type
func (m *model) handlePressTrigger(trigger *triggerView) {
	switch trigger.TargetType {
	case "effect":
		m.AppendMessage("Effect trigger applied! (Effect execution requires ActiveEffect API)", "info")
	case "recipe":
		m.AppendMessage(fmt.Sprintf("Pressing this object opens recipe menu: %d", trigger.TargetID), "info")
	default:
		m.AppendMessage(fmt.Sprintf("Target type '%s' not yet implemented.", trigger.TargetType), "error")
	}
}

// handlePressWrapperCommand wraps the press command for the command registry
func (m *model) handlePressWrapperCommand(_ *model, args []string) {
	cmd := "press"
	if len(args) > 0 {
		cmd = fmt.Sprintf("press %s", strings.Join(args, " "))
	}
	m.handlePressCommand(cmd)
}
