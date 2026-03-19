package main

import (
	"testing"
)

// TestStyleMessageDamage tests the damage message style (red)
func TestStyleMessageDamage(t *testing.T) {
	msg := styleMessage("You take 25 damage!", "damage")
	
	// Should contain the damage prefix and message
	if msg == "" {
		t.Error("Expected non-empty styled message")
	}
	
	// Should contain the sword emoji prefix
	if len(msg) < len("⚔ You take 25 damage!") {
		t.Errorf("Expected message to have damage prefix, got: %s", msg)
	}
}

// TestStyleMessageHeal tests the healing message style (green)
func TestStyleMessageHeal(t *testing.T) {
	msg := styleMessage("You recover 15 HP!", "heal")
	
	// Should contain the heal prefix and message
	if msg == "" {
		t.Error("Expected non-empty styled message")
	}
	
	// Should contain the heart emoji prefix
	if len(msg) < len("♥ You recover 15 HP!") {
		t.Errorf("Expected message to have heal prefix, got: %s", msg)
	}
}

// TestStyleMessageCombatDamageToTarget tests damage messages to target
func TestStyleMessageCombatDamageToTarget(t *testing.T) {
	msg := styleMessage("Goblin hits you for 10 damage", "damage")
	
	if msg == "" {
		t.Error("Expected non-empty styled message")
	}
	
	// Verify it contains damage styling indicator (sword prefix)
	expected := "⚔ Goblin hits you for 10 damage"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}
}

// TestStyleMessageCombatHealFromPotion tests healing from items/spells
func TestStyleMessageCombatHealFromPotion(t *testing.T) {
	msg := styleMessage("Potion heals you for 20 HP", "heal")
	
	if msg == "" {
		t.Error("Expected non-empty styled message")
	}
	
	// Verify it contains healing styling indicator (heart prefix)
	expected := "♥ Potion heals you for 20 HP"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}
}

// TestStyleMessageEmpty tests that empty messages return empty
func TestStyleMessageEmpty(t *testing.T) {
	msg := styleMessage("", "damage")
	if msg != "" {
		t.Errorf("Expected empty string for empty message, got: %s", msg)
	}
}

// TestStyleMessageUnknownType tests fallback for unknown message types
func TestStyleMessageUnknownType(t *testing.T) {
	msg := styleMessage("Some random message", "unknown_type")
	if msg != "Some random message" {
		t.Errorf("Expected plain message for unknown type, got: %s", msg)
	}
}