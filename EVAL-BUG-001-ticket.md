# Bug Ticket: EVAL-BUG-001 — Items invisible in Items list for world 2

## Status
**Open — Root cause verified, fix applied. Pending user confirmation after Vite restart.**

## Phase 1: Diagnosis

### Existing Behavior
- When logged into **world 2** and navigated to `http://100.67.206.65:5173/items`, the items list appears **empty** ("No items found").
- When creating a new item template (e.g. "test"), the POST to `/api/equipment-templates` succeeds with HTTP **201** and returns a complete payload including `id: 5`.
- After creation the item is still **not visible** in the list.
- Direct SELECT from the database shows the newly created row has `world_id='2'`.

### Expected Behavior
- Items created for world 2 should be **visible** in the items list when viewing world 2.
- The list should display all templates matching the current world's `world_id`.

## Phase 2: Terrain Map

### Affected Files
1. **`server/routes/equipment_template_routes.go`** — `createEquipmentTemplate` handler (line ~181)
2. **`server/repository/equipment_template_repo.go`** — `List()` method (line ~30)
3. **`admin/src/hooks/useEquipmentTemplates.ts`** — GET query appends `world_id` correctly
4. **`admin/src/routes/_auth/items.new.tsx`** — POST body includes `world_id` correctly

### Data Flow
```
Frontend POST /api/equipment-templates
 → WorldAccessMiddleware extracts world_id from JSON body (OK)
 → createEquipmentTemplate handler resolves worldID (OK)
 → handler calls repos.EquipmentTemplate.Create(...)
   → repository.EquipmentTemplate.Create() is called
     → builder.Save(ctx) writes to DB
       → DB stores world_id='2'  ← WRONG: world_id is empty

Frontend GET /api/equipment-templates?world_id=2
 → listEquipmentTemplates filters by world_id='2'
   → ent query: WHERE world_id = '2'
     → Returns only rows WHERE world_id='2'
       → Newly created rows (with world_id='') are EXCLUDED
```

## Phase 3: Root Cause

### The Code Path
In `server/routes/equipment_template_routes.go` line 181, the `CreateEquipmentTemplateInput` struct literal **does NOT include the `WorldID` field**.

```go
t, err := repos.EquipmentTemplate.Create(c.Request.Context(), repository.CreateEquipmentTemplateInput{
    Name:      req.Name,
    // ... many fields ...
    // BUG: WorldID is missing here!
})
```

The handler correctly resolves `worldID` from the JSON body at line 106-112, but the resolved `worldID` variable is **never passed** into the input struct. This causes ent's `SetWorldID()` to not be called, resulting in the DB writing the default value (empty string).

### Regression Source
This is a **regression** from the int-PK + slug refactor commit `f6a2b6f` (May 22). The previous fix commit `aa3122d` (May 20) had added `WorldID: worldID` to this struct literal, but the refactor removed it.

## Phase 4: Fix Applied

Commit `c20dcee` (May 22) restored `WorldID: worldID` at line 181:

```go
CreateEquipmentTemplateInput{
    Slug:    slug,
    Name:    req.Name,
    // ... other fields ...
    WorldID: worldID,   // ← restored
}
```

Backend rebuilt and restarted. Vite dev server was also restarted to serve fresh frontend bundles.

## Phase 5: Verification Checklist

- [ ] Fresh item created in world 2 appears in the items list immediately
- [ ] Existing items (created with empty world_id) are visible and accessible
- [ ] Items in the "default" world are unaffected (separate world_id scope)
- [ ] Edit/Delete actions on visible items work correctly

## Phase 6: Additional Findings

The `items.tsx` component uses `instance.equipment_template_id` for instance-to-template mapping, which is **type-safe** (`id: number` → `equipment_template_id: string` will fail TypeScript). This was confirmed not to be the cause of the invisibility bug, but should be revisited if instances need to reference templates by numeric ID.

## Domain
**Backend (Go/Gin, ent ORM)** — `server/routes/equipment_template_routes.go`

## Dispatch Recommendation
If the fix does not resolve the issue after the current restart:
- **Agent**: `codemaster` (Go/ent specialist)
- **Evidence**: `server/routes/equipment_template_routes.go` line 181, `server/repository/equipment_template_repo.go` line 30
- **Mandate**: Verify struct literal field injection and database state after create
