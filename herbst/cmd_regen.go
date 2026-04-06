package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RegenConfig configures the regeneration system
type RegenConfig struct {
	TickInterval     time.Duration
	BaseHPRegenRate  int // HP regeneration base rate
	BaseSPRegenRate  int // Stamina regeneration rate (equal for all)
	BaseMPRegenRate  int // Mana regeneration rate (equal for all)
}

// DefaultRegenConfig returns default configuration
// All classless players get equal stamina and mana regen
func DefaultRegenConfig() RegenConfig {
	return RegenConfig{
		TickInterval:     time.Duration(6) * time.Second,
		BaseHPRegenRate:  1, // Modified by CON
		BaseSPRegenRate:  3, // Equal for all - 3 stamina per tick
		BaseMPRegenRate:  2, // Equal for all - 2 mana per tick
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

// performRegen regenerates HP, Stamina, and Mana for player while out of combat
// All classless players get equal stamina (3) and mana (2) regen rates
func (m *model) performRegen() {
	// Only regen when out of combat and has character loaded
	if m.inCombat {
		return
	}
	if m.currentCharacterID == 0 {
		return
	}
	if m.characterHP <= 0 {
		return
	}

	regenMessages := []string{}

	// HP Regen (Constitution-based, only if not at full HP)
	if m.characterHP < m.characterMaxHP {
		con := m.getCharacterConstitution()
		hpRegen := getConstitutionModifier(con)
		oldHP := m.characterHP
		m.characterHP += hpRegen
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, hpRegen)
		regenMessages = append(regenMessages, fmt.Sprintf("💚 +%d HP", hpRegen))

		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] HP Regen: +%d (CON %d) %d→%d",
				hpRegen, con, oldHP, m.characterHP), "info")
		}
	}

	// Stamina Regen (Equal for all classless players)
	if m.characterStamina < m.characterMaxStamina {
		oldSP := m.characterStamina
		spRegen := regenConfig.BaseSPRegenRate
		m.characterStamina += spRegen
		if m.characterStamina > m.characterMaxStamina {
			m.characterStamina = m.characterMaxStamina
		}
		// Send stamina update to server
		m.updateStaminaOnServer(spRegen)
		regenMessages = append(regenMessages, fmt.Sprintf("⚡ +%d SP", spRegen))

		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] SP Regen: +%d %d→%d",
				spRegen, oldSP, m.characterStamina), "info")
		}
	}

	// Mana Regen (Equal for all classless players)
	if m.characterMana < m.characterMaxMana {
		oldMP := m.characterMana
		mpRegen := regenConfig.BaseMPRegenRate
		m.characterMana += mpRegen
		if m.characterMana > m.characterMaxMana {
			m.characterMana = m.characterMaxMana
		}
		// Send mana update to server
		m.updateManaOnServer(mpRegen)
		regenMessages = append(regenMessages, fmt.Sprintf("💧 +%d MP", mpRegen))

		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] MP Regen: +%d %d→%d",
				mpRegen, oldMP, m.characterMana), "info")
		}
	}

	// Show combined regen message
	if len(regenMessages) > 0 {
		m.AppendMessage(fmt.Sprintf("🌿 Regen: %s", joinStrings(regenMessages, " ")), "success")
	}

	// Also heal NPCs in current room (HP only, same rate)
	con := m.getCharacterConstitution()
	hpRegen := getConstitutionModifier(con)
	m.healRoomNPCs(hpRegen)
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

// updateStaminaOnServer sends stamina regeneration to the server
func (m *model) updateStaminaOnServer(amount int) {
	url := fmt.Sprintf("%s/characters/%d/stamina", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"amount": %d}`, amount)

	resp, err := httpPost(url, payload)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to update stamina: %v", err), "error")
		}
		return
	}
	defer resp.Body.Close()
}

// updateManaOnServer sends mana regeneration to the server
func (m *model) updateManaOnServer(amount int) {
	url := fmt.Sprintf("%s/characters/%d/mana", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"amount": %d}`, amount)

	resp, err := httpPost(url, payload)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to update mana: %v", err), "error")
		}
		return
	}
	defer resp.Body.Close()
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
