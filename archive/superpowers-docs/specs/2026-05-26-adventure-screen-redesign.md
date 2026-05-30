---
title: Adventure Screen Redesign
date: 2026-05-26
---

# Adventure Screen Redesign

## Problem

The adventure (non-combat) screen in the web-client is cluttered with ability buttons in a `HotkeyBar` at the bottom. Abilities are primarily combat actions and have no place on the exploration screen. The room panel sections (exits, characters, items) are stacked vertically, wasting horizontal space and creating visual noise.

## Goal

Clean, two-handed mobile-friendly adventure screen:
- Room description stays prominent across the top
- Room entities arranged in a 2x2 grid for dense but scannable layout
- Abilities removed entirely from non-combat view (combat screen retains them)
- Text input remains at the bottom for MUD commands (`/tell`, `look`, etc.)

## Design

### Layout: 2x2 Quadrant Room Panel

```
┌─────────────────────────────────────────┐
│  Grand Hall                             │  ← Room title
│  A vast hall with stone pillars.        │  ← Description
├───────────────┬─────────────────────────┤
│  CHARACTERS   │  EXITS                  │
│  Guard ▾      │  N North Passage        │
│  Merchant ▾   │  S Courtyard            │
├───────────────┼─────────────────────────┤
│  ITEMS        │  FUNCTIONS              │
│  Iron Sword ▾ │  L look    E examine    │
│  Potion ▾     │                         │
└───────────────┴─────────────────────────┘
```

**Quadrants:**
- **Top-left**: Characters (hostile + non-hostile, expandable chips)
- **Top-right**: Exits (direction initial + label, tap-to-move)
- **Bottom-left**: Items (takeable + non-takeable, expandable chips)
- **Bottom-right**: Functions (`look`, `examine`, and future non-combat actions)

### Abilities Removal

- `HotkeyBar` component is **only rendered during combat** (`inCombat === true`)
- Keyboard shortcuts 1-4 still work during combat (handled in `GameScreen` keydown handler)
- Abilities remain visible/equippable in the `CharacterPanel` side panel
- No ability UI on the adventure screen whatsoever

### Files to Change

| File | Change |
|------|--------|
| `web-client/src/components/RoomScreen.tsx` | Restructure from stacked sections to CSS Grid 2x2. Add `Functions` quadrant with `look`/`examine` buttons. Keep room title/description spanning full width above grid. |
| `web-client/src/components/GameScreen.tsx` | Move `HotkeyBar` inside the `inCombat` branch so it only renders during combat. Remove unconditional `<HotkeyBar>` at bottom. |
| `web-client/src/components/HotkeyBar.tsx` | Optionally: remove static `look`/`examine`/`use potion` slots since those move to Functions quadrant (or keep them for combat-only display). |

### Styling Details

- Grid: `grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; gap: 1px;` with `bg-border` as grid lines
- Each quadrant: `bg-surface`, `p-2`, scrollable if content overflows
- Quadrant labels: `text-[10px] text-muted uppercase tracking-wider`
- Functions quadrant buttons match existing exit button style (border, rounded, hover)
- Mobile: grid collapses to single column or stays 2x2 depending on viewport width. On very narrow screens (`< 400px`), consider `grid-template-columns: 1fr`.

### Keyboard Shortcuts (Unchanged)

| Key | Action | Context |
|-----|--------|---------|
| `L` | `look` | Non-combat |
| `E` | `examine` | Non-combat |
| `1-4` | Ability slot | Combat only |
| `5` / `R` | Use potion | Combat only |
| `F` | Flee | Combat only |

### Data Flow

No backend changes. Purely frontend layout refactor:
1. `useMUDSocket` still emits `roomScreen` payloads
2. `RoomScreen` receives same props, just renders differently
3. `GameScreen` conditionally renders `HotkeyBar` based on `inCombat`
4. No API contract changes

### Out of Scope

- Combat screen layout (already has its own `CombatActionBar`)
- CharacterPanel abilities tab (unchanged)
- Adding new non-combat abilities/functions (future work)
- Scrollback behavior (unchanged)
