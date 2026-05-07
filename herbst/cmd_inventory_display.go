package main

import "github.com/charmbracelet/lipgloss"

// getItemIcon returns an emoji icon based on item type.
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

// getItemRarityColor returns a lipgloss color based on item rarity.
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