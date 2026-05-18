# SSH Client UI Overhaul — Design Spec

## Overview

Polish the SSH MUD client UI with a fantasy-themed color palette, clearer screen flows, better onboarding, and a more organized playing screen. Auth screens get a consistent split-panel layout with styled headers and prompts. The playing screen gets a header bar, improved content area, and cleaner input section.

## Theme & Visual Identity

- **Palette**: Deep purples (borders, backgrounds), rich golds (headers, highlights), bright blues (links, accents, magic)
- **Style**: Modern terminal app — clean panels, organized layouts, clear visual hierarchy
- **Layout**: All auth screens use a consistent split-panel (styled output on top, input prompt on bottom)
- **Typography**: Bold colored headers, consistent spacing, color-coded information types

## Screen Flow

```
Login (account) → World Select → Character Select → Playing
```

### 1. Login Screen

Styled account login with two prompts:
- Username input (masked or plain depending on field)
- Password input (masked with `●` characters)
- Registration option from this screen
- Quit option
- **Onboarding**: Clear instructions visible at all times ("Enter your username", "Type 'register' to create an account")

### 2. World Select Screen

- Lists available worlds from the server
- Numbered selection (1, 2, 3...) and name-based selection
- Shows which world is currently selected
- Back option to return to login
- **Onboarding**: Brief instruction text at the bottom

### 3. Character Select Screen

- Shows characters in the selected world with: name, level, race, class, HP
- Numbered selection and name-based selection
- "Create new character" option
- Back option to return to world select
- **Onboarding**: Instructions for selecting vs creating

### 4. Playing Screen (Game View)

Three-zone layout:

- **Header bar** (top): Room/area name, brief status summary
- **Main content** (middle, scrollable): Room description, exits (color-coded by visited/known/new), items, NPCs/players, message history
- **Status bar** (middle-bottom): Compact HP/Stamina/Mana with progress bars
- **Input area** (bottom): Clean prompt line with command hints

#### Onboarding for new characters:
- Brief welcome message with essential commands on first entry
- Visible command hints: "Try: look, north/south/east/west, say hello, help"

## What Stays the Same

- Account auth model (username + password via REST API)
- Character creation flow (name, password, race, class)
- All game commands and mechanics
- Existing CombatScreen and SkillSelect screens (visual polish only)
- Split-panel rendering approach for auth screens
- The underlying Update/View bubbletea lifecycle

## Files to Modify

- `herbst/ui_screens.go` — Rewrite screen rendering functions with new theme
- `herbst/style.go` — Update color palette to fantasy theme, add new styles
- `herbst/auth.go` — Improve message formatting, add onboarding text
- `herbst/game_model.go` — Restructure playing screen layout (header + content + status + input)
- `herbst/model.go` — Add any new state fields if needed
- `herbst/ui_messages.go` — Update message styling if needed

## Acceptance Criteria

1. Login screen shows styled username/password prompts with masked password input
2. World select screen lists worlds with numbered selection
3. Character select screen shows characters with stats
4. Playing screen has: header bar, room content, status bar, input area
5. Fantasy color palette applied consistently across all screens
6. New characters see onboarding hints on first entry
7. All existing functionality continues to work (combat, movement, etc.)
