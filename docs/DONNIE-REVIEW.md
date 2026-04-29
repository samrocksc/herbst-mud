# Donnie's Architecture Review — Faction Skills + XP RFCs

**Reviewer:** Donnie (Technical/QA) 🐢🔴
**Date:** 2026-04-28
**Status:** REVIEWED — flags and ticket structure below

---

## Key Codebase Findings

### ent ORM
The project uses **ent** (code-first ORM). All schema changes go in `server/db/schema/*.go`, then `go generate` produces the full DAO layer (create/update/query/delete). **Never edit generated files directly.**

### NPCTemplate — CLEAN ✓
`server/db/schema/npc_template.go` already has `field.Int("level").Default(1)`. The XP RFC can add `xp_value` as a clean add to this schema.

### Character — CLEAN ✓
`server/db/schema/character.go` has `level` field, NO `xp` field. The XP RFC adds `xp` cleanly. The existing flat `skill_blades`, `skill_staves`, etc. fields are **starting proficiencies** — these are separate from both the `CharacterSkill` join table (hotbar slots 1-5) and the new competency XP system.

### Skill Schema — ONE FLAG
`server/db/schema/skill.go` already has `mana_cost`, `stamina_cost`, `hp_cost`, `cooldown` (in ticks), `effect_type`, `effect_value`, `effect_duration`, `scaling_stat`, `scaling_percent_per_point`. 

**Flag:** `cooldown` is stored as ticks, not seconds. The XP RFC faction skills spec says `cooldown_seconds`. We need to decide: store seconds and convert at runtime, or use ticks? **Recommendation: store seconds in a new `cooldown_seconds` field, keep `cooldown` as legacy ticks during migration.**

**Missing fields for faction skills** (need to add):
- `slug` — globally unique skill identifier
- `faction_id` — FK to new factions table (null = universal)
- `required_tag` — tag required to unlock (beyond faction membership)
- `skill_type_enum` — `passive` vs `active` (current `skill_type` is free-form string)
- `proc_chance` — float for passives
- `proc_event` — string: `on_hit`, `on_hit_received`, `on_crit`, `on_kill`

### GameConfig — REUSE THIS
`server/db/schema/game_config.go` already implements a key/value/config table with `key` (unique) + `value` (string). The XP RFC's `CharacterConfig` should either:
(a) reuse `GameConfig` by adding a `category` or `description` field to differentiate, OR
(b) create `CharacterConfig` as a separate table.

**Recommendation: Option (a)** — add `description` to `GameConfig`. Use key prefix `character.` for character-specific config. No new table needed.

### Combat System — NO DEATH EVENT EXISTS
`herbst/combat/manager.go` and `herbst/combat/tick.go` drive combat. There is **no death event emitted today**. When an NPC's HP reaches 0, it is removed from the room — but no API call, no event. The XP RFC's event bus needs to integrate here: after NPC removal, emit `npc_killed`. This is a thin addition to the room removal path.

**Flag:** We need to find where NPC HP→0 triggers room removal. Search `herbst/` for where this happens — it's likely in `game_room.go` or the room's character management. This is the integration point for `npc_killed` emission.

### No Service Layer — CRITICAL REFACTOR
Every route in `server/routes/*.go` has business logic in anonymous closures. No service packages exist. The XP RFC correctly identifies this as a prerequisite. The refactor pattern:

```
// server/services/character.go
package services

type CharacterService struct { client *db.Client }

func (s *CharacterService) AwardXP(ctx context.Context, charID int, amount int) (*XPAwardResult, error) {
    // business logic + DB update + event emission
}

// server/routes/character_routes.go (after refactor)
func RegisterCharacterRoutes(router *gin.Engine, client *db.Client) {
    svc := &services.CharacterService{client: client}
    router.POST("/characters/:id/xp/award", func(c *gin.Context) {
        var req XPAwardRequest
        if err := c.ShouldBindJSON(&req); err != nil { ... }
        result, err := svc.AwardXP(c.Request.Context(), charID, req.XP)
        // ...
    })
}
```

### No Event Bus
No pub/sub system exists. Recommend: simple in-process event bus in `server/events/events.go`. No external message queue needed for a MUD server.

---

## Faction RFC Flags

1. **Skills `slug` uniqueness** — faction skills need globally unique slugs. Add unique constraint in schema.
2. **Character creation tag auto-grant** — need to find where character creation happens (`server/routes/character_routes.go` POST /characters) and add `first_class` tag auto-grant there.
3. **Faction category `max_memberships`** — verify no existing `character_classes` table. There isn't — `Character.class` is a single string. The category slot system is net-new.
4. **Content externalization** — faction data (categories, factions, required tags) should be DB-backed, seeded from YAML. This is part of #282.
5. **Proc event `on_kill`** — this ties XP (kill → npc_killed event → XP subscriber) to faction passive skills. The proc system and XP event bus are independent but both need `on_kill` event.

---

## XP RFC Flags

1. **CharacterConfig vs GameConfig** — see recommendation above. Reuse `GameConfig` with a description field.
2. **NPC death integration point** — find where NPC HP→0 triggers removal, emit `npc_killed` there. This is the hook for XP and competency award.
3. **Character competency `category_id` FK** — competency_category table should use string PK (`blades`, `staves`) to match the RFC design.
4. **Cooldown units** — ticks vs seconds. Recommend storing seconds in new field.
5. **Level threshold table seeding** — seed data needed on first run. Can use an init migration or seed SQL file.
6. **Resurrection sickness debuff** — no status effect/debuff system exists. Basic implementation: reduce stats for duration. Combat prevention is v2.

---

## Recommended Ticket Structure

### Service Layer & Event Bus (Prerequisite for both RFCs)
1. **SVC-001: Extract service layer from route handlers** — refactor character routes into `server/services/character.go`. Do for all routes in `server/routes/`. Emit events from services, not handlers.

2. **SVC-002: Implement in-process event bus** — `server/events/events.go`. Support `npc_killed`, `quest_completed`, `character_died`, `character_leveled_up`. Subscribers register on init.

### XP System (depends on SVC-001 + SVC-002)
3. **XP-001: Add xp field to Character + xp_value to NPCTemplate** — ent schema migrations. Seed LevelThreshold table with defaults (levels 2-50).

4. **XP-002: Implement CharacterConfig using GameConfig table** — add `description` field to GameConfig. Insert keys: `max_level=50`, `xp_level_diff_scale_factor=0.10`, all `death_penalty_*` keys.

5. **XP-003: Implement XP award service + REST endpoint** — award XP, check level-up, apply stat bonuses from ClassLevelBonus. Emit `character_leveled_up` event.

6. **XP-004: Implement NPC kill → npc_killed event emission** — find NPC HP→0 removal point, emit event. Wire XP subscriber to event bus.

7. **XP-005: Competency system** — CompetencyCategory table, CompetencyLevelThreshold table, CharacterCompetency join table. Competency XP award on npc_killed.

8. **XP-006: Death penalty subscriber** — subscribe to `character_died`, apply configured penalties from CharacterConfig.

9. **XP-007: Quest XP event subscriber** — subscribe to `quest_completed`, award XP from QuestConfig.

10. **XP-008: Admin UI + TUI for XP content** — NPC template editor (with xp_value), level threshold editor, CharacterConfig editor, competency category/threshold editors.

### Faction Skills (can run parallel to XP, shares service layer)
11. **FACTION-001: Schema for factions** — faction_categories, factions (add category_id), character_tags, faction_required_tags. Update character_factions with status field.

12. **FACTION-002: Enhance Skill schema** — add slug (unique), faction_id, required_tag, skill_type_enum, proc_chance, proc_event fields.

13. **FACTION-003: Faction REST endpoints** — join/leave/expel/reputation, faction eligibility check, tag CRUD.

14. **FACTION-004: Skill eligibility logic** — update skill listing to check faction membership + tags. Integrate proc system (on_hit, on_kill, etc.) into combat.

15. **FACTION-005: Tag auto-grant on character creation** — grant `first_class` tag when character is created.

16. **FACTION-006: Admin UI + TUI for factions** — faction category management, faction membership viewer, tag grant/revoke, reputation editor.

### YAML → DB Migration (shared, blocks both)
17. **#282: Migrate content YAML to database** — already created. Blocking for NPC template xp_value and faction content seeding.

---

## Implementation Priority Order

1. **SVC-001** (service layer) — blocks everything
2. **SVC-002** (event bus) — blocks XP and combat events
3. **XP-001** + **XP-002** (schema + config) — can run together
4. **XP-003** + **XP-004** (XP award + NPC kill event) — core loop
5. **FACTION-001** + **FACTION-002** (schema) — can run parallel to XP-001/002
6. **XP-005** through **XP-008** — XP remaining
7. **FACTION-003** through **FACTION-006** — faction remaining
8. **#282** — YAML→DB content migration (continuous)

---

*Reviewed by Donnie 🐢🔴*
