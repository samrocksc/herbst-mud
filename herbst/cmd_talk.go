package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// handleTalkCommand starts or resumes a conversation with an NPC.
func (m *model) handleTalkCommand(cmd string) {
	parts := strings.Fields(cmd)
	if m.conversation != nil {
		// Already in conversation — resume
		m.renderCurrentNode()
		return
	}
	if len(parts) < 2 {
		m.AppendMessage("Talk to whom? (talk <npc_name>)", "error")
		return
	}
	targetName := strings.Join(parts[1:], " ")
	// Find NPC in room
	var found *RoomCharacter
	for i := range m.roomCharacters {
		rc := &m.roomCharacters[i]
		if rc.IsNPC && fuzzyWordMatch(rc.Name, targetName) {
			found = rc
			break
		}
	}
	if found == nil {
		m.AppendMessage(fmt.Sprintf("You don't see %s here.", targetName), "error")
		return
	}
	if found.NpcTemplateID == "" {
		m.AppendMessage(fmt.Sprintf("%s doesn't want to talk.", found.Name), "info")
		return
	}
	// Fetch dialog nodes for this NPC template
	nodes, err := m.fetchDialogNodes(found.NpcTemplateID)
	if err != nil || len(nodes) == 0 {
		m.AppendMessage(fmt.Sprintf("%s has nothing to say right now.", found.Name), "info")
		return
	}
	// Find entry node
	nodeMap := make(map[string]*DialogNode, len(nodes))
	var entryNode *DialogNode
	for i := range nodes {
		nodeMap[nodes[i].ID] = &nodes[i]
		if nodes[i].IsEntry {
			entryNode = &nodes[i]
		}
	}
	if entryNode == nil {
		// Use first node as fallback
		entryNode = &nodes[0]
	}
	m.conversation = &ConversationState{
		NPCTemplateID: found.NpcTemplateID,
		NPCName:       found.Name,
		CurrentNodeID: entryNode.ID,
		Nodes:         nodeMap,
		StartedAt:     time.Now(),
	}
	m.renderCurrentNode()
}

// handleDialogChoice processes a numbered response choice (1-9).
func (m *model) handleDialogChoice(choice int) {
	if m.conversation == nil {
		return
	}
	node, ok := m.conversation.Nodes[m.conversation.CurrentNodeID]
	if !ok || choice < 1 || choice > len(node.Responses) {
		m.AppendMessage("Invalid choice. Pick a number from the list.", "error")
		return
	}
	resp := node.Responses[choice-1]
	// Apply effects from the response via event hooks
	for _, effectID := range resp.Effects {
		m.effectsService.FireEvent("on_dialog_response", m.currentCharacterID, m.conversation.NPCTemplateID, map[string]interface{}{"effect_id": effectID})
	}
	// Handle quest offer
	if resp.QuestOfferID != "" {
		questID, err := strconv.Atoi(resp.QuestOfferID)
		if err == nil {
			_, err := m.questService.AcceptQuest(m.currentCharacterID, questID)
			if err != nil {
				m.AppendMessage(fmt.Sprintf("Quest acceptance failed: %v", err), "error")
			} else {
				m.AppendMessage("Quest accepted!", "success")
			}
		}
	}
	// End conversation if no next node
	if resp.NextNodeID == "" {
		m.endConversation()
		return
	}
	// Navigate to next node
	nextNode, ok := m.conversation.Nodes[resp.NextNodeID]
	if !ok {
		m.AppendMessage("The conversation ends abruptly.", "info")
		m.endConversation()
		return
	}
	// Apply on-enter effects via event hooks
	for _, effectID := range nextNode.OnEnterEffects {
		m.effectsService.FireEvent("on_dialog_enter", m.currentCharacterID, m.conversation.NPCTemplateID, map[string]interface{}{"effect_id": effectID})
	}
	m.conversation.CurrentNodeID = nextNode.ID
	m.renderCurrentNode()
}

// endConversation clears the conversation state.
func (m *model) endConversation() {
	if m.conversation != nil {
		m.AppendMessage(fmt.Sprintf("You end your conversation with %s.", m.conversation.NPCName), "info")
	}
	m.conversation = nil
}

// renderCurrentNode displays the current dialog node.
func (m *model) renderCurrentNode() {
	if m.conversation == nil {
		return
	}
	node, ok := m.conversation.Nodes[m.conversation.CurrentNodeID]
	if !ok {
		m.AppendMessage("The conversation has ended.", "info")
		m.conversation = nil
		return
	}
	var b strings.Builder
	b.WriteString(dialogNPCStyle.Render(fmt.Sprintf("%s:", m.conversation.NPCName)) + "\n")
	b.WriteString(fmt.Sprintf("  \"%s\"\n\n", node.NPCText))
	if len(node.Responses) > 0 {
		for i, resp := range node.Responses {
			label := resp.Label
			if label == "" {
				label = "[Leave]"
			}
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, label))
		}
		b.WriteString("\n" + infoStyle.Render("  Choose a number, or 0/leave/bye to end."))
	}
	m.AppendMessage(b.String(), "info")
}

// fetchDialogNodes loads all dialog nodes for an NPC template from the REST API.
func (m *model) fetchDialogNodes(templateID string) ([]DialogNode, error) {
	url := fmt.Sprintf("%s/npc-templates/%s/dialog-nodes", RESTAPIBase, templateID)
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %d", url, resp.StatusCode)
	}
	var result struct {
		Nodes []DialogNode `json:"nodes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Nodes, nil
}