package herbst

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestLipGlossStylesExist verifies that all style definitions are properly created
// This tests that the Lip Gloss styling work from PR #205 is correctly implemented
func TestLipGlossStylesExist(t *testing.T) {
	// Test that titleStyle is defined and renders
	titleRendered := titleStyle.Render("Test Title")
	if lipgloss.Width(titleRendered) == 0 {
		t.Error("titleStyle should render non-empty content")
	}

	// Test that headerStyle is defined and renders
	headerRendered := headerStyle.Render("Test Header")
	if lipgloss.Width(headerRendered) == 0 {
		t.Error("headerStyle should render non-empty content")
	}

	// Test that boxStyle is defined and renders with border
	boxRendered := boxStyle.Render("Test Box")
	if lipgloss.Width(boxRendered) == 0 {
		t.Error("boxStyle should render non-empty content")
	}

	// Test successStyle
	successRendered := successStyle.Render("Success")
	if lipgloss.Width(successRendered) == 0 {
		t.Error("successStyle should render non-empty content")
	}

	// Test errorStyle
	errorRendered := errorStyle.Render("Error")
	if lipgloss.Width(errorRendered) == 0 {
		t.Error("errorStyle should render non-empty content")
	}

	// Test infoStyle
	infoRendered := infoStyle.Render("Info")
	if lipgloss.Width(infoRendered) == 0 {
		t.Error("infoStyle should render non-empty content")
	}

	// Test menu styles
	menuSelected := menuSelectedStyle.Render("Selected")
	menuNormal := menuNormalStyle.Render("Normal")
	if lipgloss.Width(menuSelected) == 0 || lipgloss.Width(menuNormal) == 0 {
		t.Error("menu styles should render non-empty content")
	}

	// Test promptStyle
	promptRendered := promptStyle.Render("> ")
	if lipgloss.Width(promptRendered) == 0 {
		t.Error("promptStyle should render non-empty content")
	}
}

// TestProgressBar verifies the ProgressBar function works correctly with lipgloss
func TestProgressBar(t *testing.T) {
	// Test basic progress bar
	result := ProgressBar(50, 100, 20, "█", " ", green, gray)
	if lipgloss.Width(result) == 0 {
		t.Error("ProgressBar should render non-empty content")
	}

	// Test with zero max (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProgressBar panicked with zero max: %v", r)
		}
	}()
	result = ProgressBar(0, 0, 10, "█", " ", green, gray)
	_ = result // Just verify no panic
}

// TestProgressBarWidth verifies that progress bar respects the width parameter
func TestProgressBarWidth(t *testing.T) {
	widths := []int{10, 20, 30, 50}
	for _, width := range widths {
		result := ProgressBar(50, 100, width, "█", " ", green, gray)
		visualWidth := lipgloss.Width(result)
		if visualWidth != width {
			t.Errorf("ProgressBar width should be %d, got %d", width, visualWidth)
		}
	}
}

// TestExitColorCoding verifies exit colors are properly defined
func TestExitColorCoding(t *testing.T) {
	visitedStyle := lipgloss.NewStyle().Foreground(exitVisitedColor)
	knownStyle := lipgloss.NewStyle().Foreground(exitKnownColor)
	newStyle := lipgloss.NewStyle().Foreground(exitNewColor)

	visited := visitedStyle.Render("north")
	known := knownStyle.Render("east")
	new := newStyle.Render("south")

	// Verify styles render correctly
	if lipgloss.Width(visited) == 0 || lipgloss.Width(known) == 0 || lipgloss.Width(new) == 0 {
		t.Error("Exit color styles should render non-empty content")
	}
}

// TestStyledMessage verifies the styledMessage function
func TestStyledMessage(t *testing.T) {
	m := &model{} // Empty model for testing

	// Test different message types
	testCases := []struct {
		msg     string
		msgType string
	}{
		{"Hello world", "info"},
		{"Success!", "success"},
		{"Error occurred", "error"},
		{"", "info"},
	}

	for _, tc := range testCases {
		result := m.styledMessage(tc.msg)
		if tc.msg != "" && lipgloss.Width(result) == 0 {
			t.Errorf("styledMessage(%q, %q) should render non-empty content", tc.msg, tc.msgType)
		}
	}
}

// TestStyleMessage verifies the standalone styleMessage function
func TestStyleMessage(t *testing.T) {
	testCases := []struct {
		msg     string
		msgType string
	}{
		{"Hello", "info"},
		{"Success!", "success"},
		{"Error!", "error"},
		{"Warning", "warning"},
		{"Plain text", ""},
	}

	for _, tc := range testCases {
		result := styleMessage(tc.msg, tc.msgType)
		if tc.msg != "" && lipgloss.Width(result) == 0 {
			t.Errorf("styleMessage(%q, %q) should render non-empty content", tc.msg, tc.msgType)
		}
	}
}

// TestBoxStyleBorder verifies boxStyle has correct border styling
func TestBoxStyleBorder(t *testing.T) {
	content := "Test Content"
	rendered := boxStyle.Render(content)

	// Box should have more than just the content (borders add characters)
	if lipgloss.Width(rendered) <= lipgloss.Width(content) {
		t.Error("boxStyle should add border characters to content")
	}
}

// TestANSIWidthNotCounted verifies that ANSI escape codes don't affect width calculation
// This is critical for the terminal width fix
func TestANSIWidthNotCounted(t *testing.T) {
	// Style some text with lipgloss
	styled := successStyle.Render("Success")

	// The visual width should be the length of "Success" (7)
	// not including any ANSI escape codes
	visualWidth := lipgloss.Width(styled)

	if visualWidth != 7 {
		t.Errorf("Expected visual width 7 for 'Success', got %d", visualWidth)
	}
}

// TestColorVariablesExist verifies color variables are properly defined
func TestColorVariablesExist(t *testing.T) {
	// Verify all color variables can be used to create styles
	colors := []lipgloss.Color{red, green, yellow, blue, purple, white, gray, pink, cyan}

	for i, color := range colors {
		style := lipgloss.NewStyle().Foreground(color)
		rendered := style.Render("Test")
		if lipgloss.Width(rendered) == 0 {
			t.Errorf("Color at index %d should create valid style", i)
		}
	}
}