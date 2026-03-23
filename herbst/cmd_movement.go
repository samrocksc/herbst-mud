package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// MOVEMENT + PEER
// ============================================================

func (m *model) handleMovement(cmd string) bool {
	directionMap := map[string]string{
		"n": "north", "north": "north",
		"s": "south", "south": "south",
		"e": "east", "east": "east",
		"w": "west", "west": "west",
		"ne": "northeast", "northeast": "northeast",
		"se": "southeast", "southeast": "southeast",
		"sw": "southwest", "southwest": "southwest",
		"nw": "northwest", "northwest": "northwest",
		"u": "up", "up": "up",
		"d": "down", "down": "down",
	}

	direction, ok := directionMap[cmd]
	if !ok {
		return false
	}

	nextRoomID, ok := m.exits[direction]
	if !ok {
		m.AppendMessage("You can't go that way.", "error")
		return true
	}

	m.knownExits[direction] = true

	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error moving: %v", err), "error")
			return true
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits

		m.loadRoomItems()
		m.loadRoomCharacters()

		wasVisited := m.visitedRooms[m.currentRoom]
		m.visitedRooms[m.currentRoom] = true

		for dir := range m.exits {
			m.knownExits[dir] = true
		}

		roomDisplay := fmt.Sprintf("\n\nExits: %s%s%s",
			m.formatExitsWithColor(),
			m.formatRoomItems(),
			m.formatRoomCharacters())

		if wasVisited {
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		} else {
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		}
	}

	return true
}

func (m *model) handlePeerCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: peer <direction>\nDirections: n/north, s/south, e/east, w/west, ne/northeast, se/southeast, sw/southwest, nw/northwest, u/up, d/down", "error")
		return
	}
	direction := strings.ToLower(parts[1])

	validDirs := map[string]string{
		"n": "north", "north": "north",
		"s": "south", "south": "south",
		"e": "east", "east": "east",
		"w": "west", "west": "west",
		"ne": "northeast", "northeast": "northeast",
		"se": "southeast", "southeast": "southeast",
		"sw": "southwest", "southwest": "southwest",
		"nw": "northwest", "northwest": "northwest",
		"u": "up", "up": "up",
		"d": "down", "down": "down",
	}
	dir, ok := validDirs[direction]
	if !ok {
		m.AppendMessage("Invalid direction. Use: n, s, e, w, ne, se, sw, nw, u, d (or full direction names)", "error")
		return
	}

	nextRoomID, ok := m.exits[dir]
	if !ok {
		m.AppendMessage("You can't peer that way — there's no exit.", "error")
		return
	}

	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error looking: %v", err), "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("You peer %s...\n\n[%s]\n%s",
			dir,
			lipgloss.NewStyle().Bold(true).Foreground(blue).Render(room.Name),
			room.Description), "info")
	}
}
