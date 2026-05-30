package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/db"
	"herbst-server/db/socialcommand"
	"herbst-server/middleware"
)

// RegisterSocialRoutes registers REST endpoints for social commands.
func RegisterSocialRoutes(r *gin.Engine, client *db.Client) {
	socials := r.Group("/api")
	socials.Use(middleware.AuthMiddleware(nil))
	socials.Use(middleware.AdminMiddleware())
	{
		socials.GET("/socials", listSocials(client))
		socials.GET("/socials/:id", getSocial(client))
		socials.POST("/socials", createSocial(client))
		socials.PUT("/socials/:id", updateSocial(client))
		socials.DELETE("/socials/:id", deleteSocial(client))
	}
}

// ─── Social Command CRUD ─────────────────────────────────────────────────────

func listSocials(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := client.SocialCommand.Query()
		if search := c.Query("search"); search != "" {
			query = query.Where(socialcommand.NameContains(search))
		}
		socials, err := query.All(c.Request.Context())
		if err != nil {
			dblog.Error("failed to list social commands", err, slog.String("service", "socials"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]gin.H, len(socials))
		for i, s := range socials {
			result[i] = socialCommandToJSON(s)
		}
		c.JSON(http.StatusOK, result)
	}
}

func getSocial(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid social command id", slog.String("error", err.Error()), slog.String("service", "socials"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid social command id"})
			return
		}
		s, err := client.SocialCommand.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "social command not found"})
			return
		}
		c.JSON(http.StatusOK, socialCommandToJSON(s))
	}
}

func createSocial(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name           string `json:"name" binding:"required"`
			DisplayName    string `json:"display_name" binding:"required"`
			SelfText       string `json:"self_text" binding:"required"`
			RoomText       string `json:"room_text" binding:"required"`
			TargetSelfText string `json:"target_self_text"`
			TargetText     string `json:"target_text"`
			TargetRoomText string `json:"target_room_text"`
			RequiresTarget bool   `json:"requires_target"`
			IsEmote        bool   `json:"is_emote"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid create social command request", slog.String("error", err.Error()), slog.String("service", "socials"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		builder := client.SocialCommand.Create().
			SetName(req.Name).
			SetDisplayName(req.DisplayName).
			SetSelfText(req.SelfText).
			SetRoomText(req.RoomText).
			SetTargetSelfText(req.TargetSelfText).
			SetTargetText(req.TargetText).
			SetTargetRoomText(req.TargetRoomText).
			SetRequiresTarget(req.RequiresTarget).
			SetIsEmote(req.IsEmote)
		created, err := builder.Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to create social command", err, slog.String("service", "socials"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("social command created", slog.Int("social_id", created.ID), slog.String("service", "socials"))
		c.JSON(http.StatusCreated, socialCommandToJSON(created))
	}
}

func updateSocial(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid social command id", slog.String("error", err.Error()), slog.String("service", "socials"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid social command id"})
			return
		}
		var req struct {
			Name           *string `json:"name"`
			DisplayName    *string `json:"display_name"`
			SelfText       *string `json:"self_text"`
			RoomText       *string `json:"room_text"`
			TargetSelfText *string `json:"target_self_text"`
			TargetText     *string `json:"target_text"`
			TargetRoomText *string `json:"target_room_text"`
			RequiresTarget *bool   `json:"requires_target"`
			IsEmote        *bool   `json:"is_emote"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			slog.Warn("invalid update social command request", slog.String("error", err.Error()), slog.String("service", "socials"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		builder := client.SocialCommand.UpdateOneID(id)
		if req.Name != nil {
			builder.SetName(*req.Name)
		}
		if req.DisplayName != nil {
			builder.SetDisplayName(*req.DisplayName)
		}
		if req.SelfText != nil {
			builder.SetSelfText(*req.SelfText)
		}
		if req.RoomText != nil {
			builder.SetRoomText(*req.RoomText)
		}
		if req.TargetSelfText != nil {
			builder.SetTargetSelfText(*req.TargetSelfText)
		}
		if req.TargetText != nil {
			builder.SetTargetText(*req.TargetText)
		}
		if req.TargetRoomText != nil {
			builder.SetTargetRoomText(*req.TargetRoomText)
		}
		if req.RequiresTarget != nil {
			builder.SetRequiresTarget(*req.RequiresTarget)
		}
		if req.IsEmote != nil {
			builder.SetIsEmote(*req.IsEmote)
		}
		updated, err := builder.Save(c.Request.Context())
		if err != nil {
			dblog.Error("failed to update social command", err, slog.String("service", "socials"), slog.Int("social_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("social command updated", slog.Int("social_id", updated.ID), slog.String("service", "socials"))
		c.JSON(http.StatusOK, socialCommandToJSON(updated))
	}
}

func deleteSocial(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			slog.Warn("invalid social command id", slog.String("error", err.Error()), slog.String("service", "socials"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid social command id"})
			return
		}
		if err := client.SocialCommand.DeleteOneID(id).Exec(c.Request.Context()); err != nil {
			dblog.Error("failed to delete social command", err, slog.String("service", "socials"), slog.Int("social_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("social command deleted", slog.Int("social_id", id), slog.String("service", "socials"))
		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── JSON Helper ──────────────────────────────────────────────────────────────

func socialCommandToJSON(s *db.SocialCommand) gin.H {
	return gin.H{
		"id":               s.ID,
		"name":             s.Name,
		"display_name":     s.DisplayName,
		"self_text":        s.SelfText,
		"room_text":        s.RoomText,
		"target_self_text": s.TargetSelfText,
		"target_text":      s.TargetText,
		"target_room_text": s.TargetRoomText,
		"requires_target":  s.RequiresTarget,
		"is_emote":         s.IsEmote,
	}
}
