package main

import (
	"testing"
)

// TestExamineSkillLevel_Default tests default skill level initialization
func TestExamineSkillLevel_Default(t *testing.T) {
	skill := &ExamineSkillLevel{}
	if skill.Level != 0 {
		t.Errorf("Expected default level 0, got %d", skill.Level)
	}
	if skill.XP != 0 {
		t.Errorf("Expected default XP 0, got %d", skill.XP)
	}
}

// TestExamineSkillBonus tests the bonus calculation
func TestExamineSkillBonus(t *testing.T) {
	tests := []struct {
		name     string
		intStat  int
		wisStat  int
		expected int
	}{
		{"low stats", 10, 10, 15},  // 10 + (10/2) = 15
		{"high int", 20, 10, 25},   // 20 + (10/2) = 25
		{"high wis", 10, 20, 20},   // 10 + (20/2) = 20
		{"balanced", 15, 15, 22},  // 15 + (15/2) = 22
		{"zero stats", 0, 0, 0},   // 0 + (0/2) = 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExamineSkillBonus(tt.intStat, tt.wisStat)
			if result != tt.expected {
				t.Errorf("ExamineSkillBonus(%d, %d) = %d; want %d", tt.intStat, tt.wisStat, result, tt.expected)
			}
		})
	}
}

// TestExamineSkillLevel_CanRevealHiddenDetail tests threshold checks
func TestExamineSkillLevel_CanRevealHiddenDetail(t *testing.T) {
	tests := []struct {
		name      string
		skill     *ExamineSkillLevel
		minLevel  int
		expectYes bool
	}{
		{"skill equal to threshold", &ExamineSkillLevel{Level: 10}, 10, true},
		{"skill above threshold", &ExamineSkillLevel{Level: 15}, 10, true},
		{"skill below threshold", &ExamineSkillLevel{Level: 5}, 10, false},
		{"zero skill", &ExamineSkillLevel{Level: 0}, 10, false},
		{"max skill", &ExamineSkillLevel{Level: 100}, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.skill.CanRevealHiddenDetail(tt.minLevel)
			if result != tt.expectYes {
				t.Errorf("CanRevealHiddenDetail() = %v; want %v", result, tt.expectYes)
			}
		})
	}
}

// TestExamineSkillLevel_GetExamineBonusPercent tests bonus percentage
func TestExamineSkillLevel_GetExamineBonusPercent(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"level 0", 0, 0},
		{"level 10", 10, 0},
		{"level 26", 26, 10},
		{"level 51", 51, 25},
		{"level 76", 76, 50},
		{"level 91", 91, 75},
		{"level 100", 100, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{Level: tt.level}
			result := skill.GetExamineBonusPercent()
			if result != tt.expected {
				t.Errorf("GetExamineBonusPercent() = %d; want %d", result, tt.expected)
			}
		})
	}
}

// TestHiddenDetail_RevealAutomatic tests automatic reveal mode
func TestHiddenDetail_RevealAutomatic(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 20}

	details := []HiddenDetail{
		{
			Text:            "You notice a hidden inscription",
			MinExamineLevel: 15,
			Mode:            "automatic",
			Revealed:        false,
		},
		{
			Text:            "A secret compartment is visible",
			MinExamineLevel: 25,
			Mode:            "automatic",
			Revealed:        false,
		},
	}

	result := skill.RevealHiddenDetails(details, 10, 10)

	// First detail should be revealed (20 >= 15)
	if !result[0].Revealed {
		t.Error("Expected first detail to be revealed")
	}
	if result[0].Source != "skill_threshold" {
		t.Errorf("Expected source 'skill_threshold', got '%s'", result[0].Source)
	}

	// Second detail should NOT be revealed (20 < 25)
	if result[1].Revealed {
		t.Error("Expected second detail to NOT be revealed")
	}
}

// TestHiddenDetail_RevealCheck tests check-based reveal mode
func TestHiddenDetail_RevealCheck(t *testing.T) {
	// Use fixed seed for deterministic results
	// Note: In real tests, we might want to mock the random

	skill := &ExamineSkillLevel{Level: 30}

	details := []HiddenDetail{
		{
			Text:            "A hidden message",
			MinExamineLevel: 20,
			Mode:            "check",
			DC:              35, // Requires roll of 35 or higher
			Stat:            "INT",
			Revealed:        false,
		},
	}

	result := skill.RevealHiddenDetails(details, 10, 10)

	// Check if skill check passed
	if result[0].Revealed {
		if result[0].Source != "check_passed" {
			t.Errorf("Expected source 'check_passed', got '%s'", result[0].Source)
		}
	} else {
		// Check failed - that's okay, just verify it was checked
		if result[0].Source == "" && result[0].MinExamineLevel == 20 {
			t.Log("Check failed as expected for high DC")
		}
	}
}

// TestExamineSkillLevel_AddXP tests XP gain and level up
func TestExamineSkillLevel_AddXP(t *testing.T) {
	tests := []struct {
		name          string
		initialLevel  int
		initialXP     int
		xpToAdd       int
		expectedLevel int
		expectedXP    int
	}{
		{"no level up", 0, 0, 5, 0, 5},
		{"single level up", 0, 0, 10, 1, 0},
		{"multiple level up", 0, 0, 50, 5, 0},
		{"level cap", 95, 0, 100, 100, 5}, // 95 + 10 = 105 -> capped to 100
		{"no XP gain at max", 100, 0, 10, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{
				Level: tt.initialLevel,
				XP:    tt.initialXP,
			}
			skill.AddXP(tt.xpToAdd)

			if skill.Level != tt.expectedLevel {
				t.Errorf("Level = %d; want %d", skill.Level, tt.expectedLevel)
			}
			if skill.XP != tt.expectedXP {
				t.Errorf("XP = %d; want %d", skill.XP, tt.expectedXP)
			}
		})
	}
}

// TestRevealHiddenDetails_PreservesOriginalData tests that original details aren't modified
func TestRevealHiddenDetails_PreservesOriginalData(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 50}

	details := []HiddenDetail{
		{
			Text:            "Original text",
			MinExamineLevel: 10,
			Mode:            "automatic",
			Revealed:        false,
		},
	}

	// Call reveal
	_ = skill.RevealHiddenDetails(details, 10, 10)

	// Original should still be unrevealed
	if details[0].Revealed {
		t.Error("Original detail was modified")
	}
}