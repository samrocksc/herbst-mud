package configuration

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Configuration represents a game configuration
type Configuration struct {
	Schema string `json:"$schema"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
}

// LoadConfigurationFromJSON loads a configuration from a JSON file
func LoadConfigurationFromJSON(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadAllConfigurationsFromDirectory loads all configurations from JSON files in a directory
func LoadAllConfigurationsFromDirectory(directory string) (map[string]*Configuration, error) {
	configs := make(map[string]*Configuration)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" && file.Name() == "configuration.json" {
			filename := filepath.Join(directory, file.Name())
			config, err := LoadConfigurationFromJSON(filename)
			if err != nil {
				return nil, err
			}

			// Use the file name without extension as the key
			key := file.Name()[:len(file.Name())-5] // Remove .json extension
			configs[key] = config
		}
	}

	return configs, nil
}