package main

import (
	"testing"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/rooms"
	"github.com/sam/makeathing/internal/combat"
)

func TestCharacterCreation(t *testing.T) {
	char := &characters.Character{
		Name:  "TestCharacter",
		Race:  characters.Human,
		Class: characters.Warrior,
		Stats: characters.Stats{
			Strength:     15,
			Intelligence: 10,
			Dexterity:    12,
		},
		Health: 100,
		Mana:   50,
	}

	if char.Name != "TestCharacter" {
		t.Errorf("Expected name TestCharacter, got %s", char.Name)
	}

	if char.Race != characters.Human {
		t.Errorf("Expected race Human, got %s", char.Race)
	}
}

func TestRoomCreation(t *testing.T) {
	room := &rooms.Room{
		ID:          "test_room",
		Description: "A test room",
		Exits: map[rooms.Direction]string{
			rooms.North: "north_room",
			rooms.South: "south_room",
		},
	}

	if room.ID != "test_room" {
		t.Errorf("Expected ID test_room, got %s", room.ID)
	}

	if len(room.Exits) != 2 {
		t.Errorf("Expected 2 exits, got %d", len(room.Exits))
	}
}

func TestCombatHitCalculation(t *testing.T) {
	// Test low dexterity
	hits := combat.CalculateHitsPerRound(5)
	if hits != 1 {
		t.Errorf("Expected 1 hit for low dexterity, got %d", hits)
	}

	// Test high dexterity
	hits = combat.CalculateHitsPerRound(20)
	if hits != 3 {
		t.Errorf("Expected 3 hits for high dexterity, got %d", hits)
	}

	// Test mid dexterity
	hits = combat.CalculateHitsPerRound(12)
	// Should be 1 + (12-8)/4 = 1 + 1 = 2
	if hits != 2 {
		t.Errorf("Expected 2 hits for mid dexterity, got %d", hits)
	}
}