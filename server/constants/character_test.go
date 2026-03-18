package constants

import (
	"testing"
)

func TestValidClasses(t *testing.T) {
	// Test that ValidClasses is not empty and contains expected Phase 1 classes
	if len(ValidClasses) == 0 {
		t.Error("ValidClasses should not be empty")
	}

	// Phase 1 classes that must be present
	phase1Classes := []string{"warrior", "chef", "mystic"}
	classSet := make(map[string]bool)
	for _, c := range ValidClasses {
		classSet[c] = true
	}

	for _, c := range phase1Classes {
		if !classSet[c] {
			t.Errorf("Phase 1 class %s should be in ValidClasses", c)
		}
	}
}

func TestValidRaces(t *testing.T) {
	// Test that ValidRaces is not empty
	if len(ValidRaces) == 0 {
		t.Error("ValidRaces should not be empty")
	}

	// Core races that must be present
	coreRaces := []string{"human", "mutant"}
	raceSet := make(map[string]bool)
	for _, r := range ValidRaces {
		raceSet[r] = true
	}

	for _, r := range coreRaces {
		if !raceSet[r] {
			t.Errorf("Core race %s should be in ValidRaces", r)
		}
	}
}

func TestDefaultStats(t *testing.T) {
	if DefaultStats.Strength != 10 {
		t.Errorf("Expected default strength 10, got %d", DefaultStats.Strength)
	}
	if DefaultStats.Dexterity != 10 {
		t.Errorf("Expected default dexterity 10, got %d", DefaultStats.Dexterity)
	}
	if DefaultStats.Constitution != 10 {
		t.Errorf("Expected default constitution 10, got %d", DefaultStats.Constitution)
	}
	if DefaultStats.Intelligence != 10 {
		t.Errorf("Expected default intelligence 10, got %d", DefaultStats.Intelligence)
	}
	if DefaultStats.Wisdom != 10 {
		t.Errorf("Expected default wisdom 10, got %d", DefaultStats.Wisdom)
	}
	if DefaultStats.Charisma != 10 {
		t.Errorf("Expected default charisma 10, got %d", DefaultStats.Charisma)
	}
	if DefaultStats.Hitpoints != 100 {
		t.Errorf("Expected default hitpoints 100, got %d", DefaultStats.Hitpoints)
	}
	if DefaultStats.MaxHitpoints != 100 {
		t.Errorf("Expected default max hitpoints 100, got %d", DefaultStats.MaxHitpoints)
	}
}

func TestClassStatBonuses(t *testing.T) {
	// Test that all ValidClasses have stat bonuses defined
	for _, class := range ValidClasses {
		if _, ok := ClassStatBonuses[class]; !ok {
			t.Errorf("Missing stat bonuses for class %s", class)
		}
	}

	// Test specific class bonus examples
	tinkererBonus := ClassStatBonuses["tinkerer"]
	if tinkererBonus.Intelligence != 3 {
		t.Errorf("Expected tinkerer to have +3 Intelligence, got %d", tinkererBonus.Intelligence)
	}
	if tinkererBonus.Dexterity != 1 {
		t.Errorf("Expected tinkerer to have +1 Dexterity, got %d", tinkererBonus.Dexterity)
	}

	// Test warrior bonus
	warriorBonus := ClassStatBonuses["warrior"]
	if warriorBonus.Strength != 3 {
		t.Errorf("Expected warrior to have +3 Strength, got %d", warriorBonus.Strength)
	}
}