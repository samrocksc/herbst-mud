# CLAUDE.md — Claude Code Project Context for herbst-mud

## What This Project Is
herbst-mud is a Go/PostgreSQL MUD game engine designed to run multiple MUDs
and storylines. SSH client on port 4444, REST API on port 8080, admin panel
on port 3000.

## Critical Rules (NEVER SKIP)
1. **Dual ent generate**: After ANY change to `server/db/schema/` or
   `herbst/db/schema/`, run `ent generate` in BOTH directories:
   ```bash
   cd /home/sam/GitHub/herbst-mud/server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   cd /home/sam/GitHub/herbst-mud/herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   ```
2. **Rebuild after changes**: After editing any `.go` file, rebuild:
   ```bash
   cd /home/sam/GitHub/herbst-mud && make build-all
   ```
3. **One ticket at a time**: Never work on multiple tickets in parallel.
   Use `--worktree` for each ticket to keep isolation.
4. **Always run tests** before pushing:
   ```bash
   make test && cd server && go test -v
   ```
5. **Restart services** after schema or route changes:
   ```bash
   cd /home/sam/GitHub/herbst-mud && make stop && make dev
   ```

## Architecture
- `herbst/` — SSH client (bubbletea TUI), separate Go module
- `server/` — REST API (Gin), separate Go module, ent ORM
- `admin/` — Vite/React/TanStack admin panel
- `content/` — YAML data-driven content (skills, NPCs, items, rooms)

## Key Directories
- `server/db/schema/` — ent schema definitions (26 entities including AbilityEffect)
- `herbst/db/schema/` — ent schema definitions (8 entities, subset)
- `server/routes/` — Gin route handlers
- `herbst/cmd_*.go` — MUD command handlers
- `features/` — Gherkin BDD feature files

## Admin Routing (TanStack Router file-based)
- `resource.tsx` → list page, `resource.$id.tsx` → detail/edit, `resource.new.tsx` → create page
- Create forms are standalone pages at `/resource/new` (NOT modals/inline toggles)
- List pages use location gating: `pathname === '/resource' ? <List/> : <Outlet/>`
- After tsr generate: `admin/src/routeTree.gen.ts` is auto-generated, do not edit manually

## Domain Model (Ability/Skill/Stat/Effect)
- **Abilities** = actions characters perform (Concentrate, Haymaker, Fireball). Entity: `Ability`
- **Skills** = leveled proficiencies (Blades, Staves). Stored as flat Character columns, future: `Skill` entity
- **Stats** = numeric attributes (Strength, Dexterity, etc.). Fields on `Character`
- **Effects** = what happens when an ability fires. Entity: `AbilityEffect` linked to `Ability` via `effects` edge
- Abilities use `ability_class` field: "active" or "passive" (formerly "Talents")
- API paths: `/api/abilities` (NOT `/api/skills`), `/api/abilities/:id/effects` for effects

## Code Style
- Files MUST NOT exceed 100 lines. Break into new files.
- Functional over OOP. Keep code simple and modular.
- Use JSDoc-style comments, avoid inline comments.
- Sign commits with team badge emoji: 🟣 Donatello, 🔴 Raphael, 🐀 Splinter

## Service Management
- Start all: `make dev-all`
- Start backend only: `make dev`
- Stop all: `make stop`
- SSH logs: `tail -f /tmp/herbst-ssh.log`
- Web logs: `tail -f /tmp/herbst-web.log`
- Health check: `curl -s http://localhost:8080/healthz`

## Database
- PostgreSQL 15, ent ORM with auto-migration
- Dev connection: host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable
- Production: uses DATABASE_URL env var with sslmode=require

## Current State (2026-05)
- Ability/Skill/Effect refactoring: Phase 1 (rename) and Phase 2 (Effect entity) complete
- Skill entity renamed to Ability throughout DB, API, and admin
- AbilityEffect entity created for multi-effect support
- Talents not yet merged into Abilities (Phase 3 pending)
- Hardcoded ClasslessSkills not yet converted to generic effects (Phase 4 pending)
- Code health: 91/100
- ~52 stale merged branches need cleanup

## Common Failure Modes (WATCH FOR THESE)
1. Forgot ent generate → compile errors about missing generated code
2. Edited schema in server/ but not herbst/ → client can't deserialize
3. Didn't rebuild after Go changes → running stale binary
4. Didn't restart services → API returns old responses
5. Untracked files in main working tree → git pollution
6. Used "Skill" when you mean "Ability" → API is `/api/abilities`, entity is `Ability`
7. Used old handler keys (concentrate, haymaker) → use generic effect types (buff, damage, stun)