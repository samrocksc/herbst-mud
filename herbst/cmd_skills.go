package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleSkillsCommand displays the player's skills
func (m *model) handleSkillsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/skills", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching skills: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load skills", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing skills: %v", err), "error")
		return
	}

	skills, ok := result["skills"].(map[string]interface{})
	if !ok {
		m.AppendMessage("Error: skills data not found", "error")
		return
	}

	output := "=== Your Skills ===\n\n"
	for skillName, skillData := range skills {
		data := skillData.(map[string]interface{})
		level := int(data["level"].(float64))
		bonus := data["bonus"].(string)
		output += fmt.Sprintf("%-15s Lv: %2d  %s\n", skillName+":", level, bonus)
	}

	output += "\nSkills are always active and provide passive bonuses."
	m.AppendMessage(output, "info")
}

// handleTalentsCommand displays the player's talents
func (m *model) handleTalentsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching talents: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load talents", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing talents: %v", err), "error")
		return
	}

	output := "=== Your Talents ===\n\n"
	slots, ok := result["slots"].([]interface{})
	if !ok {
		output += "No talents equipped.\n\n"
		output += "Use: talent equip <talent_id> <slot>\n"
		output += "Slots: 1-4 (quick access keys)"
		m.AppendMessage(output, "info")
		return
	}

	emptySlots := 0
	for i := 1; i <= 4; i++ {
		if i < len(slots) && slots[i] != nil {
			slot := slots[i].(map[string]interface{})
			name := slot["name"].(string)
			desc := slot["description"].(string)
			output += fmt.Sprintf("[%d] %s\n     %s\n\n", i, name, desc)
		} else {
			output += fmt.Sprintf("[%d] (empty)\n\n", i)
			emptySlots++
		}
	}

	if emptySlots == 4 {
		output += "No talents equipped. Use 'talent equip <id> <slot>' to equip."
	}

	m.AppendMessage(output, "info")
}

// handleSkillEquipCommand handles skill equip commands
func (m *model) handleSkillEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}
	m.AppendMessage("Skills are always active and cannot be unequipped.\nThey provide passive bonuses based on your skill level.", "info")
}

// handleTalentEquipCommand handles talent equip/unequip/swap commands
func (m *model) handleTalentEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage:\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
		return
	}

	action := parts[1]

	switch action {
	case "equip":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent equip <talent_id> <slot>\nExample: talent equip 1 2", "error")
			return
		}
		talentID := parts[2]
		slot := parts[3]

		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"talent_id":%s,"slot":%s}`, talentID, slot)
		resp, err := httpPost(url, reqBody)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error equipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			m.AppendMessage("Failed to equip talent", "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("Talent equipped in slot %s", slot), "success")

	case "unequip":
		if len(parts) != 3 {
			m.AppendMessage("Usage: talent unequip <slot>\nExample: talent unequip 2", "error")
			return
		}
		slot := parts[2]

		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		url := fmt.Sprintf("%s/characters/%d/talents/%s", RESTAPIBase, m.currentCharacterID, slot)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talent unequipped from slot %s", slot), "success")

	case "swap":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent swap <slot1> <slot2>\nExample: talent swap 1 2", "error")
			return
		}
		slot1 := parts[2]
		slot2 := parts[3]

		slot1Num, slot2Num := 0, 0
		fmt.Sscanf(slot1, "%d", &slot1Num)
		fmt.Sscanf(slot2, "%d", &slot2Num)
		if slot1Num < 1 || slot1Num > 4 || slot2Num < 1 || slot2Num > 4 {
			m.AppendMessage("Slots must be between 1 and 4", "error")
			return
		}

		url := fmt.Sprintf("%s/characters/%d/talents/swap", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"slot1":%s,"slot2":%s}`, slot1, slot2)
		req, err := http.NewRequest("PUT", url, strings.NewReader(reqBody))
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talents swapped between slot %s and %s", slot1, slot2), "success")

	default:
		m.AppendMessage("Usage:\n  talent - Show talents\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
	}
}