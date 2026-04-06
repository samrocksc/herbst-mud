package types

import "time"

// Version is the current backup format version
const Version = "1.0"

// Manifest describes the backup metadata
type Manifest struct {
	Version       string            `json:"version"`
	CreatedAt     time.Time         `json:"created_at"`
	ServerVersion string            `json:"server_version"`
	Counts        map[string]int    `json:"counts"`
	Checksums     map[string]string `json:"checksums"`
}

// ValidationError represents a validation problem found in a backup
type ValidationError struct {
	File     string `json:"file"`
	Severity string `json:"severity"` // "critical" or "warning"
	Message  string `json:"message"`
	Field    string `json:"field,omitempty"`
}

// IDMapping tracks old -> new ID mappings during restore
type IDMapping struct {
	Users        map[int]int       `json:"users"`
	Rooms        map[int]int       `json:"rooms"`
	Skills       map[int]int       `json:"skills"`
	Talents      map[int]int       `json:"talents"`
	NPCTemplates  map[string]string `json:"npc_templates"`
	Characters   map[int]int       `json:"characters"`
	Equipment    map[int]int       `json:"equipment"`
}

// NewIDMapping creates a new IDMapping with initialized maps
func NewIDMapping() *IDMapping {
	return &IDMapping{
		Users:        make(map[int]int),
		Rooms:        make(map[int]int),
		Skills:       make(map[int]int),
		Talents:      make(map[int]int),
		NPCTemplates: make(map[string]string),
		Characters:   make(map[int]int),
		Equipment:    make(map[int]int),
	}
}

// BackupResult contains information about a created backup
type BackupResult struct {
	Path     string   `json:"path"`
	Manifest Manifest `json:"manifest"`
}

// ValidationResult contains backup validation results
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors"`
	Warnings []ValidationError `json:"warnings"`
}

// EntityFileNames maps entity names to their JSON file names
var EntityFileNames = map[string]string{
	"users":              "users.json",
	"rooms":              "rooms.json",
	"skills":             "skills.json",
	"talents":            "talents.json",
	"npc_templates":      "npc_templates.json",
	"equipment":          "equipment.json",
	"characters":         "characters.json",
	"character_skills":   "character_skills.json",
	"character_talents":  "character_talents.json",
	"available_talents": "available_talents.json",
}

// ImportOrder defines the order entities should be imported (dependencies first)
var ImportOrder = []string{
	"users",
	"rooms",
	"skills",
	"talents",
	"npc_templates",
	"equipment",
	"characters",
	"character_skills",
	"character_talents",
	"available_talents",
}