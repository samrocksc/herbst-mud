package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"herbst-server/backup/export"
	"herbst-server/backup/types"
	"herbst-server/db"
)

// CreateBackup creates a backup of all game data
func CreateBackup(client *db.Client, destDir string) (*types.BackupResult, error) {
	ctx := context.Background()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(destDir, fmt.Sprintf("backup_%s", timestamp))

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	manifest := types.Manifest{
		Version:       types.Version,
		CreatedAt:     time.Now().UTC(),
		ServerVersion: "1.0.0",
		Counts:        make(map[string]int),
		Checksums:     make(map[string]string),
	}

	exporters := []struct {
		name string
		fn   func(context.Context, *db.Client, string) (int, error)
	}{
		{"users", export.Users},
		{"rooms", export.Rooms},
		{"skills", export.Skills},
		{"talents", export.Talents},
		{"npc_templates", export.NPCTemplates},
		{"equipment", export.Equipment},
		{"characters", export.Characters},
		{"character_skills", export.CharacterSkills},
		{"character_talents", export.CharacterTalents},
		{"available_talents", export.AvailableTalents},
	}

	for _, exp := range exporters {
		filePath := filepath.Join(backupPath, types.EntityFileNames[exp.name])
		count, err := exp.fn(ctx, client, filePath)
		if err != nil {
			os.RemoveAll(backupPath)
			return nil, fmt.Errorf("failed to export %s: %w", exp.name, err)
		}
		manifest.Counts[exp.name] = count

		checksum, err := calculateChecksum(filePath)
		if err != nil {
			os.RemoveAll(backupPath)
			return nil, fmt.Errorf("failed to calculate checksum for %s: %w", exp.name, err)
		}
		manifest.Checksums[types.EntityFileNames[exp.name]] = checksum
	}

	manifestPath := filepath.Join(backupPath, "manifest.json")
	if err := writeJSON(manifestPath, manifest); err != nil {
		os.RemoveAll(backupPath)
		return nil, fmt.Errorf("failed to write manifest: %w", err)
	}

	return &types.BackupResult{Path: backupPath, Manifest: manifest}, nil
}