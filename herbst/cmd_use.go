package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ============================================================
// USE COMMANDS — interact with objects
// ============================================================

// handleUseCommand handles the "use" command for objects
func (m *model) handleUseCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use objects.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: use <object_name>\nType 'look' to see available objects.", "error")
		return
	}

	objectName := strings.Join(parts[1:], " ")

	// First, get triggers for the current room
	roomTriggers, err := m.getTriggersForRoom(m.currentRoom)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching room triggers: %v", err), "error")
		return
	}

	// Check if there's a matching trigger for this object
	var matchedTrigger *triggerView
	for _, t := range roomTriggers {
		if fuzzyWordMatch(t.Name, objectName) {
			matchedTrigger = t
			break
		}
	}

	if matchedTrigger == nil {
		// Try to find in inventory items too
		inventoryTriggers, err := m.getTriggersForCharacterInventory(m.currentCharacterID)
		if err == nil {
			for _, t := range inventoryTriggers {
				if fuzzyWordMatch(t.Name, objectName) {
					matchedTrigger = t
					break
				}
			}
		}
	}

	if matchedTrigger == nil {
		m.AppendMessage(fmt.Sprintf("No 'use' interaction found for '%s'.", objectName), "error")
		return
	}

	// Process based on trigger type
	switch matchedTrigger.TriggerType {
	case "use":
		m.handleUseTrigger(matchedTrigger)
	default:
		m.AppendMessage(fmt.Sprintf("Trigger type '%s' not yet implemented.", matchedTrigger.TriggerType), "error")
	}
}

// handleUseTrigger processes a use trigger based on target type
func (m *model) handleUseTrigger(trigger *triggerView) {
	switch trigger.TargetType {
	case "recipe":
		m.showRecipeMenu(trigger.TargetID)
	case "effect":
		// For now, just inform the user that effects will be applied
		m.AppendMessage("Effect trigger applied! (Effect execution requires ActiveEffect API)", "info")
	default:
		m.AppendMessage(fmt.Sprintf("Target type '%s' not yet implemented.", trigger.TargetType), "error")
	}
}

// showRecipeMenu displays recipes when a use trigger targets recipes
func (m *model) showRecipeMenu(recipeID int) {
	// Get the recipe details
	url := fmt.Sprintf("%s/api/recipes/%d", RESTAPIBase, recipeID)
	resp, err := httpGet(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching recipe: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Could not load recipe.", "error")
		return
	}

	var recipe struct {
		Name         string   `json:"name"`
		DisplayName  string   `json:"display_name"`
		StationTags  []string `json:"station_tags"`
		OutputItem   string   `json:"output_item"`
		OutputAmount int      `json:"output_amount"`
		RequiredClass string  `json:"required_class"`
		RequiredSkill string   `json:"required_skill"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
		m.AppendMessage("Error reading recipe data", "error")
		return
	}

	outputName := recipe.DisplayName
	if outputName == "" {
		outputName = recipe.Name
	}

	output := fmt.Sprintf("=== Crafting: %s ===\n", outputName)
	output += fmt.Sprintf("Output: %s ×%d\n", outputName, recipe.OutputAmount)

	if len(recipe.StationTags) > 0 {
		output += fmt.Sprintf("Station: %s\n", strings.Join(recipe.StationTags, ", "))
	}

	m.AppendMessage(output, "info")

	// Show available ingredients and check requirements
	m.checkCraftingRequirements(recipe)
}

// checkCraftingRequirements shows what ingredients are available and checks class/skill requirements
func (m *model) checkCraftingRequirements(recipe struct {
	Name          string   `json:"name"`
	DisplayName   string   `json:"display_name"`
	StationTags   []string `json:"station_tags"`
	OutputItem    string   `json:"output_item"`
	OutputAmount  int      `json:"output_amount"`
	RequiredClass string   `json:"required_class"`
	RequiredSkill string   `json:"required_skill"`
}) {
	// Get character inventory for ingredient check
	url := fmt.Sprintf("%s/api/characters/%d/equipment", RESTAPIBase, m.currentCharacterID)
	resp, err := httpGet(url)
	if err != nil {
		m.AppendMessage("Error fetching inventory.", "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Could not load inventory.", "error")
		return
	}

	var inventory []struct {
		Name                string `json:"name"`
		EquipmentTemplateID *int   `json:"equipment_template_id"`
		Quantity            int    `json:"quantity"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		m.AppendMessage("Error reading inventory", "error")
		return
	}

	// Check if character has required class
	if recipe.RequiredClass != "" && recipe.RequiredClass != m.characterClass {
		m.AppendMessage(fmt.Sprintf("You don't have the required class: %s", recipe.RequiredClass), "error")
		return
	}

	// Check skill requirements if any
	if recipe.RequiredSkill != "" {
		// For now, just warn - full skill check would need competency API
		m.AppendMessage(fmt.Sprintf("Skill check for '%s' requires skill level", recipe.RequiredSkill), "info")
	}

	// Show available ingredients (simplified)
	m.AppendMessage("Ingredients available in inventory:", "info")
	for _, item := range inventory {
		if item.Quantity > 0 {
			m.AppendMessage(fmt.Sprintf("  • %s (x%d)", item.Name, item.Quantity), "info")
		}
	}

	// Show confirmation prompt
	m.AppendMessage(fmt.Sprintf("Type 'craft %s' to create %s!", recipe.Name, recipe.OutputItem), "success")
}

// getTriggersForRoom fetches triggers for a room from the API
func (m *model) getTriggersForRoom(roomID int) ([]*triggerView, error) {
	url := fmt.Sprintf("%s/api/rooms/%d/triggers", RESTAPIBase, roomID)
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not load room triggers")
	}

	var result struct {
		Triggers []*triggerView `json:"triggers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Triggers, nil
}

// getTriggersForCharacterInventory fetches triggers for character's equipped/owned items
func (m *model) getTriggersForCharacterInventory(charID int) ([]*triggerView, error) {
	// Get character's equipment first
	equipmentURL := fmt.Sprintf("%s/api/characters/%d/equipment", RESTAPIBase, charID)
	resp, err := httpGet(equipmentURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var inventory []struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return nil, err
	}

	// For each item, fetch its triggers
	var allTriggers []*triggerView
	for _, item := range inventory {
		if item.ID == 0 {
			continue
		}
		triggers, err := m.getTriggersForEquipment(item.ID)
		if err == nil {
			allTriggers = append(allTriggers, triggers...)
		}
	}
	return allTriggers, nil
}

// getTriggersForEquipment fetches triggers for a specific equipment item
func (m *model) getTriggersForEquipment(equipmentID int) ([]*triggerView, error) {
	url := fmt.Sprintf("%s/api/equipment/%d/triggers", RESTAPIBase, equipmentID)
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not load equipment triggers")
	}

	var result struct {
		Triggers []*triggerView `json:"triggers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Triggers, nil
}

// triggerView represents the JSON response from the triggers API
type triggerView struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	WorldID     string `json:"world_id"`
	TriggerType string `json:"trigger_type"`
	TargetType  string `json:"target_type"`
	TargetID    int    `json:"target_id"`
	RoomID      *int   `json:"room_id,omitempty"`
	EquipmentID *int   `json:"equipment_id,omitempty"`
	Condition   string `json:"condition,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// handleUseWrapperCommand wraps the use command for the command registry
func (m *model) handleUseWrapperCommand(_ *model, args []string) {
	cmd := "use"
	if len(args) > 0 {
		cmd = fmt.Sprintf("use %s", strings.Join(args, " "))
	}
	m.handleUseCommand(cmd)
}
