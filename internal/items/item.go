package items

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Item represents an item in the game
type Item struct {
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

// LoadItemFromJSON loads an item from a JSON file
func LoadItemFromJSON(filename string) (*Item, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var item Item
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	return &item, nil
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

			// Use the item name as the key (converted to lowercase and spaces replaced with underscores)
			key := file.Name()[:len(file.Name())-5] // Remove .json extension
			items[key] = item
		}
	}

	return items, nil
}
