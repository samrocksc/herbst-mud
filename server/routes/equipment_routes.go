package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/equipment"
	"herbst-server/db/room"
	"herbst-server/middleware"
)

// RegisterEquipmentRoutes registers all equipment-related routes
// Public routes:
//   - GET /equipment - Get all equipment (for game display)
//   - GET /equipment/:id - Get equipment by ID
//   - GET /rooms/:id/equipment - Get equipment in a room
//
// Protected routes (auth required):
//   - POST /equipment - Create new equipment
//   - PUT /equipment/:id - Update equipment
//   - DELETE /equipment/:id - Delete equipment
func RegisterEquipmentRoutes(router *gin.Engine, client *db.Client) {
	// === PUBLIC ROUTES ===

	// Get all equipment (public)
	router.GET("/equipment", func(c *gin.Context) {
		items, err := client.Equipment.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	})

	// Get equipment in a room (public - for game client)
	router.GET("/rooms/:id/equipment", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		items, err := client.Equipment.Query().
			Where(equipment.HasRoomWith(room.IDEQ(id))).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	})

	// Get a single equipment item by ID (public)
	router.GET("/equipment/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
			return
		}

		eq, err := client.Equipment.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}

		c.JSON(http.StatusOK, eq)
	})

	// === PROTECTED ROUTES (authentication required) ===

	protected := router.Group("/equipment")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create a new equipment item (protected)
		protected.POST("", func(c *gin.Context) {
			var req struct {
				Name        string `json:"name" binding:"required"`
				Description string `json:"description"`
				Slot        string `json:"slot"`
				Level       int     `json:"level"`
				Weight      int     `json:"weight"`
				IsEquipped  bool    `json:"isEquipped"`
				// Item system fields (GitHub #89)
				IsImmovable bool   `json:"isImmovable"`
				Color       string `json:"color"`
				IsVisible   bool   `json:"isVisible"`
				ItemType    string `json:"itemType"`
				RoomID      int    `json:"roomId"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set defaults
			if req.ItemType == "" {
				req.ItemType = "misc"
			}
			if req.Slot == "" {
				req.Slot = "none"
			}

			builder := client.Equipment.
				Create().
				SetName(req.Name).
				SetDescription(req.Description).
				SetSlot(req.Slot).
				SetLevel(req.Level).
				SetWeight(req.Weight).
				SetIsEquipped(req.IsEquipped).
				SetIsImmovable(req.IsImmovable).
				SetIsVisible(req.IsVisible).
				SetItemType(req.ItemType)

			if req.Color != "" {
				builder.SetColor(req.Color)
			}

			// Set room if provided
			if req.RoomID > 0 {
				builder.SetRoomID(req.RoomID)
			}

			eq, err := builder.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, eq)
		})

		// Update equipment by ID (protected)
		protected.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
				return
			}

			var req struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Slot        string `json:"slot"`
				Level       *int   `json:"level"`
				Weight      *int   `json:"weight"`
				IsEquipped  *bool  `json:"isEquipped"`
				IsImmovable *bool  `json:"isImmovable"`
				Color       string `json:"color"`
				IsVisible   *bool  `json:"isVisible"`
				ItemType    string `json:"itemType"`
				RoomID      *int   `json:"roomId"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			updater := client.Equipment.UpdateOneID(id)

			if req.Name != "" {
				updater.SetName(req.Name)
			}
			if req.Description != "" {
				updater.SetDescription(req.Description)
			}
			if req.Slot != "" {
				updater.SetSlot(req.Slot)
			}
			if req.Level != nil {
				updater.SetLevel(*req.Level)
			}
			if req.Weight != nil {
				updater.SetWeight(*req.Weight)
			}
			if req.IsEquipped != nil {
				updater.SetIsEquipped(*req.IsEquipped)
			}
			if req.IsImmovable != nil {
				updater.SetIsImmovable(*req.IsImmovable)
			}
			if req.Color != "" {
				updater.SetColor(req.Color)
			}
			if req.IsVisible != nil {
				updater.SetIsVisible(*req.IsVisible)
			}
			if req.ItemType != "" {
				updater.SetItemType(req.ItemType)
			}
			if req.RoomID != nil {
				if *req.RoomID == 0 {
					updater.ClearRoom()
				} else {
					updater.SetRoomID(*req.RoomID)
				}
			}

			eq, err := updater.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
				return
			}

			c.JSON(http.StatusOK, eq)
		})

		// Delete equipment by ID (protected)
		protected.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
				return
			}

			err = client.Equipment.DeleteOneID(id).Exec(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
				return
			}

			c.JSON(http.StatusNoContent, nil)
		})
	}
}