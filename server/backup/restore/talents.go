package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/talent"
)

// Talents imports talents from backup
func Talents(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "talents.json"))
	if err != nil {
		return err
	}

	var talents []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Requirements string `json:"requirements"`
	}
	if err := json.Unmarshal(data, &talents); err != nil {
		return err
	}

	for _, t := range talents {
		existing, err := client.Talent.Query().Where(talent.Name(t.Name)).Only(ctx)
		if err == nil {
			mapping.Talents[t.ID] = existing.ID
			continue
		}

		created, err := client.Talent.Create().
			SetName(t.Name).SetDescription(t.Description).
			SetNillableRequirements(&t.Requirements).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Talents[t.ID] = created.ID
	}
	return nil
}