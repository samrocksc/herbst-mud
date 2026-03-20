package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"herbst/db"
	"herbst/dbinit"
)

func init() {
	os.Setenv("TERM", "xterm-256color")
}

func getDBConfig() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "herbst")
	password := getEnv("DB_PASSWORD", "herbst_password")
	dbname := getEnv("DB_NAME", "herbst_mud")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	client, err := db.Open("postgres", getDBConfig())
	if err != nil {
		log.Printf("Warning: failed connecting to postgres: %v", err)
	} else {
		defer client.Close()
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Printf("Warning: failed creating schema resources: %v", err)
		} else {
			log.Println("Database initialized successfully")
		}
		if err := dbinit.InitAdminUser(client); err != nil {
			log.Printf("Warning: failed to initialize admin user: %v", err)
		}
		if err := dbinit.InitGizmo(client); err != nil {
			log.Printf("Warning: failed to initialize Gizmo: %v", err)
		}
		if err := dbinit.InitWeapons(client); err != nil {
			log.Printf("Warning: failed to initialize weapons: %v", err)
		}
	}

	srv, err := wish.NewServer(
		wish.WithAddress(":4444"),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			logging.Middleware(),
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					log.Printf("New connection from %s", s.RemoteAddr().String())
					lipgloss.SetColorProfile(termenv.TrueColor)

					pty, winCh, ok := s.Pty()
					var initialWidth, initialHeight int
					if ok {
						initialWidth = pty.Window.Width
						initialHeight = pty.Window.Height
					} else {
						initialWidth = 80
						initialHeight = 24
					}

					ti := textinput.New()
					ti.Placeholder = "Enter choice..."
					ti.Focus()

					sp := spinner.New()
					sp.Spinner = spinner.Dot
					sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

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

					p := tea.NewProgram(m,
						tea.WithInput(s),
						tea.WithOutput(s),
						tea.WithAltScreen(),
					)

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
