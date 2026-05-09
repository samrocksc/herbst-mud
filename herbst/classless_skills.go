package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// AbilityData represents an ability fetched from the server
type AbilityData struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Slot        int          `json:"slot"`
	Cooldown    int          `json:"cooldown"`
	ManaCost    int          `json:"manaCost"`
	StaminaCost int          `json:"staminaCost"`
	Effects     []EffectData `json:"effects"`
}

// EffectData represents a single effect from the AbilityEffect entity
type EffectData struct {
	EffectType    string  `json:"effectType"`
	DamageSubtype string  `json:"damageSubtype"`
	Target        string  `json:"target"`
	Value         int     `json:"value"`
	Duration      int     `json:"duration"`
	ScalingStat   string  `json:"scalingStat"`
	ScalingRatio  float64 `json:"scalingRatio"`
	SortOrder     int     `json:"sortOrder"`
}

// ActiveEffect tracks an active effect in combat
type ActiveEffect struct {
	AbilityID   int    `json:"abilityId"`
	Name        string `json:"name"`
	EffectType  string `json:"effectType"`
	Target       string `json:"target"`
	Duration    int    `json:"duration"`
	Modifier    int    `json:"modifier"`
	AppliedAt   int    `json:"appliedAt"`
}

// CombatSkillState tracks active effects and cooldowns
type CombatSkillState struct {
	ActiveEffects []ActiveEffect `json:"activeEffects"`
	Cooldowns     map[int]int    `json:"cooldowns"`
	EquippedSkill [5]AbilityData `json:"equippedSkills"`
}

// initCombatSkillState initializes combat ability state
func (m *model) initCombatSkillState() {
	if m.combatSkills == nil {
		m.combatSkills = &CombatSkillState{
			ActiveEffects: make([]ActiveEffect, 0),
			Cooldowns:     make(map[int]int),
		}
		m.loadEquippedAbilities()
	}
}

// loadEquippedAbilities fetches the character's equipped abilities from the server
func (m *model) loadEquippedAbilities() {
	if m.currentCharacterID == 0 {
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/classless-skills", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var result struct {
		Skills []AbilityData `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	for _, ability := range result.Skills {
		if ability.Slot >= 1 && ability.Slot <= 5 {
			m.combatSkills.EquippedSkill[ability.Slot-1] = ability
		}
	}
}

// useAbilitySlot executes a combat ability in the given slot
func (m *model) useAbilitySlot(slot int) bool {
	if slot < 1 || slot > 5 {
		m.addCombatLog("Invalid slot")
		return false
	}

	ability := m.combatSkills.EquippedSkill[slot-1]
	if ability.ID == 0 {
		m.addCombatLog(fmt.Sprintf("Slot %d is empty", slot))
		return false
	}

	if cd, ok := m.combatSkills.Cooldowns[ability.ID]; ok && cd > 0 {
		m.addCombatLog(fmt.Sprintf("%s is on cooldown (%d rounds)", ability.Name, cd))
		return false
	}

	if ability.ManaCost > 0 && m.characterMana < ability.ManaCost {
		m.addCombatLog(fmt.Sprintf("Not enough mana for %s", ability.Name))
		return false
	}
	if ability.StaminaCost > 0 && m.characterStamina < ability.StaminaCost {
		m.addCombatLog(fmt.Sprintf("Not enough stamina for %s", ability.Name))
		return false
	}

	m.characterMana -= ability.ManaCost
	m.characterStamina -= ability.StaminaCost

	m.resolveAbilityEffects(ability)
	m.combatSkills.Cooldowns[ability.ID] = ability.Cooldown
	return true
}

// resolveAbilityEffects applies all effects of an ability generically
func (m *model) resolveAbilityEffects(ability AbilityData) {
	for _, effect := range ability.Effects {
		modifier := m.calcScalingModifier(effect.ScalingStat, effect.ScalingRatio, effect.Value)

		active := ActiveEffect{
			AbilityID:  ability.ID,
			Name:       ability.Name,
			EffectType: effect.EffectType,
			Target:      effect.Target,
			Duration:   effect.Duration,
			Modifier:   modifier,
			AppliedAt:  m.currentTick(),
		}

		switch effect.EffectType {
		case "damage":
			m.applyDamageEffect(ability, effect, modifier)
		case "heal":
			m.applyHealEffect(ability, effect, modifier)
		case "stun":
			m.applyStunEffect(ability, effect)
		case "set_bind_point":
			m.setBindPoint()
		default:
			m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, active)
			m.addCombatLog(fmt.Sprintf("%s: %s applied for %d rounds", ability.Name, effect.EffectType, effect.Duration))
		}
	}
}

// calcScalingModifier computes the stat-scaled modifier for an effect
func (m *model) calcScalingModifier(stat string, ratio float64, baseValue int) int {
	if stat == "" || ratio == 0 {
		return baseValue
	}
	statVal := m.getStatValue(stat)
	statMod := (statVal - 10) / 2
	if statMod < 0 {
		statMod = 0
	}
	scaled := int(float64(baseValue) + float64(statMod)*ratio*float64(baseValue))
	if scaled < 1 && baseValue > 0 {
		scaled = baseValue
	}
	return scaled
}

// getStatValue returns a stat value by name
func (m *model) getStatValue(stat string) int {
	switch stat {
	case "strength":
		return m.getCharacterStrength()
	case "dexterity":
		return m.characterLevel*2 + 8
	case "constitution":
		return 10 + m.characterLevel
	case "intelligence":
		return 10 + m.characterLevel/2
	case "wisdom":
		return m.getCharacterWisdom()
	default:
		return 10
	}
}

// currentTick returns the current combat tick
func (m *model) currentTick() int {
	if c, ok := m.combatManager.GetCombat(m.combatID); ok {
		return c.CurrentTick
	}
	return 0
}

// applyDamageEffect applies a damage effect to the target
func (m *model) applyDamageEffect(ability AbilityData, effect EffectData, damage int) {
	if effect.Target == "enemy" && m.combatTarget != nil {
		applyDamageToCharacter(m.combatTarget.ID, damage)
		m.combatTarget.HP -= damage
		if m.combatTarget.HP < 0 {
			m.combatTarget.HP = 0
		}
		m.addCombatLog(fmt.Sprintf("%s deals %d damage!", ability.Name, damage))
		if m.combatTarget.HP <= 0 {
			m.handleTargetDefeat()
		}
	} else {
		m.addCombatLog(fmt.Sprintf("%s: %d damage to %s", ability.Name, damage, effect.Target))
	}
}

// applyHealEffect applies a heal effect
func (m *model) applyHealEffect(ability AbilityData, effect EffectData, amount int) {
	if effect.Target == "self" {
		m.characterHP += amount
		if m.characterHP > m.characterMaxHP {
			m.characterHP = m.characterMaxHP
		}
		healCharacter(m.currentCharacterID, amount)
		m.addCombatLog(fmt.Sprintf("%s heals %d HP!", ability.Name, amount))
	} else {
		m.addCombatLog(fmt.Sprintf("%s: %d heal to %s", ability.Name, amount, effect.Target))
	}
}

// applyStunEffect applies a stun effect with a stat contest
func (m *model) applyStunEffect(ability AbilityData, effect EffectData) {
	if effect.Target != "enemy" {
		m.addCombatLog(fmt.Sprintf("%s: stun applied to %s", ability.Name, effect.Target))
		return
	}

	dexMod := m.getDexModifier()
	targetCon := 10
	if m.combatTarget != nil {
		targetCon = m.getTargetConstitution()
	}

	success := dexMod + rollDie(6) > (targetCon-10)/2 + rollDie(6)

	if success {
		active := ActiveEffect{
			AbilityID:  ability.ID,
			Name:       ability.Name,
			EffectType: "stun",
			Target:      "enemy",
			Duration:   effect.Duration,
			AppliedAt:  m.currentTick(),
		}
		m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, active)
		m.addCombatLog(fmt.Sprintf("%s! Target is stunned for %d round!", ability.Name, effect.Duration))
	} else {
		m.addCombatLog(fmt.Sprintf("%s missed! Target resisted.", ability.Name))
	}
}

// getCharacterWisdom fetches WIS from server
func (m *model) getCharacterWisdom() int {
	if m.currentCharacterID == 0 {
		return 10
	}
	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/stats", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return 10
	}
	defer resp.Body.Close()
	var stats struct {
		Wisdom int `json:"wisdom"`
	}
	if json.NewDecoder(resp.Body).Decode(&stats) != nil {
		return 10
	}
	return stats.Wisdom
}

// getTargetConstitution estimates target CON
func (m *model) getTargetConstitution() int {
	if m.combatTarget == nil {
		return 10
	}
	return 10 + m.combatTarget.Level
}

// rollDie rolls a simple dN die
func rollDie(sides int) int {
	return 1 + rand.Intn(sides)
}

// setBindPoint updates the character's respawn room to the current room
func (m *model) setBindPoint() {
	url := fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"respawnRoomId": %d}`, m.currentRoom)

	resp, err := httpPut(url, payload)
	if err != nil {
		m.AppendMessage("Failed to set bind point.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		m.AppendMessage("Failed to set bind point.", "error")
		return
	}

	m.AppendMessage(fmt.Sprintf("✦ Bind point set at %s! You will respawn here.", m.roomName), "success")
}