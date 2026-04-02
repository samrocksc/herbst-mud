package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// passiveHealNPCs performs passive healing on NPCs in the current room
// This simulates NPCs naturally recovering while players aren't present
// Call this whenever entering a room or after significant time passes
func (m *model) passiveHealNPCs() {
	if m.currentRoom == 0 {
		return
	}

	// Heal NPCs for 10% of their max HP (simulating time passing)
	// This runs independently of player regen
	url := fmt.Sprintf("%s/rooms/%d/npcs/passive-heal", RESTAPIBase, m.currentRoom)

	resp, err := httpPost(url, `{}`)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Passive heal error: %v", err), "error")
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			Healed  int `json:"healed"`
			FullyHealed int `json:"fullyHealed"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if m.debugMode && result.Healed > 0 {
				m.AppendMessage(fmt.Sprintf("[DEBUG] %d NPCs recovered naturally (%d fully healed)", 
					result.Healed, result.FullyHealed), "info")
			}
		}
	}
}
