package database

import (
	"database/sql"
	"time"
)

// GlobalStateCharacter represents a character state in the database
type GlobalStateCharacter struct {
	ID               int       `json:"id"`
	CharacterID      string    `json:"character_id"`
	CurrentRoomID    string    `json:"current_room_id"`
	Health           int       `json:"health"`
	Status           string    `json:"status"`
	LastUpdated      time.Time `json:"last_updated"`
	CreatedAt        time.Time `json:"created_at"`
}

// GlobalStateCharacterRepository provides methods for working with character state
type GlobalStateCharacterRepository struct {
	db *DB
}

// NewGlobalStateCharacterRepository creates a new character state repository
func NewGlobalStateCharacterRepository(db *DB) *GlobalStateCharacterRepository {
	return &GlobalStateCharacterRepository{db: db}
}

// Create creates a new character state
func (r *GlobalStateCharacterRepository) Create(state *GlobalStateCharacter) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO global_state_characters (character_id, current_room_id, health, status)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		state.CharacterID,
		state.CurrentRoomID,
		state.Health,
		state.Status,
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

// GetByID retrieves a character state by its ID
func (r *GlobalStateCharacterRepository) GetByID(id int) (*GlobalStateCharacter, error) {
	row := r.db.QueryRow(`
		SELECT id, character_id, current_room_id, health, status, last_updated, created_at
		FROM global_state_characters
		WHERE id = ?
	`, id)

	state := &GlobalStateCharacter{}
	err := row.Scan(
		&state.ID,
		&state.CharacterID,
		&state.CurrentRoomID,
		&state.Health,
		&state.Status,
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

// GetByCharacterID retrieves a character state by its character ID
func (r *GlobalStateCharacterRepository) GetByCharacterID(characterID string) (*GlobalStateCharacter, error) {
	row := r.db.QueryRow(`
		SELECT id, character_id, current_room_id, health, status, last_updated, created_at
		FROM global_state_characters
		WHERE character_id = ?
	`, characterID)

	state := &GlobalStateCharacter{}
	err := row.Scan(
		&state.ID,
		&state.CharacterID,
		&state.CurrentRoomID,
		&state.Health,
		&state.Status,
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

// Update updates a character state
func (r *GlobalStateCharacterRepository) Update(state *GlobalStateCharacter) error {
	stmt, err := r.db.Prepare(`
		UPDATE global_state_characters
		SET current_room_id = ?, health = ?, status = ?, last_updated = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		state.CurrentRoomID,
		state.Health,
		state.Status,
		time.Now(),
		state.ID,
	)
	return err
}

// UpdateRoom updates only the room for a character
func (r *GlobalStateCharacterRepository) UpdateRoom(characterID, roomID string) error {
	_, err := r.db.Exec(`
		UPDATE global_state_characters
		SET current_room_id = ?, last_updated = ?
		WHERE character_id = ?
	`, roomID, time.Now(), characterID)
	return err
}

// UpdateHealth updates only the health for a character
func (r *GlobalStateCharacterRepository) UpdateHealth(characterID string, health int) error {
	_, err := r.db.Exec(`
		UPDATE global_state_characters
		SET health = ?, last_updated = ?
		WHERE character_id = ?
	`, health, time.Now(), characterID)
	return err
}

// UpdateStatus updates only the status for a character
func (r *GlobalStateCharacterRepository) UpdateStatus(characterID, status string) error {
	_, err := r.db.Exec(`
		UPDATE global_state_characters
		SET status = ?, last_updated = ?
		WHERE character_id = ?
	`, status, time.Now(), characterID)
	return err
}

// Delete deletes a character state by ID
func (r *GlobalStateCharacterRepository) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM global_state_characters WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// DeleteByCharacterID deletes a character state by character ID
func (r *GlobalStateCharacterRepository) DeleteByCharacterID(characterID string) error {
	stmt, err := r.db.Prepare("DELETE FROM global_state_characters WHERE character_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(characterID)
	return err
}

// GetAll retrieves all character states
func (r *GlobalStateCharacterRepository) GetAll() ([]*GlobalStateCharacter, error) {
	rows, err := r.db.Query(`
		SELECT id, character_id, current_room_id, health, status, last_updated, created_at
		FROM global_state_characters
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*GlobalStateCharacter
	for rows.Next() {
		state := &GlobalStateCharacter{}
		err := rows.Scan(
			&state.ID,
			&state.CharacterID,
			&state.CurrentRoomID,
			&state.Health,
			&state.Status,
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

// GetCharactersInRoom retrieves all character states in a specific room
func (r *GlobalStateCharacterRepository) GetCharactersInRoom(roomID string) ([]*GlobalStateCharacter, error) {
	rows, err := r.db.Query(`
		SELECT id, character_id, current_room_id, health, status, last_updated, created_at
		FROM global_state_characters
		WHERE current_room_id = ?
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*GlobalStateCharacter
	for rows.Next() {
		state := &GlobalStateCharacter{}
		err := rows.Scan(
			&state.ID,
			&state.CharacterID,
			&state.CurrentRoomID,
			&state.Health,
			&state.Status,
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

// InitializeCharacterState initializes a new character state
func (r *GlobalStateCharacterRepository) InitializeCharacterState(characterID, roomID string, health int) error {
	state := &GlobalStateCharacter{
		CharacterID:   characterID,
		CurrentRoomID: roomID,
		Health:        health,
		Status:        "active",
		LastUpdated:   time.Now(),
		CreatedAt:     time.Now(),
	}
	return r.Create(state)
}