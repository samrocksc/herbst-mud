package main

import (
	"fmt"
	"strings"
)

// handleIgnoreCommand manages the ignore list
// ignore <player>   - add player to ignore list
// unignore <player> - remove player from ignore list
func (m *model) handleIgnoreCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to manage ignore list.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: ignore <player> or unignore <player>", "error")
		return
	}

	subcommand := strings.ToLower(args[0])
	if subcommand == "unignore" {
		if len(args) < 2 {
			m.AppendMessage("Usage: unignore <player>", "error")
			return
		}
		targetPlayer := args[1]
		m.removeIgnored(targetPlayer)
		return
	}

	// Default to ignore
	if len(args) < 2 {
		m.AppendMessage("Usage: ignore <player>", "error")
		return
	}
	targetPlayer := args[1]
	m.addIgnored(targetPlayer)
}

// addIgnored adds a player to the ignore list
func (m *model) addIgnored(playerName string) {
	if playerName == "" {
		m.AppendMessage("Usage: ignore <player>", "error")
		return
	}

	playerName = strings.TrimSpace(playerName)
	if playerName == m.currentCharacterName {
		m.AppendMessage("You cannot ignore yourself.", "error")
		return
	}

	// Check if already ignored
	for _, ignored := range m.ignoredPlayers {
		if strings.EqualFold(ignored, playerName) {
			m.AppendMessage(fmt.Sprintf("%s is already ignored.", playerName), "info")
			return
		}
	}

	m.ignoredPlayers = append(m.ignoredPlayers, playerName)
	m.AppendMessage(fmt.Sprintf("You are now ignoring %s.", playerName), "success")
}

// removeIgnored removes a player from the ignore list
func (m *model) removeIgnored(playerName string) {
	if playerName == "" {
		m.AppendMessage("Usage: unignore <player>", "error")
		return
	}

	playerName = strings.TrimSpace(playerName)

	for i, ignored := range m.ignoredPlayers {
		if strings.EqualFold(ignored, playerName) {
			m.ignoredPlayers = append(m.ignoredPlayers[:i], m.ignoredPlayers[i+1:]...)
			m.AppendMessage(fmt.Sprintf("You no longer ignore %s.", playerName), "success")
			return
		}
	}
	m.AppendMessage(fmt.Sprintf("%s is not on your ignore list.", playerName), "info")
}
