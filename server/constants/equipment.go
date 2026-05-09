package constants

// SlotCatalog is the master list of all valid equipment slot names.
var SlotCatalog = []string{
	"head", "neck", "chest", "back", "hands", "legs", "feet",
	"finger_left", "finger_right",
	"main_hand", "off_hand",
	"tail", "horn", "wings", "shell",
}

// IsValidSlot returns true if the given slot name exists in SlotCatalog.
func IsValidSlot(slot string) bool {
	for _, s := range SlotCatalog {
		if s == slot {
			return true
		}
	}
	return false
}