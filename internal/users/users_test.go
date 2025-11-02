package users

import (
	"testing"
)

func TestLoadUserFromJSON(t *testing.T) {
	// Load the user file
	user, err := LoadUserFromJSON("../../data/users/user_1.json")
	if err != nil {
		t.Fatalf("Error loading user: %v", err)
	}

	// Check that the user was loaded correctly
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	if user.CharacterID != "char_nelly" {
		t.Errorf("Expected character ID 'char_nelly', got '%s'", user.CharacterID)
	}

	if user.RoomID != "start" {
		t.Errorf("Expected room ID 'start', got '%s'", user.RoomID)
	}
}

func TestLoadAllUsersFromDirectory(t *testing.T) {
	// Load all users from the data directory
	users, err := LoadAllUsersFromDirectory("../../data/users")
	if err != nil {
		t.Fatalf("Error loading users: %v", err)
	}

	// Check that we loaded exactly one user
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Check that the user was loaded correctly
	user, exists := users["user_1"]
	if !exists {
		t.Fatal("Expected user with key 'user_1' not found")
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	if user.CharacterID != "char_nelly" {
		t.Errorf("Expected character ID 'char_nelly', got '%s'", user.CharacterID)
	}

	if user.RoomID != "start" {
		t.Errorf("Expected room ID 'start', got '%s'", user.RoomID)
	}
}