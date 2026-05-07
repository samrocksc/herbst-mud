package main

import (
	"herbst/dice"
)

// WeaponDamageResult holds the result of a weapon damage calculation.
type WeaponDamageResult struct {
	TotalDamage int
	RawRoll     int
	WeaponName  string
	IsUntrained bool
	OffHandDmg  int
}

// calculateWeaponDamage computes damage for an attack with equipped weapons.
// If no weapon is equipped, falls back to bare fists (1d6 + STR mod).
func (m *model) calculateWeaponDamage(strMod int) WeaponDamageResult {
	items := m.fetchEquippedCombatItems(m.currentCharacterID)
	mainWeapon := findMainHandWeapon(items)

	if mainWeapon == nil {
		// Bare fists: 1d6 + STR mod
		roll, total := dice.Roll(6, 1, strMod)
		if total < 1 {
			total = 1
		}
		return WeaponDamageResult{
			TotalDamage: total,
			RawRoll:     roll,
			WeaponName:  "fists",
			IsUntrained: false,
		}
	}

	// Fetch character skills for skill penalty check
	skills := m.fetchCharacterSkills()
	damageMod := strMod
	trained := isTrainedWithWeapon(mainWeapon, skills)

	if !trained {
		damageMod = damageMod / 2 // Half STR mod when untrained
	}

	totalMod := mainWeapon.DamageBonus + damageMod
	roll, total := dice.Roll(mainWeapon.DamageDiceSides, mainWeapon.DamageDiceCount, totalMod)
	if total < 1 {
		total = 1
	}

	if !trained {
		total = total / 2 // Half damage when untrained
		if total < 1 {
			total = 1
		}
	}

	result := WeaponDamageResult{
		TotalDamage: total,
		RawRoll:     roll,
		WeaponName:  mainWeapon.Name,
		IsUntrained: !trained,
	}

	// Off-hand weapon: 50% of its damage range (rounded down)
	offHand := findOffHandWeapon(items)
	if offHand != nil && !mainWeapon.IsTwoHanded {
		result.OffHandDmg = m.calcOffHandDamage(offHand, skills, strMod)
		result.TotalDamage += result.OffHandDmg
	}

	return result
}

// calcOffHandDamage computes 50% damage contribution from an off-hand weapon.
func (m *model) calcOffHandDamage(offHand *CombatItem, skills *CharacterSkills, strMod int) int {
	offMod := strMod
	offTrained := isTrainedWithWeapon(offHand, skills)
	if !offTrained {
		offMod = offMod / 2
	}
	offTotalMod := offHand.DamageBonus + offMod
	_, offDmg := dice.Roll(offHand.DamageDiceSides, offHand.DamageDiceCount, offTotalMod)
	if !offTrained {
		offDmg = offDmg / 2
	}
	offHandContribution := offDmg / 2 // 50% of off-hand damage
	if offHandContribution < 1 {
		offHandContribution = 1
	}
	return offHandContribution
}