package routes

import (
	"context"
	"fmt"
	"time"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/quest"
	"herbst-server/db/questprogress"
	"herbst-server/repository"
)

// advanceQuestObjective increments active quests matching the given objective
// and returns player-facing completion/reward messages.
// For repeatable quests (repeat_mode = "always"), it will also reset completed quests.
func advanceQuestObjective(ctx context.Context, client *db.Client, repos *repository.Container, charID int, objectiveType, targetID string, count int) []string {
	if count <= 0 {
		count = 1
	}
	key := objectiveType + ":" + targetID

	// First, try to find an active quest progress
	progresses, err := client.QuestProgress.Query().
		Where(
			questprogress.HasCharacterWith(character.IDEQ(charID)),
			questprogress.StatusEQ(questprogress.StatusActive),
		).
		WithQuest().
		All(ctx)
	if err != nil {
		return nil
	}

	// If no active quest found, check for completed repeatable quests to reset
	if len(progresses) == 0 {
		completed, resetErr := client.QuestProgress.Query().
			Where(
				questprogress.HasCharacterWith(character.IDEQ(charID)),
				questprogress.StatusEQ(questprogress.StatusCompleted),
			).
			WithQuest().
			All(ctx)
		if resetErr == nil {
			for _, p := range completed {
				q := p.Edges.Quest
				if q == nil {
					continue
				}
				// Check if this quest matches the objective and is repeatable
				for _, obj := range q.Objectives {
					if obj.Type == objectiveType && obj.TargetID == targetID && q.RepeatMode == quest.RepeatModeAlways {
						// Reset completed quest to active
						mut := client.QuestProgress.UpdateOneID(p.ID).
							SetStatus(questprogress.StatusActive).
							SetStartedAt(time.Now()).
							SetObjectiveCounts(map[string]int{})
						if _, err := mut.Save(ctx); err == nil {
							// Refresh from DB for the reset progress
							refreshed, refreshErr := client.QuestProgress.Query().Where(questprogress.IDEQ(p.ID)).WithQuest().Only(ctx)
							if refreshErr == nil && refreshed.Edges.Quest != nil {
								progresses = append(progresses, refreshed)
							}
						}
						break
					}
				}
			}
		}
	}

	messages := []string{}
	for _, p := range progresses {
		q := p.Edges.Quest
		if q == nil {
			continue
		}
		matched := false
		for _, obj := range q.Objectives {
			if obj.Type == objectiveType && obj.TargetID == targetID {
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
		counts[key] += count
		mut := client.QuestProgress.UpdateOneID(p.ID).
			SetObjectiveCounts(counts)
		if allObjectivesComplete(q, p, counts) {
			mut = mut.SetStatus(questprogress.StatusCompleted).SetCompletedAt(time.Now())
		}
		updated, err := mut.Save(ctx)
		if err != nil {
			continue
		}
		if updated.Status == questprogress.StatusCompleted {
			messages = append(messages, fmt.Sprintf("Quest completed: %s!", q.Name))
			rewardSummary := applyQuestRewards(q.Rewards)
			messages = append(messages, fmt.Sprintf("Rewards: XP=%d, Items=%v", rewardSummary["xp"], rewardSummary["item_ids"]))
		}
	}
	return messages
}

// acceptQuestIfNotActive creates an active quest progress record for a character if none exists.
func acceptQuestIfNotActive(ctx context.Context, client *db.Client, repos *repository.Container, charID, questID int) error {
	exists, err := repos.QuestProgress.CountActiveByCharacter(ctx, charID, questID)
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	_, err = client.QuestProgress.Create().
		SetCharacterID(charID).
		SetQuestID(questID).
		SetStatus(questprogress.StatusActive).
		SetStartedAt(time.Now()).
		SetObjectiveCounts(map[string]int{}).
		Save(ctx)
	return err
}

// activeQuestsMatchingObjective returns active quests that match an objective type/target.
func activeQuestsMatchingObjective(ctx context.Context, client *db.Client, charID int, objectiveType, targetID string) ([]*db.Quest, error) {
	progresses, err := client.QuestProgress.Query().
		Where(
			questprogress.HasCharacterWith(character.IDEQ(charID)),
			questprogress.StatusEQ(questprogress.StatusActive),
		).
		WithQuest().
		All(ctx)
	if err != nil {
		return nil, err
	}
	var matched []*db.Quest
	for _, p := range progresses {
		q := p.Edges.Quest
		if q == nil {
			continue
		}
		for _, obj := range q.Objectives {
			if obj.Type == objectiveType && obj.TargetID == targetID {
				matched = append(matched, q)
				break
			}
		}
	}
	return matched, nil
}
