package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/backup/types"
	"herbst-server/db"
)

// CharacterAbilities imports character abilities from backup
func CharacterAbilities(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "character_skills.json"))
	if err != nil {
		return err
	}

	var skills []struct {
		ID          int `json:"id"`
		Level       int `json:"level"`
		Experience  int `json:"experience"`
		CharacterID int `json:"character_id"`
		SkillID     int `json:"skill_id"`
	}
	if err := json.Unmarshal(data, &skills); err != nil {
		return err
	}

	for _, s := range skills {
		newCharID := mapping.Characters[s.CharacterID]
		newSkillID := mapping.Skills[s.SkillID]

		_, err := client.CharacterAbility.Create().
			SetSlot(0).
			SetCharacterID(newCharID).SetAbilityID(newSkillID).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}