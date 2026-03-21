package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/ssh"
	"herbst/db"
)

type model struct {
	connectedAt time.Time
	session     ssh.Session
	client      *db.Client
	width       int
	height      int
	err         error

	// Screen state
	screen string

	// Auth state
	currentUserID   int
	currentUserName string

	// Input handling
	textInput textinput.Model

	// Login/Register input state
	inputField    string // "username" or "password"
	loginUsername string
	loginPassword string

	// Player state
	currentRoom int
	roomName    string
	roomDesc    string
	exits       map[string]int
	inputBuffer string
	message     string
	messageType string // "success", "error", "info"

	// Menu navigation state
	menuCursor int
	menuItems  []string

	// Character state
	currentCharacterID   int
	currentCharacterName string
	characterGender      string
	characterDescription string
	characterHP          int
	characterMaxHP       int
	characterStamina     int
	characterMaxStamina  int
	characterMana        int
	characterMaxMana     int
	characterLevel       int
	characterExperience  int

	// Profile editing state
	editField string // "gender" or "description"

	// Loading state
	spinner        spinner.Model
	isLoading      bool
	loadingMessage string

	// Scrollable viewport
	viewport viewport.Model

	// Room tracking
	visitedRooms map[int]bool
	knownExits   map[string]bool

	// Room items & characters
	roomItems      []RoomItem
	roomCharacters []roomCharacter

	// Debug mode
	debugMode bool

	// Message history
	messageHistory  []string
	messageTypes    []string
	historyOffset   int
	maxHistory      int
	isScrolling     bool
}

// RoomItem represents an item in a room for display
type RoomItem struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ExamineDesc     string         `json:"examineDesc"`
	HiddenDetails   []HiddenDetail `json:"hiddenDetails"`
	HiddenThreshold int            `json:"hiddenThreshold"`
	IsImmovable     bool           `json:"isImmovable"`
	Color           string         `json:"color"`
	IsVisible       bool           `json:"isVisible"`
	ItemType        string         `json:"itemType"`
	Weight          int            `json:"weight"`
	ItemDamage      int            `json:"itemDamage"`
	ItemDurability  int            `json:"itemDurability"`
	RevealCondition map[string]any `json:"revealCondition"`
}

// roomCharacter represents a character (NPC or player) in a room
type roomCharacter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsNPC  bool   `json:"isNPC"`
	Level  int    `json:"level"`
	Class  string `json:"class"`
	Race   string `json:"race"`
	UserID int    `json:"userId"`
}

// inventoryItem represents an item in the player's inventory
type inventoryItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ItemType    string `json:"itemType"`
	IsEquipped  bool   `json:"isEquipped"`
	Rarity      string `json:"rarity"`
}

// RESTAPIBase is the base URL for the REST API
var RESTAPIBase = "http://localhost:8080"

// StartingRoomID is the ID of the room players start in
const StartingRoomID = 5

// Screen states
const (
	ScreenWelcome   = "welcome"
	ScreenLogin     = "login"
	ScreenRegister  = "register"
	ScreenPlaying   = "playing"
	ScreenProfile   = "profile"
	ScreenEditField = "edit_field"
)
