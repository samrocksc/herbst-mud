package main

import (
	"context"
	"log"
	"net/http"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dbinit"
	"herbst-server/routes"
)

func main() {
	// Initialize database
	client, err := db.Open("postgres", "host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable")
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

	// Set up Gin router
	router := gin.Default()

	// Register room routes
	routes.RegisterRoomRoutes(router, client)

	// Register user routes
	routes.RegisterUserRoutes(router, client)

	// Register character routes
	routes.RegisterCharacterRoutes(router, client)

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
	router.Run(":8080")
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
			},
		},
	}
}