package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	ScreenWelcome   = "welcome"
	ScreenLogin     = "login"
	ScreenRegister  = "register"
	ScreenPlaying   = "playing"
	ScreenProfile   = "profile"
	ScreenEditField = "edit_field"
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

		// Initialize Gizmo NPC and fountain room
		if err := dbinit.InitGizmo(client); err != nil {
			log.Printf("Warning: failed to initialize Gizmo: %v", err)
		}

		// Initialize starter weapons
		if err := dbinit.InitWeapons(client); err != nil {
			log.Printf("Warning: failed to initialize weapons: %v", err)
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
	visitedRooms map[int]bool
	knownExits   map[string]bool // For color-coded exits

	// Room items (GitHub #89 - Item system)
	roomItems []RoomItem

	// Room characters (GitHub #145 - Look command room display)
	roomCharacters []roomCharacter

	// Debug mode - shows room ID in status bar
	debugMode bool
}

// RoomItem represents an item in a room for display
type RoomItem struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	ExamineDesc    string         `json:"examineDesc"`
	HiddenDetails  []HiddenDetail `json:"hiddenDetails"`
	HiddenThreshold int           `json:"hiddenThreshold"`
	IsImmovable    bool           `json:"isImmovable"`
	Color          string         `json:"color"`
	IsVisible      bool           `json:"isVisible"`
	ItemType       string         `json:"itemType"`
	Weight         int            `json:"weight"`
	ItemDamage     int            `json:"itemDamage"`
	ItemDurability int            `json:"itemDurability"`
	// Container fields
	IsContainer bool   `json:"isContainer,omitempty"`
	Capacity    int    `json:"capacity,omitempty"`
	IsLocked    bool   `json:"isLocked,omitempty"`
	// Weapon fields
	MinDamage        int    `json:"minDamage,omitempty"`
	MaxDamage        int    `json:"maxDamage,omitempty"`
	WeaponType       string `json:"weaponType,omitempty"`
	ClassRestriction string `json:"classRestriction,omitempty"`
}

// roomCharacter represents a character (NPC or player) in a room for display
type roomCharacter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsNPC    bool   `json:"isNPC"`
	Level    int    `json:"level"`
	Class    string `json:"class"`
	Race     string `json:"race"`
	UserID   int    `json:"userId"`
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
	purple = lipgloss.Color("141")
	white  = lipgloss.Color("15")
	gray   = lipgloss.Color("8")
	pink   = lipgloss.Color("219")
	cyan   = lipgloss.Color("51")

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
	return styleMessage(msg, m.messageType)
}

// styleMessage returns a styled message based on message type (standalone helper)
func styleMessage(msg string, msgType string) string {
	if msg == "" {
		return ""
	}

	switch msgType {
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
}

// ============================================================
// ROOM ITEM DISPLAY (GitHub #89 - Item system)
// ============================================================

// Item colors for display
var (
	itemColorGold   = lipgloss.Color("220") // Gold for immovable items
	itemColorWeapon = lipgloss.Color("196") // Red for weapons
	itemColorArmor  = lipgloss.Color("75")  // Blue for armor
	itemColorMisc   = lipgloss.Color("242") // Gray for misc
)

// formatRoomItems returns a formatted string of items in the room
func (m *model) formatRoomItems() string {
	if len(m.roomItems) == 0 {
		return ""
	}

	var items []string
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue // Skip invisible items
		}

		var style lipgloss.Style
		if item.IsImmovable {
			// Immobile items get gold color by default or custom color
			if item.Color != "" {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color))
			} else {
				style = lipgloss.NewStyle().Foreground(itemColorGold)
			}
		} else {
			// Regular items get color based on type
			switch item.ItemType {
			case "weapon":
				style = lipgloss.NewStyle().Foreground(itemColorWeapon)
			case "armor":
				style = lipgloss.NewStyle().Foreground(itemColorArmor)
			default:
				style = lipgloss.NewStyle().Foreground(itemColorMisc)
			}
		}

		// Add special marker for immovable items
		if item.IsImmovable {
			items = append(items, style.Render("⬥ "+item.Name))
		} else {
			items = append(items, style.Render(item.Name))
		}
	}

	if len(items) == 0 {
		return ""
	}
	return "\n\nYou see: " + strings.Join(items, ", ")
}

// loadRoomItems fetches items for the current room from the API
func (m *model) loadRoomItems() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment", RESTAPIBase, m.currentRoom))
	if err != nil {
		log.Printf("Error fetching room items: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var items []RoomItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		log.Printf("Error decoding room items: %v", err)
		return
	}

	m.roomItems = items
}

// loadRoomCharacters fetches characters (NPCs and players) in the current room from the API
func (m *model) loadRoomCharacters() {
	if m.currentRoom == 0 {
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/characters", RESTAPIBase, m.currentRoom))
	if err != nil {
		log.Printf("Error fetching room characters: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var characters []roomCharacter
	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		log.Printf("Error decoding room characters: %v", err)
		return
	}

	m.roomCharacters = characters
}

// formatRoomCharacters returns a formatted string of characters (NPCs and players) in the room
func (m *model) formatRoomCharacters() string {
	if len(m.roomCharacters) == 0 {
		return ""
	}

	var npcs []string
	var players []string

	for _, char := range m.roomCharacters {
		if char.IsNPC {
			// NPCs in red
			style := lipgloss.NewStyle().Foreground(red)
			npcs = append(npcs, style.Render(char.Name))
		} else {
			// Players in green
			style := lipgloss.NewStyle().Foreground(green)
			players = append(players, style.Render(char.Name))
		}
	}

	var parts []string
	if len(npcs) > 0 {
		parts = append(parts, "NPCs: "+strings.Join(npcs, ", "))
	}
	if len(players) > 0 {
		parts = append(parts, "Players: "+strings.Join(players, ", "))
	}

	if len(parts) == 0 {
		return ""
	}
	return "\n\n" + strings.Join(parts, " | ")
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
  look/l - Look around (shows items)
  exits/x - Show exits  
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  open <container> - Open a container
  take <item> from <container> - Take item from container
  put <item> in <container> - Put item in container
  inventory/i - Show your inventory
  whoami - Show your info
  profile/p - Edit character profile
  clear/cls - Clear screen
  quit - Exit game`
		m.messageType = "info"
	case "look", "l":
		m.loadRoomItems()
		m.loadRoomCharacters()
		m.message = fmt.Sprintf("[%s]\n%s\n\nExits: %s%s%s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor(),
			m.formatRoomItems(),
			m.formatRoomCharacters())
		m.messageType = "info"
	case "exits", "x":
		m.message = fmt.Sprintf("Exits: %s", m.formatExitsWithColor())
		m.messageType = "info"
	case "examine", "ex", "inspect":
		m.handleExamineCommand(cmd)
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
	case "debug":
		m.handleDebugCommand(cmd)
		return
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
		// Check for take/get command (GitHub #89 - Item system)
		if strings.HasPrefix(cmd, "take ") || strings.HasPrefix(cmd, "get ") {
			m.handleTakeCommand(cmd)
			return
		}
		// Check for drop command
		if strings.HasPrefix(cmd, "drop ") {
			m.handleDropCommand(cmd)
			return
		}
		// Check for inventory command
		if cmd == "inventory" || cmd == "i" || cmd == "inv" {
			m.handleInventoryCommand()
			return
		}
		// Container commands (GitHub #143)
		if strings.HasPrefix(cmd, "open ") {
			m.handleOpenContainerCommand(cmd)
			return
		}
		if strings.HasPrefix(cmd, "take ") && strings.Contains(cmd, " from ") {
			m.handleTakeFromContainerCommand(cmd)
			return
		}
		if strings.HasPrefix(cmd, "put ") && strings.Contains(cmd, " in ") {
			m.handlePutInContainerCommand(cmd)
			return
		}
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

		// Load items and characters for the new room
		m.loadRoomItems()
		m.loadRoomCharacters()

		// Mark new room as visited
		wasVisited := m.visitedRooms[m.currentRoom]
		m.visitedRooms[m.currentRoom] = true

		// Mark new exits as known
		for dir := range m.exits {
			m.knownExits[dir] = true
		}

		// Format room display with items and characters
		roomDisplay := fmt.Sprintf("\n\nExits: %s%s%s",
			m.formatExitsWithColor(),
			m.formatRoomItems(),
			m.formatRoomCharacters())

		if wasVisited {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				roomDisplay)
		} else {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				roomDisplay)
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

func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		// Show current debug status
		if m.debugMode {
			m.message = "Debug mode: ON (Room ID visible in status bar)"
		} else {
			m.message = "Debug mode: OFF\nUsage: debug on | debug off"
		}
		m.messageType = "info"
		return
	}

	subCmd := parts[1]
	switch subCmd {
	case "on", "true", "1", "yes":
		m.debugMode = true
		m.message = "Debug mode: ON (Room ID will show in status bar)"
		m.messageType = "success"
	case "off", "false", "0", "no":
		m.debugMode = false
		m.message = "Debug mode: OFF"
		m.messageType = "info"
	default:
		m.message = "Usage: debug on | debug off"
		m.messageType = "error"
	}
}

// ============================================================
// ITEM COMMANDS (GitHub #89 - Item system)
// ============================================================

// handleTakeCommand handles the take/get command
func (m *model) handleTakeCommand(cmd string) {
	// Extract item name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Take what? Usage: take <item name>"
		m.messageType = "error"
		return
	}
	itemName := strings.Join(parts[1:], " ")

	// Load room items to find the item
	m.loadRoomItems()

	// Find item by name (case-insensitive partial match)
	var targetItem *RoomItem
	for i := range m.roomItems {
		if strings.Contains(strings.ToLower(m.roomItems[i].Name), strings.ToLower(itemName)) {
			targetItem = &m.roomItems[i]
			break
		}
	}

	if targetItem == nil {
		m.message = fmt.Sprintf("You don't see any %s here.", itemName)
		m.messageType = "error"
		return
	}

	// Check if immovable
	if targetItem.IsImmovable {
		var colorStyle lipgloss.Style
		if targetItem.Color != "" {
			colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(targetItem.Color))
		} else {
			colorStyle = lipgloss.NewStyle().Foreground(itemColorGold)
		}
		m.message = fmt.Sprintf("You can't take the %s. It's firmly fixed in place.", colorStyle.Render(targetItem.Name))
		m.messageType = "error"
		return
	}

	// Take the item - move it to player's inventory (roomId = 0 or null)
	url := fmt.Sprintf("%s/equipment/%d", RESTAPIBase, targetItem.ID)
	jsonData, _ := json.Marshal(map[string]interface{}{"roomId": nil})
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		m.message = fmt.Sprintf("Error picking up item: %v", err)
		m.messageType = "error"
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.message = fmt.Sprintf("Error picking up item: %v", err)
		m.messageType = "error"
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.message = fmt.Sprintf("Failed to pick up %s.", targetItem.Name)
		m.messageType = "error"
		return
	}

	m.message = fmt.Sprintf("You pick up the %s.", targetItem.Name)
	m.messageType = "success"
}

// handleDropCommand handles the drop command
func (m *model) handleDropCommand(cmd string) {
	// Extract item name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Drop what? Usage: drop <item name>"
		m.messageType = "error"
		return
	}
	itemName := strings.Join(parts[1:], " ")

	// For now, show a message that inventory is not fully implemented
	// This would need player inventory tracking
	m.message = fmt.Sprintf("You don't have any %s to drop.", itemName)
	m.messageType = "error"
}

// handleExamineCommand handles the examine/ex/inspect/i command
func (m *model) handleExamineCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Examine what? Usage: examine <item>"
		m.messageType = "error"
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(target)

	// First check room items
	for _, item := range m.roomItems {
		if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Then check inventory
	resp, err := http.Get(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err == nil {
		defer resp.Body.Close()
		var items []RoomItem
		if json.NewDecoder(resp.Body).Decode(&items) == nil {
			for _, item := range items {
				if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
					m.displayItemDetails(item)
					return
				}
			}
		}
	}

	// Check if it's an NPC
	if m.currentRoom > 0 {
		resp, err := http.Get(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			var npcs []struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Level       int    `json:"level"`
				Disposition string `json:"disposition"`
			}
			if json.NewDecoder(resp.Body).Decode(&npcs) == nil {
				for _, npc := range npcs {
					if strings.Contains(strings.ToLower(npc.Name), target) || strings.ToLower(npc.Name) == target {
						m.message = fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
							npc.Name, npc.Description, npc.Level, npc.Disposition)
						m.messageType = "info"
						return
					}
				}
			}
		}
	}

	m.message = fmt.Sprintf("You don't see '%s' here.", target)
	m.messageType = "error"
}

// displayItemDetails shows detailed info about an item
func (m *model) displayItemDetails(item RoomItem) {
	var details strings.Builder

	// Title with color if applicable
	if item.Color != "" {
		details.WriteString(fmt.Sprintf("[%s]\n", item.Name))
	} else {
		details.WriteString(fmt.Sprintf("[%s]\n", item.Name))
	}

	// Use examine description if available, otherwise fall back to description
	desc := item.ExamineDesc
	if desc == "" {
		desc = item.Description
	}
	details.WriteString(desc + "\n")

	// Show container info if applicable (GitHub #143)
	if item.IsContainer {
		details.WriteString("\n--- Container ---\n")
		details.WriteString(fmt.Sprintf("  Capacity: %d items\n", item.Capacity))
		if item.IsLocked {
			details.WriteString("  Status: 🔒 Locked\n")
		} else {
			details.WriteString("  Status: 🔓 Unlocked\n")
		}
	}

	// Show stats if it's equipment
	if item.ItemType == "weapon" || item.ItemType == "armor" {
		details.WriteString("\n--- Stats ---\n")
		if item.Weight > 0 {
			details.WriteString(fmt.Sprintf("  Weight: %d\n", item.Weight))
		}
		if item.ItemDamage > 0 {
			details.WriteString(fmt.Sprintf("  Damage: %d\n", item.ItemDamage))
		}
		if item.ItemDurability > 0 {
			details.WriteString(fmt.Sprintf("  Durability: %d\n", item.ItemDurability))
		}
		details.WriteString(fmt.Sprintf("  Type: %s\n", item.ItemType))
	}

	// Show hidden details if player has high enough examine skill
	// For now, we'll show all hidden details (skill check deferred)
	if len(item.HiddenDetails) > 0 && item.HiddenThreshold > 0 {
		// TODO: Fetch player's examine skill and compare to threshold
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	} else if len(item.HiddenDetails) > 0 {
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	}

	m.message = details.String()
	m.messageType = "info"
}

// handleInventoryCommand handles the inventory/i command
func (m *model) handleInventoryCommand() {
	// Fetch player's inventory from API
	// For now, show a placeholder message
	resp, err := http.Get(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.message = fmt.Sprintf("Error fetching inventory: %v", err)
		m.messageType = "error"
		return
	}
	defer resp.Body.Close()

	// Parse inventory items
	var items []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ItemType    string `json:"itemType"`
		IsEquipped  bool   `json:"isEquipped"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		m.message = "You aren't carrying anything."
		m.messageType = "info"
		return
	}

	if len(items) == 0 {
		m.message = "You aren't carrying anything."
		m.messageType = "info"
		return
	}

	// Format inventory display
	var inv strings.Builder
	inv.WriteString("=== INVENTORY ===\n\n")
	for _, item := range items {
		equipped := ""
		if item.IsEquipped {
			equipped = " [equipped]"
		}
		inv.WriteString(fmt.Sprintf("  %s%s\n", item.Name, equipped))
		if item.Description != "" {
			inv.WriteString(fmt.Sprintf("    %s\n", item.Description))
		}
	}
	m.message = inv.String()
	m.messageType = "info"
}

// handleOpenContainerCommand handles the open <container> command
func (m *model) handleOpenContainerCommand(cmd string) {
	// Extract container name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.message = "Open what? Usage: open <container name>"
		m.messageType = "error"
		return
	}
	containerName := strings.Join(parts[1:], " ")

	// Load room items to find the container
	m.loadRoomItems()

	// Find container by name (case-insensitive partial match)
	var targetContainer *RoomItem
	for i := range m.roomItems {
		if m.roomItems[i].IsContainer && strings.Contains(strings.ToLower(m.roomItems[i].Name), strings.ToLower(containerName)) {
			targetContainer = &m.roomItems[i]
			break
		}
	}

	if targetContainer == nil {
		m.message = fmt.Sprintf("You don't see any %s here.", containerName)
		m.messageType = "error"
		return
	}

	// Check if it's actually a container
	if !targetContainer.IsContainer {
		m.message = fmt.Sprintf("You can't open the %s.", targetContainer.Name)
		m.messageType = "error"
		return
	}

	// Check if locked
	if targetContainer.IsLocked {
		m.message = fmt.Sprintf("The %s is locked.", targetContainer.Name)
		m.messageType = "error"
		return
	}

	// TODO: Fetch container contents from API when container inventory is implemented
	// For now, show basic container info
	m.message = fmt.Sprintf("You open the %s.\n\nIt appears to be empty for now.", targetContainer.Name)
	m.messageType = "info"
}

// handleTakeFromContainerCommand handles the take <item> from <container> command
func (m *model) handleTakeFromContainerCommand(cmd string) {
	// Parse: take <item> from <container>
	fromIdx := strings.Index(cmd, " from ")
	if fromIdx == -1 {
		m.message = "Usage: take <item> from <container>"
		m.messageType = "error"
		return
	}

	itemName := strings.TrimSpace(cmd[5:fromIdx])
	containerName := strings.TrimSpace(cmd[fromIdx+6:])

	if itemName == "" || containerName == "" {
		m.message = "Usage: take <item> from <container>"
		m.messageType = "error"
		return
	}

	// Load room items to find the container
	m.loadRoomItems()

	// Find container
	var targetContainer *RoomItem
	for i := range m.roomItems {
		if m.roomItems[i].IsContainer && strings.Contains(strings.ToLower(m.roomItems[i].Name), strings.ToLower(containerName)) {
			targetContainer = &m.roomItems[i]
			break
		}
	}

	if targetContainer == nil {
		m.message = fmt.Sprintf("You don't see any %s here.", containerName)
		m.messageType = "error"
		return
	}

	// Check if container is locked
	if targetContainer.IsLocked {
		m.message = fmt.Sprintf("The %s is locked.", targetContainer.Name)
		m.messageType = "error"
		return
	}

	// TODO: Verify item is actually in container when container inventory is implemented
	m.message = fmt.Sprintf("You take the %s from the %s.", itemName, targetContainer.Name)
	m.messageType = "success"
}

// handlePutInContainerCommand handles the put <item> in <container> command
func (m *model) handlePutInContainerCommand(cmd string) {
	// Parse: put <item> in <container>
	inIdx := strings.Index(cmd, " in ")
	if inIdx == -1 {
		m.message = "Usage: put <item> in <container>"
		m.messageType = "error"
		return
	}

	itemName := strings.TrimSpace(cmd[4:inIdx])
	containerName := strings.TrimSpace(cmd[inIdx+4:])

	if itemName == "" || containerName == "" {
		m.message = "Usage: put <item> in <container>"
		m.messageType = "error"
		return
	}

	// Load room items to find the container
	m.loadRoomItems()

	// Find container
	var targetContainer *RoomItem
	for i := range m.roomItems {
		if m.roomItems[i].IsContainer && strings.Contains(strings.ToLower(m.roomItems[i].Name), strings.ToLower(containerName)) {
			targetContainer = &m.roomItems[i]
			break
		}
	}

	if targetContainer == nil {
		m.message = fmt.Sprintf("You don't see any %s here.", containerName)
		m.messageType = "error"
		return
	}

	// Check if container is locked
	if targetContainer.IsLocked {
		m.message = fmt.Sprintf("The %s is locked.", targetContainer.Name)
		m.messageType = "error"
		return
	}

	// TODO: Implement actual put in container when container inventory is implemented
	m.message = fmt.Sprintf("You put the %s in the %s.", itemName, targetContainer.Name)
	m.messageType = "success"
}

func (m *model) loadOrCreateCharacter() {
	// Query for existing character for this user
	ctx := context.Background()
	chars, err := m.client.Character.Query().Where(character.HasUserWith(user.IDEQ(m.currentUserID))).All(ctx)
	if err != nil || len(chars) == 0 {
		// No character yet - use defaults
		m.currentCharacterName = m.currentUserName
		m.characterGender = "unspecified"
		m.characterDescription = "A mysterious figure."
		// Default stats for new characters (Level 1)
		m.characterHP = 100
		m.characterMaxHP = 100
		m.characterStamina = 50
		m.characterMaxStamina = 50
		m.characterMana = 25
		m.characterMaxMana = 25
		m.characterLevel = 1
		m.characterExperience = 0
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
		// Build input with menu
		var inputContent strings.Builder
		inputContent.WriteString(promptStyle.Render("> "))
		inputContent.WriteString(m.textInput.View())
		inputContent.WriteString("\n\n")
		inputContent.WriteString(lipgloss.NewStyle().Foreground(gray).Render("Press 1/2/3 or type login/register/quit"))
		return welcomeScreen(m.width, m.height, inputContent.String())

	case ScreenLogin:
		// Build input prompt
		promptText := "> "
		if m.inputField == "username" {
			promptText = "Username: "
		} else if m.inputField == "password" {
			promptText = "Password: "
		}
		inputContent := promptStyle.Render(promptText) + m.textInput.View()
		return loginScreen(m.width, m.height, m.message, m.messageType, inputContent)

	case ScreenRegister:
		// Build input prompt
		promptText := "> "
		if m.inputField == "username" {
			promptText = "Username: "
		} else if m.inputField == "password" {
			promptText = "Password: "
		}
		inputContent := promptStyle.Render(promptText) + m.textInput.View()
		return registerScreen(m.width, m.height, m.message, m.messageType, inputContent)

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

		// Debug mode - show room ID when enabled
		debugInfo := ""
		if m.debugMode {
			debugInfo = " " + lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("[Room: %d]", m.currentRoom))
		}

		// Room info at top with styling (only in output viewport, no stats)
		roomInfo := fmt.Sprintf("[%s]\n%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor())

		// Show message if any
		if m.message != "" {
			roomInfo += "\n\n" + m.styledMessage(m.message)
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
		s.WriteString(separatorStyle.Align(lipgloss.Center).Render(statsLine + debugInfo))
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

		// ScreenPlaying uses full-width panels - don't center, just clear message and return
		m.message = ""
		m.messageType = ""
		return s.String()
	}

	// Center in terminal (optional - can be disabled if causing issues)
	// Use lipgloss.Width() to correctly handle ANSI escape codes (fixes issue #75)
	if m.width > 0 && m.height > 0 && m.width > 60 {
		lines := strings.Split(s.String(), "\n")
		var centered []string
		for _, line := range lines {
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

func welcomeScreen(width, height int, inputView string) string {
	// Calculate proportional heights matching game screen
	// Output: ~70%, Input: ~30%
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	// Output pane (top) - shows logo and menu
	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(outputHeight - 2)

	// Build output content - lipgloss adds the border, so just content here
	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("        Welcome Adventurer!        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  1. Login"))
	outputContent.WriteString("      - Log in to your existing account\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  2. Register"))
	outputContent.WriteString("   - Create a new character\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(cyan).Render("  3. Quit"))
	outputContent.WriteString("       - Exit the game\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Foreground(gray).Render("  Use arrow keys or type number/command"))

	// Input pane (bottom)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(inputHeight - 2)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}

func loginScreen(width, height int, message, messageType string, inputView string) string {
	// Calculate proportional heights matching game screen
	// Output: ~70%, Input: ~30%
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	// Output pane (top) - shows logo and prompts
	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(outputHeight - 2)

	// Build output content - lipgloss adds the border, so just content here
	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("            LOGIN            "))
	outputContent.WriteString("\n\n")

	// Show message/prompt
	if message != "" {
		outputContent.WriteString(styleMessage(message, messageType))
		outputContent.WriteString("\n")
	}

	// Input pane (bottom)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(inputHeight - 2)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}

func registerScreen(width, height int, message, messageType string, inputView string) string {
	// Calculate proportional heights matching game screen
	// Output: ~70%, Input: ~30%
	inputHeight := height * 30 / 100
	if inputHeight < 5 {
		inputHeight = 5
	}
	outputHeight := height - inputHeight
	if outputHeight < 10 {
		outputHeight = 10
	}

	// Output pane (top) - shows logo and prompts
	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(outputHeight - 2)

	// Build output content - lipgloss adds the border, so just content here
	var outputContent strings.Builder
	outputContent.WriteString("\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(green).Render("        🐢 HERBST MUD 🐢        "))
	outputContent.WriteString("\n\n")
	outputContent.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("        CREATE ACCOUNT        "))
	outputContent.WriteString("\n\n")

	// Show message/prompt
	if message != "" {
		outputContent.WriteString(styleMessage(message, messageType))
		outputContent.WriteString("\n")
	}

	// Input pane (bottom)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(0, 1).
		Width(width - 2).
		Height(inputHeight - 2)

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}
