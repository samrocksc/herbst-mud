package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/character"
)

// Characters imports characters from backup
func Characters(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "characters.json"))
	if err != nil {
		return err
	}

	var characters []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Password       string `json:"password"`
		IsNPC          bool   `json:"isNPC"`
		CurrentRoomID   int    `json:"currentRoomId"`
		StartingRoomID  int    `json:"startingRoomId"`
		UserID          int    `json:"user_id"`
		IsAdmin         bool   `json:"is_admin"`
		Hitpoints       int    `json:"hitpoints"`
		MaxHitpoints    int    `json:"max_hitpoints"`
		Stamina         int    `json:"stamina"`
		MaxStamina      int    `json:"max_stamina"`
		Mana            int    `json:"mana"`
		MaxMana         int    `json:"max_mana"`
		Race            string `json:"race"`
		Class           string `json:"class"`
		Specialty       string `json:"specialty"`
		Level           int    `json:"level"`
		Constitution    int    `json:"constitution"`
		Gender          string `json:"gender"`
		Description     string `json:"description"`
		Strength        int    `json:"strength"`
		Dexterity       int    `json:"dexterity"`
		Intelligence    int    `json:"intelligence"`
		Wisdom          int    `json:"wisdom"`
		NPCTemplateID   string `json:"npc_template_id"`
	}
	if err := json.Unmarshal(data, &characters); err != nil {
		return err
	}

	for _, c := range characters {
		existing, err := client.Character.Query().Where(character.Name(c.Name)).Only(ctx)
		if err == nil {
			mapping.Characters[c.ID] = existing.ID
			continue
		}

		builder := client.Character.Create().
			SetName(c.Name).SetPassword(c.Password).SetIsNPC(c.IsNPC).
			SetCurrentRoomId(mapping.Rooms[c.CurrentRoomID]).
			SetStartingRoomId(mapping.Rooms[c.StartingRoomID]).
			SetIsAdmin(c.IsAdmin).SetHitpoints(c.Hitpoints).
			SetMaxHitpoints(c.MaxHitpoints).SetStamina(c.Stamina).
			SetMaxStamina(c.MaxStamina).SetMana(c.Mana).SetMaxMana(c.MaxMana).
			SetRace(c.Race).SetClass(c.Class).SetLevel(c.Level).
			SetConstitution(c.Constitution).SetStrength(c.Strength).
			SetDexterity(c.Dexterity).SetIntelligence(c.Intelligence).
			SetWisdom(c.Wisdom)

		if c.Specialty != "" {
			builder = builder.SetNillableSpecialty(&c.Specialty)
		}
		if c.Gender != "" {
			builder = builder.SetNillableGender(&c.Gender)
		}
		if c.Description != "" {
			builder = builder.SetNillableDescription(&c.Description)
		}
		if c.UserID != 0 {
			if newUserID, ok := mapping.Users[c.UserID]; ok {
				builder = builder.SetUserID(newUserID)
			}
		}
		if c.NPCTemplateID != "" {
			builder = builder.SetNillableNpcTemplateID(&c.NPCTemplateID)
		}

		created, err := builder.Save(ctx)
		if err != nil {
			return err
		}
		mapping.Characters[c.ID] = created.ID
	}
	return nil
}