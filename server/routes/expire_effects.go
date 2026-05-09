package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/activeeffect"
)

// RegisterExpiryRoutes registers the endpoint that finds and deactivates
// expired active effects. Called by the herbst effects service on a timer.
func RegisterExpiryRoutes(r *gin.Engine, client *db.Client) {
	r.GET("/api/effects/expired", expireEffects(client))
}

func expireEffects(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		// Find all active effects that have expired
		expired, err := client.ActiveEffect.Query().
			Where(
				activeeffect.IsActiveEQ(true),
				activeeffect.ExpiresAtLTE(now),
			).
			WithEffect().
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type expiredItem struct {
			ID          int    `json:"id"`
			CharacterID int   `json:"character_id"`
			EffectID    int    `json:"effect_id"`
			EffectType  string `json:"effect_type"`
		}

		result := make([]expiredItem, 0, len(expired))
		for _, ae := range expired {
			// Deactivate the effect
			client.ActiveEffect.UpdateOneID(ae.ID).
				SetIsActive(false).
				Exec(c.Request.Context())

			effectType := ""
			if ae.Edges.Effect != nil {
				effectType = ae.Edges.Effect.EffectType
			}

			result = append(result, expiredItem{
				ID:           ae.ID,
				CharacterID:  ae.CharacterID,
				EffectID:     ae.EffectID,
				EffectType:   effectType,
			})
		}

		c.JSON(http.StatusOK, gin.H{"expired": result})
	}
}