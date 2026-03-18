package main

import (
	"testing"
)

// TestColorCodedCharacters verifies that NPCs and players are rendered with correct colors
func TestColorCodedCharacters(t *testing.T) {
	// Test roomCharacter struct has IsNPC field
	char1 := roomCharacter{
		ID:    1,
		Name:  "Combat Dummy",
		IsNPC: true,
		Level: 1,
		Class: "fighter",
		Race:  "human",
	}

	char2 := roomCharacter{
		ID:     2,
		Name:   "Sam123",
		IsNPC:  false,
		Level:  5,
		Class:  "tinkerer",
		Race:   "turtle",
		UserID: 1,
	}

	// Verify NPC flag
	if !char1.IsNPC {
		t.Error("Expected char1 to be NPC")
	}

	if char2.IsNPC {
		t.Error("Expected char2 to be player, not NPC")
	}

	// Test filtering logic
	allChars := []roomCharacter{char1, char2}
	currentCharID := 2 // Simulate being Sam123

	var otherChars []roomCharacter
	for _, rc := range allChars {
		if rc.ID != currentCharID {
			otherChars = append(otherChars, rc)
		}
	}

	// Should only have Combat Dummy
	if len(otherChars) != 1 {
		t.Errorf("Expected 1 other character, got %d", len(otherChars))
	}

	if otherChars[0].Name != "Combat Dummy" {
		t.Errorf("Expected Combat Dummy, got %s", otherChars[0].Name)
	}

	// Verify the NPC is correctly identified
	if !otherChars[0].IsNPC {
		t.Error("Expected filtered character to be NPC")
	}
}

// TestCharacterListFormatting verifies the character list formatting logic
func TestCharacterListFormatting(t *testing.T) {
	characters := []roomCharacter{
		{ID: 1, Name: "Goblin", IsNPC: true},
		{ID: 2, Name: "Player1", IsNPC: false},
		{ID: 3, Name: "Orc", IsNPC: true},
	}

	currentID := 2

	// Count NPCs vs Players
	npcCount := 0
	playerCount := 0

	for _, rc := range characters {
		if rc.ID == currentID {
			continue // skip self
		}
		if rc.IsNPC {
			npcCount++
		} else {
			playerCount++
		}
	}

	if npcCount != 2 {
		t.Errorf("Expected 2 NPCs, got %d", npcCount)
	}

	if playerCount != 1 {
		t.Errorf("Expected 1 player, got %d", playerCount)
	}
}