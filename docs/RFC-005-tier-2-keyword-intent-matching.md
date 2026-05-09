# RFC-005: NPC Dialog System — Tier 2 (Keyword-Intent Matching)

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Builds on:** RFC-004 (Tier 1 Dialog Trees)
**Next:** RFC-006 (Tier 3 LLM-Powered NPCs)
**Related:** RFC-002 (Effects & Hooks), RFC-003 (Quest System)

---

## 1. Executive Summary

This RFC extends the Tier 1 dialog tree system with **free-form text input**
and **keyword-intent matching**. Instead of numbered choices, the player types
natural language ("tell me about the goblins"), and the system matches it
against a database of `DialogIntent` records defined per NPC. This is Tier 2
of the three-tier NPC dialog architecture.

**Why this matters:** It bridges the gap between rigid dialog trees (Tier 1)
and full LLM conversation (Tier 3). Dungeon Masters can pre-author
information-rich NPCs that respond to a broad vocabulary of player questions
without any AI dependency. When Tier 3 LLM arrives, this system serves as the
deterministic fallback.

---

## 2. Core Concepts

### 2.1 DialogIntent

A `DialogIntent` is a keyword-to-response mapping for an NPC. It defines a
topic the NPC can talk about and the trigger words that invoke it.

```go
// server/db/schema/dialog_intent.go
func (DialogIntent) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").Unique(),
        field.String("npc_template_id"),          // which NPC knows this
        field.String("intent").                   // "query_caves", "query_goblins", "greeting"
        field.Strings("trigger_words"),            // ["cave", "north", "dark"]
        field.String("response_text"),             // what the NPC says
        field.String("fallback_text").             // shown when this intent is active but input doesn't advance
            Default("I'm not sure what you mean."),
        field.Int("priority").Default(0),          // match tiebreaker
        field.Bool("one_shot").Default(false),     // if true, removed from pool after firing once
        field.JSON("tags_required", []string{}),   // character must have these tags to trigger
        field.JSON("tags_forbidden", []string{}),  // character must NOT have these
        field.Int("cooldown_messages").Optional(), // how many player messages before this fires again
        field.Bool("is_greeting").Default(false),  // fire automatically when conversation starts
        field.JSON("on_match_effects", []int{}),   // effects applied on match
        field.JSON("set_current_node_id", ""),     // optional fallback to Tier 1 node after intent fires (for structured follow-ups like quest acceptance)
    }
}
```

### 2.2 Matching Algorithm

When the player types free-form text during a conversation:

```
1. Tokenize input → lowercase, strip punctuation, extract keywords
2. For each DialogIntent for this NPC:
   a. If tags_required are not all present on character → skip
   b. If tags_forbidden contains any tag the character has → skip
   c. Score = count of trigger_words that appear in input
   d. If score == 0 → skip
3. Sort candidates by: score DESC, then priority DESC
4. Pick the top candidate
   - If score > 0 → return response_text, apply on_match_effects
   - If score == 0 AND this intent is the active fallback → return fallback_text
   - If no intent matches → return generic fallback: "Elder Myrddin looks confused."
5. If intent.one_shot → remove from pool for this conversation
```

### 2.3 Greeting Intents

When `talk <npc>` begins, the system fires all intents with `is_greeting = true`
(in priority order). This replaces the Tier 1 entry node for free-form mode.

```
> talk myrddin

Elder Myrddin looks up from his tome.
"Welcome, traveler. These are dark times."

> what's happening
Elder Myrddin: "Goblin shamans in the northern caves. I've been tracking them."

> goblins?
Elder Myrddin: "Three of them carry cursed wands. They drain life itself.
Few come back."

> i'll go
Elder Myrddin: "Then take this. You'll need it."

[Quest Available: The Shadows Grow Long]
1. Accept Quest
2. Decline
```

### 2.4 Fallback to Tier 1 Nodes

Some conversation moments need structure (quest acceptance, specific choices).
The `set_current_node_id` field on `DialogIntent` lets an intent transition
the conversation into a **Tier 1 dialog node** after firing:

```json
{
  "intent": "accept_quest_offer",
  "trigger_words": ["i'll go", "i'll handle it", "let's do this", "i accept"],
  "response_text": "Then take this. You'll need it.",
  "set_current_node_id": "node_005_accept",
  "on_match_effects": []
}
```

After this intent fires, the conversation drops into Tier 1 mode at
`node_005_accept`, which shows the structured quest offer popup with
Accept/Decline options. Once the quest is accepted or declined, free-form
mode resumes.

This creates a **fluid hybrid**: free-form exploration → intent match →
structured quest interaction → back to free-form.

### 2.5 Cooldown

`cooldown_messages` prevents the same intent from firing on every message.
For example, a greeting intent fires once. A "tell me more" intent with
`cooldown_messages: 3` won't fire again for 3 more player messages. This
creates natural conversational flow instead of repetition.

---

## 3. Commands (unchanged from Tier 1)

```
talk <npc_name>    — start conversation (free-form mode)
leave / bye         — end conversation
```

During conversation: all input is treated as free-form text. No numbered choices.
Exception: when in a Tier 1 dialog node (via `set_current_node_id`), numbered
choices appear.

---

## 4. Intent Authoring — Admin UI

### 4.1 Intents Tab on NPC Template Detail

```
NPC: "Elder Myrddin"
  [Dialog Tree] | [Intents]

Intents (free-form conversation mode)

┌────────────────────────────────────────────────────────┐
│ [greeting] priority: 10                                │
│ Triggers: ["hello", "hi", ""]                          │
│ Response: "Welcome, traveler. These are dark times."   │
│ One-shot: ✓  Cooldown: 0                               │
│ Tags req: []  Tags forbid: []                          │
│ Effects: []  Fallback node: [none]                     │
│ [✎ Edit] [🗑 Delete]                                   │
├────────────────────────────────────────────────────────┤
│ [query_caves] priority: 5                              │
│ Triggers: ["cave", "north", "dark", "temple"]          │
│ Response: "The caves to the north have grown dark..."  │
│ One-shot: ✗  Cooldown: 3                               │
│ Tags req: []  Tags forbid: []                          │
│ Effects: []  Fallback node: [none]                     │
│ [✎ Edit] [🗑 Delete]                                   │
├────────────────────────────────────────────────────────┤
│ [query_goblins] priority: 5                            │
│ Triggers: ["goblin", "shaman", "wand", "drain"]        │
│ Response: "Three of them, by my count..."              │
│ One-shot: ✗  Cooldown: 2                               │
│ Tags req: []  Tags forbid: []                          │
│ Effects: []  Fallback node: [none]                     │
│ [✎ Edit] [🗑 Delete]                                   │
├────────────────────────────────────────────────────────┤
│ [accept_offer] priority: 8                             │
│ Triggers: ["i'll go", "handle it", "accept"]           │
│ Response: "Then take this. You'll need it."            │
│ One-shot: ✓  Cooldown: 0                               │
│ Tags req: [quest_shadows_available]                    │
│ Effects: []  Fallback node: [node_005_accept]          │
│ [✎ Edit] [🗑 Delete]                                   │
└────────────────────────────────────────────────────────┘
[+ Add Intent]
```

### 4.2 Intent Tester

A "Test Intents" panel where the DM types a message ("what can you tell me about the goblins") and sees which intent matched and why (score, trigger words hit).

---

## 5. Ent Schema Summary

**New schema:** `dialog_intent`

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (unique) | Intent identifier |
| `npc_template_id` | string (FK) | Parent NPC |
| `intent` | string | Topic label |
| `trigger_words` | []string | Words that trigger this |
| `response_text` | string | What NPC says |
| `fallback_text` | string | Default when no advancement |
| `priority` | int | Match tiebreaker |
| `one_shot` | bool | Remove after one use |
| `tags_required` | []string | Gate on character tags |
| `tags_forbidden` | []string | Negation gate |
| `cooldown_messages` | int? | Messages before reusable |
| `is_greeting` | bool | Fire on conversation start |
| `on_match_effects` | []int | Effects applied on match |
| `set_current_node_id` | string? | Drop into Tier 1 node after |

**Modified schema:** `npc_template`

Add edge:
```go
edge.To("dialog_intents", DialogIntent.Type),
```

---

## 6. Transition: Tier 1 ↔ Tier 2

Both modes coexist. The NPC template has a `dialog_mode` field:

```go
field.Enum("dialog_mode").
    Values("tree", "intent", "llm").
    Default("tree")
```

- `tree` → Tier 1 numbered choices
- `intent` → Tier 2 free-form keyword matching
- `llm` → Tier 3 LLM-powered (future)

The `talk` command reads `dialog_mode` and routes accordingly. If an intent
sets `set_current_node_id`, the mode temporarily switches to `tree` for that
node, then returns to `intent` after the node resolves (quest accepted/declined,
choice made).

---

## 7. Acceptance Criteria

- [ ] `dialog_intent.go` ent schema with all fields
- [ ] `NPCTemplate` has `dialog_intents` edge and `dialog_mode` field
- [ ] `ent generate` clean
- [ ] Tokenizer + matcher algorithm in `herbst/dialog/matcher.go`
- [ ] Free-form input mode during `talk` when `dialog_mode = "intent"`
- [ ] Greeting intents fire automatically on conversation start
- [ ] Intent matching with scoring, priority, and cooldowns
- [ ] One-shot intents removed from pool after firing
- [ ] Tag gates (`tags_required`, `tags_forbidden`) enforced
- [ ] `set_current_node_id` transitions to Tier 1 node, then back
- [ ] Generic fallback when no intents match
- [ ] Intents tab with CRUD on NPC template detail in admin UI
- [ ] Intent tester panel in admin UI
- [ ] `npm run build` passes
- [ ] Unit tests for matching algorithm edge cases

---
