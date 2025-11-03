package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sam/makeathing/internal/configuration"
	"github.com/sam/makeathing/internal/rooms"
	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestDBAdapter_LoadActualDataIntoGlobalState(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "mud-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	adapter, err := NewDBAdapter(dbPath)
	assert.NoError(t, err)
	defer adapter.Close()

	// Load configuration from JSON
	configJSON, err := configuration.LoadConfigurationFromJSON("../../data/configuration.json")
	assert.NoError(t, err)

	// Create configuration in database
	err = adapter.SetConfiguration(configJSON.Name, configJSON.Name)
	assert.NoError(t, err)

	// Load users from JSON directory and create them
	usersMap, err := users.LoadAllUserJSONsFromDirectory("../../data/users")
	assert.NoError(t, err)

	for _, userJSON := range usersMap {
		// Convert UserJSON to User and create
		user := &User{
			ID:          userJSON.ID,
			CharacterID: userJSON.CharacterID,
			RoomID:      userJSON.RoomID,
		}
		err = adapter.userRepo.Create(user)
		assert.NoError(t, err)
	}

	// Load rooms from JSON directory and create them
	roomsDir := "../../data/rooms"
	roomFiles, err := filepath.Glob(filepath.Join(roomsDir, "*.json"))
	assert.NoError(t, err)

	for _, roomFile := range roomFiles {
		roomJSON, err := rooms.LoadRoomJSONFromJSON(roomFile)
		if err != nil {
			t.Logf("Warning: Could not load room %s: %v", roomFile, err)
			continue
		}

		err = adapter.CreateRoom(roomJSON)
		assert.NoError(t, err)
	}

	// Load characters from JSON directory and create them
	charactersDir := "../../data/characters"
	characterFiles, err := filepath.Glob(filepath.Join(charactersDir, "*.json"))
	assert.NoError(t, err)

	for _, characterFile := range characterFiles {
		characterJSON, err := characters.LoadCharacterJSONFromJSON(characterFile)
		if err != nil {
			t.Logf("Warning: Could not load character %s: %v", characterFile, err)
			continue
		}

		err = adapter.CreateCharacter(characterJSON)
		assert.NoError(t, err)
	}

	// Now test the global state initialization
	err = adapter.InitializeGlobalState()
	assert.NoError(t, err)

	// Verify that global state was created for all rooms
	allRoomStates, err := adapter.globalStateRoomRepo.GetAll()
	assert.NoError(t, err)
	t.Logf("Created %d room states", len(allRoomStates))
	assert.True(t, len(allRoomStates) > 0, "Should have created room states")

	// Verify that global state was created for all characters
	allCharacterStates, err := adapter.globalStateCharacterRepo.GetAll()
	assert.NoError(t, err)
	t.Logf("Created %d character states", len(allCharacterStates))
	assert.True(t, len(allCharacterStates) > 0, "Should have created character states")

	// Verify the content of created states
	for _, roomState := range allRoomStates {
		t.Logf("Room State: %s - Player Count: %d", roomState.RoomID, roomState.PlayerCount)
	}

	for _, characterState := range allCharacterStates {
		t.Logf("Character State: %s in room %s - Health: %d - Status: %s", 
			characterState.CharacterID, characterState.CurrentRoomID, characterState.Health, characterState.Status)
	}

	// Test getting specific states
	if len(allRoomStates) > 0 {
		firstRoomState, err := adapter.GetRoomState(allRoomStates[0].RoomID)
		assert.NoError(t, err)
		assert.NotNil(t, firstRoomState)
	}

	if len(allCharacterStates) > 0 {
		firstCharacterState, err := adapter.GetCharacterState(allCharacterStates[0].CharacterID)
		assert.NoError(t, err)
		assert.NotNil(t, firstCharacterState)
	}

	// Test the bootstrap utility function
	err = BootstrapGlobalState(adapter)
	assert.NoError(t, err)

	// Test idempotency - calling it again should not cause issues
	err = adapter.InitializeGlobalState()
	assert.NoError(t, err)

	// Verify no duplicate states were created
	finalRoomStates, err := adapter.globalStateRoomRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, len(allRoomStates), len(finalRoomStates))

	finalCharacterStates, err := adapter.globalStateCharacterRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, len(allCharacterStates), len(finalCharacterStates))
}