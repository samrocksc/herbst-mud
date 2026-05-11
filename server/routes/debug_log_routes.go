package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type debugLogInput struct {
	CharacterID int    `json:"character_id"`
	RoomID     int    `json:"room_id,omitempty"`
	Message     string `json:"message"`
}

// RegisterDebugLogRoutes adds the POST /api/debug-log endpoint.
// This lets the SSH client emit structured debug logs that flow through
// the server's slog pipeline into the applogs table and SSE stream.
func RegisterDebugLogRoutes(protected *gin.RouterGroup) {
	protected.POST("/debug-log", func(c *gin.Context) {
		var input debugLogInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if input.CharacterID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "character_id is required"})
			return
		}
		if input.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
			return
		}

		attrs := []any{
			slog.Int("character_id", input.CharacterID),
			slog.String("service", "debug"),
		}
		if input.RoomID != 0 {
			attrs = append(attrs, slog.Int("room_id", input.RoomID))
		}

		slog.Info(input.Message, attrs...)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}