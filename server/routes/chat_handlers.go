package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendSay(c.Request.Context(), input.CharacterID, input.RoomID, input.Message)
		if err != nil {
			dblog.Error("failed to send say", err, slog.String("service", "chat"), slog.Int("character_id", input.CharacterID), slog.Int("room_id", input.RoomID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("say sent", slog.Int("character_id", input.CharacterID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendYell(c.Request.Context(), input.CharacterID, input.RoomID, input.Message)
		if err != nil {
			dblog.Error("failed to send yell", err, slog.String("service", "chat"), slog.Int("character_id", input.CharacterID), slog.Int("room_id", input.RoomID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("yell sent", slog.Int("character_id", input.CharacterID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendShout(c.Request.Context(), input.CharacterID, input.Message)
		if err != nil {
			dblog.Error("failed to send shout", err, slog.String("service", "chat"), slog.Int("character_id", input.CharacterID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("shout sent", slog.Int("character_id", input.CharacterID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
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
			dblog.Error("failed to send tell", err, slog.String("service", "chat"), slog.Int("from_id", input.FromID), slog.Int("to_id", input.ToID))
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

		slog.Info("tell sent", slog.Int("from_id", input.FromID), slog.Int("to_id", input.ToID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendWhisper(c.Request.Context(), input.FromID, input.ToID, input.Message)
		if err != nil {
			dblog.Error("failed to send whisper", err, slog.String("service", "chat"), slog.Int("from_id", input.FromID), slog.Int("to_id", input.ToID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("whisper sent", slog.Int("from_id", input.FromID), slog.Int("to_id", input.ToID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendEmote(c.Request.Context(), input.CharacterID, input.Action)
		if err != nil {
			dblog.Error("failed to send emote", err, slog.String("service", "chat"), slog.Int("character_id", input.CharacterID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("emote sent", slog.Int("character_id", input.CharacterID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := svc.Chat.SendChannel(c.Request.Context(), input.Channel, input.Message, input.CharacterID)
		if err != nil {
			dblog.Error("failed to send channel message", err, slog.String("service", "chat"), slog.Int("character_id", input.CharacterID), slog.String("channel", input.Channel))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("channel message sent", slog.Int("character_id", input.CharacterID), slog.String("channel", input.Channel), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			dblog.Error("failed to get channels", err, slog.String("service", "chat"), slog.Int("character_id", charID))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := svc.Chat.SetChannelEnabled(c.Request.Context(), charID, channel, input.Enabled)
		if err != nil {
			dblog.Error("failed to set channel enabled", err, slog.String("service", "chat"), slog.Int("character_id", charID), slog.String("channel", channel))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("channel enabled updated", slog.Int("character_id", charID), slog.String("channel", channel), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := svc.Chat.SetChannelColor(c.Request.Context(), charID, channel, input.Color)
		if err != nil {
			dblog.Error("failed to set channel color", err, slog.String("service", "chat"), slog.Int("character_id", charID), slog.String("channel", channel))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("channel color updated", slog.Int("character_id", charID), slog.String("channel", channel), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			dblog.Error("failed to ignore player", err, slog.String("service", "chat"), slog.Int("character_id", charID), slog.Int("ignored_id", ignoredID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("player ignored", slog.Int("character_id", charID), slog.Int("ignored_id", ignoredID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			dblog.Error("failed to unignore player", err, slog.String("service", "chat"), slog.Int("character_id", charID), slog.Int("ignored_id", ignoredID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("player unignored", slog.Int("character_id", charID), slog.Int("ignored_id", ignoredID), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			dblog.Error("failed to get ignored players", err, slog.String("service", "chat"), slog.Int("character_id", charID))
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
			slog.Warn("bad request", slog.String("service", "chat"), slog.String("reason", "invalid input"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Chat.QueueOfflineTell(c.Request.Context(), input.FromID, input.RecipientName, input.Message); err != nil {
			dblog.Error("failed to queue offline tell", err, slog.String("service", "chat"), slog.Int("from_id", input.FromID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("offline tell queued", slog.Int("from_id", input.FromID), slog.String("recipient", input.RecipientName), slog.String("user_email", c.GetString("email")), slog.String("service", "chat"))
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
			dblog.Error("failed to deliver queued tells", err, slog.String("service", "chat"), slog.Int("character_id", charID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tells)
	}
}
