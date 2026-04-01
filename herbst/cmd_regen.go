package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RegenConfig configures the regeneration system
type RegenConfig struct {
	TickInterval  time.Duration
	BaseRegenRate int
}

// DefaultRegenConfig returns default configuration
func DefaultRegenConfig() RegenConfig {
	return RegenConfig{
		TickInterval:  time.Duration(6) * time.Second,
		BaseRegenRate: 1,
	}
}

// RegenState tracks regeneration timing
var (
	lastRegenTick time.Time
	regenConfig   = DefaultRegenConfig()
)

// initRegen initializes the regeneration system
func init() {
	lastRegenTick = time.Now()
}

// getConstitutionModifier calculates HP regen from CON
// D&D style: 10 CON = 2 HP per tick, 12 CON = 3 HP, 14 CON = 4 HP, etc.
func getConstitutionModifier(constitution int) int {
	if constitution <= 0 {
		return 1 // Minimum regen
	}
	// Every 2 CON above 8 adds 1 to base regen
	mod := (constitution - 8) / 2
	if mod < 1 {
		mod = 1
	}
	return mod
}

// getCharacterConstitution fetches CON from the server
func (m *model) getCharacterConstitution() int {
	if m.currentCharacterID == 0 {
		return 10 // default
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/stats", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return 10 // default
	}
	defer resp.Body.Close()

	var stats struct {
		Constitution int `json:"constitution"`
	}
	if json.NewDecoder(resp.Body).Decode(&stats) != nil {
		return 10 // default
	}

	return stats.Constitution
}

// performRegen regenerates HP for player and all NPCs in current room
func (m *model) performRegen() {
	// Only regen when:
	// - Out of combat
	// - Not at full health  
	// - Alive
	// - Has a character loaded
	if m.inCombat {
		return
	}
	if m.characterHP >= m.characterMaxHP {
		return
	}
	if m.characterHP <= 0 {
		return
	}
	if m.currentCharacterID == 0 {
		return
	}

	// Get CON and calculate regen amount
	con := m.getCharacterConstitution()
	regenAmount := getConstitutionModifier(con)

	// Apply regen to player
	oldHP := m.characterHP
	m.characterHP += regenAmount
	if m.characterHP > m.characterMaxHP {
		m.characterHP = m.characterMaxHP
	}

	// Send player heal to server
	healCharacter(m.currentCharacterID, regenAmount)

	// Also heal NPCs in current room (they regen at same rate)
	m.healRoomNPCs(regenAmount)

	// Show message every time
	m.AppendMessage(fmt.Sprintf("💚 +%d HP regen (CON %d)", regenAmount, con), "success")

	if m.debugMode {
		m.AppendMessage(
			fmt.Sprintf("[DEBUG] Regen: +%d HP (CON %d) %d to %d", 
				regenAmount, con, oldHP, m.characterHP), 
			"info",
		)
	}
}

// healRoomNPCs sends a heal request for all NPCs in the current room
func (m *model) healRoomNPCs(amount int) {
	if m.currentRoom == 0 {
		return
	}

	url := fmt.Sprintf("%s/rooms/%d/npcs/heal", RESTAPIBase, m.currentRoom)
	payload := fmt.Sprintf(`{"amount": %d}`, amount)

	resp, err := httpPost(url, payload)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to heal NPCs: %v", err), "error")
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			Healed int `json:"healed"`
			Amount int `json:"amount"`
		}
		if json.NewDecoder(resp.Body).Decode(&result) == nil && m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Healed %d NPCs for +%d HP", result.Healed, result.Amount), "info")
		}
	}
}
