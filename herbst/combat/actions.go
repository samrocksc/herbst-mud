package combat

// ActionType defines how an action executes over time
type ActionType string

const (
	// ActionInstant executes immediately in the current tick
	// These actions have their full effect applied immediately
	ActionInstant ActionType = "INSTANT"
	
	// ActionChannel requires the participant to remain active for multiple ticks
	// Effect is applied each tick during the channel
	ActionChannel ActionType = "CHANNEL"
	
	// ActionCharge requires buildup ticks before the effect executes
	// Effect is applied only after all charge ticks complete
	ActionCharge ActionType = "CHARGE"
)

// ActionCategory groups related actions for organization
type ActionCategory string

const (
	CategoryBasic    ActionCategory = "BASIC"    // Basic actions (attack, defend, flee)
	CategorySkill    ActionCategory = "SKILL"    // Skill actions (class-specific)
	CategoryTalent   ActionCategory = "TALENT"   // Talent actions (learned abilities)
	CategoryItem     ActionCategory = "ITEM"     // Item usage actions
	CategoryReaction ActionCategory = "REACTION" // Reactive actions (parry, counter)
)

// ActionDefinition defines an action's properties and tick costs
type ActionDefinition struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    ActionCategory `json:"category"`
	Type        ActionType     `json:"type"`
	
	// Tick costs
	TickCost     int `json:"tickCost"`     // Ticks until action executes (for Instant/Charge)
	ChannelTicks int `json:"channelTicks"` // Total ticks for channeling (for Channel)
	ChargeTicks  int `json:"chargeTicks"`  // Buildup ticks before execution (for Charge)
	
	// Resource costs
	ManaCost    int `json:"manaCost"`
	StaminaCost int `json:"staminaCost"`
	
	// Combat effect
	BaseDamage int `json:"baseDamage"`
	BaseHeal   int `json:"baseHeal"`
	
	// Action properties
	Cooldown      int  `json:"cooldown"`      // Ticks before can use again
	MaxTargets     int  `json:"maxTargets"`    // Maximum targets (1 for single, -1 for AoE)
	Range          int  `json:"range"`         // 0 = melee, 1+ = ranged tiles
	RequiresTarget bool `json:"requiresTarget"`
	IsDefensive   bool `json:"isDefensive"`
	IsOffensive   bool `json:"isOffensive"`
	IsSupport     bool `json:"isSupport"`
	
	// Priority for execution order (lower = earlier)
	Priority ActionPriority `json:"priority"`
	
	// Tags for special handling
	Tags []string `json:"tags"`
}

// BasicActions defines the fundamental combat actions
var BasicActions = map[string]*ActionDefinition{
	"attack": {
		ID:             "attack",
		Name:           "Attack",
		Description:    "A basic attack dealing damage based on your weapon.",
		Category:       CategoryBasic,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    0,
		BaseDamage:     10,
		Cooldown:       0,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PriorityNormal,
		Tags:           []string{"physical", "melee"},
	},
	"defend": {
		ID:             "defend",
		Name:           "Defend",
		Description:    "Take a defensive stance, reducing incoming damage by 50% this tick.",
		Category:       CategoryBasic,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    5,
		BaseDamage:     0,
		BaseHeal:       0,
		Cooldown:       0,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    true,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityFast,
		Tags:           []string{"defensive", "stance"},
	},
	"flee": {
		ID:             "flee",
		Name:           "Flee",
		Description:    "Attempt to escape from combat.",
		Category:       CategoryBasic,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    10,
		BaseDamage:     0,
		Cooldown:       0,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityLast,
		Tags:           []string{"escape"},
	},
	"item": {
		ID:             "item",
		Name:           "Use Item",
		Description:    "Use an item from your inventory.",
		Category:       CategoryBasic,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    0,
		BaseDamage:     0,
		Cooldown:       0,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      true,
		Priority:       PriorityNormal,
		Tags:           []string{"item", "inventory"},
	},
	"wait": {
		ID:             "wait",
		Name:           "Wait",
		Description:    "Skip your turn and recover 10 stamina.",
		Category:       CategoryBasic,
		Type:           ActionInstant,
		TickCost:       0,
		ManaCost:       0,
		StaminaCost:    0,
		BaseDamage:     0,
		BaseHeal:       0,
		Cooldown:       0,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityLast,
		Tags:           []string{"skip"},
	},
}

// SkillActions defines class-based combat skills
var SkillActions = map[string]*ActionDefinition{
	// Warrior skills
	"slash": {
		ID:             "slash",
		Name:           "Slash",
		Description:    "A quick slash attack dealing 1.5x weapon damage.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    8,
		BaseDamage:     15,
		Cooldown:       2,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PriorityNormal,
		Tags:           []string{"physical", "melee", "warrior"},
	},
	"heavy_strike": {
		ID:             "heavy_strike",
		Name:           "Heavy Strike",
		Description:    "A powerful strike that takes an extra tick to wind up, dealing 2.5x damage.",
		Category:       CategorySkill,
		Type:           ActionCharge,
		TickCost:       1,
		ChargeTicks:    1,
		ManaCost:       0,
		StaminaCost:    15,
		BaseDamage:     25,
		Cooldown:       3,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PrioritySlow,
		Tags:           []string{"physical", "melee", "warrior", "charge"},
	},
	"parry": {
		ID:             "parry",
		Name:           "Parry",
		Description:    "Reactively block the next incoming attack and counter for half damage.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       0,
		ManaCost:       0,
		StaminaCost:    12,
		BaseDamage:     5,
		Cooldown:       4,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    true,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityImmediate,
		Tags:           []string{"defensive", "reaction", "warrior"},
	},
	"shield_bash": {
		ID:             "shield_bash",
		Name:           "Shield Bash",
		Description:    "Strike with your shield, stunning the target for 1 tick.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    10,
		BaseDamage:     8,
		Cooldown:       3,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PriorityFast,
		Tags:           []string{"physical", "melee", "warrior", "stun"},
	},
	
	// Mage skills
	"fireball": {
		ID:             "fireball",
		Name:           "Fireball",
		Description:    "Hurl a ball of fire at the enemy, dealing fire damage.",
		Category:       CategorySkill,
		Type:           ActionCharge,
		TickCost:       1,
		ChargeTicks:    2,
		ManaCost:       20,
		StaminaCost:    0,
		BaseDamage:     35,
		Cooldown:       4,
		MaxTargets:     1,
		Range:          3,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PrioritySlow,
		Tags:           []string{"magic", "fire", "mage", "charge"},
	},
	"ice_shield": {
		ID:             "ice_shield",
		Name:           "Ice Shield",
		Description:    "Create a protective barrier of ice that absorbs damage.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       15,
		StaminaCost:    0,
		BaseDamage:     0,
		Cooldown:       5,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    true,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityNormal,
		Tags:           []string{"magic", "ice", "mage", "shield"},
	},
	"channel_mana": {
		ID:             "channel_mana",
		Name:           "Channel Mana",
		Description:    "Focus to regenerate mana over 3 ticks.",
		Category:       CategorySkill,
		Type:           ActionChannel,
		ChannelTicks:   3,
		ManaCost:       0,
		StaminaCost:    5,
		BaseHeal:       10, // Per tick
		Cooldown:       0,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      true,
		Priority:       PriorityNormal,
		Tags:           []string{"magic", "mage", "channel", "regen"},
	},
	
	// Rogue skills
	"backstab": {
		ID:             "backstab",
		Name:           "Backstab",
		Description:    "Strike from behind for 3x damage. Requires being behind the target.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    12,
		BaseDamage:     30,
		Cooldown:       3,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PriorityFast,
		Tags:           []string{"physical", "melee", "rogue", "stealth"},
	},
	"poison_blade": {
		ID:             "poison_blade",
		Name:           "Poison Blade",
		Description:    "Coat your weapon in poison, applying DoT for 3 ticks.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    8,
		BaseDamage:     5,
		Cooldown:       4,
		MaxTargets:     1,
		Range:          0,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    true,
		IsSupport:      false,
		Priority:       PriorityNormal,
		Tags:           []string{"physical", "melee", "rogue", "poison", "dot"},
	},
	"smoke_bomb": {
		ID:             "smoke_bomb",
		Name:           "Smoke Bomb",
		Description:    "Create a smoke cloud, increasing dodge chance for 2 ticks.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       0,
		StaminaCost:    10,
		BaseDamage:     0,
		Cooldown:       5,
		MaxTargets:     0,
		Range:          0,
		RequiresTarget: false,
		IsDefensive:    true,
		IsOffensive:    false,
		IsSupport:      false,
		Priority:       PriorityFast,
		Tags:           []string{"defensive", "rogue", "evasion"},
	},
	
	// Cleric skills
	"heal": {
		ID:             "heal",
		Name:           "Heal",
		Description:    "Restore health to a target ally.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       15,
		StaminaCost:    0,
		BaseDamage:     0,
		BaseHeal:       20,
		Cooldown:       2,
		MaxTargets:     1,
		Range:          2,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      true,
		Priority:       PriorityNormal,
		Tags:           []string{"magic", "holy", "cleric", "healing"},
	},
	"greater_heal": {
		ID:             "greater_heal",
		Name:           "Greater Heal",
		Description:    "Channel powerful healing over 2 ticks.",
		Category:       CategorySkill,
		Type:           ActionChannel,
		ChannelTicks:   2,
		ManaCost:       25,
		StaminaCost:    0,
		BaseDamage:     0,
		BaseHeal:       15, // Per tick
		Cooldown:       4,
		MaxTargets:     1,
		Range:          2,
		RequiresTarget: true,
		IsDefensive:    false,
		IsOffensive:    false,
		IsSupport:      true,
		Priority:       PriorityNormal,
		Tags:           []string{"magic", "holy", "cleric", "healing", "channel"},
	},
	"shield_of_faith": {
		ID:             "shield_of_faith",
		Name:           "Shield of Faith",
		Description:    "Grant a protective shield for 3 ticks.",
		Category:       CategorySkill,
		Type:           ActionInstant,
		TickCost:       1,
		ManaCost:       20,
		StaminaCost:    0,
		BaseDamage:     0,
		Cooldown:       6,
		MaxTargets:     1,
		Range:          2,
		RequiresTarget: true,
		IsDefensive:    true,
		IsOffensive:    false,
		IsSupport:      true,
		Priority:       PriorityNormal,
		Tags:           []string{"magic", "holy", "cleric", "shield"},
	},
}

// GetActionDefinition retrieves an action definition by ID
func GetActionDefinition(id string) (*ActionDefinition, bool) {
	// Check basic actions first
	if action, exists := BasicActions[id]; exists {
		return action, true
	}
	// Check skill actions
	if action, exists := SkillActions[id]; exists {
		return action, true
	}
	return nil, false
}

// GetAllActionDefinitions returns all action definitions
func GetAllActionDefinitions() map[string]*ActionDefinition {
	result := make(map[string]*ActionDefinition)
	for k, v := range BasicActions {
		result[k] = v
	}
	for k, v := range SkillActions {
		result[k] = v
	}
	return result
}

// GetActionsByCategory returns all actions of a specific category
func GetActionsByCategory(category ActionCategory) []*ActionDefinition {
	var result []*ActionDefinition
	for _, action := range BasicActions {
		if action.Category == category {
			result = append(result, action)
		}
	}
	for _, action := range SkillActions {
		if action.Category == category {
			result = append(result, action)
		}
	}
	return result
}

// GetActionsByType returns all actions of a specific type
func GetActionsByType(actionType ActionType) []*ActionDefinition {
	var result []*ActionDefinition
	for _, action := range BasicActions {
		if action.Type == actionType {
			result = append(result, action)
		}
	}
	for _, action := range SkillActions {
		if action.Type == actionType {
			result = append(result, action)
		}
	}
	return result
}