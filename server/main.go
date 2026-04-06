package main

import (
	"context"
	"fmt"
	"log"
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

	// Healthz endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
			"db":     "connected",
		})
	})

	// OpenAPI specification endpoint
	router.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, getOpenAPISpec())
	})

	// Start the server
	router.Run("0.0.0.0:8080")
}

func getOpenAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "Herbst MUD API",
			"description": "API for the Herbst MUD game server",
			"version":     "1.0.0",
		},
		"servers": []map[string]interface{}{
			{
				"url": "http://localhost:8080",
			},
		},
		"paths": map[string]interface{}{
			"/healthz": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Health check endpoint",
					"description": "Returns the health status of the server",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"status": map[string]interface{}{
												"type": "string",
											},
											"ssh": map[string]interface{}{
												"type": "string",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/openapi.json": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "OpenAPI specification",
					"description": "Returns the OpenAPI specification for this API",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
									},
								},
							},
						},
					},
				},
			},
			"/rooms": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get all rooms",
					"description": "Returns a list of all rooms in the game",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/Room",
										},
									},
								},
							},
						},
					},
				},
				"post": map[string]interface{}{
					"summary": "Create a new room",
					"description": "Creates a new room with the provided details",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/RoomInput",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Room created successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Room",
									},
								},
							},
						},
					},
				},
			},
			"/rooms/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get a room by ID",
					"description": "Returns a single room by its ID",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Room",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Room not found",
						},
					},
				},
				"put": map[string]interface{}{
					"summary": "Update a room",
					"description": "Updates an existing room with the provided details",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/RoomInput",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Room updated successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/Room",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Room not found",
						},
					},
				},
				"delete": map[string]interface{}{
					"summary": "Delete a room",
					"description": "Deletes a room by its ID",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"responses": map[string]interface{}{
						"204": map[string]interface{}{
							"description": "Room deleted successfully",
						},
						"404": map[string]interface{}{
							"description": "Room not found",
						},
					},
				},
			},
			"/users": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get all users",
					"description": "Returns a list of all users in the game",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/User",
										},
									},
								},
							},
						},
					},
				},
				"post": map[string]interface{}{
					"summary": "Create a new user",
					"description": "Creates a new user with the provided details",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/UserInput",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "User created successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
					},
				},
			},
			"/users/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get a user by ID",
					"description": "Returns a single user by their ID",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Successful response",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "User not found",
						},
					},
				},
				"put": map[string]interface{}{
					"summary": "Update a user",
					"description": "Updates an existing user with the provided details",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/UserInput",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "User updated successfully",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/User",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "User not found",
						},
					},
				},
				"delete": map[string]interface{}{
					"summary": "Delete a user",
					"description": "Deletes a user by their ID",
					"parameters": []map[string]interface{}{
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"responses": map[string]interface{}{
						"204": map[string]interface{}{
							"description": "User deleted successfully",
						},
						"404": map[string]interface{}{
							"description": "User not found",
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"Room": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "integer",
						},
						"name": map[string]interface{}{
							"type": "string",
						},
						"description": map[string]interface{}{
							"type": "string",
						},
						"isStartingRoom": map[string]interface{}{
							"type": "boolean",
						},
						"exits": map[string]interface{}{
							"type": "object",
							"additionalProperties": map[string]interface{}{
								"type": "integer",
							},
						},
					},
				},
				"RoomInput": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"description": map[string]interface{}{
							"type": "string",
						},
						"isStartingRoom": map[string]interface{}{
							"type": "boolean",
						},
						"exits": map[string]interface{}{
							"type": "object",
							"additionalProperties": map[string]interface{}{
								"type": "integer",
							},
						},
					},
				},
				"User": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "integer",
						},
						"email": map[string]interface{}{
							"type": "string",
						},
						"is_admin": map[string]interface{}{
							"type": "boolean",
						},
					},
				},
				"UserInput": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type": "string",
						},
						"password": map[string]interface{}{
							"type": "string",
						},
						"isAdmin": map[string]interface{}{
							"type": "boolean",
						},
					},
				},
			},
		},
	}
}