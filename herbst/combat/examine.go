package combat

import (
	"math/rand"
)

// HiddenDetailMode defines how a hidden detail is revealed
type HiddenDetailMode string

const (
	// ModeAutomatic reveals the detail when examine skill >= min_examine_level
	ModeAutomatic HiddenDetailMode = "automatic"
	// ModeCheck requires a roll against DC (d20 + stat bonus vs DC)
	ModeCheck HiddenDetailMode = "check"
)

// HiddenDetail represents a hidden detail on an item or NPC
type HiddenDetail struct {
	Text            string          `json:"text"`
	MinExamineLevel int             `json:"min_examine_level,omitempty"`
	Mode            HiddenDetailMode `json:"mode"`
	DC              int             `json:"dc,omitempty"`   // Difficulty class for check mode
	Stat            string          `json:"stat,omitempty"` // INT or WIS for check mode
}

// HiddenDetailResult represents the result of revealing a hidden detail
type HiddenDetailResult struct {
	Text          string `json:"text"`
	Revealed      bool   `json:"revealed"`
	Source        string `json:"source,omitempty"`        // "skill_threshold" or "check_passed"
	Roll          int    `json:"roll,omitempty"`          // The d20 roll (for check mode)
	DC            int    `json:"dc,omitempty"`            // The DC (for check mode)
	RequiredLevel int    `json:"required_level,omitempty"` // Required level (if not revealed)
}

// ExamineResult represents the complete result of an examine action
type ExamineResult struct {
	VisibleDescription string              `json:"visible_description"`
	HiddenDetails      []HiddenDetailResult `json:"hidden_details"`
	ExamineXP          int                 `json:"examine_xp"`
}

// ExamineOptions contains options for the examine action
type ExamineOptions struct {
	ExamineSkillLevel int    // Player's examine skill level (0-100)
	Intelligence      int    // Player's intelligence stat
	Wisdom            int    // Player's wisdom stat
}

// Examine performs an examine action on an item with hidden details
// Returns the examine result with revealed details based on skill and rolls
func Examine(hiddenDetails []HiddenDetail, opts ExamineOptions) *ExamineResult {
	result := &ExamineResult{
		HiddenDetails: make([]HiddenDetailResult, 0, len(hiddenDetails)),
		ExamineXP:     0,
	}

	for _, detail := range hiddenDetails {
		detailResult := HiddenDetailResult{
			Text: detail.Text,
		}

		switch detail.Mode {
		case ModeAutomatic:
			// Automatic mode: reveal if skill >= threshold
			if opts.ExamineSkillLevel >= detail.MinExamineLevel {
				detailResult.Revealed = true
				detailResult.Source = "skill_threshold"
				result.ExamineXP++ // XP for successful reveal
			} else {
				detailResult.Revealed = false
				detailResult.RequiredLevel = detail.MinExamineLevel
			}

		case ModeCheck:
			// Check mode: roll d20 + stat bonus vs DC
			roll := rand.Intn(20) + 1 // d20 roll (1-20)
			detailResult.Roll = roll
			detailResult.DC = detail.DC

			// Get stat bonus
			var statBonus int
			switch detail.Stat {
			case "INT":
				// +1 per 2 points above 10, -1 per 2 below 10
				statBonus = (opts.Intelligence - 10) / 2
			case "WIS":
				statBonus = (opts.Wisdom - 10) / 2
			default:
				statBonus = (opts.Wisdom-10)/2 + (opts.Intelligence-10)/2
			}

			totalRoll := roll + statBonus
			if totalRoll >= detail.DC {
				detailResult.Revealed = true
				detailResult.Source = "check_passed"
				result.ExamineXP++ // XP for successful check
			} else {
				detailResult.Revealed = false
				detailResult.Source = "check_failed"
			}

		default:
			// Unknown mode - treat as not revealed
			detailResult.Revealed = false
		}

		result.HiddenDetails = append(result.HiddenDetails, detailResult)
	}

	// Cap XP from examine at 3 (prevents farming)
	if result.ExamineXP > 3 {
		result.ExamineXP = 3
	}

	return result
}

// ParseHiddenDetails converts raw JSON data to HiddenDetail structs
// This is useful for loading from database JSON fields
func ParseHiddenDetails(data []map[string]interface{}) []HiddenDetail {
	details := make([]HiddenDetail, 0, len(data))
	for _, d := range data {
		detail := HiddenDetail{}

		if text, ok := d["text"].(string); ok {
			detail.Text = text
		}
		if minLevel, ok := d["min_examine_level"].(float64); ok {
			detail.MinExamineLevel = int(minLevel)
		}
		if mode, ok := d["mode"].(string); ok {
			detail.Mode = HiddenDetailMode(mode)
		}
		if dc, ok := d["dc"].(float64); ok {
			detail.DC = int(dc)
		}
		if stat, ok := d["stat"].(string); ok {
			detail.Stat = stat
		}

		details = append(details, detail)
	}
	return details
}

// FormatExamineOutput creates a human-readable string from an ExamineResult
func FormatExamineOutput(result *ExamineResult) string {
	output := ""

	// Show revealed details
	revealed := false
	for _, detail := range result.HiddenDetails {
		if detail.Revealed {
			if !revealed {
				output += "\n  You notice:\n"
				revealed = true
			}
			output += "  - " + detail.Text + "\n"
			if detail.Source == "check_passed" {
				output += "    [REVEALED via perception check: roll " + string(rune('0'+detail.Roll/10)) + string(rune('0'+detail.Roll%10)) + " vs DC " + string(rune('0'+detail.DC/10)) + string(rune('0'+detail.DC%10)) + "]\n"
			}
		}
	}

	// Show hints for unrevealed automatic details
	for _, detail := range result.HiddenDetails {
		if !detail.Revealed && detail.RequiredLevel > 0 {
			output += "\n  [Hidden detail requires examine level " + string(rune('0'+detail.RequiredLevel/10)) + string(rune('0'+detail.RequiredLevel%10)) + "]\n"
		}
	}

	return output
}