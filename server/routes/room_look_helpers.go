package routes

import (
	"slices"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// partitionCharacters splits characters into NPC and player lists.
func partitionCharacters(c *gin.Context, characters []*db.Character) (npcs, players []interface{}) {
	for _, ch := range characters {
		entry := map[string]interface{}{
			"id":    ch.ID,
			"name":  ch.Name,
			"level": ch.Level,
			"class": ch.Class,
			"race":  ch.Race,
			"hp":    ch.Hitpoints,
			"maxHp": ch.MaxHitpoints,
		}
		if ch.IsNPC {
			if ch.NpcTemplateID != "" {
				entry["npcTemplateId"] = ch.NpcTemplateID
			}
			npcs = append(npcs, entry)
		} else {
			players = append(players, entry)
		}
	}
	return npcs, players
}

// filterVisibleItems returns visible items in a given room from the equipment list.
func filterVisibleItems(equipments []*db.Equipment, roomID int) []interface{} {
	visible := slices.DeleteFunc(equipments, func(e *db.Equipment) bool {
		return e.Edges.Room == nil || e.Edges.Room.ID != roomID || !e.IsVisible
	})
	items := make([]interface{}, len(visible))
	for i, item := range visible {
		items[i] = map[string]interface{}{
			"id":          item.ID,
			"name":        item.Name,
			"slot":        item.Slot,
			"itemType":    item.ItemType,
			"description": item.Description,
		}
	}
	return items
}