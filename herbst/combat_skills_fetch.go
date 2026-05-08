package main

import (
	"encoding/json"
	"fmt"
)

// CharacterSkills holds weapon and armor skill levels for a character.
// Sources data from the competency system, falls back to flat columns.
type CharacterSkills struct {
	Blades      int `json:"blades"`
	Staves      int `json:"staves"`
	Knives      int `json:"knives"`
	Martial     int `json:"martial"`
	Brawling    int `json:"brawling"`
	Tech        int `json:"tech"`
	LightArmor  int `json:"light_armor"`
	ClothArmor  int `json:"cloth_armor"`
	HeavyArmor  int `json:"heavy_armor"`
}

// CompetencyEntry represents one competency from the API
type CompetencyEntry struct {
	CategoryID        string  `json:"category_id"`
	CategoryName      string  `json:"category_name"`
	Xp                int     `json:"xp"`
	Level             int     `json:"level"`
	XpMultiplier      float64 `json:"xp_multiplier"`
	DamageMultiplier  float64 `json:"damage_multiplier"`
	DefenseMultiplier float64 `json:"defense_multiplier"`
}

// fetchCharacterSkills retrieves weapon and armor skill levels from the competency API,
// falling back to the flat-column skills API if competency data is unavailable.
func (m *model) fetchCharacterSkills() *CharacterSkills {
	if m.currentCharacterID == 0 {
		return &CharacterSkills{}
	}

	// Try competency system first
	skills := m.fetchCompetencies()
	if skills != nil {
		return skills
	}

	// Fallback to flat columns
	return m.fetchFlatSkills()
}

// fetchCompetencies loads skill levels from the competency API
func (m *model) fetchCompetencies() *CharacterSkills {
	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/competencies", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	var entries []CompetencyEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil
	}

	if len(entries) == 0 {
		return nil
	}

	skills := &CharacterSkills{}
	for _, e := range entries {
		switch e.CategoryID {
		case "blades":
			skills.Blades = e.Level
		case "staves":
			skills.Staves = e.Level
		case "knives":
			skills.Knives = e.Level
		case "martial":
			skills.Martial = e.Level
		case "brawling":
			skills.Brawling = e.Level
		case "tech":
			skills.Tech = e.Level
		case "light_armor":
			skills.LightArmor = e.Level
		case "cloth_armor":
			skills.ClothArmor = e.Level
		case "heavy_armor":
			skills.HeavyArmor = e.Level
		}
	}
	return skills
}

// fetchFlatSkills loads skill levels from the legacy flat-column API
func (m *model) fetchFlatSkills() *CharacterSkills {
	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/skills", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return &CharacterSkills{}
	}
	defer resp.Body.Close()

	var result struct {
		Skills struct {
			Blades     struct{ Level int `json:"level"` } `json:"blades"`
			Staves     struct{ Level int `json:"level"` } `json:"staves"`
			Knives     struct{ Level int `json:"level"` } `json:"knives"`
			Martial    struct{ Level int `json:"level"` } `json:"martial"`
			Brawling   struct{ Level int `json:"level"` } `json:"brawling"`
			Tech       struct{ Level int `json:"level"` } `json:"tech"`
			LightArmor struct{ Level int `json:"level"` } `json:"light_armor"`
			ClothArmor struct{ Level int `json:"level"` } `json:"cloth_armor"`
			HeavyArmor struct{ Level int `json:"level"` } `json:"heavy_armor"`
		} `json:"skills"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return &CharacterSkills{}
	}

	return &CharacterSkills{
		Blades:      result.Skills.Blades.Level,
		Staves:      result.Skills.Staves.Level,
		Knives:      result.Skills.Knives.Level,
		Martial:     result.Skills.Martial.Level,
		Brawling:    result.Skills.Brawling.Level,
		Tech:        result.Skills.Tech.Level,
		LightArmor:  result.Skills.LightArmor.Level,
		ClothArmor:  result.Skills.ClothArmor.Level,
		HeavyArmor:  result.Skills.HeavyArmor.Level,
	}
}