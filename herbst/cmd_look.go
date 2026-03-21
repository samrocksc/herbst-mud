package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ============================================================
// LOOK COMMAND — look, examine, search
// ============================================================

// handleLookCommand handles the look/l command
func (m *model) handleLookCommand(cmd string) {
	m.loadRoomItems()
	m.loadRoomCharacters()

	parts := strings.Fields(cmd)
	if len(parts) == 1 {
		// Plain look — show room
		m.AppendMessage(m.formatRoomDisplay(), "info")
		return
	}

	// Strip "at" if present: "look at Gandalf" → "Gandalf"
	var target string
	if len(parts) >= 2 && strings.ToLower(parts[1]) == "at" {
		target = strings.Join(parts[2:], " ")
	} else {
		target = strings.Join(parts[1:], " ")
	}
	target = strings.ToLower(strings.TrimSpace(target))

	m.handleLookAt(target)
}

// handleLookAt handles "look <target>" — items or characters
func (m *model) handleLookAt(target string) {
	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Looking for: '%s'", target), "info")
		m.AppendMessage(fmt.Sprintf("[DEBUG] Room items: %d, Room characters: %d", len(m.roomItems), len(m.roomCharacters)), "info")
	}

	// Room items
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue
		}
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Checking item: '%s'", item.Name), "info")
		}
		if fuzzyWordMatch(item.Name, target) || strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Room characters
	for _, char := range m.roomCharacters {
		charNameLower := strings.ToLower(char.Name)
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Checking character: '%s' (IsNPC: %v)", char.Name, char.IsNPC), "info")
		}
		if fuzzyWordMatch(char.Name, target) || strings.Contains(charNameLower, target) || charNameLower == target {
			if m.debugMode {
				m.AppendMessage(fmt.Sprintf("[DEBUG] Matched character: '%s'", char.Name), "info")
			}
			if char.IsNPC {
				resp, err := httpGet(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
				if err == nil {
					defer resp.Body.Close()
					var npcs []struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						Description string `json:"description"`
						Level       int    `json:"level"`
						Disposition string `json:"disposition"`
					}
					if json.NewDecoder(resp.Body).Decode(&npcs) == nil {
						if m.debugMode {
							m.AppendMessage(fmt.Sprintf("[DEBUG] NPCs from API: %d", len(npcs)), "info")
						}
						for _, npc := range npcs {
							if fuzzyWordMatch(npc.Name, target) || strings.ToLower(npc.Name) == charNameLower || strings.Contains(strings.ToLower(npc.Name), target) {
								m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
									npc.Name, npc.Description, npc.Level, npc.Disposition), "info")
								return
							}
						}
					}
				}
				m.AppendMessage(fmt.Sprintf("[%s]\nAn NPC you can see here.\n\nLevel: %d", char.Name, char.Level), "info")
				return
			}
			m.AppendMessage(fmt.Sprintf("[%s]\nA player adventurer.\n\nLevel: %d", char.Name, char.Level), "info")
			return
		}
	}

	// Hidden items reveal on examine
	if m.currentRoom > 0 {
		resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var allItems []RoomItem
				if json.NewDecoder(resp.Body).Decode(&allItems) == nil {
					for _, item := range allItems {
						if !item.IsVisible && item.RevealCondition != nil {
							revealType, _ := item.RevealCondition["type"].(string)
							revealTarget, _ := item.RevealCondition["target"].(string)
							if revealType == "examine" && strings.ToLower(revealTarget) == target {
								revealResp, err := httpPost(
									fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
									fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, revealTarget, m.characterLevel),
								)
								if err == nil {
									defer revealResp.Body.Close()
									if revealResp.StatusCode == http.StatusOK {
										m.loadRoomItems()
										for _, ri := range m.roomItems {
											if strings.Contains(strings.ToLower(ri.Name), target) || strings.ToLower(ri.Name) == target {
												m.AppendMessage("✨ You discovered something hidden!\n\n", "info")
												m.displayItemDetails(ri)
												return
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	m.AppendMessage(fmt.Sprintf("You don't see any '%s' here.", target), "error")
}

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
		resp, err := httpGet(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var allItems []RoomItem
				if json.NewDecoder(resp.Body).Decode(&allItems) == nil {
					for _, item := range allItems {
						if !item.IsVisible && item.RevealCondition != nil {
							revealType, _ := item.RevealCondition["type"].(string)
							revealTarget, _ := item.RevealCondition["target"].(string)
							if revealType == "examine" && strings.ToLower(revealTarget) == target {
								revealResp, err := httpPost(
									fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
									fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, revealTarget, m.characterLevel),
								)
								if err == nil {
									defer revealResp.Body.Close()
									if revealResp.StatusCode == http.StatusOK {
										m.loadRoomItems()
										for _, ri := range m.roomItems {
											if strings.Contains(strings.ToLower(ri.Name), target) || strings.ToLower(ri.Name) == target {
												m.AppendMessage("✨ You discovered something hidden!\n\n", "info")
												m.displayItemDetails(ri)
												return
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Inventory
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err == nil {
		defer resp.Body.Close()
		var items []RoomItem
		if json.NewDecoder(resp.Body).Decode(&items) == nil {
			for _, item := range items {
				if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
					m.displayItemDetails(item)
					return
				}
			}
		}
	}

	// NPCs
	if m.currentRoom > 0 {
		resp, err := httpGet(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			var npcs []struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Level       int    `json:"level"`
				Disposition string `json:"disposition"`
			}
			if json.NewDecoder(resp.Body).Decode(&npcs) == nil {
				for _, npc := range npcs {
					if strings.Contains(strings.ToLower(npc.Name), target) || strings.ToLower(npc.Name) == target {
						m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
							npc.Name, npc.Description, npc.Level, npc.Disposition), "info")
						return
					}
				}
			}
		}
	}

	m.AppendMessage(fmt.Sprintf("You don't see '%s' here.", target), "error")
}

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

// displayItemDetails shows detailed info about an item
func (m *model) displayItemDetails(item RoomItem) {
	var details strings.Builder

	details.WriteString(fmt.Sprintf("[%s]\n", item.Name))

	desc := item.ExamineDesc
	if desc == "" {
		desc = item.Description
	}
	details.WriteString(desc + "\n")

	if item.ItemType == "weapon" || item.ItemType == "armor" {
		details.WriteString("\n--- Stats ---\n")
		if item.Weight > 0 {
			details.WriteString(fmt.Sprintf("  Weight: %d\n", item.Weight))
		}
		if item.ItemDamage > 0 {
			details.WriteString(fmt.Sprintf("  Damage: %d\n", item.ItemDamage))
		}
		if item.ItemDurability > 0 {
			details.WriteString(fmt.Sprintf("  Durability: %d\n", item.ItemDurability))
		}
		details.WriteString(fmt.Sprintf("  Type: %s\n", item.ItemType))
	}

	if len(item.HiddenDetails) > 0 {
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	}

	m.AppendMessage(details.String(), "info")
}

// fuzzyWordMatch returns true if all words in target appear as substrings in name.
// "grand man" matches "Grand Ol' Man". Case-insensitive.
func fuzzyWordMatch(name, target string) bool {
	nameLower := strings.ToLower(name)
	for _, word := range strings.Fields(strings.ToLower(target)) {
		if !strings.Contains(nameLower, word) {
			return false
		}
	}
	return true
}

// httpGet is a helper for making GET requests
func httpGet(url string) (*http.Response, error) {
	return http.Get(url)
}

// httpPost is a helper for making POST requests with a body
func httpPost(url, body string) (*http.Response, error) {
	return http.Post(url, "application/json", strings.NewReader(body))
}

// ioReadAll is exported for use by other files
func ioReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
