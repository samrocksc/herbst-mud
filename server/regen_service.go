package main

import (
	"context"
	"log"
	"time"

	"herbst-server/db"
	"herbst-server/repository"
	"herbst-server/routes"
	"herbst-server/service"
)


// getConstitutionModifier calculates HP regen from CON (matching client logic)
func getConstitutionModifier(constitution int) int {
	if constitution <= 0 {
		return 1 // Minimum regen
	}
	// Every 2 CON above 8 adds 1 to base regen
	mod := (constitution - 8) / 2
	if mod < 1 {
		mod = 1
	}
	return mod
}

// performRegen calculates and applies regeneration for a character
func performRegen(ctx context.Context, char *db.Character, combatSvc service.CombatService, charRepo repository.CharacterRepo) (*routes.VitalsPayload, error) {
	// Skip if dead or in combat (for now, just skip if dead)
	if char.Hitpoints <= 0 {
		return nil, nil
	}

	// Calculate derived stats (matching character_stats.go)
	maxHP := char.Constitution*10 + char.Level*10
	maxStamina := char.Constitution*5 + char.Level*5
	maxMana := char.Intelligence*5 + char.Level*5

	// Initialize regen amounts
	hpRegen := 0
	staminaRegen := 0
	manaRegen := 0

	// HP Regen (Constitution-based)
	if char.Hitpoints < maxHP {
		hpRegen = getConstitutionModifier(char.Constitution)
	}

	// Stamina Regen (Equal for all - 3 per tick)
	if char.Stamina < maxStamina {
		staminaRegen = 3
	}

	// Mana Regen (Equal for all - 2 per tick)
	if char.Mana < maxMana {
		manaRegen = 2
	}

	// Apply regen if needed
	updated := false
	newHP := char.Hitpoints
	newStamina := char.Stamina
	newMana := char.Mana

	if hpRegen > 0 {
		newHP += hpRegen
		if newHP > maxHP {
			newHP = maxHP
		}
		updated = true
	}

	if staminaRegen > 0 {
		newStamina += staminaRegen
		if newStamina > maxStamina {
			newStamina = maxStamina
		}
		updated = true
	}

	if manaRegen > 0 {
		newMana += manaRegen
		if newMana > maxMana {
			newMana = maxMana
		}
		updated = true
	}

	// Update database if any stats changed
	if updated {
		updates := repository.CharacterUpdates{}
		if newHP != char.Hitpoints {
			updates.Hitpoints = &newHP
		}
		if newStamina != char.Stamina {
			updates.Stamina = &newStamina
		}
		if newMana != char.Mana {
			updates.Mana = &newMana
		}

		_, err := charRepo.Update(ctx, char.ID, updates)
		if err != nil {
			return nil, err
		}
	}

	// Return updated vitals payload
	return &routes.VitalsPayload{
		HP:         newHP,
		MaxHP:      maxHP,
		Stamina:    newStamina,
		MaxStamina: maxStamina,
		Mana:       newMana,
		MaxMana:    maxMana,
	}, nil
}

// regenTick performs one regeneration tick for all connected characters
func regenTick(client *db.Client, repos *repository.Container, services *service.Container) {
	ctx := context.Background()

	// Get all connected characters
	connections := routes.GetConnections()

	// Track rooms that have active players for NPC healing
	activeRooms := make(map[int]bool)

	for userID, wsc := range connections {
		// Get character from DB
		char, err := repos.Character.Get(ctx, wsc.CharacterID)
		if err != nil {
			log.Printf("[regen] failed to get character %d for user %d: %v", wsc.CharacterID, userID, err)
			continue
		}

		// Track room for NPC healing
		activeRooms[char.CurrentRoomId] = true

		// Perform regen
		vitals, err := performRegen(ctx, char, services.Combat, repos.Character)
		if err != nil {
			log.Printf("[regen] failed to regen character %d: %v", char.ID, err)
			continue
		}

		// Send vitals update if there were changes
		if vitals != nil {
			routes.SendVitalsToCharacter(wsc.CharacterID, *vitals)
		}
	}

	// Heal NPCs in rooms with active players
	for roomID := range activeRooms {
		_, err := services.Combat.PassiveHealNPCsInRoom(ctx, roomID)
		if err != nil {
			log.Printf("[regen] failed to heal NPCs in room %d: %v", roomID, err)
		}
	}
}

// StartRegenService launches the background regeneration service
func StartRegenService(repos *repository.Container, services *service.Container, client *db.Client) {
	log.Println("[regen] starting regeneration service")

	go func() {
		ticker := time.NewTicker(6 * time.Second) // Match client interval
		defer ticker.Stop()

		for range ticker.C {
			regenTick(client, repos, services)
		}
	}()
}