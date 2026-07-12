package worldexport

import (
	"context"

	"herbst-server/db"
	"herbst-server/db/zone"
)

func exportZones(ctx context.Context, client *db.Client, worldID string) ([]map[string]interface{}, error) {
	zones, err := client.Zone.Query().
		Where(zone.WorldIDEQ(worldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return zoneSliceToMaps(zones), nil
}
