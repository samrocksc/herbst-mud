package routes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/competencycategory"
	"herbst-server/db/competencylevelthreshold"
	"herbst-server/events"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
)

// RegisterEventRoutes wires the event bus HTTP bridge into the router.
func RegisterEventRoutes(router *gin.Engine, client *db.Client, logger *slog.Logger) {
	bus := events.Default()

	// Wire up subscribers — these stay in memory for the lifetime of the process.
	xpSvc := newXPService(client, logger)
	bus.Subscribe(events.EventNPCDefeated, events.XPSubscriber(xpSvc, logger))
	bus.Subscribe(events.EventCharacterDied, events.DeathPenaltySubscriber(xpSvc, logger, 10)) // 10% death penalty
	bus.Subscribe(events.EventQuestComplete, events.QuestXPSubscriber(xpSvc, logger))

	// POST /api/events — the bridge from the game server to the event bus.
	router.POST("/api/events", handleEvent(logger))
	// GET /api/events — health/debug endpoint listing active subscriber counts.
	router.GET("/api/events", handleEventDebug(logger))
}

func newXPService(client *db.Client, logger *slog.Logger) *xpServiceWrapper {
	return &xpServiceWrapper{client: client, logger: logger}
}

// xpServiceWrapper wraps the XP service for use by event subscribers.
// (We import services here to avoid a cycle: routes -> events -> services -> db.)
type xpServiceWrapper struct {
	client *db.Client
	logger *slog.Logger
}

func (w *xpServiceWrapper) AwardXP(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error) {
	// Delegate to the XP service logic directly.
	return w.awardXPImpl(ctx, characterID, xpGained)
}

func (w *xpServiceWrapper) ApplyDeathPenalty(ctx context.Context, characterID, penaltyPercent int) (xpLost, newXP int, err error) {
	char, err := w.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, err
	}
	xpLost = (char.Xp * penaltyPercent) / 100
	newXP = char.Xp - xpLost
	_, err = w.client.Character.UpdateOne(char).SetXp(newXP).Save(ctx)
	return xpLost, newXP, err
}

// AwardCompetencyXP awards XP to a character's competency in a category.
// It applies the category's xp_multiplier, upserts the character_competency record,
// and recomputes the cached level based on thresholds.
func (w *xpServiceWrapper) AwardCompetencyXP(ctx context.Context, characterID int, categoryID string, rawXP int) error {
	cat, err := w.client.CompetencyCategory.Get(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get competency category %s: %w", categoryID, err)
	}

	multiplied := int(float64(rawXP) * cat.XpMultiplier)

	cc, err := w.client.CharacterCompetency.Query().
		Where(charactercompetency.HasCharacterWith(character.ID(characterID))).
		Where(charactercompetency.HasCategoryWith(competencycategory.ID(categoryID))).
		Only(ctx)
	if err != nil {
		// Record doesn't exist — create it
		_, err = w.client.CharacterCompetency.Create().
			SetXp(multiplied).
			SetLevel(1).
			SetCharacterID(characterID).
			SetCategoryID(categoryID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create character competency: %w", err)
		}
		w.logger.Info("competency started",
			"character_id", characterID, "category", categoryID, "xp", multiplied)
		return nil
	}

	// Update XP
	cc.Xp += multiplied

	// Recompute level from thresholds
	thresholds, err := w.client.CompetencyLevelThreshold.Query().
		Where(competencylevelthreshold.HasCategoryWith(competencycategory.ID(categoryID))).
		Order(competencylevelthreshold.ByLevel(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query thresholds: %w", err)
	}

	newLevel := cc.Level
	for _, t := range thresholds {
		if cc.Xp >= t.XpRequired {
			newLevel = t.Level
		}
	}

	cc.Level = newLevel

	_, err = w.client.CharacterCompetency.UpdateOne(cc).SetXp(cc.Xp).SetLevel(cc.Level).Save(ctx)
	if err != nil {
		return fmt.Errorf("update character competency: %w", err)
	}

	w.logger.Info("competency xp awarded",
		"character_id", characterID, "category", categoryID,
		"raw_xp", rawXP, "multiplied", multiplied, "total_xp", cc.Xp, "level", cc.Level)
	return nil
}

// --- event handler ---

type eventRequest struct {
	Type      string                 `json:"type" binding:"required"`
	Payload   map[string]interface{} `json:"payload" binding:"required"`
	Timestamp int64                  `json:"timestamp"`
}

func handleEvent(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req eventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Timestamp == 0 {
			req.Timestamp = time.Now().UnixMilli()
		}

		event := events.Event{
			Type:      events.EventType(req.Type),
			Payload:   req.Payload,
			Timestamp: req.Timestamp,
		}

		if err := event.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logger.Info("event received",
			"type", req.Type,
			"payload_keys", func() []string {
				keys := make([]string, 0, len(req.Payload))
				for k := range req.Payload {
					keys = append(keys, k)
				}
				return keys
			}(),
		)

		events.Default().Publish(event)

		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	}
}

func handleEventDebug(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"registered_events": []string{
				string(events.EventNPCDefeated),
				string(events.EventCharacterDied),
				string(events.EventLevelUp),
				string(events.EventQuestComplete),
				string(events.EventSkillLearned),
			},
		})
	}
}

// --- inline XP award logic (avoids import cycle) ---

func (w *xpServiceWrapper) awardXPImpl(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error) {
	const defaultXPPerLevel = 200

	char, err := w.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, false, err
	}

	oldLevel := char.Level

	// Atomically add XP at the DB level.
	_, err = w.client.Character.UpdateOne(char).
		AddXp(xpGained).
		Save(ctx)
	if err != nil {
		return 0, 0, false, err
	}

	// Re-read to get the new XP total, then compute the correct level.
	char, err = w.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, false, err
	}

	newXP = char.Xp

	// Linear fallback: level N requires N * defaultXPPerLevel total XP.
	newLevel = oldLevel
XPLoop:
	for {
		needed := (newLevel + 1) * defaultXPPerLevel
		if newXP >= needed {
			newLevel++
		} else {
			break XPLoop
		}
	}

	// Update level only if it changed.
	if newLevel != oldLevel {
		_, err = w.client.Character.UpdateOne(char).
			SetLevel(newLevel).
			Save(ctx)
		if err != nil {
			return 0, 0, false, err
		}
		return newXP, newLevel, true, nil
	}

	return newXP, oldLevel, false, nil
}
