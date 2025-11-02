package items

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Item represents an item in the game
type Item struct {
	ID          string
	Name        string
	Description string
	Type        ItemType
	Stats       ItemStats
	IsMagical   bool
}

// ItemType represents the type of item
type ItemType string

const (
	Weapon    ItemType = "weapon"
	Wearable  ItemType = "wearable"
	Movable   ItemType = "movable"
	Immovable ItemType = "immovable"
)

// ItemStats represents stats for an item
type ItemStats struct {
	Strength     int
	Intelligence int
	Dexterity    int
}

// ItemJSON represents an item in JSON format
type ItemJSON struct {
	Schema      string    `json:"$schema"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Stats       ItemStats `json:"stats"`
	IsMagical   bool      `json:"isMagical"`
}

// LoadItemJSONFromJSON loads an ItemJSON from a JSON file
func LoadItemJSONFromJSON(filename string) (*ItemJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var itemJSON ItemJSON
	if err := json.Unmarshal(data, &itemJSON); err != nil {
		return nil, err
	}

	return &itemJSON, nil
}

// LoadAllItemJSONsFromDirectory loads all ItemJSONs from JSON files in a directory
func LoadAllItemJSONsFromDirectory(directory string) (map[string]*ItemJSON, error) {
	items := make(map[string]*ItemJSON)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			itemJSON, err := LoadItemJSONFromJSON(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load item JSON from %s: %w", filename, err)
			}
			// Use the item ID as the key
			items[itemJSON.ID] = itemJSON
		}
	}

	return items, nil
}

// LoadItemFromJSON loads an item from a JSON file
func LoadItemFromJSON(filename string) (*Item, error) {
	itemJSON, err := LoadItemJSONFromJSON(filename)
	if err != nil {
		return nil, err
	}

	item := &Item{
		ID:          itemJSON.ID,
		Name:        itemJSON.Name,
		Description: itemJSON.Description,
		Type:        ItemType(itemJSON.Type),
		Stats:       itemJSON.Stats,
		IsMagical:   itemJSON.IsMagical,
	}

	return item, nil
}

// LoadAllItemsFromDirectory loads all items from JSON files in a directory
func LoadAllItemsFromDirectory(directory string) (map[string]*Item, error) {
	items := make(map[string]*Item)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			item, err := LoadItemFromJSON(filename)
			if err != nil {
				return nil, err
			}

			// Store the item by its ID
			items[item.ID] = item
		}
	}

	return items, nil
}

// FindItemByID finds an item by its ID in a map of items
func FindItemByID(items map[string]*Item, id string) *Item {
	if item, ok := items[id]; ok {
		return item
	}
	return nil
}