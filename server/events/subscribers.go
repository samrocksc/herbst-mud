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
