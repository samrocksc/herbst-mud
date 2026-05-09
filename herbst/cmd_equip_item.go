package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleEquipItemCommand equips an inventory item by name.
func (m *model) handleEquipItemCommand(itemName string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	// Fetch unequipped items owned by this character
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching inventory: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	var items []equipItemData
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		m.AppendMessage("Error reading inventory.", "error")
		return
	}

	// Find matching item (unequipped only)
	var target *equipItemData
	for i := range items {
		if items[i].IsEquipped {
			continue
		}
		nameLower := strings.ToLower(items[i].Name)
		if fuzzyWordMatch(items[i].Name, itemName) ||
			strings.Contains(nameLower, itemName) ||
			nameLower == itemName {
			target = &items[i]
			break
		}
	}

	if target == nil {
		m.AppendMessage(fmt.Sprintf("You don't have any '%s' to equip.", itemName), "error")
		return
	}

	if target.Slot == "" || target.Slot == "none" {
		m.AppendMessage(fmt.Sprintf("%s cannot be equipped.", target.Name), "error")
		return
	}

	m.callEquipAPI(target.ID, target.Name, target.Slot)
}

// callEquipAPI calls the REST equip endpoint and handles the response.
func (m *model) callEquipAPI(itemID int, itemName, slot string) {
	url := fmt.Sprintf("%s/equipment/%d/equip", RESTAPIBase, itemID)
	body := fmt.Sprintf(`{"character_id":%d}`, m.currentCharacterID)
	equipResp, err := httpPut(url, body)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error equipping item: %v", err), "error")
		return
	}
	defer equipResp.Body.Close()

	var result struct {
		Message  string   `json:"message"`
		Messages []string `json:"messages"`
		Error    string   `json:"error"`
	}
	json.NewDecoder(equipResp.Body).Decode(&result)

	if equipResp.StatusCode != http.StatusOK {
		errMsg := result.Error
		if errMsg == "" {
			errMsg = "Failed to equip item"
		}
		m.AppendMessage(errMsg, "error")
		return
	}

	m.effectsService.FireEvent("on_equip", m.currentCharacterID, "", map[string]interface{}{
		"item_id": itemID,
		"slot":    slot,
	})

	slotName := formatSlotName(slot)
	if len(result.Messages) > 0 {
		output := ""
		for _, msg := range result.Messages {
			output += msg + ". "
		}
		output += fmt.Sprintf("You equip the %s in your %s.", itemName, slotName)
		m.AppendMessage(output, "success")
	} else {
		m.AppendMessage(fmt.Sprintf("You equip the %s in your %s.", itemName, slotName), "success")
	}
}