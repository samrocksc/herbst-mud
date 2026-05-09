package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/activeeffect"
	"herbst-server/db/character"
	"herbst-server/db/effect"
)

func listActiveEffects(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		effects, err := client.ActiveEffect.Query().
			Where(
				activeeffect.CharacterIDEQ(charID),
				activeeffect.IsActiveEQ(true),
			).
			WithEffect().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]activeEffectView, len(effects))
		for i, ae := range effects {
			result[i] = activeEffectToView(ae)
		}
		c.JSON(http.StatusOK, gin.H{"active_effects": result})
	}
}

func removeActiveEffect(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		effectID, err := strconv.Atoi(c.Param("effect_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid effect id"})
			return
		}
		err = client.ActiveEffect.UpdateOneID(effectID).
			SetIsActive(false).
			Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "active effect not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func applyEffect(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}
		var input applyEffectInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Verify character exists
		_, err = client.Character.Get(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
			return
		}
		// Load effect definition
		eff, err := client.Effect.Get(c.Request.Context(), input.EffectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "effect not found"})
			return
		}
		// Handle stack mode
		existing, _ := client.ActiveEffect.Query().
			Where(
				activeeffect.HasCharacterWith(character.IDEQ(charID)),
				activeeffect.HasEffectWith(effect.IDEQ(input.EffectID)),
				activeeffect.IsActiveEQ(true),
			).
			Only(c.Request.Context())

		if existing != nil {
			switch eff.StackMode {
			case "replace":
				client.ActiveEffect.UpdateOne(existing).
					SetStackCount(1).
					SetStartedAt(time.Now()).
					Save(c.Request.Context())
			case "refresh":
				if !eff.IsPermanent && eff.DurationSecs > 0 {
					client.ActiveEffect.UpdateOne(existing).
						SetExpiresAt(time.Now().Add(time.Duration(eff.DurationSecs) * time.Second)).
						Save(c.Request.Context())
				}
			case "stack":
				if existing.StackCount < eff.StackLimit {
					client.ActiveEffect.UpdateOne(existing).
						SetStackCount(existing.StackCount + 1).
						Save(c.Request.Context())
				}
			}
			// Return updated
			updated, _ := client.ActiveEffect.Query().
				Where(activeeffect.IDEQ(existing.ID)).
				WithEffect().
				Only(c.Request.Context())
			c.JSON(http.StatusOK, activeEffectToView(updated))
			return
		}
		// Create new active effect
		mut := client.ActiveEffect.Create().
			SetCharacterID(charID).
			SetEffectID(input.EffectID).
			SetAppliedByID(input.AppliedByID)
		if input.DurationSecs != nil && *input.DurationSecs > 0 {
			mut.SetExpiresAt(time.Now().Add(time.Duration(*input.DurationSecs) * time.Second))
		} else if !eff.IsPermanent && eff.DurationSecs > 0 {
			mut.SetExpiresAt(time.Now().Add(time.Duration(eff.DurationSecs) * time.Second))
		}
		ae, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ae, _ = client.ActiveEffect.Query().
			Where(activeeffect.IDEQ(ae.ID)).
			WithEffect().
			Only(c.Request.Context())
		c.JSON(http.StatusCreated, activeEffectToView(ae))
	}
}