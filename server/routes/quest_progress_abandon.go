package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
	"herbst-server/repository"
)

// abandonQuest marks an active quest progress as abandoned.
// TODO: migrate to fully use repos once QuestProgressRepo supports complex queries
func abandonQuest(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "quests"), slog.String("reason", "invalid character id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		questID, err := strconv.Atoi(c.Param("questId"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "quests"), slog.String("reason", "invalid quest id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		_, err = repos.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		// TODO: replace with repo method once QuestProgressRepo supports filtered queries
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
		updated, err := client.QuestProgress.UpdateOne(progress).
			SetStatus(questprogress.StatusAbandoned).
			Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to abandon quest", err, slog.String("service", "quests"), slog.Int("quest_id", questID), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updated, _ = client.QuestProgress.Query().
			Where(questprogress.IDEQ(updated.ID)).
			WithQuest().WithCharacter().
			Only(c.Request.Context())
		slog.Info("quest abandoned", slog.Int("quest_id", questID), slog.Int("character_id", charID), slog.String("user_email", c.GetString("email")), slog.String("service", "quests"))
		c.JSON(http.StatusOK, questProgressToView(updated))
	}
}