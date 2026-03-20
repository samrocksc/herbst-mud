package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// STYLING - Lipgloss styles for UI elements
// ============================================================

var (
	// Colors
	red    = lipgloss.Color("196")
	green  = lipgloss.Color("46")
	yellow = lipgloss.Color("226")
	blue   = lipgloss.Color("75")
	purple = lipgloss.Color("141")
	white  = lipgloss.Color("15")
	gray   = lipgloss.Color("8")
	pink   = lipgloss.Color("219")
	cyan   = lipgloss.Color("51")

	// Raw ANSI for direct terminal output (when lipgloss fails)
	pinkAnsi  = "\033[38;5;219m"
	pinkReset = "\033[0m"

	// Exit colors for visited/known/new
	exitVisitedColor = lipgloss.Color("46")  // Green
	exitKnownColor   = lipgloss.Color("226") // Yellow
	exitNewColor     = lipgloss.Color("15")  // White

	// Quest tracker panel colors
	questTitleColor     = lipgloss.Color("75")    // Blue
	questProgressColor  = lipgloss.Color("226")   // Yellow
	questCompletedColor = lipgloss.Color("46")    // Green
	questAvailableColor = lipgloss.Color("141")  // Purple

	// Quest tracker panel styles
	questTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(questTitleColor)

	questBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2)

	questProgressStyle = lipgloss.NewStyle().
				Foreground(questProgressColor)

	questCompletedStyle = lipgloss.NewStyle().
				Foreground(questCompletedColor).
				Strikethrough(true)

	questAvailableStyle = lipgloss.NewStyle().
				Foreground(questAvailableColor)

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(yellow)

	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(green).
				Bold(true).
				Padding(0, 0, 0, 2)

	menuNormalStyle = lipgloss.NewStyle().
			Foreground(gray).
			Padding(0, 0, 0, 2)

	promptStyle = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)
)

// ============================================================
// PROGRESS BARS - HP/Stamina/Mana using characters
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
	// Use Unicode block characters for progress
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

	var hpColor lipgloss.Color
	if hpPercent > 60 {
		hpColor = green
	} else if hpPercent > 30 {
		hpColor = yellow
	} else {
		hpColor = red
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

// Item colors for display
var (
	itemColorGold   = lipgloss.Color("220") // Gold for immovable items
	itemColorWeapon = lipgloss.Color("196") // Red for weapons
	itemColorArmor  = lipgloss.Color("75")  // Blue for armor
	itemColorMisc   = lipgloss.Color("242") // Gray for misc
)

// combatDamageStyle for damage messages (red)
var combatDamageStyle = lipgloss.NewStyle().
	Foreground(red).
	Bold(true)

// combatHealStyle for healing messages (green)
var combatHealStyle = lipgloss.NewStyle().
	Foreground(green).
	Bold(true)

// getItemIcon returns an emoji icon based on item type
func getItemIcon(itemType string) string {
	switch itemType {
	case "weapon":
		return "⚔️"
	case "armor":
		return "🛡️"
	case "potion":
		return "🧪"
	case "food":
		return "🍖"
	case "scroll":
		return "📜"
	case "key":
		return "🔑"
	case "treasure":
		return "💎"
	case "quest":
		return "📋"
	default:
		return "📦"
	}
}

// getItemRarityColor returns a lipgloss color based on item rarity
func getItemRarityColor(rarity string) lipgloss.Color {
	switch rarity {
	case "rare":
		return lipgloss.Color("51") // Blue
	case "epic":
		return lipgloss.Color("201") // Magenta
	case "legendary":
		return lipgloss.Color("220") // Gold
	default:
		return lipgloss.Color("white")
	}
}
