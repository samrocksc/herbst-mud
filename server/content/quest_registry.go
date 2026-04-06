package content

import (
	"fmt"
	"strings"
	"sync"
)

// QuestRegistry manages quest definitions
type QuestRegistry struct {
	mu     sync.RWMutex
	quests map[string]*QuestDef
}

// NewQuestRegistry creates a new quest registry
func NewQuestRegistry() *QuestRegistry {
	return &QuestRegistry{
		quests: make(map[string]*QuestDef),
	}
}

// QuestDef represents a quest definition
type QuestDef struct {
	ID            string         `yaml:"id" json:"id"`
	Name          string         `yaml:"name" json:"name"`
	Description   string         `yaml:"description" json:"description"`
	Type          string         `yaml:"type" json:"type"`
	Difficulty    string         `yaml:"difficulty" json:"difficulty"`
	LevelRequired int            `yaml:"level_required" json:"level_required"`
	Steps         []QuestStep    `yaml:"steps" json:"steps"`
	Rewards       QuestRewards   `yaml:"rewards" json:"rewards"`
	QuestGiver    string         `yaml:"quest_giver" json:"quest_giver"`
	AutoComplete  bool           `yaml:"auto_complete" json:"auto_complete"`
	Repeatable    bool           `yaml:"repeatable" json:"repeatable"`
	RepeatCooldown int         `yaml:"repeat_cooldown,omitempty" json:"repeat_cooldown,omitempty"`
	Flavor        QuestFlavor    `yaml:"flavor,omitempty" json:"flavor,omitempty"`
}

// QuestStep represents a single step in a quest
type QuestStep struct {
	ID          string `yaml:"id" json:"id"`
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
	Target      string `yaml:"target,omitempty" json:"target,omitempty"`
	TargetRoom  string `yaml:"target_room,omitempty" json:"target_room,omitempty"`
	ItemID      string `yaml:"item_id,omitempty" json:"item_id,omitempty"`
	Count       int    `yaml:"count,omitempty" json:"count,omitempty"`
	Order       int    `yaml:"order" json:"order"`
	Requires    string `yaml:"requires,omitempty" json:"requires,omitempty"`
}

// QuestRewards represents quest completion rewards
type QuestRewards struct {
	Experience int           `yaml:"experience" json:"experience"`
	Items      []QuestReward `yaml:"items,omitempty" json:"items,omitempty"`
	Reputation *RepReward  `yaml:"reputation,omitempty" json:"reputation,omitempty"`
}

// QuestReward represents an item reward
type QuestReward struct {
	ItemID string `yaml:"item_id" json:"item_id"`
	Count  int    `yaml:"count" json:"count"`
}

// RepReward represents reputation reward
type RepReward struct {
	Faction string `yaml:"faction" json:"faction"`
	Amount  int    `yaml:"amount" json:"amount"`
}

// QuestFlavor contains flavor text
type QuestFlavor struct {
	Start     string `yaml:"start,omitempty" json:"start,omitempty"`
	Progress  string `yaml:"progress,omitempty" json:"progress,omitempty"`
	Complete  string `yaml:"complete,omitempty" json:"complete,omitempty"`
}

// Register adds a quest to the registry
func (r *QuestRegistry) Register(quest *QuestDef) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if quest.ID == "" {
		return fmt.Errorf("quest ID cannot be empty")
	}
	
	id := strings.ToLower(quest.ID)
	r.quests[id] = quest
	return nil
}

// Get retrieves a quest by ID
func (r *QuestRegistry) Get(id string) (*QuestDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	quest, exists := r.quests[strings.ToLower(id)]
	return quest, exists
}

// GetAll returns all quests
func (r *QuestRegistry) GetAll() []*QuestDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*QuestDef, 0, len(r.quests))
	for _, quest := range r.quests {
		result = append(result, quest)
	}
	return result
}

// GetByDifficulty returns quests of a specific difficulty
func (r *QuestRegistry) GetByDifficulty(difficulty string) []*QuestDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*QuestDef
	diff := strings.ToLower(difficulty)
	
	for _, quest := range r.quests {
		if strings.ToLower(quest.Difficulty) == diff {
			result = append(result, quest)
		}
	}
	return result
}

// GetByType returns quests of a specific type
func (r *QuestRegistry) GetByType(questType string) []*QuestDef {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*QuestDef
	qt := strings.ToLower(questType)
	
	for _, quest := range r.quests {
		if strings.ToLower(quest.Type) == qt {
			result = append(result, quest)
		}
	}
	return result
}

// Clear removes all quests
func (r *QuestRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.quests = make(map[string]*QuestDef)
}

// Count returns number of registered quests
func (r *QuestRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.quests)
}

// Validate checks quests
func (r *QuestRegistry) Validate(skills *SkillRegistry, npcs *NPCRegistry, items *ItemRegistry) []ValidationError {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var errors []ValidationError
	
	for id, quest := range r.quests {
		// Validate quest giver exists
		if quest.QuestGiver != "" && quest.QuestGiver != "system" {
			if _, exists := npcs.Get(quest.QuestGiver); !exists {
				errors = append(errors, ValidationError{
					Type:    "quest",
					ID:      id,
					Field:   "quest_giver",
					Message: fmt.Sprintf("NPC '%s' not found", quest.QuestGiver),
				})
			}
		}
		
		// Validate quest targets exist
		for _, step := range quest.Steps {
			if step.Type == "kill" || step.Type == "talk" {
				if _, exists := npcs.Get(step.Target); !exists {
					errors = append(errors, ValidationError{
						Type:    "quest",
						ID:      id,
						Field:   fmt.Sprintf("steps.%s.target", step.ID),
						Message: fmt.Sprintf("NPC '%s' not found", step.Target),
					})
				}
			}
			if step.Type == "fetch" {
				if _, exists := items.Get(step.ItemID); !exists {
					errors = append(errors, ValidationError{
						Type:    "quest",
						ID:      id,
						Field:   fmt.Sprintf("steps.%s.item_id", step.ID),
						Message: fmt.Sprintf("Item '%s' not found", step.ItemID),
					})
				}
			}
		}
		
		// Validate reward items exist
		for _, reward := range quest.Rewards.Items {
			if _, exists := items.Get(reward.ItemID); !exists {
				errors = append(errors, ValidationError{
					Type:    "quest",
					ID:      id,
					Field:   "rewards.items",
					Message: fmt.Sprintf("Reward item '%s' not found", reward.ItemID),
				})
			}
		}
	}
	
	return errors
}
