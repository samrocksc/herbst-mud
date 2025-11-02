package configuration

import (
	"testing"
)

func TestLoadConfigurationFromJSON(t *testing.T) {
	// Load the configuration file
	config, err := LoadConfigurationFromJSON("../../data/configuration.json")
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	// Check that the configuration was loaded correctly
	if config.Name != "herbst" {
		t.Errorf("Expected name 'herbst', got '%s'", config.Name)
	}

	if config.ID != 1 {
		t.Errorf("Expected ID 1, got %d", config.ID)
	}

	if config.Schema != "../schemas/configuration.schema.json" {
		t.Errorf("Expected schema '../schemas/configuration.schema.json', got '%s'", config.Schema)
	}
}

func TestLoadAllConfigurationsFromDirectory(t *testing.T) {
	// Load all configurations from the data directory
	configs, err := LoadAllConfigurationsFromDirectory("../../data")
	if err != nil {
		t.Fatalf("Error loading configurations: %v", err)
	}

	// Check that we loaded exactly one configuration
	if len(configs) != 1 {
		t.Errorf("Expected 1 configuration, got %d", len(configs))
	}

	// Check that the configuration was loaded correctly
	config, exists := configs["configuration"]
	if !exists {
		t.Fatal("Expected configuration with key 'configuration' not found")
	}

	if config.Name != "herbst" {
		t.Errorf("Expected name 'herbst', got '%s'", config.Name)
	}

	if config.ID != 1 {
		t.Errorf("Expected ID 1, got %d", config.ID)
	}
}