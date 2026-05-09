package effects

import "strconv"

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
		if id := intFromExtras(extras, "room_id"); id > 0 {
			return []int{id}
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