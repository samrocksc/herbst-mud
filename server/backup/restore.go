package backup

import (
	"context"
	"fmt"

	"herbst-server/backup/restore"
	"herbst-server/backup/types"
	"herbst-server/backup/validate"
	"herbst-server/db"
)

// RestoreBackup restores data from a backup directory
func RestoreBackup(client *db.Client, backupDir string) error {
	result, err := validate.Backup(backupDir)
	if err != nil {
		return fmt.Errorf("backup validation failed: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("backup validation failed with critical errors")
	}

	mapping := types.NewIDMapping()
	ctx := context.Background()

	if err := restore.Users(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import users: %w", err)
	}
	if err := restore.Rooms(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import rooms: %w", err)
	}
	if err := restore.Skills(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import skills: %w", err)
	}
	if err := restore.Talents(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import talents: %w", err)
	}
	if err := restore.NPCTemplates(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import NPC templates: %w", err)
	}
	if err := restore.Equipment(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import equipment: %w", err)
	}
	if err := restore.Characters(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import characters: %w", err)
	}
	if err := restore.CharacterSkills(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import character skills: %w", err)
	}
	if err := restore.CharacterTalents(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import character talents: %w", err)
	}
	if err := restore.AvailableTalents(ctx, client, backupDir, mapping); err != nil {
		return fmt.Errorf("failed to import available talents: %w", err)
	}

	return nil
}