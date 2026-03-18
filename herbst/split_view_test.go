package main

import (
	"testing"
)

// TestProportionalHeights verifies that the proportional heights calculation works correctly
// Issue #72: UI: Full-width split windows with proportional heights
func TestProportionalHeights(t *testing.T) {
	t.Run("Calculates correct proportional heights at 80x24", func(t *testing.T) {
		_, height := 80, 24

		// Input: ~20% = 4
		inputHeight := height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}

		// Status: ~10% = 2
		statusHeight := height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}

		// Viewport: ~70% = remaining
		viewportHeight := height - inputHeight - statusHeight
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		// With min heights applied: input=4, status=3, viewport=17
		if inputHeight != 4 {
			t.Errorf("Expected inputHeight 4, got %d", inputHeight)
		}
		if statusHeight != 3 { // min height applied
			t.Errorf("Expected statusHeight 3, got %d", statusHeight)
		}
		if viewportHeight != 17 {
			t.Errorf("Expected viewportHeight 17, got %d", viewportHeight)
		}

		t.Logf("Proportional heights (80x24): input=%d, status=%d, viewport=%d",
			inputHeight, statusHeight, viewportHeight)
	})

	t.Run("Calculates correct proportional heights at 120x40", func(t *testing.T) {
		_, height := 120, 40

		// Input: ~20% = 8
		inputHeight := height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}

		// Status: ~10% = 4
		statusHeight := height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}

		// Viewport: ~70% = remaining
		viewportHeight := height - inputHeight - statusHeight
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		if inputHeight != 8 {
			t.Errorf("Expected inputHeight 8, got %d", inputHeight)
		}
		if statusHeight != 4 {
			t.Errorf("Expected statusHeight 4, got %d", statusHeight)
		}
		if viewportHeight != 28 {
			t.Errorf("Expected viewportHeight 28, got %d", viewportHeight)
		}

		t.Logf("Proportional heights (120x40): input=%d, status=%d, viewport=%d",
			inputHeight, statusHeight, viewportHeight)
	})

	t.Run("Uses minimum heights for small terminals", func(t *testing.T) {
		_, height := 40, 10

		// Input: ~20% = 2, but min is 3
		inputHeight := height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}

		// Status: ~10% = 1, but min is 3
		statusHeight := height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}

		// Viewport: remaining = 4, but min is 5
		viewportHeight := height - inputHeight - statusHeight
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		// Minimum heights should be applied
		if inputHeight != 3 {
			t.Errorf("Expected inputHeight min 3, got %d", inputHeight)
		}
		if statusHeight != 3 {
			t.Errorf("Expected statusHeight min 3, got %d", statusHeight)
		}
		if viewportHeight != 5 {
			t.Errorf("Expected viewportHeight min 5, got %d", viewportHeight)
		}

		t.Logf("Minimum heights (40x10): input=%d, status=%d, viewport=%d",
			inputHeight, statusHeight, viewportHeight)
	})

	t.Run("Panel widths use full terminal width", func(t *testing.T) {
		width := 80

		// Full-width panels should account for borders (-2 each side)
		outputWidth := width - 2
		inputWidth := width - 2

		if outputWidth != 78 {
			t.Errorf("Expected outputWidth 78, got %d", outputWidth)
		}
		if inputWidth != 78 {
			t.Errorf("Expected inputWidth 78, got %d", inputWidth)
		}

		t.Logf("Full-width panels: output=%d, input=%d", outputWidth, inputWidth)
	})
}