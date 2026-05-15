package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
	"herbst-server/service"
)

// RegisterQuestProgressRoutes registers endpoints for quest progress
// (character-quest interactions). Game-facing routes are public so the
// SSH client can call them without admin auth.
func RegisterQuestProgressRoutes(r *gin.Engine, repos *repository.Container, svc *service.Container, client *db.Client) {
	// Public routes — game client (SSH) calls these directly
	qp := r.Group("/api/characters/:id/quests")
	{
		qp.GET("", listCharacterQuests(svc))
		qp.POST("", acceptQuest(svc))
		// TODO: migrate checkProgress/abandonQuest/checkAllQuests to fully use repos once
		// QuestProgressRepo has methods for complex multi-table queries
		qp.PUT("/:questId/check", checkProgress(repos, client))
		qp.PUT("/:questId/abandon", abandonQuest(repos, client))
	}

	// Admin-only bulk check endpoint (requires auth)
	admin := r.Group("/api/characters/:id/quests")
	admin.Use(middleware.AuthMiddleware(nil))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/check-all", checkAllQuests(repos, client))
	}
}