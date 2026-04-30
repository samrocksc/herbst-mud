package events

import (
	"context"
	"fmt"
	"log/slog"
)

// XPSubscriber returns a subscriber that awards XP on npc_defeated events
// and checks for level-up.
func XPSubscriber(xpSvc XPAwarder, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		if event.Type != EventNPCDefeated {
			return nil
		}

		charID, ok := event.Payload["character_id"].(float64)
		if !ok {
			return fmt.Errorf("missing or invalid character_id in npc_defeated event")
		}
		xpGained, ok := event.Payload["xp_value"].(float64)
		if !ok {
			xpGained, ok = event.Payload["xp"].(float64)
			if !ok {
				return fmt.Errorf("missing or invalid xp/xp_value in npc_defeated event")
			}
		}

		ctx := context.Background()
		newXP, newLevel, leveledUp, err := xpSvc.AwardXP(ctx, int(charID), int(xpGained))
		if err != nil {
			return fmt.Errorf("failed to award XP: %w", err)
		}

		if leveledUp {
			logger.Info("character leveled up",
				"character_id", int(charID),
				"new_level", newLevel,
				"total_xp", newXP,
			)
			Publish(Event{
				Type: EventLevelUp,
				Payload: map[string]interface{}{
					"character_id": int(charID),
					"new_level":    newLevel,
					"total_xp":     newXP,
				},
				Timestamp: event.Timestamp,
			})
		}

		return nil
	}
}

// DeathPenaltySubscriber returns a subscriber that penalises XP on character death.
func DeathPenaltySubscriber(xpSvc XPAwarder, logger *slog.Logger, penaltyPercent int) Subscriber {
	return func(event Event) error {
		if event.Type != EventCharacterDied {
			return nil
		}

		charID, ok := event.Payload["character_id"].(float64)
		if !ok {
			return fmt.Errorf("missing or invalid character_id in character_died event")
		}

		ctx := context.Background()
		xpLost, newXP, err := xpSvc.ApplyDeathPenalty(ctx, int(charID), penaltyPercent)
		if err != nil {
			return fmt.Errorf("failed to apply death penalty: %w", err)
		}

		logger.Info("death penalty applied",
			"character_id", int(charID),
			"xp_lost", xpLost,
			"remaining_xp", newXP,
		)
		return nil
	}
}

// QuestXPAwarder extends XPAwarder with competency XP support for quest completion.
type QuestXPAwarder interface {
	XPAwarder
	AwardCompetencyXP(ctx context.Context, characterID int, categoryID string, rawXP int) error
}

// QuestXPSubscriber returns a subscriber that awards XP on quest_complete events.
// It reads quest_xp_enabled from the event payload; if false or missing, it skips.
// It awards base XP via AwardXP and, if a competency_category is set in the payload,
// awards competency XP via AwardCompetencyXP.
func QuestXPSubscriber(xpSvc QuestXPAwarder, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		if event.Type != EventQuestComplete {
			return nil
		}

		// Check if quest XP is enabled for this character
		if enabled, ok := event.Payload["quest_xp_enabled"].(bool); ok && !enabled {
			logger.Debug("quest xp skipped, disabled for character")
			return nil
		}

		charID, ok := event.Payload["character_id"].(float64)
		if !ok {
			return fmt.Errorf("missing or invalid character_id in quest_complete event")
		}

		xpReward, ok := event.Payload["xp_reward"].(float64)
		if !ok {
			return fmt.Errorf("missing or invalid xp_reward in quest_complete event")
		}

		ctx := context.Background()
		charIDInt := int(charID)
		xpRewardInt := int(xpReward)

		// Award base XP for quest completion
		newXP, newLevel, leveledUp, err := xpSvc.AwardXP(ctx, charIDInt, xpRewardInt)
		if err != nil {
			return fmt.Errorf("failed to award quest XP: %w", err)
		}

		logger.Info("quest xp awarded",
			"character_id", charIDInt,
			"xp_reward", xpRewardInt,
			"total_xp", newXP,
		)

		if leveledUp {
			logger.Info("character leveled up from quest",
				"character_id", charIDInt,
				"new_level", newLevel,
				"total_xp", newXP,
			)
			Publish(Event{
				Type: EventLevelUp,
				Payload: map[string]interface{}{
					"character_id": charIDInt,
					"new_level":    newLevel,
					"total_xp":     newXP,
				},
				Timestamp: event.Timestamp,
			})
		}

		// Award competency XP if a competency category is specified
		if category, ok := event.Payload["competency_category"].(string); ok && category != "" {
			if err := xpSvc.AwardCompetencyXP(ctx, charIDInt, category, xpRewardInt); err != nil {
				logger.Error("failed to award quest competency XP",
					"character_id", charIDInt,
					"category", category,
					"error", err,
				)
				// Non-fatal: base XP was already awarded
			}
		}

		return nil
	}
}
