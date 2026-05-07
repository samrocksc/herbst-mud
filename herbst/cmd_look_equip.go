package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// equippedItemData is the API response shape for equipped look items.
type equippedItemData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slot        string `json:"slot"`
	IsEquipped  bool   `json:"isEquipped"`
	ItemType    string `json:"itemType"`
	IsTwoHanded bool   `json:"is_two_handed"`
}

// fetchEquippedItems retrieves all equipped items for a character.
func fetchEquippedItems(charID int) []equippedItemData {
	url := fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, charID)
	resp, err := httpGet(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var items []equippedItemData
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil
	}

	var equipped []equippedItemData
	for _, item := range items {
		if item.IsEquipped {
			equipped = append(equipped, item)
		}
	}
	return equipped
}

// formatCharacterEquipment builds Wearing/Wielding lines for look output.
// Weapons (main_hand/off_hand slots) go on the Wielding line.
// All other equipped items go on the Wearing line.
func formatCharacterEquipment(items []equippedItemData) string {
	if len(items) == 0 {
		return ""
	}

	var wearing []string
	var mainHand, offHand string

	for _, item := range items {
		switch item.Slot {
		case "main_hand":
			mainHand = item.Name
		case "off_hand":
			offHand = item.Name
		default:
			wearing = append(wearing, item.Name)
		}
	}

	var lines []string

	if len(wearing) > 0 {
		lines = append(lines, fmt.Sprintf("Wearing: %s",
			strings.Join(wearing, ", ")))
	}

	var wielding string
	if mainHand != "" && offHand != "" {
		wielding = fmt.Sprintf("Wielding: %s and %s", mainHand, offHand)
	} else if mainHand != "" {
		wielding = fmt.Sprintf("Wielding: %s", mainHand)
	} else if offHand != "" {
		wielding = fmt.Sprintf("Wielding: %s", offHand)
	}
	if wielding != "" {
		lines = append(lines, wielding)
	}

	if len(lines) == 0 {
		return ""
	}

	return "\n" + strings.Join(lines, "\n")
}