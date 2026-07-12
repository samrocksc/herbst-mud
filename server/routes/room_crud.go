package routes

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/room"
	"herbst-server/dblog"
	"herbst-server/service"
	"log/slog"
)

// validateExits checks if all exits point to existing rooms.
func validateExits(ctx context.Context, client *db.Client, exits map[string]int, worldID string) []string {
	var errors []string
	query := client.Room.Query()
	if worldID != "" {
		query = query.Where(room.WorldID(worldID))
	}
	rooms, err := query.All(ctx)
	if err != nil {
		dblog.Error("failed to list rooms for exit validation", err, slog.String("service", "rooms"))
		errors = append(errors, "failed to validate exits: could not list rooms")
		return errors
	}
	roomIDs := make(map[int]bool)
	for _, r := range rooms {
		roomIDs[r.ID] = true
	}
	for dir, targetID := range exits {
		if targetID == 0 {
			continue
		}
		if !roomIDs[targetID] {
			errors = append(errors, fmt.Sprintf("exit '%s' points to non-existent room %d", dir, targetID))
		}
	}
	return errors
}

// createRoom creates a new room.
func createRoom(svc *service.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name           string         `json:"name" binding:"required"`
			Description    string         `json:"description" binding:"required"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			IsRootRoom     bool           `json:"isRootRoom"`
			Exits          map[string]int `json:"exits"`
			Atmosphere     string         `json:"atmosphere"`
			PosZ           int            `json:"posZ"`
			ZoneIDs        []string       `json:"zoneIds"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate exits point to existing rooms
		if len(input.Exits) > 0 {
			if errors := validateExits(c.Request.Context(), client, input.Exits, c.Query("world_id")); len(errors) > 0 {
				slog.Warn("invalid exits", slog.String("service", "rooms"), slog.Any("errors", errors), slog.String("client_ip", c.ClientIP()))
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exits: " + strings.Join(errors, "; ")})
				return
			}
		}
		room, err := svc.Room.CreateRoom(c.Request.Context(), service.CreateRoomInput{
			Name:           input.Name,
			Description:    input.Description,
			IsStartingRoom: input.IsStartingRoom,
			IsRootRoom:     input.IsRootRoom,
			Exits:          input.Exits,
			Atmosphere:     input.Atmosphere,
			PosZ:           input.PosZ,
			WorldID:        c.Query("world_id"),
			ZoneIDs:        input.ZoneIDs,
		})
		if err != nil {
			dblog.Error("create room failed", err, slog.String("service", "rooms"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("room created", slog.Int("room_id", room.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "rooms"))
		c.JSON(http.StatusCreated, room)
	}
}

// listRooms returns all rooms, optionally filtered by name and world_id.
func listRooms(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get world_id from query params for world-scoped listing
		worldID := c.Query("world_id")
		rooms, err := svc.Room.ListRooms(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("list rooms failed", err, slog.String("service", "rooms"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if search := c.Query("search"); search != "" {
			s := strings.ToLower(search)
			filtered := make([]*db.Room, 0, len(rooms))
			for _, r := range rooms {
				if strings.Contains(strings.ToLower(r.Name), s) {
					filtered = append(filtered, r)
				}
			}
			rooms = filtered
		}
		if zoneID := c.Query("zone_id"); zoneID != "" {
			filtered := make([]*db.Room, 0, len(rooms))
			for _, r := range rooms {
				if slices.Contains(r.ZoneIds, zoneID) {
					filtered = append(filtered, r)
				}
			}
			rooms = filtered
		}
		c.JSON(http.StatusOK, rooms)
	}
}

// getRoom returns a single room by ID.
func getRoom(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid room id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
			return
		}
		room, err := svc.Room.GetRoom(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusOK, room)
	}
}

// updateRoom updates an existing room with optimistic locking.
func updateRoom(svc *service.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid room id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
			return
		}
		var input struct {
			Name           *string         `json:"name"`
			Description    *string         `json:"description"`
			IsStartingRoom *bool           `json:"isStartingRoom"`
			IsRootRoom     *bool           `json:"isRootRoom"`
			Exits          *map[string]int `json:"exits"`
			Atmosphere     *string         `json:"atmosphere"`
			PosZ           *int            `json:"posZ"`
			Version        *int            `json:"version"`
			ZoneIDs        *[]string       `json:"zoneIds"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate exits point to existing rooms (if exits are being updated)
		if input.Exits != nil {
			if errors := validateExits(c.Request.Context(), client, *input.Exits, c.Query("world_id")); len(errors) > 0 {
				slog.Warn("invalid exits", slog.String("service", "rooms"), slog.Any("errors", errors), slog.String("client_ip", c.ClientIP()))
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exits: " + strings.Join(errors, "; ")})
				return
			}
		}
		room, err := svc.Room.UpdateRoom(c.Request.Context(), id, service.UpdateRoomInput{
			Name:           input.Name,
			Description:    input.Description,
			IsStartingRoom: input.IsStartingRoom,
			IsRootRoom:     input.IsRootRoom,
			Exits:          input.Exits,
			Atmosphere:     input.Atmosphere,
			PosZ:           input.PosZ,
			Version:        input.Version,
			ZoneIDs:        input.ZoneIDs,
		})
		if err != nil {
			if err.Error() == "version conflict" || (err.Error() != "" && len(err.Error()) > 16 && err.Error()[:16] == "version conflict") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", err.Error()), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Info("room updated", slog.Int("room_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "rooms"))
		c.JSON(http.StatusOK, room)
	}
}

// deleteRoom deletes a room and relocates characters.
func deleteRoom(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid room id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
			return
		}
		if err := svc.Room.DeleteRoom(c.Request.Context(), id); err != nil {
			dblog.Error("delete room failed", err, slog.String("service", "rooms"), slog.Int("room_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("room deleted", slog.Int("room_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "rooms"))
		c.Status(http.StatusNoContent)
	}
}
