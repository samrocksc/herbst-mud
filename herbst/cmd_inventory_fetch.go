package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// fetchRaceSlots returns equipment slots for the character's race.
func (m *model) fetchRaceSlots() []string {
	defaultSlots := []string{
		"head", "neck", "chest", "back", "hands",
		"legs", "feet", "finger_left", "finger_right",
		"main_hand", "off_hand",
	}

	raceName := m.characterRace
	if raceName == "" {
		raceName = "human"
	}

	resp, err := httpGet(fmt.Sprintf("%s/races", RESTAPIBase))
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return defaultSlots
	}
	defer resp.Body.Close()

	var races []raceData
	if err := json.NewDecoder(resp.Body).Decode(&races); err != nil {
		return defaultSlots
	}

	for _, r := range races {
		if r.Name == raceName && len(r.EquipmentSlots) > 0 {
			return r.EquipmentSlots
		}
	}
	return defaultSlots
}

// fetchInventoryItems loads all equipment owned by the character.
func (m *model) fetchInventoryItems() []inventoryItem {
	url := fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID)
	resp, err := httpGet(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return nil
	}
	defer resp.Body.Close()

	var items []inventoryItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil
	}
	return items
}