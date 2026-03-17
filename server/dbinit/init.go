package dbinit

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/npctemplate"
	"herbst-server/db/room"
	"herbst-server/db/skill"
	"herbst-server/db/talent"
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

// InitSkillsAndTalents seeds the skills and talents master data
func InitSkillsAndTalents(client *db.Client) error {
	ctx := context.Background()

	// Seed Skills
	skills := []struct {
		name        string
		description string
		skillType   string
		cost        int
	}{
		{"blades", "Proficiency with blade weapons (swords, etc.)", "weapon", 0},
		{"staves", "Proficiency with staff weapons", "weapon", 0},
		{"knives", "Proficiency with knife weapons", "weapon", 0},
		{"martial", "Martial arts proficiency (unarmed)", "weapon", 0},
		{"brawling", "Brawling and close combat", "weapon", 0},
		{"tech", "Technical devices and gadgets", "weapon", 0},
		{"fire_magic", "Fire magic proficiency", "magic", 0},
		{"water_magic", "Water magic proficiency", "magic", 0},
		{"wind_magic", "Wind magic proficiency", "magic", 0},
	}

	for _, s := range skills {
		existing, err := client.Skill.Query().Where(skill.Name(s.name)).Exist(ctx)
		if err == nil && existing {
			continue
		}
		_, err = client.Skill.Create().
			SetName(s.name).
			SetDescription(s.description).
			SetSkillType(s.skillType).
			SetCost(s.cost).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create skill %s: %v", s.name, err)
		}
	}

	// Seed Talents
	talents := []struct {
		name        string
		description string
		requirements string
	}{
		{"slash", "A quick blade slash attack", `{"required_skill": "blades"}`},
		{"parry", "Block an incoming attack", "{}"},
		{"smash", "A powerful crushing blow", `{"required_skill": "staves"}`},
		{"crash", "Body slam attack", "{}"},
		{"shield_bash", "Strike with your shield", "{}"},
		{"battle_cry", "Boost your combat morale", "{}"},
		{"second_wind", "Recover stamina in combat", "{}"},
		{"hail_storm", "Magic ice attack", `{"required_skill": "wind_magic"}`},
		{"iron_will", "Mental defense boost", "{}"},
		{"heavy_strike", "A powerful weapon strike", `{"required_skill": "blades"}`},
	}

	for _, t := range talents {
		existing, err := client.Talent.Query().Where(talent.Name(t.name)).Exist(ctx)
		if err == nil && existing {
			continue
		}
		_, err = client.Talent.Create().
			SetName(t.name).
			SetDescription(t.description).
			SetRequirements(t.requirements).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create talent %s: %v", t.name, err)
		}
	}

	log.Println("Skills and talents seeded successfully")
	return nil
}

// InitNPCTemplates creates the NPC templates and spawns NPCs in rooms
func InitNPCTemplates(client *db.Client) error {
	ctx := context.Background()

	// Create Gizmo NPC template if not exists
	existingTemplate, err := client.NPCTemplate.Get(ctx, "gizmo")
	if err == nil && existingTemplate != nil {
		log.Println("NPC template 'gizmo' already exists, skipping...")
	} else {
		_, err = client.NPCTemplate.
			Create().
			SetID("gizmo").
			SetName("Gizmo").
			SetDescription("A friendly half-dog mutant with floppy ears and wagging tail. He seems eager to help newcomers.").
			SetRace("half-dog").
			SetDisposition(npctemplate.DispositionFriendly).
			SetLevel(1).
			SetGreeting("Welcome, new traveler! I'm Gizmo, here to help you get started. Type 'help' to see what you can do!").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create gizmo npc template: %w", err)
		}
		log.Println("NPC template 'gizmo' created successfully")
	}

	// Spawn Gizmo in the Fountain room
	fountainRoom, err := client.Room.Query().Where(room.NameEQ("The Fountain")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find fountain room: %w", err)
	}

	existingGizmo, err := client.Character.Query().Where(character.HasNpcTemplateWith(npctemplate.IDEQ("gizmo"))).Only(ctx)
	if err == nil && existingGizmo != nil {
		log.Println("NPC Gizmo already spawned, skipping...")
	} else {
		_, err = client.Character.
			Create().
			SetName("Gizmo").
			SetRace("half-dog").
			SetClass("tinkerer").
			SetGender("non-binary").
			SetDescription("A friendly half-dog mutant with floppy ears and wagging tail.").
			SetLevel(1).
			SetCurrentRoomId(fountainRoom.ID).
			SetStartingRoomId(fountainRoom.ID).
			SetIsNPC(true).
			SetNpcTemplateID("gizmo").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to spawn gizmo npc: %w", err)
		}
		log.Printf("NPC Gizmo spawned in room %s (ID: %d)", fountainRoom.Name, fountainRoom.ID)
	}

	log.Println("NPC templates and spawns initialized successfully")
	return nil
}

// InitFountainItem creates the Stone Fountain as an immovable room item
func InitFountainItem(client *db.Client) error {
	ctx := context.Background()

	// Find the starting room (try ID 5 first, then by name)
	var startingRoom *db.Room
	var err error

	startingRoom, err = client.Room.Get(ctx, 5)
	if err != nil {
		startingRoom, err = client.Room.Query().Where(room.NameEQ("Fountain Courtyard")).Only(ctx)
		if err != nil {
			startingRoom, err = client.Room.Query().Where(room.NameEQ("Fountain Plaza")).Only(ctx)
			if err != nil {
				startingRoom, err = client.Room.Query().Where(room.NameEQ("The Fountain")).Only(ctx)
				if err != nil {
					log.Println("Warning: could not find starting/fountain room for fountain item")
					return nil
				}
			}
		}
	}

	// Check if fountain item already exists in the room
	items, err := client.Room.QueryEquipment(startingRoom).All(ctx)
	if err == nil {
		for _, item := range items {
			if item.Name == "Stone Fountain" {
				log.Println("Fountain item already exists, skipping...")
				return nil
			}
		}
	}

	// Create the fountain as an immovable room item
	_, err = client.Equipment.Create().
		SetName("Stone Fountain").
		SetDescription("A weathered stone fountain with crystal-clear water bubbling gently. Strange runes are carved into its basin, glowing faintly with an eerie green light.").
		SetSlot("room").
		SetLevel(1).
		SetWeight(0).
		SetIsEquipped(false).
		SetIsImmovable(true).
		SetColor("220").
		SetIsVisible(true).
		SetItemType("furniture").
		SetMinDamage(0).
		SetMaxDamage(0).
		SetWeaponType("none").
		SetClassRestriction("all").
		SetIsDroppable(false).
		SetGuaranteedDrop(false).
		SetRoom(startingRoom).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create fountain item: %w", err)
	}

	log.Printf("Fountain item created in room %s (ID: %d)", startingRoom.Name, startingRoom.ID)
	return nil
}

// InitJunkyard creates the Junkyard area - a newbie-friendly zone east of the fountain
func InitJunkyard(client *db.Client) error {
	ctx := context.Background()

	// Check if junkyard already exists (check for entrance room)
	existingRoom, err := client.Room.Query().Where(room.NameEQ("Junkyard Entrance")).Only(ctx)
	if err == nil && existingRoom != nil {
		log.Println("Junkyard area already initialized, skipping...")
		return nil
	}

	log.Println("Creating Junkyard area (5x5 grid, 25 rooms)...")

	// First, find the fountain room to connect entrance
	var fountainRoom *db.Room
	fountainRoom, err = client.Room.Query().Where(room.NameContains("Fountain")).Only(ctx)
	if err != nil {
		log.Printf("Warning: could not find Fountain room, creating unconnected junkyard: %v", err)
		// Continue anyway - junkyard will exist but not be connected
	}

	// Create 25 rooms for the 5x5 junkyard grid
	// We'll store them in a map for exit connection
	roomMap := make(map[string]*db.Room)

	roomTypes := []string{
		"Scrap Heap", "Scrap Heap", "Scrap Heap", "Scrap Heap", "Scrap Heap",
		"Golem Nest", "Golem Nest", "Broken Equipment", "Broken Equipment", "Hidden Cache",
		"Scrap Heap", "Broken Equipment", "Golem Nest", "Broken Equipment", "Scrap Heap",
		"Hidden Cache", "Broken Equipment", "Scrap Heap", "Golem Nest", "Scrap Heap",
		"Exit Corridor", "Scrap Heap", "Broken Equipment", "Scrap Heap", "Scrap Heap",
	}

	descriptions := map[string]string{
		"Scrap Heap":      "Piles of rusted metal and broken machinery rise around you. Twisted pipes hang from the ceiling, dripping slowly. The air is thick with the smell of old oil and decay.",
		"Golem Nest":      "A cavernous area filled with dormant Rust Bucket Golems. Sparks occasionally fly from their cracked exteriors. Not a good place to linger.",
		"Broken Equipment": "Dead-end passage blocked by massive pieces of broken equipment. You can see some useful scrap here if you look carefully.",
		"Hidden Cache":    "A rarely-visited corner of the junkyard. The walls are lined with old storage containers, some still sealed.",
		"Exit Corridor":   "A passage leading back toward the Fountain. You can see light filtering in from the east.",
	}

	// Create all rooms
	for i := 0; i < 25; i++ {
		roomType := roomTypes[i]
		desc, _ := descriptions[roomType]
		
		room, err := client.Room.Create().
			SetName("Junkyard - " + roomType).
			SetDescription(desc).
			SetExits(map[string]int{}). // Will be populated after all rooms created
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create junkyard room %d: %v", i, err)
			continue
		}
		
		// Store room with coordinate key
		row, col := i/5, i%5
		key := string(rune('0'+row)) + string(rune('0'+col))
		roomMap[key] = room
		
		log.Printf("Created junkyard room %d: %s (row=%d, col=%d)", i, roomType, row, col)
	}

	// Now connect exits (each room connects to adjacent rooms)
	directions := []struct {
		deltaRow, deltaCol int
		dir                string
		opp                string
	}{
		{-1, 0, "north", "south"},
		{1, 0, "south", "north"},
		{0, -1, "west", "east"},
		{0, 1, "east", "west"},
	}

	for i := 0; i < 25; i++ {
		row, col := i/5, i%5
		key := string(rune('0'+row)) + string(rune('0'+col))
		room, ok := roomMap[key]
		if !ok || room == nil {
			continue
		}

		exits := make(map[string]int)
		
		for _, d := range directions {
			newRow, newCol := row+d.deltaRow, col+d.deltaCol
			if newRow >= 0 && newRow < 5 && newCol >= 0 && newCol < 5 {
				adjKey := string(rune('0'+newRow)) + string(rune('0'+newCol))
				adjRoom, ok := roomMap[adjKey]
				if ok && adjRoom != nil {
					exits[d.dir] = adjRoom.ID
				}
			}
		}

		// Special case: Entrance room (row 4, col 0) connects to Fountain
		if row == 4 && col == 0 && fountainRoom != nil {
			exits["west"] = fountainRoom.ID
			_, err = room.Update().
				SetDescription("You stand at the entrance to the Junkyard. To the west, you can see the Fountain area. The twisted metal landscape stretches out before you.").
				SetName("Junkyard Entrance").
				SetExits(exits).
				Save(ctx)
		} else {
			_, err = room.Update().SetExits(exits).Save(ctx)
		}
		if err != nil {
			log.Printf("Warning: failed to update exits for room %s: %v", room.Name, err)
		}
	}

	// Find or create NPC template for Rust Bucket Golem
	var npcTemplate *db.NPCTemplate
	npcTemplate, err = client.NPCTemplate.Query().Where(npctemplate.NameEQ("Rust Bucket Golem")).Only(ctx)
	if err != nil {
		// Create the NPC template
		npcTemplate, err = client.NPCTemplate.Create().
			SetID("rust_bucket_golem").
			SetName("Rust Bucket Golem").
			SetDescription("A hulking construct of rusted metal and broken machinery. Its eyes glow with a faint orange light. Despite its dilapidated appearance, it's still functional - and aggressive.").
			SetRace("golem").
			SetDisposition(npctemplate.DispositionHostile).
			SetLevel(1).
			SetSkills(map[string]int{"basic_combat": 1}).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create Rust Bucket Golem template: %v", err)
		} else {
			log.Printf("Created NPC template: %s", npcTemplate.Name)
		}
	}

	// Spawn Rust Bucket Golems in Golem Nest rooms
	gNestRooms := []string{}
	for k, r := range roomMap {
		if len(k) == 2 && r != nil && len(r.Name) > 20 && r.Name[len(r.Name)-10:] == "Golem Nest" {
			gNestRooms = append(gNestRooms, r.Name)
		}
	}

	for _, roomName := range gNestRooms {
		room, err := client.Room.Query().Where(room.NameEQ(roomName)).Only(ctx)
		if err != nil {
			continue
		}
		if npcTemplate != nil {
			_, err = client.Character.Create().
				SetName("Rust Bucket Golem").
				SetLevel(1).
				SetHitpoints(10).
				SetMaxHitpoints(10).
				SetClass("npc").
				SetRace("golem").
				SetGender("none").
				SetCurrentRoomId(room.ID).
				SetNpcTemplate(npcTemplate).
				SetIsNPC(true).
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to spawn golem in %s: %v", roomName, err)
			}
		}
	}

	log.Printf("Junkyard area created successfully with %d rooms", len(roomMap))
	return nil
}
