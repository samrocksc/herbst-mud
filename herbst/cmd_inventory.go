package main

import (
	"fmt"
	"sort"
	"strings"
)

// handleInventoryCommand displays equipped items by slot and backpack items.
func (m *model) handleInventoryCommand() {
	if m.currentCharacterID == 0 {
		m.AppendMessage("You need to be playing to use this command.", "error")
		return
	}

	slots := m.fetchRaceSlots()
	items := m.fetchInventoryItems()

	var output strings.Builder
	output.WriteString("EQUIPMENT\n")

	equippedBySlot := mapSlotsToItems(slots, items)
	maxLabelLen := maxSlotLabelLen(slots)

	for _, slot := range slots {
		label := formatSlotName(slot)
		padded := fmt.Sprintf("%-*s", maxLabelLen, label)
		if item, ok := equippedBySlot[slot]; ok {
			output.WriteString(fmt.Sprintf("  %s %s\n", padded, formatItemStats(item)))
		} else {
			output.WriteString(fmt.Sprintf("  %s [empty]\n", padded))
		}
	}

	backpack := filterUnequipped(items)
	output.WriteString(fmt.Sprintf("\nBACKPACK (%d items)\n", len(backpack)))
	if len(backpack) == 0 {
		output.WriteString("  (empty)\n")
	} else {
		sort.Slice(backpack, func(i, j int) bool {
			return backpack[i].Name < backpack[j].Name
		})
		for _, item := range backpack {
			output.WriteString(fmt.Sprintf("  %s\n", item.Name))
		}
	}

	m.AppendMessage(output.String(), "info")
}