// Package constants holds game constants for character creation and validation.
package constants

// ValidRaces is the list of allowed character races.
var ValidRaces = []string{
	"human",
	"mutant",
	"android",
	"escaped_slave",
}

// DefaultStats contains the default stat values for new characters.
var DefaultStats = struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
	Hitpoints    int
	MaxHitpoints int
}{
	Strength:     10,
	Dexterity:    10,
	Constitution: 10,
	Intelligence: 10,
	Wisdom:       10,
	Charisma:     10,
	Hitpoints:    100,
	MaxHitpoints: 100,
}