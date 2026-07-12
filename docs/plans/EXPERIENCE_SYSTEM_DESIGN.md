# Experience & Progression System — Design Doc

> **Status:** STASHED — design complete, implementation deferred.
> **Created:** 2026-07-12
> **Pick up:** After skill column refactor is complete (hardcoded `skill_*` columns → DB-driven `character_skills` join table).

---

## Overview

A comprehensive, DB-driven experience system that spans multiple MUDs. Each world configures its own XP curves, skills, stat scaling, and progression rules. The engine provides the *mechanisms*; each world defines the *content*.

**Core principle:** Everything is world-specific. No hardcoded class data, skill lists, or stat scaling in the engine code. Consistent with the v0.39.14 pattern (factions, abilities, classes all DB-driven).

---

## Dual-Track XP System

### Track 1: Level XP (Character Progression)
- Single scalar XP value per character (existing `xp` field)
- Level-up triggered by reaching XP threshold
- Level scales core stats: HP, mana, stamina, saving throws
- Level does NOT directly scale resistances (those come from race + equipment — see below)
- Max level configurable per world
- Two level curve modes, selectable per world:
  1. **Percentage increase:** each level requires N% more XP than the previous (default ~50%)
  2. **Hand-coded:** explicit XP threshold array per level (JSON on world config)

### Track 2: Skill XP (Per-Skill Advancement)
- Each character has per-skill XP and level via a `character_skills` join table
- Skills gain XP from **both**:
  - **Usage** (swing sword → sword skill XP) with diminishing returns to prevent macro-grinding
  - **Quest rewards** (specific skill or player-allocated pool)
- Hitting a skill level threshold unlocks abilities (stored as `required_skill_id` + `required_skill_level` on the ability)

---

## Stat Model

### Primary Stats (Core Attributes) — Engine-Level
Fixed columns on the character table. These are the fundamental attributes all MUDs share.

| Stat | Field | Default | Notes |
|------|-------|---------|-------|
| Strength | `strength` | 10 | Physical power |
| Dexterity | `dexterity` | 10 | Agility, reflexes |
| Constitution | `constitution` | 10 | Health, endurance |
| Intelligence | `intelligence` | 10 | Spell power, learning |
| Wisdom | `wisdom` | 10 | Spell effectiveness, perception |
| Charisma | `charisma` | 10 | Social, leadership — **NEEDS TO BE ADDED** to character schema |

### Derived Stats (Vitals) — Engine-Level
Fixed columns on character table. Already exist.

| Stat | Field | Default |
|------|-------|---------|
| HP | `hitpoints` | 100 |
| MaxHP | `max_hitpoints` | 100 |
| Stamina | `stamina` | 50 |
| MaxStamina | `max_stamina` | 50 |
| Mana | `mana` | 50 |
| MaxMana | `max_mana` | 50 |

### Combat Stats — Derived at Runtime
Calculated from core attributes + equipment + skills. NOT stored on character.

| Stat | Formula (initial) |
|------|-------------------|
| hit_roll | DEX + skill_bonus |
| damage_roll | STR + skill_bonus |
| armor_class | equipment AC + DEX modifier |
| dodge | DEX + dodge_skill |
| parry | DEX + relevant_weapon_skill |

### Resistances — DB-Driven (Race + Equipment)
Percentage-based (0-100%), not binary flags.

| Type | Source | Storage |
|------|--------|---------|
| fire, cold, poison, acid, lightning, energy, mental, disease, holy, unholy | Race base values | JSON on `races` table: `resistances: {fire: 20, poison: 15}` |
| slash, pierce, blunt | Race + armor | JSON on `races` table + equipment stats |
| Vulnerabilities | Race negative values | JSON: `vulnerabilities: {cold: -30}` |

**Scaling:** Resistances are fixed from race at creation, modified by equipment and abilities at runtime. They do NOT scale with level (per MUD convention — SMAUG, ROM, Iron Realms all use this model).

> **Open question for Sam:** Earlier you mentioned "levels scale against base resistances and resiliences." Research shows MUDs use fixed racial resistances + equipment mods. Do you want to override this and have resistances scale with level, or follow convention?

### Saving Throws — Derived or JSON
Three categories (simplified from SMAUG's 5):
- `save_physical` (vs poison, disease, blunt)
- `save_mental` (vs charm, sleep, mental)
- `save_spell` (vs direct spell damage)

Derived from WIS + level + equipment, or stored as JSON on character.

---

## NPC XP Values

Stored on NPC template (world-specific):

| Field | Type | Notes |
|-------|------|-------|
| `base_xp` | int | Fixed base XP for killing this NPC |
| `xp_multiplier` | float | Multiplied against base_xp for difficulty scaling |

**XP formula:** `XP = base_xp × clamp((mob_level - player_level + 10) / 10, 0, 1.3)`

**Anti-grinding:**
- Track kill counts per NPC template per player
- After N kills of same template, XP drops to near-0
- Configurable threshold per world

**Group XP:**
- Divide by group member count, weighted by level share (ROM model)
- Members not in combat get halved share

**XP sources:** Combat kills, quests, exploration (room discovery), crafting — all configurable per world.

---

## Skills (World-Specific)

### New Tables

**`skills` table:**
| Column | Type | Notes |
|--------|------|-------|
| id | int (PK) | Auto-increment |
| world_id | string | World scope |
| name | string | e.g., "blades", "heavy_armor" |
| display_name | string | "Blades", "Heavy Armor" |
| description | text | |
| category | string | "weapon", "armor", "craft", "magic" |
| parent_skill_id | int (FK→skills) | For skill tree structure (nullable) |
| max_level | int | Cap for this skill (default 100) |
| xp_curve_mode | enum | "percentage" or "hand_coded" |
| xp_curve_data | JSON | Percentage factor or explicit thresholds |

**`character_skills` join table:**
| Column | Type | Notes |
|--------|------|-------|
| id | int (PK) | |
| character_id | int (FK→characters) | |
| skill_id | int (FK→skills) | |
| level | int | Current skill level |
| xp | int | Current skill XP |

**`abilities` table changes:**
| New Column | Type | Notes |
|-----------|------|-------|
| required_skill_id | int (FK→skills, nullable) | Which skill must be leveled |
| required_skill_level | int | Minimum skill level to unlock |

### Migration from Hardcoded Columns
Current hardcoded skill columns on `characters`:
- `skill_blades`, `skill_staves`, `skill_knives`, `skill_martial`, `skill_brawling`, `skill_tech`
- `skill_light_armor`, `skill_cloth_armor`, `skill_heavy_armor`

**Migration plan:**
1. Create `skills` table, populate with the 9 existing skills per world
2. Create `character_skills` join table
3. Read existing character skill values, populate join table
4. Update all 30+ code references (service/ability_equip.go, routes/character_stats.go, routes/ws_routes.go, repository/interface.go, repository/character_create.go)
5. Remove 9 skill columns from character schema
6. Regenerate ent code

### Skill XP Sources
- **Usage-based:** performing an action that uses the skill grants XP
  - Diminishing returns: XP gain decreases per use within a time window
  - Cap per session to prevent macro-grinding
- **Quest rewards:** quest can grant skill XP to specific skills or to a pool the player allocates
- **Crafting:** successful craft grants XP to related skill

### Ability Unlocks
- Abilities have `required_skill_id` + `required_skill_level` (nullable, defaults to no requirement)
- When a character's skill level meets the threshold, the ability becomes available
- Can optionally require visiting a trainer NPC to "learn" the ability

---

## Level Scaling

### Stat Growth on Level-Up
Three layers combine:

1. **Base growth per level** (engine default or world config):
   - HP: +10/level, Mana: +5/level, Stamina: +5/level (configurable)

2. **Class/faction modifier** (DB-driven, from `faction.stat_bonuses`):
   - Trash Mage: HP +3/level, Mana +12/level
   - Foot Clank: HP +15/level, Mana +2/level

3. **Racial multiplier** (DB-driven, from `races` table):
   - Add `stat_growth_multipliers` JSON field to races: `{"hp": 1.1, "mana": 0.9}`
   - Final HP growth = base × class_modifier × racial_multiplier

### Level Curve Modes

**Mode 1: Percentage Increase (default)**
```
xp_for_level(n) = xp_for_level(n-1) × (1 + percentage/100)
```
- Default percentage: 50%
- Example: L2=1000, L3=1500, L4=2250, L5=3375

**Mode 2: Hand-Coded**
```json
[0, 1000, 2500, 5000, 10000, 20000, 40000, ...]
```
- Explicit array on world config

**Config storage:** JSON field on `worlds` table:
```json
{
  "level_curve": {
    "mode": "percentage",
    "base_xp": 1000,
    "percentage": 50,
    "max_level": 50
  }
}
```

---

## Re-classing & Re-racing

### Re-classing
- Engine-level feature, configurable per world
- Character retains: level, level XP
- Character loses: access to old class abilities (become ineligible)
- Character gains: new class abilities at base skill levels
- Skill levels may be partially retained (world config: `reclass_skill_retention: 0.5` = keep 50% of skill levels)
- Cost: currency, XP penalty, quest, or NPC service
- Config: `allow_reclass: true, reclass_cost: 1000, reclass_min_level: 10, reclass_cooldown: 3600`

### Re-racing
- Rare, costly operation
- Recalculate base stats and resistances from new race
- Keep level, class, and skills
- Config: `allow_rerace: false, rerace_cost: 5000`
- May require quest, ritual, or world-specific mechanism

### History Tracking
- `character_class_history` table: character_id, faction_id, joined_at, left_at
- `character_race_history` table: character_id, race_id, changed_at

---

## Event System

Emit events on an internal event bus:

| Event | Payload | Subscribers |
|-------|---------|-------------|
| `xp.gained` | {type, amount, source, character_id} | Achievements, quest triggers |
| `xp.lost` | {amount, reason, character_id} | Death penalty logging |
| `level.up` | {character_id, new_level, old_level} | Stat recalculation, announcements, achievements |
| `skill.xp.gained` | {character_id, skill_id, amount, source} | Achievements |
| `skill.leveled_up` | {character_id, skill_id, new_level} | Ability unlock checks, announcements |
| `reclass` | {character_id, old_faction, new_faction} | Skill reset, history logging |
| `rerace` | {character_id, old_race, new_race} | Stat recalculation, history logging |

Achievement system subscribes to these events. Quest system can hook `xp.gained` to track objectives. World-specific triggers can listen without modifying engine code.

---

## World Configuration

New JSON config field on `worlds` table (or a `world_config` table):

```json
{
  "level_curve": {
    "mode": "percentage",
    "base_xp": 1000,
    "percentage": 50,
    "max_level": 50
  },
  "stat_growth": {
    "hp_per_level": 10,
    "mana_per_level": 5,
    "stamina_per_level": 5
  },
  "skill_xp": {
    "usage_diminishing_returns": true,
    "usage_cap_per_hour": 100,
    "anti_grind_kill_threshold": 20
  },
  "reclass": {
    "allowed": true,
    "cost": 1000,
    "min_level": 10,
    "cooldown_seconds": 3600,
    "skill_retention": 0.5
  },
  "rerace": {
    "allowed": false,
    "cost": 5000
  }
}
```

---

## Open Questions (Need Sam's Input)

1. **Resistances scaling with level?** Research says fixed-from-race + equipment. Sam mentioned scaling with level earlier. Which approach?

2. **Charisma — add to character schema?** It's in constants but not persisted. SMAUG has it, ROM doesn't. Need it for social/leadership mechanics?

3. **Luck stat?** SMAUG has it as a 7th core attribute. Optional per world, or skip?

4. **Combat stats (hit_roll, damage_roll, AC, dodge, parry) — stored or derived?** Recommendation: derived at runtime from attrs + equipment + skills. Don't store.

5. **Skill column refactor — do first before XP implementation?** The 9 hardcoded `skill_*` columns must become a `character_skills` join table. This is a prerequisite to the XP system.

---

## Implementation Phases

### Phase 0: Skill Column Refactor (PREREQUISITE)
1. Create `skills` table (world-scoped)
2. Create `character_skills` join table
3. Migrate existing skill data
4. Update all code references (30+ files)
5. Remove hardcoded skill columns
6. Regenerate ent code

### Phase 1: Level XP System
1. Add `world_config` JSON to worlds table (level curve, stat growth)
2. Implement level-up logic (threshold check → stat scaling)
3. Add `stat_growth_multipliers` to races table
4. Implement XP gain from kills (NPC base_xp + level diff formula)
5. Implement anti-grind kill tracking
6. Implement group XP distribution

### Phase 2: Skill XP System
1. Add `required_skill_id` + `required_skill_level` to abilities table
2. Implement skill XP gain from usage (with diminishing returns)
3. Implement skill XP from quest rewards
4. Implement skill level-up → ability unlock check
5. Implement skill XP curves (per-skill, configurable)

### Phase 3: Resistances
1. Add `resistances` + `vulnerabilities` JSON to races table
2. Add resistance modifiers to equipment templates
3. Implement resistance calculation at runtime (race base + equipment mods)
4. Update combat damage to apply resistances

### Phase 4: Re-classing & Re-racing
1. Add reclass/rerace config to world_config
2. Implement reclass command (skill reset, ability eligibility update)
3. Implement rerace command (stat recalculation)
4. Add character_class_history and character_race_history tables
5. Emit events

### Phase 5: Event System
1. Implement internal event bus
2. Wire XP gain, level-up, skill level-up, reclass, rerace events
3. Achievement system subscription
4. Quest trigger subscription

### Phase 6: Charisma + Polish
1. Add `charisma` column to character schema (if decided)
2. Update character creation to include charisma
3. Update stat display routes
4. World export/import for all new tables

---

## Research Sources
- SMAUG source code (smaugfuss/fight.c, mud.h, update.c)
- ROM 2.4 source code (rom24-quickmud/fight.c, merc.h, skills.c)
- Achaea Wiki (wiki.achaea.com)
- Discworld MUD Wikipedia
- Iron Realms Entertainment Wikipedia
- Alter Aeon help site
- CircleMUD website
- LPMud Wikipedia