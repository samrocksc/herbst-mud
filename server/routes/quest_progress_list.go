package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/questprogress"
)

// listCharacterQuests returns all quest progress records for a character.
func listCharacterQuests(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		// Verify character exists
		_, err = client.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		progresses, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
			).
			WithQuest().
			WithCharacter().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]questProgressView, len(progresses))
		for i, p := range progresses {
			result[i] = questProgressToView(p)
		}
		c.JSON(http.StatusOK, gin.H{"quest_progress": result})
	}
}