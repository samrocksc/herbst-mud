package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
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
	case "4", "world", "w":
		m.screen = ScreenWorldSelect
		m.fetchWorlds()
		m.AppendMessage(m.displayWorlds(), "info")
	default:
		if input != "" {
			m.AppendMessage("Invalid choice. Type 1, 2, 3, or 4", "error")
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
	if token, ok := result["token"].(string); ok {
		m.characterToken = token
	}
	// Go to world select instead of auto-loading character
	m.screen = ScreenWorldSelect
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.AppendMessage(fmt.Sprintf("Welcome back, %s!", m.currentUserName), "success")
	m.fetchWorlds()
	m.AppendMessage(m.displayWorlds(), "info")
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
	// Go to world select instead of auto-loading character
	m.screen = ScreenWorldSelect
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.AppendMessage(fmt.Sprintf("Account created! Welcome to Herbst MUD, %s!", m.currentUserName), "success")
	m.fetchWorlds()
	m.AppendMessage(m.displayWorlds(), "info")
}

// worlds holds the list of available worlds
var availableWorlds []string

// fetchWorlds retrieves the list of available worlds from the server
func (m *model) fetchWorlds() {
	m.isLoading = true
	m.loadingMessage = "Fetching worlds..."

	resp, err := http.Get(RESTAPIBase + "/admin/export/worlds")
	m.isLoading = false

	if err != nil {
		m.AppendMessage(fmt.Sprintf("Cannot connect to server at %s. Is the server running?", RESTAPIBase), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to fetch worlds from server.", "error")
		return
	}

	var result struct {
		Count   int      `json:"count"`
		Default string   `json:"default"`
		Worlds  []string `json:"worlds"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing worlds: %v", err), "error")
		return
	}
	availableWorlds = result.Worlds
	if len(availableWorlds) > 0 {
		// Set current world to first available world if not already set
		if m.currentWorld == "" {
			m.currentWorld = result.Worlds[0]
		}
	}
}

// handleWorldSelectInput handles input for the world selection screen
func (m *model) handleWorldSelectInput(input string) {
	input = strings.ToLower(input)

	// Check if we have worlds loaded
	if len(availableWorlds) == 0 {
		m.fetchWorlds()
		// Don't process this input yet, wait for worlds to load
		m.AppendMessage("Loading available worlds...", "info")
		return
	}

	switch input {
	case "b", "back", "q", "quit":
		m.screen = ScreenWelcome
		m.textInput.SetValue("")
		m.inputBuffer = ""
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Parse world index (1-based)
		if idx := parseWorldIndex(input, len(availableWorlds)); idx >= 0 {
			m.currentWorld = availableWorlds[idx]
			m.AppendMessage(fmt.Sprintf("Selected world: %s", m.currentWorld), "success")
			m.fetchCharactersByWorld()
			m.screen = ScreenCharacterSelect
			m.textInput.SetValue("")
			m.inputBuffer = ""
			return
		}
	default:
		// Check if input matches a world name exactly
		for _, world := range availableWorlds {
			if strings.ToLower(world) == input {
				m.currentWorld = world
				m.AppendMessage(fmt.Sprintf("Selected world: %s", m.currentWorld), "success")
				m.fetchCharactersByWorld()
				m.screen = ScreenCharacterSelect
				m.textInput.SetValue("")
				m.inputBuffer = ""
				return
			}
		}
		if input != "" {
			m.AppendMessage(fmt.Sprintf("Invalid choice. Type 1-%d to select a world, or 'b' to go back.", len(availableWorlds)), "error")
		}
		return
	}
}

// parseWorldIndex parses a string input to a world index (0-based)
func parseWorldIndex(input string, numWorlds int) int {
	var idx int
	fmt.Sscanf(input, "%d", &idx)
	if idx > 0 && idx <= numWorlds {
		return idx - 1
	}
	return -1
}

// displayWorlds returns the formatted world selection menu
func (m *model) displayWorlds() string {
	var buf bytes.Buffer

	if len(availableWorlds) == 0 {
		buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("Fetching available worlds..."))
		buf.WriteString("\n\n")
	} else {
		for idx, world := range availableWorlds {
			numStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(fmt.Sprintf("%d.", idx+1))
			nameStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(world)
			line := fmt.Sprintf("  %s  %s", numStyle, nameStyle)
			if world == m.currentWorld {
				line += lipgloss.NewStyle().Foreground(PrimaryGold).Render("  [ACTIVE]")
			}
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// ============================================================
// CHARACTER SELECTION
// ============================================================

// handleCharacterSelectInput handles input for the character selection screen
func (m *model) handleCharacterSelectInput(input string) {
	// Handle character creation flow
	if m.isCreatingCharacter {
		m.handleCharacterCreationInput(input)
		return
	}

	input = strings.ToLower(input)

	switch input {
	case "b", "back", "q", "quit":
		m.screen = ScreenWorldSelect
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.AppendMessage("Select a world:", "info")
	case "r", "refresh":
		m.fetchCharactersByWorld()
		m.AppendMessage("Characters refreshed.", "info")
	case "n", "new", "create":
		m.startCharacterCreation()
	default:
		// Try to parse as character number (1-9)
		if idx := parseWorldIndex(input, len(m.selectedWorldCharacters)); idx >= 0 {
			m.selectCharacter(idx)
			return
		}
		// Try exact name match
		for _, char := range m.selectedWorldCharacters {
			if strings.ToLower(char.Name) == input {
				m.selectCharacterByName(char.Name)
				return
			}
		}
		if input != "" {
			numChars := len(m.selectedWorldCharacters)
			if numChars > 0 {
				m.AppendMessage(fmt.Sprintf("Invalid choice. Type 1-%d to select, 'n' for new, or 'b' to go back.", numChars), "error")
			} else {
				m.AppendMessage("Type 'n' to create your first character, or 'b' to go back.", "info")
			}
		}
	}
}

// fetchCharactersByWorld fetches characters for the current world from the API
func (m *model) fetchCharactersByWorld() {
	if m.currentWorld == "" || m.currentUserID == 0 {
		return
	}

	m.isLoading = true
	m.loadingMessage = "Fetching characters..."

	// Use /user-characters/:id and filter by world
	resp, err := http.Get(fmt.Sprintf("%s/user-characters/%d", RESTAPIBase, m.currentUserID))
	m.isLoading = false

	if err != nil {
		m.AppendMessage(fmt.Sprintf("Cannot connect to server at %s.", RESTAPIBase), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to fetch characters.", "error")
		return
	}

	var characters []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing characters: %v", err), "error")
		return
	}

	// Filter characters by current world
	m.selectedWorldCharacters = []CharacterInfo{}
	for _, char := range characters {
		// Get full character details to check world
		charName, _ := char["name"].(string)
		charID, _ := char["id"].(float64)
		if charName == "" {
			continue
		}

		// Fetch full character to get world
		charResp, err := http.Get(fmt.Sprintf("%s/characters/%d", RESTAPIBase, int(charID)))
		if err != nil {
			continue
		}
		defer charResp.Body.Close()

		if charResp.StatusCode != http.StatusOK {
			continue
		}

		var fullChar map[string]interface{}
		if err := json.NewDecoder(charResp.Body).Decode(&fullChar); err != nil {
			continue
		}

		// Check if character belongs to current world
		charWorld, _ := fullChar["currentWorld"].(string)
		if charWorld != m.currentWorld {
			continue
		}

		ci := CharacterInfo{
			ID: int(charID),
		}
		if name, ok := char["name"].(string); ok {
			ci.Name = name
		}
		if race, ok := fullChar["race"].(string); ok {
			ci.Race = race
		}
		if class, ok := fullChar["class"].(string); ok {
			ci.Class = class
		}
		if level, ok := fullChar["level"].(float64); ok {
			ci.Level = int(level)
		}
		if gender, ok := fullChar["gender"].(string); ok {
			ci.Gender = gender
		}
		if hp, ok := fullChar["hitpoints"].(float64); ok {
			ci.Hitpoints = int(hp)
		}
		if maxHP, ok := fullChar["max_hitpoints"].(float64); ok {
			ci.MaxHitpoints = int(maxHP)
		}

		m.selectedWorldCharacters = append(m.selectedWorldCharacters, ci)
	}
}

// displayCharacters returns the formatted character selection menu
func (m *model) displayCharacters() string {
	var buf bytes.Buffer

	worldLabel := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render("World:")
	worldVal := lipgloss.NewStyle().Foreground(TextWhite).Render(m.currentWorld)
	buf.WriteString(fmt.Sprintf("  %s  %s\n\n", worldLabel, worldVal))

	if len(m.selectedWorldCharacters) == 0 {
		buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("  No characters in this world."))
		buf.WriteString("\n")
		buf.WriteString(lipgloss.NewStyle().Foreground(TextGray).Render("  Type 'n' to create a new character."))
		buf.WriteString("\n\n")
	} else {
		for idx, char := range m.selectedWorldCharacters {
			numStyle := lipgloss.NewStyle().Foreground(AccentBlue).Bold(true).Render(fmt.Sprintf("%d.", idx+1))
			nameStyle := lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(char.Name)
			details := fmt.Sprintf("Lvl %d %s %s", char.Level, char.Race, char.Class)
			detailsStyle := lipgloss.NewStyle().Foreground(TextGray).Render(details)
			hpLabel := lipgloss.NewStyle().Foreground(StatusRed).Render("HP")
			hpStyle := lipgloss.NewStyle().Foreground(TextWhite).Render(fmt.Sprintf("%d/%d", char.Hitpoints, char.MaxHitpoints))
			buf.WriteString(fmt.Sprintf("  %s  %s\n", numStyle, nameStyle))
			buf.WriteString(fmt.Sprintf("       %s  %s  %s\n", detailsStyle, hpLabel+":", hpStyle))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// selectCharacter selects a character by index
func (m *model) selectCharacter(idx int) {
	if idx < 0 || idx >= len(m.selectedWorldCharacters) {
		return
	}

	char := m.selectedWorldCharacters[idx]
	m.loadCharacter(char.ID)
	m.AppendMessage(fmt.Sprintf("Selected character: %s", char.Name), "success")
}

// selectCharacterByName selects a character by name
func (m *model) selectCharacterByName(name string) {
	for _, char := range m.selectedWorldCharacters {
		if char.Name == name {
			m.loadCharacter(char.ID)
			m.AppendMessage(fmt.Sprintf("Selected character: %s", char.Name), "success")
			return
		}
	}
}

// loadCharacter loads character data and transitions to playing screen
func (m *model) loadCharacter(charID int) {
	resp, err := http.Get(fmt.Sprintf("%s/characters/%d", RESTAPIBase, charID))
	if err != nil {
		m.AppendMessage("Failed to load character.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Character not found.", "error")
		return
	}

	var char map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		m.AppendMessage("Error loading character data.", "error")
		return
	}

	// Populate model fields
	if id, ok := char["id"].(float64); ok {
		m.currentCharacterID = int(id)
	}
	if name, ok := char["name"].(string); ok {
		m.currentCharacterName = name
	}
	if race, ok := char["race"].(string); ok {
		m.characterRace = race
	}
	if class, ok := char["class"].(string); ok {
		m.characterClass = class
	}
	if gender, ok := char["gender"].(string); ok {
		m.characterGender = gender
	}
	if desc, ok := char["description"].(string); ok {
		m.characterDescription = desc
	}
	if hp, ok := char["hitpoints"].(float64); ok {
		m.characterHP = int(hp)
	}
	if maxHP, ok := char["max_hitpoints"].(float64); ok {
		m.characterMaxHP = int(maxHP)
	}
	if stamina, ok := char["stamina"].(float64); ok {
		m.characterStamina = int(stamina)
	}
	if maxStamina, ok := char["max_stamina"].(float64); ok {
		m.characterMaxStamina = int(maxStamina)
	}
	if mana, ok := char["mana"].(float64); ok {
		m.characterMana = int(mana)
	}
	if maxMana, ok := char["max_mana"].(float64); ok {
		m.characterMaxMana = int(maxMana)
	}
	if level, ok := char["level"].(float64); ok {
		m.characterLevel = int(level)
	}

	// Transition to playing screen
	m.screen = ScreenPlaying
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.selectedWorldCharacters = []CharacterInfo{}
	m.isCreatingCharacter = false
	// Check if this is a new character (first room visit)
	isNewCharacter := !m.visitedRooms[m.currentRoom]
	m.visitedRooms[m.currentRoom] = true

	// Determine reconnect room
	targetRoomID := m.determineReconnectRoom()
	m.loadRoom(targetRoomID)

	// Show onboarding for new characters
	if isNewCharacter {
		m.AppendMessage(m.getWelcomeMessage(), "info")
	}

	// Update last seen
	m.updateLastSeenAt()

	m.effectsService.FireEvent("on_login", m.currentCharacterID, "", map[string]interface{}{
		"room_id": m.currentRoom,
	})
}

// ============================================================
// CHARACTER CREATION
// ============================================================

// Character creation input state
var createCharName string
var createCharPassword string
var createCharRace string

func (m *model) startCharacterCreation() {
	m.isCreatingCharacter = true
	m.inputField = "char_name"
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.AppendMessage("Creating new character in: "+m.currentWorld, "info")
	m.AppendMessage("Enter character name (letters only, 1-23 chars):", "info")
	m.textInput.Focus()
}

func (m *model) handleCharacterCreationInput(input string) {
	// Handle cancel
	if strings.ToLower(input) == "cancel" || strings.ToLower(input) == "c" {
		m.cancelCharacterCreation()
		return
	}

	switch m.inputField {
	case "char_name":
		if input == "" {
			m.AppendMessage("Name cannot be empty. Enter character name:", "error")
			return
		}
		if len(input) > 23 {
			m.AppendMessage("Name too long (max 23 characters). Try again:", "error")
			return
		}
		// Validate letters only
		for _, ch := range input {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				m.AppendMessage("Name can only contain letters (a-z, A-Z). Try again:", "error")
				return
			}
		}
		createCharName = input
		m.inputField = "char_password"
		m.AppendMessage("Enter character password:", "info")
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.SetValue("")
		m.textInput.Focus()
	case "char_password":
		if input == "" {
			m.AppendMessage("Password cannot be empty. Enter password:", "error")
			return
		}
		createCharPassword = input
		m.inputField = "char_confirm_password"
		m.AppendMessage("Confirm password:", "info")
		m.textInput.SetValue("")
		m.textInput.Focus()
	case "char_confirm_password":
		if input != createCharPassword {
			m.AppendMessage("Passwords do not match. Try again:", "error")
			m.inputField = "char_password"
			createCharPassword = ""
			m.textInput.EchoMode = textinput.EchoPassword
			m.textInput.SetValue("")
			m.textInput.Focus()
			return
		}
		m.inputField = "char_race"
		m.textInput.EchoMode = textinput.EchoNormal
		m.AppendMessage("Select race (or press Enter for human):", "info")
		m.AppendMessage("Available: human, elf, dwarf, halfling, orc", "info")
		m.textInput.SetValue("")
		m.textInput.Focus()
	case "char_race":
		race := strings.ToLower(input)
		if race == "" {
			race = "human"
		}
		validRaces := map[string]bool{"human": true, "elf": true, "dwarf": true, "halfling": true, "orc": true}
		if !validRaces[race] {
			m.AppendMessage("Invalid race. Choose: human, elf, dwarf, halfling, orc", "error")
			return
		}
		createCharRace = race
		m.inputField = "char_class"
		m.AppendMessage(fmt.Sprintf("Race selected: %s", race), "success")
		m.AppendMessage("Select class (or press Enter for adventurer):", "info")
		m.AppendMessage("Available: adventurer, warrior, mage, rogue, cleric", "info")
		m.textInput.SetValue("")
		m.textInput.Focus()
	case "char_class":
		class := strings.ToLower(input)
		if class == "" {
			class = "adventurer"
		}
		validClasses := map[string]bool{"adventurer": true, "warrior": true, "mage": true, "rogue": true, "cleric": true}
		if !validClasses[class] {
			m.AppendMessage("Invalid class. Choose: adventurer, warrior, mage, rogue, cleric", "error")
			return
		}
		// Create the character
		m.createCharacter(createCharName, createCharPassword, createCharRace, class)
	default:
		// Unknown state, reset
		m.cancelCharacterCreation()
	}
}

func (m *model) createCharacter(name, password, race, class string) {
	m.isLoading = true
	m.loadingMessage = "Creating character..."

	jsonData, _ := json.Marshal(map[string]interface{}{
		"name":     name,
		"password": password,
		"race":     race,
		"class":    class,
		"world":    m.currentWorld,
	})

	resp, err := http.Post(fmt.Sprintf("%s/user-characters/%d", RESTAPIBase, m.currentUserID),
		"application/json", bytes.NewBuffer(jsonData))
	m.isLoading = false

	if err != nil {
		m.AppendMessage(fmt.Sprintf("Connection error: %v", err), "error")
		m.cancelCharacterCreation()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		m.AppendMessage("Character name already taken. Try a different name.", "error")
		m.inputField = "char_name"
		createCharName = ""
		createCharPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		m.textInput.SetValue("")
		m.textInput.Focus()
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok {
			m.AppendMessage(errMsg, "error")
		} else {
			m.AppendMessage("Failed to create character.", "error")
		}
		m.cancelCharacterCreation()
		return
	}

	if resp.StatusCode != http.StatusCreated {
		m.AppendMessage(fmt.Sprintf("Server error (status %d). Try again.", resp.StatusCode), "error")
		m.cancelCharacterCreation()
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage("Error processing response.", "error")
		m.cancelCharacterCreation()
		return
	}

	// Reset creation state
	createCharName = ""
	createCharPassword = ""
	createCharRace = ""
	m.isCreatingCharacter = false
	m.inputField = ""
	m.textInput.EchoMode = textinput.EchoNormal

	// Load the new character
	if id, ok := result["id"].(float64); ok {
		m.AppendMessage(fmt.Sprintf("Character '%s' created successfully!", name), "success")
		m.loadCharacter(int(id))
	}
}

// getWelcomeMessage returns an onboarding message for new characters
func (m *model) getWelcomeMessage() string {
	msg := lipgloss.NewStyle().Bold(true).Foreground(PrimaryGold).Render("Welcome to Herbst MUD!")
	msg += "\n\n"
	msg += lipgloss.NewStyle().Foreground(AccentBlue).Render("Essential Commands:")
	msg += "\n"
	type cmdHelp struct{ cmd, desc string }
	commands := []cmdHelp{
		{"look", "Examine your surroundings"},
		{"north / south / east / west", "Move between rooms"},
		{"say <text>", "Speak to others in the room"},
		{"help", "Show all available commands"},
		{"who", "See who's online"},
	}
	for _, c := range commands {
		cmdStyle := lipgloss.NewStyle().Foreground(PrimaryGold).Bold(true).Render(c.cmd)
		msg += fmt.Sprintf("  %s — %s\n", cmdStyle, c.desc)
	}
	msg += "\n"
	msg += lipgloss.NewStyle().Foreground(TextGray).Italic(true).Render("Tip: try 'look' to see where you are!")
	return msg
}

func (m *model) cancelCharacterCreation() {
	createCharName = ""
	createCharPassword = ""
	createCharRace = ""
	m.isCreatingCharacter = false
	m.inputField = ""
	m.textInput.EchoMode = textinput.EchoNormal
	m.screen = ScreenCharacterSelect
	m.textInput.SetValue("")
	m.inputBuffer = ""
}
