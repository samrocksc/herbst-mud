package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
	"herbst-server/db/skill"
)

// InitDefaultSkills seeds the 9 default skills for each world that doesn't have any skills yet.
func InitDefaultSkills(client *db.Client) error {
	ctx := context.Background()

	worlds, err := client.World.Query().All(ctx)
	if err != nil {
		return err
	}

	defaultSkills := []struct {
		Name        string
		DisplayName string
		Category    string
		Description string
	}{
		{"blades", "Blades", "weapon", "Skill with bladed weapons like swords and axes."},
		{"staves", "Staves", "weapon", "Skill with staff weapons."},
		{"knives", "Knives", "weapon", "Skill with knives and daggers."},
		{"martial", "Martial", "weapon", "Martial arts and unarmed combat techniques."},
		{"brawling", "Brawling", "weapon", "Raw street fighting and brawling."},
		{"tech", "Tech", "weapon", "Technical knowledge and device operation."},
		{"light_armor", "Light Armor", "armor", "Proficiency with light armor."},
		{"cloth_armor", "Cloth Armor", "armor", "Proficiency with cloth armor."},
		{"heavy_armor", "Heavy Armor", "armor", "Proficiency with heavy armor."},
	}

	for _, w := range worlds {
		// Check if skills already exist for this world
		existing, err := client.Skill.Query().
			Where(skill.WorldIDEQ(w.ID)).
			Limit(1).
			All(ctx)
		if err != nil {
			log.Printf("Warning: failed to check skills for world %s: %v", w.Name, err)
			continue
		}
		if len(existing) > 0 {
			continue
		}

		// Create default skills for this world
		for _, ds := range defaultSkills {
			_, err := client.Skill.Create().
				SetWorldID(w.ID).
				SetName(ds.Name).
				SetDisplayName(ds.DisplayName).
				SetCategory(ds.Category).
				SetDescription(ds.Description).
				SetMaxLevel(100).
				SetXpCurveMode("percentage").
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to create skill %s for world %s: %v", ds.Name, w.Name, err)
				continue
			}
		}
		log.Printf("Seeded 9 default skills for world %s", w.Name)
	}

	return nil
}