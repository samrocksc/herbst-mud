package combat

import (
	"math/rand"
	"testing"
)

func TestExamine_AutomaticMode_RevealedAtThreshold(t *testing.T) {
	details := []HiddenDetail{
		{
			Text:            "Coins glint at the bottom",
			MinExamineLevel: 0,
			Mode:            ModeAutomatic,
		},
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 5,
		Intelligence:      10,
		Wisdom:            10,
	}

	result := Examine(details, opts)

	if len(result.HiddenDetails) != 1 {
		t.Fatalf("Expected 1 hidden detail, got %d", len(result.HiddenDetails))
	}

	if !result.HiddenDetails[0].Revealed {
		t.Error("Expected detail to be revealed at skill level 5 with threshold 0")
	}
	if result.HiddenDetails[0].Source != "skill_threshold" {
		t.Errorf("Expected source 'skill_threshold', got '%s'", result.HiddenDetails[0].Source)
	}
	if result.ExamineXP != 1 {
		t.Errorf("Expected 1 XP for reveal, got %d", result.ExamineXP)
	}
}

func TestExamine_AutomaticMode_NotRevealedBelowThreshold(t *testing.T) {
	details := []HiddenDetail{
		{
			Text:            "A hidden compartment behind the statue",
			MinExamineLevel: 75,
			Mode:            ModeAutomatic,
		},
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 50, // Below threshold
		Intelligence:      10,
		Wisdom:            10,
	}

	result := Examine(details, opts)

	if result.HiddenDetails[0].Revealed {
		t.Error("Expected detail NOT to be revealed below threshold")
	}
	if result.HiddenDetails[0].RequiredLevel != 75 {
		t.Errorf("Expected required level 75, got %d", result.HiddenDetails[0].RequiredLevel)
	}
	if result.ExamineXP != 0 {
		t.Errorf("Expected 0 XP without reveal, got %d", result.ExamineXP)
	}
}

func TestExamine_AutomaticMode_ExactThreshold(t *testing.T) {
	details := []HiddenDetail{
		{
			Text:            "Faint writing: WISH UPON THE OOZE",
			MinExamineLevel: 50,
			Mode:            ModeAutomatic,
		},
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 50, // Exactly at threshold
		Intelligence:      10,
		Wisdom:            10,
	}

	result := Examine(details, opts)

	if !result.HiddenDetails[0].Revealed {
		t.Error("Expected detail to be revealed at exact threshold")
	}
}

func TestExamine_CheckMode_PassWithRoll(t *testing.T) {
	// Use fixed seed for deterministic testing
	rand.Seed(42)

	details := []HiddenDetail{
		{
			Text: "A crack in the base",
			Mode: ModeCheck,
			DC:   10, // Low DC
			Stat: "INT",
		},
	}

	// High INT gives +5 bonus (INT 20: (20-10)/2 = 5)
	opts := ExamineOptions{
		ExamineSkillLevel: 0,
		Intelligence:      20,
		Wisdom:            10,
	}

	// Run multiple times to account for randomness
	successCount := 0
	for i := 0; i < 100; i++ {
		result := Examine(details, opts)
		if result.HiddenDetails[0].Revealed {
			successCount++
		}
	}

	// With INT 20 (+5 bonus) and DC 10, even rolling a 1 gives us 6 (1+5=6)
	// So we should pass when roll >= 5, which is 80% of the time (16/20 rolls)
	// With 100 runs, expect ~80 successes
	if successCount < 60 {
		t.Errorf("Expected at least 60 successes with high INT, got %d", successCount)
	}
}

func TestExamine_CheckMode_FailWithLowRoll(t *testing.T) {
	// Use fixed seed for deterministic testing
	rand.Seed(42)

	details := []HiddenDetail{
		{
			Text: "A hidden compartment",
			Mode: ModeCheck,
			DC:   30, // Very high DC
			Stat: "INT",
		},
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 0,
		Intelligence:      10, // No bonus
		Wisdom:            10,
	}

	// Run multiple times - should almost always fail
	successCount := 0
	for i := 0; i < 100; i++ {
		result := Examine(details, opts)
		if result.HiddenDetails[0].Revealed {
			successCount++
		}
	}

	// DC 30 with no bonus requires rolling 30 on a d20, which is impossible
	// Natural 20 is only 20, so should fail 100% of the time
	if successCount > 0 {
		t.Errorf("Expected 0 successes with impossible DC, got %d", successCount)
	}
}

func TestExamine_CheckMode_WisdomBonus(t *testing.T) {
	details := []HiddenDetail{
		{
			Text: "A crack in the base",
			Mode: ModeCheck,
			DC:   15,
			Stat: "WIS",
		},
	}

	// WIS 18 gives +4 bonus: (18-10)/2 = 4
	opts := ExamineOptions{
		ExamineSkillLevel: 0,
		Intelligence:      10,
		Wisdom:            18,
	}

	successCount := 0
	for i := 0; i < 100; i++ {
		result := Examine(details, opts)
		if result.HiddenDetails[0].Revealed {
			successCount++
		}
	}

	// DC 15 with +4 bonus requires roll >= 11 (roll 11-20 = 10/20 = 50%)
	// With 100 runs, expect ~50 successes
	if successCount < 35 || successCount > 65 {
		t.Errorf("Expected ~50 successes with WIS bonus, got %d", successCount)
	}
}

func TestExamine_MultipleDetails(t *testing.T) {
	details := []HiddenDetail{
		{
			Text:            "Coins glint at the bottom",
			MinExamineLevel: 0,
			Mode:            ModeAutomatic,
		},
		{
			Text:            "A crack in the base",
			Mode:            ModeCheck,
			DC:              10,
			Stat:            "WIS",
		},
		{
			Text:            "A hidden compartment behind the statue",
			MinExamineLevel: 75,
			Mode:            ModeAutomatic,
		},
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 50,
		Intelligence:      14, // +2 bonus
		Wisdom:            14, // +2 bonus
	}

	result := Examine(details, opts)

	if len(result.HiddenDetails) != 3 {
		t.Fatalf("Expected 3 hidden details, got %d", len(result.HiddenDetails))
	}

	// First detail (automatic, threshold 0) - should be revealed
	if !result.HiddenDetails[0].Revealed {
		t.Error("Expected first detail to be revealed")
	}

	// Second detail (check, DC 10) - check result should be present
	if result.HiddenDetails[1].Roll < 1 || result.HiddenDetails[1].Roll > 20 {
		t.Errorf("Expected roll between 1-20, got %d", result.HiddenDetails[1].Roll)
	}
	if result.HiddenDetails[1].DC != 10 {
		t.Errorf("Expected DC 10, got %d", result.HiddenDetails[1].DC)
	}

	// Third detail (automatic, threshold 75) - should NOT be revealed
	if result.HiddenDetails[2].Revealed {
		t.Error("Expected third detail NOT to be revealed")
	}
	if result.HiddenDetails[2].RequiredLevel != 75 {
		t.Errorf("Expected required level 75, got %d", result.HiddenDetails[2].RequiredLevel)
	}

	// XP should be at least 1 (first detail always revealed)
	if result.ExamineXP < 1 {
		t.Errorf("Expected at least 1 XP, got %d", result.ExamineXP)
	}
}

func TestExamine_XPCapped(t *testing.T) {
	// Create many details that would give lots of XP
	details := make([]HiddenDetail, 10)
	for i := range details {
		details[i] = HiddenDetail{
			Text:            "Detail " + string(rune('0'+i)),
			MinExamineLevel: 0,
			Mode:            ModeAutomatic,
		}
	}

	opts := ExamineOptions{
		ExamineSkillLevel: 100,
		Intelligence:      10,
		Wisdom:            10,
	}

	result := Examine(details, opts)

	// XP should be capped at 3
	if result.ExamineXP > 3 {
		t.Errorf("Expected XP capped at 3, got %d", result.ExamineXP)
	}
}

func TestParseHiddenDetails(t *testing.T) {
	data := []map[string]interface{}{
		{
			"text":              "Coins glint at the bottom",
			"min_examine_level": float64(0),
			"mode":              "automatic",
		},
		{
			"text": "A crack in the base",
			"mode": "check",
			"dc":   float64(30),
			"stat": "INT",
		},
		{
			"text":              "A hidden compartment",
			"min_examine_level": float64(75),
			"mode":              "automatic",
		},
	}

	details := ParseHiddenDetails(data)

	if len(details) != 3 {
		t.Fatalf("Expected 3 details, got %d", len(details))
	}

	// First detail
	if details[0].Text != "Coins glint at the bottom" {
		t.Errorf("Expected text 'Coins glint at the bottom', got '%s'", details[0].Text)
	}
	if details[0].MinExamineLevel != 0 {
		t.Errorf("Expected min_examine_level 0, got %d", details[0].MinExamineLevel)
	}
	if details[0].Mode != ModeAutomatic {
		t.Errorf("Expected mode 'automatic', got '%s'", details[0].Mode)
	}

	// Second detail
	if details[1].Text != "A crack in the base" {
		t.Errorf("Expected text 'A crack in the base', got '%s'", details[1].Text)
	}
	if details[1].Mode != ModeCheck {
		t.Errorf("Expected mode 'check', got '%s'", details[1].Mode)
	}
	if details[1].DC != 30 {
		t.Errorf("Expected DC 30, got %d", details[1].DC)
	}
	if details[1].Stat != "INT" {
		t.Errorf("Expected stat 'INT', got '%s'", details[1].Stat)
	}

	// Third detail
	if details[2].MinExamineLevel != 75 {
		t.Errorf("Expected min_examine_level 75, got %d", details[2].MinExamineLevel)
	}
}

func TestExamine_EmptyDetails(t *testing.T) {
	details := []HiddenDetail{}

	opts := ExamineOptions{
		ExamineSkillLevel: 50,
		Intelligence:      10,
		Wisdom:            10,
	}

	result := Examine(details, opts)

	if len(result.HiddenDetails) != 0 {
		t.Errorf("Expected 0 hidden details, got %d", len(result.HiddenDetails))
	}
	if result.ExamineXP != 0 {
		t.Errorf("Expected 0 XP for empty details, got %d", result.ExamineXP)
	}
}

func TestExamine_DefaultStatBonus(t *testing.T) {
	details := []HiddenDetail{
		{
			Text: "Something suspicious",
			Mode: ModeCheck,
			DC:   15,
			// No stat specified - should use default (INT+WIS bonus)
		},
	}

	// INT 14 (+2) + WIS 14 (+2) = +4 total bonus
	opts := ExamineOptions{
		ExamineSkillLevel: 0,
		Intelligence:      14,
		Wisdom:            14,
	}

	// Run multiple times to check for successful reveals
	successCount := 0
	for i := 0; i < 100; i++ {
		result := Examine(details, opts)
		if result.HiddenDetails[0].Revealed {
			successCount++
		}
	}

	// DC 15 with +4 bonus requires roll >= 11 (50% chance)
	// With 100 runs, expect ~50 successes
	if successCount < 35 || successCount > 65 {
		t.Errorf("Expected ~50 successes with default stat, got %d", successCount)
	}
}