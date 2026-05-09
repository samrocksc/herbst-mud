package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
)

// RegisterQuestRoutes registers CRUD endpoints for Quest definitions.
// All quest routes require admin authentication.
func RegisterQuestRoutes(r *gin.Engine, client *db.Client) {
	quests := r.Group("/api/quests")
	quests.Use(middleware.AuthMiddleware())
	quests.Use(middleware.AdminMiddleware())
	{
		quests.GET("", listQuests(client))
		quests.POST("", createQuest(client))
		quests.GET("/:id", getQuest(client))
		quests.PUT("/:id", updateQuest(client))
		quests.DELETE("/:id", deleteQuest(client))
	}
}