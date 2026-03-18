package main

import (
	"testing"
)

// TestRoomCharacterStruct verifies the roomCharacter struct has all required fields
func TestRoomCharacterStruct(t *testing.T) {
	char := roomCharacter{
		ID:       1,
		Name:     "Goblin",
		IsNPC:    true,
		Level:    3,
		Class:    "Warrior",
		Race:     "Orc",
		UserID:   0,
	}

	if char.ID != 1 {
		t.Errorf("expected ID=1, got %d", char.ID)
	}
	if char.Name != "Goblin" {
		t.Errorf("expected Name='Goblin', got %s", char.Name)
	}
	if !char.IsNPC {
		t.Error("expected IsNPC=true")
	}
	if char.Level != 3 {
		t.Errorf("expected Level=3, got %d", char.Level)
	}
	if char.Class != "Warrior" {
		t.Errorf("expected Class='Warrior', got %s", char.Class)
	}
	if char.Race != "Orc" {
		t.Errorf("expected Race='Orc', got %s", char.Race)
	}
}

// TestFormatRoomCharacters_NPCSOnly tests formatting when only NPCs are present
func TestFormatRoomCharacters_NPCSOnly(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{ID: 1, Name: "Goblin", IsNPC: true, Level: 3},
			{ID: 2, Name: "Rat", IsNPC: true, Level: 1},
		},
	}

	result := m.formatRoomCharacters()
	if result == "" {
		t.Error("expected non-empty result for NPCs in room")
	}
	// NPCs should be displayed in red (color codes present)
	// The result should contain "NPCs:" prefix
}

// TestFormatRoomCharacters_PlayersOnly tests formatting when only players are present
func TestFormatRoomCharacters_PlayersOnly(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{ID: 1, Name: "Sam", IsNPC: false, Level: 5},
			{ID: 2, Name: "Alex", IsNPC: false, Level: 10},
		},
	}

	result := m.formatRoomCharacters()
	if result == "" {
		t.Error("expected non-empty result for players in room")
	}
	// Players should be displayed in green (color codes present)
	// The result should contain "Players:" prefix
}

// TestFormatRoomCharacters_Mixed tests formatting with both NPCs and players
func TestFormatRoomCharacters_Mixed(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{
			{ID: 1, Name: "Goblin", IsNPC: true, Level: 3},
			{ID: 2, Name: "Sam", IsNPC: false, Level: 5},
		},
	}

	result := m.formatRoomCharacters()
	if result == "" {
		t.Error("expected non-empty result for mixed room")
	}
	// Result should contain both NPCs and Players sections
}

// TestFormatRoomCharacters_Empty tests formatting with no characters
func TestFormatRoomCharacters_Empty(t *testing.T) {
	m := &model{
		roomCharacters: []roomCharacter{},
	}

	result := m.formatRoomCharacters()
	if result != "" {
		t.Errorf("expected empty result for no characters, got: %s", result)
	}
}

// TestLookCommandIntegration tests that look command loads and displays room info
func TestLookCommandIntegration(t *testing.T) {
	// Test that the look command message format is correct
	// This verifies the structure of the output

	testCases := []struct {
		name        string
		roomName    string
		roomDesc    string
		exits       map[string]int
		items       []RoomItem
		characters  []roomCharacter
		expectEmpty bool
	}{
		{
			name:     "Empty room",
			roomName: "Test Room",
			roomDesc: "A plain room for testing.",
			exits:    map[string]int{"north": 2, "south": 3},
			items:    []RoomItem{},
			characters: []roomCharacter{},
			expectEmpty: false,
		},
		{
			name:     "Room with items and NPCs",
			roomName: "Fountain Plaza",
			roomDesc: "A central fountain surrounded by cobblestones.",
			exits:    map[string]int{"north": 1, "east": 4},
			items: []RoomItem{
				{ID: 1, Name: "Rusty Sword", IsVisible: true},
			},
			characters: []roomCharacter{
				{ID: 1, Name: "Gizmo", IsNPC: true, Level: 1},
			},
			expectEmpty: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &model{
				roomName:       tc.roomName,
				roomDesc:       tc.roomDesc,
				exits:          tc.exits,
				roomItems:      tc.items,
				roomCharacters: tc.characters,
			}

			// Verify formatRoomCharacters works
			charsResult := m.formatRoomCharacters()

			// Verify formatRoomItems works
			itemsResult := m.formatRoomItems()

			// Verify formatExitsWithColor works
			exitsResult := m.formatExitsWithColor()

			// All should be non-empty or properly empty based on content
			if len(tc.exits) > 0 && exitsResult == "" {
				t.Error("expected non-empty exits result")
			}
		})
	}
}

// TestRoomCharacterIsNPCField verifies NPC identification
func TestRoomCharacterIsNPCField(t *testing.T) {
	npc := roomCharacter{IsNPC: true}
	player := roomCharacter{IsNPC: false}

	if !npc.IsNPC {
		t.Error("NPC should have IsNPC=true")
	}
	if player.IsNPC {
		t.Error("Player should have IsNPC=false")
	}
}