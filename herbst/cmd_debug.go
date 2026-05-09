package main

import (
	"fmt"
	"sort"
	"strings"
)

// HandleDebug handles the debug command
func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		if m.debugMode {
			m.AppendMessage("Debug mode: ON (Room ID visible in status bar)", "info")
		} else {
			m.AppendMessage("Debug mode: OFF\nUsage: debug on | debug off | debug info", "info")
		}
		return
	}

	switch parts[1] {
	case "on", "true", "1", "yes":
		if !m.isTest && !m.debugMode {
			m.AppendMessage("Debug mode is only available to test characters.", "error")
			return
		}
		m.debugMode = true
		m.AppendMessage("Debug mode: ON (Room ID will show in status bar)", "success")
	case "off", "false", "0", "no":
		m.debugMode = false
		m.AppendMessage("Debug mode: OFF", "info")
	case "info":
		if !m.isTest && !m.debugMode {
			m.AppendMessage("Debug info is only available to test characters. Type 'debug on' first.", "error")
			return
		}
		m.printDebugInfo()
	default:
		m.AppendMessage("Usage: debug on | debug off | debug info", "error")
	}
}

// printDebugInfo dumps a full copyable debug block
func (m *model) printDebugInfo() {
	var b strings.Builder

	b.WriteString("═══ DEBUG INFO ═══\n")
	b.WriteString(fmt.Sprintf("Character: %s (#%d) | Race: %s | Class: adventurer | Level: %d\n",
		m.currentCharacterName, m.currentCharacterID, m.characterRace, m.characterLevel))
	b.WriteString(fmt.Sprintf("HP: %d/%d | Stamina: %d/%d | Mana: %d/%d\n",
		m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana))
	b.WriteString(fmt.Sprintf("STR: %d | DEX: %d | CON: %d | INT: %d | WIS: %d\n",
		0, 0, 0, 0, 0))
	b.WriteString(fmt.Sprintf("XP: %d | Room: %d | Respawn: %d | Admin: %v | Test: %v\n",
		m.characterExperience, m.currentRoom, m.respawnRoom, false, m.isTest))

	b.WriteString("\n═══ ROOM ═══\n")
	b.WriteString(fmt.Sprintf("Room #%d: %s\n", m.currentRoom, m.roomName))
	if m.roomDesc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", m.roomDesc))
	}

	// Exits
	if len(m.exits) > 0 {
		dirs := make([]string, 0, len(m.exits))
		for dir := range m.exits {
			dirs = append(dirs, dir)
		}
		sort.Strings(dirs)
		exitParts := make([]string, 0, len(dirs))
		for _, dir := range dirs {
			exitParts = append(exitParts, fmt.Sprintf("%s→#%d", dir, m.exits[dir]))
		}
		b.WriteString(fmt.Sprintf("Exits: %s\n", strings.Join(exitParts, ", ")))
	} else {
		b.WriteString("Exits: (none)\n")
	}

	// Items
	if len(m.roomItems) > 0 {
		itemNames := make([]string, 0, len(m.roomItems))
		for _, item := range m.roomItems {
			itemNames = append(itemNames, fmt.Sprintf("%s (#%d)", item.Name, item.ID))
		}
		b.WriteString(fmt.Sprintf("Items: %s\n", strings.Join(itemNames, ", ")))
	} else {
		b.WriteString("Items: (none)\n")
	}

	// Characters
	if len(m.roomCharacters) > 0 {
		charNames := make([]string, 0, len(m.roomCharacters))
		for _, c := range m.roomCharacters {
			npcTag := ""
			if c.IsNPC {
				npcTag = " [NPC]"
			}
			charNames = append(charNames, fmt.Sprintf("%s (#%d L%d%s)", c.Name, c.ID, c.Level, npcTag))
		}
		b.WriteString(fmt.Sprintf("Characters: %s\n", strings.Join(charNames, ", ")))
	} else {
		b.WriteString("Characters: (none)\n")
	}

	b.WriteString("═══ END ═══")

	m.AppendMessage(b.String(), "info")
}