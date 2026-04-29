package main

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/ssh"
	"herbst/combat"
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
	roomCharacters []RoomCharacter

	// Debug mode
	debugMode bool

	// Message history
	messageHistory  []string
	messageTypes    []string
	historyOffset   int
	maxHistory      int
	isScrolling     bool

	// Command history (for up/down arrow navigation)
	commandHistory []string
	historyIndex   int

	// Combat state
	inCombat           bool
	combatTarget       *RoomCharacter
	combatManager      *combat.CombatManager
	combatID           int
	combatLog          []string // Combat messages for display
	combatQueuedAction string   // Action queued for next tick
	combatJustStarted  bool     // Flag to start tick timer

	// Combat talents (slots 1-4)
	combatTalents []EquippedTalent

	// Classless combat skills (slots 1-5)
	combatSkills *CombatSkillState

	// Command registry
	commands *CommandRegistry

	// NPC skill cooldown (for enemy skills)
	npcSkillCooldown int

	// Skill selection state
	skillSelectSlot   int                 // Which slot we're selecting for (1-5)
	skillSelectCursor int                 // Cursor position in the list

	// Equipped potion (R slot)
	equippedPotion *EquippedPotion
}

// EquippedTalent represents a talent equipped in a combat slot
type EquippedTalent struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Slot          int    `json:"slot"`
	EffectType    string `json:"effectType"`   // heal|damage|dot|buff_armor|buff_dodge|buff_crit|debuff
	EffectValue   int    `json:"effectValue"`  // Amount: HP healed, damage, etc.
	EffectDuration int   `json:"effectDuration"` // Duration in ticks (0 = instant)
	Cooldown      int    `json:"cooldown"`
	ManaCost      int    `json:"manaCost"`
	StaminaCost   int    `json:"staminaCost"`
}

// EquippedPotion represents a potion equipped in the R slot
type EquippedPotion struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	EffectType    string `json:"effectType"`
	EffectValue   int    `json:"effectValue"`
	EffectDuration int   `json:"effectDuration"`
}

// combatTickMsg is sent on each combat tick interval
type combatTickMsg time.Time

// regenTickMsg is sent on each regeneration tick interval  
type regenTickMsg time.Time

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

// RoomCharacter represents a character in a room
type RoomCharacter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsNPC    bool   `json:"isNPC"`
	Level    int    `json:"level"`
	Class    string `json:"class"`
	Race     string `json:"race"`
	UserID   int    `json:"userId"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"maxHp"`
	XpValue  int    `json:"xpValue"` // XP awarded when this NPC is defeated
}

// Note: HiddenDetail is defined separately in examine_skill.go for that system
type InventoryItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ItemType    string `json:"itemType"`
	IsEquipped  bool   `json:"isEquipped"`
	Rarity      string `json:"rarity"`
}

var RESTAPIBase string

func init() {
	RESTAPIBase = os.Getenv("API_BASE_URL")
	if RESTAPIBase == "" {
		RESTAPIBase = "http://localhost:8080"
	}
}

// StartingRoomID is the ID of the room players start in
const StartingRoomID = 5

// Screen states
const (
	ScreenWelcome   = "welcome"
	ScreenLogin     = "login"
	ScreenRegister  = "register"
	ScreenPlaying   = "playing"
	ScreenProfile    = "profile"
	ScreenEditField  = "edit_field"
	ScreenCombat     = "combat"
	ScreenSkillSelect = "skill_select"
)
