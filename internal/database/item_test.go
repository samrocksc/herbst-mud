package database

import (
	"testing"
	"time"

	"github.com/sam/makeathing/internal/items"
)

func TestItem(t *testing.T) {
	// Create a new database in memory for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewItemRepository(db)

	// Create test data
	testItem := &Item{
		ID:          "test_item",
		Name:        "Test Item",
		Description: "A test item",
		Type:        "weapon",
		StatsJSON:   `{"strength": 3, "intelligence": 0, "dexterity": 1}`,
		IsMagical:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test create
	err = repo.Create(testItem)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Test get by ID
	retrieved, err := repo.GetByID("test_item")
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected item, got nil")
	}

	if retrieved.ID != testItem.ID {
		t.Errorf("Expected ID '%s', got '%s'", testItem.ID, retrieved.ID)
	}

	if retrieved.Name != testItem.Name {
		t.Errorf("Expected name '%s', got '%s'", testItem.Name, retrieved.Name)
	}

	// Test update
	retrieved.Name = "Updated Test Item"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}

	// Verify update
	updated, err := repo.GetByID("test_item")
	if err != nil {
		t.Fatalf("Failed to get updated item: %v", err)
	}

	if updated == nil {
		t.Fatal("Expected updated item, got nil")
	}

	if updated.Name != "Updated Test Item" {
		t.Errorf("Expected updated name 'Updated Test Item', got '%s'", updated.Name)
	}

	// Test get all
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all items: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 item, got %d", len(all))
	}

	// Test delete
	err = repo.Delete("test_item")
	if err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}

	// Verify deletion
	deleted, err := repo.GetByID("test_item")
	if err != nil {
		t.Fatalf("Failed to check deleted item: %v", err)
	}

	if deleted != nil {
		t.Error("Expected nil after deletion, got item")
	}
}

func TestItemConversion(t *testing.T) {
	// Test ItemFromJSONItem and ToJSONItem conversion
	jsonItem := &items.ItemJSON{
		Schema:      "../schemas/item.schema.json",
		ID:          "test_item",
		Name:        "Test Item",
		Description: "A test item",
		Type:        "weapon",
		Stats: items.ItemStats{
			Strength:     3,
			Intelligence: 0,
			Dexterity:    1,
		},
		IsMagical: false,
	}

	// Convert to database item
	dbItem, err := ItemFromJSONItem(jsonItem)
	if err != nil {
		t.Fatalf("Failed to convert JSON item to database item: %v", err)
	}

	if dbItem.ID != jsonItem.ID {
		t.Errorf("Expected ID '%s', got '%s'", jsonItem.ID, dbItem.ID)
	}

	if dbItem.Name != jsonItem.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonItem.Name, dbItem.Name)
	}

	// Convert back to JSON item
	converted, err := dbItem.ToJSONItem()
	if err != nil {
		t.Fatalf("Failed to convert database item to JSON item: %v", err)
	}

	if converted.ID != jsonItem.ID {
		t.Errorf("Expected ID '%s', got '%s'", jsonItem.ID, converted.ID)
	}

	if converted.Name != jsonItem.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonItem.Name, converted.Name)
	}

	if converted.Stats.Strength != jsonItem.Stats.Strength {
		t.Errorf("Expected strength %d, got %d", jsonItem.Stats.Strength, converted.Stats.Strength)
	}
}