# HerbSt MUD — Production Deployment

One-droplet, zero-dependency production setup using Docker Compose + Cloudflare Tunnel.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Cloudflare Tunnel (HTTPS, free)                            │
│  game.yourdomain.com  → http://web-client:80               │
│  admin.yourdomain.com → http://admin:80                    │
└─────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            ▼                 ▼                 ▼
    ┌────────────┐   ┌────────────┐   ┌────────────┐
    │ web-client │   │   admin    │   │  mud-ssh   │
    │  (nginx)   │   │  (nginx)   │   │  (Go SSH)  │
    │   :80      │   │   :80      │   │   :4444    │
    └─────┬──────┘   └─────┬──────┘   └────────────┘
          │                │                │
          └────────┬───────┘                │
                   ▼                        │
            ┌────────────┐                  │
            │    web     │◄─────────────────┘
            │  (Go/Gin)  │    API_BASE_URL
            │   :8080    │
            └─────┬──────┘
                  ▼
            ┌────────────┐
            │  postgres  │
            │    :5432   │
            │  ./data/   │
            └────────────┘
```

## Quick Start (DigitalOcean)

### 1. Create a droplet
- **Image:** Ubuntu 22.04 (LTS)
- **Size:** 2 vCPU / 2 GB RAM minimum (1 GB will OOM during Go builds)
- **Region:** Closest to your players
- **SSH key:** Your usual key
- **Firewall:** Only open **22 (SSH)** and **4444 (MUD)** inbound

### 2. Copy the setup script

```bash
scp deploy/setup.sh root@DROPLET_IP:/root/
ssh root@DROPLET_IP 'bash /root/setup.sh'
```

The script will:
- Install Docker + Docker Compose
- Clone (or update) the repo to `/opt/herbst-mud`
- Check for `.env` — if missing, creates a template and exits so you fill it in
- Build all containers
- Start everything with restart policies

### 3. Fill in `.env`

First run will create a template and stop. Edit `/opt/herbst-mud/.env`:

```bash
ssh root@DROPLET_IP
nano /opt/herbst-mud/.env
```

| Variable | What |
|---|---|
| `DB_PASSWORD` | Postgres password (make it strong) |
| `JWT_SECRET` | Random 32+ char string for auth tokens |
| `CORS_ORIGINS` | Your public domains, comma-separated |
| `CF_TUNNEL_TOKEN` | From Cloudflare Zero Trust dashboard |

Then re-run:
```bash
bash /root/setup.sh
```

### 4. Configure Cloudflare Tunnel

1. Go to [Cloudflare Zero Trust](https://one.dash.cloudflare.com/) → Access → Tunnels
2. Create a tunnel, copy the token → paste into `.env`
3. Add public hostnames:
   - `game.yourdomain.com` → `http://web-client:80`
   - `admin.yourdomain.com` → `http://admin:80`
4. Update `CORS_ORIGINS` in `.env` to match your domains
5. `cd /opt/herbst-mud && docker compose -f docker-compose.prod.yml up -d`

### 5. Test

```bash
# API health
curl https://game.yourdomain.com/healthz

# SSH game
ssh -p 4444 game@yourdomain.com

# Browser client
open https://game.yourdomain.com
```

## Data Persistence

Postgres data lives in `./data/postgres` on the host. To migrate or back up:

```bash
# Backup
docker exec herbst-postgres pg_dump -U herbst herbst_mud > backup.sql

# Restore
docker exec -i herbst-postgres psql -U herbst herbst_mud < backup.sql
```

If you destroy and recreate the droplet:
1. Mount the old volume (or copy `data/`)
2. `docker compose -f docker-compose.prod.yml up -d`
3. Your world, characters, and abilities are intact

## Updating the Game

```bash
cd /opt/herbst-mud
git fetch origin main
git reset --hard origin/main
docker compose -f docker-compose.prod.yml build --no-cache
docker compose -f docker-compose.prod.yml up -d
```

## Useful Commands

```bash
# View all logs
docker compose -f docker-compose.prod.yml logs -f

# View one service
docker compose -f docker-compose.prod.yml logs -f web

# Restart one service
docker compose -f docker-compose.prod.yml restart web

# Full stop (data survives)
docker compose -f docker-compose.prod.yml down

# Wipe everything including data (DANGER)
docker compose -f docker-compose.prod.yml down -v
rm -rf data/
```

## Security Notes

- **No ports 80/443 exposed on the host** — Cloudflare Tunnel brings HTTPS in
- **SSH 4444 is the only exposed port** — consider restricting to your home IP via DO firewall
- **.env file is `chmod 600`** — never commit secrets
- **JWT_SECRET must be unique per deployment** — shared secret = account takeover
- **Postgres is internal-only** — not exposed to the internet

## Troubleshooting

| Symptom | Fix |
|---|---|
| `CF_TUNNEL_TOKEN required` | Fill `.env` with real token |
| `unhealthy` containers | Check `logs -f` — usually DB connection failure |
| CORS errors in browser | Update `CORS_ORIGINS` to match public domain |
| WebSocket fails | Ensure `/ws` route is in Cloudflare Tunnel (not just `/`) |
| Build OOMs | Resize droplet to 2 GB+ or add swap |
