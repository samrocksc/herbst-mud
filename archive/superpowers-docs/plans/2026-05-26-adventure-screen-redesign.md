# Adventure Screen Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restructure the web-client adventure screen into a 2x2 quadrant room panel and remove abilities from non-combat view.

**Architecture:** Pure frontend layout refactor. No backend changes. `RoomScreen` gets a CSS Grid 2x2 layout with a new `Functions` quadrant. `GameScreen` conditionally renders `HotkeyBar` only during combat.

**Tech Stack:** React, Tailwind CSS, TypeScript

---

### Task 1: Restructure RoomScreen into 2x2 Quadrant Grid

**Files:**
- Modify: `web-client/src/components/RoomScreen.tsx`

- [ ] **Step 1: Read current RoomScreen.tsx**

- [ ] **Step 2: Rewrite the layout to 2x2 grid**

Replace the stacked sections (exits, characters, items) with a CSS Grid layout. Room title/description stays across the top. Add a new `Functions` quadrant with `look` and `examine` buttons.

```tsx
import { useState } from "react";
import { type RoomScreenPayload } from "../lib/types";
import { ActionExpand, type EntityAction } from "../ui";

type Props = {
  room: RoomScreenPayload;
  onTapExit: (exit: { direction: string; label: string }) => void;
  onCommand: (cmd: string) => void;
  expandedId: string | null;
  onToggleExpand: (id: string) => void;
  pendingTargets: Set<number>;
  onTogglePending: (id: number) => void;
  onConfirmAttack: (char: RoomScreenPayload["characters"][number]) => void;
};

export default function RoomScreen({
  room,
  onTapExit,
  onCommand,
  expandedId,
  onToggleExpand,
  pendingTargets,
  onTogglePending,
  onConfirmAttack,
}: Props) {
  const [descHidden, setDescHidden] = useState(false);

  const characterActions = (char: RoomScreenPayload["characters"][number]): EntityAction[] => {
    const actions: EntityAction[] = [];
    if (char.hostile) {
      actions.push({ label: "Attack", variant: "danger", onClick: () => onConfirmAttack(char) });
    } else {
      if (pendingTargets.has(char.id)) {
        actions.push({ label: "Confirm", variant: "danger", onClick: () => onConfirmAttack(char) });
      } else {
        actions.push({ label: "Attack", variant: "danger", onClick: () => onTogglePending(char.id) });
      }
      actions.push({ label: "Talk", variant: "success", onClick: () => onCommand(`talk ${char.name}`) });
    }
    actions.push({ label: "Examine", variant: "secondary", onClick: () => onCommand(`examine ${char.name}`) });
    return actions;
  };

  const itemActions = (item: RoomScreenPayload["items"][number]): EntityAction[] => {
    const actions: EntityAction[] = [
      { label: "Examine", variant: "secondary", onClick: () => onCommand(`examine ${item.name}`) },
    ];
    if (item.takeable) {
      actions.push({ label: "Take", variant: "primary", onClick: () => onCommand(`take ${item.name}`) });
    }
    return actions;
  };

  return (
    <div className="shrink-0 bg-surface border-t border-border">
      {/* Room title + description — full width */}
      <div className="px-3 py-2 border-b border-border">
        <div className="flex items-center justify-between">
          <h2 className="font-bold text-sm text-accent">{room.title}</h2>
          <button
            type="button"
            onClick={() => setDescHidden((v) => !v)}
            className="text-[10px] text-muted hover:text-foreground px-1"
            title={descHidden ? "Show description" : "Hide description"}
          >
            {descHidden ? "[+]" : "[-]"}
          </button>
        </div>
        {!descHidden && <p className="text-[11px] text-muted leading-relaxed">{room.description}</p>}
      </div>

      {/* 2x2 Quadrant Grid */}
      <div className="grid grid-cols-2 grid-rows-2 gap-px bg-border">
        {/* Top Left: Characters */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "140px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Characters</div>
          {room.characters.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.characters.map((char) => (
                <ActionExpand
                  key={char.id}
                  isOpen={expandedId === `char-${char.id}`}
                  onToggle={() => onToggleExpand(`char-${char.id}`)}
                  label={char.name}
                  triggerClassName={char.hostile ? "text-danger" : "text-foreground"}
                  actions={characterActions(char)}
                />
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Top Right: Exits */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "140px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Exits</div>
          {room.exits.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.exits.map((exit) => (
                <button
                  key={exit.direction}
                  onClick={() => onTapExit(exit)}
                  className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
                >
                  <span className="font-bold text-accent">{exit.direction.charAt(0).toUpperCase()}</span>
                  <span className="text-muted">{exit.label}</span>
                </button>
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Bottom Left: Items */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "140px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Items</div>
          {room.items.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.items.map((item) => (
                <ActionExpand
                  key={item.id}
                  isOpen={expandedId === `item-${item.id}`}
                  onToggle={() => onToggleExpand(`item-${item.id}`)}
                  label={item.name}
                  triggerClassName={item.takeable ? "text-success" : "text-muted"}
                  actions={itemActions(item)}
                />
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Bottom Right: Functions */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "140px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Functions</div>
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => onCommand("look")}
              className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
            >
              <span className="font-bold text-accent">L</span>
              <span className="text-muted">look</span>
            </button>
            <button
              onClick={() => onCommand("examine")}
              className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
            >
              <span className="font-bold text-accent">E</span>
              <span className="text-muted">examine</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Verify no type errors**

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add web-client/src/components/RoomScreen.tsx
git commit -m "🟣 feat(web-client): restructure RoomScreen into 2x2 quadrant grid with Functions"
```

---

### Task 2: Move HotkeyBar into Combat-Only Branch

**Files:**
- Modify: `web-client/src/components/GameScreen.tsx`

- [ ] **Step 1: Read current GameScreen.tsx around line 519**

- [ ] **Step 2: Remove unconditional HotkeyBar, ensure it's only inside inCombat**

Find the unconditional HotkeyBar at the bottom of GameScreen (~line 519):
```tsx
{!inCombat && (
  <HotkeyBar onActivate={handleHotkey} skills={skills} />
)}
```

Delete this block entirely. The combat branch already has `CombatScreen` which includes its own `CombatActionBar`. There is no need for `HotkeyBar` outside combat.

Also remove the `handleHotkey` callback if it becomes unused. Check if `skills` and `HOTKEY_BINDINGS` are referenced elsewhere in the file before removing.

- [ ] **Step 3: Verify no type errors**

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add web-client/src/components/GameScreen.tsx
git commit -m "🟣 feat(web-client): remove HotkeyBar from non-combat adventure screen"
```

---

### Task 3: Verify Build and Visual Check

- [ ] **Step 1: Build the web-client**

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npm run build`
Expected: Build succeeds with no errors

- [ ] **Step 2: Start the dev server and visually verify**

Run: `cd /home/sam/GitHub/herbst-mud && make dev` (or just the web-client dev server)

Open the web client in browser. Check:
- Adventure screen shows 2x2 grid (Characters, Exits, Items, Functions)
- Room title/description spans full width above grid
- `look` and `examine` buttons work in Functions quadrant
- No ability buttons visible on non-combat screen
- Combat screen still shows ability action bar
- Text input remains at the bottom

- [ ] **Step 3: Commit any final fixes**

If any visual tweaks needed (spacing, colors, etc.), commit them.

---

### Self-Review Checklist

| Spec Requirement | Task |
|---|---|
| Room description across the top | Task 1 |
| 2x2 quadrant grid | Task 1 |
| Characters top-left | Task 1 |
| Exits top-right | Task 1 |
| Items bottom-left | Task 1 |
| Functions bottom-right | Task 1 |
| Abilities removed from non-combat | Task 2 |
| Text input stays at bottom | Unchanged (no task needed) |
| No backend changes | Confirmed |

**Placeholder scan:** None found. All code is explicit.

**Type consistency:** `RoomScreenPayload` types used match existing codebase. `onCommand` prop is `string => void` in both old and new.
