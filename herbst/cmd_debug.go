package main

import (
	"strings"
)

// HandleDebug handles the debug command
func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		if m.debugMode {
			m.AppendMessage("Debug mode: ON (Room ID visible in status bar)", "info")
		} else {
			m.AppendMessage("Debug mode: OFF\nUsage: debug on | debug off", "info")
		}
		return
	}

	switch parts[1] {
	case "on", "true", "1", "yes":
		m.debugMode = true
		m.AppendMessage("Debug mode: ON (Room ID will show in status bar)", "success")
	case "off", "false", "0", "no":
		m.debugMode = false
		m.AppendMessage("Debug mode: OFF", "info")
	default:
		m.AppendMessage("Usage: debug on | debug off", "error")
	}
}