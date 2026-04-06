package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/backup/types"
	"herbst-server/db"
)

// CharacterSkills imports character skills from backup
func CharacterSkills(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
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

		_, err := client.CharacterSkill.Create().
			SetLevel(s.Level).SetExperience(s.Experience).
			SetCharacterID(newCharID).SetSkillID(newSkillID).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// CharacterTalents imports character talents from backup
func CharacterTalents(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "character_talents.json"))
	if err != nil {
		return err
	}

	var talents []struct {
		ID          int `json:"id"`
		Slot        int `json:"slot"`
		CharacterID int `json:"character_id"`
		TalentID    int `json:"talent_id"`
	}
	if err := json.Unmarshal(data, &talents); err != nil {
		return err
	}

	for _, t := range talents {
		newCharID := mapping.Characters[t.CharacterID]
		newTalentID := mapping.Talents[t.TalentID]

		_, err := client.CharacterTalent.Create().
			SetSlot(t.Slot).SetCharacterID(newCharID).
			SetTalentID(newTalentID).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// AvailableTalents imports available talents from backup
func AvailableTalents(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "available_talents.json"))
	if err != nil {
		return err
	}

	var talents []struct {
		ID              int    `json:"id"`
		UnlockReason    string `json:"unlock_reason"`
		UnlockedAtLevel int    `json:"unlocked_at_level"`
		CharacterID     int    `json:"character_id"`
		TalentID        int    `json:"talent_id"`
	}
	if err := json.Unmarshal(data, &talents); err != nil {
		return err
	}

	for _, t := range talents {
		newCharID := mapping.Characters[t.CharacterID]
		newTalentID := mapping.Talents[t.TalentID]

		_, err := client.AvailableTalent.Create().
			SetUnlockReason(t.UnlockReason).
			SetUnlockedAtLevel(t.UnlockedAtLevel).
			SetCharacterID(newCharID).SetTalentID(newTalentID).Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}