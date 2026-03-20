package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ============================================================
// QUESTS
// ============================================================

func (m *model) handleQuestsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.displayQuestTrackerPlaceholder()
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/characters/%d/quests", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.displayQuestTrackerPlaceholder()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.displayQuestTrackerPlaceholder()
		return
	}

	var questResp struct {
		Quests []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
			Objectives  []struct {
				Description string `json:"description"`
				Current     int    `json:"current"`
				Total       int    `json:"total"`
			} `json:"objectives"`
			Giver   string `json:"giver"`
			Rewards string `json:"rewards"`
		} `json:"quests"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&questResp); err != nil || len(questResp.Quests) == 0 {
		m.displayQuestTrackerPlaceholder()
		return
	}

	var quests strings.Builder

	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	activeCount := 0
	availableCount := 0
	completedCount := 0

	for _, quest := range questResp.Quests {
		switch quest.Status {
		case "in_progress":
			activeCount++
		case "available":
			availableCount++
		case "completed":
			completedCount++
		}

		quests.WriteString(questBoxStyle.Render("") + "\n")

		statusColor := questAvailableStyle
		statusText := "Available"
		if quest.Status == "in_progress" {
			statusColor = questProgressStyle
			statusText = "In Progress"
		} else if quest.Status == "completed" {
			statusColor = questCompletedStyle
			statusText = "Completed"
		}

		quests.WriteString(fmt.Sprintf("  %s [%s]\n", questTitleStyle.Render(quest.Name), statusColor.Render(statusText)))

		if quest.Description != "" {
			quests.WriteString(fmt.Sprintf("    %s\n", quest.Description))
		}

		if len(quest.Objectives) > 0 {
			quests.WriteString("\n  Objectives:\n")
			for _, obj := range quest.Objectives {
				progress := fmt.Sprintf("%d/%d", obj.Current, obj.Total)
				if obj.Current >= obj.Total {
					quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", obj.Description, questCompletedStyle.Render("("+progress+")")))
				} else {
					quests.WriteString(fmt.Sprintf("    ○ %s %s\n", obj.Description, questProgressStyle.Render("("+progress+")")))
				}
			}
		}

		if quest.Giver != "" {
			quests.WriteString(fmt.Sprintf("\n  Giver: %s\n", quest.Giver))
		}
		if quest.Rewards != "" {
			quests.WriteString(fmt.Sprintf("  Reward: %s\n", quest.Rewards))
		}

		quests.WriteString("\n")
	}

	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString(fmt.Sprintf("  Active: %d  |  Available: %d  |  Completed: %d\n",
		activeCount, availableCount, completedCount))
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	m.AppendMessage(quests.String(), "info")
}

func (m *model) displayQuestTrackerPlaceholder() {
	var quests strings.Builder

	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Prove Yourself"),
		questProgressStyle.Render("In Progress")))

	quests.WriteString("    The Scrapyard ain't for the weak. Kill 3 Scrap Rats\n")
	quests.WriteString("    and I'll let you into New Venice proper.\n\n")
	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Kill Scrap Rat", questProgressStyle.Render("(2/3)")))
	quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", "Find Guard Marco at Foggy Gate", questCompletedStyle.Render("(done)")))

	quests.WriteString("\n  Giver: Guard Marco\n")
	quests.WriteString("  Reward: 10 coins\n\n")

	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Ooze Samples"),
		questAvailableStyle.Render("Available")))

	quests.WriteString("    Jane needs Ooze samples for her research.\n")
	quests.WriteString("    The Leaking Pipes have plenty.\n\n")
	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Collect glowing goo", questProgressStyle.Render("(0/5)")))

	quests.WriteString("\n  Giver: Scavenger Jane\n")
	quests.WriteString("  Reward: repair_kit, scavenge skill\n\n")

	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString("  Active: 1  |  Available: 1  |  Completed: 0\n")
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	quests.WriteString("\n" + infoStyle.Render("  Use 'quest <name>' for details, 'accept <quest>' to begin."))

	m.AppendMessage(quests.String(), "info")
}

// ============================================================
// SKILLS & TALENTS
// ============================================================

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
		output += "Slots: 1-4 (quick access keys)\n"
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

func (m *model) handleSkillEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}
	m.AppendMessage("Skills are always active and cannot be unequipped.\nThey provide passive bonuses based on your skill level.", "info")
}

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

// ============================================================
// DEBUG
// ============================================================

func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		if m.debugMode {
			m.AppendMessage("Debug mode: ON (Room ID visible in status bar)", "info")
		} else {
			m.AppendMessage("Debug mode: OFF\nUsage: debug on | debug off", "info")
		}
		return
	}

	switch parts[1] {
	case "on", "true", "1", "yes":
		m.debugMode = true
		m.AppendMessage("Debug mode: ON (Room ID will show in status bar)", "success")
	case "off", "false", "0", "no":
		m.debugMode = false
		m.AppendMessage("Debug mode: OFF", "info")
	default:
		m.AppendMessage("Usage: debug on | debug off", "error")
	}
}
