package main

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel(screen string) *model {
	ti := textinput.New()
	ti.Focus()
	m := &model{
		screen:       screen,
		textInput:    ti,
		connectedAt:  time.Now(),
		visitedRooms: make(map[int]bool),
		knownExits:   make(map[string]bool),
		width:        80,
		height:       24,
		maxHistory:   50,
	}
	m.Init()
	return m
}

// TestWelcomeScreenNavigation verifies navigation from welcome screen
func TestWelcomeScreenNavigation(t *testing.T) {
	t.Run("typing 1 goes to login screen", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.processInput("1")
		if m.screen != ScreenLogin {
			t.Errorf("Expected screen %q, got %q", ScreenLogin, m.screen)
		}
		if m.inputField != "username" {
			t.Errorf("Expected inputField 'username', got %q", m.inputField)
		}
	})

	t.Run("typing 2 goes to register screen", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.processInput("2")
		if m.screen != ScreenRegister {
			t.Errorf("Expected screen %q, got %q", ScreenRegister, m.screen)
		}
		if m.inputField != "username" {
			t.Errorf("Expected inputField 'username', got %q", m.inputField)
		}
	})

	t.Run("typing 3 quits", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.processInput("3")
		if len(m.messageHistory) == 0 {
			t.Error("Expected message in history")
		}
	})

	t.Run("invalid choice shows error", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.processInput("9")
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
	})
}

// TestLoginFlow tests the login input state machine
func TestLoginFlow(t *testing.T) {
	t.Run("enter username transitions to password", func(t *testing.T) {
		m := newTestModel(ScreenLogin)
		m.inputField = "username"
		m.handleLoginInput("player@example.com")
		if m.inputField != "password" {
			t.Errorf("Expected inputField 'password', got %q", m.inputField)
		}
		if m.loginUsername != "player@example.com" {
			t.Errorf("Expected loginUsername 'player@example.com', got %q", m.loginUsername)
		}
	})

	t.Run("empty username transitions to password", func(t *testing.T) {
		m := newTestModel(ScreenLogin)
		m.inputField = "username"
		m.handleLoginInput("")
		if m.inputField != "password" {
			t.Errorf("Expected inputField 'password', got %q", m.inputField)
		}
	})
}

// TestNavigationKeys tests the j/k menu navigation on welcome screen
func TestNavigationKeys(t *testing.T) {
	t.Run("pressing j moves cursor down", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if m.menuCursor != 1 {
			t.Errorf("Expected menuCursor 1 after j, got %d", m.menuCursor)
		}
	})

	t.Run("pressing k moves cursor up", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 2
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		if m.menuCursor != 1 {
			t.Errorf("Expected menuCursor 1 after k, got %d", m.menuCursor)
		}
	})

	t.Run("j wraps cursor to first at bottom", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 2
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if m.menuCursor != 0 {
			t.Errorf("Expected menuCursor wrap to 0, got %d", m.menuCursor)
		}
	})

	t.Run("k wraps cursor to last at top", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		if m.menuCursor != 2 {
			t.Errorf("Expected menuCursor wrap to 2, got %d", m.menuCursor)
		}
	})

	t.Run("downarrow works same as j", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
		m.Update(tea.KeyMsg{Type: tea.KeyDown})
		if m.menuCursor != 1 {
			t.Errorf("Expected menuCursor 1 after down, got %d", m.menuCursor)
		}
	})

	t.Run("uparrow works same as k", func(t *testing.T) {
		m := newTestModel(ScreenWelcome)
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 1
		m.Update(tea.KeyMsg{Type: tea.KeyUp})
		if m.menuCursor != 0 {
			t.Errorf("Expected menuCursor 0 after up, got %d", m.menuCursor)
		}
	})
}

// TestEscapeKey tests that escape returns to welcome screen from login/register
func TestEscapeKey(t *testing.T) {
	t.Run("escape from login goes to welcome", func(t *testing.T) {
		m := newTestModel(ScreenLogin)
		m.loginUsername = "player@example.com"
		m.loginPassword = "password123"
		m.handleEscape()
		if m.screen != ScreenWelcome {
			t.Errorf("Expected screen %q after escape, got %q", ScreenWelcome, m.screen)
		}
		if m.loginUsername != "" {
			t.Errorf("Expected loginUsername cleared, got %q", m.loginUsername)
		}
		if m.loginPassword != "" {
			t.Errorf("Expected loginPassword cleared, got %q", m.loginPassword)
		}
	})
}

// TestScreenStateTransitions tests the full login flow state machine
func TestScreenStateTransitions(t *testing.T) {
	t.Run("register transitions through all fields", func(t *testing.T) {
		m := newTestModel(ScreenRegister)
		m.inputField = "username"

		m.handleRegisterInput("newuser")
		if m.inputField != "password" {
			t.Errorf("Expected password after username, got %q", m.inputField)
		}

		m.handleRegisterInput("secret123")
		if m.inputField != "confirm_password" {
			t.Errorf("Expected confirm_password after password, got %q", m.inputField)
		}

		m.handleRegisterInput("secret123")
		if m.inputField != "email" {
			t.Errorf("Expected email after confirm, got %q", m.inputField)
		}
	})

	t.Run("register password mismatch resets to password", func(t *testing.T) {
		m := newTestModel(ScreenRegister)
		m.inputField = "confirm_password"
		m.loginPassword = "secret123"

		m.handleRegisterInput("wrongpassword")

		if m.inputField != "password" {
			t.Errorf("Expected back to password on mismatch, got %q", m.inputField)
		}
		if len(m.messageTypes) == 0 {
			t.Fatal("Expected message in history")
		}
		if lastType := m.messageTypes[len(m.messageTypes)-1]; lastType != "error" {
			t.Errorf("Expected error type, got %q", lastType)
		}
	})
}
