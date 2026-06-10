# COMPLETION REPORT — Sub-ticket A: TICKET-PLAYER-COMBAT-QA-001

## Summary
Fixed Findings #4 and #7. `make build-all` passes. Server restarted and healthy.

---

## Finding #4 — Default ability equip for fresh characters ✅

**Files modified:**
- `server/service/character.go`

**Change:** After character creation, `equipDefaultAbilities()` is called:
1. Queries abilities by class (e.g. `ListByClass(ctx, worldID, class)`)
2. Falls back to `ListClassless(ctx, worldID)` if class yields nothing
3. Creates `character_abilities` rows for up to 4 abilities (slots 1–4)
4. Non-fatal — logs warnings on failure, character still playable with auto-attack

**Key code added:**
```go
func (s *characterService) equipDefaultAbilities(ctx context.Context, charID int, class, worldID string) {
    abilities, err := s.repos.Ability.ListByClass(ctx, worldID, class)
    if err != nil || len(abilities) == 0 {
        abilities, err = s.repos.Ability.ListClassless(ctx, worldID)
    }
    // Limit to 4, create character_abilities rows
}
```

**Constraint check:**
- ✅ No new abilities added to `abilities` table
- ✅ Max 4 abilities (slot 1–4)
- ✅ Uses existing repo (`CharacterAbility.Create`)
- ✅ Server restarted after change

---

## Finding #7 — Damage_logs persistence on every combat hit ✅

**Files modified:**
- `server/service/interface.go`
- `server/service/combat_service_impl.go`
- `server/routes/character_combat.go`

**Change:**
1. `ApplyDamage` signature extended to accept `attackerID int` as the new first param
2. Inside `ApplyDamage`, after HP update: `s.LogDamage(ctx, attackerID, targetID, damage)` called for every `damage > 0`
3. Route updated — attackerID passed from request body, `LogDamage` no longer called separately

**Key code added to `ApplyDamage`:**
```go
// Persist damage to damage_logs table on every hit
if damage > 0 {
    s.LogDamage(ctx, attackerID, targetID, damage)
}
```

**Constraint check:**
- ✅ Uses existing `LogDamage` method (direct insert, fast)
- ✅ Schema for `damage_logs` already exists (attacker_id, target_id, damage, created_at)
- ✅ Logs on every hit — not just kills
- ✅ Route caller updated (attackerID now passed through to `ApplyDamage`)

---

## Verification steps

```sql
-- Finding #4: After creating a character
SELECT ca.slot, a.name FROM character_abilities ca
JOIN abilities a ON ca.ability_id = a.id
WHERE character_id = <new_char_id> ORDER BY slot;

-- Finding #7: After combat
SELECT count(*) FROM damage_logs;
SELECT * FROM damage_logs ORDER BY created_at DESC LIMIT 5;
```

---

## Lines modified
- `server/service/character.go`: +30 lines (equipDefaultAbilities + call in CreateCharacter)
- `server/service/combat_service_impl.go`: +2 lines (LogDamage call inside ApplyDamage)
- `server/service/interface.go`: +1 line (signature update)
- `server/routes/character_combat.go`: +1 line, -1 line (attackerID pass-through, removed duplicate LogDamage)

**Total: ~35 lines added/changed**
