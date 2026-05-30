package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/service"
	"log/slog"
)

// createRoom creates a new room.
func createRoom(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name           string         `json:"name" binding:"required"`
			Description    string         `json:"description" binding:"required"`
			IsStartingRoom bool           `json:"isStartingRoom"`
			IsRootRoom     bool           `json:"isRootRoom"`
			Exits          map[string]int `json:"exits"`
			Atmosphere     string         `json:"atmosphere"`
			PosX           int            `json:"posX"`
			PosY           int            `json:"posY"`
			PosZ           int            `json:"posZ"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room, err := svc.Room.CreateRoom(c.Request.Context(), service.CreateRoomInput{
			Name:           input.Name,
			Description:    input.Description,
			IsStartingRoom: input.IsStartingRoom,
			IsRootRoom:     input.IsRootRoom,
			Exits:          input.Exits,
			Atmosphere:     input.Atmosphere,
			PosX:           input.PosX,
			PosY:           input.PosY,
			PosZ:           input.PosZ,
			WorldID:        c.Query("world_id"),
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
func updateRoom(svc *service.Container) gin.HandlerFunc {
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
			PosX           *int            `json:"posX"`
			PosY           *int            `json:"posY"`
			PosZ           *int            `json:"posZ"`
			Version        *int            `json:"version"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "rooms"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room, err := svc.Room.UpdateRoom(c.Request.Context(), id, service.UpdateRoomInput{
			Name:           input.Name,
			Description:    input.Description,
			IsStartingRoom: input.IsStartingRoom,
			IsRootRoom:     input.IsRootRoom,
			Exits:          input.Exits,
			Atmosphere:     input.Atmosphere,
			PosX:           input.PosX,
			PosY:           input.PosY,
			PosZ:           input.PosZ,
			Version:        input.Version,
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
