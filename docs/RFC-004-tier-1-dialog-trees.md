# RFC-004: NPC Dialog System — Tier 1 (Dialog Trees)

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Supersedes:** N/A
**Related:** RFC-002 (Effects & Hooks), RFC-003 (Quest System)
**Next:** RFC-005 (Tier 2 — Keyword-Intent Matching)

---

## 1. Executive Summary

This RFC proposes a **branching dialog tree system** for NPCs — the foundation
of NPC interaction in herbst-mud. Players use a `talk` command to initiate
conversation with an NPC, then navigate a tree of numbered responses. Dialogs
are 100% data-driven: admins define dialog nodes and response options via the
web admin UI. Dialog nodes can trigger quest offers, apply effects, and
conditionalize options based on game state.

This is **Tier 1** of a three-tier NPC dialog architecture. It prioritizes
simplicity, determinism, and admin tooling. Tier 2 (keyword-intent matching)
and Tier 3 (LLM-powered conversation) build on this foundation.

---

## 2. Design Philosophy

- **Data, not code.** Admins should be able to create and edit dialog trees
  without touching the Go codebase.
- **Numbered choices.** Classic MUD interaction — clear, predictable, testable.
- **Extensible.** Each dialog node is a stepping stone to Tier 2 intent matching
  and Tier 3 LLM integration.
- **Hooks into existing systems.** Dialog nodes trigger quest offers and effects
  using the primitives already defined in RFC-002 and RFC-003.

---

## 3. Core Concepts

### 3.1 DialogNode

A **DialogNode** is a single point in an NPC's conversation tree. It contains
what the NPC says and the player's available responses.

```go
// server/db/schema/dialog_node.go
func (DialogNode) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").Unique(),
        field.String("npc_template_id"),       // which NPC template this dialog belongs to
        field.String("npc_text").              // what the NPC says at this node
        field.JSON("responses", []DialogResponse{}), // player's options
        field.Bool("is_entry").Default(false),  // first node when conversation starts
        field.String("entry_condition").        // SPICE expression: character.tags.has("wizard_complete")
            Optional(),
        field.JSON("on_enter_effects", []int{}), // effect IDs applied when this node is reached
    }
}

// Edges
func (DialogNode) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("npc_template", NPCTemplate.Type).
            Ref("dialog_nodes").
            Unique(),
    }
}
```

### 3.2 DialogResponse

A player's response option:

```go
type DialogResponse struct {
    Label         string `json:"label"`          // display text: "What troubles you, elder?"
    NextNodeID    string `json:"next_node_id"`   // dialog node this leads to (empty = end conversation)
    Condition     string `json:"condition"`      // SPICE: only show if condition met (optional)
    QuestOfferID  string `json:"quest_offer_id"` // quest to offer when selected (optional)
    Effects       []int  `json:"effects"`        // effects applied when selected (optional)
}
```

### 3.3 Conversation State (runtime)

Per-player conversation state held in memory (not DB) while a conversation is
active:

```go
// herbst/dialog/state.go
type ConversationState struct {
    NPCName       string    // "Elder Myrddin"
    NPCTemplateID string    // "npc_elder_myrddin"
    CurrentNodeID string    // "node_001_greeting"
    StartedAt     time.Time
}
```

Stored on the character's `tea.Model` as `conversation *ConversationState`.

---

## 4. Player Commands

```
talk <npc_name>    — start conversation with an NPC in the current room
talk               — resume conversation with the last NPC if still in room
<1-9>              — select response option
0 / leave / bye    — end conversation
```

**Example session:**

```
> look
You are in a dimly lit chamber. Exits: north, south.
Elder Myrddin is here, studying an ancient tome.

> talk myrddin

Elder Myrddin studies you with ancient eyes.
1. "What troubles you, elder?"
2. "I'm looking for work."
3. [Leave]

> 1

Elder Myrddin sighs heavily.
"The shadows grow long in the northern caves. Goblin shamans
have been seen near the old temple."
1. "Tell me more about these shamans."
2. "I'll handle it." (Accept quest)
3. "Not my problem." [Leave]
```

**When `talk` is entered without a name:** if a conversation is already active
with an NPC in the room, resume it. If no conversation is active, prompt:
"Talk to whom?"

---

## 5. Dialog Tree Resolution

When the player types `talk <npc_name>`:

1. Look up the NPC in the current room by name match
2. Load the NPC's template → find all `DialogNode` records for that template
3. Find the **entry node** (`is_entry = true`)
4. If multiple entry nodes exist (with different `entry_condition`s):
   evaluate conditions; pick the first matching one
5. Set `conversation.CurrentNodeID = entryNode.ID`
6. Render: NPC text + numbered responses

When the player selects `1-9`:

1. Look up `CurrentNode.responses[index]`
2. If `next_node_id` is empty: end conversation, apply any effects
3. If `next_node_id` is set:
   - Find the target dialog node by ID
   - Evaluate `entry_condition` if present; skip if false → fallback to default
   - Apply `on_enter_effects`
   - If the response has `quest_offer_id`: offer the quest (accept/decline prompt)
   - Set `CurrentNodeID = next_node_id`
   - Render new node

---

## 6. Quest Integration

A dialog response can offer a quest. When the player selects such an option:

```
> 2

Elder Myrddin: "Brave words. Defeat the three shamans and explore the
chamber they guard. Return to me when it's done."

[Quest Available: The Shadows Grow Long]
1. Accept Quest
2. Decline
```

This uses the `QuestOfferID` field on `DialogResponse`. Accept triggers
`questservice.AcceptQuest()` from RFC-003. Decline advances to an optional
`decline_node_id` or ends the conversation.

Add to `DialogResponse`:
```go
type DialogResponse struct {
    ...
    DeclineNodeID  string `json:"decline_node_id"`  // node shown after declining quest
}
```

---

## 7. Effects Integration

Dialog nodes and responses can apply effects (RFC-002):

- `on_enter_effects` on `DialogNode` — fire when the node is reached
- `effects` on `DialogResponse` — fire when the player selects that response

This allows:

- Quest completion triggers: `tag_add("quest_shadows_complete")` on a node
- XP rewards for dialogue: `xp_gain(50)` on a response
- State changes mid-conversation: NPC becomes hostile, room changes, etc.

---

## 8. Admin UI

### 8.1 Dialog Editor (on NPC Template detail)

```
NPC: "Elder Myrddin"
  └─ Dialog Tree

Entry Node: [node_001_greeting ▾]
[+ Add Entry Node]

Nodes:
┌─────────────────────────────────────────────────────────┐
│ [node_001_greeting] ← entry                             │
│ NPC says: "Welcome, traveler. What brings you here?"   │
│                                                          │
│ Responses:                                              │
│ 1. "What troubles you?"  →  [node_002_troubles]        │
│    [cond: ]  [quest: ]  [effects: ]                    │
│ 2. "I'm looking for work." → [node_003_work]           │
│ 3. [Leave]  ← end conversation                         │
│                                                          │
│ On Enter Effects: [none]                                │
│                                                          │
│ [✎ Edit Node] [🗑 Delete] [+ Add Response]              │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ [node_002_troubles]                                      │
│ NPC says: "The shadows grow long in the northern..."    │
│                                                          │
│ Responses:                                              │
│ 1. "Tell me more." → [node_003_shaman_details]         │
│ 2. "I'll handle it." → [node_004_accept]               │
│    [quest_offer: quest_shadows_grow_long]               │
│ 3. "Not my problem." → end conversation                │
│                                                          │
│ On Enter Effects: [none]                                │
│                                                          │
│ [✎ Edit Node] [🗑 Delete] [+ Add Response]              │
└─────────────────────────────────────────────────────────┘
[+ Add Node]
```

### 8.2 Visual Design

- Nodes are listed in a vertical accordion/card stack
- Each node shows: NPC text at top, responses as a numbered list below
- Each response shows: label → target node name (or "[end]"), any conditions/quests/effects
- Add Node / Add Response buttons are prominent
- Delete confirmation is required (per Sam's deletion boundary)

### 8.3 Dialog Testing

A "Test Dialog" button opens a simulated conversation in a slide-out panel
so admins can walk through the tree without launching the game.

---

## 9. Ent Schema Summary

**New schema:** `dialog_node`

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (unique) | Node identifier |
| `npc_template_id` | string (FK) | Parent NPC template |
| `npc_text` | string | What NPC says |
| `responses` | JSON | Array of DialogResponse |
| `is_entry` | bool | Is this the conversation start node? |
| `entry_condition` | string? | SPICE gate for entry |
| `on_enter_effects` | JSON | Effect IDs to apply |

**Modified schema:** `npc_template`

Add edge:
```go
edge.To("dialog_nodes", DialogNode.Type),
```

---

## 10. Acceptance Criteria

- [ ] `dialog_node.go` ent schema with all fields
- [ ] `NPCTemplate` has `dialog_nodes` edge
- [ ] `ent generate` clean
- [ ] `talk <npc_name>` command in herbst TUI — starts conversation
- [ ] Numbered response selection (1-9) — navigates tree
- [ ] `0` / `leave` / `bye` — ends conversation
- [ ] `ConversationState` stored on model in memory
- [ ] Multiple entry nodes with conditions resolved correctly
- [ ] Quest offer flow from dialog responses
- [ ] Effects application from dialog nodes and responses
- [ ] Dialog editor in admin UI on NPC template detail
- [ ] Node CRUD with response management
- [ ] Dialog test simulator in admin UI
- [ ] `npm run build` passes
- [ ] Unit tests for dialog tree traversal

---

## 11. Example: "Elder Myrddin" Full Dialog Tree

```
[node_001_greeting] ← entry
  NPC: "Elder Myrddin studies you with ancient eyes."
  R1: "What troubles you?" → node_002_troubles
  R2: "I'm looking for work." → node_003_work
  R3: [Leave]

[node_002_troubles]
  NPC: "The shadows grow long in the northern caves. Goblin shamans
       have been seen near the old temple."
  R1: "Tell me more about these shamans." → node_004_shamans
  R2: "I'll handle it." → node_005_accept
      [quest_offer: quest_shadows_grow_long]
  R3: "Not my problem." → end

[node_003_work]
  NPC: "There is always work for willing hands. What skills do you possess?"
  R1: "I can fight." → node_002_troubles
  R2: "I prefer words to steel." → node_006_diplomacy
  R3: [Leave]

[node_004_shamans]
  NPC: "Three of them, by my count. They carry cursed wands that drain
       the very life force from their victims. Face them prepared."
  R1: "What else should I know?" → node_007_preparation
  R2: "I'll handle it." → node_005_accept
      [quest_offer: quest_shadows_grow_long]
  R3: [Leave]

[node_005_accept]
  NPC: "Brave words. Defeat the three shamans and explore the chamber
       they guard. Return to me when it's done."
  [auto: quest offer popup]
  On accept: apply effect tag_add("quest_shadows_accepted")
  → end conversation

[node_006_diplomacy]
  NPC: "Words can be sharper than swords in the right hands. But these
       goblins listen only to violence, I fear."
  R1: "Very well. I'll go." → node_002_troubles
  R2: [Leave]

[node_007_preparation]
  NPC: "Their wands cannot be parried by steel alone. Carry something
       blessed if you can find it. And watch your back."
  R1: "I'll handle it." → node_005_accept
      [quest_offer: quest_shadows_grow_long]
  R2: [Leave]
```
