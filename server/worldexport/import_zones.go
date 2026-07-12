package worldexport

import (
	"context"
	"fmt"

	"herbst-server/db"
)

func importZones(ctx context.Context, client *db.Client, zones []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, z := range zones {
		oldID := strVal(z["id"])
		if oldID == "" {
			continue
		}
		created, err := client.Zone.Create().
			SetID(oldID).
			SetName(strVal(z["name"])).
			SetWorldID(newWorldID).
			SetMinLevel(intValOr(z["min_level"], 1)).
			SetColor(strVal(z["color"])).
			SetNillableDescription(strPtr(z["description"])).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("zone %s: %w", oldID, err)
		}
		maps.Zones[oldID] = created.ID
		count++
	}

	// Second pass: remap parent_zone_id to new IDs.
	for _, z := range zones {
		oldID := strVal(z["id"])
		newID := maps.Zones[oldID]
		parentOldID := strVal(z["parent_zone_id"])
		if parentOldID == "" || newID == "" {
			continue
		}
		parentNewID := maps.Zones[parentOldID]
		if parentNewID == "" {
			continue
		}
		_, err := client.Zone.UpdateOneID(newID).SetParentZoneID(parentNewID).Save(ctx)
		if err != nil {
			return count, fmt.Errorf("update parent zone %s: %w", oldID, err)
		}
	}
	return count, nil
}
