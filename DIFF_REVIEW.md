# DIFF_REVIEW — Abilities Form RESTful Route Refactoring

## Summary

Refactored the abilities admin forms to follow the established patterns used by other resources (items, characters). Also fixed a nested `<form>` bug that prevented abilities from saving correctly when adding effects.

## Changes

### 1. New: `admin/src/routes/_auth/abilities.new.tsx` (standalone create page)
- Extracted `CreateAbilityForm` from the inline form in `abilities.tsx` into a full-page create route at `/abilities/new`
- Follows the same pattern as `items.new.tsx`
- Navigates to `/abilities` list on success
- Includes all ability fields: name, description, type, tags, costs, cooldown, proc settings, class

### 2. New: `admin/src/routes/_auth/-abilities.$abilityId.editForm.tsx` (private edit component)
- Extracted edit form from `abilities.$abilityId.tsx` into a reusable private component
- Follows the same pattern as `-items.$itemId.editForm.tsx`
- Imports `DeleteConfirmation` component (instead of inline modal)
- Imports `EffectsSubForm` for managing ability effects
- Props: `ability`, `abilityId`, `onDone` callback

### 3. New: `admin/src/routes/_auth/-abilities.$abilityId.detailView.tsx` (private detail component)
- Extracted read-only detail view from `abilities.$abilityId.tsx`
- Follows the same pattern as `-items.$itemId.detailView.tsx`
- Shows all ability fields in a grid layout with `DetailField` sub-component

### 4. Updated: `admin/src/routes/_auth/abilities.tsx` (list page)
- Replaced inline `CreateAbilityForm` toggle with navigation to `/abilities/new`
- Changed "+ Add Ability" button to `navigate({ to: '/abilities/new' })`
- Added location gating: `pathname !== '/abilities' ? <Outlet/> : <List/>`
- Removed unused imports (`useCreateAbility`, `useDeleteAbility`, `AbilityInput`, `useTags`, `TagInput`, form field components, `showToast`)

### 5. Updated: `admin/src/routes/_auth/abilities.$abilityId.tsx` (detail page)
- Replaced inline edit form and detail view with imported `AbilityEditForm` and `AbilityDetailView`
- Added `editing` state toggle (same pattern as items detail page)
- Edit button label changes between "Edit" and "Cancel"

### 6. Fixed: `admin/src/components/EffectsSubForm.tsx` (nested form bug)
- Changed `NewEffectForm` from `<form>` to `<div>` to fix invalid nested HTML forms
- Changed button from `type="submit"` with `onSubmit` handler to `onClick` handler
- **Root cause**: The effect form was nested inside the parent ability `<form>`, which is invalid HTML. Browsers handle nested forms inconsistently — clicking "Add Effect" could submit the parent ability form instead of (or in addition to) the effect form.

## Route Tree

New routes registered via `tsr generate`:
- `/abilities` — list page (unchanged path)
- `/abilities/new` — standalone create page (NEW)
- `/abilities/$abilityId` — detail/edit page (unchanged path, refactored internals)

## REST API Endpoints (unchanged)

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/abilities` | List abilities (filters: type, ability_class, search) |
| POST | `/api/abilities` | Create ability |
| GET | `/api/abilities/:id` | Get single ability |
| PUT | `/api/abilities/:id` | Update ability |
| DELETE | `/api/abilities/:id` | Delete ability |
| GET | `/api/abilities/:id/effects` | List effects for ability |
| POST | `/api/abilities/:id/effects` | Create effect |
| PUT | `/api/ability-effects/:id` | Update effect |
| DELETE | `/api/ability-effects/:id` | Delete effect |

The server-side routes were already fully RESTful and required no changes.

## Testing Notes
- TypeScript compiles cleanly (`npx tsc --noEmit` passes)
- Regenerate route tree: `npx tsr generate`
- If running: restart admin dev server to pick up new routes
