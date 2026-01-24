package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress(":4444"),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			logging.Middleware(),
			bubbleteaMiddleware,
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

// bubbleteaMiddleware handles the bubbletea UI for SSH sessions
func bubbleteaMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		// Log connection information
		log.Printf("New connection from %s", s.RemoteAddr().String())

		// Create and run the bubbletea program
		p := tea.NewProgram(
			&model{
				connectedAt: time.Now(),
				session:     s,
			},
			tea.WithInput(s),
			tea.WithOutput(s),
		)

		// Run the program and handle errors
		if _, err := p.Run(); err != nil {
			log.Printf("Bubbletea error: %v", err)
		}

		log.Printf("Connection from %s closed", s.RemoteAddr().String())
	}
}

type model struct {
	connectedAt time.Time
	session     ssh.Session
	width       int
	height      int
	err         error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress 'q' or Ctrl+C to quit\n", m.err)
	}

	s := "Welcome to Herbst MUD!\n\n"
	s += fmt.Sprintf("Connected at: %s\n", m.connectedAt.Format(time.RFC1123))
	s += fmt.Sprintf("Client: %s\n", m.session.RemoteAddr().String())
	s += "\nPress 'q' or Ctrl+C to quit\n"

	// Center the text in the terminal
	if m.width > 0 && m.height > 0 {
		lines := []string{}
		for _, line := range []string{s} {
			padding := (m.width - len(line)) / 2
			if padding > 0 {
				lines = append(lines, fmt.Sprintf("%*s%s", padding, "", line))
			} else {
				lines = append(lines, line)
			}
		}
		s = ""
		for _, line := range lines {
			s += line + "\n"
		}
	}

	return s
}