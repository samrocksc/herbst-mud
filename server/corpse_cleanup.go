package main

import (
	"context"
	"log"
	"time"

	"herbst-server/db"
	"herbst-server/db/equipment"
)

// corpseRotConfig holds the rotting configuration.
// A corpse expires after CorpseRotMinutes minutes (nil = never rots).
var CorpseRotMinutes *int

// cleanupExpiredCorpses deletes expired corpses from the database.
// Runs every minute from the background goroutine started by startCorpseCleanup.
func cleanupExpiredCorpses(client *db.Client) error {
	ctx := context.Background()
	now := time.Now()

	count, err := client.Equipment.Delete().
		Where(equipment.ExpiresAtLTE(now)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("[corpse-rot] removed %d expired corpse/placed item(s)", count)
	}
	return nil
}

// startCorpseCleanup launches the background corpse cleanup goroutine.
// Called from main() after the DB client is initialized.
func startCorpseCleanup(client *db.Client) {
	interval := 1 * time.Minute
	if CorpseRotMinutes != nil && *CorpseRotMinutes > 0 {
		// If rotting is configured with a max age, use that as a safety check
		log.Printf("[corpse-rot] configured: corpses rot after %d minute(s)", *CorpseRotMinutes)
	} else {
		log.Printf("[corpse-rot] running: corpses expire via expires_at (configurable via CorpseRotMinutes)")
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := cleanupExpiredCorpses(client); err != nil {
				log.Printf("[corpse-rot] cleanup error: %v", err)
			}
		}
	}()
}
