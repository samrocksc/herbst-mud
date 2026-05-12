# RFC-008: Config System Revamp

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-09
**Related:** ADMIN-UX tickets #308, #309, #316 (config page improvements)

---

## 1. Executive Summary

The current `GameConfig` system is a bare key/value string store with no type
awareness, no descriptions, and no documentation on how config keys tie into
game systems. This RFC proposes expanding the schema with metadata (type,
description, category, defaults) and building a first-class config admin page
that makes the config system the **control panel** for the entire game engine.

---

## 2. Current State

### 2.1 Schema

```go
// server/db/schema/game_config.go — current
GameConfig {
    key    string (unique)  // "fountain_room_id", "xp_thresholds"
    value  string            // "5", '{"1":200,"2":400}'
}
```

### 2.2 Admin UI

`admin/src/routes/_auth/config.tsx` — basic key/value table with:
- List of keys and values
- Edit inline or modal
- No type-aware editors (JSON strings shown as raw text, numbers as strings)
- No descriptions, no categories, no documentation

### 2.3 Known Config Keys (from codebase audit)

| Key | Type | Used By | Description |
|-----|------|---------|-------------|
| `fountain_room_id` | int | `dbinit/races_genders.go` | Room ID for the character creation fountain |
| `xp_thresholds` | JSON | `server/services/xp.go` | Level → XP required mapping |
| `starting_room_id` | int | `herbst/game_room.go` | Default spawn room (fallback to root room) |
| `death_penalty_xp_percent` | int | `server/routes/character_routes.go` | % XP lost on death |
| `death_penalty_currency_percent` | int | death penalty system | % currency lost on death |
| `death_penalty_item_drop` | bool | death penalty system | Whether items drop on death |
| `death_penalty_sickness_duration` | int | death penalty system | Sickness duration in seconds |
| `skill_max_level` | int | skill system | Global skill level cap |

### 2.4 How Config Is Consumed

```go
// Pattern 1: Direct query by key
cfg, _ := client.GameConfig.Query().Where(gameconfig.KeyEQ("fountain_room_id")).Only(ctx)
roomID, _ := strconv.Atoi(cfg.Value)

// Pattern 2: JSON unmarshal
cfg, _ := client.GameConfig.Query().Where(gameconfig.KeyEQ("xp_thresholds")).Only(ctx)
var thresholds map[int]int
json.Unmarshal([]byte(cfg.Value), &thresholds)
```

Every consumer manually parses strings — no type safety, no defaults, no validation.

---

## 3. Proposed Schema

### 3.1 Expanded GameConfig

```go
// server/db/schema/game_config.go — proposed
GameConfig {
    key          string (unique)  // "fountain_room_id"
    value        string            // stored value
    value_type   enum               // string | int | float | bool | json | room_id | item_id | npc_id
    default_value string            // "5" — fallback when value is unset
    description  string             // "Room where new characters are created at the fountain"
    category     string             // "world" | "combat" | "characters" | "economy" | "skills" | "death" | "logging"
    is_required  bool (default: false)  // engine fails startup if unset
    updated_at   timestamp          // last modified
}
```

### 3.2 All Config Keys — Full Catalog (Documented)

**World / Navigation**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `fountain_room_id` | room_id | null | Room for character creation fountain |
| `starting_room_id` | room_id | null | Default spawn room (falls back to root room if unset) |
| `root_room_id` | room_id | auto | The root room (managed automatically, read-only in config) |
| `default_world` | string | "default" | Default world for new connections |
| `max_rooms` | int | 10000 | Hard limit on room count |

**Characters / Progression**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `xp_thresholds` | json | `{"1":200,"2":400,...}` | Level → XP required to advance |
| `xp_multiplier` | float | 1.0 | Global XP multiplier |
| `max_level` | int | 50 | Hard level cap |
| `starting_xp` | int | 0 | XP for new characters |
| `starting_level` | int | 1 | Level for new characters |
| `starting_hp` | int | 100 | HP for new characters |
| `starting_stamina` | int | 50 | Stamina for new characters |
| `starting_mana` | int | 25 | Mana for new characters |
| `max_characters_per_user` | int | 3 | Character slot limit per account |

**Combat**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `combat_tick_interval_ms` | int | 1500 | Combat tick interval in milliseconds |
| `base_hit_chance` | float | 0.50 | Base hit chance before modifiers |
| `crit_multiplier` | float | 2.0 | Critical hit damage multiplier |
| `fumble_chance` | float | 0.05 | Natural 1 probability (always 1/20, informational) |
| `flee_dc` | int | 12 | Difficulty check to flee combat |
| `npc_aggro_range` | int | 0 | Rooms of proximity for NPC aggro (0 = same room only) |

**Death**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `death_penalty_xp_percent` | int | 0 | % of current XP lost on death |
| `death_penalty_currency_percent` | int | 0 | % of currency lost on death |
| `death_penalty_item_drop` | bool | false | Drop items on death? |
| `death_penalty_sickness_duration` | int | 0 | Post-respawn sickness in seconds |
| `corpse_rot_seconds` | int | 600 | Seconds before corpse decays |
| `respawn_room_id` | room_id | null | Global respawn room (overridden by character bind point) |

**Skills / Abilities**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `skill_max_level` | int | 10 | Global skill level cap |
| `skill_xp_multiplier` | float | 1.0 | Global skill XP multiplier |
| `ability_slots` | int | 5 | Number of ability slots per character |
| `default_ability_cooldown_seconds` | int | 30 | Default cooldown for abilities without specific config |

**Economy**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `currency_name` | string | "gold" | Display name for currency |
| `starting_currency` | int | 0 | Currency for new characters |
| `sell_price_multiplier` | float | 0.5 | Sell price as fraction of buy price |
| `npc_drop_multiplier` | float | 1.0 | Global loot drop rate multiplier |

**Logging / Debug**

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `log_min_level` | string | "INFO" | Minimum log level persisted to DB |
| `log_retention_days` | int | 3 | Days to keep logs before cleanup |
| `debug_enabled` | bool | false | Global debug mode (verbose logging) |

---

## 4. Admin UI — Config Page Revamp

### 4.1 Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  Config                                [+ Add Config Key]       │
│                                                                  │
│  [All] [World] [Characters] [Combat] [Death] [Skills] [Economy] │
│                                                                  │
│  🔍 Search config keys...                                        │
│                                                                  │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ WORLD & NAVIGATION                                 5 keys    │ │
│ │                                                              │ │
│ │ fountain_room_id     room_id    5              "The Fountain" │ │
│ │ starting_room_id     room_id    1              "Crossroads"   │ │
│ │ root_room_id         room_id    1    (auto)    "Crossroads"   │ │
│ │ default_world        string     "default"                     │ │
│ │ max_rooms            int        10000                         │ │
│ ├─────────────────────────────────────────────────────────────┤ │
│ │ CHARACTERS & PROGRESSION                           4 keys    │ │
│ │ ...                                                          │ │
│ └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Edit Config Key (slide-over or inline)

```
┌──────────────────────────────────────────────┐
│  Edit: xp_thresholds                          │
│                                               │
│  Key:     xp_thresholds                       │
│  Type:    [json ▾]                            │
│  Category: [characters ▾]                     │
│  Default:  {"1": 200, "2": 400, ...}          │
│  Required: [✓]                                │
│                                               │
│  Description:                                  │
│  ┌──────────────────────────────────────────┐ │
│  │ Maps character level to XP required to    │ │
│  │ advance. Used by XPAwardService.          │ │
│  │ Example: {"1": 200, "2": 400, "3": 800}   │ │
│  └──────────────────────────────────────────┘ │
│                                               │
│  Value:                                        │
│  ┌──────────────────────────────────────────┐ │
│  │ {                                        │ │
│  │   "1": 200,                              │ │
│  │   "2": 400,                              │ │
│  │   "3": 800,                              │ │
│  │   ...                                    │ │
│  │ }                                        │ │
│  └──────────────────────────────────────────┘ │
│                                               │
│  Used by: server/services/xp.go (line 129)     │
│                                               │
│  [Save]  [Reset to Default]  [Cancel]         │
└──────────────────────────────────────────────┘
```

### 4.3 Type-Aware Editors

| Type | Editor |
|------|--------|
| `string` | Text input |
| `int` | Number input (min/max from metadata) |
| `float` | Number input with decimal |
| `bool` | Toggle switch |
| `json` | Code editor with syntax highlighting + validation |
| `room_id` | Dropdown with room name lookup |
| `item_id` | Dropdown with item name lookup |
| `npc_id` | Dropdown with NPC name lookup |

---

## 5. Runtime Integration

### 5.1 Config Service (new)

```go
// server/services/config.go
type ConfigService struct {
    client *db.Client
    cache  map[string]GameConfig  // loaded at startup, refreshed on write
    mu     sync.RWMutex
}

func (s *ConfigService) GetString(key string) (string, error)
func (s *ConfigService) GetInt(key string) (int, error)
func (s *ConfigService) GetFloat(key string) (float64, error)
func (s *ConfigService) GetBool(key string) (bool, error)
func (s *ConfigService) GetJSON(key string, target interface{}) error
func (s *ConfigService) GetRoomID(key string) (int, error)
```

Every consumer that currently does `client.GameConfig.Query().Where(...)` 
manually should use `ConfigService` instead. The service handles:
- Caching (single DB query at startup, invalidated on write)
- Type parsing (no more `strconv.Atoi` scattered everywhere)
- Default fallback (returns default_value if key is unset AND not required)
- Required-key validation at startup

### 5.2 Startup Validation

```go
func (s *ConfigService) ValidateRequired() error {
    for _, key := range s.RequiredKeys() {
        cfg, err := s.client.GameConfig.Query().Where(gameconfig.KeyEQ(key)).Only(ctx)
        if err != nil || cfg.Value == "" {
            return fmt.Errorf("required config key %q is unset", key)
        }
    }
    return nil
}
```

Called in `main.go` after DB init, before route registration. Engine fails fast
if `is_required=true` keys are missing.

---

## 6. Migration Path

### Phase 1: Schema + Service (backend only)
- Add `value_type`, `default_value`, `description`, `category`, `is_required`, `updated_at` to `GameConfig` ent schema
- Create `ConfigService` with typed getters
- Migrate existing consumers from raw queries → ConfigService
- Seed default config keys in `dbinit/`
- Backward compatible: `value` field unchanged, existing `key`/`value` queries still work

### Phase 2: Admin UI Revamp
- Replace `admin/src/routes/_auth/config.tsx` with category-grouped, type-aware editor
- Each config key shows: current value, type, description, where it's used in code
- Type-appropriate editors (toggle for bool, JSON viewer for json, room picker for room_id)
- Search + category filter tabs

### Phase 3: Documentation
- Generate a config reference doc page at `/docs/config-reference` (or update the existing one)
- Include all keys, types, defaults, descriptions, and system linkages
- Auto-generate from the schema (or source of truth)

---

## 7. Existing ADMIN-UX Tickets Superseded

This RFC consolidates and expands:

| Ticket | Scope | Superseded By |
|--------|-------|---------------|
| #308 | ADMIN-UX-001: Config page needs formatted JSON viewer | This RFC — type-aware editors for all types |
| #309 | ADMIN-UX-002: Config page should show human-readable descriptions | This RFC — `description` field per key |
| #316 | ADMIN-UX-005: Research collapsible JSON formatter | This RFC — native JSON editor |

---

## 8. Acceptance Criteria

- [ ] `GameConfig` schema expanded with all new fields
- [ ] `ConfigService` with typed getters, caching, and startup validation
- [ ] All existing raw config queries migrated to ConfigService
- [ ] Default config keys seeded in `dbinit/`
- [ ] Admin config page revamped with categories, type-aware editors, descriptions
- [ ] JSON values use syntax-highlighted editor
- [ ] Room/item/NPC ID fields use lookup dropdowns
- [ ] Each config key documents which system consumes it
- [ ] `docs/config-reference` page updated with full key catalog
- [ ] `npm run build` passes
