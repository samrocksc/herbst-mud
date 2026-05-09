package main

import (
	"fmt"
	"strconv"
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

	// Handle look/l commands with or without targets
	if cmd == "look" || cmd == "l" || strings.HasPrefix(cmd, "look ") || strings.HasPrefix(cmd, "l ") {
		m.handleLookCommand(cmd)
		return
	}

	// Handle debug commands (debug, debug on, debug off, debug info)
	if cmd == "debug" || strings.HasPrefix(cmd, "debug ") {
		m.handleDebugCommand(cmd)
		return
	}

	// Handle dialog choices when in a conversation
	if m.conversation != nil {
		if cmd == "0" || cmd == "leave" || cmd == "bye" {
			m.endConversation()
			return
		}
		choice, err := strconv.Atoi(cmd)
		if err == nil && choice >= 1 && choice <= 9 {
			m.handleDialogChoice(choice)
			return
		}
		// Fall through to normal commands if not a valid choice
	}

	// Try command registry first
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		if m.commands.Execute(m, parts[0], parts) {
			return
		}
	}

	switch cmd {
	case "exits", "x":
		m.AppendMessage(fmt.Sprintf("Exits: %s", m.formatExitsWithColor()), "info")

	case "examine", "ex", "inspect":
		m.handleExamineCommand(cmd)

	case "search", "perception":
		m.handleSearchCommand(cmd)


	case "peer":
		m.AppendMessage("Usage: peer <direction>", "error")

	default:
		// attack/kill/a/fight - handle both bare commands and commands with targets
		if cmd == "attack" || cmd == "kill" || cmd == "a" || cmd == "fight" ||
			strings.HasPrefix(cmd, "attack ") || strings.HasPrefix(cmd, "kill ") || strings.HasPrefix(cmd, "a ") || strings.HasPrefix(cmd, "fight ") {
			m.handleAttackCommand(cmd)
			return
		}
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
		// loot
		if cmd == "loot" || strings.HasPrefix(cmd, "loot ") {
			m.handleLootCommand(cmd)
			return
		}
		// quests
		if cmd == "quests" || cmd == "q" || cmd == "quest" {
			m.handleQuestsCommand(cmd)
			return
		}
		if strings.HasPrefix(cmd, "quest accept ") || strings.HasPrefix(cmd, "quest abandon ") {
			m.handleQuestSubcommand(cmd)
			return
		}
		// skills (classless combat skills)
		if cmd == "skills" {
			if m.combatSkills == nil {
				m.combatSkills = &CombatSkillState{}
				m.initCombatSkillState()
			}
			m.showEquippedAbilities()
			return
		}
		// skill (classless skills) - equip, swap, show
		if strings.HasPrefix(cmd, "skill ") {
			if m.combatSkills == nil {
				m.combatSkills = &CombatSkillState{}
				m.initCombatSkillState()
			}
			m.handleAbilityCommand(cmd[6:])
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