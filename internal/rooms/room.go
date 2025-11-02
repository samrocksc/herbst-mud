package rooms

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
)

// Room represents a room in the game
type Room struct {
	ID               string
	Description      string
	Exits            map[Direction]string // Direction to Room ID
	ImmovableObjects []items.Item
	MovableObjects   []items.Item
	Smells           string
	NPCs             []characters.Character
}

// RoomJSON represents a room in JSON format
type RoomJSON struct {
	Schema           string                 `json:"$schema"`
	ID               string                 `json:"id"`
	Description      string                 `json:"description"`
	Exits            map[string]string      `json:"exits"`
	ImmovableObjects []items.Item           `json:"immovableObjects"`
	MovableObjects   []items.Item           `json:"movableObjects"`
	Smells           string                 `json:"smells"`
	NPCs             []characters.Character `json:"npcs"`
}

// Direction represents cardinal directions
type Direction string

const (
	North     Direction = "north"
	South     Direction = "south"
	East      Direction = "east"
	West      Direction = "west"
	Northeast Direction = "northeast"
	Northwest Direction = "northwest"
	Southeast Direction = "southeast"
	Southwest Direction = "southwest"
	Up        Direction = "up"
	Down      Direction = "down"
)

// LoadRoomFromJSON loads a room from a JSON file
func LoadRoomFromJSON(filename string) (*Room, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var roomJSON RoomJSON
	if err := json.Unmarshal(data, &roomJSON); err != nil {
		return nil, err
	}

	// Convert exits from string keys to Direction keys
	exits := make(map[Direction]string)
	for k, v := range roomJSON.Exits {
		exits[Direction(k)] = v
	}

	room := &Room{
		ID:               roomJSON.ID,
		Description:      roomJSON.Description,
		Exits:            exits,
		ImmovableObjects: roomJSON.ImmovableObjects,
		MovableObjects:   roomJSON.MovableObjects,
		Smells:           roomJSON.Smells,
		NPCs:             roomJSON.NPCs,
	}

	return room, nil
}

// LoadAllRoomsFromDirectory loads all rooms from JSON files in a directory
func LoadAllRoomsFromDirectory(directory string) (map[string]*Room, error) {
	rooms := make(map[string]*Room)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			room, err := LoadRoomFromJSON(filename)
			if err != nil {
				return nil, err
			}
			rooms[room.ID] = room
		}
	}

	return rooms, nil
}
