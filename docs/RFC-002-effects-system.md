# RFC-002: Effects & Hooks System

**Status:** Implemented
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Implemented:** 2026-05-12
**Supersedes:** N/A
**Related:** RFC-001 (Web Client Architecture)

---

## 1. Executive Summary

This RFC proposes a flexible **Effects & Hooks system** that allows abilities,
items, and character templates to apply **state changes** to characters. Effects
can be triggered immediately (one-shot), attached to an event via a **hook**, or
applied directly by an ability. The system is designed to support use cases
including XP drain, bind point modification, DoT/HoT afflictions, and
administrative stat changes.

The system lives in **data** (`herbst-server/`, the REST API + DB) so that
effects and hooks are admin-configurable without code changes. The `herbst/`
game server reads effect and hook definitions from the DB and executes them at
appropriate game events.

---

## 2. Core Concepts

### 2.1 Events

Events are things that happen in the game world. They are the **trigger points**
where hooks can fire.

**Proposed event list (MVP):**

| Event | When it fires |
|-------|---------------|
| `on_death` | A character dies |
| `on_kill` | A character kills another character |
| `on_hit_received` | A character is hit by an attack |
| `on_hit_dealt` | A character lands an attack |
| `on_enter_room` | A character enters a room |
| `on_leave_room` | A character leaves a room |
| `on_equip` | A character equips an item |
| `on_unequip` | A character unequips an item |
| `on_effect_start` | An effect begins (including the one being defined) |
| `on_effect_end` | An effect expires or is removed |
| `on_login` | A character logs in |
| `on_logout` | A character logs out |
| `on_timer` | Every N seconds while the effect is active (deferred) |

Events are **code-defined constants** in `herbst/` — new events require a code
change. The event name is a string identifier stored on hook records.

### 2.2 Effects

An **effect** is the actual state change applied to a character. Effects are
**data** — defined in the DB, not in code.

**Effect types (MVP):**

| Effect Type | Parameters | Description |
|-------------|------------|-------------|
| `xp_drain` | `amount` (flat integer) | Subtracts `amount` XP from the target |
| `xp_gain` | `amount` (flat integer) | Adds `amount` XP to the target |
| `xp_set` | `amount` (flat integer) | Sets target XP to exactly `amount` |
| `bind_point_set` | `room_id` (integer) | Changes target's bind point to `room_id` |
| `hp_change` | `amount` (+/- integer) | Modifies current HP |
| `stamina_change` | `amount` (+/- integer) | Modifies current stamina |
| `mana_change` | `amount` (+/- integer) | Modifies current mana |
| `stat_mod` | `stat`, `modifier`, `duration_seconds?` | Temporarily modifies a stat multiplier (deferred) |
| `message` | `message_text`, `message_type` | Sends a text message to the character (`info`, `success`, `error`, `warn`) |
| `teleport` | `room_id` | Teleports target to `room_id` |
| `apply_effect` | `effect_id` | Applies another effect by ID to the target (effect chaining) |

**Effect model (DB):**

```go
// herbst-server/db/effect.go (ent schema)
Effect implements ent.Schema {
    Fields: {
        name          string      // "XP Drain", "Bind Point Reset"
        description   string
        effect_type   string      // "xp_drain", "bind_point_set", etc.
        parameters    JSON        // {"amount": 500}, {"room_id": 12}
        stack_mode    string      // "replace", "refresh", "stack"
        stack_limit   int         // max stacks (default 1)
        is_permanent  bool        // if true, effect does not auto-expire
        duration_secs int         // 0 = instant/one-shot; >0 = expires after N seconds
        messages      JSON        // {"on_start": "You feel weakened!", "on_end": "The curse lifts."}
    }
    Edges: {
        hooks []Hook  // hooks that apply this effect
    }
}
```

`parameters` is a JSON object. The `effect_type` field determines which
parameters are expected:

- `xp_drain` → `{"amount": 500}`
- `xp_gain` → `{"amount": 500}`
- `xp_set` → `{"amount": 0}`
- `bind_point_set` → `{"room_id": 12}`
- `hp_change` → `{"amount": -50}` (negative = damage, positive = heal)
- `stamina_change` → `{"amount": -20}`
- `mana_change` → `{"amount": 10}`
- `message` → `{"text": "You have been poisoned!", "message_type": "warn"}`
- `teleport` → `{"room_id": 5}`
- `apply_effect` → `{"effect_id": 42}`

### 2.3 Hooks

A **hook** is the binding between an event and an effect. It attaches an effect
to a character template (NPC or player template) so that when the specified
event fires on that character, the effect is applied.

**Hook model (DB):**

```go
// herbst-server/db/hook.go (ent schema)
Hook implements ent.Schema {
    Fields: {
        name          string   // "Death Drain — XP from killer"
        event         string   // "on_death", "on_hit_received", etc.
        effect_id     int      // FK to Effect
        target        string   // "self", "attacker", "room", "killer"
        condition     string   // optional SPICE condition expression
        enabled       bool
    }
    Edges: {
        characterTemplate *CharacterTemplate  // which template this hook is attached to
    }
}
```

**Target resolution:**

| Target | Meaning |
|--------|---------|
| `self` | The character the event fired on |
| `attacker` | The character who hit `self` (for `on_hit_received`) |
| `killer` | The character who killed `self` (for `on_death`) |
| `room` | All characters in the same room |
| `owner` | The character who owns the item/NPC |

**Condition (optional, deferred):**

Hooks may carry a SPICE condition expression (`herbst/` evaluates it). Example:
`character.hp < character.max_hp * 0.3` — only fires if target is below 30% HP.

### 2.4 ActiveEffect (runtime state)

When an effect is applied to a live character in-game, a record is created to
track it at runtime:

```go
// herbst-server/db/activeeffect.go (ent schema)
ActiveEffect implements ent.Schema {
    Fields: {
        character_id    int       // who it's applied to
        effect_id        int       // FK to Effect (the definition)
        applied_by_id    int       // character ID of who applied it (attacker, NPC, etc.)
        stack_count      int       // current stack count
        started_at       time.Time
        expires_at       time.Time // null if permanent
        is_active        bool
    }
}
```

ActiveEffects are stored in the DB and also cached in `herbst/` memory for
fast lookup during combat/game loops.

### 2.5 Abilities vs Effects

An **ability** is what a character can do. When an ability is used, it may:

1. **Deal damage / heal** directly (existing combat logic)
2. **Apply an effect directly** — the ability references an `effect_id`
3. **Fire a hook** — the ability is the event source; when cast, it fires a
   named event (e.g., `on_ability_cast`) which triggers attached hooks

The ability → effect relationship is **one ability can apply one or more
effects**.

```go
// Extended Ability model (existing field extended)
Ability {
    ...
    effect_ids   []int   // list of Effect IDs to apply on successful cast
    hook_event   string  // optional: fire this event name when ability is used
}
```

**Example — the "XP Drain Wand":**

```
Item: "Wand of Drain"
  → on_hit_received hook (attached to NPC wielding wand)
     → event: "on_hit_received"
     → effect: "XP Drain"
     → target: "attacker"
     → effect_type: "xp_drain"
     → parameters: {"amount": 500}
```

When NPC hits player → `on_hit_received` fires on NPC → hook resolves target
(`attacker` = player) → applies `xp_drain` effect → player's XP drops by 500.

---

## 3. Effect Lifecycle

### 3.1 Application

```
Event fires (e.g., character takes a hit)
  → Find all Hooks on that character matching event type
    → For each hook:
       → Evaluate condition (if any) — skip if false
       → Resolve target (self, attacker, room, etc.)
       → Apply effect to target
          → Check stack_mode on Effect:
              replace:  set stack_count=1, reset expires_at
              refresh:  update expires_at, increment stack_count up to limit
              stack:    increment stack_count (up to stack_limit)
          → If effect has duration > 0: create ActiveEffect record
          → Fire "on_effect_start" event for the applied effect (for chaining)
          → Send "on_start" message to target (if defined)
```

### 3.2 Effect Chaining

Effects can apply other effects via `apply_effect` type. This allows chains:

```
"Poison Curse" effect
  → on_effect_start:
      apply_effect(Effect: "Poison Damage DOT")
      apply_effect(Effect: "Bind Point Lock")

"Poison Damage DOT" effect
  → on_effect_end:
      apply_effect(Effect: "Poison Cleanse Message")
```

The chaining limit is **3 levels deep** to prevent infinite loops.

### 3.3 Expiration

For effects with `duration_secs > 0`:

```
Timer fires (checked every 1 second in herbst game loop)
  → Find all ActiveEffects where expires_at <= now AND is_active = true
    → Mark is_active = false
    → Fire "on_effect_end" event on the character
      → Any hooks on the character for "on_effect_end" → apply their effects
    → Send "on_end" message to character (if defined)
    → Delete or archive ActiveEffect record
```

### 3.4 Removal

Effects can be removed by:
- **Explicit ability/item** — applies `apply_effect` of a "Cleanse" effect
  that targets and removes ActiveEffects matching specific criteria
- **Expiration** — duration expires
- **Character death** — configurable per-effect; if `persist_on_death = false`,
  ActiveEffects are cleared on character death

---

## 4. Backend Architecture

### 4.1 Data Layer (herbst-server)

```
herbst-server/
  db/
    effect.go         — Effect ent schema
    hook.go           — Hook ent schema
    activeeffect.go   — ActiveEffect ent schema
    hook.go           — CharacterTemplate hooks edge
    ...
  services/
    effectservice.go  — ApplyEffect, GetActiveEffects, RemoveEffect, etc.
    hookservice.go    — FireEvent, EvaluateHooks, ResolveTarget
  routes/
    effect_routes.go  — CRUD for Effect definitions (admin API)
    hook_routes.go    — CRUD for Hooks (admin API)
    character_routes.go — GET /admin/characters/:id/effects (active effects)
```

### 4.2 Game Logic Layer (herbst/)

```
herbst/
  effects/
    effects.go        — ApplyEffect(), FireEvent(), EvaluateHooks()
    resolver.go       — Target resolution (self, attacker, killer, room)
    conditions.go     — SPICE condition evaluator (deferred)
    messages.go       — on_start / on_effect_end message dispatch
  model.go            — add ActiveEffects cache to character state
  combat.go           — emit on_hit_received, on_hit_dealt events
  game_model.go       — emit on_death, on_kill events
  room.go             — emit on_enter_room, on_leave_room events
```

**Key design constraint:** The herbst/ game server reads effect and hook
definitions from the DB on startup and caches them. It does **not** modify
effect definitions — only reads and executes. The admin UI manages all CRUD
on effects and hooks via the `herbst-server/` REST API.

### 4.3 Game Loop Integration Points

| Event | Where fired in herbst/ |
|-------|------------------------|
| `on_hit_received` | `combat.go` — `processHit()` after damage calculated |
| `on_hit_dealt` | `combat.go` — `processHit()` after damage confirmed |
| `on_death` | `game_model.go` — after death resolution |
| `on_kill` | `game_model.go` — after kill confirmation |
| `on_enter_room` | `game_room.go` — after room change complete |
| `on_leave_room` | `game_room.go` — before room change |
| `on_equip` | `cmd_equip.go` — after equip confirmed |
| `on_unequip` | `cmd_unequip.go` — after unequip confirmed |
| `on_login` | `auth.go` — after character auth confirmed |
| `on_effect_start` | `effects.go` — immediately after effect applied |
| `on_effect_end` | `effects.go` — when ActiveEffect expires or is removed |

---

## 5. Admin UI

### 5.1 Effects Page (`/effects`)

Full CRUD for Effect definitions. Alphabetized list.

| Column | Description |
|--------|-------------|
| Name | Effect display name |
| Type | Effect type (`xp_drain`, `bind_point_set`, etc.) |
| Parameters | Summary of parameters JSON |
| Duration | `permanent` or `Xs` |
| Stack Mode | `replace` / `refresh` / `stack` |
| Hooks | Count of hooks referencing this effect |

**Create/Edit Effect form:**

```
Name:          [text input]
Description:   [text input]
Effect Type:   [dropdown — populates params below based on selection]
Parameters:    [dynamic form based on effect type]
Duration:      [number] seconds  [checkbox: permanent]
Stack Mode:    [dropdown: replace | refresh | stack]
Stack Limit:   [number, default 1]
Messages:
  On Start:    [text input]
  On End:      [text input]
```

### 5.2 Hooks

Hooks are managed on the **Character Templates page** (NPC template detail view).

Each character template (NPC or player class template) has an expandable
**"Hooks" section**:

```
NPC: "Goblin Shaman"
  └─ Hooks
      ├─ [on_hit_received] → XP Drain (attacker) [enabled ✓]
      ├─ [on_death]       → Resurrection (self)   [enabled ✓]
      └─ [+ Add Hook]
```

**Add/Edit Hook form:**

```
Hook Name:     [text]
Event:         [dropdown — on_hit_received, on_death, on_kill, etc.]
Effect:        [dropdown — select from Effect list]
Target:        [dropdown — self | attacker | killer | room | owner]
Condition:     [text input — SPICE expression, optional, deferred]
Enabled:       [checkbox]
```

### 5.3 Character Effects Management (`/characters/:id/effects`)

Shows all **ActiveEffects** currently on a character (runtime state, not
definitions). Admin can:

- View active effects with remaining duration
- Manually remove an active effect
- Force-apply an effect to a character

---

## 6. Open Questions & Deferred Items

| Item | Status | Notes |
|------|--------|-------|
| `on_timer` event (tick-based effects) | Deferred | Duration-based is MVP; tick handled later |
| SPICE condition expressions | Deferred | MVP hooks have no condition; stub for future |
| Effect immunity/cooldown | Deferred | Not in MVP |
| Effect templates (reusable param sets) | Deferred | Each hook links directly to an Effect |
| WebSocket real-time updates to admin UI | Deferred | Page refresh polling is fine for MVP |
| DoT/HoT tick loop integration | Deferred | `duration_secs > 0` stored; tick processing deferred |
| herbst-server REST endpoint for `FireEvent` | Deferred | Events fire from herbst/ game loop; herbst-server provides data only |

---

## 7. Example Use Cases

### Use Case 1: XP Drain Wand (NPC)

```
Effect:
  name: "XP Drain"
  effect_type: "xp_drain"
  parameters: {"amount": 500}
  stack_mode: "refresh"
  is_permanent: false
  duration_secs: 0  (one-shot)

Hook (on NPC "Goblin Sorcerer" template):
  event: "on_hit_received"
  effect: "XP Drain"
  target: "attacker"
```

When Goblin Sorcerer hits a player → `on_hit_received` fires →
hook resolves target=attacker → effect drains 500 XP from player.

### Use Case 2: Bind Point Curse

```
Effect:
  name: "Shadow Bind"
  effect_type: "bind_point_set"
  parameters: {"room_id": 99}  -- the cursed chamber
  stack_mode: "replace"
  duration_secs: 0

Hook (on room "Cursed Chamber" or on NPC "Shadow Mage"):
  event: "on_enter_room"  (or "on_hit_received" from the NPC)
  effect: "Shadow Bind"
  target: "self"
```

### Use Case 3: Death Drain (vampire XP leech)

```
Effect:
  name: "Life Steal"
  effect_type: "xp_drain"
  parameters: {"amount": 200}

Hook (on NPC "Vampire Lord"):
  event: "on_kill"
  effect: "Life Steal"
  target: "self"
```

On kill → drain 200 XP from killer → XP goes to... the NPC? (Question: does the
NPC gain the XP, or does it vanish? **Sam: it just vanishes.**)

---

## 8. Acceptance Criteria (MVP)

- [ ] `Effect` ent schema with all MVP effect types
- [ ] `Hook` ent schema with event/target/effect binding
- [ ] `ActiveEffect` ent schema for runtime tracking
- [ ] `CharacterTemplate` has `hooks` edge
- [ ] CRUD routes for Effect definitions (admin API)
- [ ] CRUD routes for Hook management (admin API)
- [ ] GET `/characters/:id/effects` — list active effects on a character
- [ ] `effects.go` service in herbst/ — `ApplyEffect()`, `FireEvent()`, `ResolveTarget()`
- [ ] All MVP event integration points in herbst/ game loop
- [ ] Effects page in admin UI (alphabetized list, create/edit)
- [ ] Hooks section on character/NPC template detail view
- [ ] Character effects management panel
- [ ] Messages (`on_start`, `on_effect_end`) dispatch to character

---

## 9. TL;DR

**Effects** = what changes (XP drain 500, bind to room 99, -50 HP)  
**Hooks** = event → effect bindings attached to character templates  
**Events** = trigger points (on_death, on_hit_received, etc.)  
**Abilities** = can cast spells (which apply effects) or directly apply effects  
**ActiveEffect** = runtime record of an effect currently on a character

Data lives in `herbst-server/` (DB + REST). Game logic lives in `herbst/`
(reads DB, executes effects). Admin UI manages everything. One-shot by default,
duration/stacking configurable per-effect.
