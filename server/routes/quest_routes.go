package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/middleware"
	"herbst-server/service"
)

// RegisterQuestRoutes registers CRUD endpoints for Quest definitions.
// All quest routes require admin authentication.
func RegisterQuestRoutes(r *gin.Engine, svc *service.Container) {
	quests := r.Group("/api/quests")
	quests.Use(middleware.AuthMiddleware())
	quests.Use(middleware.AdminMiddleware())
	{
		quests.GET("", listQuests(svc))
		quests.POST("", createQuest(svc))
		quests.GET("/:id", getQuest(svc))
		quests.PUT("/:id", updateQuest(svc))
		quests.DELETE("/:id", deleteQuest(svc))
	}
}