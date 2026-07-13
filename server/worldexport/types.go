// Package worldexport provides world-scoped JSON export/import for herbst-mud.
package worldexport

// CurrentVersion is the export format version accepted by ImportWorld.
const CurrentVersion = "2.0"

// WorldSnapshot is the serialized form of an entire game world.
type WorldSnapshot struct {
	Version     string                   `json:"version"`
	ExportedAt  string                   `json:"exported_at"`
	World       map[string]interface{}   `json:"world"`
	Rooms       []map[string]interface{} `json:"rooms"`
	NPCs        []map[string]interface{} `json:"npcs"`
	Abilities   []map[string]interface{} `json:"abilities"`
	Effects     []map[string]interface{} `json:"effects"`
	NPCTemplates []map[string]interface{} `json:"npc_templates"`
	Races       []map[string]interface{} `json:"races"`
	Genders     []map[string]interface{} `json:"genders"`
	Factions    []map[string]interface{} `json:"factions"`
	FactionCategories []map[string]interface{} `json:"faction_categories"`
	Tags        []map[string]interface{} `json:"tags"`
	Equipment   []map[string]interface{} `json:"equipment"`
	CharacterAbilities []map[string]interface{} `json:"character_abilities"`
	Recipes     []map[string]interface{} `json:"recipes"`
	Quests      []map[string]interface{} `json:"quests"`
	Triggers    []map[string]interface{} `json:"triggers"`
	DialogNodes []map[string]interface{} `json:"dialog_nodes"`
	EffectHooks []map[string]interface{} `json:"effect_hooks"`
	SocialCommands []map[string]interface{} `json:"social_commands"`
	ShopTemplates []map[string]interface{} `json:"shop_templates"`
	Zones        []map[string]interface{} `json:"zones"`
	Skills        []map[string]interface{} `json:"skills"`
}

// ImportResult reports what was created during import.
type ImportResult struct {
	WorldID             int `json:"world_id"`
	Rooms               int `json:"rooms"`
	NPCs                int `json:"npcs"`
	Abilities           int `json:"abilities"`
	Effects             int `json:"effects"`
	NPCTemplates        int `json:"npc_templates"`
	Races               int `json:"races"`
	Genders             int `json:"genders"`
	Factions            int `json:"factions"`
	FactionCategories   int `json:"faction_categories"`
	Tags                int `json:"tags"`
	Equipment           int `json:"equipment"`
	CharacterAbilities  int `json:"character_abilities"`
	Recipes             int `json:"recipes"`
	Quests              int `json:"quests"`
	Triggers            int `json:"triggers"`
	DialogNodes         int `json:"dialog_nodes"`
	EffectHooks         int `json:"effect_hooks"`
	SocialCommands      int `json:"social_commands"`
	ShopTemplates       int `json:"shop_templates"`
	Zones               int `json:"zones"`
	Skills              int `json:"skills"`
}
