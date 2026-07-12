package worldexport

import (
	"context"
	"fmt"
	"strconv"

	"herbst-server/db"
)

// ImportWorld creates a new world and imports the full snapshot into it.
func ImportWorld(ctx context.Context, client *db.Client, snap *WorldSnapshot, newName, newTitle string) (*ImportResult, error) {
	if snap.World == nil {
		return nil, fmt.Errorf("snapshot missing world metadata")
	}
	if newTitle == "" {
		newTitle = newName
	}

	w, err := client.World.Create().
		SetName(newName).
		SetTitle(newTitle).
		SetNillableDescription(strPtr(snap.World["description"])).
		SetActive(false).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create world: %w", err)
	}
	newWorldID := strconv.Itoa(w.ID)
	res := &ImportResult{WorldID: w.ID}

	maps := newIDMaps()

	// Phase 1: Foundation entities (no FK deps)
	res.Races, err = importRaces(ctx, client, snap.Races, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import races: %w", err)
	}

	res.Genders, err = importGenders(ctx, client, snap.Genders, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import genders: %w", err)
	}

	// Phase 2: Rooms + Zones (rooms may reference zones)
	res.Rooms, err = importRooms(ctx, client, snap.Rooms, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import rooms: %w", err)
	}

	res.Zones, err = importZones(ctx, client, snap.Zones, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import zones: %w", err)
	}

	// Phase 3: Faction categories → Factions (factions depend on categories)
	res.FactionCategories, err = importFactionCategories(ctx, client, snap.FactionCategories, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import faction_categories: %w", err)
	}

	res.Factions, err = importFactions(ctx, client, snap.Factions, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import factions: %w", err)
	}

	// Phase 4: Tags (no deps)
	res.Tags, err = importTags(ctx, client, snap.Tags, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import tags: %w", err)
	}

	// Phase 5: Abilities → Effects (effects depend on abilities)
	res.Abilities, err = importAbilities(ctx, client, snap.Abilities, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import abilities: %w", err)
	}

	res.Effects, err = importEffects(ctx, client, snap.Effects, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import effects: %w", err)
	}

	// Phase 6: Equipment templates (no deps beyond world)
	res.Equipment, err = importEquipment(ctx, client, snap.Equipment, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import equipment: %w", err)
	}

	// Phase 7: NPC templates → Dialog nodes (dialog nodes reference npc templates)
	res.NPCTemplates, err = importNPCTemplates(ctx, client, snap.NPCTemplates, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import npc_templates: %w", err)
	}

	res.DialogNodes, err = importDialogNodes(ctx, client, snap.DialogNodes, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import dialog_nodes: %w", err)
	}

	// Phase 8: NPCs (depend on races, rooms)
	res.NPCs, err = importNPCs(ctx, client, snap.NPCs, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import npcs: %w", err)
	}

	// Phase 9: Quests, Recipes (no FK deps beyond world)
	res.Quests, err = importQuests(ctx, client, snap.Quests, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import quests: %w", err)
	}

	res.Recipes, err = importRecipes(ctx, client, snap.Recipes, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import recipes: %w", err)
	}

	// Phase 10: Triggers (depend on rooms, equipment, recipes, dialog nodes)
	res.Triggers, err = importTriggers(ctx, client, snap.Triggers, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import triggers: %w", err)
	}

	// Phase 11: Social commands, shop templates, effect hooks
	res.SocialCommands, err = importSocialCommands(ctx, client, snap.SocialCommands, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import social_commands: %w", err)
	}

	res.ShopTemplates, err = importShopTemplates(ctx, client, snap.ShopTemplates, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import shop_templates: %w", err)
	}

	res.EffectHooks, err = importEffectHooks(ctx, client, snap.EffectHooks, newWorldID, maps)
	if err != nil {
		return nil, fmt.Errorf("import effect_hooks: %w", err)
	}

	return res, nil
}