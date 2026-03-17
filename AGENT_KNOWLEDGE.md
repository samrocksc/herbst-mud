# AGENT_KNOWLEDGE.md - Herbst-MUD Project Knowledge

This file contains essential project knowledge for agents working on Herbst-MUD.

## Project Overview

**Herbst-MUD** is a Multi-User Dungeon game with:
- Go backend (server/, herbst/)
- React admin frontend (admin/)
- Terminal/TUI client (herbst/ - Go-based SSH client)
- PostgreSQL database (Ent ORM)

## Architecture

```
herbst-mud/
├── admin/           # React admin dashboard (Vite, TanStack Router)
├── server/          # Go REST API server
├── herbst/          # Go TUI client (SSH-based)
├── docs/            # Project documentation
├── features/        # Feature specs (Gherkin format)
└── Makefile         # Development commands
```

## Tech Stack

- **Backend:** Go, Ent ORM, PostgreSQL
- **Frontend:** React, Vite, TanStack Router
- **TUI:** Go, Bubble Tea, Lipgloss
- **Testing:** Gherkin BDD tests
- **Container:** Docker Compose

## Key Files

| File | Purpose |
|------|---------|
| `Makefile` | `make dev` starts servers, `make test` runs tests |
| `docker-compose.yml` | PostgreSQL, SSH server, admin UI |
| `docs/DEVELOPMENT.md` | Setup instructions |
| `docs/CODE.md` | Code standards |
| `docs/DATABASES.md` | Database schema |
| `docs/API.md` | REST API endpoints |
| `docs/TESTING.md` | Testing guidelines |

## Important Conventions

### GitHub Projects
- Project board: "Turtle Time" (#2)
- Use labels: donnie, raph, mikey, leo

### Code Standards
- Use Ent ORM for database
- Write Gherkin feature tests in `features/`
- Semantic versioning for releases
- Emoji badges in commits: 🔵🟣🔴🎨

### Database
- Ent schema in `server/db/schema/` and `herbst/db/schema/`
- Run code generation after schema changes:
  ```bash
  cd herbst && go generate ./...
  cd server && go generate ./...
  ```

## Starting Development

```bash
# Clone and start
git clone https://github.com/samrocksc/herbst-mud.git
cd herbst-mud

# Start all services
make dev

# Run tests
make test
```

## Common Tasks

| Task | Command |
|------|---------|
| Start dev server | `make dev` |
| Stop servers | `make stop` |
| Run tests | `make test` |
| Start admin | `make start-admin` |
| Run BDD tests | `make test-bdd` |

## Known Issues

- NPC Template schema needs `go generate` to build
- Split-screen UI requires terminal dimensions passed to screens

## Last Updated

2026-03-17