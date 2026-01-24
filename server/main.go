package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up Gin router
	router := gin.Default()

	// Healthz endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"ssh":    "running",
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
		},
	}
}