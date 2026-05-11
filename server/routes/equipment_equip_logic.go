package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/repository"
)

// handleEquipSlotLogic handles slot occupation, two-handed logic, and final equip.
func handleEquipSlotLogic(c *gin.Context, repos *repository.Container, id int, item *db.Equipment, char *db.Character, slot string, raceObj *db.Race) {
	messages := []string{}
	falseVal := false

	// If two-handed weapon going to main_hand, auto-unequip off_hand
	if item.IsTwoHanded && slot == "main_hand" {
		allItems, err := repos.Equipment.ListByOwner(c.Request.Context(), char.ID)
		if err == nil {
			for _, offItem := range allItems {
				if offItem.IsEquipped && offItem.Slot == "off_hand" {
					repos.Equipment.Update(c.Request.Context(), offItem.ID, repository.EquipmentUpdates{IsEquipped: &falseVal})
					messages = append(messages, "Unequipped "+offItem.Name+" from off_hand")
				}
			}
		}
	}

	// If equipping to off_hand when main_hand has two-handed weapon, block
	if slot == "off_hand" {
		allItems, err := repos.Equipment.ListByOwner(c.Request.Context(), char.ID)
		if err == nil {
			for _, mainItem := range allItems {
				if mainItem.IsEquipped && mainItem.Slot == "main_hand" && mainItem.IsTwoHanded {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Cannot equip off-hand item while wielding a two-handed weapon",
					})
					return
				}
			}
		}
	}

	// Auto-unequip any existing item in the same slot
	allItems, err := repos.Equipment.ListByOwner(c.Request.Context(), char.ID)
	if err == nil {
		for _, existing := range allItems {
			if existing.IsEquipped && existing.Slot == slot && existing.ID != id {
				repos.Equipment.Update(c.Request.Context(), existing.ID, repository.EquipmentUpdates{IsEquipped: &falseVal})
				messages = append(messages, "Unequipped "+existing.Name+" from "+slot)
			}
		}
	}

	// Equip the new item
	trueVal := true
	ownerID := char.ID
	_, err = repos.Equipment.Update(c.Request.Context(), id, repository.EquipmentUpdates{
		IsEquipped: &trueVal,
		OwnerID:    &ownerID,
	})
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