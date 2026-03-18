package nomenclature

import "testing"

func TestSkillsNomenclature(t *testing.T) {
	requiredSkills := map[string]string{
		"blades":     "STR",
		"staves":     "DEX",
		"knives":     "DEX",
		"martial":    "DEX",
		"brawling":   "STR",
		"tech":       "INT",
		"fire_magic": "INT",
		"water_magic": "INT",
		"wind_magic": "WIS",
	}

	t.Logf("Total skills: %d", len(requiredSkills))
	if len(requiredSkills) != 9 {
		t.Errorf("expected 9 skills, got %d", len(requiredSkills))
	}
}

func TestTalentsNomenclature(t *testing.T) {
	requiredTalents := []string{
		"slash", "parry", "smash", "crash", "shield_bash",
		"battle_cry", "second_wind", "hail_storm", "iron_will", "heavy_strike",
	}

	t.Logf("Total talents: %d", len(requiredTalents))
	if len(requiredTalents) != 10 {
		t.Errorf("expected 10 talents, got %d", len(requiredTalents))
	}
}

func TestStatsNomenclature(t *testing.T) {
	stats := []string{"STR", "CON", "WIS", "DEX", "INT"}
	t.Logf("Total stats: %d", len(stats))
	if len(stats) != 5 {
		t.Errorf("expected 5 stats, got %d", len(stats))
	}
}

func TestFighterLoadout(t *testing.T) {
	skills := []string{"blades", "brawling"}
	talents := []string{"slash", "parry", "smash", "crash"}
	
	t.Logf("Fighter skills: %v", skills)
	t.Logf("Fighter talents: %v", talents)
	
	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
	if len(talents) != 4 {
		t.Errorf("expected 4 talents, got %d", len(talents))
	}
}