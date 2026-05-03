# Sidebar Collapse Design — Icons-Only Mode (#314)

## Overview

The sidebar transitions between an **expanded** state (220px, icons + labels) and a **collapsed** state (64px, icons only) via a toggle button, similar to VS Code's sidebar behavior.

## State Management

- **State**: `collapsed: boolean` — stored in React `useState`, initialized from `localStorage`
- **Persistence key**: `sidebar-collapsed` in `localStorage`
- **Default**: `false` (expanded) — localStorage read is wrapped in try/catch for embedded browser contexts
- **Effect**: A `useEffect` syncs state changes back to localStorage on every toggle

## Width Strategy

Width is controlled by **conditional Tailwind classes**, NOT inline `style={}`:
- Expanded: `w-[220px]` with `max-w-[220px]`
- Collapsed: `w-[64px]` with `max-w-[64px]`
- Transition: `transition-all duration-300 ease-in-out`

This avoids the specificity problem where a global `.sidebar` CSS rule with `min-width: 220px` would override class-based width.

## Collapsed Behavior

| Property | Expanded | Collapsed |
|---|---|---|
| Width | 220px | 64px |
| Header title | visible, opacity-100 | hidden, opacity-0, select-none |
| Nav item labels | visible | hidden (opacity-0, pointer-events-none, w-0, overflow-hidden) |
| Nav item layout | `gap-3 px-3` | `justify-center px-0` |
| Toggle icon | ChevronLeft (‹) | ChevronRight (›) |
| Tooltip on nav items | none | `title={item.label}` for hover hint |

## Collapse Toggle Button

- **Component**: `SidebarCollapseToggle` — named function (not inline), for DevTools readability
- **Position**: Inside the header row, right-aligned
- **Icons**: Uses `ChevronLeftIcon` / `ChevronRightIcon` from `./icons/ChevronIcons`
- **Styling via Button component**: `variant="ghost"`, `size="sm"`, explicit `style={{ color: '#646cff' }}` AND `stroke="#646cff"` on SVGs (avoids specificity battle with global button reset)
- **Accessibility**: `aria-label` toggles between "Expand sidebar" / "Collapse sidebar"

## Main Content Area

The `<main>` in `__root.tsx` uses `flex-1` which automatically expands when the sidebar shrinks — no layout changes needed. The sidebar is inside a flex row, so width transitions on the sidebar element naturally push/pull the main content.

## CSS Rules (from collapsible-sidebar-pattern.md)

1. No `overflow-hidden` on the `<nav>` — it creates a collapsed scrollbar. Use `overflow-y-auto` on the inner nav items div only.
2. No `flex-shrink-0` on Link elements — the flex column must be allowed to compress.
3. `min-w-0` on the header text div cancels inherited min-width.
4. `whitespace-nowrap block overflow-hidden text-ellipsis` on the title prevents wrapping.
5. `pointer-events-none` as Tailwind class, not inline style.
6. Global `.sidebar { min-width: 220px }` must be removed if present (currently not present in this project).

## Implementation Plan

1. Update `Sidebar.tsx` — replace text-based toggle (`'›'`/`'‹'`) with SVG icons from `ChevronIcons.tsx`, update toggle styling per pattern doc rules
2. Update `Sidebar.test.tsx` — add tests for collapse toggle, localStorage persistence, collapsed rendering
3. Verify no `min-width` global CSS rules block the transition
4. Commit with message: `feat(UI): collapsible sidebar — icons-only mode (#314)`