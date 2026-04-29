# RFC-XP: Experience Points System

**Status:** Draft
**Author:** Leonardo (PM)
**Created:** 2026-04-28
**Related:** RFC-FACTION-SKILLS (#219, approved)

---

## 1. Overview

This RFC defines an experience points (XP) system for herbst-mud. Characters earn XP by defeating NPCs. XP accumulates, triggers level-ups at configurable thresholds, and is tracked per-character. A separate competencies/proficiencies XP system tracks weapon and armor mastery over time.

All configuration is in the database — no YAML.

---

## 2. Mob XP Values

### 2.1 Content-Defined XP Per NPC Template

Each NPC template (in DB, replacing YAML) carries an `xp_value` field:

| Field | Type | Default | Notes |
|---|---|---|---|
| `xp_value` | int | 0 | Base XP awarded when this NPC is killed |

This means designers have full control: a weak goblin might be 10 XP, a boss dragon might be 1000 XP.

### 2.2 Level-Differential Multiplier

XP is scaled by the difference between the NPC's level and the character's level:

```
XP_gained = npc_template.xp_value × max(1, 1 + (npc_template.level − character_level) × scale_factor)
```

- If NPC level == character level → `×1.0` (full XP)
- If NPC level > character level → `×1.0 +` bonus per level above
- If NPC level < character level → `×1.0` (no penalty — players don't feel punished for outleveling content)
- `scale_factor` is a game config value (default: 0.10 = +10% per level difference). Stored in `CharacterConfig.xp_level_diff_scale_factor`. Admins can tune it without restarting.

**Examples with scale_factor=0.10:**
- Character level 5 kills level 5 NPC (xp_value=50) → `50 × 1.0 = 50 XP`
- Character level 5 kills level 8 NPC (xp_value=50) → `50 × (1 + (8−5)×0.10) = 50 × 1.30 = 65 XP`
- Character level 10 kills level 2 NPC (xp_value=50) → `50 × 1.0 = 50 XP` (no penalty)

The formula is configurable at the game config level.

---

## 3. Character XP and Leveling

### 3.1 Character Fields

Add to the `Character` entity:

| Field | Type | Default | Notes |
|---|---|---|---|
| `xp` | int | 0 | Current accumulated XP |

`level` already exists. `xp` is new. `max_level` is not on Character — it comes from `CharacterConfig`.

### 3.2 Level Cap — `CharacterConfig` Table

Max character level is configurable, not hardcoded. A `CharacterConfig` DB table holds server-wide character settings:

| Field | Type | Default | Notes |
|---|---|---|---|
| `key` | string | PK | Config key name |
| `value` | string | | Config value (JSON or plain) |
| `description` | string | | Admin-facing description |

**Default entries:**

| key | value | description |
|---|---|---|
| `max_level` | `50` | Maximum achievable character level |
| `xp_level_diff_scale_factor` | `0.10` | XP formula scale factor (+10% per level difference) |
| `death_penalty_enabled` | `false` | Toggle death penalty system |
| `death_penalty_xp_percent` | `0.0` | % of current-level XP lost on death |
| `death_penalty_currency_percent` | `0.0` | % of gold lost on death |
| `death_penalty_item_drop` | `false` | Drop one random equipped item on death |
| `death_penalty_sickness_duration` | `0` | Resurrection sickness duration (seconds) |
| `quest_xp_enabled` | `true` | Award XP on quest completion |

Admins can edit these via the admin UI without restarting. The event bus and XP system read these values at runtime — no restart required.

### 3.3 Level-Up Thresholds

Level-up thresholds are stored in a new `LevelThreshold` DB table (game config):

| Field | Type | Notes |
|---|---|---|
| `level` | int | Target level (2, 3, … N). Unique. Primary key. |
| `xp_required` | int | XP needed to reach this level |

**Default threshold table:**

| Level | XP Required |
|---|---|
| 2 | 100 |
| 3 | 300 |
| 4 | 600 |
| 5 | 1,000 |
| 6 | 1,500 |
| 7 | 2,100 |
| 8 | 2,800 |
| 9 | 3,600 |
| 10 | 4,500 |
| 11 | 5,500 |
| 12 | 6,600 |
| 13 | 7,800 |
| 14 | 9,100 |
| 15 | 10,500 |
| … | … |

Formula (configurable): `XP_n = 50 × n × (n − 1)` for level `n`.

### 3.4 Level-Up Process

When XP is awarded to a character:
1. Add XP to `character.xp`
2. Look up `max_level` from `CharacterConfig`
3. Check if `character.level >= max_level` — if so, stop (capped)
4. Check if `character.xp >= LevelThreshold[level + 1].xp_required`
5. If yes: increment `character.level`, apply stat bonuses, emit `character_leveled_up` event
6. Repeat until no further level-up qualifies

### 3.5 Stat Bonuses Per Level

A new `ClassLevelBonus` DB table maps `(class, level) → stat bonuses`:

| Field | Type | Notes |
|---|---|---|
| `class` | string | e.g., "warrior", "survivor" |
| `level` | int | Level reached |
| `hp_bonus` | int | Extra max HP at this level |
| `mana_bonus` | int | Extra max mana at this level |
| `stamina_bonus` | int | Extra max stamina at this level |
| `stat_bonus_str` | int | Bonus STR at this level |
| `stat_bonus_dex` | int | Bonus DEX at this level |
| `stat_bonus_con` | int | Bonus CON at this level |
| `stat_bonus_int` | int | Bonus INT at this level |
| `stat_bonus_wis` | int | Bonus WIS at this level |

On level-up: apply the bonuses for the new level, update `max_hitpoints`, `max_mana`, `max_stamina`, and stat fields.

---

## 4. Competencies (Proficiencies XP)

### 4.1 Overview

Separate from character XP, competencies track mastery in specific weapon/armor categories. As a character uses weapons and armor, they earn competency XP in those categories. At certain thresholds, competency bonuses apply (e.g., +5% damage with blades at competency level 3).

This is distinct from the existing flat `skill_blades`, `skill_staves` etc. fields on Character — those are static proficiency levels set at character creation. The new competencies system is a dynamic XP-based progression tracked over time.

### 4.2 Competency Categories

Defined in a `CompetencyCategory` DB table:

| Field | Type | Notes |
|---|---|---|
| `id` | string | Unique ID: "blades", "staves", "fire_magic", "light_armor" |
| `name` | string | Display name: "Blades", "Staves", "Fire Magic" |
| `stat_key` | string | Which character field tracks XP: "skill_blades", "skill_staves" |
| `bonus_per_level` | string | JSON: {stat: "strength", per_point: 0.05} |

The `stat_key` maps to the existing flat proficiency fields on Character (e.g., a character's `skill_blades` value is their competency level in blades).

### 4.3 Competency XP — Data Model

Use a `CharacterCompetency` join table for extensibility (vs. flat columns per category):

| Field | Type | Notes |
|---|---|---|
| `id` | int PK | |
| `character_id` | int FK | Character |
| `category_id` | string FK | CompetencyCategory |
| `xp` | int | Current XP in this category |
| `level` | int | Derived from XP (cached, recomputed on award) |

`level` is derived from `xp` on each award — it is a cached computation, not independently stored. This makes level lookups O(1) without table scans.

> **Note:** A join table is more extensible than flat columns — new categories require only a `CompetencyCategory` DB entry, not a schema migration. This matches the "simple but extensible" guideline.

### 4.4 Competency Level Thresholds

A `CompetencyLevelThreshold` DB table:

| Field | Type | Notes |
|---|---|---|
| `category_id` | string | FK |
| `level` | int | Competency level (1, 2, 3…) |
| `xp_required` | int | XP needed |
| `damage_multiplier` | float | e.g., 1.05 = +5% damage |
| `defense_multiplier` | float | e.g., 1.03 = +3% defense |

**Example thresholds for blades:**

| Competency Level | XP Required | Damage Multiplier |
|---|---|---|
| 1 | 0 | 1.00 |
| 2 | 100 | 1.03 |
| 3 | 300 | 1.05 |
| 4 | 600 | 1.08 |
| 5 | 1,000 | 1.10 |

### 4.5 Earning Competency XP

When a character defeats a NPC:
1. Calculate NPC XP (Section 2)
2. For each weapon/armor category the character used in combat, award `NPC_xp × 0.20` competency XP in that category

Example: Character kills a level 5 NPC worth 50 base XP. Character used blades in combat. Award `50 × 0.20 = 10` competency XP in blades.

### 4.6 Applying Competency Bonuses

During combat damage calculation:
- Look up the character's competency level for the weapon used
- Apply `damage_multiplier` from `CompetencyLevelThreshold`

---

## 5. Passive/Proc Skills

Skills (from RFC-FACTION-SKILLS) are one-shot unlocks with damage multipliers. The competencies system described here is a separate, parallel XP-based progression. They do not conflict — a character can have both competency bonuses (from XP) and faction skills (from one-shot unlocks).

---

## 6. REST API Endpoints

### 6.1 Architecture — Handlers, Services, and Events

**Handlers** receive HTTP requests, parse input, and call **service functions**. **Service functions** contain business logic, update the database, and **emit events** to the event bus. Handlers do NOT emit events directly — this keeps the HTTP layer separate from game logic and makes the service callable from anywhere (HTTP, CLI commands, the event bus itself, tests).

```
HTTP Request → Handler → Service Function → DB Update
                                          → Event Emission → Event Bus → Subscribers
```

**Subscribers** listen to the event bus and trigger side effects. For example, the `OnNPCKilled` subscriber listens for `npc_killed` events, calculates XP, and calls the XP award service.

This also means NPC kill events trigger XP naturally — no direct coupling between combat and the XP system.

### 6.2 XP Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/xp` | Get XP and level |
| POST | `/characters/:id/xp/award` | Award XP to character (admin/system) |
| GET | `/characters/:id/xp/level-threshold` | Get next level threshold |
| GET | `/characters/:id/xp/history` | XP gain history (optional, for player view) |

**Award XP request:**
```json
{ "xp": 50, "source": "NPC_kill", "NPC_id": 12, "NPC_name": "Goblin" }
```

**Response:**
```json
{
  "xp_gained": 50,
  "total_xp": 150,
  "level": 2,
  "leveled_up": true,
  "new_stats": { "max_hitpoints": 112, "max_mana": 54, "max_stamina": 50 }
}
```

### 6.3 Competency Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/competencies` | All competency XP and levels |
| GET | `/characters/:id/competencies/:category` | One category |
| POST | `/characters/:id/competencies/:category/award` | Award competency XP |

### 6.4 Admin Endpoints

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/admin/game-config/level-thresholds` | List level thresholds |
| POST | `/admin/game-config/level-thresholds` | Set threshold for a level |
| GET | `/admin/game-config/competency-categories` | List categories |
| POST | `/admin/game-config/competency-categories` | Create category |
| GET | `/admin/game-config/competency-thresholds/:category` | List thresholds for category |
| POST | `/admin/game-config/competency-thresholds/:category` | Set threshold |
| GET | `/admin/game-config/character-config` | Get all character config |
| PUT | `/admin/game-config/character-config/:key` | Update a CharacterConfig value |
| GET | `/npc-templates` | List all NPC templates |
| POST | `/npc-templates` | Create NPC template |
| PUT | `/npc-templates/:id` | Update template (including xp_value) |
| GET | `/npc-templates/:id` | Get one template |
| DELETE | `/npc-templates/:id` | Delete template |

---

## 7. Database Schema Additions

### New Tables

**`character`** — add `xp` field:
- `xp` int default 0

**`character_competency`**:
- `id` int PK
- `character_id` int FK → Character
- `category_id` string FK → CompetencyCategory
- `xp` int default 0
- `level` int (cached, recomputed on award)
- UNIQUE(`character_id`, `category_id`)

**`level_threshold`** (game config):
- `level` int PK
- `xp_required` int

**`class_level_bonus`**:
- `class` string
- `level` int
- `hp_bonus`, `mana_bonus`, `stamina_bonus` int
- `stat_bonus_str/dex/con/int/wis` int

**`competency_category`**:
- `id` string PK
- `name` string
- `stat_key` string
- `xp_multiplier` float (default 0.20)

**`competency_level_threshold`**:
- `category_id` string FK
- `level` int
- `xp_required` int
- `damage_multiplier` float
- `defense_multiplier` float

**`xp_award_log`** (history, optional):
- `id` int PK
- `character_id` int FK
- `xp_awarded` int
- `source` string (NPC_kill, quest, admin, etc.)
- `NPC_id` int nullable
- `NPC_name` string nullable
- `created_at` timestamp

### Modified Tables

**`npc_template`** — add `xp_value` int default 0

---

## 8. Web Admin UI

Screens needed:
1. **NPC Template Editor** — CRUD for NPC templates including `xp_value` field
2. **Level Thresholds Editor** — Table editor for level → XP required mapping
3. **Class Level Bonuses Editor** — Table editor for per-class per-level stat/HP bonuses
4. **Competency Categories Editor** — CRUD for competency categories
5. **Competency Thresholds Editor** — Per-category threshold editor with multiplier fields
6. **Character XP Viewer** — View XP, level, competency XP for any character

---

## 9. TUI (Admin) Screens

Mirror all Web Admin UI screens in the Go/Bubble Tea admin TUI:
- `/npc-templates` — list, view, edit
- `/level-thresholds` — list, edit
- `/competencies` — list categories and thresholds

---

## 10. Implementation Order

1. Add `xp_value` to `npc_template` — schema migration only
2. Add `xp` field to `Character` — schema migration
3. Create `CharacterConfig` table with `max_level`, `xp_level_diff_scale_factor`, and `death_penalty_*` entries
4. Create `LevelThreshold` table and seed default values
5. **Refactor routes to service layer** — extract business logic from route handlers into service packages. Each service function calls DB and emits events to the event bus. Handlers become thin: parse request → call service → return response. (See Section 6.1 architecture diagram.)
6. Implement event bus (in-memory pub/sub; emit `NPC_killed`, `quest_completed`, `character_died`, `character_leveled_up`)
7. Implement `/characters/:id/xp/award` endpoint with NPC XP formula
8. Implement level-up logic with stat bonuses from `class_level_bonus`
9. Create `CompetencyCategory` and `CompetencyLevelThreshold` tables
10. Create `CharacterCompetency` join table
11. Implement `/characters/:id/competencies/:category/award` endpoint
12. Wire competency XP award into `NPC_killed` event subscriber
13. Implement XP subscriber for `quest_completed` event
14. Implement death penalty subscriber on `character_died` event
15. Integrate competency multipliers into combat damage calculation
16. Admin UI and TUI for all new content (NPC templates, level thresholds, competency categories/thresholds, CharacterConfig)
17. Remove YAML NPC template loading (covered by #282)

---

## 11. Quest XP — Event Hook System

Quests do not hardcode XP logic. Instead, quest completion emits a completion event. Any system can subscribe to these events. The XP award is one possible subscriber.

**Event: `quest_completed`**
```json
{
  "event": "quest_completed",
  "character_id": 42,
  "quest_id": "rescue_the_innkeeper",
  "quest_title": "Rescue the Innkeeper",
  "timestamp": "2026-04-28T..."
}
```

**XP Hook subscriber** (registered by default):
- On `quest_completed`, looks up `QuestConfig.xp_reward` and awards that amount of character XP
- Can be disabled per-quest or globally via `CharacterConfig.quest_xp_enabled`

This pattern extends to achievements, boss kills, exploration milestones, etc. — any event can have zero or more subscribers. The event bus is the integration point, not hardcoded calls.

**QuestConfig DB table additions:**
| Field | Type | Notes |
|---|---|---|
| `xp_reward` | int | XP to award on completion (0 = none) |
| `competency_category` | string | Optional: award competency XP in this category |
| `competency_xp_multiplier` | float | Fraction of NPC XP to award (e.g., 0.5) |

## 12. Death Penalty

**Design:** Configurable via `CharacterConfig` keys. Default: **disabled (0%)**.

| Config Key | Type | Default | Notes |
|---|---|---|---|
| `death_penalty_enabled` | bool | false | Toggle death penalty system |
| `death_penalty_xp_percent` | float | 0.0 | % of current-level XP lost on death (0.0 = off) |
| `death_penalty_currency_percent` | float | 0.0 | % of gold lost on death |
| `death_penalty_item_drop` | bool | false | Drop one random equipped item on death |
| `death_penalty_duration_seconds` | int | 0 | Resurrection sickness debuff duration |

**On death:**
1. If `death_penalty_enabled = false` → no penalty
2. If `death_penalty_xp_percent > 0`: `character.xp = character.xp × (1 − death_penalty_xp_percent)`
3. Apply other configured penalties (currency, item drop)
4. Emit `character_died` event (for achievement/event subscribers)

**Extensibility:** The penalty model is just a subscriber on the `character_died` event. New penalty types (debuffs, stat drain, temporary resurrection sickness) can be added as additional subscribers without modifying core death logic.

## 13. Open Questions

1. **Resurrection sickness** — The death penalty config includes `sickness_duration` seconds. Basic implementation: stat reduction for that duration. Should it also prevent combat participation during the debuff? (Stat reduction only for v1 — confirm?)
