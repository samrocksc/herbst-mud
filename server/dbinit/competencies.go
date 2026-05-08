package dbinit

import (
	"context"
	"log"
	"strconv"

	"herbst-server/db"
)

// InitCompetencies seeds the 9 weapon/armor proficiency categories with level thresholds
func InitCompetencies(client *db.Client) error {
	ctx := context.Background()

	existing, err := client.CompetencyCategory.Query().Count(ctx)
	if err != nil {
		return err
	}
	if existing > 0 {
		log.Println("Competency categories already exist, skipping seed")
		return nil
	}

	categories := []struct {
		ID           string
		Name         string
		XpMultiplier float64
	}{
		{"blades", "Blades", 0.20},
		{"staves", "Staves", 0.20},
		{"knives", "Knives", 0.25},
		{"martial", "Martial Arts", 0.20},
		{"brawling", "Brawling", 0.30},
		{"tech", "Tech Weapons", 0.15},
		{"light_armor", "Light Armor", 0.20},
		{"cloth_armor", "Cloth Armor", 0.25},
		{"heavy_armor", "Heavy Armor", 0.15},
	}

	// Level thresholds: 10 levels per category
	// Each level requires progressively more XP and provides better multipliers
	type threshold struct {
		level             int
		xpRequired       int
		damageMultiplier float64
		defenseMultiplier float64
	}

	defaultThresholds := []threshold{
		{1, 0, 1.0, 1.0},
		{2, 100, 1.05, 1.05},
		{3, 300, 1.10, 1.10},
		{4, 600, 1.15, 1.15},
		{5, 1000, 1.20, 1.20},
		{6, 1500, 1.30, 1.25},
		{7, 2200, 1.40, 1.30},
		{8, 3000, 1.50, 1.35},
		{9, 4000, 1.65, 1.40},
		{10, 5500, 1.80, 1.50},
	}

	for _, cat := range categories {
		created, err := client.CompetencyCategory.Create().
			SetID(cat.ID).
			SetName(cat.Name).
			SetXpMultiplier(cat.XpMultiplier).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create competency category %s: %v", cat.ID, err)
			continue
		}

		for _, t := range defaultThresholds {
			_, err := client.CompetencyLevelThreshold.Create().
				SetID(cat.ID + "-" + strconv.Itoa(t.level)).
				SetLevel(t.level).
				SetXpRequired(t.xpRequired).
				SetDamageMultiplier(t.damageMultiplier).
				SetDefenseMultiplier(t.defenseMultiplier).
				SetCategoryID(created.ID).
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to create threshold %s-%d: %v", cat.ID, t.level, err)
			}
		}
	}

	log.Printf("Seeded %d competency categories with level thresholds", len(categories))
	return nil
}