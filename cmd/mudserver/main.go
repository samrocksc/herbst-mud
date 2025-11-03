package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/sam/makeathing/internal/actions"
	"github.com/sam/makeathing/internal/adapters"
	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/database"
	"github.com/sam/makeathing/internal/items"
	"github.com/sam/makeathing/internal/rooms"
	"github.com/sam/makeathing/internal/users"
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
	DBAdapter  *database.DBAdapter
	UseDB      bool
}

func main() {
	// Parse command line flags
	var useDB bool
	var mergeJSON bool
	flag.BoolVar(&useDB, "db", true, "Use database storage (false for JSON only)")
	flag.BoolVar(&mergeJSON, "merge-json", false, "Merge JSON files into the database (only applicable if --db is true)")
	flag.Parse()

	// Create the database adapter if using DB
	var dbAdapter *database.DBAdapter
	if useDB {
		var err error
		dbAdapter, err = database.NewDBAdapter("./data/mud.db")
		if err != nil {
			log.Fatalf("Failed to create database adapter: %v", err)
		}
		defer dbAdapter.Close()

		// Set default configuration
		err = dbAdapter.SetConfiguration("mud_name", "Makeathing MUD")
		if err != nil {
			log.Fatalf("Failed to set configuration: %v", err)
		}

		// If the merge-json flag is set, load JSON files and merge them into the database
		if mergeJSON {
			infoLog("Merging JSON files into the database...")
			err = mergeJSONDataToDB(dbAdapter)
			if err != nil {
				log.Fatalf("Failed to merge JSON data into database: %v", err)
			}
			infoLog("JSON merge complete.")
		}
	} else if mergeJSON {
		log.Fatalln("Cannot use --merge-json without --db flag. Please enable database storage with --db.")
	}

	// Create the game instance
	game := &Game{
		Characters: make(map[string]*characters.Character),
		Rooms:      make(map[string]*rooms.Room),
		DBAdapter:  dbAdapter,
		UseDB:      useDB,
	}

	// Initialize the game world
	initializeGameWorld(game)

	// Create the SSH adapter with the game
	sshAdapter := &adapters.SSHAdapter{Game: game, DBAdapter: dbAdapter}

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
	if useDB {
		infoLog("Using database storage")
	} else {
		infoLog("Using JSON storage only")
	}

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

// Ensure Game implements GameDBInterface
var _ database.GameDBInterface = (*Game)(nil)

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

	// Load rooms with resolved item/character references from JSON files
	rooms, err := rooms.LoadAllRoomsItemsAndCharactersWithReferences("./data/rooms", "./data/items", "./data/characters")
	if err != nil {
		log.Fatalf("Failed to load rooms with references from JSON: %v", err)
	}

	// Verify that all required rooms are loaded
	requiredRooms := []string{"start", "up_room", "nw_room", "e_room"}
	for _, roomID := range requiredRooms {
		if _, exists := rooms[roomID]; !exists {
			log.Fatalf("Required room '%s' not found in JSON files", roomID)
		}
	}

	game.Rooms = rooms
}

// mergeJSONDataToDB loads all JSON data and upserts it into the database
func mergeJSONDataToDB(dbAdapter *database.DBAdapter) error {
	// Load all rooms, items, and characters from JSON files
	roomsMap, err := rooms.LoadAllRoomsItemsAndCharactersWithReferences("./data/rooms", "./data/items", "./data/characters")
	if err != nil {
		return fmt.Errorf("failed to load rooms with references from JSON: %w", err)
	}

	// Extract items and characters from the loaded rooms for upsertion
	// This approach assumes items and characters are uniquely identified by ID across all rooms
	// and that rooms.LoadAllRoomsItemsAndCharactersWithReferences correctly resolves and provides them.
	allLoadedItemsMap := make(map[string]*items.Item)
	allLoadedCharactersMap := make(map[string]*characters.Character)

	for _, room := range roomsMap {
		for _, item := range room.ImmovableObjects {
			allLoadedItemsMap[item.ID] = &item
		}
		for _, item := range room.MovableObjects {
			allLoadedItemsMap[item.ID] = &item
		}
		for _, character := range room.NPCs {
			allLoadedCharactersMap[character.ID] = &character
		}
	}
	
	// Upsert items
	for _, item := range allLoadedItemsMap {
		existingItem, err := dbAdapter.GetItem(item.ID)
		if err != nil {
			return fmt.Errorf("failed to get item %s from DB: %w", item.ID, err)
		}
		if existingItem == nil {
			infoLog("Creating item %s in DB", item.ID)
			err = dbAdapter.CreateItem(item.ToJSON())
			if err != nil {
				return fmt.Errorf("failed to create item %s in DB: %w", item.ID, err)
			}
		} else {
			infoLog("Updating item %s in DB", item.ID)
			err = dbAdapter.UpdateItem(item.ToJSON())
			if err != nil {
				return fmt.Errorf("failed to update item %s in DB: %w", item.ID, err)
			}
		}
	}

	// Upsert characters
	for _, character := range allLoadedCharactersMap {
		existingCharacter, err := dbAdapter.GetCharacter(character.ID)
		if err != nil {
			return fmt.Errorf("failed to get character %s from DB: %w", character.ID, err)
		}
		if existingCharacter == nil {
			infoLog("Creating character %s in DB", character.ID)
			err = dbAdapter.CreateCharacter(character.ToJSON())
			if err != nil {
				return fmt.Errorf("failed to create character %s in DB: %w", character.ID, err)
			}
		} else {
			infoLog("Updating character %s in DB", character.ID)
			err = dbAdapter.UpdateCharacter(character.ToJSON())
			if err != nil {
				return fmt.Errorf("failed to update character %s in DB: %w", character.ID, err)
			}
		}
	}

	// Upsert rooms
	for _, room := range roomsMap {
		roomJSON := room.ToJSON()
		existingRoom, err := dbAdapter.GetRoom(room.ID)
		if err != nil {
			return fmt.Errorf("failed to get room %s from DB: %w", room.ID, err)
		}
		if existingRoom == nil {
			infoLog("Creating room %s in DB", room.ID)
			err = dbAdapter.CreateRoom(roomJSON)
			if err != nil {
				return fmt.Errorf("failed to create room %s in DB: %w", room.ID, err)
			}
		} else {
			infoLog("Updating room %s in DB", room.ID)
			err = dbAdapter.UpdateRoom(roomJSON)
			if err != nil {
				return fmt.Errorf("failed to update room %s in DB: %w", room.ID, err)
			}
		}
	}

	// Load and upsert users from JSON files
	usersMap, err := users.LoadAllUserJSONsFromDirectory("./data/users")
	if err != nil {
		return fmt.Errorf("failed to load users from JSON: %w", err)
	}

	for _, userJSON := range usersMap {
		// Check if user already exists by username
		existingUser, err := dbAdapter.GetUserByUsername(userJSON.Username)
		if err != nil {
			return fmt.Errorf("failed to get user %s from DB: %w", userJSON.Username, err)
		}
		if existingUser == nil {
			infoLog("Creating user %s in DB", userJSON.Username)
			err = dbAdapter.CreateUserFromJSON(userJSON)
			if err != nil {
				return fmt.Errorf("failed to create user %s in DB: %w", userJSON.Username, err)
			}
		} else {
			infoLog("Updating user %s in DB", userJSON.Username)
			err = dbAdapter.CreateUserFromJSON(userJSON)
			if err != nil {
				return fmt.Errorf("failed to update user %s in DB: %w", userJSON.Username, err)
			}
		}
	}

	// Load and upsert actions from JSON file
	actionsMap, err := actions.LoadAllActionsFromDirectory("./data")
	if err != nil {
		return fmt.Errorf("failed to load actions from JSON: %w", err)
	}

	// Process actions from the actions.json file
	if actionsJSON, exists := actionsMap["actions"]; exists {
		for _, action := range actionsJSON.Actions {
			existingAction, err := dbAdapter.GetAction(action.Name)
			if err != nil {
				return fmt.Errorf("failed to get action %s from DB: %w", action.Name, err)
			}
			if existingAction == nil {
				infoLog("Creating action %s in DB", action.Name)
				err = dbAdapter.CreateAction(&action)
				if err != nil {
					return fmt.Errorf("failed to create action %s in DB: %w", action.Name, err)
				}
			} else {
				infoLog("Updating action %s in DB", action.Name)
				err = dbAdapter.UpdateAction(&action)
				if err != nil {
					return fmt.Errorf("failed to update action %s in DB: %w", action.Name, err)
				}
			}
		}
	}

	return nil
}
