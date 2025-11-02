package database

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
)

func TestRoom(t *testing.T) {
	// Create a new database in memory for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewRoomRepository(db)

	// Create test data
	testRoom := &Room{
		ID:          "test_room",
		Description: "A test room for testing",
		Smells:      "It smells testy",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Convert some test data to JSON for the complex fields
	exits := map[string]string{"north": "north_room", "south": "south_room"}
	exitsJSON, err := json.Marshal(exits)
	if err != nil {
		t.Fatalf("Failed to marshal exits: %v", err)
	}
	testRoom.ExitsJSON = string(exitsJSON)

	// Create a test item
	testItem := items.Item{
		ID:          "test_item",
		Name:        "Test Item",
		Description: "An item for testing",
		Type:        items.Movable,
		Stats: items.ItemStats{
			Strength:     1,
			Intelligence: 1,
			Dexterity:    1,
		},
		IsMagical: false,
	}
	movableObjects := []items.Item{testItem}
	movableObjectsJSON, err := json.Marshal(movableObjects)
	if err != nil {
		t.Fatalf("Failed to marshal movable objects: %v", err)
	}
	testRoom.MovableObjectsJSON = string(movableObjectsJSON)

	// Create a test character
	testCharacter := characters.Character{
		ID:          "test_npc",
		Name:        "Test NPC",
		Race:        characters.Human,
		Class:       characters.Warrior,
		Stats: characters.Stats{
			Strength:     10,
			Intelligence: 10,
			Dexterity:    10,
		},
		Health:     100,
		Mana:       50,
		Experience: 0,
		Level:      1,
		IsVendor:   false,
		IsNpc:      true,
		Inventory:  []items.Item{},
		Skills:     []characters.Skill{},
	}
	npcs := []characters.Character{testCharacter}
	npcsJSON, err := json.Marshal(npcs)
	if err != nil {
		t.Fatalf("Failed to marshal NPCs: %v", err)
	}
	testRoom.NPCsJSON = string(npcsJSON)

	// Test create
	err = repo.Create(testRoom)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	// Test get by ID
	retrieved, err := repo.GetByID("test_room")
	if err != nil {
		t.Fatalf("Failed to get room: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected room, got nil")
	}

	if retrieved.ID != testRoom.ID {
		t.Errorf("Expected ID '%s', got '%s'", testRoom.ID, retrieved.ID)
	}

	if retrieved.Description != testRoom.Description {
		t.Errorf("Expected description '%s', got '%s'", testRoom.Description, retrieved.Description)
	}

	if retrieved.Smells != testRoom.Smells {
		t.Errorf("Expected smells '%s', got '%s'", testRoom.Smells, retrieved.Smells)
	}

	// Test update
	retrieved.Description = "Updated test room"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update room: %v", err)
	}

	// Verify update
	updated, err := repo.GetByID("test_room")
	if err != nil {
		t.Fatalf("Failed to get updated room: %v", err)
	}

	if updated == nil {
		t.Fatal("Expected updated room, got nil")
	}

	if updated.Description != "Updated test room" {
		t.Errorf("Expected updated description 'Updated test room', got '%s'", updated.Description)
	}

	// Test get all
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all rooms: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 room, got %d", len(all))
	}

	// Test delete
	err = repo.Delete("test_room")
	if err != nil {
		t.Fatalf("Failed to delete room: %v", err)
	}

	// Verify deletion
	deleted, err := repo.GetByID("test_room")
	if err != nil {
		t.Fatalf("Failed to check deleted room: %v", err)
	}

	if deleted != nil {
		t.Error("Expected nil after deletion, got room")
	}
}

func TestRoomConversion(t *testing.T) {
	// Test RoomFromJSONRoom and ToJSONRoom conversion
	exits := map[string]string{"north": "north_room", "south": "south_room"}
	
	// Create test items
	testItems := []items.Item{
		{
			ID:          "test_item_1",
			Name:        "Test Item 1",
			Description: "First test item",
			Type:        items.Movable,
			Stats: items.ItemStats{
				Strength:     1,
				Intelligence: 1,
				Dexterity:    1,
			},
			IsMagical: false,
		},
	}

	// Create test characters
	testCharacters := []characters.Character{
		{
			ID:          "test_npc_1",
			Name:        "Test NPC 1",
			Race:        characters.Human,
			Class:       characters.Warrior,
			Stats: characters.Stats{
				Strength:     10,
				Intelligence: 10,
				Dexterity:    10,
			},
			Health:     100,
			Mana:       50,
			Experience: 0,
			Level:      1,
			IsVendor:   false,
			IsNpc:      true,
			Inventory:  []items.Item{},
			Skills:     []characters.Skill{},
		},
	}

	// Create a JSON room
	jsonRoom := &struct {
		Schema           string                 `json:"$schema"`
		ID               string                 `json:"id"`
		Description      string                 `json:"description"`
		Exits            map[string]string      `json:"exits"`
		ImmovableObjects []items.Item           `json:"immovableObjects"`
		MovableObjects   []items.Item           `json:"movableObjects"`
		Smells           string                 `json:"smells"`
		NPCs             []characters.Character `json:"npcs"`
	}{
		Schema:           "../schemas/room.schema.json",
		ID:               "test_room",
		Description:      "A test room",
		Exits:            exits,
		ImmovableObjects: []items.Item{},
		MovableObjects:   testItems,
		Smells:           "Test smell",
		NPCs:             testCharacters,
	}

	// Convert to database room (manual conversion since we don't have access to rooms.RoomJSON)
	// This is a simplified test for the JSON conversion functions
	exitsJSON, _ := json.Marshal(exits)
	movableObjectsJSON, _ := json.Marshal(testItems)
	npcsJSON, _ := json.Marshal(testCharacters)

	dbRoom := &Room{
		ID:                   jsonRoom.ID,
		Description:          jsonRoom.Description,
		Smells:               jsonRoom.Smells,
		ExitsJSON:            string(exitsJSON),
		ImmovableObjectsJSON: "[]",
		MovableObjectsJSON:   string(movableObjectsJSON),
		NPCsJSON:             string(npcsJSON),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Test conversion back to JSON
	converted, err := dbRoom.ToJSONRoom()
	if err != nil {
		t.Fatalf("Failed to convert room to JSON: %v", err)
	}

	if converted.ID != dbRoom.ID {
		t.Errorf("Expected ID '%s', got '%s'", dbRoom.ID, converted.ID)
	}

	if converted.Description != dbRoom.Description {
		t.Errorf("Expected description '%s', got '%s'", dbRoom.Description, converted.Description)
	}
}