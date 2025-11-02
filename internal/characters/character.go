package characters

import (
	"encoding/json"
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

// LoadCharacterFromJSON loads a character from a JSON file
func LoadCharacterFromJSON(filename string) (*Character, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var character Character
	if err := json.Unmarshal(data, &character); err != nil {
		return nil, err
	}

	return &character, nil
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