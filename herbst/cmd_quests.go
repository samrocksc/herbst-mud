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

	resp, err := httpGet(fmt.Sprintf("%s/api/characters/%d/quests", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage("Could not reach quest server.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.displayQuestTrackerPlaceholder()
		return
	}

	var apiResp struct {
		QuestProgress []questProgressAPI `json:"quest_progress"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		m.AppendMessage("Error reading quest data.", "error")
		return
	}
	if len(apiResp.QuestProgress) == 0 {
		m.AppendMessage("You have no quests yet. Explore the world to find them!", "info")
		return
	}

	var b strings.Builder
	b.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	b.WriteString(questTitleStyle.Render("  QUEST LOG") + "\n")
	b.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	active := 0
	completed := 0
	abandoned := 0

	for _, qp := range apiResp.QuestProgress {
		statusColor := questAvailableStyle
		statusText := strings.ToUpper(qp.Status)
		switch qp.Status {
		case "active":
			active++
			statusColor = questProgressStyle
			statusText = "Active"
		case "completed":
			completed++
			statusColor = questCompletedStyle
			statusText = "Completed"
		case "abandoned":
			abandoned++
			statusColor = questAvailableStyle
		}

		b.WriteString(fmt.Sprintf("  %s [%s]\n",
			questTitleStyle.Render(qp.QuestName), statusColor.Render(statusText)))

		// Show objectives from quest cache if available
		qDef, found := m.questService.GetQuest(qp.QuestID)
		if found && len(qDef.Objectives) > 0 {
			b.WriteString("\n  Objectives:\n")
			for i, obj := range qDef.Objectives {
				key := obj.Type + ":" + obj.TargetID
				current := qp.ObjectiveCounts[key]
				if current >= obj.Count {
					b.WriteString(fmt.Sprintf("    ✓ %s %s\n",
						obj.Label, questCompletedStyle.Render(fmt.Sprintf("(%d/%d)", current, obj.Count))))
				} else if current > 0 {
					b.WriteString(fmt.Sprintf("    ○ %s %s\n",
						obj.Label, questProgressStyle.Render(fmt.Sprintf("(%d/%d)", current, obj.Count))))
				} else {
					b.WriteString(fmt.Sprintf("    ○ %s (%d)\n", obj.Label, obj.Count))
				}
				_ = i
			}
		}

		// Show rewards from quest cache
		if found && qDef.Rewards.XP > 0 {
			b.WriteString(fmt.Sprintf("\n  Reward: %d XP\n", qDef.Rewards.XP))
		}

		b.WriteString("\n")
	}

	b.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	b.WriteString(fmt.Sprintf("  Active: %d  |  Completed: %d  |  Abandoned: %d\n",
		active, completed, abandoned))
	b.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	m.AppendMessage(b.String(), "info")
}

type questProgressAPI struct {
	ID              int            `json:"id"`
	CharacterID     int            `json:"character_id"`
	QuestID         int            `json:"quest_id"`
	QuestName       string         `json:"quest_name"`
	Status          string         `json:"status"`
	CurrentStep     int            `json:"current_step"`
	ObjectiveCounts map[string]int `json:"objective_counts"`
	StartedAt       string         `json:"started_at"`
	CompletedAt     *string        `json:"completed_at,omitempty"`
}

func (m *model) displayQuestTrackerPlaceholder() {
	var b strings.Builder
	b.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	b.WriteString(questTitleStyle.Render("  QUEST LOG") + "\n")
	b.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")
	b.WriteString("  No quest data available.\n")
	b.WriteString("  Explore the world to discover quests!\n\n")
	b.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	m.AppendMessage(b.String(), "info")
}