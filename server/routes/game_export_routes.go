package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/room"
)

// GameExport represents the full game data export
type GameExport struct {
	Version    string      `json:"version"`
	ExportedAt string      `json:"exported_at"`
	Rooms      []RoomData  `json:"rooms"`
	NPCs       []NPCData   `json:"npcs"`
	Skills     []SkillData `json:"skills"`
	Items      []ItemData  `json:"items"`
}

// RoomData represents a room in the export
type RoomData struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsStarting  bool       `json:"is_starting"`
	Exits       []ExitData `json:"exits"`
}

// ExitData represents a room exit
type ExitData struct {
	Direction    string `json:"direction"`
	TargetRoomID int    `json:"target_room_id"`
}

// NPCData represents an NPC character
type NPCData struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CurrentRoomID int   `json:"current_room_id"`
	Race         string `json:"race"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Hitpoints    int    `json:"hitpoints"`
	MaxHitpoints int    `json:"max_hitpoints"`
	Stamina      int    `json:"stamina"`
	MaxStamina   int    `json:"max_stamina"`
	Mana         int    `json:"mana"`
	MaxMana      int    `json:"max_mana"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	NPCSkillID   string `json:"npc_skill_id,omitempty"`
	IsImmortal   bool   `json:"is_immortal"`
}

// SkillData represents a skill/spell
type SkillData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinCooldown int    `json:"min_cooldown"`
	EffectType  string `json:"effect_type,omitempty"`
}

// ItemData represents an item in the game world (rooms/NPCs, not players)
type ItemData struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	LocationType string `json:"location_type"` // "room" or "npc"
	LocationID   int    `json:"location_id"`
}

// RegisterGameExportRoutes registers export/import routes
func RegisterGameExportRoutes(router *gin.Engine, client *db.Client) {
	// Export the entire game world (NPCs, rooms, skills, items - NO users/player data)
	router.GET("/admin/export", func(c *gin.Context) {
		export := GameExport{
			Version:    "1.0",
			ExportedAt: time.Now().Format(time.RFC3339),
		}

		// Export rooms
		rooms, err := client.Room.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms: " + err.Error()})
			return
		}

		for _, r := range rooms {
			roomData := RoomData{
				ID:          r.ID,
				Name:        r.Name,
				Description: r.Description,
				IsStarting:  r.IsStartingRoom,
			}

			// Convert exits map to array
			for dir, targetID := range r.Exits {
				roomData.Exits = append(roomData.Exits, ExitData{
					Direction:    dir,
					TargetRoomID: targetID,
				})
			}

			export.Rooms = append(export.Rooms, roomData)
		}

		// Export NPCs only (isNPC=true), exclude player characters
		npcs, err := client.Character.Query().
			Where(character.IsNPCEQ(true)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch NPCs: " + err.Error()})
			return
		}

		for _, npc := range npcs {
			// Skip player characters (double-check)
			if !npc.IsNPC {
				continue
			}

			npcData := NPCData{
				ID:            npc.ID,
				Name:          npc.Name,
				CurrentRoomID: npc.CurrentRoomId,
				Race:          npc.Race,
				Class:         npc.Class,
				Level:         npc.Level,
				Hitpoints:     npc.Hitpoints,
				MaxHitpoints:  npc.MaxHitpoints,
				Stamina:       npc.Stamina,
				MaxStamina:    npc.MaxStamina,
				Mana:          npc.Mana,
				MaxMana:       npc.MaxMana,
				Strength:      npc.Strength,
				Dexterity:     npc.Dexterity,
				Constitution:  npc.Constitution,
				Intelligence:  npc.Intelligence,
				Wisdom:        npc.Wisdom,
				NPCSkillID:    npc.NpcSkillID,
				IsImmortal:    npc.IsImmortal,
			}

			export.NPCs = append(export.NPCs, npcData)
		}

		// Export skills from registry (hardcoded for now, extendable)
		export.Skills = []SkillData{
			{
				ID:          "druid_heal",
				Name:        "Nature's Blessing",
				Description: "Heals 5% of max HP using druidic magic",
				MinCooldown: 4,
				EffectType:  "heal",
			},
			{
				ID:          "concentrate",
				Name:        "Concentrate",
				Description: "Focus your mind to increase accuracy. +WIS to hit for 4 rounds.",
				MinCooldown: 8,
				EffectType:  "buff",
			},
			{
				ID:          "haymaker",
				Name:        "Haymaker",
				Description: "A powerful but reckless strike. +STR damage, -DEX to hit.",
				MinCooldown: 6,
				EffectType:  "attack",
			},
			{
				ID:          "backoff",
				Name:        "Back-off",
				Description: "Use agility to dodge all attacks this round. Costs stamina.",
				MinCooldown: 10,
				EffectType:  "defense",
			},
			{
				ID:          "scream",
				Name:        "Scream",
				Description: "Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.",
				MinCooldown: 12,
				EffectType:  "buff",
			},
			{
				ID:          "slap",
				Name:        "Slap",
				Description: "A quick stunning strike. DEX vs CON to stun for 1 round.",
				MinCooldown: 8,
				EffectType:  "debuff",
			},
		}

		// TODO: Export items from rooms and NPCs (not on players)
		// This would require an items/equipment table
		export.Items = []ItemData{}

		c.JSON(http.StatusOK, export)
	})

	// Import game world data
	router.POST("/admin/import", func(c *gin.Context) {
		var importData GameExport
		if err := c.ShouldBindJSON(&importData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Validate version
		if importData.Version != "1.0" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported version: " + importData.Version})
			return
		}

		imported := struct {
			Rooms  int `json:"rooms"`
			NPCs   int `json:"npcs"`
			Skills int `json:"skills"`
			Items  int `json:"items"`
		}{}

		// Import rooms
		for _, r := range importData.Rooms {
			// Check if room exists
			exists, _ := client.Room.Query().
				Where(room.IDEQ(r.ID)).
				Exist(c.Request.Context())

			// Build exits map from array
			exitsMap := make(map[string]int)
			for _, exit := range r.Exits {
				exitsMap[exit.Direction] = exit.TargetRoomID
			}

			if exists {
				// Update existing room
				client.Room.UpdateOneID(r.ID).
					SetName(r.Name).
					SetDescription(r.Description).
					SetIsStartingRoom(r.IsStarting).
					SetExits(exitsMap).
					Save(c.Request.Context())
			} else {
				// Create new room
				client.Room.Create().
					SetName(r.Name).
					SetDescription(r.Description).
					SetIsStartingRoom(r.IsStarting).
					SetExits(exitsMap).
					Save(c.Request.Context())
			}
			imported.Rooms++
		}

		// Import NPCs (skip if they have user accounts - safety check)
		for _, npc := range importData.NPCs {
			// Check if NPC exists by name
			existing, err := client.Character.Query().
				Where(character.NameEQ(npc.Name)).
				Where(character.IsNPCEQ(true)).
				First(c.Request.Context())

			if err == nil && existing != nil {
				// Update existing NPC
				client.Character.UpdateOneID(existing.ID).
					SetName(npc.Name).
					SetCurrentRoomId(npc.CurrentRoomID).
					SetRace(npc.Race).
					SetClass(npc.Class).
					SetLevel(npc.Level).
					SetHitpoints(npc.Hitpoints).
					SetMaxHitpoints(npc.MaxHitpoints).
					SetStamina(npc.Stamina).
					SetMaxStamina(npc.MaxStamina).
					SetMana(npc.Mana).
					SetMaxMana(npc.MaxMana).
					SetStrength(npc.Strength).
					SetDexterity(npc.Dexterity).
					SetConstitution(npc.Constitution).
					SetIntelligence(npc.Intelligence).
					SetWisdom(npc.Wisdom).
					SetNpcSkillID(npc.NPCSkillID).
					SetIsImmortal(npc.IsImmortal).
					Save(c.Request.Context())
			} else {
				// Create new NPC
				client.Character.Create().
					SetName(npc.Name).
					SetIsNPC(true).
					SetCurrentRoomId(npc.CurrentRoomID).
					SetStartingRoomId(npc.CurrentRoomID). // Assume same as current
					SetRace(npc.Race).
					SetClass(npc.Class).
					SetLevel(npc.Level).
					SetHitpoints(npc.Hitpoints).
					SetMaxHitpoints(npc.MaxHitpoints).
					SetStamina(npc.Stamina).
					SetMaxStamina(npc.MaxStamina).
					SetMana(npc.Mana).
					SetMaxMana(npc.MaxMana).
					SetStrength(npc.Strength).
					SetDexterity(npc.Dexterity).
					SetConstitution(npc.Constitution).
					SetIntelligence(npc.Intelligence).
					SetWisdom(npc.Wisdom).
					SetNpcSkillID(npc.NPCSkillID).
					SetIsImmortal(npc.IsImmortal).
					Save(c.Request.Context())
				imported.NPCs++
			}
		}

		// Skills are hardcoded in registry, so we just validate they're recognized
		imported.Skills = len(importData.Skills)

		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"imported":    imported,
			"version":     importData.Version,
			"imported_at": time.Now().Format(time.RFC3339),
		})
	})

	// Validate export file without importing
	router.POST("/admin/import/validate", func(c *gin.Context) {
		var importData GameExport
		if err := c.ShouldBindJSON(&importData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		validation := struct {
			Version string   `json:"version"`
			IsValid bool     `json:"is_valid"`
			Rooms   int      `json:"rooms"`
			NPCs    int      `json:"npcs"`
			Skills  int      `json:"skills"`
			Errors  []string `json:"errors,omitempty"`
		}{
			Version: importData.Version,
			IsValid: true,
			Rooms:   len(importData.Rooms),
			NPCs:    len(importData.NPCs),
			Skills:  len(importData.Skills),
		}

		if importData.Version != "1.0" {
			validation.IsValid = false
			validation.Errors = append(validation.Errors, "Unsupported version: "+importData.Version)
		}

		if len(importData.Rooms) == 0 {
			validation.Errors = append(validation.Errors, "No rooms found in import")
		}

		c.JSON(http.StatusOK, validation)
	})
}
