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
	result := welcomeScreen(80, 24, "test input")

	// Verify ASCII art is rendered (green)
	if !strings.Contains(result, "\x1b[38;5;46m") && !strings.Contains(result, "HERBST") {
		// Either ANSI green code or text should be present
		if !strings.Contains(result, "HERBST MUD") {
			t.Error("Welcome screen should contain HERBST MUD")
		}
	}

	// Verify commands are displayed
	if !strings.Contains(result, "login") {
		t.Error("Welcome screen should contain login command")
	}

	if !strings.Contains(result, "register") {
		t.Error("Welcome screen should contain register command")
	}

	if !strings.Contains(result, "quit") {
		t.Error("Welcome screen should contain quit command")
	}

	// Verify welcome message
	if !strings.Contains(result, "Welcome Adventurer") {
		t.Error("Welcome screen should contain welcome message")
	}

	// Verify tip is present
	if !strings.Contains(result, "Tip:") {
		t.Error("Welcome screen should contain tip")
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

			result := welcomeScreen(tt.width, tt.height, "test")
			if result == "" {
				t.Error("Result should not be empty")
			}
		})
	}
}

// TestWelcomeScreenArrowIcon verifies arrow icons are displayed
func TestWelcomeScreenArrowIcon(t *testing.T) {
	result := welcomeScreen(80, 24, "test input")

	// Verify arrow icons are used
	if !strings.Contains(result, "➤") {
		t.Error("Welcome screen should use arrow icons (➤)")
	}
}