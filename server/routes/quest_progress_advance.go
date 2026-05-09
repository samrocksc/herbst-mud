package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/questprogress"
)

// progressResult holds the result of advancing a quest objective.
type progressResult struct {
	view     questProgressView
	err      error
	notFound bool
}

// advanceObjective increments counts and checks for completion.
func advanceObjective(client *db.Client, c *gin.Context, progress *db.QuestProgress, input questCheckInput, questID int) progressResult {
	counts := progress.ObjectiveCounts
	if counts == nil {
		counts = map[string]int{}
	}
	increment := input.Count
	if increment <= 0 {
		increment = 1
	}
	counts[input.ObjectiveKey] += increment
	q := progress.Edges.Quest
	if q == nil {
		var err error
		q, err = client.Quest.Get(c.Request.Context(), questID)
		if err != nil {
			return progressResult{err: err, notFound: true}
		}
	}
	mut := client.QuestProgress.UpdateOne(progress).
		SetObjectiveCounts(counts)
	if q != nil && progress.CurrentStep < len(q.Objectives) {
		obj := q.Objectives[progress.CurrentStep]
		key := obj.Type + ":" + obj.TargetID
		if counts[key] >= obj.Count {
			mut.SetCurrentStep(progress.CurrentStep + 1)
		}
	}
	if allObjectivesComplete(q, progress, counts) {
		now := time.Now()
		mut.SetStatus(questprogress.StatusCompleted).SetCompletedAt(now)
	}
	updated, err := mut.Save(c.Request.Context())
	if err != nil {
		return progressResult{err: err}
	}
	updated, _ = client.QuestProgress.Query().
		Where(questprogress.IDEQ(updated.ID)).
		WithQuest().WithCharacter().
		Only(c.Request.Context())
	view := questProgressToView(updated)
	if q != nil && updated.Status == questprogress.StatusCompleted {
		view.RewardsApplied = applyQuestRewards(q.Rewards)
	}
	return progressResult{view: view}
}

// allObjectivesComplete checks if all quest objectives are met.
func allObjectivesComplete(q *db.Quest, progress *db.QuestProgress, counts map[string]int) bool {
	if q == nil {
		return false
	}
	if progress.CurrentStep+1 < len(q.Objectives) {
		return false
	}
	if progress.CurrentStep < len(q.Objectives) {
		obj := q.Objectives[progress.CurrentStep]
		key := obj.Type + ":" + obj.TargetID
		return counts[key] >= obj.Count
	}
	return true
}