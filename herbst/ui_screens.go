package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// STATIC SCREENS
// ============================================================

func welcomeScreen(width, height int, inputView string) string {
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(outputHeight - 2)

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("        Welcome Adventurer!        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  1. Login"))
	outputContent.WriteString("      - Log in to your existing account\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  2. Register"))
	outputContent.WriteString("   - Create a new character\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  3. Quit"))
	outputContent.WriteString("       - Exit the game\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(gray).Render("  Use arrow keys or type number/command"))

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(inputHeight - 2)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}

func loginScreen(width, height int, message, messageType string, inputView string) string {
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(outputHeight - 2)

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("            LOGIN            "))
	outputContent.WriteString("\n\n")

	if message != "" {
		outputContent.WriteString(styleMessage(message, messageType))
		outputContent.WriteString("\n")
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(inputHeight - 2)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}

func registerScreen(width, height int, message, messageType string, inputView string) string {
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(outputHeight - 2)

	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("        CREATE ACCOUNT        "))
	outputContent.WriteString("\n\n")

	if message != "" {
		outputContent.WriteString(styleMessage(message, messageType))
		outputContent.WriteString("\n")
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width).
		Height(inputHeight - 2)

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
func MiniStatusBar(hp, maxHP, stamina, maxStamina, mana, maxMana int) string {
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

	return fmt.Sprintf("[%s%s%% %s%s%% %s%s%%]",
		lipgloss.NewStyle().Foreground(hpColor).Render("❤️"),
		hpStr,
		lipgloss.NewStyle().Foreground(yellow).Render("💪"),
		staStr,
		lipgloss.NewStyle().Foreground(blue).Render("✨"),
		manaStr)
}
