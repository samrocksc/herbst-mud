package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
)

// SeedDefaultResistances updates existing races with default resistance and vulnerability values.
// This is idempotent — it only sets resistances on races that have empty resistances/vulnerabilities.
func SeedDefaultResistances(client *db.Client) error {
	ctx := context.Background()

	// Default resistance/vulnerability values by race name
	defaults := map[string]struct {
		resistances     map[string]int
		vulnerabilities map[string]int
	}{
		"human": {
			resistances:     map[string]int{},
			vulnerabilities: map[string]int{},
		},
		"turtle": {
			resistances:     map[string]int{"slashing": 15, "piercing": 10, "bludgeoning": 5},
			vulnerabilities: map[string]int{},
		},
		"mutant": {
			resistances:     map[string]int{"poison": 20, "acid": 10},
			vulnerabilities: map[string]int{"fire": -10},
		},
		"faerie": {
			resistances:     map[string]int{"fire": 10, "cold": 10, "magic": 15},
			vulnerabilities: map[string]int{"iron": -10},
		},
		"ooze": {
			resistances:     map[string]int{"acid": 50, "poison": 30},
			vulnerabilities: map[string]int{"fire": -25, "cold": -15},
		},
	}

	// Get all races across all worlds
	allRaces, err := client.Race.Query().All(ctx)
	if err != nil {
		return err
	}

	updated := 0
	for _, r := range allRaces {
		// Skip if already has resistances set
		if len(r.Resistances) > 0 || len(r.Vulnerabilities) > 0 {
			continue
		}

		def, ok := defaults[r.Name]
		if !ok {
			continue
		}

		builder := client.Race.UpdateOneID(r.ID)
		if def.resistances != nil {
			builder = builder.SetResistances(def.resistances)
		}
		if def.vulnerabilities != nil {
			builder = builder.SetVulnerabilities(def.vulnerabilities)
		}

		_, err := builder.Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed resistances for race %s (world %s): %v", r.Name, r.WorldID, err)
		} else {
			updated++
		}
	}

	if updated > 0 {
		log.Printf("Default resistances seeded for %d races", updated)
	}

	return nil
}