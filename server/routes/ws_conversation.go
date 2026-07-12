package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
)

// sendConversationScreen pushes a conversation view to the WebSocket client.
func sendConversationScreen(wsc *WSConn, npcName, templateID string, nodes []*db.DialogNode, currentNodeID string) {
	nodeMap := make(map[string]dialogNodeView, len(nodes))
	for _, n := range nodes {
		nodeMap[n.ID] = dialogNodeToViewSimpleWithTemplate(n, templateID)
	}
	wsc.send(ServerMessage{
		Type: MsgScreen,
		Text: npcName,
		Data: gin.H{
			"view_type":       "conversation",
			"npc_name":        npcName,
			"npc_template_id": templateID,
			"nodes":           nodeMap,
			"current_node_id": currentNodeID,
		},
		Timestamp: time.Now().UnixMilli(),
	})
}
