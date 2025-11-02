package database

import (
	"testing"
	"time"
)

func TestDatabase(t *testing.T) {
	// Create a new database in memory for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test configuration repository
	t.Run("Configuration", func(t *testing.T) {
		repo := NewConfigurationRepository(db)

		// Create a new configuration
		config, err := repo.Create("Test MUD")
		if err != nil {
			t.Fatalf("Failed to create configuration: %v", err)
		}

		if config.Name != "Test MUD" {
			t.Errorf("Expected name 'Test MUD', got '%s'", config.Name)
		}

		// Retrieve configuration by name
		retrieved, err := repo.GetByName("Test MUD")
		if err != nil {
			t.Fatalf("Failed to get configuration: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected configuration, got nil")
		}

		if retrieved.ID != config.ID {
			t.Errorf("Expected ID %d, got %d", config.ID, retrieved.ID)
		}

		if retrieved.Name != config.Name {
			t.Errorf("Expected name '%s', got '%s'", config.Name, retrieved.Name)
		}

		// Update configuration
		retrieved.Name = "Updated MUD"
		err = repo.Update(retrieved)
		if err != nil {
			t.Fatalf("Failed to update configuration: %v", err)
		}

		// Verify update
		updated, err := repo.GetByName("Updated MUD")
		if err != nil {
			t.Fatalf("Failed to get updated configuration: %v", err)
		}

		if updated == nil {
			t.Fatal("Expected updated configuration, got nil")
		}

		if updated.Name != "Updated MUD" {
			t.Errorf("Expected name 'Updated MUD', got '%s'", updated.Name)
		}

		// Get all configurations
		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get all configurations: %v", err)
		}

		if len(all) != 1 {
			t.Errorf("Expected 1 configuration, got %d", len(all))
		}

		// Delete configuration
		err = repo.Delete(updated.ID)
		if err != nil {
			t.Fatalf("Failed to delete configuration: %v", err)
		}

		// Verify deletion
		deleted, err := repo.GetByName("Updated MUD")
		if err != nil {
			t.Fatalf("Failed to check deleted configuration: %v", err)
		}

		if deleted != nil {
			t.Error("Expected nil after deletion, got configuration")
		}
	})

	// Test user repository
	t.Run("User", func(t *testing.T) {
		repo := NewUserRepository(db)

		// Create a new user
		user := &User{
			CharacterID: "char123",
			RoomID:      "room456",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.Create(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Retrieve user by ID
		retrieved, err := repo.GetByID(user.ID)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected user, got nil")
		}

		if retrieved.CharacterID != user.CharacterID {
			t.Errorf("Expected character ID '%s', got '%s'", user.CharacterID, retrieved.CharacterID)
		}

		if retrieved.RoomID != user.RoomID {
			t.Errorf("Expected room ID '%s', got '%s'", user.RoomID, retrieved.RoomID)
		}

		// Retrieve user by character ID
		retrievedByChar, err := repo.GetByCharacterID(user.CharacterID)
		if err != nil {
			t.Fatalf("Failed to get user by character ID: %v", err)
		}

		if retrievedByChar == nil {
			t.Fatal("Expected user, got nil")
		}

		if retrievedByChar.ID != user.ID {
			t.Errorf("Expected ID %d, got %d", user.ID, retrievedByChar.ID)
		}

		// Update user
		retrieved.RoomID = "room789"
		err = repo.Update(retrieved)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify update
		updated, err := repo.GetByID(user.ID)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		if updated == nil {
			t.Fatal("Expected updated user, got nil")
		}

		if updated.RoomID != "room789" {
			t.Errorf("Expected room ID 'room789', got '%s'", updated.RoomID)
		}

		// Get all users
		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get all users: %v", err)
		}

		if len(all) != 1 {
			t.Errorf("Expected 1 user, got %d", len(all))
		}

		// Delete user
		err = repo.Delete(updated.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify deletion
		deleted, err := repo.GetByID(user.ID)
		if err != nil {
			t.Fatalf("Failed to check deleted user: %v", err)
		}

		if deleted != nil {
			t.Error("Expected nil after deletion, got user")
		}
	})

	// Test session repository
	t.Run("Session", func(t *testing.T) {
		repo := NewSessionRepository(db)

		// Create a new session
		session := &Session{
			ID:         "session123",
			UserID:     1,
			CharacterID: "char456",
			RoomID:     "room789",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := repo.Create(session)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Retrieve session by ID
		retrieved, err := repo.GetByID(session.ID)
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected session, got nil")
		}

		if retrieved.ID != session.ID {
			t.Errorf("Expected ID '%s', got '%s'", session.ID, retrieved.ID)
		}

		if retrieved.UserID != session.UserID {
			t.Errorf("Expected user ID %d, got %d", session.UserID, retrieved.UserID)
		}

		if retrieved.CharacterID != session.CharacterID {
			t.Errorf("Expected character ID '%s', got '%s'", session.CharacterID, retrieved.CharacterID)
		}

		if retrieved.RoomID != session.RoomID {
			t.Errorf("Expected room ID '%s', got '%s'", session.RoomID, retrieved.RoomID)
		}

		// Update session
		retrieved.RoomID = "room000"
		err = repo.Update(retrieved)
		if err != nil {
			t.Fatalf("Failed to update session: %v", err)
		}

		// Verify update
		updated, err := repo.GetByID(session.ID)
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}

		if updated == nil {
			t.Fatal("Expected updated session, got nil")
		}

		if updated.RoomID != "room000" {
			t.Errorf("Expected room ID 'room000', got '%s'", updated.RoomID)
		}

		// Get all sessions
		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get all sessions: %v", err)
		}

		if len(all) != 1 {
			t.Errorf("Expected 1 session, got %d", len(all))
		}

		// Delete session
		err = repo.Delete(session.ID)
		if err != nil {
			t.Fatalf("Failed to delete session: %v", err)
		}

		// Verify deletion
		deleted, err := repo.GetByID(session.ID)
		if err != nil {
			t.Fatalf("Failed to check deleted session: %v", err)
		}

		if deleted != nil {
			t.Error("Expected nil after deletion, got session")
		}
	})
}