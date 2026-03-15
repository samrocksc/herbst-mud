package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"herbst/db"
	"herbst/dbinit"
)

// StartingRoomID is the ID of the room players start in
const StartingRoomID = 5

// LoginState represents the authentication state
type LoginState int

const (
	StateUsername LoginState = iota
	StatePassword
	StateAuthenticated
)

func main() {
	// Initialize database
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
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

					// Create program with shared client
					p := tea.NewProgram(
						&model{
							connectedAt: time.Now(),
							session:     s,
							client:      client,
							currentRoom: StartingRoomID,
							loginState:  StateUsername,
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

	// Login state
	loginState      LoginState
	username        string
	password        string
	loginError      string

	// Player state
	currentRoom   int
	roomName      string
	roomDesc      string
	exits         map[string]int
	inputBuffer   string
	message       string
}

func (m model) Init() tea.Cmd {
	// If authenticated, load starting room info
	if m.loginState == StateAuthenticated && m.client != nil {
		room, err := m.client.Room.Get(context.Background(), StartingRoomID)
		if err != nil {
			m.err = fmt.Errorf("failed to load starting room: %v", err)
			return nil
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits
	}
	return nil
}

// authenticateUser verifies username and password against the database
func (m *model) authenticateUser() bool {
	if m.client == nil {
		m.loginError = "Database not connected"
		return false
	}

	ctx := context.Background()

	// Find user by email (username)
	user, err := m.client.User.Query().Where(func(s *sql.Selector) {
		s.Where(sql.EQ("email", m.username))
	}).Only(ctx)
	if err != nil {
		m.loginError = "Invalid username or password"
		return false
	}

	// Check password (plaintext comparison for now - in production use bcrypt)
	if user.Password != m.password {
		m.loginError = "Invalid username or password"
		return false
	}

	// Authentication successful
	m.loginState = StateAuthenticated
	m.loginError = ""
	return true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		key := msg.String()

		// Handle Enter key - process command
		if key == "enter" {
			input := m.inputBuffer
			m.inputBuffer = ""

			// Handle login flow
			if m.loginState == StateUsername {
				m.username = input
				m.loginState = StatePassword
				return m, nil
			} else if m.loginState == StatePassword {
				m.password = input
				// Attempt authentication
				if m.authenticateUser() {
					// Load starting room after successful auth
					if m.client != nil {
						room, err := m.client.Room.Get(context.Background(), StartingRoomID)
						if err != nil {
							m.err = fmt.Errorf("failed to load starting room: %v", err)
							return m, nil
						}
						m.currentRoom = room.ID
						m.roomName = room.Name
						m.roomDesc = room.Description
						m.exits = room.Exits
					}
				} else {
					// Login failed - reset to username
					m.loginState = StateUsername
					m.username = ""
					m.password = ""
				}
				return m, nil
			}

			// Handle game commands (only when authenticated)
			if m.loginState == StateAuthenticated {
				m.processCommand(input)
			}
			return m, nil
		}

		// Handle backspace
		if key == "backspace" {
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
			return m, nil
		}

		// Handle Ctrl+C or q
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle regular character input
		if len(key) == 1 {
			m.inputBuffer += key
		}
	}
	return m, nil
}

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
		m.message = "Commands:\n  n/north - Move north\n  s/south - Move south\n  e/east - Move east\n  w/west - Move west\n  look - Look around\n  exits - Show exits\n  quit - Exit game"
	case "look", "l":
		m.message = fmt.Sprintf("[%s]\n%s\n\nExits: %s", m.roomName, m.roomDesc, m.formatExits())
	case "exits", "x":
		m.message = fmt.Sprintf("Exits: %s", m.formatExits())
	case "quit", "q":
		// Note: q is handled in Update, this is just fallback
		m.message = "Type 'q' or Ctrl+C to quit"
	default:
		m.message = fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd)
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
		return true
	}

	// Move to new room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.message = fmt.Sprintf("Error moving: %v", err)
			return true
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits
		m.message = fmt.Sprintf("You go %s.\n\n[%s]\n%s\n\nExits: %s", direction, m.roomName, m.roomDesc, m.formatExits())
	}

	return true
}

func (m *model) formatExits() string {
	if len(m.exits) == 0 {
		return "none"
	}
	var dirs []string
	for dir := range m.exits {
		dirs = append(dirs, dir)
	}
	return strings.Join(dirs, ", ")
}

func (m model) View() string {
	// Build the view
	var s strings.Builder

	// Show login screen if not authenticated
	if m.loginState != StateAuthenticated {
		s.WriteString("================================\n")
		s.WriteString("     🐢 Welcome to Herbst MUD 🐢\n")
		s.WriteString("================================\n\n")

		if m.loginState == StateUsername {
			s.WriteString("Enter your username: ")
			s.WriteString(m.inputBuffer)
			s.WriteString("_")
		} else if m.loginState == StatePassword {
			s.WriteString(fmt.Sprintf("Username: %s\n", m.username))
			s.WriteString("Enter your password: ")
			// Mask password
			for range m.inputBuffer {
				s.WriteString("*")
			}
			s.WriteString("_")

			// Show error if any
			if m.loginError != "" {
				s.WriteString(fmt.Sprintf("\n\nError: %s\n", m.loginError))
				s.WriteString("Press Enter to try again...")
			}
		}

		// Center in terminal
		if m.width > 0 && m.height > 0 {
			lines := strings.Split(s.String(), "\n")
			var centered []string
			for _, line := range lines {
				padding := (m.width - len(line)) / 2
				if padding > 0 {
					centered = append(centered, fmt.Sprintf("%*s%s", padding, "", line))
				} else {
					centered = append(centered, line)
				}
			}
			return strings.Join(centered, "\n")
		}

		return s.String()
	}

	// Game view (authenticated)
	s.WriteString(fmt.Sprintf("[%s]\n", m.roomName))
	s.WriteString(fmt.Sprintf("%s\n\n", m.roomDesc))
	s.WriteString(fmt.Sprintf("Exits: %s\n\n", m.formatExits()))

	// Show message if any
	if m.message != "" {
		s.WriteString(m.message)
		s.WriteString("\n\n")
		m.message = "" // Clear after showing
	}

	// Input prompt
	s.WriteString("\n> " + m.inputBuffer + "_")

	// Center in terminal
	if m.width > 0 && m.height > 0 {
		lines := strings.Split(s.String(), "\n")
		var centered []string
		for _, line := range lines {
			padding := (m.width - len(line)) / 2
			if padding > 0 {
				centered = append(centered, fmt.Sprintf("%*s%s", padding, "", line))
			} else {
				centered = append(centered, line)
			}
		}
		return strings.Join(centered, "\n")
	}

	return s.String()
}