package routes

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/room"
	"herbst-server/db/user"
)

// RegisterCharacterRoutes registers all character-related routes
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
	// Create a new character
	router.POST("/characters", func(c *gin.Context) {
		var req struct {
			Name         string `json:"name" binding:"required"`
			Password     string `json:"password"`
			IsNPC        bool   `json:"isNPC"`
			CurrentRoom  int    `json:"currentRoomId"`
			StartingRoom int    `json:"startingRoomId"`
			UserID       int    `json:"userId"`
			IsAdmin      bool   `json:"isAdmin"`
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

		c.JSON(http.StatusCreated, gin.H{
			"id":             character.ID,
			"name":           character.Name,
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
			Name         string `json:"name"`
			IsNPC        *bool  `json:"isNPC"`
			CurrentRoom  *int   `json:"currentRoomId"`
			StartingRoom *int   `json:"startingRoomId"`
			IsAdmin      *bool  `json:"isAdmin"`
			Gender       string `json:"gender"`
			Description  string `json:"description"`
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
		if req.IsAdmin != nil {
			updater.SetIsAdmin(*req.IsAdmin)
		}
		if req.Gender != "" {
			updater.SetGender(req.Gender)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
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
			"access_token":  accessToken,
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
				"id":             char.ID,
				"name":           char.Name,
				"isNPC":          char.IsNPC,
				"is_admin":       char.IsAdmin,
				"currentRoomId":  char.CurrentRoomId,
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
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if character name already exists
		existingChar, err := client.Character.Query().Where(character.NameEQ(req.Name)).Only(c.Request.Context())
		if err == nil && existingChar != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Character name already exists"})
			return
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPassword := string(hash)

		// Set default stats based on class (if provided)
		hitpoints := 100
		maxHitpoints := 100
		stamina := 50
		maxStamina := 50
		mana := 50
		maxMana := 50

		// Get starting room (first room by default)
		startingRooms, err := client.Room.Query().Where(room.IsStartingRoom(true)).All(c.Request.Context())
		if err != nil {
			log.Printf("Warning: failed to get starting room: %v", err)
		}
		var startingRoomID int
		if err == nil && len(startingRooms) > 0 {
			startingRoomID = startingRooms[0].ID
		}

		// Set race (default to human)
		race := "human"
		if req.Race != "" {
			race = req.Race
		}

		// Set class (default to survivor)
		class := "survivor"
		if req.Class != "" {
			class = req.Class
		}

		// Create the character
		char, err := client.Character.
			Create().
			SetName(req.Name).
			SetPassword(hashedPassword).
			SetUserID(userId).
			SetHitpoints(hitpoints).
			SetMaxHitpoints(maxHitpoints).
			SetStamina(stamina).
			SetMaxStamina(maxStamina).
			SetMana(mana).
			SetMaxMana(maxMana).
			SetCurrentRoomId(startingRoomID).
			SetStartingRoomId(startingRoomID).
			SetRace(race).
			SetClass(class).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":             char.ID,
			"name":           char.Name,
			"isNPC":          char.IsNPC,
			"is_admin":       char.IsAdmin,
			"currentRoomId":  char.CurrentRoomId,
			"startingRoomId": char.StartingRoomId,
			"hitpoints":      char.Hitpoints,
			"max_hitpoints":  char.MaxHitpoints,
			"stamina":        char.Stamina,
			"max_stamina":    char.MaxStamina,
			"mana":           char.Mana,
			"max_mana":       char.MaxMana,
			"race":           char.Race,
			"class":          char.Class,
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

		// Validate class using constants
		classValid := false
		for _, c := range constants.ValidClasses {
			if c == req.Class {
				classValid = true
				break
			}
		}
		if !classValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid class. Valid classes: %v", constants.ValidClasses),
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

	// Update character race
	router.PUT("/characters/:id/race", func(c *gin.Context) {
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

		// Validate race using constants
		raceValid := false
		for _, r := range constants.ValidRaces {
			if r == req.Race {
				raceValid = true
				break
			}
		}
		if !raceValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid race. Valid races: %v", constants.ValidRaces)})
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
			"max_hp":       char.Constitution*10 + char.Level*10,
			"max_stamina":  char.Constitution*5 + char.Level*5,
			"max_mana":     char.Intelligence*5 + char.Level*5,
			"carry_weight": char.Strength * 10,
			"dodge_chance": char.Dexterity,
			"crit_chance":  char.Dexterity * 5 / 10, // 0.5% per DEX
		}

		c.JSON(http.StatusOK, gin.H{
			"id":           char.ID,
			"name":         char.Name,
			"strength":     char.Strength,
			"dexterity":    char.Dexterity,
			"constitution": char.Constitution,
			"intelligence": char.Intelligence,
			"wisdom":       char.Wisdom,
			"derived":      derivedStats,
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
			"id":           char.ID,
			"name":         char.Name,
			"strength":     char.Strength,
			"dexterity":    char.Dexterity,
			"constitution": char.Constitution,
			"intelligence": char.Intelligence,
			"wisdom":       char.Wisdom,
		})
	})

	// Get character skills
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

		c.JSON(http.StatusOK, gin.H{
			"id":   char.ID,
			"name": char.Name,
			"skills": gin.H{
				"blades":      gin.H{"level": char.SkillBlades, "bonus": calcBonus(char.SkillBlades)},
				"staves":      gin.H{"level": char.SkillStaves, "bonus": calcBonus(char.SkillStaves)},
				"knives":      gin.H{"level": char.SkillKnives, "bonus": calcBonus(char.SkillKnives)},
				"martial":     gin.H{"level": char.SkillMartial, "bonus": calcBonus(char.SkillMartial)},
				"brawling":    gin.H{"level": char.SkillBrawling, "bonus": calcBonus(char.SkillBrawling)},
				"tech":        gin.H{"level": char.SkillTech, "bonus": calcBonus(char.SkillTech)},
				"light_armor": gin.H{"level": char.SkillLightArmor, "bonus": calcBonus(char.SkillLightArmor)},
				"cloth_armor": gin.H{"level": char.SkillClothArmor, "bonus": calcBonus(char.SkillClothArmor)},
				"heavy_armor": gin.H{"level": char.SkillHeavyArmor, "bonus": calcBonus(char.SkillHeavyArmor)},
			},
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
			Blades     *int `json:"blades"`
			Staves     *int `json:"staves"`
			Knives     *int `json:"knives"`
			Martial    *int `json:"martial"`
			Brawling   *int `json:"brawling"`
			Tech       *int `json:"tech"`
			LightArmor *int `json:"light_armor"`
			ClothArmor *int `json:"cloth_armor"`
			HeavyArmor *int `json:"heavy_armor"`
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
			"id":          char.ID,
			"name":        char.Name,
			"blades":      char.SkillBlades,
			"staves":      char.SkillStaves,
			"knives":      char.SkillKnives,
			"martial":     char.SkillMartial,
			"brawling":    char.SkillBrawling,
			"tech":        char.SkillTech,
			"light_armor": char.SkillLightArmor,
			"cloth_armor": char.SkillClothArmor,
			"heavy_armor": char.SkillHeavyArmor,
		})
	})

	// Get available talents
	router.GET("/talents", func(c *gin.Context) {
		talents, err := client.Talent.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var result []gin.H
		for _, t := range talents {
			result = append(result, gin.H{
				"id":           t.ID,
				"name":         t.Name,
				"description":  t.Description,
				"requirements": t.Requirements,
			})
		}
		c.JSON(http.StatusOK, result)
	})

	// Get character's equipped talents
	router.GET("/characters/:id/talents", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		charTalentEdges, err := client.CharacterTalent.Query().
			Where(charactertalent.HasCharacterWith(character.IDEQ(id))).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Slots 1-4
		slots := [5]*gin.H{} // index 1-4
		for _, ct := range charTalentEdges {
			if ct.Slot >= 1 && ct.Slot <= 4 {
				talent, _ := client.Talent.Get(c.Request.Context(), ct.TalentID)
				if talent != nil {
					slots[ct.Slot] = gin.H{
						"slot":        ct.Slot,
						"talent_id":   talent.ID,
						"name":        talent.Name,
						"description": talent.Description,
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"character_id": id,
			"slots":        slots,
		})
	})

	// Equip/unequip talents
	router.PUT("/characters/:id/talents", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Slot     int `json:"slot"`
			TalentID int `json:"talent_id"` // 0 to unequip
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Slot < 1 || req.Slot > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slot must be between 1 and 4"})
			return
		}

		ctx := c.Request.Context()

		// Remove existing talent in this slot
		existing, err := client.CharacterTalent.Query().
			Where(
				charactertalent.HasCharacterWith(character.IDEQ(id)),
				charactertalent.SlotEQ(req.Slot),
			).
			All(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, ct := range existing {
			client.CharacterTalent.DeleteOne(ct).Exec(ctx)
		}

		// If talent_id > 0, equip new talent
		if req.TalentID > 0 {
			// Verify talent exists
			_, err := client.Talent.Get(ctx, req.TalentID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
				return
			}

			// Verify character exists
			char, err := client.Character.Get(ctx, id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
				return
			}

			// Check if character has this talent available
			available, err := client.AvailableTalent.Query().
				Where(
					availabletalent.HasCharacterWith(character.IDEQ(id)),
					availabletalent.HasTalentWith(talent.IDEQ(req.TalentID)),
				).
				Exist(ctx)
			if err != nil || !available {
				// Auto-grant talent if not available (for testing)
				client.AvailableTalent.Create().
					SetCharacterID(char.ID).
					SetTalentID(req.TalentID).
					Save(ctx)
			}

			// Create new character_talent
			client.CharacterTalent.Create().
				SetCharacterID(id).
				SetTalentID(req.TalentID).
				SetSlot(req.Slot).
				Save(ctx)
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "slot": req.Slot})
	})
}
