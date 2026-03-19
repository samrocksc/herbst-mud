package main

import (
	"testing"
)

// TestHandleAttackCommand tests the attack command handler
func TestHandleAttackCommand(t *testing.T) {
	m := model{
		roomCharacters: []roomCharacter{
			{Name: "Combat Dummy", IsNPC: true, ShortName: "dummy"},
			{Name: "Evil Rat", IsNPC: true, ShortName: "rat"},
			{Name: "TestPlayer", IsNPC: false},
		},
		currentCharacterName: "TestPlayer",
	}

	// Test 1: No target specified
	m.handleAttackCommand("attack")
	if m.message != "Attack what? Usage: attack <target name>" {
		t.Errorf("Expected no target error message, got: %s", m.message)
	}
	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got: %s", m.messageType)
	}

	// Test 2: Target not found
	m.handleAttackCommand("attack nonexistent")
	expectedMsg := "You don't see any 'nonexistent' here to attack."
	if m.message != expectedMsg {
		t.Errorf("Expected '%s', got: %s", expectedMsg, m.message)
	}
	if m.messageType != "error" {
		t.Errorf("Expected messageType 'error', got: %s", m.messageType)
	}

	// Test 3: Exact match on NPC
	m.handleAttackCommand("attack Combat Dummy")
	if m.messageType != "combat" {
		t.Errorf("Expected messageType 'combat', got: %s", m.messageType)
	}

	// Test 4: Partial match (short name "dummy" matches "Combat Dummy")
	m.handleAttackCommand("attack dummy")
	if m.messageType != "combat" {
		t.Errorf("Expected messageType 'combat' for partial match, got: %s", m.messageType)
	}

	// Test 5: Case-insensitive match
	m.handleAttackCommand("attack EVIL RAT")
	if m.messageType != "combat" {
		t.Errorf("Expected messageType 'combat' for case-insensitive match, got: %s", m.messageType)
	}

	// Test 6: Attack player (PvP)
	m.handleAttackCommand("attack TestPlayer")
	if m.messageType != "combat" {
		t.Errorf("Expected messageType 'combat' for PvP, got: %s", m.messageType)
	}
}