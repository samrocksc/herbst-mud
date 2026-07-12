package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/service"
)

// RegisterRoomRoutes registers all room-related routes.
func RegisterRoomRoutes(router *gin.Engine, client *db.Client, svc *service.Container) {
	rooms := router.Group("/api")
	rooms.Use(middleware.AuthMiddleware(nil))
	rooms.Use(middleware.AdminMiddleware())
	rooms.Use(middleware.WorldAccessMiddleware())
	{
		rooms.GET("/rooms", listRooms(svc))
		rooms.POST("/rooms", createRoom(svc, client))
		rooms.GET("/rooms/:id", getRoom(svc))
		rooms.PUT("/rooms/:id", updateRoom(svc, client))
		rooms.DELETE("/rooms/:id", deleteRoom(svc))
		rooms.POST("/rooms/cleanup-orphan-exits", cleanupOrphanExits(svc))
		rooms.POST("/rooms/:id/exits/bidirectional", createBidirectionalExit(svc))
		rooms.DELETE("/rooms/:id/exits/bidirectional", deleteBidirectionalExit(svc))
	}
	// Look/characters endpoints still need direct DB access until CharacterService is migrated
	rc := &roomClient{svc: svc, db: client}
	rooms.GET("/rooms/:id/characters", rc.getCharacters)
	rooms.GET("/rooms/:id/look", rc.getLook)
}