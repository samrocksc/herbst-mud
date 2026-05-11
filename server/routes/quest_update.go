package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db/schema"
	"herbst-server/service"
)

// updateQuest updates an existing quest definition.
func updateQuest(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		var input questInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.RepeatMode != nil && !validRepeatModes[*input.RepeatMode] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repeat_mode"})
			return
		}
		updateInput := service.UpdateQuestInput{
			Name:                 input.Name,
			Description:          input.Description,
			PrerequisiteQuestIDs: input.PrerequisiteQuestIDs,
			CooldownHours:       input.CooldownHours,
			IsActive:            input.IsActive,
		}
		if input.RepeatMode != nil {
			rm := *input.RepeatMode
			updateInput.RepeatMode = &rm
		}
		if input.Objectives != nil {
			objs := objectivesToSchema(*input.Objectives)
			updateInput.Objectives = &objs
		}
		if input.Rewards != nil {
			rwds := rewardsToSchema(*input.Rewards)
			updateInput.Rewards = &rwds
		}
		updated, err := svc.Quest.UpdateQuest(c.Request.Context(), id, updateInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, questToView(updated))
	}
}

// objectivesToSchema converts route input types to schema types.
func objectivesToSchema(objs []questObjectiveInput) []schema.QuestObjective {
	result := make([]schema.QuestObjective, len(objs))
	for i, o := range objs {
		result[i] = schema.QuestObjective{
			Type:     o.Type,
			TargetID: o.TargetID,
			Count:    o.Count,
			Label:    o.Label,
			Hint:     o.Hint,
		}
	}
	return result
}

// rewardsToSchema converts route input type to schema type.
func rewardsToSchema(r questRewardsInput) schema.QuestRewards {
	return schema.QuestRewards{
		XP:             r.XP,
		ItemIDs:        r.ItemIDs,
		EffectIDs:      r.EffectIDs,
		TagAdds:        r.TagAdds,
		TagRemoves:     r.TagRemoves,
		AchievementIDs: r.AchievementIDs,
	}
}