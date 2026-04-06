package restore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"herbst-server/db"
	"herbst-server/backup/types"
)

// Equipment imports equipment from backup
func Equipment(ctx context.Context, client *db.Client, backupDir string, mapping *types.IDMapping) error {
	data, err := os.ReadFile(filepath.Join(backupDir, "equipment.json"))
	if err != nil {
		return err
	}

	var equipmentList []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		Slot              string `json:"slot"`
		Level             int    `json:"level"`
		Weight            int    `json:"weight"`
		IsEquipped        bool   `json:"isEquipped"`
		IsImmovable       bool   `json:"isImmovable"`
		Color             string `json:"color"`
		IsVisible         bool   `json:"isVisible"`
		ItemType          string `json:"itemType"`
		IsContainer       bool   `json:"isContainer"`
		ContainerCapacity int    `json:"containerCapacity"`
		IsLocked          bool   `json:"isLocked"`
		KeyItemID         string `json:"keyItemID"`
		ContainedItems    string `json:"containedItems"`
		RoomID            int    `json:"room_id"`
	}
	if err := json.Unmarshal(data, &equipmentList); err != nil {
		return err
	}

	for _, e := range equipmentList {
		newRoomID := mapping.Rooms[e.RoomID]

		created, err := client.Equipment.Create().
			SetName(e.Name).SetDescription(e.Description).SetSlot(e.Slot).
			SetLevel(e.Level).SetWeight(e.Weight).SetIsEquipped(e.IsEquipped).
			SetIsImmovable(e.IsImmovable).SetColor(e.Color).SetIsVisible(e.IsVisible).
			SetItemType(e.ItemType).SetIsContainer(e.IsContainer).
			SetContainerCapacity(e.ContainerCapacity).SetIsLocked(e.IsLocked).
			SetNillableKeyItemID(&e.KeyItemID).SetContainedItems(e.ContainedItems).
			SetRoomID(newRoomID).Save(ctx)
		if err != nil {
			return err
		}
		mapping.Equipment[e.ID] = created.ID
	}
	return nil
}