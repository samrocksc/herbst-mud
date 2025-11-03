package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
	"github.com/sam/makeathing/internal/rooms"
	"github.com/stretchr/testify/assert"
)

func TestDBAdapter_InitializeGlobalState(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "mud-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	adapter, err := NewDBAdapter(dbPath)
	assert.NoError(t, err)
	defer adapter.Close()

	// Create a test room
	testRoom := &rooms.RoomJSON{
		Schema:     "../schemas/room.schema.json",
		ID:         "test-room-1",
		Description: "A test room",
		Exits:      map[string]string{"north": "test-room-2"},
		Smells:     "Fresh air",
	}

	err = adapter.CreateRoom(testRoom)
	assert.NoError(t, err)

	// Create a test character
	testCharacter := &characters.CharacterJSON{
		Schema:     "../schemas/character.schema.json",
		ID:         "test-character-1",
		Name:       "Test Character",
		Race:       "Human",
		Class:      "Warrior",
		Stats:      characters.Stats{Strength: 15, Intelligence: 10, Dexterity: 12},
		Health:     100,
		Mana:       50,
		Experience: 0,
		Level:      1,
		IsVendor:   false,
		IsNpc:      false,
		Inventory:  []items.Item{},
		Skills:     []characters.Skill{},
	}

	err = adapter.CreateCharacter(testCharacter)
	assert.NoError(t, err)

	// Create a user to link the character to a room
	userID, err := adapter.CreateUser(testCharacter.ID, testRoom.ID)
	assert.NoError(t, err)
	assert.True(t, userID > 0)

	// Initialize global state
	err = adapter.InitializeGlobalState()
	assert.NoError(t, err)

	// Verify room state was created
	roomState, err := adapter.GetRoomState("test-room-1")
	assert.NoError(t, err)
	assert.NotNil(t, roomState)
	assert.Equal(t, "test-room-1", roomState.RoomID)
	assert.Equal(t, 0, roomState.PlayerCount)

	// Verify character state was created
	characterState, err := adapter.GetCharacterState("test-character-1")
	assert.NoError(t, err)
	assert.NotNil(t, characterState)
	assert.Equal(t, "test-character-1", characterState.CharacterID)
	assert.Equal(t, "test-room-1", characterState.CurrentRoomID)
	assert.Equal(t, 100, characterState.Health)
	assert.Equal(t, "active", characterState.Status)

	// Test InitializeGlobalStateForCharacter
	// Create another character and test individual initialization
	testCharacter2 := &characters.CharacterJSON{
		Schema:     "../schemas/character.schema.json",
		ID:         "test-character-2",
		Name:       "Test Character 2",
		Race:       "Dwarf",
		Class:      "Mage",
		Stats:      characters.Stats{Strength: 8, Intelligence: 18, Dexterity: 10},
		Health:     80,
		Mana:       100,
		Experience: 0,
		Level:      1,
		IsVendor:   false,
		IsNpc:      false,
		Inventory:  []items.Item{},
		Skills:     []characters.Skill{},
	}

	err = adapter.CreateCharacter(testCharacter2)
	assert.NoError(t, err)

	// Initialize global state for just this character
	err = adapter.InitializeGlobalStateForCharacter("test-character-2")
	assert.NoError(t, err)

	// Verify the second character state was created
	characterState2, err := adapter.GetCharacterState("test-character-2")
	assert.NoError(t, err)
	assert.NotNil(t, characterState2)
	assert.Equal(t, "test-character-2", characterState2.CharacterID)
	assert.Equal(t, "starting_room", characterState2.CurrentRoomID) // Should default to starting room
	assert.Equal(t, 80, characterState2.Health)
	assert.Equal(t, "active", characterState2.Status)

	// Test InitializeGlobalStateForRoom
	// Create another room and test individual initialization
	testRoom2 := &rooms.RoomJSON{
		Schema:     "../schemas/room.schema.json",
		ID:         "test-room-2",
		Description: "Another test room",
		Exits:      map[string]string{"south": "test-room-1"},
		Smells:     "Musty air",
	}

	err = adapter.CreateRoom(testRoom2)
	assert.NoError(t, err)

	// Initialize global state for just this room
	err = adapter.InitializeGlobalStateForRoom("test-room-2")
	assert.NoError(t, err)

	// Verify the second room state was created
	roomState2, err := adapter.GetRoomState("test-room-2")
	assert.NoError(t, err)
	assert.NotNil(t, roomState2)
	assert.Equal(t, "test-room-2", roomState2.RoomID)
	assert.Equal(t, 0, roomState2.PlayerCount)

	// Test idempotency - calling InitializeGlobalState again should not fail
	err = adapter.InitializeGlobalState()
	assert.NoError(t, err)

	// Verify all states still exist and haven't been duplicated
	allRoomStates, err := adapter.globalStateRoomRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allRoomStates))

	allCharacterStates, err := adapter.globalStateCharacterRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allCharacterStates))
}