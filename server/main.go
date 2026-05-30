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
	"herbst-server/db/applog"
	"herbst-server/dblog"
	"herbst-server/dbinit"
	"herbst-server/events"
	"herbst-server/middleware"
	"herbst-server/repository"
	"herbst-server/routes"
	"herbst-server/service"
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

	// Initialize repository and service containers (3-tier architecture)
	repos := repository.NewContainer(client)
	services := service.NewContainer(client, repos, slog.Default())
	_ = services // Will be wired to routes in Phase 2

	// Initialize async log handler (LOGS-002, LOGS-003)
	dbLogHandler := dblog.NewDBHandler(client, nil)
	dbLogHandler.SetBroadcastFunc(routes.BroadcastLogLine)
	minLevel := slog.LevelError
	if isDev := os.Getenv("DATABASE_URL") == ""; isDev {
		minLevel = slog.LevelDebug
	}
	multiHandler := slogmulti{
		stdout: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: minLevel}),
		db:     dbLogHandler,
	}
	multiLogger := slog.New(multiHandler)
	slog.SetDefault(multiLogger)
	dblog.Logger = multiLogger
	defer dbLogHandler.GracefulShutdown()

	// Apply database fixes (converts old data types, sets invincible NPCs, etc.)
	if err := dbinit.ApplyDatabaseFixes(client); err != nil {
		log.Printf("Warning: failed to apply database fixes: %v", err)
	}

	// Initialize default admin user (required for login)
	if err := dbinit.InitAdminUser(client); err != nil {
		log.Printf("Warning: failed to initialize admin user: %v", err)
	}

	// Initialize worlds (creates default "Herbst MUD" world if none exists)
	if err := dbinit.InitWorlds(client); err != nil {
		log.Printf("Warning: failed to initialize worlds: %v", err)
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


	// Set up Gin router
	router := gin.Default()
	
	// Prevent CDN/proxy caching of dynamic API responses (fixes #317:
	// NPC endpoint returns 304 from non-local due to DO App Platform CDN)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Next()
	})

	// Log server errors (5xx) to slog so they appear in admin logs page
	router.Use(func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 500 {
			msg := fmt.Sprintf("%s %s returned %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
			slog.Error(msg,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"status", c.Writer.Status(),
			)
		}
	})

	// CORS middleware - configurable origins for security
	// Empty CORS_ORIGINS = development mode: mirror any origin back
	allowedOrigins := getEnv("CORS_ORIGINS", "")
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		if allowedOrigins == "" {
			// Dev mode: echo back whatever origin the client sends
			allowed = origin != ""
		} else {
			for _, o := range strings.Split(allowedOrigins, ",") {
				if strings.TrimSpace(o) == origin || origin == "" {
					allowed = true
					break
				}
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
	routes.RegisterRoomRoutes(router, client, services)

	// Register ability routes
	routes.RegisterAbilityRoutes(router, repos, client)

	// Register effect routes (ability effects)
	routes.RegisterEffectRoutes(router, repos, client)

	// Register effect definition routes (EFF-003)
	routes.RegisterEffectDefRoutes(router, repos)

	// Register hook routes (EFF-004)
	routes.RegisterHookRoutes(router, repos)

	// Register active effect routes (EFF-005)
	routes.RegisterActiveEffectRoutes(router, repos)

	// Register effect expiry routes
	routes.RegisterExpiryRoutes(router, repos)

	// Register user routes
	routes.RegisterUserRoutes(router, repos)

	// Register character routes (public endpoints)
	routes.RegisterCharacterRoutes(router, services, repos)

	// Protected routes - require authentication
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(client))
	{
		routes.RegisterMeRoutes(protected, services, repos)
	}

	// Register equipment routes (GitHub #89 - Item system)
	routes.RegisterEquipmentRoutes(router, repos, client)

	// Register equipment equip/unequip routes (EQUIP-003)
	routes.RegisterEquipmentEquipRoutes(router, repos)

	// Register backup routes
	routes.RegisterBackupRoutes(router, client)

	// Initialize the event bus.
	events.Init(slog.Default())

	// Register faction routes
	routes.RegisterFactionRoutes(router, repos, client)

	// Register crafting recipe routes (CRAFT-003)
	routes.RegisterCraftingRecipeRoutes(router, repos, client)

	// Register craft endpoint (CRAFT-004)
	routes.RegisterCraftRoutes(router, repos)

	// Register social command routes
	routes.RegisterSocialRoutes(router, client)
	// Register channel config routes
	routes.RegisterChannelRoutes(router, client)

	// Register race routes (RACES-001)
	routes.RegisterRaceRoutes(router, repos, client)

	// Register playable race routes (public endpoint for character creation)
	routes.RegisterPlayableRaceRoutes(router, repos)

	// Register gender routes (GENDERS-001)
	routes.RegisterGenderRoutes(router, repos, client)

	// Register character tag routes
	routes.RegisterCharacterTagRoutes(router, repos)

	// Register tag routes (TAG-001)
	routes.RegisterTagRoutes(router, repos, client)

	// Register achievement routes (ACH-001)
	routes.RegisterAchievementRoutes(router, repos)

	// Register NPC template routes (XP-008)
	// Register trigger routes (TRIGGERS-001)
	routes.RegisterTriggerRoutes(router, repos)
	routes.RegisterNPCTemplateRoutes(router, repos)

	// Register NPC instance routes (NPC-004)
	routes.RegisterNPCInstanceRoutes(router, repos, client)
	// Register item instance routes (NPC-005)
	routes.RegisterItemInstanceRoutes(router, repos, client)
	// Register equipment template routes
	routes.RegisterEquipmentTemplateRoutes(router, repos)

	// Register competency routes (XP-005)
	routes.RegisterCompetencyRoutes(router, repos, client)

	// Register quest routes
	routes.RegisterQuestRoutes(router, services)
	routes.RegisterQuestProgressRoutes(router, repos, services, client)

	// Register chat/messaging routes (RFC-009)
	routes.RegisterChatRoutes(router, services)

	// Register dialog node routes
	routes.RegisterDialogNodeRoutes(router, repos, client)

	// Register event routes (HTTP bridge for game server → event bus)
	routes.RegisterEventRoutes(router, client, slog.Default())

	// Start the respawn ticker (checks dead NPCs every 10s)
	respawnSvc := events.NewRespawnService(client, slog.Default())
	respawnSvc.Start()

	// Register game config routes (protected — admin management)
	routes.RegisterGameConfigRoutes(protected, repos)

	// Register game export/import routes
	routes.RegisterGameExportRoutes(router, client)

	// Register content routes (Week 2: Content Externalization)
	if contentManager != nil {
		routes.RegisterContentRoutes(router, contentManager)
	}

	// Register DB-backed world CRUD routes (always, even without worldManager)
	routes.RegisterWorldCRUDRoutes(router, repos)

	// Register world-scoped content routes (requires worldManager)
	if worldManager != nil {
		routes.RegisterWorldRoutes(router, worldManager, repos)
	}

	// Register admin wipe/reload routes
	routes.RegisterAdminWipeRoutes(router, client)

	// Register log routes (LOGS-004)
	routes.RegisterLogRoutes(router, protected, client)

	// Register debug log routes — SSH client debug messages flow to applogs
	routes.RegisterDebugLogRoutes(protected)

	// Start daily log cleanup goroutine (LOGS-005)
	go startLogCleanup(client)

	// Start corpse cleanup background goroutine (GitHub #22)
	startCorpseCleanup(client)

	// Start regeneration service
	StartRegenService(repos, services, client)

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

	// Inline Swagger UI for /docs (no swagger/ directory needed on disk)
	const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Herbst MUD — API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: '/openapi.json', dom_id: '#swagger-ui', deepLinking: true });
  </script>
</body>
</html>`

	router.GET("/openapi.json", func(c *gin.Context) {
		c.File(staticPath + "/openapi.json")
	})
	router.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUI)
	})

	// Register WebSocket endpoint (Phase 4)
	routes.RegisterWSRoutes(router, repos, client)

	// Start the server
	router.Run("0.0.0.0:8080")
}

// slogmulti fans log records to multiple slog.Handler implementations.
type slogmulti struct {
	stdout slog.Handler
	db    slog.Handler
}

func (h slogmulti) Enabled(ctx context.Context, level slog.Level) bool {
	return h.stdout.Enabled(ctx, level) || h.db.Enabled(ctx, level)
}

func (h slogmulti) Handle(ctx context.Context, r slog.Record) error {
	if err := h.stdout.Handle(ctx, r); err != nil {
		return err
	}
	return h.db.Handle(ctx, r)
}

func (h slogmulti) WithAttrs(attrs []slog.Attr) slog.Handler {
	return slogmulti{
		stdout: h.stdout.WithAttrs(attrs),
		db:    h.db.WithAttrs(attrs),
	}
}

func (h slogmulti) WithGroup(name string) slog.Handler {
	return slogmulti{
		stdout: h.stdout.WithGroup(name),
		db:    h.db.WithGroup(name),
	}
}

// startLogCleanup runs a daily goroutine that prunes applog entries older than
// LOG_RETENTION_DAYS (default 3). Runs once immediately on startup, then every 24h.
func startLogCleanup(client *db.Client) {
	retentionDays := 3
	if v := os.Getenv("LOG_RETENTION_DAYS"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d > 0 {
			retentionDays = d
		}
	}

	runCleanup := func() {
		cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
		count, err := client.AppLog.Delete().Where(applog.CreatedAtLT(cutoff)).Exec(context.Background())
		if err != nil {
			slog.Warn("log cleanup failed", "error", err)
		} else if count > 0 {
			slog.Info("log cleanup complete", "deleted", count, "retention_days", retentionDays)
		}
	}

	runCleanup() // run once on startup
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		runCleanup()
	}
}
