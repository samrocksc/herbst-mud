package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db/schema"
	"herbst-server/service"
)

// createQuest creates a new quest definition.
func createQuest(svc *service.Container) gin.HandlerFunc {
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
		name := *input.Name
		description := ""
		if input.Description != nil {
			description = *input.Description
		}
		var prereqs []string
		if input.PrerequisiteQuestIDs != nil {
			prereqs = *input.PrerequisiteQuestIDs
		}
		var objectives []schema.QuestObjective
		if input.Objectives != nil {
			objectives = objectivesToSchema(*input.Objectives)
		}
		var rewards schema.QuestRewards
		if input.Rewards != nil {
			rewards = rewardsToSchema(*input.Rewards)
		}
		repeatMode := "none"
		if input.RepeatMode != nil {
			if !validRepeatModes[*input.RepeatMode] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repeat_mode"})
				return
			}
			repeatMode = *input.RepeatMode
		}
		cooldownHours := 0
		if input.CooldownHours != nil {
			cooldownHours = *input.CooldownHours
		}
		isActive := true
		if input.IsActive != nil {
			isActive = *input.IsActive
		}
		q, err := svc.Quest.CreateQuest(c.Request.Context(), service.CreateQuestInput{
			Name:                 name,
			Description:          description,
			PrerequisiteQuestIDs: prereqs,
			Objectives:           objectives,
			Rewards:              rewards,
			RepeatMode:           repeatMode,
			CooldownHours:        cooldownHours,
			IsActive:             isActive,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, questToView(q))
	}
}