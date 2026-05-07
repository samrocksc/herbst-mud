package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// handleTakeCommand picks up an item from the room.
func (m *model) handleTakeCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Take what? Usage: take <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")

	m.loadRoomItems()

	var targetItem *RoomItem
	for i := range m.roomItems {
		if strings.Contains(strings.ToLower(m.roomItems[i].Name), strings.ToLower(itemName)) {
			targetItem = &m.roomItems[i]
			break
		}
	}

	if targetItem == nil {
		m.AppendMessage(fmt.Sprintf("You don't see any %s here.", itemName), "error")
		return
	}

	if targetItem.IsImmovable {
		colorStyle := lipgloss.NewStyle()
		if targetItem.Color != "" {
			colorStyle = colorStyle.Foreground(lipgloss.Color(targetItem.Color))
		} else {
			colorStyle = colorStyle.Foreground(itemColorGold)
		}
		m.AppendMessage(fmt.Sprintf("You can't take the %s.", colorStyle.Render(targetItem.Name)), "error")
		return
	}

	url := fmt.Sprintf("%s/equipment/%d", RESTAPIBase, targetItem.ID)
	jsonData, _ := json.Marshal(map[string]interface{}{"roomId": nil})
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error picking up item: %v", err), "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * 1000000000}
	resp, err := client.Do(req)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error picking up item: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage(fmt.Sprintf("Failed to pick up %s.", targetItem.Name), "error")
		return
	}

	m.AppendMessage(fmt.Sprintf("You pick up the %s.", targetItem.Name), "success")
}

// handleDropCommand drops an item from inventory.
func (m *model) handleDropCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Drop what? Usage: drop <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")
	m.AppendMessage(fmt.Sprintf("You don't have any %s to drop.", itemName), "error")
}