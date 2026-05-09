package questservice

import "time"

// QuestDef is a cached quest definition from the REST API.
type QuestDef struct {
	ID                    int              `json:"id"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	PrerequisiteQuestIDs  []string         `json:"prerequisite_quest_ids"`
	Objectives            []QuestObjective `json:"objectives"`
	Rewards               QuestRewards     `json:"rewards"`
	RepeatMode            string           `json:"repeat_mode"`
	CooldownHours         int              `json:"cooldown_hours"`
	IsActive              bool             `json:"is_active"`
}

// QuestObjective is a single objective within a quest.
type QuestObjective struct {
	Type     string `json:"type"`
	TargetID string `json:"target_id"`
	Count    int    `json:"count"`
	Label    string `json:"label"`
	Hint     string `json:"hint"`
}

// QuestRewards describes what a character receives on quest completion.
type QuestRewards struct {
	XP             int      `json:"xp"`
	ItemIDs        []string `json:"item_ids"`
	EffectIDs      []int    `json:"effect_ids"`
	TagAdds        []string `json:"tag_adds"`
	TagRemoves     []string `json:"tag_removes"`
	AchievementIDs []int    `json:"achievement_ids"`
}

// QuestProgress represents a character's progress on a quest.
type QuestProgress struct {
	ID                int                      `json:"id"`
	QuestID           int                      `json:"quest_id"`
	QuestName         string                   `json:"quest_name"`
	QuestDescription  string                   `json:"quest_description"`
	Status            string                   `json:"status"`
	CurrentStep       int                      `json:"current_step"`
	Objectives        []QuestProgressObjective `json:"objectives"`
	StartedAt         time.Time                `json:"started_at"`
	CompletedAt       *time.Time               `json:"completed_at"`
}

// QuestProgressObjective shows progress on a single objective.
type QuestProgressObjective struct {
	Type     string `json:"type"`
	TargetID string `json:"target_id"`
	Count    int    `json:"count"`
	Current  int    `json:"current"`
	Label    string `json:"label"`
	Complete bool   `json:"complete"`
}