package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleUnequipItemCommand unequips an item by slot name or item name.
func (m *model) handleUnequipItemCommand(arg string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	// Fetch equipped items owned by this character
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

	// Find equipped item matching slot name or item name
	var target *equipItemData
	argLower := strings.ToLower(arg)
	for i := range items {
		if !items[i].IsEquipped {
			continue
		}
		// Match by slot name
		if items[i].Slot == argLower {
			target = &items[i]
			break
		}
		// Match by item name
		nameLower := strings.ToLower(items[i].Name)
		if fuzzyWordMatch(items[i].Name, arg) ||
			strings.Contains(nameLower, argLower) ||
			nameLower == argLower {
			target = &items[i]
			break
		}
	}

	if target == nil {
		m.AppendMessage(fmt.Sprintf("No equipped item found matching '%s'.", arg), "error")
		return
	}

	m.callUnequipAPI(target.ID, target.Name, target.Slot)
}

// callUnequipAPI calls the REST unequip endpoint and handles the response.
func (m *model) callUnequipAPI(itemID int, itemName, slot string) {
	url := fmt.Sprintf("%s/equipment/%d/unequip", RESTAPIBase, itemID)
	body := fmt.Sprintf(`{"character_id":%d}`, m.currentCharacterID)
	unequipResp, err := httpPut(url, body)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error unequipping item: %v", err), "error")
		return
	}
	defer unequipResp.Body.Close()

	var result struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	json.NewDecoder(unequipResp.Body).Decode(&result)

	if unequipResp.StatusCode != http.StatusOK {
		errMsg := result.Error
		if errMsg == "" {
			errMsg = "Failed to unequip item"
		}
		m.AppendMessage(errMsg, "error")
		return
	}

	slotName := formatSlotName(slot)
	m.AppendMessage(fmt.Sprintf("You unequip the %s from your %s.", itemName, slotName), "success")
}