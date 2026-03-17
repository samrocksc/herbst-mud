package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	_ "github.com/lib/pq"
	"github.com/muesli/termenv"
	"herbst/db"
	"herbst/db/character"
	"herbst/db/user"
	"herbst/dbinit"
)

func init() {
	// Force TERM for proper color support - always set regardless of client value
	os.Setenv("TERM", "xterm-256color")
}

// RESTAPIBase is the base URL for the REST API
const RESTAPIBase = "http://localhost:8080"

// StartingRoomID is the ID of the room players start in
const StartingRoomID = 5

// Screen states
const (
	ScreenWelcome        = "welcome"
	ScreenLogin          = "login"
	ScreenRegister       = "register"
	ScreenPlaying        = "playing"
	ScreenProfile        = "profile"
	ScreenEditField      = "edit_field"
	ScreenFountainWake   = "fountain_wake"
	ScreenFountainWash   = "fountain_wash"
	ScreenCharacterCreate = "character_create"
)

// Menu selection constants for vim-style navigation
const (
	MenuWelcome = iota
	MenuProfile
)

// getDBConfig returns database connection config from environment variables
func getDBConfig() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "herbst")
	password := getEnv("DB_PASSWORD", "herbst_password")
	dbname := getEnv("DB_NAME", "herbst_mud")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// getEnv returns environment variable or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize database
	client, err := db.Open("postgres", getDBConfig())
	if err != nil {
		log.Printf("Warning: failed connecting to postgres: %v", err)
	} else {
		defer client.Close()

		// Run auto migration tool
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Printf("Warning: failed creating schema resources: %v", err)
		} else {
			log.Println("Database initialized successfully")
		}

		// Initialize default admin user
		if err := dbinit.InitAdminUser(client); err != nil {
			log.Printf("Warning: failed to initialize admin user: %v", err)
		}
	}

	// Pass client to server options
	srv, err := wish.NewServer(
		wish.WithAddress(":4444"),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			logging.Middleware(),
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					log.Printf("New connection from %s", s.RemoteAddr().String())

					// Force TrueColor for this session
					lipgloss.SetColorProfile(termenv.TrueColor)

					// Create initial text input for login/register
					ti := textinput.New()
					ti.Placeholder = "Enter choice..."
					ti.Focus()

					// Create loading spinner
					sp := spinner.New()
					sp.Spinner = spinner.Dot
					sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

					// Create program with shared client
					p := tea.NewProgram(
						&model{
							connectedAt:  time.Now(),
							session:      s,
							client:       client,
							screen:       ScreenWelcome,
							currentRoom:  StartingRoomID,
							textInput:    ti,
							spinner:      sp,
							visitedRooms: make(map[int]bool),
							knownExits:   make(map[string]bool),
						},
						tea.WithInput(s),
						tea.WithOutput(s),
					)

					if _, err := p.Run(); err != nil {
						log.Printf("Bubbletea error: %v", err)
					}

					log.Printf("Connection from %s closed", s.RemoteAddr().String())
				}
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting SSH server on :4444")
	if err = srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

type model struct {
	connectedAt time.Time
	session     ssh.Session
	client      *db.Client
	width       int
	height      int
	err         error

	// Screen state
	screen string

	// Auth state
	currentUserID   int
	currentUserName string

	// Input handling (replaces manual inputBuffer)
	textInput textinput.Model

	// Login/Register input state
	inputField    string // "username" or "password"
	loginUsername string
	loginPassword string

	// Player state
	currentRoom int
	roomName    string
	roomDesc    string
	exits       map[string]int
	inputBuffer string
	message     string
	messageType string // "success", "error", "info" for styling

	// Menu navigation state (vim-style)
	menuCursor int
	menuItems  []string

	// Character state (for whoami/profile)
	currentCharacterID   int
	currentCharacterName string
	characterGender      string
	characterDescription string
	characterHP          int
	characterMaxHP       int
	characterStamina     int
	characterMaxStamina  int
	characterMana        int
	characterMaxMana     int
	characterLevel       int
	characterExperience  int

	// Profile editing state
	editField string // "gender" or "description"

	// Loading state
	spinner        spinner.Model
	isLoading      bool
	loadingMessage string

	// Room tracking
	visitedRooms   map[int]bool
	knownExits     map[string]bool // For color-coded exits
	roomCharacters []roomCharacter // Characters in current room
}

// roomCharacter represents a character in the room (NPC or player)
type roomCharacter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsNPC    bool   `json:"isNPC"`
	Level    int    `json:"level"`
	Class    string `json:"class"`
	Race     string `json:"race"`
	UserID   int    `json:"userId"`
}

	// Render state (prevents double rendering of messages in same frame)
	renderDone bool
}

// ============================================================
// STYLING - Lipgloss styles for UI elements
// ============================================================

var (
	// Colors
	red    = lipgloss.Color("196")
	green  = lipgloss.Color("46")
	yellow = lipgloss.Color("226")
	blue   = lipgloss.Color("75")
	cyan   = lipgloss.Color("81")  // Lighter blue for login
	purple = lipgloss.Color("141")
	white  = lipgloss.Color("15")
	gray   = lipgloss.Color("8")
	pink   = lipgloss.Color("219")

	// Raw ANSI for direct terminal output (when lipgloss fails)
	pinkAnsi  = "\033[38;5;219m"
	pinkReset = "\033[0m"

	// Exit colors for visited/known/new
	exitVisitedColor = lipgloss.Color("46")  // Green
	exitKnownColor   = lipgloss.Color("226") // Yellow
	exitNewColor     = lipgloss.Color("15")  // White

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(yellow)

	menuSelectedStyle = lipgloss.NewStyle().
				Foreground(green).
				Bold(true).
				Padding(0, 0, 0, 2)

	menuNormalStyle = lipgloss.NewStyle().
			Foreground(gray).
			Padding(0, 0, 0, 2)

	promptStyle = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)
)

// ============================================================
// PROGRESS BARS - HP/Stamina/Mana using characters
// ============================================================

// ProgressBar creates a text-based progress bar
func ProgressBar(current, max, width int, filledChar, emptyChar string, fillColor, emptyColor lipgloss.Color) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	filledWidth := int(float64(current) / float64(max) * float64(width))
	emptyWidth := width - filledWidth

	filledStyle := lipgloss.NewStyle().Foreground(fillColor)
	emptyStyle := lipgloss.NewStyle().Foreground(emptyColor)

	return filledStyle.Render(strings.Repeat(filledChar, filledWidth)) +
		emptyStyle.Render(strings.Repeat(emptyChar, emptyWidth))
}

// StatusBar creates colorful status bars with progress bars
func StatusBar(hp, maxHP, stamina, maxStamina, mana, maxMana int) string {
	// Use Unicode block characters for progress
	hpBar := ProgressBar(hp, maxHP, 20, "█", "░", red, gray)
	staminaBar := ProgressBar(stamina, maxStamina, 20, "▓", "░", yellow, gray)
	manaBar := ProgressBar(mana, maxMana, 20, "✨", "○", blue, gray)

	return fmt.Sprintf(" %s  HP: %s %d/%d\n %s  STA: %s %d/%d\n %s  MANA: %s %d/%d",
		lipgloss.NewStyle().Foreground(red).Render("❤️"),
		hpBar, hp, maxHP,
		lipgloss.NewStyle().Foreground(yellow).Render("💪"),
		staminaBar, stamina, maxStamina,
		lipgloss.NewStyle().Foreground(blue).Render("✨"),
		manaBar, mana, maxMana)
}

// MiniStatusBar creates a compact inline status bar
func MiniStatusBar(hp, maxHP, stamina, maxStamina, mana, maxMana int) string {
	hpPercent := float64(hp) / float64(maxHP) * 100
	staminaPercent := float64(stamina) / float64(maxStamina) * 100
	manaPercent := float64(mana) / float64(maxMana) * 100

	var hpColor lipgloss.Color
	if hpPercent > 60 {
		hpColor = green
	} else if hpPercent > 30 {
		hpColor = yellow
	} else {
		hpColor = red
	}

	hpStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", hpPercent))
	staStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", staminaPercent))
	manaStr := lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%.0f", manaPercent))

	return fmt.Sprintf("[%s%s%% %s%s%% %s%s%%]",
		lipgloss.NewStyle().Foreground(hpColor).Render("❤️"),
		hpStr,
		lipgloss.NewStyle().Foreground(yellow).Render("💪"),
		staStr,
		lipgloss.NewStyle().Foreground(blue).Render("✨"),
		manaStr)
}

// ============================================================
// STYLING HELPERS
// ============================================================

// styledMessage returns a styled message based on messageType
func (m *model) styledMessage(msg string) string {
	if msg == "" {
		return ""
	}

	switch m.messageType {
	case "success":
		return successStyle.Render("✓ ") + msg
	case "error":
		return errorStyle.Render("✗ ") + msg
	case "info":
		return infoStyle.Render("ℹ ") + msg
	default:
		return msg
	}
}

// formatExitsWithColor returns color-coded exits
func (m *model) formatExitsWithColor() string {
	if len(m.exits) == 0 {
		return lipgloss.NewStyle().Foreground(gray).Render("none")
	}

	var dirs []string
	for dir, roomID := range m.exits {
		var exitStyle lipgloss.Style
		if m.visitedRooms[roomID] {
			// Green = visited
			exitStyle = lipgloss.NewStyle().Foreground(exitVisitedColor)
		} else if m.knownExits[dir] {
			// Yellow = known but not visited
			exitStyle = lipgloss.NewStyle().Foreground(exitKnownColor)
		} else {
			// White = new
			m.knownExits[dir] = true
			exitStyle = lipgloss.NewStyle().Foreground(exitNewColor)
		}
		dirs = append(dirs, exitStyle.Render(dir))
	}

	return strings.Join(dirs, ", ")
}

// ============================================================
// MODEL LIFECYCLE
// ============================================================

func (m model) Init() tea.Cmd {
	// Don't load room info yet - we need to login first
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		// Update spinner if loading
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		key := msg.String()
		// Debug: log.Printf("KeyMsg received: %q, screen: %s", key, m.screen)

		// Global quit handling
		if key == "ctrl+c" || key == "ctrl+q" {
			return m, tea.Quit
		}

		// Handle loading state - block most input
		if m.isLoading {
			return m, nil
		}

		// Vim-style navigation and arrow keys for menu selection
		if m.screen == ScreenWelcome || m.screen == ScreenProfile {
			if key == "j" || key == "down" {
				m.menuCursor++
				if m.menuCursor >= len(m.menuItems) {
					m.menuCursor = 0
				}
				return m, nil
			}
			if key == "k" || key == "up" {
				m.menuCursor--
				if m.menuCursor < 0 {
					m.menuCursor = len(m.menuItems) - 1
				}
				return m, nil
			}
		}

		// Handle Enter key - process command based on screen
		if key == "enter" || key == "ctrl+j" || key == "ctrl+m" {
			input := m.textInput.Value()
			m.textInput.SetValue("")
			// Debug: log.Printf("Enter pressed, calling processInput with: %q", input)
			m.processInput(input)
			// Debug: log.Printf("After processInput, screen is: %s", m.screen)
			return m, nil
		}

		// Handle Escape - go back to previous screen
		if key == "esc" {
			m.handleEscape()
			return m, nil
		}

		// Let text input handle regular typing
		m.textInput, cmd = m.textInput.Update(msg)
		m.inputBuffer = m.textInput.Value()

		return m, cmd
	}

	return m, nil
}

func (m *model) handleEscape() {
	switch m.screen {
	case ScreenLogin, ScreenRegister:
		m.screen = ScreenWelcome
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.loginUsername = ""
		m.loginPassword = ""
		m.inputField = "username"
		m.message = ""
		m.messageType = "info"
		// Re-initialize menu items for welcome screen
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	case ScreenProfile, ScreenEditField:
		m.screen = ScreenPlaying
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = "info"
	case ScreenPlaying:
		// Could add a "really quit?" confirmation
		m.message = "Type 'quit' or press Ctrl+C to exit"
		m.messageType = "info"
	}
}

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
	case "1", "login", "l":
		m.screen = ScreenLogin
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.message = "Enter your username:"
		m.messageType = "info"
		m.textInput.Focus()
	case "2", "register", "r", "create":
		m.screen = ScreenRegister
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.message = "Choose a username:"
		m.messageType = "info"
		m.textInput.Focus()
	case "3", "quit", "q":
		m.message = "Goodbye! Thanks for playing Herbst MUD."
		m.messageType = "success"
		m.inputBuffer = ""
		return
	default:
		if input != "" {
			m.message = "Invalid choice. Type 1, 2, or 3"
			m.messageType = "error"
		}
	}
}

func (m *model) handleLoginInput(input string) {
	if m.inputField == "username" {
		m.loginUsername = input
		m.inputField = "password"
		m.message = "Enter your password:"
		m.messageType = "info"
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.Focus()
	} else if m.inputField == "password" {
		m.loginPassword = input
		m.textInput.EchoMode = textinput.EchoNormal
		m.attemptLogin()
	}
}

func (m *model) attemptLogin() {
	// Debug: log.Printf("attemptLogin called with username: %q", m.loginUsername)

	// Start loading state
	m.isLoading = true
	m.loadingMessage = "Logging in..."

	// Use REST API for authentication (simulated - in real use would need async)
	jsonData, _ := json.Marshal(map[string]string{
		"email":    m.loginUsername,
		"password": m.loginPassword,
	})

	// Debug: log.Printf("Sending auth request to %s/users/auth", RESTAPIBase)
	resp, err := http.Post(RESTAPIBase+"/users/auth", "application/json", bytes.NewBuffer(jsonData))
	m.isLoading = false

	if err != nil {
		// Debug: log.Printf("Connection error: %v", err)
		m.message = fmt.Sprintf("Connection error: %v", err)
		m.messageType = "error"
		return
	}
	defer resp.Body.Close()

	// Debug: log.Printf("Auth response status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		m.message = "Invalid username or password. Try again."
		m.messageType = "error"
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.textInput.EchoMode = textinput.EchoNormal
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.message = fmt.Sprintf("Login error: %v", err)
		m.messageType = "error"
		return
	}

	// Login successful
	if id, ok := result["id"].(float64); ok {
		m.currentUserID = int(id)
	}
	if email, ok := result["email"].(string); ok {
		m.currentUserName = email
	}
	m.screen = ScreenPlaying
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.message = fmt.Sprintf("Welcome back, %s!", m.currentUserName)
	m.messageType = "success"

	// Load or create character for this user
	m.loadOrCreateCharacter()

	// Mark starting room as visited
	m.visitedRooms[m.currentRoom] = true

	// Load starting room info (can still use direct DB for this)
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

		// Mark exits as known
		for dir := range m.exits {
			m.knownExits[dir] = true
		}

		// Fetch characters in the room (NPCs and players)
		m.fetchRoomCharacters()
	}
}

func (m *model) handleRegisterInput(input string) {
	if m.inputField == "username" {
		if input == "" {
			m.message = "Username cannot be empty. Try again:"
			m.messageType = "error"
			return
		}
		m.loginUsername = input
		m.inputField = "password"
		m.message = "Choose a password:"
		m.messageType = "info"
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.Focus()
	} else if m.inputField == "password" {
		if input == "" {
			m.message = "Password cannot be empty. Try again:"
			m.messageType = "error"
			return
		}
		m.loginPassword = input
		m.inputField = "confirm_password"
		m.message = "Confirm your password:"
		m.messageType = "info"
		m.textInput.Focus()
	} else if m.inputField == "confirm_password" {
		if input != m.loginPassword {
			m.message = "Passwords do not match. Try again:"
			m.messageType = "error"
			m.inputField = "password"
			m.loginPassword = ""
			m.textInput.EchoMode = textinput.EchoPassword
			m.textInput.Focus()
			return
		}
		m.inputField = "email"
		m.message = "Enter your email (optional, press enter to skip):"
		m.messageType = "info"
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

func (m *model) attemptRegistration(email string) {
	// Use REST API for user creation
	jsonData, _ := json.Marshal(map[string]string{
		"email":    m.loginUsername,
		"password": m.loginPassword,
	})

	resp, err := http.Post(RESTAPIBase+"/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		m.message = fmt.Sprintf("Connection error: %v", err)
		m.messageType = "error"
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok && (strings.Contains(errMsg, "unique") || strings.Contains(errMsg, "already exists")) {
			m.message = "Username already taken. Choose a different one."
			m.messageType = "error"
			m.inputField = "username"
			m.loginUsername = ""
			m.loginPassword = ""
			m.textInput.EchoMode = textinput.EchoNormal
			return
		}
		m.message = "Failed to create account. Please try again."
		m.messageType = "error"
		return
	}

	if resp.StatusCode != http.StatusCreated {
		m.message = "Failed to create account. Please try again."
		m.messageType = "error"
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.message = fmt.Sprintf("Error processing response: %v", err)
		m.messageType = "error"
		return
	}

	// Auto-login after registration
	if id, ok := result["id"].(float64); ok {
		m.currentUserID = int(id)
	}
	if email, ok := result["email"].(string); ok {
		m.currentUserName = email
	}
	m.screen = ScreenPlaying
	m.textInput.SetValue("")
	m.inputBuffer = ""
	m.message = fmt.Sprintf("Account created! Welcome to Herbst MUD, %s!", m.currentUserName)
	m.messageType = "success"

	// Load or create character for this user
	m.loadOrCreateCharacter()

	// Mark starting room as visited
	m.visitedRooms[StartingRoomID] = true

	// Load starting room info
	room, err := m.client.Room.Get(context.Background(), StartingRoomID)
	if err != nil {
		m.err = fmt.Errorf("failed to load starting room: %v", err)
		return
	}
	m.currentRoom = room.ID
	m.roomName = room.Name
	m.roomDesc = room.Description
	m.exits = room.Exits

	// Mark exits as known
	for dir := range m.exits {
		m.knownExits[dir] = true
	}

	// Fetch characters in the room (NPCs and players)
	m.fetchRoomCharacters()
}

// ============================================================
// GAME COMMAND PROCESSING
// ============================================================

func (m *model) processCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))

	if cmd == "" {
		return
	}

	// Handle movement commands
	if m.handleMovement(cmd) {
		return
	}

	// Handle other commands
	switch cmd {
	case "help", "?":
		m.message = `Commands:
  n/north, s/south, e/east, w/west - Move
  look/l - Look around
  exits/x - Show exits  
  peer <dir> - Peek at adjacent room
  whoami - Show your info
  profile/p - Edit character profile
  clear/cls - Clear screen
  quit - Exit game`
		m.messageType = "info"
	case "look", "l":
		m.message = fmt.Sprintf("[%s]\n%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor())
		m.messageType = "info"
	case "exits", "x":
		m.message = fmt.Sprintf("Exits: %s", m.formatExitsWithColor())
		m.messageType = "info"
	case "whoami":
		// Show character info including level with progress bars
		m.message = fmt.Sprintf("=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
			m.currentUserName, m.currentUserID, m.roomName,
			m.characterLevel, m.characterExperience,
			StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana))
		m.messageType = "info"
	case "profile", "p":
		m.screen = ScreenProfile
		m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
		m.menuCursor = 0
		m.message = ""
		m.messageType = "info"
	case "peer":
		m.handlePeerCommand(cmd)
	case "clear", "cls":
		// Clear the terminal screen - reset message buffer
		m.message = ""
		m.messageType = ""
		m.inputBuffer = ""
		return
	case "quit", "q":
		m.message = "Thanks for playing! Goodbye!"
		m.messageType = "success"
		m.inputBuffer = ""
		return
	default:
		m.message = fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd)
		m.messageType = "error"
	}
}

func (m *model) handleMovement(cmd string) bool {
	directionMap := map[string]string{
		"n": "north", "north": "north",
		"s": "south", "south": "south",
		"e": "east", "east": "east",
		"w": "west", "west": "west",
	}

	direction, ok := directionMap[cmd]
	if !ok {
		return false
	}

	// Check if exit exists
	nextRoomID, ok := m.exits[direction]
	if !ok {
		m.message = "You can't go that way."
		m.messageType = "error"
		return true
	}

	// Mark exit as known
	m.knownExits[direction] = true

	// Move to new room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.message = fmt.Sprintf("Error moving: %v", err)
			m.messageType = "error"
			return true
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits

		// Mark new room as visited
		wasVisited := m.visitedRooms[m.currentRoom]
		m.visitedRooms[m.currentRoom] = true

		// Mark new exits as known
		for dir := range m.exits {
			m.knownExits[dir] = true
		}

		// Fetch characters in the room (NPCs and players)
		m.fetchRoomCharacters()

		if wasVisited {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s\n\nExits: %s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				m.formatExitsWithColor())
		} else {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s\n\nExits: %s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				m.formatExitsWithColor())
		}
		m.messageType = "success"
	}

	return true
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
		m.messageType = "info"
	case "2":
		m.editField = "description"
		m.screen = ScreenEditField
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.message = ""
		m.messageType = "info"
	case "3", "back", "b", "esc":
		m.screen = ScreenPlaying
		m.message = ""
		m.messageType = "info"
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

// createCharacterInDB creates a character in the database via the API
func (m *model) createCharacterInDB() {
	if m.currentUserID == 0 {
		log.Printf("Cannot create character: no user logged in")
		return
	}

	// Use the REST API to create the character
	jsonData, err := json.Marshal(map[string]interface{}{
		"name":    m.currentUserName,
		"userId":  m.currentUserID,
		"isNPC":   false,
	})
	if err != nil {
		log.Printf("Error marshaling character data: %v", err)
		return
	}

	resp, err := http.Post(RESTAPIBase+"/characters", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating character: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create character, status: %d", resp.StatusCode)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding character response: %v", err)
		return
	}
	if _, ok := result["id"]; !ok {
		log.Printf("Character created but no ID returned")
		return
	}

	if id, ok := result["id"].(float64); ok {
		m.currentCharacterID = int(id)
		log.Printf("Character created successfully with ID: %d", m.currentCharacterID)
	}

	// Set gender/description to default values that will be saved to DB
	m.characterGender = "unspecified"
	m.characterDescription = "A mysterious figure."

	// Now save these defaults to the DB
	m.saveProfileToDB()
}

// saveProfileToDB sends profile updates (gender, description) to the server
func (m *model) saveProfileToDB() {
	if m.currentCharacterID == 0 {
		return
	}

	jsonData, err := json.Marshal(map[string]string{
		"gender":      m.characterGender,
		"description": m.characterDescription,
	})
	if err != nil {
		log.Printf("Error marshaling profile data: %v", err)
		return
	}

	url := fmt.Sprintf("%s/characters/%d", RESTAPIBase, m.currentCharacterID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating profile update request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending profile update: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Profile update failed with status: %d", resp.StatusCode)
	}
}

func (m *model) handlePeerCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Usage: peer <direction>\nDirections: north, south, east, west, up, down"
		m.messageType = "error"
		return
	}
	direction := strings.ToLower(parts[1])

	// Validate direction
	validDirs := map[string]string{"north": "north", "south": "south", "east": "east", "west": "west", "up": "up", "down": "down"}
	dir, ok := validDirs[direction]
	if !ok {
		m.message = "Invalid direction. Use: north, south, east, west, up, down"
		m.messageType = "error"
		return
	}

	// Check if exit exists
	nextRoomID, ok := m.exits[dir]
	if !ok {
		m.message = "You can't peer that way — there's no exit."
		m.messageType = "error"
		return
	}

	// Get the room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.message = fmt.Sprintf("Error looking: %v", err)
			m.messageType = "error"
			return
		}

		m.message = fmt.Sprintf("You peer %s...\n\n[%s]\n%s",
			dir,
			lipgloss.NewStyle().Bold(true).Foreground(blue).Render(room.Name),
			room.Description)
		m.messageType = "info"
	}
}

func (m *model) loadOrCreateCharacter() {
	// Query for existing character for this user
	ctx := context.Background()
	chars, err := m.client.Character.Query().Where(character.HasUserWith(user.IDEQ(m.currentUserID))).All(ctx)
	if err != nil || len(chars) == 0 {
		// No character in DB yet - try to create one via API
		log.Printf("No character found for user %d, attempting to create...", m.currentUserID)
		m.createCharacterInDB()
		return
	}

	// Use first character
	char := chars[0]
	m.currentCharacterID = char.ID
	m.currentCharacterName = char.Name
	// Gender/Description - use from DB if available, else defaults
// 	m.characterGender = char.Gender
// 	m.characterDescription = char.Description
// 
// 	// Calculate HP/Stamina/Mana from stats (constitution → HP, dexterity → stamina, intelligence → mana)
// 	// Base: 50 + (stat * 5)
// 	m.characterMaxHP = 50 + (char.Constitution * 5)
// 	m.characterMaxStamina = 50 + (char.Dexterity * 5)
// 	m.characterMaxMana = 50 + (char.Intelligence * 5)
// 	m.characterHP = m.characterMaxHP
// 	m.characterStamina = m.characterMaxStamina
// 	m.characterMana = m.characterMaxMana
// 
// 	// Level/Experience - based on total stats for now
// 	// Level/Experience - defaults since stats not available
	m.characterLevel = 1
	m.characterExperience = 0
}

// ============================================================
// VIEW RENDERING
// ============================================================

func (m *model) View() string {
	var s strings.Builder

	// Debug: log.Printf("View() called, screen: %s, inputBuffer: %q", m.screen, m.inputBuffer)

	// Reset render state at start of each render cycle
	m.renderDone = false

	// Clear screen before each render to prevent previous state showing below new content
	// This fixes the "output re-rendering glitch" where previous state appears below new content
	s.WriteString("\033[2J\033[H")

	// Show loading spinner if loading
	if m.isLoading {
		s.WriteString(m.spinner.View())
		s.WriteString(" " + m.loadingMessage)
		// Clear message after rendering
		m.message = ""
		m.messageType = ""
		return s.String()
	}

	switch m.screen {
	case ScreenWelcome:
		s.WriteString(welcomeScreen())
		s.WriteString("\n\n")

		// Vim-style menu rendering
		for i, item := range m.menuItems {
			cursor := "  "
			if i == m.menuCursor {
				cursor = "▶ "
				s.WriteString(menuSelectedStyle.Render(cursor + item))
			} else {
				s.WriteString(menuNormalStyle.Render(cursor + item))
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")
		if m.message != "" {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n")
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

	case ScreenLogin:
		// Use split-screen layout like ScreenPlaying
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cyan).
			Padding(0, 1)

		loginContent := loginScreen(m.width, m.height)
		if m.message != "" {
			loginContent += "\n\n" + m.styledMessage(m.message)
		}

		s.WriteString(outputStyle.Render(loginContent))
		s.WriteString("\n\n")

		// Separator
		separatorStyle := lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)
		s.WriteString(separatorStyle.Render(strings.Repeat("─", int(math.Max(0, float64(m.width-4))))))
		s.WriteString("\n")

		// Input area
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cyan).
			Padding(0, 1)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

	case ScreenRegister:
		// Use split-screen layout like ScreenPlaying
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 1)

		registerContent := registerScreen(m.width, m.height)
		if m.message != "" {
			registerContent += "\n\n" + m.styledMessage(m.message)
		}

		s.WriteString(outputStyle.Render(registerContent))
		s.WriteString("\n\n")

		// Separator
		separatorStyle := lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)
		s.WriteString(separatorStyle.Render(strings.Repeat("─", int(math.Max(0, float64(m.width-4))))))
		s.WriteString("\n")

		// Input area
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 1)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

	case ScreenProfile:
		s.WriteString("=== CHARACTER PROFILE ===\n\n")
		s.WriteString(fmt.Sprintf("Name: %s\n", lipgloss.NewStyle().Bold(true).Render(m.currentCharacterName)))
		s.WriteString(fmt.Sprintf("Gender: %s\n", m.characterGender))
		s.WriteString(fmt.Sprintf("Description: %s\n\n", m.characterDescription))
		s.WriteString("Stats:\n")
		s.WriteString(StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana))
		s.WriteString("\n\n")

		// Vim-style menu for profile
		for i, item := range m.menuItems {
			cursor := "  "
			if i == m.menuCursor {
				cursor = "▶ "
				s.WriteString(menuSelectedStyle.Render(cursor + item))
			} else {
				s.WriteString(menuNormalStyle.Render(cursor + item))
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")
		if m.message != "" {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n")
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

	case ScreenEditField:
		if m.editField == "gender" {
			s.WriteString("Enter your gender (e.g., he/him, she/her, they/them):\n\n")
		} else {
			s.WriteString("Enter your description (what people see when they look at you):\n\n")
		}
		if m.message != "" {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n\n")
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

	case ScreenPlaying:
		// Track if we've already rendered the message in this frame
		messageRendered := false

		// Ensure we have valid dimensions - if not, use defaults
		width := m.width
		height := m.height
		if width < 40 {
			width = 80
		}
		if height < 10 {
			height = 24
		}

		// Calculate proportional heights
		// Input: ~20%, Status bar: ~10%, Viewport: ~70%
		inputHeight := height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}
		statusHeight := height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}
		viewportHeight := height - inputHeight - statusHeight
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		// Full-width output viewport (top ~70%)
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Padding(0, 1).
			Width(width - 2).          // Account for border
			Height(viewportHeight - 2) // Account for border

		// Colorful status bar with mini progress bars
		statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)

		// Build character list with color coding (NPCs = red, Players = green)
		var characterList string
		otherChars := make([]string, 0)
		for _, rc := range m.roomCharacters {
			// Skip own character
			if rc.ID == m.currentCharacterID {
				continue
			}
			var charLine string
			if rc.IsNPC {
				charLine = lipgloss.NewStyle().Foreground(red).Render("• " + rc.Name + " (NPC)")
			} else {
				charLine = lipgloss.NewStyle().Foreground(green).Render("• " + rc.Name + " (Player)")
			}
			otherChars = append(otherChars, charLine)
		}
		if len(otherChars) > 0 {
			characterList = "\nYou see here:\n  " + lipgloss.NewStyle().Bold(true).Render(strings.Join(otherChars, "\n  ")) + "\n"
		}

		// Room info at top with styling (only in output viewport, no stats)
		roomInfo := fmt.Sprintf("[%s]\n%s%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			characterList,
			m.formatExitsWithColor())

		// Show message if any (only once!)
		if m.message != "" && !messageRendered {
			roomInfo += "\n\n" + m.styledMessage(m.message)
			messageRendered = true
		}

		// Render output viewport with pink border (room info only, no stats)
		s.WriteString(outputStyle.Render(roomInfo))
		s.WriteString("\n")

		// Full-width status bar separator (middle ~10%)
		separatorStyle := lipgloss.NewStyle().
			Foreground(pink).
			Bold(true).
			Width(width)
		separatorLine := separatorStyle.Render(strings.Repeat("─", width-2))
		s.WriteString(separatorLine)
		s.WriteString("\n")
		// Stats go in the actual status bar (middle panel)
		s.WriteString(separatorStyle.Align(lipgloss.Center).Render(statsLine))
		s.WriteString("\n")
		s.WriteString(separatorLine)
		s.WriteString("\n")

		// Full-width input area (bottom ~20%)
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Padding(0, 1).
			Width(width - 2).
			Height(inputHeight - 2)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

	case ScreenFountainWake:
		s.WriteString(fountainWakeScreen())
		s.WriteString("\n\n")
		if m.message != "" && !messageRendered {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n\n")
			messageRendered = true
		}
		s.WriteString(promptStyle.Render("Press ENTER to continue..."))

	case ScreenFountainWash:
		s.WriteString(fountainWashScreen())
		s.WriteString("\n\n")
		if m.message != "" && !messageRendered {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n\n")
			messageRendered = true
		}
		s.WriteString(promptStyle.Render("Press ENTER to wash your face and remember who you are..."))

	case ScreenCharacterCreate:
		s.WriteString(characterCreateScreen())
		s.WriteString("\n\n")
		if m.message != "" && !messageRendered {
			s.WriteString(m.styledMessage(m.message))
			s.WriteString("\n\n")
			messageRendered = true
		}
		s.WriteString(promptStyle.Render("> "))
		s.WriteString(m.textInput.View())

		// ScreenPlaying uses full-width panels - don't center, just clear message and return
		m.message = ""
		m.messageType = ""
		return s.String()
	}

	// Center in terminal using proper visual width calculation
	// lipgloss.Width() correctly handles ANSI escape codes (they don't take visual space)
	if m.width > 0 && m.height > 0 && m.width > 60 {
		lines := strings.Split(s.String(), "\n")
		var centered []string
		for _, line := range lines {
			// Use lipgloss.Width for proper visual width (ignores ANSI codes)
			visualWidth := lipgloss.Width(line)
			padding := (m.width - visualWidth) / 2
			if padding > 0 && visualWidth < m.width-10 {
				centered = append(centered, fmt.Sprintf("%*s%s", padding, "", line))
			} else {
				centered = append(centered, line)
			}
		}
		// Clear message after rendering to prevent accumulation
		m.message = ""
		m.messageType = ""
		return strings.Join(centered, "\n")
	}

	// Clear message after rendering to prevent accumulation on next tick
	// This fixes the "jumbled text" issue during combat and other rapid updates
	m.message = ""
	m.messageType = ""

	return s.String()
}

// ============================================================
// STATIC SCREENS
// ============================================================

func welcomeScreen() string {
	// Use a flexible layout that works at various terminal widths
	// The ASCII art logo is designed to be ~60 chars wide, so we let it be
	// but wrap it in a bordered box that adapts
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(1, 2).
		MaxWidth(80). // Limit max width to prevent overflow on wide terminals
		Render(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║    ██████╗ ███████╗████████╗██████╗  ██████╗                ║
║    ██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗               ║
║    ██████╔╝█████╗     ██║   ██████╔╝██║   ██║               ║
║    ██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║               ║
║    ██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝               ║
║    ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝                ║
║                                                            ║
║           ██████╗  █████╗  ██████╗ ██████╗                  ║
║           ██╔══██╗██╔══██╗██╔════╝██╔═══██╗                 ║
║           ██████╔╝███████║██║     ██║   ██║                 ║
║           ██╔══██╗██╔══██║██║     ██║   ██║                 ║
║           ██║  ██║██║  ██║╚██████╗╚██████╔╝                 ║
║           ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝                  ║
║                                                            ║
║                    Welcome to Herbst MUD!                  ║
║                    Das Text-Adventure                      ║
║                                                            ║
╠════════════════════════════════════════════════════════════╣
║                                                            ║
║   1. Login      - Log in to your existing account         ║
║   2. Register   - Create a new player account             ║
║   3. Quit       - Exit the game                           ║
║                                                            ║
║   Use ↑/↓ or j/k to navigate, Enter to select             ║
║   Press ESC to go back                                    ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
`)
}

func loginScreen(width, height int) string {
	// Calculate dynamic dimensions
	boxWidth := 60
	if width > 70 {
		boxWidth = width - 20
	}
	if boxWidth > 100 {
		boxWidth = 100
	}

	verticalPadding := 2
	if height > 20 {
		verticalPadding = (height - 16) / 2
	}
	if verticalPadding > 10 {
		verticalPadding = 10
	}

	// Build dynamic login screen
	horizontalBorder := strings.Repeat("═", boxWidth-2)

	var sb strings.Builder
	// Top padding for vertical centering
	sb.WriteString(strings.Repeat("\n", verticalPadding))
	// Box
	sb.WriteString(lipgloss.NewStyle().
		Width(boxWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("75")).
		Padding(1, 2).
		Render(fmt.Sprintf(`
╔%s╗
║%s║
║%s║
║%s║
╚%s╝
`, horizontalBorder,
			strings.Repeat(" ", boxWidth-2),
			"                      LOGIN                          ",
			"  Enter your credentials to continue your adventure.  ",
			horizontalBorder)))

	return sb.String()
}

func registerScreen(width, height int) string {
	// Calculate dynamic dimensions
	boxWidth := 60
	if width > 70 {
		boxWidth = width - 20
	}
	if boxWidth > 100 {
		boxWidth = 100
	}

	verticalPadding := 2
	if height > 20 {
		verticalPadding = (height - 16) / 2
	}
	if verticalPadding > 10 {
		verticalPadding = 10
	}

	// Build dynamic register screen
	var sb strings.Builder
	// Top padding for vertical centering
	sb.WriteString(strings.Repeat("\n", verticalPadding))
	// Box - split screen style with purple border
	sb.WriteString(lipgloss.NewStyle().
		Width(boxWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(1, 2).
		Render("CREATE ACCOUNT\n\nCreate a new account to begin your adventure!\nPress ESC to go back to the main menu."))

	return sb.String()
}

// ============================================================
// FOUNTAIN CHARACTER CREATION SCREENS
// ============================================================

func fountainWakeScreen() string {
	return lipgloss.NewStyle().
		Foreground(green).
		Bold(true).
		MaxWidth(80).
		Render(`
╔══════════════════════════════════════════════════════════════════════╗
║                         ♨ THE FOUNTAIN ♨                             ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                      ║
║    You wake up at a murky fountain, covered in sticky mutant mud.    ║
║    The water glows faintly with an eerie green Ooze color.          ║
║    Your head throbs - you have no memory of how you got here.       ║
║    Something glints in the mud near your hand...                    ║
║                                                                      ║
║    The world around you is strange. Mutant weeds push through       ║
║    cracked cobblestones. The air smells of pizza and ooze.          ║
║                                                                      ║
║    You reach down and pick up the glinting object - a small        ║
║    copper coin with a turtle symbol on it.                          ║
║                                                                      ║
║    As you touch it, visions flash through your mind:               ║
║    → Mutant turtles trained by a wise rat master                    ║
║    → A city ruined by a strange Ooze                                ║
║    → Your own face, now covered in fur and scales...                ║
║                                                                      ║
║    You remember now. You ARE a turtle! And there's a whole          ║
║    world out there to explore. First, you need to wash up.          ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝
`)
}

func fountainWashScreen() string {
	return lipgloss.NewStyle().
		Foreground(cyan).
		Bold(true).
		MaxWidth(80).
		Render(`
╔══════════════════════════════════════════════════════════════════════╗
║                      ♨ WASHING AT THE FOUNTAIN ♨                     ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                      ║
║    You lean over the fountain and splash the cool, glowing water    ║
║    on your face. The mutant mud washes away, revealing your true     ║
║    form - a green turtle shell, scaly skin, and determined eyes.    ║
║                                                                      ║
║    As the mud clears, so do your memories...                        ║
║                                                                      ║
║    You are a Mutant Turtle, trained in the martial arts by your     ║
║    sensei, Splinter. The Great Mutagen Spill transformed you        ║
║    from a ordinary turtle into a thinking, speaking being.          ║
║                                                                      ║
║    Your sensei taught you well. You know:                           ║
║    → Ninjutsu - the way of the shadow warrior                       ║
║    → Survival - how to live in this post-Ooze world                ║
║    → Pizza - the most important food in existence                   ║
║                                                                      ║
║    The fountain water shows your reflection - you're ready for      ║
║    your next adventure!                                             ║
║                                                                      ║
║    A path leads north to the Crossroads.                            ║
║    To the east, you see signs for the "Canal District".            ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝
`)
}

func characterCreateScreen() string {
	return lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true).
		MaxWidth(80).
		Render(`
╔══════════════════════════════════════════════════════════════════════╗
║                    ✦ CHARACTER CREATION ✦                            ║
╠══════════════════════════════════════════════════════════════════════╣
║                                                                      ║
║    Welcome, young turtle! Time to create your hero!                  ║
║                                                                      ║
║    Available Options:                                                ║
║                                                                      ║

║    [1] Name     - What shall we call you?                           ║
║    [2] Race     - Human, Turtle, Rabbit, Rat, Rhino                  ║
║    [3] Gender   - Male, Female, Other                                ║
║    [4] Class    - Warrior, Chef, Mystic                              ║
║    [5] Size     - Small, Medium, Large (affects combat)              ║
║                                                                      ║
║    Type a number to select, or 'done' when finished.                ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝
`)
}

// ============================================================
// INPUT HANDLERS FOR FOUNTAIN FLOW
// ============================================================

func (m *model) handleFountainWakeInput(input string) {
	// Any input advances to the wash screen
	m.screen = ScreenFountainWash
	m.message = "You wash the mud from your face. The cool water feels refreshing. Your memories start to return..."
}

func (m *model) handleFountainWashInput(input string) {
	// Any input advances to character creation
	m.screen = ScreenCharacterCreate
	m.message = "Now to create your character!\n\nSelect an option (1-5) or type 'done' when finished."
}

func (m *model) handleCharacterCreateInput(input string) {
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "1", "name":
		m.message = "Enter your character name:"
		// Could switch to edit field mode for name input
	case "2", "race":
		m.message = "Select your race:\n1. Human - Balanced stats\n2. Turtle - High defense, low speed\n3. Rabbit - High speed, low defense\n4. Rat - High agility, stealthy\n5. Rhino - High strength, slow"
	case "3", "gender":
		m.message = "Select your gender:\n1. Male\n2. Female\n3. Other"
	case "4", "class":
		m.message = "Select your class:\n1. Warrior - Strong melee fighter\n2. Chef - Pizza-powered combat\n3. Mystic - Uses Ooze energy"
	case "5", "size":
		m.message = "Select your size:\n1. Small - Fast, less HP\n2. Medium - Balanced\n3. Large - Slow, more HP"
	case "done", "finished", "complete":
		m.message = "Character creation complete! Welcome to Herbst MUD!"
		m.screen = ScreenPlaying
	default:
		m.message = "Invalid choice. Select 1-5 or 'done' when finished."
	}
}

// fetchRoomCharacters fetches characters (NPCs and players) in the current room
func (m *model) fetchRoomCharacters() {
	url := fmt.Sprintf("%s/rooms/%d/characters", RESTAPIBase, m.currentRoom)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching room characters: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching room characters: status %d", resp.StatusCode)
		return
	}

	var characters []roomCharacter
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		log.Printf("Error decoding room characters: %v", err)
		return
	}

	m.roomCharacters = characters
}

// formatRoomCharacters returns a formatted string of characters in the room with color coding
// NPCs show in RED, players show in GREEN
func (m *model) formatRoomCharacters() string {
	if len(m.roomCharacters) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "\nYou see here:")

	for _, char := range m.roomCharacters {
		// Skip the current player (their own character)
		if char.ID == m.currentCharacterID {
			continue
		}

		var styledName string
		if char.IsNPC {
			// NPCs in RED
			styledName = lipgloss.NewStyle().Foreground(red).Render(char.Name)
		} else {
			// Players in GREEN
			styledName = lipgloss.NewStyle().Foreground(green).Render(char.Name)
		}

		// Add level and class info
		info := fmt.Sprintf("  - %s (Lv.%d %s %s)", styledName, char.Level, char.Race, char.Class)
		lines = append(lines, info)
	}

	return strings.Join(lines, "\n")
}
