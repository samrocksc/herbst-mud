package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/questprogress"
	"herbst-server/repository"
)

// checkAllInput is the JSON request body for bulk quest progress checking.
type checkAllInput struct {
	ObjectiveType string `json:"objective_type" binding:"required"`
	TargetID      string `json:"target_id" binding:"required"`
	Count         int    `json:"count"`
}

// checkAllQuests finds all active quests for a character that match the
// given objective type and target, then increments progress on each.
// TODO: migrate to fully use repos once QuestProgressRepo supports complex queries
func checkAllQuests(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "quests"), slog.String("reason", "invalid character id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		var input checkAllInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "quests"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		messages := advanceQuestObjective(c.Request.Context(), client, repos, charID, input.ObjectiveType, input.TargetID, input.Count)

		progresses, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.StatusEQ(questprogress.StatusActive),
			).
			WithQuest().
			All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list active quest progress", err, slog.String("service", "quests"), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type updatedQuest struct {
			QuestID   int            `json:"quest_id"`
			QuestName string         `json:"quest_name"`
			Status    string         `json:"status"`
			Counts    map[string]int `json:"objective_counts"`
			Messages  []string       `json:"messages,omitempty"`
		}
		var results []updatedQuest
		for _, p := range progresses {
			q := p.Edges.Quest
			if q == nil {
				continue
			}
			results = append(results, updatedQuest{
				QuestID:   q.ID,
				QuestName: q.Name,
				Status:    string(p.Status),
				Counts:    p.ObjectiveCounts,
			})
		}

		slog.Info("quest progress checked", slog.Int("character_id", charID), slog.Int("updated_count", len(results)), slog.String("user_email", c.GetString("email")), slog.String("service", "quests"))
		c.JSON(http.StatusOK, gin.H{"updated": results, "count": len(results), "messages": messages})
	}
}
