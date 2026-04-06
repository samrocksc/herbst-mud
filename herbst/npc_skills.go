package main

import (
	"fmt"
	"math"
	"math/rand"
)

// NPCSkill represents an NPC's special ability
type NPCSkill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinCooldown int    `json:"minCooldown"` // Minimum ticks between uses
}

// Skill registry - all NPC skills
var NPCSkillRegistry = map[string]NPCSkill{
	"druid_heal": {
		ID:          "druid_heal",
		Name:        "Nature's Blessing",
		Description: "Heals 5% of max HP using druidic magic",
		MinCooldown: 4, // Minimum 4 ticks between uses
	},
}

// NPCSkillUsageChance calculates the probability of an NPC using their skill
// Hybrid algorithm: base chance + bonus based on missing health
func NPCSkillUsageChance(currentHP, maxHP int) float64 {
	if maxHP <= 0 {
		return 0
	}
	
	// Base chance: 10% per tick
	baseChance := 0.10
	
	// Health percentage (0.0 to 1.0 where 1.0 = full health)
	healthPercent := float64(currentHP) / float64(maxHP)
	
	// Missing health multiplier (0.0 to 1.0 where 1.0 = dead)
	missingHealth := 1.0 - healthPercent
	
	// Health component: scales from 0% to 40% based on missing health
	// At 100% HP: 0% bonus
	// At 50% HP: 20% bonus  
	// At 0% HP: 40% bonus (but can't use skill when defeated)
	healthComponent := missingHealth * 0.40
	
	// Total chance: base + health component
	// Max: 50% at very low health
	// Min: 10% at full health
	totalChance := baseChance + healthComponent
	
	return math.Min(totalChance, 0.50) // Cap at 50% chance per tick
}

// executeNPCSkill performs the skill effect
func (m *model) executeNPCSkill(skillID string, targetNPC *RoomCharacter) bool {
	skill, exists := NPCSkillRegistry[skillID]
	if !exists {
		return false
	}
	
	switch skill.ID {
	case "druid_heal":
		// Calculate 5% of max HP
		healAmount := int(float64(targetNPC.MaxHP) * 0.05)
		if healAmount < 1 {
			healAmount = 1
		}
		
		// Apply heal via server
		healCharacter(targetNPC.ID, healAmount)
		targetNPC.HP = min(targetNPC.HP+healAmount, targetNPC.MaxHP)
		
		m.addCombatLog(fmt.Sprintf("🌿 %s uses %s! Heals for %d HP (%d/%d)", 
			targetNPC.Name, skill.Name, healAmount, targetNPC.HP, targetNPC.MaxHP))
		
		return true
	}
	
	return false
}

// attemptNPCSkill checks if NPC should use skill this tick
// Returns true if skill was used
func (m *model) attemptNPCSkill(npcChar *RoomCharacter) bool {
	// Get full character data from server to check skill ID and cooldown
	// For now, we assume the NPC data is already loaded
	
	// Check if NPC has a skill - we'll need to fetch from character data
	// This is a simplified version - in reality we'd check the character's npc_skill_id field
	
	// Only Aragorn has druid heal for now
	if npcChar.Name != "Aragorn" {
		return false
	}
	
	skillID := "druid_heal"
	skill, exists := NPCSkillRegistry[skillID]
	if !exists {
		return false
	}
	
	// Check cooldown
	if m.npcSkillCooldown > 0 {
		return false
	}
	
	// Calculate usage chance
	chance := NPCSkillUsageChance(npcChar.HP, npcChar.MaxHP)
	
	// Roll for skill usage
	if rand.Float64() < chance {
		// Execute the skill
		if m.executeNPCSkill(skillID, npcChar) {
			// Set cooldown after successful use
			m.npcSkillCooldown = skill.MinCooldown
			return true
		}
	}
	
	return false
}

// decrementNPCSkillCooldown decrements the cooldown each tick
func (m *model) decrementNPCSkillCooldown() {
	if m.npcSkillCooldown > 0 {
		m.npcSkillCooldown--
	}
}
