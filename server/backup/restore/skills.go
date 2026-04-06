package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/skill"
)

// Skills imports skills from backup
func Skills(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "skills.json"))
	if err != nil {
		return err
	}

	var skills []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		SkillType    string `json:"skill_type"`
		Cost         int    `json:"cost"`
		Cooldown     int    `json:"cooldown"`
		Requirements string `json:"requirements"`
	}
	if err := json.Unmarshal(data, &skills); err != nil {
		return err
	}

	for _, s := range skills {
		existing, err := client.Skill.Query().Where(skill.Name(s.Name)).Only(ctx)
		if err == nil {
			mapping.Skills[s.ID] = existing.ID
			continue
		}

		created, err := client.Skill.Create().
			SetName(s.Name).SetDescription(s.Description).
			SetSkillType(s.SkillType).SetCost(s.Cost).
			SetCooldown(s.Cooldown).SetNillableRequirements(&s.Requirements).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Skills[s.ID] = created.ID
	}
	return nil
}