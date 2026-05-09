package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
)

// createDialogNode creates a new dialog node definition.
func createDialogNode(client *db.Client) gin.HandlerFunc {
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
		mut := client.DialogNode.Create().
			SetID(*input.ID).
			SetNpcText(*input.NpcText).
			SetNpcTemplateID(*input.NpcTemplateID)
		if input.IsEntry != nil {
			mut.SetIsEntry(*input.IsEntry)
		}
		if input.EntryCondition != nil {
			mut.SetEntryCondition(*input.EntryCondition)
		}
		if input.Responses != nil {
			mut.SetResponses(responsesToSchema(*input.Responses))
		}
		if input.OnEnterEffects != nil {
			mut.SetOnEnterEffects(*input.OnEnterEffects)
		}
		dn, err := mut.Save(c.Request.Context())
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