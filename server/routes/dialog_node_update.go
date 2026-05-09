package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
)

// updateDialogNode updates an existing dialog node definition.
func updateDialogNode(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dialog node id"})
			return
		}
		existing, err := client.DialogNode.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog node not found"})
			return
		}
		var input dialogNodeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mut := client.DialogNode.UpdateOne(existing)
		if input.NpcText != nil {
			mut.SetNpcText(*input.NpcText)
		}
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
		if input.NpcTemplateID != nil {
			mut.SetNpcTemplateID(*input.NpcTemplateID)
		}
		updated, err := mut.Save(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with npc_template edge for response.
		updated, err = client.DialogNode.Query().
			Where(dialognode.IDEQ(updated.ID)).
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dialogNodeToView(updated))
	}
}