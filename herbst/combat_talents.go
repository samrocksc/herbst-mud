package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// useHealthPotion attempts to use a health potion during combat
func (m *model) useHealthPotion() {
	if m.equippedPotion != nil && m.equippedPotion.ID != 0 {
		m.useEquippedPotion()
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d&itemType=potion", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.addCombatLog("Can't access inventory")
		return
	}
	defer resp.Body.Close()

	var items []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		ItemType       string `json:"itemType"`
		EffectType     string `json:"effectType"`
		EffectValue    int    `json:"effectValue"`
		EffectDuration int    `json:"effectDuration"`
		Healing        int    `json:"healing"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		m.addCombatLog("No potions available")
		return
	}

	if len(items) == 0 {
		m.addCombatLog("No health potions in inventory")
		return
	}

	potion := items[0]
	m.applyPotionEffect(potion.ID, potion.Name, potion.EffectType, potion.EffectValue, potion.Healing)
}

// useEquippedPotion uses the potion equipped in the R slot
func (m *model) useEquippedPotion() {
	if m.equippedPotion == nil || m.equippedPotion.ID == 0 {
		m.addCombatLog("No potion equipped in R slot")
		return
	}
	m.applyPotionEffect(m.equippedPotion.ID, m.equippedPotion.Name,
		m.equippedPotion.EffectType, m.equippedPotion.EffectValue, 0)
}

// applyPotionEffect applies a potion's effect and consumes it
func (m *model) applyPotionEffect(potionID int, potionName, effectType string, effectValue, fallbackHealing int) {
	value := effectValue
	if value <= 0 {
		value = fallbackHealing
	}
	if value <= 0 {
		value = 25
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/equipment/%d", RESTAPIBase, potionID), nil)
	if err != nil {
		m.addCombatLog("Can't use potion")
		return
	}
	client := &http.Client{Timeout: 5e9}
	resp, err := client.Do(req)
	if err != nil {
		m.addCombatLog("Can't use potion")
		return
	}
	resp.Body.Close()

	switch effectType {
	case "heal", "":
		m.characterHP += value
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, value)
		m.addCombatLog(fmt.Sprintf("Used %s! Restored %d HP", potionName, value))
		m.addCombatLog(fmt.Sprintf("HP: %d/%d", m.characterHP, m.characterMaxHP))
	default:
		m.characterHP += value
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, value)
		m.addCombatLog(fmt.Sprintf("Used %s! Effect applied", potionName))
	}

	if m.equippedPotion != nil && m.equippedPotion.ID == potionID {
		m.equippedPotion = nil
	}
}

// getPotionSlotName returns a display name for the potion slot
func (m *model) getPotionSlotName() string {
	if m.equippedPotion == nil || m.equippedPotion.ID == 0 {
		return "[Potion]"
	}
	name := m.equippedPotion.Name
	if len(name) > 10 {
		name = name[:10]
	}
	return fmt.Sprintf("[%s]", name)
}