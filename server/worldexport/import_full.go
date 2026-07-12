package worldexport

import (
	"context"
	"encoding/json"
	"fmt"

	"herbst-server/db"
	"herbst-server/db/factioncategory"
	"herbst-server/db/npctemplate"
	"herbst-server/db/quest"
	"herbst-server/db/schema"
	"herbst-server/db/tag"
)

// importFactionCategories creates FactionCategory records in the new world.
func importFactionCategories(ctx context.Context, client *db.Client, cats []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, c := range cats {
		oldID := intVal(c["id"])
		name := strVal(c["name"])
		// Check if a category with this name already exists in the new world
		existing, err := client.FactionCategory.Query().
			Where(factioncategory.NameEQ(name), factioncategory.WorldIDEQ(newWorldID)).
			Only(ctx)
		if err == nil && existing != nil {
			maps.FactionCategories[oldID] = existing.ID
			continue
		}
		created, err := client.FactionCategory.Create().
			SetWorldID(newWorldID).
			SetName(name).
			SetDisplayName(strVal(c["display_name"])).
			SetNillableDescription(strPtr(c["description"])).
			SetMaxMemberships(intValOr(c["max_memberships"], 1)).
			SetAutoJoin(boolVal(c["auto_join"])).
			SetInitialConfig(boolVal(c["initial_config"])).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("faction_category %d: %w", oldID, err)
		}
		maps.FactionCategories[oldID] = created.ID
		count++
	}
	return count, nil
}

// importFactions creates Faction records, remapping category FK via idMaps.
func importFactions(ctx context.Context, client *db.Client, factions []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, f := range factions {
		oldID := intVal(f["id"])
		name := strVal(f["name"])
		// Faction names are globally unique, so suffix with the new world ID
		// to avoid collisions with the source world.
		importedName := name + "_w" + newWorldID
		builder := client.Faction.Create().
			SetName(importedName).
			SetWorldID(newWorldID).
			SetDisplayName(strVal(f["display_name"])).
			SetNillableDescription(strPtr(f["description"]))
		// Set member_tags if present
		if mt := strSlicePtrVal(f["member_tags"]); mt != nil {
			builder = builder.SetMemberTags(*mt)
		}
		// Set stat_bonuses if present
		if sb := parseStatBonuses(f["stat_bonuses"]); sb != nil {
			builder = builder.SetStatBonuses(*sb)
		}
		// Set specialties if present
		if sp := parseSpecialties(f["specialties"]); sp != nil {
			builder = builder.SetSpecialties(*sp)
		}
		// Remap category FK — check both direct field and edge
		catOldID := intVal(f["faction_category_factions"])
		if catOldID == 0 {
			catOldID = edgeID(f["edges"], "category")
		}
		if newCatID, ok := maps.FactionCategories[catOldID]; ok {
			builder = builder.SetCategoryID(newCatID)
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("faction %d: %w", oldID, err)
		}
		maps.Factions[oldID] = created.ID
		count++
	}
	return count, nil
}

// importAbilities creates Ability records, remapping faction FK via idMaps.
func importAbilities(ctx context.Context, client *db.Client, abilities []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, a := range abilities {
		oldID := intVal(a["id"])
		name := strVal(a["name"])
		builder := client.Ability.Create().
			SetName(name).
			SetDescription(strVal(a["description"])).
			SetAbilityType(strVal(a["ability_type"])).
			SetCost(intVal(a["cost"])).
			SetCooldown(intVal(a["cooldown"])).
			SetManaCost(intVal(a["mana_cost"])).
			SetStaminaCost(intVal(a["stamina_cost"])).
			SetHpCost(intVal(a["hp_cost"])).
			SetRequirements(strValOr(a["requirements"], "")).
			SetSlug(strValOr(a["slug"], "")).
			SetAbilityClass(strValOr(a["ability_class"], "active")).
			SetWorldID(newWorldID)
		// Optional fields
		if rt := strVal(a["required_tag"]); rt != "" {
			builder = builder.SetRequiredTag(rt)
		}
		if pc := floatValOr(a["proc_chance"], 0); pc != 0 {
			builder = builder.SetProcChance(pc)
		}
		if pe := strVal(a["proc_event"]); pe != "" {
			builder = builder.SetProcEvent(pe)
		}
		// Remap faction edge — check edge object first, then direct field
		factionOldID := edgeID(a["edges"], "faction")
		if factionOldID == 0 {
			factionOldID = intVal(a["faction_abilities"])
		}
		if factionOldID > 0 {
			if newFactionID, ok := maps.Factions[factionOldID]; ok {
				builder = builder.SetFactionID(newFactionID)
			}
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("ability %d (%s): %w", oldID, name, err)
		}
		maps.Abilities[oldID] = created.ID
		count++
	}
	return count, nil
}

// importEffects creates AbilityEffect records, remapping ability FK via idMaps.
func importEffects(ctx context.Context, client *db.Client, effects []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, e := range effects {
		// Find the ability ID from the edge or direct field
		abilityOldID := edgeID(e["edges"], "ability")
		if abilityOldID == 0 {
			abilityOldID = intVal(e["ability_effects"])
		}
		newAbilityID, ok := maps.Abilities[abilityOldID]
		if !ok {
			continue // skip effects for abilities that weren't imported
		}
		builder := client.AbilityEffect.Create().
			SetEffectType(strVal(e["effect_type"])).
			SetTarget(strValOr(e["target"], "enemy")).
			SetValue(intVal(e["value"])).
			SetAbilityID(newAbilityID)
		// Optional fields
		if ds := strVal(e["damage_subtype"]); ds != "" {
			builder = builder.SetDamageSubtype(ds)
		}
		if d := intVal(e["duration"]); d != 0 {
			builder = builder.SetDuration(d)
		}
		if ss := strVal(e["scaling_stat"]); ss != "" {
			builder = builder.SetScalingStat(ss)
		}
		if sr := floatValOr(e["scaling_ratio"], 0); sr != 0 {
			builder = builder.SetScalingRatio(sr)
		}
		if so := intVal(e["sort_order"]); so != 0 {
			builder = builder.SetSortOrder(so)
		}
		// No effect_message field on AbilityEffect schema
		_ = e["effect_message"]
		_, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("effect for ability %d: %w", abilityOldID, err)
		}
		count++
	}
	return count, nil
}

// importTags creates Tag records in the new world.
func importTags(ctx context.Context, client *db.Client, tags []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, t := range tags {
		oldID := intVal(t["id"])
		name := strVal(t["name"])
		existing, err := client.Tag.Query().Where(tag.NameEQ(name), tag.WorldIDEQ(newWorldID)).Only(ctx)
		if err == nil && existing != nil {
			maps.Tags[oldID] = existing.ID
			continue
		}
		builder := client.Tag.Create().
			SetWorldID(newWorldID).
			SetName(name)
		if c := strVal(t["color"]); c != "" {
			builder = builder.SetColor(c)
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("tag %d: %w", oldID, err)
		}
		maps.Tags[oldID] = created.ID
		count++
	}
	return count, nil
}

// importEquipment creates EquipmentTemplate records in the new world.
func importEquipment(ctx context.Context, client *db.Client, items []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, eq := range items {
		oldID := intVal(eq["id"])
		builder := client.EquipmentTemplate.Create().
			SetSlug(strVal(eq["slug"])).
			SetWorldID(newWorldID).
			SetName(strVal(eq["name"])).
			SetDescription(strVal(eq["description"])).
			SetSlot(strVal(eq["slot"])).
			SetItemType(strValOr(eq["item_type"], "weapon"))
		// Optional numeric fields
		if v := intVal(eq["level"]); v != 0 {
			builder = builder.SetLevel(v)
		}
		if v := intVal(eq["weight"]); v != 0 {
			builder = builder.SetWeight(v)
		}
		// Stats JSON (map[string]int)
		if s := mapValInt(eq["stats"], nil); s != nil {
			builder = builder.SetStats(s)
		}
		// Color
		if c := strVal(eq["color"]); c != "" {
			builder = builder.SetColor(c)
		}
		// Visibility / container / locked
		builder = builder.SetIsVisible(boolValOr(eq["is_visible"], true))
		if v := boolVal(eq["is_immovable"]); v {
			builder = builder.SetIsImmovable(v)
		}
		if v := boolVal(eq["is_container"]); v {
			builder = builder.SetIsContainer(v)
		}
		if v := intVal(eq["container_capacity"]); v != 0 {
			builder = builder.SetContainerCapacity(v)
		}
		if v := boolVal(eq["is_locked"]); v {
			builder = builder.SetIsLocked(v)
		}
		// Combat fields
		if v := strVal(eq["effect_type"]); v != "" {
			builder = builder.SetEffectType(v)
		}
		if v := intVal(eq["effect_value"]); v != 0 {
			builder = builder.SetEffectValue(v)
		}
		if v := intVal(eq["effect_duration"]); v != 0 {
			builder = builder.SetEffectDuration(v)
		}
		if v := intVal(eq["armor_rating"]); v != 0 {
			builder = builder.SetArmorRating(v)
		}
		if v := strVal(eq["armor_type"]); v != "" {
			builder = builder.SetArmorType(v)
		}
		if v := strVal(eq["rarity"]); v != "" {
			builder = builder.SetRarity(v)
		}
		if v := strVal(eq["skill_requirement"]); v != "" {
			builder = builder.SetSkillRequirement(v)
		}
		if v := intVal(eq["skill_requirement_level"]); v != 0 {
			builder = builder.SetSkillRequirementLevel(v)
		}
		if v := intVal(eq["damage_dice_count"]); v != 0 {
			builder = builder.SetDamageDiceCount(v)
		}
		if v := intVal(eq["damage_dice_sides"]); v != 0 {
			builder = builder.SetDamageDiceSides(v)
		}
		if v := intVal(eq["damage_bonus"]); v != 0 {
			builder = builder.SetDamageBonus(v)
		}
		if v := strVal(eq["damage_type"]); v != "" {
			builder = builder.SetDamageType(v)
		}
		if v := strVal(eq["weapon_type"]); v != "" {
			builder = builder.SetWeaponType(v)
		}
		if v := boolVal(eq["is_two_handed"]); v {
			builder = builder.SetIsTwoHanded(v)
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("equipment %d: %w", oldID, err)
		}
		maps.Equipment[oldID] = created.ID
		count++
	}
	return count, nil
}

// importNPCTemplates creates NPCTemplate records, remapping race_id via idMaps.
func importNPCTemplates(ctx context.Context, client *db.Client, templates []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, t := range templates {
		oldSlug := strVal(t["id"])
		// NPC template IDs and slugs are globally unique, so generate new ones.
		newSlug := oldSlug + "_w" + newWorldID
		builder := client.NPCTemplate.Create().
			SetID(newSlug).
			SetWorldID(newWorldID).
			SetName(strVal(t["name"])).
			SetDescription(strVal(t["description"]))
		// Set slug if present (suffix to avoid collision)
		if oldSlugVal := strVal(t["slug"]); oldSlugVal != "" {
			builder = builder.SetSlug(oldSlugVal + "_w" + newWorldID)
		} else {
			builder = builder.SetSlug(newSlug)
		}
		// Remap race_id
		if oldRaceID := intVal(t["race_id"]); oldRaceID > 0 {
			if newRaceID, ok := maps.Races[oldRaceID]; ok {
				builder = builder.SetRaceID(newRaceID)
			}
		}
		// Disposition is an enum (hostile/friendly/neutral). Fall back to neutral for unknown values.
		if v := strVal(t["disposition"]); v != "" {
			switch v {
			case "hostile", "friendly", "neutral":
				builder = builder.SetDisposition(npctemplate.Disposition(v))
			default:
				builder = builder.SetDisposition(npctemplate.DispositionNeutral)
			}
		}
		if v := intVal(t["level"]); v != 0 {
			builder = builder.SetLevel(v)
		}
		if v := intVal(t["xp_value"]); v != 0 {
			builder = builder.SetXpValue(v)
		}
		// Skills is a required JSON field, always set it
		skills := mapValInt(t["skills"], map[string]int{})
		builder = builder.SetSkills(skills)
		// trades_with is required (field.Strings)
		builder = builder.SetTradesWith(strSliceVal(t["trades_with"], []string{}))
		// greeting is required (field.Text)
		builder = builder.SetGreeting(strVal(t["greeting"]))
		// respawn_rooms is required (field.JSON)
		builder = builder.SetRespawnRooms(strSliceVal(t["respawn_rooms"], []string{}))
		if v := intVal(t["respawn_cooldown"]); v != 0 {
			builder = builder.SetRespawnCooldown(v)
		}
		if v := strVal(t["roam_pattern"]); v != "" {
			switch v {
			case "static", "wander", "patrol", "return_home":
				builder = builder.SetRoamPattern(npctemplate.RoamPattern(v))
			default:
				builder = builder.SetRoamPattern(npctemplate.RoamPatternStatic)
			}
		}
		if s := strSliceVal(t["roam_zone_ids"], nil); s != nil {
			builder = builder.SetRoamZoneIds(s)
		}
		if v := intVal(t["roam_interval_seconds"]); v != 0 {
			builder = builder.SetRoamIntervalSeconds(v)
		}
		if v := intVal(t["roam_pause_min_seconds"]); v != 0 {
			builder = builder.SetRoamPauseMinSeconds(v)
		}
		if v := intVal(t["roam_pause_max_seconds"]); v != 0 {
			builder = builder.SetRoamPauseMaxSeconds(v)
		}
		_, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("npc_template %s: %w", oldSlug, err)
		}
		maps.NPCTemplates[oldSlug] = newSlug
		count++
	}
	return count, nil
}

// importDialogNodes creates DialogNode records, remapping npc_template edge.
func importDialogNodes(ctx context.Context, client *db.Client, nodes []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, n := range nodes {
		oldID := strVal(n["id"])
		// Dialog node IDs are globally unique, suffix with world ID
		newID := oldID + "_w" + newWorldID
		builder := client.DialogNode.Create().
			SetID(newID).
			SetWorldID(newWorldID).
			SetNpcText(strVal(n["npc_text"]))
		// Responses JSON → []schema.DialogResponse
		if rBytes, err := json.Marshal(n["responses"]); err == nil {
			var responses []schema.DialogResponse
			if json.Unmarshal(rBytes, &responses) == nil && len(responses) > 0 {
				builder = builder.SetResponses(responses)
			}
		}
		// Entry node
		if v := boolVal(n["is_entry"]); v {
			builder = builder.SetIsEntry(v)
		}
		// Entry condition
		if v := strVal(n["entry_condition"]); v != "" {
			builder = builder.SetEntryCondition(v)
		}
		// On enter effects JSON → []int
		if effBytes, err := json.Marshal(n["on_enter_effects"]); err == nil {
			var effs []int
			if json.Unmarshal(effBytes, &effs) == nil && len(effs) > 0 {
				builder = builder.SetOnEnterEffects(effs)
			}
		}
		// Remap npc_template edge
		tmplSlug := edgeStrID(n["edges"], "npc_template")
		if tmplSlug == "" {
			tmplSlug = strVal(n["npc_template_id"])
		}
		if tmplSlug != "" {
			if newSlug, ok := maps.NPCTemplates[tmplSlug]; ok {
				builder = builder.SetNpcTemplateID(newSlug)
			} else {
				builder = builder.SetNpcTemplateID(tmplSlug)
			}
		}
		_, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("dialog_node %s: %w", oldID, err)
		}
		maps.DialogNodes[oldID] = newID
		count++
	}
	return count, nil
}

// importQuests creates Quest records in the new world.
func importQuests(ctx context.Context, client *db.Client, quests []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, q := range quests {
		oldID := intVal(q["id"])
		// Quest names are globally unique, suffix with world ID
		importedName := strVal(q["name"]) + "_w" + newWorldID
		builder := client.Quest.Create().
			SetName(importedName).
			SetWorldID(newWorldID).
			SetDescription(strVal(q["description"]))
		// Prerequisite quest IDs (optional)
		if s := strSliceVal(q["prerequisite_quest_ids"], nil); s != nil {
			builder = builder.SetPrerequisiteQuestIds(s)
		}
		// Objectives JSON → []schema.QuestObjective (required, has default)
		if objBytes, err := json.Marshal(q["objectives"]); err == nil {
			var objectives []schema.QuestObjective
			if json.Unmarshal(objBytes, &objectives) == nil {
				builder = builder.SetObjectives(objectives)
			}
		}
		// Rewards JSON → schema.QuestRewards (required, has default)
		if rewBytes, err := json.Marshal(q["rewards"]); err == nil {
			var rewards schema.QuestRewards
			if json.Unmarshal(rewBytes, &rewards) == nil {
				builder = builder.SetRewards(rewards)
			}
		}
		// Repeat mode (enum: none/cooldown/always)
		if v := strVal(q["repeat_mode"]); v != "" {
			switch v {
			case "none", "cooldown", "always":
				builder = builder.SetRepeatMode(quest.RepeatMode(v))
			default:
				builder = builder.SetRepeatMode(quest.RepeatModeNone)
			}
		}
		// Main type (enum: hunter/collector/explorer/general)
		if v := strVal(q["main_type"]); v != "" {
			switch v {
			case "hunter", "collector", "explorer", "general":
				builder = builder.SetMainType(quest.MainType(v))
			default:
				builder = builder.SetMainType(quest.MainTypeGeneral)
			}
		}
		// Cooldown hours
		if v := intVal(q["cooldown_hours"]); v != 0 {
			builder = builder.SetCooldownHours(v)
		}
		// Is active
		builder = builder.SetIsActive(boolValOr(q["is_active"], true))
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("quest %d: %w", oldID, err)
		}
		maps.Quests[oldID] = created.ID
		count++
	}
	return count, nil
}

// importRecipes creates CraftingRecipe records in the new world.
func importRecipes(ctx context.Context, client *db.Client, recipes []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, r := range recipes {
		oldID := intVal(r["id"])
		builder := client.CraftingRecipe.Create().
			SetName(strVal(r["name"])).
			SetDisplayName(strVal(r["display_name"])).
			SetWorldID(newWorldID)
		if v := strVal(r["description"]); v != "" {
			builder = builder.SetDescription(v)
		}
		if v := strVal(r["required_station_tag"]); v != "" {
			builder = builder.SetRequiredStationTag(v)
		}
		if v := strVal(r["required_class"]); v != "" {
			builder = builder.SetRequiredClass(v)
		}
		if v := strVal(r["required_skill"]); v != "" {
			builder = builder.SetRequiredSkill(v)
		}
		if v := intVal(r["required_skill_level"]); v != 0 {
			builder = builder.SetRequiredSkillLevel(v)
		}
		if v := intVal(r["craft_time_secs"]); v != 0 {
			builder = builder.SetCraftTimeSecs(v)
		}
		// Inputs JSON → schema.CraftingInput
		if inBytes, err := json.Marshal(r["inputs"]); err == nil {
			var inputs []schema.CraftingInput
			if json.Unmarshal(inBytes, &inputs) == nil && len(inputs) > 0 {
				builder = builder.SetInputs(inputs)
			}
		}
		// Outputs JSON → schema.CraftingOutput
		if outBytes, err := json.Marshal(r["outputs"]); err == nil {
			var outputs []schema.CraftingOutput
			if json.Unmarshal(outBytes, &outputs) == nil && len(outputs) > 0 {
				builder = builder.SetOutputs(outputs)
			}
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("recipe %d: %w", oldID, err)
		}
		maps.Recipes[oldID] = created.ID
		count++
	}
	return count, nil
}

// importTriggers creates Trigger records, remapping room, equipment, effect, recipe, dialog_node FKs.
func importTriggers(ctx context.Context, client *db.Client, triggers []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, t := range triggers {
		oldID := intVal(t["id"])
		builder := client.Trigger.Create().
			SetName(strVal(t["name"])).
			SetWorldID(newWorldID).
			SetTriggerType(strVal(t["trigger_type"])).
			SetExamineWeight(intVal(t["examine_weight"])).
			SetTargetType(strValOr(t["target_type"], "")).
			SetCondition(strValOr(t["condition"], "")).
			SetEnabled(boolValOr(t["enabled"], true))
		// Target ID
		if v := intVal(t["target_id"]); v != 0 {
			builder = builder.SetTargetID(v)
		}
		// Remap room_id
		if oldRoomID := intVal(t["room_id"]); oldRoomID > 0 {
			if newRoomID, ok := maps.Rooms[oldRoomID]; ok {
				builder = builder.SetRoomID(newRoomID)
			}
		}
		// Remap equipment_id
		if oldEqID := intVal(t["equipment_id"]); oldEqID > 0 {
			if newEqID, ok := maps.Equipment[oldEqID]; ok {
				builder = builder.SetEquipmentID(newEqID)
			}
		}
		// Effect edge — effects are global, ID stays the same
		if effID := edgeID(t["edges"], "effect"); effID > 0 {
			builder = builder.SetEffectID(effID)
		}
		// Recipe edge
		if oldRecipeID := edgeID(t["edges"], "recipe"); oldRecipeID > 0 {
			if newRecipeID, ok := maps.Recipes[oldRecipeID]; ok {
				builder = builder.SetRecipeID(newRecipeID)
			}
		}
		// Dialog node edge
		if dnID := edgeStrID(t["edges"], "dialog_node"); dnID != "" {
			if newID, ok := maps.DialogNodes[dnID]; ok {
				builder = builder.SetDialogNodeID(newID)
			} else {
				builder = builder.SetDialogNodeID(dnID)
			}
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("trigger %d: %w", oldID, err)
		}
		maps.Triggers[oldID] = created.ID
		count++
	}
	return count, nil
}

// importSocialCommands creates SocialCommand records in the new world.
func importSocialCommands(ctx context.Context, client *db.Client, cmds []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, c := range cmds {
		oldID := intVal(c["id"])
		builder := client.SocialCommand.Create().
			SetWorldID(newWorldID).
			SetName(strVal(c["name"])).
			SetDisplayName(strVal(c["display_name"])).
			SetSelfText(strVal(c["self_text"])).
			SetRoomText(strVal(c["room_text"])).
			SetTargetSelfText(strVal(c["target_self_text"])).
			SetTargetText(strVal(c["target_text"])).
			SetTargetRoomText(strVal(c["target_room_text"])).
			SetRequiresTarget(boolValOr(c["requires_target"], false)).
			SetIsEmote(boolValOr(c["is_emote"], false))
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("social_command %d: %w", oldID, err)
		}
		maps.SocialCommands[oldID] = created.ID
		count++
	}
	return count, nil
}

// importShopTemplates creates ShopTemplate records, remapping npc_template edge.
func importShopTemplates(ctx context.Context, client *db.Client, shops []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, s := range shops {
		oldID := intVal(s["id"])
		builder := client.ShopTemplate.Create().
			SetName(strVal(s["name"])).
			SetWorldID(newWorldID)
		// Remap npc_template_id
		tmplSlug := strVal(s["npc_template_id"])
		if tmplSlug != "" {
			if newSlug, ok := maps.NPCTemplates[tmplSlug]; ok {
				builder = builder.SetNpcTemplateID(newSlug)
			} else {
				builder = builder.SetNpcTemplateID(tmplSlug)
			}
		}
		if v := intVal(s["currency_item_type"]); v != 0 {
			builder = builder.SetCurrencyItemType(v)
		}
		if v := intVal(s["max_inventory"]); v != 0 {
			builder = builder.SetMaxInventory(v)
		}
		if v := intVal(s["gold_reserves"]); v != 0 {
			builder = builder.SetGoldReserves(v)
		}
		builder = builder.SetIsActive(boolValOr(s["is_active"], true))
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("shop_template %d: %w", oldID, err)
		}
		maps.ShopTemplates[oldID] = created.ID
		count++
	}
	return count, nil
}

// importEffectHooks creates EffectHook records, remapping npc_template edge.
func importEffectHooks(ctx context.Context, client *db.Client, hooks []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	count := 0
	for _, h := range hooks {
		oldID := intVal(h["id"])
		builder := client.EffectHook.Create().
			SetWorldID(newWorldID).
			SetName(strVal(h["name"])).
			SetEvent(strVal(h["event"])).
			SetEnabled(boolValOr(h["enabled"], true))
		if v := strVal(h["target"]); v != "" {
			builder = builder.SetTarget(v)
		}
		if v := strVal(h["condition"]); v != "" {
			builder = builder.SetCondition(v)
		}
		// Effect edge — global ID, no remap
		if effID := edgeID(h["edges"], "effect"); effID > 0 {
			builder = builder.SetEffectID(effID)
		}
		// Remap npc_template edge
		tmplSlug := edgeStrID(h["edges"], "npc_template")
		if tmplSlug == "" {
			tmplSlug = strVal(h["npc_template_id"])
		}
		if tmplSlug != "" {
			if newSlug, ok := maps.NPCTemplates[tmplSlug]; ok {
				builder = builder.SetNpcTemplateID(newSlug)
			} else {
				builder = builder.SetNpcTemplateID(tmplSlug)
			}
		}
		created, err := builder.Save(ctx)
		if err != nil {
			return count, fmt.Errorf("effect_hook %d: %w", oldID, err)
		}
		maps.EffectHooks[oldID] = created.ID
		count++
	}
	return count, nil
}

// --- Helpers ---

func boolValOr(v interface{}, def bool) bool {
	if v == nil {
		return def
	}
	return boolVal(v)
}

// parseStatBonuses converts a JSON map to schema.StatBonuses.
func parseStatBonuses(v interface{}) *schema.StatBonuses {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var sb schema.StatBonuses
	if json.Unmarshal(b, &sb) != nil {
		return nil
	}
	return &sb
}

// parseSpecialties converts a JSON array to []schema.ClassSpecialty.
func parseSpecialties(v interface{}) *[]schema.ClassSpecialty {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var sp []schema.ClassSpecialty
	if json.Unmarshal(b, &sp) != nil {
		return nil
	}
	return &sp
}

// Ensure unused imports are referenced
var _ = npctemplate.WorldIDEQ