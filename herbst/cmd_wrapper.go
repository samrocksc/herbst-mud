package main

import (
	"fmt"
	"strings"
)

// initCommands registers all game commands
func (m *model) initCommands() {
	if m.commands == nil {
		m.commands = NewCommandRegistry()
	}

	// Help
	m.commands.Register("help", m.handleHelpCommand, "?")

	// Dialog
	m.commands.Register("talk", m.handleTalkWrapperCommand)

	// Info commands
	m.commands.Register("whoami", m.handleWhoamiCommand)
	m.commands.Register("profile", m.handleProfileCommand, "p")

	// Action commands
	m.commands.Register("inventory", m.handleInventoryWrapperCommand, "inv", "i")
	m.commands.Register("equip", m.handleEquipWrapperCommand)
	m.commands.Register("unequip", m.handleUnequipWrapperCommand)
	m.commands.Register("loot", m.handleLootWrapperCommand)
	m.commands.Register("clear", m.handleClearCommand, "cls")
	m.commands.Register("quit", m.handleQuitCommand, "q")
}

func (m *model) handleInventoryWrapperCommand(_ *model, args []string) {
	m.handleInventoryCommand()
}

func (m *model) handleEquipWrapperCommand(_ *model, args []string) {
	m.handleEquipCommand(fmt.Sprintf("equip %s", strings.Join(args, " ")))
}

func (m *model) handleUnequipWrapperCommand(_ *model, args []string) {
	if len(args) == 0 {
		m.AppendMessage("Usage: unequip <slot|item>", "error")
		return
	}
	m.handleUnequipItemCommand(strings.Join(args, " "))
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

func (m *model) handleTalkWrapperCommand(_ *model, args []string) {
	cmd := "talk"
	if len(args) > 1 {
		cmd = "talk " + strings.Join(args[1:], " ")
	}
	m.handleTalkCommand(cmd)
}