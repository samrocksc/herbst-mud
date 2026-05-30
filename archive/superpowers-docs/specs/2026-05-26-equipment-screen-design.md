# Equipment Screen — Web Client Design Spec

## Overview

A full-screen equipment overlay accessible by clicking the character name in the top header of `GameScreen`. Displays all equipment slots available to the character's race and allows equip/unequip actions from inventory.

## Goals

- Show character equipment at a glance
- Allow equip/unequip without typing MUD commands
- Respect race-specific slot availability
- Stay consistent with existing web-client UI (Tailwind, font-mono, MUD theme)

## Non-Goals

- Drag-and-drop (click-based selection only)
- Equipment comparison / stat diff previews
- Tooltip hover details (out of scope)

## Architecture

### Components

| Component | File | Responsibility |
|-----------|------|--------------|
| `EquipmentScreen` | `components/EquipmentScreen.tsx` | Full-screen overlay, slot grid, item list |
| `GameScreen` (modified) | `components/GameScreen.tsx` | Toggle `equipmentOpen` state on name click |

### State

- `equipmentOpen: boolean` — lives in `GameScreen`, passed to overlay
- `selectedSlot: string \| null` — lives in `EquipmentScreen`
- `equippedItems: EquipmentItem[]`, `inventoryItems: EquipmentItem[]` — fetched from API
- `error: string \| null` — local API error banner

### API Surface (web-client)

```typescript
// lib/api.ts additions
type EquipmentItem = {
  id: number;
  name: string;
  description: string;
  slot: string;
  level: number;
  weight: number;
  isEquipped: boolean;
  rarity: string;
  color: string;
  armor_rating: number;
  damage_dice_count: number;
  damage_dice_sides: number;
  damage_bonus: number;
  ownerId: number | null;
};

export async function getCharacterEquipment(charID: number): Promise<readonly EquipmentItem[]>;
export async function equipItem(itemID: number, charID: number): Promise<void>;
export async function unequipItem(itemID: number, charID: number): Promise<void>;
export async function getRaceByName(raceName: string): Promise<Race | null>;
```

### Data Flow

1. User clicks character name → `setEquipmentOpen(true)`
2. `EquipmentScreen` mounts → parallel fetch:
   - `GET /equipment?ownerId={charID}` → split into `equipped` and `unequipped`
   - `GET /races` → find race, extract `equipment_slots`
3. Render grid: one card per slot from race definition
4. User clicks slot → `selectedSlot = slotName`
5. Expand list below card: filter `unequipped` where `item.slot === selectedSlot`
6. User clicks item → `PUT /equipment/{id}/equip` → refetch all items → re-render
7. User clicks ✕ on equipped item → `PUT /equipment/{id}/unequip` → refetch → re-render
8. Press ESC or click Close → `setEquipmentOpen(false)`

## Layout

### Overlay Container

- `fixed inset-0 z-50 bg-background flex flex-col`
- Header bar with character name, level, race, class, close button
- Body: scrollable slot grid + expanded item list area

### Slot Grid

- `grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 p-4`
- Each card:
  - Border `border-border`, rounded, padding
  - Slot label: `text-[10px] uppercase tracking-wider text-muted`
  - Equipped item name (or "Empty" italic muted)
  - Rarity color dot or badge (if equipped)
  - Selected state: `border-accent bg-accent/10`

### Expanded Item List

- Appears immediately below the selected slot card (using grid placement or a full-width sub-panel)
- Each row: item name, level, weight, rarity dot. Click to equip.
- "No items available" if filter returns empty.

## Race Slot Fallback

If `/races` fails or the race has no `equipment_slots`, fall back to:
```
["head", "body", "hands", "legs", "feet", "main_hand"]
```
This ensures the screen never breaks.

## Error Handling

- API errors render as an inline banner at the top of the overlay: `bg-danger/20 text-danger text-xs px-3 py-2`
- Equip/unequip failures show the error inline and do not close the overlay.
- **All server-side errors (5xx, 4xx, network failures) are logged to `console.error` with structured context `{ endpoint, status, body, characterID }` so they can be forwarded to the admin for triage.**

## Keyboard

- `Escape` closes the overlay (handled in `EquipmentScreen` useEffect)
- No other keyboard shortcuts (out of scope)

## Server Changes

**None required.** `GET /equipment?ownerId={charID}` already exists. We filter `isEquipped` client-side to avoid backend changes. If performance becomes an issue with large inventories, we can add a server-side `isEquipped` filter later.

## Testing (Manual)

1. Open overlay via character name click
2. Verify all race slots appear
3. Click a slot with no matching inventory → see "No items available"
4. Click a slot with matching items → equip one
5. Verify item appears in slot, disappears from inventory list
6. Click ✕ → item unequips, returns to inventory list
7. Press ESC → overlay closes
8. Disconnect race endpoint → overlay still renders with fallback slots
