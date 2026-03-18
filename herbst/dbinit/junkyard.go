package dbinit

import (
	"context"
	"fmt"
	"log"

	"herbst/db"
	"herbst/db/character"
	"herbst/db/ent/npctemplate"
	"herbst/db/room"
)

// InitJunkyard creates the 5x5 Junkyard area connected to the Fountain Courtyard
func InitJunkyard(client *db.Client) error {
	ctx := context.Background()

	// Check if Junkyard already exists
	existingJunkyard, err := client.Room.Query().Where(room.NameEQ("Junkyard Entrance")).Only(ctx)
	if err == nil && existingJunkyard != nil {
		log.Println("Junkyard area already exists, skipping...")
		return nil
	}

	// Find the Fountain Courtyard
	fountainRoom, err := client.Room.Query().Where(room.NameEQ("Fountain Courtyard")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find Fountain Courtyard: %w", err)
	}

	// Create 5x5 grid of rooms (25 rooms total)
	// We'll create them with east-west connections to the fountain and inter-junkyard connections
	
	// First row (entrance) - connects to fountain
	entranceRoom, err := createJunkyardRoom(client, "Junkyard Entrance", 
		"You stand at the entrance to the Junkyard. The smell of rust and old metal fills the air. " +
		"Twisted pipes hang from the ceiling, and dim light filters in from above. To the WEST is the Fountain Courtyard.",
		0, 0)
	if err != nil {
		return fmt.Errorf("failed to create entrance: %w", err)
	}

	// Connect entrance to fountain (west)
	err = updateRoomExits(entranceRoom, map[string]int{"west": fountainRoom.ID}, client)
	if err != nil {
		return fmt.Errorf("failed to connect entrance to fountain: %w", err)
	}

	// Connect fountain to entrance (east) - preserve existing exits
	currentFountainExits := fountainRoom.Exits
	if currentFountainExits == nil {
		currentFountainExits = make(map[string]int)
	}
	currentFountainExits["east"] = entranceRoom.ID
	err = client.Room.UpdateOne(fountainRoom).SetExits(currentFountainExits).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect fountain to entrance: %w", err)
	}

	// Create 4 more rows (20 rooms)
	grid := make([][]int, 5)
	for y := 0; y < 5; y++ {
		grid[y] = make([]int, 5)
	}
	grid[0][0] = entranceRoom.ID

	// Room names for the grid (excluding entrance)
	roomConfigs := []struct {
		name        string
		description string
		hasScrap    bool
	}{
		{"Scrap Heap", "Piles of twisted metal and rusted machinery rise like small hills. Old cobwebs drape from ceiling pipes.", true},
		{"Golem Nest", "A hollow area filled with discarded machine parts. The ground is covered in oil stains.", false},
		{"Broken Equipment", "Dead ends marked by towering walls of crushed cars. Something glints in the shadows.", false},
		{"Hidden Cache", "A rare open area with ancient machinery standing in a circle, humming with mysterious energy.", false},
		{"Twisted Alley", "Narrow passages between piles of junk. Twisted pipes create a claustrophobic maze.", false},
	}

	roomCount := 1
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if y == 0 && x == 0 { // Entrance already created
				continue
			}

			cfg := roomConfigs[(roomCount-1)%len(roomConfigs)]
			
			// Enhance description for variety
			desc := cfg.description
			if y == 4 { // Bottom - golem territory
				desc += " Rust Bucket Golems roam this area."
			} else if cfg.hasScrap {
				desc += " You can SEARCH these piles for hidden treasures."
			}

			r, err := createJunkyardRoom(client, cfg.name, desc, x, y)
			if err != nil {
				return fmt.Errorf("failed to create room at (%d,%d): %w", x, y, err)
			}
			grid[y][x] = r.ID
			roomCount++
		}
	}

	// Connect all rooms with exits
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			currentID := grid[y][x]
			if currentID == 0 {
				continue
			}
			
			exits := make(map[string]int)

			// West exit
			if x > 0 && grid[y][x-1] != 0 {
				exits["west"] = grid[y][x-1]
			}
			// East exit
			if x < 4 && grid[y][x+1] != 0 {
				exits["east"] = grid[y][x+1]
			}
			// North exit
			if y > 0 && grid[y-1][x] != 0 {
				exits["north"] = grid[y-1][x]
			}
			// South exit
			if y < 4 && grid[y+1][x] != 0 {
				exits["south"] = grid[y+1][x]
			}

			currentRoom, err := client.Room.Get(ctx, currentID)
			if err != nil {
				continue
			}
			err = updateRoomExits(currentRoom, exits, client)
			if err != nil {
				log.Printf("Warning: failed to set exits for room %d: %v", currentID, err)
			}
		}
	}

	log.Printf("Junkyard area created with %d rooms", roomCount)
	return nil
}

// createJunkyardRoom creates a single junkyard room
func createJunkyardRoom(client *db.Client, name, desc string, x, y int) (*db.Room, error) {
	ctx := context.Background()

	room, err := client.Room.
		Create().
		SetName(name).
		SetDescription(desc).
		SetExits(map[string]int{}).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create junkyard room: %w", err)
	}

	return room, nil
}

// updateRoomExits updates a room's exits
func updateRoomExits(r *db.Room, exits map[string]int, client *db.Client) error {
	ctx := context.Background()
	// Merge with existing exits
	currentExits := r.Exits
	if currentExits == nil {
		currentExits = make(map[string]int)
	}
	for k, v := range exits {
		currentExits[k] = v
	}
	return client.Room.UpdateOne(r).SetExits(currentExits).Exec(ctx)
}

// InitRustBucketGolem creates the Rust Bucket Golem NPC template
func InitRustBucketGolem(client *db.Client) error {
	ctx := context.Background()

	// Check if template already exists
	existing, err := client.NPCTemplate.Get(ctx, "rust_bucket")
	if err == nil && existing != nil {
		log.Println("Rust Bucket Golem template already exists, skipping...")
		return ensureRustBucketCharacters(client)
	}

	// Create the template
	_, err = client.NPCTemplate.
		Create().
		SetID("rust_bucket").
		SetName("Rust Bucket Golem").
		SetDescription("A hulking construct of rusted metal and broken machinery. " +
			"Its eyes glow with a faint orange light. It shuffles around, " +
			"adding scrap to its growing body.").
		SetRace("construct").
		SetDisposition(npctemplate.DispositionHostile).
		SetLevel(2).
		SetMaxHealth(15).
		SetCurrentHealth(15).
		SetSkills(map[string]int{
			"bash": 1,
		}).
		SetTradesWith([]string{}).
		SetGreeting("The Rust Bucket Golem clanks threateningly at you!").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Rust Bucket Golem template: %w", err)
	}

	log.Println("Rust Bucket Golem NPC template created")
	return ensureRustBucketCharacters(client)
}

// ensureRustBucketCharacters spawns Rust Bucket Golems in random rooms
func ensureRustBucketCharacters(client *db.Client) error {
	ctx := context.Background()

	// Find all junkyard rooms
	junkyardRooms, err := client.Room.Query().Where(
		room.NameContains("Junkyard")).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to find junkyard rooms: %w", err)
	}

	if len(junkyardRooms) == 0 {
		log.Println("No junkyard rooms found, skipping golem spawn")
		return nil
	}

	// Check how many golem characters already exist
	golemCount, err := client.Character.Query().Where(character.NpcTemplateIDEQ("rust_bucket")).Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count rust bucket golems: %w", err)
	}

	// Spawn up to 3 golems
	targetGolems := 3
	if len(junkyardRooms) < targetGolems {
		targetGolems = len(junkyardRooms)
	}

	// Only spawn if we don't have enough
	if golemCount >= targetGolems {
		log.Printf("Rust Bucket Golems already spawned (%d)", golemCount)
		return nil
	}

	template, err := client.NPCTemplate.Get(ctx, "rust_bucket")
	if err != nil {
		return fmt.Errorf("failed to get rust_bucket template: %w", err)
	}

	// Spawn golems in junkyard rooms
	spawned := 0
	for _, r := range junkyardRooms {
		if spawned >= targetGolems {
			break
		}
		
		// Check if a golem already exists in this room
		existingInRoom, err := client.Character.Query().
			Where(character.NpcTemplateIDEQ("rust_bucket"), character.CurrentRoomIdEQ(r.ID)).
			Exist(ctx)
		if err != nil || existingInRoom {
			continue
		}

		_, err = client.Character.
			Create().
			SetName("Rust Bucket Golem").
			SetIsNPC(true).
			SetCurrentRoomId(r.ID).
			SetStartingRoomId(r.ID).
			SetNpcTemplate(template).
			SetLevel(2).
			SetMaxHealth(15).
			SetCurrentHealth(15).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to spawn golem in room %d: %v", r.ID, err)
			continue
		}
		log.Printf("Spawned Rust Bucket Golem in %s", r.Name)
		spawned++
	}

	return nil
}

// InitJunkyardArea initializes the complete Junkyard area
func InitJunkyardArea(client *db.Client) error {
	log.Println("Initializing Junkyard area...")

	// Create the area
	if err := InitJunkyard(client); err != nil {
		return fmt.Errorf("failed to initialize junkyard: %w", err)
	}

	// Create the enemy
	if err := InitRustBucketGolem(client); err != nil {
		return fmt.Errorf("failed to initialize rust bucket golem: %w", err)
	}

	log.Println("Junkyard area initialization complete!")
	return nil
}