package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/db/questprogress"
	"herbst-server/db/character"
)

// QuestTriggerSubscriber returns a subscriber that checks active quest progress against
// game events. It subscribes to xp.gained, level.up, and skill.leveled_up events.
//
// Quest trigger criteria are stored in the quest's objectives as QuestObjective entries.
// The subscriber matches event-to-quest-trigger by checking if the event implies a
// threshold was reached (e.g., reach_level:25, gain_xp:1000, skill_level:blades:30).
//
// Trigger matching is simple field-based threshold checking — not a full rules engine.
//
// When a trigger matches, quest progress is updated:
//   - For reach_level triggers: if current_level >= trigger_level, mark objective complete
//   - For gain_xp triggers: if total_xp >= trigger_xp, mark objective complete
//   - For skill_level triggers: if skill_level >= trigger_level, mark objective complete
func QuestTriggerSubscriber(client *db.Client, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		ctx := context.Background()

		charID, ok := event.Payload["character_id"].(float64)
		if !ok {
			return nil
		}
		characterID := int(charID)

		// Query active quest progress for this character, eager-loading the quest
		progressList, err := client.QuestProgress.Query().
			Where(questprogress.StatusEQ(questprogress.StatusActive)).
			Where(questprogress.HasCharacterWith(character.ID(characterID))).
			WithQuest().
			All(ctx)
		if err != nil {
			dblog.Error("failed to query active quest progress", err,
				slog.String("service", "events"),
				slog.String("subscriber", "quest_trigger"),
				slog.Int("character_id", characterID),
			)
			return fmt.Errorf("query active quest progress: %w", err)
		}

		if len(progressList) == 0 {
			return nil
		}

		for _, qp := range progressList {
			quest := qp.Edges.Quest
			if quest == nil {
				continue
			}

			// Check each objective in the quest
			updated := false
			objCounts := qp.ObjectiveCounts
			if objCounts == nil {
				objCounts = make(map[string]int)
			}

			for i, obj := range quest.Objectives {
				key := fmt.Sprintf("%s:%s", obj.Type, obj.TargetID)

				// Skip if already at required count
				if obj.Count > 0 && objCounts[key] >= obj.Count {
					continue
				}

				matched := matchEventToObjective(event, obj)
				if !matched {
					continue
				}

				// Increment the objective count
				objCounts[key]++
				updated = true

				logger.Info("quest trigger matched",
					"character_id", characterID,
					"quest_id", quest.ID,
					"quest_name", quest.Name,
					"objective_index", i,
					"objective_type", obj.Type,
					"objective_target", obj.TargetID,
					"new_count", objCounts[key],
					"required_count", obj.Count,
				)

				// Check if this objective is now complete
				if obj.Count > 0 && objCounts[key] >= obj.Count {
					// Advance current step if this is the current objective
					if i == qp.CurrentStep {
						newStep := qp.CurrentStep + 1
						if newStep >= len(quest.Objectives) {
							// All objectives complete — mark quest as completed
							now := time.Now()
							_, err := client.QuestProgress.UpdateOne(qp).
								SetStatus(questprogress.StatusCompleted).
								SetCompletedAt(now).
								SetCurrentStep(newStep).
								SetObjectiveCounts(objCounts).
								Save(ctx)
							if err != nil {
								dblog.Error("failed to complete quest", err,
									slog.String("service", "events"),
									slog.Int("quest_progress_id", qp.ID),
								)
							} else {
								logger.Info("quest completed via trigger",
									"character_id", characterID,
									"quest_id", quest.ID,
									"quest_name", quest.Name,
								)
								// Emit quest_complete event
								Publish(Event{
									Type: EventQuestComplete,
									Payload: map[string]interface{}{
										"character_id": characterID,
										"quest_id":     quest.ID,
										"quest_name":   quest.Name,
									},
									Timestamp: event.Timestamp,
								})
							}
						} else {
							// Advance to next objective
							_, err := client.QuestProgress.UpdateOne(qp).
								SetCurrentStep(newStep).
								SetObjectiveCounts(objCounts).
								Save(ctx)
							if err != nil {
								dblog.Error("failed to advance quest step", err,
									slog.String("service", "events"),
									slog.Int("quest_progress_id", qp.ID),
								)
							}
						}
					}
				}
			}

			// If only counts changed (no step advancement), save the updated counts
			if updated && qp.Status == questprogress.StatusActive {
				// Re-check: only save if we didn't already save above
				// (the step-advancement saves already set objective_counts)
				_, err := client.QuestProgress.UpdateOne(qp).
					SetObjectiveCounts(objCounts).
					Save(ctx)
				if err != nil {
					dblog.Error("failed to save quest objective counts", err,
						slog.String("service", "events"),
						slog.Int("quest_progress_id", qp.ID),
					)
				}
			}
		}

		return nil
	}
}

// matchEventToObjective checks if an event matches a quest objective's trigger criteria.
// This is simple threshold-based matching:
//   - "reach_level" objectives match level.up events where new_level >= obj.Count
//   - "gain_xp" objectives match xp.gained events where total_xp >= obj.Count
//   - "skill_level" objectives match skill.leveled_up events where skill matches and new_level >= obj.Count
//   - "kill" objectives match npc.defeated events where target matches
func matchEventToObjective(event Event, obj interface{}) bool {
	// QuestObjective fields: Type, TargetID, TagFilter, Count, Labels, Hint
	// We use the struct from the schema package via reflection on the map
	// Since we have the struct imported via db/schema, we can type-assert
	objMap, ok := obj.(map[string]interface{})
	if !ok {
		// Try the actual schema.QuestObjective type — it's a struct, not map
		// We need to handle this via JSON marshal/unmarshal for flexibility
		objBytes, err := json.Marshal(obj)
		if err != nil {
			return false
		}
		if err := json.Unmarshal(objBytes, &objMap); err != nil {
			return false
		}
	}

	objType, _ := objMap["type"].(string)
	targetID, _ := objMap["target_id"].(string)
	countF, _ := objMap["count"].(float64)
	requiredCount := int(countF)
	if requiredCount == 0 {
		requiredCount = 1
	}

	switch event.Type {
	case EventLevelUp:
		if objType != "reach_level" && objType != "level" {
			return false
		}
		newLevelF, ok := event.Payload["new_level"].(float64)
		if !ok {
			return false
		}
		return int(newLevelF) >= requiredCount

	case EventXPGained:
		if objType != "gain_xp" && objType != "xp" {
			return false
		}
		totalXPF, ok := event.Payload["total_xp"].(float64)
		if !ok {
			return false
		}
		return int(totalXPF) >= requiredCount

	case EventSkillLeveledUp:
		if objType != "skill_level" {
			return false
		}
		// Check skill name matches target_id if specified
		if targetID != "" {
			skillName, _ := event.Payload["skill_name"].(string)
			skillID, _ := event.Payload["skill_id"].(string)
			if !strings.EqualFold(skillName, targetID) && !strings.EqualFold(skillID, targetID) {
				return false
			}
		}
		newLevelF, ok := event.Payload["new_level"].(float64)
		if !ok {
			return false
		}
		return int(newLevelF) >= requiredCount

	case EventNPCDefeated:
		if objType != "kill" {
			return false
		}
		// Check target matches if target_id is specified
		if targetID != "" {
			npcTemplateID, _ := event.Payload["npc_template_id"].(string)
			if npcTemplateID != targetID {
				return false
			}
		}
		return true

	default:
		return false
	}
}

// RegisterQuestTriggerSubscriber subscribes the quest trigger subscriber to all relevant event types.
func RegisterQuestTriggerSubscriber(client *db.Client, logger *slog.Logger) {
	sub := QuestTriggerSubscriber(client, logger)
	Subscribe(EventXPGained, sub)
	Subscribe(EventLevelUp, sub)
	Subscribe(EventSkillLeveledUp, sub)
	Subscribe(EventNPCDefeated, sub)
}