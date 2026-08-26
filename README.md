# Herbst MUD

A modern Multi-User Dungeon (MUD) engine built with Go, TypeScript, and PostgreSQL.
Herbst is an **engine** that can run multiple worlds and storylines, not a single
game — content is data-driven through the admin panel and the database.

## Architecture

Herbst MUD is split into four services plus a database:

| Service      | Stack                  | Port | Purpose                                    |
|--------------|------------------------|------|--------------------------------------------|
| `herbst/`    | Go + SSH               | 4444 | Classic SSH MUD client (Bubble Tea TUI)    |
| `server/`    | Go + Gin + ent ORM     | 8080 | REST API + WebSocket game backend          |
| `admin/`     | Vite + React + TanStack| 5173 | Admin panel: build worlds, NPCs, items     |
| `web-client/`| Vite + React           | 5174 | Browser MUD client                         |
| PostgreSQL   | postgres:15            | 5432 | All game/engine data                       |

- The `server/` backend owns the REST API (`/api/*` via the Vite proxy) and the
  game WebSocket (`/ws`). Schema is managed with the **ent ORM** and migrated
  with **Atlas**.
- Worlds (e.g. "Ooze Surfers") are **data in the database**, scoped by
  `world_id`. Each world is an isolated set of rooms, NPCs, items, and quests.
- `admin/` and `web-client/` are separate Vite apps; the admin panel drives the
  live backend through the REST API.
- `features/` holds Gherkin BDD specifications used by the server test suite.

## Quick Start (local dev)

Requires Go 1.25+, Node.js 18+, Docker + Docker Compose, and a PostgreSQL 15.

```bash
# 1. Bring up the full stack (SSH + web + admin + web client)
make dev-all

# 2. The database schema auto-migrates on first server start (Atlas).

# 3. Connect
#    SSH client:
ssh -p 4444 localhost
#    Web client:  http://localhost:5174
#    Admin panel: http://localhost:5173
#    API:         http://localhost:8080
```

Key `make` commands:

```bash
make dev          # SSH + web + web client
make dev-all      # everything, including admin panel
make test         # TUI/SSH tests
make test-server  # backend (REST) tests
make test-server-bdd  # backend Gherkin BDD tests
make reload       # rebuild + restart SSH server
make reload-web   # rebuild + restart web API server
make stop         # stop all services
```

### Database

The database schema is defined by the ent schemas in `server/db/schema/`.
On startup the web server auto-runs schema migration. For explicit,
versioned migrations use the Atlas CLI (see the `atlas.hcl` in the repo root
and the upgrade procedures below).

## Deployment

The production stack runs on **Docker Compose with images pulled from GitHub
Container Registry** (`docker-compose.prod.yml`). `docker-compose.yml` is the
local build-based compose for development.

```bash
# Production: images from GHCR, tunneled behind Cloudflare
cp .env.example .env   # set DB_PASSWORD, JWT_SECRET, CF_TUNNEL_TOKEN
make deploy-up

# Local dev: build images from source
make deploy-build && make deploy-up
```

## API Reference

The complete REST API is documented in **[`docs/API_REFERENCE.md`](docs/API_REFERENCE.md)**.

## Contributing

- Write code with tests. Backend tests run via `make test-server`; the Gherkin
  specs live in `features/`.
- Keep world data in the database, never hardcoded in the engine.
- This repo is public — never commit secrets or credentials. `.env` and all
  `.env` files are gitignored; a `gitleaks` pre-commit hook enforces this.
