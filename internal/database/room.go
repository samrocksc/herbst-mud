package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
	"github.com/sam/makeathing/internal/rooms"
)

// Room represents a room in the database
type Room struct {
	ID                   string                    `json:"id"`
	Description          string                    `json:"description"`
	Smells               string                    `json:"smells"`
	ExitsJSON            string                    `json:"exits_json"`              // JSON serialized map[string]string
	ImmovableObjectsJSON string                    `json:"immovable_objects_json"`  // JSON serialized []items.Item
	MovableObjectsJSON   string                    `json:"movable_objects_json"`    // JSON serialized []items.Item
	NPCsJSON             string                    `json:"npcs_json"`               // JSON serialized []characters.Character
	CreatedAt            time.Time                 `json:"created_at"`
	UpdatedAt            time.Time                 `json:"updated_at"`
}

// RoomRepository provides methods for working with rooms
type RoomRepository struct {
	db *DB
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create creates a new room
func (r *RoomRepository) Create(room *Room) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO rooms (id, description, smells, exits_json, immovable_objects_json, movable_objects_json, npcs_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		room.ID,
		room.Description,
		room.Smells,
		room.ExitsJSON,
		room.ImmovableObjectsJSON,
		room.MovableObjectsJSON,
		room.NPCsJSON,
		room.CreatedAt,
		room.UpdatedAt,
	)
	return err
}

// GetByID retrieves a room by its ID
func (r *RoomRepository) GetByID(id string) (*Room, error) {
	row := r.db.QueryRow(`
		SELECT id, description, smells, exits_json, immovable_objects_json, movable_objects_json, npcs_json, created_at, updated_at
		FROM rooms
		WHERE id = ?
	`, id)

	room := &Room{}
	err := row.Scan(
		&room.ID,
		&room.Description,
		&room.Smells,
		&room.ExitsJSON,
		&room.ImmovableObjectsJSON,
		&room.MovableObjectsJSON,
		&room.NPCsJSON,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return room, nil
}

// Update updates a room
func (r *RoomRepository) Update(room *Room) error {
	stmt, err := r.db.Prepare(`
		UPDATE rooms
		SET description = ?, smells = ?, exits_json = ?, immovable_objects_json = ?, movable_objects_json = ?, npcs_json = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		room.Description,
		room.Smells,
		room.ExitsJSON,
		room.ImmovableObjectsJSON,
		room.MovableObjectsJSON,
		room.NPCsJSON,
		time.Now(),
		room.ID,
	)
	return err
}

// Delete deletes a room by ID
func (r *RoomRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM rooms WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll retrieves all rooms
func (r *RoomRepository) GetAll() ([]*Room, error) {
	rows, err := r.db.Query(`
		SELECT id, description, smells, exits_json, immovable_objects_json, movable_objects_json, npcs_json, created_at, updated_at
		FROM rooms
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		room := &Room{}
		err := rows.Scan(
			&room.ID,
			&room.Description,
			&room.Smells,
			&room.ExitsJSON,
			&room.ImmovableObjectsJSON,
			&room.MovableObjectsJSON,
			&room.NPCsJSON,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// Helper functions to convert between JSON room and database room

// RoomFromJSONRoom converts a JSON room to a database room
func RoomFromJSONRoom(jsonRoom *rooms.RoomJSON) (*Room, error) {
	// Convert exits map to JSON
	exitsJSON, err := json.Marshal(jsonRoom.Exits)
	if err != nil {
		return nil, err
	}

	// Convert immovable objects to JSON
	immovableObjectsJSON, err := json.Marshal(jsonRoom.ImmovableObjects)
	if err != nil {
		return nil, err
	}

	// Convert movable objects to JSON
	movableObjectsJSON, err := json.Marshal(jsonRoom.MovableObjects)
	if err != nil {
		return nil, err
	}

	// Convert NPCs to JSON
	npcsJSON, err := json.Marshal(jsonRoom.NPCs)
	if err != nil {
		return nil, err
	}

	return &Room{
		ID:                   jsonRoom.ID,
		Description:          jsonRoom.Description,
		Smells:               jsonRoom.Smells,
		ExitsJSON:            string(exitsJSON),
		ImmovableObjectsJSON: string(immovableObjectsJSON),
		MovableObjectsJSON:   string(movableObjectsJSON),
		NPCsJSON:             string(npcsJSON),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

// ToJSONRoom converts a database room to a JSON room
func (r *Room) ToJSONRoom() (*rooms.RoomJSON, error) {
	// Convert exits JSON to map
	var exits map[string]string
	if r.ExitsJSON != "" {
		if err := json.Unmarshal([]byte(r.ExitsJSON), &exits); err != nil {
			return nil, err
		}
	}

	// Convert immovable objects JSON to slice
	var immovableObjects []items.Item
	if r.ImmovableObjectsJSON != "" {
		if err := json.Unmarshal([]byte(r.ImmovableObjectsJSON), &immovableObjects); err != nil {
			return nil, err
		}
	}

	// Convert movable objects JSON to slice
	var movableObjects []items.Item
	if r.MovableObjectsJSON != "" {
		if err := json.Unmarshal([]byte(r.MovableObjectsJSON), &movableObjects); err != nil {
			return nil, err
		}
	}

	// Convert NPCs JSON to slice
	var npcs []characters.Character
	if r.NPCsJSON != "" {
		if err := json.Unmarshal([]byte(r.NPCsJSON), &npcs); err != nil {
			return nil, err
		}
	}

	return &rooms.RoomJSON{
		Schema:           "../schemas/room.schema.json",
		ID:               r.ID,
		Description:      r.Description,
		Exits:            exits,
		ImmovableObjects: immovableObjects,
		MovableObjects:   movableObjects,
		Smells:           r.Smells,
		NPCs:             npcs,
	}, nil
}