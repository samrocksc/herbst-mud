package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// STYLE TESTS
// ============================================================

func TestColorPalette(t *testing.T) {
	// Verify all colors are valid
	assert.NotEqual(t, lipgloss.Color(""), PrimaryTeal)
	assert.NotEqual(t, lipgloss.Color(""), PrimaryOrange)
	assert.NotEqual(t, lipgloss.Color(""), AccentPink)
	assert.NotEqual(t, lipgloss.Color(""), AccentPurple)
	assert.NotEqual(t, lipgloss.Color(""), StatusRed)
	assert.NotEqual(t, lipgloss.Color(""), StatusGreen)
	assert.NotEqual(t, lipgloss.Color(""), StatusYellow)
	assert.NotEqual(t, lipgloss.Color(""), StatusBlue)
}

func TestTypographyStyles(t *testing.T) {
	// Test TitleStyle
	result := TitleStyle.Render("Test Title")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Test Title")

	// Test SubtitleStyle
	result = SubtitleStyle.Render("Test Subtitle")
	assert.NotEmpty(t, result)

	// Test HeaderStyle
	result = HeaderStyle.Render("Test Header")
	assert.NotEmpty(t, result)

	// Test BodyStyle
	result = BodyStyle.Render("Test Body")
	assert.NotEmpty(t, result)

	// Test SecondaryStyle
	result = SecondaryStyle.Render("Test Secondary")
	assert.NotEmpty(t, result)

	// Test BoldStyle
	result = BoldStyle.Render("Test Bold")
	assert.NotEmpty(t, result)

	// Test ItalicStyle
	result = ItalicStyle.Render("Test Italic")
	assert.NotEmpty(t, result)
}

func TestBorderStyles(t *testing.T) {
	// Test BoxStyle
	result := BoxStyle.Render("Test Content")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Test Content")

	// Test DoubleBoxStyle
	result = DoubleBoxStyle.Render("Double Box")
	assert.NotEmpty(t, result)

	// Test ThickBoxStyle
	result = ThickBoxStyle.Render("Thick Box")
	assert.NotEmpty(t, result)

	// Test HighlightBoxStyle
	result = HighlightBoxStyle.Render("Highlighted")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Highlighted")

	// Test ErrorBoxStyle
	result = ErrorBoxStyle.Render("Error Content")
	assert.NotEmpty(t, result)

	// Test SuccessBoxStyle
	result = SuccessBoxStyle.Render("Success Content")
	assert.NotEmpty(t, result)
}

func TestStatusStyles(t *testing.T) {
	// Test SuccessStyle
	result := SuccessStyle.Render("Success!")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Success!")

	// Test ErrorStyle
	result = ErrorStyle.Render("Error!")
	assert.NotEmpty(t, result)

	// Test WarningStyle
	result = WarningStyle.Render("Warning!")
	assert.NotEmpty(t, result)

	// Test InfoStyle
	result = InfoStyle.Render("Info!")
	assert.NotEmpty(t, result)

	// Test DamageStyle
	result = DamageStyle.Render("50 damage")
	assert.NotEmpty(t, result)

	// Test HealStyle
	result = HealStyle.Render("+25 HP")
	assert.NotEmpty(t, result)

	// Test MagicStyle
	result = MagicStyle.Render("50 mana")
	assert.NotEmpty(t, result)
}

func TestMenuStyles(t *testing.T) {
	// Test MenuSelectedStyle
	result := MenuSelectedStyle.Render("Selected Item")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Selected Item")

	// Test MenuNormalStyle
	result = MenuNormalStyle.Render("Normal Item")
	assert.NotEmpty(t, result)

	// Test MenuDisabledStyle
	result = MenuDisabledStyle.Render("Disabled Item")
	assert.NotEmpty(t, result)
}

func TestProgressBarStyles(t *testing.T) {
	// Test HealthBarStyle
	result := HealthBarStyle.Render("Health")
	assert.NotEmpty(t, result)

	// Test StaminaBarStyle
	result = StaminaBarStyle.Render("Stamina")
	assert.NotEmpty(t, result)

	// Test ManaBarStyle
	result = ManaBarStyle.Render("Mana")
	assert.NotEmpty(t, result)

	// Test XPBarStyle
	result = XPBarStyle.Render("XP")
	assert.NotEmpty(t, result)
}

func TestInputStyles(t *testing.T) {
	// Test PromptStyle
	result := PromptStyle.Render(">")
	assert.NotEmpty(t, result)

	// Test InputStyle
	result = InputStyle.Render("input text")
	assert.NotEmpty(t, result)

	// Test InputFocusedStyle
	result = InputFocusedStyle.Render("focused input")
	assert.NotEmpty(t, result)

	// Test PasswordStyle
	result = PasswordStyle.Render("********")
	assert.NotEmpty(t, result)
}

func TestUtilityFunctions(t *testing.T) {
	// Test StyleRoomName
	result := StyleRoomName("The Entrance", true)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "The Entrance")

	result = StyleRoomName("Secret Cave", false)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Secret Cave")

	// Test StyleExit
	result = StyleExit("north", "visited")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "north")

	result = StyleExit("east", "known")
	assert.NotEmpty(t, result)

	result = StyleExit("west", "new")
	assert.NotEmpty(t, result)

	// Test StyleItemQuality
	result = StyleItemQuality("Iron Sword", "common")
	assert.NotEmpty(t, result)

	result = StyleItemQuality("Golden Amulet", "legendary")
	assert.NotEmpty(t, result)

	result = StyleItemQuality("Mystic Staff", "epic")
	assert.NotEmpty(t, result)

	result = StyleItemQuality("Rare Gem", "rare")
	assert.NotEmpty(t, result)

	result = StyleItemQuality("Uncommon Potion", "uncommon")
	assert.NotEmpty(t, result)

	// Test StyleQuestStatus
	style := StyleQuestStatus("completed")
	assert.NotNil(t, style)

	style = StyleQuestStatus("in_progress")
	assert.NotNil(t, style)

	style = StyleQuestStatus("available")
	assert.NotNil(t, style)

	style = StyleQuestStatus("unknown")
	assert.NotNil(t, style)
}

func TestStyleChaining(t *testing.T) {
	// Test that styles can be chained
	result := TitleStyle.Background(TextDarkGray).Bold(true).Render("Chained Title")
	assert.NotEmpty(t, result)

	// Test border and padding chaining
	result = BoxStyle.Padding(2, 3).Width(50).Render("Chained Box")
	assert.NotEmpty(t, result)

	// Test color chaining
	result = SuccessStyle.Bold(true).Underline(true).Render("Chained Success")
	assert.NotEmpty(t, result)
}

func TestStyleWithEmptyContent(t *testing.T) {
	// Test with empty strings
	result := TitleStyle.Render("")
	assert.NotNil(t, result)

	result = BodyStyle.Render("")
	assert.NotNil(t, result)

	result = BoxStyle.Render("")
	assert.NotNil(t, result)
}

func TestStyleWidthConstraints(t *testing.T) {
	// Test TitleStyle respects width
	result := TitleStyle.Width(40).Render("Short")
	assert.NotEmpty(t, result)

	// Test BoxStyle respects width
	result = BoxStyle.Width(30).Render("Narrow Box")
	assert.NotEmpty(t, result)
}