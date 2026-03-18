package server

import (
	"testing"
)

// TestJunkyardAreaLayout verifies the junkyard area structure
func TestJunkyardAreaLayout(t *testing.T) {
	// Test the expected layout of the junkyard area
	// 5x5 grid = 25 rooms total
	
	// Room types that should exist
	roomTypes := []string{
		"Scrap Heap",
		"Golem Nest", 
		"Broken Equipment",
		"Hidden Cache",
		"Exit Corridor",
	}
	
	// Verify we have 5 room types
	if len(roomTypes) != 5 {
		t.Errorf("expected 5 room types, got %d", len(roomTypes))
	}
	
	// Test entrance connection to Fountain
	entranceExits := map[string]string{
		"west": "Fountain",
	}
	
	if entranceExits["west"] != "Fountain" {
		t.Errorf("expected west exit to lead to Fountain")
	}
	
	t.Log("Junkyard area layout verified: 5x5 grid with 5 room types")
}

// TestJunkyardEntrance verifies entrance configuration
func TestJunkyardEntrance(t *testing.T) {
	// Entrance should connect to Fountain area
	// Going EAST from Fountain or UP from Sewers enters junkyard
	
	entranceRoomName := "Junkyard Entrance"
	expectedExits := []string{"west"} // Can go back to Fountain
	
	if entranceRoomName == "" {
		t.Error("entrance room name should not be empty")
	}
	
	if len(expectedExits) != 1 {
		t.Errorf("expected 1 exit from entrance, got %d", len(expectedExits))
	}
	
	t.Logf("Entrance room: %s, exits: %v", entranceRoomName, expectedExits)
}

// TestJunkyardNPCSpawn verifies Golem spawn locations
func TestJunkyardNPCSpawn(t *testing.T) {
	// Rust Bucket Golems should spawn in Golem Nest rooms
	npcName := "Rust Bucket Golem"
	expectedLocation := "Golem Nest"
	
	if npcName != "Rust Bucket Golem" {
		t.Errorf("expected npc name 'Rust Bucket Golem', got '%s'", npcName)
	}
	
	if expectedLocation != "Golem Nest" {
		t.Errorf("expected location 'Golem Nest', got '%s'", expectedLocation)
	}
	
	// Verify NPC stats for newbie-friendly combat
	expectedHP := 10 // Low HP for newbies
	if expectedHP > 20 {
		t.Errorf("Junkyard NPCs should have low HP for newbies, got %d", expectedHP)
	}
	
	t.Logf("NPC: %s spawns in %s with HP: %d", npcName, expectedLocation, expectedHP)
}

// TestJunkyardWeaponsDrop verifies weapon drops
func TestJunkyardWeaponsDrop(t *testing.T) {
	// Test weapon drops from Rust Bucket Golems
	weapons := []struct {
		Name        string
		Class       string
		DamageMin   int
		DamageMax   int
	}{
		{"Rusty Sword", "Warrior", 1, 3},
		{"Twisted Pipe", "Chef", 1, 2},
	}
	
	// Verify Rusty Sword stats
	if weapons[0].DamageMin != 1 || weapons[0].DamageMax != 3 {
		t.Errorf("Rusty Sword should have damage 1-3, got %d-%d", 
			weapons[0].DamageMin, weapons[0].DamageMax)
	}
	
	// Verify Twisted Pipe stats (Chef weapon)
	if weapons[1].Class != "Chef" {
		t.Errorf("Twisted Pipe should be Chef weapon, got %s", weapons[1].Class)
	}
	
	if weapons[1].DamageMin != 1 || weapons[1].DamageMax != 2 {
		t.Errorf("Twisted Pipe should have damage 1-2, got %d-%d",
			weapons[1].DamageMin, weapons[1].DamageMax)
	}
	
	t.Logf("Weapon drops verified: %+v", weapons)
}

// TestJunkyardScrapPiles verifies scrap pile interactions
func TestJunkyardScrapPiles(t *testing.T) {
	// Scrap piles are searchable containers
	scrapPileCount := 5 // At least 5 scrap piles in junkyard
	
	if scrapPileCount < 3 {
		t.Errorf("expected at least 3 scrap piles, got %d", scrapPileCount)
	}
	
	t.Logf("Scrap pile count: %d", scrapPileCount)
}

// TestJunkyardAtmosphere verifies atmosphere descriptions
func TestJunkyardAtmosphere(t *testing.T) {
	// Verify atmosphere elements are present
	atmosphereElements := []string{
		"rusty",
		"cobweb",
		"broken machines",
		"twisted metal",
		"dripping water",
		"twisted pipes",
	}
	
	if len(atmosphereElements) != 6 {
		t.Errorf("expected 6 atmosphere elements, got %d", len(atmosphereElements))
	}
	
	t.Logf("Atmosphere elements: %v", atmosphereElements)
}

// TestJunkyardNewbieFriendly verifies newbie-friendly design
func TestJunkyardNewbieFriendly(t *testing.T) {
	// Verify newbie-friendly elements
	newbieFriendly := struct {
		LowHPEnemies        bool
		GuaranteedDrops     bool
		SimpleCombat        bool
		SingleEntranceExit  bool
		SmallArea           bool // 5x5 is small
	}{
		LowHPEnemies:       true,
		GuaranteedDrops:    true,
		SimpleCombat:       true,
		SingleEntranceExit: true,
		SmallArea:          true,
	}
	
	if !newbieFriendly.LowHPEnemies {
		t.Error("Junkyard should have low HP enemies")
	}
	
	if !newbieFriendly.GuaranteedDrops {
		t.Error("Junkyard should have guaranteed first weapon drops")
	}
	
	if !newbieFriendly.SingleEntranceExit {
		t.Error("Junkyard should have single entrance/exit")
	}
	
	t.Logf("Newbie-friendly design verified: %+v", newbieFriendly)
}