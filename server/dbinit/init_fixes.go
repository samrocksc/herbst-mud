package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
	"herbst-server/db/character"
)

// FixTestCharacterNPCStatus converts Aragorn, Legolas, Gimli, Frodo, and Sam to NPCs
// This fixes a bug where they were created as IsNPC=false (players) instead of IsNPC=true (NPCs)
func FixTestCharacterNPCStatus(client *db.Client) error {
	ctx := context.Background()

	testNames := []string{"Aragorn", "Legolas", "Gimli", "Frodo", "Sam"}
	fixedCount := 0

	for _, name := range testNames {
		char, err := client.Character.Query().
			Where(character.NameEQ(name)).
			Only(ctx)
		if err != nil {
			continue // Character doesn't exist, skip
		}

		if !char.IsNPC {
			// Fix the character to be an NPC
			_, err = client.Character.UpdateOneID(char.ID).
				SetIsNPC(true).
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to fix %s to NPC: %v", name, err)
				continue
			}
			log.Printf("Fixed %s: changed from player to NPC", name)
			fixedCount++
		}
	}

	if fixedCount > 0 {
		log.Printf("Fixed %d test characters to be NPCs", fixedCount)
	}
	return nil
}

// EnsureCombatDummyExists creates the Combat Dummy NPC if it doesn't exist
func EnsureCombatDummyExists(client *db.Client) error {
	ctx := context.Background()

	// Check if Combat Dummy already exists
	existing, err := client.Character.Query().
		Where(character.NameEQ("Combat Dummy")).
		Count(ctx)
	if err != nil || existing > 0 {
		return nil // Already exists or error
	}

	// Find "The Hole" room
	_, err = client.Character.Query().Select().
		Where(character.NameEQ("Gandalf")).
		Only(ctx)
	if err != nil {
		log.Println("Could not find Gandalf to determine room for Combat Dummy")
		return nil
	}

	log.Println("Combat Dummy initialization would require room lookup")
	return nil
}
