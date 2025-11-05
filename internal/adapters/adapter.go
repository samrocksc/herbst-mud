package adapters

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/sam/makeathing/internal/database"
	"github.com/sam/makeathing/internal/rooms"
)

// Adapter represents an interface for different connection types
type Adapter interface {
	HandleConnection(sess ssh.Session)
	SendMessage(sess ssh.Session, message string)
	GetInput(sess ssh.Session) string
}

// SSHAdapter implements the Adapter interface for SSH connections
type SSHAdapter struct {
	Game           GameInterface
	SessionManager *SessionManager
	DBAdapter      *database.DBAdapter
}

// HandleConnection handles an SSH connection
func (s *SSHAdapter) HandleConnection(sess ssh.Session) {
	infoLog("Handling connection for user: %s", sess.User())

	// Initialize session manager if not already done
	if s.SessionManager == nil {
		s.SessionManager = NewSessionManager(s.Game)
	}

	// AUTHENTICATION FLOW - Ask for username and password before creating session
	var authenticatedUser *database.User
	var authenticated bool
	
	if s.DBAdapter != nil {
		// Prompt for username
		s.SendMessage(sess, "Welcome to the MUD game!\n")
		s.SendMessage(sess, "Username: ")
		username := s.GetInput(sess)
		
		if username == "" {
			s.SendMessage(sess, "Invalid username. Disconnecting.\n")
			infoLog("Empty username provided, disconnecting")
			return
		}
		
		// Prompt for password
		s.SendMessage(sess, "Password: ")
		password := s.GetInput(sess)
		
		if password == "" {
			s.SendMessage(sess, "Invalid password. Disconnecting.\n")
			infoLog("Empty password provided for username %s, disconnecting", username)
			return
		}
		
		// Authenticate user
		infoLog("Authenticating user: %s", username)
		user, err := s.DBAdapter.AuthenticateUser(username, password)
		if err != nil {
			infoLog("Authentication error: %v", err)
			s.SendMessage(sess, "Authentication error. Disconnecting.\n")
			return
		}
		
		if user == nil {
			// User doesn't exist
			s.SendMessage(sess, "Invalid username or password. Disconnecting.\n")
			infoLog("Authentication failed for username: %s", username)
			return
		}
		
		// Authentication successful
		authenticatedUser = user
		authenticated = true
		infoLog("Authentication successful for user: %s (ID: %d, Character: %s)", 
			user.Username, user.ID, user.CharacterID)
	} else {
		infoLog("No database adapter available, proceeding without authentication")
		authenticated = true // Allow connection without auth if no DB
	}
	
	if !authenticated {
		return
	}

	// AUTHENTICATION PASSED - Create session and proceed to game
	sessionID := sess.Context().SessionID()
	debugLog("Creating session for ID: %s", sessionID)
	
	var characterID string
	if authenticatedUser != nil {
		characterID = authenticatedUser.CharacterID
	} else {
		// Default character if no authentication
		characterID = "char_nelly"
	}
	
	// Create session in database if we have user info and DB adapter
	if authenticatedUser != nil && s.DBAdapter != nil {
		err := s.DBAdapter.CreateSession(sessionID, authenticatedUser.ID, characterID, authenticatedUser.RoomID)
		if err != nil {
			infoLog("Failed to create database session: %v", err)
			s.SendMessage(sess, "Failed to create session. Disconnecting.\n")
			return
		}
		infoLog("Created database session for user %d, character %s", authenticatedUser.ID, characterID)
	}
	
	playerSession := s.SessionManager.CreatePlayerSession(sessionID)
	debugLog("Session created successfully")

	// Check if PTY is requested and handle terminal settings
	pty, _, isPty := sess.Pty()
	if isPty {
		// Log PTY information for debugging
		debugLog("PTY requested: %v, term: %s", isPty, pty.Term)
	}

	// Send welcome message with username if authenticated
	if authenticatedUser != nil {
		s.SendMessage(sess, fmt.Sprintf("Welcome %s!\n", authenticatedUser.Username))
	}
	s.SendMessage(sess, fmt.Sprintf("You are in: %s\n", playerSession.CurrentRoom.Description))
	s.SendMessage(sess, "Type 'help' for available commands.\n")
	s.SendMessage(sess, "\n> ")
	debugLog("Welcome message sent")

	// Handle user input (existing code from here on)
	debugLog("Starting input loop")
	reader := bufio.NewReader(sess)
	commandCount := 0
	for {
		debugLog("Waiting for input (command #%d)", commandCount+1)

		// Custom line reading to handle both \r and \n
		var input strings.Builder
		for {
			char, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					debugLog("EOF received, closing connection")
					return
				}
				infoLog("Error reading input: %v", err)
				return
			}

			// Echo the character back to the user (except for line terminators)
			if char != '\n' && char != '\r' {
				fmt.Fprintf(sess, "%c", char)
			}

			// Log each character for debugging (only in debug mode)
			debugLog("Received character: %q (byte: %d)", char, char)

			// Check for line terminators
			if char == '\n' {
				// Unix line ending
				debugLog("Received Unix line ending (LF)")
				// Echo a newline to the user
				fmt.Fprintf(sess, "\n")
				break
			} else if char == '\r' {
				// Windows line ending or Mac classic line ending
				debugLog("Received carriage return (CR)")

				// Try to peek at the next character to see if it's \n
				nextByte, err := reader.Peek(1)
				if err != nil {
					// If we can't peek (EOF or other error), treat \r as end of line
					debugLog("Error peeking after CR (%v), treating CR as end of line", err)
					// Echo a newline to the user
					fmt.Fprintf(sess, "\n")
					break
				}

				// If the next byte is \n, consume it
				if len(nextByte) > 0 && nextByte[0] == '\n' {
					// Consume the \n
					_, err := reader.ReadByte() // consume the \n
					if err != nil {
						debugLog("Error consuming LF after CR: %v", err)
					} else {
						debugLog("Consumed CRLF sequence")
					}
				}
				// Echo a newline to the user
				fmt.Fprintf(sess, "\n")
				break
			}

			// Add character to input
			input.WriteByte(char)
		}

		// Get the final input string
		inputStr := input.String()

		// Log the received input for debugging
		debugLog("Received input string: %q", inputStr)
		debugLog("Input length: %d", len(inputStr))

		if inputStr == "quit" || inputStr == "exit" {
			debugLog("Quit command received, closing connection")
			s.SendMessage(sess, "Goodbye!\n")
			break
		}

		// Process command
		debugLog("Processing command: %q", inputStr)
		s.processCommand(sess, inputStr, sessionID)
		debugLog("Command processed")
		s.SendMessage(sess, "\n> ")
		commandCount++
		debugLog("Prompt sent (command #%d completed)", commandCount)
	}

	// Clean up session
	debugLog("Cleaning up session")
	s.SessionManager.RemovePlayerSession(sessionID)
	debugLog("Session cleaned up")
}

// SendMessage sends a message to the client
func (s *SSHAdapter) SendMessage(sess ssh.Session, message string) {
	debugLog("Sending message: %q", message)
	fmt.Fprint(sess, message)
	debugLog("Message sent")
}

// GetInput gets input from the client
func (s *SSHAdapter) GetInput(sess ssh.Session) string {
	debugLog("Getting input from client")
	reader := bufio.NewReader(sess)

	// Read character by character and echo back
	var input strings.Builder
	for {
		char, err := reader.ReadByte()
		if err != nil {
			infoLog("Error reading input: %v", err)
			return ""
		}

		// Echo the character back to the user (except for line terminators)
		if char != '\n' && char != '\r' {
			fmt.Fprintf(sess, "%c", char)
		}

		// Check for line terminators
		if char == '\n' || char == '\r' {
			// Echo a newline to the user
			fmt.Fprintf(sess, "\n")

			// Handle CRLF sequence
			if char == '\r' {
				// Peek to see if next char is \n
				nextByte, err := reader.Peek(1)
				if err == nil && len(nextByte) > 0 && nextByte[0] == '\n' {
					// Consume the \n
					reader.ReadByte()
				}
			}
			break
		}

		// Add character to input
		input.WriteByte(char)
	}

	trimmed := strings.TrimSpace(input.String())
	debugLog("Input received: %q", trimmed)
	return trimmed
}

// processCommand processes user commands
func (s *SSHAdapter) processCommand(sess ssh.Session, command string, sessionID string) {
	// Debug: log that we're processing a command
	debugLog("Processing command: %q", command)

	// Handle empty commands
	if command == "" {
		debugLog("Empty command received, ignoring")
		return
	}

	// Handle abbreviated movement commands
	switch command {
	case "n":
		command = "north"
	case "s":
		command = "south"
	case "e":
		command = "east"
	case "w":
		command = "west"
	case "ne":
		command = "northeast"
	case "nw":
		command = "northwest"
	case "se":
		command = "southeast"
	case "sw":
		command = "southwest"
	case "u":
		command = "up"
	case "d":
		command = "down"
	case "l", "look":
		command = "look"
	}

	// Check if it's a movement command
	switch command {
	case "north", "south", "east", "west", "northeast", "northwest", "southeast", "southwest", "up", "down":
		// Convert string to Direction type
		direction := rooms.Direction(command)

		// Try to move the player
		debugLog("Attempting to move %s", command)
		nextRoom, err := s.SessionManager.MovePlayer(sessionID, direction)
		if err != nil {
			debugLog("Move failed: %v", err)
			s.SendMessage(sess, fmt.Sprintf("You cannot go %s.\n", command))
		} else {
			debugLog("Move successful to room: %s", nextRoom.ID)
			s.SendMessage(sess, fmt.Sprintf("You move %s.\n", command))
			s.SendMessage(sess, fmt.Sprintf("You are now in: %s\n", nextRoom.Description))
		}
	case "help":
		debugLog("Help command received")
		s.SendMessage(sess, "Available commands:\n")
		s.SendMessage(sess, "- help: Show this help message\n")
		s.SendMessage(sess, "- look, l: Look around the room\n")
		s.SendMessage(sess, "- n/s/e/w/ne/nw/se/sw/u/d or north/south/east/west/northeast/northwest/southeast/southwest/up/down: Move in a direction\n")
		s.SendMessage(sess, "- quit/exit: Exit the game\n")
	case "look":
		debugLog("Look command received")
		// Show current room description
		playerSession := s.SessionManager.GetPlayerSession(sessionID)
		if playerSession != nil {
			s.SendMessage(sess, fmt.Sprintf("You are in: %s\n", playerSession.CurrentRoom.Description))
			s.SendMessage(sess, fmt.Sprintf("Smells: %s\n", playerSession.CurrentRoom.Smells))

			// List available exits
			if len(playerSession.CurrentRoom.Exits) > 0 {
				s.SendMessage(sess, "Exits: ")
				exitList := []string{}
				for dir := range playerSession.CurrentRoom.Exits {
					exitList = append(exitList, string(dir))
				}
				s.SendMessage(sess, strings.Join(exitList, ", ")+"\n")
			}
		}
	case "echo":
		debugLog("Echo command received")
		// Debug command to echo back what was typed
		s.SendMessage(sess, fmt.Sprintf("You typed: %s\n", command))
	default:
		debugLog("Unknown command received: %s", command)
		s.SendMessage(sess, fmt.Sprintf("Unknown command: %s\n", command))
	}
}
