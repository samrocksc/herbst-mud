package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// INVENTORY — take, drop, inventory
// ============================================================

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
		m.AppendMessage(fmt.Sprintf("You can't take the %s. It's firmly fixed in place.", colorStyle.Render(targetItem.Name)), "error")
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

func (m *model) handleDropCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) < 2 {
		m.AppendMessage("Drop what? Usage: drop <item name>", "error")
		return
	}
	itemName := strings.Join(parts[1:], " ")
	m.AppendMessage(fmt.Sprintf("You don't have any %s to drop.", itemName), "error")
}

func (m *model) handleInventoryCommand() {
	resp, err := httpGet(fmt.Sprintf("%s/equipment?ownerId=%d", RESTAPIBase, m.currentCharacterID))
	if err != nil {
		m.AppendMessage(fmt.Sprintf("Error fetching inventory: %v", err), "error")
		return
	}
	defer resp.Body.Close()

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

	if len(rawItems) == 0 {
		m.AppendMessage("Your pockets are empty. Time to loot some stuff!", "info")
		return
	}

	items := make([]InventoryItem, len(rawItems))
	for i, raw := range rawItems {
		items[i] = InventoryItem(raw)
	}

	var inv strings.Builder
	inv.WriteString(lipgloss.NewStyle().Bold(true).Foreground(pink).Render("🎒 INVENTORY"))
	inv.WriteString("\n")
	inv.WriteString(strings.Repeat("─", 30))
	inv.WriteString("\n\n")

	typeGroups := make(map[string][]InventoryItem)
	for _, item := range items {
		typeGroups[item.ItemType] = append(typeGroups[item.ItemType], item)
	}

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

// getItemIcon returns an emoji icon based on item type
func getItemIcon(itemType string) string {
	switch itemType {
	case "weapon":
		return "⚔️"
	case "armor":
		return "🛡️"
	case "potion":
		return "🧪"
	case "food":
		return "🍖"
	case "scroll":
		return "📜"
	case "key":
		return "🔑"
	case "treasure":
		return "💎"
	case "quest":
		return "📋"
	default:
		return "📦"
	}
}

// getItemRarityColor returns a lipgloss color based on item rarity
func getItemRarityColor(rarity string) lipgloss.Color {
	switch rarity {
	case "rare":
		return lipgloss.Color("51")
	case "epic":
		return lipgloss.Color("201")
	case "legendary":
		return lipgloss.Color("220")
	default:
		return lipgloss.Color("white")
	}
}
