package schema

// QuestObjective is a single objective within a quest.
type QuestObjective struct {
	Type     string `json:"type"`      // kill | explore | collect
	TargetID string `json:"target_id"` // NPC template ID, room ID, or item template ID
	Count    int    `json:"count"`     // required count (default 1)
	Label    string `json:"label"`     // display: "Kill 3 Goblin Shamans"
	Hint     string `json:"hint"`      // optional hint shown to player
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