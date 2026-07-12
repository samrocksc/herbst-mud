package worldexport

import (
	"context"
	"fmt"

	"herbst-server/db"
	"herbst-server/db/room"
)

func importRooms(ctx context.Context, client *db.Client, rooms []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, r := range rooms {
		oldID := intVal(r["id"])
		exits := mapVal(r["exits"])
		if exits == nil {
			exits = map[string]interface{}{}
		}

		created, err := client.Room.Create().
			SetName(strVal(r["name"])).
			SetWorldID(newWorldID).
			SetDescription(strVal(r["description"])).
			SetIsStartingRoom(boolVal(r["isStartingRoom"])).
			SetIsRootRoom(boolVal(r["isRootRoom"])).
			SetAtmosphere(room.Atmosphere(strValOr(r["atmosphere"], "air"))).
			SetExits(map[string]int{}).
			SetNillablePosZ(intPtrVal(r["posZ"])).
			SetNillableVersion(intPtrVal(r["version"])).
			SetTags(strSliceVal(r["tags"], []string{})).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("room %d: %w", oldID, err)
		}
		maps.Rooms[oldID] = created.ID
		count++
	}

	// Second pass: remap exits to new room IDs.
	for _, r := range rooms {
		oldID := intVal(r["id"])
		newID := maps.Rooms[oldID]
		exits := mapVal(r["exits"])
		if exits == nil {
			continue
		}
		newExits := make(map[string]int)
		for dir, target := range exits {
			targetID := intVal(target)
			if remapped, ok := maps.Rooms[targetID]; ok {
				newExits[dir] = remapped
			} else {
				newExits[dir] = targetID
			}
		}
		if len(newExits) > 0 {
			_, err := client.Room.UpdateOneID(newID).SetExits(newExits).Save(ctx)
			if err != nil {
				return count, fmt.Errorf("update exits room %d: %w", oldID, err)
			}
		}
	}
	return count, nil
}
