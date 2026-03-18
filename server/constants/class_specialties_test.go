package constants

import (
	"testing"
)

func TestGetClassConfig_WarriorFighter(t *testing.T) {
	config := GetClassConfig("warrior", "fighter")

	if config.Class != "warrior" {
		t.Errorf("Expected class 'warrior', got '%s'", config.Class)
	}
	if config.Specialty != "fighter" {
		t.Errorf("Expected specialty 'fighter', got '%s'", config.Specialty)
	}

	// Check starting skills
	if config.StartingSkills["blades"] != 1 {
		t.Errorf("Expected blades level 1, got %d", config.StartingSkills["blades"])
	}
	if config.StartingSkills["brawling"] != 1 {
		t.Errorf("Expected brawling level 1, got %d", config.StartingSkills["brawling"])
	}

	// Check starting talents
	expectedTalents := []string{"slash", "parry", "shield_bash", "heavy_strike"}
	if len(config.StartingTalents) != len(expectedTalents) {
		t.Errorf("Expected %d talents, got %d", len(expectedTalents), len(config.StartingTalents))
	}
	for i, talent := range expectedTalents {
		if config.StartingTalents[i] != talent {
			t.Errorf("Expected talent[%d] '%s', got '%s'", i, talent, config.StartingTalents[i])
		}
	}

	// Check stat bonuses
	if config.StatBonuses.Strength != 3 {
		t.Errorf("Expected Strength bonus 3, got %d", config.StatBonuses.Strength)
	}
	if config.StatBonuses.Constitution != 1 {
		t.Errorf("Expected Constitution bonus 1, got %d", config.StatBonuses.Constitution)
	}
}

func TestGetClassConfig_FallbackToSurvivor(t *testing.T) {
	// Unknown class should fallback to survivor:generalist
	config := GetClassConfig("unknown_class", "unknown_specialty")

	if config.Class != "survivor" {
		t.Errorf("Expected fallback to survivor class, got '%s'", config.Class)
	}
	if config.Specialty != "generalist" {
		t.Errorf("Expected fallback to generalist specialty, got '%s'", config.Specialty)
	}
}

func TestGetClassConfig_ClassWithDefaultSpecialty(t *testing.T) {
	// Requesting just "warrior" should return first specialty (fighter)
	config := GetClassConfig("warrior", "")

	// Should still get fighter config
	if config.Class != "warrior" {
		t.Errorf("Expected class 'warrior', got '%s'", config.Class)
	}
	// Fighter is the first specialty for warrior
	if config.Specialty != "fighter" {
		t.Errorf("Expected specialty 'fighter', got '%s'", config.Specialty)
	}
}

func TestGetSpecialty_Warrior(t *testing.T) {
	specialty := GetSpecialty("warrior")
	if specialty != "fighter" {
		t.Errorf("Expected 'fighter' as default warrior specialty, got '%s'", specialty)
	}
}

func TestGetSpecialty_UnknownClass(t *testing.T) {
	specialty := GetSpecialty("unknown_class")
	if specialty != "generalist" {
		t.Errorf("Expected 'generalist' as fallback specialty, got '%s'", specialty)
	}
}

func TestClassSpecialties_WarriorHasFighter(t *testing.T) {
	specialties, ok := ClassSpecialties["warrior"]
	if !ok {
		t.Fatal("Expected warrior class to have specialties")
	}

	found := false
	for _, s := range specialties {
		if s.ID == "fighter" {
			found = true
			if s.Name != "Fighter" {
				t.Errorf("Expected Fighter name, got '%s'", s.Name)
			}
			break
		}
	}
	if !found {
		t.Error("Expected fighter specialty in warrior class")
	}
}

func TestClassSpecialties_AllClassesHaveSpecialties(t *testing.T) {
	classes := []string{"warrior", "chef", "mystic", "tinkerer", "trader", "brawler", "vine_climber", "survivor"}

	for _, class := range classes {
		specialties, ok := ClassSpecialties[class]
		if !ok {
			t.Errorf("Class '%s' missing from ClassSpecialties", class)
			continue
		}
		if len(specialties) == 0 {
			t.Errorf("Class '%s' has no specialties defined", class)
		}
	}
}

func TestStartingConfigs_WarriorFighterStatBonuses(t *testing.T) {
	config := GetClassConfig("warrior", "fighter")

	// Warrior fighter should have STR+3, CON+1
	if config.StatBonuses.Strength != 3 {
		t.Errorf("Expected Strength bonus 3, got %d", config.StatBonuses.Strength)
	}
	if config.StatBonuses.Constitution != 1 {
		t.Errorf("Expected Constitution bonus 1, got %d", config.StatBonuses.Constitution)
	}
	// DEX, INT, WIS should be 0
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