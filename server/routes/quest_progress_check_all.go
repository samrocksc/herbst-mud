package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		var input checkAllInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		key := input.ObjectiveType + ":" + input.TargetID
		increment := input.Count
		if increment <= 0 {
			increment = 1
		}
		progresses, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.StatusEQ(questprogress.StatusActive),
			).
			WithQuest().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		type updatedQuest struct {
			QuestID   int              `json:"quest_id"`
			QuestName string           `json:"quest_name"`
			Status    string           `json:"status"`
			Counts    map[string]int   `json:"objective_counts"`
		}
		var results []updatedQuest
		for _, p := range progresses {
			q := p.Edges.Quest
			if q == nil {
				continue
			}
			matched := false
			for _, obj := range q.Objectives {
				if obj.Type == input.ObjectiveType && obj.TargetID == input.TargetID {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			counts := p.ObjectiveCounts
			if counts == nil {
				counts = map[string]int{}
			}
			counts[key] += increment
			mut := client.QuestProgress.UpdateOneID(p.ID).
				SetObjectiveCounts(counts)
			if allObjectivesComplete(q, p, counts) {
				mut = mut.SetStatus(questprogress.StatusCompleted)
			}
			updated, err := mut.Save(c.Request.Context())
			if err != nil {
				continue
			}
			results = append(results, updatedQuest{
				QuestID:   q.ID,
				QuestName: q.Name,
				Status:     string(updated.Status),
				Counts:     counts,
			})
		}
		c.JSON(http.StatusOK, gin.H{"updated": results, "count": len(results)})
	}
}