package worldexport

import (
	"context"
	"encoding/json"
	"fmt"

	"herbst-server/db"
	"herbst-server/db/gender"
	"herbst-server/db/race"
)

func importRaces(ctx context.Context, client *db.Client, races []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, r := range races {
		oldID := intVal(r["id"])
		name := strVal(r["name"])
		existing, err := client.Race.Query().Where(race.NameEQ(name)).Only(ctx)
		if err == nil && existing != nil {
			maps.Races[oldID] = existing.ID
			continue
		}
		created, err := client.Race.Create().
			SetWorldID(newWorldID).
			SetName(name).
			SetDisplayName(strVal(r["display_name"])).
			SetDescription(strVal(r["description"])).
			SetNillableStatModifiers(jsonStrPtr(r["stat_modifiers"])).
			SetNillableSkillGrants(jsonStrPtr(r["skill_grants"])).
			SetEquipmentSlots(strSliceVal(r["equipment_slots"], []string{})).
			SetRequirementTags(strSliceVal(r["requirement_tags"], []string{})).
			SetNillableColor(strPtr(r["color"])).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("race %d: %w", oldID, err)
		}
		maps.Races[oldID] = created.ID
		count++
	}
	return count, nil
}

func importGenders(ctx context.Context, client *db.Client, genders []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, g := range genders {
		oldID := intVal(g["id"])
		name := strVal(g["name"])
		existing, err := client.Gender.Query().Where(gender.NameEQ(name)).Only(ctx)
		if err == nil && existing != nil {
			maps.Genders[oldID] = existing.ID
			continue
		}
		created, err := client.Gender.Create().
			SetName(name).
			SetDisplayName(strVal(g["display_name"])).
			SetSubjectPronoun(strVal(g["subject_pronoun"])).
			SetObjectPronoun(strVal(g["object_pronoun"])).
			SetPossessivePronoun(strVal(g["possessive_pronoun"])).
			SetWorldID(newWorldID).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("gender %d: %w", oldID, err)
		}
		maps.Genders[oldID] = created.ID
		count++
	}
	return count, nil
}

func jsonStrPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	switch s := v.(type) {
	case string:
		if s == "" {
			return nil
		}
		return &s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	s := string(b)
	if s == "null" || s == "{}" || s == "[]" {
		return nil
	}
	return &s
}
