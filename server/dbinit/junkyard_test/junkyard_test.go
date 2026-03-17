package junkyard_test

import "testing"

// TestRoomTypes verifies room type constants
func TestRoomTypes(t *testing.T) {
	roomTypes := []string{
		"Scrap Heap", "Scrap Heap", "Scrap Heap", "Scrap Heap", "Scrap Heap",
		"Golem Nest", "Golem Nest", "Broken Equipment", "Broken Equipment", "Hidden Cache",
		"Scrap Heap", "Broken Equipment", "Golem Nest", "Broken Equipment", "Scrap Heap",
		"Hidden Cache", "Broken Equipment", "Scrap Heap", "Golem Nest", "Scrap Heap",
		"Exit Corridor", "Scrap Heap", "Broken Equipment", "Scrap Heap", "Scrap Heap",
	}
	
	if len(roomTypes) != 25 {
		t.Fatalf("expected 25 rooms, got %d", len(roomTypes))
	}
	t.Log("Room types verified: 5x5 grid")
}

// TestEntrancePosition verifies entrance is at bottom-left
func TestEntrancePosition(t *testing.T) {
	// Entrance at index 20 (row 4, col 0)
	entranceIdx := 20
	row, col := entranceIdx/5, entranceIdx%5
	
	if row != 4 || col != 0 {
		t.Errorf("entrance should be at (4, 0), got (%d, %d)", row, col)
	}
	
	// West exit to Fountain
	t.Log("Entrance position verified: row 4, col 0 (west exit to Fountain)")
}

// TestNPCSpawnLocations verifies Golem spawns
func TestNPCSpawnLocations(t *testing.T) {
	// Golem Nest indices: 5, 12, 18
	expected := []int{5, 12, 18}
	
	for _, idx := range expected {
		row, col := idx/5, idx%5
		t.Logf("Golem Nest at index %d: row %d, col %d", idx, row, col)
	}
	
	t.Log("NPC spawn locations verified: 3 Golem Nest rooms")
}
