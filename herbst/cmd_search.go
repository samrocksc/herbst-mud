package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// handleSearchCommand handles the search/perception command
func (m *model) handleSearchCommand(cmd string) {
	if m.currentRoom == 0 {
		m.AppendMessage("You can't search here.", "error")
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error searching: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Error searching the area.", "error")
		return
	}

	var allItems []RoomItem
	if err := json.NewDecoder(resp.Body).Decode(&allItems); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing items: %v", err), "error")
		return
	}

	var found []string
	revealed := 0

	for _, item := range allItems {
		if item.IsVisible {
			continue
		}
		if item.RevealCondition != nil {
			revealType, _ := item.RevealCondition["type"].(string)
			if revealType == "perception_check" {
				revealResp, err := httpPost(
					fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
					fmt.Sprintf(`{"revealType":"perception_check","skillLevel":%d}`, m.characterLevel),
				)
				if err == nil {
					defer revealResp.Body.Close()
					if revealResp.StatusCode == http.StatusOK {
						revealed++
						found = append(found, item.Name)
					}
				}
			}
		}
	}

	m.loadRoomItems()

	if revealed > 0 {
		m.AppendMessage(fmt.Sprintf("🔍 You search the area carefully...\n\n✨ You discovered %d hidden item(s): %s",
			revealed, strings.Join(found, ", ")), "success")
	} else {
		m.AppendMessage("🔍 You search the area carefully...\n\nYou find nothing of interest.", "info")
	}
}