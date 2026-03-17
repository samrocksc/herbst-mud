package dbinit

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"herbst-server/db"
	"herbst-server/db/room"
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

// InitCharacters creates initial characters including test characters and Gandalf the admin NPC
func InitCharacters(client *db.Client) error {
	ctx := context.Background()

	// Check if characters already exist
	existingChars, err := client.Character.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count existing characters: %w", err)
	}

	if existingChars > 0 {
		log.Println("Characters already exist, skipping seed...")
		return nil
	}

	// Get all rooms to assign random rooms to test characters
	rooms, err := client.Room.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rooms: %w", err)
	}

	if len(rooms) == 0 {
		log.Println("No rooms available for character placement, skipping...")
		return nil
	}

	// Create 5 test characters in random rooms
	testCharacterNames := []string{"Aragorn", "Legolas", "Gimli", "Frodo", "Sam"}
	for _, name := range testCharacterNames {
		randomRoom := rooms[rand.Intn(len(rooms))]
		_, err := client.Character.
			Create().
			SetName(name).
			SetIsNPC(false).
			SetCurrentRoomId(randomRoom.ID).
			SetStartingRoomId(randomRoom.ID).
			SetIsAdmin(false).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create character %s: %v", name, err)
		}
	}

	// Find or create the "hole" room (center room with "The Hole" name)
	holeRoom, err := client.Room.Query().Where(room.NameEQ("The Hole")).Only(ctx)
	if err != nil {
		log.Printf("Warning: could not find 'The Hole' room: %v", err)
		// Use center room as fallback
		if len(rooms) > 0 {
			holeRoom = rooms[0]
		}
	}

	// Create Gandalf as admin NPC in the "hole" room
	if holeRoom != nil {
		_, err = client.Character.
			Create().
			SetName("Gandalf").
			SetIsNPC(true).
			SetCurrentRoomId(holeRoom.ID).
			SetStartingRoomId(holeRoom.ID).
			SetIsAdmin(true).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Gandalf: %w", err)
		}
		log.Println("Created Gandalf NPC in 'The Hole' room with admin privileges")
	}

	log.Println("Character seed data initialized successfully")
	return nil
}

// InitFountain creates the Fountain starting area for new character creation
func InitFountain(client *db.Client) error {
	ctx := context.Background()

	// Check if fountain room already exists
	existingRooms, err := client.Room.Query().Where(room.NameEQ("The Fountain")).Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for fountain room: %w", err)
	}

	if existingRooms > 0 {
		log.Println("Fountain already initialized, skipping...")
		return nil
	}

	// Create The Fountain room
	fountainRoom, err := client.Room.
		Create().
		SetName("The Fountain").
		SetDescription("You wake up at a murky fountain, covered in sticky mutant mud. The water glows faintly with an eerie green Ooze color. Your head throbs - you have no memory of how you got here. Something glints in the mud near your hand.").
		SetIsStartingRoom(false).
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create fountain room: %w", err)
	}

	log.Printf("Fountain room created with ID: %d", fountainRoom.ID)

	// Update the main crossway center to point to fountain
	err = client.Room.Update().
		Where(room.NameEQ("The Hole")).
		SetExits(map[string]int{
			"north": 1,
			"south": 2,
			"east":  4,
			"west":  5,
		}).Exec(ctx)
	if err != nil {
		log.Printf("Warning: failed to update center room exits: %v", err)
	}

	// Create the starting room (New Venice - Fountain Plaza)
	// This is where players end up after washing at the fountain
	startingRoom, err := client.Room.
		Create().
		SetName("Fountain Plaza").
		SetDescription("You stand in a dusty plaza dominated by a large stone fountain at its center. The water glows with a faint green Ooze color - the result of the Great Mutagen Spill. Mutant weeds push through cracked cobblestones. The Canal District lies to the east. A path leads north toward the Crossroads.").
		SetIsStartingRoom(true).
		SetExits(map[string]int{
			"east":  0, // Will be canal district
			"north": 3, // Crossroads
		}).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create starting room: %w", err)
	}

	log.Printf("Starting room (Fountain Plaza) created with ID: %d", startingRoom.ID)

	// Note: StartingRoomID constant in herbst/main.go will need updating
	log.Println("Fountain and Fountain Plaza rooms initialized successfully")
	return nil
}
