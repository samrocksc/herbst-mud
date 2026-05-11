package routes

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
	"herbst-server/repository"
)

// validationError is a simple error for validation failures.
type validationError struct {
	msg string
}

func (e *validationError) Error() string { return e.msg }

// validatePrerequisites verifies all prerequisite quests are completed.
// TODO: migrate to use repos once QuestProgressRepo supports complex filtered count queries
func validatePrerequisites(client *db.Client, repos *repository.Container, c *gin.Context, charID int, q *db.Quest) error {
	for _, prereqStr := range q.PrerequisiteQuestIds {
		prereqID, err := strconv.Atoi(prereqStr)
		if err != nil {
			continue // skip invalid IDs
		}
		completed, err := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.HasQuestWith(quest.IDEQ(prereqID)),
				questprogress.StatusEQ(questprogress.StatusCompleted),
			).
			Count(c.Request.Context())
		if err != nil {
			return err
		}
		if completed == 0 {
			return &validationError{msg: "prerequisite quest not completed: " + prereqStr}
		}
	}
	return nil
}

// validateCooldown ensures the character is not within the cooldown period.
// TODO: migrate to use repos once QuestProgressRepo supports complex filtered queries
func validateCooldown(client *db.Client, repos *repository.Container, c *gin.Context, charID int, q *db.Quest) error {
	if q.RepeatMode == quest.RepeatModeNone {
		return nil
	}
	pastProgress, err := client.QuestProgress.Query().
		Where(
			questprogress.HasCharacterWith(character.IDEQ(charID)),
			questprogress.HasQuestWith(quest.IDEQ(q.ID)),
		).
		All(c.Request.Context())
	if err != nil {
		return err
	}
	for _, p := range pastProgress {
		if p.CompletedAt != nil && q.RepeatMode == quest.RepeatModeCooldown {
			cooldownEnd := p.CompletedAt.Add(
				time.Duration(q.CooldownHours) * time.Hour,
			)
			if time.Now().Before(cooldownEnd) {
				return &validationError{msg: "quest is on cooldown"}
			}
		}
	}
	return nil
}