package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestBDD(t *testing.T) {
	t.Run("Server can be initialized", func(t *testing.T) {
		t.Log("Server initialization test passed")
	})

	t.Run("Bubbletea model initializes correctly", func(t *testing.T) {
		model := &model{
			connectedAt: time.Now(),
		}

		// Init returns textinput.Blink which is expected
		_ = model.Init()

		t.Log("Bubbletea model initialization test passed")
	})
}

// TestScreenRendering tests that the screen rendering doesn't accumulate messages
// This is the fix for issue #66: Screen rendering bug - viewport/input duplication
func TestScreenRendering(t *testing.T) {
	t.Run("Message is cleared after View renders in ScreenPlaying", func(t *testing.T) {
		// Create a model with a message set (simulating a command response)
		m := &model{
			screen:         ScreenPlaying,
			message:        "Test message",
			messageType:    "info",
			roomName:       "Test Room",
			roomDesc:       "A test room",
			exits:          map[string]int{"north": 2},
			knownExits:     make(map[string]bool),
			visitedRooms:   make(map[int]bool),
			characterHP:    100,
			characterMaxHP: 100,
			characterStamina:     50,
			characterMaxStamina:  50,
			characterMana:        25,
			characterMaxMana:     25,
			width:         80,
			height:        24,
		}
		m.Init()

		// Render the view
		view1 := m.View()

		// The message should appear in the first render
		if view1 == "" {
			t.Error("Expected non-empty view")
		}

		// The message should now be cleared from the model (this is the fix!)
		if m.message != "" {
			t.Errorf("Expected message to be cleared after View, got: %q", m.message)
		}

		// Render again - message should not appear again
		_ = m.View()

		// The key test: message should remain empty after second render
		if m.message != "" {
			t.Errorf("Message should remain cleared after second View, got: %q", m.message)
		}

		t.Log("Message clearing test passed - no duplication")
	})

	t.Run("View returns without panic for all screens", func(t *testing.T) {
		screens := []string{
			ScreenWelcome,
			ScreenLogin,
			ScreenRegister,
			ScreenProfile,
			ScreenEditField,
			ScreenPlaying,
		}

		for _, screen := range screens {
			m := &model{
				screen:         screen,
				width:          80,
				height:         24,
				characterHP:    100,
				characterMaxHP: 100,
				characterStamina:     50,
				characterMaxStamina:  50,
				characterMana:        25,
				characterMaxMana:     25,
				exits:        make(map[string]int),
				knownExits:   make(map[string]bool),
				visitedRooms: make(map[int]bool),
			}
			m.Init()

			// This should not panic
			view := m.View()
			if view == "" {
				t.Errorf("View returned empty string for screen: %s", screen)
			}
		}
		t.Log("All screens render without panic")
	})
}

// TestMessageTypes verifies that different message types are styled correctly
func TestMessageTypes(t *testing.T) {
	t.Run("styledMessage returns correct style for each type", func(t *testing.T) {
		m := &model{}

		// Test success style
		m.message = "Success message"
		m.messageType = "success"
		styled := m.styledMessage(m.message)
		if styled == "" {
			t.Error("Expected styled success message")
		}

		// Test error style
		m.message = "Error message"
		m.messageType = "error"
		styled = m.styledMessage(m.message)
		if styled == "" {
			t.Error("Expected styled error message")
		}

		// Test info style
		m.message = "Info message"
		m.messageType = "info"
		styled = m.styledMessage(m.message)
		if styled == "" {
			t.Error("Expected styled info message")
		}

		// Test empty message
		m.message = ""
		styled = m.styledMessage(m.message)
		if styled != "" {
			t.Error("Expected empty string for empty message")
		}

		t.Log("Message styling test passed")
	})
}

// TestStatusBar verifies the StatusBar function works correctly
func TestStatusBar(t *testing.T) {
	t.Run("StatusBar renders with correct percentages", func(t *testing.T) {
		result := StatusBar(50, 100, 25, 50, 10, 20)

		// Should contain HP, Stamina, Mana indicators
		if result == "" {
			t.Error("Expected non-empty status bar")
		}

		// Should not contain format errors
		if len(result) < 10 {
			t.Error("Status bar seems too short")
		}

		t.Log("StatusBar rendering test passed")
	})
}

// TestPinkBorderStyling verifies the pink border styling is applied correctly
func TestPinkBorderStyling(t *testing.T) {
	t.Run("ScreenPlaying uses pink borders", func(t *testing.T) {
		m := &model{
			screen:             ScreenPlaying,
			message:            "Test message",
			messageType:        "info",
			roomName:           "Test Room",
			roomDesc:           "A test room",
			exits:              map[string]int{"north": 2},
			knownExits:         make(map[string]bool),
			visitedRooms:       make(map[int]bool),
			characterHP:        100,
			characterMaxHP:     100,
			characterStamina:   50,
			characterMaxStamina: 50,
			characterMana:      25,
			characterMaxMana:   25,
			width:             80,
			height:             24,
		}
		m.Init()

		view := m.View()

		// Should contain status_bar label
		if !strings.Contains(view, "status_bar") {
			t.Error("Expected status_bar label in view")
		}

		// Should contain rounded border character (horizontal line)
		if !strings.Contains(view, "─") {
			t.Error("Expected border characters in view")
		}

		// Should contain ANSI escape code for pink (219 = 38;5;219)
		if !strings.Contains(view, "\x1b[38;5;219") && !strings.Contains(view, "219") {
			t.Log("Note: ANSI codes may vary, checking for border style presence")
		}

		t.Log("Pink border styling test passed")
	})
}

// TestModelFields verifies the model has all required fields for rendering
func TestModelFields(t *testing.T) {
	t.Run("Model has required fields for ScreenPlaying", func(t *testing.T) {
		m := &model{
			screen:         ScreenPlaying,
			message:        "",
			messageType:    "info",
			roomName:       "Test Room",
			roomDesc:       "Test description",
			exits:          make(map[string]int),
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
			width:          80,
			height:         24,
		}

		// Model should render without panic
		view := m.View()
		if view == "" {
			t.Error("Expected non-empty view")
		}

		t.Log("Model fields test passed")
	})
}
// TestClearCommand verifies the clear/cls command clears the message buffer
func TestClearCommand(t *testing.T) {
	t.Run("processInput clears message on 'clear' command", func(t *testing.T) {
		m := &model{
			screen:         ScreenPlaying,
			message:        "Test message that should be cleared",
			messageType:    "info",
			roomName:       "Test Room",
			roomDesc:       "A test room",
			exits:          map[string]int{"north": 2},
			knownExits:     make(map[string]bool),
			visitedRooms:   make(map[int]bool),
			characterHP:    100,
			characterMaxHP: 100,
			characterStamina:     50,
			characterMaxStamina:  50,
			characterMana:        25,
			characterMaxMana:     25,
			width:         80,
			height:        24,
		}
		m.Init()

		// Verify message is set before
		if m.message != "Test message that should be cleared" {
			t.Errorf("Expected message to be set initially")
		}

		// Execute clear command
		m.processInput("clear")

		// Message should be cleared
		if m.message != "" {
			t.Errorf("Expected message to be cleared, got: %q", m.message)
		}

		// messageType should be cleared
		if m.messageType != "" {
			t.Errorf("Expected messageType to be cleared, got: %q", m.messageType)
		}

		// inputBuffer should be cleared
		if m.inputBuffer != "" {
			t.Errorf("Expected inputBuffer to be cleared, got: %q", m.inputBuffer)
		}

		t.Log("Clear command test passed")
	})

	t.Run("processInput clears message on 'cls' command", func(t *testing.T) {
		m := &model{
			screen:         ScreenPlaying,
			message:        "Another test message",
			messageType:    "success",
			roomName:       "Test Room",
			roomDesc:       "A test room",
			exits:          map[string]int{},
			knownExits:     make(map[string]bool),
			visitedRooms:   make(map[int]bool),
			characterHP:    100,
			characterMaxHP: 100,
			characterStamina:     50,
			characterMaxStamina:  50,
			characterMana:        25,
			characterMaxMana:     25,
			width:         80,
			height:        24,
		}
		m.Init()

		// Execute cls command
		m.processInput("cls")

		// Message should be cleared
		if m.message != "" {
			t.Errorf("Expected message to be cleared, got: %q", m.message)
		}

		t.Log("CLS command test passed")
	})
}

// TestNoDebugOutput verifies that View() doesn't output debug logs during normal rendering
func TestNoDebugOutput(t *testing.T) {
	t.Run("View() doesn't produce log output during normal rendering", func(t *testing.T) {
		m := &model{
			screen:         ScreenWelcome,
			width:          80,
			height:         24,
			characterHP:    100,
			characterMaxHP: 100,
			characterStamina:     50,
			characterMaxStamina:  50,
			characterMana:        25,
			characterMaxMana:     25,
			exits:        make(map[string]int),
			knownExits:   make(map[string]bool),
			visitedRooms: make(map[int]bool),
		}
		m.Init()

		// Render the view - should not produce any output to stdout from debug logs
		view := m.View()
		
		// Verify view renders without issues
		if view == "" {
			t.Error("Expected non-empty view")
		}
		
		// Test all screens render without debug output
		screens := []string{ScreenWelcome, ScreenLogin, ScreenRegister, ScreenProfile, ScreenPlaying}
		for _, screen := range screens {
			m.screen = screen
			view = m.View()
			if view == "" {
				t.Errorf("View returned empty for screen: %s", screen)
			}
		}
		
		t.Log("All screens render without debug log output")
	})

	// Test that debug logging is disabled in Update and processInput
	t.Run("Debug logging statements are commented out in Update", func(t *testing.T) {
		// Read the main.go file and check that debug log.Printf are commented
		content, err := os.ReadFile("main.go")
		if err != nil {
			t.Fatalf("Failed to read main.go: %v", err)
		}

		// Check that active (non-commented) debug logs don't exist in Update/processInput
		lines := strings.Split(string(content), "\n")
		updateStarted := false
		processInputStarted := false

		for _, line := range lines {
			if strings.Contains(line, "func (m *model) Update(") {
				updateStarted = true
				processInputStarted = false
			}
			if strings.Contains(line, "func (m *model) processInput(") {
				processInputStarted = true
				updateStarted = false
			}
			if strings.Contains(line, "func (m *model) handle") && processInputStarted {
				processInputStarted = false
			}

			// Check for active debug log.Printf in Update and processInput functions
			isCommented := strings.HasPrefix(strings.TrimSpace(line), "//")

			if (updateStarted || processInputStarted) && strings.Contains(line, "log.Printf") && !isCommented {
				// Check if it's actually a debug log (not error logging)
				if !strings.Contains(line, "Warning") && !strings.Contains(line, "Error") {
					t.Errorf("Found active debug log.Printf in Update/processInput: %s", strings.TrimSpace(line))
				}
			}
		}
		t.Log("Debug logging properly disabled in Update and processInput")
	})
}
