# Herbst MUD — Developer Guide

> **Audience:** Engineers contributing to the codebase.
> **Goal:** Understand the architecture, build the project, and follow conventions.

---

## Architecture Overview

```
┌─────────────┐  SSH 4444    ┌──────────────┐  HTTP 8080  ┌─────────────┐
│   herbst/   │─────────────▶│   server/    │◀───────────│   admin/    │ 3000
│ SSH TUI     │              │ REST API     │            │ React Admin │
│ (BubbleTea) │              │ (Gin + ent)  │            │ (Vite)      │
└─────────────┘              └──────────────┘            └─────────────┘
                                    │                            ▲
                                    ▼                            │
                           ┌─────────────────┐     ┌─────────────┐
                           │   PostgreSQL    │     │ web-client/ │ 3001
                           │   (ent ORM)     │     │ Browser MUD │
                           └─────────────────┘     └─────────────┘
```

**Two Go binaries + two React frontends:**

| Module | Directory | Entry Point | Role |
|--------|-----------|-------------|------|
| SSH Client | `herbst/` | `herbst/main.go` | TUI MUD client (Go) |
| REST Server | `server/` | `server/main.go` | REST API + WebSocket (Go) |
| Admin Panel | `admin/` | `admin/src/main.tsx` | Vite/React SPA |
| Web Client | `web-client/` | `web-client/src/main.tsx` | Browser game client |

---

## Quick Build

```bash
# Full build
cd /home/sam/GitHub/herbst-mud
make build-all

# Server only
cd server && go build -o bin/server .

# Admin only
cd admin && npm run build
```

---

## The Golden Rules

1. **Dual ent generate** — after ANY schema change in `server/db/schema/` OR `herbst/db/schema/`, regenerate BOTH:
   ```bash
   cd server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   cd herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   ```

2. **Rebuild after Go changes** — `make build-all`

3. **Restart after schema/route changes** — `make stop && make dev`

4. **One ticket at a time** — use worktrees for isolation

5. **Tests before push** — `make test && cd server && go test ./...`

---

## Tech Stack

| Layer | Technology | Where |
|-------|-----------|-------|
| Backend API | Go + Gin | `server/` |
| SSH Game | Go + BubbleTea + Lipgloss | `herbst/` |
| Database | PostgreSQL + ent ORM | `server/db/`, `herbst/db/` |
| Admin UI | React 19 + Vite + Tailwind | `admin/src/` |
| Admin Router | TanStack Router v1 (file-based) | `admin/src/routes/` |
| Admin Data | TanStack Query v5 | `admin/src/hooks/` |
| Web Client | React 19 + Vite | `web-client/src/` |
| Tests | Vitest v4 + jsdom + MSW | `admin/src/` |

---

## Directory Map

```
herbst-mud/
├── herbst/              # SSH TUI client (separate Go module)
│   ├── cmd_*.go         # MUD commands (look, attack, say, etc.)
│   ├── game_*.go        # Game logic (combat, inventory, etc.)
│   ├── db/              # ent generated code
│   └── dbinit/          # Seeding + initial world setup
├── server/              # REST API server (separate Go module)
│   ├── routes/          # Gin HTTP handlers
│   ├── service/         # Business logic (CombatService, etc.)
│   ├── repository/      # Repository interfaces + implementations
│   ├── db/              # ent generated code
│   ├── dblog/           # Structured logging to app_logs table
│   └── middleware/      # Auth, world_id filtering
├── admin/               # Web admin panel (React/Vite)
│   ├── src/routes/      # TanStack Router file-based routes
│   ├── src/components/  # Shared UI (FormFields, SearchableSelect, etc.)
│   └── src/hooks/       # TanStack Query hooks
├── web-client/          # Browser game client (React/Vite)
│   └── src/             # Game screens, input bar, hotkey bar
├── content/             # YAML/JSON world content files
├── docs/                # Documentation hierarchy
├── tickets/             # Implementation tickets
└── archive/             # Stale content (features, old plans, etc.)
```

---

## File Size Convention

**The 100-line rule** — if a file exceeds 100 lines, split it by subdomain:

```
routes/character.go           → routes/character_list.go
                              → routes/character_crud.go
                              → routes/character_stats.go
```

Exceptions: generated code (`routeTree.gen.ts`, ent files) and deeply-nested React components that are structurally flat.

---

## Frontend Conventions

### Admin Routing (TanStack Router)
- `resource.tsx` → list page
- `resource.$id.tsx` → detail/edit
- `resource.new.tsx` → create page (standalone, **not modal**)
- Location gating: `pathname === '/resource' ? <List/> : <Outlet/>`

### Form Components
- Relational fields → `SearchableSelect`, `ResourceSearchSelect`, `ResourceMultiSelect`
- ID fields → `ResourceIdField` (with inline validation)
- Tags → `TagInput` (only when free creation is intended)
- Raw text/number inputs for foreign keys are **forbidden** — they create broken tails.

### Query Hooks
All API calls go through TanStack Query hooks in `admin/src/hooks/`:
- `useQuery` for reads
- `useMutation` for writes, with `queryClient.invalidateQueries()` on success

---

## Testing

```bash
cd admin && npm test          # Vitest
cd server && go test ./...     # Go unit tests
```

**Test stack:** Vitest + jsdom + MSW (Mock Service Worker) for API mocking. MSW must be installed before writing route/hook tests.

---

## Contributing

1. Pick a ticket from `tickets/`
2. Create a branch or worktree
3. Implement with tests
4. Run `make build-all && make test`
5. Open PR

See `docs/OPERATIONS/INDEX.md` for deployment operations and `docs/OPERATIONS/INSTALL.md` for installation and upgrades.
