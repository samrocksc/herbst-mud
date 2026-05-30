package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/dialognode"
	"herbst-server/db/schema"
	"herbst-server/dblog"
	"herbst-server/repository"
	"log/slog"
)

// createDialogNode creates a new dialog node definition.
func createDialogNode(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dialogNodeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "invalid json"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.ID == nil || *input.ID == "" {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "missing id"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		if input.NpcText == nil || *input.NpcText == "" {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "missing npc_text"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "npc_text is required"})
			return
		}
		if input.NpcTemplateID == nil || *input.NpcTemplateID == "" {
			slog.Warn("bad request", slog.String("service", "dialog_nodes"), slog.String("reason", "missing npc_template_id"), slog.String("client_ip", c.ClientIP()))
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
			dblog.Error("create dialog node failed", err, slog.String("service", "dialog_nodes"), slog.String("dialog_node_id", *input.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Reload with npc_template edge for response.
		dn, err = client.DialogNode.Query().
			Where(dialognode.IDEQ(dn.ID)).
			WithNpcTemplate().
			Only(c.Request.Context())
		if err != nil {
			dblog.Error("reload dialog node failed", err, slog.String("service", "dialog_nodes"), slog.String("dialog_node_id", dn.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("dialog node created", slog.String("dialog_node_id", dn.ID), slog.String("user_email", c.GetString("email")), slog.String("service", "dialog_nodes"))
		c.JSON(http.StatusCreated, dialogNodeToView(dn))
	}
}