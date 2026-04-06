package content

import (
	"fmt"
	"strings"
	"sync"
)

// SkillRegistry manages skill definitions
type SkillRegistry struct {
	mu     sync.RWMutex
	skills map[string]*SkillDef
}

// NewSkillRegistry creates a new skill registry
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]*SkillDef),
	}
}

// SkillDef represents a skill definition from YAML
type SkillDef struct {
	ID               string            `yaml:"id" json:"id"`
	Name             string            `yaml:"name" json:"name"`
	Description      string            `yaml:"description" json:"description"`
	Type             string            `yaml:"type" json:"type"`
	Tags             []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	LevelRequirement int               `yaml:"level_requirement" json:"level_requirement"`
	ClassRequirement string            `yaml:"class_requirement,omitempty" json:"class_requirement,omitempty"`
	Prerequisites    []SkillPrereq     `yaml:"prerequisites,omitempty" json:"prerequisites,omitempty"`
	Effects          []EffectDef       `yaml:"effects" json:"effects"`
	Cooldown         int               `yaml:"cooldown" json:"cooldown"`
	ManaCost         int               `yaml:"mana_cost" json:"mana_cost"`
	StaminaCost      int               `yaml:"stamina_cost" json:"stamina_cost"`
	HealthCost       int               `yaml:"health_cost,omitempty" json:"health_cost,omitempty"`
	Visual           VisualDef         `yaml:"visual,omitempty" json:"visual,omitempty"`
	AIBehavior       AIBehaviorDef     `yaml:"ai_behavior,omitempty" json:"ai_behavior,omitempty"`
}

// SkillPrereq represents a skill prerequisite
type SkillPrereq struct {
	SkillID string `yaml:"skill_id" json:"skill_id"`
	Level   int    `yaml:"level" json:"level"`
}

// EffectDef represents a skill effect
type EffectDef struct {
	Type     string      `yaml:"type" json:"type"`
	Target   string      `yaml:"target" json:"target"`
	Value    interface{} `yaml:"value" json:"value"`
	Scaling  *ScalingDef `yaml:"scaling,omitempty" json:"scaling,omitempty"`
	Duration int         `yaml:"duration" json:"duration"`
}

// ScalingDef represents stat scaling
type ScalingDef struct {
	Stat  string  `yaml:"stat" json:"stat"`
	Ratio float64 `yaml:"ratio" json:"ratio"`
}

// VisualDef represents visual/sound properties
type VisualDef struct {
	Icon      string `yaml:"icon,omitempty" json:"icon,omitempty"`
	Color     string `yaml:"color,omitempty" json:"color,omitempty"`
	Animation string `yaml:"animation,omitempty" json:"animation,omitempty"`
	Sound     string `yaml:"sound,omitempty" json:"sound,omitempty"`
}

// AIBehaviorDef defines how NPCs use a skill
type AIBehaviorDef struct {
	CanUse          bool              `yaml:"can_use" json:"can_use"`
	UseChance       float64           `yaml:"use_chance" json:"use_chance"`
	HealthThreshold *HealthThreshold  `yaml:"health_threshold,omitempty" json:"health_threshold,omitempty"`
}

// HealthThreshold defines when AI uses skill based on health
type HealthThreshold struct {
	Below float64 `yaml:"below,omitempty" json:"below,omitempty"`
	Above float64 `yaml:"above,omitempty" json:"above,omitempty"`
}

// Register adds a skill to the registry
func (r *SkillRegistry) Register(skill *SkillDef) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if skill.ID == "" {
		return fmt.Errorf("skill ID cannot be empty")
	}
	
	id := strings.ToLower(skill.ID)
	if _, exists := r.skills[id]; exists {
		return fmt.Errorf("skill '%s' already registered", id)
	}
	
	r.skills[id] = skill
	return nil
}

// Get retrieves a skill by ID
func (r *SkillRegistry) Get(id string) (*SkillDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	skill, exists := r.skills[strings.ToLower(id)]
	return skill, exists
}

// GetAll returns all skills
func (r *SkillRegistry) GetAll() []*SkillDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*SkillDef, 0, len(r.skills))
	for _, skill := range r.skills {
		result = append(result, skill)
	}
	return result
}

// GetByTag returns skills matching any tag
func (r *SkillRegistry) GetByTag(tags ...string) []*SkillDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*SkillDef
	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[strings.ToLower(t)] = true
	}
	
	for _, skill := range r.skills {
		for _, skillTag := range skill.Tags {
			if tagSet[strings.ToLower(skillTag)] {
				result = append(result, skill)
				break
			}
		}
	}
	return result
}

// GetByClass returns skills for a class
func (r *SkillRegistry) GetByClass(class string) []*SkillDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*SkillDef
	class = strings.ToLower(class)
	
	for _, skill := range r.skills {
		if skill.ClassRequirement == "" {
			// Classless
			if class == "" || class == "classless" || class == "survivor" {
				result = append(result, skill)
			}
		} else if strings.ToLower(skill.ClassRequirement) == class {
			result = append(result, skill)
		}
	}
	return result
}

// Clear removes all skills
func (r *SkillRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills = make(map[string]*SkillDef)
}

// Count returns number of registered skills
func (r *SkillRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.skills)
}

// Validate checks all skills for consistency
func (r *SkillRegistry) Validate() []ValidationError {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var errors []ValidationError
	
	for id, skill := range r.skills {
		// Check required fields
		if skill.Name == "" {
			errors = append(errors, ValidationError{
				Type:    "skill",
				ID:      id,
				Field:   "name",
				Message: "name is required",
			})
		}
		
		// Validate effects
		if len(skill.Effects) == 0 {
			errors = append(errors, ValidationError{
				Type:    "skill",
				ID:      id,
				Field:   "effects",
				Message: "at least one effect required",
			})
		}
		
		// Validate prerequisites exist
		for _, prereq := range skill.Prerequisites {
			if _, exists := r.skills[strings.ToLower(prereq.SkillID)]; !exists {
				errors = append(errors, ValidationError{
					Type:    "skill",
					ID:      id,
					Field:   "prerequisites",
					Message: fmt.Sprintf("prereq skill '%s' not found", prereq.SkillID),
				})
			}
		}
	}
	
	return errors
}
