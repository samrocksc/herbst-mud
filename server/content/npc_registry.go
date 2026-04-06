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
	ID             string            `json:"id"`
	TemplateName   string            `json:"template_name"`
	Description    string            `json:"description"`
	Level          int               `json:"level"`
	Classification string            `json:"classification"`
	Stats          StatsDef          `json:"stats"`
	HP             HPCalc            `json:"hp"`
	Resources      ResourcesDef      `json:"resources,omitempty"`
	Skills         NPCSkills         `json:"skills,omitempty"`
	AI             AIDef             `json:"ai"`
	Equipment      EquipmentDef      `json:"equipment,omitempty"`
	Loot           []LootEntry       `json:"loot,omitempty"`
	Visual         VisualDef         `json:"visual,omitempty"`
	Flags          []string          `json:"flags,omitempty"`
	Reputation     *ReputationDef    `json:"reputation,omitempty"`
	QuestsOffered  []string          `json:"quests_offered,omitempty"`
}

// StatsDef represents NPC base stats
type StatsDef struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// HPCalc defines how to calculate max HP
type HPCalc struct {
	Base          int     `json:"base"`
	ConMultiplier float64 `json:"con_multiplier"`
	LevelBonus    int     `json:"level_bonus"`
}

// ResourcesDef defines stamina/mana
type ResourcesDef struct {
	Stamina ResourceCalc `json:"stamina"`
	Mana    ResourceCalc `json:"mana"`
}

// ResourceCalc defines resource calculation
type ResourceCalc struct {
	Base   int `json:"base"`
	Bonus  int `json:"bonus"`
}

// NPCSkills defines equipped skills
type NPCSkills struct {
	Classless map[string]string `json:"classless,omitempty"` // key: slot, value: skill_id
	Special   []string          `json:"special,omitempty"`
}

// AIDef defines AI behavior
type AIDef struct {
	Type        string          `json:"type"`
	Aggression  string          `json:"aggression"`
	Abilities   []AIAbility     `json:"abilities,omitempty"`
	Speech      SpeechDef       `json:"speech,omitempty"`
}

// AIAbility defines an AI skill use
type AIAbility struct {
	SkillID               string  `json:"skill_id"`
	UseChance             float64 `json:"use_chance"`
	HealthThresholdBelow  float64 `json:"health_threshold_below,omitempty"`
	HealthThresholdAbove  float64 `json:"health_threshold_above,omitempty"`
}

// SpeechDef defines NPC dialogue
type SpeechDef struct {
	Greeting string   `json:"greeting,omitempty"`
	Attack   string   `json:"attack,omitempty"`
	Defeat   string   `json:"defeat,omitempty"`
	Idle     []string `json:"idle,omitempty"`
}

// EquipmentDef defines equipped gear
type EquipmentDef struct {
	Weapon      string `json:"weapon,omitempty"`
	Armor       string `json:"armor,omitempty"`
	Accessory1  string `json:"accessory_1,omitempty"`
	Accessory2  string `json:"accessory_2,omitempty"`
}

// LootEntry defines a loot table entry
type LootEntry struct {
	ItemID   string  `json:"item_id"`
	Chance   float64 `json:"chance"`
	CountMin int     `json:"count_min"`
	CountMax int     `json:"count_max"`
}

// ReputationDef defines faction reputation
type ReputationDef struct {
	Faction  string `json:"faction"`
	Standing string `json:"standing"`
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
