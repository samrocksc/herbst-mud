package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
)

// RegisterQuestProgressRoutes registers endpoints for quest progress
// (character-quest interactions). Game-facing routes are public so the
// SSH client can call them without admin auth.
func RegisterQuestProgressRoutes(r *gin.Engine, client *db.Client) {
	// Public routes — game client (SSH) calls these directly
	qp := r.Group("/api/characters/:id/quests")
	{
		qp.GET("", listCharacterQuests(client))
		qp.POST("", acceptQuest(client))
		qp.PUT("/:questId/check", checkProgress(client))
		qp.PUT("/:questId/abandon", abandonQuest(client))
	}

	// Admin-only bulk check endpoint (requires auth)
	admin := r.Group("/api/characters/:id/quests")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/check-all", checkAllQuests(client))
	}
}