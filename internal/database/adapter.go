package database

import (
	"fmt"

	"github.com/sam/makeathing/internal/adapters"
	"github.com/sam/makeathing/internal/rooms"
)

// DBAdapter implements the Adapter interface for database operations
type DBAdapter struct {
	db *DB
	configRepo *ConfigurationRepository
	sessionRepo *SessionRepository
	userRepo    *UserRepository
	roomRepo    *RoomRepository
}

// NewDBAdapter creates a new database adapter
func NewDBAdapter(dbPath string) (*DBAdapter, error) {
	db, err := New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &DBAdapter{
		db:          db,
		configRepo:  NewConfigurationRepository(db),
		sessionRepo: NewSessionRepository(db),
		userRepo:    NewUserRepository(db),
		roomRepo:    NewRoomRepository(db),
	}, nil
}

// Close closes the database connection
func (d *DBAdapter) Close() error {
	return d.db.Close()
}

// GetConfiguration retrieves the game configuration
func (d *DBAdapter) GetConfiguration(name string) (*Configuration, error) {
	return d.configRepo.GetByName(name)
}

// SetConfiguration sets the game configuration
func (d *DBAdapter) SetConfiguration(name, value string) error {
	config, err := d.configRepo.GetByName(value)
	if err != nil {
		return err
	}

	if config == nil {
		_, err := d.configRepo.Create(value)
		return err
	}

	// If a configuration with this name already exists, we don't need to create another one
	return nil
}

// CreateSession creates a new session
func (d *DBAdapter) CreateSession(sessionID string, userID int, characterID, roomID string) error {
	session := &Session{
		ID:          sessionID,
		UserID:      userID,
		CharacterID: characterID,
		RoomID:      roomID,
	}
	
	return d.sessionRepo.Create(session)
}

// GetSession retrieves a session by ID
func (d *DBAdapter) GetSession(sessionID string) (*Session, error) {
	return d.sessionRepo.GetByID(sessionID)
}

// UpdateSession updates a session
func (d *DBAdapter) UpdateSession(session *Session) error {
	return d.sessionRepo.Update(session)
}

// DeleteSession deletes a session
func (d *DBAdapter) DeleteSession(sessionID string) error {
	return d.sessionRepo.Delete(sessionID)
}

// CreateUser creates a new user
func (d *DBAdapter) CreateUser(characterID, roomID string) (int, error) {
	user := &User{
		CharacterID: characterID,
		RoomID:      roomID,
	}
	
	err := d.userRepo.Create(user)
	if err != nil {
		return 0, err
	}
	
	return user.ID, nil
}

// GetUser retrieves a user by ID
func (d *DBAdapter) GetUser(userID int) (*User, error) {
	return d.userRepo.GetByID(userID)
}

// GetUserByCharacterID retrieves a user by character ID
func (d *DBAdapter) GetUserByCharacterID(characterID string) (*User, error) {
	return d.userRepo.GetByCharacterID(characterID)
}

// UpdateUser updates a user
func (d *DBAdapter) UpdateUser(user *User) error {
	return d.userRepo.Update(user)
}

// DeleteUser deletes a user
func (d *DBAdapter) DeleteUser(userID int) error {
	return d.userRepo.Delete(userID)
}

// Room operations

// CreateRoom creates a new room from a JSON room
func (d *DBAdapter) CreateRoom(jsonRoom *rooms.RoomJSON) error {
	room, err := RoomFromJSONRoom(jsonRoom)
	if err != nil {
		return err
	}
	
	return d.roomRepo.Create(room)
}

// GetRoom retrieves a room by ID
func (d *DBAdapter) GetRoom(roomID string) (*rooms.RoomJSON, error) {
	room, err := d.roomRepo.GetByID(roomID)
	if err != nil {
		return nil, err
	}
	
	if room == nil {
		return nil, nil
	}
	
	return room.ToJSONRoom()
}

// UpdateRoom updates a room
func (d *DBAdapter) UpdateRoom(jsonRoom *rooms.RoomJSON) error {
	room, err := RoomFromJSONRoom(jsonRoom)
	if err != nil {
		return err
	}
	
	return d.roomRepo.Update(room)
}

// DeleteRoom deletes a room
func (d *DBAdapter) DeleteRoom(roomID string) error {
	return d.roomRepo.Delete(roomID)
}

// GetAllRooms retrieves all rooms
func (d *DBAdapter) GetAllRooms() ([]*rooms.RoomJSON, error) {
	dbRooms, err := d.roomRepo.GetAll()
	if err != nil {
		return nil, err
	}
	
	var jsonRooms []*rooms.RoomJSON
	for _, room := range dbRooms {
		jsonRoom, err := room.ToJSONRoom()
		if err != nil {
			return nil, err
		}
		jsonRooms = append(jsonRooms, jsonRoom)
	}
	
	return jsonRooms, nil
}

// GameDBInterface defines the methods needed from the game engine for database operations
type GameDBInterface interface {
	GetRoom(roomID string) *rooms.Room
	GetStartingRoom() *rooms.Room
}

// SessionManagerWithDB wraps the existing SessionManager to add database persistence
type SessionManagerWithDB struct {
	*adapters.SessionManager
	dbAdapter *DBAdapter
	game      GameDBInterface
}

// NewSessionManagerWithDB creates a new session manager with database persistence
func NewSessionManagerWithDB(game GameDBInterface, dbAdapter *DBAdapter) *SessionManagerWithDB {
	return &SessionManagerWithDB{
		SessionManager: adapters.NewSessionManager(game),
		dbAdapter:      dbAdapter,
		game:           game,
	}
}