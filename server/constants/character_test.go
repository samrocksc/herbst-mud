package constants

import (
	"testing"
)

func TestValidRaces(t *testing.T) {
	validRaces := []string{
		"human",
		"mutant",
		"android",
		"escaped_slave",
	}

	if len(validRaces) != len(ValidRaces) {
		t.Errorf("Expected %d valid races, got %d", len(validRaces), len(ValidRaces))
	}

	raceSet := make(map[string]bool)
	for _, r := range ValidRaces {
		raceSet[r] = true
	}

	for _, r := range validRaces {
		if !raceSet[r] {
			t.Errorf("Expected race %s to be valid", r)
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