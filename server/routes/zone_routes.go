package routes

import (
	"context"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
	"herbst-server/service"
)

// RegisterZoneRoutes registers REST endpoints for zones.
func RegisterZoneRoutes(r *gin.Engine, svc *service.Container, repos *repository.Container) {
	g := r.Group("/api")
	g.Use(middleware.AuthMiddleware(nil))
	g.Use(middleware.AdminMiddleware())
	{
		g.GET("/zones", listZones(svc))
		g.POST("/zones", createZone(svc, repos))
		g.GET("/zones/:id", getZone(svc))
		g.PUT("/zones/:id", updateZone(svc))
		g.DELETE("/zones/:id", deleteZone(svc))
		g.GET("/zones/:id/rooms", listZoneRooms(svc, repos))
	}
}

type zoneInput struct {
	ID           string `json:"id" binding:"required"`
	WorldID      string `json:"world_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	MinLevel     int    `json:"min_level"`
	ParentZoneID string `json:"parent_zone_id"`
	Color        string `json:"color"`
	RoomIDs      []int  `json:"room_ids"`
}

func listZones(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		if worldID == "" {
			worldID = middleware.GetWorldID(c)
		}
		zones, err := svc.Zone.ListZonesByWorld(c.Request.Context(), worldID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"zones": zones})
	}
}

func createZone(svc *service.Container, repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input zoneInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		zone, err := svc.Zone.CreateZone(c.Request.Context(), repository.CreateZoneInput{
			ID:           input.ID,
			WorldID:      input.WorldID,
			Name:         input.Name,
			Description:  input.Description,
			MinLevel:     input.MinLevel,
			ParentZoneID: input.ParentZoneID,
			Color:        input.Color,
			RoomIDs:      input.RoomIDs,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, zone)
	}
}

func getZone(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		zone, err := svc.Zone.GetZone(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "zone not found"})
			return
		}
		c.JSON(http.StatusOK, zone)
	}
}

func updateZone(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Name         *string `json:"name"`
			Description  *string `json:"description"`
			MinLevel     *int    `json:"min_level"`
			ParentZoneID *string `json:"parent_zone_id"`
			Color        *string `json:"color"`
			RoomIDs      *[]int  `json:"room_ids"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		zone, err := svc.Zone.UpdateZone(c.Request.Context(), id, repository.ZoneUpdates{
			Name:         input.Name,
			Description:  input.Description,
			MinLevel:     input.MinLevel,
			ParentZoneID: input.ParentZoneID,
			Color:        input.Color,
			RoomIDs:      input.RoomIDs,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, zone)
	}
}

func deleteZone(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := svc.Zone.DeleteZone(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "zone deleted"})
	}
}

// zoneRoomView is the JSON shape for a single room in the zone list.
type zoneRoomView struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Exists  bool   `json:"exists"`
	Message string `json:"message,omitempty"`
}

// listZoneRooms returns the rooms in a zone, marking rooms that have been
// removed from the world as `exists: false`. The list is sorted by ID.
func listZoneRooms(svc *service.Container, repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")
		zone, err := svc.Zone.GetZone(ctx, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "zone not found"})
			return
		}
		// Build a set of existing room IDs in this zone's world.
		existing, err := repos.Room.List(ctx, zone.WorldID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		exists := make(map[int]string, len(existing))
		for _, rm := range existing {
			exists[rm.ID] = rm.Name
		}
		// Build the response from zone.room_ids (persistent list).
		views := make([]zoneRoomView, 0, len(zone.RoomIds))
		for _, rid := range zone.RoomIds {
			v := zoneRoomView{ID: rid, Exists: false}
			if name, ok := exists[rid]; ok {
				v.Exists = true
				v.Name = name
			} else {
				v.Message = "room has been removed from this world"
			}
			views = append(views, v)
		}
		sort.Slice(views, func(i, j int) bool { return views[i].ID < views[j].ID })
		c.JSON(http.StatusOK, gin.H{"zone_id": zone.ID, "rooms": views})
	}
}

// Compile-time check: ensure context is used (suppresses unused import if
// future refactors drop the import).
var _ = context.Background

// ZoneView is the JSON shape for zone API responses.
type ZoneView struct {
	ID           string `json:"id"`
	WorldID      string `json:"world_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	MinLevel     int    `json:"min_level"`
	ParentZoneID string `json:"parent_zone_id,omitempty"`
	Color        string `json:"color,omitempty"`
}

func zoneToView(z *db.Zone) ZoneView {
	return ZoneView{
		ID:           z.ID,
		WorldID:      z.WorldID,
		Name:         z.Name,
		Description:  z.Description,
		MinLevel:     z.MinLevel,
		ParentZoneID: z.ParentZoneID,
		Color:        z.Color,
	}
}
