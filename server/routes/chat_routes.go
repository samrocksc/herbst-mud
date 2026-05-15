package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/middleware"
	"herbst-server/service"
)

// RegisterChatRoutes registers chat and messaging endpoints.
func RegisterChatRoutes(r *gin.Engine, svc *service.Container) {
	chat := r.Group("/api/chat")
	chat.Use(middleware.AuthMiddleware(nil))
	{
		// messaging
		chat.POST("/say", sendSay(svc))
		chat.POST("/yell", sendYell(svc))
		chat.POST("/shout", sendShout(svc))
		chat.POST("/tell", sendTell(svc))
		chat.POST("/whisper", sendWhisper(svc))
		chat.POST("/emote", sendEmote(svc))
		chat.POST("/channel", sendChannel(svc))

		// channels
		chat.GET("/channels", getChannels(svc))
		chat.PUT("/channels/:channel/enabled", setChannelEnabled(svc))
		chat.PUT("/channels/:channel/color", setChannelColor(svc))

		// ignore
		chat.POST("/ignore/:characterId", ignorePlayer(svc))
		chat.DELETE("/ignore/:characterId", unignorePlayer(svc))
		chat.GET("/ignored", getIgnoredPlayers(svc))

		// offline tells
		chat.POST("/tell/offline", queueOfflineTell(svc))
		chat.POST("/tells/deliver", deliverQueuedTells(svc))
	}
}
