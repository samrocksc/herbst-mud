package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/room"
)

// Rooms imports rooms from backup
func Rooms(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "rooms.json"))
	if err != nil {
		return err
	}

	var rooms []struct {
		ID             int            `json:"id"`
		Name           string         `json:"name"`
		Description    string         `json:"description"`
		IsStartingRoom bool           `json:"isStartingRoom"`
		Exits          map[string]int `json:"exits"`
		Atmosphere     string         `json:"atmosphere"`
	}
	if err := json.Unmarshal(data, &rooms); err != nil {
		return err
	}

	for _, r := range rooms {
		existing, err := client.Room.Query().Where(room.Name(r.Name)).Only(ctx)
		if err == nil {
			mapping.Rooms[r.ID] = existing.ID
			continue
		}

		created, err := client.Room.Create().
			SetName(r.Name).SetDescription(r.Description).
			SetIsStartingRoom(r.IsStartingRoom).SetExits(r.Exits).
			SetAtmosphere(room.Atmosphere(r.Atmosphere)).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Rooms[r.ID] = created.ID
	}
	return nil
}