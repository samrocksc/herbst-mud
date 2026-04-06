package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
	"herbst-server/db/user"
)

// Users imports users from backup
func Users(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "users.json"))
	if err != nil {
		return err
	}

	var users []struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"is_admin"`
		GodMode  bool   `json:"god_mode"`
	}
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	for _, u := range users {
		existing, err := client.User.Query().Where(user.Email(u.Email)).Only(ctx)
		if err == nil {
			mapping.Users[u.ID] = existing.ID
			continue
		}

		created, err := client.User.Create().
			SetEmail(u.Email).SetPassword(u.Password).
			SetIsAdmin(u.IsAdmin).SetGodMode(u.GodMode).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Users[u.ID] = created.ID
	}
	return nil
}