package dbinit

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"herbst-server/db"
	"herbst-server/db/availabletalent"
	"herbst-server/db/character"
	"herbst-server/db/equipment"
	"herbst-server/db/npctemplate"
	"herbst-server/db/room"
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
		
		create := client.Character.
			Create().
			SetName(name).
			SetIsNPC(true).
			SetCurrentRoomId(randomRoom.ID).
			SetStartingRoomId(randomRoom.ID).
			SetIsAdmin(false)
		
		// Aragorn gets Druid healing skill
		if name == "Aragorn" {
			create = create.SetNpcSkillID("druid_heal")
		}
		
		_, err := create.Save(ctx)
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

		// Create Combat Dummy - immortal training dummy (takes damage but never dies)
		_, err = client.Character.
			Create().
			SetName("Combat Dummy").
			SetIsNPC(true).
			SetIsImmortal(true).
			SetCurrentRoomId(holeRoom.ID).
			SetStartingRoomId(holeRoom.ID).
			SetIsAdmin(false).
			SetHitpoints(100).
			SetMaxHitpoints(100).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to create Combat Dummy: %v", err)
		} else {
			log.Println("Created immortal Combat Dummy in 'The Hole' for training")
		}
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

	// Store fountain room ID in game_config
	if err := SetFountainRoomID(ctx, client, fountainRoom.ID); err != nil {
		log.Printf("Warning: failed to set fountain_room_id in game_config: %v", err)
	} else {
		log.Printf("Fountain room ID %d stored in game_config", fountainRoom.ID)
	}

	// Note: StartingRoomID constant in herbst/main.go will need updating
	log.Println("Fountain and Fountain Plaza rooms initialized successfully")
	return nil
}

// InitGizmoNPC creates the Gizmo NPC template and spawns Gizmo in the fountain room
func InitGizmoNPC(client *db.Client) error {
	ctx := context.Background()

	// Check if gizmo NPC already exists
	existingNPCs, err := client.Character.Query().Where(character.IsNPC(true)).All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query existing NPCs: %w", err)
	}
	for _, npc := range existingNPCs {
		if npc.Name == "Gizmo" {
			log.Println("Gizmo NPC already exists, skipping...")
			return nil
		}
	}

	// Get the fountain room
	fountainRoom, err := client.Room.Query().Where(room.NameEQ("The Fountain")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find fountain room: %w", err)
	}

	// Create NPC template for Gizmo
	_, err = client.NPCTemplate.Create().
		SetID("gizmo").
		SetName("Gizmo").
		SetDescription("A friendly half-dog creature with soulful eyes and wagging tail. Looks eager to help.").
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
	log.Println("Gizmo NPC template created")

	// Get the NPC template
	gizmoTemplate, err := client.NPCTemplate.Get(ctx, "gizmo")
	if err != nil {
		return fmt.Errorf("failed to get Gizmo template: %w", err)
	}

	// Create Gizmo character in the fountain room
	_, err = client.Character.Create().
		SetName("Gizmo").
		SetIsNPC(true).
		SetCurrentRoomId(fountainRoom.ID).
		SetStartingRoomId(fountainRoom.ID).
		SetLevel(1).
		SetRace("half-dog").
		SetClass("adventurer").
		SetNpcTemplate(gizmoTemplate).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Gizmo character: %w", err)
	}

	log.Printf("Gizmo NPC spawned in fountain room (ID: %d)", fountainRoom.ID)
	log.Println("Gizmo NPC initialization complete")
	return nil
}

// InitJunkyard creates the Junkyard area - a newbie-friendly zone east of the Fountain
func InitJunkyard(client *db.Client) error {
	ctx := context.Background()

	// Check if junkyard already exists (check for "Junkyard Entrance" room)
	existingRooms, err := client.Room.Query().Where(room.NameEQ("Junkyard Entrance")).Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for junkyard rooms: %w", err)
	}

	if existingRooms > 0 {
		log.Println("Junkyard already initialized, skipping...")
		return nil
	}

	// Get the Fountain Plaza room to connect the exit
	fountainPlaza, err := client.Room.Query().Where(room.NameEQ("Fountain Plaza")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to find Fountain Plaza: %w", err)
	}

	// Room type descriptions for randomization
	roomTypes := []struct {
		name        string
		description string
	 атмосфера   string
	}{
		{
			name:        "Scrap Heap",
			description: "Piles of rusted metal and twisted machinery rise like mountains. Old wires hang from the ceiling.",
			атмосфера:   "The air smells of oil and decay. You hear the creak of metal in the wind.",
		},
		{
			name:        "Golem Nest",
			description: "A dark cavern filled with glowing red eyes. Scrap metal lines the walls like strange art.",
			атмосфера:   "Heat radiates from dormant Rust Bucket Golems. Sparks occasionally fly from their joints.",
		},
		{
			name:        "Broken Equipment Graveyard",
			description: "Rows of defunct machines stretch into the darkness. A graveyard of pre-Ooze technology.",
			атмосфера:   "Echoes bounce off the metallic surfaces. Dripping water echoes from somewhere nearby.",
		},
		{
			name:        "Hidden Cache",
			description: "A small, cramped space tucked behind collapsed shelving. Someone stored supplies here once.",
			атмосфера:   "Dust motes dance in the dim light. The walls are covered in faded warning signs.",
		},
		{
			name:        "Exit Corridor",
			description: "A passage leading back toward the entrance. Faint light filters in from the east.",
			атмосфера:   "A cool breeze carries the scent of fresh air. You're getting close to the exit.",
		},
	}

	// Create 5x5 = 25 rooms in a grid
	// Layout: rows 0-4, columns 0-4
	// Entrance is at (2, 0) - middle of west wall, connects to Fountain Plaza
	type gridRoom struct {
		id       int
		row      int
		col      int
		name     string
		desc     string
		atmosphere string
	}

	grid := make([][]*gridRoom, 5)
	for i := range grid {
		grid[i] = make([]*gridRoom, 5)
	}

	// Create all rooms in the grid
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			roomType := roomTypes[rand.Intn(len(roomTypes))]
			
			// Customize based on position
			name := roomType.name
			desc := roomType.description
			atmosphere := roomType.атмосфера

			// Make the entrance room distinct
			if row == 2 && col == 0 {
				name = "Junkyard Entrance"
				desc = "The mouth of the Junkyard gapes before you. Twisted pipes and rusted metal form an archway. Old signs warning of 'DANGER' hang at odd angles. To the EAST, the Junkyard stretches on. The Fountain Plaza lies to the WEST."
				atmosphere = "The smell of rust and old oil fills your nostrils. Distant clanks and groans echo from deep within."
			}

			// Make the exit corridor distinctive
			if row == 2 && col == 4 {
				name = "Junkyard Exit"
				desc = "The eastern edge of the Junkyard. Faint daylight filters down from above. A ladder leads UP to the surface sewers. WEST leads deeper into the Junkyard."
				atmosphere = "A cool breeze flows from above. You can almost breathe freely here."
			}

			room, err := client.Room.
				Create().
				SetName(name).
				SetDescription(desc).
				SetAtmosphere(room.AtmosphereWind).
				SetIsStartingRoom(false).
				SetExits(map[string]int{}).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create room at (%d,%d): %w", row, col, err)
			}

			grid[row][col] = &gridRoom{
				id:       room.ID,
				row:      row,
				col:      col,
				name:     name,
				desc:     desc,
				atmosphere: atmosphere,
			}
		}
	}

	// Connect rooms with exits
	directions := []struct {
		dx, dy     int
		dir        string
		opp        string
	}{
		{0, -1, "north", "south"},
		{0, 1, "south", "north"},
		{-1, 0, "west", "east"},
		{1, 0, "east", "west"},
	}

	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if grid[row][col] == nil {
				continue
			}

			exits := make(map[string]int)

			for _, d := range directions {
				newRow, newCol := row+d.dx, col+d.dy
				if newRow >= 0 && newRow < 5 && newCol >= 0 && newCol < 5 && grid[newRow][newCol] != nil {
					exits[d.dir] = grid[newRow][newCol].id
				}
			}

			// Entrance room (2,0) also connects to Fountain Plaza to the west
			if row == 2 && col == 0 {
				exits["west"] = fountainPlaza.ID
			}

			// Exit room (2,4) has an up exit to surface
			if row == 2 && col == 4 {
				// "up" will be handled as a special case in navigation
			}

			err = client.Room.UpdateOneID(grid[row][col].id).
				SetExits(exits).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to set exits for room at (%d,%d): %w", row, col, err)
			}
		}
	}

	log.Printf("Junkyard area created: 25 rooms (5x5 grid)")
	log.Println("Junkyard initialization complete")
	return nil
}

// InitSkills seeds the master skills table from the spike design
func InitSkills(client *db.Client) error {
	ctx := context.Background()

	// Check if skills already exist
	existingSkills, err := client.Skill.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count existing skills: %w", err)
	}

	if existingSkills > 0 {
		log.Println("Skills already exist, skipping seed...")
		return nil
	}

	// Skills from SKILLS_SPIKE.md and CLASS_SYSTEM_SPIKE.md
	skills := []struct {
		name        string
		description string
		skillType   string
	}{
		// Weapon skills
		{"blades", "Proficiency with swords, machetes, cleavers - affects damage and accuracy with blade weapons", "weapon"},
		{"staves", "Proficiency with polearms, spears, bows - affects damage and range with pole weapons", "weapon"},
		{"knives", "Proficiency with daggers, sais, small blades - affects damage and critical hits with knives", "weapon"},
		{"martial", "Proficiency with nunchukus, shuriken, tonfas - affects damage and special moves with exotic weapons", "weapon"},
		{"brawling", "Proficiency with fists, improvised weapons - fallback combat skill", "weapon"},
		{"tech", "Proficiency with laser weapons, gadgets - affects accuracy and damage with tech weapons", "weapon"},
		// Magic skills
		{"fire_magic", "Proficiency with fire spells", "magic"},
		{"water_magic", "Proficiency with water spells", "magic"},
		{"wind_magic", "Proficiency with wind spells", "magic"},
		// Armor skills
		{"light_armor", "Proficiency with light armor - affects dodge bonus while wearing light armor", "armor"},
		{"cloth_armor", "Proficiency with cloth armor - minimal protection but wide availability", "armor"},
		{"heavy_armor", "Proficiency with heavy armor - affects defense while reducing mobility", "armor"},
	}

	for _, s := range skills {
		_, err := client.Skill.Create().
			SetName(s.name).
			SetDescription(s.description).
			SetSkillType(s.skillType).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create skill %s: %w", s.name, err)
		}
	}

	log.Printf("Seeded %d skills", len(skills))
	return nil
}

// InitTalents seeds the master talents table from the spike design
func InitTalents(client *db.Client) error {
	ctx := context.Background()

	// Check if talents already exist
	existingTalents, err := client.Talent.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count existing talents: %w", err)
	}

	if existingTalents > 0 {
		log.Println("Talents already exist, skipping seed...")
		return nil
	}

	// Talents from CLASS_SYSTEM_SPIKE.md with unified effect system
	talents := []struct {
		name          string
		description   string
		requirements  string
		effectType    string
		effectValue   int
		effectDuration int
		cooldown      int
		manaCost      int
		staminaCost   int
	}{
		// Combat talents - damage dealers
		{"slash", "Basic sword/blade attack", `{"skills":["blades","knives"]}`, "damage", 10, 0, 1, 0, 5},
		{"smash", "Powerful blunt attack", `{"skills":["staves","martial"]}`, "damage", 15, 0, 2, 0, 10},
		{"heavy_strike", "Strong but slow attack", `{"skills":["blades","staves"]}`, "damage", 20, 0, 3, 0, 15},
		{"quick_slash", "Fast attack with lower damage", `{"skills":["knives"]}`, "damage", 6, 0, 0, 0, 3},
		{"crash", "Body damage based on weight (STR), no weapon required", "[]", "damage", 12, 0, 2, 0, 8},

		// Defensive talents - buffs
		{"parry", "Deflect incoming attacks", "[]", "buff_dodge", 20, 2, 3, 0, 5},
		{"shield_wall", "Increase defense for one round", "[]", "buff_armor", 10, 2, 4, 0, 10},
		{"dodge", "Increase dodge chance for one round", "[]", "buff_dodge", 30, 2, 3, 0, 5},
		{"iron_will", "Resists stun/blind effects (passive)", "[]", "buff_armor", 5, 999, 0, 0, 0},

		// Combat talents - offensive buffs
		{"focus", "Increase critical chance for next attack", "[]", "buff_crit", 15, 2, 3, 0, 5},
		{"hail_storm", "Double attacks for 2 cycles", "[]", "buff_crit", 50, 2, 10, 5, 20},
		{"battle_cry", "Demoralize enemies, reduces their accuracy", "[]", "debuff", 10, 2, 5, 0, 15},

		// Healing talents
		{"second_wind", "Recover HP when low", "[]", "heal", 25, 0, 5, 0, 15},

		// Damage over time
		{"shield_bash", "Bash with shield, has stun chance", "[]", "dot", 5, 3, 4, 0, 12},
	}

	for _, t := range talents {
		_, err := client.Talent.Create().
			SetName(t.name).
			SetDescription(t.description).
			SetRequirements(t.requirements).
			SetEffectType(t.effectType).
			SetEffectValue(t.effectValue).
			SetEffectDuration(t.effectDuration).
			SetCooldown(t.cooldown).
			SetManaCost(t.manaCost).
			SetStaminaCost(t.staminaCost).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create talent %s: %w", t.name, err)
		}
	}

	log.Printf("Seeded %d talents", len(talents))
	return nil
}

// InitAvailableTalentsForCharacter unlocks default talents for a new character based on their class
func InitAvailableTalentsForCharacter(client *db.Client, charID int, charClass string) error {
	ctx := context.Background()

	// Map of class -> default available talents
	classTalents := map[string][]string{
		"warrior":    {"slash", "parry", "smash", "crash"},
		"chef":       {"slash", "second_wind", "battle_cry", "dodge"},
		"mystic":     {"parry", "second_wind", "iron_will", "focus"},
		"survivor":   {"slash", "parry", "crash", "dodge"},
		"brawler":    {"crash", "parry", "dodge", "heavy_strike"},
		"tinkerer":   {"parry", "dodge", "focus", "shield_bash"},
		"trader":     {"parry", "battle_cry", "second_wind", "dodge"},
		"vine_climber": {"parry", "dodge", "quick_slash", "focus"},
	}

	// Get talents for this class (default to survivor if not found)
	talentNames, ok := classTalents[charClass]
	if !ok {
		talentNames = classTalents["survivor"]
	}

	// Get character level
	char, err := client.Character.Get(ctx, charID)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// Add each available talent
	for _, talentName := range talentNames {
		// Find the talent
		talentObj, err := client.Talent.Query().Where(talent.NameEQ(talentName)).Only(ctx)
		if err != nil {
			log.Printf("Warning: talent %s not found, skipping", talentName)
			continue
		}

		// Check if already available
		existing, err := client.AvailableTalent.Query().
			Where(availabletalent.HasCharacterWith(character.ID(charID))).
			Where(availabletalent.HasTalentWith(talent.IDEQ(talentObj.ID))).
			Exist(ctx)
		if err == nil && existing {
			continue
		}

		// Create available talent
		_, err = client.AvailableTalent.Create().
			SetCharacterID(charID).
			SetTalentID(talentObj.ID).
			SetUnlockReason("class_default").
			SetUnlockedAtLevel(char.Level).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to add available talent %s: %v", talentName, err)
		}
	}

	log.Printf("Initialized available talents for character %d (class: %s)", charID, charClass)
	return nil
}

// InitCharacterHealth ensures all characters have valid HP values
// Heals any characters with HP <= 0 back to their max HP
func InitCharacterHealth(client *db.Client) error {
	ctx := context.Background()

	// Find all characters with HP <= 0
	characters, err := client.Character.Query().
		Where(character.HitpointsLTE(0)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query characters with low HP: %w", err)
	}

	if len(characters) == 0 {
		log.Println("All characters have valid HP")
		return nil
	}

	// Heal each character to max HP (skip immortal characters and invincible 0 HP characters)
	for _, char := range characters {
		// Skip immortal characters (cannot die, special handling)
		if char.IsImmortal {
			continue
		}

		// Skip invincible characters (0 max HP means invincible - old system)
		if char.MaxHitpoints == 0 {
			continue
		}

		maxHP := char.MaxHitpoints
		if maxHP <= 0 {
			maxHP = 100 // Default max HP
		}

		_, err := client.Character.UpdateOneID(char.ID).
			SetHitpoints(maxHP).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to heal character %s (ID: %d): %v", char.Name, char.ID, err)
			continue
		}
		log.Printf("Healed character %s (ID: %d) to %d HP", char.Name, char.ID, maxHP)
	}

	log.Printf("Healed %d characters with invalid HP", len(characters))
	return nil
}

// InitConsumables creates consumable items like health potions
func InitConsumables(client *db.Client) error {
	ctx := context.Background()

	// Define all potions
	potions := []struct {
		name          string
		description   string
		effectType    string
		effectValue   int
		effectDuration int
	}{
		{
			name:          "tiny green potion",
			description:   "A tiny green glass vial containing a weak healing elixir. Restores 15 HP when consumed. Press R in combat to use.",
			effectType:    "heal",
			effectValue:   15,
			effectDuration: 0,
		},
		{
			name:          "small green potion",
			description:   "A small green glass vial containing a healing elixir. Restores 25 HP when consumed. Press R in combat to use.",
			effectType:    "heal",
			effectValue:   25,
			effectDuration: 0,
		},
	}

	for _, p := range potions {
		// Check if this potion already exists
		existing, err := client.Equipment.Query().
			Where(equipment.NameEQ(p.name)).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("failed to check for existing potion %s: %w", p.name, err)
		}

		if existing > 0 {
			continue
		}

		// Create the potion
		_, err = client.Equipment.Create().
			SetName(p.name).
			SetDescription(p.description).
			SetItemType("potion").
			SetSlot("consumable").
			SetEffectType(p.effectType).
			SetEffectValue(p.effectValue).
			SetEffectDuration(p.effectDuration).
			SetIsVisible(true).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", p.name, err)
		}

		log.Printf("Created %s consumable", p.name)
	}

	return nil
}

// GivePotionToCharacter gives a health potion to a specific character
func GivePotionToCharacter(client *db.Client, characterID int) error {
	ctx := context.Background()

	// Check if character already has a potion
	existingPotions, err := client.Equipment.Query().
		Where(
			equipment.OwnerIdEQ(characterID),
			equipment.ItemTypeEQ("potion"),
		).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for existing potions: %w", err)
	}

	if existingPotions > 0 {
		log.Printf("Character %d already has potions, skipping...", characterID)
		return nil
	}

	// Create a potion for the character
	_, err = client.Equipment.Create().
		SetName("small green potion").
		SetDescription("A small green glass vial containing a healing elixir. Restores 25 HP when consumed. Press R in combat to use.").
		SetItemType("potion").
		SetSlot("consumable").
		SetEffectType("heal").
		SetEffectValue(25).
		SetEffectDuration(0).
		SetIsVisible(true).
		SetOwnerId(characterID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create potion for character: %w", err)
	}

	log.Printf("Created small green potion for character %d", characterID)
	return nil
}
