package items

import (
	"encoding/json"
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