package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// handleSkillEquipCommand handles skill equip commands
func (m *model) handleSkillEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}
	m.AppendMessage("Skills are always active and cannot be unequipped.\nThey provide passive bonuses based on your skill level.", "info")
}

