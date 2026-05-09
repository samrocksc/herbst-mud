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
		objs[i] = questObjectiveInput{
			Type: o.Type, TargetID: o.TargetID,
			Count: o.Count, Label: o.Label, Hint: o.Hint,
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

// questObjectivesToSchema converts input objectives to schema types.
func questObjectivesToSchema(objs []questObjectiveInput) []schema.QuestObjective {
	result := make([]schema.QuestObjective, len(objs))
	for i, o := range objs {
		result[i] = schema.QuestObjective{
			Type: o.Type, TargetID: o.TargetID,
			Count: o.Count, Label: o.Label, Hint: o.Hint,
		}
	}
	return result
}

// questRewardsToSchema converts input rewards to schema type.
func questRewardsToSchema(r questRewardsInput) schema.QuestRewards {
	return schema.QuestRewards{
		XP: r.XP, ItemIDs: r.ItemIDs, EffectIDs: r.EffectIDs,
		TagAdds: r.TagAdds, TagRemoves: r.TagRemoves,
		AchievementIDs: r.AchievementIDs,
	}
}

// applyQuestRewards returns a summary of rewards that would be applied.
// Actual reward application is handled by the game engine events system.
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