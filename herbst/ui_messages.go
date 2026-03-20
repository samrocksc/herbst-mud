package main

import (
	"strings"
)

// ============================================================
// MESSAGES & OUTPUT
// ============================================================

// styledMessage returns a styled message based on messageType
func (m *model) styledMessage(msg string) string {
	return styleMessage(msg, m.messageType)
}

// styleMessage returns a styled message based on message type
func styleMessage(msg string, msgType string) string {
	if msg == "" {
		return ""
	}

	switch msgType {
	case "success":
		return successStyle.Render("✓ ") + msg
	case "error":
		return errorStyle.Render("✗ ") + msg
	case "info":
		return infoStyle.Render("ℹ ") + msg
	case "damage":
		return combatDamageStyle.Render("⚔ ") + msg
	case "heal":
		return combatHealStyle.Render("♥ ") + msg
	default:
		return msg
	}
}

// AppendMessage adds a message to the history buffer
func (m *model) AppendMessage(text, msgType string) {
	m.messageHistory = append(m.messageHistory, text)
	m.messageTypes = append(m.messageTypes, msgType)
	if len(m.messageHistory) > m.maxHistory {
		m.messageHistory = m.messageHistory[len(m.messageHistory)-m.maxHistory:]
		m.messageTypes = m.messageTypes[len(m.messageTypes)-m.maxHistory:]
	}
	m.historyOffset = 0
	m.isScrolling = false
}

// buildOutputContent constructs the message history content for display
func (m *model) buildOutputContent() string {
	total := len(m.messageHistory)
	if total == 0 {
		return ""
	}

	var lines []string

	if !m.isScrolling {
		start := 0
		if total > 3 {
			start = total - 3
		}
		for i := start; i < total; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
	} else {
		for i := m.historyOffset; i < total-1; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
		lines = append(lines, "─── NEWEST ───")
		lines = append(lines, styleMessage(m.messageHistory[total-1], m.messageTypes[total-1]))
	}

	return strings.Join(lines, "\n\n")
}
