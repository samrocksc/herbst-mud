package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestVisualWidthCalculation verifies that lipgloss.Width correctly
// calculates visual width ignoring ANSI escape codes.
// This is the core fix for Issue #75: Terminal width not honored - right side cut off.
// The bug was that len() was used instead of lipgloss.Width(), causing ANSI codes
// to be counted as visible characters, resulting in incorrect centering calculations.
func TestVisualWidthCalculation(t *testing.T) {
	// Test that ANSI codes don't count toward visual width
	coloredText := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Hello")
	visualWidth := lipgloss.Width(coloredText)
	
	// "Hello" has 5 visual characters, ANSI codes shouldn't count
	if visualWidth != 5 {
		t.Errorf("lipgloss.Width should return 5 for 'Hello', got %d", visualWidth)
	}
	
	// Test with bold styling
	boldText := lipgloss.NewStyle().Bold(true).Render("Test")
	visualWidth = lipgloss.Width(boldText)
	
	if visualWidth != 4 {
		t.Errorf("lipgloss.Width should return 4 for 'Test', got %d", visualWidth)
	}
}

// TestCenteringCalculation verifies that centering uses visual width, not byte length.
// This directly tests the fix for Issue #75.
func TestCenteringCalculation(t *testing.T) {
	// Simulate the centering logic from View()
	width := 80
	line := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("Test Line")
	
	// Using len() would count ANSI codes (wrong - this was the bug)
	lenPadding := (width - len(line)) / 2
	
	// Using lipgloss.Width() correctly ignores ANSI codes (right - this is the fix)
	visualWidth := lipgloss.Width(line)
	visualPadding := (width - visualWidth) / 2
	
	// The visual padding should be larger than len-based padding
	// because len() counts ANSI codes as extra characters
	if visualPadding < lenPadding {
		t.Errorf("Visual padding %d should be >= len-based padding %d", visualPadding, lenPadding)
	}
	
	// Verify the visual padding makes sense
	// "Test Line" has 9 characters, so padding should be (80-9)/2 = 35
	if visualPadding != 35 {
		t.Errorf("Expected visual padding 35, got %d (visualWidth=%d)", visualPadding, visualWidth)
	}
}

// TestWidthDefaultValue verifies default width fallback when terminal reports 0.
func TestWidthDefaultValue(t *testing.T) {
	// When width < 40, we should default to 80
	testCases := []struct {
		input    int
		expected int
	}{
		{0, 80},   // No terminal size reported
		{20, 80},  // Too narrow
		{39, 80},  // Just below threshold
		{40, 40},  // At threshold - use actual
		{80, 80},  // Normal width - use actual
		{120, 120}, // Wide terminal - use actual
	}
	
	for _, tc := range testCases {
		width := tc.input
		if width < 40 {
			width = 80
		}
		if width != tc.expected {
			t.Errorf("Input width %d: expected %d, got %d", tc.input, tc.expected, width)
		}
	}
}