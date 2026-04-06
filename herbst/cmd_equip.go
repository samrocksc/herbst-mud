package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleEquipCommand shows the equip screen with slots for talents and potion
func (m *model) handleEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.showEquipScreen()
		return
	}

	action := parts[1]
	switch action {
	case "talent":
		m.handleEquipTalent(parts[2:])
	case "potion":
		m.handleEquipPotion(parts[2:])
	default:
		m.showEquipScreen()
	}
}

// showEquipScreen displays the equipment slots UI
func (m *model) showEquipScreen() {
	output := "=== Equipment Slots ===\n\n"
	output += "Combat Talents (1-4):\n"

	// Load current talents
	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID))
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if json.NewDecoder(resp.Body).Decode(&result) == nil {
				slots, ok := result["slots"].([]interface{})
				if ok {
					for i := 1; i <= 4; i++ {
						if i < len(slots) && slots[i] != nil {
							slot := slots[i].(map[string]interface{})
							name := slot["name"].(string)
							effectType := ""
							if et, ok := slot["effectType"].(string); ok && et != "" {
								effectType = fmt.Sprintf(" [%s]", et)
							}
							output += fmt.Sprintf("  [%d] %s%s\n", i, name, effectType)
						} else {
							output += fmt.Sprintf("  [%d] (empty)\n", i)
						}
					}
				}
			}
		}
	}

	output += "\nPotion Slot (R):\n"

	// Show equipped potion
	if m.equippedPotion != nil && m.equippedPotion.ID != 0 {
		effectType := m.equippedPotion.EffectType
		if effectType == "" {
			effectType = "heal"
		}
		output += fmt.Sprintf("  [R] %s [%s +%d]\n", m.equippedPotion.Name, effectType, m.equippedPotion.EffectValue)
	} else {
		// Check inventory for available potions
		resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d&itemType=potion", RESTAPIBase, m.currentCharacterID))
		if err == nil {
			defer resp.Body.Close()
			var potions []struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				EffectType  string `json:"effectType"`
				EffectValue int    `json:"effectValue"`
				Healing     int    `json:"healing"`
			}
			if json.NewDecoder(resp.Body).Decode(&potions) == nil {
				if len(potions) > 0 {
					output += "  [R] (no potion equipped)\n"
					output += fmt.Sprintf("      %d potion(s) in inventory\n", len(potions))
				} else {
					output += "  [R] (no potion equipped)\n"
					output += "      No potions in inventory\n"
				}
			}
		}
	}

	output += "\nCommands:\n"
	output += "  equip talent <id> <slot> - Equip talent to slot 1-4\n"
	output += "  equip potion <id>        - Equip potion to R slot\n"
	output += "  equip clear <slot>        - Clear slot (1-4 or R)\n"
	output += "  talents                   - View available talents\n"
	output += "  inventory                 - View inventory items\n"

	m.AppendMessage(output, "info")
}

// handleEquipTalent equips a talent to a slot
func (m *model) handleEquipTalent(args []string) {
	if len(args) < 2 {
		m.AppendMessage("Usage: equip talent <talent_id> <slot>\nSlots: 1-4", "error")
		return
	}

	talentID := args[0]
	slot := args[1]

	slotNum := 0
	fmt.Sscanf(slot, "%d", &slotNum)
	if slotNum < 1 || slotNum > 4 {
		m.AppendMessage("Slot must be between 1 and 4", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
	reqBody := fmt.Sprintf(`{"talent_id":%s,"slot":%s}`, talentID, slot)
	resp, err := httpPost(url, reqBody)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error equipping talent: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		m.AppendMessage("Failed to equip talent", "error")
		return
	}

	m.AppendMessage(fmt.Sprintf("Talent equipped in slot %s", slot), "success")
	m.loadCombatTalents()
}

// handleEquipPotion equips a potion to the R slot
func (m *model) handleEquipPotion(args []string) {
	if len(args) < 1 {
		// Show available potions
		resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d&itemType=potion", RESTAPIBase, m.currentCharacterID))
		if err != nil {
			m.AppendMessage("Error fetching potions", "error")
			return
		}
		defer resp.Body.Close()

		var potions []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			EffectType  string `json:"effectType"`
			EffectValue int    `json:"effectValue"`
			Healing     int    `json:"healing"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&potions); err != nil {
			m.AppendMessage("Error parsing potions", "error")
			return
		}

		if len(potions) == 0 {
			m.AppendMessage("No potions in inventory.\nUse 'take <potion>' to pick up a potion.", "info")
			return
		}

		output := "Available Potions:\n"
		for _, p := range potions {
			effectType := p.EffectType
			if effectType == "" {
				effectType = "heal"
			}
			effectValue := p.EffectValue
			if effectValue == 0 {
				effectValue = p.Healing
			}
			output += fmt.Sprintf("  ID %d: %s [%s +%d]\n", p.ID, p.Name, effectType, effectValue)
		}
		output += "\nUse: equip potion <id>"
		m.AppendMessage(output, "info")
		return
	}

	// Get potion by ID
	potionID := args[0]
	var potionIDNum int
	fmt.Sscanf(potionID, "%d", &potionIDNum)

	// Fetch potion details
	resp, err := httpGet(fmt.Sprintf("%s/equipment/%s", RESTAPIBase, potionID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching potion: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	var potion struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		EffectType    string `json:"effectType"`
		EffectValue   int    `json:"effectValue"`
		EffectDuration int   `json:"effectDuration"`
		ItemType      string `json:"itemType"`
		OwnerID       int    `json:"ownerId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&potion); err != nil {
		m.AppendMessage("Error parsing potion data", "error")
		return
	}

	// Verify it's a potion
	if potion.ItemType != "potion" {
		m.AppendMessage("That item is not a potion", "error")
		return
	}

	// Verify ownership
	if potion.OwnerID != m.currentCharacterID {
		m.AppendMessage("You don't own that potion", "error")
		return
	}

	// Set as equipped potion
	m.equippedPotion = &EquippedPotion{
		ID:            potion.ID,
		Name:          potion.Name,
		Description:   potion.Description,
		EffectType:    potion.EffectType,
		EffectValue:   potion.EffectValue,
		EffectDuration: potion.EffectDuration,
	}

	// Default effect type if empty
	if m.equippedPotion.EffectType == "" {
		m.equippedPotion.EffectType = "heal"
	}
	if m.equippedPotion.EffectValue == 0 {
		m.equippedPotion.EffectValue = 25 // Default healing
	}

	m.AppendMessage(fmt.Sprintf("Equipped %s in R slot. Press R in combat to use.", potion.Name), "success")
}

// handleEquipClear clears a slot
func (m *model) handleEquipClear(slot string) {
	if slot == "R" || slot == "r" {
		m.equippedPotion = nil
		m.AppendMessage("Potion slot cleared", "success")
		return
	}

	// Clear talent slot
	slotNum := 0
	fmt.Sscanf(slot, "%d", &slotNum)
	if slotNum < 1 || slotNum > 4 {
		m.AppendMessage("Slot must be 1-4 or R", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/talents/%d", RESTAPIBase, m.currentCharacterID, slotNum)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error clearing slot: %v", err), "error")
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error clearing slot: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	m.AppendMessage(fmt.Sprintf("Slot %d cleared", slotNum), "success")
	m.loadCombatTalents()
}