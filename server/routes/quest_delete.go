package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
)

// deleteQuest removes a quest by ID.
// Fails with 409 if the quest has progress records referencing it.
func deleteQuest(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		// Check if any quest progress references this quest
		progressCount, err := client.QuestProgress.Query().
			Where(questprogress.HasQuestWith(quest.IDEQ(id))).
			Count(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if progressCount > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":       "cannot delete quest: referenced by quest progress",
				"progress_count": progressCount,
			})
			return
		}
		err = client.Quest.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "quest not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}