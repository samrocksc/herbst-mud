# Combat Screen Design — Web Client

## Overview
Replace the web client's toggle-based combat/adventure modes with a proper **Combat Screen** that mirrors the SSH client's `ScreenCombat`. When in combat, the room description is replaced entirely by a combat panel showing targets, player vitals, action buttons, and a combat log.

## Goals
- Remove combat commands from the adventure screen's HotkeyBar
- Support multi-target selection with a Confirm button before entering combat
- Auto-enter combat when a hostile mob is present in the room
- Run combat rounds on the same 3-second tick interval as the SSH client
- Display health/mana/stamina bars in the combat panel (header shows raw numbers)

## State Machine

```
[ADVENTURE] --attack hostile--> [COMBAT]
[ADVENTURE] --select target + confirm--> [COMBAT]
[COMBAT] --target defeated / flee / player defeat--> [ADVENTURE]
```

Once in combat, the user stays in combat view until combat ends. The Tab key toggles combat/adventure **only when not actively in combat**.

## GameScreen State Additions

```ts
type CombatTarget = {
  id: number;
  name: string;
  hp: number;
  maxHp: number;
  level?: number;
};

type CombatLogEntry = {
  timestamp: number;
  text: string;
  kind: "hit" | "miss" | "crit" | "heal" | "system" | "queue" | "flee";
};

// New GameScreen state:
- inCombat: boolean
- combatTargets: CombatTarget[]
- combatLog: CombatLogEntry[]
- combatRound: number
- queuedAction: string | null
- pendingTargets: Set<number>     // selected but not yet confirmed
- tickTimer: number | null          // setInterval handle
```

## Combat Flow (Client-Side Engine)

### Tick Interval
3 seconds per round (`herbst/combat/config.go:DefaultTickInterval`). Matches the SSH client exactly.

### Each Tick
1. Decrement ability cooldowns
2. Execute queued action (or auto-attack if none queued)
3. Poll target HP via `GET /characters/:id/combat-status`
4. Simulate enemy turn (calculate damage, apply to player HP)
5. Check for combat end (target HP ≤ 0, player HP ≤ 0, or flee success)

### Combat Start
1. Fetch fresh target HP via `/characters/:id/combat-status`
2. Push "Combat started" entry to combat log
3. Start 3-second interval timer
4. Set `inCombat = true`, render CombatScreen

### Combat End
1. Clear tick interval
2. Set `inCombat = false`
3. Fetch fresh room screen from server (`look`)
4. Render RoomScreen

### Damage Application (Client-Side)
Until the backend supports server-driven combat rounds, the web client mirrors the SSH client's dice logic:
- Player attack: d20 + DEX mod vs target AC, weapon damage + STR mod
- Enemy attack: d20 + level/3 vs player AC (base 10 + level/2 + armor)
- Flee: d20 + level/2 vs DC 12

All damage is sent to the server via existing `POST /characters/:id/damage` endpoint.

## CombatScreen Component Layout

Replaces `RoomScreen` when `inCombat === true`.

```
┌─────────────────────────────────┐
│ COMBAT — Round 3                 │  header bar
├─────────────────────────────────┤
│ 🎯 Goblin Scout                  │  targets
│ HP ████████░░░░ 45/60           │
│ 🎯 Goblin Warrior                │
│ HP ██████░░░░░░ 30/60           │
├─────────────────────────────────┤
│ ❤  HP  ████████░░ 80/100        │  player vitals
│ ⚡ STA ██████████  50/50         │
│ 💧 MANA ██████░░░░ 30/50        │
├─────────────────────────────────┤
│ [1] Slash      [2] Stab         │  action buttons
│ [3] —          [4] —             │
│ [5] Potion(3)     [F] Flee      │
├─────────────────────────────────┤
│ [14:32:05] Queued: Slash        │  combat log
│ [14:32:03] You slash Goblin for 12 dmg
│ [14:32:02] Goblin hits you for 5 dmg
└─────────────────────────────────┘
```

### Sub-Components
- `CombatHeader`: round counter + combat status
- `CombatTargetList`: each target shows name, level, HP bar
- `CombatVitals`: player HP/STA/MANA bars
- `CombatActionBar`: 4 ability slots + potion + flee button. Same cooldown overlay as current `CombatHUD`
- `CombatLog`: auto-scrolling timestamped entries

### Stat Bars
- HP bar: red-to-green gradient based on percentage
- STA bar: yellow
- MANA bar: blue
- Only rendered in CombatScreen. Header shows raw numbers always.

## RoomScreen Changes (Attack/Confirm)

### Character Actions
- **Hostile NPC**: "Attack" button immediately enters combat (red, no confirm)
- **Neutral/Friendly NPC**: "Attack" → replaces with red "Confirm" button. Multiple targets can be in confirm state simultaneously. "Confirm" sends `attack <name>` and enters combat.
- **All NPCs**: "Talk" and "Examine" remain

### Implementation
Replace the current `ActionExpand` actions array with stateful logic:

```ts
const [pendingTargets, setPendingTargets] = useState<Set<number>>(new Set());

// In characterActions:
if (char.hostile) {
  actions.push({ label: "Attack", variant: "danger", onClick: () => enterCombat([char]) });
} else {
  if (pendingTargets.has(char.id)) {
    actions.push({ label: "Confirm", variant: "danger", onClick: () => enterCombat([char]) });
  } else {
    actions.push({ label: "Attack", variant: "danger", onClick: () => addPendingTarget(char.id) });
  }
  actions.push({ label: "Talk", variant: "success", onClick: () => onCommand(`talk ${char.name}`) });
}
actions.push({ label: "Examine", variant: "secondary", onClick: () => onCommand(`examine ${char.name}`) });
```

Clicking "Confirm" on multiple targets sends `attack <first>` for now. Backend multi-target support is Phase 2.

## HotkeyBar Cleanup (Adventure Mode Only)

Remove combat commands from adventure-mode HotkeyBar:

| Key | Current | New |
|---|---|---|
| 1 | attack | *(removed)* |
| 2 | flee | *(removed)* |
| 3 | hide | *(removed)* |
| 4 | backstab | *(removed)* |
| 5 | concentrate | use potion |
| L | look | look |
| E | examine | examine |
| R | use potion | use potion |
| Q | quit | *(removed — UI button only)* |

Combat shortcuts (1-4 abilities, 5 potion, F flee) only work when `inCombat === true`.

## Keyboard Handler (GameScreen)

```ts
// Only when NOT in an <input>
"tab"    → toggle combatMode (only if !inCombat)
"1"-"4"  → use equipped ability slot (combat only)
"5"/"r"  → use potion (always)
"f"      → flee (combat only)
"l"      → look / refresh room (adventure only)
"e"      → examine (adventure only)
"i"      → open inventory panel
"s"      → open skills panel
"a"      → open abilities panel
```

## API Endpoints Used

| Endpoint | Method | Purpose |
|---|---|---|
| `/characters/:id/combat-status` | GET | Fetch target HP/MaxHP |
| `/characters/:id/damage` | POST | Apply damage to target |
| `/characters/:id/heal` | POST | Heal target (defeat/respawn) |
| `/characters/:id` | GET | Refresh player stats after combat |

## WebSocket Protocol (Phase 2 — Server-Driven)

When the backend supports server-driven combat rounds, add new `ServerMessage` types:

```ts
type CombatStartMessage = {
  type: "combat_start";
  payload: { targets: CombatTarget[] };
};

type CombatTickMessage = {
  type: "combat_tick";
  payload: { round: number };
};

type CombatEventMessage = {
  type: "combat_event";
  payload: {
    actor: string;
    action: string;
    target?: string;
    damage?: number;
    healing?: number;
  };
};

type CombatEndMessage = {
  type: "combat_end";
  payload: { reason: "victory" | "defeat" | "flee" };
};
```

This is a non-breaking addition. The current `output`/`system`/`screen` flow continues.

## Files Changed

| File | Change |
|---|---|
| `web-client/src/components/GameScreen.tsx` | Add combat state, tick timer, keyboard gating, conditionally render CombatScreen |
| `web-client/src/components/RoomScreen.tsx` | Add pending target state, confirm button logic |
| `web-client/src/components/HotkeyBar.tsx` | Remove combat commands from adventure mode |
| `web-client/src/components/CombatHUD.tsx` | *(kept as-is — used inside CombatScreen)* |
| `web-client/src/components/CombatScreen.tsx` | **NEW** — full combat view component |
| `web-client/src/lib/types.ts` | Add `CombatTarget`, `CombatLogEntry`, `CombatScreenPayload` |
| `web-client/src/lib/api.ts` | Add `getCombatStatus`, `applyDamage`, `healCharacter` helpers |
| `web-client/src/index.css` | Add HP/STA/MANA bar CSS variables |

## Aesthetic Direction
- Dark, high-contrast combat panel with a subtle red border glow
- Action buttons have tactile press-down animations
- Combat log uses monospace timestamps, colored by event type
- Health bars use CSS gradients that shift from green → yellow → red as HP drops

## Open Questions
- Should the combat log persist after combat ends (visible in scrollback)? → Yes, push each entry to the main scrollback as "combat" styled lines.
- Should auto-attack be the default when no action is queued? → Yes, matches SSH client.
- What happens if the user disconnects during combat? → Combat state is lost. On reconnect, the room screen is fresh.
