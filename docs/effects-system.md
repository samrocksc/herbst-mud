# Effects System Documentation

## Overview

The herbst-mud effect system defines **what happens** when an ability is used, an item is equipped, or a passive triggers. Effects are the building blocks that make abilities and items do things.

## Four-Domain Model

The system separates concerns across four domains:

| Domain | What it is | Examples | Entity |
|--------|-----------|----------|--------|
| **Abilities** | Actions characters can perform | Concentrate, Haymaker, Fireball | `Ability` |
| **Skills** | Leveled proficiencies | Blades, Staves, Light Armor | Character columns (future: `Skill` + `CharacterSkillLevel`) |
| **Stats** | Numeric attributes | Strength, Dexterity, Wisdom | `Character` fields |
| **Effects** | What actually happens | Damage, Heal, Buff, Stun | `AbilityEffect` |

### Abilities vs Skills vs Stats

- **Abilities** are things you **do** — active combat moves, passive triggers, spells. Each ability has one or more effects that define its outcome.
- **Skills** are things you **train** — weapon proficiencies and armor mastery that improve with use. They gate what equipment you can use and provide passive bonuses. Currently stored as flat Character columns (`skill_blades`, `skill_staves`, etc.). Phase 6 will migrate these to a proper `Skill` entity.
- **Stats** are things you **are** — numeric attributes like Strength (10+), Dexterity (10+), Constitution (10+). They scale ability effects.

## Ability Entity

Each `Ability` row represents a combat action or passive trigger:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique name (e.g., "Concentrate") |
| `description` | string | Flavor text shown to players |
| `ability_type` | string | Category: combat, magic, utility, healing, support, defensive |
| `ability_class` | string | active, passive, or toggle (passive = former talents) |
| `mana_cost` | int | MP drained on activation |
| `stamina_cost` | int | SP drained on activation |
| `hp_cost` | int | HP sacrificed to use |
| `cooldown_seconds` | int | Seconds before reuse |
| `proc_chance` | float | For passives: 0.0–1.0 trigger probability |
| `proc_event` | string | For passives: on_hit, on_hit_received, on_crit, on_kill |

Flat fields `effect_type`, `effect_value`, `effect_duration` remain for backward compatibility but are deprecated in favor of the `effects` edge.

## Effect Entity Schema

Each `AbilityEffect` row belongs to one `Ability` and defines a single outcome:

| Field | Type | Description |
|-------|------|-------------|
| `effect_type` | string | The category of effect (see below) |
| `damage_subtype` | string | For damage effects: slashing, piercing, bludgeoning, fire, cold, lightning, poison, psychic |
| `target` | string | Who receives the effect: self, enemy, ally, area, random_enemy |
| `value` | int | Base magnitude (damage amount, heal amount, buff duration, etc.) |
| `duration` | int | Duration in combat ticks (0 = instant) |
| `scaling_stat` | string | Which stat modifies this effect: strength, dexterity, constitution, intelligence, wisdom |
| `scaling_ratio` | float | Multiplier per point of scaling stat (0.5 = +50% per point) |
| `sort_order` | int | Order within the ability (0 = first effect) |

## Effect Types

### Primary Types

| Type | Description | Example |
|------|-------------|---------|
| `damage` | Direct HP reduction to the target | Fireball deals 30 fire damage |
| `heal` | HP restoration to the target | Heal restores 25 HP |
| `buff` | Positive status effect on the target | Concentrate boosts accuracy for 4 ticks |
| `debuff` | Negative status effect on the target | Scream reduces WIS/INT for 2 ticks |
| `dot` | Damage over time (ticks) | Poison deals 5 damage per tick for 3 ticks |
| `hot` | Heal over time (ticks) | Regeneration heals 3 HP per tick for 5 ticks |
| `stun` | Target skips their next turn | Slap stuns for 1 tick |
| `accuracy_boost` | Increases hit chance | Concentrate boosts accuracy |
| `dodge_all` | Avoids all attacks for duration | Back-off dodges everything for 1 tick |

### Damage Subtypes

Damage effects can specify a subtype for resistance/weakness calculations:

| Subtype | Description |
|---------|-------------|
| `slashing` | Swords, axes — countered by heavy armor |
| `piercing` | Daggers, arrows — countered by cloth armor |
| `bludgeoning` | Clubs, fists — countered by light armor |
| `fire` | Fireballs, dragon breath |
| `cold` | Ice spells, frost weapons |
| `lightning` | Shock spells, storm abilities |
| `poison` | Toxins, venoms |
| `psychic` | Mind-based attacks |

## How Effects Resolve

### 1. Ability Activation

When a character uses an ability:

1. **Resource check** — Verify the character has enough mana, stamina, and HP
2. **Cooldown check** — Verify the ability isn't on cooldown
3. **Resource deduction** — Subtract mana_cost, stamina_cost, hp_cost
4. **Effect resolution** — For each effect (sorted by sort_order):
   - Determine the target based on `target` field
   - Calculate the magnitude: `final_value = value + (scaling_stat_value × scaling_ratio × value)`
   - Apply the effect for `duration` ticks (0 = instant)
5. **Cooldown start** — Set cooldown for `cooldown_seconds`

### 2. Scaling Formula

```
final_value = base_value + (stat_modifier × scaling_ratio × base_value)
```

Where `stat_modifier = (stat - 10) / 2` (standard D&D modifier).

Example: Haymaker has `scaling_stat=strength`, `scaling_ratio=0.5`, `base_value=15`.
A character with Strength 18 (+4 modifier):
```
final_value = 15 + (4 × 0.5 × 15) = 15 + 30 = 45 damage
```

### 3. Passive Abilities

Passive abilities (`ability_class=passive`) use `proc_chance` and `proc_event` instead of manual activation:

- `proc_event` defines when to check: `on_hit`, `on_hit_received`, `on_crit`, `on_kill`
- `proc_chance` (0.0–1.0) defines the probability of triggering
- When triggered, the ability's effects are applied to the appropriate target

## Classless Abilities (Starting Abilities)

The five starting abilities and their effects, seeded in `server/dbinit/skills.go`:

### Concentrate
- **ability_type**: combat | **ability_class**: active | **cooldown_seconds**: 8 | **mana_cost**: 10
- **Effects**: buff (accuracy_boost, target=self, duration=4, scaling_stat=wisdom, scaling_ratio=0.5)

### Haymaker
- **ability_type**: combat | **ability_class**: active | **cooldown_seconds**: 6 | **stamina_cost**: 15
- **Effects**: damage (target=enemy, value=15, scaling_stat=strength, scaling_ratio=0.5) + debuff (target=self, duration=1, value=2)

### Back-off
- **ability_type**: defensive | **ability_class**: active | **cooldown_seconds**: 10 | **stamina_cost**: 25
- **Effects**: buff (dodge_all, target=self, duration=1, value=999)

### Scream
- **ability_type**: support | **ability_class**: active | **cooldown_seconds**: 12 | **mana_cost**: 5, **stamina_cost**: 10
- **Effects**: buff (target=self, scaling_stat=constitution, duration=2) + debuff (target=enemy, scaling_stat=constitution, duration=2)

### Slap
- **ability_type**: combat | **ability_class**: active | **cooldown_seconds**: 8 | **stamina_cost**: 12
- **Effects**: stun (target=enemy, duration=1, scaling_stat=dexterity, scaling_ratio=0.3)

## Client-Side Effect Resolution

The herbst SSH client resolves effects generically via `resolveAbilityEffects()` in `herbst/classless_skills.go`:

1. Fetches ability data from `/api/characters/:id/classless-skills` (includes effects array)
2. For each effect, calculates scaling using `calcScalingModifier()`
3. Dispatches by type: `damage` → `applyDamageEffect()`, `heal` → `applyHealEffect()`, `stun` → `applyStunEffect()`, others → tracked as `ActiveEffect`

No handler-key dispatch (the old `switch skill.EffectType` pattern is removed). All ability behavior is data-driven through the effects array.

## API Endpoints

### Abilities

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/abilities` | List all abilities (supports `?type=` and `?ability_class=` filters) |
| GET | `/api/abilities/:id` | Get a single ability |
| POST | `/api/abilities` | Create a new ability |
| PUT | `/api/abilities/:id` | Update an ability |
| DELETE | `/api/abilities/:id` | Delete an ability |

### Effects

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/abilities/:id/effects` | List effects for an ability |
| POST | `/api/abilities/:id/effects` | Add an effect to an ability |
| PUT | `/api/effects/:id` | Update an effect |
| DELETE | `/api/effects/:id` | Delete an effect |

### Character Abilities (Combat Slots)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/characters/:id/classless-skills` | Get equipped abilities with effects |
| POST | `/api/characters/:id/classless-skills` | Equip ability to slot (1-5) |
| PUT | `/api/characters/:id/classless-skills/swap` | Swap two ability slots |

## Admin Panel

The abilities management page (`/abilities`) supports:
- **Add/Edit form** with all ability fields
- **Effects sub-form** (appears when editing an existing ability) — add/remove AbilityEffect rows
- **Filters** by ability type and ability class (active/passive/toggle)
- **Class column** in the table showing active/passive status

## Migration Notes

### From Skills to Abilities (May 2026)

The `Skill` entity was renamed to `Ability` to clarify the domain separation:
- **Ability** = something you do (Concentrate, Haymaker)
- **Skill** = something you train (Blades, Staves)

Renamed tables: `skills` → `abilities`, `character_skills` → `character_abilities`, `npc_skills` → `npc_abilities`
Renamed columns: `skill_type` → `ability_type`, `skill_class` → `ability_class`

### From Handler Keys to Generic Effects

Previously, each ability had a single `effect_type` field that acted as a handler key (e.g., "concentrate" routed to the `applyConcentrate()` function). Now, effects are generic types (damage, heal, buff, etc.) stored in the `AbilityEffect` entity, with each ability able to have multiple effects. The client resolves effects generically — no more per-ability handler functions.

### From Flat Fields to Effect Entity

The flat fields `effect_type`, `effect_value`, `effect_duration` remain on the `Ability` entity for backward compatibility but are deprecated in favor of the `effects` edge to `AbilityEffect`. The admin panel uses the effects sub-form for creating/editing ability effects.

### Talents → Passive Abilities

The `Talent` entity has been merged into `Ability` using `ability_class='passive'` as the differentiator. Former talents (heal, damage, dot, buff_armor, buff_dodge, buff_crit, debuff) are now passive abilities with effect types matching the generic set.

### Phase 6 (Deferred): Skill Entity for Proficiencies

The flat Character columns (`skill_blades`, `skill_staves`, etc.) will eventually be migrated to a proper `Skill` entity with a `CharacterSkillLevel` join table. This is deferred because the flat columns work correctly.