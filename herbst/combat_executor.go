package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SkillBridge bridges content YAML skills to the existing combat system
// This is a transitional bridge during Week 2 migration
type SkillBridge struct{}

// NewSkillBridge creates a new bridge
func NewSkillBridge() *SkillBridge {
	return &SkillBridge{}
}

// SkillDef matches server/content.SkillDef for API responses
type ContentSkill struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Effects     []EffectDef   `json:"effects"`
	Cooldown    int           `json:"cooldown"`
	ManaCost    int           `json:"mana_cost"`
	StaminaCost int           `json:"stamina_cost"`
	Visual      VisualDef     `json:"visual"`
}

// EffectDef represents a skill effect from content
type EffectDef struct {
	Type     string      `json:"type"`
	Target   string      `json:"target"`
	Value    interface{} `json:"value"`
	Duration int         `json:"duration"`
}

// VisualDef from content
type VisualDef struct {
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// fetchContentSkill fetches skill from content API
func (sb *SkillBridge) fetchContentSkill(skillID string) (*ContentSkill, error) {
	resp, err := http.Get(fmt.Sprintf("%s/content/skills/%s", RESTAPIBase, skillID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("skill not found: %s", skillID)
	}

	var skill ContentSkill
	if err := json.NewDecoder(resp.Body).Decode(&skill); err != nil {
		return nil, err
	}
	return &skill, nil
}

// loadClasslessSkillsFromContent fetches all classless skills from content API
func (sb *SkillBridge) loadClasslessSkillsFromContent() ([]ContentSkill, error) {
	resp, err := http.Get(fmt.Sprintf("%s/content/skills/tag/classless", RESTAPIBase))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch classless skills")
	}

	var result struct {
		Skills []ContentSkill `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Skills, nil
}
