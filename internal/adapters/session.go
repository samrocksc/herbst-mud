package adapters

import (
	"fmt"
	"sync"

	"github.com/sam/makeathing/internal/rooms"
)

// PlayerSession tracks a player's state during their session
type PlayerSession struct {
	CurrentRoom *rooms.Room
	// Add other player state here as needed
}

// SessionManager manages all active player sessions
type SessionManager struct {
	sessions map[string]*PlayerSession // session ID to player session
	game     GameInterface
	mu       sync.RWMutex
}

// GameInterface defines the methods needed from the game engine
type GameInterface interface {
	GetRoom(roomID string) *rooms.Room
	GetStartingRoom() *rooms.Room
}

// NewSessionManager creates a new session manager
func NewSessionManager(game GameInterface) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*PlayerSession),
		game:     game,
	}
}

// GetPlayerSession gets a player's session by session ID
func (sm *SessionManager) GetPlayerSession(sessionID string) *PlayerSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[sessionID]
}

// CreatePlayerSession creates a new player session
func (sm *SessionManager) CreatePlayerSession(sessionID string) *PlayerSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &PlayerSession{
		CurrentRoom: sm.game.GetStartingRoom(),
	}
	sm.sessions[sessionID] = session
	return session
}

// RemovePlayerSession removes a player's session
func (sm *SessionManager) RemovePlayerSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
}

// MovePlayer moves a player in the specified direction
func (sm *SessionManager) MovePlayer(sessionID string, direction rooms.Direction) (*rooms.Room, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("player session not found")
	}

	// Check if there's an exit in that direction
	nextRoomID, exists := session.CurrentRoom.Exits[direction]
	if !exists {
		return nil, fmt.Errorf("you cannot go %s", direction)
	}

	// Get the next room
	nextRoom := sm.game.GetRoom(nextRoomID)
	if nextRoom == nil {
		return nil, fmt.Errorf("destination room not found")
	}

	// Move the player
	session.CurrentRoom = nextRoom
	return nextRoom, nil
}
