package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// STATIC SCREENS
// ============================================================

func welcomeScreen(width, height int, cursor int, inputView string) string {
	if width < 40 {
		width = 80
	}

	// Calculate split
	outputHeight := height * 60 / 100
	if outputHeight < 8 {
		outputHeight = 8
	}

	var outputContent strings.Builder

	// Title banner
	outputContent.WriteString("\n")
	titleLine := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryGold).
		Render("✦ HERBST MUD ✦")
	outputContent.WriteString(lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(titleLine))
	outputContent.WriteString("\n\n")

	// Subtitle
	subtitleLine := lipgloss.NewStyle().
		Foreground(PrimaryPurple).
		Italic(true).
		Render("A World of Adventure Awaits")
	outputContent.WriteString(lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(subtitleLine))
	outputContent.WriteString("\n\n")

	// Menu
	menuItems := []struct {
		key  string
		desc string
	}{
		{"1", "Login"},
		{"2", "Register"},
		{"3", "Quit"},
	}
	for i, item := range menuItems {
		cursorStr := " "
		keyStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(item.key + ".")
		nameStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(item.desc)
		if i == cursor {
			cursorStr = lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render("▸")
			keyStyle = lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(item.key + ".")
			nameStyle = lipgloss.NewStyle().Foreground(PrimaryGold).Render(item.desc)
		}
		outputContent.WriteString(fmt.Sprintf("  %s  %s  %s\n", cursorStr, keyStyle, nameStyle))
	}
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().
		Foreground(TextGray).
		Width(width).
		Align(lipgloss.Center).
		Render("Type a number or command"))

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryPurple).
		Padding(0, 1).
		Width(width).
		Height(outputHeight)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryGold).
		Padding(0, 1).
		Width(width).
		Height(height - outputHeight - 1)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))
	return sb.String()
}

func loginScreen(width, height int, message, messageType string, inputView string) string {
	if width < 40 {
		width = 80
	}

	outputHeight := height * 55 / 100
	if outputHeight < 8 {
		outputHeight = 8
	}

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryGold).
		Width(width).
		Align(lipgloss.Center).
		Render("✦ ACCOUNT LOGIN ✦")
	outputContent.WriteString(title)
	outputContent.WriteString("\n\n")

	if message != "" {
		styled := styleMessage(message, messageType)
		outputContent.WriteString(styled)
		outputContent.WriteString("\n\n")
	}

	instructions := lipgloss.NewStyle().
		Foreground(TextGray).
		Render("Type 'register' to create a new account")
	outputContent.WriteString(instructions)
	outputContent.WriteString("\n")
	quitText := lipgloss.NewStyle().
		Foreground(TextGray).
		Render("Type 'quit' or press Ctrl+C to exit")
	outputContent.WriteString(quitText)

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryPurple).
		Padding(0, 1).
		Width(width).
		Height(outputHeight)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryGold).
		Padding(0, 1).
		Width(width).
		Height(height - outputHeight - 1)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))
	return sb.String()
}

func registerScreen(width, height int, message, messageType string, inputView string) string {
	if width < 40 {
		width = 80
	}

	outputHeight := height * 55 / 100
	if outputHeight < 8 {
		outputHeight = 8
	}

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryGold).
		Width(width).
		Align(lipgloss.Center).
		Render("✦ CREATE ACCOUNT ✦")
	outputContent.WriteString(title)
	outputContent.WriteString("\n\n")

	if message != "" {
		styled := styleMessage(message, messageType)
		outputContent.WriteString(styled)
		outputContent.WriteString("\n\n")
	}

	instructions := lipgloss.NewStyle().
		Foreground(TextGray).
		Render("Press Esc to go back")
	outputContent.WriteString(instructions)

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryPurple).
		Padding(0, 1).
		Width(width).
		Height(outputHeight)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryGold).
		Padding(0, 1).
		Width(width).
		Height(height - outputHeight - 1)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))
	return sb.String()
}

func worldSelectScreen(width, height int, displayContent string, inputView string) string {
	if width < 40 {
		width = 80
	}

	outputHeight := height * 55 / 100
	if outputHeight < 8 {
		outputHeight = 8
	}

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryGold).
		Width(width).
		Align(lipgloss.Center).
		Render("✦ SELECT WORLD ✦")
	outputContent.WriteString(title)
	outputContent.WriteString("\n\n")
	outputContent.WriteString(displayContent)
	outputContent.WriteString("\n")
	hint := lipgloss.NewStyle().
		Foreground(TextGray).
		Render("j/k navigate · enter/1-9 select · 'b' back")
	outputContent.WriteString(hint)

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryPurple).
		Padding(0, 1).
		Width(width).
		Height(outputHeight)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryGold).
		Padding(0, 1).
		Width(width).
		Height(height - outputHeight - 1)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))
	return sb.String()
}

func characterSelectScreen(width, height int, displayContent string, inputView string) string {
	if width < 40 {
		width = 80
	}

	outputHeight := height * 55 / 100
	if outputHeight < 8 {
		outputHeight = 8
	}

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryGold).
		Width(width).
		Align(lipgloss.Center).
		Render("✦ SELECT CHARACTER ✦")
	outputContent.WriteString(title)
	outputContent.WriteString("\n\n")

	outputContent.WriteString(displayContent)
	outputContent.WriteString("\n")
	hint := lipgloss.NewStyle().
		Foreground(TextGray).
		Render("j/k navigate · enter/1-9 select · 'n' new · 'b' back")
	outputContent.WriteString(hint)

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryPurple).
		Padding(0, 1).
		Width(width).
		Height(outputHeight)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryGold).
		Padding(0, 1).
		Width(width).
		Height(height - outputHeight - 1)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))
	return sb.String()
}

// ============================================================
// STATUS BARS
// ============================================================

// ProgressBar creates a text-based progress bar
func ProgressBar(current, max, width int, filledChar, emptyChar string, fillColor, emptyColor lipgloss.Color) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	filledWidth := int(float64(current) / float64(max) * float64(width))
	emptyWidth := width - filledWidth

	filledStyle := lipgloss.NewStyle().Foreground(fillColor)
	emptyStyle := lipgloss.NewStyle().Foreground(emptyColor)

	return filledStyle.Render(strings.Repeat(filledChar, filledWidth)) +
		emptyStyle.Render(strings.Repeat(emptyChar, emptyWidth))
}

// StatusBar creates colorful status bars with progress bars
func StatusBar(hp, maxHP, stamina, maxStamina, mana, maxMana int) string {
	hpBar := ProgressBar(hp, maxHP, 20, "█", "░", red, gray)
	staminaBar := ProgressBar(stamina, maxStamina, 20, "▓", "░", yellow, gray)
	manaBar := ProgressBar(mana, maxMana, 20, "✨", "○", blue, gray)

	return fmt.Sprintf(" %s  HP: %s %d/%d\n %s  STA: %s %d/%d\n %s  MANA: %s %d/%d",
		lipgloss.NewStyle().Foreground(red).Render("❤️"),
		hpBar, hp, maxHP,
		lipgloss.NewStyle().Foreground(yellow).Render("💪"),
		staminaBar, stamina, maxStamina,
		lipgloss.NewStyle().Foreground(blue).Render("✨"),
		manaBar, mana, maxMana)
}

// MiniStatusBar creates a compact inline status bar
func MiniStatusBar(hp, maxHP, stamina, maxStamina, mana, maxMana int, regenActive bool) string {
	hpPercent := float64(hp) / float64(maxHP) * 100
	staminaPercent := float64(stamina) / float64(maxStamina) * 100
	manaPercent := float64(mana) / float64(maxMana) * 100

	hpColor := green
	if hpPercent <= 30 {
		hpColor = red
	} else if hpPercent <= 60 {
		hpColor = yellow
	}

	hpStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", hpPercent))
	staStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", staminaPercent))
	manaStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", manaPercent))

	// Use single-width characters to avoid alignment issues
	var hpIcon, staIcon, manaIcon string
	if regenActive {
		hpIcon = lipgloss.NewStyle().Foreground(green).Render("+") // + for regen
	} else {
		hpIcon = lipgloss.NewStyle().Foreground(hpColor).Render("H") // H for HP
	}
	staIcon = lipgloss.NewStyle().Foreground(yellow).Render("S") // S for Stamina
	manaIcon = lipgloss.NewStyle().Foreground(blue).Render("M") // M for Mana

	// For users who prefer emoji, we can use these with proper spacing
	// Note: Emojis are double-width, so we don't add extra spaces
	// hpIcon = lipgloss.NewStyle().Foreground(hpColor).Render("❤")
	// staIcon = lipgloss.NewStyle().Foreground(yellow).Render("⚡")
	// manaIcon = lipgloss.NewStyle().Foreground(blue).Render("✦")

	return fmt.Sprintf("[%s%s%% %s%s%% %s%s%%]",
		hpIcon, hpStr,
		staIcon, staStr,
		manaIcon, manaStr)
}
