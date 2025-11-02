package database

import (
	"testing"
)

func TestDBAdapter(t *testing.T) {
	// Create a new database adapter in memory for testing
	adapter, err := NewDBAdapter(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database adapter: %v", err)
	}
	defer adapter.Close()

	// Test configuration operations
	t.Run("Configuration", func(t *testing.T) {
		// Set configuration
		err := adapter.SetConfiguration("test_mud", "Test MUD Server")
		if err != nil {
			t.Fatalf("Failed to set configuration: %v", err)
		}

		// Get configuration
		config, err := adapter.GetConfiguration("Test MUD Server")
		if err != nil {
			t.Fatalf("Failed to get configuration: %v", err)
		}

		if config == nil {
			t.Fatal("Expected configuration, got nil")
		}

		if config.Name != "Test MUD Server" {
			t.Errorf("Expected name 'Test MUD Server', got '%s'", config.Name)
		}
	})

	// Test user operations
	t.Run("User", func(t *testing.T) {
		// Create user
		userID, err := adapter.CreateUser("char123", "room456")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if userID <= 0 {
			t.Errorf("Expected user ID > 0, got %d", userID)
		}

		// Get user by ID
		user, err := adapter.GetUser(userID)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user, got nil")
		}

		if user.ID != userID {
			t.Errorf("Expected ID %d, got %d", userID, user.ID)
		}

		if user.CharacterID != "char123" {
			t.Errorf("Expected character ID 'char123', got '%s'", user.CharacterID)
		}

		if user.RoomID != "room456" {
			t.Errorf("Expected room ID 'room456', got '%s'", user.RoomID)
		}

		// Get user by character ID
		userByChar, err := adapter.GetUserByCharacterID("char123")
		if err != nil {
			t.Fatalf("Failed to get user by character ID: %v", err)
		}

		if userByChar == nil {
			t.Fatal("Expected user, got nil")
		}

		if userByChar.ID != userID {
			t.Errorf("Expected ID %d, got %d", userID, userByChar.ID)
		}

		// Update user
		user.RoomID = "room789"
		err = adapter.UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify update
		updatedUser, err := adapter.GetUser(userID)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		if updatedUser == nil {
			t.Fatal("Expected updated user, got nil")
		}

		if updatedUser.RoomID != "room789" {
			t.Errorf("Expected room ID 'room789', got '%s'", updatedUser.RoomID)
		}

		// Delete user
		err = adapter.DeleteUser(userID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify deletion
		deletedUser, err := adapter.GetUser(userID)
		if err != nil {
			t.Fatalf("Failed to check deleted user: %v", err)
		}

		if deletedUser != nil {
			t.Error("Expected nil after deletion, got user")
		}
	})

	// Test session operations
	t.Run("Session", func(t *testing.T) {
		// Create session
		sessionID := "session123"
		err := adapter.CreateSession(sessionID, 1, "char456", "room789")
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		// Get session
		session, err := adapter.GetSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		if session == nil {
			t.Fatal("Expected session, got nil")
		}

		if session.ID != sessionID {
			t.Errorf("Expected ID '%s', got '%s'", sessionID, session.ID)
		}

		if session.UserID != 1 {
			t.Errorf("Expected user ID 1, got %d", session.UserID)
		}

		if session.CharacterID != "char456" {
			t.Errorf("Expected character ID 'char456', got '%s'", session.CharacterID)
		}

		if session.RoomID != "room789" {
			t.Errorf("Expected room ID 'room789', got '%s'", session.RoomID)
		}

		// Update session
		session.RoomID = "room000"
		err = adapter.UpdateSession(session)
		if err != nil {
			t.Fatalf("Failed to update session: %v", err)
		}

		// Verify update
		updatedSession, err := adapter.GetSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}

		if updatedSession == nil {
			t.Fatal("Expected updated session, got nil")
		}

		if updatedSession.RoomID != "room000" {
			t.Errorf("Expected room ID 'room000', got '%s'", updatedSession.RoomID)
		}

		// Delete session
		err = adapter.DeleteSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to delete session: %v", err)
		}

		// Verify deletion
		deletedSession, err := adapter.GetSession(sessionID)
		if err != nil {
			t.Fatalf("Failed to check deleted session: %v", err)
		}

		if deletedSession != nil {
			t.Error("Expected nil after deletion, got session")
		}
	})
}