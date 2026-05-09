package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
)

// checkProgress increments objective counts and advances quest progress.
// If all objectives are complete, the quest is marked completed and rewards applied.
func checkProgress(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		questID, err := strconv.Atoi(c.Param("questId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		_, err = client.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		var input questCheckInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		progress, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.HasQuestWith(quest.IDEQ(questID)),
				questprogress.StatusEQ(questprogress.StatusActive),
			).
			WithQuest().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "active quest progress not found"})
			return
		}
		result := advanceObjective(client, c, progress, input, questID)
		if result.err != nil {
			status := http.StatusInternalServerError
			if result.notFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": result.err.Error()})
			return
		}
		c.JSON(http.StatusOK, result.view)
	}
}