package rooms

import (
	"encoding/json"
	"fmt"
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

// RoomWithReferences represents a room that needs to resolve references
type RoomWithReferences struct {
	ID               string               `json:"id"`
	Description      string               `json:"description"`
	Exits            map[string]string    `json:"exits"`
	ImmovableObjects []ItemReference      `json:"immovableObjects"`
	MovableObjects   []ItemReference      `json:"movableObjects"`
	Smells           string               `json:"smells"`
	NPCs             []CharacterReference `json:"npcs"`
}

// ItemReference represents a reference to an item by ID
type ItemReference struct {
	ID string `json:"id"`
}

// CharacterReference represents a reference to a character by ID
type CharacterReference struct {
	ID string `json:"id"`
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

// LoadRoomFromJSONWithReferences loads a room from a JSON file and resolves item/character references
func LoadRoomFromJSONWithReferences(filename string, itemsMap map[string]*items.Item, charactersMap map[string]*characters.Character) (*Room, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var roomWithRefs RoomWithReferences
	if err := json.Unmarshal(data, &roomWithRefs); err != nil {
		// If JSON unmarshaling as RoomWithReferences fails, try unmarshaling as RoomJSON (backward compatibility)
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

	// Resolve item references
	immovableObjects := make([]items.Item, 0)
	for _, ref := range roomWithRefs.ImmovableObjects {
		// First, check if it's a reference by ID
		if item := items.FindItemByID(itemsMap, ref.ID); item != nil {
			// Create a copy of the item with the reference ID
			itemCopy := *item
			itemCopy.ID = ref.ID
			immovableObjects = append(immovableObjects, itemCopy)
		} else {
			// If no matching item is found in the items map, use the original reference object
			// This is for backward compatibility with inline item definitions
			immovableObjects = append(immovableObjects, items.Item{ID: ref.ID})
		}
	}

	movableObjects := make([]items.Item, 0)
	for _, ref := range roomWithRefs.MovableObjects {
		// First, check if it's a reference by ID
		if item := items.FindItemByID(itemsMap, ref.ID); item != nil {
			// Create a copy of the item with the reference ID
			itemCopy := *item
			itemCopy.ID = ref.ID
			movableObjects = append(movableObjects, itemCopy)
		} else {
			// If no matching item is found in the items map, use the original reference object
			// This is for backward compatibility with inline item definitions
			movableObjects = append(movableObjects, items.Item{ID: ref.ID})
		}
	}

	// Resolve character references
	npcs := make([]characters.Character, 0)
	for _, ref := range roomWithRefs.NPCs {
		// First, check if it's a reference by ID
		if char := characters.FindCharacterByID(charactersMap, ref.ID); char != nil {
			// Create a copy of the character with the reference ID
			charCopy := *char
			charCopy.ID = ref.ID
			npcs = append(npcs, charCopy)
		} else {
			// If no matching character is found in the characters map, use the original reference object
			// This is for backward compatibility with inline character definitions
			npcs = append(npcs, characters.Character{ID: ref.ID})
		}
	}

	// Convert exits from string keys to Direction keys
	exits := make(map[Direction]string)
	for k, v := range roomWithRefs.Exits {
		exits[Direction(k)] = v
	}

	room := &Room{
		ID:               roomWithRefs.ID,
		Description:      roomWithRefs.Description,
		Exits:            exits,
		ImmovableObjects: immovableObjects,
		MovableObjects:   movableObjects,
		Smells:           roomWithRefs.Smells,
		NPCs:             npcs,
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

// LoadRoomJSONFromJSON loads a RoomJSON from a JSON file
func LoadRoomJSONFromJSON(filename string) (*RoomJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var roomJSON RoomJSON
	if err := json.Unmarshal(data, &roomJSON); err != nil {
		return nil, err
	}

	return &roomJSON, nil
}

// LoadAllRoomJSONsFromDirectory loads all RoomJSONs from JSON files in a directory
func LoadAllRoomJSONsFromDirectory(directory string) (map[string]*RoomJSON, error) {
	rooms := make(map[string]*RoomJSON)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			roomJSON, err := LoadRoomJSONFromJSON(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load room JSON from %s: %w", filename, err)
			}
			// Use the room ID as the key
			rooms[roomJSON.ID] = roomJSON
		}
	}

	return rooms, nil
}

// LoadAllRoomsItemsAndCharactersWithReferences loads all rooms from JSON files in a directory with resolved item/character references
func LoadAllRoomsItemsAndCharactersWithReferences(roomDir string, itemsDir string, charactersDir string) (map[string]*Room, error) {
	// Load all items
	itemsMap, err := items.LoadAllItemsFromDirectory(itemsDir)
	if err != nil {
		return nil, err
	}

	// Load all characters
	charactersMap, err := characters.LoadAllCharactersFromDirectory(charactersDir)
	if err != nil {
		return nil, err
	}

	rooms := make(map[string]*Room)

	files, err := os.ReadDir(roomDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(roomDir, file.Name())
			room, err := LoadRoomFromJSONWithReferences(filename, itemsMap, charactersMap)
			if err != nil {
				return nil, err
			}
			rooms[room.ID] = room
		}
	}

	return rooms, nil
}
