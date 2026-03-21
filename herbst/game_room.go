package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Item colors for room display (defined in style.go for compatibility)

// formatRoomDisplay returns the full room display string
func (m *model) formatRoomDisplay() string {
	return fmt.Sprintf("[%s]\n%s\n\nExits: %s%s%s",
		lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
		m.roomDesc,
		m.formatExitsWithColor(),
		m.formatRoomItems(),
		m.formatRoomCharacters())
}

// formatExitsWithColor returns color-coded exits
func (m *model) formatExitsWithColor() string {
	if len(m.exits) == 0 {
		return lipgloss.NewStyle().Foreground(gray).Render("none")
	}

	// Sort exit directions for consistent ordering
	dirs := make([]string, 0, len(m.exits))
	for dir := range m.exits {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	var formatted []string
	for _, dir := range dirs {
		roomID := m.exits[dir]
		var exitStyle lipgloss.Style
		if m.visitedRooms[roomID] {
			exitStyle = lipgloss.NewStyle().Foreground(exitVisitedColor)
		} else if m.knownExits[dir] {
			exitStyle = lipgloss.NewStyle().Foreground(exitKnownColor)
		} else {
			m.knownExits[dir] = true
			exitStyle = lipgloss.NewStyle().Foreground(exitNewColor)
		}
		formatted = append(formatted, exitStyle.Render(dir))
	}

	return strings.Join(formatted, ", ")
}

// formatRoomItems returns a formatted string of items in the room
func (m *model) formatRoomItems() string {
	if len(m.roomItems) == 0 {
		return ""
	}

	var items []string
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue
		}

		var style lipgloss.Style
		if item.IsImmovable {
			if item.Color != "" {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color))
			} else {
				style = lipgloss.NewStyle().Foreground(itemColorGold)
			}
		} else {
			switch item.ItemType {
			case "weapon":
				style = lipgloss.NewStyle().Foreground(itemColorWeapon)
			case "armor":
				style = lipgloss.NewStyle().Foreground(itemColorArmor)
			default:
				style = lipgloss.NewStyle().Foreground(itemColorMisc)
			}
		}

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

// loadRoomItems fetches items for the current room from the API
func (m *model) loadRoomItems() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/equipment", RESTAPIBase, m.currentRoom))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var items []RoomItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return
	}

	m.roomItems = items
}

// loadRoomCharacters fetches characters in the current room from the API
func (m *model) loadRoomCharacters() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/characters", RESTAPIBase, m.currentRoom))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var characters []roomCharacter
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		return
	}

	m.roomCharacters = characters
}

// formatRoomCharacters returns a formatted string of characters in the room
func (m *model) formatRoomCharacters() string {
	if len(m.roomCharacters) == 0 {
		return ""
	}

	var npcs []string
	var players []string

	for _, char := range m.roomCharacters {
		if char.IsNPC {
			style := lipgloss.NewStyle().Foreground(red)
			npcs = append(npcs, style.Render(char.Name))
		} else {
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
