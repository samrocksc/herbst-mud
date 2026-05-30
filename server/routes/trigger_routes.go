package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterTriggerRoutes registers CRUD endpoints for Trigger definitions.
func RegisterTriggerRoutes(r *gin.Engine, repos *repository.Container) {
	triggers := r.Group("/api/triggers")
	triggers.Use(middleware.AuthMiddleware(nil))
	triggers.Use(middleware.AdminMiddleware())
	{
		triggers.GET("", listTriggers(repos))
		triggers.POST("", createTrigger(repos))
		triggers.GET("/:id", getTrigger(repos))
		triggers.PUT("/:id", updateTrigger(repos))
		triggers.DELETE("/:id", deleteTrigger(repos))
	}

	// Public: game client fetches triggers for a room
	r.GET("/api/rooms/:id/triggers", middleware.AuthMiddleware(nil), middleware.AdminMiddleware(), getRoomTriggers(repos))

	// Public: game client fetches triggers for an equipment item
	r.GET("/api/equipment/:id/triggers", middleware.AuthMiddleware(nil), middleware.AdminMiddleware(), getEquipmentTriggers(repos))
}

// triggerView is the JSON shape returned by the API.
type triggerView struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	WorldID     string `json:"world_id"`
	TriggerType string `json:"trigger_type"`
	TargetType  string `json:"target_type"`
	TargetID    int    `json:"target_id"`
	RoomID      *int   `json:"room_id,omitempty"`
	EquipmentID *int   `json:"equipment_id,omitempty"`
	Condition   string `json:"condition,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// triggerInput is the request body for create and update.
type triggerInput struct {
	Name        string  `json:"name"`
	WorldID     string  `json:"world_id"`
	TriggerType string  `json:"trigger_type"`
	TargetType  string  `json:"target_type"`
	TargetID    int     `json:"target_id"`
	RoomID      *int    `json:"room_id,omitempty"`
	EquipmentID *int    `json:"equipment_id,omitempty"`
	Condition   string  `json:"condition"`
	Enabled     bool    `json:"enabled"`
}

func triggerToView(t *db.Trigger) triggerView {
	return triggerView{
		ID:          t.ID,
		Name:        t.Name,
		WorldID:     t.WorldID,
		TriggerType: t.TriggerType,
		TargetType:  t.TargetType,
		TargetID:    t.TargetID,
		RoomID:      t.RoomID,
		EquipmentID: t.EquipmentID,
		Condition:   t.Condition,
		Enabled:     t.Enabled,
	}
}

func listTriggers(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		triggers, err := repos.Trigger.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]triggerView, len(triggers))
		for i, t := range triggers {
			result[i] = triggerToView(t)
		}
		c.JSON(http.StatusOK, gin.H{"triggers": result})
	}
}

func createTrigger(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input triggerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		if input.TriggerType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "trigger_type is required"})
			return
		}

		if input.TargetType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target_type is required"})
			return
		}

		created, err := repos.Trigger.Create(c.Request.Context(), repository.CreateTriggerInput{
			Name:        input.Name,
			WorldID:     input.WorldID,
			TriggerType: input.TriggerType,
			TargetType:  input.TargetType,
			TargetID:    input.TargetID,
			RoomID:      input.RoomID,
			EquipmentID: input.EquipmentID,
			Condition:   input.Condition,
			Enabled:     input.Enabled,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, triggerToView(created))
	}
}

func getTrigger(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trigger id"})
			return
		}
		t, err := repos.Trigger.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}
		c.JSON(http.StatusOK, triggerToView(t))
	}
}

func updateTrigger(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trigger id"})
			return
		}

		_, err = repos.Trigger.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}

		var input triggerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := repository.TriggerUpdates{
			Name:        &input.Name,
			WorldID:     &input.WorldID,
			TriggerType: &input.TriggerType,
			TargetType:  &input.TargetType,
			TargetID:    &input.TargetID,
			RoomID:      input.RoomID,
			EquipmentID: input.EquipmentID,
			Condition:   &input.Condition,
			Enabled:     &input.Enabled,
		}

		updated, err := repos.Trigger.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, triggerToView(updated))
	}
}

func deleteTrigger(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trigger id"})
			return
		}
		err = repos.Trigger.Delete(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// getRoomTriggers returns all triggers for a room (public access).
func getRoomTriggers(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
			return
		}

		triggers, err := repos.Trigger.ListByRoom(c.Request.Context(), roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]triggerView, len(triggers))
		for i, t := range triggers {
			result[i] = triggerToView(t)
		}
		c.JSON(http.StatusOK, gin.H{"triggers": result})
	}
}

// getEquipmentTriggers returns all triggers for an equipment item (public access).
func getEquipmentTriggers(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		equipmentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment id"})
			return
		}

		triggers, err := repos.Trigger.ListByEquipment(c.Request.Context(), equipmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]triggerView, len(triggers))
		for i, t := range triggers {
			result[i] = triggerToView(t)
		}
		c.JSON(http.StatusOK, gin.H{"triggers": result})
	}
}
