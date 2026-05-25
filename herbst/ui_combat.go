package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// getSkillCooldown returns the remaining cooldown for a given skill slot,
// or 0 if not on cooldown.
func (m *model) getSkillCooldown(slot int) int {
	if m.combatSkills == nil || slot < 1 || slot > 5 {
		return 0
	}
	skill := m.combatSkills.EquippedSkill[slot-1]
	if skill.ID == 0 {
		return 0
	}
	if cd, ok := m.combatSkills.Cooldowns[skill.ID]; ok {
		return cd
	}
	return 0
}

// formatSkillSlot builds a display string for one skill slot in the HUD.
// Shows skill name, cooldown indicator, and resource cost.
func (m *model) formatSkillSlot(slot int) string {
	if m.combatSkills == nil || slot < 1 || slot > 5 {
		return "[Attack]"
	}
	skill := m.combatSkills.EquippedSkill[slot-1]
	if skill.ID == 0 {
		return "[Empty]"
	}

	name := skill.Name
	if len(name) > 10 {
		name = name[:10]
	}

	cd := m.getSkillCooldown(slot)
	if cd > 0 {
		return fmt.Sprintf("%s (%d)", name, cd)
	}
	return name
}

// renderCombatScreen renders the combat UI with a bottom HUD for
// 4 skill slots + 1 potion slot, and a scrollable combat log above.
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

	inputHeight := height * 15 / 100
	if inputHeight < 3 {
		inputHeight = 3
	}
	hudHeight := 7
	vpHeight := height - hudHeight - inputHeight
	if vpHeight < 5 {
		vpHeight = 5
	}

	// Build combat log content for viewport
	var content strings.Builder

	// Combat header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff6b6b")).
		Background(lipgloss.Color("#2d2d2d")).
		Padding(0, 1).
		Width(width)
	content.WriteString(headerStyle.Render(" COMBAT "))
	content.WriteString("\n\n")

	// Target info
	if m.combatTarget != nil {
		targetName := m.combatTarget.Name
		targetLevel := m.combatTarget.Level
		targetHP := m.combatTarget.HP
		targetMaxHP := m.combatTarget.MaxHP

		targetStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9f43"))
		content.WriteString(fmt.Sprintf("%s (Level %d)\n",
			targetStyle.Render(targetName), targetLevel))

		// Target HP bar
		hpBar := m.renderHPBar(targetHP, targetMaxHP, width-24)
		content.WriteString(fmt.Sprintf("HP: %s %d/%d\n", hpBar, targetHP, targetMaxHP))
	}

	content.WriteString("\n")

	// Player stats
	playerHP := m.characterHP
	playerMaxHP := m.characterMaxHP
	playerStamina := m.characterStamina
	playerMaxStamina := m.characterMaxStamina
	playerMana := m.characterMana
	playerMaxMana := m.characterMaxMana

	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#74b9ff"))
	hpBar := m.renderHPBar(playerHP, playerMaxHP, width-24)
	staBar := m.renderHPBar(playerStamina, playerMaxStamina, width-24)
	manaBar := m.renderHPBar(playerMana, playerMaxMana, width-24)
	content.WriteString(fmt.Sprintf("HP:   %s %d/%d\n", hpBar, playerHP, playerMaxHP))
	content.WriteString(fmt.Sprintf("STA:  %s %d/%d\n", staBar, playerStamina, playerMaxStamina))
	content.WriteString(fmt.Sprintf("MANA: %s %d/%d\n", statsStyle.Render(manaBar), playerMana, playerMaxMana))
	content.WriteString("\n")

	// Queued action indicator
	if m.combatQueuedAction != "" {
		queueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00b894")).Bold(true)
		content.WriteString(queueStyle.Render(fmt.Sprintf("> Queued: %s", m.combatQueuedAction)))
		content.WriteString("\n\n")
	}

	// Combat log
	if len(m.combatLog) > 0 {
		logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#dfe6e9"))
		content.WriteString(logStyle.Render("Combat Log:"))
		content.WriteString("\n")
		maxLines := 10
		startIdx := len(m.combatLog) - maxLines
		if startIdx < 0 {
			startIdx = 0
		}
		for i := startIdx; i < len(m.combatLog); i++ {
			content.WriteString(fmt.Sprintf("  %s\n", m.combatLog[i]))
		}
	}

	// Render viewport
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

	// Skill HUD bar
	s.WriteString(m.renderCombatHUD(width))

	// Input area
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(inputHeight - 2)
	actionPrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6b6b")).Render("Action: ")
	s.WriteString(inputStyle.Render(actionPrompt + m.textInput.View()))

	return s.String()
}

// renderCombatHUD builds the skill + potion HUD bar shown during combat.
// Layout: [1 Skill1] [2 Skill2] [3 Skill3] [4 Skill4] | [5 Potion]
func (m *model) renderCombatHUD(width int) string {
	hudStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Width(width).
		Padding(0, 1)

	// Build each skill slot
	slots := make([]string, 4)
	for i := 1; i <= 4; i++ {
		name := m.formatSkillSlot(i)
		cd := m.getSkillCooldown(i)
		keyLabel := fmt.Sprintf("%d", i)

		var slotStr string
		if cd > 0 {
			// On cooldown: dim the slot
			slotStr = fmt.Sprintf("[%s] %s (CD:%d)", keyLabel, lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(name), cd)
		} else if m.combatSkills != nil && m.combatSkills.EquippedSkill[i-1].ID != 0 {
			// Available: bright
			slotStr = fmt.Sprintf("[%s] %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#74b9ff")).Bold(true).Render(keyLabel), lipgloss.NewStyle().Foreground(lipgloss.Color("#00b894")).Render(name))
		} else {
			// Empty
			slotStr = fmt.Sprintf("[%s] %s", lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(keyLabel), lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(name))
		}
		slots[i-1] = slotStr
	}

	// Potion slot (key 5)
	potionName := m.getPotionSlotName()
	potionStr := fmt.Sprintf("[5] %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#fdcb6e")).Render(potionName))

	// Separator
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("|")

	// Assemble: skills | potion | Tab hint
	hudLine := fmt.Sprintf("%s  %s  %s  %s  %s  %s",
		slots[0], slots[1], slots[2], slots[3], sep, potionStr)

	// Tab toggle hint
	tabHint := lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render("Tab:room")
	hudLine += "  " + tabHint

	// Flee hint
	fleeHint := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6b6b")).Render("Q:flee")
	hudLine += " " + fleeHint

	// Second line: prompt hint
	promptLine := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("1-5:use skill  R:potion  Tab:switch view  Q:flee")

	return hudStyle.Render(hudLine) + "\n" + hudStyle.Render(promptLine)
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
