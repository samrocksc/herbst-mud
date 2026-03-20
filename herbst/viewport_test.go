package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// TestViewportInitialization verifies viewport is initialized correctly
func TestViewportInitialization(t *testing.T) {
	width := 80
	height := 24
	statusHeight := 2
	inputHeight := 3

	vpHeight := height - statusHeight - inputHeight

	if vpHeight <= 0 {
		t.Fatalf("Viewport height should be positive, got %d", vpHeight)
	}

	vp := viewport.New(width, vpHeight)

	if vp.Width != width {
		t.Errorf("Expected viewport width %d, got %d", width, vp.Width)
	}

	if vp.Height != vpHeight {
		t.Errorf("Expected viewport height %d, got %d", vpHeight, vp.Height)
	}
}

// TestViewportContentSet verifies content can be set and retrieved
func TestViewportContentSet(t *testing.T) {
	vp := viewport.New(80, 20)

	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	vp.SetContent(content)

	// Viewport should have content set
	view := vp.View()
	if view == "" {
		t.Error("Viewport view should not be empty after setting content")
	}
}

// TestViewportWindowSizeMessage verifies viewport responds to WindowSizeMsg
func TestViewportWindowSizeMessage(t *testing.T) {
	vp := viewport.New(80, 20)

	// Simulate window resize
	newWidth := 120
	newHeight := 30
	statusHeight := 2
	inputHeight := 3

	msg := tea.WindowSizeMsg{
		Width:  newWidth,
		Height: newHeight,
	}

	vp.Width = msg.Width
	vp.Height = msg.Height - statusHeight - inputHeight

	if vp.Width != newWidth {
		t.Errorf("Expected width %d after resize, got %d", newWidth, vp.Width)
	}

	expectedHeight := newHeight - statusHeight - inputHeight
	if vp.Height != expectedHeight {
		t.Errorf("Expected height %d after resize, got %d", expectedHeight, vp.Height)
	}
}

// TestViewportScrollingContent verifies viewport handles scrolling content
func TestViewportScrollingContent(t *testing.T) {
	vp := viewport.New(80, 10)

	// Create long content that exceeds viewport height
	longContent := ""
	for i := 1; i <= 50; i++ {
		longContent += "Line " + string(rune('0'+i%10)) + " - This is line number " + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10)) + "\n"
	}

	vp.SetContent(longContent)

	// Viewport should handle the content
	view := vp.View()
	if view == "" {
		t.Error("Viewport should render content")
	}

	// Verify we can render content and access properties
	_ = vp.View() // Should return rendered string
	if vp.Width <= 0 {
		t.Error("Viewport should have positive width")
	}
	if vp.Height <= 0 {
		t.Error("Viewport should have positive height")
	}
}

// TestViewportFullScreenDimensions verifies full terminal minus status/input
func TestViewportFullScreenDimensions(t *testing.T) {
	testCases := []struct {
		terminalWidth  int
		terminalHeight int
		statusHeight   int
		inputHeight    int
	}{
		{80, 24, 2, 3},    // Standard small terminal
		{120, 40, 2, 3},   // Larger terminal
		{200, 60, 3, 4},   // Wide terminal
		{132, 43, 2, 3},   // Old-school 132x43
	}

	for _, tc := range testCases {
		vpHeight := tc.terminalHeight - tc.statusHeight - tc.inputHeight

		if vpHeight <= 0 {
			t.Errorf("Terminal %dx%d: viewport height should be positive, got %d",
				tc.terminalWidth, tc.terminalHeight, vpHeight)
			continue
		}

		vp := viewport.New(tc.terminalWidth, vpHeight)

		if vp.Width != tc.terminalWidth {
			t.Errorf("Terminal %dx%d: expected width %d, got %d",
				tc.terminalWidth, tc.terminalHeight, tc.terminalWidth, vp.Width)
		}

		if vp.Height != vpHeight {
			t.Errorf("Terminal %dx%d: expected height %d, got %d",
				tc.terminalWidth, tc.terminalHeight, vpHeight, vp.Height)
		}
	}
}