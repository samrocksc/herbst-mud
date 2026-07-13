package dbinit

import (
	"context"
	"encoding/json"
	"log"

	"herbst-server/db"
)

// defaultWorldConfig is the default config seeded for worlds that don't have one.
var defaultWorldConfig = map[string]interface{}{
	"level_curve": map[string]interface{}{
		"mode":      "percentage",
		"base_xp":   1000,
		"percentage": 50,
		"max_level": 50,
	},
	"stat_growth": map[string]interface{}{
		"hp_per_level":      10,
		"mana_per_level":    5,
		"stamina_per_level": 5,
	},
	"skill_xp": map[string]interface{}{
		"usage_diminishing_returns":  true,
		"usage_cap_per_hour":         100,
		"anti_grind_kill_threshold": 20,
	},
}

// SeedDefaultWorldConfig seeds the default world config for worlds that
// don't have a config set (nil or empty). Safe to call on every startup.
func SeedDefaultWorldConfig(client *db.Client) error {
	ctx := context.Background()

	worlds, err := client.World.Query().All(ctx)
	if err != nil {
		return err
	}

	if len(worlds) == 0 {
		log.Println("No worlds found, skipping world config seed...")
		return nil
	}

	configJSON, err := json.Marshal(defaultWorldConfig)
	if err != nil {
		return err
	}

	for _, w := range worlds {
		if w.Config != nil && len(w.Config) > 0 {
			continue // already has a config
		}

		// Parse a fresh copy of the default config for each world
		var cfg map[string]interface{}
		if err := json.Unmarshal(configJSON, &cfg); err != nil {
			log.Printf("Warning: failed to unmarshal default world config: %v", err)
			continue
		}

		_, err := client.World.UpdateOne(w).
			SetConfig(cfg).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed world config for world %d (%s): %v", w.ID, w.Name, err)
			continue
		}

		log.Printf("Seeded default world config for world %d (%s)", w.ID, w.Name)
	}

	return nil
}