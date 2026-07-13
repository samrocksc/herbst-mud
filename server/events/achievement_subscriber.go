package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/db/systemlog"
)

// AchievementSubscriber returns a subscriber that checks achievement criteria against
// game events. It subscribes to xp.gained, level.up, skill.leveled_up, reclass, and rerace events.
// For each event, it queries achievements whose criteria match the event type and checks
// if the character has met the threshold. Completed achievements are logged to system_logs
// (action="achievement.completed") to avoid duplicate awards.
//
// Criteria JSON format (stored in Achievement.Criteria field):
//   {"type": "level", "threshold": 10}        — reach level 10
//   {"type": "xp", "threshold": 1000}        — accumulate 1000 total XP
//   {"type": "skill_level", "threshold": 30} — level a skill to 30
//   {"type": "reclass", "threshold": 1}      — reclass N times
//   {"type": "rerace", "threshold": 1}       — rerace N times
func AchievementSubscriber(client *db.Client, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		ctx := context.Background()

		charID, ok := event.Payload["character_id"].(float64)
		if !ok {
			return nil // no character to credit
		}
		characterID := int(charID)

		// Map event type to criteria type
		var criteriaType string
		switch event.Type {
		case EventLevelUp:
			criteriaType = "level"
		case EventXPGained:
			criteriaType = "xp"
		case EventSkillLeveledUp:
			criteriaType = "skill_level"
		case EventReclass:
			criteriaType = "reclass"
		case EventRerace:
			criteriaType = "rerace"
		default:
			return nil
		}

		// Extract the current value from the event payload
		var currentValue int
		switch event.Type {
		case EventLevelUp:
			if v, ok := event.Payload["new_level"].(float64); ok {
				currentValue = int(v)
			}
		case EventXPGained:
			if v, ok := event.Payload["total_xp"].(float64); ok {
				currentValue = int(v)
			} else if v, ok := event.Payload["xp"].(float64); ok {
				currentValue = int(v)
			}
		case EventSkillLeveledUp:
			if v, ok := event.Payload["new_level"].(float64); ok {
				currentValue = int(v)
			} else if v, ok := event.Payload["skill_level"].(float64); ok {
				currentValue = int(v)
			}
		case EventReclass:
			if v, ok := event.Payload["reclass_count"].(float64); ok {
				currentValue = int(v)
			} else {
				currentValue = 1 // at least one reclass happened
			}
		case EventRerace:
			if v, ok := event.Payload["rerace_count"].(float64); ok {
				currentValue = int(v)
			} else {
				currentValue = 1
			}
		}

		// Query all achievements
		achievements, err := client.Achievement.Query().All(ctx)
		if err != nil {
			dblog.Error("failed to query achievements", err,
				slog.String("service", "events"),
				slog.String("subscriber", "achievement"),
			)
			return fmt.Errorf("query achievements: %w", err)
		}

		for _, ach := range achievements {
			if ach.Criteria == "" {
				continue
			}

			var criteria map[string]interface{}
			if err := json.Unmarshal([]byte(ach.Criteria), &criteria); err != nil {
				logger.Warn("achievement has invalid criteria JSON",
					"achievement_id", ach.ID,
					"name", ach.Name,
					"criteria", ach.Criteria,
				)
				continue
			}

			cType, ok := criteria["type"].(string)
			if !ok || cType != criteriaType {
				continue
			}

			thresholdF, ok := criteria["threshold"].(float64)
			if !ok {
				continue
			}
			threshold := int(thresholdF)

			if currentValue < threshold {
				continue
			}

			// Check if this achievement was already completed by this character
			alreadyCompleted, err := isAchievementCompleted(ctx, client, characterID, ach.ID)
			if err != nil {
				logger.Warn("failed to check achievement completion status",
					"achievement_id", ach.ID,
					"character_id", characterID,
					"error", err,
				)
				continue
			}
			if alreadyCompleted {
				continue
			}

			// Mark achievement as completed by logging to system_logs
			details := fmt.Sprintf(`{"achievement_id":%d,"name":%q,"xp_reward":%d,"criteria_type":%q,"threshold":%d,"value":%d}`,
				ach.ID, ach.Name, ach.XpReward, criteriaType, threshold, currentValue)

			_, err = client.SystemLog.Create().
				SetAction("achievement.completed").
				SetCharacterID(characterID).
				SetDetails(details).
				SetTimestamp(time.Now()).
				Save(ctx)
			if err != nil {
				dblog.Error("failed to log achievement completion", err,
					slog.String("service", "events"),
					slog.Int("achievement_id", ach.ID),
					slog.Int("character_id", characterID),
				)
				continue
			}

			logger.Info("achievement completed",
				"achievement_id", ach.ID,
				"name", ach.Name,
				"character_id", characterID,
				"criteria_type", criteriaType,
				"threshold", threshold,
				"current_value", currentValue,
				"xp_reward", ach.XpReward,
			)

			// Emit achievement.completed event
			Publish(Event{
				Type: EventType("achievement.completed"),
				Payload: map[string]interface{}{
					"character_id":  characterID,
					"achievement_id": ach.ID,
					"name":          ach.Name,
					"xp_reward":     ach.XpReward,
				},
				Timestamp: event.Timestamp,
			})
		}

		return nil
	}
}

// isAchievementCompleted checks the system_logs table for a prior achievement.completed
// entry for this character + achievement combination.
func isAchievementCompleted(ctx context.Context, client *db.Client, characterID, achievementID int) (bool, error) {
	// Query system_logs for action="achievement.completed" with this character_id
	logs, err := client.SystemLog.Query().
		Where(systemlog.ActionEQ("achievement.completed")).
		Where(systemlog.CharacterIDEQ(characterID)).
		All(ctx)
	if err != nil {
		return false, err
	}

	target := fmt.Sprintf(`"achievement_id":%d`, achievementID)
	for _, log := range logs {
		if strings.Contains(log.Details, target) {
			return true, nil
		}
	}
	return false, nil
}

// RegisterAchievementSubscriber subscribes the achievement subscriber to all relevant event types.
func RegisterAchievementSubscriber(client *db.Client, logger *slog.Logger) {
	sub := AchievementSubscriber(client, logger)
	Subscribe(EventXPGained, sub)
	Subscribe(EventLevelUp, sub)
	Subscribe(EventSkillLeveledUp, sub)
	Subscribe(EventReclass, sub)
	Subscribe(EventRerace, sub)
}

// strconv import guard — used for potential future int parsing
var _ = strconv.Itoa