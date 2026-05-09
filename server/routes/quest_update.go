package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/quest"
)

// updateQuest updates an existing quest definition.
func updateQuest(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		existing, err := client.Quest.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "quest not found"})
			return
		}
		var input questInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mut := client.Quest.UpdateOne(existing)
		if input.Name != nil {
			mut.SetName(*input.Name)
		}
		if input.Description != nil {
			mut.SetDescription(*input.Description)
		}
		if input.RepeatMode != nil {
			if !validRepeatModes[*input.RepeatMode] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repeat_mode"})
				return
			}
			mut.SetRepeatMode(quest.RepeatMode(*input.RepeatMode))
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
		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, questToView(updated))
	}
}