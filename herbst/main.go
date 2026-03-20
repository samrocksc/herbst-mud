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
	"github.com/charmbracelet/bubbles/viewport"
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

					// Get terminal size from SSH session
					pty, winCh, ok := s.Pty()
					var initialWidth, initialHeight int
					if ok {
						initialWidth = pty.Window.Width
						initialHeight = pty.Window.Height
						log.Printf("PTY size: %dx%d", initialWidth, initialHeight)
					} else {
						// Fallback if no PTY
						initialWidth = 80
						initialHeight = 24
						log.Printf("No PTY, using default size: %dx%d", initialWidth, initialHeight)
					}

					// Create initial text input for login/register
					ti := textinput.New()
					ti.Placeholder = "Enter choice..."
					ti.Focus()

					// Create loading spinner
					sp := spinner.New()
					sp.Spinner = spinner.Dot
					sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

					// Create model with initial size
					m := &model{
						connectedAt:  time.Now(),
						session:      s,
						client:       client,
						screen:       ScreenWelcome,
						currentRoom:  StartingRoomID,
						textInput:    ti,
						spinner:      sp,
						visitedRooms: make(map[int]bool),
						knownExits:   make(map[string]bool),
						width:        initialWidth,
						height:       initialHeight,
						maxHistory:   50,
					}

					// Create program with shared client
					p := tea.NewProgram(
						m,
						tea.WithInput(s),
						tea.WithOutput(s),
						tea.WithAltScreen(), // Full-screen mode
					)

					// Handle window resize events from SSH (only if PTY was allocated)
					if ok && winCh != nil {
						go func() {
							for win := range winCh {
								p.Send(tea.WindowSizeMsg{Width: win.Width, Height: win.Height})
							}
						}()
					}

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

	// Scrollable viewport for game output (ScreenPlaying)
	viewport viewport.Model

	// Room tracking
	visitedRooms map[int]bool
	knownExits   map[string]bool // For color-coded exits

	// Room items (GitHub #89 - Item system)
	roomItems []RoomItem

	// Room characters (GitHub #145 - Look command room display)
	roomCharacters []roomCharacter

	// Debug mode - shows room ID in status bar
	debugMode bool

	// Message history buffer for output pane (UI-21)
	messageHistory  []string // all messages (oldest → newest)
	messageTypes    []string // parallel array with message types
	historyOffset   int      // 0 = at latest (pinned), scrolls up into history
	maxHistory      int      // max messages to keep
	isScrolling     bool     // true when user has scrolled away from bottom
}

// RoomItem represents an item in a room for display
type RoomItem struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ExamineDesc     string         `json:"examineDesc"`
	HiddenDetails   []HiddenDetail `json:"hiddenDetails"`
	HiddenThreshold int            `json:"hiddenThreshold"`
	IsImmovable     bool           `json:"isImmovable"`
	Color           string         `json:"color"`
	IsVisible       bool           `json:"isVisible"`
	ItemType        string         `json:"itemType"`
	Weight          int            `json:"weight"`
	ItemDamage      int            `json:"itemDamage"`
	ItemDurability  int            `json:"itemDurability"`
	RevealCondition map[string]any `json:"revealCondition"`
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

	// Quest tracker panel colors
	questTitleColor     = lipgloss.Color("75")    // Blue
	questProgressColor  = lipgloss.Color("226")  // Yellow
	questCompletedColor = lipgloss.Color("46")   // Green
	questAvailableColor = lipgloss.Color("141")  // Purple

	// Quest tracker panel styles
	questTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(questTitleColor)

	questBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2)

	questProgressStyle = lipgloss.NewStyle().
				Foreground(questProgressColor)

	questCompletedStyle = lipgloss.NewStyle().
				Foreground(questCompletedColor).
				Strikethrough(true)

	questAvailableStyle = lipgloss.NewStyle().
				Foreground(questAvailableColor)

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
// combatDamageStyle for damage messages (red)
var combatDamageStyle = lipgloss.NewStyle().
	Foreground(red).
	Bold(true)

// combatHealStyle for healing messages (green)
var combatHealStyle = lipgloss.NewStyle().
	Foreground(green).
	Bold(true)

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
	case "damage":
		// Combat damage: red text with ⚔ prefix
		return combatDamageStyle.Render("⚔ ") + msg
	case "heal":
		// Combat healing: green text with ♥ prefix
		return combatHealStyle.Render("♥ ") + msg
	default:
		return msg
	}
}

// AppendMessage adds a message to the history buffer (UI-21 message history)
func (m *model) AppendMessage(text, msgType string) {
	m.messageHistory = append(m.messageHistory, text)
	m.messageTypes = append(m.messageTypes, msgType)
	if len(m.messageHistory) > m.maxHistory {
		m.messageHistory = m.messageHistory[len(m.messageHistory)-m.maxHistory:]
		m.messageTypes = m.messageTypes[len(m.messageTypes)-m.maxHistory:]
	}
	m.historyOffset = 0
	m.isScrolling = false
}

// buildOutputContent constructs the message history content for display
func (m *model) buildOutputContent() string {
	total := len(m.messageHistory)
	if total == 0 {
		return ""
	}

	var lines []string

	if !m.isScrolling {
		// Show last 3 messages (or fewer)
		start := 0
		if total > 3 {
			start = total - 3
		}
		for i := start; i < total; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
	} else {
		// Scrolled up: show from historyOffset toward newest (excluding latest, which is pinned)
		for i := m.historyOffset; i < total-1; i++ {
			lines = append(lines, styleMessage(m.messageHistory[i], m.messageTypes[i]))
		}
		// Latest pinned at bottom with ─── NEWEST ─── marker
		lines = append(lines, "─── NEWEST ───")
		lines = append(lines, styleMessage(m.messageHistory[total-1], m.messageTypes[total-1]))
	}

	return strings.Join(lines, "\n\n")
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
		// Initialize or resize the viewport for scrollable output
		inputHeight := m.height * 20 / 100
		if inputHeight < 3 {
			inputHeight = 3
		}
		statusHeight := m.height * 10 / 100
		if statusHeight < 3 {
			statusHeight = 3
		}
		vpHeight := m.height - inputHeight - statusHeight
		if vpHeight < 5 {
			vpHeight = 5
		}
		// vpWidth includes the border chars (2), so use full width
		m.viewport = viewport.New(msg.Width, vpHeight)
		// Debug logging only in debug mode to avoid noise
		if m.debugMode {
			log.Printf("DEBUG: Window size changed: %dx%d", m.width, m.height)
		}

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

		// Let viewport handle scrolling keys first (before vim/enter/escape)
		if m.screen == ScreenPlaying {
			vp, vpCmd := m.viewport.Update(msg)
			m.viewport = vp
			if cmd == nil {
				cmd = vpCmd
			}
			// Viewport handles up/down/pgup/pgdn for scrolling; don't return here —
			// continue so textinput and other screen-specific handlers also get the event.
		}

		// Message history scrolling (ctrl+k up, ctrl+j down) - must intercept BEFORE Enter
		if m.screen == ScreenPlaying {
			switch key {
			case "ctrl+k":
				// Scroll up (older messages)
				if !m.isScrolling {
					m.isScrolling = true
					m.historyOffset = 1
				} else {
					m.historyOffset++
				}
				maxOffset := len(m.messageHistory) - 1
				if m.historyOffset > maxOffset {
					m.historyOffset = maxOffset
				}
				return m, nil
			case "ctrl+j":
				// Scroll down (newer messages) - DO NOT process as Enter here
				if !m.isScrolling {
					return m, nil
				}
				m.historyOffset--
				if m.historyOffset < 0 {
					m.historyOffset = 0
					m.isScrolling = false
				}
				return m, nil
			}
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
		if key == "enter" || key == "ctrl+m" {
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
		m.AppendMessage("", "")
		// Re-initialize menu items for welcome screen
		m.menuItems = []string{"Login", "Register", "Quit"}
		m.menuCursor = 0
	case ScreenProfile, ScreenEditField:
		m.screen = ScreenPlaying
		m.textInput.SetValue("")
		m.inputBuffer = ""
		m.AppendMessage("", "")
	case ScreenPlaying:
		// Could add a "really quit?" confirmation
		m.AppendMessage("Type 'quit' or press Ctrl+C to exit", "info")
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
		m.AppendMessage(fmt.Sprintf("Connection error: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	// Debug: log.Printf("Auth response status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Invalid username or password. Try again.", "error")
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
	m.AppendMessage(fmt.Sprintf("Welcome back, %s!", m.currentUserName), "success")

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

func (m *model) attemptRegistration(email string) {
	// Use REST API for user creation
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
	m.AppendMessage(fmt.Sprintf("Account created! Welcome to Herbst MUD, %s!", m.currentUserName), "success")

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
		m.AppendMessage(`Commands:
  n/north, s/south, e/east, w/west - Move
  look/l - Look around (shows items)
  exits/x - Show exits  
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  inventory/i - Show your inventory
  quests/q - Show your quest log
  whoami - Show your info
  profile/p - Edit character profile
  clear/cls - Clear screen
  quit - Exit game`, "info")
	case "look", "l":
		m.loadRoomItems()
		m.loadRoomCharacters()
		m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nExits: %s%s%s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor(),
			m.formatRoomItems(),
			m.formatRoomCharacters()), "info")
	case "exits", "x":
		m.AppendMessage(fmt.Sprintf("Exits: %s", m.formatExitsWithColor()), "info")
	case "examine", "ex", "inspect":
		m.handleExamineCommand(cmd)
	case "search", "perception":
		// GitHub #12 - Perception check to reveal hidden items
		m.handleSearchCommand(cmd)
	case "whoami":
		// Show character info including level with progress bars
		m.AppendMessage(fmt.Sprintf("=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
			m.currentUserName, m.currentUserID, m.roomName,
			m.characterLevel, m.characterExperience,
			StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)), "info")
	case "profile", "p":
		m.screen = ScreenProfile
		m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
		m.menuCursor = 0
		m.AppendMessage("", "")
	case "peer":
		m.handlePeerCommand(cmd)
	case "debug":
		m.handleDebugCommand(cmd)
		return
	case "clear", "cls":
		// Clear message history (UI-21)
		m.messageHistory = nil
		m.messageTypes = nil
		m.historyOffset = 0
		m.isScrolling = false
		m.inputBuffer = ""
		return
	case "quit", "q":
		m.AppendMessage("Thanks for playing! Goodbye!", "success")
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
		// Check for quests command
		if cmd == "quests" || cmd == "q" || cmd == "quest" {
			m.handleQuestsCommand(cmd)
			return
		}
		// Check for skills command
		if cmd == "skills" {
			m.handleSkillsCommand(cmd)
			return
		}
		// Check for talents command
		if cmd == "talents" {
			m.handleTalentsCommand(cmd)
			return
		}
		// Check for skill equip command
		if strings.HasPrefix(cmd, "skill ") {
			m.handleSkillEquipCommand(cmd)
			return
		}
		// Check for talent equip/unequip/swap commands
		if strings.HasPrefix(cmd, "talent ") {
			m.handleTalentEquipCommand(cmd)
			return
		}
		m.AppendMessage(fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd), "error")
	}
}

// handleSkillsCommand displays character skills
func (m *model) handleSkillsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/skills", RESTAPIBase, m.currentCharacterID)
	resp, err := http.Get(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching skills: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load skills", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing skills: %v", err), "error")
		return
	}

	skills, ok := result["skills"].(map[string]interface{})
	if !ok {
		m.AppendMessage("Error: skills data not found", "error")
		return
	}

	// Format skills display
	output := "=== Your Skills ===\n\n"
	for skillName, skillData := range skills {
		data := skillData.(map[string]interface{})
		level := int(data["level"].(float64))
		bonus := data["bonus"].(string)
		output += fmt.Sprintf("%-15s Lv: %2d  %s\n", skillName+":", level, bonus)
	}

	output += "\nSkills are always active and provide passive bonuses."
	m.AppendMessage(output, "info")
}

// handleTalentsCommand displays equipped talents
func (m *model) handleTalentsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
	resp, err := http.Get(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching talents: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load talents", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing talents: %v", err), "error")
		return
	}

	// Format talents display with slots
	output := "=== Your Talents ===\n\n"
	slots, ok := result["slots"].([]interface{})
	if !ok {
		// No talents equipped yet
		output += "No talents equipped.\n\n"
		output += "Use: talent equip <talent_id> <slot>\n"
		output += "Slots: 1-4 (quick access keys)\n"
		m.AppendMessage(output, "info")
		return
	}

	emptySlots := 0
	for i := 1; i <= 4; i++ {
		if i < len(slots) && slots[i] != nil {
			slot := slots[i].(map[string]interface{})
			name := slot["name"].(string)
			desc := slot["description"].(string)
			output += fmt.Sprintf("[%d] %s\n     %s\n\n", i, name, desc)
		} else {
			output += fmt.Sprintf("[%d] (empty)\n\n", i)
			emptySlots++
		}
	}

	if emptySlots == 4 {
		output += "No talents equipped. Use 'talent equip <id> <slot>' to equip."
	}

	m.AppendMessage(output, "info")
}

// handleSkillEquipCommand handles skill equip command
func (m *model) handleSkillEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	// Skills are always active, no equip needed
	m.AppendMessage("Skills are always active and cannot be unequipped.\nThey provide passive bonuses based on your skill level.", "info")
}

// handleTalentEquipCommand handles talent equip/unequip/swap commands
func (m *model) handleTalentEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage:\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
		return
	}

	action := parts[1]

	switch action {
	case "equip":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent equip <talent_id> <slot>\nExample: talent equip 1 2", "error")
			return
		}
		talentID := parts[2]
		slot := parts[3]

		// Validate slot is 1-4
		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		// Call API to equip talent
		url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"talent_id":%s,"slot":%s}`, talentID, slot)
		resp, err := http.Post(url, "application/json", strings.NewReader(reqBody))
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error equipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			m.AppendMessage("Failed to equip talent", "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("Talent equipped in slot %s", slot), "success")

	case "unequip":
		if len(parts) != 3 {
			m.AppendMessage("Usage: talent unequip <slot>\nExample: talent unequip 2", "error")
			return
		}
		slot := parts[2]

		// Validate slot
		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		// Call API to unequip talent
		url := fmt.Sprintf("%s/characters/%d/talents/%s", RESTAPIBase, m.currentCharacterID, slot)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talent unequipped from slot %s", slot), "success")

	case "swap":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent swap <slot1> <slot2>\nExample: talent swap 1 2", "error")
			return
		}
		slot1 := parts[2]
		slot2 := parts[3]

		// Validate slots
		slot1Num, slot2Num := 0, 0
		fmt.Sscanf(slot1, "%d", &slot1Num)
		fmt.Sscanf(slot2, "%d", &slot2Num)
		if slot1Num < 1 || slot1Num > 4 || slot2Num < 1 || slot2Num > 4 {
			m.AppendMessage("Slots must be between 1 and 4", "error")
			return
		}

		// Call API to swap talents
		url := fmt.Sprintf("%s/characters/%d/talents/swap", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"slot1":%s,"slot2":%s}`, slot1, slot2)
		req, err := http.NewRequest("PUT", url, strings.NewReader(reqBody))
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talents swapped between slot %s and %s", slot1, slot2), "success")

	default:
		m.AppendMessage("Usage:\n  talent - Show talents\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
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
		m.AppendMessage("You can't go that way.", "error")
		return true
	}

	// Mark exit as known
	m.knownExits[direction] = true

	// Move to new room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error moving: %v", err), "error")
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
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		} else {
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		}
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
		m.AppendMessage("Usage: peer <direction>\nDirections: north, south, east, west, up, down", "error")
		return
	}
	direction := strings.ToLower(parts[1])

	// Validate direction
	validDirs := map[string]string{"north": "north", "south": "south", "east": "east", "west": "west", "up": "up", "down": "down"}
	dir, ok := validDirs[direction]
	if !ok {
		m.AppendMessage("Invalid direction. Use: north, south, east, west, up, down", "error")
		return
	}

	// Check if exit exists
	nextRoomID, ok := m.exits[dir]
	if !ok {
		m.AppendMessage("You can't peer that way — there's no exit.", "error")
		return
	}

	// Get the room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error looking: %v", err), "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("You peer %s...\n\n[%s]\n%s",
			dir,
			lipgloss.NewStyle().Bold(true).Foreground(blue).Render(room.Name),
			room.Description), "info")
	}
}

// handleSearchCommand handles the search/perception command to reveal hidden items
// GitHub #12 - Look System: Hidden Items and Reveal Conditions
func (m *model) handleSearchCommand(cmd string) {
	if m.currentRoom == 0 {
		m.AppendMessage("You can't search here.", "error")
		return
	}

	// Fetch all items (including hidden) for this room
	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error searching: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Error searching the area.", "error")
		return
	}

	var allItems []RoomItem
	if err := json.NewDecoder(resp.Body).Decode(&allItems); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing items: %v", err), "error")
		return
	}

	var found []string
	revealed := 0

	for _, item := range allItems {
		// Skip already visible items
		if item.IsVisible {
			continue
		}

		// Check if this is a hidden item that can be revealed by perception
		if item.RevealCondition != nil {
			revealType, _ := item.RevealCondition["type"].(string)
			if revealType == "perception_check" {
				// Try to reveal the item with perception check
				// Use character level as skill level for now
				revealResp, err := http.Post(
					fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
					"application/json",
					strings.NewReader(fmt.Sprintf(`{"revealType":"perception_check","skillLevel":%d}`, m.characterLevel)),
				)
				if err == nil {
					defer revealResp.Body.Close()
					if revealResp.StatusCode == http.StatusOK {
						revealed++
						found = append(found, item.Name)
					}
				}
			}
		}
	}

	// Reload room items to see the newly revealed
	m.loadRoomItems()

	if revealed > 0 {
		m.AppendMessage(fmt.Sprintf("🔍 You search the area carefully...\n\n✨ You discovered %d hidden item(s): %s",
			revealed, strings.Join(found, ", ")), "success")
	} else {
		m.AppendMessage("🔍 You search the area carefully...\n\nYou find nothing of interest.", "info")
	}
}

func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		// Show current debug status
		if m.debugMode {
			m.AppendMessage("Debug mode: ON (Room ID visible in status bar)", "info")
		} else {
			m.AppendMessage("Debug mode: OFF\nUsage: debug on | debug off", "info")
		}
		return
	}

	subCmd := parts[1]
	switch subCmd {
	case "on", "true", "1", "yes":
		m.debugMode = true
		m.AppendMessage("Debug mode: ON (Room ID will show in status bar)", "success")
	case "off", "false", "0", "no":
		m.debugMode = false
		m.AppendMessage("Debug mode: OFF", "info")
	default:
		m.AppendMessage("Usage: debug on | debug off", "error")
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
		m.AppendMessage("Take what? Usage: take <item name>", "error")
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
		m.AppendMessage(fmt.Sprintf("You don't see any %s here.", itemName), "error")
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
		m.AppendMessage(fmt.Sprintf("You can't take the %s. It's firmly fixed in place.", colorStyle.Render(targetItem.Name)), "error")
		return
	}

	// Take the item - move it to player's inventory (roomId = 0 or null)
	url := fmt.Sprintf("%s/equipment/%d", RESTAPIBase, targetItem.ID)
	jsonData, _ := json.Marshal(map[string]interface{}{"roomId": nil})
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error picking up item: %v", err), "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error picking up item: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage(fmt.Sprintf("Failed to pick up %s.", targetItem.Name), "error")
		return
	}

	m.AppendMessage(fmt.Sprintf("You pick up the %s.", targetItem.Name), "success")
}

// handleDropCommand handles the drop command
func (m *model) handleDropCommand(cmd string) {
	// Extract item name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Drop what? Usage: drop <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")

	// For now, show a message that inventory is not fully implemented
	// This would need player inventory tracking
	m.AppendMessage(fmt.Sprintf("You don't have any %s to drop.", itemName), "error")
}

// handleExamineCommand handles the examine/ex/inspect/i command
func (m *model) handleExamineCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Examine what? Usage: examine <item>", "error")
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(target)

	// First check room items (only visible ones for display)
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue // Skip hidden items - they'll be handled separately
		}
		if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Check for hidden items that could be revealed by examining this target
	// (GitHub #12 - Hidden Items and Reveal Conditions)
	if m.currentRoom > 0 {
		// Fetch all items (including hidden) for this room
		resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var allItems []RoomItem
				if json.NewDecoder(resp.Body).Decode(&allItems) == nil {
					for _, item := range allItems {
						// Check if this is a hidden item that reveals on examine
						if !item.IsVisible && item.RevealCondition != nil {
							revealType, _ := item.RevealCondition["type"].(string)
							revealTarget, _ := item.RevealCondition["target"].(string)
							if revealType == "examine" && strings.ToLower(revealTarget) == target {
								// Try to reveal the item
								revealResp, err := http.Post(
									fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
									"application/json",
									strings.NewReader(fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, revealTarget, m.characterLevel)),
								)
								if err == nil {
									defer revealResp.Body.Close()
									if revealResp.StatusCode == http.StatusOK {
										// Item revealed! Reload room items and try again
										m.loadRoomItems()
										// Re-check with now-visible item
										for _, ri := range m.roomItems {
											if strings.Contains(strings.ToLower(ri.Name), target) || strings.ToLower(ri.Name) == target {
												m.AppendMessage("✨ You discovered something hidden!\n\n", "info")
												m.displayItemDetails(ri)
												return
											}
										}
									}
								}
							}
						}
					}
				}
			}
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
						m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
							npc.Name, npc.Description, npc.Level, npc.Disposition), "info")
						return
					}
				}
			}
		}
	}

	m.AppendMessage(fmt.Sprintf("You don't see '%s' here.", target), "error")
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

	m.AppendMessage(details.String(), "info")
}

// getItemIcon returns an emoji icon based on item type
func getItemIcon(itemType string) string {
	switch itemType {
	case "weapon":
		return "⚔️"
	case "armor":
		return "🛡️"
	case "potion":
		return "🧪"
	case "food":
		return "🍖"
	case "scroll":
		return "📜"
	case "key":
		return "🔑"
	case "treasure":
		return "💎"
	case "quest":
		return "📋"
	default:
		return "📦"
	}
}

// getItemRarityColor returns a lipgloss color based on item rarity
func getItemRarityColor(rarity string) lipgloss.Color {
	switch rarity {
	case "rare":
		return lipgloss.Color("51") // Blue
	case "epic":
		return lipgloss.Color("201") // Magenta
	case "legendary":
		return lipgloss.Color("220") // Gold
	default:
		return lipgloss.Color("white")
	}
}

// inventoryItem represents an item in the player's inventory
type inventoryItem struct {
	ID          int
	Name        string
	Description string
	ItemType    string
	IsEquipped  bool
	Rarity      string
}

// handleInventoryCommand handles the inventory/i command
func (m *model) handleInventoryCommand() {
	// Fetch player's inventory from API
	resp, err := http.Get(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching inventory: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	// Parse inventory items
	var rawItems []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ItemType    string `json:"itemType"`
		IsEquipped  bool   `json:"isEquipped"`
		Rarity      string `json:"rarity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawItems); err != nil {
		m.AppendMessage("You aren't carrying anything.", "info")
		return
	}

	// Convert to typed items
	items := make([]inventoryItem, len(rawItems))
	for i, raw := range rawItems {
		items[i] = inventoryItem(raw)
	}

	if len(items) == 0 {
		m.AppendMessage("Your pockets are empty. Time to loot some stuff!", "info")
		return
	}

	// Format inventory display with icons and styling
	var inv strings.Builder
	inv.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("🎒 INVENTORY"))
	inv.WriteString("\n")
	inv.WriteString(strings.Repeat("─", 30))
	inv.WriteString("\n\n")

	// Group items by type for better organization
	typeGroups := make(map[string][]inventoryItem)

	for _, item := range items {
		typeGroups[item.ItemType] = append(typeGroups[item.ItemType], item)
	}

	// Display items grouped by type with icons
	for itemType, groupItems := range typeGroups {
		icon := getItemIcon(itemType)
		typeLabel := strings.ToUpper(itemType)
		inv.WriteString(lipgloss.NewStyle().Bold(true).Foreground(cyan).Render(fmt.Sprintf("%s %s", icon, typeLabel)))
		inv.WriteString("\n")

		for _, invItem := range groupItems {
			rarityColor := getItemRarityColor(invItem.Rarity)
			itemStyle := lipgloss.NewStyle().Foreground(rarityColor)

			equipped := ""
			if invItem.IsEquipped {
				equipped = " " + lipgloss.NewStyle().Bold(true).Foreground(green).Render("⚡ equipped")
			}

			inv.WriteString(fmt.Sprintf("  %s %s%s\n", icon, itemStyle.Render(invItem.Name), equipped))
			if invItem.Description != "" {
				inv.WriteString(fmt.Sprintf("     %s\n", invItem.Description))
			}
		}
		inv.WriteString("\n")
	}

	m.AppendMessage(inv.String(), "info")
}

// handleQuestsCommand handles the quests/q command to display quest tracker
func (m *model) handleQuestsCommand(cmd string) {
	// Show placeholder when no character is selected
	if m.currentCharacterID == 0 {
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Fetch quests from API
	// For now, we'll return mock data until the full quest system is implemented
	// In production, this would call: GET /characters/:id/quests
	resp, err := http.Get(fmt.Sprintf("%s/characters/%d/quests", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		// Network error - show placeholder (expected in dev mode without server)
		m.displayQuestTrackerPlaceholder()
		return
	}
	defer resp.Body.Close()

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		// If no quests endpoint exists yet, show placeholder message
		// This allows the feature to work before the full quest system is built
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Parse quest response
	var questResp struct {
		Quests []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
			Objectives  []struct {
				Description string `json:"description"`
				Current     int    `json:"current"`
				Total       int    `json:"total"`
			} `json:"objectives"`
			Giver  string `json:"giver"`
			Rewards string `json:"rewards"`
		} `json:"quests"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&questResp); err != nil || len(questResp.Quests) == 0 {
		// No quests available - show placeholder
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Format quest tracker display with Lip Gloss styling
	var quests strings.Builder

	// Title
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	activeCount := 0
	availableCount := 0
	completedCount := 0

	for _, quest := range questResp.Quests {
		status := quest.Status
		switch status {
		case "in_progress":
			activeCount++
		case "available":
			availableCount++
		case "completed":
			completedCount++
		}

		// Quest box with styled border
		quests.WriteString(questBoxStyle.Render("") + "\n")

		// Quest name with status color
		statusColor := questAvailableStyle
		statusText := "Available"
		if status == "in_progress" {
			statusColor = questProgressStyle
			statusText = "In Progress"
		} else if status == "completed" {
			statusColor = questCompletedStyle
			statusText = "Completed"
		}

		quests.WriteString(fmt.Sprintf("  %s [%s]\n", questTitleStyle.Render(quest.Name), statusColor.Render(statusText)))

		// Description
		if quest.Description != "" {
			quests.WriteString(fmt.Sprintf("    %s\n", quest.Description))
		}

		// Objectives with progress
		if len(quest.Objectives) > 0 {
			quests.WriteString("\n  Objectives:\n")
			for _, obj := range quest.Objectives {
				progress := fmt.Sprintf("%d/%d", obj.Current, obj.Total)
				if obj.Current >= obj.Total {
					quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", obj.Description, questCompletedStyle.Render("("+progress+")")))
				} else {
					quests.WriteString(fmt.Sprintf("    ○ %s %s\n", obj.Description, questProgressStyle.Render("("+progress+")")))
				}
			}
		}

		// Giver
		if quest.Giver != "" {
			quests.WriteString(fmt.Sprintf("\n  Giver: %s\n", quest.Giver))
		}

		// Rewards
		if quest.Rewards != "" {
			quests.WriteString(fmt.Sprintf("  Reward: %s\n", quest.Rewards))
		}

		quests.WriteString("\n")
	}

	// Summary footer
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString(fmt.Sprintf("  Active: %d  |  Available: %d  |  Completed: %d\n",
		activeCount, availableCount, completedCount))
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	m.AppendMessage(quests.String(), "info")
}

// displayQuestTrackerPlaceholder shows a placeholder quest tracker
// when no quests are available (before full quest system is implemented)
func (m *model) displayQuestTrackerPlaceholder() {
	var quests strings.Builder

	// Title with Lip Gloss styling
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	// Placeholder quests from the quest system spec
	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Prove Yourself"),
		questProgressStyle.Render("In Progress")))

	quests.WriteString("    The Scrapyard ain't for the weak. Kill 3 Scrap Rats\n")
	quests.WriteString("    and I'll let you into New Venice proper.\n\n")

	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Kill Scrap Rat", questProgressStyle.Render("(2/3)")))
	quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", "Find Guard Marco at Foggy Gate", questCompletedStyle.Render("(done)")))

	quests.WriteString("\n  Giver: Guard Marco\n")
	quests.WriteString("  Reward: 10 coins\n\n")

	// Second placeholder quest
	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Ooze Samples"),
		questAvailableStyle.Render("Available")))

	quests.WriteString("    Jane needs Ooze samples for her research.\n")
	quests.WriteString("    The Leaking Pipes have plenty.\n\n")

	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Collect glowing goo", questProgressStyle.Render("(0/5)")))

	quests.WriteString("\n  Giver: Scavenger Jane\n")
	quests.WriteString("  Reward: repair_kit, scavenge skill\n\n")

	// Summary footer
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString("  Active: 1  |  Available: 1  |  Completed: 0\n")
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	quests.WriteString("\n" + infoStyle.Render("  Use 'quest <name>' for details, 'accept <quest>' to begin."))

	m.AppendMessage(quests.String(), "info")
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

		// Debug: log the actual dimensions being used
		if m.debugMode {
			log.Printf("ScreenPlaying: terminal dimensions: %dx%d (raw: %dx%d)", width, height, m.width, m.height)
		}

		if width < 40 {
			log.Printf("WARNING: width too small (%d), defaulting to 80", width)
			width = 80
		}
		if height < 10 {
			log.Printf("WARNING: height too small (%d), defaulting to 24", height)
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

		// Build viewport content (room info + message history)
		roomInfo := fmt.Sprintf("[%s]\n%s\n\nExits: %s",
			lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
			m.roomDesc,
			m.formatExitsWithColor())

		// Append message history if any
		if len(m.messageHistory) > 0 {
			msgContent := m.buildOutputContent()
			if msgContent != "" {
				roomInfo += "\n\n" + msgContent
			}
		}

		// Update viewport content and size
		if m.viewport.Width != width {
			m.viewport.Width = width
		}
		if m.viewport.Height != height-statusHeight-inputHeight {
			m.viewport.Height = height - statusHeight - inputHeight
		}
		m.viewport.SetContent(roomInfo)

		// Scrollable viewport (top ~70%) - with pink border
		viewportStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Width(width)
		s.WriteString(viewportStyle.Render(m.viewport.View()))
		s.WriteString("\n")

		// Full-width status bar (middle ~10%)
		statsLine := MiniStatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)
		debugInfo := ""
		if m.debugMode {
			debugInfo = " " + lipgloss.NewStyle().Foreground(yellow).Bold(true).Render(fmt.Sprintf("[Room: %d]", m.currentRoom))
		}
		statusBarStyle := lipgloss.NewStyle().
			Foreground(pink).
			Background(lipgloss.Color("235")).
			Bold(true).
			Width(width).
			Padding(0, 1)
		s.WriteString(statusBarStyle.Render(statsLine + debugInfo))
		s.WriteString("\n")

		// Full-width input area (bottom ~20%)
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pink).
			Padding(0, 1).
			Width(width).
			Height(inputHeight - 2) // Subtract 2 for border top/bottom
		s.WriteString(inputStyle.Render(promptStyle.Render("> ") + m.textInput.View()))

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
		return strings.Join(centered, "\n")
	}

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
		Width(width).
		Height(outputHeight - 2) // Subtract 2 for border top/bottom

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
		Width(width).
		Height(inputHeight - 2) // Subtract 2 for border top/bottom

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
		Width(width).
		Height(outputHeight - 2) // Subtract 2 for border top/bottom

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
		Width(width).
		Height(inputHeight - 2) // Subtract 2 for border top/bottom

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
		Width(width).
		Height(outputHeight - 2) // Subtract 2 for border top/bottom

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
		Width(width).
		Height(inputHeight - 2) // Subtract 2 for border top/bottom

	var sb strings.Builder
	sb.WriteString(outputStyle.Render(outputContent.String()))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(inputView))

	return sb.String()
}

