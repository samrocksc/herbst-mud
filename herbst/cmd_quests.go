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
