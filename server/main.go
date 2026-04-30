package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	_ "github.com/lib/pq"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"herbst-server/content"
	"herbst-server/db"
	"herbst-server/dbinit"
	"herbst-server/events"
	"herbst-server/middleware"
	"herbst-server/routes"
)

// getDBConfig returns database connection config from environment variables
// Supports Neon DB via DATABASE_URL or individual variables
func getDBConfig() string {
	// Neon DB and many managed Postgres providers set DATABASE_URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	// Build from individual env vars (required for production)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check if this is development mode (no env vars set)
	isDev := host == "" && user == ""
	
	// For development only - use defaults if all env vars are empty
	if isDev {
		log.Println("Warning: Using development database defaults. Set DATABASE_URL or DB_* env vars for production.")
		host = "localhost"
		port = "5432"
		user = "herbst"
		password = "herbst_password"
		dbname = "herbst_mud"
	}

	// SSL mode: Neon requires 'require', but local dev can use 'disable'
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		if isDev {
			sslMode = "disable" // Local dev default
		} else {
			sslMode = "require" // Production default (Neon compatible)
		}
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)
}

// getEnv returns environment variable or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize database
	client, err := db.Open("postgres", getDBConfig())
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()

	// Run auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Database initialized successfully")

	// Apply database fixes (converts old data types, sets invincible NPCs, etc.)
	if err := dbinit.ApplyDatabaseFixes(client); err != nil {
		log.Printf("Warning: failed to apply database fixes: %v", err)
	}

	// Initialize cross-shaped rooms
	if err := dbinit.InitCrossWay(client); err != nil {
		log.Printf("Warning: failed to initialize cross-shaped rooms: %v", err)
	}

	// Initialize default admin user
	if err := dbinit.InitAdminUser(client); err != nil {
		log.Printf("Warning: failed to initialize admin user: %v", err)
	}

	// Initialize characters (test characters + Gandalf)
	if err := dbinit.InitCharacters(client); err != nil {
		log.Printf("Warning: failed to initialize characters: %v", err)
	}

	// Initialize fountain for new character creation flow
	if err := dbinit.InitFountain(client); err != nil {
		log.Printf("Warning: failed to initialize fountain: %v", err)
	}

	// Initialize Gizmo NPC
	if err := dbinit.InitGizmoNPC(client); err != nil {
		log.Printf("Warning: failed to initialize Gizmo NPC: %v", err)
	}

	// Seed races and genders
	if err := dbinit.InitRaces(client); err != nil {
		log.Printf("Warning: failed to initialize races: %v", err)
	}
	if err := dbinit.InitGenders(client); err != nil {
		log.Printf("Warning: failed to initialize genders: %v", err)
	}

	// Initialize the Junkyard newbie zone
	if err := dbinit.InitJunkyard(client); err != nil {
		log.Printf("Warning: failed to initialize Junkyard: %v", err)
	}

	// Heal all characters with invalid HP (startup fix)
	if err := dbinit.InitCharacterHealth(client); err != nil {
		log.Printf("Warning: failed to initialize character health: %v", err)
	}

	// Apply database fixes again (after InitCharacterHealth which might restore 0 HP chars)
	if err := dbinit.EnsureCombatDummyImmortal(client); err != nil {
		log.Printf("Warning: failed to ensure Combat Dummy immortality: %v", err)
	}

	// Initialize world manager (Week 8: Multi-World Support)
	worldManager, err := content.NewWorldManager("/home/sam/GitHub/herbst-mud/content")
	if err != nil {
		log.Printf("Warning: failed to load world manager: %v", err)
	} else {
		log.Printf("Worlds loaded: %d", len(worldManager.GetAllWorlds()))
		for _, world := range worldManager.GetAllWorlds() {
			log.Printf("  - %s (%s)", world.Name, world.Status)
			if stats, ok := worldManager.GetWorldStats(world.ID); ok {
				log.Printf("    Content: %d skills, %d items, %d npcs, %d rooms, %d quests",
					stats.Skills, stats.Items, stats.NPCs, stats.Rooms, stats.Quests)
			}
		}
	}

	// Keep legacy content manager for backward compatibility (default world)
	var contentManager *content.Manager
	if worldManager != nil {
		defaultWorld := worldManager.GetDefaultWorld()
		contentManager, _ = worldManager.GetWorldManager(defaultWorld)
		log.Printf("Default world: %s", defaultWorld)
	}

	// Initialize consumables (health potions, etc.)
	if err := dbinit.InitConsumables(client); err != nil {
		log.Printf("Warning: failed to initialize consumables: %v", err)
	}

	// Give starting characters health potions
	if err := dbinit.GivePotionToCharacter(client, 9); err != nil { // sma
		log.Printf("Warning: failed to give potion to character: %v", err)
	}

	// Set up Gin router
	router := gin.Default()
	
	// CORS middleware - configurable origins for security
	allowedOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, o := range strings.Split(allowedOrigins, ",") {
			if strings.TrimSpace(o) == origin || origin == "" {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Rate limiting middleware - prevents DoS/brute force
	rate := getEnv("RATE_LIMIT", "100") // requests per minute
	window := getEnv("RATE_WINDOW", "60") // seconds
	rateInt, _ := strconv.Atoi(rate)
	if rateInt == 0 {
		rateInt = 100
	}
	windowInt, _ := strconv.Atoi(window)
	if windowInt == 0 {
		windowInt = 60
	}
	
	limiterStore := memory.NewStore()
	limiterRate := limiter.Rate{
		Period: time.Duration(windowInt) * time.Second,
		Limit:  int64(rateInt),
	}
	rateLimiter := limiter.New(limiterStore, limiterRate)
	router.Use(func(c *gin.Context) {
		context, err := rateLimiter.Get(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if context.Reached {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	})

	// Register room routes
	routes.RegisterRoomRoutes(router, client)

	// Register skill routes
	routes.RegisterSkillRoutes(router, client)

	// Register talent routes
	routes.RegisterTalentRoutes(router, client)

	// Register user routes
	routes.RegisterUserRoutes(router, client)

	// Register character routes (public endpoints)
	routes.RegisterCharacterRoutes(router, client)

	// Protected routes - require authentication
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Add protected character routes here
		// Example: protected.PUT("/characters/:id", ...)
	}

	// Register equipment routes (GitHub #89 - Item system)
	routes.RegisterEquipmentRoutes(router, client)

	// Register backup routes
	routes.RegisterBackupRoutes(router, client)

	// Initialize the event bus.
	events.Init(slog.Default())

	// Register faction routes
	routes.RegisterFactionRoutes(router, client)

	// Register NPC template routes (XP-008)
	routes.RegisterNPCTemplateRoutes(router, client)

	// Register competency routes (XP-005)
	routes.RegisterCompetencyRoutes(router, client)

	// Register event routes (HTTP bridge for game server → event bus)
	routes.RegisterEventRoutes(router, client, slog.Default())

	// Register game config routes (protected — admin management)
	routes.RegisterGameConfigRoutes(protected, client)

	// Register game export/import routes
	routes.RegisterGameExportRoutes(router, client)

	// Register content routes (Week 2: Content Externalization)
	if contentManager != nil {
		routes.RegisterContentRoutes(router, contentManager)
	}
	if worldManager != nil {
		routes.RegisterWorldRoutes(router, worldManager)
	}

	// Register admin wipe/reload routes
	routes.RegisterAdminWipeRoutes(router, client)

	// Start corpse cleanup background goroutine (GitHub #22)
	startCorpseCleanup(client)

	// Healthz endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
			"db":     "connected",
		})
	})

	// OpenAPI spec — served from static file (generated; see tools/openapi-gen/)
	// Binary runs from repo root, so use ./server/static/ relative to that
	staticPath := "./server/static"
	router.GET("/openapi.json", func(c *gin.Context) {
		c.File(staticPath + "/openapi.json")
	})
	router.GET("/docs", func(c *gin.Context) {
		c.File(staticPath + "/swagger/index.html")
	})

	// Start the server
	router.Run("0.0.0.0:8080")
}
