package database

import (
	"testing"
	"time"

	"github.com/sam/makeathing/internal/actions"
	"github.com/sam/makeathing/internal/characters"
)

func TestAction(t *testing.T) {
	// Create a new database in memory for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	repo := NewActionRepository(db)

	// Create test data
	testAction := &Action{
		Name:                 "test_action",
		Type:                 "combat",
		Description:          "A test action",
		MinLevel:             1,
		RequiredStatsJSON:    `{"strength": 3, "intelligence": 1, "dexterity": 2}`,
		RequiredSkillsJSON:   `[]`,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Test create
	err = repo.Create(testAction)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	// Test get by name
	retrieved, err := repo.GetByName("test_action")
	if err != nil {
		t.Fatalf("Failed to get action: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected action, got nil")
	}

	if retrieved.Name != testAction.Name {
		t.Errorf("Expected name '%s', got '%s'", testAction.Name, retrieved.Name)
	}

	if retrieved.Type != testAction.Type {
		t.Errorf("Expected type '%s', got '%s'", testAction.Type, retrieved.Type)
	}

	// Test get by ID
	retrievedByID, err := repo.GetByID(retrieved.ID)
	if err != nil {
		t.Fatalf("Failed to get action by ID: %v", err)
	}

	if retrievedByID == nil {
		t.Fatal("Expected action, got nil")
	}

	if retrievedByID.Name != testAction.Name {
		t.Errorf("Expected name '%s', got '%s'", testAction.Name, retrievedByID.Name)
	}

	// Test update
	retrieved.Description = "Updated test action"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update action: %v", err)
	}

	// Verify update
	updated, err := repo.GetByName("test_action")
	if err != nil {
		t.Fatalf("Failed to get updated action: %v", err)
	}

	if updated == nil {
		t.Fatal("Expected updated action, got nil")
	}

	if updated.Description != "Updated test action" {
		t.Errorf("Expected updated description 'Updated test action', got '%s'", updated.Description)
	}

	// Test get all
	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all actions: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 action, got %d", len(all))
	}

	// Test delete
	err = repo.Delete(updated.ID)
	if err != nil {
		t.Fatalf("Failed to delete action: %v", err)
	}

	// Verify deletion
	deleted, err := repo.GetByName("test_action")
	if err != nil {
		t.Fatalf("Failed to check deleted action: %v", err)
	}

	if deleted != nil {
		t.Error("Expected nil after deletion, got action")
	}
}

func TestActionConversion(t *testing.T) {
	// Test ActionFromJSONAction and ToJSONAction conversion
	jsonAction := &actions.Action{
		Name:        "test_action",
		Type:        characters.Spell,
		Description: "A test action",
		Requirements: actions.ActionRequirements{
			MinLevel: 1,
			RequiredStats: characters.Stats{
				Strength:     3,
				Intelligence: 1,
				Dexterity:    2,
			},
			RequiredSkills: []string{},
		},
	}

	// Convert to database action
	dbAction, err := ActionFromJSONAction(jsonAction)
	if err != nil {
		t.Fatalf("Failed to convert JSON action to database action: %v", err)
	}

	if dbAction.Name != jsonAction.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonAction.Name, dbAction.Name)
	}

	if dbAction.Type != string(jsonAction.Type) {
		t.Errorf("Expected type '%s', got '%s'", string(jsonAction.Type), dbAction.Type)
	}

	// Convert back to JSON action
	converted, err := dbAction.ToJSONAction()
	if err != nil {
		t.Fatalf("Failed to convert database action to JSON action: %v", err)
	}

	if converted.Name != jsonAction.Name {
		t.Errorf("Expected name '%s', got '%s'", jsonAction.Name, converted.Name)
	}

	if converted.Type != jsonAction.Type {
		t.Errorf("Expected type '%s', got '%s'", jsonAction.Type, converted.Type)
	}

	if converted.Requirements.MinLevel != jsonAction.Requirements.MinLevel {
		t.Errorf("Expected min level %d, got %d", jsonAction.Requirements.MinLevel, converted.Requirements.MinLevel)
	}
}