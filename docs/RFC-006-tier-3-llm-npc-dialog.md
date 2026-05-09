# RFC-006: NPC Dialog System — Tier 3 (LLM-Powered Conversational NPCs)

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Builds on:** RFC-004 (Tier 1 Dialog Trees), RFC-005 (Tier 2 Keyword-Intent Matching)
**Related:** RFC-002 (Effects & Hooks), RFC-003 (Quest System)

---

## 1. Executive Summary

This RFC proposes an **LLM-powered conversational layer** for NPCs — the crown
of the three-tier dialog architecture. Each NPC can be backed by a language
model that generates dynamic, in-character responses to anything the player
types. The LLM does not replace the structured game systems; it **wraps them
in natural language**. When the player says something that triggers a quest
offer, item trade, or effect, the LLM responds with character-appropriate
dialogue AND a structured JSON action that the game engine executes.

Key design: **DM-first, LLM-second.** Dungeon Masters build quests, dialog
nodes, and intents in the admin UI (Tiers 1+2). The LLM enriches the
conversational surface but defers to those structured definitions for game
state changes. The LLM cannot invent quests, give items, or modify character
state — it can only suggest actions that the game engine validates against
the DM's authored definitions.

---

## 2. Architecture

```
Player types message
        │
        ▼
┌───────────────────────────────────────────────────┐
│              Dialog Router (herbst/dialog/)        │
│                                                    │
│  If dialog_mode == "llm" → LLM Pipeline            │
│  If error/timeout      → Tier 2 Intent Matcher     │
│  If dialog_mode != "llm" → Tier 1 or Tier 2        │
└───────────────────┬───────────────────────────────┘
                    │
                    ▼
┌───────────────────────────────────────────────────┐
│              LLM Pipeline                          │
│                                                    │
│  1. Build prompt (NPC persona + state + history)   │
│  2. Send to LLM API                                │
│  3. Parse response (prose + structured action)     │
│  4. Validate action against DM definitions         │
│  5. Execute valid actions; discard invalid ones    │
│  6. Return prose to player                         │
└───────────────────────────────────────────────────┘
```

---

## 3. Prompt Construction

Each LLM call is constructed from:

### 3.1 System Prompt (per NPC)

```yaml
You are an NPC in a text-based multiplayer RPG. Stay in character at all times.

NPC ID: npc_elder_myrddin
Name: Elder Myrddin
Race: Human
Disposition: friendly
Level: 15

## PERSONALITY
{prompt_personality}

## KNOWLEDGE
{prompt_knowledge}

## CONSTRAINTS
- You are a character in a MUD. You do not know you're in a game.
- You only know what an elder human wizard in this world would reasonably know.
- You cannot give the player items, XP, or quests on your own — only through
  the DM-authored systems below.
- Do not break character. Do not mention the game, players, or mechanics.
- Keep responses concise (1-3 sentences unless explaining something complex).

## AVAILABLE ACTIONS
You may include ONE of these actions in your response. If none apply, omit the action.

### Offer Quest
Offer a quest from this NPC's available quests:
{available_quests_json}

### Trade
Initiate a trade with items from this NPC's trades_with list:
{trades_with_json}

### Apply Effect
Apply an effect from this list (DM-authored, validated by game engine):
{available_effects_json}

### No Action
If the player is just making conversation, respond in character with no action.

## RESPONSE FORMAT
Always respond with valid JSON. The "dialogue" field is shown to the player.
The "action" field is optional and executed by the game engine.

{
  "dialogue": "Your in-character response here",
  "internal_thought": "Brief note on why you responded this way (not shown to player)",
  "action": {
    "type": "offer_quest" | "trade" | "apply_effect" | null,
    "payload": { ... }
  }
}
```

### 3.2 NPC-Specific Prompt Fields

New fields on `NPCTemplate`:

```go
field.Text("llm_personality").Optional().
    Comment("LLM personality prompt: 'You are a weary old wizard who has seen too much...'"),
field.Text("llm_knowledge").Optional().
    Comment("LLM knowledge prompt: 'You know about: the northern caves, goblin shamans, the old temple...'"),
field.String("llm_model").Optional().
    Comment("Override model for this NPC: 'claude-sonnet-4-20250514'"),
field.Float("llm_temperature").Default(0.8),
field.Int("llm_max_tokens").Default(256),
```

These are DM-authored in the admin UI. If left blank, the system uses sensible
defaults generated from the NPC's template fields (race, description, quests,
trades_with, etc.).

### 3.3 Context Injection (dynamic, per message)

Injected at the top of each call:

```
## CURRENT GAME STATE
Room: A dimly lit chamber. Stone walls covered in ancient runes.
Time: Evening
Characters present: Gorthak (level 8 warrior)

## PLAYER STATE
Player: Gorthak
Race: Orc
Level: 8
Faction: None
Active quests: none
Completed quests: none
Tags: [new_arrival]

## CONVERSATION HISTORY
Player: "Elder Myrddin? Are you there?"
Elder Myrddin: "Ah, Gorthak. Yes — I'm here. Lost in thought, I'm afraid."
Player: "What's on your mind?"
```

Context is limited to the last 8 exchanges to manage token costs.

---

## 4. Structured Action Validation

The LLM returns a JSON action alongside prose. The **game engine validates**
every action before execution:

```go
func validateAction(npc *NPCTemplate, action *LLMAction, player *Character) error {
    switch action.Type {
    case "offer_quest":
        // Only quests on this NPC's available_quests list
        if !npc.HasQuest(action.Payload.QuestID) {
            return ErrQuestNotAvailable
        }
        // Check prerequisites, repeat mode
        if !questservice.CanAccept(player.ID, action.Payload.QuestID) {
            return ErrQuestCannotAccept
        }
    case "trade":
        // Only items on this NPC's trades_with list
        if !npc.TradesItem(action.Payload.ItemID) {
            return ErrItemNotTraded
        }
    case "apply_effect":
        // Only effects on this NPC's allowed_effects list (new field)
        if !npc.AllowsEffect(action.Payload.EffectID) {
            return ErrEffectNotAllowed
        }
    }
    return nil
}
```

**The LLM cannot invent game actions.** It can only reference IDs from the
DM-authored lists injected into its prompt. If the LLM hallucinates an action
the NPC doesn't have, the game engine silently discards the action and returns
only the dialogue. The player never sees the rejection.

### 3.4 Allowed Effects (new field on NPCTemplate)

```go
field.JSON("allowed_effects", []int{}).Optional().
    Comment("Effect IDs this NPC can apply during LLM conversation")
```

This gates which effects the LLM can trigger. The DM selects specific effects
in the admin UI. If empty, the `apply_effect` action type is not included in
the LLM's available actions list.

---

## 5. Fallback Strategy

```
LLM call
  ├─ Success → return validated prose + action
  ├─ Timeout (5s) → log warning → fall back to Tier 2 intent matching
  ├─ HTTP error (5xx) → log error → fall back to Tier 2
  ├─ Rate limit (429) → queue, retry once → fall back to Tier 2
  └─ Malformed JSON → strip action, return prose only
```

When falling back: the system sends the same player input through the Tier 2
intent matcher. The transition is seamless — the player sees only the NPC's
response. The conversation remains in free-form mode.

In the admin UI, an "LLM Status" indicator per NPC shows: `online`, `degraded`
(fallback active), `disabled`.

---

## 6. Cost Controls

### 6.1 Per-NPC Rate Limiting

```go
type LLMRateLimiter struct {
    MaxRequestsPerMinute int  // default 20
    MaxTokensPerSession  int  // default 4096
    CooldownSeconds      int  // default 0 (no cooldown)
}
```

### 6.2 Session Token Budget

Each conversation tracks cumulative tokens. When the budget is exhausted,
the system falls back to Tier 2:

```
"Elder Myrddin grows tired of talking. He nods at you briefly."
```

### 6.3 Model Selection

| NPC Priority | Recommended Model | Cost Profile |
|-------------|-------------------|--------------|
| Quest-giver NPCs (critical) | Claude Sonnet / GPT-4o-mini | Medium |
| Lore NPCs (important) | GPT-4o-mini / local model | Low |
| Ambient NPCs (town guards, merchants) | Tier 2 only (no LLM) | Zero |

Model is per-NPC configurable via `llm_model`. If unset, uses the server default.

---

## 7. Player Experience

```
> talk myrddin

Elder Myrddin looks up from a dusty tome, his eyes weary but sharp.
"You've returned. The caves have grown darker since you left.
What did you find?"

> goblins. three of them. dead now.
Elder Myrddin: *He exhales slowly, a weight lifting from his shoulders.*
"Three less shadows in this world. I am in your debt, Gorthak.
The chamber they guarded — did you enter it?"

> yes. cursed place. something's still in there though.
Elder Myrddin: "The Heart of Shadow. I feared as much. It cannot be
destroyed by steel alone. You'll need something blessed — and
someone who knows the old rites. That... would be me, I suppose."

[Quest Available: The Heart of Shadow]
1. Accept Quest
2. Decline
```

The quest offer appears as a structured popup **alongside** the LLM's dialogue.
The LLM doesn't present the accept/decline UI — it just creates the narrative
context. The game engine handles the mechanical offer.

---

## 8. Admin UI — LLM Configuration

### 8.1 NPC Template — LLM Tab

```
NPC: "Elder Myrddin"
  [Dialog Tree] | [Intents] | [LLM]

LLM Configuration
┌────────────────────────────────────────────────────────┐
│ LLM Enabled: [✓]                                        │
│                                                          │
│ Model: [claude-sonnet-4-20250514 ▾]                     │
│ Temperature: [0.8    ▸]                                  │
│ Max Tokens:  [256]                                      │
│                                                          │
│ Rate Limit:   [20] requests/min                         │
│ Token Budget: [4096] per session                        │
│                                                          │
│ Personality Prompt:                                      │
│ ┌──────────────────────────────────────────────────┐   │
│ │ You are Elder Myrddin, a weary but wise human    │   │
│ │ wizard. You have watched over this valley for     │   │
│ │ decades. You speak slowly, choosing words with    │   │
│ │ care. You are haunted by past failures but remain │   │
│ │ hopeful. You call the player "child" or by name.  │   │
│ └──────────────────────────────────────────────────┘   │
│                                                          │
│ Knowledge Prompt:                                        │
│ ┌──────────────────────────────────────────────────┐   │
│ │ You know about: the northern caves, the goblin    │   │
│ │ shamans, the Heart of Shadow (an ancient evil     │   │
│ │ artifact), the old temple, the blessed water of   │   │
│ │ the Moonwell, the dragon that once lived here.     │   │
│ │ You do NOT know about: the eastern mountains,      │   │
│ │ the king's court, modern politics.                 │   │
│ └──────────────────────────────────────────────────┘   │
│                                                          │
│ Allowed Effects:                                         │
│ [tag_add: quest_shadows_complete]  [✕]                  │
│ [xp_gain: 50]                       [✕]                  │
│ [+ Add Allowed Effect]                                  │
│                                                          │
│ LLM Status: ● online                                     │
│                                                          │
│ [Test LLM] [Save]                                        │
└────────────────────────────────────────────────────────┘
```

### 8.2 LLM Test Console

A panel that lets the DM simulate a conversation with the LLM-powered NPC:
type a message, see the raw LLM response (prose + action JSON), and how the
game engine validated/executed any action. Critical for prompt tuning.

---

## 9. LLM Provider Abstraction

The LLM backend is pluggable via an interface:

```go
// herbst/dialog/llm.go
type LLMProvider interface {
    Generate(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
    Name() string
    IsAvailable() bool
}

type LLMRequest struct {
    SystemPrompt string
    UserMessage  string
    History      []Message
    MaxTokens    int
    Temperature  float64
}

type LLMResponse struct {
    Dialogue       string
    InternalThought string
    Action         *LLMAction
    TokensUsed     int
}
```

Providers:
- `anthropic` — Claude models via Anthropic API
- `openai` — GPT models via OpenAI API
- `ollama` — Local models via Ollama (for offline play)
- `mock` — Returns deterministic responses for testing

Provider is server-level config (env var `LLM_PROVIDER`). API keys are env vars
(`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`). The `ollama` provider requires no key
and runs locally.

Model is per-NPC (the model name string is passed to the provider's API call).

---

## 10. Safety Guardrails

### 10.1 Content Filtering

The system prompt includes:
```
You are a character in a text RPG suitable for all audiences. Keep content
appropriate. Do not use profanity, graphic violence, or adult themes.
```

Additionally, the response text is run through a simple word filter (DM-configurable
blocklist in game config) before display.

### 10.2 Prompt Injection Prevention

All LLM outputs are treated as untrusted. Player input is never injected into
the system prompt — it goes only into the user message turn. The structured
action JSON is validated server-side before execution.

### 10.3 Conversation Logging

All LLM interactions are logged to `AppLog` (RFC-002 LOGS tickets) for DM review
and debugging. Logs include: NPC ID, player ID, prompt sent, raw response, action
validation result, tokens used.

---

## 11. Implementation Sequence

| Phase | What | Depends On |
|-------|------|------------|
| Phase 1 | `llm_personality`, `llm_knowledge`, `llm_model`, `llm_temperature`, `llm_max_tokens`, `allowed_effects` fields on `NPCTemplate` | RFC-004 |
| Phase 2 | `LLMProvider` interface + `anthropic` implementation | Phase 1 |
| Phase 3 | Prompt builder (system prompt + context injection) | Phase 2 |
| Phase 4 | Response parser (prose + JSON action extraction) | Phase 3 |
| Phase 5 | Action validator (quests, trades, effects) | Phase 4 + RFC-003 |
| Phase 6 | Fallback to Tier 2 on error/timeout | Phase 5 + RFC-005 |
| Phase 7 | Rate limiter + token budget | Phase 6 |
| Phase 8 | Content filter + logging | Phase 7 |
| Phase 9 | Admin UI: LLM tab + test console | Phase 8 |
| Phase 10 | `ollama` provider (local models) | Phase 2 |
| Phase 11 | `openai` provider | Phase 2 |
| Phase 12 | Per-NPC LLM status indicator | Phase 6 |

---

## 12. Acceptance Criteria

### Core
- [ ] `NPCTemplate` extended with all LLM fields
- [ ] `ent generate` clean
- [ ] `LLMProvider` interface with `anthropic` implementation
- [ ] Prompt builder combines: system prompt, NPC personality, NPC knowledge, game context, player state, conversation history
- [ ] LLM responses parsed into {dialogue, action}
- [ ] Action validation gates all game state changes
- [ ] Invalid/hallucinated actions silently discarded
- [ ] Fallback to Tier 2 intent matching on any LLM error
- [ ] Rate limiter per NPC
- [ ] Session token budget enforcement

### Admin UI
- [ ] LLM tab on NPC template detail with all config fields
- [ ] Personality/knowledge prompt editors (textareas)
- [ ] Allowed effects multi-select
- [ ] LLM test console with raw response + validation display
- [ ] Per-NPC LLM status: online / degraded / disabled

### Quality
- [ ] Content filter applied to all LLM output
- [ ] All LLM interactions logged to AppLog
- [ ] Unit tests for prompt builder, action validator, fallback logic
- [ ] Integration test: full conversation flow with mock LLM provider
- [ ] `npm run build` passes
