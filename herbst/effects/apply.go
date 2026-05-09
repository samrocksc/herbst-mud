package effects

import (
	"fmt"
	"strconv"
)

// ApplyEffect applies an effect definition to a target character.
// It modifies character state via the REST API and dispatches messages.
// depth limits recursive apply_effect chains to 3 levels.
func (s *Service) ApplyEffect(effectID int, targetCharID int, appliedBy int, depth int) error {
	if depth > 3 {
		s.logger.Warn("apply_effect chain depth limit reached", "effect_id", effectID, "target", targetCharID, "depth", depth)
		return nil
	}
	eff, ok := s.GetEffect(effectID)
	if !ok {
		return fmt.Errorf("effect %d not found in cache", effectID)
	}
	switch eff.EffectType {
	case "xp_drain":
		amount := intParam(eff.Parameters, "amount")
		return s.applyXPChange(targetCharID, -amount)
	case "xp_gain":
		amount := intParam(eff.Parameters, "amount")
		return s.applyXPChange(targetCharID, amount)
	case "xp_set":
		amount := intParam(eff.Parameters, "amount")
		return s.applyXPSet(targetCharID, amount)
	case "hp_change":
		amount := intParam(eff.Parameters, "amount")
		return s.applyHPChange(targetCharID, amount)
	case "stamina_change":
		amount := intParam(eff.Parameters, "amount")
		return s.applyStaminaChange(targetCharID, amount)
	case "mana_change":
		amount := intParam(eff.Parameters, "amount")
		return s.applyManaChange(targetCharID, amount)
	case "bind_point_set":
		roomID := intParam(eff.Parameters, "room_id")
		return s.applyBindPointSet(targetCharID, roomID)
	case "teleport":
		roomID := intParam(eff.Parameters, "room_id")
		return s.applyTeleport(targetCharID, roomID)
	case "message":
		text := strParam(eff.Parameters, "text")
		msgType := strParam(eff.Parameters, "message_type")
		s.dispatchMessage(targetCharID, text, msgType)
		return nil
	case "tag_add":
		tag := strParam(eff.Parameters, "tag_name")
		return s.applyTagAdd(targetCharID, tag)
	case "tag_remove":
		tag := strParam(eff.Parameters, "tag_name")
		return s.applyTagRemove(targetCharID, tag)
	case "apply_effect":
		nestedID := intParam(eff.Parameters, "effect_id")
		return s.ApplyEffect(nestedID, targetCharID, appliedBy, depth+1)
	default:
		s.logger.Warn("unknown effect type", "type", eff.EffectType, "effect_id", eff.ID)
		return nil
	}
}

func intParam(params map[string]interface{}, key string) int {
	v, ok := params[key]
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

func strParam(params map[string]interface{}, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}