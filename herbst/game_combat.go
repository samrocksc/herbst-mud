package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"herbst/combat"
	"herbst/dice"
)

// getCharacterStrength returns the character's strength stat
func (m *model) getCharacterStrength() int {
	if m.currentCharacterID == 0 {
		return 10 // default
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/stats", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return 10 // default
	}
	defer resp.Body.Close()

	var stats struct {
		Strength int `json:"strength"`
	}
	if json.NewDecoder(resp.Body).Decode(&stats) != nil {
		return 10 // default
	}

	return stats.Strength
}

// startCombat initializes a combat encounter
func (m *model) startCombat(target *RoomCharacter) {
	m.inCombat = true
	m.combatTarget = target
	m.combatLog = []string{}
	m.combatQueuedAction = ""
	m.combatJustStarted = true // Signal Update to start tick
	m.screen = ScreenCombat

	// Initialize combat manager if not already done
	if m.combatManager == nil {
		tickLoop := combat.NewTickLoop(combat.DefaultTickInterval)
		m.combatManager = combat.NewCombatManager(tickLoop)
	}

	// Create a new combat encounter
	m.combatID = m.combatManager.CreateCombat(0)

	// Add player as participant
	playerParticipant := combat.NewParticipant(
		m.currentCharacterID,
		m.currentCharacterName,
		m.characterHP,
	)
	m.combatManager.AddParticipant(m.combatID, playerParticipant)

	// Add target as participant
	targetParticipant := combat.NewParticipant(
		target.ID,
		target.Name,
		target.HP,
	)
	m.combatManager.AddParticipant(m.combatID, targetParticipant)

	// Load classless skills for combat (slots 1-5)
	m.initCombatSkillState()

	m.AppendMessage(fmt.Sprintf("⚔ You enter combat with %s!", target.Name), "combat")
	m.addCombatLog(fmt.Sprintf("Combat started with %s (Level %d)", target.Name, target.Level))
	m.addCombatLog("⏱ Tick combat begins - actions queue for next tick")
}

// combatAction queues a combat action for the next tick
func (m *model) combatAction(action string) {
	if !m.inCombat || m.combatTarget == nil {
		m.exitCombat()
		return
	}

	m.queueCombatAction(action)
}

// performCombatAttack handles the attack action with dice rolls
func (m *model) performCombatAttack() {
	if m.combatTarget == nil || m.combatTarget.HP <= 0 {
		return
	}

	// Get modifiers
	dexMod := m.getDexModifier()
	strMod := m.getStrModifier()

	// Roll to hit (d20 + DEX vs target AC)
	roll, toHit, isCrit, isFumble := dice.RollWithCrit(dexMod)
	targetAC := m.calculateAC(m.combatTarget)

	if isFumble {
		// Critical miss - automatic failure
		m.addCombatLog(fmt.Sprintf("🎲 FUMBLE! Natural 1 - You stumble badly!"))
		return
	}

	if toHit < targetAC && !isCrit {
		// Miss
		m.addCombatLog(fmt.Sprintf("🎲 Miss! (d20=%d + %d DEX = %d vs AC %d)",
			roll, dexMod, toHit, targetAC))
		return
	}

	// Hit - roll damage (d6 + STR)
	damageRoll, damage := dice.Roll(6, 1, strMod)
	if damage < 1 {
		damage = 1
	}

	if isCrit {
		damage *= 2
		m.addCombatLog(fmt.Sprintf("🎲 CRITICAL HIT! Natural 20!"))
	}

	// Apply damage
	applyDamageToCharacter(m.combatTarget.ID, damage)
	m.combatTarget.HP -= damage
	if m.combatTarget.HP < 0 {
		m.combatTarget.HP = 0
	}

	// Log the hit
	if isCrit {
		m.addCombatLog(fmt.Sprintf("⚔ Critical hit! %d damage (d6=%d + %d STR, doubled!)",
			damage, damageRoll, strMod))
	} else {
		m.addCombatLog(fmt.Sprintf("⚔ Hit! %d damage (d6=%d + %d STR)",
			damage, damageRoll, strMod))
	}

	// Check if target is defeated
	if m.combatTarget.HP <= 0 {
		m.addCombatLog(fmt.Sprintf("✦ %s has been defeated!", m.combatTarget.Name))
		m.AppendMessage(fmt.Sprintf("⚔ You defeated %s!", m.combatTarget.Name), "success")
		
		// Generate corpse
		// Generate corpse with victim's equipment
		m.generateCorpse(m.combatTarget)
		
		
		// Heal NPC back to max HP so they can be fought again (respawn)
		if m.combatTarget.IsNPC {
			healCharacter(m.combatTarget.ID, m.combatTarget.MaxHP)
		}
		m.exitCombat()
	}
}

// getDexModifier returns the DEX modifier for the player
func (m *model) getDexModifier() int {
	// Simple formula: level-based for now, can be enhanced later
	return m.characterLevel / 3
}

// getStrModifier returns the STR modifier for the player
func (m *model) getStrModifier() int {
	strength := m.getCharacterStrength()
	// D&D-style: 10 is 0, each +2 is +1 modifier
	return (strength - 10) / 2
}

// calculateAC returns target's armor class
func (m *model) calculateAC(target *RoomCharacter) int {
	// Base AC = 10 + level/2 (simple formula)
	baseAC := 10 + target.Level/2
	return baseAC
}

// performCombatDefend handles the defend action
func (m *model) performCombatDefend() {
	m.addCombatLog("🛡 You take a defensive stance!")
	// Defense reduces damage on the next enemy hit
	// This is handled in enemyTurnWithDice via a defending flag
}

// performCombatSkill handles skill actions (placeholder)
func (m *model) performCombatSkill() {
	m.addCombatLog("⚡ You attempt to use a skill... (not implemented)")
}

// performCombatItem handles item usage in combat (placeholder)
func (m *model) performCombatItem() {
	m.addCombatLog("🎒 You reach for an item... (not implemented)")
}

// enemyTurn handles the enemy's combat turn (legacy, use enemyTurnWithDice)
func (m *model) enemyTurn() {
	m.enemyTurnWithDice()
}

// enemyTurnWithDefense handles enemy turn with defense consideration
func (m *model) enemyTurnWithDefense(defending bool) {
	if m.combatTarget == nil || m.combatTarget.HP <= 0 {
		return
	}

	// Enemy rolls to hit
	enemyDexMod := m.combatTarget.Level / 3
	roll, toHit, isCrit, isFumble := dice.RollWithCrit(enemyDexMod)
	playerAC := 10 + m.getDexModifier()

	if defending {
		playerAC += 5 // Defense bonus
	}

	if isFumble {
		m.addCombatLog(fmt.Sprintf("🎲 %s FUMBLES! (rolled 1)", m.combatTarget.Name))
		return
	}

	if toHit < playerAC && !isCrit {
		m.addCombatLog(fmt.Sprintf("🎲 %s misses! (d20=%d + %d = %d vs AC %d)",
			m.combatTarget.Name, roll, enemyDexMod, toHit, playerAC))
		return
	}

	// Enemy hits - roll damage
	baseDamage := m.combatTarget.Level + 2
	damage := baseDamage
	if isCrit {
		damage *= 2
		m.addCombatLog(fmt.Sprintf("🎲 %s CRITS!", m.combatTarget.Name))
	}

	// Apply damage
	applyDamageToCharacter(m.currentCharacterID, damage)
	m.characterHP -= damage
	if m.characterHP < 0 {
		m.characterHP = 0
	}

	if defending {
		m.addCombatLog(fmt.Sprintf("⚔ %s attacks! Blocked for %d damage!",
			m.combatTarget.Name, damage))
	} else if isCrit {
		m.addCombatLog(fmt.Sprintf("⚔ %s critical hit! %d damage!",
			m.combatTarget.Name, damage))
	} else {
		m.addCombatLog(fmt.Sprintf("⚔ %s hits for %d damage!",
			m.combatTarget.Name, damage))
	}

	// Check if player is defeated
	if m.characterHP <= 0 {
		m.addCombatLog("☠ You have been defeated!")
		m.AppendMessage("☠ You have been defeated! Respawning...", "error")
		
		// Generate player corpse in current room with equipment
		playerAsChar := &RoomCharacter{
			ID:    m.currentCharacterID,
			Name:  m.currentCharacterName,
			HP:    0,
			MaxHP: m.characterMaxHP,
			IsNPC: false,
		}
		m.generateCorpse(playerAsChar)
		
		m.exitCombat()
		m.respawnPlayer()
	}
}

// combatFlee attempts to flee from combat
func (m *model) combatFlee() {
	if !m.inCombat {
		return
	}

	// Roll to escape - d20 vs difficulty
	roll, total := dice.D20(m.characterLevel / 2)
	escapeDC := 12 // Escape difficulty

	if total >= escapeDC {
		m.addCombatLog(fmt.Sprintf("🏃 Escape successful! (d20=%d + %d = %d vs DC %d)",
			roll, m.characterLevel/2, total, escapeDC))
		m.AppendMessage("🏃 You fled from combat!", "success")
		m.exitCombat()
	} else {
		m.addCombatLog(fmt.Sprintf("🏃 Escape failed! (d20=%d + %d = %d vs DC %d)",
			roll, m.characterLevel/2, total, escapeDC))
		m.AppendMessage("🏃 You couldn't escape!", "error")
	}
}

// processCombatTick executes at each combat tick interval
func (m *model) processCombatTick() {
	if !m.inCombat || m.combatTarget == nil {
		return
	}

	// Decrement cooldowns at the start of each tick
	m.decrementCooldowns()

	// Execute queued action or auto-attack
	if m.combatQueuedAction != "" {
		m.executeCombatAction(m.combatQueuedAction)
		m.combatQueuedAction = ""
	} else {
		// No action = auto-attack
		m.addCombatLog("⏱ No action queued - auto-attacking!")
		m.performCombatAttack()
	}

	// Check if combat ended after player action
	if !m.inCombat {
		return
	}

	// Enemy turn (with dice rolls)
	m.enemyTurnWithDice()

	// Check if combat ended after enemy action
	m.checkCombatEnd()
}

// executeCombatAction executes a queued combat action
func (m *model) executeCombatAction(action string) {
	// Check if it's a talent name (from equipped slots)
	for _, talent := range m.combatTalents {
		if talent.Name == action {
			m.performTalentAction(talent)
			return
		}
	}

	// Default actions
	switch action {
	case "attack":
		m.performCombatAttack()
	case "defend":
		m.performCombatDefend()
	default:
		// Unknown action - default to attack
		m.performCombatAttack()
	}
}

// performTalentAction executes a talent's effect in combat
func (m *model) performTalentAction(talent EquippedTalent) {
	if m.combatTarget == nil || m.combatTarget.HP <= 0 {
		return
	}

	// Log the talent use
	m.addCombatLog(fmt.Sprintf("⚔ Using %s!", talent.Name))

	// Determine talent effect based on type
	switch talent.EffectType {
	case "damage", "":
		// Damage talent - apply damage to target
		damage := talent.EffectValue
		if damage <= 0 {
			// Default attack if no damage value
			m.performCombatAttack()
			return
		}
		applyDamageToCharacter(m.combatTarget.ID, damage)
		m.combatTarget.HP -= damage
		if m.combatTarget.HP < 0 {
			m.combatTarget.HP = 0
		}
		m.addCombatLog(fmt.Sprintf("⚔ %s deals %d damage!", talent.Name, damage))
		if m.combatTarget.HP <= 0 {
			m.addCombatLog(fmt.Sprintf("✦ %s has been defeated!", m.combatTarget.Name))
			m.AppendMessage(fmt.Sprintf("⚔ You defeated %s!", m.combatTarget.Name), "success")
			
			// Generate corpse with victim's equipment
			m.generateCorpse(m.combatTarget)
			
			if m.combatTarget.IsNPC {
				healCharacter(m.combatTarget.ID, m.combatTarget.MaxHP)
			}
			m.exitCombat()
		}
	case "heal":
		// Healing talent
		healAmt := talent.EffectValue
		if healAmt <= 0 {
			healAmt = 15 // Default heal
		}
		m.characterHP += healAmt
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, healAmt)
		m.addCombatLog(fmt.Sprintf("💚 Healed for %d HP!", healAmt))
	case "dot":
		// Damage over time - instant for now, can be expanded later
		damage := talent.EffectValue
		if damage > 0 {
			applyDamageToCharacter(m.combatTarget.ID, damage)
			m.combatTarget.HP -= damage
			if m.combatTarget.HP < 0 {
				m.combatTarget.HP = 0
			}
			m.addCombatLog(fmt.Sprintf("☠ %s applies DoT: %d damage!", talent.Name, damage))
		}
	case "buff_armor":
		// Armor buff - would need buff system
		m.addCombatLog(fmt.Sprintf("🛡 %s increases defense by %d!", talent.Name, talent.EffectValue))
	case "buff_dodge":
		// Dodge buff - would need buff system
		m.addCombatLog(fmt.Sprintf("💨 %s increases dodge by %d%%!", talent.Name, talent.EffectValue))
	case "buff_crit":
		// Crit buff - would need buff system
		m.addCombatLog(fmt.Sprintf("🎯 %s increases critical chance by %d%%!", talent.Name, talent.EffectValue))
	case "debuff":
		// Debuff enemy
		m.addCombatLog(fmt.Sprintf("🔻 %s weakens the enemy!", talent.Name))
	default:
		// Unknown type - default to attack
		m.performCombatAttack()
	}
}

// queueCombatAction queues an action for the next tick
func (m *model) queueCombatAction(action string) {
	m.combatQueuedAction = action
	m.addCombatLog(fmt.Sprintf("⏱ Queued: %s", action))
}

// enemyTurnWithDice handles enemy turn with dice rolls
func (m *model) enemyTurnWithDice() {
	if m.combatTarget == nil || m.combatTarget.HP <= 0 {
		return
	}

	// Attempt NPC skill usage (happens before regular attack)
	// This makes combat more dynamic - NPC might heal instead of attacking
	if m.attemptNPCSkill(m.combatTarget) {
		// NPC used a skill (like healing) - check if combat ended
		if !m.inCombat || m.combatTarget.HP <= 0 {
			return
		}
		// After using skill, NPC might still attack this tick
		// 50% chance to still attack after healing
		if rand.Float64() < 0.5 {
			return
		}
	}

	// Enemy rolls to hit player
	enemyDexMod := m.combatTarget.Level / 3 // Simple enemy stat approximation
	roll, toHit, isCrit, isFumble := dice.RollWithCrit(enemyDexMod)
	playerAC := 10 + m.getDexModifier() // Player AC = 10 + DEX mod

	if isFumble {
		m.addCombatLog(fmt.Sprintf("🎲 %s FUMBLES! (rolled 1)", m.combatTarget.Name))
		return
	}

	if toHit < playerAC && !isCrit {
		m.addCombatLog(fmt.Sprintf("🎲 %s misses! (d20=%d + %d = %d vs AC %d)",
			m.combatTarget.Name, roll, enemyDexMod, toHit, playerAC))
		return
	}

	// Enemy hits - roll damage
	baseDamage := m.combatTarget.Level + 2
	damage := baseDamage
	if isCrit {
		damage *= 2
		m.addCombatLog(fmt.Sprintf("🎲 %s CRITS!", m.combatTarget.Name))
	}

	// Apply damage
	applyDamageToCharacter(m.currentCharacterID, damage)
	m.characterHP -= damage
	if m.characterHP < 0 {
		m.characterHP = 0
	}

	if isCrit {
		m.addCombatLog(fmt.Sprintf("⚔ %s critical hit! %d damage!",
			m.combatTarget.Name, damage))
	} else {
		m.addCombatLog(fmt.Sprintf("⚔ %s hits for %d damage!",
			m.combatTarget.Name, damage))
	}

	// Check if player is defeated
	if m.characterHP <= 0 {
		m.addCombatLog("☠ You have been defeated!")
		m.AppendMessage("☠ You have been defeated! Respawning...", "error")
		
		// Generate player corpse in current room with equipment
		playerAsChar := &RoomCharacter{
			ID:    m.currentCharacterID,
			Name:  m.currentCharacterName,
			HP:    0,
			MaxHP: m.characterMaxHP,
			IsNPC: false,
		}
		m.generateCorpse(playerAsChar)
		
		m.exitCombat()
		m.respawnPlayer()
	}
}

// checkCombatEnd checks if combat should end
func (m *model) checkCombatEnd() {
	if m.combatTarget == nil || m.combatTarget.HP <= 0 {
		m.exitCombat()
	}
}

// exitCombat cleans up combat state
// decrementCooldowns reduces all cooldowns at the start of a tick
func (m *model) decrementCooldowns() {
	// Decrement player skill cooldowns
	if m.combatSkills != nil {
		for skillID, cd := range m.combatSkills.Cooldowns {
			if cd > 0 {
				m.combatSkills.Cooldowns[skillID] = cd - 1
			}
		}
	}
	
	// Decrement NPC skill cooldown
	if m.npcSkillCooldown > 0 {
		m.npcSkillCooldown--
	}
}

func (m *model) exitCombat() {
	if m.combatManager != nil && m.combatID > 0 {
		m.combatManager.EndCombat(m.combatID)
	}
	m.inCombat = false
	m.combatTarget = nil
	m.combatID = 0
	m.combatQueuedAction = ""
	m.combatJustStarted = false
	m.screen = ScreenPlaying
}

// addCombatLog adds a message to the combat log
func (m *model) addCombatLog(msg string) {
	timestamp := time.Now().Format("15:04:05")
	m.combatLog = append(m.combatLog, fmt.Sprintf("[%s] %s", timestamp, msg))
	// Keep log to last 50 messages
	if len(m.combatLog) > 50 {
		m.combatLog = m.combatLog[len(m.combatLog)-50:]
	}
}

// respawnPlayer handles player defeat
func (m *model) respawnPlayer() {
	// Reset HP to max
	m.characterHP = m.characterMaxHP

	// Fetch respawn room from server
	respawnRoomID := StartingRoomID // default fallback
	if m.currentCharacterID != 0 {
		resp, err := httpGet(fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID))
		if err == nil {
			defer resp.Body.Close()
			var char struct {
				RespawnRoomId int `json:"respawnRoomId"`
			}
			if json.NewDecoder(resp.Body).Decode(&char) == nil && char.RespawnRoomId > 0 {
				respawnRoomID = char.RespawnRoomId
			}
		}
	}

	// Move player to respawn room
	m.currentRoom = respawnRoomID

	// Reload room data from server
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), respawnRoomID)
		if err == nil {
			m.roomName = room.Name
			m.roomDesc = room.Description
			m.exits = room.Exits
		}
	}

	// Reload room contents
	m.loadRoomItems()
	m.loadRoomCharactersWithHP()

	// Heal player on server
	healCharacter(m.currentCharacterID, m.characterMaxHP)
	
	m.AppendMessage(fmt.Sprintf("☠ You respawn at %s!", m.roomName), "success")
}

// applyDamageToCharacter sends damage to the server
func applyDamageToCharacter(characterID, damage int) {
	url := fmt.Sprintf("%s/characters/%d/damage", RESTAPIBase, characterID)
	payload := fmt.Sprintf(`{"damage": %d}`, damage)

	resp, err := httpPost(url, payload)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// healCharacter sends heal request to the server
func healCharacter(characterID, amount int) {
	url := fmt.Sprintf("%s/characters/%d/heal", RESTAPIBase, characterID)
	payload := fmt.Sprintf(`{"amount": %d}`, amount)

	resp, err := httpPost(url, payload)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// getCharacterCombatStatus fetches combat status from server
func getCharacterCombatStatus(characterID int) (int, int, bool) {
	url := fmt.Sprintf("%s/characters/%d/combat-status", RESTAPIBase, characterID)

	resp, err := httpGet(url)
	if err != nil {
		return 0, 0, false
	}
	defer resp.Body.Close()

	var status struct {
		HP     int `json:"hp"`
		MaxHP  int `json:"maxHp"`
		IsNPC  bool `json:"isNPC"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return 0, 0, false
	}

	return status.HP, status.MaxHP, status.IsNPC
}

// loadRoomCharactersWithHP loads characters in the room with HP info
func (m *model) loadRoomCharactersWithHP() {
	m.loadRoomCharacters()
	// Refresh HP from server for NPCs
	for i := range m.roomCharacters {
		if m.roomCharacters[i].IsNPC {
			hp, maxHP, _ := getCharacterCombatStatus(m.roomCharacters[i].ID)
			m.roomCharacters[i].HP = hp
			m.roomCharacters[i].MaxHP = maxHP
		}
	}
}

// handleCombatInput handles keyboard input during combat
func (m *model) handleCombatInput(key string) {
	switch strings.ToLower(key) {
	case "1":
		m.useCombatTalent(1)
	case "2":
		m.useCombatTalent(2)
	case "3":
		m.useCombatTalent(3)
	case "4":
		m.useCombatTalent(4)
	case "r":
		m.useHealthPotion()
	case "q":
		m.combatFlee()
	}
}