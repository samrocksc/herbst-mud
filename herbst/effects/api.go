package effects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// applyXPChange modifies a character's XP via the REST API.
func (s *Service) applyXPChange(charID int, delta int) error {
	return s.patchCharacter(charID, map[string]interface{}{"xp_delta": delta})
}

// applyXPSet sets a character's XP to an absolute value via the REST API.
func (s *Service) applyXPSet(charID int, value int) error {
	return s.patchCharacter(charID, map[string]interface{}{"xp": value})
}

// applyHPChange modifies a character's HP via the REST API.
func (s *Service) applyHPChange(charID int, delta int) error {
	return s.patchCharacter(charID, map[string]interface{}{"hp_delta": delta})
}

// applyStaminaChange modifies a character's stamina via the REST API.
func (s *Service) applyStaminaChange(charID int, delta int) error {
	return s.patchCharacter(charID, map[string]interface{}{"stamina_delta": delta})
}

// applyManaChange modifies a character's mana via the REST API.
func (s *Service) applyManaChange(charID int, delta int) error {
	return s.patchCharacter(charID, map[string]interface{}{"mana_delta": delta})
}

// applyBindPointSet changes a character's respawn room via the REST API.
func (s *Service) applyBindPointSet(charID int, roomID int) error {
	return s.patchCharacter(charID, map[string]interface{}{"respawn_room_id": roomID})
}

// applyTeleport moves a character to a new room via the REST API.
func (s *Service) applyTeleport(charID int, roomID int) error {
	return s.patchCharacter(charID, map[string]interface{}{"current_room_id": roomID})
}

// applyTagAdd adds a tag to a character via the REST API.
func (s *Service) applyTagAdd(charID int, tag string) error {
	url := fmt.Sprintf("%s/api/characters/%d/tags", s.restBase, charID)
	body := map[string]interface{}{"tag": tag, "source": "effect"}
	return s.postJSON(url, body)
}

// applyTagRemove removes a tag from a character via the REST API.
func (s *Service) applyTagRemove(charID int, tag string) error {
	url := fmt.Sprintf("%s/api/characters/%d/tags/%s", s.restBase, charID, tag)
	return s.deleteJSON(url)
}

// patchCharacter sends a PATCH request to update a character.
func (s *Service) patchCharacter(charID int, fields map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/characters/%d", s.restBase, charID)
	return s.patchJSON(url, fields)
}

func (s *Service) patchJSON(url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("PATCH %s: %d", url, resp.StatusCode)
	}
	return nil
}

func (s *Service) postJSON(url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("POST %s: %d", url, resp.StatusCode)
	}
	return nil
}

func (s *Service) deleteJSON(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("DELETE %s: %d", url, resp.StatusCode)
	}
	return nil
}

// MessageBus allows the effects service to send messages to characters.
// The game model registers a handler via RegisterMessageHandler.
type MessageBus struct {
	handler func(charID int, text string, msgType string)
}

func NewMessageBus() *MessageBus {
	return &MessageBus{}
}

func (mb *MessageBus) Send(charID int, text string, msgType string) {
	if mb.handler != nil {
		mb.handler(charID, text, msgType)
	}
}

func (mb *MessageBus) RegisterHandler(fn func(charID int, text string, msgType string)) {
	mb.handler = fn
}

// dispatchMessage sends a text message to a character via the game's
// message system. This is called when a message effect fires.
func dispatchMessage(s *Service, charID int, text string, msgType string) {
	if text == "" {
		return
	}
	// Game chat commands (say, yell, shout, whisper) should be passed through as-is
	// since they represent game commands that will be sent to the server.
	// Other valid message types for in-game display: info, success, error, damage, heal
	if msgType == "" {
		msgType = "info"
	}
	// Messages are dispatched through the game model's message system.
	// The game model will pick up pending messages via the MessageBus.
	s.messageBus.Send(charID, text, msgType)
}

// FireEvent looks up hooks for the given event and NPC template,
// resolves targets, and applies the linked effects.
func (s *Service) FireEvent(eventName string, sourceCharID int, npcTemplateID string, extras map[string]interface{}) {
	s.logger.Debug("FireEvent", "event", eventName, "source_char_id", sourceCharID, "npc_template_id", npcTemplateID, "extras", extras)
	hooks := s.GetHooksForEvent(eventName)
	s.logger.Debug("FireEvent", "hooks_found", len(hooks))
	for _, hook := range hooks {
		if !hook.Enabled {
			s.logger.Debug("FireEvent: skipping disabled hook", "hook_id", hook.ID)
			continue
		}
		if hook.NPCTemplateID != "" && npcTemplateID != "" && hook.NPCTemplateID != npcTemplateID {
			s.logger.Debug("FireEvent: NPC template mismatch", "hook_template", hook.NPCTemplateID, "event_template", npcTemplateID)
			continue
		}
		targets := ResolveTarget(hook.Target, sourceCharID, extras)
		s.logger.Debug("FireEvent: resolved targets", "hook_id", hook.ID, "target", hook.Target, "targets", targets)
		eff, ok := s.GetEffect(hook.EffectID)
		if !ok {
			s.logger.Warn("hook references missing effect", "hook_id", hook.ID, "effect_id", hook.EffectID)
			continue
		}
		s.logger.Debug("FireEvent: applying effect", "hook_id", hook.ID, "effect_id", hook.EffectID, "effect_type", eff.EffectType)
		dispatchStartMessage(eff.Messages, s.messageBus, sourceCharID)
		for _, targetID := range targets {
			// Check if this is a room target (special marker)
			if roomID, isRoom := ParseRoomTarget(targetID); isRoom {
				s.logger.Debug("FireEvent: room target, will dispatch to all room members", "room_id", roomID)
				// For room target, the effect type determines how we handle it
				if eff.EffectType == "message" {
					// For message effects on room target, send to all characters in room
					// The actual room member lookup should be done by the game server
					// For now, we dispatch to sourceCharID as a fallback
					s.logger.Debug("FireEvent: message effect on room, dispatching to source", "source_char_id", sourceCharID)
					dispatchMessage(s, sourceCharID, strParam(eff.Parameters, "text"), strParam(eff.Parameters, "message_type"))
				} else {
					s.logger.Debug("FireEvent: non-message room target, applying to source", "target_id", sourceCharID)
					if err := s.ApplyEffect(hook.EffectID, sourceCharID, sourceCharID, 0); err != nil {
						s.logger.Error("apply effect failed", "effect_id", hook.EffectID, "target", sourceCharID, "error", err)
					}
				}
			} else {
				if err := s.ApplyEffect(hook.EffectID, targetID, sourceCharID, 0); err != nil {
					s.logger.Error("apply effect failed", "effect_id", hook.EffectID, "target", targetID, "error", err)
				}
			}
		}
	}
}

// parseInt extracts an int from extras map.
func parseInt(v interface{}) int {
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