package events

import (
	"herbst-server/dblog"
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/damagelog"
)

// DefeatXPSubscriber returns a subscriber that awards split XP on npc_defeated events.
// It queries the damage_log to determine proportional damage contribution per attacker,
// then applies a level-gap penalty before awarding XP.
func DefeatXPSubscriber(xpSvc XPAwarder, client *db.Client, logger *slog.Logger) Subscriber {
	return func(event Event) error {
		if event.Type != EventNPCDefeated {
			return nil
		}

		npcIDf, ok := event.Payload["npc_id"].(float64)
		if !ok {
			return fmt.Errorf("missing or invalid npc_id in npc_defeated event")
		}
		npcID := int(npcIDf)

		baseXPf, ok := event.Payload["base_xp"].(float64)
		if !ok {
			baseXPf, ok = event.Payload["xp_value"].(float64)
			if !ok {
				baseXPf, ok = event.Payload["xp"].(float64)
				if !ok {
					return fmt.Errorf("missing or invalid base_xp/xp_value in npc_defeated event")
				}
			}
		}
		baseXP := int(baseXPf)

		// Apply NPC xp_multiplier from the event payload (if present)
		if mult, ok := event.Payload["xp_multiplier"].(float64); ok && mult > 0 {
			baseXP = int(float64(baseXP) * mult)
		}

		npcLevelf, ok := event.Payload["npc_level"].(float64)
		if !ok {
			npcLevelf = 1
		}
		npcLevel := int(npcLevelf)

		// Get NPC template ID for anti-grind tracking
		npcTemplateID, _ := event.Payload["npc_template_id"].(string)

		// Get anti-grind kill threshold from world config (if present)
		antiGrindThreshold := 20 // default
		if worldCfgJSON, ok := event.Payload["world_config"].(map[string]interface{}); ok {
			if sx, ok := worldCfgJSON["skill_xp"].(map[string]interface{}); ok {
				if t, ok := sx["anti_grind_kill_threshold"].(float64); ok && t > 0 {
					antiGrindThreshold = int(t)
				}
			}
		}

		ctx := context.Background()

		logs, err := client.DamageLog.Query().
			Where(damagelog.TargetIDEQ(npcID)).
			Order(damagelog.ByCreatedAt(sql.OrderAsc())).
			All(ctx)
		if err != nil {
			dblog.Error("failed to query damage log", err, slog.Int("npc_id", npcID), slog.String("service", "events"))
			return fmt.Errorf("query damage log: %w", err)
		}

		if len(logs) == 0 {
			logger.Warn("npc defeated but no damage log entries found", "npc_id", npcID)
			return nil
		}

		attackerDamage := make(map[int]int)
		for _, l := range logs {
			attackerDamage[l.AttackerID] += l.Damage
		}

		totalDamage := 0
		for _, dmg := range attackerDamage {
			totalDamage += dmg
		}

		for attackerID, dmg := range attackerDamage {
			share := float64(dmg) / float64(totalDamage) * float64(baseXP)

			attacker, err := client.Character.Get(ctx, attackerID)
			if err != nil {
				logger.Warn("failed to get attacker character, skipping XP award",
					slog.String("error", err.Error()), slog.String("service", "events"), slog.Int("attacker_id", attackerID))
				continue
			}

			levelDiff := attacker.Level - npcLevel
			penaltyPercent := 0.0
			if levelDiff > 3 {
				penaltyPercent = math.Min(0.9, float64(levelDiff)*0.1)
			}

			adjustedXP := int(share * (1.0 - penaltyPercent))

			// Anti-grind: check kill counts for this NPC template
			if npcTemplateID != "" && attacker.KillCounts != nil {
				if kills, ok := attacker.KillCounts[npcTemplateID]; ok && kills >= antiGrindThreshold {
					// Reduce XP to 10% of normal
					adjustedXP = int(float64(adjustedXP) * 0.1)
					logger.Info("anti-grind: reduced XP for over-killed NPC",
						"attacker_id", attackerID,
						"npc_template_id", npcTemplateID,
						"kill_count", kills,
						"threshold", antiGrindThreshold,
						"reduced_xp", adjustedXP,
					)
				}
			}

			if adjustedXP < 1 {
				adjustedXP = 1
			}

			newXP, newLevel, leveledUp, err := xpSvc.AwardXPWithSource(ctx, attackerID, adjustedXP, "kill")
			if err != nil {
				dblog.Error("failed to award defeat XP",
					err, slog.String("service", "events"), slog.Int("attacker_id", attackerID))
				continue
			}

			// Increment kill_counts for anti-grind tracking
			if npcTemplateID != "" {
				killCounts := attacker.KillCounts
				if killCounts == nil {
					killCounts = make(map[string]int)
				}
				killCounts[npcTemplateID]++
				_, _ = client.Character.UpdateOne(attacker).
					SetKillCounts(killCounts).
					Save(ctx)
			}

			if leveledUp {
				logger.Info("character leveled up from defeat",
					"attacker_id", attackerID,
					"new_level", newLevel,
					"total_xp", newXP,
				)
				Publish(Event{
					Type: EventLevelUp,
					Payload: map[string]interface{}{
						"character_id": attackerID,
						"new_level":    newLevel,
						"total_xp":     newXP,
					},
					Timestamp: event.Timestamp,
				})
			}

			logger.Info("defeat xp awarded",
				"attacker_id", attackerID,
				"npc_id", npcID,
				"damage_dealt", dmg,
				"damage_share_pct", fmt.Sprintf("%.1f", float64(dmg)/float64(totalDamage)*100),
				"level_diff", levelDiff,
				"penalty_pct", fmt.Sprintf("%.0f", penaltyPercent*100),
				"base_share", int(share),
				"final_xp", adjustedXP,
			)
		}

		_, err = client.DamageLog.Delete().
			Where(damagelog.TargetIDEQ(npcID)).
			Exec(ctx)
		if err != nil {
			logger.Warn("failed to clean up damage log", "npc_id", npcID, "error", err)
		}

		return nil
	}
}

// RespawnService handles periodic respawn of dead NPCs.
type RespawnService struct {
	client   *db.Client
	logger   *slog.Logger
	interval time.Duration
}

// NewRespawnService creates a respawn ticker.
func NewRespawnService(client *db.Client, logger *slog.Logger) *RespawnService {
	return &RespawnService{
		client:   client,
		logger:   logger,
		interval: 10 * time.Second,
	}
}

// Start begins the respawn ticker loop. Call once at startup.
func (s *RespawnService) Start() {
	go s.tickLoop()
}

func (s *RespawnService) tickLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := s.processRespawns(); err != nil {
			s.logger.Error("respawn tick failed", "error", err)
		}
	}
}

func (s *RespawnService) processRespawns() error {
	ctx := context.Background()

	deadNPCs, err := s.client.Character.Query().
		Where(character.IsNPCEQ(true)).
		Where(character.HitpointsEQ(0)).
		Where(character.DiedAtNotNil()).
		WithNpcTemplate().
		All(ctx)
	if err != nil {
		return fmt.Errorf("query dead npcs: %w", err)
	}

	for _, npc := range deadNPCs {
		template := npc.Edges.NpcTemplate
		if template == nil {
			s.logger.Warn("dead NPC has no template, skipping respawn", "npc_id", npc.ID)
			continue
		}

		cooldown := template.RespawnCooldown
		if cooldown == 0 {
			continue // 0 = no respawn
		}

		elapsed := time.Since(*npc.DiedAt).Seconds()
		if elapsed < float64(cooldown) {
			continue // still on cooldown
		}

		rooms := template.RespawnRooms
		if len(rooms) == 0 {
			continue
		}

		// Pick random room from respawn_rooms
		roomID := 0
		if len(rooms) == 1 {
			roomID, _ = strconv.Atoi(rooms[0])
		} else {
			idx := rand.Intn(len(rooms))
			roomID, _ = strconv.Atoi(rooms[idx])
		}

		if roomID == 0 {
			s.logger.Warn("invalid respawn room for NPC", "npc_id", npc.ID, "rooms", rooms)
			continue
		}

		// Restore NPC (and clear died_at)
		_, err = s.client.Character.UpdateOne(npc).
			SetHitpoints(npc.MaxHitpoints).
			SetCurrentRoomId(roomID).
			ClearDiedAt().
			Save(ctx)
		if err != nil {
			dblog.Error("failed to respawn NPC", err, slog.String("service", "events"), slog.Int("npc_id", npc.ID))
			continue
		}

		s.logger.Info("NPC respawned",
			"npc_id", npc.ID,
			"name", npc.Name,
			"room", roomID,
			"hp", npc.MaxHitpoints,
		)
	}

	return nil
}

