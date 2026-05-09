package schema

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