package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterEffectDefRoutes registers CRUD endpoints for Effect definitions.
// These are the data-driven state-change effects (xp_drain, hp_change, etc.),
// separate from AbilityEffect (combat-only effects tied to abilities).
func RegisterEffectDefRoutes(r *gin.Engine, repos *repository.Container) {
	effects := r.Group("/api/effects")
	effects.Use(middleware.AuthMiddleware(nil))
	effects.Use(middleware.AdminMiddleware())
	{
		effects.GET("", listEffectDefs(repos))
		effects.POST("", createEffectDef(repos))
		effects.GET("/:id", getEffectDef(repos))
		effects.PUT("/:id", updateEffectDef(repos))
		effects.DELETE("/:id", deleteEffectDef(repos))
	}
}

type effectDefView struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	EffectType   string                 `json:"effect_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	StackMode    string                 `json:"stack_mode"`
	StackLimit   int                    `json:"stack_limit"`
	IsPermanent  bool                   `json:"is_permanent"`
	DurationSecs int                    `json:"duration_secs"`
	Messages     map[string]string      `json:"messages"`
	HookCount    int                    `json:"hook_count"`
}

type effectDefInput struct {
	Name         *string                 `json:"name"`
	Description  *string                 `json:"description"`
	EffectType   *string                 `json:"effect_type"`
	Parameters   *map[string]interface{} `json:"parameters"`
	StackMode    *string                 `json:"stack_mode"`
	StackLimit   *int                    `json:"stack_limit"`
	IsPermanent  *bool                   `json:"is_permanent"`
	DurationSecs *int                    `json:"duration_secs"`
	Messages     *map[string]string      `json:"messages"`
}

func effectDefToView(e *db.Effect) effectDefView {
	hookCount := 0
	if e.Edges.Hooks != nil {
		hookCount = len(e.Edges.Hooks)
	}
	return effectDefView{
		ID:           e.ID,
		Name:         e.Name,
		Description:  e.Description,
		EffectType:   e.EffectType,
		Parameters:   e.Parameters,
		StackMode:    e.StackMode,
		StackLimit:   e.StackLimit,
		IsPermanent:  e.IsPermanent,
		DurationSecs: e.DurationSecs,
		Messages:     e.Messages,
		HookCount:    hookCount,
	}
}