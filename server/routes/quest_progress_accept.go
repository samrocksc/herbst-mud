package routes

import (
	"herbst-server/dblog"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/service"
)

// acceptQuest creates a new QuestProgress record for a character.
func acceptQuest(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		var input questAcceptInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.QuestProgress.Accept(c.Request.Context(), charID, input.QuestID)
		if err != nil {
			dblog.Error("accept quest failed", err, slog.String("service", "quests"))
			switch {
			case isNotFound(err):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case isConflict(err):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

// isNotFound checks if an error indicates a resource was not found.
func isNotFound(err error) bool {
	return err != nil && (err.Error() == "character not found" || err.Error() == "quest not found" ||
		err.Error() == "quest not found: not found")
}

// isConflict checks if an error indicates a conflict (e.g., already active).
func isConflict(err error) bool {
	return err != nil && err.Error() == "quest already active for this character"
}