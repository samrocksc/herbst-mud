package main

// ============================================================
// COMPREHENSIVE STYLING SYSTEM - Post-Apocalyptic MUD Theme
// ============================================================
// Following the style guide from charm-ui-tickets.md:
// - Dark theme with neon accents
// - Color scheme: teal (#00AAAA), orange (#FF6600), white text
// - Use borders to separate panels
// - Progress bars for HP/Stamina/Mana

import "github.com/charmbracelet/lipgloss"

// ============================================================
// COLOR PALETTE - Post-apocalyptic neon theme
// ============================================================

var (
	// Primary colors (neon accents for dark background)
	PrimaryTeal    = lipgloss.Color("51")   // #00CCCC - Teal/Cyan
	PrimaryOrange  = lipgloss.Color("202")  // #FF6600 - Orange
	AccentPink     = lipgloss.Color("219")  // #FF99CC - Neon pink
	AccentPurple   = lipgloss.Color("141")  // #CC99FF - Purple

	// Status colors
	StatusRed     = lipgloss.Color("196") // Red - danger/low HP
	StatusGreen   = lipgloss.Color("46")  // Green - healthy/good
	StatusYellow  = lipgloss.Color("226") // Yellow - warning
	StatusBlue    = lipgloss.Color("75")  // Blue - mana/magic

	// Neutral colors
	TextWhite    = lipgloss.Color("15")  // Primary text
	TextGray     = lipgloss.Color("8")   // Secondary text
	TextDarkGray = lipgloss.Color("236") // Backgrounds
	BorderColor  = lipgloss.Color("240") // Borders

	// Exit colors (for navigation)
	ExitVisited = lipgloss.Color("46") // Green - visited
	ExitKnown   = lipgloss.Color("226") // Yellow - known
	ExitNew     = lipgloss.Color("15") // White - new
)

// ============================================================
// TYPOGRAPHY STYLES
// ============================================================

var (
	// Title style - bold, teal on dark
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryTeal).
			Background(TextDarkGray).
			Bold(true).
			Padding(0, 1).
			Width(60).
			Align(lipgloss.Center)

	// Subtitle style - orange
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryOrange).
			Padding(0, 1).
			Width(60).
			Align(lipgloss.Center)

	// Header style - bold blue
	HeaderStyle = lipgloss.NewStyle().
			Foreground(StatusBlue).
			Bold(true).
			Padding(0, 1)

	// Body text - white
	BodyStyle = lipgloss.NewStyle().
			Foreground(TextWhite).
			Padding(0, 1)

	// Secondary text - gray
	SecondaryStyle = lipgloss.NewStyle().
			Foreground(TextGray).
			Padding(0, 1)

	// Bold text
	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextWhite)

	// Italic text
	ItalicStyle = lipgloss.NewStyle().
			Foreground(TextGray).
			Italic(true)
)

// ============================================================
// BORDER STYLES
// ============================================================

var (
	// Rounded border box - for main panels
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Width(58).
			Align(lipgloss.Left)

	// Double border - for important sections
	DoubleBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(PrimaryTeal).
			Padding(1, 2).
			Width(58)

	// Thick border - for headers
	ThickBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(PrimaryOrange).
			Padding(1, 2).
			Width(58)

	// Left border highlight - for highlighted items
	HighlightBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(PrimaryTeal).
			Padding(0, 0, 0, 1)

	// Error/warning box
	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(StatusRed).
			Padding(1, 2).
			Width(58)

	// Success box
	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(StatusGreen).
			Padding(1, 2).
			Width(58)
)

// ============================================================
// STATUS STYLES - For messages and feedback
// ============================================================

var (
	// Success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(StatusGreen).
			Bold(true)

	// Error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(StatusRed).
			Bold(true)

	// Warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(StatusYellow).
			Bold(true)

	// Info messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(StatusBlue)

	// Combat damage - red
	DamageStyle = lipgloss.NewStyle().
			Foreground(StatusRed)

	// Combat heal - green
	HealStyle = lipgloss.NewStyle().
			Foreground(StatusGreen)

	// Magic/mana - blue
	MagicStyle = lipgloss.NewStyle().
			Foreground(StatusBlue)
)

// ============================================================
// MENU STYLES
// ============================================================

var (
	// Selected menu item - teal background
	MenuSelectedStyle = lipgloss.NewStyle().
				Foreground(PrimaryTeal).
				Bold(true).
				Padding(0, 0, 0, 2)

	// Normal menu item - gray
	MenuNormalStyle = lipgloss.NewStyle().
			Foreground(TextGray).
			Padding(0, 0, 0, 2)

	// Disabled menu item - dark gray
	MenuDisabledStyle = lipgloss.NewStyle().
			Foreground(TextDarkGray).
			Padding(0, 0, 0, 2)
)

// ============================================================
// PROGRESS BAR STYLES
// ============================================================

var (
	// Health bar - red filled
	HealthBarStyle = lipgloss.NewStyle().
			Foreground(StatusRed).
			Background(lipgloss.Color("52"))

	// Stamina bar - yellow filled
	StaminaBarStyle = lipgloss.NewStyle().
			Foreground(StatusYellow).
			Background(lipgloss.Color("58"))

	// Mana bar - blue filled
	ManaBarStyle = lipgloss.NewStyle().
			Foreground(StatusBlue).
			Background(lipgloss.Color("24"))

	// Experience bar - purple filled
	XPBarStyle = lipgloss.NewStyle().
			Foreground(AccentPurple).
			Background(lipgloss.Color("56"))
)

// ============================================================
// INPUT/OUTPUT STYLES
// ============================================================

var (
	// Prompt style - teal
	PromptStyle = lipgloss.NewStyle().
			Foreground(PrimaryTeal).
			Bold(true)

	// Input field
	InputStyle = lipgloss.NewStyle().
			Foreground(TextWhite).
			Background(TextDarkGray).
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	// Input field focused
	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(PrimaryTeal).
				Background(TextDarkGray).
				Border(lipgloss.NormalBorder()).
				BorderForeground(PrimaryTeal).
				Padding(0, 1)

	// Password input (masked)
	PasswordStyle = lipgloss.NewStyle().
			Foreground(PrimaryOrange).
			Background(TextDarkGray).
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)
)

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// StyleRoomName applies the appropriate style to a room name
func StyleRoomName(name string, visited bool) string {
	if visited {
		return lipgloss.NewStyle().Foreground(PrimaryTeal).Bold(true).Render(name)
	}
	return lipgloss.NewStyle().Foreground(TextWhite).Bold(true).Render(name)
}

// StyleExit applies color based on exit status
func StyleExit(direction string, status string) string {
	var color lipgloss.Color
	switch status {
	case "visited":
		color = ExitVisited
	case "known":
		color = ExitKnown
	default:
		color = ExitNew
	}
	return lipgloss.NewStyle().Foreground(color).Render(direction)
}

// StyleItemQuality applies color based on item rarity
func StyleItemQuality(name string, quality string) string {
	var color lipgloss.Color
	switch quality {
	case "common":
		color = TextWhite
	case "uncommon":
		color = StatusGreen
	case "rare":
		color = StatusBlue
	case "epic":
		color = AccentPurple
	case "legendary":
		color = PrimaryOrange
	default:
		color = TextGray
	}
	return lipgloss.NewStyle().Foreground(color).Render(name)
}

// StyleQuestStatus applies color based on quest status
func StyleQuestStatus(status string) lipgloss.Style {
	switch status {
	case "completed":
		return lipgloss.NewStyle().Foreground(StatusGreen).Bold(true)
	case "in_progress":
		return lipgloss.NewStyle().Foreground(PrimaryOrange).Bold(true)
	case "available":
		return lipgloss.NewStyle().Foreground(StatusBlue)
	default:
		return lipgloss.NewStyle().Foreground(TextGray)
	}
}

// Backward-compatibility aliases for code that uses lowercase names
var (
	pink                = AccentPink
	cyan                = PrimaryTeal
	gray                = TextGray
	green               = StatusGreen
	yellow              = StatusYellow
	blue                = StatusBlue
	red                 = StatusRed
	white               = TextWhite
	purple              = AccentPurple
	exitVisitedColor    = ExitVisited
	exitKnownColor      = ExitKnown
	exitNewColor        = ExitNew
	questTitleStyle     = TitleStyle
	questBoxStyle       = BoxStyle
	questProgressStyle  = WarningStyle
	questCompletedStyle = SuccessStyle
	questAvailableStyle = InfoStyle
	titleStyle          = TitleStyle
	boxStyle            = BoxStyle
	successStyle        = SuccessStyle
	errorStyle          = ErrorStyle
	infoStyle           = InfoStyle
	menuSelectedStyle   = MenuSelectedStyle
	menuNormalStyle     = MenuNormalStyle
	promptStyle         = PromptStyle
	combatDamageStyle   = DamageStyle
	combatHealStyle     = HealStyle
	itemColorGold       = PrimaryOrange
	itemColorWeapon     = StatusRed
	itemColorArmor      = StatusBlue
	itemColorMisc       = TextGray
)