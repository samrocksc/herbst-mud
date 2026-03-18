package main

import (
	"testing"
)

// TestExamineSkillLevel_ExamineCheck tests the examine check formula
func TestExamineSkillLevel_ExamineCheck(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		intStat  int
		wisStat  int
		minScore int
		maxScore int
	}{
		{
			name:     "low skill level",
			level:    10,
			intStat:  10,
			wisStat:  10,
			minScore: 10 + 10 + 5 + 1, // level + INT + WIS/2 + random(1,10) = 26 at minimum
			maxScore: 10 + 10 + 5 + 10,
		},
		{
			name:     "high skill level",
			level:    50,
			intStat:  20,
			wisStat:  20,
			minScore: 50 + 20 + 10 + 1,
			maxScore: 50 + 20 + 10 + 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{Level: tt.level}
			// Run multiple times to test range
			for i := 0; i < 100; i++ {
				score := skill.ExamineCheck(tt.intStat, tt.wisStat)
				if score < tt.minScore || score > tt.maxScore {
					t.Errorf("ExamineCheck() = %d, want between %d and %d", score, tt.minScore, tt.maxScore)
				}
			}
		})
	}
}

// TestExamineSkillLevel_CanRevealHiddenDetail tests threshold checks
func TestExamineSkillLevel_CanRevealHiddenDetail(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		minLevel      int
		wantRevealed  bool
	}{
		{"skill meets threshold", 50, 50, true},
		{"skill above threshold", 75, 50, true},
		{"skill below threshold", 25, 50, false},
		{"skill exactly below", 49, 50, false},
		{"level 0 meets requirement", 0, 0, true},
		{"high level, low requirement", 100, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{Level: tt.level}
			got := skill.CanRevealHiddenDetail(tt.minLevel)
			if got != tt.wantRevealed {
				t.Errorf("CanRevealHiddenDetail(%d) = %v, want %v", tt.minLevel, got, tt.wantRevealed)
			}
		})
	}
}

// TestExamineSkillLevel_GetExamineBonusPercent tests bonus percentages
func TestExamineSkillLevel_GetExamineBonusPercent(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		wantBonus     int
	}{
		{"level 0", 0, 0},
		{"level 10", 10, 0},
		{"level 25", 25, 0},
		{"level 26", 26, 10},
		{"level 50", 50, 10},
		{"level 51", 51, 25},
		{"level 75", 75, 25},
		{"level 76", 76, 50},
		{"level 90", 90, 50},
		{"level 91", 91, 75},
		{"level 100", 100, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{Level: tt.level}
			got := skill.GetExamineBonusPercent()
			if got != tt.wantBonus {
				t.Errorf("GetExamineBonusPercent() = %d, want %d", got, tt.wantBonus)
			}
		})
	}
}

// TestExamineSkillLevel_AddXP tests XP gain and leveling
func TestExamineSkillLevel_AddXP(t *testing.T) {
	tests := []struct {
		name          string
		initialLevel  int
		initialXP     int
		xpToAdd       int
		wantLevel     int
		wantXP        int
	}{
		{
			name:          "level up from 1 to 2",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       10,
			wantLevel:     2,
			wantXP:        0,
		},
		{
			name:          "level up from 10 to 15",
			initialLevel:  10,
			initialXP:     0,
			xpToAdd:       50,
			wantLevel:     15,
			wantXP:        0,
		},
		{
			name:          "partial XP",
			initialLevel:  5,
			initialXP:     0,
			xpToAdd:       5,
			wantLevel:     5,
			wantXP:        5,
		},
		{
			name:          "capped at 100",
			initialLevel:  98,
			initialXP:     0,
			xpToAdd:       50,
			wantLevel:     100,
			wantXP:        0, // Excess XP lost when capping at max level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skill := &ExamineSkillLevel{
				Level: tt.initialLevel,
				XP:    tt.initialXP,
			}
			skill.AddXP(tt.xpToAdd)
			if skill.Level != tt.wantLevel {
				t.Errorf("AddXP() level = %d, want %d", skill.Level, tt.wantLevel)
			}
			if skill.XP != tt.wantXP {
				t.Errorf("AddXP() xp = %d, want %d", skill.XP, tt.wantXP)
			}
		})
	}
}

// TestRevealHiddenDetails tests automatic and check modes
func TestRevealHiddenDetails(t *testing.T) {
	skill := &ExamineSkillLevel{Level: 50}
	intStat, wisStat := 10, 10

	tests := []struct {
		name          string
		details       []HiddenDetail
		wantRevealed  []bool
	}{
		{
			name: "automatic mode - skill meets threshold",
			details: []HiddenDetail{
				{Text: "coins at bottom", MinExamineLevel: 25, Mode: "automatic"},
			},
			wantRevealed: []bool{true},
		},
		{
			name: "automatic mode - skill below threshold",
			details: []HiddenDetail{
				{Text: "hidden compartment", MinExamineLevel: 75, Mode: "automatic"},
			},
			wantRevealed: []bool{false},
		},
		{
			name: "check mode - passed",
			details: []HiddenDetail{
				{Text: "faint writing", MinExamineLevel: 0, Mode: "check", DC: 10, Stat: "INT"},
			},
			wantRevealed: []bool{true}, // Will pass with skill 50 + bonus 15 + roll >= 10
		},
		{
			name: "check mode - failed",
			details: []HiddenDetail{
				{Text: "secret door", MinExamineLevel: 0, Mode: "check", DC: 200, Stat: "INT"},
			},
			wantRevealed: []bool{false}, // DC too high
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skill.RevealHiddenDetails(tt.details, intStat, wisStat)
			for i, want := range tt.wantRevealed {
				if result[i].Revealed != want {
					t.Errorf("detail[%d].Revealed = %v, want %v (text: %s)", i, result[i].Revealed, want, result[i].Text)
				}
			}
		})
	}
}

// TestExamineSkillBonus tests INT/WIS bonuses
func TestExamineSkillBonus(t *testing.T) {
	tests := []struct {
		name     string
		intStat  int
		wisStat  int
		wantBonus int
	}{
		{"both stats 0", 0, 0, 0},
		{"int 10, wis 0", 10, 0, 10},
		{"int 0, wis 10", 0, 10, 5},
		{"int 10, wis 10", 10, 10, 15},
		{"int 20, wis 20", 20, 20, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExamineSkillBonus(tt.intStat, tt.wisStat)
			if got != tt.wantBonus {
				t.Errorf("ExamineSkillBonus(%d, %d) = %d, want %d", tt.intStat, tt.wisStat, got, tt.wantBonus)
			}
		})
	}
}