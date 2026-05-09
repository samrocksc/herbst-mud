package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/quest"
)

// createQuest creates a new quest definition.
func createQuest(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input questInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name == nil || *input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		repeatMode := quest.RepeatModeNone
		if input.RepeatMode != nil {
			if !validRepeatModes[*input.RepeatMode] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repeat_mode"})
				return
			}
			repeatMode = quest.RepeatMode(*input.RepeatMode)
		}
		mut := client.Quest.Create().
			SetName(*input.Name).
			SetRepeatMode(repeatMode)
		if input.Description != nil {
			mut.SetDescription(*input.Description)
		}
		if input.PrerequisiteQuestIDs != nil {
			mut.SetPrerequisiteQuestIds(*input.PrerequisiteQuestIDs)
		}
		if input.Objectives != nil {
			mut.SetObjectives(questObjectivesToSchema(*input.Objectives))
		}
		if input.Rewards != nil {
			mut.SetRewards(questRewardsToSchema(*input.Rewards))
		}
		if input.CooldownHours != nil {
			mut.SetCooldownHours(*input.CooldownHours)
		}
		if input.IsActive != nil {
			mut.SetIsActive(*input.IsActive)
		}
		q, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, questToView(q))
	}
}