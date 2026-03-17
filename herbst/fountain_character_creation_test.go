package main

import (
	"testing"
)

// TestFountainScreensExist verifies the fountain character creation screens are defined
func TestFountainScreensExist(t *testing.T) {
	// Test that the screens are accessible
	_ = fountainWakeScreen()
	_ = fountainWashScreen()
	_ = characterCreateScreen()
}

// TestFountainWakeInput tests the wake screen input handler
func TestFountainWakeInput(t *testing.T) {
	m := &model{
		screen: ScreenFountainWake,
	}

	// Should transition to wash screen
	m.handleFountainWakeInput("anything")

	if m.screen != ScreenFountainWash {
		t.Errorf("Expected screen to be %s, got %s", ScreenFountainWash, m.screen)
	}
}

// TestFountainWashInput tests the wash screen input handler
func TestFountainWashInput(t *testing.T) {
	m := &model{
		screen: ScreenFountainWash,
	}

	// Should transition to character creation
	m.handleFountainWashInput("anything")

	if m.screen != ScreenCharacterCreate {
		t.Errorf("Expected screen to be %s, got %s", ScreenCharacterCreate, m.screen)
	}
}

// TestCharacterCreateInput tests various character creation inputs
func TestCharacterCreateInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"name option", "1", "Enter your character name:"},
		{"race option", "2", "Select your race:"},
		{"gender option", "3", "Select your gender:"},
		{"class option", "4", "Select your class:"},
		{"size option", "5", "Select your size:"},
		{"done", "done", "Character creation complete!"},
		{"invalid", "invalid", "Invalid choice"},
		{"lowercase name", "name", "Enter your character name:"},
		{"lowercase race", "race", "Select your race:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{
				screen:       ScreenCharacterCreate,
				textInput:    textinput.New(),
			}
			m.handleCharacterCreateInput(tt.input)

			if m.message != tt.expected {
				t.Errorf("Expected message %q, got %q", tt.expected, m.message)
			}
		})
	}
}

// TestScreenConstants verifies all screen constants are defined
func TestScreenConstants(t *testing.T) {
	screens := []string{
		ScreenWelcome,
		ScreenLogin,
		ScreenRegister,
		ScreenPlaying,
		ScreenProfile,
		ScreenEditField,
		ScreenFountainWake,
		ScreenFountainWash,
		ScreenCharacterCreate,
	}

	expected := []string{
		"welcome",
		"login",
		"register",
		"playing",
		"profile",
		"edit_field",
		"fountain_wake",
		"fountain_wash",
		"character_create",
	}

	for i, screen := range screens {
		if screen != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], screen)
		}
	}
}