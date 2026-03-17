package constants

import (
	"testing"
)

func TestWarriorFighterStartingConfig(t *testing.T) {
	// Test that warrior class gets fighter specialty by default
	config := GetClassConfig("warrior", "")

	if config.Class != "warrior" {
		t.Errorf("Expected class 'warrior', got '%s'", config.Class)
	}
	if config.Specialty != "fighter" {
		t.Errorf("Expected specialty 'fighter', got '%s'", config.Specialty)
	}

	// Test starting skills
	if config.StartingSkills["blades"] != 1 {
		t.Errorf("Expected blades level 1, got %d", config.StartingSkills["blades"])
	}
	if config.StartingSkills["brawling"] != 1 {
		t.Errorf("Expected brawling level 1, got %d", config.StartingSkills["brawling"])
	}

	// Test starting talents
	expectedTalents := []string{"slash", "parry", "smash", "crash"}
	if len(config.StartingTalents) != len(expectedTalents) {
		t.Errorf("Expected %d talents, got %d", len(expectedTalents), len(config.StartingTalents))
	}
	for i, talent := range expectedTalents {
		if config.StartingTalents[i] != talent {
			t.Errorf("Expected talent[%d] '%s', got '%s'", i, talent, config.StartingTalents[i])
		}
	}

	// Test stat bonuses
	if config.StatBonuses.Strength != 3 {
		t.Errorf("Expected Strength bonus 3, got %d", config.StatBonuses.Strength)
	}
	if config.StatBonuses.Constitution != 1 {
		t.Errorf("Expected Constitution bonus 1, got %d", config.StatBonuses.Constitution)
	}
	if config.StatBonuses.Dexterity != 0 {
		t.Errorf("Expected Dexterity bonus 0, got %d", config.StatBonuses.Dexterity)
	}
	if config.StatBonuses.Intelligence != 0 {
		t.Errorf("Expected Intelligence bonus 0, got %d", config.StatBonuses.Intelligence)
	}
	if config.StatBonuses.Wisdom != 0 {
		t.Errorf("Expected Wisdom bonus 0, got %d", config.StatBonuses.Wisdom)
	}
}

func TestWarriorFighterExplicitSpecialty(t *testing.T) {
	// Test that explicitly requesting fighter specialty works
	config := GetClassConfig("warrior", "fighter")

	if config.Class != "warrior" {
		t.Errorf("Expected class 'warrior', got '%s'", config.Class)
	}
	if config.Specialty != "fighter" {
		t.Errorf("Expected specialty 'fighter', got '%s'", config.Specialty)
	}
}

func TestWarriorOtherSpecialties(t *testing.T) {
	// Test that warrior has multiple specialties defined
	specialties, ok := ClassSpecialties["warrior"]
	if !ok {
		t.Fatal("Expected warrior class to have specialties defined")
	}

	// Should have at least fighter (more specialties like knight, berserker planned)
	if len(specialties) < 1 {
		t.Error("Warrior should have at least one specialty")
	}

	// First specialty should be fighter
	if specialties[0].ID != "fighter" {
		t.Errorf("Expected first specialty 'fighter', got '%s'", specialties[0].ID)
	}
}