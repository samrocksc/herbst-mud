package routes

import (
	"context"
	"log/slog"
	"strconv"

	"herbst-server/db"
	"herbst-server/db/schema"
	"herbst-server/dblog"
	"herbst-server/repository"
)

// handleDialogChoice processes a player's dialog response. The command format is
// "dialog <template_id> <node_id> [<choice_index>]". choice_index is 1-based.
func handleDialogChoice(templateID, nodeID, choiceStr string, wsc *WSConn, repos *repository.Container, client *db.Client) string {
	ctx := context.Background()

	node, err := repos.DialogNode.Get(ctx, nodeID)
	if err != nil {
		dblog.Error("handleDialogChoice: failed to get dialog node", err, slog.String("node_id", nodeID))
		return "Failed to load dialog node."
	}

	tmpl, err := repos.NPCTemplate.Get(ctx, templateID)
	if err != nil {
		dblog.Error("handleDialogChoice: failed to get NPC template", err, slog.String("template_id", templateID))
		return "Invalid NPC template."
	}

	slog.Info("dialog choice", slog.Int("character_id", wsc.CharacterID), slog.String("template_id", templateID), slog.String("node_id", nodeID), slog.String("choice", choiceStr))

	resp := pickResponse(node.Responses, choiceStr)
	if resp == nil {
		return "That is not a valid response."
	}

	applyDialogEffects(ctx, wsc, repos, client, resp.Effects)

	if resp.NextNodeID == "" {
		return "You end the conversation."
	}

	nextNode, err := repos.DialogNode.Get(ctx, resp.NextNodeID)
	if err != nil {
		dblog.Error("handleDialogChoice: failed to get next node", err, slog.String("node_id", resp.NextNodeID))
		return "The conversation trails off."
	}

	applyDialogEffects(ctx, wsc, repos, client, nextNode.OnEnterEffects)

	nodes, _ := repos.DialogNode.ListByTemplate(ctx, templateID)
	sendConversationScreen(wsc, tmpl.Name, templateID, nodes, nextNode.ID)
	return ""
}

// pickResponse selects the response chosen by the player. choiceStr is 1-based.
func pickResponse(responses []schema.DialogResponse, choiceStr string) *schema.DialogResponse {
	if len(responses) == 0 {
		return nil
	}
	if choiceStr == "" {
		return &responses[0]
	}
	idx, err := strconv.Atoi(choiceStr)
	if err != nil || idx < 1 || idx > len(responses) {
		return nil
	}
	return &responses[idx-1]
}
