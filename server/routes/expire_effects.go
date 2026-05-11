package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/repository"
)

// RegisterExpiryRoutes registers the endpoint that finds and deactivates
// expired active effects. Called by the herbst effects service on a timer.
func RegisterExpiryRoutes(r *gin.Engine, repos *repository.Container) {
	r.GET("/api/effects/expired", expireEffects(repos))
}

func expireEffects(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		deactivated, err := repos.ActiveEffect.DeactivateExpired(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type expiredItem struct {
			ID          int    `json:"id"`
			CharacterID int    `json:"character_id"`
			EffectID    int    `json:"effect_id"`
			EffectType  string `json:"effect_type"`
		}

		result := make([]expiredItem, 0, len(deactivated))
		for _, ae := range deactivated {
			effectType := ""
			if ae.Edges.Effect != nil {
				effectType = ae.Edges.Effect.EffectType
			}

			result = append(result, expiredItem{
				ID:          ae.ID,
				CharacterID: ae.CharacterID,
				EffectID:    ae.EffectID,
				EffectType:  effectType,
			})
		}

		c.JSON(http.StatusOK, gin.H{"expired": result})
	}
}