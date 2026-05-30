# Character Registry System — Deep Dive

**Date:** 2026-04-28
**Status:** Research Phase
**File:** `docs/CHARACTER-SYSTEM-DEEPDIVE.md`

---

## 1. Core Character Entity

**Location:** `server/db/schema/character.go`

The `Character` entity is the central player/NPC model. It stores all persistent state.

### Fields

| Field | Type | Default | Notes |
|---|---|---|---|
| `name` | string | — | 1-23 chars, letters only. Unique constraint not enforced at DB level |
| `password` | string | — | bcrypt-hashed. Optional for NPCs |
| `isNPC` | bool | false | True = NPC/mob |
| `is_admin` | bool | false | Admin flag |
| `is_immortal` | bool | false | Cannot die — takes damage but never dies |
| `currentRoomId` | int | — | Active room |
| `startingRoomId` | int | — | Respawn point |
| `respawnRoomId` | int | 5 | "The Hole" default |
| `hitpoints` / `max_hitpoints` | int | 100 | HP pool |
| `stamina` / `max_stamina` | int | 50 | Combat/action resource |
| `mana` / `max_mana` | int | 25 | Spell resource |
| `race` | string | "human" | human, turtle, mutant (from DB Race table) |
| `class` | string | "adventurer" | Free-form string (warrior, survivor, etc.) |
| `specialty` | string | — | Class specialty (e.g., "fighter" for warrior) |
| `level` | int | 1 | Current level |
| `strength` / `dexterity` / `intelligence` / `wisdom` / `constitution` | int | 10 | Core stats |
| `gender` | string | — | From DB Gender table (he_him, she_her, they_them) |
| `description` | string | — | Free text |

### Combat Skill Proficiencies (flat integer fields on Character)

These represent weapon/proficiency levels, not linked to the Skill table:

- `skill_blades` — swords, daggers
- `skill_staves` — staves, wands
- `skill_knives` — knives, short blades
- `skill_martial` — martial arts, unarmed
- `skill_brawling` — brawling, fists
- `skill_tech` — technology, gadgets
- `skill_light_armor` / `skill_cloth_armor` / `skill_heavy_armor` — armor proficiency

### NPC-specific Fields

- `npc_skill_id` — identifier string (e.g., "druid_heal"), NOT a FK to Skill table
- `npc_skill_cooldown` — current cooldown ticks

### Character Edges (Relationships)

```
Character
├── user → User (1:1, optional)
├── room → Room (1:1, required) — current room
├── npcTemplate → NPCTemplate (1:1, optional)
├── skills → CharacterSkill (1:many) — learned combat skills (slot 1-5)
├── talents → CharacterTalent (1:many) — equipped talents (slot 0-3)
└── available_talents → AvailableTalent (1:many) — unlocked but not equipped
```

---

## 2. Combat Skill System (CharacterSkill)

**Location:** `server/db/schema/character_skill.go`

A join table between Character and Skill with one extra field:

| Field | Type | Notes |
|---|---|---|
| `slot` | int | Skill slot 1-5 (matching classless skill hotbar system) |

**Edge:** Character (1) ↔ Skill (1) via CharacterSkill

This is NOT the same as the combat proficiency fields (`skill_blades`, etc.). The proficiency fields are unlocked weapon categories; `CharacterSkill` entries are actual learned skill abilities.

---

## 3. Talent System

There are THREE talent-related entities:

### 3a. Talent (the template)

**Location:** `server/db/schema/talent.go`

A globally-defined ability template stored in the DB. Created/administered via REST.

| Field | Type | Notes |
|---|---|---|
| `name` | string | Unique |
| `description` | string | |
| `requirements` | string | JSON string (skill prereqs, level prereqs) |
| `effect_type` | string | heal, damage, dot, buff_armor, buff_dodge, buff_crit, debuff |
| `effect_value` | int | Magnitude (HP healed, damage dealt, etc.) |
| `effect_duration` | int | Ticks. 0 = instant |
| `cooldown` | int | Ticks between uses |
| `mana_cost` / `stamina_cost` | int | |

**REST endpoints (talents_routes.go):**
- `GET /talents` — list all
- `GET /talents/:id` — get one
- `POST /talents` — create
- `PUT /talents/:id` — update
- `DELETE /talents/:id` — delete
- `GET /talents/effect/:effectType` — filter by effect type

### 3b. CharacterTalent (equipped)

**Location:** `server/db/schema/character_talent.go`

Join table: Character ↔ Talent

| Field | Type | Notes |
|---|---|---|
| `slot` | int | Equipment slot 0-3 (quick-access bar) |

### 3c. AvailableTalent (unlocked)

**Location:** `server/db/schema/available_talent.go`

Tracks talents a character has unlocked but NOT yet equipped.

| Field | Type | Notes |
|---|---|---|
| `unlock_reason` | string | level_up, quest, skill_trainer, item (default: level_up) |
| `unlocked_at_level` | int | Character level when unlocked (default: 1) |

**Edge:** Character (1) ↔ AvailableTalent (1) ↔ Talent (1)

---

## 4. Skill Content Registry (YAML-based, not DB)

**Location:** `server/content/skill_registry.go`

Skills are defined in YAML files under `content/<world>/skills/`, loaded into an in-memory registry (not the DB). These are content definitions — the `Skill` DB table is a separate, runtime-created record.

### SkillDef (in-memory registry struct)

```go
type SkillDef struct {
    ID               string
    Name             string
    Description      string
    Type             string          // combat, magic, utility
    Tags             []string        // classless, damage, buff, cyberpunk, tech, etc.
    LevelRequirement int
    ClassRequirement string          // empty = classless
    Prerequisites    []SkillPrereq   // {SkillID, Level}
    Effects          []EffectDef     // damage, buff, debuff with stat scaling
    Cooldown         int             // seconds
    ManaCost         int
    StaminaCost      int
    HealthCost       int             // HP sacrificed
    Visual           VisualDef       // icon, color, animation, sound
    AIBehavior       AIBehaviorDef   // can_use, use_chance, health_threshold
}
```

### EffectDef

```go
type EffectDef struct {
    Type     string      // damage, buff, debuff
    Target   string      // self, target
    Value    interface{} // base number or complex value
    Scaling  *ScalingDef // {Stat, Ratio} — e.g., {strength, 1.0}
    Duration int         // ticks
}
```

### Scaling example (Haymaker skill)

```yaml
effects:
  - type: damage
    target: self
    value: 5
    scaling:
      stat: strength
      ratio: 1.0    # +1 damage per STR point
  - type: debuff
    target: self
    value: -50      # -50% (percentage reduction)
    scaling:
      stat: strength
      ratio: 0.5    # -0.5% per STR point (so DEX penalty scales with STR oddly)
    stat: dexterity
```

### Available skills (YAML)

| Path | Skills |
|---|---|
| `content/default/skills/classless/` | back_off, concentrate, haymaker, scream, slap |
| `content/cyberpunk/skills/` | hack |
| `content/default/skills/magic/` | (empty) |
| `content/default/skills/combat/` | (empty) |
| `content/default/skills/pc/` | (empty) |

**Registry methods:**
- `Get(id)` / `GetAll()` / `Count()`
- `GetByTag(tags...)` — OR match
- `GetByClass(class)` — also returns classless skills for "", "classless", "survivor"
- `Validate()` — checks required fields, prerequisites exist

### Content Routes (read-only registry access)

- `GET /content/skills` — all skills
- `GET /content/skills/:id` — one skill
- `GET /content/skills/tag/:tag` — by tag
- `GET /content/skills/class/:class` — by class
- `GET /content/validate` — validate all content

---

## 5. Skill DB Table (Runtime)

**Location:** `server/db/schema/skill.go`

Separate from the YAML registry. This is the runtime/DB representation.

| Field | Type | Notes |
|---|---|---|
| `name` | string | Unique. Skill name |
| `description` | string | |
| `skill_type` | string | combat, magic, utility |
| `cost` | int | Points to learn/unlearn |
| `cooldown` | int | Ticks |
| `requirements` | string | JSON prereqs |
| `effect_type` | string | Handler key: concentrate, haymaker, backoff, scream, slap |
| `effect_value` | int | Base damage/heal |
| `effect_duration` | int | Ticks. 0 = instant |
| `scaling_stat` | string | wisdom, strength, dexterity, constitution, intelligence |
| `scaling_percent_per_point` | float | e.g., 0.05 = +5% per stat point |
| `mana_cost` / `stamina_cost` / `hp_cost` | int | |

**Edge:** Skill (1) → CharacterSkill (many) → Character (many)

The YAML skill definitions and the DB Skill table appear to be two parallel systems. The YAML registry is for content/AI; the DB Skill table is for runtime character learning.

---

## 6. NPC Templates

**Location:** `server/db/schema/npc_template.go`

Templates for spawning NPCs/mobs.

| Field | Type | Notes |
|---|---|---|
| `id` | string | Unique identifier (e.g., "aragorn") |
| `name` | string | Display name |
| `description` | string | |
| `race` | string | |
| `disposition` | enum | hostile, friendly, neutral (default: neutral) |
| `level` | int | Default: 1 |
| `skills` | JSON map[string]int | Skill ID → level (NPC skill levels) |
| `trades_with` | []string | Faction tags for trading |
| `greeting` | string | |

**Edge:** NPCTemplate (1) → NPCSkill (many) → Skill (many)

**Content templates** (LotR-themed example):
- `content/default/npcs/templates/aragorn.yaml` — level 10, warrior class, sword+shield equipped, has combat skills
- `content/default/npcs/templates/frodo.yaml` — level 2, survivor class
- `content/default/npcs/templates/gimli.yaml` — level 8, warrior
- `content/default/npcs/templates/legolas.yaml` — level 9, hunter
- `content/default/npcs/templates/sam.yaml` — level 2, survivor

### NPC Spawning

NPCs are spawned from templates via `dbinit/init.go`. There's no separate spawn YAML file — the spawn logic is Go code that reads templates and creates Character records with `isNPC=true`. The Gizmo NPC is spawned in the fountain room via `InitGizmoNPC`. The junkyard golem spawns are tested in `junkyard_test.go`.

**NPC fields set at spawn:**
- `isNPC = true`
- `is_immortal = true` (for some NPCs)
- `npc_skill_id` — NPC ability identifier
- `npc_skill_cooldown`
- Template stats override character defaults

---

## 7. Class System & Constants

**Location:** `server/constants/character.go` + `server/constants/class_specialties.go`

### ClassConfig structure

```go
type ClassConfig struct {
    Name            string
    Specialty       string
    StatBonuses     StatBonuses  // +/- to STR, DEX, CON, INT, WIS
    StartingSkills  map[string]int  // skill_X → level (proficiency fields)
    PrimaryStat     string
    HealthPerLevel  float
    ManaPerLevel    float
}
```

### Warrior fighter example

```go
{
    Name:       "warrior",
    Specialty:  "fighter",
    StatBonuses: StatBonuses{Strength: 4, Constitution: 3, Dexterity: 1},
    StartingSkills: map[string]int{
        "blades":       3,
        "heavy_armor":  2,
        "brawling":     1,
    },
    PrimaryStat: "strength",
    HealthPerLevel: 12,
    ManaPerLevel: 2,
}
```

### Survivor (default class)

```go
{
    Name:       "survivor",
    Specialty:  "",
    StatBonuses: StatBonuses{},
    StartingSkills: map[string]int{
        "brawling": 1,
        "knives":  1,
    },
    PrimaryStat: "dexterity",
    HealthPerLevel: 8,
    ManaPerLevel: 4,
}
```

---

## 8. Character Routes (REST API)

**Location:** `server/routes/character_routes.go` (2160 lines)

### Authentication

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/characters/authenticate` | bcrypt verify → 64-char hex token (not stored, stateless) |

### Character CRUD

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/characters` | Create NPC or PC (no user association by default) |
| GET | `/characters` | List all |
| GET | `/characters/:id` | Get one |
| PUT | `/characters/:id` | Update name, room, admin, gender, description |
| DELETE | `/characters/:id` | Delete |

### Character-by-User

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/user-characters/:id` | List user's characters (no passwords) |
| POST | `/user-characters/:id` | Create character for user |

The `POST /user-characters/:id` endpoint:
1. Validates name (1-23 letters only)
2. Enforces max 3 characters per user
3. Hashes password with bcrypt
4. Resolves race from DB (human, turtle, mutant — must be `is_playable=true`)
5. Resolves gender from DB (he_him, she_her, they_them)
6. Applies class stat bonuses from `constants.GetClassConfig(class, specialty)`
7. Sets starting room from `room.is_starting_room=true`
8. Applies race stat modifiers via `dbinit.ApplyRaceToCharacter()`
9. Sets class proficiency fields (skill_blades, etc.) from class config

### Stat/Attribute Routes

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/stats` | Full stats (see below) |
| PUT | `/characters/:id/stats` | Update STR/DEX/INT/WIS/CON |
| GET | `/characters/:id/attributes` | List computed attributes |
| PUT | `/characters/:id/attributes` | Update attributes (gender, description, etc.) |
| GET | `/characters/:id/status` | HP, stamina, mana, level, class, race, room |
| PUT | `/characters/:id/status` | Modify HP/stamina/mana (admin) |
| GET | `/characters/:id/combat` | Combat info: defenses, damage, resistances |
| PUT | `/characters/:id/combat` | Update combat stats |

### Movement & Rooms

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/characters/:id/move` | Move to room |
| GET | `/characters/:id/room` | Current room info |
| GET | `/characters/:id/room/history` | Visit history |

### Skills & Proficiencies

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/skills` | All skills (CharacterSkill edges + proficiency fields) |
| PUT | `/characters/:id/skills` | Update skill levels (proficiency fields: blades, staves, etc.) |
| GET | `/characters/:id/skills/info` | Skill info for client |
| POST | `/characters/:id/skills/learn` | Learn a skill (add to CharacterSkill) |
| DELETE | `/characters/:id/skills/:skillName` | Forget a skill |
| GET | `/characters/:id/skills/available` | Skills available to learn |
| PUT | `/characters/:id/skills/progress` | Update proficiency XP |

### Talents

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/talents` | Equipped + available talents |
| POST | `/characters/:id/talents/equip` | Equip from available to slot |
| DELETE | `/characters/:id/talents/:slot` | Unequip from slot |

### Equipment

| Method | Endpoint | Notes |
|---|---|---|
| GET | `/characters/:id/equipment` | Full equipment loadout |
| PUT | `/characters/:id/equipment` | Equip/unequip items |
| GET | `/characters/:id/equipment/slots` | All equipment slots with items |

### Death & Respawn

| Method | Endpoint | Notes |
|---|---|---|
| POST | `/characters/:id/die` | Handle death |
| POST | `/characters/:id/respawn` | Respawn at respawnRoomId |
| PUT | `/characters/:id/respawn-point` | Set respawn room |

---

## 9. Equipment System

**Location:** `server/db/schema/equipment.go`

Item instances that can be in rooms or owned by characters.

| Field | Type | Notes |
|---|---|---|
| `name` / `description` | string | |
| `slot` | string | head, chest, weapon, legs, etc. |
| `level` | int | Required level |
| `weight` | int | |
| `isEquipped` | bool | |
| `isImmovable` | bool | Cannot be picked up |
| `isVisible` | bool | Shown in room list |
| `color` | string | Display color (e.g., gold for immovable) |
| `itemType` | string | weapon, armor, consumable, quest, misc, container, potion |
| `ownerId` | int? | Character ID who owns this item (nil = in room) |
| `effect_type` | string | heal, damage, dot, buff_armor, buff_dodge, buff_crit, debuff |
| `effect_value` / `effect_duration` | int | |
| `isContainer` | bool | Can hold items |
| `containerCapacity` | int | Max items |
| `isLocked` / `keyItemID` | | Container locking |
| `containedItems` | string | JSON array of item IDs |
| `revealCondition` | string | JSON: examine, perception_check, use_item, event |
| `expiresAt` | time? | Corpse rotting — nil = never |

**Equipment Routes (equipment_routes.go):**
- `GET /equipment` — all items
- `GET /equipment/:id` — one item
- `POST /equipment` — create
- `PUT /equipment/:id` — update
- `DELETE /equipment/:id` — delete
- `GET /rooms/:id/equipment` — items in a room

---

## 10. Race & Gender (Lookup Tables)

**Location:** `server/db/schema/race.go` + `server/db/schema/gender.go`

These are seed-data lookup tables, not content YAML.

### Race Fields

| Field | Type | Notes |
|---|---|---|
| `name` | string | human, turtle, mutant |
| `is_playable` | bool | Only playable races can be used for player characters |
| `stat_modifiers` | JSON | {strength, dexterity, constitution, intelligence, wisdom} |
| `description` | string | |

### Gender Fields

| Field | Type | Notes |
|---|---|---|
| `name` | string | he_him, she_her, they_them |
| `description` | string | |

---

## 11. Key Architectural Observations

### Two Parallel Skill Systems

There are TWO distinct skill systems that need to be understood together:

1. **YAML Content Registry** (`content/<world>/skills/`) — Designer's tool. Defines skill effects, scaling, AI behavior, visual feedback. Loaded at startup into `SkillRegistry` in-memory map. Accessed via `/content/skills/*` routes. This is the content definition.

2. **DB Skill Table** (`server/db/schema/skill.go`) — Runtime skill definitions. The `effect_type` field maps to handler names (concentrate, haymaker, backoff, scream, slap). These are what characters actually "learn" via the `/characters/:id/skills/learn` endpoint.

The YAML registry and DB Skill table appear to be two parallel representations. The YAML is richer (has AI behavior, complex scaling, visual effects); the DB is simpler (handler key + base values).

### NPC System

NPCs are `Character` records with `isNPC=true`. They're spawned from `NPCTemplate` records (LotR examples) via Go code in `dbinit/`, not YAML spawn files. Each NPC template has a `skills` JSON map (skill IDs → levels) and a disposition (hostile/friendly/neutral).

### Level & XP

The `Character.level` field exists (default: 1). There is NO XP field on the Character entity. The current system has levels but no XP tracking. The XP RFC will add:
- Mob XP values (content-defined `xp_value` on NPC templates or spawn entries)
- Level-differential multiplier
- Character XP accumulation
- Level-up thresholds
- Competencies/proficiencies XP (separate tracking for weapon categories)

### No Faction/Tags System Yet

The RFC-FACTION-SKILLS is approved but not yet implemented. The current system has no:
- `FactionCategories` table
- `CharacterTags` table
- `CharacterFactions` join table
- Faction-gated skills

### Talent System is Separate from Skills

Talents (DB table, CRUD via REST) and Skills (YAML + DB, learnable by characters) are two distinct systems:
- **Talents**: heal/damage/dot/buff effects with cooldown, mana/stamina cost, stored in DB
- **Skills**: more complex YAML-defined abilities with stat scaling, AI behavior, visual feedback

### What "Skill" Means in Three Contexts

1. **Character skill proficiency fields** (`skill_blades`, `skill_staves`, etc.) — flat integer levels on the Character table. These are weapon/proficiency categories, NOT linked to the Skill DB table or CharacterSkill join.

2. **CharacterSkill join table** — links Character to Skill (DB). Has a `slot` field (1-5) for hotbar.

3. **YAML SkillDef** — content definition with complex effects, scaling, AI behavior. Used by the content registry.

---

## 12. Entity Relationship Summary

```
User (1) ──→ (many) Character
                        │
                        ├─── (1:1) ──→ Room (current)
                        │
                        ├─── (1:many) ──→ CharacterSkill ── (1:1) ──→ Skill (DB)
                        │                    └── slot (1-5 hotbar)
                        │
                        ├─── (1:many) ──→ CharacterTalent ── (1:1) ──→ Talent
                        │                    └── slot (0-3 quick bar)
                        │
                        ├─── (1:many) ──→ AvailableTalent ── (1:1) ──→ Talent
                        │                    └── unlock_reason, unlocked_at_level
                        │
                        ├─── (1:1) ──→ NPCTemplate (if isNPC=true)
                        │                    └── skills: JSON map[skillID]level
                        │                    └── disposition: hostile/friendly/neutral
                        │
                        └─── (1:many, via ownerId) ──→ Equipment
                                                          └─── itemType, slot, effect_type

Skill (DB) ◄─── (1:many) ─── NPCSkill ── (many:1) ─── NPCTemplate
        │
        └─── Also loaded from YAML into SkillRegistry (in-memory)
                └─── content/<world>/skills/*.yaml

Talent (DB) ── CRUD via REST, no YAML content file
Race (DB) ── seed data: human, turtle, mutant
Gender (DB) ── seed data: he_him, she_her, they_them
```

---

## 13. REST Endpoint Summary

### Characters
- `POST /characters` — create
- `GET /characters` — list all
- `GET /characters/:id` — get one
- `PUT /characters/:id` — update
- `DELETE /characters/:id` — delete
- `POST /characters/authenticate` — login
- `GET /user-characters/:userId` — by user
- `POST /user-characters/:userId` — create for user

### Character Sub-resources
- `GET|PUT /characters/:id/stats` — core stats
- `GET|PUT /characters/:id/attributes` — gender, description
- `GET|PUT /characters/:id/status` — HP/stamina/mana
- `GET|PUT /characters/:id/combat` — defenses, damage
- `POST /characters/:id/move` — movement
- `GET /characters/:id/room` / `room/history`
- `GET|PUT /characters/:id/respawn-point`
- `POST /characters/:id/die` / `respawn`

### Skills
- `GET /characters/:id/skills` — learned skills + proficiencies
- `PUT /characters/:id/skills` — update proficiencies
- `GET /characters/:id/skills/info` — for client
- `POST /characters/:id/skills/learn` — learn skill
- `DELETE /characters/:id/skills/:skillName` — forget
- `GET /characters/:id/skills/available` — learnable
- `PUT /characters/:id/skills/progress` — proficiency XP

### Talents
- `GET /characters/:id/talents` — equipped + available
- `POST /characters/:id/talents/equip` — equip
- `DELETE /characters/:id/talents/:slot` — unequip

### Equipment
- `GET /characters/:id/equipment`
- `PUT /characters/:id/equipment`
- `GET /characters/:id/equipment/slots`

### Talents (master list)
- `GET /talents`
- `GET /talents/:id`
- `POST /talents`
- `PUT /talents/:id`
- `DELETE /talents/:id`
- `GET /talents/effect/:effectType`

### Content (read-only registry)
- `GET /content/skills` / `:id` / `tag/:tag` / `class/:class`
- `GET /content/items` / `:id`
- `GET /content/npcs` / `:id`
- `GET /content/rooms` / `:id` / `:id/exits`
- `GET /content/quests` / `:id` / `difficulty/:d` / `type/:t`
- `GET /content/stats`
- `GET /content/validate`

### Equipment (global)
- `GET|POST /equipment`
- `GET|PUT|DELETE /equipment/:id`
- `GET /rooms/:id/equipment`
