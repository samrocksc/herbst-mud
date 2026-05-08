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
	Tags             []string `json:"tags,omitempty"`
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

	var tagNames []string
	if r.Edges.Tags != nil {
		for _, t := range r.Edges.Tags {
			tagNames = append(tagNames, t.Name)
		}
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
		Tags:           tagNames,
	}
}