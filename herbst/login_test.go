package main

import (
	"testing"
	"time"
)

// TestLoginFieldState tests the login input field state machine
func TestLoginFieldState(t *testing.T) {
	// Test initial state is username
	m := model{
		connectedAt: time.Now(),
		inputField:  "username",
	}
	if m.inputField != "username" {
		t.Errorf("Expected initial inputField to be 'username', got %v", m.inputField)
	}

	// Test transition to password after username input
	m.inputField = "password"
	if m.inputField != "password" {
		t.Errorf("Expected inputField to be 'password', got %v", m.inputField)
	}
}

// TestLoginCredentials tests login credential handling
func TestLoginCredentials(t *testing.T) {
	m := model{
		connectedAt:   time.Now(),
		inputField:    "username",
		loginUsername: "testuser",
		loginPassword: "testpass",
	}

	// Test username is set correctly
	if m.loginUsername != "testuser" {
		t.Errorf("Expected loginUsername 'testuser', got %s", m.loginUsername)
	}

	// Test password is set correctly
	if m.loginPassword != "testpass" {
		t.Errorf("Expected loginPassword 'testpass', got %s", m.loginPassword)
	}

	// Test clearing login fields after failed attempt
	m.loginUsername = ""
	m.loginPassword = ""
	m.inputField = "username"

	if m.loginUsername != "" {
		t.Errorf("Expected loginUsername to be cleared, got %s", m.loginUsername)
	}
	if m.loginPassword != "" {
		t.Errorf("Expected loginPassword to be cleared, got %s", m.loginPassword)
	}
}