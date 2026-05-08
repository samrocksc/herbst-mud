package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleAbilityCommand handles skill equip/swap/show commands
func (m *model) handleAbilityCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 1 {
		m.showAbilityHelp()
		return
	}

	action := parts[0]
	switch action {
	case "show", "list":
		m.showEquippedAbilities()
	case "slot":
		if len(parts) != 2 {
			m.AppendMessage("Usage: skill slot <1-5>", "error")
			return
		}
		slot := 0
		fmt.Sscanf(parts[1], "%d", &slot)
		if slot < 1 || slot > 5 {
			m.AppendMessage("Slot must be between 1 and 5", "error")
			return
		}
		m.startAbilitySelection(slot)
	case "equip":
		if len(parts) != 3 {
			m.AppendMessage("Usage: skill equip <skill_name> <slot>", "error")
			return
		}
		skillName := parts[1]
		slot := 0
		fmt.Sscanf(parts[2], "%d", &slot)
		m.equipAbility(skillName, slot)
	case "swap":
		if len(parts) != 3 {
			m.AppendMessage("Usage: skill swap <slot1> <slot2>", "error")
			return
		}
		slot1, slot2 := 0, 0
		fmt.Sscanf(parts[1], "%d", &slot1)
		fmt.Sscanf(parts[2], "%d", &slot2)
		m.swapAbilities(slot1, slot2)
	case "all":
		m.showAllAvailableAbilities()
	default:
		m.showAbilityHelp()
	}
}

// showEquippedAbilities displays currently equipped abilities
func (m *model) showEquippedAbilities() {
	output := "═══════════════════════════════════════════\n"
	output += "           Combat Abilities\n"
	output += "═══════════════════════════════════════════\n\n"

	if m.combatSkills == nil {
		m.initCombatSkillState()
	}

	output += "[ Classless ] — Available to all characters\n\n"
	for i := 0; i < 5; i++ {
		ability := m.combatSkills.EquippedSkill[i]
		if ability.ID == 0 {
			output += fmt.Sprintf("  [%d] - (empty)\n", i+1)
			continue
		}
		output += fmt.Sprintf("  [%d] + %s\n", i+1, ability.Name)
		output += fmt.Sprintf("       | %s\n", ability.Description)
		if ability.ManaCost > 0 || ability.StaminaCost > 0 {
			output += fmt.Sprintf("       + Cost: %d MP %d SP | CD: %d rounds\n",
				ability.ManaCost, ability.StaminaCost, ability.Cooldown)
		} else {
			output += fmt.Sprintf("       + CD: %d rounds\n", ability.Cooldown)
		}
		output += "\n"
	}

	output += "-------------------------------------------\n"
	output += "In combat: Press 1-5 to activate\n"
	output += "To change: skill slot <1-5>"
	m.AppendMessage(output, "info")
}

// showAllAvailableAbilities fetches and displays all classless abilities from server
func (m *model) showAllAvailableAbilities() {
	resp, err := httpGet(fmt.Sprintf("%s/abilities", RESTAPIBase))
	if err != nil {
		m.AppendMessage("Error fetching abilities", "error")
		return
	}
	defer resp.Body.Close()

	var result struct {
		Abilities []AbilityData `json:"abilities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage("Error parsing abilities", "error")
		return
	}

	output := "═══════════════════════════════════════════\n"
	output += "       All Classless Combat Abilities\n"
	output += "═══════════════════════════════════════════\n\n"

	for _, ability := range result.Abilities {
		if ability.Slot < 0 {
			continue
		}
		output += fmt.Sprintf("+- %s\n", ability.Name)
		output += fmt.Sprintf("| %s\n", ability.Description)
		costParts := []string{}
		if ability.ManaCost > 0 {
			costParts = append(costParts, fmt.Sprintf("%d MP", ability.ManaCost))
		}
		if ability.StaminaCost > 0 {
			costParts = append(costParts, fmt.Sprintf("%d SP", ability.StaminaCost))
		}
		costStr := strings.Join(costParts, " ")
		if costStr == "" {
			costStr = "Free"
		}
		output += fmt.Sprintf("+ Cost: %s | CD: %d rounds\n\n", costStr, ability.Cooldown)
	}

	output += "-------------------------------------------\n"
	output += "Any 5 of these can be equipped to your combat slots."
	m.AppendMessage(output, "info")
}

// startAbilitySelection enters skill selection mode for a slot
func (m *model) startAbilitySelection(slot int) {
	if m.combatSkills == nil {
		m.initCombatSkillState()
	}
	m.skillSelectSlot = slot
	m.skillSelectCursor = 0
	m.screen = ScreenSkillSelect
	m.renderAbilitySelection()
}

// renderAbilitySelection displays the ability selection UI
func (m *model) renderAbilitySelection() {
	output := "═══════════════════════════════════════════\n"
	output += fmt.Sprintf("      Choose Ability for Slot %d\n", m.skillSelectSlot)
	output += "═══════════════════════════════════════════\n\n"

	currentAbility := m.combatSkills.EquippedSkill[m.skillSelectSlot-1]
	if currentAbility.ID != 0 {
		output += fmt.Sprintf("Currently Equipped: %s\n\n", currentAbility.Name)
	} else {
		output += "Slot Empty - Choose an ability:\n\n"
	}

	// Fetch abilities from server for display
	resp, err := httpGet(fmt.Sprintf("%s/abilities", RESTAPIBase))
	var availableAbilities []AbilityData
	if err == nil {
		defer resp.Body.Close()
		var result struct {
			Abilities []AbilityData `json:"abilities"`
		}
		if json.NewDecoder(resp.Body).Decode(&result) == nil {
			availableAbilities = result.Abilities
		}
	}

	output += "+- Classless Abilities (Available to All) -+\n"
	for i, ability := range availableAbilities {
		cursor := "  "
		if i == m.skillSelectCursor {
			cursor = "> "
		}
		costStr := ""
		if ability.ManaCost > 0 {
			costStr += fmt.Sprintf(" %d MP", ability.ManaCost)
		}
		if ability.StaminaCost > 0 {
			costStr += fmt.Sprintf(" %d SP", ability.StaminaCost)
		}
		if costStr == "" {
			costStr = " Free"
		}
		output += fmt.Sprintf("%s%d. %-15s |%s | CD:%d\n",
			cursor, i+1, ability.Name, costStr, ability.Cooldown)
	}
	output += "+------------------------------------------+\n\n"

	if m.skillSelectCursor >= 0 && m.skillSelectCursor < len(availableAbilities) {
		selected := availableAbilities[m.skillSelectCursor]
		output += fmt.Sprintf("> %s\n", selected.Name)
		output += fmt.Sprintf("  %s\n\n", selected.Description)
	}

	output += "Commands: 1-5 select | enter confirm | q cancel"
	m.AppendMessage(output, "info")
}

// handleAbilitySelectionInput processes input in ability selection mode
func (m *model) handleAbilitySelectionInput(key string) bool {
	switch strings.ToLower(key) {
	case "up", "k":
		if m.skillSelectCursor > 0 {
			m.skillSelectCursor--
		}
		m.renderAbilitySelection()
		return true
	case "down", "j":
		m.skillSelectCursor++
		m.renderAbilitySelection()
		return true
	case "enter", "b", " ":
		// Fetch abilities to get the one at cursor
		resp, err := httpGet(fmt.Sprintf("%s/abilities", RESTAPIBase))
		if err == nil {
			defer resp.Body.Close()
			var result struct {
				Abilities []AbilityData `json:"abilities"`
			}
			if json.NewDecoder(resp.Body).Decode(&result) == nil && m.skillSelectCursor < len(result.Abilities) {
				selected := result.Abilities[m.skillSelectCursor]
				m.equipAbility(selected.Name, m.skillSelectSlot)
				m.screen = ScreenPlaying
				m.AppendMessage(fmt.Sprintf("Equipped %s to slot %d!", selected.Name, m.skillSelectSlot), "success")
				return false
			}
		}
		m.AppendMessage("Error selecting ability", "error")
		return false
	case "q", "esc", "cancel":
		m.screen = ScreenPlaying
		m.AppendMessage("Ability selection cancelled.", "info")
		return false
	default:
		m.AppendMessage("Use up/down to select, enter to confirm, q to cancel", "info")
		return true
	}
}

// showAbilityHelp displays help text
func (m *model) showAbilityHelp() {
	help := `Ability Commands:
  skills               - Show equipped combat abilities
  skill slot <1-5>    - Select an ability for a slot
  skill all            - Show available classless abilities
  skill equip <n> <s>  - Equip ability to slot 1-5 (quick)
  skill swap <s1> <s2> - Swap abilities between slots

In combat: press 1-5 to activate the ability in that slot.`
	m.AppendMessage(help, "info")
}

// equipAbility equips an ability to a slot via the server
func (m *model) equipAbility(abilityName string, slot int) {
	if slot < 1 || slot > 5 {
		m.AppendMessage("Slot must be 1-5", "error")
		return
	}

	// Find ability ID by name from server
	resp, err := httpGet(fmt.Sprintf("%s/abilities", RESTAPIBase))
	if err != nil {
		m.AppendMessage("Error fetching abilities", "error")
		return
	}
	defer resp.Body.Close()

	var result struct {
		Abilities []AbilityData `json:"abilities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage("Error parsing abilities", "error")
		return
	}

	var foundAbility *AbilityData
	for i := range result.Abilities {
		if strings.EqualFold(result.Abilities[i].Name, abilityName) {
			foundAbility = &result.Abilities[i]
			break
		}
	}

	if foundAbility == nil {
		m.AppendMessage(fmt.Sprintf("Ability '%s' not found. Use 'skill all' to see available abilities.", abilityName), "error")
		return
	}

	// Check if already equipped in another slot
	for i, equipped := range m.combatSkills.EquippedSkill {
		if equipped.ID == foundAbility.ID && i != slot-1 {
			m.AppendMessage(fmt.Sprintf("%s is already equipped in slot %d", foundAbility.Name, i+1), "error")
			return
		}
	}

	// Send to server
	payload := fmt.Sprintf(`{"skill_id":%d,"slot":%d}`, foundAbility.ID, slot)
	postResp, err := httpPost(
		fmt.Sprintf("%s/characters/%d/classless-skills", RESTAPIBase, m.currentCharacterID),
		payload)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error saving ability: %v", err), "error")
		return
	}
	postResp.Body.Close()

	// Update local state
	foundAbility.Slot = slot
	m.combatSkills.EquippedSkill[slot-1] = *foundAbility
	m.AppendMessage(fmt.Sprintf("Equipped %s in slot %d", foundAbility.Name, slot), "success")
}

// swapAbilities swaps two ability slots
func (m *model) swapAbilities(slot1, slot2 int) {
	if slot1 < 1 || slot1 > 5 || slot2 < 1 || slot2 > 5 {
		m.AppendMessage("Slots must be 1-5", "error")
		return
	}

	m.combatSkills.EquippedSkill[slot1-1], m.combatSkills.EquippedSkill[slot2-1] =
		m.combatSkills.EquippedSkill[slot2-1], m.combatSkills.EquippedSkill[slot1-1]

	m.combatSkills.EquippedSkill[slot1-1].Slot = slot1
	m.combatSkills.EquippedSkill[slot2-1].Slot = slot2

	url := fmt.Sprintf("%s/characters/%d/classless-skills/swap", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"slot1":%d,"slot2":%d}`, slot1, slot2)
	req, err := http.NewRequest("PUT", url, strings.NewReader(payload))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error creating swap request: %v", err), "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error swapping abilities: %v", err), "error")
		return
	}
	resp.Body.Close()

	m.AppendMessage(fmt.Sprintf("Swapped abilities in slots %d and %d", slot1, slot2), "success")
}

// getAbilitySlotName returns a display name for a combat ability slot
func (m *model) getAbilitySlotName(slot int) string {
	if slot < 1 || slot > 5 {
		return "[Attack]"
	}
	if m.combatSkills != nil {
		ability := m.combatSkills.EquippedSkill[slot-1]
		if ability.ID != 0 {
			name := ability.Name
			if len(name) > 12 {
				name = name[:12]
			}
			return fmt.Sprintf("[%s]", name)
		}
	}
	return "[Attack]"
}