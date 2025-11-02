package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGlobalStateRoomRepository(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "mud-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := New(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := NewGlobalStateRoomRepository(db)

	// Test Create
	state := &GlobalStateRoom{
		RoomID:         "test-room",
		PlayerCount:    2,
		NPCStateJSON:   `[]`,
		ItemStateJSON:  `[]`,
	}

	err = repo.Create(state)
	assert.NoError(t, err)
	assert.True(t, state.ID > 0)

	// Test GetByID
	retrievedState, err := repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedState)
	assert.Equal(t, "test-room", retrievedState.RoomID)
	assert.Equal(t, 2, retrievedState.PlayerCount)

	// Test GetByRoomID
	retrievedState, err = repo.GetByRoomID("test-room")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedState)
	assert.Equal(t, "test-room", retrievedState.RoomID)

	// Test Update
	state.PlayerCount = 3
	state.NPCStateJSON = `[{"npc_id":"npc1","health":50,"current_room_id":"test-room","status":"active"}]`
	err = repo.Update(state)
	assert.NoError(t, err)

	// Verify the update
	updatedState, err := repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedState)
	assert.Equal(t, 3, updatedState.PlayerCount)

	// Test UpdatePlayerCount
	err = repo.UpdatePlayerCount("test-room", 4)
	assert.NoError(t, err)

	playerCountUpdatedState, err := repo.GetByRoomID("test-room")
	assert.NoError(t, err)
	assert.Equal(t, 4, playerCountUpdatedState.PlayerCount)

	// Test UpdateNPCState
	npcState := []NPCState{
		{
			NpcID:          "npc1",
			Health:         100,
			CurrentRoomID:  "test-room",
			Status:         "active",
		},
	}

	err = repo.UpdateNPCState("test-room", npcState)
	assert.NoError(t, err)

	npcStateUpdatedState, err := repo.GetByRoomID("test-room")
	assert.NoError(t, err)

	// Verify the NPC state was saved correctly
	var savedNPCState []NPCState
	err = json.Unmarshal([]byte(npcStateUpdatedState.NPCStateJSON), &savedNPCState)
	assert.NoError(t, err)
	assert.Len(t, savedNPCState, 1)
	assert.Equal(t, "npc1", savedNPCState[0].NpcID)
	assert.Equal(t, 100, savedNPCState[0].Health)

	// Test UpdateItemState
	itemState := []ItemState{
		{
			ItemID:         "item1",
			CurrentRoomID:  "test-room",
			Status:         "available",
		},
	}

	err = repo.UpdateItemState("test-room", itemState)
	assert.NoError(t, err)

	itemStateUpdatedState, err := repo.GetByRoomID("test-room")
	assert.NoError(t, err)

	// Verify the item state was saved correctly
	var savedItemState []ItemState
	err = json.Unmarshal([]byte(itemStateUpdatedState.ItemStateJSON), &savedItemState)
	assert.NoError(t, err)
	assert.Len(t, savedItemState, 1)
	assert.Equal(t, "item1", savedItemState[0].ItemID)

	// Create another state for testing GetAll
	state2 := &GlobalStateRoom{
		RoomID:         "test-room-2",
		PlayerCount:    1,
		NPCStateJSON:   `[]`,
		ItemStateJSON:  `[]`,
		LastUpdated:    time.Now(),
		CreatedAt:      time.Now(),
	}

	err = repo.Create(state2)
	assert.NoError(t, err)

	// Test GetAll
	allStates, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, allStates, 2)

	// Test InitializeRoomState
	err = repo.InitializeRoomState("new-room")
	assert.NoError(t, err)

	newState, err := repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.NotNil(t, newState)
	assert.Equal(t, "new-room", newState.RoomID)
	assert.Equal(t, 0, newState.PlayerCount)

	// Test GetNPCState
	npcStateFromDB, err := repo.GetNPCState("test-room")
	assert.NoError(t, err)
	assert.Len(t, npcStateFromDB, 1)
	assert.Equal(t, "npc1", npcStateFromDB[0].NpcID)
	assert.Equal(t, 100, npcStateFromDB[0].Health)

	// Test GetItemState
	itemStateFromDB, err := repo.GetItemState("test-room")
	assert.NoError(t, err)
	assert.Len(t, itemStateFromDB, 1)
	assert.Equal(t, "item1", itemStateFromDB[0].ItemID)

	// Test IncrementPlayerCount
	err = repo.IncrementPlayerCount("new-room")
	assert.NoError(t, err)

	incrementedState, err := repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.Equal(t, 1, incrementedState.PlayerCount)

	err = repo.IncrementPlayerCount("new-room")
	assert.NoError(t, err)

	incrementedState, err = repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.Equal(t, 2, incrementedState.PlayerCount)

	// Test DecrementPlayerCount
	err = repo.DecrementPlayerCount("new-room")
	assert.NoError(t, err)

	decrementedState, err := repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.Equal(t, 1, decrementedState.PlayerCount)

	err = repo.DecrementPlayerCount("new-room")
	assert.NoError(t, err)

	decrementedState, err = repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.Equal(t, 0, decrementedState.PlayerCount)

	// Test DeleteByRoomID
	err = repo.DeleteByRoomID("new-room")
	assert.NoError(t, err)

	deletedState, err := repo.GetByRoomID("new-room")
	assert.NoError(t, err)
	assert.Nil(t, deletedState)

	// Test Delete
	err = repo.Delete(state.ID)
	assert.NoError(t, err)

	deletedState, err = repo.GetByID(state.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedState)
}