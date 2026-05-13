package routes

import (
	"herbst-server/db"
	"herbst-server/db/schema"
)

// questToView converts a db.Quest entity to a questView.
func questToView(q *db.Quest) questView {
	prereqs := q.PrerequisiteQuestIds
	if prereqs == nil {
		prereqs = []string{}
	}
	objs := make([]questObjectiveInput, len(q.Objectives))
	for i, o := range q.Objectives {
		labels := []string{}
		if o.Labels != nil && len(o.Labels) > 0 {
			labels = o.Labels
		}
		objs[i] = questObjectiveInput{
			Type:      o.Type, TargetID: o.TargetID,
			TagFilter: o.TagFilter, Count: o.Count, Labels: labels, Hint: o.Hint,
		}
	}
	r := q.Rewards
	itemIDs := r.ItemIDs
	if itemIDs == nil {
		itemIDs = []string{}
	}
	tagAdds := r.TagAdds
	if tagAdds == nil {
		tagAdds = []string{}
	}
	tagRemoves := r.TagRemoves
	if tagRemoves == nil {
		tagRemoves = []string{}
	}
	effIDs := r.EffectIDs
	if effIDs == nil {
		effIDs = []int{}
	}
	achIDs := r.AchievementIDs
	if achIDs == nil {
		achIDs = []int{}
	}
	return questView{
		ID: q.ID, Name: q.Name, Description: q.Description,
		PrerequisiteQuestIDs: prereqs, Objectives: objs,
		Rewards: questRewardsInput{
			XP: r.XP, ItemIDs: itemIDs, EffectIDs: effIDs,
			TagAdds: tagAdds, TagRemoves: tagRemoves,
			AchievementIDs: achIDs,
		},
		RepeatMode:    string(q.RepeatMode),
		CooldownHours: q.CooldownHours, IsActive: q.IsActive,
	}
}

// applyQuestRewards returns a summary of rewards that would be applied.
func applyQuestRewards(rewards schema.QuestRewards) map[string]interface{} {
	return map[string]interface{}{
		"xp":              rewards.XP,
		"item_ids":        rewards.ItemIDs,
		"effect_ids":      rewards.EffectIDs,
		"tag_adds":        rewards.TagAdds,
		"tag_removes":     rewards.TagRemoves,
		"achievement_ids": rewards.AchievementIDs,
	}
}