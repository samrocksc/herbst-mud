package content

import (
	"fmt"
	"strings"
	"sync"
)

// ItemRegistry manages item definitions
type ItemRegistry struct {
	mu    sync.RWMutex
	items map[string]*ItemDef
}

// NewItemRegistry creates a new item registry
func NewItemRegistry() *ItemRegistry {
	return &ItemRegistry{
		items: make(map[string]*ItemDef),
	}
}

// ItemDef represents an item definition
type ItemDef struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Type              string            `json:"type"`
	Slot              string            `json:"slot"`
	LevelRequirement  int               `json:"level_requirement"`
	ClassRequirement  []string          `json:"class_requirement,omitempty"`
	StatRequirements  StatReqDef        `json:"stat_requirements,omitempty"`
	Damage            *DamageDef        `json:"damage,omitempty"`
	Armor             int               `json:"armor,omitempty"`
	Stats             StatBonusDef      `json:"stats,omitempty"`
	Effects           []ItemEffectDef   `json:"effects,omitempty"`
	Uses              int               `json:"uses"`
	Stackable         bool              `json:"stackable"`
	MaxStack          int               `json:"max_stack"`
	Weight            float64           `json:"weight"`
	Value             ValueDef          `json:"value"`
	Rarity            string            `json:"rarity"`
	Durability        *DurabilityDef    `json:"durability,omitempty"`
	Visual            VisualDef         `json:"visual"`
	DropTable         *DropTableDef     `json:"drop_table,omitempty"`
}

// StatReqDef represents stat requirements
type StatReqDef struct {
	Strength     int `json:"strength,omitempty"`
	Dexterity    int `json:"dexterity,omitempty"`
	Constitution int `json:"constitution,omitempty"`
	Intelligence int `json:"intelligence,omitempty"`
	Wisdom       int `json:"wisdom,omitempty"`
	Charisma     int `json:"charisma,omitempty"`
}

// DamageDef represents weapon damage
type DamageDef struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Speed float64 `json:"speed"`
	Type  string  `json:"type"`
}

// StatBonusDef represents stat bonuses from equipment
type StatBonusDef struct {
	StatsDef ``
	HP       int `json:"hp,omitempty"`
	Mana     int `json:"mana,omitempty"`
	Stamina  int `json:"stamina,omitempty"`
}

// ItemEffectDef represents item effect
type ItemEffectDef struct {
	Type   string         `json:"type"`
	Effect EffectDataDef  `json:"effect"`
}

// EffectDataDef represents effect data
type EffectDataDef struct {
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
	Duration int     `json:"duration"`
	Chance   float64 `json:"chance"`
}

// ValueDef represents item value
type ValueDef struct {
	Buy  int `json:"buy"`
	Sell int `json:"sell"`
}

// DurabilityDef represents item durability
type DurabilityDef struct {
	Current    int  `json:"current"`
	Max        int  `json:"max"`
	Repairable bool `json:"repairable"`
}

// DropTableDef defines where item drops from
type DropTableDef struct {
	FromNPC  string  `json:"from_npc,omitempty"`
	FromArea string  `json:"from_area,omitempty"`
	Chance   float64 `json:"chance"`
}

// Register adds an item to the registry
func (r *ItemRegistry) Register(item *ItemDef) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}

	id := strings.ToLower(item.ID)
	if _, exists := r.items[id]; exists {
		return fmt.Errorf("item '%s' already registered", id)
	}

	r.items[id] = item
	return nil
}

// Get retrieves an item by ID
func (r *ItemRegistry) Get(id string) (*ItemDef, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[strings.ToLower(id)]
	return item, exists
}

// GetByType returns items of a specific type
func (r *ItemRegistry) GetByType(itemType string) []*ItemDef {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ItemDef
	itemType = strings.ToLower(itemType)

	for _, item := range r.items {
		if strings.ToLower(item.Type) == itemType {
			result = append(result, item)
		}
	}
	return result
}

// GetBySlot returns items equippable in a slot
func (r *ItemRegistry) GetBySlot(slot string) []*ItemDef {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*ItemDef
	slot = strings.ToLower(slot)

	for _, item := range r.items {
		if strings.ToLower(item.Slot) == slot {
			result = append(result, item)
		}
	}
	return result
}

// GetAll returns all items
func (r *ItemRegistry) GetAll() []*ItemDef {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ItemDef, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result
}

// Clear removes all items
func (r *ItemRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]*ItemDef)
}

// Count returns number of registered items
func (r *ItemRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

// Validate checks all items
func (r *ItemRegistry) Validate() []ValidationError {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var errors []ValidationError

	for id, item := range r.items {
		if item.Name == "" {
			errors = append(errors, ValidationError{
				Type:    "item",
				ID:      id,
				Field:   "name",
				Message: "name is required",
			})
		}

		// Validate item-specific constraints
		if item.Type == "weapon" && item.Damage == nil {
			errors = append(errors, ValidationError{
				Type:    "item",
				ID:      id,
				Field:   "damage",
				Message: "weapons require damage definition",
			})
		}

		if item.Stackable && item.MaxStack < 1 {
			errors = append(errors, ValidationError{
				Type:    "item",
				ID:      id,
				Field:   "max_stack",
				Message: "stackable items must have max_stack >= 1",
			})
		}
	}

	return errors
}

// GetDropTable returns items for a specific NPC
func (r *ItemRegistry) GetDropTable(npcID string) []LootEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []LootEntry
	for _, item := range r.items {
		if item.DropTable != nil {
			if strings.ToLower(item.DropTable.FromNPC) == strings.ToLower(npcID) {
				result = append(result, LootEntry{
					ItemID:   item.ID,
					Chance:   item.DropTable.Chance,
					CountMin: 1,
					CountMax: 1,
				})
			}
		}
	}
	return result
}

// CanEquip checks if character meets requirements
func (r *ItemRegistry) CanEquip(item *ItemDef, level int, class string, stats StatsDef) (bool, string) {
	if item.LevelRequirement > level {
		return false, fmt.Sprintf("Requires level %d", item.LevelRequirement)
	}

	if len(item.ClassRequirement) > 0 {
		found := false
		for _, c := range item.ClassRequirement {
			if strings.ToLower(c) == strings.ToLower(class) {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Sprintf("Requires class: %s", strings.Join(item.ClassRequirement, ", "))
		}
	}

	// Check stat requirements
	req := []struct {
		name  string
		req   int
		value int
	}{
		{"Strength", item.StatRequirements.Strength, stats.Strength},
		{"Dexterity", item.StatRequirements.Dexterity, stats.Dexterity},
		{"Constitution", item.StatRequirements.Constitution, stats.Constitution},
		{"Intelligence", item.StatRequirements.Intelligence, stats.Intelligence},
		{"Wisdom", item.StatRequirements.Wisdom, stats.Wisdom},
		{"Charisma", item.StatRequirements.Charisma, stats.Charisma},
	}

	for _, check := range req {
		if check.req > check.value {
			return false, fmt.Sprintf("Requires %d %s", check.req, check.name)
		}
	}

	return true, ""
}
