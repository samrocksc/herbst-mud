package main

import (
	"testing"
)

// ChefClassConfig defines the Chef class configuration
type ChefClassConfig struct {
	Name         string
	StatsFocus   string // INT for mana, DEX for speed, CON for durability
	StartingHP   int
	StartingMana int
}

// ChefTalent represents a Chef-specific talent
type ChefTalent struct {
	ID          string
	Name        string
	Type        string // Attack, Defense, Buff, Utility
	MPCost      int
	Description string
	Damage      int // 0 for non-attacks
}

// ChefSkill represents a Chef skill category
type ChefSkill struct {
	ID          string
	Name        string
	PrimaryStat string // STR, DEX, INT, WIS
	Description string
}

// ChefTalentBook defines the chef's available talents
var ChefTalents = []ChefTalent{
	// Pizzaiolo (Pizza Combat)
	{ID: "dough_ball", Name: "Dough Ball", Type: "Attack", MPCost: 5, Description: "Flour-based ranged attack", Damage: 8},
	{ID: "sauce_splash", Name: "Sauce Splash", Type: "Attack", MPCost: 8, Description: "Hot sauce blinds enemies", Damage: 12},
	{ID: "pizza_cutter_dash", Name: "Pizza Cutter Dash", Type: "Attack", MPCost: 15, Description: "Spin attack with pizza cutter", Damage: 20},
	{ID: "pizza_meteor", Name: "Pizza Meteor", Type: "Ultimate", MPCost: 30, Description: "Ultimate giant pizza slam", Damage: 35},

	// Cooking (Potions & Buffs)
	{ID: "recipe_book", Name: "Recipe Book", Type: "Utility", MPCost: 10, Description: "Learn new dishes"},
	{ID: "mutant_seasoning", Name: "Mutant Seasoning", Type: "Buff", MPCost: 5, Description: "Ooze-tinged ingredients boost"},
	{ID: "serving_size", Name: "Serving Size", Type: "Buff", MPCost: 8, Description: "Buff more people at once"},

	// Foraging & Brewing
	{ID: "foraging", Name: "Foraging", Type: "Utility", MPCost: 3, Description: "Find mutant ingredients"},
	{ID: "brewing", Name: "Brewing", Type: "Utility", MPCost: 10, Description: "Create mutant potions"},
	{ID: "food_preservation", Name: "Food Preservation", Type: "Passive", MPCost: 0, Description: "Keep supplies fresh longer"},

	// The Pizzeria
	{ID: "open_pizza_stall", Name: "Open Pizza Stall", Type: "Utility", MPCost: 0, Description: "Sell slices for profit (passive income)"},
	{ID: "signature_pie", Name: "Signature Pie", Type: "Buff", MPCost: 20, Description: "Unique dish with bonus effects"},
	{ID: "food_fight", Name: "Food Fight", Type: "AoE", MPCost: 15, Description: "AoE attack with food projectiles", Damage: 18},
}

// ChefSkills defines available skill categories for Chef
var ChefSkills = []ChefSkill{
	{ID: "pizzaiolo", Name: "Pizzaiolo", PrimaryStat: "DEX", Description: "Pizza combat proficiency"},
	{ID: "cooking", Name: "Cooking", PrimaryStat: "INT", Description: "Alchemy via cooking"},
	{ID: "foraging", Name: "Foraging", PrimaryStat: "WIS", Description: "Finding mutant ingredients"},
	{ID: "brewing", Name: "Brewing", PrimaryStat: "INT", Description: "Creating mutant potions"},
}

// ChefClass is the Chef class configuration
var ChefClass = ChefClassConfig{
	Name:         "Chef",
	StatsFocus:   "INT",
	StartingHP:   80,
	StartingMana: 100,
}

func TestChefClassConfig(t *testing.T) {
	if ChefClass.Name != "Chef" {
		t.Errorf("Expected Chef class name, got %s", ChefClass.Name)
	}
	if ChefClass.StatsFocus != "INT" {
		t.Errorf("Expected INT stat focus, got %s", ChefClass.StatsFocus)
	}
	if ChefClass.StartingHP != 80 {
		t.Errorf("Expected 80 starting HP, got %d", ChefClass.StartingHP)
	}
	if ChefClass.StartingMana != 100 {
		t.Errorf("Expected 100 starting mana, got %d", ChefClass.StartingMana)
	}
}

func TestChefTalentsCount(t *testing.T) {
	if len(ChefTalents) != 13 {
		t.Errorf("Expected 13 Chef talents, got %d", len(ChefTalents))
	}
}

func TestChefTalentsHaveIDs(t *testing.T) {
	for _, talent := range ChefTalents {
		if talent.ID == "" {
			t.Errorf("Talent %s has empty ID", talent.Name)
		}
		if talent.Name == "" {
			t.Errorf("Talent with ID %s has empty Name", talent.ID)
		}
	}
}

func TestChefTalentTypes(t *testing.T) {
	validTypes := map[string]bool{
		"Attack":    true,
		"Defense":   true,
		"Buff":      true,
		"Utility":   true,
		"Passive":   true,
		"Ultimate":  true,
		"AoE":       true,
	}
	for _, talent := range ChefTalents {
		if !validTypes[talent.Type] {
			t.Errorf("Talent %s has invalid type: %s", talent.Name, talent.Type)
		}
	}
}

func TestChefSkillsCount(t *testing.T) {
	if len(ChefSkills) != 4 {
		t.Errorf("Expected 4 Chef skill categories, got %d", len(ChefSkills))
	}
}

func TestChefSkillsHaveValidStats(t *testing.T) {
	validStats := map[string]bool{
		"STR": true,
		"DEX": true,
		"INT": true,
		"WIS": true,
		"CON": true,
	}
	for _, skill := range ChefSkills {
		if !validStats[skill.PrimaryStat] {
			t.Errorf("Skill %s has invalid stat: %s", skill.Name, skill.PrimaryStat)
		}
	}
}

func TestPizzaCombatTalents(t *testing.T) {
	pizzaTalents := []string{"dough_ball", "sauce_splash", "pizza_cutter_dash", "pizza_meteor"}
	for _, id := range pizzaTalents {
		found := false
		for _, talent := range ChefTalents {
			if talent.ID == id {
				found = true
				if talent.Type != "Attack" && talent.Type != "Ultimate" {
					t.Errorf("Pizza talent %s should be Attack or Ultimate, got %s", id, talent.Type)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected pizza combat talent %s not found", id)
		}
	}
}

func TestChefClassHasPizzaCombat(t *testing.T) {
	pizzaTalents := 0
	for _, talent := range ChefTalents {
		if talent.Type == "Attack" || talent.Type == "Ultimate" || talent.Type == "AoE" {
			pizzaTalents++
		}
	}
	if pizzaTalents < 4 {
		t.Errorf("Chef class should have at least 4 attack talents, got %d", pizzaTalents)
	}
}