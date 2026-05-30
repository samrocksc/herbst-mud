# Herbst MUD — Operations Manual

> **Audience:** DevOps, server admins, anyone running the stack.
> **Goal:** Start the stack, keep it running, know where to look when it breaks.

---

## Prerequisites

- Docker + docker-compose
- PostgreSQL 15+ (or use the compose-provided db container)
- Ports: `4444` (SSH MUD), `8080` (REST API), `3000` (admin panel), `5173` (web client dev)

---

## Quick Start

```bash
# Start everything (db + server + admin + web-client + ssh)
make dev-all

# Or start just the backend stack
make dev
```

**Services started:**
| Service | URL | Port |
|---------|-----|------|
| PostgreSQL | `localhost:5432` | 5432 |
| REST API | `http://localhost:8080` | 8080 |
| Admin Panel | `http://localhost:3000` | 3000 |
| Web Client | `http://localhost:5173` | 5173 |
| SSH MUD | `ssh -p 4444 player@localhost` | 4444 |

---

## Make Targets

| Command | What it does |
|---------|-------------|
| `make dev` | Start DB + server + admin (docker-compose) |
| `make dev-all` | Start DB + server + admin + web-client + ssh |
| `make stop` | Stop all docker-compose services |
| `make build-all` | Build server binary, admin bundle, web-client bundle |
| `make test` | Run frontend tests (`cd admin && npm test`) |
| `make deploy` | Deploy to Digital Ocean App Platform (uses `.do/app.yaml`) |

---

## Environment Variables

Create `.env` in repo root. Required vars:

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/herbst?sslmode=disable
JWT_SECRET=changeme-at-least-32-chars-long
ADMIN_USERNAME=sma
ADMIN_PASSWORD=sma
```

Optional:
```bash
WORLD_ID=2          # Default world filter for admin panel
LOG_LEVEL=info      # slog level: debug, info, warn, error
```

---

## Logs & Debugging

### View live logs
```bash
docker-compose logs -f server
```

### Admin Logs SSE Stream
Open `http://localhost:8080/api/logs/stream` in a browser or curl:
```bash
curl -N http://localhost:8080/api/logs/stream
```

### Database connection issues
1. Check `DATABASE_URL` in `.env`
2. Verify PostgreSQL is running: `docker-compose ps db`
3. Check server logs for `connection refused` or `password authentication failed`

---

## Backup & Restore

### Create a backup
```bash
curl -X POST http://localhost:8080/api/backup \
  -H "Authorization: Bearer $TOKEN"
```

### List backups
```bash
curl http://localhost:8080/api/backups \
  -H "Authorization: Bearer $TOKEN"
```

### Restore from backup
```bash
curl -X POST http://localhost:8080/api/restore \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"backup_id": "backup_2026-04-29_11-50-58"}'
```

---

## Health Checks

| Endpoint | Expected |
|----------|----------|
| `GET /health` | `{"status":"ok"}` |
| `GET /api/health` | `{"status":"ok"}` |
| `GET /api/logs/services` | List of known log services |

If `/health` returns 500, check:
1. Database connectivity
2. Ent schema version mismatch (rebuild with `make build-all`)
3. Missing migrations

---

## Restarting After Schema Changes

**Never restart services blindly after ent schema edits.**

1. Regenerate ent code in BOTH directories:
   ```bash
   cd server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   cd herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   ```
2. Rebuild: `make build-all`
3. Stop: `make stop`
4. Start: `make dev`

---

## Production Deployment

The project deploys to Digital Ocean App Platform via `deploy.sh`.

1. Ensure `ghcr.io` token is valid: `docker login ghcr.io`
2. Run `make deploy` or `./deploy.sh`
3. Verify health: `curl https://your-app.ondigitalocean.app/health`

See `deploy/README.md` for platform-specific configuration.

---

## Next Steps

- **[Developer Guide](DEVELOPER-GUIDE/INDEX.md)** — build, test, contribute
- **[Installation & Upgrade](INSTALL.md)** — Docker Compose install, upgrades, migrations
