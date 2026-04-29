# RFC: Faction-Based Skills System

**Status:** Approved for Implementation  \
**Author:** Sam (via Mikey 🐢🟠 + Leonardo 🐢🔵 interview)  \
**Date:** 2026-04-15  \
**Interview:** 2026-04-28 — All open questions resolved  \
**Related:** object-model.md, WORLD_BIBLE.md, skill-talent-db-schema.feature

---

## Summary

Add a **faction affiliation** dimension to the skills system. Skills are currently granted by:
1. **base_skills** — available to all characters (classless)
2. **class_skills** — granted by your character class

This RFC proposes adding:
3. **faction_skills** — granted by your faction affiliation

---

## Motivation

In Surf-Punk 2052, faction identity is core to the world:
- **Surf Wardens** protect the coast, master wave magic
- **Dune Traders** traverse the desert, know sand-step and caravan survival
- **Gondoliers** control trade routes, have canal-specific abilities
- **Tinkerers** salvage dead-zone tech, have gadget proficiencies
- **Vine Climb** live in the Overgrown Metropolis, have vine-climbing abilities

Currently, a Surf Warden Warrior and a Dune Trader Warrior have identical skill sets beyond class. Faction skills would differentiate them and reinforce faction identity.

---

## Proposed Entities

### Faction Categories

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| name | string | e.g., `class`, `alignment` |
| display_name | string | e.g., "Class", "Alignment" |
| description | string | Flavor text |
| max_memberships | int | How many factions in this category a character can hold (e.g., 2 for class, 1 for alignment) |
| auto_join | bool | If true, earning required tag auto-joins faction without confirmation (default false) |

### Factions Table

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| category_id | UUID | Foreign key to FactionCategory |
| name | string | e.g., `surf_warden`, `dune_trader`, `foot_clan`, `ninja` |
| display_name | string | e.g., "Surf Warden", "Foot Clan" |
| description | string | Flavor text |

### Character Tags (existing or new)

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| character_id | UUID | Foreign key to Character |
| tag | string | Tag identifier, e.g., `first_class`, `chef_learned`, `wizard_complete` |
| source | string | How the tag was earned: `system`, `quest`, `achievement`, `admin` |
| earned_at | timestamp | When tag was granted |

### Faction Required Tags

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| faction_id | UUID | Foreign key to Faction |
| required_tag | string | Tag the character must have to pledge to this faction |

### Character Factions (existing or enhanced)

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| character_id | UUID | Foreign key to Character |
| faction_id | UUID | Foreign key to Faction |
| reputation | int | Faction reputation (0-100), default 0 |
| joined_at | timestamp | When character joined |
| status | string | `active`, `expelled`, `voluntarily_left` |

### Skill Enhancement

Skills are extended with optional gating:

|| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key (existing) |
| slug | string | Globally unique skill identifier, e.g., `foot_clan_power_strike` |
| name | string | Display name, e.g., "Power Strike" (must be unique globally) |
| faction_id | UUID? | Faction that grants this skill, or null if universal |
| required_tag | string? | Additional tag required to unlock this skill beyond faction membership |
| type | string | `passive` or `active` |
| proc_chance | float? | For passives: % chance to proc (e.g., 0.15 = 15%) |
| proc_event | string? | What triggers the proc: `on_hit`, `on_hit_received`, `on_crit`, `on_kill` |
| mana_cost | int? | For actives: mana cost |
| cooldown_seconds | int? | For actives: cooldown duration |

---

## Example Skills by Faction

### Surf Warden Skills (faction_id = surf_warden)

| Skill ID | Name | Type | Description |
|----------|------|------|-------------|
| `tidal_weave` | Tidal Weave | passive | +WIS scaling near water |
| `wave_reader` | Wave Reader | passive | Detect currents, +DEX in water |
| `surf_rescue` | Surf Rescue | active | Save allies from drowning (talent) |

### Dune Trader Skills (faction_id = dune_trader)

| Skill ID | Name | Type | Description |
|----------|------|------|-------------|
| `sandstep` | Sand Step | passive | +DEX on sand/dunes |
| `beast_caller` | Beast Caller | passive | Can call mutant iguana rides |
| `oasis_sense` | Oasis Sense | passive | Detect water nearby |

### Tinkerer Skills (faction_id = tinkerer)

| Skill ID | Name | Type | Description |
|----------|------|------|-------------|
| `salvage` | Salvage | passive | Better loot from dead zones |
| `jury_rig` | Jury Rig | active | Temporarily repair broken items (talent) |
| `bio_hack` | Bio Hack | active | Use bio-tech implants (talent) |

---

## Resolved Design Decisions

All open questions answered by Sam on 2026-04-28.

### 1. Faction Categories with Slot Limits

Factions belong to a **category** (e.g., `class`, `alignment`). Each category has a configurable `max_memberships` limit. Characters can belong to up to N factions within that category simultaneously.

**Example:**
- `class` category → max 2 memberships → dual-classing allowed (e.g., `ninja` + `pizza_chef`)
- `alignment` category → max 1 membership → must choose side (e.g., `long_island_brutes` OR `foot_clan`)

Categories are defined in config/content, not hardcoded.

### 2. Tags as the Universal Gating System

**Tags** are achievement/milestone flags on a character. They gate both faction membership AND individual skill access.

- Tags are earned from: quest completion, achievements, admin commands, story events, character creation
- Every character auto-receives tag `first_class` on creation (system-granted, not quest-earned)
- Factions declare required tags to attempt joining
- Faction skills may additionally require specific tags to unlock

**Example flow:**
1. New character: 0 class slots unlocked
2. System grants `first_class` tag → unlocks class category slot 1
3. Character pledges to `ninja` faction (requires `first_class`) → joins
4. Complete `chef_questline` → earns tag `chef_learned` → unlocks class slot 2
5. Character pledges to `pizza_chef` faction (requires `chef_learned`) → dual-classed

### 3. Auto-Join is Explicit, Not Automatic

**Gating is automatic; joining is intentional.** Earning the required tag unlocks the ability to pledge, but the player must take a deliberate confirmation step (visit NPC or use interface) to formally join.

Auto-join can exist as a configurable toggle per category/faction (for seamless onboarding categories), but the default is formal pledge with confirmation.

### 4. Both Active Talents and Passive Proc Skills

Factions grant both:
## Example Factions

### Class Factions (category: `class`, max_memberships: 2)

| Faction | Required Tag | Passive Skills | Active Talents |
|---------|-------------|---------------|----------------|
| `ninja` | `first_class` | `shadow_step` (+15% dodge), `silent_killer` (+25% crit damage when untracked) | `foot_flip_kick` (30 mana, 10s CD) |
| `pizza_chef` | `chef_learned` | `dough_hands` (+DEX), `spice_sense` (detect poison) | `pizza_toss` (stun, 20 mana, 15s CD) |
| `wizard` | `wizard_complete` | `mana_flow` (+20% mana regen), `arcane_aura` (passive magic resist) | `arcane_bolt` (damage spell, 25 mana, 5s CD) |

### Alignment Factions (category: `alignment`, max_memberships: 1)

| Faction | Required Tag | Passive Skills | Active Talents |
|---------|-------------|---------------|----------------|
| `foot_clan` | `first_class` | `foot_loyalty` (can't be attacked by other Foot Clan), `shadow_walk` (+stealth) | `foot_fist` (unarmed strike, 15 mana) |
| `long_island_brutes` | `first_class` | `brute_strength` (+carry capacity), `thick_skin` (-10% physical damage) | `brute_slam` (AoE damage, 40 mana, 20s CD) |

---

## Implementation Notes for Donnie

### Schema Changes Needed

1. Create `faction_categories` table (name, display_name, max_memberships, auto_join)
2. Rename/update `factions` table to add `category_id` foreign key
3. Create `character_tags` table (character_id, tag, source, earned_at)
4. Create `faction_required_tags` table (faction_id, required_tag)
5. Update `character_factions` table: add `status` field
6. Update `skills` table: add `slug`, `required_tag`, `type`, `proc_chance`, `proc_event`, `mana_cost`, `cooldown_seconds`
7. Add unique constraint on `skills.slug`
8. Add unique constraint on `(character_id, tag)` in character_tags
9. Seed: auto-grant `first_class` tag on character creation
10. Seed: default faction categories (`class` max=2, `alignment` max=1, `auto_join=false`)

### API Changes

**Character Factions:**
- `GET /characters/:id/factions` — list affiliations and reputation with status
- `POST /characters/:id/factions` — pledge to join a faction (body: `{faction_id}`)
- `DELETE /characters/:id/factions/:faction_id` — leave a faction
- `PATCH /characters/:id/factions/:faction_id` — modify reputation (body: `{reputation}`, `{status}`)

**Factions:**
- `GET /factions` — list all factions with eligibility status for calling character
- `GET /factions/:id` — single faction with required tags
- `GET /characters/:id/faction-eligibility` — check which factions character can join given tags + current load

**Tags:**
- `GET /characters/:id/tags` — list all tags on character
- `POST /characters/:id/tags` — grant a tag (admin only, body: `{tag, source}`)
- `DELETE /characters/:id/tags/:tag` — remove a tag

**Skills (enhanced):**
- `GET /characters/:id/skills` — include faction_skills with eligibility check (tag + membership + reputation gates)
- `GET /skills` — list all skills with faction and tag requirements visible

### Eligibility Check Logic

When a character attempts to pledge to a faction:
1. Does character have `required_tag` for this faction? (check `faction_required_tags`)
2. Does the category this faction belongs to have an open slot? (`character_factions` count for category < `category.max_memberships`)
3. Is character already a member of this faction? (if yes, reject)

When listing available skills for a character:
1. Character must be `active` member of the skill's faction
2. If skill has `required_tag`, character must have that tag
3. Return skill with `eligible: true/false` and `reason` if ineligible

### Content Externalization

All factions, categories, tags, and faction-required-tags should be content-externalized (see `docs/CONTENT_EXTERNALIZATION_COMPLETE.md`). Seed from YAML/JSON content files.

### Test Files

- Update `skill-talent-db-schema.feature` with faction scenarios
- Add `faction-skills.feature` for faction-specific tests (join, leave, eligibility, proc passthrough)
- Add `character-tags.feature` for tag grant/revoke scenarios

### Frontend Changes (admin/ web UI + admin-tui TUI)

Both admin interfaces must support:
- Viewing all factions by category
- Viewing a character's tags, faction memberships, and reputation
- Granting/revoking tags
- Joining and leaving factions (with confirmation step)
- Adjusting reputation
- Setting faction required tags when creating/editing factions

---

## Open Questions (Resolved)

~~1. **Can a character have multiple faction affiliations?**~~ → YES, via category-based slot limits
~~2. **Reputation gate — Should faction skills require minimum reputation?**~~ → YES, reputation field exists; skill-level rep gates can be added per-skill if needed
~~3. **Joining a faction — How do characters join?**~~ → Formal pledge with confirmation; auto-join toggleable per category
~~4. **Skill vs Talent — Should faction abilities be passive, active, or both?**~~ → BOTH — passives with proc %, actives with mana/cooldown
~~5. **Leaving a faction — what happens to skills?**~~ → Skills revoked, reputation preserved, slot locked, REST-controlled

---

## Alignment with Surf-Punk Ideals

The faction system reinforces the post-collapse world:
- **Community** — People identify with their faction, not just their class
- **Rebellion** — Surf Wardens keeping the waves safe, Dune Traders keeping trade alive
- **Rebuilding** — Tinkerers restoring civilization, Vine Climb reclaiming the metropolis
- **Dual identity** — A character can be a `ninja` AND a `foot_clan` member — class + alignment as two independent axes

Skills by affiliation make the world feel lived-in and purposeful. Tags as universal gates create a flexible achievement system that extends far beyond factions.

---

*RFC created 2026-04-15 by Michaelangelo 🐢🟠*
*Interview update 2026-04-28 by Leonardo 🐢🔵 — all open questions resolved*
