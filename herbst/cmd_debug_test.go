package main

import (
	"strings"
	"testing"
)

func TestHandleDebugCommand(t *testing.T) {
	tests := []struct {
		name          string
		initialMode   bool
		command       string
		expectedMode  bool
		expectedInMsg string
	}{
		{
			name:          "debug without argument shows status off",
			initialMode:   false,
			command:       "debug",
			expectedMode:  false,
			expectedInMsg: "Debug mode: OFF",
		},
		{
			name:          "debug without argument shows status on",
			initialMode:   true,
			command:       "debug",
			expectedMode:  true,
			expectedInMsg: "Debug mode: ON",
		},
		{
			name:          "debug on enables debug mode",
			initialMode:   false,
			command:       "debug on",
			expectedMode:  true,
			expectedInMsg: "Debug mode: ON",
		},
		{
			name:          "debug off disables debug mode",
			initialMode:   true,
			command:       "debug off",
			expectedMode:  false,
			expectedInMsg: "Debug mode: OFF",
		},
		{
			name:          "debug true enables debug mode",
			initialMode:   false,
			command:       "debug true",
			expectedMode:  true,
			expectedInMsg: "Debug mode: ON",
		},
		{
			name:          "debug false disables debug mode",
			initialMode:   true,
			command:       "debug false",
			expectedMode:  false,
			expectedInMsg: "Debug mode: OFF",
		},
		{
			name:          "debug 1 enables debug mode",
			initialMode:   false,
			command:       "debug 1",
			expectedMode:  true,
			expectedInMsg: "Debug mode: ON",
		},
		{
			name:          "debug 0 disables debug mode",
			initialMode:   true,
			command:       "debug 0",
			expectedMode:  false,
			expectedInMsg: "Debug mode: OFF",
		},
		{
			name:          "debug yes enables debug mode",
			initialMode:   false,
			command:       "debug yes",
			expectedMode:  true,
			expectedInMsg: "Debug mode: ON",
		},
		{
			name:          "debug no disables debug mode",
			initialMode:   true,
			command:       "debug no",
			expectedMode:  false,
			expectedInMsg: "Debug mode: OFF",
		},
		{
			name:          "debug with invalid argument shows usage",
			initialMode:   false,
			command:       "debug invalid",
			expectedMode:  false,
			expectedInMsg: "Usage: debug on | debug off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{
				debugMode:      tt.initialMode,
				messageHistory: []string{},
				messageTypes:   []string{},
				maxHistory:     100,
			}
			m.handleDebugCommand(tt.command)

			if m.debugMode != tt.expectedMode {
				t.Errorf("Expected debugMode=%v, got %v", tt.expectedMode, m.debugMode)
			}

			// Check that the expected message was added
			if len(m.messageHistory) == 0 {
				t.Errorf("Expected message containing %q but messageHistory is empty", tt.expectedInMsg)
			} else if !strings.Contains(m.messageHistory[0], tt.expectedInMsg) {
				t.Errorf("Expected message containing %q, got %q", tt.expectedInMsg, m.messageHistory[0])
			}
		})
	}
}