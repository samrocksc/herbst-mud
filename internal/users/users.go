package users

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// User represents a user in the game
type User struct {
	ID          int
	CharacterID string
	RoomID      string
}

// UserJSON represents a user in JSON format
type UserJSON struct {
	Schema      string `json:"$schema"`
	ID          int    `json:"id"`
	CharacterID string `json:"characterId"`
	RoomID      string `json:"roomId"`
}

// LoadUserJSONFromJSON loads a UserJSON from a JSON file
func LoadUserJSONFromJSON(filename string) (*UserJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var userJSON UserJSON
	if err := json.Unmarshal(data, &userJSON); err != nil {
		return nil, err
	}

	return &userJSON, nil
}

// LoadAllUserJSONsFromDirectory loads all UserJSONs from JSON files in a directory
func LoadAllUserJSONsFromDirectory(directory string) (map[string]*UserJSON, error) {
	users := make(map[string]*UserJSON)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			userJSON, err := LoadUserJSONFromJSON(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load user JSON from %s: %w", filename, err)
			}
			// Use the user ID as the key
			key := fmt.Sprintf("user_%d", userJSON.ID)
			users[key] = userJSON
		}
	}

	return users, nil
}

// LoadUserFromJSON loads a user from a JSON file
func LoadUserFromJSON(filename string) (*User, error) {
	userJSON, err := LoadUserJSONFromJSON(filename)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:          userJSON.ID,
		CharacterID: userJSON.CharacterID,
		RoomID:      userJSON.RoomID,
	}

	return user, nil
}

// LoadAllUsersFromDirectory loads all users from JSON files in a directory
func LoadAllUsersFromDirectory(directory string) (map[string]*User, error) {
	users := make(map[string]*User)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filename := filepath.Join(directory, file.Name())
			user, err := LoadUserFromJSON(filename)
			if err != nil {
				return nil, err
			}
			// Use the user ID as the key
			key := fmt.Sprintf("user_%d", user.ID)
			users[key] = user
		}
	}

	return users, nil
}