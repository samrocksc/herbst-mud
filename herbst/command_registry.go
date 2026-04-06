package main

import (
	"fmt"
	"strings"
)

// CommandHandler is a function that handles a command
type CommandHandler func(*model, []string)

// CommandRegistry manages command handlers and aliases
type CommandRegistry struct {
	handlers map[string]CommandHandler
	aliases  map[string]string
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
		aliases:  make(map[string]string),
	}
}

// Register adds a command handler with optional aliases
func (r *CommandRegistry) Register(name string, handler CommandHandler, aliases ...string) {
	r.handlers[name] = handler
	for _, alias := range aliases {
		r.aliases[alias] = name
	}
}

// Execute runs a command handler if found
// Returns true if command was handled, false otherwise
func (r *CommandRegistry) Execute(m *model, cmd string, args []string) bool {
	canonical := r.resolveCommand(cmd)
	if handler, ok := r.handlers[canonical]; ok {
		handler(m, args)
		return true
	}
	return false
}

// resolveCommand converts alias to canonical name
func (r *CommandRegistry) resolveCommand(cmd string) string {
	if canonical, ok := r.aliases[cmd]; ok {
		return canonical
	}
	return cmd
}

// IsRegistered checks if a command is registered
func (r *CommandRegistry) IsRegistered(cmd string) bool {
	canonical := r.resolveCommand(cmd)
	_, ok := r.handlers[canonical]
	return ok
}

// initCommands registers all game commands
func (m *model) initCommands() {
	if m.commands == nil {
		m.commands = NewCommandRegistry()
	}

	// Help
	m.commands.Register("help", m.handleHelpCommand, "?")

	// Info commands
	m.commands.Register("whoami", m.handleWhoamiCommand)
	m.commands.Register("profile", m.handleProfileCommand, "p")

	// Action commands
	m.commands.Register("equip", m.handleEquipWrapperCommand)
	m.commands.Register("loot", m.handleLootWrapperCommand)
	m.commands.Register("clear", m.handleClearCommand, "cls")
	m.commands.Register("quit", m.handleQuitCommand, "q")
}

// Command handler implementations
func (m *model) handleHelpCommand(_ *model, args []string) {
	m.AppendMessage(`Commands:
  n/north, s/south, e/east, w/west - Move (cardinal)
  ne/northeast, se/southeast, sw/southwest, nw/northwest - Move (ordinal)
  u/up, d/down - Move (vertical)
  look/l [target] - Look around (or examine: look <target>, look at <target>)
  attack/a/kill/fight <target> - Attack a target
  loot [corpse] - Loot a corpse or all corpses in room
  ctrl+p - Scroll output up (older messages)
  ctrl+n - Scroll output down (newer messages)
  exits/x - Show exits
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  inventory/i - Show your inventory
  equip - Manage equipment slots (talents 1-4, potion R)
  quests/q - Show your quest log
  whoami - Show your info
  profile/p - Edit character profile
  skills - Show your equipped combat skills
  skill slot <1-5> - Select a skill for a slot
  skill all - Show available classless skills
  skill swap <s1> <s2> - Swap skills between slots
  talents - Show your talents
  debug - Toggle debug mode
  clear/cls - Clear screen
  quit - Exit game`, "info")
}

func (m *model) handleWhoamiCommand(_ *model, args []string) {
	m.AppendMessage(fmt.Sprintf("=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
		m.currentUserName, m.currentUserID, m.roomName,
		m.characterLevel, m.characterExperience,
		StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)), "info")
}

func (m *model) handleProfileCommand(_ *model, args []string) {
	m.screen = ScreenProfile
	m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
	m.menuCursor = 0
	m.AppendMessage("", "")
}

func (m *model) handleEquipWrapperCommand(_ *model, args []string) {
	m.handleEquipCommand(fmt.Sprintf("equip %s", strings.Join(args, " ")))
}

func (m *model) handleLootWrapperCommand(_ *model, args []string) {
	cmd := "loot"
	if len(args) > 0 {
		cmd = fmt.Sprintf("loot %s", strings.Join(args, " "))
	}
	m.handleLootCommand(cmd)
}

func (m *model) handleClearCommand(_ *model, args []string) {
	m.messageHistory = nil
	m.messageTypes = nil
	m.historyOffset = 0
	m.isScrolling = false
	m.inputBuffer = ""
}

func (m *model) handleQuitCommand(_ *model, args []string) {
	m.AppendMessage("Thanks for playing! Goodbye!", "success")
	m.inputBuffer = ""
}
