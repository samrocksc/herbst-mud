package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sam/makeathing/internal/characters"
)

// Action represents an action that can be performed
type Action struct {
	Name         string
	Type         characters.SkillType
	Description  string
	Requirements ActionRequirements
}

// ActionRequirements represents what is needed to perform an action
type ActionRequirements struct {
	MinLevel       int
	RequiredStats  characters.Stats
	RequiredSkills []string
}

// ActionsJSON represents the JSON structure for actions
type ActionsJSON struct {
	Schema  string   `json:"$schema"`
	Actions []Action `json:"actions"`
}

// LoadActionsFromJSON loads actions from a JSON file
func LoadActionsFromJSON(filename string) (*ActionsJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var actionsJSON ActionsJSON
	if err := json.Unmarshal(data, &actionsJSON); err != nil {
		return nil, err
	}

	return &actionsJSON, nil
}

// LoadAllActionsFromDirectory loads all actions from JSON files in a directory
func LoadAllActionsFromDirectory(directory string) (map[string]*ActionsJSON, error) {
	actions := make(map[string]*ActionsJSON)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" && file.Name() == "actions.json" {
			filename := filepath.Join(directory, file.Name())
			actionsJSON, err := LoadActionsFromJSON(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load actions JSON from %s: %w", filename, err)
			}
			// Use the file name without extension as the key
			key := file.Name()[:len(file.Name())-5] // Remove .json extension
			actions[key] = actionsJSON
		}
	}

	return actions, nil
}