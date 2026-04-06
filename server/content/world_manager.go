package content

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// WorldManager manages multiple MUD worlds
type WorldManager struct {
	mu     sync.RWMutex
	worlds map[string]*World
	managers map[string]*Manager  // world_id -> content manager
	defaultWorld string
}

// World represents a MUD world configuration
type World struct {
	ID             string            `yaml:"id" json:"id"`
	Name           string            `yaml:"name" json:"name"`
	Description    string            `yaml:"description" json:"description"`
	Status         string            `yaml:"status" json:"status"` // active, development, maintenance
	ContentPath    string            `yaml:"content_path" json:"content_path"`
	DatabasePrefix string            `yaml:"database_prefix" json:"database_prefix"`
	Features       []string          `yaml:"features" json:"features"`
	Settings       WorldSettings     `yaml:"settings" json:"settings"`
}

// WorldSettings contains world-specific game settings
type WorldSettings struct {
	PvPEnabled       bool    `yaml:"pvp_enabled" json:"pvp_enabled"`
	Permadeath       bool    `yaml:"permadeath" json:"permadeath"`
	XPMultiplier     float64 `yaml:"xp_multiplier" json:"xp_multiplier"`
	GoldMultiplier   float64 `yaml:"gold_multiplier" json:"gold_multiplier"`
}

// WorldRegistry represents the worlds.yaml structure
type WorldRegistry struct {
	Version      string   `yaml:"version"`
	Worlds       []World  `yaml:"worlds"`
	DefaultWorld string   `yaml:"default_world"`
	Shared       []string `yaml:"shared"`
}

// NewWorldManager creates a new world manager
func NewWorldManager(basePath string) (*WorldManager, error) {
	wm := &WorldManager{
		worlds:   make(map[string]*World),
		managers: make(map[string]*Manager),
	}

	// Load world registry
	registryPath := basePath + "/worlds.yaml"
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read worlds registry: %w", err)
	}

	var registry WorldRegistry
	if err := yaml.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse worlds registry: %w", err)
	}

	wm.defaultWorld = registry.DefaultWorld

	// Load each world
	for _, world := range registry.Worlds {
		if world.Status == "active" || world.Status == "development" {
			worldCopy := world
			wm.worlds[world.ID] = &worldCopy

			// Create content manager for this world
			contentMgr := NewManager(basePath + "/" + world.ContentPath)
			if err := contentMgr.LoadAll(); err != nil {
				return nil, fmt.Errorf("failed to load world %s: %w", world.ID, err)
			}
			wm.managers[world.ID] = contentMgr
		}
	}

	return wm, nil
}

// GetWorld returns a world by ID
func (wm *WorldManager) GetWorld(id string) (*World, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	world, exists := wm.worlds[id]
	return world, exists
}

// GetWorldManager returns the content manager for a world
func (wm *WorldManager) GetWorldManager(worldID string) (*Manager, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	mgr, exists := wm.managers[worldID]
	return mgr, exists
}

// GetDefaultWorld returns the default world ID
func (wm *WorldManager) GetDefaultWorld() string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.defaultWorld
}

// GetAllWorlds returns all loaded worlds
func (wm *WorldManager) GetAllWorlds() []*World {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	result := make([]*World, 0, len(wm.worlds))
	for _, world := range wm.worlds {
		result = append(result, world)
	}
	return result
}

// GetActiveWorlds returns worlds with active status
func (wm *WorldManager) GetActiveWorlds() []*World {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	var result []*World
	for _, world := range wm.worlds {
		if world.Status == "active" {
			result = append(result, world)
		}
	}
	return result
}

// ValidateWorldAccess checks if a user can access a world
func (wm *WorldManager) ValidateWorldAccess(worldID string) bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	world, exists := wm.worlds[worldID]
	if !exists {
		return false
	}

	// Only active worlds are accessible
	return world.Status == "active"
}

// GetWorldStats returns content stats for a specific world
func (wm *WorldManager) GetWorldStats(worldID string) (*Stats, bool) {
	mgr, exists := wm.GetWorldManager(worldID)
	if !exists {
		return nil, false
	}
	stats := mgr.GetStats()
	return &stats, true
}

// ReloadWorld reloads a world's content
func (wm *WorldManager) ReloadWorld(worldID string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	world, exists := wm.worlds[worldID]
	if !exists {
		return fmt.Errorf("world %s not found", worldID)
	}

	// Clear and reload
	if mgr, ok := wm.managers[worldID]; ok {
		mgr.Clear()
	}

	mgr := NewManager(wm.getBasePath() + "/" + world.ContentPath)
	if err := mgr.LoadAll(); err != nil {
		return fmt.Errorf("failed to reload world %s: %w", worldID, err)
	}

	wm.managers[worldID] = mgr
	return nil
}

// getBasePath returns the base path (simplified - should come from config)
func (wm *WorldManager) getBasePath() string {
	return "/home/sam/GitHub/herbst-mud/content"
}
