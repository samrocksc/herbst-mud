package routes

import (
	"herbst-server/db/quest"
)

// questView is the JSON response shape for a Quest entity.
type questView struct {
	ID                   int                  `json:"id"`
	Name                 string               `json:"name"`
	WorldID              string               `json:"world_id"`
	Description          string               `json:"description"`
	PrerequisiteQuestIDs []string             `json:"prerequisite_quest_ids"`
	Objectives           []questObjectiveInput `json:"objectives"`
	Rewards              questRewardsInput    `json:"rewards"`
	RepeatMode           string               `json:"repeat_mode"`
	CooldownHours        int                  `json:"cooldown_hours"`
	IsActive             bool                 `json:"is_active"`
}

// questObjectiveInput mirrors schema.QuestObjective for JSON input/output.
type questObjectiveInput struct {
	Type      string   `json:"type"`
	TargetID  string   `json:"target_id"`
	TagFilter string   `json:"tag_filter"`
	Count     int      `json:"count"`
	Labels    []string `json:"labels"`
	Hint      string   `json:"hint"`
}

// questRewardsInput mirrors schema.QuestRewards for JSON input/output.
type questRewardsInput struct {
	XP             int      `json:"xp"`
	ItemIDs        []string `json:"item_ids"`
	EffectIDs      []int    `json:"effect_ids"`
	TagAdds        []string `json:"tag_adds"`
	TagRemoves     []string `json:"tag_removes"`
	AchievementIDs []int    `json:"achievement_ids"`
}

// questInput is the JSON request shape for creating/updating a Quest.
type questInput struct {
	Name                 *string                `json:"name"`
	Description          *string                `json:"description"`
	WorldID              *string                `json:"world_id"`
	PrerequisiteQuestIDs *[]string              `json:"prerequisite_quest_ids"`
	Objectives           *[]questObjectiveInput `json:"objectives"`
	Rewards              *questRewardsInput     `json:"rewards"`
	RepeatMode           *string                `json:"repeat_mode"`
	CooldownHours        *int                   `json:"cooldown_hours"`
	IsActive             *bool                  `json:"is_active"`
}

// validRepeatModes lists allowed repeat_mode enum values.
var validRepeatModes = map[string]bool{
	string(quest.RepeatModeNone):     true,
	string(quest.RepeatModeCooldown): true,
	string(quest.RepeatModeAlways):   true,
}