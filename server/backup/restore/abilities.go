package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/ability"
)

// Abilities imports abilities from backup (formerly "skills")
func Abilities(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "skills.json"))
	if err != nil {
		return err
	}

	var abilities []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		AbilityType  string `json:"ability_type"`
		Cost         int    `json:"cost"`
		Cooldown     int    `json:"cooldown"`
		Requirements string `json:"requirements"`
	}
	if err := json.Unmarshal(data, &abilities); err != nil {
		return err
	}

	for _, s := range abilities {
		existing, err := client.Ability.Query().Where(ability.Name(s.Name)).Only(ctx)
		if err == nil {
			mapping.Skills[s.ID] = existing.ID
			continue
		}

		created, err := client.Ability.Create().
			SetName(s.Name).SetDescription(s.Description).
			SetAbilityType(s.AbilityType).SetCost(s.Cost).
			SetCooldown(s.Cooldown).SetNillableRequirements(&s.Requirements).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Skills[s.ID] = created.ID
	}
	return nil
}