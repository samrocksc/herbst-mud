package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
	"herbst-server/db/factioncategory"
)

// classDefinitions holds the seed data for each class.
// Classes are now dynamic - stored as factions in the "class" category.
// Only Survivor is seeded as default; other classes are created via admin panel.
var classDefinitions = []struct {
	name        string
	displayName string
	description string
}{
	{name: "survivor", displayName: "Survivor", description: "Versatile generalist, adaptable to any situation."},
}

// InitClasses seeds class faction categories and factions for each known world.
// Idempotent — skips if the class category already exists for a world.
func InitClasses(client *db.Client) error {
	ctx := context.Background()

	// Known world IDs to seed classes for.
	// world "1" = herbst-mud, world "2" = Ooze Surfers.
	worldIDs := []string{"1", "2"}

	for _, worldID := range worldIDs {
		if err := seedClassesForWorld(ctx, client, worldID); err != nil {
			return err
		}
	}
	return nil
}

func seedClassesForWorld(ctx context.Context, client *db.Client, worldID string) error {
	// Check if class category already exists for this world.
	existing, err := client.FactionCategory.Query().
		Where(
			factioncategory.Name("class"),
			factioncategory.WorldID(worldID),
		).
		Count(ctx)
	if err != nil {
		return err
	}
	if existing > 0 {
		log.Printf("Class category already exists for world %s, skipping seed", worldID)
		return nil
	}

	// Create the "class" faction category.
	cat, err := client.FactionCategory.Create().
		SetName("class").
		SetDisplayName("Class").
		SetDescription("Character class — defines a character's role and abilities.").
		SetMaxMemberships(1).
		SetInitialConfig(true).
		SetWorldID(worldID).
		Save(ctx)
	if err != nil {
		return err
	}

	// Create one faction per class.
	for _, def := range classDefinitions {
		_, err := client.Faction.Create().
			SetName(def.name).
			SetDisplayName(def.displayName).
			SetDescription(def.description).
			SetWorldID(worldID).
			SetCategory(cat).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed class %s for world %s: %v", def.name, worldID, err)
		}
	}

	log.Printf("Classes seeded successfully for world %s", worldID)
	return nil
}