package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CorpseData represents loot and state of a corpse
type CorpseData struct {
	OriginalName string    `json:"originalName"`
	NPCID        int       `json:"npcId,omitempty"`
	CharacterID  int       `json:"characterId,omitempty"`
	DefeatedAt   string    `json:"defeatedAt"`
	Gold         int       `json:"gold"`
	LootItems    []LootItem `json:"lootItems"` // Equipment details for looting
}

// LootItem represents an item that can be looted from a corpse
type LootItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ItemType    string `json:"itemType"`
	Slot        string `json:"slot"`
}

// CorpseRotDuration is how long corpses stick around before rotting away.
// Players can loot them before they expire.
const CorpseRotDuration = 10 * time.Minute

// generateCorpse creates a corpse item in the room when someone dies
// and transfers all equipment from the defeated character to the corpse
func (m *model) generateCorpse(defeated *RoomCharacter) {
	if m.currentRoom == 0 {
		return
	}

	// Fetch all equipment from the defeated character
	lootItems := m.fetchCharacterEquipment(defeated.ID)

	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Fetched %d items from %s's inventory", len(lootItems), defeated.Name), "info")
	}

	// Build corpse description
	corpseName := fmt.Sprintf("corpse of %s", defeated.Name)
	description := fmt.Sprintf("The lifeless body of %s.", defeated.Name)

	// Create corpse data
	corpseData := CorpseData{
		OriginalName: defeated.Name,
		DefeatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Gold:         0, // TODO: Add gold system
		LootItems:    lootItems,
	}
	if defeated.IsNPC {
		corpseData.NPCID = defeated.ID
	} else {
		corpseData.CharacterID = defeated.ID
	}

	dataJSON, _ := json.Marshal(corpseData)

	// Set expiry time for corpse rotting (GitHub #22)
	expiresAt := time.Now().Add(CorpseRotDuration)

	// Build the request payload
	payload := struct {
		Name              string    `json:"name"`
		Description       string    `json:"description"`
		Slot              string    `json:"slot"`
		ItemType          string    `json:"itemType"`
		IsImmovable       bool      `json:"isImmovable"`
		IsVisible         bool      `json:"isVisible"`
		Weight            int       `json:"weight"`
		Color             string    `json:"color"`
		RoomID            int       `json:"roomId"`
		ExamineDesc       string    `json:"examineDesc"`
		IsContainer       bool      `json:"isContainer"`
		ContainerCapacity int       `json:"containerCapacity"`
		ExpiresAt         *time.Time `json:"expiresAt,omitempty"`
	}{
		Name:              corpseName,
		Description:       description,
		Slot:              "corpse",
		ItemType:          "corpse",
		IsImmovable:       false,
		IsVisible:         true,
		Weight:            200,
		Color:             "purple",
		RoomID:            m.currentRoom,
		ExamineDesc:       string(dataJSON),
		IsContainer:       true,
		ContainerCapacity: 100,
		ExpiresAt:         &expiresAt,
	}

	// Serialize payload
	payloadBytes, _ := json.Marshal(payload)

	// Send request to create corpse
	resp, err := httpPost(fmt.Sprintf("%s/equipment", RESTAPIBase), string(payloadBytes))
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to create corpse: %v", err), "error")
		}
		return
	}
	defer resp.Body.Close()

	var corpseResult struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&corpseResult); err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to parse corpse response: %v", err), "error")
		}
		return
	}

	// Now transfer all equipment to the room (effectively "inside" the corpse)
	for _, item := range lootItems {
		m.transferItemToCorpse(item.ID, m.currentRoom)
	}

	if m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Corpse created for %s with %d items", defeated.Name, len(lootItems)), "info")
	}

	// Show message to players in room
	m.AppendMessage(fmt.Sprintf("☠ The corpse of %s falls to the ground.", defeated.Name), "combat")

	// Reload room items to show the corpse
	m.loadRoomItems()
}

// fetchCharacterEquipment gets all equipment owned by a character
func (m *model) fetchCharacterEquipment(characterID int) []LootItem {
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, characterID))
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to fetch equipment: %v", err), "error")
		}
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var items []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ItemType    string `json:"itemType"`
		Slot        string `json:"slot"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil
	}

	// Convert to LootItem slice
	lootItems := make([]LootItem, 0, len(items))
	for _, item := range items {
		lootItems = append(lootItems, LootItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			ItemType:    item.ItemType,
			Slot:        item.Slot,
		})
	}

	return lootItems
}

// transferItemToCorpse transfers an item from a character to a room (dropping it)
func (m *model) transferItemToCorpse(itemID int, roomID int) {
	// Update the item to be in the room with no owner
	payload := fmt.Sprintf(`{"roomId": %d, "ownerId": null, "isEquipped": false}`, roomID)
	
	resp, err := httpPut(fmt.Sprintf("%s/equipment/%d", RESTAPIBase, itemID), payload)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to transfer item %d: %v", itemID, err), "error")
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && m.debugMode {
		m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to transfer item %d: %d", itemID, resp.StatusCode), "error")
	}
}

// isCorpse checks if an item is a corpse
func isCorpse(item RoomItem) bool {
	return item.ItemType == "corpse" || strings.Contains(strings.ToLower(item.Name), "corpse of")
}

// lootCorpse handles looting a corpse - transfers all items to player
func (m *model) lootCorpse(item RoomItem) {
	if !isCorpse(item) {
		return
	}

	// Parse corpse data from examineDesc
	var data CorpseData
	if err := json.Unmarshal([]byte(item.ExamineDesc), &data); err != nil {
		m.AppendMessage("You can't seem to find anything valuable.", "error")
		return
	}

	// Give gold
	if data.Gold > 0 {
		m.AppendMessage(fmt.Sprintf("💰 You loot %d gold from the corpse of %s!", data.Gold, data.OriginalName), "success")
		// TODO: Add gold to character inventory
	}

	// Transfer all loot items to player
	if len(data.LootItems) > 0 {
		for _, lootItem := range data.LootItems {
			m.transferItemToPlayer(lootItem.ID)
			m.AppendMessage(fmt.Sprintf("📦 You take %s from the corpse.", lootItem.Name), "success")
		}
	} else {
		m.AppendMessage("The corpse has nothing of value.", "info")
	}

	// Remove the empty corpse
	m.emptyCorpse(item.ID)
}

// transferItemToPlayer transfers an equipment item to the player
func (m *model) transferItemToPlayer(itemID int) {
	// Update item ownership to current player
	payload := fmt.Sprintf(`{"ownerId": %d, "roomId": null}`, m.currentCharacterID)
	
	resp, err := httpPut(fmt.Sprintf("%s/equipment/%d", RESTAPIBase, itemID), payload)
	if err != nil {
		if m.debugMode {
			m.AppendMessage(fmt.Sprintf("[DEBUG] Failed to take item %d: %v", itemID, err), "error")
		}
		return
	}
	defer resp.Body.Close()
}

// emptyCorpse removes a looted corpse from the room
func (m *model) emptyCorpse(itemID int) {
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s/equipment/%d", RESTAPIBase, itemID), nil)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 5e9}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()

	// Reload items
	m.loadRoomItems()
}

// getCurrentTimestamp returns current timestamp as string
func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// handleLootCommand handles the "loot" command
func (m *model) handleLootCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		// If no target specified, look for corpses in the room
		m.loadRoomItems()
		
		var corpses []RoomItem
		for _, item := range m.roomItems {
			if isCorpse(item) {
				corpses = append(corpses, item)
			}
		}
		
		if len(corpses) == 0 {
			m.AppendMessage("There are no corpses here to loot.", "error")
			return
		}
		
		if len(corpses) == 1 {
			m.lootCorpse(corpses[0])
			return
		}
		
		// Multiple corpses - ask which one
		m.AppendMessage(fmt.Sprintf("Which corpse? %d found:", len(corpses)), "info")
		for _, c := range corpses {
			m.AppendMessage(fmt.Sprintf("  - %s", c.Name), "info")
		}
		return
	}
	
	// Find specified corpse
	target := strings.ToLower(strings.Join(parts[1:], " "))
	m.loadRoomItems()
	
	for _, item := range m.roomItems {
		if isCorpse(item) && strings.Contains(strings.ToLower(item.Name), target) {
			m.lootCorpse(item)
			return
		}
	}
	
	m.AppendMessage(fmt.Sprintf("You don't see a '%s' here to loot.", target), "error")
}
