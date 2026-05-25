title: HerbSt MUD — Build & Push to GHCR
category: infrastructure
description: Build all four herbst-mud service images and push to GitHub Container Registry.

---

## When to use

- After any code changes to SSH server, API, web-client, or admin
- Before deploying to a droplet (so production pulls pre-built images)
- When testing a fresh image build from a worktree or branch
- As a manual backup when CI is down

## Prerequisites

```bash
# Docker + Buildx installed
docker buildx ls

# Logged into GHCR with a PAT that has write:packages scope
echo "${GITHUB_TOKEN}" | docker login ghcr.io -u samrocksc --password-stdin

# Running from repo root (${HERBST_HOME}/herbst-mud)
cd /home/sam/GitHub/herbst-mud
```

## Quick commands

```bash
# Build + push all four images
make deploy-build-push

# Build one image only
bash scripts/build-push-ghcr.sh ssh
bash scripts/build-push-ghcr.sh api
bash scripts/build-push-ghcr.sh web
bash scripts/build-push-ghcr.sh admin

# Verify on GHCR
docker pull ghcr.io/samrocksc/herbst-mud-ssh:latest
docker images | grep ghcr.io/samrocksc/herbst-mud
```

## What gets pushed

| Image | Dockerfile | Context |
|-------|-----------|---------|
| `ghcr.io/samrocksc/herbst-mud-ssh` | `Dockerfile` | `.` |
| `ghcr.io/samrocksc/herbst-mud-api` | `server/Dockerfile` | `server/` |
| `ghcr.io/samrocksc/herbst-mud-web` | `web-client/Dockerfile` | `web-client/` |
| `ghcr.io/samrocksc/herbst-mud-admin` | `admin/Dockerfile` | `admin/` |

Tags: `latest` (movable) + `sha-abc1234` (immutable).

## Makefile targets

| Target | What it does |
|--------|--------------|
| `make deploy-build-push` | Build + push all four images |
| `make deploy-push-ssh` | Build + push only SSH |
| `make deploy-push-api` | Build + push only API |
| `make deploy-push-web` | Build + push only web client |
| `make deploy-push-admin` | Build + push only admin |

## Pitfall: GHCR authentication

If you get `denied: unauthenticated`:

1. Your .netrc PAT may be expired. Generate a new one at https://github.com/settings/tokens
2. The token needs `write:packages` + `read:packages` + `repo` scopes.
3. Store it somewhere safe (not in ~/.netrc). Docker credential helper preferred, or export as `GITHUB_TOKEN` env var per session.

## Pitfall: package visibility

Packages created on first push default to **Private**. After first push:
1. Go to https://github.com/samrocksc?tab=packages
2. Click each package → Package Settings → Change Visibility → Public
3. This is a one-time step per package.

## Pitfall: image tag collision

`latest` is overwritten on every push. If you need rollback safety, pin to the SHA tag:
```yaml
# In docker-compose.prod.yml
image: ghcr.io/samrocksc/herbst-mud-api:sha-abc1234
```

## Related

- `.github/workflows/docker.yml` — Same build steps, automated on GitHub Actions
- `scripts/build-push-ghcr.sh` — The script that performs the build + push
- `docker-compose.prod.yml` — Consumer of these images
- RFC-0003: GHCR Deployment Pipeline
