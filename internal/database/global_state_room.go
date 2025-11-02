package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

// GlobalStateRoom represents a room state in the database
type GlobalStateRoom struct {
	ID             int               `json:"id"`
	RoomID         string            `json:"room_id"`
	PlayerCount    int               `json:"player_count"`
	NPCStateJSON   string            `json:"npc_state_json"`
	ItemStateJSON  string            `json:"item_state_json"`
	LastUpdated    time.Time         `json:"last_updated"`
	CreatedAt      time.Time         `json:"created_at"`
}

// NPCState represents the state of NPCs in a room
type NPCState struct {
	NpcID          string            `json:"npc_id"`
	Health         int               `json:"health"`
	CurrentRoomID  string            `json:"current_room_id"`
	Status         string            `json:"status"`
}

// ItemState represents the state of items in a room
type ItemState struct {
	ItemID         string            `json:"item_id"`
	CurrentRoomID  string            `json:"current_room_id"`
	Status         string            `json:"status"`
}

// GlobalStateRoomRepository provides methods for working with room state
type GlobalStateRoomRepository struct {
	db *DB
}

// NewGlobalStateRoomRepository creates a new room state repository
func NewGlobalStateRoomRepository(db *DB) *GlobalStateRoomRepository {
	return &GlobalStateRoomRepository{db: db}
}

// Create creates a new room state
func (r *GlobalStateRoomRepository) Create(state *GlobalStateRoom) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO global_state_rooms (room_id, player_count, npc_state_json, item_state_json)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		state.RoomID,
		state.PlayerCount,
		state.NPCStateJSON,
		state.ItemStateJSON,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	state.ID = int(id)

	// Get the created state to populate the timestamps
	createdState, err := r.GetByID(state.ID)
	if err != nil {
		return err
	}
	if createdState != nil {
		state.LastUpdated = createdState.LastUpdated
		state.CreatedAt = createdState.CreatedAt
	}

	return err
}

// GetByID retrieves a room state by its ID
func (r *GlobalStateRoomRepository) GetByID(id int) (*GlobalStateRoom, error) {
	row := r.db.QueryRow(`
		SELECT id, room_id, player_count, npc_state_json, item_state_json, last_updated, created_at
		FROM global_state_rooms
		WHERE id = ?
	`, id)

	state := &GlobalStateRoom{}
	err := row.Scan(
		&state.ID,
		&state.RoomID,
		&state.PlayerCount,
		&state.NPCStateJSON,
		&state.ItemStateJSON,
		&state.LastUpdated,
		&state.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return state, nil
}

// GetByRoomID retrieves a room state by its room ID
func (r *GlobalStateRoomRepository) GetByRoomID(roomID string) (*GlobalStateRoom, error) {
	row := r.db.QueryRow(`
		SELECT id, room_id, player_count, npc_state_json, item_state_json, last_updated, created_at
		FROM global_state_rooms
		WHERE room_id = ?
	`, roomID)

	state := &GlobalStateRoom{}
	err := row.Scan(
		&state.ID,
		&state.RoomID,
		&state.PlayerCount,
		&state.NPCStateJSON,
		&state.ItemStateJSON,
		&state.LastUpdated,
		&state.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return state, nil
}

// Update updates a room state
func (r *GlobalStateRoomRepository) Update(state *GlobalStateRoom) error {
	stmt, err := r.db.Prepare(`
		UPDATE global_state_rooms
		SET player_count = ?, npc_state_json = ?, item_state_json = ?, last_updated = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		state.PlayerCount,
		state.NPCStateJSON,
		state.ItemStateJSON,
		time.Now(),
		state.ID,
	)
	return err
}

// UpdatePlayerCount updates only the player count for a room
func (r *GlobalStateRoomRepository) UpdatePlayerCount(roomID string, count int) error {
	_, err := r.db.Exec(`
		UPDATE global_state_rooms
		SET player_count = ?, last_updated = ?
		WHERE room_id = ?
	`, count, time.Now(), roomID)
	return err
}

// UpdateNPCState updates only the NPC state for a room
func (r *GlobalStateRoomRepository) UpdateNPCState(roomID string, npcState []NPCState) error {
	npcStateJSON, err := json.Marshal(npcState)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE global_state_rooms
		SET npc_state_json = ?, last_updated = ?
		WHERE room_id = ?
	`, string(npcStateJSON), time.Now(), roomID)
	return err
}

// UpdateItemState updates only the item state for a room
func (r *GlobalStateRoomRepository) UpdateItemState(roomID string, itemState []ItemState) error {
	itemStateJSON, err := json.Marshal(itemState)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE global_state_rooms
		SET item_state_json = ?, last_updated = ?
		WHERE room_id = ?
	`, string(itemStateJSON), time.Now(), roomID)
	return err
}

// Delete deletes a room state by ID
func (r *GlobalStateRoomRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM global_state_rooms WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// DeleteByRoomID deletes a room state by room ID
func (r *GlobalStateRoomRepository) DeleteByRoomID(roomID string) error {
	stmt, err := r.db.Prepare("DELETE FROM global_state_rooms WHERE room_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roomID)
	return err
}

// GetAll retrieves all room states
func (r *GlobalStateRoomRepository) GetAll() ([]*GlobalStateRoom, error) {
	rows, err := r.db.Query(`
		SELECT id, room_id, player_count, npc_state_json, item_state_json, last_updated, created_at
		FROM global_state_rooms
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*GlobalStateRoom
	for rows.Next() {
		state := &GlobalStateRoom{}
		err := rows.Scan(
			&state.ID,
			&state.RoomID,
			&state.PlayerCount,
			&state.NPCStateJSON,
			&state.ItemStateJSON,
			&state.LastUpdated,
			&state.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}

	return states, nil
}

// InitializeRoomState initializes a new room state
func (r *GlobalStateRoomRepository) InitializeRoomState(roomID string) error {
	state := &GlobalStateRoom{
		RoomID:         roomID,
		PlayerCount:    0,
		NPCStateJSON:   "[]",
		ItemStateJSON:  "[]",
		LastUpdated:    time.Now(),
		CreatedAt:      time.Now(),
	}
	return r.Create(state)
}

// GetNPCState retrieves and parses the NPC state for a room
func (r *GlobalStateRoomRepository) GetNPCState(roomID string) ([]NPCState, error) {
	state, err := r.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return []NPCState{}, nil
	}

	var npcState []NPCState
	if state.NPCStateJSON != "" {
		if err := json.Unmarshal([]byte(state.NPCStateJSON), &npcState); err != nil {
			return nil, err
		}
	}

	return npcState, nil
}

// GetItemState retrieves and parses the item state for a room
func (r *GlobalStateRoomRepository) GetItemState(roomID string) ([]ItemState, error) {
	state, err := r.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return []ItemState{}, nil
	}

	var itemState []ItemState
	if state.ItemStateJSON != "" {
		if err := json.Unmarshal([]byte(state.ItemStateJSON), &itemState); err != nil {
			return nil, err
		}
	}

	return itemState, nil
}

// IncrementPlayerCount increments the player count for a room
func (r *GlobalStateRoomRepository) IncrementPlayerCount(roomID string) error {
	state, err := r.GetByRoomID(roomID)
	if err != nil {
		return err
	}
	if state == nil {
		return r.InitializeRoomState(roomID)
	}

	return r.UpdatePlayerCount(roomID, state.PlayerCount+1)
}

// DecrementPlayerCount decrements the player count for a room
func (r *GlobalStateRoomRepository) DecrementPlayerCount(roomID string) error {
	state, err := r.GetByRoomID(roomID)
	if err != nil {
		return err
	}
	if state == nil || state.PlayerCount <= 0 {
		return nil
	}

	return r.UpdatePlayerCount(roomID, state.PlayerCount-1)
}