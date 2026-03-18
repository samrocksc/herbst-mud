// Package constants holds game constants for character creation and validation.
package constants

// ClassSpecialty represents a specialty within a class.
type ClassSpecialty struct {
	ID          string
	Name        string
	Description string
}

// ClassConfig holds the starting configuration for a class.
type ClassConfig struct {
	Class              string
	Specialty          string
	Description        string
	StartingSkills     map[string]int  // skill_id -> level
	StartingTalents    []string        // talent_ids (max 4)
	StatBonuses        struct {
		Strength     int
		Dexterity    int
		Constitution int
		Intelligence int
		Wisdom       int
	}
}

// ClassSpecialties maps classes to their available specialties.
var ClassSpecialties = map[string][]ClassSpecialty{
	"warrior": {
		{ID: "fighter", Name: "Fighter", Description: "The bread-and-butter warrior. Solid, reliable combat with blades and brawling."},
		{ID: "knight", Name: "Knight", Description: "Heavy armor specialist with defensive tactics."},
		{ID: "berserker", Name: "Berserker", Description: "Raw aggression and devastating offense."},
	},
	"chef": {
		{ID: "pizzaiolo", Name: "Pizzaiolo", Description: "Master of pizza combat and mutant cuisine."},
	},
	"mystic": {
		{ID: "elementalist", Name: "Elementalist", Description: "Wielder of elemental magic."},
	},
	"tinkerer": {
		{ID: "mechanic", Name: "Mechanic", Description: "Tech salvage and gadget specialist."},
	},
	"trader": {
		{ID: "merchant", Name: "Merchant", Description: "Economy and NPC interactions."},
	},
	"brawler": {
		{ID: "street_fighter", Name: "Street Fighter", Description: "Unarmed combat specialist."},
	},
	"vine_climber": {
		{ID: "scout", Name: "Scout", Description: "Stealth and climbing."},
	},
	"survivor": {
		{ID: "generalist", Name: "Generalist", Description: "Jack of all trades."},
	},
}

// StartingConfigs maps class+specialty to starting configuration.
var StartingConfigs = map[string]ClassConfig{
	"warrior:fighter": {
		Class:        "warrior",
		Specialty:    "fighter",
		Description:  "Defenders of survivor enclaves, battle-hardened veterans.",
		StartingSkills: map[string]int{
			"blades":   1,
			"brawling": 1,
		},
		StartingTalents: []string{"slash", "parry", "shield_bash", "heavy_strike"},
		StatBonuses: struct {
			Strength     int
			Dexterity    int
			Constitution int
			Intelligence int
			Wisdom       int
		}{
			Strength:     3,
			Constitution: 1,
		},
	},
	"survivor:generalist": {
		Class:        "survivor",
		Specialty:    "generalist",
		Description:  "A versatile survivor who can adapt to any situation.",
		StartingSkills: map[string]int{
			"brawling": 1,
		},
		StartingTalents: []string{"crash"},
		StatBonuses: struct {
			Strength     int
			Dexterity    int
			Constitution int
			Intelligence int
			Wisdom       int
		}{
			Constitution: 2,
			Wisdom:       2,
		},
	},
}

// GetClassConfig returns the starting configuration for a class+specialty.
// Falls back to survivor:generalist if class not found.
func GetClassConfig(class, specialty string) ClassConfig {
	key := class + ":" + specialty
	if config, ok := StartingConfigs[key]; ok {
		return config
	}
	// Try just class with default specialty
	if specialties, ok := ClassSpecialties[class]; ok && len(specialties) > 0 {
		key = class + ":" + specialties[0].ID
		if config, ok := StartingConfigs[key]; ok {
			return config
		}
	}
	// Fallback to survivor
	return StartingConfigs["survivor:generalist"]
}

// GetSpecialty returns the specialty for a class, or the first available one.
func GetSpecialty(class string) string {
	if specialties, ok := ClassSpecialties[class]; ok && len(specialties) > 0 {
		return specialties[0].ID
	}
	return "generalist"
}