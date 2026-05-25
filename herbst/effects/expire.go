package effects

import (
	"context"
	"fmt"
	"time"
)

// expiredEffect represents an active effect that has expired.
type expiredEffect struct {
	ID          int    `json:"id"`
	CharacterID int    `json:"character_id"`
	EffectID    int    `json:"effect_id"`
	EffectType  string `json:"effect_type"`
}

// CheckExpiredEffects queries the REST API for active effects whose
// expires_at is in the past, deactivates them, and fires on_effect_end events.
func (s *Service) CheckExpiredEffects(ctx context.Context) error {
	// Fetch expired effects for all characters
	var result struct {
		Expired []expiredEffect `json:"expired"`
	}
	if err := s.getJSON(ctx, "/api/effects/expired", &result); err != nil {
		return fmt.Errorf("check expired effects: %w", err)
	}

	for _, exp := range result.Expired {
		// Fire on_effect_end event for this expiration
		s.FireEvent("on_effect_end", exp.CharacterID, "", map[string]interface{}{
			"effect_id":   exp.EffectID,
			"effect_type": exp.EffectType,
		})

		s.logger.Info("expired effect deactivated",
			"effect_id", exp.EffectID,
			"character_id", exp.CharacterID,
			"effect_type", exp.EffectType,
		)
	}

	return nil
}

// StartExpiryLoop starts a background goroutine that checks for expired effects.
func (s *Service) StartExpiryLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.CheckExpiredEffects(context.Background()); err != nil {
				s.logger.Error("expired effects check failed", "error", err)
			}
		}
	}()
}