package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderCombatScreen renders the combat UI with layout matching the playing screen
func (m *model) renderCombatScreen() string {
	var s strings.Builder

	width := m.width
	height := m.height
	if width < 40 {
		width = 80
	}
	if height < 10 {
		height = 24
	}

	inputHeight := height * 20 / 100
	if inputHeight < 3 {
		inputHeight = 3
	}
	statusHeight := height * 10 / 100
	if statusHeight < 3 {
		statusHeight = 3
	}
	vpHeight := height - statusHeight - inputHeight
	if vpHeight < 5 {
		vpHeight = 5
	}

	// Build combat content for viewport
	var content strings.Builder

	// Combat header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff6b6b")).
		Background(lipgloss.Color("#2d2d2d")).
		Padding(0, 1).
		Width(width)
	content.WriteString(headerStyle.Render("[ COMBAT ]"))
	content.WriteString("\n\n")

	// Target info
	if m.combatTarget != nil {
		targetName := m.combatTarget.Name
		targetLevel := m.combatTarget.Level
		targetHP := m.combatTarget.HP
		targetMaxHP := m.combatTarget.MaxHP

		targetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9f43"))
		content.WriteString(fmt.Sprintf("You are fighting: %s (Level %d)\n",
			targetStyle.Render(targetName), targetLevel))

		// Target HP bar
		hpBar := m.renderHPBar(targetHP, targetMaxHP, width-20)
		content.WriteString(fmt.Sprintf("HP: %s %d/%d\n", hpBar, targetHP, targetMaxHP))
	}

	content.WriteString("\n")

	// Player stats
	playerHP := m.characterHP
	playerMaxHP := m.characterMaxHP
	playerStamina := m.characterStamina
	playerMaxStamina := m.characterMaxStamina

	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#74b9ff"))
	content.WriteString(fmt.Sprintf("Your HP: %s\n", statsStyle.Render(fmt.Sprintf("%d/%d", playerHP, playerMaxHP))))
	content.WriteString(fmt.Sprintf("Stamina: %s\n\n", statsStyle.Render(fmt.Sprintf("%d/%d", playerStamina, playerMaxStamina))))

	// Combat actions
	actionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a29bfe"))
	actionBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#636e72")).
		Padding(0, 1).
		Width(width - 4)

	// Build action display with equipped talents
	slot1 := m.getTalentSlotName(1)
	slot2 := m.getTalentSlotName(2)
	slot3 := m.getTalentSlotName(3)
	slot4 := m.getTalentSlotName(4)
	potionSlot := m.getPotionSlotName()

	actions := fmt.Sprintf("[1] %-10s  [2] %-10s\n[3] %-10s  [4] %-10s\n[R] %-12s  [Q] Flee",
		slot1, slot2, slot3, slot4, potionSlot)
	content.WriteString(actionBox.Render(actionStyle.Render(actions)))
	content.WriteString("\n")

	// Queued action indicator
	if m.combatQueuedAction != "" {
		queueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00b894")).Bold(true)
		content.WriteString(queueStyle.Render(fmt.Sprintf("  > Queued: %s", m.combatQueuedAction)))
		content.WriteString("\n")
	}

	// Tick indicator
	tickStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fdcb6e"))
	content.WriteString(tickStyle.Render("  > Tick combat active (1.5s intervals)"))
	content.WriteString("\n\n")

	// Combat log
	if len(m.combatLog) > 0 {
		logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#dfe6e9"))
		content.WriteString(logStyle.Render("Combat Log:"))
		content.WriteString("\n")
		for i := len(m.combatLog) - 1; i >= 0 && i >= len(m.combatLog)-8; i-- {
			content.WriteString(fmt.Sprintf("  %s\n", m.combatLog[i]))
		}
	}

	// Render viewport with combat content
	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Width(width)

	if m.viewport.Width != width {
		m.viewport.Width = width
	}
	if m.viewport.Height != vpHeight {
		m.viewport.Height = vpHeight
	}
	m.viewport.SetContent(content.String())
	s.WriteString(viewportStyle.Render(m.viewport.View()))
	s.WriteString("\n")

	// Status bar (use red heart always in combat)
	statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana, false)
	debugInfo := ""
	if m.debugMode {
		debugInfo = " " + lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("[Combat: %d]", m.combatID))
	}
	statusBarStyle := lipgloss.NewStyle().
		Foreground(pink).
		Background(lipgloss.Color("235")).
		Bold(true).
		Width(width).
		Padding(0, 1)
	s.WriteString(statusBarStyle.Render(statsLine + debugInfo))
	s.WriteString("\n")

	// Input area
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(inputHeight - 2)
	s.WriteString(inputStyle.Render(promptStyle.Render("Action: ") + m.textInput.View()))

	return s.String()
}

// renderHPBar renders a visual HP bar
func (m *model) renderHPBar(current, max, width int) string {
	if max == 0 {
		return ""
	}

	percentage := float64(current) / float64(max)
	filled := int(percentage * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	empty := width - filled

	var bar strings.Builder
	bar.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00b894")).Render(strings.Repeat("|", filled)))
	bar.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#636e72")).Render(strings.Repeat("-", empty)))

	return bar.String()
}
