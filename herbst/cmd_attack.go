package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	character "herbst/db/character"
)

// handleAttackCommand handles the attack/kill/a command
func (m *model) handleAttackCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Attack what? Usage: attack <target name>", "error")
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(strings.TrimSpace(target))

	// Load room characters if needed
	m.loadRoomCharacters()

	// Find target in room
	var foundChar *roomCharacter
	for i := range m.roomCharacters {
		char := &m.roomCharacters[i]
		charNameLower := strings.ToLower(char.Name)
		if fuzzyWordMatch(char.Name, target) || strings.Contains(charNameLower, target) || charNameLower == target {
			foundChar = char
			break
		}
	}

	if foundChar == nil {
		m.AppendMessage(fmt.Sprintf("You don't see any '%s' here to attack.", target), "error")
		return
	}

	// Attack the target
	m.performAttack(foundChar)
}

// performAttack executes an attack against a target
func (m *model) performAttack(target *roomCharacter) {
	// Get character stats for damage calculation
	strength := m.getCharacterStrength()

	// Calculate base damage (strength-based, classless)
	// Base damage: 1-3 + strength modifier
	// Every 5 strength adds 1 base damage
	baseDamage := 1 + (strength / 5)
	if baseDamage < 1 {
		baseDamage = 1
	}

	// Get equipped weapon damage bonus (for future expansion)
	weaponDamage := m.getWeaponDamage()

	// Total damage range
	minDamage := baseDamage + weaponDamage
	maxDamage := minDamage + 3 // 3 damage variance

	// Simple deterministic "random" based on character level for damage variance
	damage := minDamage
	if maxDamage > minDamage && m.characterLevel > 0 {
		damage = minDamage + (m.characterLevel % (maxDamage - minDamage + 1))
		if damage > maxDamage {
			damage = maxDamage
		}
	}

	// Apply damage to target
	if target.IsNPC {
		m.attackNPC(target, damage)
	} else {
		m.attackPlayer(target, damage)
	}
}

// getCharacterStrength returns the character's strength stat
func (m *model) getCharacterStrength() int {
	// Try to get strength from API
	if m.currentCharacterID == 0 {
		return 10 // default
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return 10 // default
	}
	defer resp.Body.Close()

	var charData struct {
		Strength int `json:"strength"`
	}
	if json.NewDecoder(resp.Body).Decode(&charData) != nil {
		return 10 // default
	}

	return charData.Strength
}

// getWeaponDamage returns the damage bonus from equipped weapon
func (m *model) getWeaponDamage() int {
	// For now, return 0 - weapon damage would be fetched from equipped items
	// This can be enhanced later to query equipped weapons
	return 0
}

// attackNPC handles attacking an NPC target
func (m *model) attackNPC(target *roomCharacter, damage int) {
	// Get NPC stats from API
	resp, err := httpGet(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("⚔ You attack %s but nothing happens.", target.Name), "combat")
		return
	}
	defer resp.Body.Close()

	var npcs []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		HP          int    `json:"hp"`
		MaxHP       int    `json:"maxHp"`
		Disposition string `json:"disposition"`
	}
	if json.NewDecoder(resp.Body).Decode(&npcs) != nil {
		m.AppendMessage(fmt.Sprintf("⚔ You attack %s but nothing happens.", target.Name), "combat")
		return
	}

	// Find the specific NPC
	var npc *struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		HP          int    `json:"hp"`
		MaxHP       int    `json:"maxHp"`
		Disposition string `json:"disposition"`
	}
	for i := range npcs {
		if strings.ToLower(npcs[i].Name) == strings.ToLower(target.Name) {
			npc = &npcs[i]
			break
		}
	}

	if npc == nil {
		m.AppendMessage(fmt.Sprintf("⚔ You attack %s dealing %d damage!", target.Name, damage), "combat")
		return
	}

	// Display combat result
	if npc.HP > 0 {
		remainingHP := npc.HP - damage
		if remainingHP < 0 {
			remainingHP = 0
		}
		m.AppendMessage(fmt.Sprintf("⚔ You attack %s (Level %d) for %d damage! [%d/%d HP]",
			target.Name, npc.Level, damage, remainingHP, npc.MaxHP), "combat")
	} else {
		m.AppendMessage(fmt.Sprintf("⚔ %s is already defeated!", target.Name), "combat")
	}
}

// attackPlayer handles attacking another player (PvP)
func (m *model) attackPlayer(target *roomCharacter, damage int) {
	// PvP is enabled - show combat message
	m.AppendMessage(fmt.Sprintf("⚔ You attack %s for %d damage! (PvP)", target.Name, damage), "combat")

	// In a full implementation, this would:
	// 1. Send PvP attack to API
	// 2. Update target's HP
	// 3. Handle target's response (counter-attack, flee, etc.)
}

// getCharacterFromDB fetches character from database (for future use)
func (m *model) getCharacterFromDB() (int, error) {
	if m.client == nil || m.currentCharacterID == 0 {
		return 0, fmt.Errorf("no character loaded")
	}

	ctx := context.Background()
	char, err := m.client.Character.Query().
		Where(character.IDEQ(m.currentCharacterID)).
		Only(ctx)
	if err != nil {
		return 0, err
	}

	return char.Strength, nil
}