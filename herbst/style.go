package main

// ============================================================
// COMPREHENSIVE STYLING SYSTEM - Fantasy MUD Theme
// ============================================================
// - Dark theme with fantasy colors
// - Color scheme: purple (#AF00FF), gold (#FFD700), blue (#5F87FF)
// - Use borders to separate panels
// - Progress bars for HP/Stamina/Mana

import "github.com/charmbracelet/lipgloss"

// ============================================================
// COLOR PALETTE - Fantasy theme
// ============================================================

var (
	// Fantasy primary colors
	PrimaryPurple = lipgloss.Color("129") // #AF00FF - Deep purple
	PrimaryGold   = lipgloss.Color("220") // #FFD700 - Rich gold
	AccentBlue    = lipgloss.Color("75")  // #5F87FF - Bright blue

	// Status colors (keep same semantics, adjust shades)
	StatusRed     = lipgloss.Color("196") // Red - danger/low HP
	StatusGreen   = lipgloss.Color("46")  // Green - healthy/good
	StatusYellow  = lipgloss.Color("226") // Yellow - warning
	StatusBlue    = lipgloss.Color("75")  // Blue - mana/magic

	// Neutral colors
	TextWhite    = lipgloss.Color("15")  // Primary text
	TextGray     = lipgloss.Color("8")   // Secondary text
	TextDarkGray = lipgloss.Color("236") // Dark backgrounds
	BorderColor  = lipgloss.Color("240") // Borders

	// Exit colors (for navigation)
	ExitVisited = lipgloss.Color("46")  // Green - visited
	ExitKnown   = lipgloss.Color("220") // Gold - known
	ExitNew     = lipgloss.Color("15")  // White - new
)

// Header bar colors
var (
	HeaderBg    = lipgloss.Color("54")  // Dark purple background
	HeaderFg    = lipgloss.Color("220") // Gold text
	StatusBarBg = lipgloss.Color("236") // Dark bg for status
)

// ============================================================
// TYPOGRAPHY STYLES
// ============================================================

var (
	// Title style - bold, gold on dark
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryGold).
			Background(TextDarkGray).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center)

	// Subtitle style - purple
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryPurple).
			Padding(0, 1).
			Align(lipgloss.Center)

	// Header style - bold blue
	HeaderStyle = lipgloss.NewStyle().
			Foreground(AccentBlue).
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
// HEADER BAR STYLES
// ============================================================

var (
	HeaderBarStyle = lipgloss.NewStyle().
			Background(HeaderBg).
			Foreground(HeaderFg).
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center)

	HeaderBarRoomStyle = lipgloss.NewStyle().
			Foreground(HeaderFg).
			Bold(true)

	HeaderBarInfoStyle = lipgloss.NewStyle().
			Foreground(AccentBlue)
)

// ============================================================
// BORDER STYLES
// ============================================================

var (
	// Rounded border box - for main panels
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryPurple).
			Padding(1, 2).
			Align(lipgloss.Left)

	// Double border - for important sections
	DoubleBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(PrimaryGold).
			Padding(1, 2)

	// Thick border - for headers
	ThickBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(PrimaryGold).
			Padding(1, 2)

	// Left border highlight - for highlighted items
	HighlightBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(PrimaryPurple).
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
	// Selected menu item - purple
	MenuSelectedStyle = lipgloss.NewStyle().
				Foreground(PrimaryPurple).
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
			Foreground(PrimaryPurple).
			Background(lipgloss.Color("56"))
)

// ============================================================
// INPUT/OUTPUT STYLES
// ============================================================

var (
	// Prompt style - gold
	PromptStyle = lipgloss.NewStyle().
			Foreground(PrimaryGold).
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
				Foreground(PrimaryPurple).
				Background(TextDarkGray).
				Border(lipgloss.NormalBorder()).
				BorderForeground(PrimaryPurple).
				Padding(0, 1)

	// Password input (masked)
	PasswordStyle = lipgloss.NewStyle().
			Foreground(PrimaryGold).
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
		return lipgloss.NewStyle().Foreground(PrimaryPurple).Bold(true).Render(name)
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
		color = PrimaryPurple
	case "legendary":
		color = PrimaryGold
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
		return lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true)
	case "available":
		return lipgloss.NewStyle().Foreground(StatusBlue)
	default:
		return lipgloss.NewStyle().Foreground(TextGray)
	}
}

// Backward-compatibility aliases for code that uses lowercase names
var (
	pink                = PrimaryPurple // was AccentPink, now purple
	cyan                = AccentBlue    // was PrimaryTeal, now blue
	gray                = TextGray
	green               = StatusGreen
	yellow              = StatusYellow
	blue                = StatusBlue
	red                 = StatusRed
	white               = TextWhite
	purple              = PrimaryPurple
	exitVisitedColor    = ExitVisited
	exitKnownColor      = ExitKnown
	exitNewColor        = ExitNew
	dialogNPCStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#87CEEB"))
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
	itemColorGold       = PrimaryGold
	itemColorWeapon     = StatusRed
	itemColorArmor      = StatusBlue
	itemColorMisc       = TextGray
)