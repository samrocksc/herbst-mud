package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ============================================================
// CRAFTING COMMANDS — craft, recipes, stations
// ============================================================

// handleCraftCommand crafts an item using a recipe
func (m *model) handleCraftCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to craft items.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: craft <recipe_name>\nType 'recipes' to see available recipes.", "error")
		return
	}

	recipeName := strings.Join(parts[1:], " ")
	url := fmt.Sprintf("%s/api/characters/%d/craft", RESTAPIBase, m.currentCharacterID)
	payload := fmt.Sprintf(`{"recipe": "%s"}`, recipeName)

	resp, err := httpPost(url, payload)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error crafting: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		if json.NewDecoder(resp.Body).Decode(&errResp) == nil && errResp.Error != "" {
			m.AppendMessage(errResp.Error, "error")
		} else {
			m.AppendMessage(fmt.Sprintf("Crafting failed (status %d)", resp.StatusCode), "error")
		}
		return
	}

	var craftResp struct {
		Success bool `json:"success"`
		Outputs []struct {
			Name       string `json:"name"`
			InstanceID int    `json:"instance_id"`
		} `json:"outputs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&craftResp); err != nil {
		m.AppendMessage("Error reading crafting response", "error")
		return
	}

	if !craftResp.Success {
		m.AppendMessage("Crafting was not successful.", "error")
		return
	}

	output := "Crafted items:\n"
	for _, item := range craftResp.Outputs {
		output += fmt.Sprintf("  • %s\n", item.Name)
	}
	m.AppendMessage(output, "success")
}

// handleRecipesCommand lists available recipes
func (m *model) handleRecipesCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to view recipes.", "error")
		return
	}

	parts := strings.Fields(cmd)
	var url string
	if len(parts) >= 2 {
		stationTag := parts[1]
		url = fmt.Sprintf("%s/api/recipes?station_tag=%s", RESTAPIBase, stationTag)
	} else {
		url = fmt.Sprintf("%s/api/recipes", RESTAPIBase)
	}

	resp, err := httpGet(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching recipes: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage(fmt.Sprintf("Could not load recipes (status %d)", resp.StatusCode), "error")
		return
	}

	var recipes []struct {
		Name         string   `json:"name"`
		DisplayName  string   `json:"display_name"`
		StationTags  []string `json:"station_tags"`
		OutputItem   string   `json:"output_item"`
		OutputAmount int      `json:"output_amount"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&recipes); err != nil {
		m.AppendMessage("Error reading recipes", "error")
		return
	}

	if len(recipes) == 0 {
		m.AppendMessage("No recipes found.", "info")
		return
	}

	output := "=== Recipes ===\n\n"
	for _, r := range recipes {
		stationList := ""
		if len(r.StationTags) > 0 {
			stationList = " [" + strings.Join(r.StationTags, ", ") + "]"
		}
		amount := r.OutputAmount
		if amount == 0 {
			amount = 1
		}
		outputName := r.DisplayName
		if outputName == "" {
			outputName = r.Name
		}
		output += fmt.Sprintf("  %s%s\n    → %s ×%d\n", outputName, stationList, outputName, amount)
	}
	m.AppendMessage(output, "info")
}

// handleStationsCommand shows crafting stations in the current room
func (m *model) handleStationsCommand() {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to see stations.", "error")
		return
	}

	if m.currentRoom == 0 {
		m.AppendMessage("You are not in a room.", "error")
		return
	}

	resp, err := httpGet(fmt.Sprintf("%s/rooms/%d", RESTAPIBase, m.currentRoom))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching room: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage(fmt.Sprintf("Could not load room (status %d)", resp.StatusCode), "error")
		return
	}

	var room struct {
		ID          int      `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		m.AppendMessage("Error reading room data", "error")
		return
	}

	if len(room.Tags) == 0 {
		m.AppendMessage("This room has no crafting stations.", "info")
		return
	}

	output := "=== Crafting Stations ===\n\n"
	output += fmt.Sprintf("Room: %s\n\n", room.Name)
	for _, tag := range room.Tags {
		// Convert snake_case to Title Case for display
		displayTag := strings.Title(strings.ReplaceAll(tag, "_", " "))
		output += fmt.Sprintf("  • %s\n", displayTag)
	}
	m.AppendMessage(output, "info")
}