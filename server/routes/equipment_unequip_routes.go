package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// unequipItem handles PUT /equipment/:id/unequip
func unequipItem(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Verify ownership
		if item.OwnerId == nil || *item.OwnerId != req.CharacterID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You do not own this item",
			})
			return
		}

		if !item.IsEquipped {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": item.Name + " is not equipped",
			})
			return
		}

		_, err = client.Equipment.UpdateOneID(id).
			SetIsEquipped(false).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unequip item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Unequipped " + item.Name + " from " + item.Slot,
			"item_id": id,
			"slot":    item.Slot,
		})
	}
}