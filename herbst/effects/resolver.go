package effects

import (
	"strconv"
)

// ResolveTarget resolves a hook target string into character IDs.
// The extras map carries contextual IDs: attacker_id, killer_id, room_id, etc.
func ResolveTarget(target string, sourceCharID int, extras map[string]interface{}) []int {
	switch target {
	case "self":
		return []int{sourceCharID}
	case "attacker":
		if id := intFromExtras(extras, "attacker_id"); id > 0 {
			return []int{id}
		}
		return nil
	case "killer":
		if id := intFromExtras(extras, "killer_id"); id > 0 {
			return []int{id}
		}
		return nil
	case "room":
		// For "room" target, we need to return all character IDs in the room
		// Room ID is passed in extras
		if roomID := intFromExtras(extras, "room_id"); roomID > 0 {
			// Return a special marker that means "all characters in room"
			// The FireEvent caller should handle this by getting room members
			return []int{roomID | 0x80000000} // Use high bit as marker
		}
		return nil
	case "owner":
		// Owner resolves to the same as self for player characters;
		// for NPC instances, it could resolve to the NPC template owner.
		return []int{sourceCharID}
	default:
		return []int{sourceCharID}
	}
}

// ParseRoomTarget checks if the target is a room marker and returns the room ID
func ParseRoomTarget(targetID int) (int, bool) {
	if targetID&0x80000000 != 0 {
		return targetID & 0x7FFFFFFF, true
	}
	return 0, false
}

func intFromExtras(extras map[string]interface{}, key string) int {
	v, ok := extras[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, _ := strconv.Atoi(n)
		return i
	default:
		return 0
	}
}