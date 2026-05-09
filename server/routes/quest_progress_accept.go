package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
)

// acceptQuest creates a new QuestProgress record for a character.
// Validates: quest is active, prerequisites met, not already active, cooldown check.
func acceptQuest(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		_, err = client.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		var input questAcceptInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		q, err := client.Quest.Get(c.Request.Context(), input.QuestID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "quest not found"})
			return
		}
		if !q.IsActive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quest is not active"})
			return
		}
		if err := validateNotAlreadyActive(client, c, charID, input.QuestID); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err := validatePrerequisites(client, c, charID, q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateCooldown(client, c, charID, q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		progress, err := client.QuestProgress.Create().
			SetCharacterID(charID).
			SetQuestID(input.QuestID).
			SetStatus(questprogress.StatusActive).
			SetStartedAt(time.Now()).
			SetCurrentStep(0).
			SetObjectiveCounts(map[string]int{}).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		progress, _ = client.QuestProgress.Query().
			Where(questprogress.IDEQ(progress.ID)).
			WithQuest().WithCharacter().
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, questProgressToView(progress))
	}
}

// validateNotAlreadyActive checks character doesn't already have this quest active.
func validateNotAlreadyActive(client *db.Client, c *gin.Context, charID, questID int) error {
	activeCount, err := client.QuestProgress.Query().
		Where(
			questprogress.HasCharacterWith(character.IDEQ(charID)),
			questprogress.HasQuestWith(quest.IDEQ(questID)),
			questprogress.StatusEQ(questprogress.StatusActive),
		).
		Count(c.Request.Context())
	if err != nil {
		return err
	}
	if activeCount > 0 {
		return &validationError{msg: "quest already active for this character"}
	}
	return nil
}