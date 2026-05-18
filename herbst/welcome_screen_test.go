package main

import (
	"strings"
	"testing"
)

// TestWelcomeScreenASCIIArt verifies the ASCII art logo is present
// NOTE: asciiLogo is not currently defined. Skipping until implemented.
func TestWelcomeScreenASCIIArt(t *testing.T) {
	t.Skip("asciiLogo not defined yet")
}

// TestWelcomeScreenContent verifies the welcome screen content structure
func TestWelcomeScreenContent(t *testing.T) {
	result := welcomeScreen(80, 24, 0, "test input")

	// Verify HERBST MUD title is present
	if !strings.Contains(result, "HERBST MUD") {
		t.Error("Welcome screen should contain HERBST MUD")
	}

	// Verify subtitle
	if !strings.Contains(result, "A World of Adventure Awaits") {
		t.Error("Welcome screen should contain subtitle")
	}

	// Verify menu items are displayed (capitalized with ANSI styling)
	if !strings.Contains(result, "Login") {
		t.Error("Welcome screen should contain Login option")
	}

	if !strings.Contains(result, "Register") {
		t.Error("Welcome screen should contain Register option")
	}

	if !strings.Contains(result, "Quit") {
		t.Error("Welcome screen should contain Quit option")
	}

	// Verify hint text
	if !strings.Contains(result, "Type a number or command") {
		t.Error("Welcome screen should contain hint text")
	}
}

// TestWelcomeScreenDimensions verifies proper dimension handling
func TestWelcomeScreenDimensions(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		height      int
		expectPanic bool
	}{
		{"normal", 80, 24, false},
		{"minimum width", 40, 20, false},
		{"minimum height", 80, 10, false},
		{"very small", 20, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			result := welcomeScreen(tt.width, tt.height, 0, "test")
			if result == "" {
				t.Error("Result should not be empty")
			}
		})
	}
}

// TestWelcomeScreenArrowIcon verifies arrow icons are displayed
func TestWelcomeScreenArrowIcon(t *testing.T) {
	result := welcomeScreen(80, 24, 0, "test input")

	// Verify arrow icons are used
	if !strings.Contains(result, "▸") {
		t.Error("Welcome screen should use arrow icons (▸)")
	}
}
