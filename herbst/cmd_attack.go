package main

import (
	"fmt"
	"strings"
)

// handleAttackCommand handles the attack/kill/a/fight command
func (m *model) handleAttackCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Attack what? Usage: attack <target name>", "error")
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(strings.TrimSpace(target))

	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Attack command received, target: '%s'", target), "info")
	}

	// Load room characters if needed
	m.loadRoomCharacters()

	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Room has %d characters", len(m.roomCharacters)), "info")
		for _, c := range m.roomCharacters {
			m.AppendMessage(fmt.Sprintf("[DEBUG]   - %s (ID: %d, HP: %d/%d)", c.Name, c.ID, c.HP, c.MaxHP), "info")
		}
	}

	// Find target in room
	var foundChar *RoomCharacter
	for i := range m.roomCharacters {
		char := &m.roomCharacters[i]
		charNameLower := strings.ToLower(char.Name)
		if fuzzyWordMatch(char.Name, target) || strings.Contains(charNameLower, target) || charNameLower == target {
			foundChar = char
			break
		}
	}

	if foundChar == nil {
		m.AppendMessage(fmt.Sprintf("You don't see any '%s' here to attack.", target), "error")
		return
	}

	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Found target: %s (ID: %d)", foundChar.Name, foundChar.ID), "info")
	}

	// Fetch fresh HP from server for combat
	hp, maxHP, _ := getCharacterCombatStatus(foundChar.ID)
	foundChar.HP = hp
	foundChar.MaxHP = maxHP

	// Check if target is alive
	if hp <= 0 {
		m.AppendMessage(fmt.Sprintf("%s is already defeated. They need time to recover.", foundChar.Name), "error")
		return
	}

	// Start combat mode
	m.AppendMessage(fmt.Sprintf("⚔ Starting combat with %s...", foundChar.Name), "combat")
	m.startCombat(foundChar)
}