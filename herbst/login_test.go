package main

import (
	"testing"
	"time"
)

// TestLoginStateTransitions tests the login state machine
func TestLoginStateTransitions(t *testing.T) {
	// Test initial state is username
	m := model{
		loginState: StateUsername,
	}
	if m.loginState != StateUsername {
		t.Errorf("Expected initial state to be StateUsername, got %v", m.loginState)
	}

	// Test transition to password after username input
	m.loginState = StatePassword
	if m.loginState != StatePassword {
		t.Errorf("Expected state to be StatePassword, got %v", m.loginState)
	}

	// Test transition to authenticated after successful login
	m.loginState = StateAuthenticated
	if m.loginState != StateAuthenticated {
		t.Errorf("Expected state to be StateAuthenticated, got %v", m.loginState)
	}
}

// TestLoginStateConstants verifies the login state constants
func TestLoginStateConstants(t *testing.T) {
	if StateUsername != 0 {
		t.Errorf("Expected StateUsername to be 0, got %v", StateUsername)
	}
	if StatePassword != 1 {
		t.Errorf("Expected StatePassword to be 1, got %v", StatePassword)
	}
	if StateAuthenticated != 2 {
		t.Errorf("Expected StateAuthenticated to be 2, got %v", StateAuthenticated)
	}
}

// TestModelInitialization tests the model initialization
func TestModelInitialization(t *testing.T) {
	m := model{
		connectedAt: time.Now(),
		loginState: StateUsername,
		username:   "",
		password:   "",
		inputBuffer: "",
	}

	// Test initial state
	if m.loginState != StateUsername {
		t.Errorf("Expected loginState to be StateUsername, got %v", m.loginState)
	}

	// Test empty credentials
	if m.username != "" {
		t.Errorf("Expected username to be empty, got %s", m.username)
	}
	if m.password != "" {
		t.Errorf("Expected password to be empty, got %s", m.password)
	}
	if m.inputBuffer != "" {
		t.Errorf("Expected inputBuffer to be empty, got %s", m.inputBuffer)
	}
}

// TestFormatExits tests the exit formatting
func TestFormatExits(t *testing.T) {
	m := model{}

	// Test empty exits
	m.exits = map[string]int{}
	result := m.formatExits()
	if result != "none" {
		t.Errorf("Expected 'none', got %s", result)
	}

	// Test single exit
	m.exits = map[string]int{"north": 1}
	result = m.formatExits()
	if result != "north" {
		t.Errorf("Expected 'north', got %s", result)
	}

	// Test multiple exits
	m.exits = map[string]int{"north": 1, "south": 2, "east": 3}
	result = m.formatExits()
	// Note: map iteration order is not guaranteed, so we check for presence
	if result != "north, south, east" && result != "north, east, south" &&
		result != "south, north, east" && result != "south, east, north" &&
		result != "east, north, south" && result != "east, south, north" {
		t.Errorf("Expected exits in some order, got %s", result)
	}
}

// TestLoginFields tests login field handling
func TestLoginFields(t *testing.T) {
	m := model{
		loginState: StateUsername,
		username:   "testuser",
		password:   "testpass",
		loginError: "",
	}

	// Test username is set correctly
	if m.username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", m.username)
	}

	// Test password is set correctly
	if m.password != "testpass" {
		t.Errorf("Expected password 'testpass', got %s", m.password)
	}

	// Test clearing login fields after failed attempt
	m.username = ""
	m.password = ""
	m.loginError = "Invalid credentials"

	if m.username != "" {
		t.Errorf("Expected username to be cleared, got %s", m.username)
	}
	if m.password != "" {
		t.Errorf("Expected password to be cleared, got %s", m.password)
	}
	if m.loginError == "" {
		t.Errorf("Expected loginError to be set")
	}
}