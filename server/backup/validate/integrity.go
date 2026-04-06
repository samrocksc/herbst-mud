package validate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/backup/types"
)

func referentialIntegrity(backupDir string) []types.ValidationError {
	var errors []types.ValidationError

	userIDs := make(map[int]bool)
	roomIDs := make(map[int]bool)
	skillIDs := make(map[int]bool)
	talentIDs := make(map[int]bool)
	npcTemplateIDs := make(map[string]bool)
	characterIDs := make(map[int]bool)

	loadEntityIDs(backupDir, "users.json", userIDs)
	loadEntityIDs(backupDir, "rooms.json", roomIDs)
	loadEntityIDs(backupDir, "skills.json", skillIDs)
	loadEntityIDs(backupDir, "talents.json", talentIDs)
	loadStringIDs(backupDir, "npc_templates.json", npcTemplateIDs)

	validateCharacters(backupDir, userIDs, roomIDs, characterIDs, &errors)
	validateEquipment(backupDir, roomIDs, &errors)

	return errors
}

func loadEntityIDs(backupDir, filename string, ids map[int]bool) {
	data, err := os.ReadFile(filepath.Join(backupDir, filename))
	if err != nil {
		return
	}

	var entities []map[string]interface{}
	if err := json.Unmarshal(data, &entities); err != nil {
		return
	}

	for _, e := range entities {
		if id, ok := e["id"].(float64); ok {
			ids[int(id)] = true
		}
	}
}

func loadStringIDs(backupDir, filename string, ids map[string]bool) {
	data, err := os.ReadFile(filepath.Join(backupDir, filename))
	if err != nil {
		return
	}

	var entities []map[string]interface{}
	if err := json.Unmarshal(data, &entities); err != nil {
		return
	}

	for _, e := range entities {
		if id, ok := e["id"].(string); ok {
			ids[id] = true
		}
	}
}