package database

import (
	"testing"
	"time"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
)

func TestCharacter(t *testing.T) {
	// Create a new database in memory for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewCharacterRepository(db)

	// Create test data
	testCharacter := &Character{
		ID:         "test_character",
		Name:       "Test Character",
		Race:       "Human",
		Class:      "Warrior",
		StatsJSON:  `{"strength": 10, "intelligence": 5, "dexterity": 8}`,
		Health:     100,
		Mana:       50,
		Experience: 0,
		Level:      1,
		IsVendor:   false,
		IsNpc:      false,
		InventoryJSON: `[]`,
		SkillsJSON: `[]`,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test create
	err = repo.Create(testCharacter)
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// Test get by ID
	retrieved, err := repo.GetByID("test_character")
	if err != nil {
		t.Fatalf("Failed to get character: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected character, got nil")
	}

	if retrieved.ID != testCharacter.ID {
		t.Errorf("Expected ID '%s', got '%s'", testCharacter.ID, retrieved.ID)
	}

	if retrieved.Name != testCharacter.Name {
		t.Errorf("Expected name '%s', got '%s'", testCharacter.Name, retrieved.Name)
	}

	// Test update
	retrieved.Name = "Updated Test Character"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update character: %v", err)
	}

	// Verify update
	updated, err := repo.GetByID("test_character")
	if err != nil {
		t.Fatalf("Failed to get updated character: %v", err)
	}

	if updated == nil {
		t.Fatal("Expected updated character, got nil")
	}

	if updated.Name != "Updated Test Character" {
		t.Errorf("Expected updated name 'Updated Test Character', got '%s'", updated.Name)
	}

	// Test get all
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all characters: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 character, got %d", len(all))
	}

	// Test delete
	err = repo.Delete("test_character")
	if err != nil {
		t.Fatalf("Failed to delete character: %v", err)
	}

	// Verify deletion
	deleted, err := repo.GetByID("test_character")
	if err != nil {
		t.Fatalf("Failed to check deleted character: %v", err)
	}

	if deleted != nil {
		t.Error("Expected nil after deletion, got character")
	}
}

func TestCharacterConversion(t *testing.T) {
	// Test CharacterFromJSONCharacter and ToJSONCharacter conversion
	jsonCharacter := &characters.CharacterJSON{
		Schema:     "../schemas/character.schema.json",
		ID:         "test_char",
		Name:       "Test Character",
		Race:       "Human",
		Class:      "Warrior",
		Stats: characters.Stats{
			Strength:     10,
			Intelligence: 5,
			Dexterity:    8,
		},
		Health:     100,
		Mana:       50,
		Experience: 0,
		Level:      1,
		IsVendor:   false,
		IsNpc:      false,
		Inventory:  []items.Item{},
		Skills:     []characters.Skill{},
	}

	// Convert to database character
	dbCharacter, err := CharacterFromJSONCharacter(jsonCharacter)
	if err != nil {
		t.Fatalf("Failed to convert JSON character to database character: %v", err)
	}

	if dbCharacter.ID != jsonCharacter.ID {
		t.Errorf("Expected ID '%s', got '%s'", jsonCharacter.ID, dbCharacter.ID)
	}

	if dbCharacter.Name != jsonCharacter.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonCharacter.Name, dbCharacter.Name)
	}

	// Convert back to JSON character
	converted, err := dbCharacter.ToJSONCharacter()
	if err != nil {
		t.Fatalf("Failed to convert database character to JSON character: %v", err)
	}

	if converted.ID != jsonCharacter.ID {
		t.Errorf("Expected ID '%s', got '%s'", jsonCharacter.ID, converted.ID)
	}

	if converted.Name != jsonCharacter.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonCharacter.Name, converted.Name)
	}

	if converted.Stats.Strength != jsonCharacter.Stats.Strength {
		t.Errorf("Expected strength %d, got %d", jsonCharacter.Stats.Strength, converted.Stats.Strength)
	}
}