package main

// CharacterSkills holds weapon and armor skill levels for a character.
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

// isTrainedWithWeapon checks if a character has the required skill for a weapon.
func isTrainedWithWeapon(weapon *CombatItem, skills *CharacterSkills) bool {
	if weapon.SkillRequirement == "" {
		return true
	}

	var skillLevel int
	switch weapon.SkillRequirement {
	case "blades":
		skillLevel = skills.Blades
	case "staves":
		skillLevel = skills.Staves
	case "knives":
		skillLevel = skills.Knives
	case "martial":
		skillLevel = skills.Martial
	case "brawling":
		skillLevel = skills.Brawling
	case "tech":
		skillLevel = skills.Tech
	default:
		return true
	}

	requiredLevel := weapon.SkillRequirementLevel
	if requiredLevel <= 0 {
		requiredLevel = 1
	}

	return skillLevel >= requiredLevel
}

// isTrainedWithArmor checks if a character has the required skill for an armor piece.
func isTrainedWithArmor(item *CombatItem, skills *CharacterSkills) bool {
	if item.SkillRequirement == "" {
		return true
	}

	var skillLevel int
	switch item.SkillRequirement {
	case "light_armor":
		skillLevel = skills.LightArmor
	case "cloth_armor":
		skillLevel = skills.ClothArmor
	case "heavy_armor":
		skillLevel = skills.HeavyArmor
	default:
		return true
	}

	requiredLevel := item.SkillRequirementLevel
	if requiredLevel <= 0 {
		requiredLevel = 1
	}

	return skillLevel >= requiredLevel
}