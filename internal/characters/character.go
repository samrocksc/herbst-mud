package characters

import (
	"github.com/sam/makeathing/internal/items"
)

// Character represents a character in the game
type Character struct {
	Name       string
	Race       Race
	Class      Class
	Stats      Stats
	Health     int
	Mana       int
	Experience int
	Level      int
	IsVendor   bool
	IsNpc      bool
	Inventory  []items.Item
	Skills     []Skill
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
