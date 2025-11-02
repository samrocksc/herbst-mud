package main

import (
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/logging"
	"github.com/sam/makeathing/internal/adapters"
	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/rooms"
	"github.com/sam/makeathing/internal/items"
)

// debugMode checks if debug mode is enabled
var debugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"

// infoLog always logs informational messages
func infoLog(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Game represents the main game state
type Game struct {
	Characters map[string]*characters.Character
	Rooms      map[string]*rooms.Room
	Admin      *characters.Character
}

func main() {
	// Create the game instance
	game := &Game{
		Characters: make(map[string]*characters.Character),
		Rooms:      make(map[string]*rooms.Room),
	}

	// Initialize the game world
	initializeGameWorld(game)

	// Create the SSH adapter with the game
	sshAdapter := &adapters.SSHAdapter{Game: game}

	// Create the wish server
	s, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath("./.ssh/term_info_ed25519"),
		wish.WithMiddleware(
			logging.Middleware(),
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					// Verbose connection logging (only in debug mode)
					if debugMode {
						infoLog("=== NEW CONNECTION ===")
						infoLog("Remote Address: %s", sess.RemoteAddr())
						infoLog("Local Address: %s", sess.LocalAddr())
						infoLog("User: %s", sess.User())
						infoLog("Session ID: %s", sess.Context().SessionID())
						
						// Check if PTY is requested
						pty, winCh, isPty := sess.Pty()
						if isPty {
							infoLog("PTY Requested - Terminal: %s, Columns: %d, Rows: %d", pty.Term, pty.Window.Width, pty.Window.Height)
							
							// Log window size changes
							go func() {
								for win := range winCh {
									infoLog("Window Resize - Columns: %d, Rows: %d", win.Width, win.Height)
								}
							}()
						} else {
							infoLog("No PTY requested")
						}
						
						// Log environment variables
						env := sess.Environ()
						if len(env) > 0 {
							infoLog("Environment Variables:")
							for _, e := range env {
								infoLog("  %s", e)
							}
						}
						
						infoLog("Raw Command: %s", sess.RawCommand())
						infoLog("Command: %v", sess.Command())
					}
					
					// Handle the connection
					sshAdapter.HandleConnection(sess)
					next(sess)
				}
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	infoLog("Starting SSH server on :2222")
	if err = s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

// GetRoom returns a room by its ID
func (g *Game) GetRoom(roomID string) *rooms.Room {
	return g.Rooms[roomID]
}

// GetStartingRoom returns the starting room
func (g *Game) GetStartingRoom() *rooms.Room {
	return g.Rooms["start"]
}

// initializeGameWorld sets up the initial game world
func initializeGameWorld(game *Game) {
	// Create the admin character
	admin := &characters.Character{
		Name:  "Admin",
		Race:  characters.Human,
		Class: characters.Warrior,
		Stats: characters.Stats{
			Strength:     20,
			Intelligence: 15,
			Dexterity:    18,
		},
		Health:   100,
		Mana:     50,
		Level:    10,
		IsVendor: false,
		IsNpc:    false,
	}
	game.Admin = admin

	// Create the first 4 rooms as specified in the rules
	// Starting room with exits to up, northwest, and east
	startRoom := &rooms.Room{
		ID:          "start",
		Description: "You are in the starting room. There are exits to the up, northwest, and east.",
		Exits: map[rooms.Direction]string{
			rooms.Up:        "up_room",
			rooms.Northwest: "nw_room",
			rooms.East:      "e_room",
		},
		ImmovableObjects: []items.Item{},
		MovableObjects:   []items.Item{},
		Smells:           "The air smells fresh and clean.",
		NPCs:             []characters.Character{},
	}

	upRoom := &rooms.Room{
		ID:          "up_room",
		Description: "You are in the room above the starting room. There is an exit back down.",
		Exits: map[rooms.Direction]string{
			rooms.Down: "start",
		},
		ImmovableObjects: []items.Item{},
		MovableObjects:   []items.Item{},
		Smells:           "The air smells musty.",
		NPCs:             []characters.Character{},
	}

	nwRoom := &rooms.Room{
		ID:          "nw_room",
		Description: "You are in the northwest room. There is an exit back southeast.",
		Exits: map[rooms.Direction]string{
			rooms.Southeast: "start",
		},
		ImmovableObjects: []items.Item{},
		MovableObjects:   []items.Item{},
		Smells:           "The air smells damp.",
		NPCs:             []characters.Character{},
	}

	eRoom := &rooms.Room{
		ID:          "e_room",
		Description: "You are in the east room. There is an exit back west.",
		Exits: map[rooms.Direction]string{
			rooms.West: "start",
		},
		ImmovableObjects: []items.Item{},
		MovableObjects:   []items.Item{},
		Smells:           "The air smells of incense.",
		NPCs:             []characters.Character{},
	}

	// Add rooms to the game
	game.Rooms["start"] = startRoom
	game.Rooms["up_room"] = upRoom
	game.Rooms["nw_room"] = nwRoom
	game.Rooms["e_room"] = eRoom
}