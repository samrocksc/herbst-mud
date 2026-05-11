package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

func listHooks(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		hooks, err := repos.EffectHook.ListWithEdges(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]hookView, len(hooks))
		for i, h := range hooks {
			result[i] = hookToView(h)
		}
		c.JSON(http.StatusOK, gin.H{"hooks": result})
	}
}

func getHook(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		h, err := repos.EffectHook.GetWithEdges(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "hook not found"})
			return
		}
		c.JSON(http.StatusOK, hookToView(h))
	}
}

func listTemplateHooks(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("id")
		if templateID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}
		hooks, err := repos.EffectHook.ListByTemplateWithEdges(c.Request.Context(), templateID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]hookView, len(hooks))
		for i, h := range hooks {
			result[i] = hookToView(h)
		}
		c.JSON(http.StatusOK, gin.H{"hooks": result})
	}
}

func createTemplateHook(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("id")
		if templateID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
			return
		}
		var input hookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name == nil || *input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if input.Event == nil || !validHookEvents[*input.Event] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event"})
			return
		}
		if input.EffectID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "effect_id is required"})
			return
		}
		enabled := true
		if input.Enabled != nil {
			enabled = *input.Enabled
		}
		target := ""
		if input.Target != nil {
			if !validHookTargets[*input.Target] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target"})
				return
			}
			target = *input.Target
		}
		condition := ""
		if input.Condition != nil {
			condition = *input.Condition
		}
		h, err := repos.EffectHook.Create(c.Request.Context(), repository.CreateEffectHookInput{
			Name:          *input.Name,
			Event:         *input.Event,
			Target:        target,
			Condition:     condition,
			Enabled:       enabled,
			EffectID:      *input.EffectID,
			NPCTemplateID: &templateID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Re-fetch with edges for response
		h, _ = repos.EffectHook.GetWithEdges(c.Request.Context(), h.ID)
		c.JSON(http.StatusCreated, hookToView(h))
	}
}

func updateHook(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		// Verify hook exists
		_, err = repos.EffectHook.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "hook not found"})
			return
		}
		var input hookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := repository.EffectHookUpdates{
			Name:      input.Name,
			Event:     input.Event,
			Target:    input.Target,
			Condition: input.Condition,
			Enabled:   input.Enabled,
			EffectID:  input.EffectID,
		}

		_, err = repos.EffectHook.Update(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Re-fetch with edges for response
		updated, _ := repos.EffectHook.GetWithEdges(c.Request.Context(), id)
		c.JSON(http.StatusOK, hookToView(updated))
	}
}

func deleteHook(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		err = repos.EffectHook.Delete(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "hook not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func parseHookID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hook id"})
		return 0, err
	}
	return id, nil
}