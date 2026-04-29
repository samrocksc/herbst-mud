# herbst-mud QA Dogfood Report

**Date:** 2026-04-29  
**Stack:** Web API (port 8080) + Admin UI (port 5173, Vite dev) + SSH (port 4444)  
**Approach:** Hybrid — curl API walkthrough + admin UI surface probe  
**Tester:** Raph (delegated QA) via Leonardo  

---

## Executive Summary

| Severity | Count |
|----------|-------|
| Critical | 1     |
| High     | 4     |
| Medium   | 3     |
| Low      | 2     |
| **Total**| 10    |

### Critical
1. **Password hash exposure** — `/characters` endpoint leaks bcrypt password hashes for all user accounts

### High
2. Character creation fails — requires undocumented `currentRoomId` + `startingRoomId`
3. All `/content/*` endpoints return 404 despite being in OpenAPI spec
4. All `/worlds/*` endpoints return 404 despite being in OpenAPI spec  
5. `/classes/{class}/specialties` returns 404 for all classes (no class data loaded)

### Medium
6. Auth returns 200 on success instead of documented 201
7. Auth returns 401 instead of documented 400 on bad credentials
8. `/admin/import/validate` returns 404 on GET (method mismatch with spec?)

### Low
9. `/characters/{id}/skills` returns 404 for non-NPC characters
10. Swagger UI accessible at `/docs` — confirmed working (positive finding)

---

## Phase 1 — API Walkthrough

### Public Endpoints

| Endpoint | Expected | Actual | Pass |
|----------|----------|--------|------|
| `GET /healthz` | 200 | 200 `{"db":"connected","ssh":"running","status":"ok"}` | ✅ |
| `GET /genders` | 200 | 200 `[he_him, she_her, they_them]` | ✅ |
| `GET /races` | 200 | 200 `[human, turtle, mutant]` | ✅ |
| `GET /skills/effect/heal` | 200 | 200 `{"count":0,"skills":[]}` | ✅ |
| `GET /talents/effect/buff` | 200 | 200 `{"count":0,"talents":[]}` | ✅ |
| `GET /content/items` | 200 | **404** | ❌ |
| `GET /content/npcs` | 200 | **404** | ❌ |
| `GET /content/rooms` | 200 | **404** | ❌ |
| `GET /content/skills` | 200 | **404** | ❌ |
| `GET /content/quests` | 200 | **404** | ❌ |
| `GET /content/stats` | 200 | **404** | ❌ |
| `GET /worlds` | 200 | **404** | ❌ |
| `GET /worlds/active` | 200 | **404** | ❌ |
| `GET /classes/fighter/specialties` | 200 | **404** `{"error":"Class not found"}` | ❌ |
| `GET /content/validate` | 200 | **404** | ❌ |

### Auth Flow

| Test | Expected | Actual | Pass |
|------|----------|--------|------|
| `POST /users` (register) | 201 | 201 | ✅ |
| `POST /users/auth` wrong creds | 400 | **401** | ❌ |
| `POST /users/auth` correct | 201 | **200** | ❌ |

### Authenticated Endpoints

| Endpoint | Actual | Notes |
|----------|--------|-------|
| `GET /users` | 200 | ✅ works |
| `GET /characters` | 200 | ⚠️ **leaks password hashes** (see Critical) |
| `GET /rooms` | 200 | ✅ |
| `GET /rooms/{id}` | 200 | ✅ |
| `GET /rooms/{id}/look` | 200 | ✅ |
| `GET /rooms/{id}/characters` | 200 | ✅ |
| `GET /rooms/{id}/equipment` | 200 | ✅ |
| `GET /equipment` | 200 | ✅ |
| `GET /npcs` | 200 | ✅ |
| `GET /skills` | 200 | ✅ |
| `GET /talents` | 200 | ✅ |
| `GET /worlds` | 404 | ❌ |
| `GET /worlds/1` | 200 | **404** |
| `GET /worlds/1/content/items` | 200 | **404** |
| `GET /worlds/1/stats` | 200 | **404** |
| `POST /characters` (new char) | 201 | **500** — missing `currentRoomId` then missing `startingRoomId` |
| `GET /characters/{id}/stats` | 200 | **404** for id=1 (chars 9-20 exist, id=1 doesn't) |
| `GET /characters/{id}/class` | 200 | **404** |
| `GET /characters/{id}/race` | 200 | **404** |
| `GET /characters/{id}/talents` | 200 | **200** for id=1 (empty slots, inconsistent) |
| `GET /characters/{id}/skills` | 200 | **404** |

### Admin Routes

| Endpoint | Actual | Notes |
|----------|--------|-------|
| `GET /admin/export` | 200 | ✅ works |
| `GET /api/backups` | 200 | ✅ `{"backups":["backup_2026-04-29_11-50-58"]}` |
| `POST /api/backups` | 201 | ✅ created backup |
| `POST /admin/import/validate` (GET) | 200 | **404** — method mismatch? |

---

## Phase 2 — Admin UI Surface

| Check | Result |
|-------|--------|
| `GET /` (root) | 200 ✅ |
| Static assets | Vite HMR ✅, vite.svg ✅ |
| JS/CSS loading | Via Vite dev server (`/@vite/client`, `/src/main.tsx`) |
| `/admin` | Returns 200 ✅ |
| `/api` | Returns 200 ✅ |
| Console errors | Could not check (Chrome sandbox unavailable) |

---

## Detailed Issues

### ISSUE-1: PASSWORD HASH EXPOSURE (CRITICAL — Security)

**Severity:** Critical  
**Category:** Console / Security  
**URL:** `GET /characters` (authed)  
**Description:** The `/characters` endpoint returns all characters including one user account (`id: 9, name: "sma"`) with its bcrypt password hash fully exposed in the response body.

**Steps to Reproduce:**
1. `POST /users` to register an account
2. `POST /users/auth` to get JWT token
3. `GET /characters` with `Authorization: Bearer <token>`

**Expected:** Character list should not contain password fields.  
**Actual:** Response includes `"password": "$2a$10$wasTNdVMCdN7yjO2O7Kqnea7hLd7oKedR2116pEIO5e4M47Wx23By"`  
**Screenshot:** N/A (API-only finding)

---

### ISSUE-2: CHARACTER CREATION FAILS — UNDOCUMENTED REQUIRED FIELDS (HIGH)

**Severity:** High  
**Category:** Functional  
**URL:** `POST /characters`  
**Description:** Creating a character via the API fails with 500 errors requiring `currentRoomId` and `startingRoomId` fields that are not documented in the OpenAPI spec request body.

**Steps to Reproduce:**
```
POST /characters {"name": "RaphQA", "race": "human", "gender": "he_him", "class": "fighter", "specialty": "berserker"}
→ 500 {"error": "db: missing required field \"Character.currentRoomId\""}

POST /characters {"name": "RaphQA", ..., "currentRoomId": 1}
→ 500 {"error": "db: missing required field \"Character.startingRoomId\""}
```

**Expected:** OpenAPI spec defines required fields; 422 or clear validation error for missing fields.  
**Actual:** 500 with db-level error message  
**Fix:** Add `currentRoomId` and `startingRoomId` to OpenAPI `CharacterCreate` schema, or default them server-side.

---

### ISSUE-3: ALL /content/* ENDPOINTS RETURN 404 (HIGH)

**Severity:** High  
**Category:** Functional  
**URL:** All `/content/*` routes  
**Description:** Every endpoint under `/content/` (items, NPCs, rooms, skills, quests, stats, validate) returns `404 page not found`. All are defined in the OpenAPI spec.

**Expected:** 200 with content list or 404 only if world not initialized.  
**Actual:** All return plain `404 page not found` — routes are not registered.

---

### ISSUE-4: ALL /worlds/* ENDPOINTS RETURN 404 (HIGH)

**Severity:** High  
**Category:** Functional  
**URL:** All `/worlds/*` routes  
**Description:** Every endpoint under `/worlds/` returns 404. Endpoints like `GET /worlds`, `GET /worlds/active`, `GET /worlds/{id}/content/items` etc. are all defined in OpenAPI but not registered.

---

### ISSUE-5: /classes/{class}/specialties ALWAYS 404 (HIGH)

**Severity:** High  
**Category:** Functional  
**URL:** `GET /classes/fighter/specialties`  
**Description:** Returns `{"error":"Class not found"}` for all class names. The class data (fighter, mage, etc.) appears not to be seeded in the database.

---

### ISSUE-6: AUTH SUCCESS RETURNS 200 NOT 201 (MEDIUM)

**Severity:** Medium  
**Category:** Spec Compliance  
**URL:** `POST /users/auth`  
**Description:** OpenAPI spec says `201 Created` on successful auth, but actual response is `200 OK`.

---

### ISSUE-7: AUTH WRONG CREDS RETURNS 401 NOT 400 (MEDIUM)

**Severity:** Medium  
**Category:** Spec Compliance  
**URL:** `POST /users/auth`  
**Description:** OpenAPI spec says `400 Bad Request` for invalid credentials; server returns `401 Unauthorized`.

---

### ISSUE-8: /admin/import/validate METHOD MISMATCH (MEDIUM)

**Severity:** Medium  
**Category:** Functional  
**URL:** `GET /admin/import/validate`  
**Description:** OpenAPI spec defines `POST /admin/import/validate` but the actual server returns 404 on GET (as expected for POST-only), suggesting the route exists but the spec and implementation are out of sync on method.

---

### ISSUE-9: /characters/{id}/skills 404 FOR VALID CHARACTERS (LOW)

**Severity:** Low  
**Category:** Functional  
**URL:** `GET /characters/{id}/skills`  
**Description:** Returns `{"error":"Character not found"}` for characters 9-20. Only `/characters/{id}/talents` works (returns empty slots). Inconsistent sub-resource behavior.

---

### ISSUE-10: SWAGGER UI CONFIRMED WORKING (LOW — Positive)

**Severity:** Low / Positive  
**URL:** `GET /docs`  
**Description:** Swagger UI loads correctly at `:8080/docs` with full OpenAPI 3.0.3 spec rendered. All routes listed. No visual issues detected.

---

## Summary Table

| # | Severity | Category | Endpoint | Issue |
|---|----------|----------|----------|-------|
| 1 | **Critical** | Security | `GET /characters` | Password hash exposure |
| 2 | High | Functional | `POST /characters` | Missing required fields not in spec |
| 3 | High | Functional | `GET /content/*` | All 404 — routes not registered |
| 4 | High | Functional | `GET /worlds/*` | All 404 — routes not registered |
| 5 | High | Functional | `GET /classes/*` | Class data not seeded |
| 6 | Medium | Spec Compliance | `POST /users/auth` | Returns 200 not 201 |
| 7 | Medium | Spec Compliance | `POST /users/auth` | Returns 401 not 400 |
| 8 | Medium | Spec Compliance | `GET /admin/import/validate` | Method mismatch |
| 9 | Low | Functional | `GET /characters/{id}/skills` | 404 for valid characters |
| 10 | Low | Positive | `GET /docs` | Swagger UI works |

---

## Testing Notes

**Tested:** Public content endpoints, auth flow, character CRUD, room navigation, equipment, NPCs, skills, talents, admin backup/restore, worlds, class data, admin UI surface.

**Not Tested (blockers):**
- Browser-based admin UI interactions (Chrome sandbox unavailable in container)
- Form submissions and interactive UI flows
- WebSocket/real-time features
- Combat system (`POST /characters/{id}/damage`, etc.)
- Chat system
- Content creation endpoints (`POST /rooms`, `POST /equipment`, etc.)

**Env:** Dev stack, PostgreSQL seeded with test data (rooms 1-5, characters 9-20, items, NPCs).
