package constants

import (
	"testing"
)

// Test that Chef Pizza Combat talents exist and are properly configured
func TestChefPizzaCombatTalentsExist(t *testing.T) {
	// These talents should be defined in dbinit/init.go
	// This test verifies the Chef class configuration includes them
	config := GetClassConfig("chef", "pizzaiolo")

	talentSet := make(map[string]bool)
	for _, talent := range config.StartingTalents {
		talentSet[talent] = true
	}

	// Verify at least some pizza combat talents are in starting talents
	pizzaTalents := []string{"dough_ball", "sauce_splash", "pizza_cutter_dash", "pizza_meteor"}
	foundPizzaTalents := 0
	for _, pt := range pizzaTalents {
		if talentSet[pt] {
			foundPizzaTalents++
		}
	}

	if foundPizzaTalents < 2 {
		t.Errorf("Expected at least 2 pizza combat talents in starting talents, found %d", foundPizzaTalents)
	}
}

// Test Chef skills exist
func TestChefSkillsExist(t *testing.T) {
	config := GetClassConfig("chef", "pizzaiolo")

	// Chef should have cooking and pizza_combat skills
	if config.StartingSkills["cooking"] != 1 {
		t.Errorf("Expected cooking skill level 1, got %d", config.StartingSkills["cooking"])
	}
	if config.StartingSkills["pizza_combat"] != 1 {
		t.Errorf("Expected pizza_combat skill level 1, got %d", config.StartingSkills["pizza_combat"])
	}
	if config.StartingSkills["foraging"] != 1 {
		t.Errorf("Expected foraging skill level 1, got %d", config.StartingSkills["foraging"])
	}
}

// Test Chef has correct stat distribution
func TestChefStatDistribution(t *testing.T) {
	config := GetClassConfig("chef", "pizzaiolo")

	// Chef should be DEX/INT focused (support + ranged)
	if config.StatBonuses.Dexterity != 2 {
		t.Errorf("Expected DEX bonus 2 for Chef, got %d", config.StatBonuses.Dexterity)
	}
	if config.StatBonuses.Intelligence != 2 {
		t.Errorf("Expected INT bonus 2 for Chef, got %d", config.StatBonuses.Intelligence)
	}
	// Should NOT have high STR (not a melee class)
	if config.StatBonuses.Strength > 1 {
		t.Errorf("Expected STR bonus <= 1 for Chef, got %d", config.StatBonuses.Strength)
	}
}