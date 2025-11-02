package database

import (
	"fmt"

	"github.com/sam/makeathing/internal/characters"
	"github.com/sam/makeathing/internal/items"
	"github.com/sam/makeathing/internal/rooms"
)

// DBAdapter implements the Adapter interface for database operations
type DBAdapter struct {
	db            *DB
	configRepo    *ConfigurationRepository
	sessionRepo   *SessionRepository
	userRepo      *UserRepository
	roomRepo      *RoomRepository
	characterRepo *CharacterRepository
	itemRepo      *ItemRepository
	actionRepo    *ActionRepository
	globalStateCharacterRepo *GlobalStateCharacterRepository
	globalStateRoomRepo      *GlobalStateRoomRepository
}

// NewDBAdapter creates a new database adapter
func NewDBAdapter(dbPath string) (*DBAdapter, error) {
	db, err := New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &DBAdapter{
		db:            db,
		configRepo:    NewConfigurationRepository(db),
		sessionRepo:   NewSessionRepository(db),
		userRepo:      NewUserRepository(db),
		roomRepo:      NewRoomRepository(db),
		characterRepo: NewCharacterRepository(db),
		itemRepo:      NewItemRepository(db),
		actionRepo:    NewActionRepository(db),
		globalStateCharacterRepo: NewGlobalStateCharacterRepository(db),
		globalStateRoomRepo:      NewGlobalStateRoomRepository(db),
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

// Character operations

// CreateCharacter creates a new character from a JSON character
func (d *DBAdapter) CreateCharacter(jsonCharacter *characters.CharacterJSON) error {
	character, err := CharacterFromJSONCharacter(jsonCharacter)
	if err != nil {
		return err
	}

	return d.characterRepo.Create(character)
}

// GetCharacter retrieves a character by ID
func (d *DBAdapter) GetCharacter(characterID string) (*characters.CharacterJSON, error) {
	character, err := d.characterRepo.GetByID(characterID)
	if err != nil {
		return nil, err
	}

	if character == nil {
		return nil, nil
	}

	return character.ToJSONCharacter()
}

// UpdateCharacter updates a character
func (d *DBAdapter) UpdateCharacter(jsonCharacter *characters.CharacterJSON) error {
	character, err := CharacterFromJSONCharacter(jsonCharacter)
	if err != nil {
		return err
	}

	return d.characterRepo.Update(character)
}

// DeleteCharacter deletes a character
func (d *DBAdapter) DeleteCharacter(characterID string) error {
	return d.characterRepo.Delete(characterID)
}

// Item operations

// CreateItem creates a new item from a JSON item
func (d *DBAdapter) CreateItem(jsonItem *items.ItemJSON) error {
	item, err := ItemFromJSONItem(jsonItem)
	if err != nil {
		return err
	}

	return d.itemRepo.Create(item)
}

// GetItem retrieves an item by ID
func (d *DBAdapter) GetItem(itemID string) (*items.ItemJSON, error) {
	item, err := d.itemRepo.GetByID(itemID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	return item.ToJSONItem()
}

// UpdateItem updates an item
func (d *DBAdapter) UpdateItem(jsonItem *items.ItemJSON) error {
	item, err := ItemFromJSONItem(jsonItem)
	if err != nil {
		return err
	}

	return d.itemRepo.Update(item)
}

// DeleteItem deletes an item
func (d *DBAdapter) DeleteItem(itemID string) error {
	return d.itemRepo.Delete(itemID)
}

// GetAllItems retrieves all items
func (d *DBAdapter) GetAllItems() ([]*items.ItemJSON, error) {
	dbItems, err := d.itemRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var jsonItems []*items.ItemJSON
	for _, item := range dbItems {
		jsonItem, err := item.ToJSONItem()
		if err != nil {
			return nil, err
		}
		jsonItems = append(jsonItems, jsonItem)
	}

	return jsonItems, nil
}

// GetAllCharacters retrieves all characters
func (d *DBAdapter) GetAllCharacters() ([]*characters.CharacterJSON, error) {
	dbCharacters, err := d.characterRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var jsonCharacters []*characters.CharacterJSON
	for _, character := range dbCharacters {
		jsonCharacter, err := character.ToJSONCharacter()
		if err != nil {
			return nil, err
		}
		jsonCharacters = append(jsonCharacters, jsonCharacter)
	}

	return jsonCharacters, nil
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

// Global State Character operations

// InitializeCharacterState initializes a new character state
func (d *DBAdapter) InitializeCharacterState(characterID, roomID string, health int) error {
	return d.globalStateCharacterRepo.InitializeCharacterState(characterID, roomID, health)
}

// GetCharacterState retrieves a character's state by character ID
func (d *DBAdapter) GetCharacterState(characterID string) (*GlobalStateCharacter, error) {
	return d.globalStateCharacterRepo.GetByCharacterID(characterID)
}

// UpdateCharacterRoom updates the room for a character
func (d *DBAdapter) UpdateCharacterRoom(characterID, roomID string) error {
	return d.globalStateCharacterRepo.UpdateRoom(characterID, roomID)
}

// UpdateCharacterHealth updates the health for a character
func (d *DBAdapter) UpdateCharacterHealth(characterID string, health int) error {
	return d.globalStateCharacterRepo.UpdateHealth(characterID, health)
}

// UpdateCharacterStatus updates the status for a character
func (d *DBAdapter) UpdateCharacterStatus(characterID, status string) error {
	return d.globalStateCharacterRepo.UpdateStatus(characterID, status)
}

// GetCharactersInRoom retrieves all character states in a specific room
func (d *DBAdapter) GetCharactersInRoom(roomID string) ([]*GlobalStateCharacter, error) {
	return d.globalStateCharacterRepo.GetCharactersInRoom(roomID)
}

// Global State Room operations

// InitializeRoomState initializes a new room state
func (d *DBAdapter) InitializeRoomState(roomID string) error {
	return d.globalStateRoomRepo.InitializeRoomState(roomID)
}

// GetRoomState retrieves a room's state by room ID
func (d *DBAdapter) GetRoomState(roomID string) (*GlobalStateRoom, error) {
	return d.globalStateRoomRepo.GetByRoomID(roomID)
}

// UpdateRoomNPCState updates the NPC state for a room
func (d *DBAdapter) UpdateRoomNPCState(roomID string, npcState []NPCState) error {
	return d.globalStateRoomRepo.UpdateNPCState(roomID, npcState)
}

// UpdateRoomItemState updates the item state for a room
func (d *DBAdapter) UpdateRoomItemState(roomID string, itemState []ItemState) error {
	return d.globalStateRoomRepo.UpdateItemState(roomID, itemState)
}

// GetRoomNPCState retrieves the NPC state for a room
func (d *DBAdapter) GetRoomNPCState(roomID string) ([]NPCState, error) {
	return d.globalStateRoomRepo.GetNPCState(roomID)
}

// GetRoomItemState retrieves the item state for a room
func (d *DBAdapter) GetRoomItemState(roomID string) ([]ItemState, error) {
	return d.globalStateRoomRepo.GetItemState(roomID)
}

// IncrementRoomPlayerCount increments the player count for a room
func (d *DBAdapter) IncrementRoomPlayerCount(roomID string) error {
	return d.globalStateRoomRepo.IncrementPlayerCount(roomID)
}

// DecrementRoomPlayerCount decrements the player count for a room
func (d *DBAdapter) DecrementRoomPlayerCount(roomID string) error {
	return d.globalStateRoomRepo.DecrementPlayerCount(roomID)
}

// GameDBInterface defines the methods needed from the game engine for database operations
type GameDBInterface interface {
	GetRoom(roomID string) *rooms.Room
	GetStartingRoom() *rooms.Room
}
