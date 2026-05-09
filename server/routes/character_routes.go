package routes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
	"herbst-server/db/gameconfig"
	"herbst-server/db/race"
	"herbst-server/db/user"
	"herbst-server/dbinit"
	"herbst-server/events"
	"herbst-server/services"
	"golang.org/x/crypto/bcrypt"
)

// RegisterCharacterRoutes registers all character-related routes
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
	// Create a new character
	router.POST("/characters", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			Password    string `json:"password"`
			IsNPC       bool   `json:"isNPC"`
			CurrentRoom int    `json:"currentRoomId"`
			StartingRoom int   `json:"startingRoomId"`
			UserID      int    `json:"userId"`
			IsAdmin     bool   `json:"isAdmin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password if provided
		var hashedPassword string
		if req.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			hashedPassword = string(hash)
		}

		// Name validation: 1-23 chars, letters only
		if len(req.Name) < 1 || len(req.Name) > 23 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Character name must be 1-23 characters"})
			return
		}
		for _, ch := range req.Name {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Character name can only contain letters (a-z, A-Z)"})
				return
			}
		}

		// Max 3 characters per user
		if req.UserID > 0 {
			userChars, _ := client.Character.Query().Where(character.HasUserWith(user.IDEQ(req.UserID))).Count(c.Request.Context())
			if userChars >= 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 3 characters per user reached"})
				return
			}
		}

		builder := client.Character.
			Create().
			SetName(req.Name).
			SetIsNPC(req.IsNPC).
			SetIsAdmin(req.IsAdmin)

		if hashedPassword != "" {
			builder.SetPassword(hashedPassword)
		}

		if req.CurrentRoom > 0 {
			builder.SetCurrentRoomId(req.CurrentRoom)
		}
		if req.StartingRoom > 0 {
			builder.SetStartingRoomId(req.StartingRoom)
		}
		if req.UserID > 0 {
			builder.SetUserID(req.UserID)
		}

		character, err := builder.Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Apply race stat modifiers (default to human)
		character, err = dbinit.ApplyRaceToCharacter(c.Request.Context(), client, character)
		if err != nil {
			log.Printf("Warning: failed to apply race modifiers: %v", err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":             character.ID,
			"name":           character.Name,
			"race":           character.Race,
			"gender":         character.Gender,
			"class":          character.Class,
			"strength":       character.Strength,
			"dexterity":      character.Dexterity,
			"constitution":   character.Constitution,
			"intelligence":   character.Intelligence,
			"wisdom":         character.Wisdom,
			"isNPC":          character.IsNPC,
			"is_admin":       character.IsAdmin,
			"currentRoomId":  character.CurrentRoomId,
			"startingRoomId": character.StartingRoomId,
		})
	})

	// Get all characters
	router.GET("/characters", func(c *gin.Context) {
		characters, err := client.Character.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, characters)
	})

	// Get a single character by ID
	router.GET("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		character, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, character)
	})

	// Update a character by ID
	router.PUT("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Name         string  `json:"name"`
			IsNPC        *bool   `json:"isNPC"`
			CurrentRoom  *int    `json:"currentRoomId"`
			StartingRoom *int    `json:"startingRoomId"`
			RespawnRoom  *int    `json:"respawnRoomId"`
			IsAdmin      *bool   `json:"isAdmin"`
			IsTest       *bool   `json:"isTest"`
			Gender       string  `json:"gender"`
			Description  string  `json:"description"`
			LastSeenAt   *string `json:"lastSeenAt"`
			Level        *int    `json:"level"`
			XP           *int    `json:"xp"`
			HP           *int    `json:"hitpoints"`
			MaxHP        *int    `json:"maxHitpoints"`
			Stamina      *int    `json:"stamina"`
			MaxStamina   *int    `json:"maxStamina"`
			Mana         *int    `json:"mana"`
			MaxMana      *int    `json:"maxMana"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Character.UpdateOneID(id)

		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.IsNPC != nil {
			updater.SetIsNPC(*req.IsNPC)
		}
		if req.CurrentRoom != nil {
			updater.SetCurrentRoomId(*req.CurrentRoom)
		}
		if req.StartingRoom != nil {
			updater.SetStartingRoomId(*req.StartingRoom)
		}
		if req.RespawnRoom != nil {
			updater.SetRespawnRoomId(*req.RespawnRoom)
		}
		if req.IsAdmin != nil {
			updater.SetIsAdmin(*req.IsAdmin)
		}
		if req.IsTest != nil {
			updater.SetIsTest(*req.IsTest)
		}
		if req.Gender != "" {
			updater.SetGender(req.Gender)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		if req.LastSeenAt != nil {
			t, err := time.Parse(time.RFC3339, *req.LastSeenAt)
			if err == nil {
				updater.SetLastSeenAt(t)
			}
		}
		if req.Level != nil {
			updater.SetLevel(*req.Level)
		}
		if req.XP != nil {
			updater.SetXp(*req.XP)
		}
		if req.HP != nil {
			updater.SetHitpoints(*req.HP)
		}
		if req.MaxHP != nil {
			updater.SetMaxHitpoints(*req.MaxHP)
		}
		if req.Stamina != nil {
			updater.SetStamina(*req.Stamina)
		}
		if req.MaxStamina != nil {
			updater.SetMaxStamina(*req.MaxStamina)
		}
		if req.Mana != nil {
			updater.SetMana(*req.Mana)
		}
		if req.MaxMana != nil {
			updater.SetMaxMana(*req.MaxMana)
		}

		character, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, character)
	})

	// Delete a character by ID
	router.DELETE("/characters/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		err = client.Character.DeleteOneID(id).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Authenticate a character
	router.POST("/characters/authenticate", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find character by name
		char, err := client.Character.Query().Where(character.NameEQ(req.Name)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character not found"})
			return
		}

		// Check if character has a password set
		if char.Password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Character has no password set"})
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(char.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// Generate access token
		tokenBytes := make([]byte, 32)
		rand.Read(tokenBytes)
		accessToken := hex.EncodeToString(tokenBytes)

		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"access_token": accessToken,
			"character": gin.H{
				"id":       char.ID,
				"name":     char.Name,
				"is_admin": char.IsAdmin,
			},
		})
	})

	// Get characters for a specific user
	router.GET("/user-characters/:id", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Check if user exists
		_, err = client.User.Get(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Get characters for this user
		characters, err := client.Character.Query().
			Where(character.HasUserWith(user.IDEQ(userId))).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return characters without passwords
		result := make([]gin.H, len(characters))
		for i, char := range characters {
			result[i] = gin.H{
				"id":              char.ID,
				"name":            char.Name,
				"isNPC":           char.IsNPC,
				"is_admin":        char.IsAdmin,
				"currentRoomId":   char.CurrentRoomId,
				"startingRoomId": char.StartingRoomId,
				"hitpoints":      char.Hitpoints,
				"max_hitpoints":  char.MaxHitpoints,
				"stamina":        char.Stamina,
				"max_stamina":    char.MaxStamina,
				"mana":           char.Mana,
				"max_mana":       char.MaxMana,
			}
		}

		c.JSON(http.StatusOK, result)
	})

	// Create a character for a specific user
	router.POST("/user-characters/:id", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Check if user exists
		_, err = client.User.Get(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var req struct {
			Name     string `json:"name" binding:"required"`
			Password string `json:"password" binding:"required"`
			Class    string `json:"class"`
			Race     string `json:"race"`
			Gender   string `json:"gender"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password
		hashedPassword, err := services.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Delegate to character service
		charSvc := services.NewCharacterService(client)
		char, err := charSvc.CreateCharacter(c.Request.Context(), services.CreateCharacterInput{
			UserID:   userId,
			Name:     req.Name,
			Password: hashedPassword,
			Class:    req.Class,
			Race:     req.Race,
			Gender:   req.Gender,
		})
		if err != nil {
			switch {
			case errors.Is(err, services.ErrCharacterNameTaken):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			case errors.Is(err, services.ErrTooManyCharacters):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, services.ErrInvalidRace):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race"})
			case errors.Is(err, services.ErrInvalidGender):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":              char.ID,
			"name":            char.Name,
			"isNPC":           char.IsNPC,
			"is_admin":        char.IsAdmin,
			"currentRoomId":   char.CurrentRoomId,
			"startingRoomId":  char.StartingRoomId,
			"hitpoints":       char.Hitpoints,
			"max_hitpoints":   char.MaxHitpoints,
			"stamina":         char.Stamina,
			"max_stamina":     char.MaxStamina,
			"mana":            char.Mana,
			"max_mana":        char.MaxMana,
			"race":            char.Race,
			"class":           char.Class,
			"specialty":       char.Specialty,
			"strength":        char.Strength,
			"dexterity":       char.Dexterity,
			"constitution":    char.Constitution,
			"intelligence":    char.Intelligence,
			"wisdom":          char.Wisdom,
			"level":           char.Level,
			"xp":              char.Xp,
		})
	})

	// Get character class
	router.GET("/characters/:id/class", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    char.ID,
			"name":  char.Name,
			"class": char.Class,
		})
	})

	// Update character class
	router.PUT("/characters/:id/class", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Class string `json:"class" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate class value
		validClasses := map[string]bool{
			"tinkerer":     true,
			"trader":       true,
			"warrior":     true,
			"brawler":     true,
			"mystic":      true,
			"chef":         true,
			"vine_climber": true,
			"survivor":    true,
		}

		if !validClasses[req.Class] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid class. Valid classes: tinkerer, trader, warrior, brawler, mystic, chef, vine_climber, survivor",
			})
			return
		}

		char, err := client.Character.UpdateOneID(id).SetClass(req.Class).Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    char.ID,
			"name":  char.Name,
			"class": char.Class,
		})
	})

	// Get character specialty
	router.GET("/characters/:id/specialty", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        char.ID,
			"name":      char.Name,
			"class":     char.Class,
			"specialty": char.Specialty,
		})
	})

	// Update character specialty
	router.PUT("/characters/:id/specialty", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Specialty string `json:"specialty" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		char, err := client.Character.UpdateOneID(id).SetSpecialty(req.Specialty).Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        char.ID,
			"name":      char.Name,
			"class":     char.Class,
			"specialty": char.Specialty,
		})
	})

	// Get available specialties for a class
	router.GET("/classes/:class/specialties", func(c *gin.Context) {
		class := c.Param("class")

		specialties, ok := constants.ClassSpecialties[class]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"class":       class,
			"specialties": specialties,
		})
	})

	// Check if user needs to create a character
	router.GET("/user-characters/:id/needed", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Check if user exists
		_, err = client.User.Get(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Count characters for this user
		count, err := client.Character.Query().
			Where(character.HasUserWith(user.IDEQ(userId))).
			Count(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"needs_character": count == 0,
			"character_count": count,
		})
	})

	// Get character race
	router.GET("/characters/:id/race", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   char.ID,
			"race": char.Race,
		})
	})

	// GET /races — list all playable races
	router.GET("/races", func(c *gin.Context) {
		races, err := dbinit.GetPlayableRaces(c.Request.Context(), client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(races))
		for i, r := range races {
			result[i] = gin.H{
				"name":           r.Name,
				"display_name":   r.DisplayName,
				"description":    r.Description,
				"stat_modifiers": r.StatModifiers,
				"skill_grants":   r.SkillGrants,
				"equipment_slots":  r.EquipmentSlots,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET /genders — list all genders
	router.GET("/genders", func(c *gin.Context) {
		genders, err := dbinit.GetAllGenders(c.Request.Context(), client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(genders))
		for i, g := range genders {
			result[i] = gin.H{
				"name":               g.Name,
				"display_name":       g.DisplayName,
				"subject_pronoun":    g.SubjectPronoun,
				"object_pronoun":     g.ObjectPronoun,
				"possessive_pronoun": g.PossessivePronoun,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET /game-config/:key
	router.GET("/game-config/:key", func(c *gin.Context) {
		key := c.Param("key")
		cfg, err := client.GameConfig.Query().Where(gameconfig.KeyEQ(key)).Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Config key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": cfg.Key, "value": cfg.Value})
	})

	// PUT /game-config/:key (admin-set fountain room, etc.)
	router.PUT("/game-config/:key", func(c *gin.Context) {
		key := c.Param("key")
		var req struct{ Value string `json:"value" binding:"required"` }
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		existing, err := client.GameConfig.Query().Where(gameconfig.KeyEQ(key)).Only(c.Request.Context())
		if err == nil && existing != nil {
			updated, err := client.GameConfig.UpdateOne(existing).SetValue(req.Value).Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"key": updated.Key, "value": updated.Value})
			return
		}
		created, err := client.GameConfig.Create().SetKey(key).SetValue(req.Value).Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"key": created.Key, "value": created.Value})
	})

	router.PUT("/characters/:id/race", func(c *gin.Context) {
	// Update character race
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Race string `json:"race" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate race against DB
		existingRace, err := client.Race.Query().Where(race.NameEQ(req.Race)).Only(c.Request.Context())
		if err != nil || !existingRace.IsPlayable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-playable race"})
			return
		}

		char, err := client.Character.UpdateOneID(id).SetRace(req.Race).Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   char.ID,
			"name": char.Name,
			"race": char.Race,
		})
	})

	// Get character stats
	router.GET("/characters/:id/stats", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Calculate derived stats based on the spike
		derivedStats := gin.H{
			"max_hp":         char.Constitution*10 + char.Level*10,
			"max_stamina":    char.Constitution*5 + char.Level*5,
			"max_mana":       char.Intelligence*5 + char.Level*5,
			"carry_weight":  char.Strength * 10,
			"dodge_chance":  char.Dexterity,
			"crit_chance":   char.Dexterity * 5 / 10, // 0.5% per DEX
		}

		c.JSON(http.StatusOK, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"strength":      char.Strength,
			"dexterity":     char.Dexterity,
			"constitution":  char.Constitution,
			"intelligence":  char.Intelligence,
			"wisdom":        char.Wisdom,
			"derived":       derivedStats,
		})
	})

	// Update character stats
	router.PUT("/characters/:id/stats", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Strength     *int `json:"strength"`
			Dexterity    *int `json:"dexterity"`
			Constitution *int `json:"constitution"`
			Intelligence *int `json:"intelligence"`
			Wisdom       *int `json:"wisdom"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate stat ranges (1-30)
		validateStat := func(stat *int, name string) error {
			if stat != nil && (*stat < 1 || *stat > 30) {
				return fmt.Errorf("%s must be between 1 and 30", name)
			}
			return nil
		}

		if err := validateStat(req.Strength, "strength"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateStat(req.Dexterity, "dexterity"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateStat(req.Constitution, "constitution"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateStat(req.Intelligence, "intelligence"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateStat(req.Wisdom, "wisdom"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Character.UpdateOneID(id)

		if req.Strength != nil {
			updater.SetStrength(*req.Strength)
		}
		if req.Dexterity != nil {
			updater.SetDexterity(*req.Dexterity)
		}
		if req.Constitution != nil {
			updater.SetConstitution(*req.Constitution)
		}
		if req.Intelligence != nil {
			updater.SetIntelligence(*req.Intelligence)
		}
		if req.Wisdom != nil {
			updater.SetWisdom(*req.Wisdom)
		}

		char, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"strength":      char.Strength,
			"dexterity":     char.Dexterity,
			"constitution":  char.Constitution,
			"intelligence":  char.Intelligence,
			"wisdom":        char.Wisdom,
		})
	})

	// Get character skills (with faction eligibility info)
	router.GET("/characters/:id/skills", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Calculate skill bonuses based on the spike
		calcBonus := func(skill int) string {
			switch {
			case skill >= 91:
				return "+75%"
			case skill >= 76:
				return "+50%"
			case skill >= 51:
				return "+25%"
			case skill >= 26:
				return "+10%"
			default:
				return "+0%"
			}
		}

		// Get eligibility info for faction-based abilities
		abilitySvc := services.NewAbilityEligibilityService(client)
		abilitiesWithElig, err := abilitySvc.AbilitiesForCharacterWithEligibility(c.Request.Context(), id)
		if err != nil {
			log.Printf("[abilities] eligibility check failed for character %d: %v", id, err)
			// Don't fail the entire request — just omit faction abilities eligibility
		}

		// Build faction abilities list with eligibility
		factionAbilities := make([]gin.H, 0)
		if err == nil {
			for _, swe := range abilitiesWithElig {
				sk := swe.Ability
				el := swe.Eligibility
				entry := gin.H{
					"id":             sk.ID,
					"name":           sk.Name,
					"slug":           sk.Slug,
					"ability_type":   sk.AbilityType,
					"ability_class":  sk.AbilityClass,
					"proc_chance":    sk.ProcChance,
					"proc_event":     sk.ProcEvent,
					"cooldown_seconds": sk.CooldownSeconds,
					"mana_cost":      sk.ManaCost,
					"stamina_cost":   sk.StaminaCost,
					"hp_cost":        sk.HpCost,
					"required_tag":  sk.RequiredTag,
					"eligible":       el.Eligible,
					"reason":         el.Reason,
				}
				if sk.Edges.Faction != nil {
					entry["faction_id"] = sk.Edges.Faction.ID
					entry["faction_name"] = sk.Edges.Faction.Name
				}
				factionAbilities = append(factionAbilities, entry)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   char.ID,
			"name": char.Name,
			"skills": gin.H{
				"blades":       gin.H{"level": char.SkillBlades, "bonus": calcBonus(char.SkillBlades)},
				"staves":       gin.H{"level": char.SkillStaves, "bonus": calcBonus(char.SkillStaves)},
				"knives":       gin.H{"level": char.SkillKnives, "bonus": calcBonus(char.SkillKnives)},
				"martial":      gin.H{"level": char.SkillMartial, "bonus": calcBonus(char.SkillMartial)},
				"brawling":     gin.H{"level": char.SkillBrawling, "bonus": calcBonus(char.SkillBrawling)},
				"tech":         gin.H{"level": char.SkillTech, "bonus": calcBonus(char.SkillTech)},
				"light_armor":  gin.H{"level": char.SkillLightArmor, "bonus": calcBonus(char.SkillLightArmor)},
				"cloth_armor":  gin.H{"level": char.SkillClothArmor, "bonus": calcBonus(char.SkillClothArmor)},
				"heavy_armor":  gin.H{"level": char.SkillHeavyArmor, "bonus": calcBonus(char.SkillHeavyArmor)},
			},
			"faction_abilities": factionAbilities,
		})
	})

	// Update character skills
	router.PUT("/characters/:id/skills", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Blades      *int `json:"blades"`
			Staves      *int `json:"staves"`
			Knives      *int `json:"knives"`
			Martial     *int `json:"martial"`
			Brawling    *int `json:"brawling"`
			Tech        *int `json:"tech"`
			LightArmor  *int `json:"light_armor"`
			ClothArmor  *int `json:"cloth_armor"`
			HeavyArmor  *int `json:"heavy_armor"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate skill ranges (0-100)
		validateSkill := func(skill *int, name string) error {
			if skill != nil && (*skill < 0 || *skill > 100) {
				return fmt.Errorf("%s must be between 0 and 100", name)
			}
			return nil
		}

		if err := validateSkill(req.Blades, "blades"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.Staves, "staves"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.Knives, "knives"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.Martial, "martial"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.Brawling, "brawling"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.Tech, "tech"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.LightArmor, "light_armor"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.ClothArmor, "cloth_armor"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validateSkill(req.HeavyArmor, "heavy_armor"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater := client.Character.UpdateOneID(id)

		if req.Blades != nil {
			updater.SetSkillBlades(*req.Blades)
		}
		if req.Staves != nil {
			updater.SetSkillStaves(*req.Staves)
		}
		if req.Knives != nil {
			updater.SetSkillKnives(*req.Knives)
		}
		if req.Martial != nil {
			updater.SetSkillMartial(*req.Martial)
		}
		if req.Brawling != nil {
			updater.SetSkillBrawling(*req.Brawling)
		}
		if req.Tech != nil {
			updater.SetSkillTech(*req.Tech)
		}
		if req.LightArmor != nil {
			updater.SetSkillLightArmor(*req.LightArmor)
		}
		if req.ClothArmor != nil {
			updater.SetSkillClothArmor(*req.ClothArmor)
		}
		if req.HeavyArmor != nil {
			updater.SetSkillHeavyArmor(*req.HeavyArmor)
		}

		char, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":            char.ID,
			"name":          char.Name,
			"blades":        char.SkillBlades,
			"staves":        char.SkillStaves,
			"knives":        char.SkillKnives,
			"martial":       char.SkillMartial,
			"brawling":      char.SkillBrawling,
			"tech":          char.SkillTech,
			"light_armor":  char.SkillLightArmor,
			"cloth_armor":  char.SkillClothArmor,
			"heavy_armor":  char.SkillHeavyArmor,
		})
	})

	// Get NPCs in a specific room
	router.GET("/npcs/room/:id", func(c *gin.Context) {
		roomId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Query for NPCs in this room
		npcs, err := client.Character.Query().
			Where(
				character.IsNPC(true),
				character.CurrentRoomId(roomId),
			).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build response with NPC info
		result := make([]gin.H, len(npcs))
		for i, npc := range npcs {
			// Default XP reward: level * 10, but look up template xp_value if available
			xpValue := npc.Level * 10
			if tmpl, err := client.NPCTemplate.Get(c.Request.Context(), npc.Name); err == nil && tmpl.XpValue > 0 {
				xpValue = tmpl.XpValue
			}
			result[i] = gin.H{
				"id":              npc.ID,
				"name":            npc.Name,
				"isNPC":           npc.IsNPC,
				"currentRoomId":   npc.CurrentRoomId,
				"race":            npc.Race,
				"class":           npc.Class,
				"level":           npc.Level,
				"hitpoints":       npc.Hitpoints,
				"max_hitpoints":   npc.MaxHitpoints,
				"xpValue":         xpValue,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"roomId": roomId,
			"npcs":   result,
			"count":  len(result),
		})
	})

	// Get all NPCs (non-player characters)
	router.GET("/npcs", func(c *gin.Context) {
		npcs, err := client.Character.Query().
			Where(character.IsNPC(true)).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(npcs))
		for i, npc := range npcs {
			result[i] = gin.H{
				"id":              npc.ID,
				"name":            npc.Name,
				"isNPC":           npc.IsNPC,
				"currentRoomId":   npc.CurrentRoomId,
				"race":            npc.Race,
				"class":           npc.Class,
				"level":           npc.Level,
				"hitpoints":       npc.Hitpoints,
				"max_hitpoints":   npc.MaxHitpoints,
				"stamina":         npc.Stamina,
				"max_stamina":     npc.MaxStamina,
				"mana":            npc.Mana,
				"max_mana":        npc.MaxMana,
				"constitution":    npc.Constitution,
				"strength":        npc.Strength,
				"dexterity":        npc.Dexterity,
				"intelligence":   npc.Intelligence,
				"wisdom":          npc.Wisdom,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"npcs":  result,
			"count": len(result),
		})
	})

	// Get character abilities (equipped abilities with slot info)
	router.GET("/characters/:id/abilities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Query abilities for this character
		charAbilities, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slots := make([]map[string]interface{}, 6)
		for i := range slots {
			slots[i] = nil
		}

		for _, ca := range charAbilities {
			abilityName := ""
			abilityDesc := ""
			cooldown := 0
			manaCost := 0
			staminaCost := 0
			effects := make([]gin.H, 0)
			if ca.Edges.Ability != nil {
				abilityName = ca.Edges.Ability.Name
				abilityDesc = ca.Edges.Ability.Description
				cooldown = ca.Edges.Ability.Cooldown
				manaCost = ca.Edges.Ability.ManaCost
				staminaCost = ca.Edges.Ability.StaminaCost
				for _, e := range ca.Edges.Ability.Edges.Effects {
					effects = append(effects, gin.H{
						"effectType":    e.EffectType,
						"damageSubtype": e.DamageSubtype,
						"target":        e.Target,
						"value":         e.Value,
						"duration":      e.Duration,
						"scalingStat":   e.ScalingStat,
						"scalingRatio":  e.ScalingRatio,
						"sortOrder":     e.SortOrder,
					})
				}
			}
			slots[ca.Slot] = map[string]interface{}{
				"slot":           ca.Slot,
				"ability_id":     ca.Edges.Ability.ID,
				"name":           abilityName,
				"description":    abilityDesc,
				"effects":        effects,
				"cooldown":       cooldown,
				"manaCost":       manaCost,
				"staminaCost":    staminaCost,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"character_id": id,
			"slots":        slots,
		})
	})

	// Equip ability to slot
	router.POST("/characters/:id/abilities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			AbilityID int `json:"ability_id"`
			Slot      int `json:"slot"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Slot < 1 || req.Slot > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slot must be between 1 and 5"})
			return
		}

		// Verify ability exists
		abilityObj, err := client.Ability.Get(c.Request.Context(), req.AbilityID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ability not found"})
			return
		}

		// Check if character already has 5 abilities equipped
		existingAbilities, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			All(c.Request.Context())
		if err == nil {
			slotsUsed := make(map[int]bool)
			for _, ca := range existingAbilities {
				if ca.Slot != req.Slot {
					slotsUsed[ca.Slot] = true
				}
			}
			if len(slotsUsed) >= 5 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot equip more than 5 abilities"})
				return
			}
		}

		// Validate ability requirements (skill levels)
		if abilityObj.Requirements != "" {
			requirements := map[string]int{}
			if err := json.Unmarshal([]byte(abilityObj.Requirements), &requirements); err == nil {
				char, err := client.Character.Get(c.Request.Context(), id)
				if err == nil {
					skillLevels := map[string]int{
						"blades":   char.SkillBlades,
						"staves":   char.SkillStaves,
						"knives":   char.SkillKnives,
						"martial":  char.SkillMartial,
						"brawling": char.SkillBrawling,
						"tech":     char.SkillTech,
					}
					for skillName, requiredLevel := range requirements {
						if skillLevels[skillName] < requiredLevel {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot equip ability: requires " + skillName + " level " + strconv.Itoa(requiredLevel)})
							return
						}
					}
				}
			}
		}

		// Remove any existing ability in this slot
		existing, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			Where(characterability.SlotEQ(req.Slot)).
			All(c.Request.Context())

		if err == nil && len(existing) > 0 {
			for _, ca := range existing {
				client.CharacterAbility.DeleteOne(ca).Exec(c.Request.Context())
			}
		}

		// Create new character ability
		charAbility, err := client.CharacterAbility.Create().
			SetCharacterID(id).
			SetAbilityID(req.AbilityID).
			SetSlot(req.Slot).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get the ability name for response
		abilityResp, _ := client.Ability.Get(c.Request.Context(), req.AbilityID)
		abilityName := ""
		if abilityResp != nil {
			abilityName = abilityResp.Name
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"slot":         charAbility.Slot,
			"ability_id":   req.AbilityID,
			"ability_name": abilityName,
		})
	})

	// Unequip ability from slot
	router.DELETE("/characters/:id/abilities/:slot", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		slot, err := strconv.Atoi(c.Param("slot"))
		if err != nil || slot < 1 || slot > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot"})
			return
		}

		// Find and delete the ability in this slot
		charAbilities, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			Where(characterability.SlotEQ(slot)).
			All(c.Request.Context())

		if err != nil || len(charAbilities) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No ability in this slot"})
			return
		}

		for _, ca := range charAbilities {
			client.CharacterAbility.DeleteOne(ca).Exec(c.Request.Context())
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "slot": slot})
	})

	// Swap two ability slots
	router.PUT("/characters/:id/abilities/swap", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Slot1 int `json:"slot1"`
			Slot2 int `json:"slot2"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Slot1 < 1 || req.Slot1 > 5 || req.Slot2 < 1 || req.Slot2 > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slots must be between 1 and 5"})
			return
		}

		// Get abilities in both slots
		getAbilityInSlot := func(slot int) (int, string) {
			cas, err := client.Character.Query().
				Where(character.ID(id)).
				QueryAbilities().
				Where(characterability.SlotEQ(slot)).
				WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
				All(c.Request.Context())

			if err != nil || len(cas) == 0 {
				return 0, ""
			}
			abilityName := ""
			if cas[0].Edges.Ability != nil {
				abilityName = cas[0].Edges.Ability.Name
			}
			return cas[0].Edges.Ability.ID, abilityName
		}

		ability1ID, ability1Name := getAbilityInSlot(req.Slot1)
		ability2ID, ability2Name := getAbilityInSlot(req.Slot2)

		// Clear both slots
		for _, slot := range []int{req.Slot1, req.Slot2} {
			cas, _ := client.Character.Query().
				Where(character.ID(id)).
				QueryAbilities().
				Where(characterability.SlotEQ(slot)).
				All(c.Request.Context())

			for _, ca := range cas {
				client.CharacterAbility.DeleteOne(ca).Exec(c.Request.Context())
			}
		}

		// Swap: assign ability1 to slot2 and ability2 to slot1
		if ability1ID > 0 {
			client.CharacterAbility.Create().
				SetCharacterID(id).
				SetAbilityID(ability1ID).
				SetSlot(req.Slot2).
				Save(c.Request.Context())
		}

		if ability2ID > 0 {
			client.CharacterAbility.Create().
				SetCharacterID(id).
				SetAbilityID(ability2ID).
				SetSlot(req.Slot1).
				Save(c.Request.Context())
		}

		c.JSON(http.StatusOK, gin.H{
			"success":       true,
			"slot1":          req.Slot1,
			"slot2":         req.Slot2,
			"ability1_id":   ability1ID,
			"ability1_name": ability1Name,
			"ability2_id":   ability2ID,
			"ability2_name": ability2Name,
		})
	})

	// Get passive abilities available to a character
	router.GET("/characters/:id/passive-abilities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Verify character exists
		_, err = client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Query passive abilities
		passives, err := client.Ability.Query().
			Where(ability.AbilityClassEQ("passive")).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make([]gin.H, len(passives))
		for i, a := range passives {
			result[i] = gin.H{
				"id":          a.ID,
				"name":        a.Name,
				"description": a.Description,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"character_id":       id,
			"passive_abilities":  result,
			"count":              len(result),
		})
	})

	// Unlock a passive ability for a character
	router.POST("/characters/:id/passive-abilities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Verify character exists
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		var req struct {
			AbilityID int `json:"ability_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify ability exists and is passive
		abilityObj, err := client.Ability.Get(c.Request.Context(), req.AbilityID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ability not found"})
			return
		}

		if abilityObj.AbilityClass != "passive" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ability is not a passive ability"})
			return
		}

		// Check if already equipped
		existing, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			Where(characterability.HasAbilityWith(ability.IDEQ(req.AbilityID))).
			Exist(c.Request.Context())

		if err == nil && existing {
			c.JSON(http.StatusConflict, gin.H{"error": "Ability already equipped for this character"})
			return
		}

		// Find next available slot
		charAbilities, _ := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			All(c.Request.Context())
		usedSlots := make(map[int]bool)
		for _, ca := range charAbilities {
			usedSlots[ca.Slot] = true
		}
		slot := 0
		for s := 1; s <= 5; s++ {
			if !usedSlots[s] {
				slot = s
				break
			}
		}
		if slot == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No available slots"})
			return
		}

		// Create character ability
		charAbility, err := client.CharacterAbility.Create().
			SetCharacterID(char.ID).
			SetAbilityID(req.AbilityID).
			SetSlot(slot).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":      true,
			"id":           charAbility.ID,
			"ability_id":   req.AbilityID,
			"ability_name": abilityObj.Name,
			"slot":         charAbility.Slot,
		})
	})

	// Remove passive ability from character
	router.DELETE("/characters/:id/passive-abilities/:abilityId", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		abilityId, err := strconv.Atoi(c.Param("abilityId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ability ID"})
			return
		}

		// Find and delete the character ability
		charAbility, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			Where(characterability.HasAbilityWith(ability.IDEQ(abilityId))).
			Only(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ability not equipped for this character"})
			return
		}

		err = client.CharacterAbility.DeleteOne(charAbility).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "ability_id": abilityId})
	})

	// Apply damage to a character (used by combat system)
	router.POST("/characters/:id/damage", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Damage     int `json:"damage" binding:"required"`
			AttackerID int `json:"attacker_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Damage < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Damage must be non-negative"})
			return
		}

		// Get the character
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Immortal characters take damage but never die (HP stays at minimum 1)
		if char.IsImmortal {
			newHP := char.Hitpoints - req.Damage
			if newHP < 1 {
				newHP = 1 // Immortal - cannot die
			}

			// Update HP
			updatedChar, err := client.Character.UpdateOneID(id).
				SetHitpoints(newHP).
				Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":        updatedChar.ID,
				"hp":        updatedChar.Hitpoints,
				"maxHp":     updatedChar.MaxHitpoints,
				"defeated":  false,
				"immortal":  true,
				"message":   "Took damage but cannot be killed",
			})
			return
		}

		// Calculate new HP for mortal characters
		newHP := char.Hitpoints - req.Damage
		if newHP < 0 {
			newHP = 0
		}

		// Update character HP (and died_at if NPC is defeated)
		builder := client.Character.UpdateOneID(id).
			SetHitpoints(newHP)
		if newHP == 0 && char.IsNPC {
			builder.SetDiedAt(time.Now())
		}
		updatedChar, err := builder.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log damage contribution for XP split later
		if req.AttackerID > 0 {
			_, err = client.DamageLog.Create().
				SetAttackerID(req.AttackerID).
				SetTargetID(id).
				SetDamage(req.Damage).
				Save(c.Request.Context())
			if err != nil {
				log.Printf("ERROR: failed to log damage: %v", err)
				// non-fatal — combat still proceeds
			}
		}

		defeated := newHP == 0

		// Publish defeat event when an NPC hits 0 HP
		if defeated && updatedChar.IsNPC {
			// Default: derive xp from level * 100
			baseXP := updatedChar.Level * 100
			events.Publish(events.Event{
				Type:    events.EventNPCDefeated,
				Payload: map[string]interface{}{
					"npc_id":    updatedChar.ID,
					"npc_level": updatedChar.Level,
					"base_xp":   baseXP,
				},
				Timestamp: time.Now().UnixMilli(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       updatedChar.ID,
			"hp":       updatedChar.Hitpoints,
			"maxHp":    updatedChar.MaxHitpoints,
			"defeated": defeated,
		})
	})

	// Heal a character (used by combat system for healing)
	router.POST("/characters/:id/heal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Amount int `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Heal amount must be non-negative"})
			return
		}

		// Get the character
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Calculate new HP (don't exceed max)
		newHP := char.Hitpoints + req.Amount
		if newHP > char.MaxHitpoints {
			newHP = char.MaxHitpoints
		}

		// Update character HP
		updatedChar, err := client.Character.UpdateOneID(id).
			SetHitpoints(newHP).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    updatedChar.ID,
			"hp":    updatedChar.Hitpoints,
			"maxHp": updatedChar.MaxHitpoints,
		})
	})

	// Heal all NPCs in a room (used by regen system)
	router.POST("/rooms/:id/npcs/heal", func(c *gin.Context) {
		roomID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var req struct {
			Amount int `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Heal amount must be non-negative"})
			return
		}

		// Get all NPCs in this room that need healing
		npcs, err := client.Character.Query().
			Where(character.IsNPCEQ(true)).
			Where(character.CurrentRoomIdEQ(roomID)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Heal each NPC
		healedCount := 0
		for _, npc := range npcs {
			if npc.Hitpoints >= npc.MaxHitpoints {
				continue // Skip full HP NPCs
			}
			if npc.Hitpoints <= 0 {
				continue // Skip defeated NPCs
			}

			newHP := npc.Hitpoints + req.Amount
			if newHP > npc.MaxHitpoints {
				newHP = npc.MaxHitpoints
			}

			_, err := client.Character.UpdateOneID(npc.ID).
				SetHitpoints(newHP).
				Save(c.Request.Context())
			if err == nil {
				healedCount++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"healed": healedCount,
			"amount": req.Amount,
		})
	})

	// Get classless skills for a character (active abilities not from a faction)
	router.GET("/characters/:id/classless-skills", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Query equipped abilities for this character with effects
		charAbilities, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			WithAbility(func(q *db.AbilityQuery) {
				q.WithEffects()
			}).
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		skills := make([]gin.H, 0)
		for _, ca := range charAbilities {
			ab := ca.Edges.Ability
			if ab == nil || ab.AbilityClass != "active" {
				continue
			}
			effectList := make([]gin.H, 0)
			for _, e := range ab.Edges.Effects {
				effectList = append(effectList, gin.H{
					"effectType":    e.EffectType,
					"damageSubtype": e.DamageSubtype,
					"target":        e.Target,
					"value":         e.Value,
					"duration":      e.Duration,
					"scalingStat":   e.ScalingStat,
					"scalingRatio":  e.ScalingRatio,
					"sortOrder":     e.SortOrder,
				})
			}
			skills = append(skills, gin.H{
				"id":          ab.ID,
				"name":        ab.Name,
				"description": ab.Description,
				"slot":        ca.Slot,
				"cooldown":    ab.Cooldown,
				"manaCost":    ab.ManaCost,
				"staminaCost": ab.StaminaCost,
				"effects":     effectList,
			})
		}

		c.JSON(http.StatusOK, gin.H{"skills": skills})
	})

	// Equip a classless skill
	router.POST("/characters/:id/classless-skills", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			SkillID int `json:"skill_id" binding:"required"`
			Slot    int `json:"slot" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Slot < 1 || req.Slot > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slot must be 1-5"})
			return
		}

		// Verify ability exists
		_, err = client.Ability.Get(c.Request.Context(), req.SkillID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ability not found"})
			return
		}

		// Remove any existing ability in this slot
		existing, err := client.Character.Query().
			Where(character.ID(id)).
			QueryAbilities().
			Where(characterability.SlotEQ(req.Slot)).
			All(c.Request.Context())
		if err == nil && len(existing) > 0 {
			for _, ca := range existing {
				client.CharacterAbility.DeleteOne(ca).Exec(c.Request.Context())
			}
		}

		// Create character ability
		_, err = client.CharacterAbility.Create().
			SetCharacterID(id).
			SetAbilityID(req.SkillID).
			SetSlot(req.Slot).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Skill equipped",
			"skill_id":     req.SkillID,
			"slot":         req.Slot,
			"character_id": id,
		})
	})

	// Swap classless skills between slots
	router.PUT("/characters/:id/classless-skills/swap", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Slot1 int `json:"slot1" binding:"required"`
			Slot2 int `json:"slot2" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Slot1 < 1 || req.Slot1 > 5 || req.Slot2 < 1 || req.Slot2 > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slots must be 1-5"})
			return
		}

		// Get abilities in both slots
		getAbilityInSlot := func(slot int) (int, string) {
			cas, err := client.Character.Query().
				Where(character.ID(id)).
				QueryAbilities().
				Where(characterability.SlotEQ(slot)).
				WithAbility(func(q *db.AbilityQuery) { q.WithEffects() }).
				All(c.Request.Context())
			if err != nil || len(cas) == 0 {
				return 0, ""
			}
			name := ""
			if cas[0].Edges.Ability != nil {
				name = cas[0].Edges.Ability.Name
			}
			return cas[0].Edges.Ability.ID, name
		}

		ability1ID, _ := getAbilityInSlot(req.Slot1)
		ability2ID, _ := getAbilityInSlot(req.Slot2)

		// Clear both slots
		for _, slot := range []int{req.Slot1, req.Slot2} {
			cas, _ := client.Character.Query().
				Where(character.ID(id)).
				QueryAbilities().
				Where(characterability.SlotEQ(slot)).
				All(c.Request.Context())
			for _, ca := range cas {
				client.CharacterAbility.DeleteOne(ca).Exec(c.Request.Context())
			}
		}

		// Swap: assign ability1 to slot2 and ability2 to slot1
		if ability1ID > 0 {
			client.CharacterAbility.Create().
				SetCharacterID(id).
				SetAbilityID(ability1ID).
				SetSlot(req.Slot2).
				Save(c.Request.Context())
		}
		if ability2ID > 0 {
			client.CharacterAbility.Create().
				SetCharacterID(id).
				SetAbilityID(ability2ID).
				SetSlot(req.Slot1).
				Save(c.Request.Context())
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Skills swapped",
			"slot1":        req.Slot1,
			"slot2":        req.Slot2,
			"character_id": id,
		})
	})

	// Get character combat status (for combat system)
	router.GET("/characters/:id/combat-status", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Get the character
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    char.ID,
			"hp":    char.Hitpoints,
			"maxHp": char.MaxHitpoints,
			"isNPC": char.IsNPC,
		})
	})

	// Passive heal NPCs in a room (simulates natural recovery over time)
	router.POST("/rooms/:id/npcs/passive-heal", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Get all NPCs in this room that are hurt
		npcs, err := client.Character.Query().
			Where(character.IsNPCEQ(true)).
			Where(character.CurrentRoomId(id)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Heal each NPC for 25% of max HP (simulating time passing)
		healedCount := 0
		fullyHealedCount := 0
		for _, npc := range npcs {
			if npc.Hitpoints <= 0 {
				// Revive defeated NPCs with 50% HP
				newHP := npc.MaxHitpoints / 2
				_, err := client.Character.UpdateOneID(npc.ID).
					SetHitpoints(newHP).
					Save(c.Request.Context())
				if err == nil {
					healedCount++
				}
				continue
			}

			if npc.Hitpoints >= npc.MaxHitpoints {
				continue // Skip full HP NPCs
			}

			// Heal for 25% of max HP
			healAmount := npc.MaxHitpoints / 4
			if healAmount < 1 {
				healAmount = 1
			}

			newHP := npc.Hitpoints + healAmount
			if newHP >= npc.MaxHitpoints {
				newHP = npc.MaxHitpoints
				fullyHealedCount++
			}

			_, err := client.Character.UpdateOneID(npc.ID).
				SetHitpoints(newHP).
				Save(c.Request.Context())
			if err == nil {
				healedCount++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"healed": healedCount,
			"fullyHealed": fullyHealedCount,
			"room": id,
		})
	})

	// Regenerate stamina for a character
	router.POST("/characters/:id/stamina", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Amount int `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stamina amount must be non-negative"})
			return
		}

		// Get the character
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Calculate new stamina (don't exceed max)
		newStamina := char.Stamina + req.Amount
		if newStamina > char.MaxStamina {
			newStamina = char.MaxStamina
		}

		// Update character stamina
		updatedChar, err := client.Character.UpdateOneID(id).
			SetStamina(newStamina).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      updatedChar.ID,
			"stamina": updatedChar.Stamina,
			"maxStamina": updatedChar.MaxStamina,
		})
	})

	// Regenerate mana for a character
	router.POST("/characters/:id/mana", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Amount int `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mana amount must be non-negative"})
			return
		}

		// Get the character
		char, err := client.Character.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
			return
		}

		// Calculate new mana (don't exceed max)
		newMana := char.Mana + req.Amount
		if newMana > char.MaxMana {
			newMana = char.MaxMana
		}

		// Update character mana
		updatedChar, err := client.Character.UpdateOneID(id).
			SetMana(newMana).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":     updatedChar.ID,
			"mana":   updatedChar.Mana,
			"maxMana": updatedChar.MaxMana,
		})
	})
}
