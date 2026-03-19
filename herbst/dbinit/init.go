package dbinit

import (
	"context"
	"fmt"
	"log"

	"herbst/db"
	"herbst/db/character"
	"herbst/db/equipment"
	"herbst/db/npctemplate"
	"herbst/db/room"
)

const (
	DefaultAdminEmail    = "admin@herbstmud.local"
	DefaultAdminPassword = "herb5t2026!"
)

// InitAdminUser creates a default admin user if none exists
func InitAdminUser(client *db.Client) error {
	ctx := context.Background()

	// Check if admin user already exists
	existingUsers, err := client.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count existing users: %w", err)
	}

	if existingUsers > 0 {
		log.Println("Users already exist, skipping admin seed...")
		return nil
	}

	// Create default admin user
	_, err = client.User.
		Create().
		SetEmail(DefaultAdminEmail).
		SetPassword(DefaultAdminPassword).
		SetIsAdmin(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Default admin user created: %s / %s", DefaultAdminEmail, DefaultAdminPassword)
	return nil
}

// InitCrossWay creates the initial cross-shaped rooms
func InitCrossWay(client *db.Client) error {
	ctx := context.Background()

	// Check if rooms already exist
	existingRooms, err := client.Room.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count existing rooms: %w", err)
	}

	// If rooms already exist, don't recreate them
	if existingRooms > 0 {
		log.Println("Rooms already initialized, skipping...")
		return nil
	}

	// Create the cross-shaped rooms
	// North room
	northRoom, err := client.Room.
		Create().
		SetName("Northern Path").
		SetDescription("A path leading north from the center.").
		SetExits(map[string]int{"south": 3}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create northern room: %w", err)
	}

	// South room
	southRoom, err := client.Room.
		Create().
		SetName("Southern Path").
		SetDescription("A path leading south from the center.").
		SetExits(map[string]int{"north": 3}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create southern room: %w", err)
	}

	// East room
	eastRoom, err := client.Room.
		Create().
		SetName("Eastern Path").
		SetDescription("A path leading east from the center.").
		SetExits(map[string]int{"west": 3}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create eastern room: %w", err)
	}

	// West room
	westRoom, err := client.Room.
		Create().
		SetName("Western Path").
		SetDescription("A path leading west from the center.").
		SetExits(map[string]int{"east": 3}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create western room: %w", err)
	}

	// Center room (the hole)
	centerRoom, err := client.Room.
		Create().
		SetName("The Hole").
		SetDescription("The central hub of the cross-shaped paths. This is where everything connects.").
		SetIsStartingRoom(true).
		SetExits(map[string]int{
			"north": northRoom.ID,
			"south": southRoom.ID,
			"east":  eastRoom.ID,
			"west":  westRoom.ID,
		}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create center room: %w", err)
	}

	// Update exits for the directional rooms to point to the center
	err = client.Room.UpdateOne(northRoom).
		SetExits(map[string]int{"south": centerRoom.ID}).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update northern room exits: %w", err)
	}

	err = client.Room.UpdateOne(southRoom).
		SetExits(map[string]int{"north": centerRoom.ID}).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update southern room exits: %w", err)
	}

	err = client.Room.UpdateOne(eastRoom).
		SetExits(map[string]int{"west": centerRoom.ID}).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update eastern room exits: %w", err)
	}

	err = client.Room.UpdateOne(westRoom).
		SetExits(map[string]int{"east": centerRoom.ID}).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update western room exits: %w", err)
	}

	log.Println("Cross-shaped rooms initialized successfully")
	return nil
}

// InitFountainRoom creates the fountain room if it doesn't exist
func InitFountainRoom(client *db.Client) error {
	ctx := context.Background()

	// Check if fountain room already exists
	existingRoom, err := client.Room.Query().Where(room.NameEQ("Fountain Courtyard")).Only(ctx)
	if err == nil && existingRoom != nil {
		log.Println("Fountain Courtyard room already exists, skipping...")
		return nil
	}

	// Find the center room to connect to
	centerRoom, err := client.Room.Query().Where(room.NameEQ("The Hole")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find center room: %w", err)
	}

	// Create the fountain room
	_, err = client.Room.
		Create().
		SetName("Fountain Courtyard").
		SetDescription("A peaceful courtyard with a beautiful stone fountain in the center. The gentle sound of water is soothing.").
		SetExits(map[string]int{"west": centerRoom.ID}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create fountain room: %w", err)
	}

	// Update the center room to have an east exit to the fountain
	err = client.Room.UpdateOne(centerRoom).
		SetExits(map[string]int{"north": 1, "south": 2, "east": 6, "west": 4}).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update center room exits: %w", err)
	}

	log.Println("Fountain Courtyard room created successfully")
	return nil
}

// InitGizmo creates the Gizmo NPC template and spawns him in the fountain room
func InitGizmo(client *db.Client) error {
	ctx := context.Background()

	// First ensure fountain room exists
	if err := InitFountainRoom(client); err != nil {
		log.Printf("Warning: failed to initialize fountain room: %v", err)
	}

	// Check if Gizmo template already exists
	existingTemplate, err := client.NPCTemplate.Get(ctx, "gizmo")
	if err == nil && existingTemplate != nil {
		log.Println("Gizmo NPC template already exists, skipping...")
		// Still ensure Gizmo character exists in the fountain room
		return ensureGizmoCharacter(client, ctx)
	}

	// Create the Gizmo NPC template
	_, err = client.NPCTemplate.
		Create().
		SetID("gizmo").
		SetName("Gizmo").
		SetDescription("A friendly half-dog creature with soulful eyes and wagging tail.").
		SetRace("half-dog").
		SetDisposition(npctemplate.DispositionFriendly).
		SetLevel(1).
		SetSkills(map[string]int{}).
		SetTradesWith([]string{}).
		SetGreeting("Welcome, new traveler! I'm Gizmo, here to help you get started.").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Gizmo NPC template: %w", err)
	}

	log.Println("Gizmo NPC template created successfully")
	return ensureGizmoCharacter(client, ctx)
}

// ensureGizmoCharacter ensures the Gizmo character exists in the fountain room
func ensureGizmoCharacter(client *db.Client, ctx context.Context) error {
	// Find the fountain room
	fountainRoom, err := client.Room.Query().Where(room.NameEQ("Fountain Courtyard")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find fountain room: %w", err)
	}

	// Check if Gizmo character already exists
	existingGizmo, err := client.Character.Query().Where(character.NameEQ("Gizmo")).Only(ctx)
	if err == nil && existingGizmo != nil {
		log.Println("Gizmo character already exists, skipping...")
		return nil
	}

	// Get the Gizmo template
	gizmoTemplate, err := client.NPCTemplate.Get(ctx, "gizmo")
	if err != nil {
		return fmt.Errorf("failed to get Gizmo template: %w", err)
	}

	// Create the Gizmo character in the fountain room
	_, err = client.Character.
		Create().
		SetName("Gizmo").
		SetIsNPC(true).
		SetCurrentRoomId(fountainRoom.ID).
		SetStartingRoomId(fountainRoom.ID).
		SetNpcTemplate(gizmoTemplate).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Gizmo character: %w", err)
	}

	log.Println("Gizmo character spawned in the Fountain Courtyard")
	return nil
}

// InitWeapons seeds the database with starter weapons
func InitWeapons(client *db.Client) error {
	ctx := context.Background()

	// Check if weapons already exist (by guaranteed drop flag)
	existingWeapons, err := client.Equipment.Query().
		Where(equipment.GuaranteedDropEQ(true)).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to query existing weapons: %w", err)
	}

	if existingWeapons > 0 {
		log.Println("Starter weapons already exist, skipping...")
		return nil
	}

	// Create Rusty Sword (Warrior starter weapon)
	_, err = client.Equipment.
		Create().
		SetName("Rusty Sword").
		SetDescription("An old, worn sword with a dull blade. It's seen better days but will still do the job.").
		SetSlot("weapon").
		SetItemType("weapon").
		SetLevel(1).
		SetWeight(3).
		SetMinDamage(1).
		SetMaxDamage(3).
		SetWeaponType("sword").
		SetGuaranteedDrop(true).
		SetClassRestriction("warrior").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Rusty Sword: %w", err)
	}

	// Create Twisted Pipe (Chef starter weapon)
	_, err = client.Equipment.
		Create().
		SetName("Twisted Pipe").
		SetDescription("A crude weapon made from a twisted junkyard pipe. Dangerous in the right hands.").
		SetSlot("weapon").
		SetItemType("weapon").
		SetLevel(1).
		SetWeight(2).
		SetMinDamage(1).
		SetMaxDamage(2).
		SetWeaponType("pipe").
		SetGuaranteedDrop(true).
		SetClassRestriction("chef").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Twisted Pipe: %w", err)
	}

	log.Println("Starter weapons seeded: Rusty Sword (warrior), Twisted Pipe (chef)")
	return nil
}

// SelectWeaponForClass returns the weapon that should drop for a given character class
func SelectWeaponForClass(class string) string {
	weaponMap := map[string]string{
		"warrior": "Rusty Sword",
		"chef":    "Twisted Pipe",
	}
	if weapon, ok := weaponMap[class]; ok {
		return weapon
	}
	// Default to Rusty Sword for unknown classes
	return "Rusty Sword"
}
