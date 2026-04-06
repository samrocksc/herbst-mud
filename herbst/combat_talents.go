package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// loadCombatTalents fetches equipped talents for the current character
func (m *model) loadCombatTalents() {
	if m.currentCharacterID == 0 {
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var result struct {
		Slots []interface{} `json:"slots"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// Convert slots to equipped talents
	m.combatTalents = make([]EquippedTalent, 5) // slots 1-4 (index 0 is unused)
	for i := 1; i < len(result.Slots) && i <= 4; i++ {
		if result.Slots[i] != nil {
			slotData := result.Slots[i].(map[string]interface{})
			talent := EquippedTalent{
				Slot: i,
			}
			if id, ok := slotData["id"].(float64); ok {
				talent.ID = int(id)
			}
			if name, ok := slotData["name"].(string); ok {
				talent.Name = name
			}
			if desc, ok := slotData["description"].(string); ok {
				talent.Description = desc
			}
			// New effect system fields
			if effectType, ok := slotData["effectType"].(string); ok {
				talent.EffectType = effectType
			}
			if effectValue, ok := slotData["effectValue"].(float64); ok {
				talent.EffectValue = int(effectValue)
			}
			if effectDuration, ok := slotData["effectDuration"].(float64); ok {
				talent.EffectDuration = int(effectDuration)
			}
			if cooldown, ok := slotData["cooldown"].(float64); ok {
				talent.Cooldown = int(cooldown)
			}
			if manaCost, ok := slotData["manaCost"].(float64); ok {
				talent.ManaCost = int(manaCost)
			}
			if staminaCost, ok := slotData["staminaCost"].(float64); ok {
				talent.StaminaCost = int(staminaCost)
			}
			m.combatTalents[i] = talent
		}
	}
}

// useCombatTalent executes a talent in combat
// Now uses classless skills from slots 1-4
func (m *model) useCombatTalent(slot int) {
	if slot < 1 || slot > 4 {
		m.addCombatLog("Invalid slot")
		return
	}

	// Use classless skill instead of legacy talent
	m.useClasslessSkill(slot)
}

// useHealthPotion attempts to use a health potion during combat
func (m *model) useHealthPotion() {
	// If a potion is equipped in R slot, use that first
	if m.equippedPotion != nil && m.equippedPotion.ID != 0 {
		m.useEquippedPotion()
		return
	}

	// Fall back to inventory search
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d&itemType=potion", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.addCombatLog("❌ Can't access inventory")
		return
	}
	defer resp.Body.Close()

	var items []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		ItemType      string `json:"itemType"`
		Description   string `json:"description"`
		EffectType    string `json:"effectType"`
		EffectValue   int    `json:"effectValue"`
		EffectDuration int   `json:"effectDuration"`
		Healing       int    `json:"healing"` // Deprecated fallback
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		m.addCombatLog("❌ No potions available")
		return
	}

	if len(items) == 0 {
		m.addCombatLog("❌ No health potions in inventory")
		return
	}

	// Use the first potion
	potion := items[0]
	m.applyPotionEffect(potion.ID, potion.Name, potion.EffectType, potion.EffectValue, potion.Healing)
}

// useEquippedPotion uses the potion equipped in the R slot
func (m *model) useEquippedPotion() {
	if m.equippedPotion == nil || m.equippedPotion.ID == 0 {
		m.addCombatLog("❌ No potion equipped in R slot")
		return
	}

	m.applyPotionEffect(m.equippedPotion.ID, m.equippedPotion.Name,
		m.equippedPotion.EffectType, m.equippedPotion.EffectValue, 0)
}

// applyPotionEffect applies a potion's effect and consumes it
func (m *model) applyPotionEffect(potionID int, potionName, effectType string, effectValue, fallbackHealing int) {
	// Determine the effect value
	value := effectValue
	if value <= 0 {
		value = fallbackHealing
	}
	if value <= 0 {
		value = 25 // Default fallback
	}

	// Delete the potion from inventory
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/equipment/%d", RESTAPIBase, potionID), nil)
	if err != nil {
		m.addCombatLog("❌ Can't use potion")
		return
	}
	client := &http.Client{Timeout: 5e9}
	resp, err := client.Do(req)
	if err != nil {
		m.addCombatLog("❌ Can't use potion")
		return
	}
	resp.Body.Close()

	// Apply effect based on type
	switch effectType {
	case "heal", "": // Default to healing
		m.characterHP += value
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, value)
		m.addCombatLog(fmt.Sprintf("🧪 Used %s! Restored %d HP", potionName, value))
		m.addCombatLog(fmt.Sprintf("❤️ HP: %d/%d", m.characterHP, m.characterMaxHP))
	case "damage":
		// Damage potions would damage the target
		m.addCombatLog(fmt.Sprintf("🧪 Used %s! %d damage applied to target", potionName, value))
	case "dot":
		m.addCombatLog(fmt.Sprintf("🧪 Used %s! Applied DoT effect", potionName))
	case "buff_armor", "buff_dodge", "buff_crit":
		m.addCombatLog(fmt.Sprintf("🧪 Used %s! Applied %s +%d", potionName, effectType[5:], value))
	default:
		// Unknown effect type, treat as healing
		m.characterHP += value
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, value)
		m.addCombatLog(fmt.Sprintf("🧪 Used %s! Effect applied", potionName))
	}

	// Clear equipped potion if this was the equipped one
	if m.equippedPotion != nil && m.equippedPotion.ID == potionID {
		m.equippedPotion = nil
	}
}

// getTalentSlotName returns a display name for a skill slot
// Generic display: [Skill 1], [Skill 2], etc. - works with any skill system
func (m *model) getTalentSlotName(slot int) string {
	if slot < 1 || slot > 4 {
		return "[Attack]"
	}

	// Check if there's a skill equipped in this slot
	if m.combatSkills != nil && slot <= 5 {
		skill := m.combatSkills.EquippedSkill[slot-1]
		if skill.ID != 0 {
			// Show actual skill name (truncted if too long)
			name := skill.Name
			if len(name) > 12 {
				name = name[:12]
			}
			return fmt.Sprintf("[%s]", name)
		}
	}

	// Fallback: check legacy talents
	if slot < len(m.combatTalents) && m.combatTalents[slot].ID != 0 {
		name := m.combatTalents[slot].Name
		if len(name) > 12 {
			name = name[:12]
		}
		return fmt.Sprintf("[%s]", name)
	}

	return "[Attack]"
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