package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/dblog"
	"herbst-server/repository"
	"log/slog"
)

// updateDialogNode updates an existing dialog node definition.
func updateDialogNode(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "invalid dialog node id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dialog node id"})
			return
		}
		// Verify it exists
		_, err := repos.DialogNode.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog node not found"})
			return
		}
		var input dialogNodeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := repository.DialogNodeUpdates{
			NPCText:        input.NpcText,
			IsEntry:        input.IsEntry,
			EntryCondition: input.EntryCondition,
		}
		if input.Responses != nil {
			responses := responsesToSchema(*input.Responses)
			updates.Responses = &responses
		}
		if input.OnEnterEffects != nil {
			updates.OnEnterEffects = input.OnEnterEffects
		}
		if input.NpcTemplateID != nil {
			updates.NPCTemplateID = input.NpcTemplateID
		}

		_, err = repos.DialogNode.Update(c.Request.Context(), id, updates)
		if err != nil {
			dblog.Error("update dialog node failed", err, slog.String("service", "dialog_nodes"), slog.String("dialog_node_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with npc_template edge for response.
		updated, err := client.DialogNode.Query().
			Where(dialognode.IDEQ(id)).
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			dblog.Error("reload dialog node failed", err, slog.String("service", "dialog_nodes"), slog.String("dialog_node_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("dialog node updated", slog.String("dialog_node_id", id), slog.String("user_email", c.GetString("email")), slog.String("service", "dialog_nodes"))
		c.JSON(http.StatusOK, dialogNodeToView(updated))
	}
}