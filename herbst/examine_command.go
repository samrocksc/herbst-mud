package main

import (
	"context"
	"fmt"
	"strings"

	"herbst/quest"
)

// CharacterExamineSkill tracks examine skill per character
// In production, this would be persisted to the database
var characterExamineSkills = make(map[int]*ExamineSkillLevel)

// GetExamineSkill returns the examine skill for a character
func GetExamineSkill(characterID int) *ExamineSkillLevel {
	if skill, ok := characterExamineSkills[characterID]; ok {
		return skill
	}
	// Default new characters start at level 1 with 0 XP
	skill := &ExamineSkillLevel{Level: 1, XP: 0}
	characterExamineSkills[characterID] = skill
	return skill
}

// handleExamineCommand handles the examine/ex/inspect command
func (m *model) handleExamineCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Usage: examine <item|npc|object>\n\nExamine something in detail to reveal hidden details."
		m.messageType = "error"
		return
	}

	target := strings.Join(parts[1:], " ")

	// First check if it's a direction (peer into adjacent room)
	validDirs := map[string]string{"north": "north", "south": "south", "east": "east", "west": "west", "up": "up", "down": "down"}
	if dir, ok := validDirs[strings.ToLower(target)]; ok {
		m.handlePeerCommand("peer " + dir)
		return
	}

	// Check if the target is an item in the current room
	if m.client != nil {
		// Get current room with items
		char, err := m.client.Character.Get(context.Background(), m.currentCharacterID)
		if err != nil {
			m.message = fmt.Sprintf("Error: %v", err)
			m.messageType = "error"
			return
		}

		// Search in room items
		roomItems, err := m.client.Room.QueryEquipment(m.client.Room.GetX(context.Background(), char.CurrentRoomID)).All(context.Background())
		if err == nil && roomItems != nil {
			for _, item := range roomItems {
				if strings.Contains(strings.ToLower(item.Name), strings.ToLower(target)) ||
					strings.Contains(strings.ToLower(item.Description), strings.ToLower(target)) {
					m.displayItemExamineWithQuests(item.Name, item.Description, item.ExamineDesc, item.ID)
					return
				}
			}
		}

		// Search NPCs in the room
		characters, err := m.client.Room.QueryCharacters(m.client.Room.GetX(context.Background(), char.CurrentRoomID)).All(context.Background())
		if err == nil && characters != nil {
			for _, npc := range characters {
				if npc.IsNPC && (strings.Contains(strings.ToLower(npc.Name), strings.ToLower(target)) ||
					strings.Contains(strings.ToLower(npc.Description), strings.ToLower(target))) {
					m.displayNPCExamine(npc.Name, npc.Description)
					return
				}
			}
		}
	}

	m.message = fmt.Sprintf("You don't see '%s' here.", target)
	m.messageType = "error"
}

// displayItemExamineWithQuests displays item information and checks for quest unlocks
func (m *model) displayItemExamineWithQuests(name, description, examineDesc string, itemID int) {
	var output strings.Builder

	output.WriteString(lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(name))
	output.WriteString("\n\n")

	if description != "" {
		output.WriteString(description)
		output.WriteString("\n\n")
	}

	if examineDesc != "" {
		output.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("Examination reveals: "))
		output.WriteString(examineDesc)
		output.WriteString("\n")
	}

	// Get character's examine skill
	examineSkill := GetExamineSkill(m.currentCharacterID)

	// Check for quest unlocks
	questResult := quest.CheckExamineQuestUnlock(m.currentCharacterID, name, examineSkill.Level)

	if questResult != nil {
		if questResult.Unlocked {
			// Quest was just unlocked!
			output.WriteString("\n")
			output.WriteString(lipgloss.NewStyle().Bold(true).Foreground(purple).Render("✦ Quest Unlocked! ✦"))
			output.WriteString("\n")
			output.WriteString(lipgloss.NewStyle().Foreground(purple).Render(questResult.QuestName))
			output.WriteString("\n\n")
			output.WriteString(questResult.RevealText)
			output.WriteString("\n")

			// Grant examine XP
			examineSkill.AddXP(questResult.XPGained)
			output.WriteString(fmt.Sprintf("\n[Examine skill +%d XP]", questResult.XPGained))

			m.message = output.String()
			m.messageType = "success"
			return
		} else if questResult.AlreadyUnlocked {
			// Quest already unlocked - show hint
			output.WriteString("\n")
			output.WriteString(lipgloss.NewStyle().Foreground(gray).Render("You sense there's more to discover here..."))
		}
	}

	// Show examine skill level
	output.WriteString(fmt.Sprintf("\n[Examine skill: %d]", examineSkill.Level))

	m.message = output.String()
	m.messageType = "info"
}

// displayNPCExamine displays NPC examination (no quest integration for NPCs yet)
func (m *model) displayNPCExamine(name, description string) {
	var output strings.Builder

	output.WriteString(lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(name))
	output.WriteString("\n\n")

	if description != "" {
		output.WriteString(description)
		output.WriteString("\n")
	}

	m.message = output.String()
	m.messageType = "info"
}

// handleLookAtCommand handles "look at <target>" - same as examine
func (m *model) handleLookAtCommand(target string) {
	if target == "" {
		m.message = "Look at what?"
		m.messageType = "error"
		return
	}

	// Handle "look at me" - show character info
	if target == "me" || target == "myself" {
		m.message = fmt.Sprintf("=== Your Character ===\nName: %s\nDescription: %s",
			m.currentUserName, m.characterDescription)
		if m.characterGender != "" {
			m.message += fmt.Sprintf("\nGender: %s", m.characterGender)
		}
		// Show examine skill
		examineSkill := GetExamineSkill(m.currentCharacterID)
		m.message += fmt.Sprintf("\nExamine Skill: Level %d (%d/10 XP)", examineSkill.Level, examineSkill.XP)
		m.messageType = "info"
		return
	}

	// Otherwise, treat as examine
	m.handleExamineCommand("examine " + target)
}

// handleQuestsCommand shows the player's quest log
func (m *model) handleQuestsCommand() {
	quests := quest.GetVisibleQuests(m.currentCharacterID)

	if len(quests) == 0 {
		m.message = "You haven't discovered any quests yet."
		m.messageType = "info"
		return
	}

	var output strings.Builder
	output.WriteString("=== Your Quests ===\n\n")

	for _, q := range quests {
		status := ""
		if cq := quest.GlobalQuestStore.GetCharacterQuest(m.currentCharacterID, q.ID); cq != nil {
			status = fmt.Sprintf(" [%s]", cq.Status)
		} else if q.Hidden {
			continue // Don't show hidden quests that aren't unlocked
		}

		output.WriteString(fmt.Sprintf("%s%s\n", q.Name, status))
		output.WriteString(fmt.Sprintf("  %s\n\n", q.Description))
	}

	m.message = output.String()
	m.messageType = "info"
}