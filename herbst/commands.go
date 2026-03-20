package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// GAME COMMAND PROCESSING
// ============================================================

func (m *model) processCommand(cmd string) {
	cmd = strings.TrimSpace(strings.ToLower(cmd))

	if cmd == "" {
		return
	}

	// Handle movement commands
	if m.handleMovement(cmd) {
		return
	}

	// Handle other commands
	switch cmd {
	case "help", "?":
		m.AppendMessage(`Commands:
  n/north, s/south, e/east, w/west - Move
  look/l [target] - Look around (or examine: look <target>, look at <target>)
  ctrl+p - Scroll output up (older messages)
  ctrl+n - Scroll output down (newer messages)
  exits/x - Show exits  
  peer <dir> - Peek at adjacent room
  take/get <item> - Pick up an item
  drop <item> - Drop an item
  inventory/i - Show your inventory
  quests/q - Show your quest log
  whoami - Show your info
  profile/p - Edit character profile
  clear/cls - Clear screen
  quit - Exit game`, "info")
	case "look", "l":
		m.loadRoomItems()
		m.loadRoomCharacters()

		// Parse look arguments: "look", "l", "look <target>", "look at <target>", "l at <target>"
		parts := strings.Fields(cmd)

		var target string
		if len(parts) == 1 {
			// Plain "look" or "l" — show room
			m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nExits: %s%s%s",
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				m.formatExitsWithColor(),
				m.formatRoomItems(),
				m.formatRoomCharacters()), "info")
			return
		}

		// Has arguments — could be "look Frodo" or "look at Gandalf"
		// Strip "at" if present: "look at Gandalf" → target = "Gandalf"
		if len(parts) >= 3 && strings.ToLower(parts[2]) == "at" {
			target = strings.Join(parts[3:], " ")
		} else {
			target = strings.Join(parts[1:], " ")
		}
		target = strings.ToLower(strings.TrimSpace(target))

		// Forward to examine handler with parsed target
		m.handleLookAt(target)
	case "exits", "x":
		m.AppendMessage(fmt.Sprintf("Exits: %s", m.formatExitsWithColor()), "info")
	case "examine", "ex", "inspect":
		m.handleExamineCommand(cmd)
	case "search", "perception":
		// GitHub #12 - Perception check to reveal hidden items
		m.handleSearchCommand(cmd)
	case "whoami":
		// Show character info including level with progress bars
		m.AppendMessage(fmt.Sprintf("=== Character Status ===\nUser: %s (ID: %d)\nRoom: %s\n\n[Level %d - %d XP]\n%s",
			m.currentUserName, m.currentUserID, m.roomName,
			m.characterLevel, m.characterExperience,
			StatusBar(m.characterHP, m.characterMaxHP, m.characterStamina, m.characterMaxStamina, m.characterMana, m.characterMaxMana)), "info")
	case "profile", "p":
		m.screen = ScreenProfile
		m.menuItems = []string{"Edit Gender", "Edit Description", "Back to Game"}
		m.menuCursor = 0
		m.AppendMessage("", "")
	case "peer":
		m.handlePeerCommand(cmd)
	case "debug":
		m.handleDebugCommand(cmd)
		return
	case "clear", "cls":
		// Clear message history (UI-21)
		m.messageHistory = nil
		m.messageTypes = nil
		m.historyOffset = 0
		m.isScrolling = false
		m.inputBuffer = ""
		return
	case "quit", "q":
		m.AppendMessage("Thanks for playing! Goodbye!", "success")
		m.inputBuffer = ""
		return
	default:
		// Check for take/get command (GitHub #89 - Item system)
		if strings.HasPrefix(cmd, "take ") || strings.HasPrefix(cmd, "get ") {
			m.handleTakeCommand(cmd)
			return
		}
		// Check for drop command
		if strings.HasPrefix(cmd, "drop ") {
			m.handleDropCommand(cmd)
			return
		}
		// Check for inventory command
		if cmd == "inventory" || cmd == "i" || cmd == "inv" {
			m.handleInventoryCommand()
			return
		}
		// Check for quests command
		if cmd == "quests" || cmd == "q" || cmd == "quest" {
			m.handleQuestsCommand(cmd)
			return
		}
		// Check for skills command
		if cmd == "skills" {
			m.handleSkillsCommand(cmd)
			return
		}
		// Check for talents command
		if cmd == "talents" {
			m.handleTalentsCommand(cmd)
			return
		}
		// Check for skill equip command
		if strings.HasPrefix(cmd, "skill ") {
			m.handleSkillEquipCommand(cmd)
			return
		}
		// Check for talent equip/unequip/swap commands
		if strings.HasPrefix(cmd, "talent ") {
			m.handleTalentEquipCommand(cmd)
			return
		}
		m.AppendMessage(fmt.Sprintf("Unknown command: %s\nType 'help' for commands", cmd), "error")
	}
}

// handleSkillsCommand displays character skills
func (m *model) handleSkillsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/skills", RESTAPIBase, m.currentCharacterID)
	resp, err := http.Get(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching skills: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load skills", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing skills: %v", err), "error")
		return
	}

	skills, ok := result["skills"].(map[string]interface{})
	if !ok {
		m.AppendMessage("Error: skills data not found", "error")
		return
	}

	// Format skills display
	output := "=== Your Skills ===\n\n"
	for skillName, skillData := range skills {
		data := skillData.(map[string]interface{})
		level := int(data["level"].(float64))
		bonus := data["bonus"].(string)
		output += fmt.Sprintf("%-15s Lv: %2d  %s\n", skillName+":", level, bonus)
	}

	output += "\nSkills are always active and provide passive bonuses."
	m.AppendMessage(output, "info")
}

// handleTalentsCommand displays equipped talents
func (m *model) handleTalentsCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
	resp, err := http.Get(url)
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching talents: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.AppendMessage("Failed to load talents", "error")
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		m.AppendMessage(fmt.Sprintf("Error parsing talents: %v", err), "error")
		return
	}

	// Format talents display with slots
	output := "=== Your Talents ===\n\n"
	slots, ok := result["slots"].([]interface{})
	if !ok {
		// No talents equipped yet
		output += "No talents equipped.\n\n"
		output += "Use: talent equip <talent_id> <slot>\n"
		output += "Slots: 1-4 (quick access keys)\n"
		m.AppendMessage(output, "info")
		return
	}

	emptySlots := 0
	for i := 1; i <= 4; i++ {
		if i < len(slots) && slots[i] != nil {
			slot := slots[i].(map[string]interface{})
			name := slot["name"].(string)
			desc := slot["description"].(string)
			output += fmt.Sprintf("[%d] %s\n     %s\n\n", i, name, desc)
		} else {
			output += fmt.Sprintf("[%d] (empty)\n\n", i)
			emptySlots++
		}
	}

	if emptySlots == 4 {
		output += "No talents equipped. Use 'talent equip <id> <slot>' to equip."
	}

	m.AppendMessage(output, "info")
}

// handleSkillEquipCommand handles skill equip command
func (m *model) handleSkillEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	// Skills are always active, no equip needed
	m.AppendMessage("Skills are always active and cannot be unequipped.\nThey provide passive bonuses based on your skill level.", "info")
}

// handleTalentEquipCommand handles talent equip/unequip/swap commands
func (m *model) handleTalentEquipCommand(cmd string) {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage:\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
		return
	}

	action := parts[1]

	switch action {
	case "equip":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent equip <talent_id> <slot>\nExample: talent equip 1 2", "error")
			return
		}
		talentID := parts[2]
		slot := parts[3]

		// Validate slot is 1-4
		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		// Call API to equip talent
		url := fmt.Sprintf("%s/characters/%d/talents", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"talent_id":%s,"slot":%s}`, talentID, slot)
		resp, err := http.Post(url, "application/json", strings.NewReader(reqBody))
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error equipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			m.AppendMessage("Failed to equip talent", "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("Talent equipped in slot %s", slot), "success")

	case "unequip":
		if len(parts) != 3 {
			m.AppendMessage("Usage: talent unequip <slot>\nExample: talent unequip 2", "error")
			return
		}
		slot := parts[2]

		// Validate slot
		slotNum := 0
		fmt.Sscanf(slot, "%d", &slotNum)
		if slotNum < 1 || slotNum > 4 {
			m.AppendMessage("Slot must be between 1 and 4", "error")
			return
		}

		// Call API to unequip talent
		url := fmt.Sprintf("%s/characters/%d/talents/%s", RESTAPIBase, m.currentCharacterID, slot)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error unequipping talent: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talent unequipped from slot %s", slot), "success")

	case "swap":
		if len(parts) != 4 {
			m.AppendMessage("Usage: talent swap <slot1> <slot2>\nExample: talent swap 1 2", "error")
			return
		}
		slot1 := parts[2]
		slot2 := parts[3]

		// Validate slots
		slot1Num, slot2Num := 0, 0
		fmt.Sscanf(slot1, "%d", &slot1Num)
		fmt.Sscanf(slot2, "%d", &slot2Num)
		if slot1Num < 1 || slot1Num > 4 || slot2Num < 1 || slot2Num > 4 {
			m.AppendMessage("Slots must be between 1 and 4", "error")
			return
		}

		// Call API to swap talents
		url := fmt.Sprintf("%s/characters/%d/talents/swap", RESTAPIBase, m.currentCharacterID)
		reqBody := fmt.Sprintf(`{"slot1":%s,"slot2":%s}`, slot1, slot2)
		req, err := http.NewRequest("PUT", url, strings.NewReader(reqBody))
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error swapping talents: %v", err), "error")
			return
		}
		defer resp.Body.Close()

		m.AppendMessage(fmt.Sprintf("Talents swapped between slot %s and %s", slot1, slot2), "success")

	default:
		m.AppendMessage("Usage:\n  talent - Show talents\n  talent equip <talent_id> <slot>\n  talent unequip <slot>\n  talent swap <slot1> <slot2>", "error")
	}
}

func (m *model) handleMovement(cmd string) bool {
	directionMap := map[string]string{
		"n": "north", "north": "north",
		"s": "south", "south": "south",
		"e": "east", "east": "east",
		"w": "west", "west": "west",
	}

	direction, ok := directionMap[cmd]
	if !ok {
		return false
	}

	// Check if exit exists
	nextRoomID, ok := m.exits[direction]
	if !ok {
		m.AppendMessage("You can't go that way.", "error")
		return true
	}

	// Mark exit as known
	m.knownExits[direction] = true

	// Move to new room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error moving: %v", err), "error")
			return true
		}
		m.currentRoom = room.ID
		m.roomName = room.Name
		m.roomDesc = room.Description
		m.exits = room.Exits

		// Load items and characters for the new room
		m.loadRoomItems()
		m.loadRoomCharacters()

		// Mark new room as visited
		wasVisited := m.visitedRooms[m.currentRoom]
		m.visitedRooms[m.currentRoom] = true

		// Mark new exits as known
		for dir := range m.exits {
			m.knownExits[dir] = true
		}

		// Format room display with items and characters
		roomDisplay := fmt.Sprintf("\n\nExits: %s%s%s",
			m.formatExitsWithColor(),
			m.formatRoomItems(),
			m.formatRoomCharacters())

		if wasVisited {
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(green).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		} else {
			m.AppendMessage(fmt.Sprintf("You go %s.\n\n[%s]\n%s%s",
				direction,
				lipgloss.NewStyle().Bold(true).Foreground(yellow).Render(m.roomName),
				m.roomDesc,
				roomDisplay), "success")
		}
	}

	return true
}

func (m *model) handlePeerCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Usage: peer <direction>\nDirections: north, south, east, west, up, down", "error")
		return
	}
	direction := strings.ToLower(parts[1])

	// Validate direction
	validDirs := map[string]string{"north": "north", "south": "south", "east": "east", "west": "west", "up": "up", "down": "down"}
	dir, ok := validDirs[direction]
	if !ok {
		m.AppendMessage("Invalid direction. Use: north, south, east, west, up, down", "error")
		return
	}

	// Check if exit exists
	nextRoomID, ok := m.exits[dir]
	if !ok {
		m.AppendMessage("You can't peer that way — there's no exit.", "error")
		return
	}

	// Get the room
	if m.client != nil {
		room, err := m.client.Room.Get(context.Background(), nextRoomID)
		if err != nil {
			m.AppendMessage(fmt.Sprintf("Error looking: %v", err), "error")
			return
		}

		m.AppendMessage(fmt.Sprintf("You peer %s...\n\n[%s]\n%s",
			dir,
			lipgloss.NewStyle().Bold(true).Foreground(blue).Render(room.Name),
			room.Description), "info")
	}
}

// handleSearchCommand handles the search/perception command to reveal hidden items
// GitHub #12 - Look System: Hidden Items and Reveal Conditions
func (m *model) handleSearchCommand(cmd string) {
	if m.currentRoom == 0 {
		m.AppendMessage("You can't search here.", "error")
		return
	}

	// Fetch all items (including hidden) for this room
	resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
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
		// Skip already visible items
		if item.IsVisible {
			continue
		}

		// Check if this is a hidden item that can be revealed by perception
		if item.RevealCondition != nil {
			revealType, _ := item.RevealCondition["type"].(string)
			if revealType == "perception_check" {
				// Try to reveal the item with perception check
				// Use character level as skill level for now
				revealResp, err := http.Post(
					fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
					"application/json",
					strings.NewReader(fmt.Sprintf(`{"revealType":"perception_check","skillLevel":%d}`, m.characterLevel)),
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

	// Reload room items to see the newly revealed
	m.loadRoomItems()

	if revealed > 0 {
		m.AppendMessage(fmt.Sprintf("🔍 You search the area carefully...\n\n✨ You discovered %d hidden item(s): %s",
			revealed, strings.Join(found, ", ")), "success")
	} else {
		m.AppendMessage("🔍 You search the area carefully...\n\nYou find nothing of interest.", "info")
	}
}

func (m *model) handleDebugCommand(cmd string) {
	parts := strings.Fields(strings.ToLower(cmd))
	if len(parts) < 2 {
		// Show current debug status
		if m.debugMode {
			m.AppendMessage("Debug mode: ON (Room ID visible in status bar)", "info")
		} else {
			m.AppendMessage("Debug mode: OFF\nUsage: debug on | debug off", "info")
		}
		return
	}

	subCmd := parts[1]
	switch subCmd {
	case "on", "true", "1", "yes":
		m.debugMode = true
		m.AppendMessage("Debug mode: ON (Room ID will show in status bar)", "success")
	case "off", "false", "0", "no":
		m.debugMode = false
		m.AppendMessage("Debug mode: OFF", "info")
	default:
		m.AppendMessage("Usage: debug on | debug off", "error")
	}
}

// ============================================================
// ITEM COMMANDS (GitHub #89 - Item system)
// ============================================================

// handleTakeCommand handles the take/get command
func (m *model) handleTakeCommand(cmd string) {
	// Extract item name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Take what? Usage: take <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")

	// Load room items to find the item
	m.loadRoomItems()

	// Find item by name (case-insensitive partial match)
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

	// Check if immovable
	if targetItem.IsImmovable {
		var colorStyle lipgloss.Style
		if targetItem.Color != "" {
			colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(targetItem.Color))
		} else {
			colorStyle = lipgloss.NewStyle().Foreground(itemColorGold)
		}
		m.AppendMessage(fmt.Sprintf("You can't take the %s. It's firmly fixed in place.", colorStyle.Render(targetItem.Name)), "error")
		return
	}

	// Take the item - move it to player's inventory (roomId = 0 or null)
	url := fmt.Sprintf("%s/equipment/%d", RESTAPIBase, targetItem.ID)
	jsonData, _ := json.Marshal(map[string]interface{}{"roomId": nil})
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error picking up item: %v", err), "error")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
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

// handleDropCommand handles the drop command
func (m *model) handleDropCommand(cmd string) {
	// Extract item name from command
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Drop what? Usage: drop <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")

	// For now, show a message that inventory is not fully implemented
	// This would need player inventory tracking
	m.AppendMessage(fmt.Sprintf("You don't have any %s to drop.", itemName), "error")
}

// handleLookAt handles "look <target>" and "look at <target>" — examines items or characters
func (m *model) handleLookAt(target string) {
	// Check room items first (exact or fuzzy word match)
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue
		}
		if fuzzyWordMatch(item.Name, target) || strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Check room characters (NPCs and players) — fuzzy word match
	for _, char := range m.roomCharacters {
		charNameLower := strings.ToLower(char.Name)
		if fuzzyWordMatch(char.Name, target) || strings.Contains(charNameLower, target) || charNameLower == target {
			if char.IsNPC {
				// Fetch NPC details from API
				resp, err := http.Get(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
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
							if fuzzyWordMatch(npc.Name, target) || strings.ToLower(npc.Name) == charNameLower || strings.Contains(strings.ToLower(npc.Name), target) {
								m.AppendMessage(fmt.Sprintf("[%s]\n%s\n\nLevel: %d\nDisposition: %s",
									npc.Name, npc.Description, npc.Level, npc.Disposition), "info")
								return
							}
						}
					}
				}
				// Fallback if API fails
				m.AppendMessage(fmt.Sprintf("[%s]\nAn NPC you can see here.\n\nLevel: %d", char.Name, char.Level), "info")
				return
			} else {
				// Player character
				m.AppendMessage(fmt.Sprintf("[%s]\nA player adventurer.\n\nLevel: %d", char.Name, char.Level), "info")
				return
			}
		}
	}

	// Try hidden items reveal on examine (GitHub #12)
	if m.currentRoom > 0 {
		resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
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
								revealResp, err := http.Post(
									fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
									"application/json",
									strings.NewReader(fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, revealTarget, m.characterLevel)),
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

	// Not found
	m.AppendMessage(fmt.Sprintf("You don't see any '%s' here.", target), "error")
}

// handleExamineCommand handles the examine/ex/inspect/i command
func (m *model) handleExamineCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Examine what? Usage: examine <item>", "error")
		return
	}

	target := strings.Join(parts[1:], " ")
	target = strings.ToLower(target)

	// First check room items (only visible ones for display)
	for _, item := range m.roomItems {
		if !item.IsVisible {
			continue // Skip hidden items - they'll be handled separately
		}
		if strings.Contains(strings.ToLower(item.Name), target) || strings.ToLower(item.Name) == target {
			m.displayItemDetails(item)
			return
		}
	}

	// Check for hidden items that could be revealed by examining this target
	// (GitHub #12 - Hidden Items and Reveal Conditions)
	if m.currentRoom > 0 {
		// Fetch all items (including hidden) for this room
		resp, err := http.Get(fmt.Sprintf("%s/rooms/%d/equipment?includeHidden=true", RESTAPIBase, m.currentRoom))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var allItems []RoomItem
				if json.NewDecoder(resp.Body).Decode(&allItems) == nil {
					for _, item := range allItems {
						// Check if this is a hidden item that reveals on examine
						if !item.IsVisible && item.RevealCondition != nil {
							revealType, _ := item.RevealCondition["type"].(string)
							revealTarget, _ := item.RevealCondition["target"].(string)
							if revealType == "examine" && strings.ToLower(revealTarget) == target {
								// Try to reveal the item
								revealResp, err := http.Post(
									fmt.Sprintf("%s/equipment/%d/reveal", RESTAPIBase, item.ID),
									"application/json",
									strings.NewReader(fmt.Sprintf(`{"revealType":"examine","target":"%s","skillLevel":%d}`, revealTarget, m.characterLevel)),
								)
								if err == nil {
									defer revealResp.Body.Close()
									if revealResp.StatusCode == http.StatusOK {
										// Item revealed! Reload room items and try again
										m.loadRoomItems()
										// Re-check with now-visible item
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

	// Then check inventory
	resp, err := http.Get(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
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

	// Check if it's an NPC
	if m.currentRoom > 0 {
		resp, err := http.Get(fmt.Sprintf("%s/npc?roomId=%d", RESTAPIBase, m.currentRoom))
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

// displayItemDetails shows detailed info about an item
func (m *model) displayItemDetails(item RoomItem) {
	var details strings.Builder

	// Title with color if applicable
	if item.Color != "" {
		details.WriteString(fmt.Sprintf("[%s]\n", item.Name))
	} else {
		details.WriteString(fmt.Sprintf("[%s]\n", item.Name))
	}

	// Use examine description if available, otherwise fall back to description
	desc := item.ExamineDesc
	if desc == "" {
		desc = item.Description
	}
	details.WriteString(desc + "\n")

	// Show stats if it's equipment
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

	// Show hidden details if player has high enough examine skill
	// For now, we'll show all hidden details (skill check deferred)
	if len(item.HiddenDetails) > 0 && item.HiddenThreshold > 0 {
		// TODO: Fetch player's examine skill and compare to threshold
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	} else if len(item.HiddenDetails) > 0 {
		details.WriteString("\n--- You Notice ---\n")
		for _, hd := range item.HiddenDetails {
			details.WriteString(fmt.Sprintf("  %s\n", hd.Text))
		}
	}

	m.AppendMessage(details.String(), "info")
}

// handleInventoryCommand handles the inventory/i command
func (m *model) handleInventoryCommand() {
	// Fetch player's inventory from API
	resp, err := http.Get(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching inventory: %v", err), "error")
		return
	}
	defer resp.Body.Close()

	// Parse inventory items
	var rawItems []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ItemType    string `json:"itemType"`
		IsEquipped  bool   `json:"isEquipped"`
		Rarity      string `json:"rarity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawItems); err != nil {
		m.AppendMessage("You aren't carrying anything.", "info")
		return
	}

	// Convert to typed items
	items := make([]inventoryItem, len(rawItems))
	for i, raw := range rawItems {
		items[i] = inventoryItem(raw)
	}

	if len(items) == 0 {
		m.AppendMessage("Your pockets are empty. Time to loot some stuff!", "info")
		return
	}

	// Format inventory display with icons and styling
	var inv strings.Builder
	inv.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("🎒 INVENTORY"))
	inv.WriteString("\n")
	inv.WriteString(strings.Repeat("─", 30))
	inv.WriteString("\n\n")

	// Group items by type for better organization
	typeGroups := make(map[string][]inventoryItem)

	for _, item := range items {
		typeGroups[item.ItemType] = append(typeGroups[item.ItemType], item)
	}

	// Display items grouped by type with icons
	for itemType, groupItems := range typeGroups {
		icon := getItemIcon(itemType)
		typeLabel := strings.ToUpper(itemType)
		inv.WriteString(lipgloss.NewStyle().Bold(true).Foreground(cyan).Render(fmt.Sprintf("%s %s", icon, typeLabel)))
		inv.WriteString("\n")

		for _, invItem := range groupItems {
			rarityColor := getItemRarityColor(invItem.Rarity)
			itemStyle := lipgloss.NewStyle().Foreground(rarityColor)

			equipped := ""
			if invItem.IsEquipped {
				equipped = " " + lipgloss.NewStyle().Bold(true).Foreground(green).Render("⚡ equipped")
			}

			inv.WriteString(fmt.Sprintf("  %s %s%s\n", icon, itemStyle.Render(invItem.Name), equipped))
			if invItem.Description != "" {
				inv.WriteString(fmt.Sprintf("     %s\n", invItem.Description))
			}
		}
		inv.WriteString("\n")
	}

	m.AppendMessage(inv.String(), "info")
}

// handleQuestsCommand handles the quests/q command to display quest tracker
func (m *model) handleQuestsCommand(cmd string) {
	// Show placeholder when no character is selected
	if m.currentCharacterID == 0 {
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Fetch quests from API
	// For now, we'll return mock data until the full quest system is implemented
	// In production, this would call: GET /characters/:id/quests
	resp, err := http.Get(fmt.Sprintf("%s/characters/%d/quests", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		// Network error - show placeholder (expected in dev mode without server)
		m.displayQuestTrackerPlaceholder()
		return
	}
	defer resp.Body.Close()

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		// If no quests endpoint exists yet, show placeholder message
		// This allows the feature to work before the full quest system is built
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Parse quest response
	var questResp struct {
		Quests []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Status      string `json:"status"`
			Objectives  []struct {
				Description string `json:"description"`
				Current     int    `json:"current"`
				Total       int    `json:"total"`
			} `json:"objectives"`
			Giver  string `json:"giver"`
			Rewards string `json:"rewards"`
		} `json:"quests"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&questResp); err != nil || len(questResp.Quests) == 0 {
		// No quests available - show placeholder
		m.displayQuestTrackerPlaceholder()
		return
	}

	// Format quest tracker display with Lip Gloss styling
	var quests strings.Builder

	// Title
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	activeCount := 0
	availableCount := 0
	completedCount := 0

	for _, quest := range questResp.Quests {
		status := quest.Status
		switch status {
		case "in_progress":
			activeCount++
		case "available":
			availableCount++
		case "completed":
			completedCount++
		}

		// Quest box with styled border
		quests.WriteString(questBoxStyle.Render("") + "\n")

		// Quest name with status color
		statusColor := questAvailableStyle
		statusText := "Available"
		if status == "in_progress" {
			statusColor = questProgressStyle
			statusText = "In Progress"
		} else if status == "completed" {
			statusColor = questCompletedStyle
			statusText = "Completed"
		}

		quests.WriteString(fmt.Sprintf("  %s [%s]\n", questTitleStyle.Render(quest.Name), statusColor.Render(statusText)))

		// Description
		if quest.Description != "" {
			quests.WriteString(fmt.Sprintf("    %s\n", quest.Description))
		}

		// Objectives with progress
		if len(quest.Objectives) > 0 {
			quests.WriteString("\n  Objectives:\n")
			for _, obj := range quest.Objectives {
				progress := fmt.Sprintf("%d/%d", obj.Current, obj.Total)
				if obj.Current >= obj.Total {
					quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", obj.Description, questCompletedStyle.Render("("+progress+")")))
				} else {
					quests.WriteString(fmt.Sprintf("    ○ %s %s\n", obj.Description, questProgressStyle.Render("("+progress+")")))
				}
			}
		}

		// Giver
		if quest.Giver != "" {
			quests.WriteString(fmt.Sprintf("\n  Giver: %s\n", quest.Giver))
		}

		// Rewards
		if quest.Rewards != "" {
			quests.WriteString(fmt.Sprintf("  Reward: %s\n", quest.Rewards))
		}

		quests.WriteString("\n")
	}

	// Summary footer
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString(fmt.Sprintf("  Active: %d  |  Available: %d  |  Completed: %d\n",
		activeCount, availableCount, completedCount))
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	m.AppendMessage(quests.String(), "info")
}

// displayQuestTrackerPlaceholder shows a placeholder quest tracker
// when no quests are available (before full quest system is implemented)
func (m *model) displayQuestTrackerPlaceholder() {
	var quests strings.Builder

	// Title with Lip Gloss styling
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n")
	quests.WriteString(questTitleStyle.Render("  🤺  QUEST LOG  🤺") + "\n")
	quests.WriteString(questTitleStyle.Render("═══════════════════════════════════════") + "\n\n")

	// Placeholder quests from the quest system spec
	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Prove Yourself"),
		questProgressStyle.Render("In Progress")))

	quests.WriteString("    The Scrapyard ain't for the weak. Kill 3 Scrap Rats\n")
	quests.WriteString("    and I'll let you into New Venice proper.\n\n")

	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Kill Scrap Rat", questProgressStyle.Render("(2/3)")))
	quests.WriteString(fmt.Sprintf("    ✓ %s %s\n", "Find Guard Marco at Foggy Gate", questCompletedStyle.Render("(done)")))

	quests.WriteString("\n  Giver: Guard Marco\n")
	quests.WriteString("  Reward: 10 coins\n\n")

	// Second placeholder quest
	quests.WriteString(questBoxStyle.Render("") + "\n")
	quests.WriteString(fmt.Sprintf("  %s [%s]\n",
		questTitleStyle.Render("Ooze Samples"),
		questAvailableStyle.Render("Available")))

	quests.WriteString("    Jane needs Ooze samples for her research.\n")
	quests.WriteString("    The Leaking Pipes have plenty.\n\n")

	quests.WriteString("  Objectives:\n")
	quests.WriteString(fmt.Sprintf("    ○ %s %s\n", "Collect glowing goo", questProgressStyle.Render("(0/5)")))

	quests.WriteString("\n  Giver: Scavenger Jane\n")
	quests.WriteString("  Reward: repair_kit, scavenge skill\n\n")

	// Summary footer
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")
	quests.WriteString("  Active: 1  |  Available: 1  |  Completed: 0\n")
	quests.WriteString(questTitleStyle.Render("───────────────────────────────────────") + "\n")

	quests.WriteString("\n" + infoStyle.Render("  Use 'quest <name>' for details, 'accept <quest>' to begin."))

	m.AppendMessage(quests.String(), "info")
}
