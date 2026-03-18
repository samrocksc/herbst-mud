package main

import (
	"strings"
	"testing"
)

// TestLookCommandDisplay tests the look command displays room content correctly
func TestLookCommandDisplay(t *testing.T) {
	// Create a minimal model to test formatting
	m := &model{
		roomName:     "Test Room",
		roomDesc:     "This is a test room description.",
		exits:        map[string]int{"north": 2, "south": 3},
		roomItems:    []string{"rusty sword", "shield"},
		roomNPCs:     []string{"Guard Marco"},
		visitedRooms: make(map[int]bool),
		knownExits:   make(map[string]bool),
	}

	// Test formatRoomContent
	content := m.formatRoomContent()
	if !strings.Contains(content, "Items:") {
		t.Error("Expected content to contain Items")
	}
	if !strings.Contains(content, "rusty sword") {
		t.Error("Expected content to contain rusty sword")
	}
	if !strings.Contains(content, "NPCs:") {
		t.Error("Expected content to contain NPCs")
	}
	if !strings.Contains(content, "Guard Marco") {
		t.Error("Expected content to contain Guard Marco")
	}

	// Test formatExitsWithColor
	exits := m.formatExitsWithColor()
	if !strings.Contains(exits, "north") {
		t.Error("Expected exits to contain north")
	}
	if !strings.Contains(exits, "south") {
		t.Error("Expected exits to contain south")
	}
}

// TestLookAtCommand tests the "look at <target>" functionality
func TestLookAtCommand(t *testing.T) {
	m := &model{
		currentUserName:    "TestPlayer",
		characterLevel:    5,
		roomName:           "Test Room",
		characterHP:        50,
		characterMaxHP:     100,
		characterStamina:   30,
		characterMaxStamina: 50,
		characterMana:      20,
		characterMaxMana:   50,
		roomItems:          []string{"rusty sword", "gold coin"},
		roomNPCs:           []string{"Guard Marco", "Merchant Joe"},
		visitedRooms:       make(map[int]bool),
		knownExits:         make(map[string]bool),
	}

	// Test "look at me"
	m.handleLookAt("me")
	if !strings.Contains(m.message, "TestPlayer") {
		t.Error("Expected look at me to show player name")
	}

	// Test "look at rusty sword"
	m.handleLookAt("rusty")
	if !strings.Contains(m.message, "rusty sword") {
		t.Error("Expected look at rusty to find rusty sword")
	}

	// Test "look at guard"
	m.handleLookAt("guard")
	if !strings.Contains(m.message, "Guard Marco") {
		t.Error("Expected look at guard to find Guard Marco")
	}

	// Test "look at nonexistent"
	m.handleLookAt("nonexistent")
	if !strings.Contains(m.message, "don't see") {
		t.Error("Expected look at nonexistent to show error")
	}
}

// TestLookWithNoContent tests look when room has no items or NPCs
func TestLookWithNoContent(t *testing.T) {
	m := &model{
		roomName: "Empty Room",
		roomDesc: "An empty room.",
		exits:    map[string]int{},
		roomItems: []string{},
		roomNPCs:  []string{},
	}

	content := m.formatRoomContent()
	if content != "" {
		t.Error("Expected empty content for empty room")
	}

	exits := m.formatExitsWithColor()
	if exits == "" {
		t.Error("Expected exits to show none")
	}
}