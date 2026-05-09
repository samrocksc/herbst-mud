package main

import (
	"fmt"
	"strconv"
	"strings"
)

// handleQuestSubcommand handles "quest accept <id>" and "quest abandon <id>".
func (m *model) handleQuestSubcommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 3 {
		m.AppendMessage("Usage: quest accept <id> | quest abandon <id>", "error")
		return
	}
	subcmd := parts[1]
	questID, err := strconv.Atoi(parts[2])
	if err != nil {
		m.AppendMessage("Invalid quest ID. Usage: quest accept <number>", "error")
		return
	}
	switch subcmd {
	case "accept":
		m.handleQuestAccept(questID)
	case "abandon":
		m.handleQuestAbandon(questID)
	default:
		m.AppendMessage("Unknown quest command. Use: quest accept <id> or quest abandon <id>", "error")
	}
}

func (m *model) handleQuestAccept(questID int) {
	qp, err := m.questService.AcceptQuest(m.currentCharacterID, questID)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Failed to accept quest: %v", err), "error")
		return
	}
	name := qp.QuestName
	if name == "" {
		name = fmt.Sprintf("Quest #%d", questID)
	}
	m.AppendMessage(fmt.Sprintf("Quest accepted: %s", name), "success")
}

func (m *model) handleQuestAbandon(questID int) {
	err := m.questService.AbandonQuest(m.currentCharacterID, questID)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Failed to abandon quest: %v", err), "error")
		return
	}
	m.AppendMessage(fmt.Sprintf("Quest #%d abandoned.", questID), "info")
}