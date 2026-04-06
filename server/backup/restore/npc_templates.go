package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/npctemplate"
)

// NPCTemplates imports NPC templates from backup
func NPCTemplates(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "npc_templates.json"))
	if err != nil {
		return err
	}

	var templates []struct {
		ID          string         `json:"id"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Race        string         `json:"race"`
		Disposition string         `json:"disposition"`
		Level       int            `json:"level"`
		Skills      map[string]int `json:"skills"`
		TradesWith  []string       `json:"trades_with"`
		Greeting    string         `json:"greeting"`
	}
	if err := json.Unmarshal(data, &templates); err != nil {
		return err
	}

	for _, t := range templates {
		existing, err := client.NPCTemplate.Query().Where(npctemplate.ID(t.ID)).Only(ctx)
		if err == nil {
			mapping.NPCTemplates[t.ID] = existing.ID
			continue
		}

		_, err = client.NPCTemplate.Create().
			SetID(t.ID).SetName(t.Name).SetDescription(t.Description).
			SetRace(t.Race).SetDisposition(npctemplate.Disposition(t.Disposition)).
			SetLevel(t.Level).SetSkills(t.Skills).SetTradesWith(t.TradesWith).
			SetGreeting(t.Greeting).Save(ctx)
		if err != nil {
			return err
		}
		mapping.NPCTemplates[t.ID] = t.ID
	}
	return nil
}