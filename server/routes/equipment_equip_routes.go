package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/race"
)

// RegisterEquipmentEquipRoutes registers equip/unequip action routes.
func RegisterEquipmentEquipRoutes(router *gin.Engine, client *db.Client) {
	router.PUT("/equipment/:id/equip", equipItem(client))
	router.PUT("/equipment/:id/unequip", unequipItem(client))
}

// equipItem handles PUT /equipment/:id/equip
func equipItem(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		equipItemHandler(c, client)
	}
}

// equipItemHandler contains the core equip logic.
func equipItemHandler(c *gin.Context, client *db.Client) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
		return
	}

	var req struct {
		CharacterID int `json:"character_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := client.Equipment.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	char, err := client.Character.Get(c.Request.Context(), req.CharacterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	slot := item.Slot
	if !constants.IsValidSlot(slot) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Item has invalid slot: " + slot,
		})
		return
	}

	raceObj, err := client.Race.Query().Where(race.NameEQ(char.Race)).Only(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Race not found: " + char.Race,
		})
		return
	}

	if !slotInRace(slot, raceObj.EquipmentSlots) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your race cannot equip items in the " + slot + " slot",
		})
		return
	}

	handleEquipSlotLogic(c, client, id, item, char, slot, raceObj)
}

// slotInRace checks if a slot is available for the given race.
func slotInRace(slot string, raceSlots []string) bool {
	for _, s := range raceSlots {
		if s == slot {
			return true
		}
	}
	return false
}