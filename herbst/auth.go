package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// ============================================================
// AUTH — WELCOME, LOGIN, REGISTER
// ============================================================

func (m *model) handleWelcomeInput(input string) {
	input = strings.ToLower(input)

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

func (m *model) attemptLogin() {
	m.isLoading = true
	m.loadingMessage = "Logging in..."

	jsonData, _ := json.Marshal(map[string]string{
		"email":    m.loginUsername,
		"password": m.loginPassword,
	})

	resp, err := http.Post(RESTAPIBase+"/users/auth", "application/json", bytes.NewBuffer(jsonData))
	m.isLoading = false

	if err != nil {
		m.AppendMessage(fmt.Sprintf("Cannot connect to server at %s. Is the server running?", RESTAPIBase), "error")
		m.AppendMessage("Start the server with: cd server && go run main.go", "info")
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Success - handled below
	case http.StatusUnauthorized:
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok {
			m.AppendMessage(errMsg, "error")
		} else {
			m.AppendMessage("Invalid username or password.", "error")
		}
		m.AppendMessage("Type 'login' to try again or 'register' to create an account.", "info")
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		return
	case http.StatusInternalServerError:
		m.AppendMessage("Server error. Please try again later.", "error")
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		return
	case http.StatusBadRequest:
		m.AppendMessage("Invalid request. Please check your input.", "error")
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		return
	default:
		m.AppendMessage(fmt.Sprintf("Login failed (status %d). Please try again.", resp.StatusCode), "error")
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Login error: %v", err), "error")
		return
	}

	if id, ok := result["id"].(float64); ok {
		m.currentUserID = int(id)
	}
	if email, ok := result["email"].(string); ok {
		m.currentUserName = email
	}
	m.screen = ScreenPlaying
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.AppendMessage(fmt.Sprintf("Welcome back, %s!", m.currentUserName), "success")

	m.loadOrCreateCharacter()
	m.visitedRooms[m.currentRoom] = true

	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), StartingRoomID)
		if err != nil {
			m.err = fmt.Errorf("failed to load starting room: %v", err)
			return
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits
		for dir := range m.exits {
			m.knownExits[dir] = true
		}
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
		email := input
		if email == "" {
			email = m.loginUsername + "@herbstmud.local"
		}
		m.attemptRegistration(email)
	}
}

func (m *model) attemptRegistration(email string) {
	jsonData, _ := json.Marshal(map[string]string{
		"email":    m.loginUsername,
		"password": m.loginPassword,
	})

	resp, err := http.Post(RESTAPIBase+"/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Connection error: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok && (strings.Contains(errMsg, "unique") || strings.Contains(errMsg, "already exists")) {
			m.AppendMessage("Username already taken. Choose a different one.", "error")
			m.inputField = "username"
			m.loginUsername = ""
			m.loginPassword = ""
			m.textInput.EchoMode = textinput.EchoNormal
			return
		}
		m.AppendMessage("Failed to create account. Please try again.", "error")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		m.AppendMessage("Failed to create account. Please try again.", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error processing response: %v", err), "error")
		return
	}

	if id, ok := result["id"].(float64); ok {
		m.currentUserID = int(id)
	}
	if email, ok := result["email"].(string); ok {
		m.currentUserName = email
	}
	m.screen = ScreenPlaying
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.AppendMessage(fmt.Sprintf("Account created! Welcome to Herbst MUD, %s!", m.currentUserName), "success")

	m.loadOrCreateCharacter()
	m.visitedRooms[StartingRoomID] = true

	room, err := m.client.Room.Get(context.Background(), StartingRoomID)
	if err != nil {
		m.err = fmt.Errorf("failed to load starting room: %v", err)
		return
	}
	m.currentRoom = room.ID
	m.roomName = room.Name
	m.roomDesc = room.Description
	m.exits = room.Exits
	for dir := range m.exits {
		m.knownExits[dir] = true
	}
}
