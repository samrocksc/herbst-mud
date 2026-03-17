package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"herbst/db"
	"herbst/dbinit"
)

// StartingRoomID is the ID of the room players start in
const StartingRoomID = 5

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

		// Initialize skills and talents
		if err := dbinit.InitSkillsAndTalents(client); err != nil {
			log.Printf("Warning: failed to initialize skills and talents: %v", err)
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

	// Player state
	characterId   int
	currentRoom   int
	roomName      string
	roomDesc      string
	exits         map[string]int
	inputBuffer   string
	message       string
}

func (m model) Init() tea.Cmd {
	// Load starting room info
	if m.client != nil {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		key := msg.String()

		// Handle Enter key - process command
		if key == "enter" {
			m.processCommand(m.inputBuffer)
			m.inputBuffer = ""
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

	// Handle commands with arguments
	parts := strings.Fields(cmd)
	if len(parts) > 1 {
		baseCmd := parts[0]
		args := parts[1:]
		
		switch baseCmd {
		case "/swap-skill":
			m.handleSwapSkill(args)
			return
		case "/swap-talent":
			m.handleSwapTalent(args)
			return
		}
	}

	// Handle movement commands
	if m.handleMovement(cmd) {
		return
	}

	// Handle other commands
	switch cmd {
	case "help", "?":
		m.message = "Commands:\n  n/north - Move north\n  s/south - Move south\n  e/east - Move east\n  w/west - Move west\n  look - Look around\n  exits - Show exits\n  skills - Show your skills\n  talents - Show your talents\n  /swap-skill - Swap a skill\n  /swap-talent - Swap a talent\n  quit - Exit game"
	case "look", "l":
		m.message = fmt.Sprintf("[%s]\n%s\n\nExits: %s", m.roomName, m.roomDesc, m.formatExits())
	case "exits", "x":
		m.message = fmt.Sprintf("Exits: %s", m.formatExits())
	case "skills":
		m.handleSkillsCommand()
	case "talents":
		m.handleTalentsCommand()
	case "swap-skill", "/swap-skill":
		m.message = "Usage: /swap-skill <old_skill> <new_skill>\nUse 'skills' to see available skills"
	case "swap-talent", "/swap-talent":
		m.message = "Usage: /swap-talent <old_talent> <new_talent>\nUse 'talents' to see available talents"
	case "quit", "q":
		// Note: q is handled in Update, this is just fallback
		m.message = "Type 'q' or Ctrl+C to quit"
	default:
		m.message = fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd)
	}
}

// handleSkillsCommand shows available skills
func (m *model) handleSkillsCommand() {
	if m.client == nil {
		m.message = "Database not connected"
		return
	}

	ctx := context.Background()

	// Get all available skills
	allSkills, err := m.client.Skill.Query().All(ctx)
	if err != nil {
		m.message = fmt.Sprintf("Error fetching skills: %v", err)
		return
	}

	if len(allSkills) == 0 {
		m.message = "No skills available yet."
		return
	}

	var sb strings.Builder
	sb.WriteString("=== AVAILABLE SKILLS ===\n\n")
	for _, s := range allSkills {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", s.Name, s.Description))
		sb.WriteString(fmt.Sprintf("  Type: %s | Cost: %d | Cooldown: %ds | Power: %d\n\n", 
			s.Type, s.Cost, s.Cooldown, s.Power))
	}

	m.message = sb.String()
}

// handleTalentsCommand shows available talents
func (m *model) handleTalentsCommand() {
	if m.client == nil {
		m.message = "Database not connected"
		return
	}

	ctx := context.Background()

	// Get all available talents
	allTalents, err := m.client.Talent.Query().All(ctx)
	if err != nil {
		m.message = fmt.Sprintf("Error fetching talents: %v", err)
		return
	}

	if len(allTalents) == 0 {
		m.message = "No talents available yet."
		return
	}

	var sb strings.Builder
	sb.WriteString("=== AVAILABLE TALENTS ===\n\n")
	for _, t := range allTalents {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", t.Name, t.Description))
		if t.Requirements != nil && len(t.Requirements) > 0 {
			sb.WriteString("  Requirements: ")
			reqs := []string{}
			for k, v := range t.Requirements {
				reqs = append(reqs, fmt.Sprintf("%s %d", k, v))
			}
			sb.WriteString(strings.Join(reqs, ", "))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	m.message = sb.String()
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

// handleSwapSkill handles the /swap-skill command
func (m *model) handleSwapSkill(args []string) {
	if m.client == nil {
		m.message = "Database not connected"
		return
	}

	if len(args) != 2 {
		m.message = "Usage: /swap-skill <old_skill> <new_skill>\nUse 'skills' to see available skills"
		return
	}

	oldSkillName := args[0]
	newSkillName := args[1]

	ctx := context.Background()

	// Get all skills and find the one with matching name
	allSkills, err := m.client.Skill.Query().All(ctx)
	if err != nil {
		m.message = fmt.Sprintf("Error fetching skills: %v", err)
		return
	}

	var newSkill *db.Skill
	for _, s := range allSkills {
		if strings.EqualFold(s.Name, newSkillName) {
			newSkill = s
			break
		}
	}

	if newSkill == nil {
		m.message = fmt.Sprintf("Skill '%s' not found. Use 'skills' to see available skills.", newSkillName)
		return
	}

	// Check if character has skill points
	// For now, we just do the swap (simplified version)
	m.message = fmt.Sprintf("Skill swap feature coming soon!\nWould swap '%s' for '%s' (Cost: %d skill points)", 
		oldSkillName, newSkill.Name, newSkill.Cost)
}

// handleSwapTalent handles the /swap-talent command
func (m *model) handleSwapTalent(args []string) {
	if m.client == nil {
		m.message = "Database not connected"
		return
	}

	if len(args) != 2 {
		m.message = "Usage: /swap-talent <old_talent> <new_talent>\nUse 'talents' to see available talents"
		return
	}

	oldTalentName := args[0]
	newTalentName := args[1]

	ctx := context.Background()

	// Get all talents and find the one with matching name
	allTalents, err := m.client.Talent.Query().All(ctx)
	if err != nil {
		m.message = fmt.Sprintf("Error fetching talents: %v", err)
		return
	}

	var newTalent *db.Talent
	for _, t := range allTalents {
		if strings.EqualFold(t.Name, newTalentName) {
			newTalent = t
			break
		}
	}

	if newTalent == nil {
		m.message = fmt.Sprintf("Talent '%s' not found. Use 'talents' to see available talents.", newTalentName)
		return
	}

	// Check requirements
	if newTalent.Requirements != nil {
		m.message = fmt.Sprintf("Talent swap feature coming soon!\nWould swap '%s' for '%s' (Requirements: %v)", 
			oldTalentName, newTalent.Name, newTalent.Requirements)
		return
	}

	m.message = fmt.Sprintf("Talent swap feature coming soon!\nWould swap '%s' for '%s'", 
		oldTalentName, newTalent.Name)
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

	// Room info at top
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