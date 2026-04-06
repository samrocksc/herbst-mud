// Package content provides runtime content loading for Herbst MUD
// Supports YAML/JSON files with validation and hot-reloading
package content

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Manager is the central content registry
type Manager struct {
	basePath string

	mu sync.RWMutex

	// Content registries
	Skills  *SkillRegistry
	NPCs    *NPCRegistry
	Items   *ItemRegistry
	Rooms   *RoomRegistry
	Quests  *QuestRegistry
}

// NewManager creates a new content manager
func NewManager(basePath string) *Manager {
	return &Manager{
		basePath: basePath,
		Skills:   NewSkillRegistry(),
		NPCs:     NewNPCRegistry(),
		Items:    NewItemRegistry(),
		Rooms:    NewRoomRegistry(),
		Quests:   NewQuestRegistry(),
	}
}

// LoadAll loads all content from the base path
func (m *Manager) LoadAll() error {
	// Load in dependency order
	if err := m.loadSkills(); err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}
	if err := m.loadItems(); err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}
	if err := m.loadNPCs(); err != nil {
		return fmt.Errorf("failed to load NPCs: %w", err)
	}
	if err := m.loadRooms(); err != nil {
		return fmt.Errorf("failed to load rooms: %w", err)
	}
	if err := m.loadQuests(); err != nil {
		return fmt.Errorf("failed to load quests: %w", err)
	}
	return nil
}

// Validate runs validation on all loaded content
func (m *Manager) Validate() []ValidationError {
	var errors []ValidationError

	m.mu.RLock()
	defer m.mu.RUnlock()

	errors = append(errors, m.Skills.Validate()...)
	errors = append(errors, m.Items.Validate()...)
	errors = append(errors, m.NPCs.Validate(m.Skills, m.Items)...)
	errors = append(errors, m.Rooms.Validate()...)
	errors = append(errors, m.Quests.Validate(m.Skills, m.NPCs, m.Items)...)

	return errors
}

// Reload clears and reloads all content
func (m *Manager) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear all registries
	m.Skills.Clear()
	m.NPCs.Clear()
	m.Items.Clear()
	m.Rooms.Clear()
	m.Quests.Clear()

	return m.LoadAll()
}

// loadSkills loads skill definitions from content/default/skills/
func (m *Manager) loadSkills() error {
	skillsPath := filepath.Join(m.basePath, "skills")

	return filepath.WalkDir(skillsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var skill SkillDef
		if err := yaml.Unmarshal(data, &skill); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		if err := m.Skills.Register(&skill); err != nil {
			return fmt.Errorf("failed to register skill from %s: %w", path, err)
		}

		return nil
	})
}

// loadItems loads item definitions from content/default/items/
func (m *Manager) loadItems() error {
	itemsPath := filepath.Join(m.basePath, "items")

	return filepath.WalkDir(itemsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var item ItemDef
		if err := yaml.Unmarshal(data, &item); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		if err := m.Items.Register(&item); err != nil {
			return fmt.Errorf("failed to register item from %s: %w", path, err)
		}

		return nil
	})
}

// loadNPCs loads NPC templates from content/default/npcs/templates/
func (m *Manager) loadNPCs() error {
	npcsPath := filepath.Join(m.basePath, "npcs", "templates")

	return filepath.WalkDir(npcsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var npc NPCTemplate
		if err := yaml.Unmarshal(data, &npc); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		if err := m.NPCs.Register(&npc); err != nil {
			return fmt.Errorf("failed to register NPC from %s: %w", path, err)
		}

		return nil
	})
}

// loadRooms loads room definitions from content/default/rooms/
func (m *Manager) loadRooms() error {
	roomsPath := filepath.Join(m.basePath, "rooms")

	return filepath.WalkDir(roomsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Try to parse as individual room first
		var room RoomDef
		if err := yaml.Unmarshal(data, &room); err == nil && room.ID != "" {
			if err := m.Rooms.Register(&room); err != nil {
				return fmt.Errorf("failed to register room from %s: %w", path, err)
			}
			return nil
		}

		// Fall back to AreaDef format
		var area AreaDef
		if err := yaml.Unmarshal(data, &area); err != nil {
			return fmt.Errorf("failed to parse %s (not room or area format): %w", path, err)
		}

		for i := range area.Rooms {
			if err := m.Rooms.Register(&area.Rooms[i]); err != nil {
				return fmt.Errorf("failed to register room from %s: %w", path, err)
			}
		}

		return nil
	})
}

// loadQuests loads quest definitions from content/default/quests/
func (m *Manager) loadQuests() error {
	questsPath := filepath.Join(m.basePath, "quests")

	return filepath.WalkDir(questsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var quest QuestDef
		if err := yaml.Unmarshal(data, &quest); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		if err := m.Quests.Register(&quest); err != nil {
			return fmt.Errorf("failed to register quest from %s: %w", path, err)
		}

		return nil
	})
}

// ValidationError represents a content validation error
type ValidationError struct {
	Type    string
	ID      string
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s.%s: %s", e.Type, e.ID, e.Field, e.Message)
}

// GetStats returns content loading statistics
func (m *Manager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Stats{
		Skills: m.Skills.Count(),
		NPCs:   m.NPCs.Count(),
		Items:  m.Items.Count(),
		Rooms:  m.Rooms.Count(),
		Quests: m.Quests.Count(),
	}
}

// Stats holds content loading statistics
type Stats struct {
	Skills int
	NPCs   int
	Items  int
	Rooms  int
	Quests int
}

// ReloadSkillFile reloads a single skill file
func (m *Manager) ReloadSkillFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read skill file: %w", err)
	}

	var skill SkillDef
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return fmt.Errorf("failed to parse skill file: %w", err)
	}

	if skill.ID == "" {
		return fmt.Errorf("skill missing ID")
	}

	if err := m.Skills.Register(&skill); err != nil {
		return fmt.Errorf("failed to register skill: %w", err)
	}

	return nil
}

// ReloadNPCFile reloads a single NPC file
func (m *Manager) ReloadNPCFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read NPC file: %w", err)
	}

	var npc NPCTemplate
	if err := yaml.Unmarshal(data, &npc); err != nil {
		return fmt.Errorf("failed to parse NPC file: %w", err)
	}

	if npc.ID == "" {
		return fmt.Errorf("NPC missing ID")
	}

	if err := m.NPCs.Register(&npc); err != nil {
		return fmt.Errorf("failed to register NPC: %w", err)
	}

	return nil
}

// ReloadItemFile reloads a single item file
func (m *Manager) ReloadItemFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read item file: %w", err)
	}

	var item ItemDef
	if err := yaml.Unmarshal(data, &item); err != nil {
		return fmt.Errorf("failed to parse item file: %w", err)
	}

	if item.ID == "" {
		return fmt.Errorf("item missing ID")
	}

	if err := m.Items.Register(&item); err != nil {
		return fmt.Errorf("failed to register item: %w", err)
	}

	return nil
}

// ReloadRoomFile reloads a single room file
func (m *Manager) ReloadRoomFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read room file: %w", err)
	}

	// Try single room format first
	var room RoomDef
	if err := yaml.Unmarshal(data, &room); err == nil && room.ID != "" {
		if err := m.Rooms.Register(&room); err != nil {
			return fmt.Errorf("failed to register room: %w", err)
		}
		return nil
	}

	// Fall back to area format
	var area AreaDef
	if err := yaml.Unmarshal(data, &area); err != nil {
		return fmt.Errorf("failed to parse room file: %w", err)
	}

	for i := range area.Rooms {
		if err := m.Rooms.Register(&area.Rooms[i]); err != nil {
			return fmt.Errorf("failed to register room: %w", err)
		}
	}

	return nil
}

// ReloadQuestFile reloads a single quest file
func (m *Manager) ReloadQuestFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read quest file: %w", err)
	}

	var quest QuestDef
	if err := yaml.Unmarshal(data, &quest); err != nil {
		return fmt.Errorf("failed to parse quest file: %w", err)
	}

	if quest.ID == "" {
		return fmt.Errorf("quest missing ID")
	}

	if err := m.Quests.Register(&quest); err != nil {
		return fmt.Errorf("failed to register quest: %w", err)
	}

	return nil
}

// StartContentWatcher starts watching content files for changes
func (m *Manager) StartContentWatcher() error {
	watcher, err := NewContentWatcher(m)
	if err != nil {
		return err
	}
	return watcher.Start()
}

// ValidateContentChange validates a content change before applying
func (m *Manager) ValidateContentChange(contentType string, data json.RawMessage) (bool, []string) {
	var errors []string

	switch contentType {
	case "skill":
		var skill SkillDef
		if err := yaml.Unmarshal(data, &skill); err != nil {
			return false, []string{fmt.Sprintf("Invalid skill YAML: %v", err)}
		}
		if skill.ID == "" {
			errors = append(errors, "Skill missing ID")
		}
		if skill.Name == "" {
			errors = append(errors, "Skill missing name")
		}

	case "npc":
		var npc NPCTemplate
		if err := yaml.Unmarshal(data, &npc); err != nil {
			return false, []string{fmt.Sprintf("Invalid NPC YAML: %v", err)}
		}
		if npc.ID == "" {
			errors = append(errors, "NPC missing ID")
		}
		// Validate loot table items exist
		for _, loot := range npc.Loot {
			if _, exists := m.Items.Get(loot.ItemID); !exists {
				errors = append(errors, fmt.Sprintf("Loot item '%s' does not exist", loot.ItemID))
			}
		}

	case "item":
		var item ItemDef
		if err := yaml.Unmarshal(data, &item); err != nil {
			return false, []string{fmt.Sprintf("Invalid item YAML: %v", err)}
		}
		if item.ID == "" {
			errors = append(errors, "Item missing ID")
		}

	case "room":
		var room RoomDef
		if err := yaml.Unmarshal(data, &room); err != nil {
			return false, []string{fmt.Sprintf("Invalid room YAML: %v", err)}
		}
		if room.ID == "" {
			errors = append(errors, "Room missing ID")
		}
		// Validate exit targets exist
		for direction, targetID := range room.Exits {
			if _, exists := m.Rooms.Get(targetID); !exists {
				errors = append(errors, fmt.Sprintf("Exit '%s' points to non-existent room '%s'", direction, targetID))
			}
		}

	case "quest":
		var quest QuestDef
		if err := yaml.Unmarshal(data, &quest); err != nil {
			return false, []string{fmt.Sprintf("Invalid quest YAML: %v", err)}
		}
		if quest.ID == "" {
			errors = append(errors, "Quest missing ID")
		}
		// Validate quest targets exist
		for _, step := range quest.Steps {
			if step.Type == "kill" || step.Type == "talk" {
				if _, exists := m.NPCs.Get(step.Target); !exists {
					errors = append(errors, fmt.Sprintf("Step '%s' references non-existent NPC '%s'", step.ID, step.Target))
				}
			}
			if step.Type == "fetch" {
				if _, exists := m.Items.Get(step.ItemID); !exists {
					errors = append(errors, fmt.Sprintf("Step '%s' references non-existent item '%s'", step.ID, step.ItemID))
				}
			}
		}

	default:
		return false, []string{fmt.Sprintf("Unknown content type: %s", contentType)}
	}

	return len(errors) == 0, errors
}

// PreviewContent returns a preview of content data
func (m *Manager) PreviewContent(contentType string, data json.RawMessage) (interface{}, error) {
	switch contentType {
	case "skill":
		var skill SkillDef
		if err := yaml.Unmarshal(data, &skill); err != nil {
			return nil, err
		}
		return skill, nil
	case "npc":
		var npc NPCTemplate
		if err := yaml.Unmarshal(data, &npc); err != nil {
			return nil, err
		}
		return npc, nil
	case "item":
		var item ItemDef
		if err := yaml.Unmarshal(data, &item); err != nil {
			return nil, err
		}
		return item, nil
	case "room":
		var room RoomDef
		if err := yaml.Unmarshal(data, &room); err != nil {
			return nil, err
		}
		return room, nil
	case "quest":
		var quest QuestDef
		if err := yaml.Unmarshal(data, &quest); err != nil {
			return nil, err
		}
		return quest, nil
	default:
		return nil, fmt.Errorf("unknown content type: %s", contentType)
	}
}

// Clear removes all content from the manager
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Skills.Clear()
	m.NPCs.Clear()
	m.Items.Clear()
	m.Rooms.Clear()
	m.Quests.Clear()
}
