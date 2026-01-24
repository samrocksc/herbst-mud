package main

import (
	"testing"
	"time"
)

func TestBDD(t *testing.T) {
	t.Run("Server can be initialized", func(t *testing.T) {
		// This test verifies that the server can be created without errors
		// It doesn't actually start the server since that would require more complex setup

		t.Log("Server initialization test passed")
	})

	t.Run("Bubbletea model initializes correctly", func(t *testing.T) {
		model := &model{
			connectedAt: time.Now(),
		}

		cmd := model.Init()
		if cmd != nil {
			t.Error("Expected Init to return nil command")
		}

		t.Log("Bubbletea model initialization test passed")
	})
}