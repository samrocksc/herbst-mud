package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GameExport represents the full game data export (matches API format)
type GameExport struct {
	Version    string      `json:"version"`
	ExportedAt string      `json:"exported_at"`
	Rooms      []RoomData  `json:"rooms"`
	NPCs       []NPCData   `json:"npcs"`
	Skills     []SkillData `json:"skills"`
	Items      []ItemData  `json:"items"`
}

// RoomData represents a room in the export
type RoomData struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsStarting  bool      `json:"is_starting"`
	Exits       []ExitData `json:"exits"`
}

// ExitData represents a room exit
type ExitData struct {
	Direction    string `json:"direction"`
	TargetRoomID int    `json:"target_room_id"`
}

// NPCData represents an NPC character
type NPCData struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CurrentRoomID int    `json:"current_room_id"`
	Race          string `json:"race"`
	Class         string `json:"class"`
	Level         int    `json:"level"`
	Hitpoints     int    `json:"hitpoints"`
	MaxHitpoints  int    `json:"max_hitpoints"`
	Stamina       int    `json:"stamina"`
	MaxStamina    int    `json:"max_stamina"`
	Mana          int    `json:"mana"`
	MaxMana       int    `json:"max_mana"`
	Strength      int    `json:"strength"`
	Dexterity     int    `json:"dexterity"`
	Constitution  int    `json:"constitution"`
	Intelligence  int    `json:"intelligence"`
	Wisdom        int    `json:"wisdom"`
	NPCSkillID    string `json:"npc_skill_id,omitempty"`
	IsImmortal    bool   `json:"is_immortal"`
}

// SkillData represents a skill/spell
type SkillData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinCooldown int    `json:"min_cooldown"`
	EffectType  string `json:"effect_type,omitempty"`
}

// ItemData represents an item in the game world
type ItemData struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	LocationType string `json:"location_type"` // "room" or "npc"
	LocationID   int    `json:"location_id"`
}

// ImportResult reports what was imported
type ImportResult struct {
	Success     bool        `json:"success"`
	Imported    ImportCount `json:"imported"`
	Version     string      `json:"version"`
	ImportedAt  string      `json:"imported_at"`
}

type ImportCount struct {
	Rooms  int `json:"rooms"`
	NPCs   int `json:"npcs"`
	Skills int `json:"skills"`
	Items  int `json:"items"`
}

// WorldListResult lists available worlds
type WorldListResult struct {
	Worlds []WorldInfo `json:"worlds"`
	Default string     `json:"default"`
	Count  int        `json:"count"`
}

type WorldInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var (
	contentDir     string
	apiBaseURL     string
	worldFilter    string
	dryRun         bool
	verbose        bool
)

func main() {
	flag.StringVar(&contentDir, "content", "/home/sam/GitHub/source_worlds", "Content directory with JSON files")
	flag.StringVar(&apiBaseURL, "api", "http://localhost:8080", "API base URL")
	flag.StringVar(&worldFilter, "worlds", "", "Comma-separated world names to process")
	flag.BoolVar(&dryRun, "dry-run", false, "Don't actually import, just validate")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed output")
	flag.Parse()

	fmt.Println("=== JSON Import Tool for Herbst MUD ===")
	fmt.Printf("Content dir: %s\n", contentDir)
	fmt.Printf("API URL: %s\n", apiBaseURL)
	fmt.Printf("Dry run: %v\n", dryRun)
	fmt.Println()

	// Find JSON files
	jsonFiles, err := findJSONFiles(contentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding JSON files: %v\n", err)
		os.Exit(1)
	}

	if len(jsonFiles) == 0 {
		fmt.Println("No JSON files found in content directory")
		return
	}

	fmt.Printf("Found %d JSON files\n\n", len(jsonFiles))

	// Parse worlds filter
	var filterWorlds []string
	if worldFilter != "" {
		filterWorlds = strings.Split(worldFilter, ",")
		for i := range filterWorlds {
			filterWorlds[i] = strings.TrimSpace(filterWorlds[i])
		}
	}

	// Process each world
	var results []WorldResult
	for _, file := range jsonFiles {
		base := filepath.Base(file)
		// Extract world name from "default_world.json" or "default_export.json"
		name := strings.TrimSuffix(base, filepath.Ext(base))
		name = strings.TrimSuffix(name, "_world")
		name = strings.TrimSuffix(name, "_export")

		if len(filterWorlds) > 0 {
			found := false
			for _, f := range filterWorlds {
				if f == name {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Skipping %s (not in filter)\n", name)
				continue
			}
		}

		result, err := processWorld(file, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", name, err)
			results = append(results, WorldResult{
				World:   name,
				Success: false,
				Error:   err.Error(),
			})
			continue
		}
		results = append(results, *result)
	}

	// Print summary
	fmt.Println("\n=== Summary ===")
	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
			fmt.Printf("✓ %s: %d rooms, %d NPCs, %d skills, %d items imported\n",
				r.World, r.Imported.Rooms, r.Imported.NPCs, r.Imported.Skills, r.Imported.Items)
		} else {
			fmt.Printf("✗ %s: FAILED - %s\n", r.World, r.Error)
		}
	}
	fmt.Printf("\nSuccess: %d/%d\n", successCount, len(results))
}

// WorldResult holds result for a single world
type WorldResult struct {
	World      string      `json:"world"`
	Success    bool        `json:"success"`
	Imported   ImportCount `json:"imported"`
	Error      string      `json:"error,omitempty"`
}

func findJSONFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func processWorld(filePath, worldName string) (*WorldResult, error) {
	if verbose {
		fmt.Printf("Processing: %s (%s)\n", worldName, filePath)
	}

	// Read and parse JSON file
	export, err := loadGameExport(filePath)
	if err != nil {
		return &WorldResult{World: worldName, Success: false, Error: fmt.Sprintf("Failed to parse: %v", err)}, err
	}

	// Validate first
	if verbose {
		fmt.Println("  Validating...")
	}
	if err := validateExport(export); err != nil {
		return &WorldResult{World: worldName, Success: false, Error: fmt.Sprintf("Validation failed: %v", err)}, nil
	}

	if dryRun {
		fmt.Printf("  [DRY RUN] %s: Would import %d rooms, %d NPCs, %d skills, %d items\n",
			worldName, len(export.Rooms), len(export.NPCs), len(export.Skills), len(export.Items))
		return &WorldResult{
			World: worldName,
			Success: true,
			Imported: ImportCount{
				Rooms:  len(export.Rooms),
				NPCs:   len(export.NPCs),
				Skills: len(export.Skills),
				Items:  len(export.Items),
			},
		}, nil
	}

	// Wipe world first
	if verbose {
		fmt.Println("  Wiping world data...")
	}
	if err := wipeWorld(worldName); err != nil {
		return &WorldResult{World: worldName, Success: false, Error: fmt.Sprintf("Wipe failed: %v", err)}, nil
	}

	// Import
	if verbose {
		fmt.Println("  Importing...")
	}
	_, err = importWorld(export, worldName)
	if err != nil {
		return &WorldResult{World: worldName, Success: false, Error: fmt.Sprintf("Import failed: %v", err)}, nil
	}

	// Verify
	if verbose {
		fmt.Println("  Verifying...")
	}
	verified, err := verifyImport(worldName)
	if err != nil {
		return &WorldResult{World: worldName, Success: false, Error: fmt.Sprintf("Verification failed: %v", err)}, nil
	}

	fmt.Printf("  ✓ Imported: %d rooms, %d NPCs, %d skills, %d items\n",
		len(export.Rooms), len(export.NPCs), len(export.Skills), len(export.Items))
	fmt.Printf("  ✓ Verified: %d rooms, %d NPCs, %d skills, %d items\n",
		verified.Rooms, verified.NPCs, verified.Skills, verified.Items)

	return &WorldResult{
		World:   worldName,
		Success: true,
		Imported: ImportCount{
			Rooms:  len(export.Rooms),
			NPCs:   len(export.NPCs),
			Skills: len(export.Skills),
			Items:  len(export.Items),
		},
	}, nil
}

func loadGameExport(filePath string) (*GameExport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var export GameExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, err
	}
	return &export, nil
}

func validateExport(export *GameExport) error {
	if export.Version != "1.0" {
		return fmt.Errorf("unsupported version: %s", export.Version)
	}
	return nil
}

// WipeRequest configures what to wipe
type WipeRequest struct {
	WipeNPCs      bool `json:"wipe_npcs"`
	WipeRooms     bool `json:"wipe_rooms"`
	WipeItems     bool `json:"wipe_items"`
	WipeSkills    bool `json:"wipe_skills"`
	PreserveUsers bool `json:"preserve_users"`
}

func wipeWorld(worldName string) error {
	if dryRun {
		return nil
	}

	url := fmt.Sprintf("%s/admin/wipe", apiBaseURL)
	req := WipeRequest{
		WipeNPCs:      true,
		WipeRooms:     true,
		WipeItems:     true,
		WipeSkills:    true,
		PreserveUsers: true,
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("wipe failed with status %d", resp.StatusCode)
	}

	return nil
}

func importWorld(export *GameExport, worldName string) (ImportCount, error) {
	if dryRun {
		return ImportCount{
			Rooms:  len(export.Rooms),
			NPCs:   len(export.NPCs),
			Skills: len(export.Skills),
			Items:  len(export.Items),
		}, nil
	}

	url := fmt.Sprintf("%s/admin/import?world=%s", apiBaseURL, worldName)
	body, _ := json.Marshal(export)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return ImportCount{}, err
	}
	defer resp.Body.Close()

	var result ImportResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ImportCount{}, err
	}

	if !result.Success {
		return ImportCount{}, fmt.Errorf("import rejected")
	}

	return result.Imported, nil
}

type VerifyResult struct {
	Rooms  int `json:"rooms"`
	NPCs   int `json:"npcs"`
	Skills int `json:"skills"`
	Items  int `json:"items"`
}

func verifyImport(worldName string) (VerifyResult, error) {
	url := fmt.Sprintf("%s/admin/export?world=%s", apiBaseURL, worldName)
	resp, err := http.Get(url)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()

	var export GameExport
	if err := json.NewDecoder(resp.Body).Decode(&export); err != nil {
		return VerifyResult{}, err
	}

	return VerifyResult{
		Rooms:  len(export.Rooms),
		NPCs:   len(export.NPCs),
		Skills: len(export.Skills),
		Items:  len(export.Items),
	}, nil
}

func listWorlds() ([]WorldInfo, error) {
	url := fmt.Sprintf("%s/admin/export/worlds", apiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result WorldListResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Worlds, nil
}
