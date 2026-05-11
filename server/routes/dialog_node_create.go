package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/repository"
	"herbst-server/db/schema"
)

// createDialogNode creates a new dialog node definition.
func createDialogNode(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dialogNodeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.ID == nil || *input.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		if input.NpcText == nil || *input.NpcText == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "npc_text is required"})
			return
		}
		if input.NpcTemplateID == nil || *input.NpcTemplateID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "npc_template_id is required"})
			return
		}

		var responses []schema.DialogResponse
		if input.Responses != nil {
			responses = responsesToSchema(*input.Responses)
		}
		var onEnterEffects []int
		if input.OnEnterEffects != nil {
			onEnterEffects = *input.OnEnterEffects
		}

		dn, err := repos.DialogNode.Create(c.Request.Context(), repository.CreateDialogNodeInput{
			ID:             *input.ID,
			NPCTemplateID:  *input.NpcTemplateID,
			NPCText:        *input.NpcText,
			Responses:      responses,
			IsEntry:        input.IsEntry != nil && *input.IsEntry,
			EntryCondition: "",
			OnEnterEffects: onEnterEffects,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with npc_template edge for response.
		dn, err = client.DialogNode.Query().
			Where(dialognode.IDEQ(dn.ID)).
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, dialogNodeToView(dn))
	}
}