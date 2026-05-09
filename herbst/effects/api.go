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

// FireEvent looks up hooks for the given event and NPC template,
// resolves targets, and applies the linked effects.
func (s *Service) FireEvent(eventName string, sourceCharID int, npcTemplateID string, extras map[string]interface{}) {
	hooks := s.GetHooksForEvent(eventName)
	for _, hook := range hooks {
		if !hook.Enabled {
			continue
		}
		if hook.NPCTemplateID != "" && npcTemplateID != "" && hook.NPCTemplateID != npcTemplateID {
			continue
		}
		targets := ResolveTarget(hook.Target, sourceCharID, extras)
		eff, ok := s.GetEffect(hook.EffectID)
		if !ok {
			s.logger.Warn("hook references missing effect", "hook_id", hook.ID, "effect_id", hook.EffectID)
			continue
		}
		dispatchStartMessage(eff.Messages, s.messageBus, sourceCharID)
		for _, targetID := range targets {
			if err := s.ApplyEffect(hook.EffectID, targetID, sourceCharID, 0); err != nil {
				s.logger.Error("apply effect failed", "effect_id", hook.EffectID, "target", targetID, "error", err)
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