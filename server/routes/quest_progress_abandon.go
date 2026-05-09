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

// abandonQuest marks an active quest progress as abandoned.
func abandonQuest(client *db.Client) gin.HandlerFunc {
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
		// Verify character exists
		_, err = client.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		// Find active progress record
		progress, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.HasQuestWith(quest.IDEQ(questID)),
				questprogress.StatusEQ(questprogress.StatusActive),
			).
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "active quest progress not found"})
			return
		}
		// Mark as abandoned
		updated, err := client.QuestProgress.UpdateOne(progress).
			SetStatus(questprogress.StatusAbandoned).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updated, _ = client.QuestProgress.Query().
			Where(questprogress.IDEQ(updated.ID)).
			WithQuest().
			WithCharacter().
			Only(c.Request.Context())
		c.JSON(http.StatusOK, questProgressToView(updated))
	}
}