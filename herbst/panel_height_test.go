package main

import (
	"strings"
	"testing"
)

// TestPanelHeightsFillTerminal verifies that panels expand to fill terminal height
// instead of collapsing to content height (Height(0) bug fix)
// See: tickets/ui-20-fix-panels-not-full-height.md
func TestPanelHeightsFillTerminal(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		height        int
		minLineCount  int // minimum lines expected in output
		screenFunc    func(int, int, string) string
	}{
		{"welcomeScreen 80x24", 80, 24, 15, func(w, h int, s string) string { return welcomeScreen(w, h, s) }},
		{"welcomeScreen 120x40", 120, 40, 25, func(w, h int, s string) string { return welcomeScreen(w, h, s) }},
		{"loginScreen 80x24", 80, 24, 15, func(w, h int, s string) string { return loginScreen(w, h, "", "", s) }},
		{"loginScreen 120x40", 120, 40, 25, func(w, h int, s string) string { return loginScreen(w, h, "", "", s) }},
		{"registerScreen 80x24", 80, 24, 15, func(w, h int, s string) string { return registerScreen(w, h, "", "", s) }},
		{"registerScreen 120x40", 120, 40, 25, func(w, h int, s string) string { return registerScreen(w, h, "", "", s) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.screenFunc(tt.width, tt.height, "test input")
			lines := strings.Split(result, "\n")

			// The output should have significantly more lines than just the content
			// (which is about 10 lines) - it should fill most of the terminal
			if len(lines) < tt.minLineCount {
				t.Errorf("Panel should fill terminal height. Got %d lines, want >= %d", len(lines), tt.minLineCount)
				t.Logf("Full output:\n%s", result)
			}

			// Verify we have multiple border lines (top and bottom of each panel)
			borderCount := strings.Count(result, "─")
			if borderCount < 4 {
				t.Errorf("Should have multiple panel borders. Got %d border chars, want >= 4", borderCount)
			}
		})
	}
}

// TestExplicitHeightCalculations verifies height calculation logic
func TestExplicitHeightCalculations(t *testing.T) {
	tests := []struct {
		name         string
		terminalH    int
		wantMinLines int
		fn           string
	}{
		{"welcome 24h", 24, 18, "welcomeScreen"},
		{"welcome 40h", 40, 32, "welcomeScreen"},
		{"login 24h", 24, 18, "loginScreen"},
		{"login 40h", 40, 32, "loginScreen"},
		{"register 24h", 24, 18, "registerScreen"},
		{"register 40h", 40, 32, "registerScreen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			switch tt.fn {
			case "welcomeScreen":
				result = welcomeScreen(80, tt.terminalH, "input")
			case "loginScreen":
				result = loginScreen(80, tt.terminalH, "", "", "input")
			case "registerScreen":
				result = registerScreen(80, tt.terminalH, "", "", "input")
			}

			lines := strings.Split(result, "\n")
			totalPanelHeight := len(lines)
			
			// Panels should be within reasonable range of terminal height
			if totalPanelHeight < tt.wantMinLines {
				t.Errorf("%s: panels should fill most of terminal. Got %d lines, want >= %d", 
					tt.fn, totalPanelHeight, tt.wantMinLines)
			}
		})
	}
}
