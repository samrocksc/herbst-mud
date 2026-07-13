package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"herbst-server/db"
	"herbst-server/dblog"
)

// EventLogSubscriber returns a subscriber that logs every event to the system_logs table.
// This creates an audit trail of all game events. It is registered for ALL event types
// by calling RegisterEventLogSubscriber (which subscribes individually to each known type).
func EventLogSubscriber(client *db.Client, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		ctx := context.Background()

		// Extract character_id from payload if present
		charID := 0
		if id, ok := event.Payload["character_id"].(float64); ok {
			charID = int(id)
		}

		// Serialize payload as JSON details
		detailsBytes, err := json.Marshal(event.Payload)
		if err != nil {
			detailsBytes = []byte(fmt.Sprintf("%v", event.Payload))
		}

		_, err = client.SystemLog.Create().
			SetAction(string(event.Type)).
			SetCharacterID(charID).
			SetDetails(string(detailsBytes)).
			SetTimestamp(time.Now()).
			Save(ctx)
		if err != nil {
			dblog.Error("failed to log event to system_logs", err,
				slog.String("service", "events"),
				slog.String("event_type", string(event.Type)),
			)
			return fmt.Errorf("log event to system_logs: %w", err)
		}

		logger.Debug("event logged to system_logs",
			"event_type", string(event.Type),
			"character_id", charID,
		)

		return nil
	}
}

// RegisterEventLogSubscriber subscribes the event log subscriber to all known event types.
// Call this once from the event registration setup.
func RegisterEventLogSubscriber(client *db.Client, logger *slog.Logger) {
	sub := EventLogSubscriber(client, logger)
	allEventTypes := []EventType{
		EventXPGained,
		EventLevelUp,
		EventSkillLeveledUp,
		EventSkillAbilityUnlocked,
		EventSkillXPGained,
		EventReclass,
		EventRerace,
		EventNPCDefeated,
		EventCharacterDied,
		EventQuestComplete,
		EventSkillLearned,
	}
	for _, et := range allEventTypes {
		Subscribe(et, sub)
	}
}