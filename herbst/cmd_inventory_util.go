package main

import "fmt"

// mapSlotsToItems maps equipped items to their slot keys.
func mapSlotsToItems(slots []string, items []inventoryItem) map[string]inventoryItem {
	result := make(map[string]inventoryItem)
	for _, item := range items {
		if item.IsEquipped {
			for _, slot := range slots {
				if item.Slot == slot {
					result[slot] = item
					break
				}
			}
		}
	}
	return result
}

// filterUnequipped returns items that are not equipped.
func filterUnequipped(items []inventoryItem) []inventoryItem {
	var backpack []inventoryItem
	for _, item := range items {
		if !item.IsEquipped {
			backpack = append(backpack, item)
		}
	}
	return backpack
}

// maxSlotLabelLen returns the length of the longest slot display name.
func maxSlotLabelLen(slots []string) int {
	maxLen := 0
	for _, slot := range slots {
		label := formatSlotName(slot)
		if len(label) > maxLen {
			maxLen = len(label)
		}
	}
	return maxLen
}

// formatItemStats returns the item name with inline stats.
func formatItemStats(item inventoryItem) string {
	name := item.Name
	if item.DamageDiceCount > 0 && item.DamageDiceSides > 0 {
		bonus := ""
		if item.DamageBonus > 0 {
			bonus = fmt.Sprintf("+%d", item.DamageBonus)
		}
		name += fmt.Sprintf(" (%dd%d%s)", item.DamageDiceCount, item.DamageDiceSides, bonus)
	}
	if item.ArmorRating > 0 {
		name += fmt.Sprintf(" (+%d AC)", item.ArmorRating)
	}
	return name
}