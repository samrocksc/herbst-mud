package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/repository"
)

// RegisterEquipmentEquipRoutes registers equip/unequip action routes.
func RegisterEquipmentEquipRoutes(router *gin.Engine, repos *repository.Container) {
	router.PUT("/equipment/:id/equip", equipItem(repos))
	router.PUT("/equipment/:id/unequip", unequipItem(repos))
}

// equipItem handles PUT /equipment/:id/equip
func equipItem(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		equipItemHandler(c, repos)
	}
}

// equipItemHandler contains the core equip logic.
func equipItemHandler(c *gin.Context, repos *repository.Container) {
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

	item, err := repos.Equipment.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	char, err := repos.Character.Get(c.Request.Context(), req.CharacterID)
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

	raceObj, err := repos.Race.GetByName(c.Request.Context(), char.Race)
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

	handleEquipSlotLogic(c, repos, id, item, char, slot, raceObj)
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