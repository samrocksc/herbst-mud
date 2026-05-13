package schema

// QuestObjective is a single objective within a quest.
type QuestObjective struct {
	Type       string   `json:"type"`        // kill | explore | collect
	TargetID   string   `json:"target_id"`   // NPC template ID, room ID, or item template ID (optional if tag_filter is set)
	TagFilter  string   `json:"tag_filter"`  // tag to filter targets (e.g., "bandit" for NPCs with that tag)
	Count      int      `json:"count"`       // required count (default 1)
	Labels     []string `json:"labels"`      // labels for each target (for multi-target objectives)
	Hint       string   `json:"hint"`        // optional hint shown to player
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