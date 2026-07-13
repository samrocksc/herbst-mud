package worldexport

import (
	"context"
	"fmt"

	"herbst-server/db"
)

// importSkills creates Skill records in the new world, remapping parent_skill_id via idMaps.
// Skills are imported in two passes: first all skills without a parent, then skills
// with a parent (so the parent already exists in idMaps). If a parent ID can't be
// resolved, the skill is imported without a parent (parent_skill_id = nil).
func importSkills(ctx context.Context, client *db.Client, skills []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0

	// Pass 1: skills without parent_skill_id
	for _, s := range skills {
		oldID := intVal(s["id"])
		if oldID == 0 {
			continue
		}
		parentID := intVal(s["parent_skill_id"])
		if parentID != 0 {
			continue // defer to pass 2
		}
		created, err := createSkill(ctx, client, s, newWorldID, 0)
		if err != nil {
			return count, fmt.Errorf("skill %d: %w", oldID, err)
		}
		maps.Skills[oldID] = created.ID
		count++
	}

	// Pass 2: skills with parent_skill_id (remap via idMaps)
	for _, s := range skills {
		oldID := intVal(s["id"])
		if oldID == 0 {
			continue
		}
		parentID := intVal(s["parent_skill_id"])
		if parentID == 0 {
			continue // already handled in pass 1
		}
		// Skip if already imported (shouldn't happen, but guard against dupes)
		if _, ok := maps.Skills[oldID]; ok {
			continue
		}
		var remappedParentID int
		if newParentID, ok := maps.Skills[parentID]; ok {
			remappedParentID = newParentID
		}
		created, err := createSkill(ctx, client, s, newWorldID, remappedParentID)
		if err != nil {
			return count, fmt.Errorf("skill %d: %w", oldID, err)
		}
		maps.Skills[oldID] = created.ID
		count++
	}

	return count, nil
}

// createSkill creates a single Skill record from the exported map.
// parentID > 0 sets the parent_skill_id; parentID == 0 leaves it nil.
func createSkill(ctx context.Context, client *db.Client, s map[string]interface{}, newWorldID string, parentID int) (*db.Skill, error) {
	builder := client.Skill.Create().
		SetWorldID(atoi(newWorldID)).
		SetName(strVal(s["name"])).
		SetDisplayName(strVal(s["display_name"])).
		SetCategory(strValOr(s["category"], "weapon")).
		SetMaxLevel(intValOr(s["max_level"], 100)).
		SetXpCurveMode(strValOr(s["xp_curve_mode"], "percentage"))

	// Optional description
	if d := strVal(s["description"]); d != "" {
		builder = builder.SetDescription(d)
	}

	// Optional xp_curve_data (JSON map)
	if xcd := mapVal(s["xp_curve_data"]); xcd != nil {
		builder = builder.SetXpCurveData(xcd)
	}

	// Parent skill ID
	if parentID > 0 {
		builder = builder.SetParentSkillID(parentID)
	}

	return builder.Save(ctx)
}