package main

// ArmorResult holds the total AC contribution from equipped armor.
type ArmorResult struct {
	TotalAC      int
	ArmorDetails []ArmorDetail
}

// ArmorDetail describes one armor piece's contribution.
type ArmorDetail struct {
	Name   string
	AC     int
	Halved bool
}

// calculateArmorAC computes total AC bonus from all equipped items.
// Items with armor_rating > 0 contribute. Untrained skill = half AC.
func (m *model) calculateArmorAC(charID int, skills *CharacterSkills) ArmorResult {
	items := m.fetchEquippedCombatItems(charID)
	if len(items) == 0 {
		return ArmorResult{}
	}

	var result ArmorResult
	for i := range items {
		item := &items[i]
		if item.ArmorRating <= 0 {
			continue
		}

		ac := item.ArmorRating
		halved := false

		if !isTrainedWithArmor(item, skills) {
			ac = ac / 2 // Half AC when untrained
			if ac < 1 {
				ac = 1
			}
			halved = true
		}

		result.TotalAC += ac
		result.ArmorDetails = append(result.ArmorDetails, ArmorDetail{
			Name:   item.Name,
			AC:     ac,
			Halved: halved,
		})
	}

	return result
}

// calculatePlayerAC computes the player's total AC including armor.
func (m *model) calculatePlayerAC() int {
	baseAC := 10 + m.getDexModifier() + m.characterLevel/2

	skills := m.fetchCharacterSkills()
	armor := m.calculateArmorAC(m.currentCharacterID, skills)

	return baseAC + armor.TotalAC
}