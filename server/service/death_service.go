package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"herbst-server/events"
	"herbst-server/repository"
)

// DeathConfig holds per-world death configuration.
// In Phase 1 we read from game_configs; later phases add world-scoped config.
type DeathConfig struct {
	DeathPenaltyPercent int
	RespawnRoomID       int
}

// DeathService handles player death: XP penalty, respawn, and event emission.
type DeathService interface {
	HandlePlayerDeath(ctx context.Context, characterID, attackerID int, attackerName string) (*DeathResult, error)
	GetDeathConfig(ctx context.Context, worldID string) (DeathConfig, error)
}

// DeathResult is returned after processing a player death.
type DeathResult struct {
	CharacterID   int
	AttackerID    int
	AttackerName  string
	RespawnRoomID int
	OldHP         int
	NewHP         int
	OldRoomID     int
	NewRoomID     int
}

type deathService struct {
	charRepo repository.CharacterRepo
	logger   *slog.Logger
}

// NewDeathService creates a DeathService.
func NewDeathService(charRepo repository.CharacterRepo, logger *slog.Logger) DeathService {
	return &deathService{
		charRepo: charRepo,
		logger:   logger,
	}
}

// GetDeathConfig loads death-related config from game_configs.
// Falls back to sensible defaults if keys are missing.
// Phase 2 will wire this to the actual game_configs table.
func (s *deathService) GetDeathConfig(ctx context.Context, worldID string) (DeathConfig, error) {
	return DeathConfig{
		DeathPenaltyPercent: 10, // default; Phase 2 reads from game_configs
		RespawnRoomID:       0,  // 0 = use character's respawn_room_id
	}, nil
}

// HandlePlayerDeath is called when a player character's HP reaches 0.
// It:
//  1. Publishes EventCharacterDied (for death penalty subscriber, event logs, biography)
//  2. Applies XP penalty (capped to prevent de-leveling)
//  3. Respawns the character at their respawn room with restored HP/STA/MANA
//  4. Emits a respawn event for biography tracking
func (s *deathService) HandlePlayerDeath(ctx context.Context, characterID, attackerID int, attackerName string) (*DeathResult, error) {
	char, err := s.charRepo.Get(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}

	oldHP := char.Hitpoints
	oldRoomID := char.CurrentRoomId

	// 1. Publish character_died event — this triggers:
	//    - DeathPenaltySubscriber: applies XP penalty
	//    - EventLogSubscriber: records death in event_logs for biography
	events.Publish(events.Event{
		Type: events.EventCharacterDied,
		Payload: map[string]interface{}{
			"character_id":   float64(characterID),
			"character_name": char.Name,
			"attacker_id":    float64(attackerID),
			"attacker_name":  attackerName,
			"level":          float64(char.Level),
			"world_id":       char.CurrentWorld,
			"room_id":        float64(char.CurrentRoomId),
		},
		Timestamp: time.Now().UnixMilli(),
	})

	// 2. Respawn: move to respawn room, restore vitals.
	// XP penalty is handled by the DeathPenaltySubscriber via the event above.
	respawnRoomID := char.RespawnRoomId
	if respawnRoomID == 0 {
		respawnRoomID = char.StartingRoomId
	}

	fullHP := char.MaxHitpoints
	fullSTA := char.MaxStamina
	fullMANA := char.MaxMana

	_, err = s.charRepo.Update(ctx, characterID, repository.CharacterUpdates{
		Hitpoints:      &fullHP,
		MaxHitpoints:   &char.MaxHitpoints,
		Stamina:        &fullSTA,
		MaxStamina:     &char.MaxStamina,
		Mana:           &fullMANA,
		MaxMana:        &char.MaxMana,
		CurrentRoomID:  &respawnRoomID,
	})
	if err != nil {
		s.logger.Error("failed to respawn character", "character_id", characterID, "error", err)
		// Don't return error — the death event was already published
	}

	// 3. Emit respawn event for biography tracking
	events.Publish(events.Event{
		Type: events.EventCharacterRespawned,
		Payload: map[string]interface{}{
			"character_id":    float64(characterID),
			"character_name":  char.Name,
			"respawn_room_id": float64(respawnRoomID),
			"death_room_id":   float64(oldRoomID),
			"level":           float64(char.Level),
			"world_id":        char.CurrentWorld,
		},
		Timestamp: time.Now().UnixMilli(),
	})

	s.logger.Info("player death handled",
		"character_id", characterID,
		"character_name", char.Name,
		"attacker_id", attackerID,
		"attacker_name", attackerName,
		"respawn_room", respawnRoomID,
	)

	return &DeathResult{
		CharacterID:   characterID,
		AttackerID:    attackerID,
		AttackerName:  attackerName,
		RespawnRoomID: respawnRoomID,
		OldHP:         oldHP,
		NewHP:         fullHP,
		OldRoomID:     oldRoomID,
		NewRoomID:     respawnRoomID,
	}, nil
}