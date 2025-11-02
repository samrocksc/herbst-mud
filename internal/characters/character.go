package characters

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sam/makeathing/internal/items"
)

// Character represents a character in the game
type Character struct {
	ID          string
	Name        string
	Race        Race
	Class       Class
	Stats       Stats
	Health      int
	Mana        int
	Experience  int
	Level       int
	IsVendor    bool
	IsNpc       bool
	Inventory   []items.Item
	Skills      []Skill
}

// CharacterJSON represents a character in JSON format
type CharacterJSON struct {
	Schema     string      `json:"$schema"`
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Race       string      `json:"race"`
	Class      string      `json:"class"`
	Stats      Stats       `json:"stats"`
	Health     int         `json:"health"`
	Mana       int         `json:"mana"`
	Experience int         `json:"experience"`
	Level      int         `json:"level"`
	IsVendor   bool        `json:"isVendor"`
	IsNpc      bool        `json:"isNpc"`
	Inventory  []items.Item `json:"inventory"`
	Skills     []Skill     `json:"skills"`
}

// Race represents a character's race
type Race string

const (
	Human     Race = "Human"
	RatPeople Race = "Rat People"
	Dwarf     Race = "Dwarf"
	Dog       Race = "Dog"
)

// Class represents a character's class
type Class string

const (
	Warrior Class = "Warrior"
	Mage    Class = "Mage"
	Rogue   Class = "Rogue"
)

// Stats represents a character's stats
type Stats struct {
	Strength     int // 1-25, affects melee damage
	Intelligence int // 1-25, affects magic damage taken
	Dexterity    int // 1-25, affects hit chance and speed
}

// Skill represents a character's skills
type Skill struct {
	Name        string
	Type        SkillType
	Description string
}

// SkillType represents the type of skill
type SkillType string

const (
	RaceBased  SkillType = "race-based"
	ClassSkill SkillType = "class-skill"
	Spell      SkillType = "spell"
)

// LoadCharacterJSONFromJSON loads a CharacterJSON from a JSON file
func LoadCharacterJSONFromJSON(filename string) (*CharacterJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var characterJSON CharacterJSON
	if err := json.Unmarshal(data, &characterJSON); err != nil {
		return nil, err
	}

	return &characterJSON, nil
}

// LoadAllCharacterJSONsFromDirectory loads all CharacterJSONs from JSON files in a directory
func LoadAllCharacterJSONsFromDirectory(directory string) (map[string]*CharacterJSON, error) {
	characters := make(map[string]*CharacterJSON)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			characterJSON, err := LoadCharacterJSONFromJSON(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load character JSON from %s: %w", filename, err)
			}
			// Use the character ID as the key
			characters[characterJSON.ID] = characterJSON
		}
	}

	return characters, nil
}

// LoadCharacterFromJSON loads a character from a JSON file
func LoadCharacterFromJSON(filename string) (*Character, error) {
	characterJSON, err := LoadCharacterJSONFromJSON(filename)
	if err != nil {
		return nil, err
	}

	character := &Character{
		ID:          characterJSON.ID,
		Name:        characterJSON.Name,
		Race:        Race(characterJSON.Race),
		Class:       Class(characterJSON.Class),
		Stats:       characterJSON.Stats,
		Health:      characterJSON.Health,
		Mana:        characterJSON.Mana,
		Experience:  characterJSON.Experience,
		Level:       characterJSON.Level,
		IsVendor:    characterJSON.IsVendor,
		IsNpc:       characterJSON.IsNpc,
		Inventory:   characterJSON.Inventory,
		Skills:      characterJSON.Skills,
	}

	return character, nil
}

// LoadAllCharactersFromDirectory loads all characters from JSON files in a directory
func LoadAllCharactersFromDirectory(directory string) (map[string]*Character, error) {
	characters := make(map[string]*Character)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			character, err := LoadCharacterFromJSON(filename)
			if err != nil {
				return nil, err
			}
			characters[character.ID] = character
		}
	}

	return characters, nil
}

// FindCharacterByID finds a character by its ID in a map of characters
func FindCharacterByID(characters map[string]*Character, id string) *Character {
	if character, ok := characters[id]; ok {
		return character
	}
	return nil
}