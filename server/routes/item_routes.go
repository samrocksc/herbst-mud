package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// RegisterItemRoutes registers all item/equipment-related routes
func RegisterItemRoutes(router *gin.Engine, client *db.Client) {
	group := router.Group("/items")
	{
		// Get all items
		group.GET("", func(c *gin.Context) {
			items, err := client.Equipment.Query().All(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, items)
		})

		// Get item by ID
		group.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
				return
			}

			item, err := client.Equipment.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}

			c.JSON(http.StatusOK, item)
		})

		// Get item examine details (with hidden details revealed based on skill)
		group.GET("/:id/examine", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
				return
			}

			item, err := client.Equipment.Get(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}

			// Get character's examine skill level (default 1 if not specified)
			examineLevel := 1
			if charIDStr := c.Query("character_id"); charIDStr != "" {
				charID, err := strconv.Atoi(charIDStr)
				if err == nil {
					char, err := client.Character.Get(c.Request.Context(), charID)
					if err == nil {
						// For now, derive examine level from INT stat or default
						// In full implementation, this would query character skills
						examineLevel = 1 + (char.IntStat / 5) // Rough approximation
					}
				}
			}

			// Process hidden details based on examine level
			hiddenDetails := processHiddenDetails(item.HiddenDetails, examineLevel)

			// Check readable content
			var readableContent string
			var canRead bool = true
			if item.IsReadable && item.Content != "" {
				if item.ReadSkill != "" && item.ReadSkillLevel > 0 {
					// Check if character has required skill
					charLevel := examineLevel
					if charLevel < item.ReadSkillLevel {
						canRead = false
						readableContent = "[Requires " + item.ReadSkill + " skill level " + 
							strconv.Itoa(item.ReadSkillLevel) + " to decode]"
					}
				}
				if canRead {
					readableContent = item.Content
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"id":            item.ID,
				"name":          item.Name,
				"description":   item.Description,
				"examineDesc":   item.ExamineDesc,
				"hiddenDetails": hiddenDetails,
				"isReadable":    item.IsReadable,
				"readContent":   readableContent,
				"examineLevel":  examineLevel,
				"examineXP":     1, // XP awarded for examining
				"type":          item.Slot,
				"weight":        item.Weight,
				"level":         item.Level,
				"isImmovable":   item.IsImmovable,
				"isContainer":   item.IsContainer,
			})
		})

		// Create new item
		group.POST("", func(c *gin.Context) {
			var req struct {
				Name         string `json:"name" binding:"required"`
				Description  string `json:"description"`
				ShortDesc    string `json:"shortDesc"`
				ExamineDesc  string `json:"examineDesc"`
				Slot         string `json:"slot"`
				Level        int    `json:"level"`
				Weight       int    `json:"weight"`
				IsImmovable  bool   `json:"isImmovable"`
				IsContainer  bool   `json:"isContainer"`
				IsReadable   bool   `json:"isReadable"`
				Content      string `json:"content"`
				ReadSkill    string `json:"readSkill"`
				ReadSkillLevel int  `json:"readSkillLevel"`
				RoomID       *int   `json:"roomId"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			builder := client.Equipment.Create().
				SetName(req.Name).
				SetDescription(req.Description).
				SetSlot(req.Slot).
				SetLevel(req.Level).
				SetWeight(req.Weight).
				SetIsImmovable(req.IsImmovable).
				SetIsContainer(req.IsContainer).
				SetIsReadable(req.IsReadable).
				SetContent(req.Content).
				SetReadSkill(req.ReadSkill).
				SetReadSkillLevel(req.ReadSkillLevel)

			if req.ShortDesc != "" {
				builder.SetShortDesc(req.ShortDesc)
			}
			if req.ExamineDesc != "" {
				builder.SetExamineDesc(req.ExamineDesc)
			}
			if req.RoomID != nil {
				builder.SetRoomID(*req.RoomID)
			}

			item, err := builder.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, item)
		})

		// Update item
		group.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
				return
			}

			var req struct {
				Name         string `json:"name"`
				Description  string `json:"description"`
				ShortDesc    string `json:"shortDesc"`
				ExamineDesc  string `json:"examineDesc"`
				Slot         string `json:"slot"`
				Level        int    `json:"level"`
				Weight       int    `json:"weight"`
				IsImmovable  *bool  `json:"isImmovable"`
				IsContainer  *bool  `json:"isContainer"`
				IsReadable   *bool  `json:"isReadable"`
				Content      string `json:"content"`
				ReadSkill    string `json:"readSkill"`
				ReadSkillLevel int  `json:"readSkillLevel"`
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
			if req.ShortDesc != "" {
				updater.SetShortDesc(req.ShortDesc)
			}
			if req.ExamineDesc != "" {
				updater.SetExamineDesc(req.ExamineDesc)
			}
			if req.Slot != "" {
				updater.SetSlot(req.Slot)
			}
			if req.Level > 0 {
				updater.SetLevel(req.Level)
			}
			if req.Weight > 0 {
				updater.SetWeight(req.Weight)
			}
			if req.IsImmovable != nil {
				updater.SetIsImmovable(*req.IsImmovable)
			}
			if req.IsContainer != nil {
				updater.SetIsContainer(*req.IsContainer)
			}
			if req.IsReadable != nil {
				updater.SetIsReadable(*req.IsReadable)
			}
			if req.Content != "" {
				updater.SetContent(req.Content)
			}
			if req.ReadSkill != "" {
				updater.SetReadSkill(req.ReadSkill)
			}
			if req.ReadSkillLevel > 0 {
				updater.SetReadSkillLevel(req.ReadSkillLevel)
			}

			item, err := updater.Save(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}

			c.JSON(http.StatusOK, item)
		})

		// Delete item
		group.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
				return
			}

			err = client.Equipment.DeleteOneID(id).Exec(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}

			c.JSON(http.StatusNoContent, nil)
		})

		// Get items in a room
		group.GET("/room/:roomId", func(c *gin.Context) {
			roomID, err := strconv.Atoi(c.Param("roomId"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
				return
			}

			items, err := client.Equipment.Query().
				Where(db.HasRoomWith(db.Room.ID(roomID))).
				All(c.Request.Context())

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, items)
		})
	}
}

// processHiddenDetails filters hidden details based on examine skill level
func processHiddenDetails(details []map[string]interface{}, examineLevel int) []map[string]interface{} {
	if details == nil {
		return nil
	}

	var revealed []map[string]interface{}
	for _, detail := range details {
		minLevel := 0
		if ml, ok := detail["min_examine_level"].(float64); ok {
			minLevel = int(ml)
		}

		// Check if level is sufficient to reveal
		if examineLevel >= minLevel {
			detail["revealed"] = true
			revealed = append(revealed, detail)
		} else {
			detail["revealed"] = false
			revealed = append(revealed, detail)
		}
	}

	return revealed
}