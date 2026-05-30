# HerbSt MUD — Installation & Upgrade Guide

> **Audience:** Anyone deploying the stack via Docker Compose.
> **Goal:** Install, upgrade, and manage database migrations.

---

## Prerequisites

- Docker Engine 24.x+ and Docker Compose v2+
- 2 GB RAM minimum (4 GB recommended)
- Ports 4444, 8080, 3000, 3001 available (or remap in compose)

---

## 1. Fresh Installation

### 1.1 Clone and configure

```bash
git clone https://github.com/samrocksc/herbst-mud.git
cd herbst-mud
cp .env.example .env
```

Edit `.env`. Minimum required values:

```bash
# ── Database ──
DB_PASSWORD=CHANGE_ME_NOW
DB_NAME=herbst_mud

# ── Auth ──
JWT_SECRET=CHANGE_ME_NOW_AT_LEAST_32_CHARS

# ── Cloudflare Tunnel (optional, for public HTTPS) ──
# CF_TUNNEL_TOKEN=your_token_here
```

> **Do not** use the defaults in production. `DB_PASSWORD` and `JWT_SECRET` must be strong, unique values.

### 1.2 Choose your compose file

| File | Use case |
|------|----------|
| `docker-compose.yml` | Local development — builds from source, hot-reload friendly |
| `docker-compose.prod.yml` | Production — pulls pre-built images from GHCR, healthchecks, restart policies |

### 1.3 Start the stack (production)

```bash
docker compose -f docker-compose.prod.yml up -d
```

Verify everything is healthy:

```bash
docker compose -f docker-compose.prod.yml ps
```

Expected output:

```
NAME            SERVICE        STATUS    PORTS
herbst-api      web            running   0.0.0.0:8080->8080/tcp
herbst-ssh      mud-ssh        running   0.0.0.0:4444->4444/tcp
herbst-postgres postgres       running   0.0.0.0:5432->5432/tcp
herbst-web      web-client     running   80/tcp
herbst-admin    admin          running   80/tcp
```

### 1.4 First-time health checks

```bash
# API health
curl http://localhost:8080/healthz

# SSH MUD (interactive)
ssh -p 4444 player@localhost

# Admin panel
open http://localhost:3000

# Web client
open http://localhost:3001
```

---

## 2. First Login — Default Admin Account

On a fresh database, the API seeds a default admin user automatically:

| Field | Value |
|-------|-------|
| Email | `admin@herbstmud.local` |
| Password | `herb5t2026!` (stored in DB as plaintext — see note below) |
| Is admin | `true` |

**This account is created only if the `users` table is empty.** If you already have users, the seeder skips it.

### Reset the admin password via psql

The seeded password is stored as **plaintext**, but the login endpoint expects a **bcrypt hash**. You must reset it before logging in:

```bash
# Connect to the database
docker exec -it herbst-postgres psql -U herbst -d herbst_mud

# Generate a bcrypt hash for your desired password
# Quick one-liner (Python):
# python3 -c "import bcrypt; print(bcrypt.hashpw(b'sma', bcrypt.gensalt()).decode())"

# Update the admin password
UPDATE users
SET password = '$2b$12$YyYNBVcodA.Ax9d65xljQ.XTJLYlwCIozaw946jl6L/z303QDV9n6'
WHERE email = 'admin@herbstmud.local';

# Verify
\q
```

After the reset, log into the admin panel at `https://admin.yourdomain.com` with:
- **Email:** `admin@herbstmud.local`
- **Password:** `sma` (or whatever you generated a hash for)

> ⚠️ **Known issue:** The seed function does not bcrypt the password. This is tracked in `tickets/016-bug-multi-world-isolation-races-genders-tags.md` as a secondary finding. A proper fix requires adding `bcrypt.GenerateFromPassword` to `server/dbinit/init.go` (`InitAdminUser`) before the next release.

### Alternative: use the password-reset endpoint

If the API is running and you know the admin user's ID (usually `1`), you can reset it to `"password"` without touching the database:

```bash
# Reset admin password to "password"
curl -X POST http://localhost:8080/users/1/reset-password

# Then log in with admin@herbstmud.local / password
```

> ⚠️ **Security note:** The endpoint is not protected by auth. Anyone with network access to the API can reset any user's password. This is by design for admin recovery but should be firewalled in production.

---

## 3. Database Migrations

HerbSt uses **ent auto-migration**. On every startup, both the REST API server (`server/`) and the SSH MUD client (`herbst/`) automatically create or update tables to match the current schema.

### 2.1 How it works

The migration runs in `main.go` for both binaries:

```go
if err := client.Schema.Create(context.Background()); err != nil {
    log.Fatalf("failed creating schema resources: %v", err)
}
```

This is **idempotent** — safe to run on every restart. It:
- Creates missing tables
- Adds missing columns
- Adds missing indexes
- Does **not** drop columns or data

### 2.2 When migrations run automatically

| Event | Migration triggered? |
|-------|---------------------|
| `docker compose up` | Yes — on container start |
| `make dev` / `make dev-all` | Yes — on binary start |
| `make reload` / `make reload-web` | Yes — on binary restart |
| Schema code change without restart | No — stale binary in memory |

### 2.3 Manual migration (if needed)

If you need to force a migration check without restarting the full stack, run the API container's entrypoint directly:

```bash
docker compose -f docker-compose.prod.yml exec web /app/herbst-web
```

The binary exits after schema creation if configured to run migrations only. In practice, just restarting the `web` or `mud-ssh` service is sufficient:

```bash
docker compose -f docker-compose.prod.yml restart web
```

### 2.4 Migration safety

- **Always back up** before major version upgrades (see section 4).
- Auto-migration is additive-only. If you need to rename or drop a column, do it via a manual SQL script.
- Both `server/` and `herbst/` schemas must stay in sync. After editing `server/db/schema/`, regenerate BOTH:

```bash
cd server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
cd herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
```

Then rebuild and restart:

```bash
make build-all
make stop && make dev
```

---

## 3. Upgrading

### 3.1 Upgrade via GHCR images (production)

This is the standard path for a production deployment using `docker-compose.prod.yml`.

```bash
cd /path/to/herbst-mud

# 1. Pull latest images
docker compose -f docker-compose.prod.yml pull

# 2. Recreate containers (runs migrations on startup)
docker compose -f docker-compose.prod.yml up -d

# 3. Verify health
docker compose -f docker-compose.prod.yml ps
curl -s http://localhost:8080/healthz
```

> The `pull` + `up -d` cycle is the entire upgrade. New images contain the latest schema code; the containers auto-migrate on first boot.

### 3.2 Upgrade from source (development)

```bash
git fetch origin
git reset --hard origin/main

# Rebuild everything
make build-all

# Restart
docker compose down
docker compose up -d
```

### 3.3 One-command upgrade (Makefile)

For convenience, the Makefile includes:

```bash
make deploy-update
```

Which runs:

```bash
git fetch origin main
git reset --hard origin/main
docker compose -f docker-compose.prod.yml build --no-cache
docker compose -f docker-compose.prod.yml up -d
```

> **Note:** `deploy-update` is intended for the build host (where you have git + Docker). On a production server that only runs the compose file, use the GHCR pull method in section 3.1.

### 3.4 Rollback

If a new image fails, roll back to the previous tag:

```bash
# Edit docker-compose.prod.yml to pin the previous image tag
# e.g. image: ghcr.io/samrocksc/herbst-mud-api:v1.2.3

docker compose -f docker-compose.prod.yml up -d
```

Or, if you have the previous image locally:

```bash
docker compose -f docker-compose.prod.yml down
docker image tag ghcr.io/samrocksc/herbst-mud-api:previous ghcr.io/samrocksc/herbst-mud-api:latest
docker compose -f docker-compose.prod.yml up -d
```

---

## 4. Backup & Restore

### 4.1 Backup PostgreSQL

```bash
docker compose -f docker-compose.prod.yml exec postgres \
  pg_dump -U herbst -d herbst_mud > herbst_backup_$(date +%F_%H-%M-%S).sql
```

### 4.2 Restore PostgreSQL

```bash
# Drop and recreate the database (DESTRUCTIVE)
docker compose -f docker-compose.prod.yml exec postgres \
  psql -U herbst -c "DROP DATABASE herbst_mud; CREATE DATABASE herbst_mud;"

# Restore
docker compose -f docker-compose.prod.yml exec -T postgres \
  psql -U herbst -d herbst_mud < herbst_backup_YYYY-MM-DD_HH-MM-SS.sql
```

### 4.3 Backup data volume

```bash
# The prod compose mounts ./data/postgres on the host
tar czvf postgres_data_backup.tar.gz ./data/postgres
```

---

## 5. Troubleshooting

### Container stuck restarting

```bash
docker compose -f docker-compose.prod.yml logs web
docker compose -f docker-compose.prod.yml logs mud-ssh
```

Common causes:
- **Database unreachable** — check `DATABASE_URL` or `DB_HOST` (should be `postgres` inside Docker network)
- **Schema mismatch** — ent generated code is stale; run `ent generate` in both modules and rebuild
- **Port conflict** — something else is using 4444, 8080, or 5432

### Migration failures

If `Schema.Create` fails on startup, the binary exits and the container restarts. Check logs for:

```
failed creating schema resources: ...
```

Fix:
1. Ensure PostgreSQL container is healthy first
2. Check that `DB_PASSWORD` matches the postgres container's `POSTGRES_PASSWORD`
3. If the error mentions a missing relation or column, you likely skipped dual `ent generate`

### Data volume permissions

If the postgres container fails to start with permission errors:

```bash
sudo chown -R 999:999 ./data/postgres
```

(UID 999 is the postgres user inside the official image.)

---

## 6. Reference

| File | Purpose |
|------|---------|
| `.env.example` | Template for all environment variables |
| `docker-compose.yml` | Local dev stack (build from source) |
| `docker-compose.prod.yml` | Production stack (GHCR images) |
| `Makefile` | `deploy-up`, `deploy-down`, `deploy-update`, `deploy-logs` |
| `server/main.go` | Runs `client.Schema.Create` on startup |
| `herbst/main.go` | Runs `client.Schema.Create` on startup |
