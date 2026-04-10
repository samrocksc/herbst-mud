package routes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/availabletalent"
	"herbst-server/db/character"
	"herbst-server/db/charactertalent"
	"herbst-server/db/room"
	"herbst-server/db/talent"
	"herbst-server/db/user"
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
			"id":           character.ID,
			"name":         character.Name,
			"isNPC":        character.IsNPC,
			"is_admin":     character.IsAdmin,
			"currentRoomId": character.CurrentRoomId,
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
			Name        string `json:"name"`
			IsNPC       *bool  `json:"isNPC"`
			CurrentRoom *int   `json:"currentRoomId"`
			StartingRoom *int  `json:"startingRoomId"`
			IsAdmin     *bool  `json:"isAdmin"`
			Gender      string `json:"gender"`
			Description string `json:"description"`
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
			Name        string `json:"name" binding:"required"`
			Password    string `json:"password" binding:"required"`
			Class       string `json:"class"`
			Race        string `json:"race"`
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

		// Get class configuration with specialty
		classConfig := constants.GetClassConfig(class, "")

		// Calculate base stats with class bonuses
		baseStrength := constants.DefaultStats.Strength + classConfig.StatBonuses.Strength
		baseDexterity := constants.DefaultStats.Dexterity + classConfig.StatBonuses.Dexterity
		baseConstitution := constants.DefaultStats.Constitution + classConfig.StatBonuses.Constitution
		baseIntelligence := constants.DefaultStats.Intelligence + classConfig.StatBonuses.Intelligence
		baseWisdom := constants.DefaultStats.Wisdom + classConfig.StatBonuses.Wisdom

		// Create the character builder
		builder := client.Character.
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
			SetSpecialty(classConfig.Specialty).
			SetStrength(baseStrength).
			SetDexterity(baseDexterity).
			SetConstitution(baseConstitution).
			SetIntelligence(baseIntelligence).
			SetWisdom(baseWisdom)

		// Apply starting skills
		for skill, level := range classConfig.StartingSkills {
			switch skill {
			case "blades":
				builder.SetSkillBlades(level)
			case "staves":
				builder.SetSkillStaves(level)
			case "knives":
				builder.SetSkillKnives(level)
			case "martial":
				builder.SetSkillMartial(level)
			case "brawling":
				builder.SetSkillBrawling(level)
			case "tech":
				builder.SetSkillTech(level)
			}
		}

		// Create the character
		char, err := builder.Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		// Validate race
		validRaces := map[string]bool{
			"human":         true,
			"mutant":        true,
			"android":       true,
			"escaped_slave": true,
		}
		if !validRaces[req.Race] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race. Valid races: human, mutant, android, escaped_slave"})
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
			"id":            char.ID,
			"name":          char.Name,
			"skills": gin.H{
				"blades":        gin.H{"level": char.SkillBlades, "bonus": calcBonus(char.SkillBlades)},
				"staves":        gin.H{"level": char.SkillStaves, "bonus": calcBonus(char.SkillStaves)},
				"knives":        gin.H{"level": char.SkillKnives, "bonus": calcBonus(char.SkillKnives)},
				"martial":       gin.H{"level": char.SkillMartial, "bonus": calcBonus(char.SkillMartial)},
				"brawling":      gin.H{"level": char.SkillBrawling, "bonus": calcBonus(char.SkillBrawling)},
				"tech":          gin.H{"level": char.SkillTech, "bonus": calcBonus(char.SkillTech)},
				"light_armor":  gin.H{"level": char.SkillLightArmor, "bonus": calcBonus(char.SkillLightArmor)},
				"cloth_armor":  gin.H{"level": char.SkillClothArmor, "bonus": calcBonus(char.SkillClothArmor)},
				"heavy_armor":  gin.H{"level": char.SkillHeavyArmor, "bonus": calcBonus(char.SkillHeavyArmor)},
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

	// Get character talents (equipped talents with slot info)
	// Get character talents (equipped talents with slot info)
	router.GET("/characters/:id/talents", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// Query talents for this character
		charTalents, err := client.Character.Query().
			Where(character.ID(id)).
			QueryTalents().
			WithTalent().
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build slots array (index 0 unused, 1-4 for slots)
		slots := make([]map[string]interface{}, 5)
		for i := range slots {
			slots[i] = nil
		}

		for _, ct := range charTalents {
			talentName := ""
			talentDesc := ""
			effectType := ""
			effectValue := 0
			effectDuration := 0
			cooldown := 0
			manaCost := 0
			staminaCost := 0
			if ct.Edges.Talent != nil {
				talentName = ct.Edges.Talent.Name
				talentDesc = ct.Edges.Talent.Description
				effectType = ct.Edges.Talent.EffectType
				effectValue = ct.Edges.Talent.EffectValue
				effectDuration = ct.Edges.Talent.EffectDuration
				cooldown = ct.Edges.Talent.Cooldown
				manaCost = ct.Edges.Talent.ManaCost
				staminaCost = ct.Edges.Talent.StaminaCost
			}
			slots[ct.Slot] = map[string]interface{}{
				"slot":           ct.Slot,
				"talent_id":      ct.Edges.Talent.ID,
				"name":           talentName,
				"description":    talentDesc,
				"effectType":     effectType,
				"effectValue":    effectValue,
				"effectDuration": effectDuration,
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

	// Equip talent to slot
	router.POST("/characters/:id/talents", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			TalentID int `json:"talent_id"`
			Slot     int `json:"slot"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate slot 1-4
		if req.Slot < 1 || req.Slot > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slot must be between 1 and 4"})
			return
		}

		// Verify talent exists
		talentObj, err := client.Talent.Get(c.Request.Context(), req.TalentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
			return
		}

		// Check if character already has 4 talents equipped (excluding the slot we're replacing)
		existingTalents, err := client.Character.Query().
			Where(character.ID(id)).
			QueryTalents().
			All(c.Request.Context())
		if err == nil {
			// Count unique slots in use (not counting the slot we're about to replace)
			slotsUsed := make(map[int]bool)
			for _, ct := range existingTalents {
				if ct.Slot != req.Slot {
					slotsUsed[ct.Slot] = true
				}
			}
			if len(slotsUsed) >= 4 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot equip more than 4 talents"})
				return
			}
		}

		// Validate talent requirements (skill levels)
		if talentObj.Requirements != "" {
			requirements := map[string]int{}
			if err := json.Unmarshal([]byte(talentObj.Requirements), &requirements); err == nil {
				// Get character's skill columns (blades, staves, etc.)
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
							c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot equip talent: requires " + skillName + " level " + strconv.Itoa(requiredLevel)})
							return
						}
					}
				}
			}
		}
		if talentObj.Requirements != "" {
			requirements := map[string]int{}
			if err := json.Unmarshal([]byte(talentObj.Requirements), &requirements); err == nil {
				// Get character's skills
				charSkills, err := client.Character.Query().
					Where(character.ID(id)).
					QuerySkills().
					All(c.Request.Context())
				if err == nil {
					skillLevels := make(map[string]int)
					for _, cs := range charSkills {
						if cs.Edges.Skill != nil {
						}
					}
					for skillName, requiredLevel := range requirements {
						if skillLevels[skillName] < requiredLevel {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot equip talent: requires " + skillName + " level " + strconv.Itoa(requiredLevel)})
							return
						}
					}
				}
			}
		}

		// Remove any existing talent in this slot
		existing, err := client.Character.Query().
			Where(character.ID(id)).
			QueryTalents().
			Where(charactertalent.SlotEQ(req.Slot)).
			All(c.Request.Context())

		if err == nil && len(existing) > 0 {
			for _, ct := range existing {
				client.CharacterTalent.DeleteOne(ct).Exec(c.Request.Context())
			}
		}

		// Create new character talent
		charTalent, err := client.CharacterTalent.Create().
			SetCharacterID(id).
			SetTalentID(req.TalentID).
			SetSlot(req.Slot).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get the talent name for response
		talent, _ := client.Talent.Get(c.Request.Context(), req.TalentID)
		talentName := ""
		if talent != nil {
			talentName = talent.Name
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":     true,
			"slot":        charTalent.Slot,
			"talent_id":   req.TalentID,
			"talent_name": talentName,
		})
	})

	// Unequip talent from slot
	router.DELETE("/characters/:id/talents/:slot", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		slot, err := strconv.Atoi(c.Param("slot"))
		if err != nil || slot < 1 || slot > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot"})
			return
		}

		// Find and delete the talent in this slot
		charTalents, err := client.Character.Query().
			Where(character.ID(id)).
			QueryTalents().
			Where(charactertalent.SlotEQ(slot)).
			All(c.Request.Context())

		if err != nil || len(charTalents) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No talent in this slot"})
			return
		}

		for _, ct := range charTalents {
			client.CharacterTalent.DeleteOne(ct).Exec(c.Request.Context())
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "slot": slot})
	})

	// Swap two talent slots
	router.PUT("/characters/:id/talents/swap", func(c *gin.Context) {
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

		// Validate slots
		if req.Slot1 < 1 || req.Slot1 > 4 || req.Slot2 < 1 || req.Slot2 > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Slots must be between 1 and 4"})
			return
		}

		// Get talents in both slots
		getTalentInSlot := func(slot int) (int, string) {
			cts, err := client.Character.Query().
				Where(character.ID(id)).
				QueryTalents().
				Where(charactertalent.SlotEQ(slot)).
				WithTalent().
				All(c.Request.Context())

			if err != nil || len(cts) == 0 {
				return 0, ""
			}
			talentName := ""
			if cts[0].Edges.Talent != nil {
				talentName = cts[0].Edges.Talent.Name
			}
			return cts[0].Edges.Talent.ID, talentName
		}

		talent1ID, talent1Name := getTalentInSlot(req.Slot1)
		talent2ID, talent2Name := getTalentInSlot(req.Slot2)

		// Clear both slots
		for _, slot := range []int{req.Slot1, req.Slot2} {
			cts, _ := client.Character.Query().
				Where(character.ID(id)).
				QueryTalents().
				Where(charactertalent.SlotEQ(slot)).
				All(c.Request.Context())

			for _, ct := range cts {
				client.CharacterTalent.DeleteOne(ct).Exec(c.Request.Context())
			}
		}

		// Swap: assign talent1 to slot2 and talent2 to slot1
		if talent1ID > 0 {
			client.CharacterTalent.Create().
				SetCharacterID(id).
				SetTalentID(talent1ID).
				SetSlot(req.Slot2).
				Save(c.Request.Context())
		}

		if talent2ID > 0 {
			client.CharacterTalent.Create().
				SetCharacterID(id).
				SetTalentID(talent2ID).
				SetSlot(req.Slot1).
				Save(c.Request.Context())
		}

		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"slot1":        req.Slot1,
			"slot2":        req.Slot2,
			"talent1_id":   talent1ID,
			"talent1_name": talent1Name,
			"talent2_id":   talent2ID,
			"talent2_name": talent2Name,
		})
	})

	// Get available talents for a character (unlocked talents)
	router.GET("/characters/:id/available-talents", func(c *gin.Context) {
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

		// Query available talents for this character
		availableTalents, err := client.AvailableTalent.Query().
			Where(availabletalent.HasCharacterWith(character.ID(id))).
			WithTalent().
			All(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build response
		result := make([]gin.H, len(availableTalents))
		for i, at := range availableTalents {
			talentName := ""
			talentDesc := ""
			talentID := 0
			if at.Edges.Talent != nil {
				talentID = at.Edges.Talent.ID
				talentName = at.Edges.Talent.Name
				talentDesc = at.Edges.Talent.Description
			}
			result[i] = gin.H{
				"id":              at.ID,
				"talent_id":       talentID,
				"name":            talentName,
				"description":     talentDesc,
				"unlock_reason":   at.UnlockReason,
				"unlocked_at_level": at.UnlockedAtLevel,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"character_id":       id,
			"available_talents":  result,
			"count":              len(result),
		})
	})

	// Add available talent to character (unlock a new talent)
	router.POST("/characters/:id/available-talents", func(c *gin.Context) {
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
			TalentID     int    `json:"talent_id" binding:"required"`
			UnlockReason string `json:"unlock_reason"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate unlock reason
		if req.UnlockReason == "" {
			req.UnlockReason = "manual_unlock"
		}

		// Verify talent exists
		talentObj, err := client.Talent.Get(c.Request.Context(), req.TalentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
			return
		}

		// Check if already available
		existing, err := client.AvailableTalent.Query().
			Where(availabletalent.HasCharacterWith(character.ID(id))).
			Where(availabletalent.HasTalentWith(talent.IDEQ(talentObj.ID))).
			Exist(c.Request.Context())

		if err == nil && existing {
			c.JSON(http.StatusConflict, gin.H{"error": "Talent already available for this character"})
			return
		}

		// Create available talent
		availableTalent, err := client.AvailableTalent.Create().
			SetCharacterID(id).
			SetTalentID(req.TalentID).
			SetUnlockReason(req.UnlockReason).
			SetUnlockedAtLevel(char.Level).
			Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":            true,
			"id":                 availableTalent.ID,
			"talent_id":          req.TalentID,
			"talent_name":        talentObj.Name,
			"unlock_reason":      availableTalent.UnlockReason,
			"unlocked_at_level":  availableTalent.UnlockedAtLevel,
		})
	})

	// Remove available talent from character
	router.DELETE("/characters/:id/available-talents/:talentId", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		talentId, err := strconv.Atoi(c.Param("talentId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talent ID"})
			return
		}

		// Find and delete the available talent
		availableTalent, err := client.AvailableTalent.Query().
			Where(availabletalent.HasCharacterWith(character.ID(id))).
			Where(availabletalent.HasTalentWith(talent.IDEQ(talentId))).
			Only(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Available talent not found"})
			return
		}

		err = client.AvailableTalent.DeleteOne(availableTalent).Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "talent_id": talentId})
	})

	// Apply damage to a character (used by combat system)
	router.POST("/characters/:id/damage", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var req struct {
			Damage int `json:"damage" binding:"required"`
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

		// Update character HP
		updatedChar, err := client.Character.UpdateOneID(id).
			SetHitpoints(newHP).
			Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defeated := newHP == 0

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

	// Get classless skills for a character
	router.GET("/characters/:id/classless-skills", func(c *gin.Context) {
		_, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		// For now, return the 5 default classless skills
		// In the future, this could be stored in the database
		skills := []gin.H{
			{"id": 100, "name": "Concentrate", "description": "Focus your mind to increase accuracy. +WIS to hit for 4 rounds.", "slot": 1, "effectType": "concentrate", "cooldown": 8, "manaCost": 10, "staminaCost": 0, "baseStat": "wisdom", "duration": 4},
			{"id": 101, "name": "Haymaker", "description": "A powerful but reckless strike. +STR damage, -DEX to hit.", "slot": 2, "effectType": "haymaker", "cooldown": 6, "manaCost": 0, "staminaCost": 15, "baseStat": "strength", "duration": 1},
			{"id": 102, "name": "Back-off", "description": "Use agility to dodge all attacks this round. Costs stamina.", "slot": 3, "effectType": "backoff", "cooldown": 10, "manaCost": 0, "staminaCost": 25, "baseStat": "dexterity", "duration": 1},
			{"id": 103, "name": "Scream", "description": "Release a berserker cry. -WIS/INT, +DEX/STR for 2 rounds.", "slot": 4, "effectType": "scream", "cooldown": 12, "manaCost": 5, "staminaCost": 10, "baseStat": "constitution", "duration": 2},
			{"id": 104, "name": "Slap", "description": "A quick stunning strike. DEX vs CON to stun for 1 round.", "slot": 5, "effectType": "slap", "cooldown": 8, "manaCost": 0, "staminaCost": 12, "baseStat": "dexterity", "duration": 1},
		}

		c.JSON(http.StatusOK, gin.H{"skills": skills})
	})

	// Equip a classless skill (mock - just returns success)
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

		// For now, just return success (storage would require DB schema update)
		c.JSON(http.StatusOK, gin.H{
			"message": "Skill equipped",
			"skill_id": req.SkillID,
			"slot": req.Slot,
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

		// For now, just return success
		c.JSON(http.StatusOK, gin.H{
			"message": "Skills swapped",
			"slot1": req.Slot1,
			"slot2": req.Slot2,
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
