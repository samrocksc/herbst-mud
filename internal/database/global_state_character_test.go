package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalStateCharacterRepository(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "mud-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := New(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := NewGlobalStateCharacterRepository(db)

	// Test Create
	state := &GlobalStateCharacter{
		CharacterID:   "test-character",
		CurrentRoomID: "test-room",
		Health:        100,
		Status:        "active",
	}

	err = repo.Create(state)
	assert.NoError(t, err)
	assert.True(t, state.ID > 0)

	// Test GetByID
	retrievedState, err := repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedState)
	assert.Equal(t, "test-character", retrievedState.CharacterID)
	assert.Equal(t, "test-room", retrievedState.CurrentRoomID)
	assert.Equal(t, 100, retrievedState.Health)
	assert.Equal(t, "active", retrievedState.Status)

	// Test GetByCharacterID
	retrievedState, err = repo.GetByCharacterID("test-character")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedState)
	assert.Equal(t, "test-character", retrievedState.CharacterID)

	// Test Update
	state.Health = 90
	state.Status = "injured"
	err = repo.Update(state)
	assert.NoError(t, err)

	// Verify the update
	updatedState, err := repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedState)
	assert.Equal(t, 90, updatedState.Health)
	assert.Equal(t, "injured", updatedState.Status)

	// Test UpdateRoom - put the character back in test-room
	err = repo.UpdateRoom("test-character", "test-room")
	assert.NoError(t, err)

	roomUpdatedState, err := repo.GetByCharacterID("test-character")
	assert.NoError(t, err)
	assert.Equal(t, "test-room", roomUpdatedState.CurrentRoomID)

	// Test UpdateHealth
	err = repo.UpdateHealth("test-character", 80)
	assert.NoError(t, err)

	healthUpdatedState, err := repo.GetByCharacterID("test-character")
	assert.NoError(t, err)
	assert.Equal(t, 80, healthUpdatedState.Health)

	// Test UpdateStatus
	err = repo.UpdateStatus("test-character", "critical")
	assert.NoError(t, err)

	statusUpdatedState, err := repo.GetByCharacterID("test-character")
	assert.NoError(t, err)
	assert.Equal(t, "critical", statusUpdatedState.Status)

	// Create another state for testing GetCharactersInRoom
	state2 := &GlobalStateCharacter{
		CharacterID:   "test-character-2",
		CurrentRoomID: "test-room",
		Health:        100,
		Status:        "active",
	}

	err = repo.Create(state2)
	assert.NoError(t, err)
	assert.True(t, state2.ID > 0)

	// Verify state2 was created
	allStates, err := repo.GetAll()
	assert.NoError(t, err)
	t.Logf("Total states: %d", len(allStates))
	for i, s := range allStates {
		t.Logf("State %d: %s in room %s", i, s.CharacterID, s.CurrentRoomID)
	}

	// Test GetCharactersInRoom - both should be in test-room now
	charactersInRoom, err := repo.GetCharactersInRoom("test-room")
	assert.NoError(t, err)
	t.Logf("Characters in test-room: %d", len(charactersInRoom))
	for i, char := range charactersInRoom {
		t.Logf("Character %d: %s", i, char.CharacterID)
	}
	// Both state and state2 should be in test-room
	assert.Len(t, charactersInRoom, 2)

	// Test GetAll
	allStates = nil
	allStates, err = repo.GetAll()
	assert.NoError(t, err)
	// We now have 2 states: state and state2
	assert.Equal(t, 2, len(allStates))

	// Test InitializeCharacterState
	err = repo.InitializeCharacterState("new-character", "new-room", 100)
	assert.NoError(t, err)

	newState, err := repo.GetByCharacterID("new-character")
	assert.NoError(t, err)
	assert.NotNil(t, newState)
	assert.Equal(t, "new-character", newState.CharacterID)
	assert.Equal(t, "new-room", newState.CurrentRoomID)
	assert.Equal(t, 100, newState.Health)
	assert.Equal(t, "active", newState.Status)

	// Test DeleteByCharacterID
	err = repo.DeleteByCharacterID("new-character")
	assert.NoError(t, err)

	deletedState, err := repo.GetByCharacterID("new-character")
	assert.NoError(t, err)
	assert.Nil(t, deletedState)

	// Test Delete
	err = repo.Delete(state.ID)
	assert.NoError(t, err)

	deletedState, err = repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedState)
}