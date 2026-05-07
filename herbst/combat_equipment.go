package main

import (
	"encoding/json"
	"fmt"
)

// CombatItem holds equipment fields relevant to combat calculations.
type CombatItem struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Slot                  string `json:"slot"`
	IsEquipped            bool   `json:"isEquipped"`
	ItemType              string `json:"itemType"`
	DamageDiceCount       int    `json:"damage_dice_count"`
	DamageDiceSides       int    `json:"damage_dice_sides"`
	DamageBonus           int    `json:"damage_bonus"`
	DamageType            string `json:"damage_type"`
	WeaponType            string `json:"weapon_type"`
	IsTwoHanded           bool   `json:"is_two_handed"`
	ArmorRating           int    `json:"armor_rating"`
	ArmorType             string `json:"armor_type"`
	SkillRequirement      string `json:"skill_requirement"`
	SkillRequirementLevel int    `json:"skill_requirement_level"`
}

// fetchEquippedCombatItems retrieves all equipped items for a character.
func (m *model) fetchEquippedCombatItems(charID int) []CombatItem {
	if charID == 0 {
		return nil
	}

	url := fmt.Sprintf("%s/equipment?ownerId=%d&isEquipped=true", RESTAPIBase, charID)
	resp, err := httpGet(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var items []CombatItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil
	}
	return items
}

// findWeaponInSlot returns the first equipped weapon matching a slot.
func findWeaponInSlot(items []CombatItem, slot string) *CombatItem {
	for i := range items {
		if items[i].Slot == slot && items[i].DamageDiceCount > 0 {
			return &items[i]
		}
	}
	return nil
}

// findMainHandWeapon returns the weapon in main_hand slot, or any weapon.
func findMainHandWeapon(items []CombatItem) *CombatItem {
	if w := findWeaponInSlot(items, "main_hand"); w != nil {
		return w
	}
	for i := range items {
		if items[i].DamageDiceCount > 0 {
			return &items[i]
		}
	}
	return nil
}

// findOffHandWeapon returns a weapon in off_hand or tail slot.
func findOffHandWeapon(items []CombatItem) *CombatItem {
	if w := findWeaponInSlot(items, "off_hand"); w != nil {
		return w
	}
	if w := findWeaponInSlot(items, "tail"); w != nil {
		return w
	}
	return nil
}

// findArmorItems returns all equipped items with armor_rating > 0.
func findArmorItems(items []CombatItem) []CombatItem {
	var armor []CombatItem
	for i := range items {
		if items[i].ArmorRating > 0 {
			armor = append(armor, items[i])
		}
	}
	return armor
}