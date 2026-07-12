🔵 feat: DB-driven classes + world export/import parity

## Remove Hardcoded Class Constants

All class-specific data has been removed from server/constants/ and is
now fully DB-driven via the factions table:

### Schema Changes
- Added `stat_bonuses` (JSON StatBonuses) and `specialties` (JSON
  []ClassSpecialty) fields to the Faction entity
- Ent code regenerated
- DB migration applied automatically on server startup

### Removed Constants
- server/constants/class_specialties.go — DELETED entirely
  (ClassSpecialties, StartingConfigs, GetClassConfig, GetSpecialty,
  ClassConfig, ClassSpecialty struct)
- server/constants/class_specialties_test.go — DELETED
- server/constants/warrior_fighter_test.go — DELETED
- server/constants/character.go — Removed ValidClasses, IsValidClass,
  ClassStatBonuses. Kept ValidRaces, DefaultStats, IsValidSlot.

### Updated Callers (DB-driven)
- server/service/character.go: IsValidClass → DB faction query,
  GetClassConfig stat bonuses → faction.StatBonuses from DB,
  specialty → first specialty ID from faction.Specialties
- server/routes/character_classes.go: Removed constants fallback,
  getSpecialties reads from faction.Specialties JSON
- server/routes/character_stats.go: getSpecialtiesForClass queries
  faction from DB, getClassFactionByName helper
- server/routes/ws_dialog_effects.go: DB validation + specialty lookup

### DB Data Populated
All 10 factions (8 world-1 classes + 2 world-2 classes) now have:
- stat_bonuses: STR/DEX/CON/INT/WIS bonuses per class
- specialties: JSON array of {id, name, description}

## World Export/Import Parity

### Export Fixes
- Abilities: Added .WithFaction() to load faction edge
- Effects: Added .WithAbility() to load ability edge
- Factions: Added .WithCategory() to load category edge
- Triggers: Added .WithEffect().WithRecipe().WithDialogNode()
- DialogNodes: Added .WithNpcTemplate()
- EffectHooks: Added .WithEffect().WithNpcTemplate()

### Import Functions (new file: import_full.go)
13 new import functions covering all entity types:
- importFactionCategories, importFactions
- importAbilities, importEffects
- importTags, importEquipment
- importNPCTemplates, importDialogNodes
- importQuests, importRecipes
- importTriggers, importSocialCommands
- importShopTemplates, importEffectHooks

### Import Fixes
- Globally-unique fields suffixed with world ID (faction names, quest
  names, NPC template IDs/slugs, dialog node IDs)
- Enum validation with fallback to defaults (disposition, roam_pattern,
  repeat_mode, main_type)
- Required fields always set (NPC template skills, greeting,
  respawn_rooms, trades_with, social command fields)
- Edge FK remapping via idMaps (faction→category, ability→faction,
  effect→ability, trigger→room/equipment/recipe/dialog_node,
  dialog_node→npc_template, effect_hook→effect/npc_template)

### Round-Trip Verified
Export world 2 → Import as "Final Test" (world 16):
- 47 rooms ✅, 2 NPCs ✅, 25 abilities ✅, 17 effects ✅
- 9 NPC templates ✅, 2 factions ✅, 1 faction category ✅
- 1 tag ✅, 12 equipment ✅, 3 quests ✅, 1 trigger ✅, 1 dialog node ✅

Files:
  server/db/schema/faction.go           +StatBonuses, +ClassSpecialty, +2 JSON fields
  server/db/faction.go                  +generated code for new fields
  server/db/faction_create.go           +SetStatBonuses, +SetSpecialties
  server/db/faction_update.go           +update builders
  server/db/migrate/schema.go           +migration columns
  server/constants/character.go         -ValidClasses, -IsValidClass, -ClassStatBonuses
  server/constants/class_specialties.go DELETED
  server/constants/class_specialties_test.go DELETED
  server/constants/warrior_fighter_test.go DELETED
  server/service/character.go           DB-driven class validation + stat bonuses
  server/routes/character_stats.go      DB-driven specialties, getClassFactionByName
  server/worldexport/export.go          +edge loading for all entities
  server/worldexport/import.go          +11 new import function calls
  server/worldexport/import_full.go     NEW: 13 import functions
  server/worldexport/import_utils.go    +idMaps fields, +edgeID/edgeStrID helpers
  server/worldexport/types.go           +ImportResult fields
  features/player-class-abilities.feature  5 Gherkin scenarios