package main

import (
	"math/rand"
)

// ExamineSkillLevel represents the examine skill system
type ExamineSkillLevel struct {
	Level int `json:"level"` // 0-100
	XP    int `json:"xp"`     // Experience points
}

// ExamineSkillBonus calculates the bonus based on stats
func ExamineSkillBonus(intStat, wisStat int) int {
	return intStat + (wisStat / 2)
}

// ExamineCheck performs an examine skill check
// Formula: skill_level + (INT × 1) + random(1, 10)
func (e *ExamineSkillLevel) ExamineCheck(intStat, wisStat int) int {
	bonus := ExamineSkillBonus(intStat, wisStat)
	roll := rand.Intn(10) + 1
	return e.Level + bonus + roll
}

// CanRevealHiddenDetail checks if the skill level is sufficient
func (e *ExamineSkillLevel) CanRevealHiddenDetail(minLevel int) bool {
	return e.Level >= minLevel
}

// ExamineXP rewards
const (
	ExamineXPFirstTime   = 1
	ExamineXPDiscover    = 2
	ExamineXPReveal      = 5
	ExamineXPDecrypt     = 10
)

// GetExamineBonusPercent returns bonus percentage based on skill level
func (e *ExamineSkillLevel) GetExamineBonusPercent() int {
	switch {
	case e.Level >= 91:
		return 75
	case e.Level >= 76:
		return 50
	case e.Level >= 51:
		return 25
	case e.Level >= 26:
		return 10
	default:
		return 0
	}
}

// HiddenDetail represents a hidden detail that can be revealed
type HiddenDetail struct {
	Text            string `json:"text"`
	MinExamineLevel int    `json:"min_examine_level"`
	Mode            string `json:"mode"` // "automatic" or "check"
	DC              int    `json:"dc"`    // Difficulty class for "check" mode
	Stat            string `json:"stat"` // "INT" or "WIS" for check mode
	Revealed        bool   `json:"revealed"`
	Source          string `json:"source"`
	Roll            int    `json:"roll,omitempty"`
}

// RevealHiddenDetails reveals hidden details based on examine skill
func (e *ExamineSkillLevel) RevealHiddenDetails(details []HiddenDetail, intStat, wisStat int) []HiddenDetail {
	result := make([]HiddenDetail, len(details))
	check := e.ExamineCheck(intStat, wisStat)

	for i, detail := range details {
		result[i] = detail

		switch detail.Mode {
		case "automatic":
			// Revealed automatically if skill >= requirement
			if e.Level >= detail.MinExamineLevel {
				result[i].Revealed = true
				result[i].Source = "skill_threshold"
			}
		case "check":
			// Roll against DC
			if e.Level >= detail.MinExamineLevel {
				if check >= detail.DC {
					result[i].Revealed = true
					result[i].Source = "check_passed"
					result[i].Roll = check
				}
			}
		}
	}

	return result
}

// AddXP adds experience points and handles level ups
func (e *ExamineSkillLevel) AddXP(amount int) {
	if e.Level >= 100 {
		return // Already at max level
	}

	e.XP += amount

	// Level up every 10 XP
	levelsGained := e.XP / 10
	if levelsGained > 0 {
		newLevel := e.Level + levelsGained
		if newLevel > 100 {
			// Calculate remaining XP when capped
			excessLevels := newLevel - 100
			e.XP = e.XP - (excessLevels * 10)
			newLevel = 100
		}
		e.Level = newLevel
		e.XP = e.XP % 10
	}
}