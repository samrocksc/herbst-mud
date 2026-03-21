package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFuzzyWordMatch(t *testing.T) {
	tests := []struct {
		name     string
		nameArg  string
		target   string
		expected bool
	}{
		{
			name:     "exact match",
			nameArg:  "Gandalf",
			target:   "gandalf",
			expected: true,
		},
		{
			name:     "partial match at start",
			nameArg:  "Grand Ol' Man",
			target:   "grand",
			expected: true,
		},
		{
			name:     "fuzzy word match",
			nameArg:  "Grand Ol' Man",
			target:   "grand man",
			expected: true,
		},
		{
			name:     "fuzzy word match different order",
			nameArg:  "The Grand Hotel",
			target:   "hotel grand",
			expected: true,
		},
		{
			name:     "no match",
			nameArg:  "Gandalf",
			target:   "frodo",
			expected: false,
		},
		{
			name:     "partial word not found",
			nameArg:  "Sword of Power",
			target:   "axe",
			expected: false,
		},
		{
			name:     "case insensitive",
			nameArg:  "MIGHTY SWORD",
			target:   "mighty",
			expected: true,
		},
		{
			name:     "empty target matches",
			nameArg:  "Gandalf",
			target:   "",
			expected: true,
		},
		{
			name:     "multi-word target all match",
			nameArg:  "Rusty Iron Sword",
			target:   "rusty sword",
			expected: true,
		},
		{
			name:     "multi-word target partial match fails",
			nameArg:  "Rusty Iron Sword",
			target:   "rusty axe",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fuzzyWordMatch(tt.nameArg, tt.target)
			if result != tt.expected {
				t.Errorf("fuzzyWordMatch(%q, %q) = %v, expected %v", tt.nameArg, tt.target, result, tt.expected)
			}
		})
	}
}

func TestHTTPGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	resp, err := httpGet(server.URL)
	if err != nil {
		t.Fatalf("httpGet returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHTTPPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	resp, err := httpPost(server.URL, `{"test":"data"}`)
	if err != nil {
		t.Fatalf("httpPost returned error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestIOReadAll(t *testing.T) {
	reader := strings.NewReader("test content")
	data, err := ioReadAll(reader)
	if err != nil {
		t.Fatalf("ioReadAll returned error: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("Expected 'test content', got %q", string(data))
	}
}