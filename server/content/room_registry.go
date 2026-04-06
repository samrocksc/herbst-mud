package content

import (
	"fmt"
	"strings"
	"sync"
)

// RoomRegistry manages room definitions
type RoomRegistry struct {
	mu    sync.RWMutex
	rooms map[string]*RoomDef
}

// NewRoomRegistry creates a new room registry
func NewRoomRegistry() *RoomRegistry {
	return &RoomRegistry{
		rooms: make(map[string]*RoomDef),
	}
}

// RoomDef represents a room definition
type RoomDef struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	AreaID      string            `json:"area_id,omitempty"`
	Flags       []string          `json:"flags,omitempty"`
	Exits       map[string]string `json:"exits,omitempty"` // direction -> room_id
	Items       []string          `json:"items,omitempty"`     // item_ids
	NPCs        []string          `json:"npcs,omitempty"`     // npc_template_ids
}

// AreaDef represents an area containing rooms
type AreaDef struct {
	AreaID  string    ``
	Name    string    ``
	Rooms   []RoomDef ``
	Exits   []AreaExitDef ``
}

// AreaExitDef represents exits between areas
type AreaExitDef struct {
	From        string ``
	To          string ``
	Direction   string ``
	Description string ``
}

// Register adds a room to the registry
func (r *RoomRegistry) Register(room *RoomDef) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if room.ID == "" {
		return fmt.Errorf("room ID cannot be empty")
	}
	
	id := strings.ToLower(room.ID)
	if _, exists := r.rooms[id]; exists {
		return fmt.Errorf("room '%s' already registered", id)
	}
	
	r.rooms[id] = room
	return nil
}

// Get retrieves a room by ID
func (r *RoomRegistry) Get(id string) (*RoomDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	room, exists := r.rooms[strings.ToLower(id)]
	return room, exists
}

// GetByArea returns rooms in an area
func (r *RoomRegistry) GetByArea(areaID string) []*RoomDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*RoomDef
	areaID = strings.ToLower(areaID)
	
	for _, room := range r.rooms {
		if strings.ToLower(room.AreaID) == areaID {
			result = append(result, room)
		}
	}
	return result
}

// GetConnected returns rooms connected via exits
func (r *RoomRegistry) GetConnected(roomID string) []*RoomDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	room, exists := r.rooms[strings.ToLower(roomID)]
	if !exists {
		return nil
	}
	
	var result []*RoomDef
	for _, targetID := range room.Exits {
		if target, exists := r.rooms[strings.ToLower(targetID)]; exists {
			result = append(result, target)
		}
	}
	return result
}

// Clear removes all rooms
func (r *RoomRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms = make(map[string]*RoomDef)
}

// Count returns number of registered rooms
func (r *RoomRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.rooms)
}

// Validate checks all rooms
func (r *RoomRegistry) Validate() []ValidationError {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var errors []ValidationError
	
	for id, room := range r.rooms {
		if room.Name == "" {
			errors = append(errors, ValidationError{
				Type:    "room",
				ID:      id,
				Field:   "name",
				Message: "name is required",
			})
		}
		
		// Check exit targets exist (self-validated)
		for dir, targetID := range room.Exits {
			if _, exists := r.rooms[strings.ToLower(targetID)]; !exists {
				errors = append(errors, ValidationError{
					Type:    "room",
					ID:      id,
					Field:   fmt.Sprintf("exits.%s", dir),
					Message: fmt.Sprintf("target room '%s' not found", targetID),
				})
			}
		}
	}
	
	return errors
}

// GetAll returns all rooms
func (r *RoomRegistry) GetAll() []*RoomDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*RoomDef, 0, len(r.rooms))
	for _, room := range r.rooms {
		result = append(result, room)
	}
	return result
}
