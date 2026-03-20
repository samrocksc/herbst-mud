package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// ============================================================
// INPUT PROCESSING
// ============================================================

func (m *model) processInput(input string) {
	// Debug: log.Printf("processInput called with: %q, screen: %s", input, m.screen)
	input = strings.TrimSpace(input)

	switch m.screen {
	case ScreenWelcome:
		m.handleWelcomeInput(input)
	case ScreenLogin:
		m.handleLoginInput(input)
	case ScreenRegister:
		m.handleRegisterInput(input)
	case ScreenProfile:
		m.handleProfileInput(input)
	case ScreenEditField:
		m.handleEditFieldInput(input)
	case ScreenPlaying:
		m.processCommand(input)
	}
}

func (m *model) handleWelcomeInput(input string) {
	// Debug: log.Printf("handleWelcomeInput called with: %q", input)
	input = strings.ToLower(input)

	// Vim-style selection with numbers or j/k navigation
	switch input {
	case "1", "login":
		m.screen = ScreenLogin
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.AppendMessage("Enter your username:", "info")
		m.textInput.Focus()
	case "2", "register", "r", "create":
		m.screen = ScreenRegister
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.AppendMessage("Choose a username:", "info")
		m.textInput.Focus()
	case "3", "quit", "q":
		m.AppendMessage("Goodbye! Thanks for playing Herbst MUD.", "success")
		m.inputBuffer = ""
		return
	default:
		if input != "" {
			m.AppendMessage("Invalid choice. Type 1, 2, or 3", "error")
		}
	}
}

func (m *model) handleLoginInput(input string) {
	if m.inputField == "username" {
		m.loginUsername = input
		m.inputField = "password"
		m.AppendMessage("Enter your password:", "info")
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.Focus()
	} else if m.inputField == "password" {
		m.loginPassword = input
		m.textInput.EchoMode = textinput.EchoNormal
		m.attemptLogin()
	}
}

func (m *model) handleRegisterInput(input string) {
	if m.inputField == "username" {
		if input == "" {
			m.AppendMessage("Username cannot be empty. Try again:", "error")
			return
		}
		m.loginUsername = input
		m.inputField = "password"
		m.AppendMessage("Choose a password:", "info")
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.Focus()
	} else if m.inputField == "password" {
		if input == "" {
			m.AppendMessage("Password cannot be empty. Try again:", "error")
			return
		}
		m.loginPassword = input
		m.inputField = "confirm_password"
		m.AppendMessage("Confirm your password:", "info")
		m.textInput.Focus()
	} else if m.inputField == "confirm_password" {
		if input != m.loginPassword {
			m.AppendMessage("Passwords do not match. Try again:", "error")
			m.inputField = "password"
			m.loginPassword = ""
			m.textInput.EchoMode = textinput.EchoPassword
			m.textInput.Focus()
			return
		}
		m.inputField = "email"
		m.AppendMessage("Enter your email (optional, press enter to skip):", "info")
		m.textInput.EchoMode = textinput.EchoNormal
		m.textInput.Focus()
	} else if m.inputField == "email" {
		// Email is optional - use username if not provided
		email := input
		if email == "" {
			email = m.loginUsername + "@herbstmud.local"
		}
		m.attemptRegistration(email)
	}
}

func (m *model) handleProfileInput(input string) {
	input = strings.ToLower(input)
	switch input {
	case "1":
		m.editField = "gender"
		m.screen = ScreenEditField
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = ""
	case "2":
		m.editField = "description"
		m.screen = ScreenEditField
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = ""
	case "3", "back", "b", "esc":
		m.screen = ScreenPlaying
		m.message = ""
		m.messageType = ""
		m.menuItems = []string{}
		// Vim-style selection for welcome screen
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	default:
		m.message = "Invalid choice. Enter 1, 2, or 3"
		m.messageType = "error"
	}
}

func (m *model) handleEditFieldInput(input string) {
	if m.editField == "gender" {
		m.characterGender = input
		m.saveProfileToDB()
		m.message = "Gender updated!"
		m.messageType = "success"
	} else if m.editField == "description" {
		m.characterDescription = input
		m.saveProfileToDB()
		m.message = "Description updated!"
		m.messageType = "success"
	}
	m.screen = ScreenProfile
	m.textInput.SetValue("")
	m.inputBuffer = ""
}

// fuzzyWordMatch returns true if all words in target appear as substrings in name.
// "grand man" matches "Grand Ol' Man". "man" also matches. Case-insensitive.
func fuzzyWordMatch(name, target string) bool {
	nameLower := strings.ToLower(name)
	for _, word := range strings.Fields(strings.ToLower(target)) {
		if !strings.Contains(nameLower, word) {
			return false
		}
	}
	return true
}
