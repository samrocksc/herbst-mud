# HerbSt MUD — Production Deployment

## Quick Deploy (GHCR Edition)

The droplet pulls pre-built images from GitHub Container Registry. No Go, no Node, no building on the server.

### 1. First-Time Droplet Setup

```bash
# On your machine — create .env and copy compose file
scp .env docker-compose.prod.yml root@DROPLET_IP:/opt/herbst-mud/

# On the droplet — run the setup script
ssh root@DROPLET_IP 'bash /opt/herbst-mud/deploy/setup.sh'
```

### 2. Update to Latest Images

```bash
ssh root@DROPLET_IP 'cd /opt/herbst-mud && docker compose pull && docker compose up -d'
```

### 3. Roll Back to a Specific SHA

```bash
# Edit docker-compose.prod.yml to use a pinned tag, then:
ssh root@DROPLET_IP 'cd /opt/herbst-mud && docker compose up -d'
```

## Required .env Variables

| Variable | Required? | What it does |
|----------|-----------|--------------|
| `DB_PASSWORD` | YES | Postgres password |
| `JWT_SECRET` | YES | HS256 signing key |
| `CF_TUNNEL_TOKEN` | YES | Cloudflare Tunnel ingress |
| `API_BASE_URL` | YES | Internal Docker URL: `http://web:8080` |
| `CORS_ORIGINS` | Recommended | Comma-separated public domains |
| `ADMIN_EMAIL` | Optional | Admin events account |
| `ADMIN_PASSWORD` | Optional | Admin events password (override hardcoded default) |

## Architecture

```
Internet → Cloudflare Tunnel (HTTPS) → nginx (web/admin) → Go API → Postgres
                        ↓
                    SSH port 4444
```

- **Images**: `ghcr.io/samrocksc/herbst-mud-{ssh,api,web,admin}:latest`
- **Registry**: Public GHCR (free, no auth needed to pull)
- **CI**: `.github/workflows/docker.yml` — builds + pushes on every push to `main`
- **Data**: `./data/postgres` on host volume (survives container recreation)

## GHCR Package Visibility

Packages are public by default (no cost, no auth). To verify:

1. Visit `https://github.com/samrocksc?tab=packages`
2. Set each package visibility to **Public** if not already

## Watchtower (Optional Auto-Pull)

Add to `docker-compose.prod.yml` for automatic updates:

```yaml
watchtower:
  image: containrrr/watchtower:latest
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  environment:
    - WATCHTOWER_POLL_INTERVAL=300
    - WATCHTOWER_CLEANUP=true
```

Then `docker compose up -d watchtower`. Every 5 minutes it checks GHCR for new `latest` images and restarts services.

**Risk**: Broken images auto-deploy. Pin to SHA tags for safety.
