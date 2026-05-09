package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/room"
)

var oppositeDir = map[string]string{
	"north":     "south",
	"south":     "north",
	"east":      "west",
	"west":      "east",
	"northeast": "southwest",
	"southwest": "northeast",
	"northwest": "southeast",
	"southeast": "northwest",
	"up":        "down",
	"down":      "up",
}

// RegisterRoomRoutes registers all room-related routes
func RegisterRoomRoutes(router *gin.Engine, client *db.Client) {
	// Create a new room
	router.POST("/rooms", func(c *gin.Context) {
		var req struct {
			Name           string         `json:"name" binding:"required"`
			Description    string         `json:"description" binding:"required"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			IsRootRoom     bool           `json:"isRootRoom"`
			Exits          map[string]int `json:"exits"`
			PosX           int            `json:"posX"`
			PosY           int            `json:"posY"`
			PosZ           int            `json:"posZ"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.PosX = max(0, req.PosX)
		req.PosY = max(0, req.PosY)
		// Enforce single root room
		if req.IsRootRoom {
			client.Room.Update().
				Where(room.IsRootRoom(true)).
				SetIsRootRoom(false).
				Save(c.Request.Context())
		}

		roomBuilder := client.Room.
			Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetIsStartingRoom(req.IsStartingRoom).
			SetIsRootRoom(req.IsRootRoom).
			SetExits(req.Exits).
			SetPosX(req.PosX).
			SetPosY(req.PosY).
			SetPosZ(req.PosZ)
		room, err := roomBuilder.Save(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, room)
	})

	// Get all rooms
	router.GET("/rooms", func(c *gin.Context) {
		rooms, err := client.Room.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rooms)
	})

	// Get a single room by ID
	router.GET("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		room, err := client.Room.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.JSON(http.StatusOK, room)
	})

	// Update a room by ID (optimistic locking via version field)
	router.PUT("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var req struct {
			Name           string         `json:"name"`
			Description    string         `json:"description"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			IsRootRoom     bool           `json:"isRootRoom"`
			Exits          map[string]int `json:"exits"`
			PosX           *int           `json:"posX"`
			PosY           *int           `json:"posY"`
			PosZ           *int           `json:"posZ"`
			Version        int            `json:"version"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch current room to check version (optimistic locking)
		currentRoom, err := client.Room.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// If version is provided (>= 1), verify it matches the current version
		if req.Version > 0 && currentRoom.Version != req.Version {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "Conflict: room has been modified by another editor",
				"current":    currentRoom,
				"yourVersion": req.Version,
			})
			return
		}

		updater := client.Room.UpdateOneID(id)

		// Only update fields that are provided
		if req.Name != "" {
			updater.SetName(req.Name)
		}
		if req.Description != "" {
			updater.SetDescription(req.Description)
		}
		// For boolean and map fields, we'll always update them if provided
		updater.SetIsStartingRoom(req.IsStartingRoom)
		updater.SetIsRootRoom(req.IsRootRoom)
		if req.IsRootRoom {
			// Enforce single root room: clear isRootRoom on all other rooms
			client.Room.Update().
				Where(room.IsRootRoom(true)).
				SetIsRootRoom(false).
				Save(c.Request.Context())
		}
		if req.Exits != nil {
			updater.SetExits(req.Exits)
		}
		if req.PosX != nil {
			updater.SetPosX(max(0, *req.PosX))
		}
		if req.PosY != nil {
			updater.SetPosY(max(0, *req.PosY))
		}
		if req.PosZ != nil {
			updater.SetPosZ(*req.PosZ)
		}

		// Increment version on every save
		updater.AddVersion(1)

		room, err := updater.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		c.JSON(http.StatusOK, room)
	})

	// Delete a room by ID (transaction-wrapped with cascade cleanup)
	router.DELETE("/rooms/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		ctx := c.Request.Context()
		tx, err := client.Tx(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		txClient := tx.Client()

		// Get the root room ID for relocating characters
		rootRooms, err := txClient.Room.Query().
			Where(room.IsRootRoom(true)).
			All(ctx)
		if err != nil || len(rootRooms) == 0 {
			rootRooms = []*db.Room{{ID: 5}}
		}
		defaultRoomID := rootRooms[0].ID

		// Move any characters in this room to the starting room
		_, err = txClient.Character.Update().
			Where(character.CurrentRoomIdEQ(id)).
			SetCurrentRoomId(defaultRoomID).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to relocate characters"})
			return
		}

		// Remove this room from all other rooms' exits
		allRooms, err := txClient.Room.Query().All(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query rooms"})
			return
		}

		for _, r := range allRooms {
			if r.ID == id || r.Exits == nil {
				continue
			}
			needsUpdate := false
			newExits := make(map[string]int)
			for dir, targetID := range r.Exits {
				if targetID != id {
					newExits[dir] = targetID
				} else {
					needsUpdate = true
				}
			}
			if needsUpdate {
				_, err = txClient.Room.UpdateOneID(r.ID).
					SetExits(newExits).
					AddVersion(1).
					Save(ctx)
				if err != nil {
					_ = tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean up orphan exits"})
					return
				}
			}
		}

		// Now delete the room
		err = txClient.Room.DeleteOneID(id).Exec(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// Clean up orphan exits (exits pointing to non-existent rooms)
	router.POST("/rooms/cleanup-orphan-exits", func(c *gin.Context) {
		ctx := c.Request.Context()
		allRooms, err := client.Room.Query().All(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		roomIDs := make(map[int]bool)
		for _, r := range allRooms {
			roomIDs[r.ID] = true
		}

		cleaned := 0
		for _, r := range allRooms {
			if r.Exits == nil {
				continue
			}
			newExits := make(map[string]int)
			changed := false
			for dir, targetID := range r.Exits {
				if roomIDs[targetID] {
					newExits[dir] = targetID
				} else {
					changed = true
				}
			}
			if changed {
				_, err := client.Room.UpdateOneID(r.ID).
					SetExits(newExits).
					AddVersion(1).
					Save(ctx)
				if err != nil {
					continue
				}
				cleaned += len(r.Exits) - len(newExits)
			}
		}

		c.JSON(http.StatusOK, gin.H{"cleaned": cleaned})
	})

	// Create a bidirectional exit (transaction-wrapped)
	router.POST("/rooms/:id/exits/bidirectional", func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source room ID"})
			return
		}

		var req struct {
			Direction    string `json:"direction" binding:"required"`
			TargetRoomID int   `json:"targetRoomId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		reverseDir, ok := oppositeDir[req.Direction]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown direction: " + req.Direction})
			return
		}

		ctx := c.Request.Context()
		tx, err := client.Tx(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		txClient := tx.Client()

		// Fetch source room
		sourceRoom, err := txClient.Room.Get(ctx, sourceID)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Source room not found"})
			return
		}

		// Fetch target room
		targetRoom, err := txClient.Room.Get(ctx, req.TargetRoomID)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Target room not found"})
			return
		}

		// Add exit to source room
		sourceExits := sourceRoom.Exits
		if sourceExits == nil {
			sourceExits = make(map[string]int)
		}
		sourceExits[req.Direction] = req.TargetRoomID

		sourceRoom, err = txClient.Room.UpdateOneID(sourceID).
			SetExits(sourceExits).
			AddVersion(1).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source room"})
			return
		}

		// Add reverse exit to target room
		targetExits := targetRoom.Exits
		if targetExits == nil {
			targetExits = make(map[string]int)
		}
		targetExits[reverseDir] = sourceID

		targetRoom, err = txClient.Room.UpdateOneID(req.TargetRoomID).
			SetExits(targetExits).
			AddVersion(1).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target room"})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"source": sourceRoom, "target": targetRoom})
	})

	// Delete a bidirectional exit (transaction-wrapped)
	router.DELETE("/rooms/:id/exits/bidirectional", func(c *gin.Context) {
		sourceID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source room ID"})
			return
		}

		direction := c.Query("direction")
		if direction == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "direction query parameter is required"})
			return
		}

		reverseDir, ok := oppositeDir[direction]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown direction: " + direction})
			return
		}

		ctx := c.Request.Context()
		tx, err := client.Tx(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer func() { _ = tx.Rollback() }()

		txClient := tx.Client()

		// Fetch source room
		sourceRoom, err := txClient.Room.Get(ctx, sourceID)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Source room not found"})
			return
		}

		// Remove exit from source room
		sourceExits := sourceRoom.Exits
		if sourceExits == nil {
			sourceExits = make(map[string]int)
		}
		targetID, hasTarget := sourceExits[direction]
		delete(sourceExits, direction)

		sourceRoom, err = txClient.Room.UpdateOneID(sourceID).
			SetExits(sourceExits).
			AddVersion(1).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source room"})
			return
		}

		// Remove reverse exit from target room
		if hasTarget && targetID > 0 {
			targetRoom, err := txClient.Room.Get(ctx, targetID)
			if err == nil && targetRoom.Exits != nil {
				targetExits := targetRoom.Exits
				delete(targetExits, reverseDir)
				_, err = txClient.Room.UpdateOneID(targetID).
					SetExits(targetExits).
					AddVersion(1).
					Save(ctx)
				if err != nil {
					_ = tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target room"})
					return
				}
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"source": sourceRoom})
	})

	// Get characters in a room (for displaying NPCs vs players)
	router.GET("/rooms/:id/characters", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		characters, err := client.Character.Query().
			Where(character.CurrentRoomId(id)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return full character data including HP for combat visibility
		result := make([]gin.H, len(characters))
		for i, char := range characters {
			entry := gin.H{
				"id":         char.ID,
				"name":       char.Name,
				"isNPC":      char.IsNPC,
				"level":      char.Level,
				"class":      char.Class,
				"race":       char.Race,
				"hp":         char.Hitpoints,
				"maxHp":      char.MaxHitpoints,
				"lastSeenAt": char.LastSeenAt,
			}
			if char.IsNPC && char.NpcTemplateID != "" {
				if tmpl, err := client.NPCTemplate.Get(c.Request.Context(), char.NpcTemplateID); err == nil && tmpl.XpValue > 0 {
					entry["xpValue"] = tmpl.XpValue
				} else {
					entry["xpValue"] = char.Level * 10
				}
			}
			result[i] = entry
		}

		c.JSON(http.StatusOK, result)
	})

	// Get room with items and NPCs (look-10)
	router.GET("/rooms/:id/look", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		room, err := client.Room.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Get characters in room (NPCs)
		characters, _ := client.Character.Query().
			Where(character.CurrentRoomId(id)).
			All(c.Request.Context())

		// Separate NPCs from players
		npcs := []gin.H{}
		players := []gin.H{}
		for _, char := range characters {
			charData := gin.H{
				"id":    char.ID,
				"name":  char.Name,
				"level": char.Level,
				"class": char.Class,
				"race":  char.Race,
			}
			if char.IsNPC {
				npcs = append(npcs, charData)
			} else {
				players = append(players, charData)
			}
		}

		// Get visible items in room using edge query
		items, _ := client.Equipment.Query().
			All(c.Request.Context())

		// Filter to items in this room using edge
		visibleItems := []gin.H{}
		for _, item := range items {
			if item.IsVisible && item.Edges.Room != nil && item.Edges.Room.ID == id {
				visibleItems = append(visibleItems, gin.H{
					"id":           item.ID,
					"name":         item.Name,
					"description":  item.Description,
					"type":         item.ItemType,
					"is_immovable": item.IsImmovable,
					"color":        item.Color,
				})
			}
		}

		// Build look response
		c.JSON(http.StatusOK, gin.H{
			"id":          room.ID,
			"name":        room.Name,
			"description": room.Description,
			"exits":       room.Exits,
			"z_level":     0, // Default z-level
			"items":       visibleItems,
			"npcs":        npcs,
			"players":     players,
		})
	})
}