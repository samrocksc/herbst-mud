// Package constants holds game constants for character creation and validation.
package constants

// ValidClasses is the list of allowed character classes.
var ValidClasses = []string{
	"tinkerer",
	"trader",
	"warrior",
	"brawler",
	"mystic",
	"chef",
	"vine_climber",
	"survivor",
}

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

// ClassStatBonuses returns stat bonuses for each class.
var ClassStatBonuses = map[string]struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}{
	"tinkerer":   {Intelligence: 3, Dexterity: 1},
	"trader":     {Charisma: 3, Wisdom: 1},
	"warrior":   {Strength: 3, Constitution: 1},
	"brawler":   {Strength: 2, Dexterity: 2},
	"mystic":    {Wisdom: 3, Intelligence: 1},
	"chef":      {Constitution: 2, Charisma: 2},
	"vine_climber": {Dexterity: 3, Strength: 1},
	"survivor":  {Constitution: 2, Wisdom: 2},
}