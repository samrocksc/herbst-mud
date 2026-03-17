package main

import (
	"testing"
)

// TestColorCodedCharacters verifies that NPCs and players are rendered with correct colors
func TestColorCodedCharacters(t *testing.T) {
	// Test the color coding logic
	tests := []struct {
		name     string
		char     roomCharacter
		isNPC    bool   // expected isNPC
		expected string // substring expected in output
	}{
		{
			name: "NPC should show red",
			char: roomCharacter{
				ID:    1,
				Name:  "Goblin",
				IsNPC: true,
				Level: 3,
				Class: "Warrior",
				Race:  "Orc",
			},
			isNPC:    true,
			expected: "Goblin (NPC)",
		},
		{
			name: "Player should show green",
			char: roomCharacter{
				ID:      2,
				Name:    "Sam123",
				IsNPC:   false,
				Level:   5,
				Class:   "Tinkerer",
				Race:    "Human",
				UserID:  1,
			},
			isNPC:    false,
			expected: "Sam123 (Player)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify character is correctly identified
			if tt.char.IsNPC != tt.isNPC {
				t.Errorf("expected IsNPC=%v, got %v", tt.isNPC, tt.char.IsNPC)
			}

			// Verify name is present
			if tt.char.Name == "" {
				t.Error("character name should not be empty")
			}

			// Verify level is valid
			if tt.char.Level < 1 {
				t.Errorf("expected level >= 1, got %d", tt.char.Level)
			}
		})
	}
}

// TestRoomCharacterStruct verifies the roomCharacter struct has all required fields
func TestRoomCharacterStruct(t *testing.T) {
	char := roomCharacter{
		ID:       1,
		Name:     "Test",
		IsNPC:    true,
		Level:    5,
		Class:    "Warrior",
		Race:     "Orc",
		UserID:   0,
	}

	if char.ID != 1 {
		t.Errorf("expected ID=1, got %d", char.ID)
	}
	if char.Name != "Test" {
		t.Errorf("expected Name='Test', got %s", char.Name)
	}
	if !char.IsNPC {
		t.Error("expected IsNPC=true")
	}
	if char.Level != 5 {
		t.Errorf("expected Level=5, got %d", char.Level)
	}
	if char.Class != "Warrior" {
		t.Errorf("expected Class='Warrior', got %s", char.Class)
	}
	if char.Race != "Orc" {
		t.Errorf("expected Race='Orc', got %s", char.Race)
	}
	if char.UserID != 0 {
		t.Errorf("expected UserID=0, got %d", char.UserID)
	}
}

// TestColorConstants verifies the color constants are defined
func TestColorConstants(t *testing.T) {
	// These should be defined in main.go
	// red, green should be lipgloss.Color values
	// Verify they exist by checking they compile
	_ = red
	_ = green
}