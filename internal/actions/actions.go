package actions

import (
	"github.com/sam/makeathing/internal/characters"
)

// Action represents an action that can be performed
type Action struct {
	Name        string
	Type        characters.SkillType
	Description string
	Requirements ActionRequirements
}

// ActionRequirements represents what is needed to perform an action
type ActionRequirements struct {
	MinLevel      int
	RequiredStats characters.Stats
	RequiredSkills []string
}