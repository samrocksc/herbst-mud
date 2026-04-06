package content

import (
	"fmt"
	"strings"
	"sync"
)

// NPCRegistry manages NPC templates
type NPCRegistry struct {
	mu      sync.RWMutex
	npcs    map[string]*NPCTemplate
}

// NewNPCRegistry creates a new NPC registry
func NewNPCRegistry() *NPCRegistry {
	return &NPCRegistry{
		npcs: make(map[string]*NPCTemplate),
	}
}

// NPCTemplate defines an NPC template from YAML
type NPCTemplate struct {
	ID             string            `yaml:"id" json:"id"`
	TemplateName   string            `yaml:"template_name" json:"template_name"`
	Description    string            `yaml:"description" json:"description"`
	Level          int               `yaml:"level" json:"level"`
	Classification string            `yaml:"classification" json:"classification"`
	Stats          StatsDef          `yaml:"stats" json:"stats"`
	HP             HPCalc            `yaml:"hp" json:"hp"`
	Resources      ResourcesDef      `yaml:"resources,omitempty" json:"resources,omitempty"`
	Skills         NPCSkills         `yaml:"skills,omitempty" json:"skills,omitempty"`
	AI             AIDef             `yaml:"ai" json:"ai"`
	Equipment      EquipmentDef      `yaml:"equipment,omitempty" json:"equipment,omitempty"`
	Loot           []LootEntry       `yaml:"loot,omitempty" json:"loot,omitempty"`
	Visual         VisualDef         `yaml:"visual,omitempty" json:"visual,omitempty"`
	Flags          []string          `yaml:"flags,omitempty" json:"flags,omitempty"`
	Reputation     *ReputationDef    `yaml:"reputation,omitempty" json:"reputation,omitempty"`
	QuestsOffered  []string          `yaml:"quests_offered,omitempty" json:"quests_offered,omitempty"`
}

// StatsDef represents NPC base stats
type StatsDef struct {
	Strength     int `yaml:"strength" json:"strength"`
	Dexterity    int `yaml:"dexterity" json:"dexterity"`
	Constitution int `yaml:"constitution" json:"constitution"`
	Intelligence int `yaml:"intelligence" json:"intelligence"`
	Wisdom       int `yaml:"wisdom" json:"wisdom"`
	Charisma     int `yaml:"charisma" json:"charisma"`
}

// HPCalc defines how to calculate max HP
type HPCalc struct {
	Base          int     `yaml:"base" json:"base"`
	ConMultiplier float64 `yaml:"con_multiplier" json:"con_multiplier"`
	LevelBonus    int     `yaml:"level_bonus" json:"level_bonus"`
}

// ResourcesDef defines stamina/mana
type ResourcesDef struct {
	Stamina ResourceCalc `yaml:"stamina" json:"stamina"`
	Mana    ResourceCalc `yaml:"mana" json:"mana"`
}

// ResourceCalc defines resource calculation
type ResourceCalc struct {
	Base   int `yaml:"base" json:"base"`
	Bonus  int `yaml:"bonus" json:"bonus"`
}

// NPCSkills defines equipped skills
type NPCSkills struct {
	Classless map[string]string `yaml:"classless,omitempty" json:"classless,omitempty"` // key: slot, value: skill_id
	Special   []string          `yaml:"special,omitempty" json:"special,omitempty"`
}

// AIDef defines AI behavior
type AIDef struct {
	Type        string          `yaml:"type" json:"type"`
	Aggression  string          `yaml:"aggression" json:"aggression"`
	Abilities   []AIAbility     `yaml:"abilities,omitempty" json:"abilities,omitempty"`
	Speech      SpeechDef       `yaml:"speech,omitempty" json:"speech,omitempty"`
}

// AIAbility defines an AI skill use
type AIAbility struct {
	SkillID               string  `yaml:"skill_id" json:"skill_id"`
	UseChance             float64 `yaml:"use_chance" json:"use_chance"`
	HealthThresholdBelow  float64 `yaml:"health_threshold_below,omitempty" json:"health_threshold_below,omitempty"`
	HealthThresholdAbove  float64 `yaml:"health_threshold_above,omitempty" json:"health_threshold_above,omitempty"`
}

// SpeechDef defines NPC dialogue
type SpeechDef struct {
	Greeting string   `yaml:"greeting,omitempty" json:"greeting,omitempty"`
	Attack   string   `yaml:"attack,omitempty" json:"attack,omitempty"`
	Defeat   string   `yaml:"defeat,omitempty" json:"defeat,omitempty"`
	Idle     []string `yaml:"idle,omitempty" json:"idle,omitempty"`
}

// EquipmentDef defines equipped gear
type EquipmentDef struct {
	Weapon      string `yaml:"weapon,omitempty" json:"weapon,omitempty"`
	Armor       string `yaml:"armor,omitempty" json:"armor,omitempty"`
	Accessory1  string `yaml:"accessory_1,omitempty" json:"accessory_1,omitempty"`
	Accessory2  string `yaml:"accessory_2,omitempty" json:"accessory_2,omitempty"`
}

// LootEntry defines a loot table entry
type LootEntry struct {
	ItemID   string  `yaml:"item_id" json:"item_id"`
	Chance   float64 `yaml:"chance" json:"chance"`
	CountMin int     `yaml:"count_min" json:"count_min"`
	CountMax int     `yaml:"count_max" json:"count_max"`
}

// ReputationDef defines faction reputation
type ReputationDef struct {
	Faction  string `yaml:"faction" json:"faction"`
	Standing string `yaml:"standing" json:"standing"`
}

// Register adds an NPC template
func (r *NPCRegistry) Register(npc *NPCTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if npc.ID == "" {
		return fmt.Errorf("NPC ID cannot be empty")
	}
	
	id := strings.ToLower(npc.ID)
	if _, exists := r.npcs[id]; exists {
		return fmt.Errorf("NPC '%s' already registered", id)
	}
	
	r.npcs[id] = npc
	return nil
}

// Get retrieves an NPC template
func (r *NPCRegistry) Get(id string) (*NPCTemplate, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	npc, exists := r.npcs[strings.ToLower(id)]
	return npc, exists
}

// GetAll returns all NPCs
func (r *NPCRegistry) GetAll() []*NPCTemplate {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*NPCTemplate, 0, len(r.npcs))
	for _, npc := range r.npcs {
		result = append(result, npc)
	}
	return result
}

// Clear removes all NPCs
func (r *NPCRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.npcs = make(map[string]*NPCTemplate)
}

// Count returns number of registered NPCs
func (r *NPCRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.npcs)
}

// Validate checks NPCs against loaded skills and items
func (r *NPCRegistry) Validate(skills *SkillRegistry, items *ItemRegistry) []ValidationError {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var errors []ValidationError
	
	for id, npc := range r.npcs {
		// Validate skill references
		for slot, skillID := range npc.Skills.Classless {
			if _, exists := skills.Get(skillID); !exists {
				errors = append(errors, ValidationError{
					Type:    "npc",
					ID:      id,
					Field:   fmt.Sprintf("skills.classless.%s", slot),
					Message: fmt.Sprintf("skill '%s' not found", skillID),
				})
			}
		}
		
		for _, skillID := range npc.Skills.Special {
			if _, exists := skills.Get(skillID); !exists {
				errors = append(errors, ValidationError{
					Type:    "npc",
					ID:      id,
					Field:   "skills.special",
					Message: fmt.Sprintf("special skill '%s' not found", skillID),
				})
			}
		}
		
		// Validate equipment references
		if npc.Equipment.Weapon != "" {
			if _, exists := items.Get(npc.Equipment.Weapon); !exists {
				errors = append(errors, ValidationError{
					Type:    "npc",
					ID:      id,
					Field:   "equipment.weapon",
					Message: fmt.Sprintf("item '%s' not found", npc.Equipment.Weapon),
				})
			}
		}
		
		// Validate loot references
		for _, loot := range npc.Loot {
			if _, exists := items.Get(loot.ItemID); !exists {
				errors = append(errors, ValidationError{
					Type:    "npc",
					ID:      id,
					Field:   "loot",
					Message: fmt.Sprintf("loot item '%s' not found", loot.ItemID),
				})
			}
		}
	}
	
	return errors
}

// HasFlag checks if NPC has a flag
func (n *NPCTemplate) HasFlag(flag string) bool {
	for _, f := range n.Flags {
		if strings.ToLower(f) == strings.ToLower(flag) {
			return true
		}
	}
	return false
}

// CalculateMaxHP calculates max HP from template
func (n *NPCTemplate) CalculateMaxHP() int {
	hp := n.HP.Base
	hp += int(float64(n.Stats.Constitution) * n.HP.ConMultiplier)
	hp += n.Level * n.HP.LevelBonus
	return hp
}
