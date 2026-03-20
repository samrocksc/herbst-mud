package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// MESSAGE HANDLING - Message types and helpers
// ============================================================

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
		// Combat damage: red text with ⚔ prefix
		return combatDamageStyle.Render("⚔ ") + msg
	case "heal":
		// Combat healing: green text with ♥ prefix
		return combatHealStyle.Render("♥ ") + msg
	default:
		return msg
	}
}

// AppendMessage adds a message to the history buffer (UI-21 message history)
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
		// Show last 3 messages (or fewer)
		start := 0
		if total > 3 {
			start = total - 3
		}
		for i := start; i < total; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
	} else {
		// Scrolled up: show from historyOffset toward newest (excluding latest, which is pinned)
		for i := m.historyOffset; i < total-1; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
		// Latest pinned at bottom with ─── NEWEST ─── marker
		lines = append(lines, "─── NEWEST ───")
		lines = append(lines, styleMessage(m.messageHistory[total-1], m.messageTypes[total-1]))
	}

	return strings.Join(lines, "\n\n")
}

// formatExitsWithColor returns color-coded exits
func (m *model) formatExitsWithColor() string {
	if len(m.exits) == 0 {
		return lipgloss.NewStyle().Foreground(gray).Render("none")
	}

	var dirs []string
	for dir, roomID := range m.exits {
		var exitStyle lipgloss.Style
		if m.visitedRooms[roomID] {
			// Green = visited
			exitStyle = lipgloss.NewStyle().Foreground(exitVisitedColor)
		} else if m.knownExits[dir] {
			// Yellow = known but not visited
			exitStyle = lipgloss.NewStyle().Foreground(exitKnownColor)
		} else {
			// White = new
			m.knownExits[dir] = true
			exitStyle = lipgloss.NewStyle().Foreground(exitNewColor)
		}
		dirs = append(dirs, exitStyle.Render(dir))
	}

	return strings.Join(dirs, ", ")
}

// formatRoomItems returns a formatted string of items in the room
func (m *model) formatRoomItems() string {
	if len(m.roomItems) == 0 {
		return ""
	}

	var items []string
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue // Skip invisible items
		}

		var style lipgloss.Style
		if item.IsImmovable {
			// Immobile items get gold color by default or custom color
			if item.Color != "" {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color))
			} else {
				style = lipgloss.NewStyle().Foreground(itemColorGold)
			}
		} else {
			// Regular items get color based on type
			switch item.ItemType {
			case "weapon":
				style = lipgloss.NewStyle().Foreground(itemColorWeapon)
			case "armor":
				style = lipgloss.NewStyle().Foreground(itemColorArmor)
			default:
				style = lipgloss.NewStyle().Foreground(itemColorMisc)
			}
		}

		// Add special marker for immovable items
		if item.IsImmovable {
			items = append(items, style.Render("⬥ "+item.Name))
		} else {
			items = append(items, style.Render(item.Name))
		}
	}

	if len(items) == 0 {
		return ""
	}
	return "\n\nYou see: " + strings.Join(items, ", ")
}

// formatRoomCharacters returns a formatted string of characters (NPCs and players) in the room
func (m *model) formatRoomCharacters() string {
	if len(m.roomCharacters) == 0 {
		return ""
	}

	var npcs []string
	var players []string

	for _, char := range m.roomCharacters {
		if char.IsNPC {
			// NPCs in red
			style := lipgloss.NewStyle().Foreground(red)
			npcs = append(npcs, style.Render(char.Name))
		} else {
			// Players in green
			style := lipgloss.NewStyle().Foreground(green)
			players = append(players, style.Render(char.Name))
		}
	}

	var parts []string
	if len(npcs) > 0 {
		parts = append(parts, "NPCs: "+strings.Join(npcs, ", "))
	}
	if len(players) > 0 {
		parts = append(parts, "Players: "+strings.Join(players, ", "))
	}

	if len(parts) == 0 {
		return ""
	}
	return "\n\n" + strings.Join(parts, " | ")
}
