package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/service"
)

// cleanupOrphanExits removes exits pointing to non-existent rooms.
func cleanupOrphanExits(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		cleaned, err := svc.Room.CleanupOrphanExits(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("Failed to cleanup orphan exits", err, slog.String("service", "rooms"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Cleaned orphan exits", slog.String("service", "rooms"), slog.Int("count", cleaned), slog.String("world_id", worldID))
		c.JSON(http.StatusOK, gin.H{"cleaned": cleaned})
	}
}

// createBidirectionalExit creates an exit in both directions between two rooms.
func createBidirectionalExit(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("Invalid source room id", "error", err, slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source room id"})
			return
		}
		var input struct {
			Direction    string `json:"direction" binding:"required"`
			TargetRoomID int    `json:"targetRoomId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("Invalid bidirectional exit request", "error", err, slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Room.CreateBidirectionalExit(c.Request.Context(), sourceID, input.Direction, input.TargetRoomID)
		if err != nil {
			slog.Warn("Failed to create bidirectional exit", "error", err, slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Created bidirectional exit", slog.String("service", "rooms"), slog.Int("source", sourceID), slog.String("direction", input.Direction), slog.Int("target", input.TargetRoomID))
		c.JSON(http.StatusOK, gin.H{"source": result.Source, "target": result.Target})
	}
}

// deleteBidirectionalExit removes an exit in both directions.
func deleteBidirectionalExit(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("Invalid source room id", "error", err, slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source room id"})
			return
		}
		direction := c.Query("direction")
		if direction == "" {
			slog.Warn("Missing direction query parameter", slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "direction query parameter is required"})
			return
		}
		if err := svc.Room.DeleteBidirectionalExit(c.Request.Context(), sourceID, direction); err != nil {
			slog.Warn("Failed to delete bidirectional exit", "error", err, slog.String("service", "rooms"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Deleted bidirectional exit", slog.String("service", "rooms"), slog.Int("source", sourceID), slog.String("direction", direction))
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
