# RFC-007: Digital Ocean Deployment for herbst-mud

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-09
**Related:** Engine Audit (turtle-wiki), docker-compose.yml

---

## 1. Executive Summary

This RFC covers deploying herbst-mud as a Docker Compose stack on a single
Digital Ocean droplet. No managed services — just Docker on a VM with
Cloudflare Tunnel for routing. The primary concerns are: (1) database
persistence outside the Docker volume, (2) reliable backup strategy, and
(3) JSON export reliability for disaster recovery.

---

## 2. Current State Assessment

### 2.1 What Works

| Component | Status |
|-----------|--------|
| Docker Compose (5 services) | ✅ Working locally |
| REST API (Go/Gin, port 8080) | ✅ Builds and serves |
| SSH Game Server (port 4444) | ✅ Builds and serves |
| Admin UI (React/nginx, port 80) | ✅ Builds and serves |
| Postgres 15 | ✅ Running with named volume |
| Cloudflare Tunnel | ✅ Configured in compose |
| In-app backup system | ✅ Export + restore with checksums |
| JSON game export | ✅ 7 entity types, manifest, validation |
| Makefile (build, start, stop, dev) | ✅ Complete |

### 2.2 What Needs Fixing Before Production

| Concern | Severity | Issue |
|---------|----------|-------|
| **DB persistence** | 🔴 Critical | Postgres uses Docker named volume (`postgres_data`). If the VM is destroyed or the volume is pruned, ALL player data is lost. Named volumes are tied to Docker, not the host filesystem. |
| **Backup location** | 🔴 Critical | Backups written to `./backups/` inside the web container — ephemeral. Lost on container rebuild/restart. Not accessible from host. |
| **JWT secret** | 🟡 High | Default `development-jwt-secret-change-in-production` — must be a strong random value per deployment |
| **SSH host key** | 🟡 High | Generated at runtime inside container — lost on rebuild unless mounted |
| **JSON export completeness** | 🟡 Medium | Exports 7 entity types but newer schemas (applogs, factions, competencies, character_tags, achievements, dialogs) are NOT included in export/backup |
| **No health check orchestration** | 🟡 Medium | `depends_on` exists but no `healthcheck` blocks — services may start before Postgres is ready |

### 2.3 In-App Backup System Coverage

The `server/backup/` package exports these entities:
- users, rooms, abilities, npc_templates, equipment, characters, character_abilities

**Missing from backup (schemas that exist but are NOT exported):**
- achievements (schema exists, no backup entry)
- factions + character_factions (schema exists, no backup entry)
- character_tags (schema exists, no backup entry)
- competencies + character_competencies (schema exists, no backup entry)
- applogs (schema exists, no backup entry)
- factions, faction_categories, faction_required_tags
- game_config
- races, genders

This means `POST /api/backups/:id/restore` WILL NOT restore these entities.

---

## 3. Deployment Architecture

```
┌────────────────────────────────────────────────────────────┐
│                  Digital Ocean Droplet                      │
│  (Ubuntu 24.04 LTS, 2 vCPU, 4GB RAM, 80GB SSD — ~$24/mo)  │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              /opt/herbst-mud/                        │  │
│  │                                                     │  │
│  │  docker-compose.yml                                 │  │
│  │  .env          ← secrets (JWT, CF token, DB creds) │  │
│  │  data/                                              │  │
│  │    postgres/    ← bind-mount (persistent DB)        │  │
│  │    backups/     ← bind-mount (off-container)        │  │
│  │    ssh-keys/    ← bind-mount (persistent host key)  │  │
│  │    logs/        ← bind-mount (app logs)             │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
│  Docker Network: herbst-network (bridge)                    │
│  ┌─────────┐  ┌─────────┐  ┌──────────┐  ┌─────────────┐  │
│  │mud-ssh  │  │  web    │  │  admin   │  │  postgres   │  │
│  │:4444    │  │ :8080   │  │  :80     │  │  :5432      │  │
│  │(herbst/)│  │(server/)│  │(nginx)   │  │(pg15)       │  │
│  └─────────┘  └─────────┘  └──────────┘  └─────────────┘  │
│                                                     │      │
│  ┌─────────────────────────────────────────────┐   │      │
│  │  cloudflare-tunnel (cloudflared)            │   │      │
│  │  Routes: herbstmud.com → admin:80           │   │      │
│  │          ssh.herbstmud.com → mud-ssh:4444    │   │      │
│  └─────────────────────────────────────────────┘   │      │
│                                                     │      │
│  UFW Firewall:                                      │      │
│    Allow: 22/tcp (SSH for admin)                    │      │
│    Allow: 80/tcp, 443/tcp (if not using CF Tunnel)  │      │
│    Deny:  4444, 8080, 5432 (block direct access)    │      │
└─────────────────────────────────────────────────────┘
```

---

## 4. Required Changes

### 4.1 docker-compose.yml — Production Overrides

```yaml
# Production docker-compose.yml changes
services:
  postgres:
    volumes:
      # REPLACE named volume with bind mount
      - ./data/postgres:/var/lib/postgresql/data

  mud-ssh:
    volumes:
      # Persist SSH host key across rebuilds
      - ./data/ssh-keys:/root/.ssh

  web:
    volumes:
      # Backups accessible from host
      - ./data/backups:/root/backups
      # Logs accessible from host
      - ./data/logs:/root/logs
    environment:
      - JWT_SECRET=${JWT_SECRET}  # must be set in .env
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-herbst}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
```

### 4.2 New File: `.env.example`

```
# Postgres
DB_USER=herbst
DB_PASSWORD=<generate 32-char random>
DB_NAME=herbst_mud

# JWT
JWT_SECRET=<generate 64-char random>

# Cloudflare Tunnel
CF_TUNNEL_TOKEN=<from Cloudflare Zero Trust dashboard>

# CORS (comma-separated origins)
CORS_ORIGINS=https://herbstmud.com

# Logging
LOG_MIN_LEVEL=INFO
LOG_RETENTION_DAYS=7
```

### 4.3 Fix Backup Coverage

Add missing schemas to `server/backup/export/`:

```
New exporters needed:
  - achievements
  - factions + character_factions
  - character_tags
  - competencies + character_competencies
  - game_config
  - races
  - genders
```

And corresponding `server/backup/restore/` importers.

### 4.4 External Backup Script

Create `scripts/backup-db.sh` (runs on host via cron):

```bash
#!/bin/bash
# Runs daily on the host via cron: 0 2 * * * /opt/herbst-mud/scripts/backup-db.sh

BACKUP_DIR="/opt/herbst-mud/data/backups/pg_dump"
RETENTION_DAYS=30
mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)

# pg_dump from the running container
docker compose exec -T postgres pg_dump -U herbst herbst_mud \
  > "$BACKUP_DIR/herbst_mud_$TIMESTAMP.sql"

# Compress
gzip "$BACKUP_DIR/herbst_mud_$TIMESTAMP.sql"

# Prune old backups
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
```

This gives you two backup layers:
1. **In-app JSON export** (app-level, entity-aware, id-mapping for restore)
2. **pg_dump SQL** (database-level, complete, standard tooling)

### 4.5 Admin Dockerfile — `/api/backups` route

The nginx config is missing a location block for `/api/backups`:

```nginx
location /api/backups {
    proxy_pass http://web:8080;
    ...
}
```

(This route exists in the Go code but isn't proxied through nginx.)

---

## 5. Deployment Checklist

### Step 1: Provision Droplet

- Ubuntu 24.04 LTS
- 2 vCPU / 4GB RAM / 80GB SSD (~$24/mo at DO)
- Add SSH key during creation
- Enable monitoring (free)

### Step 2: Initial Server Setup

```bash
ssh root@<droplet-ip>

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
systemctl enable docker

# Add deploy user
useradd -m -s /bin/bash deploy
usermod -aG docker deploy

# Firewall
ufw allow 22/tcp
ufw allow 80/tcp
# If Cloudflare Tunnel handles ingress, you can deny 80/443 too
ufw enable
```

### Step 3: Clone & Configure

```bash
su - deploy
mkdir -p /opt/herbst-mud/data/{postgres,backups,ssh-keys,logs}

cd /opt/herbst-mud
git clone https://github.com/samrocksc/herbst-mud.git .

# Create .env from .env.example with real secrets
cp .env.example .env
# Edit .env — generate JWT_SECRET and DB_PASSWORD

# Generate SSH host key
ssh-keygen -t ed25519 -f data/ssh-keys/term_info_ed25519 -N ""
```

### Step 4: Start Stack

```bash
cd /opt/herbst-mud
docker compose up -d

# Verify
docker compose ps
curl http://localhost:8080/healthz
# → {"db":"connected","ssh":"running","status":"ok"}
```

### Step 5: Configure Cloudflare Tunnel

- In Cloudflare Zero Trust dashboard: create tunnel, point to `http://admin:80`
- Add public hostname: `herbstmud.com` → `http://admin:80`
- Add public hostname: `ssh.herbstmud.com` → `ssh://mud-ssh:4444` (if needed for SSH over browser via Cloudflare)
- Copy tunnel token to `.env` as `CF_TUNNEL_TOKEN`
- Restart: `docker compose up -d`

### Step 6: Daily Backup Cron (on host)

```bash
# As root or deploy user
crontab -e
# Add:
0 2 * * * /opt/herbst-mud/scripts/backup-db.sh >> /var/log/herbst-backup.log 2>&1
```

---

## 6. Cost Estimate

| Resource | Monthly |
|----------|---------|
| DO Droplet (2 vCPU, 4GB, 80GB) | ~$24 |
| Cloudflare Tunnel | Free |
| Domain (herbstmud.com) | ~$12/year |

**Total: ~$25/month**

Can scale to 4GB/2vCPU at ~$48/mo if needed, or add a managed Postgres ($15/mo) to offload DB.

---

## 7. Future Scaling

### Single Droplet (this RFC)
- Everything on one VM
- Good for: <50 concurrent players
- Risk: SPOF (single point of failure)

### Vertical Scale (6 months+)
- Bump to 4 vCPU / 8GB RAM (~$48/mo)
- Move Postgres to DO Managed Database ($15/mo)
- Add a second web container behind nginx load balancer

### Horizontal Scale (1 year+)
- Separate droplets per service (API, SSH, DB)
- DO Load Balancer ($10/mo)
- Postgres read replicas for admin queries

---

## 8. Open Questions

1. **Drop-in replacement for Cloudflare Tunnel?** If you ever want to drop CF, you'd need nginx or Caddy on the host for TLS termination + reverse proxy. Cloudflare Tunnel handles TLS automatically and hides the origin IP. Worth keeping unless there's a specific reason to ditch it.

2. **SSH game server directly exposed?** Currently the SSH server on 4444 is how players connect. If Cloudflare Tunnel's SSH support works, players could `ssh ssh.herbstmud.com` through CF's edge. Otherwise, you'd need to expose port 4444 directly (or use a non-standard port).

3. **Admin UI auth on public internet?** The admin panel at `herbstmud.com` has a login page but no rate limiting on `/users/auth`. Consider adding fail2ban or increasing the rate limiter sensitivity in production.

4. **Monitoring?** DO's free monitoring covers CPU/memory/disk. For app-level monitoring, consider adding a `/metrics` endpoint later. The applogs table already captures runtime events.

---

## 9. Acceptance Criteria

- [ ] Postgres data bind-mounted to `./data/postgres/` (not Docker volume)
- [ ] Backups written to `./data/backups/` (host-accessible)
- [ ] SSH host key mounted from host (`./data/ssh-keys/`)
- [ ] `.env.example` committed with all required vars
- [ ] Backup system covers ALL entity types (not just 7 of 14+)
- [ ] `scripts/backup-db.sh` for pg_dump cron job
- [ ] Health checks on web and postgres services
- [ ] UFW firewall rules documented
- [ ] Admin nginx config proxies `/api/backups` and `/api/logs/stream` correctly
- [ ] `docker compose up -d` starts all 5 services clean on fresh DO droplet
