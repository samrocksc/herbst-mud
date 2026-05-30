package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// unequipItem handles PUT /equipment/:id/unequip
func unequipItem(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "invalid equipment id"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		var req struct {
			CharacterID int `json:"character_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := repos.Equipment.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		// Verify ownership
		if item.OwnerId == nil || *item.OwnerId != req.CharacterID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You do not own this item",
			})
			return
		}

		if !item.IsEquipped {
			slog.Warn("bad request", slog.String("service", "equipment"), slog.String("reason", "item not equipped"))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": item.Name + " is not equipped",
			})
			return
		}

		falseVal := false
		_, err = repos.Equipment.Update(c.Request.Context(), id, repository.EquipmentUpdates{IsEquipped: &falseVal})
		if err != nil {
			dblog.Error("failed to unequip item", err, slog.String("service", "equipment"), slog.Int("item_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unequip item"})
			return
		}

		slog.Info("item unequipped", slog.String("service", "equipment"), slog.Int("item_id", id), slog.String("slot", item.Slot))
		c.JSON(http.StatusOK, gin.H{
			"message": "Unequipped " + item.Name + " from " + item.Slot,
			"item_id": id,
			"slot":    item.Slot,
		})
	}
}
