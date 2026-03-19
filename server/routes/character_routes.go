package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/constants"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/room"
	"herbst-server/db/user"
	"herbst-server/middleware"
	"golang.org/x/crypto/bcrypt"
)

// RegisterCharacterRoutes registers all character-related routes
// Public routes:
//   - POST /characters/authenticate - Character login
//
// Protected routes (auth required):
//   - All other character CRUD operations
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
	// === PUBLIC ROUTES ===

	// Authenticate a character (public - for login)
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

		// Generate access token using JWT
		token, err := middleware.GenerateToken(char.ID, char.Name, char.IsAdmin, "character")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"access_token":   token,
			"character": gin.H{
				"id":       char.ID,
				"name":     char.Name,
				"is_admin": char.IsAdmin,
			},
		})
	})

	// Check if user needs to create a character (public)
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
			"needs_character":  count == 0,
			"character_count": count,
		})
	})

	// Get available specialties for a class (public)
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

	// === PROTECTED ROUTES (authentication required) ===

	protected := router.Group("/characters")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create a new character (protected)
		protected.POST("", func(c *gin.Context) {
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

		// Get all characters (protected)
		protected.GET("", func(c *gin.Context) {
			characters, err := client.Character.Query().All(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, characters)
		})

		// Get a single character by ID (protected)
		protected.GET("/:id", func(c *gin.Context) {
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

			c.JSON(http.StatusOK, char)
		})

		// Update a character by ID (protected)
		protected.PUT("/:id", func(c *gin.Context) {
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

			char, err := updater.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
				return
			}

			c.JSON(http.StatusOK, char)
		})

		// Delete a character by ID (protected)
		protected.DELETE("/:id", func(c *gin.Context) {
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

		// Get characters for a specific user (protected)
		protected.GET("/user/:id", func(c *gin.Context) {
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

		// Create a character for a specific user (protected)
		protected.POST("/user/:id", func(c *gin.Context) {
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

		// Get character class (protected)
		protected.GET("/:id/class", func(c *gin.Context) {
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

		// Update character class (protected)
		protected.PUT("/:id/class", func(c *gin.Context) {
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
				"tinkerer":      true,
				"trader":        true,
				"warrior":       true,
				"brawler":       true,
				"mystic":        true,
				"chef":          true,
				"vine_climber": true,
				"survivor":      true,
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

		// Get character specialty (protected)
		protected.GET("/:id/specialty", func(c *gin.Context) {
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

		// Update character specialty (protected)
		protected.PUT("/:id/specialty", func(c *gin.Context) {
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

		// Get character race (protected)
		protected.GET("/:id/race", func(c *gin.Context) {
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

		// Update character race (protected)
		protected.PUT("/:id/race", func(c *gin.Context) {
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
				"human":          true,
				"mutant":         true,
				"android":        true,
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

		// Get character stats (protected)
		protected.GET("/:id/stats", func(c *gin.Context) {
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

		// Update character stats (protected)
		protected.PUT("/:id/stats", func(c *gin.Context) {
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

		// Get character skills (protected)
		protected.GET("/:id/skills", func(c *gin.Context) {
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
				"id":       char.ID,
				"name":     char.Name,
				"skills": gin.H{
					"blades":       gin.H{"level": char.SkillBlades, "bonus": calcBonus(char.SkillBlades)},
					"staves":       gin.H{"level": char.SkillStaves, "bonus": calcBonus(char.SkillStaves)},
					"knives":       gin.H{"level": char.SkillKnives, "bonus": calcBonus(char.SkillKnives)},
					"martial":      gin.H{"level": char.SkillMartial, "bonus": calcBonus(char.SkillMartial)},
					"brawling":     gin.H{"level": char.SkillBrawling, "bonus": calcBonus(char.SkillBrawling)},
					"tech":         gin.H{"level": char.SkillTech, "bonus": calcBonus(char.SkillTech)},
					"light_armor": gin.H{"level": char.SkillLightArmor, "bonus": calcBonus(char.SkillLightArmor)},
					"cloth_armor": gin.H{"level": char.SkillClothArmor, "bonus": calcBonus(char.SkillClothArmor)},
					"heavy_armor": gin.H{"level": char.SkillHeavyArmor, "bonus": calcBonus(char.SkillHeavyArmor)},
				},
			})
		})

		// Update character skills (protected)
		protected.PUT("/:id/skills", func(c *gin.Context) {
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
	}
}