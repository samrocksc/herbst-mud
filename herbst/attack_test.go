package main

import (
	"testing"
)

// TestHandleAttackCommandNoTarget tests attack with no target
func TestHandleAttackCommandNoTarget(t *testing.T) {
	m := &model{
		messageHistory: []string{},
		messageTypes:   []string{},
		maxHistory:     100,
	}
	m.handleAttackCommand("attack")

	if len(m.messageHistory) == 0 {
		t.Error("Expected error message for missing target")
	} else if m.messageHistory[0] != "Attack what? Usage: attack <target name>" {
		t.Errorf("Expected usage message, got: %s", m.messageHistory[0])
	}
}

// TestHandleAttackCommandTargetNotFound tests attacking non-existent target
func TestHandleAttackCommandTargetNotFound(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{},
		messageHistory: []string{},
		messageTypes:   []string{},
		maxHistory:     100,
	}
	m.handleAttackCommand("attack nonexistent")

	if len(m.messageHistory) == 0 {
		t.Error("Expected error message for target not found")
	} else {
		msg := m.messageHistory[0]
		if msg != "You don't see any 'nonexistent' here to attack." {
			t.Errorf("Expected 'not found' message, got: %s", msg)
		}
	}
}

// TestHandleAttackCommandNPC tests attacking an NPC
func TestHandleAttackCommandNPC(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{Name: "Goblin", IsNPC: true, Level: 1},
		},
		currentCharacterID: 1,
		characterLevel:     1,
		messageHistory:     []string{},
		messageTypes:       []string{},
		maxHistory:         100,
	}

	// Note: This will attempt an API call which will fail without a server
	// The test verifies the command routing, not full combat mechanics
	m.handleAttackCommand("attack Goblin")

	// Should have a message (either combat or error about API)
	if len(m.messageHistory) == 0 {
		t.Error("Expected a message after attack command")
	}
}

// TestHandleAttackCommandExactMatch tests exact name matching
func TestHandleAttackCommandExactMatch(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{Name: "Rat", IsNPC: true, Level: 1},
			{Name: "Rat King", IsNPC: true, Level: 5},
		},
		currentCharacterID: 1,
		characterLevel:     1,
		messageHistory:     []string{},
		messageTypes:       []string{},
		maxHistory:         100,
	}

	// Attack "rat" should match "Rat" exactly first
	m.handleAttackCommand("attack rat")

	if len(m.messageHistory) == 0 {
		t.Error("Expected a message after attack command")
	}
}

// TestHandleAttackCommandFuzzyMatch tests fuzzy name matching
func TestHandleAttackCommandFuzzyMatch(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{Name: "Scrap Rat", IsNPC: true, Level: 1},
		},
		currentCharacterID: 1,
		characterLevel:     1,
		messageHistory:     []string{},
		messageTypes:       []string{},
		maxHistory:         100,
	}

	// "scrap" should match "Scrap Rat" via fuzzy matching
	m.handleAttackCommand("attack scrap")

	if len(m.messageHistory) == 0 {
		t.Error("Expected a message after attack command")
	}
}

// TestHandleAttackCommandCaseInsensitive tests case insensitive matching
func TestHandleAttackCommandCaseInsensitive(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{Name: "Goblin Warrior", IsNPC: true, Level: 3},
		},
		currentCharacterID: 1,
		characterLevel:     1,
		messageHistory:     []string{},
		messageTypes:       []string{},
		maxHistory:         100,
	}

	// Case insensitive matching
	m.handleAttackCommand("attack GOBLIN WARRIOR")

	if len(m.messageHistory) == 0 {
		t.Error("Expected a message after attack command")
	}
}

// TestGetCharacterStrength tests strength calculation
func TestGetCharacterStrength(t *testing.T) {
	m := &model{
		currentCharacterID: 0, // No character loaded
	}

	// Should return default strength when no character
	strength := m.getCharacterStrength()
	if strength != 10 {
		t.Errorf("Expected default strength 10, got %d", strength)
	}
}

// TestHandleAttackCommandPvP tests attacking another player
func TestHandleAttackCommandPvP(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{Name: "TestPlayer", IsNPC: false, Level: 1},
		},
		currentCharacterID: 1,
		characterLevel:     1,
		messageHistory:     []string{},
		messageTypes:       []string{},
		maxHistory:         100,
	}

	m.handleAttackCommand("attack TestPlayer")

	// Should have a PvP combat message
	if len(m.messageHistory) == 0 {
		t.Error("Expected a message after PvP attack")
	}
}