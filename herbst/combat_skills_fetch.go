package main

import (
	"encoding/json"
	"fmt"
)

// fetchCharacterSkills retrieves weapon and armor skill levels from the API.
func (m *model) fetchCharacterSkills() *CharacterSkills {
	if m.currentCharacterID == 0 {
		return &CharacterSkills{}
	}

	url := fmt.Sprintf("%s/characters/%d/skills", RESTAPIBase, m.currentCharacterID)
	resp, err := httpGet(url)
	if err != nil {
		return &CharacterSkills{}
	}
	defer resp.Body.Close()

	var result struct {
		Skills struct {
			Blades      struct{ Level int `json:"level"` } `json:"blades"`
			Staves      struct{ Level int `json:"level"` } `json:"staves"`
			Knives      struct{ Level int `json:"level"` } `json:"knives"`
			Martial     struct{ Level int `json:"level"` } `json:"martial"`
			Brawling    struct{ Level int `json:"level"` } `json:"brawling"`
			Tech        struct{ Level int `json:"level"` } `json:"tech"`
			LightArmor  struct{ Level int `json:"level"` } `json:"light_armor"`
			ClothArmor  struct{ Level int `json:"level"` } `json:"cloth_armor"`
			HeavyArmor  struct{ Level int `json:"level"` } `json:"heavy_armor"`
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