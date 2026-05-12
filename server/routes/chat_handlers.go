package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/service"
)

// --- Logging helper ---
func logRequest(ctx *gin.Context, format string, args ...any) {
	slog.Debug(format, args...)
}

// --- Messaging handlers ---

func sendSay(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CharacterID int    `json:"character_id"`
			RoomID      int    `json:"room_id"`
			Message     string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendSay(c.Request.Context(), input.CharacterID, input.RoomID, input.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendYell(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CharacterID int    `json:"character_id"`
			RoomID      int    `json:"room_id"`
			Message     string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendYell(c.Request.Context(), input.CharacterID, input.RoomID, input.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendShout(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CharacterID int    `json:"character_id"`
			Message    string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendShout(c.Request.Context(), input.CharacterID, input.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendTell(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		logRequest(c, "sendTell called")

		var input struct {
			FromID       int    `json:"from_id"`
			ToID         int    `json:"to_id"`
			TargetName   string `json:"target_name"`
			Message      string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			logRequest(c, "sendTell: JSON bind error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logRequest(c, "sendTell: input=%+v", input)

		// If target_name is provided, look up the character ID
		if input.TargetName != "" && input.ToID == 0 {
			logRequest(c, "sendTell: looking up target by name: %s", input.TargetName)
			char, err := svc.Character.QueryCharacterByName(c.Request.Context(), input.TargetName)
			if err != nil {
				logRequest(c, "sendTell: character not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
				return
			}
			logRequest(c, "sendTell: found character %s (ID: %d)", char.Name, char.ID)
			input.ToID = char.ID
		} else {
			logRequest(c, "sendTell: using to_id=%d", input.ToID)
		}

		logRequest(c, "sendTell: calling SendTell with from_id=%d, to_id=%d, message=%s", input.FromID, input.ToID, input.Message)
		result, err := svc.Chat.SendTell(c.Request.Context(), input.FromID, input.ToID, input.Message)
		if err != nil {
			logRequest(c, "sendTell: SendTell error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logRequest(c, "sendTell: success: %v", result)

		// Determine the display message based on whether it's self-tell
		if input.ToID == input.FromID {
			result.DisplayMessage = fmt.Sprintf("You tell yourself, \"%s\"", input.Message)
		} else {
			result.DisplayMessage = fmt.Sprintf("You tell %s, \"%s\"", result.FromCharacterName, input.Message)
		}

		c.JSON(http.StatusOK, result)
	}
}

func sendWhisper(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FromID int    `json:"from_id"`
			ToID   int    `json:"to_id"`
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendWhisper(c.Request.Context(), input.FromID, input.ToID, input.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendEmote(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CharacterID int    `json:"character_id"`
			Action      string `json:"action"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendEmote(c.Request.Context(), input.CharacterID, input.Action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func sendChannel(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			CharacterID int    `json:"character_id"`
			Channel    string `json:"channel"`
			Message    string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendChannel(c.Request.Context(), input.Channel, input.Message, input.CharacterID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// --- Channel handlers ---

func getChannels(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		channels, err := svc.Chat.GetChannels(charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, channels)
	}
}

func setChannelEnabled(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		channel := c.Param("channel")
		var input struct {
			Enabled bool `json:"enabled"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := svc.Chat.SetChannelEnabled(c.Request.Context(), charID, channel, input.Enabled)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func setChannelColor(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		channel := c.Param("channel")
		var input struct {
			Color string `json:"color"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := svc.Chat.SetChannelColor(c.Request.Context(), charID, channel, input.Color)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// --- Ignore handlers ---

func ignorePlayer(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		ignoredID := c.GetInt("characterId")
		if err := svc.Chat.IgnorePlayer(c.Request.Context(), charID, ignoredID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func unignorePlayer(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		ignoredID := c.GetInt("characterId")
		if err := svc.Chat.UnignorePlayer(c.Request.Context(), charID, ignoredID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func getIgnoredPlayers(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		ids, err := svc.Chat.GetIgnoredPlayers(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ids)
	}
}

// --- Offline tell handlers ---

func queueOfflineTell(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FromID        int    `json:"from_id"`
			RecipientName string `json:"recipient_name"`
			Message       string `json:"message"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Chat.QueueOfflineTell(c.Request.Context(), input.FromID, input.RecipientName, input.Message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "queued"})
	}
}

func deliverQueuedTells(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID := c.GetInt("character_id")
		if charID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tells, err := svc.Chat.DeliverQueuedTells(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tells)
	}
}
