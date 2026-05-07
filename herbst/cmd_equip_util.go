package main

// equipItemData is the JSON shape for inventory items in equip/unequip.
type equipItemData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slot        string `json:"slot"`
	IsEquipped  bool   `json:"isEquipped"`
	IsTwoHanded bool   `json:"is_two_handed"`
}

// formatSlotName converts slot IDs to display names.
func formatSlotName(slot string) string {
	names := map[string]string{
		"head":         "head",
		"neck":         "neck",
		"chest":        "chest",
		"back":         "back",
		"hands":        "hands",
		"legs":         "legs",
		"feet":         "feet",
		"finger_left":  "left finger",
		"finger_right": "right finger",
		"main_hand":    "main hand",
		"off_hand":     "off hand",
		"tail":         "tail",
		"horn":         "horn",
		"wings":        "wings",
	}
	if name, ok := names[slot]; ok {
		return name
	}
	return slot
}