# RFC-TOOLTIPS: Admin Tooltips & Documentation System

**Status:** Draft  
**Author:** Leonardo (Hermes Agent)  
**Date:** 2026-05-07  
**Scope:** Admin UI (`admin/src/`) only — backend endpoints already exist  
**Depends on:** None (self-contained frontend feature)

---

## 1. Problem

The Abilities management page exposes 20+ form fields (`mana_cost`, `scaling_percent_per_point`, `proc_chance`, etc.) with zero explanation. A game designer must already know the combat engine to understand what these do. The same problem exists on Skills, NPCs, Items, and Config pages — cryptic fields everywhere with no in-app guidance.

Additionally, there is no central place in the admin UI that explains the game's mechanics. A designer shouldn't have to read Go source files to understand what "Effect Type: haymaker" means.

## 2. Goals

1. **Inline tooltips** on every cryptic form field across all management pages — explain what it does in game mechanics terms
2. **Reusable Tooltip component** — CSS-only, no external tooltip library, supports hover and focus, accessible
3. **Docs section** in sidebar — dedicated documentation pages for game systems (abilities, combat, factions, etc.)
4. **Easy to extend** — adding a tooltip to a new field is one prop, adding a Doc page is one route + markdown file

## 3. Non-Goals

- Do NOT change any backend API or database schema
- Do NOT add tooltips to the SSH/TUI client — admin web UI only
- Do NOT use a heavy tooltip library (no `@floating-ui/react`, no `react-tooltip`) — pure CSS for control and zero deps

---

## 4. Game Mechanics Research

### 4.1 Abilities (aka "Skills" in the DB)

Abilities are **combat actions** — things a player presses `1-5` to activate during combat. They are stored in the `skills` table (unfortunate name collision with trainable skills). The admin UI calls them "Abilities" to reduce confusion.

| Field | What it is | Game Mechanics Meaning |
|-------|-----------|----------------------|
| `name` | Display name | The command name the player types, e.g. `concentrate` |
| `description` | Flavor text | Shown to player when they `examine` or `help` the ability |
| `skill_type` | Category | `combat` = direct attack, `magic` = spell-based, `utility` = non-combat, `healing` = restore HP, `support` = buff allies |
| `cost` | Generic cost | Legacy field — some abilities use this as a flat energy cost before mana/stamina existed. Prefer `mana_cost`/`stamina_cost` |
| `cooldown` | Legacy cooldown | Old tick-based cooldown. Prefer `cooldown_seconds` for real-time |
| `cooldown_seconds` | Real cooldown | Seconds before ability can be used again. Back-off = 10s, Haymaker = 6s |
| `requirements` | Level gate | Minimum character level to learn this ability. e.g. `5` = level 5+ |
| `effect_type` | What happens | `damage` = subtract HP, `heal` = restore HP, `buff` = raise stat temporarily, `debuff` = lower stat, `dot` = damage over time (per tick), `hot` = heal over time, `concentrate` = +WIS to hit rolls, `haymaker` = +STR dmg / -DEX hit, `scream` = trade INT/WIS for DEX/STR, `slap` = stun check, `backoff` = guaranteed dodge |
| `effect_value` | Magnitude | How much damage/heal/stat change. Damage scales with `scaling_stat` if set |
| `effect_duration` | Ticks | How many combat ticks the effect lasts. 1 = one round, 4 = Concentrate's full duration |
| `mana_cost` | MP drain | Subtracted from character's mana pool on activation. Mages use this |
| `stamina_cost` | SP drain | Subtracted from stamina pool. Fighters use this. Haymaker = 15 |
| `hp_cost` | Self-damage | Some abilities cost HP to cast (blood magic, berserker skills) |
| `scaling_stat` | Which stat scales | `STR` = strength, `DEX` = dexterity, `INT` = intelligence, `WIS` = wisdom, `CON` = constitution. The `effect_value` is multiplied by `scaling_percent_per_point` × stat value |
| `scaling_percent_per_point` | Scaling factor | 0.05 = 5% of stat per point. If STR=10, effect_value=50, scaling=0.05 → final = 50 + (10 × 0.05 × 50) = 75 |
| `proc_chance` | Trigger chance | 0.0–1.0 probability that the ability triggers its effect on the `proc_event`. 0.3 = 30% chance |
| `proc_event` | When it triggers | `on_hit` = when attack lands, `on_crit` = critical hit, `on_dodge` = when dodging, `on_kill` = on enemy death |
| `skill_class` | Usage mode | `active` = press button to use, `passive` = always active (no button), `toggle` = on/off state |
| `required_tag` | Prerequisite tag | Comma-separated item tags. e.g. `sword,blade` means character must have an item with the `sword` or `blade` tag equipped to use this ability |
| `slug` | URL-safe name | Used in API routes and YAML references. e.g. `haymaker`, `concentrate` |

### 4.2 Trainable Skills (the `/skills` page)

These are **not** combat abilities. They are proficiency multipliers like "Blades", "Fire Magic", "Pizza Making". Characters train them over time to improve effectiveness with related abilities.

| Field | What it is | Game Mechanics Meaning |
|-------|-----------|----------------------|
| `name` | Display name | e.g. "One-Handed", "Fire Magic", "Sword Fighting" |
| `description` | Flavor text | Shown to player |
| `skill_category` | Grouping | `blades`, `knives`, `fire_magic`, etc. Determines which abilities get bonuses when this skill is high |
| `requirements` | Prerequisites | Format: `level:5,str:10` — character must be level 5 AND have 10 STR to train this |

### 4.3 NPCs

| Field | What it is | Game Mechanics Meaning |
|-------|-----------|----------------------|
| `level` | NPC power | Determines HP, damage, XP value. Used in combat calculations |
| `xp_value` | XP reward | How much XP a player gains when defeating this NPC |
| `respawn_cooldown` | Seconds | How long after death before NPC reappears in a respawn room |
| `respawn_rooms` | Room IDs | Comma-separated. NPC will randomly appear in one of these rooms after cooldown |
| `race` | NPC type | Affects base stats, resistances, available abilities. References the `races` table |

### 4.4 Items

| Field | What it is | Game Mechanics Meaning |
|-------|-----------|----------------------|
| `damage` | Base damage | Added to attack rolls when item is equipped |
| `armor` | Defense bonus | Subtracted from incoming damage. Negative = cursed item |
| `slot` | Equipment slot | `head`, `body`, `hands`, `legs`, `feet`, `main_hand`, `off_hand`, `both_hands` |
| `tags` | Item keywords | Used by abilities (`required_tag`) and for filtering. e.g. `sword,metal,sharp` |

### 4.5 Factions

| Field | What it is | Game Mechanics Meaning |
|-------|-----------|----------------------|
| `standing` | Relationship | -100 (hated) to +100 (revered). Affects NPC aggression, prices, quest availability |
| `category` | Group type | `political`, `religious`, `criminal`, `guild`, `race` |

---

## 5. Tooltip Design

### 5.1 Component: `Tooltip`

```tsx
// components/Tooltip.tsx
// CSS-only tooltip. Wraps any element. Shows on hover and focus.
// No JS positioning library — uses CSS transform + absolute positioning.

type TooltipProps = Readonly<{
  children: ReactNode        // The trigger element (label, icon, badge)
  content: string             // Tooltip text (plain string, keep under 120 chars)
  placement?: 'top' | 'bottom' | 'left' | 'right'  // default: 'top'
}>
```

**Requirements:**
- Pure CSS positioning — `position: relative` on wrapper, `position: absolute` on tooltip text
- Appears on `:hover` and `:focus-within` (accessible for keyboard users)
- Max width 240px, wraps text, `z-index: 50`
- Arrow indicator using CSS border trick (optional — a small triangle pointing to trigger)
- Dark theme: `bg-surface-muted border border-border text-text text-xs p-2 rounded shadow-lg`
- Delay: 300ms before appearing (CSS `transition-delay`)
- Does NOT use `title` attribute (avoids native browser tooltip conflict)

**Accessibility:**
- Tooltip text lives in a `<span role="tooltip">` with `id`
- Trigger element gets `aria-describedby={tooltipId}`
- Screen readers announce the description

### 5.2 Component: `TooltipIcon`

A small `?` or `ⓘ` circle icon that wraps `Tooltip` for compact use on form labels:

```tsx
// components/TooltipIcon.tsx
// Usage: <FormField label={<><span>Name</span> <TooltipIcon content="The command name..." /></>} ... />
```

Style: `w-4 h-4 rounded-full bg-primary/20 text-primary text-[10px] inline-flex items-center justify-center cursor-help`

### 5.3 Extended FormFields

Each field component (`FormField`, `NumberField`, `SelectField`, `TextareaField`, `TagInput`) gets an optional `tooltip?: string` prop. When present, the label renders with a `TooltipIcon`.

```tsx
<NumberField
  label="Mana Cost"
  tooltip="MP drained on activation. Mages rely on this pool. If mana < cost, ability fails."
  value={formData.mana_cost}
  onChange={(v) => set({ mana_cost: v })}
/>
```

**Implementation:** Modify `FieldLabel` to accept an optional `tooltip` prop and render the icon when present. No breaking changes — `tooltip` is optional.

---

## 6. Tooltip Placement Plan

### 6.1 Abilities Page (`abilities.tsx`) — Priority 1

Every form field gets a tooltip:

| Field | Tooltip Content |
|-------|----------------|
| Name | The command name players type, e.g. `concentrate`. Also used as the API slug if slug is empty. |
| Description | Flavor text shown to players when they examine or help this ability. |
| Skill Type | `combat`=direct attack, `magic`=spell, `utility`=non-combat, `healing`=restore HP, `support`=buff allies |
| Required Tag | Comma-separated item tags. Character must have an item with this tag equipped to use this ability. e.g. `sword,blade` |
| Level Req | Minimum character level to learn this ability. |
| Cost | Legacy flat energy cost. Prefer mana_cost/stamina_cost for new abilities. |
| Cooldown (s) | Seconds before ability can be reused. Back-off=10s, Haymaker=6s. |
| Effect Type | What happens on use: damage/heal/buff/debuff/dot/hot, or special: concentrate/haymaker/scream/slap/backoff |
| Effect Value | Base magnitude. Scales with `scaling_stat` if set. |
| Effect Duration | Combat ticks the effect lasts. 1=one round, 4=Concentrate's full duration. |
| Mana Cost | MP drained on activation. If mana < cost, ability fails. |
| Stamina Cost | SP drained on activation. Fighters use this resource. |
| HP Cost | Self-damage to cast. Berserker and blood magic abilities. |
| Scaling Stat | Which character stat boosts the effect. STR=damage, WIS=healing, DEX=dodge, etc. |
| Scaling %/point | Percentage of stat value added per point. 0.05=5%. Formula: final = base + (stat × pct × base) |
| Proc Chance | 0.0–1.0 chance the effect triggers on the proc_event. 0.3=30%. |
| Proc Event | When to roll proc_chance: on_hit=attack lands, on_crit=critical, on_dodge=dodging, on_kill=enemy death |
| Skill Class | `active`=press button, `passive`=always on, `toggle`=on/off switch |

### 6.2 Skills Page (`skills.tsx`) — Priority 2

| Field | Tooltip Content |
|-------|----------------|
| Name | Trainable skill name. e.g. "Blades", "Fire Magic", "Pizza Making" |
| Description | Flavor text shown to players. |
| Category | Determines which abilities get bonuses. Blades → sword abilities, Fire Magic → fire spells |
| Requirements | Format: `level:5,str:10`. Character must meet ALL conditions to train this skill. |

### 6.3 NPCs Page (`npcs.tsx`) — Priority 3

| Field | Tooltip Content |
|-------|----------------|
| Level | NPC power rating. Determines HP pool, damage output, and combat difficulty. |
| XP Value | Experience points awarded when player defeats this NPC. |
| Respawn Cooldown | Seconds after death before NPC reappears in a respawn room. |
| Respawn Rooms | Comma-separated room IDs. NPC randomly picks one after cooldown. |
| Race | Affects base stats and resistances. References the Races table. |

### 6.4 Items Page — Priority 3

| Field | Tooltip Content |
|-------|----------------|
| Damage | Added to attack rolls. Only applies when equipped in main_hand or both_hands. |
| Armor | Subtracted from incoming damage. Negative values = cursed item that hurts wearer. |
| Slot | Where item equips. `both_hands` prevents off_hand item. `off_hand` usually shields. |
| Tags | Keywords used by abilities (required_tag) and for filtering. e.g. `sword,metal,sharp` |

### 6.5 Factions Page — Priority 3

| Field | Tooltip Content |
|-------|----------------|
| Standing | -100 (kill on sight) to +100 (revered). Affects prices, quest access, NPC aggression. |
| Category | `political`=city states, `religious`=churches, `criminal`=thieves guilds, `guild`=craft, `race`=species |

### 6.6 Dashboard Cards — Priority 4

Add `title` attribute (or lightweight tooltip) to dashboard stat cards showing the count and the tool card grid.

---

## 7. Docs Section

### 7.1 Sidebar Entry

Add a "Docs" nav item to the sidebar, below "NPCs" or as a collapsible section. Clicking it shows a list of doc pages.

```tsx
// In Sidebar.tsx — new navItems entry
{ label: 'Docs', path: '/docs', Icon: DocsIcon, children: [
  { label: 'Game Mechanics', path: '/docs/game-mechanics' },
  { label: 'Ability System', path: '/docs/ability-system' },
  { label: 'Combat Guide', path: '/docs/combat-guide' },
  { label: 'Trainable Skills', path: '/docs/trainable-skills' },
  { label: 'NPC System', path: '/docs/npc-system' },
  { label: 'Item System', path: '/docs/item-system' },
  { label: 'Faction System', path: '/docs/faction-system' },
  { label: 'Examine Skill', path: '/docs/examine-skill' },
]}
```

### 7.2 Doc Page Structure

Each doc page is a static React component with markdown-like styled content. No external content management — it's compiled into the bundle.

```tsx
// routes/docs/game-mechanics.tsx
// Plain text with styled sections — no form state, no API calls
```

**Template for each page:**
- `PageHeader` with title, no back button (or back to `/docs`)
- Styled `<section>` blocks with `<h2>` headings
- `<code>` blocks for examples
- `<table>` for data reference
- Consistent CSS: prose-like reading experience on dark theme

### 7.3 Doc Pages to Create (in order)

| Page | Contents | Status |
|------|----------|--------|
| `/docs` | Index — list all doc pages with one-line summaries | TODO |
| `/docs/game-mechanics` | Overview of all game systems. Links to sub-pages. | TODO |
| `/docs/ability-system` | Full ability field reference (table from section 4.1), effect types explained, classless skills list, formula examples | TODO |
| `/docs/combat-guide` | Combat flow, tick system, damage formula, dodge/parry mechanics, classless skills in combat | TODO |
| `/docs/trainable-skills` | How skills relate to abilities, training mechanics, requirement format, skill categories | TODO |
| `/docs/npc-system` | NPC lifecycle (template → instance → respawn), level scaling, race effects | TODO |
| `/docs/item-system` | Equipment slots, damage/armor calculation, tags and abilities, item categories | TODO |
| `/docs/faction-system` | Standing mechanics, categories, how factions affect gameplay | TODO |
| `/docs/examine-skill` | Examine command, hidden details, skill levels, DC checks, XP rewards | TODO |
| `/docs/config-reference` | What each config key does (from the Config page) | TODO |

### 7.4 Routing

Use TanStack Router's flat route convention. All docs live under `admin/src/routes/docs/`.

```
routes/
  docs.tsx           ← index layout / list
  docs/
    game-mechanics.tsx
    ability-system.tsx
    combat-guide.tsx
    ...
```

After creating files, run `npm run build:routes` to regenerate `routeTree.gen.ts`.

---

## 8. Implementation Order

### Phase 1: Foundation (est. 1 dispatch)
1. Build `Tooltip.tsx` component
2. Build `TooltipIcon.tsx` component
3. Modify `FormFields.tsx` — add optional `tooltip` prop to all field components
4. Add CSS for tooltip positioning and animation
5. Add `aria-describedby` wiring

### Phase 2: Abilities Tooltips (est. 1 dispatch)
1. Add `tooltip` prop to every field in `abilities.tsx` AbilityForm
2. Verify visually — hover each field, check tooltip appears

### Phase 3: Remaining Pages (est. 1 dispatch per page)
1. Skills page tooltips
2. NPCs page tooltips
3. Items page tooltips
4. Factions page tooltips
5. Config page tooltips

### Phase 4: Docs Section (est. 1 dispatch per page, or batched)
1. Create `DocsIcon` in `components/icons/`
2. Add Docs to sidebar nav with child routes
3. Create each doc page component
4. Run `npm run build:routes` after adding routes

### Phase 5: Polish
1. Mobile: tooltips should appear on tap (or be hidden gracefully)
2. Ensure no tooltip overflows viewport (CSS `max-width` + `overflow-wrap`)
3. Verify screen reader announces all tooltip content

---

## 9. Dispatch Strategy

Each phase above maps to a single Claude dispatch with a precise prompt. Example for Phase 1:

```
Task: Build Tooltip.tsx and TooltipIcon.tsx in admin/src/components/.
Modify admin/src/components/FormFields.tsx to add optional 'tooltip' prop
to FormField, NumberField, SelectField, TextareaField.

Rules:
- Pure CSS positioning, no tooltip library
- Accessible: aria-describedby, role="tooltip"
- Dark theme matching existing admin UI
- Max width 240px, 300ms hover delay
- TooltipIcon: inline 4×4 circle with ? or i

Do NOT touch any management pages.
Do NOT add routes.
Do NOT modify the sidebar.
After writing files, run 'npm run build' in admin/ to verify.
Report: files changed, build status.
```

---

## 10. Files to Create / Modify

### New files:
- `admin/src/components/Tooltip.tsx`
- `admin/src/components/TooltipIcon.tsx`
- `admin/src/components/icons/DocsIcon.tsx`
- `admin/src/routes/docs.tsx`
- `admin/src/routes/docs/game-mechanics.tsx`
- `admin/src/routes/docs/ability-system.tsx`
- `admin/src/routes/docs/combat-guide.tsx`
- `admin/src/routes/docs/trainable-skills.tsx`
- `admin/src/routes/docs/npc-system.tsx`
- `admin/src/routes/docs/item-system.tsx`
- `admin/src/routes/docs/faction-system.tsx`
- `admin/src/routes/docs/examine-skill.tsx`
- `admin/src/routes/docs/config-reference.tsx`

### Modified files:
- `admin/src/components/FormFields.tsx` — add `tooltip` prop
- `admin/src/components/Sidebar.tsx` — add Docs nav
- `admin/src/routes/_auth/abilities.tsx` — add tooltips
- `admin/src/routes/_auth/skills.tsx` — add tooltips
- `admin/src/routes/_auth/npcs.tsx` — add tooltips
- `admin/src/routes/_auth/items.tsx` — add tooltips
- `admin/src/routes/_auth/factions.tsx` — add tooltips
- `admin/src/routes/_auth/config.tsx` — add tooltips

---

## 11. Open Questions

1. Should tooltips also appear on table column headers (DataTable)? e.g. hovering "Effect Value" column header explains what it means?
2. Should the Docs section be a collapsible group in the sidebar, or a single "Docs" page with internal navigation?
3. Should we pre-populate doc page content from existing markdown files (`SKILLS-README.md`, `AGENTS.md`)?

**Recommended answers:**
1. **No** for table headers — tooltips on form fields are higher value. Table headers are self-evident once form fields are documented.
2. **Collapsible group** — the sidebar has room, and it's faster to navigate.
3. **Yes, partially** — `SKILLS-README.md` has good content for Combat Guide and Ability System. Extract and reformat, don't duplicate maintenance.
