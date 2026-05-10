# RFC-003: Quest System

**Status:** Final
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Related:** RFC-002 (Effects & Hooks) — EFF-001 through EFF-009

---

## 1. Executive Summary

This RFC proposes a **Quest system** built as a layer on top of the existing
Effects, Hooks, Tags, and Achievement infrastructure. Quests are defined as
data in `herbst-server/` and executed by the `herbst/` game server. The system
supports kill quests, exploration quests, fetch quests, and quest chains.

**Design philosophy:** Minimal new concepts — quests are a *wrapper* around
existing primitives (hooks, effects, tags, achievements). No new runtime state
machinery if existing primitives suffice.

---

## 2. Core Concepts

### 2.1 Quest Definition (data, not code)

A **Quest** is a named, sequenced set of **Objectives** that a character
progresses through. It lives in the DB as a record.

```
Quest
  name
  description
  prerequisites     — list of quest IDs that must be completed first
  objectives[]       — ordered list of objectives
  rewards            — XP, items, effects, tags, achievements
  repeat_mode        — none | cooldown | always
  cooldown_hours      — if repeat_mode = cooldown
```

### 2.2 Objectives

Each objective has a **type** and a **target**:

| Objective Type | Target | Description |
|---------------|--------|-------------|
| `kill` | NPC template ID | Kill N of this NPC |
| `explore` | Room ID | Enter this room |
| `collect` | Item template ID | Obtain N of this item |
| `deliver` | NPC template ID | Deliver something to this NPC (fetch quest handoff) |
| `talk` | NPC template ID | Have a conversation with this NPC |
| `custom` | Event name | Fire a custom event via effects system |

**Example objectives:**

```
1. [kill]        Goblin Shaman (id: goblin_shaman_001)  × 3
2. [explore]     The Cursed Chamber (id: room_099)
3. [collect]     Ancient Key (id: item_key_ancient)     × 1
4. [deliver]     Elder Myrddin (id: npc_elder_myrrdin)  — return the key
```

### 2.3 Quest Progress (runtime state)

Each character's active quest state is stored in the DB:

```
QuestProgress
  character_id
  quest_id
  status          — active | completed | failed | abandoned
  started_at
  completed_at
  current_step    — index into objectives[]
  objective_counts — JSON map of objective ID → current count
```

A character can have **multiple active quests** simultaneously.

### 2.4 How Objectives Are Tracked

The system hooks into **existing game events** via the Effects Hooks system:

| Objective | Event Hook | Notes |
|-----------|------------|-------|
| `kill` NPC | `on_kill` → hook on NPC template | Hook fires `checkKillObjective(char, npcTemplateID)` |
| `explore` room | `on_enter_room` | Hook fires `checkExploreObjective(char, roomID)` |
| `collect` item | `on_item_pickup` | Fire from `cmd_take.go` or item pickup event |
| `deliver` to NPC | `on_npc_interact` with NPC ID | Hook on NPC for `on_npc_interact` |
| `talk` to NPC | `on_npc_interact` | Same as deliver; NPC dialogue event |
| `custom` | custom event name | Fire from effects system via `apply_effect` type `fire_custom_event` |

**For each hook**, the logic is:
```
if character has active quest with matching objective:
  increment objective_counts[objective_id]
  if objective_counts[objective_id] >= required_count:
    advance current_step
    if current_step == len(objectives):
      mark quest as completed
      apply rewards
```

This logic lives in a `questservice/` inside `herbst/`, not in the hook handler
itself. The hook is a thin bridge: it calls `questservice.CheckProgress(...)`.

---

## 3. Quest Definition Schema

```go
// herbst-server/db/schema/quest.go
func (Quest) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique(),
        field.String("description"),
        field.Strings("prerequisite_quest_ids"),    // quest IDs that must be completed first
        field.JSON("objectives", []QuestObjective{}),  // ordered objectives
        field.JSON("rewards", QuestRewards{}),
        field.Enum("repeat_mode").Values("none", "cooldown", "always").Default("none"),
        field.Int("cooldown_hours").Optional(),     // if repeat_mode = cooldown
        field.Bool("is_active").Default(true),      // can be taken or not
    }
}

type QuestObjective struct {
    Type      string `json:"type"`   // kill | explore | collect | deliver | talk | custom
    TargetID  string `json:"target_id"`  // NPC template ID, room ID, item ID
    Count     int    `json:"count"`  // required count (default 1)
    Label     string `json:"label"`  // display: "Kill 3 Goblin Shamans"
    Hint      string `json:"hint"`   // optional hint shown to player
}

type QuestRewards struct {
    XP           int      `json:"xp"`
    ItemIDs      []string `json:"item_ids"`      // items to give
    EffectIDs    []int    `json:"effect_ids"`    // effects to apply
    TagAdds      []string `json:"tag_adds"`      // tags to add to character
    TagRemoves   []string `json:"tag_removes"`   // tags to remove
    AchievementIDs []int  `json:"achievement_ids"` // achievements to unlock
}
```

---

## 4. Quest Progress Schema

```go
// herbst-server/db/schema/quest_progress.go
func (QuestProgress) Fields() []ent.Field {
    return []ent.Field{
        field.Enum("status").Values("active", "completed", "failed", "abandoned").Default("active"),
        field.Time("started_at"),
        field.Time("completed_at").Optional(),
        field.Int("current_step").Default(0),  // index into objectives[]
        field.JSON("objective_counts", map[string]int{}),  // objective_target_id → count
    }
}

func (QuestProgress) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("character", Character.Type).Ref("quest_progress").Unique(),
        edge.From("quest", Quest.Type).Ref("progress").Unique(),
    }
}
```

---

## 5. Quest Lifecycle

### 5.1 Acquiring a Quest

Player talks to an NPC or uses a command:

```
/quest accept <quest_id>
```

System checks:
1. `quest.is_active == true`
2. All `prerequisite_quest_ids` are completed by this character
3. Character does not already have this quest active
4. If `repeat_mode = cooldown`: last completion was > `cooldown_hours` ago

If all pass: create `QuestProgress` record, send quest description to player,
register hooks for this quest's objectives on the character.

### 5.2 Progress Tracking

Each objective has a hook registered at quest-accept time. The hook fires on
game events and calls into `questservice.CheckProgress()`:

```go
// In herbst/effects/hook_handlers.go
func OnKill(s *ssh.Session, killerID, victimTemplateID int) {
    questservice.CheckProgress(killerID, "kill", victimTemplateID)
}
```

`CheckProgress`:
1. Find all active `QuestProgress` for the character
2. Find all objectives in those quests matching (`type=kill`, `target_id=victimTemplateID`)
3. Increment `objective_counts`
4. If step complete: advance `current_step`, send progress message
5. If all steps done: complete quest, apply rewards

### 5.3 Completion

When `current_step` reaches `len(objectives)`:

1. Set `QuestProgress.status = "completed"`, `completed_at = now()`
2. Apply rewards (XP, items, effects, tags, achievements)
3. Send `"Quest Complete: <name>"` message to player with reward summary
4. Fire `on_quest_complete` event (can trigger hooks/effects)
5. Deregister objective hooks for this quest

### 5.4 Failure and Abandonment

- **Fail** — if implemented: some quests may have a `fail_conditions` array
  (e.g., timer expires). On fail: `status = "failed"`, no rewards.
- **Abandon** — player types `/quest abandon <quest_id>`:
  `status = "abandoned"`, hooks deregistered. Player can re-accept if repeat allows.

### 5.5 Repeat

| Mode | Behavior |
|------|---------|
| `none` | One-time only; re-accept blocked |
| `cooldown` | Can re-accept after `cooldown_hours` |
| `always` | Can re-accept immediately; fresh `QuestProgress` each time |

---

## 6. NPC Interaction — Quest Dialog

NPCs with quests need to surface available quests when talked to. The existing
`greeting` field on `NPCTemplate` can be extended or a separate `quest_greeting`
field added:

```
NPC: "Elder Myrddin"
  greeting: "Welcome, traveler."
  quest_greeting: "The shadows grow long. Will you help me?"   (if player has quest)
  available_quests: ["quest_shadows_grow_long"]                  (quests NPC offers)
```

When player `talk`s to an NPC:
1. Fire `on_npc_interact` event
2. If NPC has `available_quests` and player hasn't accepted them → offer them
3. Hook on `on_npc_interact` → effects system handles quest-offer flow

**Simple approach for MVP:** A `quest_offer` field on `NPCTemplate` listing
quest IDs. If non-empty, talking to the NPC automatically offers those quests.

```go
field.Strings("available_quests").Optional().
    Comment("Quest IDs this NPC offers")
```

---

## 7. Quest Chains

Quests can have `prerequisite_quest_ids`. If a quest requires `quest_A` to be
completed first, the player cannot accept it until `quest_A` is in
`QuestProgress.status = "completed"`.

Chains are expressed purely through the data — no special code needed.

---

## 8. Integration with Effects System (RFC-002)

Quest completion naturally applies effects:

```
QuestRewards:
  xp            → xp_gain effect
  effect_ids     → apply_effect for each
  tag_adds       → tag_add effect
  tag_removes    → tag_remove effect
  achievement_ids → unlock achievement (calls achievementservice.Unlock())
  item_ids       → add item to character inventory
```

**Hooks for objective tracking** (registered at quest accept):

| Objective | Hook event | Target entity |
|-----------|------------|---------------|
| `kill` | `on_kill` | NPC template |
| `explore` | `on_enter_room` | Room |
| `collect` | `on_item_pickup` | Item |
| `deliver` | `on_npc_interact` | NPC |
| `talk` | `on_npc_interact` | NPC |
| `custom` | `on_custom_event` (fire via effect) | Event name |

Each quest-registerable hook is a **one-line bridge** into `questservice.CheckProgress()`.
The heavy logic (find matching quests, increment counts, check completion) lives in
the quest service — not in the effect system.

---

## 9. Rewards — Items

Item rewards are applied directly by the quest service at completion time, not
as effects. The `inventory.go` or `character` service already has `AddItem()`
capability:

```go
// questservice/complete.go
func CompleteQuest(charID, questID int) error {
    // ...
    for _, itemID := range quest.Rewards.ItemIDs {
        if err := inventory.AddItem(charID, itemID); err != nil {
            return err
        }
    }
    // ...
}
```

---

## 10. Player-facing Commands

```
/quest list          — show active quests with progress
/quest log           — show completed/failed/abandoned quests
/quest accept <id>   — accept a quest (from NPC offer or direct ID)
/quest abandon <id>  — abandon active quest
/quest info <id>     — show quest description and objective details
```

NPC interaction flow:
```
> talk to Elder Myrddin
Elder Myrddin: "The shadows grow long. Will you help me?"
  [Quest Available: The Shadows Grow Long]
  1. Accept Quest
  2. Decline
```

---

## 11. Admin UI Requirements

### 11.1 Quests Page (`/quests`)

Alphabetized list of all quest definitions.

| Column | Description |
|--------|-------------|
| Name | Quest name |
| Objectives | Count + summary of objective types |
| Repeat | none / cooldown / always |
| Prerequisites | Quest names required |
| Active | yes/no toggle |

**Create/Edit Quest form:**

```
Name:              [text]
Description:       [textarea]
Repeat Mode:       [dropdown: none | cooldown | always]
Cooldown Hours:    [number] (shown if cooldown)
Prerequisites:     [multi-select of other quests]
Active:            [checkbox]

Objectives: [sortable list]
  [+ Add Objective]
  ┌──────────────────────────────────────────────┐
  │ Type: [kill]  Target: [Goblin Shaman ▾]     │
  │ Count: [3]    Label: "Kill 3 Goblin Shamans" │
  │ Hint: "They lurk in the northern caves."    │
  └──────────────────────────────────────────────┘
  [+ Add Objective]

Rewards:
  XP:             [number]
  Items:          [multi-select item templates]
  Effects:        [multi-select effects]
  Tags to Add:    [multi-select]
  Tags to Remove: [multi-select]
  Achievements:   [multi-select]
```

### 11.2 Quest Progress on Character

On the character detail screen:

```
Active Quests:
  The Shadows Grow Long  [step 2/4]  "Explore the Cursed Chamber"
  └─ [✓] Kill 3 Goblin Shamans
     [ ] Explore the Cursed Chamber
     [ ] Collect the Ancient Key
     [ ] Deliver to Elder Myrddin
```

---

## 12. MVP Scope

**In MVP:**
- Kill, explore, collect objectives
- Quest chains via prerequisite_quest_ids
- XP, item, tag, effect rewards on completion
- Repeat mode: `none` and `cooldown`
- Quest log (active / completed / abandoned)
- Admin UI full CRUD on quests

**Deferred:**
- `deliver` / `talk` objectives (need `on_npc_interact` event)
- `custom` event objectives
- Quest fail conditions
- Quest timer/deadlines
- NPC quest-greeting flow (can be added with NPC template `available_quests` field)
- Achievement rewards (already have achievement schema — just needs integration)
- Quest-level XP reward display in character sheet

---

## 13. Acceptance Criteria

- [x] `quest.go` ent schema with all fields
- [x] `quest_progress.go` ent schema with all fields
- [x] `questservice/` in `herbst/` — `AcceptQuest()`, `CheckProgress()`, `AbandonQuest()`
- [x] Hook bridge functions for each objective type in `herbst/` (on_kill, on_enter_room via FireEvent)
- [x] Quest REST API routes (admin CRUD + character quest progress)
- [x] Quests page in admin UI (list + detail/edit with objectives)
- [ ] Character quest progress panel in admin UI (deferred — accessible via API)
- [x] Player-facing `/quests`, `/quest accept <id>`, `/quest abandon <id>` commands
- [ ] Item reward application on quest completion (deferred — rewards returned in response but not auto-applied to inventory)
- [x] Effect rewards returned in completion response (applied via effects system)
- [x] `repeat_mode=cooldown` enforced on re-accept

---

## 14. Example Quest: "The Shadows Grow Long"

```json
{
  "id": "quest_shadows_grow_long",
  "name": "The Shadows Grow Long",
  "description": "Elder Myrddin has noticed the growing darkness in the northern caves. Goblin shamans have been darkening the land. Find and defeat three of them, then explore the cursed chamber they guard.",
  "repeat_mode": "cooldown",
  "cooldown_hours": 168,
  "prerequisite_quest_ids": [],
  "objectives": [
    {
      "type": "kill",
      "target_id": "goblin_shaman_001",
      "count": 3,
      "label": "Defeat Goblin Shamans",
      "hint": "They are found in the northern caves."
    },
    {
      "type": "explore",
      "target_id": "room_cursed_chamber",
      "count": 1,
      "label": "Explore the Cursed Chamber",
      "hint": "The shamans guard something deep within."
    }
  ],
  "rewards": {
    "xp": 500,
    "item_ids": ["item_shadow_key"],
    "effect_ids": [],
    "tag_adds": ["quest_shadows_complete"],
    "tag_removes": [],
    "achievement_ids": []
  }
}
```

---

## 15. Open Questions

1. **Collect objective — item template vs. instanced item?** If `collect` targets
   an item template ID, the count is "items of this type picked up." But if the
   player already has 2, picks up 1 more, the count goes from 2→3. This works.
   But does picking up ANY item of that type count, or only specific containers?
2. **NPC quest offer flow** — do you want an explicit "quest offer" dialog with
   Accept/Decline buttons, or just auto-add the quest when the player talks to
   the NPC? Auto-add is simpler; explicit offer is more RPG-flavored.
3. **Deliver objective** — does the player need to have a specific item in their
   inventory (collected via a prior objective), or does the deliver objective
   auto-complete when talking to the NPC if the item is present?
4. **How are quests discovered?** Through NPC greetings, through `/quest list`
   (shows all available quests), or both?
