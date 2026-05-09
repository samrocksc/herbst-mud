package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/effecthook"
	"herbst-server/db/npctemplate"
)

func listHooks(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		hooks, err := client.EffectHook.Query().
			WithEffect().
			WithNpcTemplate().
			All(c.Request.Context())
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

func getHook(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		h, err := client.EffectHook.Query().
			Where(effecthook.IDEQ(id)).
			WithEffect().
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "hook not found"})
			return
		}
		c.JSON(http.StatusOK, hookToView(h))
	}
}

func listTemplateHooks(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("id")
		hooks, err := client.EffectHook.Query().
			Where(effecthook.HasNpcTemplateWith(npctemplate.IDEQ(templateID))).
			WithEffect().
			WithNpcTemplate().
			All(c.Request.Context())
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

func createTemplateHook(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param("id")
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
		mut := client.EffectHook.Create().
			SetName(*input.Name).
			SetEvent(*input.Event).
			SetNpcTemplateID(templateID).
			SetEffectID(*input.EffectID)
		if input.Target != nil {
			if !validHookTargets[*input.Target] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target"})
				return
			}
			mut.SetTarget(*input.Target)
		}
		if input.Condition != nil {
			mut.SetCondition(*input.Condition)
		}
		if input.Enabled != nil {
			mut.SetEnabled(*input.Enabled)
		}
		h, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h, _ = client.EffectHook.Query().
			Where(effecthook.IDEQ(h.ID)).
			WithEffect().
			WithNpcTemplate().
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, hookToView(h))
	}
}

func updateHook(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		h, err := client.EffectHook.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "hook not found"})
			return
		}
		var input hookInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mut := client.EffectHook.UpdateOne(h)
		if input.Name != nil {
			mut.SetName(*input.Name)
		}
		if input.Event != nil {
			if !validHookEvents[*input.Event] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event"})
				return
			}
			mut.SetEvent(*input.Event)
		}
		if input.Target != nil {
			if !validHookTargets[*input.Target] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target"})
				return
			}
			mut.SetTarget(*input.Target)
		}
		if input.Condition != nil {
			mut.SetCondition(*input.Condition)
		}
		if input.Enabled != nil {
			mut.SetEnabled(*input.Enabled)
		}
		if input.EffectID != nil {
			mut.SetEffectID(*input.EffectID)
		}
		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updated, _ = client.EffectHook.Query().
			Where(effecthook.IDEQ(updated.ID)).
			WithEffect().
			WithNpcTemplate().
			Only(c.Request.Context())
		c.JSON(http.StatusOK, hookToView(updated))
	}
}

func deleteHook(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseHookID(c)
		if err != nil {
			return
		}
		err = client.EffectHook.DeleteOneID(id).Exec(c.Request.Context())
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