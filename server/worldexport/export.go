package worldexport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"herbst-server/db"
	"herbst-server/db/ability"
	"herbst-server/db/abilityeffect"
	"herbst-server/db/character"
	"herbst-server/db/characterability"
	"herbst-server/db/craftingrecipe"
	"herbst-server/db/dialognode"
	"herbst-server/db/effecthook"
	"herbst-server/db/equipmenttemplate"
	"herbst-server/db/faction"
	"herbst-server/db/factioncategory"
	"herbst-server/db/gender"
	"herbst-server/db/npctemplate"
	"herbst-server/db/quest"
	"herbst-server/db/race"
	"herbst-server/db/room"
	"herbst-server/db/shoptemplate"
	"herbst-server/db/skill"
	"herbst-server/db/socialcommand"
	"herbst-server/db/tag"
	"herbst-server/db/trigger"
	"herbst-server/db/world"
)

// ExportWorld exports a single world as a JSON snapshot.
func ExportWorld(ctx context.Context, client *db.Client, worldID string) (*WorldSnapshot, error) {
	w, err := client.World.Query().Where(world.IDEQ(atoi(worldID))).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch world: %w", err)
	}

	snap := &WorldSnapshot{
		Version:    "2.0",
		ExportedAt: time.Now().Format(time.RFC3339),
		World:      worldToMap(w),
	}

	snap.Rooms, err = loadByWorld(ctx, client.Room.Query().Where(room.WorldIDEQ(worldID)).All, roomSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.NPCs, err = loadByWorld(ctx, client.Character.Query().Where(character.IsNPCEQ(true), character.CurrentWorldEQ(worldID)).All, charSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Abilities, err = loadByWorld(ctx, client.Ability.Query().Where(ability.WorldIDEQ(worldID)).WithFaction().All, abilitySliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Effects, err = loadEffectsForWorld(ctx, client, worldID)
	if err != nil {
		return nil, err
	}
	snap.NPCTemplates, err = loadByWorld(ctx, client.NPCTemplate.Query().Where(npctemplate.WorldIDEQ(worldID)).All, npcTemplateSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Races, err = loadByWorld(ctx, client.Race.Query().Where(race.WorldIDEQ(worldID)).All, raceSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Genders, err = loadByWorld(ctx, client.Gender.Query().Where(gender.WorldIDEQ(worldID)).All, genderSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Factions, err = loadByWorld(ctx, client.Faction.Query().Where(faction.WorldIDEQ(worldID)).WithCategory().All, factionSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.FactionCategories, err = loadByWorld(ctx, client.FactionCategory.Query().Where(factioncategory.WorldIDEQ(worldID)).All, factionCategorySliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Tags, err = loadByWorld(ctx, client.Tag.Query().Where(tag.WorldIDEQ(worldID)).All, tagSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Equipment, err = loadByWorld(ctx, client.EquipmentTemplate.Query().Where(equipmenttemplate.WorldIDEQ(worldID)).All, equipmentSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.CharacterAbilities, err = loadCharacterAbilitiesForWorld(ctx, client, worldID)
	if err != nil {
		return nil, err
	}
	snap.Recipes, err = loadByWorld(ctx, client.CraftingRecipe.Query().Where(craftingrecipe.WorldIDEQ(worldID)).All, recipeSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Quests, err = loadByWorld(ctx, client.Quest.Query().Where(quest.WorldIDEQ(worldID)).All, questSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Triggers, err = loadByWorld(ctx, client.Trigger.Query().Where(trigger.WorldIDEQ(worldID)).WithEffect().WithRecipe().WithDialogNode().All, triggerSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.DialogNodes, err = loadByWorld(ctx, client.DialogNode.Query().Where(dialognode.WorldIDEQ(worldID)).WithNpcTemplate().All, dialogNodeSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.EffectHooks, err = loadByWorld(ctx, client.EffectHook.Query().Where(effecthook.WorldIDEQ(worldID)).WithEffect().WithNpcTemplate().All, effectHookSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.SocialCommands, err = loadByWorld(ctx, client.SocialCommand.Query().Where(socialcommand.WorldIDEQ(worldID)).All, socialCommandSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.ShopTemplates, err = loadByWorld(ctx, client.ShopTemplate.Query().Where(shoptemplate.WorldIDEQ(worldID)).All, shopTemplateSliceToMaps)
	if err != nil {
		return nil, err
	}
	snap.Zones, err = exportZones(ctx, client, worldID)
	if err != nil {
		return nil, err
	}
	snap.Skills, err = loadByWorld(ctx, client.Skill.Query().Where(skill.WorldIDEQ(atoi(worldID))).All, skillSliceToMaps)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

func loadByWorld[T any](ctx context.Context, fetch func(context.Context) ([]T, error), convert func([]T) []map[string]interface{}) ([]map[string]interface{}, error) {
	items, err := fetch(ctx)
	if err != nil {
		return nil, err
	}
	return convert(items), nil
}

func loadEffectsForWorld(ctx context.Context, client *db.Client, worldID string) ([]map[string]interface{}, error) {
	abilities, err := client.Ability.Query().Where(ability.WorldIDEQ(worldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	abilityIDs := make([]int, len(abilities))
	for i, a := range abilities {
		abilityIDs[i] = a.ID
	}
	if len(abilityIDs) == 0 {
		return []map[string]interface{}{}, nil
	}
	effects, err := client.AbilityEffect.Query().Where(abilityeffect.HasAbilityWith(ability.IDIn(abilityIDs...))).WithAbility().All(ctx)
	if err != nil {
		return nil, err
	}
	return abilityEffectSliceToMaps(effects), nil
}

func loadCharacterAbilitiesForWorld(ctx context.Context, client *db.Client, worldID string) ([]map[string]interface{}, error) {
	npcs, err := client.Character.Query().Where(character.IsNPCEQ(true), character.CurrentWorldEQ(worldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	charIDs := make([]int, len(npcs))
	for i, n := range npcs {
		charIDs[i] = n.ID
	}
	if len(charIDs) == 0 {
		return []map[string]interface{}{}, nil
	}
	cas, err := client.CharacterAbility.Query().Where(characterability.HasCharacterWith(character.IDIn(charIDs...))).All(ctx)
	if err != nil {
		return nil, err
	}
	return characterAbilitySliceToMaps(cas), nil
}

func worldToMap(w *db.World) map[string]interface{} {
	return map[string]interface{}{
		"id":          w.ID,
		"name":        w.Name,
		"title":       w.Title,
		"description": w.Description,
		"active":      w.Active,
	}
}

func toMaps[T any](items []T) []map[string]interface{} {
	out := make([]map[string]interface{}, len(items))
	for i, item := range items {
		b, _ := json.Marshal(item)
		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)
		out[i] = m
	}
	return out
}

func roomSliceToMaps(items []*db.Room) []map[string]interface{}     { return toMaps(items) }
func charSliceToMaps(items []*db.Character) []map[string]interface{} { return toMaps(items) }
func abilitySliceToMaps(items []*db.Ability) []map[string]interface{} { return toMaps(items) }
func abilityEffectSliceToMaps(items []*db.AbilityEffect) []map[string]interface{} { return toMaps(items) }
func npcTemplateSliceToMaps(items []*db.NPCTemplate) []map[string]interface{} { return toMaps(items) }
func raceSliceToMaps(items []*db.Race) []map[string]interface{}       { return toMaps(items) }
func genderSliceToMaps(items []*db.Gender) []map[string]interface{}   { return toMaps(items) }
func factionSliceToMaps(items []*db.Faction) []map[string]interface{} { return toMaps(items) }
func factionCategorySliceToMaps(items []*db.FactionCategory) []map[string]interface{} { return toMaps(items) }
func tagSliceToMaps(items []*db.Tag) []map[string]interface{}         { return toMaps(items) }
func equipmentSliceToMaps(items []*db.EquipmentTemplate) []map[string]interface{} { return toMaps(items) }
func characterAbilitySliceToMaps(items []*db.CharacterAbility) []map[string]interface{} { return toMaps(items) }
func recipeSliceToMaps(items []*db.CraftingRecipe) []map[string]interface{} { return toMaps(items) }
func questSliceToMaps(items []*db.Quest) []map[string]interface{}     { return toMaps(items) }
func triggerSliceToMaps(items []*db.Trigger) []map[string]interface{} { return toMaps(items) }
func dialogNodeSliceToMaps(items []*db.DialogNode) []map[string]interface{} { return toMaps(items) }
func effectHookSliceToMaps(items []*db.EffectHook) []map[string]interface{} { return toMaps(items) }
func socialCommandSliceToMaps(items []*db.SocialCommand) []map[string]interface{} { return toMaps(items) }
func shopTemplateSliceToMaps(items []*db.ShopTemplate) []map[string]interface{} { return toMaps(items) }
func skillSliceToMaps(items []*db.Skill) []map[string]interface{} { return toMaps(items) }
func zoneSliceToMaps(items []*db.Zone) []map[string]interface{} { return toMaps(items) }

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
