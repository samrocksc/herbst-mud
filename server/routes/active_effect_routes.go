package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterActiveEffectRoutes registers endpoints for managing active effects on characters.
func RegisterActiveEffectRoutes(r *gin.Engine, repos *repository.Container) {
	effects := r.Group("/api")
	effects.Use(middleware.AuthMiddleware())
	effects.Use(middleware.AdminMiddleware())
	{
		effects.GET("/characters/:id/effects", listActiveEffects(repos))
		effects.DELETE("/characters/:id/effects/:effect_id", removeActiveEffect(repos))
		effects.POST("/characters/:id/effects/apply", applyEffect(repos))
	}
}

type activeEffectView struct {
	ID           int                    `json:"id"`
	CharacterID  int                    `json:"character_id"`
	EffectID     int                    `json:"effect_id"`
	EffectName   string                 `json:"effect_name"`
	EffectType   string                 `json:"effect_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	AppliedByID  int                    `json:"applied_by_id"`
	StackCount   int                    `json:"stack_count"`
	StartedAt    time.Time              `json:"started_at"`
	ExpiresAt    *time.Time             `json:"expires_at"`
	IsActive     bool                   `json:"is_active"`
}

type applyEffectInput struct {
	EffectID     int  `json:"effect_id"`
	AppliedByID  int  `json:"applied_by_id"`
	DurationSecs *int `json:"duration_secs"`
}

func activeEffectToView(ae *db.ActiveEffect) activeEffectView {
	v := activeEffectView{
		ID:          ae.ID,
		CharacterID: ae.CharacterID,
		EffectID:    ae.EffectID,
		AppliedByID: ae.AppliedByID,
		StackCount:  ae.StackCount,
		StartedAt:   ae.StartedAt,
		ExpiresAt:   ae.ExpiresAt,
		IsActive:    ae.IsActive,
	}
	if ae.Edges.Effect != nil {
		v.EffectName = ae.Edges.Effect.Name
		v.EffectType = ae.Edges.Effect.EffectType
		v.Parameters = ae.Edges.Effect.Parameters
	}
	return v
}