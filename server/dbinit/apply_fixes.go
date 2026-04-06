package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
	"herbst-server/db/character"
)

// ApplyDatabaseFixes applies all necessary fixes to existing database data
// This should be called on server startup to ensure data consistency
func ApplyDatabaseFixes(client *db.Client) error {
	ctx := context.Background()

	// Fix 1: Convert test characters (Aragorn, Legolas, etc.) to NPCs
	testNames := []string{"Aragorn", "Legolas", "Gimli", "Frodo", "Sam"}
	for _, name := range testNames {
		char, err := client.Character.Query().
			Where(character.NameEQ(name)).
			Only(ctx)
		if err != nil {
			continue // Character doesn't exist, skip
		}

		if !char.IsNPC {
			_, err = client.Character.UpdateOneID(char.ID).
				SetIsNPC(true).
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to fix %s: %v", name, err)
			} else {
				log.Printf("Fixed %s: changed from player to NPC", name)
			}
		}
	}

	return nil
}

// EnsureCombatDummyImmortal ensures the Combat Dummy has is_immortal=true with normal HP
// This runs AFTER InitCharacterHealth to prevent that function from changing it
func EnsureCombatDummyImmortal(client *db.Client) error {
	ctx := context.Background()

	// Find Combat Dummy
	dummy, err := client.Character.Query().
		Where(character.NameEQ("Combat Dummy")).
		Only(ctx)
	if err != nil {
		// Combat Dummy doesn't exist, that's ok
		return nil
	}

	// If dummy has 0 HP (old invincible style), convert to new immortal style
	if dummy.Hitpoints == 0 && dummy.MaxHitpoints == 0 {
		_, err = client.Character.UpdateOneID(dummy.ID).
			SetHitpoints(100).
			SetMaxHitpoints(100).
			SetIsImmortal(true).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to convert Combat Dummy to immortal: %v", err)
			return err
		}
		log.Println("Converted Combat Dummy to new immortal system (100 HP, is_immortal=true)")
		return nil
	}

	// Ensure is_immortal flag is set
	if !dummy.IsImmortal {
		_, err = client.Character.UpdateOneID(dummy.ID).
			SetIsImmortal(true).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to set Combat Dummy as immortal: %v", err)
			return err
		}
		log.Println("Set Combat Dummy to immortal (is_immortal=true)")
	}

	return nil
}
