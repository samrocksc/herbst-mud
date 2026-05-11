package routes

import (
	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterHookRoutes registers CRUD endpoints for EffectHook definitions.
func RegisterHookRoutes(r *gin.Engine, repos *repository.Container) {
	hooks := r.Group("/api/hooks")
	hooks.Use(middleware.AuthMiddleware())
	hooks.Use(middleware.AdminMiddleware())
	{
		hooks.GET("", listHooks(repos))
		hooks.GET("/:id", getHook(repos))
		hooks.PUT("/:id", updateHook(repos))
		hooks.DELETE("/:id", deleteHook(repos))
	}
	// Template-scoped routes
	r.GET("/api/npc-templates/:id/hooks", listTemplateHooks(repos))
	r.POST("/api/npc-templates/:id/hooks", middleware.AuthMiddleware(), middleware.AdminMiddleware(), createTemplateHook(repos))
}

type hookView struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Event           string `json:"event"`
	Target          string `json:"target"`
	Condition       string `json:"condition"`
	Enabled         bool   `json:"enabled"`
	EffectID        int    `json:"effect_id"`
	EffectName      string `json:"effect_name,omitempty"`
	NPCTemplateID   string `json:"npc_template_id,omitempty"`
	NPCTemplateName string `json:"npc_template_name,omitempty"`
}

type hookInput struct {
	Name          *string `json:"name"`
	Event         *string `json:"event"`
	Target        *string `json:"target"`
	Condition     *string `json:"condition"`
	Enabled       *bool   `json:"enabled"`
	EffectID      *int    `json:"effect_id"`
	NPCTemplateID *string `json:"npc_template_id"`
}

var validHookEvents = map[string]bool{
	"on_death": true, "on_hit_received": true, "on_hit_dealt": true,
	"on_kill": true, "on_enter_room": true, "on_leave_room": true,
	"on_equip": true, "on_unequip": true, "on_login": true,
	"on_effect_start": true, "on_effect_end": true,
}

var validHookTargets = map[string]bool{
	"self": true, "attacker": true, "killer": true, "room": true, "owner": true,
}

func hookToView(h *db.EffectHook) hookView {
	v := hookView{
		ID:        h.ID,
		Name:      h.Name,
		Event:     h.Event,
		Target:    h.Target,
		Condition: h.Condition,
		Enabled:   h.Enabled,
	}
	if h.Edges.Effect != nil {
		v.EffectID = h.Edges.Effect.ID
		v.EffectName = h.Edges.Effect.Name
	}
	if h.Edges.NpcTemplate != nil {
		v.NPCTemplateID = h.Edges.NpcTemplate.ID
		v.NPCTemplateName = h.Edges.NpcTemplate.Name
	}
	return v
}
