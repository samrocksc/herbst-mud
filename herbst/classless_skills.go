package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// ClasslessSkill represents one of the 5 swappable classless skills
type ClasslessSkill struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Slot          int    `json:"slot"`         // 1-5 for classless skills
	EffectType    string `json:"effectType"`   // concentrate|haymaker|backoff|scream|slap
	Cooldown      int    `json:"cooldown"`     // in ticks
	ManaCost      int    `json:"manaCost"`
	StaminaCost   int    `json:"staminaCost"`
	BaseStat      string `json:"baseStat"`     // wisdom|strength|dexterity
	 Duration      int    `json:"duration"`     // effect duration in rounds
	IsPassive     bool   `json:"isPassive"`    // true for active skills
}

// SkillEffect tracks active skill effects in combat
type SkillEffect struct {
	SkillID      int       `json:"skillId"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Duration     int       `json:"duration"`     // remaining rounds
	Modifier     int       `json:"modifier"`     // bonus/penalty value
	AppliedAt    int       `json:"appliedAt"`    // tick when applied
}

// ============================================================
// SKILL DEFINITIONS
// ============================================================

// The 5 classless skills available to all characters
var ClasslessSkills = []ClasslessSkill{
	{
		ID:          100,
		Name:        "Concentrate",
		Description: "Focus your mind to increase accuracy. +WIS to hit for 4 rounds.",
		Slot:        1,
		EffectType:  "concentrate",
		Cooldown:    8,
		ManaCost:    10,
		StaminaCost: 0,
		BaseStat:    "wisdom",
		Duration:    4,
		IsPassive:   false,
	},
	{
		ID:          101,
		Name:        "Haymaker",
		Description: "A powerful but reckless strike. +STR to damage, -DEX to hit.",
		Slot:        2,
		EffectType:  "haymaker",
		Cooldown:    6,
		ManaCost:    0,
		StaminaCost: 15,
		BaseStat:    "strength",
		Duration:    1,
		IsPassive:   false,
	},
	{
		ID:          102,
		Name:        "Back-off",
		Description: "Use agility to dodge all attacks this round. Costs stamina.",
		Slot:        3,
		EffectType:  "backoff",
		Cooldown:    10,
		ManaCost:    0,
		StaminaCost: 25,
		BaseStat:    "dexterity",
		Duration:    1,
		IsPassive:   false,
	},
	{
		ID:          103,
		Name:        "Scream",
		Description: "Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.",
		Slot:        4,
		EffectType:  "scream",
		Cooldown:    12,
		ManaCost:    5,
		StaminaCost: 10,
		BaseStat:    "constitution", // affects how well you handle the rage
		Duration:    2,
		IsPassive:   false,
	},
	{
		ID:          104,
		Name:        "Slap",
		Description: "A quick stunning strike. DEX vs CON to stun for 1 round.",
		Slot:        5,
		EffectType:  "slap",
		Cooldown:    8,
		ManaCost:    0,
		StaminaCost: 12,
		BaseStat:    "dexterity",
		Duration:    1,
		IsPassive:   false,
	},
}

// CombatState tracks active skill effects and cooldowns
type CombatSkillState struct {
	ActiveEffects []SkillEffect      `json:"activeEffects"`
	Cooldowns     map[int]int        `json:"cooldowns"`     // skillID: remainingCooldown
	EquippedSkill [5]ClasslessSkill   `json:"equippedSkills"` // slots 1-5
}

// Initialize combat skill state
func (m *model) initCombatSkillState() {
	if m.combatSkills == nil {
		m.combatSkills = &CombatSkillState{
			ActiveEffects: make([]SkillEffect, 0),
			Cooldowns:     make(map[int]int),
		}
		// Load equipped skills from character data
		m.loadEquippedClasslessSkills()
	}
}

// loadEquippedClasslessSkills fetches the character's equipped classless skills
func (m *model) loadEquippedClasslessSkills() {
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
		Skills []ClasslessSkill `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// Populate equipped slots
	for _, skill := range result.Skills {
		if skill.Slot >= 1 && skill.Slot <= 5 {
			m.combatSkills.EquippedSkill[skill.Slot-1] = skill
		}
	}
}

// ============================================================
// SKILL USAGE
// ============================================================

// useClasslessSkill executes a classless skill in combat
func (m *model) useClasslessSkill(slot int) bool {
	if slot < 1 || slot > 5 {
		m.addCombatLog("Invalid skill slot")
		return false
	}

	skill := m.combatSkills.EquippedSkill[slot-1]
	if skill.ID == 0 {
		m.addCombatLog(fmt.Sprintf("Slot %d is empty", slot))
		return false
	}

	// Check cooldown
	if cd, ok := m.combatSkills.Cooldowns[skill.ID]; ok && cd > 0 {
		m.addCombatLog(fmt.Sprintf("%s is on cooldown (%d rounds)", skill.Name, cd))
		return false
	}

	// Check costs
	if skill.ManaCost > 0 && m.characterMana < skill.ManaCost {
		m.addCombatLog(fmt.Sprintf("Not enough mana for %s", skill.Name))
		return false
	}
	if skill.StaminaCost > 0 && m.characterStamina < skill.StaminaCost {
		m.addCombatLog(fmt.Sprintf("Not enough stamina for %s", skill.Name))
		return false
	}

	// Pay costs
	m.characterMana -= skill.ManaCost
	m.characterStamina -= skill.StaminaCost

	// Execute skill effect
	switch skill.EffectType {
	case "concentrate":
		m.applyConcentrate(skill)
	case "haymaker":
		m.applyHaymaker(skill)
	case "backoff":
		m.applyBackoff(skill)
	case "scream":
		m.applyScream(skill)
	case "slap":
		return m.applySlap(skill) // returns true if stun applied
	}

	// Set cooldown
	m.combatSkills.Cooldowns[skill.ID] = skill.Cooldown

	return true
}

// applyConcentrate: +WIS to accuracy for 4 rounds
func (m *model) applyConcentrate(skill ClasslessSkill) {
	wisMod := m.getWisModifier()
	effect := SkillEffect{
		SkillID:   skill.ID,
		Name:      skill.Name,
		Type:      "accuracy_boost",
		Duration:  skill.Duration,
		Modifier:  wisMod,
		AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
	}
	m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, effect)
	m.addCombatLog(fmt.Sprintf("🎯 Concentrate! +%d accuracy for %d rounds", wisMod, skill.Duration))
}

// applyHaymaker: Next attack +STR damage, -DEX to hit
func (m *model) applyHaymaker(skill ClasslessSkill) {
	strMod := m.getStrModifier()
	dexPenalty := strMod / 2 // Half STR as DEX penalty
	if dexPenalty < 1 {
		dexPenalty = 1
	}

	effect := SkillEffect{
		SkillID:   skill.ID,
		Name:      skill.Name,
		Type:      "haymaker",
		Duration:  skill.Duration,
		Modifier:  strMod, // Positive damage bonus stored here
		AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
	}
	m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, effect)
	m.addCombatLog(fmt.Sprintf("💪 Haymaker! +%d damage, -%d accuracy this attack", strMod+5, dexPenalty))
}

// applyBackoff: Dodge all attacks this round
func (m *model) applyBackoff(skill ClasslessSkill) {
	effect := SkillEffect{
		SkillID:   skill.ID,
		Name:      skill.Name,
		Type:      "dodge_all",
		Duration:  skill.Duration,
		Modifier:  999, // Guaranteed dodge
		AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
	}
	m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, effect)
	m.addCombatLog("💨 Back-off! Dodging all attacks this round!")
}

// applyScream: -WIS/INT, +DEX/STR for 2 rounds
func (m *model) applyScream(skill ClasslessSkill) {
	conMod := getConstitutionModifier(m.getCharacterConstitution())
	statShift := conMod / 2
	if statShift < 1 {
		statShift = 1
	}

	// Apply both buff and debuff
	buffEffect := SkillEffect{
		SkillID:   skill.ID,
		Name:      "Scream Buff",
		Type:      "scream_buff",
		Duration:  skill.Duration,
		Modifier:  statShift,
		AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
	}
	debuffEffect := SkillEffect{
		SkillID:   skill.ID,
		Name:      "Scream Debuff",
		Type:      "scream_debuff",
		Duration:  skill.Duration,
		Modifier:  -statShift,
		AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
	}
	m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, buffEffect, debuffEffect)
	m.addCombatLog(fmt.Sprintf("😤 SCREAM! +%d DEX/STR, -%d WIS/INT for %d rounds", statShift, statShift, skill.Duration))
}

// applySlap: Attempt to stun target for 1 round
func (m *model) applySlap(skill ClasslessSkill) bool {
	dexMod := m.getDexModifier()
	targetCon := 10 // Assume base CON for target
	if m.combatTarget != nil {
		targetCon = m.getTargetConstitution()
	}

	// Roll DEX vs CON
	success := dexMod + rollDie(6) > (targetCon-10)/2 + rollDie(6)

	if success {
		effect := SkillEffect{
			SkillID:   skill.ID,
			Name:      skill.Name,
			Type:      "stun",
			Duration:  skill.Duration,
			Modifier:  0,
			AppliedAt: func() int { if c, ok := m.combatManager.GetCombat(m.combatID); ok { return c.CurrentTick }; return 0 }(),
		}
		m.combatSkills.ActiveEffects = append(m.combatSkills.ActiveEffects, effect)
		m.addCombatLog("👋 SLAP! Target is stunned for 1 round!")
		return true
	}

	m.addCombatLog("👋 Slap missed! Target resisted.")
	return false
}

// ============================================================
// SKILL MANAGEMENT
// ============================================================

// handleClasslessSkillCommand handles skill equip/swap/show commands
func (m *model) handleClasslessSkillCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 1 {
		m.showClasslessSkillsHelp()
		return
	}

	action := parts[0]

	switch action {
	case "show", "list":
		m.showEquippedClasslessSkills()
	case "equip":
		if len(parts) != 3 {
			m.AppendMessage("Usage: skill equip <skill_name> <slot>", "error")
			return
		}
		skillName := parts[1]
		slot := 0
		fmt.Sscanf(parts[2], "%d", &slot)
		m.equipClasslessSkill(skillName, slot)
	case "swap":
		if len(parts) != 3 {
			m.AppendMessage("Usage: skill swap <slot1> <slot2>", "error")
			return
		}
		slot1, slot2 := 0, 0
		fmt.Sscanf(parts[1], "%d", &slot1)
		fmt.Sscanf(parts[2], "%d", &slot2)
		m.swapClasslessSkills(slot1, slot2)
	case "all":
		m.showAllAvailableClasslessSkills()
	default:
		m.showClasslessSkillsHelp()
	}
}

// showEquippedClasslessSkills displays currently equipped skills
func (m *model) showEquippedClasslessSkills() {
	output := "=== Classless Skills (Slots 1-5) ===\n\n"

	for i := 0; i < 5; i++ {
		skill := m.combatSkills.EquippedSkill[i]
		if skill.ID == 0 {
			output += fmt.Sprintf("[%d] (empty)\n", i+1)
		} else {
			output += fmt.Sprintf("[%d] %s - %s\n", i+1, skill.Name, skill.Description)
			if skill.ManaCost > 0 || skill.StaminaCost > 0 {
				output += fmt.Sprintf("     Cost: %d💧 %d⚡ | CD: %d rounds\n",
					skill.ManaCost, skill.StaminaCost, skill.Cooldown)
			}
		}
		output += "\n"
	}

	output += "Use 'skill swap <slot1> <slot2>' to change positions\n"
	output += "Use 'skill equip <name> <slot>' to equip a different skill\n"
	output += "Use 'skill all' to see all available skills"
	m.AppendMessage(output, "info")
}

// showAllAvailableClasslessSkills shows all 5 skills
func (m *model) showAllAvailableClasslessSkills() {
	output := "=== All Classless Skills ===\n\n"

	for _, skill := range ClasslessSkills {
		output += fmt.Sprintf("%s\n", skill.Name)
		output += fmt.Sprintf("  %s\n", skill.Description)
		output += fmt.Sprintf("  Cost: %d💧 %d⚡ | CD: %d | Dur: %d\n\n",
			skill.ManaCost, skill.StaminaCost, skill.Cooldown, skill.Duration)
	}

	output += "All characters can equip any 5 of these skills."
	m.AppendMessage(output, "info")
}

// showClasslessSkillsHelp displays help text
func (m *model) showClasslessSkillsHelp() {
	help := `Classless Skills Commands:
  skill show           - Show equipped skills
  skill all            - Show all available skills
  skill equip <n> <s>  - Equip skill name to slot 1-5
  skill swap <s1> <s2> - Swap skills between slots

In combat, press 1-5 to use the skill in that slot.`
	m.AppendMessage(help, "info")
}

// equipClasslessSkill equips a skill to a slot
func (m *model) equipClasslessSkill(skillName string, slot int) {
	if slot < 1 || slot > 5 {
		m.AppendMessage("Slot must be 1-5", "error")
		return
	}

	// Find skill by name
	var selectedSkill ClasslessSkill
	found := false
	for _, s := range ClasslessSkills {
		if strings.EqualFold(s.Name, skillName) {
			selectedSkill = s
			selectedSkill.Slot = slot
			found = true
			break
		}
	}

	if !found {
		m.AppendMessage(fmt.Sprintf("Skill '%s' not found. Use 'skill all' to see available skills.", skillName), "error")
		return
	}

	// Check if already equipped in another slot
	for i, equipped := range m.combatSkills.EquippedSkill {
		if equipped.ID == selectedSkill.ID && i != slot-1 {
			m.AppendMessage(fmt.Sprintf("%s is already equipped in slot %d", selectedSkill.Name, i+1), "error")
			return
		}
	}

	// Equip to slot
	m.combatSkills.EquippedSkill[slot-1] = selectedSkill

	// Send to server
	payload := fmt.Sprintf(`{"skill_id":%d,"slot":%d}`, selectedSkill.ID, slot)
	resp, err := httpPost(
		fmt.Sprintf("%s/characters/%d/classless-skills", RESTAPIBase, m.currentCharacterID),
		payload)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error saving skill: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	m.AppendMessage(fmt.Sprintf("Equipped %s in slot %d", selectedSkill.Name, slot), "success")
}

// swapClasslessSkills swaps two skill slots
func (m *model) swapClasslessSkills(slot1, slot2 int) {
	if slot1 < 1 || slot1 > 5 || slot2 < 1 || slot2 > 5 {
		m.AppendMessage("Slots must be 1-5", "error")
		return
	}

	// Swap in memory
	m.combatSkills.EquippedSkill[slot1-1], m.combatSkills.EquippedSkill[slot2-1] =
		m.combatSkills.EquippedSkill[slot2-1], m.combatSkills.EquippedSkill[slot1-1]

	// Update slot numbers
	m.combatSkills.EquippedSkill[slot1-1].Slot = slot1
	m.combatSkills.EquippedSkill[slot2-1].Slot = slot2

	// Send to server
	url := fmt.Sprintf("%s/characters/%d/classless-skills/swap", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"slot1":%d,"slot2":%d}`, slot1, slot2)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	m.AppendMessage(fmt.Sprintf("Swapped skills in slots %d and %d", slot1, slot2), "success")
}

// ============================================================
// STAT HELPERS
// ============================================================

// getWisModifier returns the WIS modifier
func (m *model) getWisModifier() int {
	wisdom := m.getCharacterWisdom()
	return (wisdom - 10) / 2
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

// getTargetConstitution estimates target CON (for slap)
func (m *model) getTargetConstitution() int {
	// Could fetch from server, but for now use level-based estimate
	if m.combatTarget == nil {
		return 10
	}
	return 10 + m.combatTarget.Level
}

// rollDie rolls a simple d6
func rollDie(sides int) int {
	return 1 + rand.Intn(sides)
}
