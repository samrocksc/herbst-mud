package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
)

// handleEquipSlotLogic handles slot occupation, two-handed logic, and final equip.
func handleEquipSlotLogic(c *gin.Context, client *db.Client, id int, item *db.Equipment, char *db.Character, slot string, raceObj *db.Race) {
	messages := []string{}

	// If two-handed weapon going to main_hand, auto-unequip off_hand
	if item.IsTwoHanded && slot == "main_hand" {
		offHandItems, err := client.Equipment.Query().
			Where(
				equipment.OwnerIdEQ(char.ID),
				equipment.IsEquipped(true),
				equipment.SlotEQ("off_hand"),
			).
			All(c.Request.Context())
		if err == nil {
			for _, offItem := range offHandItems {
				client.Equipment.UpdateOneID(offItem.ID).
					SetIsEquipped(false).
					Save(c.Request.Context())
				messages = append(messages, "Unequipped "+offItem.Name+" from off_hand")
			}
		}
	}

	// If equipping to off_hand when main_hand has two-handed weapon, block
	if slot == "off_hand" {
		mainHandItems, err := client.Equipment.Query().
			Where(
				equipment.OwnerIdEQ(char.ID),
				equipment.IsEquipped(true),
				equipment.SlotEQ("main_hand"),
			).
			All(c.Request.Context())
		if err == nil && len(mainHandItems) > 0 && mainHandItems[0].IsTwoHanded {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Cannot equip off-hand item while wielding a two-handed weapon",
			})
			return
		}
	}

	// Auto-unequip any existing item in the same slot
	existingItems, err := client.Equipment.Query().
		Where(
			equipment.OwnerIdEQ(char.ID),
			equipment.IsEquipped(true),
			equipment.SlotEQ(slot),
		).
		All(c.Request.Context())
	if err == nil {
		for _, existing := range existingItems {
			if existing.ID != id {
				client.Equipment.UpdateOneID(existing.ID).
					SetIsEquipped(false).
					Save(c.Request.Context())
				messages = append(messages, "Unequipped "+existing.Name+" from "+slot)
			}
		}
	}

	// Equip the new item
	_, err = client.Equipment.UpdateOneID(id).
		SetIsEquipped(true).
		SetOwnerId(char.ID).
		Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to equip item"})
		return
	}

	response := gin.H{
		"message":  "Equipped " + item.Name + " in " + slot,
		"item_id":  id,
		"slot":     slot,
		"messages": messages,
	}
	c.JSON(http.StatusOK, response)
}