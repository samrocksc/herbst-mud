package routes

import (
	"encoding/json"

	"herbst-server/db"
)

// raceView is the JSON shape returned by the API.
type raceView struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	DisplayName      string   `json:"display_name"`
	Description      string   `json:"description,omitempty"`
	StatModifiers    any      `json:"stat_modifiers,omitempty"`
	SkillGrants      []string `json:"skill_grants,omitempty"`
	EquipmentSlots   []string `json:"equipment_slots,omitempty"`
	AbilityModifiers []string `json:"ability_modifiers,omitempty"`
	IsPlayable       bool     `json:"is_playable"`
	Color            string   `json:"color,omitempty"`
}

// parseJSON safely parses a JSON string into a value. Returns nil on error.
func parseJSON[T any](v *T) any {
	if v == nil {
		return nil
	}
	var out any
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return out
}

// raceToView converts a Race ent model to a raceView.
func raceToView(r *db.Race) raceView {
	var statMod any
	if r.StatModifiers != "" {
		_ = json.Unmarshal([]byte(r.StatModifiers), &statMod)
	}

	var skillGrants []string
	if r.SkillGrants != "" {
		_ = json.Unmarshal([]byte(r.SkillGrants), &skillGrants)
	}

	return raceView{
		ID:             r.ID,
		Name:           r.Name,
		DisplayName:    r.DisplayName,
		Description:    r.Description,
		StatModifiers:  statMod,
		SkillGrants:    skillGrants,
		EquipmentSlots: r.EquipmentSlots,
		IsPlayable:     r.IsPlayable,
		Color:          r.Color,
	}
}