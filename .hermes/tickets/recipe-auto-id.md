# Ticket: Auto-generate unique IDs for recipes

**Status:** TODO
**Priority:** Medium
**Labels:** enhancement
**Reporter:** Leonardo 🐢

## Description
Recipes currently require a manually-entered ID when creating/editing. This is fragile — duplicate IDs, typos, and confusion are inevitable as the recipe catalog grows.

## Problem
- Manually entering recipe IDs is error-prone (duplicates, formatting inconsistencies)
- No validation or collision detection on the admin form
- Will cause chaos as more recipes are added across multiple worlds
- The same pattern was already a problem for equipment templates (PK collision bug) — this prevents future incidents

## Proposed Solution
Auto-generate a unique, human-readable ID for each recipe when it is created:
- Slug-based from recipe name (e.g., "Supreme Pizza" → "supreme-pizza") with dedup suffix if collision
- The form should not have a manual ID field — the ID is assigned server-side on creation
- Existing recipes retain their current IDs after migration

## Acceptance Criteria
- [ ] Recipe creation does not require a manual ID field
- [ ] Server auto-generates a unique ID when a recipe is created
- [ ] Existing recipes retain their current IDs after migration
- [ ] No regressions in recipe search, update, or display functionality
- [ ] Feature file written at `features/recipe-auto-id.feature`

## Files to Modify
- `server/routes/recipe_routes.go` — ID generation on create
- `server/ent/schema/recipe.go` — possibly auto-increment or generated ID field
- `admin/src/routes/_auth/RecipeForm.tsx` — remove manual ID field
- `server/dbinit/crafting_seed.go` — seed recipes already have IDs, no change needed

## Notes
This is the same class of bug that caused equipment templates to be invisible (manual ID entry without validation). The recipe system should use the same fix pattern: server-side ID generation with slug-based naming.
