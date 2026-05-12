package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// availableChannels lists all available chat channels
var availableChannels = []string{"chat", "newbie", "trade"}

// handleChannelCommand handles chat channel messages and management
// chat <msg>       - send to chat channel
// newbie <msg>     - send to newbie channel
// trade <msg>      - send to trade channel
// channels         - list available channels
// channel <name> on/off - toggle channel subscription
func (m *model) handleChannelCommand(args []string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You must be playing to use channels.", "error")
		return
	}

	if len(args) < 1 {
		m.AppendMessage("Usage: chat <msg>, newbie <msg>, trade <msg>, channels, or channel <name> on/off", "error")
		return
	}

	subcommand := strings.ToLower(args[0])

	switch subcommand {
	case "chat":
		if len(args) < 2 {
			m.AppendMessage("Usage: chat <message>", "error")
			return
		}
		m.sendToChannel("chat", strings.Join(args[1:], " "))

	case "newbie":
		if len(args) < 2 {
			m.AppendMessage("Usage: newbie <message>", "error")
			return
		}
		m.sendToChannel("newbie", strings.Join(args[1:], " "))

	case "trade":
		if len(args) < 2 {
			m.AppendMessage("Usage: trade <message>", "error")
			return
		}
		m.sendToChannel("trade", strings.Join(args[1:], " "))

	case "channels":
		m.listChannels()

	case "on", "off":
		if len(args) < 2 {
			m.AppendMessage("Usage: channel <name> on/off", "error")
			return
		}
		channelName := strings.ToLower(args[1])
		if subcommand == "on" {
			m.subscribeChannel(channelName)
		} else {
			m.unsubscribeChannel(channelName)
		}

	default:
		// Check if it's a channel toggle: channel <name> on/off
		if len(args) >= 2 && (strings.ToLower(args[1]) == "on" || strings.ToLower(args[1]) == "off") {
			channelName := subcommand
			if strings.ToLower(args[1]) == "on" {
				m.subscribeChannel(channelName)
			} else {
				m.unsubscribeChannel(channelName)
			}
		} else {
			m.AppendMessage("Usage: chat <msg>, newbie <msg>, trade <msg>, channels, or channel <name> on/off", "error")
		}
	}
}

// sendToChannel sends a message to a specific channel
func (m *model) sendToChannel(channelName, message string) {
	if message == "" {
		m.AppendMessage("Message cannot be empty.", "error")
		return
	}

	// Check if subscribed
	if !m.isChannelSubscribed(channelName) {
		m.AppendMessage(fmt.Sprintf("You are not subscribed to %s. Use 'channel %s on' to subscribe.", channelName, channelName), "error")
		return
	}

	// Check if ignored
	for _, ignored := range m.ignoredPlayers {
		if strings.EqualFold(ignored, channelName) {
			m.AppendMessage(fmt.Sprintf("You are ignoring %s.", channelName), "error")
			return
		}
	}

	// Call the chat API
	jsonData, _ := json.Marshal(map[string]interface{}{
		"channel":  channelName,
		"message": message,
	})
	resp, err := http.Post(RESTAPIBase+"/api/chat/channel", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		m.AppendMessage("Failed to send to channel.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if msg, ok := result["error"].(string); ok {
			m.AppendMessage(msg, "error")
		} else {
			m.AppendMessage("Failed to send to channel.", "error")
		}
		return
	}

	m.AppendMessage(fmt.Sprintf("[%s] %s", channelName, message), "info")
}

// listChannels displays available channels and subscription status
func (m *model) listChannels() {
	output := "Available Channels:\n"
	for _, ch := range availableChannels {
		status := "OFF"
		if m.isChannelSubscribed(ch) {
			status = "ON"
		}
		output += fmt.Sprintf("  %s - [%s]\n", ch, status)
	}
	output += "\nUse 'channel <name> on/off' to toggle subscriptions."
	m.AppendMessage(output, "info")
}

// subscribeChannel subscribes to a channel
func (m *model) subscribeChannel(channelName string) {
	// Validate channel name
	valid := false
	for _, ch := range availableChannels {
		if ch == channelName {
			valid = true
			break
		}
	}
	if !valid {
		m.AppendMessage(fmt.Sprintf("Unknown channel: %s. Available: chat, newbie, trade", channelName), "error")
		return
	}

	// Add to subscribed channels if not already there
	for _, ch := range m.activeChannels {
		if ch == channelName {
			m.AppendMessage(fmt.Sprintf("You are already subscribed to %s.", channelName), "info")
			return
		}
	}
	m.activeChannels = append(m.activeChannels, channelName)
	m.AppendMessage(fmt.Sprintf("Subscribed to %s.", channelName), "success")
}

// unsubscribeChannel unsubscribes from a channel
func (m *model) unsubscribeChannel(channelName string) {
	// Validate channel name
	valid := false
	for _, ch := range availableChannels {
		if ch == channelName {
			valid = true
			break
		}
	}
	if !valid {
		m.AppendMessage(fmt.Sprintf("Unknown channel: %s. Available: chat, newbie, trade", channelName), "error")
		return
	}

	// Remove from subscribed channels
	for i, ch := range m.activeChannels {
		if ch == channelName {
			m.activeChannels = append(m.activeChannels[:i], m.activeChannels[i+1:]...)
			m.AppendMessage(fmt.Sprintf("Unsubscribed from %s.", channelName), "success")
			return
		}
	}
	m.AppendMessage(fmt.Sprintf("You are not subscribed to %s.", channelName), "info")
}

// isChannelSubscribed checks if subscribed to a channel
func (m *model) isChannelSubscribed(channelName string) bool {
	for _, ch := range m.activeChannels {
		if ch == channelName {
			return true
		}
	}
	return false
}
