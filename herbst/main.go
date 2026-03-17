package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
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

		// Initialize fountain item in starting room
		if err := dbinit.InitFountainItem(client); err != nil {
			log.Printf("Warning: failed to initialize fountain item: %v", err)
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
							connectedAt:     time.Now(),
							session:         s,
							client:          client,
							screen:          ScreenWelcome,
							currentRoom:     StartingRoomID,
							textInput:       ti,
							spinner:         sp,
							visitedRooms:    make(map[int]bool),
							knownExits:      make(map[string]bool),
							roomCharacters:  make([]roomCharacter, 0),
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
	roomItems      []roomItemDisplay
	roomCharacters []roomCharacter // Characters in current room (NPCs + players)

	// Debug mode - shows room ID in status bar
	debugMode bool
}

// roomCharacter holds info about characters in a room for display
type roomCharacter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsNPC  bool   `json:"isNPC"`
	Level  int    `json:"level"`
	Class  string `json:"class"`
	Race   string `json:"race"`
	UserID int    `json:"userId,omitempty"`
}

// roomItemDisplay holds display info for items in a room
type roomItemDisplay struct {
	name        string
	description string
	color       lipgloss.Color // Display color (gold for immovable, custom otherwise)
	isImmovable bool
	isVisible   bool
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
// IMPORTANT: This function must be PURE - it should NOT modify model state
// State changes (like marking exits as known) should happen in Update(), not View()
func (m *model) formatExitsWithColor() string {
	if len(m.exits) == 0 {
		return lipgloss.NewStyle().Foreground(gray).Render("none")
	}

	// Sort exits for consistent order (Go maps are unordered)
	var dirs []string
	for dir := range m.exits {
		dirs = append(dirs, dir)
	}
	// Sort alphabetically for consistent display
	sort.Strings(dirs)

	var result []string
	for _, dir := range dirs {
		roomID := m.exits[dir]
		var exitStyle lipgloss.Style
		if m.visitedRooms[roomID] {
			// Green = visited
			exitStyle = lipgloss.NewStyle().Foreground(exitVisitedColor)
		} else if m.knownExits[dir] {
			// Yellow = known but not visited
			exitStyle = lipgloss.NewStyle().Foreground(exitKnownColor)
		} else {
			// White = new (never seen before)
			// NOTE: We do NOT mark as known here - that's done in Update()
			exitStyle = lipgloss.NewStyle().Foreground(exitNewColor)
		}
		result = append(result, exitStyle.Render(dir))
	}

	return strings.Join(result, ", ")
}

// formatRoomItemsWithColor returns a formatted string of items in the room
func (m *model) formatRoomItemsWithColor() string {
	if len(m.roomItems) == 0 {
		return ""
	}

	var items []string
	for _, item := range m.roomItems {
		if !item.isVisible {
			continue
		}
		itemStyle := lipgloss.NewStyle().Foreground(item.color)
		if item.isImmovable {
			// Immovable items shown in gold with "(fixed)" indicator
			items = append(items, itemStyle.Render(item.name+" (fixed)")+": "+item.description)
		} else {
			items = append(items, itemStyle.Render(item.name)+": "+item.description)
		}
	}

	if len(items) == 0 {
		return ""
	}

	return lipgloss.NewStyle().Bold(true).Foreground(purple).Render("Items:") + "\n  " + strings.Join(items, "\n  ")
}

// formatRoomCharactersWithColor returns a formatted string of characters in the room
// NPCs are shown in RED, players are shown in GREEN
func (m *model) formatRoomCharactersWithColor() string {
	if len(m.roomCharacters) == 0 {
		return ""
	}

	var chars []string
	for _, rc := range m.roomCharacters {
		var charStyle lipgloss.Style
		var typeLabel string
		if rc.IsNPC {
			// NPCs in RED
			charStyle = lipgloss.NewStyle().Foreground(red).Bold(true)
			typeLabel = "NPC"
		} else {
			// Players in GREEN
			charStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
			typeLabel = "Player"
		}

		// Format: "Name (NPC) - Level 5 Orc Warrior" or just "Name (Player)"
		var details []string
		if rc.Level > 0 {
			details = append(details, fmt.Sprintf("Level %d", rc.Level))
		}
		if rc.Race != "" {
			details = append(details, rc.Race)
		}
		if rc.Class != "" {
			details = append(details, rc.Class)
		}

		var detailStr string
		if len(details) > 0 {
			detailStr = " - " + strings.Join(details, " ")
		}

		chars = append(chars, charStyle.Render(rc.Name)+fmt.Sprintf(" (%s)%s", typeLabel, detailStr))
	}

	if len(chars) == 0 {
		return ""
	}

	return lipgloss.NewStyle().Bold(true).Foreground(yellow).Render("You see:") + "\n  • " + strings.Join(chars, "\n  • ")
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
		m.textInput.EchoMode = textinput.EchoNormal // Reset to normal echo
		m.textInput.EchoCharacter = 0 // Reset echo character
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
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("") // Clear any previous input
		m.textInput.Placeholder = "Enter your username..."
		m.textInput.EchoMode = textinput.EchoNormal // Ensure normal echo for username
		m.textInput.EchoCharacter = 0 // Reset echo character
		m.textInput.Focus()
	case "2", "register", "r", "create":
		m.screen = ScreenRegister
		m.inputField = "username"
		m.loginUsername = ""
		m.loginPassword = ""
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("") // Clear any previous input
		m.textInput.Placeholder = "Choose a username..."
		m.textInput.EchoMode = textinput.EchoNormal // Ensure normal echo for username
		m.textInput.EchoCharacter = 0 // Reset echo character
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
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("") // Clear previous input
		m.textInput.Placeholder = "Enter your password..."
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.EchoCharacter = '•' // Use bullet character for masking
		m.textInput.Focus()
	} else if m.inputField == "password" {
		m.loginPassword = input
		m.textInput.SetValue("")
		m.textInput.EchoMode = textinput.EchoNormal
		m.textInput.EchoCharacter = 0 // Reset to default
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
		m.textInput.SetValue("")
		m.textInput.EchoMode = textinput.EchoNormal
		m.textInput.EchoCharacter = 0
		m.textInput.Placeholder = "Enter your username..."
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
	m.textInput.EchoMode = textinput.EchoNormal // Reset to normal echo
	m.textInput.EchoCharacter = 0 // Reset echo character
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

		// Load room items
		m.loadRoomItems()
		// Load room characters
		m.loadRoomCharacters()
	}
}

func (m *model) loadRoomItems() {
	if m.client == nil {
		return
	}
	// Get the room first
	room, err := m.client.Room.Get(context.Background(), m.currentRoom)
	if err != nil || room == nil {
		return // Silently fail - items are optional
	}
	items, err := m.client.Room.QueryEquipment(room).All(context.Background())
	if err != nil {
		return // Silently fail - items are optional
	}
	m.roomItems = make([]roomItemDisplay, 0, len(items))
	for _, item := range items {
		if !item.IsVisible {
			continue
		}
		display := roomItemDisplay{
			name:        item.Name,
			description: item.Description,
			isImmovable: item.IsImmovable,
			isVisible:   item.IsVisible,
		}
		// Set color: gold for immovable, custom color if set, white otherwise
		if item.IsImmovable {
			display.color = lipgloss.Color("220") // Gold
		} else if item.Color != "" {
			display.color = lipgloss.Color(item.Color)
		} else {
			display.color = lipgloss.Color("15") // White
		}
		m.roomItems = append(m.roomItems, display)
	}
}

// loadRoomCharacters fetches characters from the server for the current room
func (m *model) loadRoomCharacters() {
	if m.currentRoom == 0 {
		return
	}

	url := fmt.Sprintf("%s/rooms/%d/characters", RESTAPIBase, m.currentRoom)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error loading room characters: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to load room characters: status %d", resp.StatusCode)
		return
	}

	var chars []roomCharacter
	if err := json.NewDecoder(resp.Body).Decode(&chars); err != nil {
		log.Printf("Error decoding room characters: %v", err)
		return
	}

	// Filter out current player
	m.roomCharacters = make([]roomCharacter, 0, len(chars))
	for _, rc := range chars {
		if rc.ID != m.currentCharacterID {
			m.roomCharacters = append(m.roomCharacters, rc)
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
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("") // Clear previous input
		m.textInput.Placeholder = "Choose a password..."
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.EchoCharacter = '•' // Use bullet character for masking
		m.textInput.Focus()
	} else if m.inputField == "password" {
		if input == "" {
			m.message = "Password cannot be empty. Try again:"
			m.messageType = "error"
			return
		}
		m.loginPassword = input
		m.inputField = "confirm_password"
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("") // Clear previous input
		m.textInput.Placeholder = "Confirm your password..."
		m.textInput.Focus()
	} else if m.inputField == "confirm_password" {
		if input != m.loginPassword {
			m.message = "Passwords do not match. Try again:"
			m.messageType = "error"
			m.inputField = "password"
			m.loginPassword = ""
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Choose a password..."
			m.textInput.EchoMode = textinput.EchoPassword
			m.textInput.EchoCharacter = '•'
			m.textInput.Focus()
			return
		}
		m.inputField = "email"
		m.message = "" // Clear message - View() will show field label
		m.messageType = ""
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Enter your email..."
		m.textInput.EchoMode = textinput.EchoNormal
		m.textInput.EchoCharacter = 0 // Reset to default
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
			m.textInput.SetValue("")
			m.textInput.EchoMode = textinput.EchoNormal
			m.textInput.EchoCharacter = 0
			m.textInput.Placeholder = "Choose a username..."
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
	m.textInput.EchoMode = textinput.EchoNormal // Reset to normal echo
	m.textInput.EchoCharacter = 0 // Reset echo character
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

	// Load room items
	m.loadRoomItems()
}

// ============================================================
// GAME COMMAND PROCESSING
// ============================================================

func (m *model) processCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))
	args := cmd // For commands that need the full string with args

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
		// Format items and characters
		itemStr := m.formatRoomItemsWithColor()
		charStr := m.formatRoomCharactersWithColor()

		var extraParts []string
		if itemStr != "" {
			extraParts = append(extraParts, itemStr)
		}
		if charStr != "" {
			extraParts = append(extraParts, charStr)
		}
		extraPart := ""
		if len(extraParts) > 0 {
			extraPart = "\n\n" + strings.Join(extraParts, "\n\n")
		}

		m.message = fmt.Sprintf("[%s]\n%s%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			extraPart,
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
	case "debug":
		m.handleDebugCommand(args)
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

		// Load room items
		m.loadRoomItems()
		// Load room characters
		m.loadRoomCharacters()

		// Format items and characters
		itemStr := m.formatRoomItemsWithColor()
		charStr := m.formatRoomCharactersWithColor()

		var extraParts []string
		if itemStr != "" {
			extraParts = append(extraParts, itemStr)
		}
		if charStr != "" {
			extraParts = append(extraParts, charStr)
		}
		extraPart := ""
		if len(extraParts) > 0 {
			extraPart = "\n\n" + strings.Join(extraParts, "\n\n")
		}

		if wasVisited {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s%s\n\nExits: %s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				extraPart,
				m.formatExitsWithColor())
		} else {
			m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s%s\n\nExits: %s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				extraPart,
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
		// Split-screen layout matching ScreenPlaying
		width := m.width
		height := m.height
		if width < 40 {
			width = 40
		}
		if height < 10 {
			height = 10
		}

		// Calculate proportional heights (same as ScreenPlaying)
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

		// Output viewport (top ~70%)
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("75")).
			Padding(1, 1).
			Width(width - 2).
			Height(viewportHeight - 2)

		content := loginScreenContent()
		s.WriteString(outputStyle.Render(content))
		s.WriteString("\n")

		// Status bar separator (middle ~10%)
		separatorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Bold(true).
			Width(width)
		separatorLine := separatorStyle.Render(strings.Repeat("─", width-2))
		s.WriteString(separatorLine)
		s.WriteString("\n")

		// Status line showing current prompt
		fieldLabel := "Username:"
		if m.inputField == "password" {
			fieldLabel = "Password:"
		}
		var statusText string
		if m.messageType != "" && m.message != "" {
			// Show error/success messages
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("LOGIN: ") + m.styledMessage(m.message)
		} else {
			// Show current field label
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("LOGIN: ") + fieldLabel
		}
		s.WriteString(separatorStyle.Align(lipgloss.Center).Render(statusText))
		s.WriteString("\n")
		s.WriteString(separatorLine)
		s.WriteString("\n")

		// Input area (bottom ~20%)
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("75")).
			Padding(0, 1).
			Width(width - 2).
			Height(inputHeight - 2)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

		// Don't center - return directly like ScreenPlaying
		return s.String()

	case ScreenRegister:
		// Split-screen layout matching ScreenPlaying
		width := m.width
		height := m.height
		if width < 40 {
			width = 40
		}
		if height < 10 {
			height = 10
		}

		// Calculate proportional heights (same as ScreenPlaying)
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

		// Output viewport (top ~70%)
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("147")).
			Padding(1, 1).
			Width(width - 2).
			Height(viewportHeight - 2)

		content := registerScreenContent()
		s.WriteString(outputStyle.Render(content))
		s.WriteString("\n")

		// Status bar separator (middle ~10%)
		separatorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("147")).
			Bold(true).
			Width(width)
		separatorLine := separatorStyle.Render(strings.Repeat("─", width-2))
		s.WriteString(separatorLine)
		s.WriteString("\n")

		// Status line showing current prompt
		fieldLabel := "Username:"
		if m.inputField == "password" {
			fieldLabel = "Password:"
		} else if m.inputField == "confirm_password" {
			fieldLabel = "Confirm password:"
		} else if m.inputField == "email" {
			fieldLabel = "Email (optional):"
		}
		var statusText string
		if m.messageType != "" && m.message != "" {
			// Show error/success messages
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("REGISTER: ") + m.styledMessage(m.message)
		} else {
			// Show current field label
			statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Render("REGISTER: ") + fieldLabel
		}
		s.WriteString(separatorStyle.Align(lipgloss.Center).Render(statusText))
		s.WriteString("\n")
		s.WriteString(separatorLine)
		s.WriteString("\n")

		// Input area (bottom ~20%)
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("147")).
			Padding(0, 1).
			Width(width - 2).
			Height(inputHeight - 2)
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

		// Don't center - return directly like ScreenPlaying
		return s.String()

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
		// CRITICAL: Clear message BEFORE rendering to prevent previous state from showing
		// This fixes the "output re-rendering glitch" bug where old content briefly appears
		m.message = ""
		m.messageType = ""

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
			Height(viewportHeight - 2). // Account for border
			Align(lipgloss.Left, lipgloss.Top) // Fill from top-left, don't center

		// Colorful status bar with mini progress bars
		statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)

		// Debug mode - show room ID when enabled
		debugInfo := ""
		if m.debugMode {
			debugInfo = " " + lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("[Room: %d]", m.currentRoom))
		}

		// Room info at top with styling (only in output viewport, no stats)
		itemStr := m.formatRoomItemsWithColor()
		charStr := m.formatRoomCharactersWithColor()

		var extraParts []string
		if itemStr != "" {
			extraParts = append(extraParts, itemStr)
		}
		if charStr != "" {
			extraParts = append(extraParts, charStr)
		}
		extraPart := ""
		if len(extraParts) > 0 {
			extraPart = "\n\n" + strings.Join(extraParts, "\n\n")
		}

		roomInfo := fmt.Sprintf("[%s]\n%s%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			extraPart,
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
		// Stats go in the actual status bar (middle panel) - include debug info if enabled
		statsLineWithDebug := statsLine + debugInfo
		s.WriteString(separatorStyle.Align(lipgloss.Center).Render(statsLineWithDebug))
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

		// ScreenPlaying uses full-width panels - don't center, return directly
		return s.String()
	}

	// Center in terminal (optional - can be disabled if causing issues)
	// Use lipgloss.Width() to correctly handle ANSI escape codes (fixes issue #75)
	// CRITICAL: Respect terminal width - don't center if content would be truncated
	// NOTE: Vertical centering is DISABLED to fix viewport height issues
	// Only horizontal centering is applied to non-fullscreen screens
	if m.width > 0 && m.width > 60 {
		lines := strings.Split(s.String(), "\n")
		var centered []string
		for _, line := range lines {
			visualWidth := lipgloss.Width(line)
			// Only center horizontally if it won't cause truncation
			padding := (m.width - visualWidth) / 2
			if padding > 0 && visualWidth > 0 && visualWidth < m.width {
				centered = append(centered, fmt.Sprintf("%*s%s", padding, "", line))
			} else {
				// Content exceeds or matches terminal width - don't pad, let it flow
				centered = append(centered, line)
			}
		}
		return strings.Join(centered, "\n")
	}

	return s.String()
}

// ============================================================
// STATIC SCREENS
// ============================================================

func welcomeScreen() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(1, 2).
		Render(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║    ██████╗ ███████╗████████╗██████╗  ██████╗                 ║
║    ██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗                ║
║    ██████╔╝█████╗     ██║   ██████╔╝██║   ██║                ║
║    ██╔══██╗██╔══╝     ██║   ██╔══██╗██║   ██║                ║
║    ██║  ██║███████╗   ██║   ██║  ██║╚██████╔╝                ║
║    ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝                 ║
║                                                            ║
║           ██████╗  █████╗  ██████╗ ██████╗                  ║
║           ██╔══██╗██╔══██╗██╔════╝██╔═══██╗                 ║
║           ██████╔╝███████║██║     ██║   ██║                ║
║           ██╔══██╗██╔══██║██║     ██║   ██║                 ║
║           ██║  ██║██║  ██║╚██████╗╚██████╔╝                 ║
║           ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝                  ║
║                                                            ║
║                    Welcome to Herbst MUD!                  ║
║                    Das Text-Adventure                       ║
║                                                            ║
╠════════════════════════════════════════════════════════════╣
║                                                            ║
║   1. Login      - Log in to your existing account         ║
║   2. Register   - Create a new player account             ║
║   3. Quit       - Exit the game                            ║
║                                                            ║
║   Use ↑/↓ or j/k to navigate, Enter to select            ║
║   Press ESC to go back                                     ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
`)
}

// loginScreenContent returns the content for the login screen (layout handled in View())
func loginScreenContent() string {
	return `
        ╔════════════════════════════════════════╗
        ║                                        ║
        ║              === LOGIN ===             ║
        ║                                        ║
        ║    Enter your credentials to           ║
        ║    continue your adventure.            ║
        ║                                        ║
        ║    Press ESC to go back to menu        ║
        ║                                        ║
        ╚════════════════════════════════════════╝
`
}

// registerScreenContent returns the content for the register screen (layout handled in View())
func registerScreenContent() string {
	return `
        ╔════════════════════════════════════════╗
        ║                                        ║
        ║           === CREATE ACCOUNT ===        ║
        ║                                        ║
        ║    Choose a username and password      ║
        ║    to begin your adventure.            ║
        ║                                        ║
        ║    Press ESC to go back to menu        ║
        ║                                        ║
        ╚════════════════════════════════════════╝
`
}


