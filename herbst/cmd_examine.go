package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleExamineCommand handles the examine/ex command
func (m *model) handleExamineCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Examine what? Usage: examine <item>", "error")
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(target)

	// Room items (visible)
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue
		}
		if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Hidden items that reveal on examine
	if m.currentRoom > 0 {
		m.examineHiddenItems(target)
		return
	}

	// Inventory
	m.examineInventory(target)

	// NPCs
	m.examineNPCs(target)

	m.AppendMessage(fmt.Sprintf("You don't see '%s' here.", target), "error")
}

// examineHiddenItems checks for hidden items that reveal on examine
func (m *model) examineHiddenItems(target string) bool {
	resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var allItems []RoomItem
	if json.NewDecoder(resp.Body).Decode(&allItems) != nil {
		return false
	}

	for _, item := range allItems {
		if !item.IsVisible && item.RevealCondition != nil {
			revealType, _ := item.RevealCondition["type"].(string)
			revealTarget, _ := item.RevealCondition["target"].(string)
			if revealType == "examine" && strings.ToLower(revealTarget) == target {
				if m.revealHiddenItem(item, revealTarget) {
					return true
				}
			}
		}
	}
	return false
}

// revealHiddenItem reveals a hidden item
func (m *model) revealHiddenItem(item RoomItem, target string) bool {
	revealResp, err := httpPost(
		fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
		fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, target, m.characterLevel),
	)
	if err != nil {
		return false
	}
	defer revealResp.Body.Close()

	if revealResp.StatusCode != http.StatusOK {
		return false
	}

	m.loadRoomItems()
	for _, ri := range m.roomItems {
		if strings.Contains(strings.ToLower(ri.Name), target) || strings.ToLower(ri.Name) == target {
			m.AppendMessage("✨ You discovered something hidden!\n\n", "info")
			m.displayItemDetails(ri)
			return true
		}
	}
	return false
}

// examineInventory checks player inventory for an item
func (m *model) examineInventory(target string) bool {
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var items []RoomItem
	if json.NewDecoder(resp.Body).Decode(&items) != nil {
		return false
	}

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return true
		}
	}
	return false
}

// examineNPCs checks NPCs in the room
func (m *model) examineNPCs(target string) bool {
	if m.currentRoom == 0 {
		return false
	}

	resp, err := httpGet(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var npcs []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		Disposition string `json:"disposition"`
	}
	if json.NewDecoder(resp.Body).Decode(&npcs) != nil {
		return false
	}

	for _, npc := range npcs {
		if strings.Contains(strings.ToLower(npc.Name), target) || strings.ToLower(npc.Name) == target {
			m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
				npc.Name, npc.Description, npc.Level, npc.Disposition), "info")
			return true
		}
	}
	return false
}