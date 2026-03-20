package main

import (
	"fmt"
	"strings"
)

// ============================================================
// GAME COMMAND PROCESSING
// ============================================================

func (m *model) processCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))
	if cmd == "" {
		return
	}

	// Handle movement commands first
	if m.handleMovement(cmd) {
		return
	}

	switch cmd {
	case "help", "?":
		m.AppendMessage(`Commands:
  n/north, s/south, e/east, w/west - Move
  look/l [target] - Look around (or examine: look <target>, look at <target>)
  ctrl+p - Scroll output up (older messages)
  ctrl+n - Scroll output down (newer messages)
  exits/x - Show exits
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  inventory/i - Show your inventory
  quests/q - Show your quest log
  whoami - Show your info
  profile/p - Edit character profile
  clear/cls - Clear screen
  quit - Exit game`, "info")

	case "look", "l":
		m.handleLookCommand(cmd)

	case "exits", "x":
		m.AppendMessage(fmt.Sprintf("Exits: %s", m.formatExitsWithColor()), "info")

	case "examine", "ex", "inspect":
		m.handleExamineCommand(cmd)

	case "search", "perception":
		m.handleSearchCommand(cmd)

	case "whoami":
		m.AppendMessage(fmt.Sprintf("=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
			m.currentUserName, m.currentUserID, m.roomName,
			m.characterLevel, m.characterExperience,
			StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)), "info")

	case "profile", "p":
		m.screen = ScreenProfile
		m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
		m.menuCursor = 0
		m.AppendMessage("", "")

	case "peer":
		m.AppendMessage("Usage: peer <direction>", "error")

	case "debug":
		m.handleDebugCommand(cmd)

	case "clear", "cls":
		m.messageHistory = nil
		m.messageTypes = nil
		m.historyOffset = 0
		m.isScrolling = false
		m.inputBuffer = ""

	case "quit", "q":
		m.AppendMessage("Thanks for playing! Goodbye!", "success")
		m.inputBuffer = ""

	default:
		// take/get
		if strings.HasPrefix(cmd, "take ") || strings.HasPrefix(cmd, "get ") {
			m.handleTakeCommand(cmd)
			return
		}
		// drop
		if strings.HasPrefix(cmd, "drop ") {
			m.handleDropCommand(cmd)
			return
		}
		// inventory
		if cmd == "inventory" || cmd == "i" || cmd == "inv" {
			m.handleInventoryCommand()
			return
		}
		// quests
		if cmd == "quests" || cmd == "q" || cmd == "quest" {
			m.handleQuestsCommand(cmd)
			return
		}
		// skills
		if cmd == "skills" {
			m.handleSkillsCommand(cmd)
			return
		}
		// talents
		if cmd == "talents" {
			m.handleTalentsCommand(cmd)
			return
		}
		// skill equip
		if strings.HasPrefix(cmd, "skill ") {
			m.handleSkillEquipCommand(cmd)
			return
		}
		// talent equip/unequip/swap
		if strings.HasPrefix(cmd, "talent ") {
			m.handleTalentEquipCommand(cmd)
			return
		}
		// peer with direction
		if strings.HasPrefix(cmd, "peer ") {
			m.handlePeerCommand(cmd)
			return
		}

		m.AppendMessage(fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd), "error")
	}
}
