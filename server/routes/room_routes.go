package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/service"
)

// RegisterRoomRoutes registers all room-related routes.
func RegisterRoomRoutes(router *gin.Engine, client *db.Client, svc *service.Container) {
	rooms := router.Group("/rooms")
	{
		rooms.POST("", createRoom(svc))
		rooms.GET("", listRooms(svc))
		rooms.GET("/:id", getRoom(svc))
		rooms.PUT("/:id", updateRoom(svc))
		rooms.DELETE("/:id", deleteRoom(svc))
		rooms.POST("/cleanup-orphan-exits", cleanupOrphanExits(svc))
		rooms.POST("/:id/exits/bidirectional", createBidirectionalExit(svc))
		rooms.DELETE("/:id/exits/bidirectional", deleteBidirectionalExit(svc))
	}
	// Look/characters endpoints still need direct DB access until CharacterService is migrated
	rc := &roomClient{svc: svc, db: client}
	rooms.GET("/:id/characters", rc.getCharacters)
	rooms.GET("/:id/look", rc.getLook)
}