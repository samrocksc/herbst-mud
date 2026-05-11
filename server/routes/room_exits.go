package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/service"
)

// cleanupOrphanExits removes exits pointing to non-existent rooms.
func cleanupOrphanExits(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		cleaned, err := svc.Room.CleanupOrphanExits(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"cleaned": cleaned})
	}
}

// createBidirectionalExit creates an exit in both directions between two rooms.
func createBidirectionalExit(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source room id"})
			return
		}
		var input struct {
			Direction    string `json:"direction" binding:"required"`
			TargetRoomID int    `json:"targetRoomId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Room.CreateBidirectionalExit(c.Request.Context(), sourceID, input.Direction, input.TargetRoomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"source": result.Source, "target": result.Target})
	}
}

// deleteBidirectionalExit removes an exit in both directions.
func deleteBidirectionalExit(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source room id"})
			return
		}
		direction := c.Query("direction")
		if direction == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "direction query parameter is required"})
			return
		}
		if err := svc.Room.DeleteBidirectionalExit(c.Request.Context(), sourceID, direction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}