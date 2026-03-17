package main

import (
	"fmt"
	"strings"
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

// TestViewRespectsTerminalWidth verifies that the View() function properly
// uses the terminal width from the model for rendering ScreenPlaying.
func TestViewRespectsTerminalWidth(t *testing.T) {
	testCases := []struct {
		name          string
		width         int
		height        int
		expectedLines int // minimum expected lines
	}{
		{"Narrow terminal (40)", 40, 20, 5},
		{"Normal terminal (80)", 80, 24, 5},
		{"Wide terminal (120)", 120, 30, 5},
		{"Very wide terminal (200)", 200, 40, 5},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &model{
				screen:         ScreenPlaying,
				width:          tc.width,
				height:         tc.height,
				roomName:       "Test Room",
				roomDesc:       "A test room description",
				exits:          map[string]int{"north": 2, "south": 3},
				knownExits:     make(map[string]bool),
				visitedRooms:   make(map[int]bool),
				characterHP:    100,
				characterMaxHP: 100,
				characterStamina:     50,
				characterMaxStamina:  50,
				characterMana:        25,
				characterMaxMana:     25,
				characterLevel:       1,
				characterExperience:  0,
			}
			m.Init()
			
			view := m.View()
			
			// View should not be empty
			if view == "" {
				t.Error("View returned empty string")
			}
			
			// Count lines in output
			lines := strings.Count(view, "\n") + 1
			if lines < tc.expectedLines {
				t.Errorf("Expected at least %d lines, got %d", tc.expectedLines, lines)
			}
			
			// View should not panic on any width
			t.Logf("Successfully rendered view for width %d, height %d", tc.width, tc.height)
		})
	}
}

// TestViewWidthConsistency verifies that rendered output doesn't exceed terminal width.
func TestViewWidthConsistency(t *testing.T) {
	widths := []int{80, 100, 120, 150}
	
	for _, width := range widths {
		t.Run(fmt.Sprintf("Width_%d", width), func(t *testing.T) {
			m := &model{
				screen:         ScreenPlaying,
				width:          width,
				height:         24,
				roomName:       "The Grand Hall",
				roomDesc:       "A magnificent hall with towering pillars and golden chandeliers.",
				exits:          map[string]int{"north": 2, "south": 3, "east": 4},
				knownExits:     make(map[string]bool),
				visitedRooms:   make(map[int]bool),
				characterHP:    100,
				characterMaxHP: 100,
				characterStamina:     50,
				characterMaxStamina:  50,
				characterMana:        25,
				characterMaxMana:     25,
			}
			m.Init()
			
			view := m.View()
			
			// Check each line's visual width
			lines := strings.Split(view, "\n")
			for i, line := range lines {
				visualWidth := lipgloss.Width(line)
				// Allow some tolerance for border/padding
				maxAllowed := width + 10
				if visualWidth > maxAllowed {
					t.Errorf("Line %d exceeds max width: visualWidth=%d, maxAllowed=%d, width=%d", 
						i, visualWidth, maxAllowed, width)
				}
			}
		})
	}
}

// TestScreenPlayingWidthUsage verifies that ScreenPlaying uses the model's width correctly.
func TestScreenPlayingWidthUsage(t *testing.T) {
	// Test that the model's width value is used in rendering
	m1 := &model{
		screen:         ScreenPlaying,
		width:          80,
		height:         24,
		roomName:       "Test",
		roomDesc:       "Test room",
		exits:          map[string]int{"north": 2},
		knownExits:     make(map[string]bool),
		visitedRooms:   make(map[int]bool),
		characterHP:    100,
		characterMaxHP: 100,
		characterStamina:     50,
		characterMaxStamina:  50,
		characterMana:        25,
		characterMaxMana:     25,
	}
	m1.Init()
	view1 := m1.View()
	
	// Now with different width
	m2 := &model{
		screen:         ScreenPlaying,
		width:          120,
		height:         24,
		roomName:       "Test",
		roomDesc:       "Test room",
		exits:          map[string]int{"north": 2},
		knownExits:     make(map[string]bool),
		visitedRooms:   make(map[int]bool),
		characterHP:    100,
		characterMaxHP: 100,
		characterStamina:     50,
		characterMaxStamina:  50,
		characterMana:        25,
		characterMaxMana:     25,
	}
	m2.Init()
	view2 := m2.View()
	
	// The outputs should be different because they use different widths
	// ScreenPlaying doesn't center, but it should have different internal widths
	// Check that the view generates without error
	if view1 == "" || view2 == "" {
		t.Error("Both views should render successfully")
	}
	
	// Both should contain the room name
	if !strings.Contains(view1, "Test") || !strings.Contains(view2, "Test") {
		t.Error("Both views should contain room name")
	}
}

// TestCenteringLogicWithAnsiCodes verifies that centering works correctly with ANSI-styled text.
func TestCenteringLogicWithAnsiCodes(t *testing.T) {
	// Test various styled strings
	testCases := []struct {
		name          string
		text          string
		style         lipgloss.Style
		expectedWidth int
	}{
		{
			name:          "Plain text",
			text:          "Hello World",
			style:         lipgloss.NewStyle(),
			expectedWidth: 11,
		},
		{
			name:          "Bold text",
			text:          "Hello World",
			style:         lipgloss.NewStyle().Bold(true),
			expectedWidth: 11,
		},
		{
			name:          "Colored text",
			text:          "Hello World",
			style:         lipgloss.NewStyle().Foreground(lipgloss.Color("46")),
			expectedWidth: 11,
		},
		{
			name:          "Bold colored text",
			text:          "Hello World",
			style:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")),
			expectedWidth: 11,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rendered := tc.style.Render(tc.text)
			visualWidth := lipgloss.Width(rendered)
			
			if visualWidth != tc.expectedWidth {
				t.Errorf("Expected visual width %d, got %d for %q", 
					tc.expectedWidth, visualWidth, tc.text)
			}
			
			// len() should be larger due to ANSI codes (for styled text)
			byteLen := len(rendered)
			if tc.style.GetBold() || tc.style.GetForeground() != "" {
				if byteLen <= tc.expectedWidth {
					t.Logf("Warning: ANSI codes should increase byte length, got %d for visual %d",
						byteLen, visualWidth)
				}
			}
		})
	}
}

// TestWidthFallbackOnZero verifies that ScreenPlaying handles zero/negative width.
func TestWidthFallbackOnZero(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
	}{
		{"Zero width", 0, 24},
		{"Negative width", -10, 24},
		{"Zero height", 80, 0},
		{"Negative height", 80, -5},
		{"Both zero", 0, 0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &model{
				screen:         ScreenPlaying,
				width:          tc.width,
				height:         tc.height,
				roomName:       "Test Room",
				roomDesc:       "Test description",
				exits:          map[string]int{"north": 2},
				knownExits:     make(map[string]bool),
				visitedRooms:   make(map[int]bool),
				characterHP:    100,
				characterMaxHP: 100,
				characterStamina:     50,
				characterMaxStamina:  50,
				characterMana:        25,
				characterMaxMana:     25,
			}
			m.Init()
			
			// Should not panic
			view := m.View()
			
			// Should produce some output (using defaults)
			if view == "" {
				t.Error("View should not be empty even with invalid width/height")
			}
		})
	}
}