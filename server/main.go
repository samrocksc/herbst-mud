package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dbinit"
	"herbst-server/routes"
)

// getDBConfig returns database connection config from environment variables
func getDBConfig() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "herbst")
	password := getEnv("DB_PASSWORD", "herbst_password")
	dbname := getEnv("DB_NAME", "herbst_mud")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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

	// Set up Gin router
	router := gin.Default()
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register room routes
	routes.RegisterRoomRoutes(router, client)

	// Register user routes
	routes.RegisterUserRoutes(router, client)

	// Register admin-only routes (requires both auth and admin role)
	routes.RegisterAdminRoutes(router, client)

	// Register character routes
	routes.RegisterCharacterRoutes(router, client)

	// Register equipment routes (GitHub #89 - Item system)
	routes.RegisterEquipmentRoutes(router, client)

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