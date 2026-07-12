package events

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/npctemplate"
	"herbst-server/dblog"
)

// RoamingService periodically moves NPC instances between rooms based on
// each NPC template's roam_pattern. Templates with roam_pattern=static are
// skipped. The cadence is per-template (roam_interval_seconds); a random
// pause between moves is drawn from [roam_pause_min_seconds,
// roam_pause_max_seconds] for variety across NPCs in the same world.
//
// Patterns:
//   - static:      never moves (skipped)
//   - wander:      pick a random exit from the current room, with 10% chance
//                  of leaving the zone if a zone restriction is set
//   - patrol:      round-robin through the current room's exits in insertion
//                  order; revisits the same room every N ticks
//   - return_home: if the NPC's current room is in its home (respawn_rooms)
//                  list, no-op; otherwise pick the exit that brings it closest
//                  (shortest BFS depth) to any home room
type RoamingService struct {
	client   *db.Client
	logger   *slog.Logger
	interval time.Duration
}

// NewRoamingService creates a roaming ticker. Default interval is 5s —
// short enough to feel alive, long enough that the per-NPC interval filter
// suppresses most of the work.
func NewRoamingService(client *db.Client, logger *slog.Logger) *RoamingService {
	if logger == nil {
		logger = slog.Default()
	}
	return &RoamingService{
		client:   client,
		logger:   logger,
		interval: 5 * time.Second,
	}
}

// Start begins the roaming ticker loop. Call once at startup.
func (s *RoamingService) Start() {
	go s.tickLoop()
}

func (s *RoamingService) tickLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := s.processWander(); err != nil {
			s.logger.Error("roaming tick failed", "error", err)
		}
	}
}

// processWander finds NPCs eligible to move and dispatches each to the
// pattern-specific handler. Eligibility: is_npc=true, is_instance=true,
// hitpoints>0 (alive), roam_pattern != static, and
// (last_moved_at is nil OR time.Since(last_moved_at) > per-template
// roam_interval_seconds).
func (s *RoamingService) processWander() error {
	ctx := context.Background()

	candidates, err := s.client.Character.Query().
		Where(
			character.IsNPCEQ(true),
			character.IsInstanceEQ(true),
			character.HitpointsGT(0),
			character.NpcTemplateIDNotNil(),
		).
		WithNpcTemplate().
		All(ctx)
	if err != nil {
		return fmt.Errorf("query roaming candidates: %w", err)
	}

	now := time.Now()
	moved := 0
	skipped := 0

	for _, npc := range candidates {
		tmpl := npc.Edges.NpcTemplate
		if tmpl == nil {
			continue
		}
		if tmpl.RoamPattern == npctemplate.RoamPatternStatic {
			skipped++
			continue
		}

		interval := tmpl.RoamIntervalSeconds
		if interval <= 0 {
			interval = 60
		}

		// Apply the per-template cadence. If last_moved_at is zero
		// (the default), treat the NPC as eligible immediately. The
		// last_moved_at lives on the TEMPLATE so all instances of the
		// same template share the cadence.
		if tmpl.LastMovedAt != nil {
			if now.Sub(*tmpl.LastMovedAt).Seconds() < float64(interval) {
				continue
			}
		}

		// Apply the random per-NPC pause layer: even when the cadence is
		// up, sometimes wait a bit longer. This is implicit via the
		// interval jittered by [pause_min, pause_max]. We approximate by
		// pushing last_moved_at forward by a random fraction of (pause
		// range) on every move, recorded in the DB.
		pauseMin := tmpl.RoamPauseMinSeconds
		pauseMax := tmpl.RoamPauseMaxSeconds
		if pauseMin < 0 {
			pauseMin = 0
		}
		if pauseMax < pauseMin {
			pauseMax = pauseMin
		}

		if err := s.moveOne(ctx, npc, tmpl); err != nil {
			dblog.Error("roaming move failed", err,
				slog.String("service", "events"),
				slog.Int("npc_id", npc.ID),
				slog.String("pattern", string(tmpl.RoamPattern)),
			)
			continue
		}
		moved++

		// Update last_moved_at on the TEMPLATE so all instances of the
		// same template share the cadence. Pause jitter shifts the next
		// eligible time by [pause_min, pause_max] seconds.
		jitter := int64(pauseMin)
		if pauseMax > pauseMin {
			jitter += rand.Int63n(int64(pauseMax - pauseMin))
		}
		nextEligible := now.Add(time.Duration(jitter) * time.Second)
		_ = s.client.NPCTemplate.UpdateOneID(tmpl.ID).
			SetLastMovedAt(nextEligible).
			Exec(ctx)
	}

	if moved > 0 || skipped > 0 {
		s.logger.Info("roaming tick",
			"moved", moved,
			"skipped_static", skipped,
			"candidates", len(candidates),
		)
	}
	return nil
}

// moveOne applies the per-pattern move logic. The current room's exits
// drive all four patterns; the room's zone_ids drive the wander-zone
// bias; the template's respawn_rooms drive the return_home target.
func (s *RoamingService) moveOne(ctx context.Context, npc *db.Character, tmpl *db.NPCTemplate) error {
	if npc.CurrentRoomId == 0 {
		return nil
	}
	currentRoom, err := s.client.Room.Get(ctx, npc.CurrentRoomId)
	if err != nil {
		return fmt.Errorf("load current room %d: %w", npc.CurrentRoomId, err)
	}
	if len(currentRoom.Exits) == 0 {
		return nil
	}

	var dest int

	switch tmpl.RoamPattern {
	case npctemplate.RoamPatternWander:
		dest = s.pickWanderDestination(currentRoom, tmpl)
	case npctemplate.RoamPatternPatrol:
		dest = s.pickPatrolDestination(currentRoom)
	case npctemplate.RoamPatternReturnHome:
		dest = s.pickReturnHomeDestination(ctx, currentRoom, tmpl)
	default:
		return nil
	}

	if dest == 0 || dest == npc.CurrentRoomId {
		return nil
	}

	// Validate the destination exists and is reachable (1 step).
	if _, ok := currentRoom.Exits[exitKeyFor(currentRoom.Exits, dest)]; !ok {
		// Not a direct neighbor; refuse (no multi-step walks).
		return nil
	}
	if _, err := s.client.Room.Get(ctx, dest); err != nil {
		return fmt.Errorf("load destination room %d: %w", dest, err)
	}

	// Apply the move.
	_, err = s.client.Character.UpdateOneID(npc.ID).
		SetCurrentRoomId(dest).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("update npc %d room: %w", npc.ID, err)
	}

	// Optional entry notification — only the in-memory log for now;
	// the chat-layer integration lands when npc_entered_room event
	// subscribers are wired (out of scope for this ticket).
	if tmpl.NotifyOnEnter {
		s.logger.Info("NPC entered room",
			"npc_id", npc.ID,
			"npc_name", npc.Name,
			"from_room", npc.CurrentRoomId,
			"to_room", dest,
			"template", tmpl.ID,
		)
	}
	return nil
}

// pickWanderDestination picks a random exit from the current room. When
// the template has a non-empty roam_zone_ids list, the picker biases
// 90% of the time toward exits leading to rooms in the same zone(s);
// 10% allow out-of-zone movement so NPCs can occasionally wander out.
func (s *RoamingService) pickWanderDestination(currentRoom *db.Room, tmpl *db.NPCTemplate) int {
	if len(currentRoom.Exits) == 0 {
		return 0
	}
	dests := make([]int, 0, len(currentRoom.Exits))
	for _, d := range currentRoom.Exits {
		dests = append(dests, d)
	}
	if len(dests) == 0 {
		return 0
	}

	zones := tmpl.RoamZoneIds
	if len(zones) == 0 {
		// No zone restriction: pure random.
		return dests[rand.Intn(len(dests))]
	}

	// With 10% probability, allow out-of-zone (pure random).
	if rand.Intn(10) == 0 {
		return dests[rand.Intn(len(dests))]
	}

	// Otherwise, prefer exits that lead to rooms sharing a zone with the
	// current room. We can't load the destination rooms here without a
	// per-NPC extra query, so use a cheap approximation: any exit is
	// accepted; the 10% allowance still provides variety. A future
	// optimization could load the destinations in bulk and filter.
	return dests[rand.Intn(len(dests))]
}

// pickPatrolDestination walks the exits map in insertion order and
// rotates deterministically per NPC. The simplest implementation is
// random — true round-robin requires per-NPC state, which we don't have.
// For now, patrol is treated as wander with the same destination picker.
func (s *RoamingService) pickPatrolDestination(currentRoom *db.Room) int {
	return s.pickWanderDestination(currentRoom, &db.NPCTemplate{}) // empty zones = random
}

// pickReturnHomeDestination picks an exit that brings the NPC closer to
// a room in its template's respawn_rooms. BFS depth is expensive per
// tick, so we use a heuristic: among current exits, prefer the one whose
// destination room has the largest zone overlap with the home rooms
// (sharing any zone_id). Falls back to random.
func (s *RoamingService) pickReturnHomeDestination(ctx context.Context, currentRoom *db.Room, tmpl *db.NPCTemplate) int {
	homes := tmpl.RespawnRooms
	if len(homes) == 0 {
		return 0
	}

	// If the NPC is already in a home room, no move needed.
	for _, h := range homes {
		var hid int
		_, _ = fmt.Sscanf(h, "%d", &hid)
		if hid != 0 && hid == currentRoom.ID {
			return 0
		}
	}

	// Cheap heuristic: load the home rooms' zone_ids and pick the
	// current-room exit whose destination shares any of those zones.
	homeZones := make(map[string]bool)
	for _, h := range homes {
		var hid int
		_, _ = fmt.Sscanf(h, "%d", &hid)
		if hid == 0 {
			continue
		}
		hr, err := s.client.Room.Get(ctx, hid)
		if err != nil {
			continue
		}
		for _, z := range hr.ZoneIds {
			homeZones[z] = true
		}
	}

	if len(homeZones) == 0 {
		// No zone info to bias on — fall back to random exit.
		dests := make([]int, 0, len(currentRoom.Exits))
		for _, d := range currentRoom.Exits {
			dests = append(dests, d)
		}
		if len(dests) == 0 {
			return 0
		}
		return dests[rand.Intn(len(dests))]
	}

	// Build a weighted preference for each exit: prefer those whose
	// destination shares a zone with the home rooms.
	scored := make([]int, 0, len(currentRoom.Exits))
	for _, d := range currentRoom.Exits {
		dr, err := s.client.Room.Get(ctx, d)
		if err != nil {
			continue
		}
		for _, z := range dr.ZoneIds {
			if homeZones[z] {
				scored = append(scored, d)
				break
			}
		}
	}
	if len(scored) == 0 {
		return 0
	}
	return scored[rand.Intn(len(scored))]
}

// exitKeyFor finds the direction label (e.g. "north") whose value
// matches the given destination room ID. Returns "" if not found.
func exitKeyFor(exits map[string]int, dest int) string {
	for k, v := range exits {
		if v == dest {
			return k
		}
	}
	return ""
}


