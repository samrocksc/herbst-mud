---
project: herbst-mud
type: skill
scope: infrastructure
---

# HerbSt MUD — Local Build & Push to GHCR

One-stop script for building all four service images locally and pushing them to GitHub Container Registry.

## Prerequisites

- Docker + Buildx installed
- GitHub Personal Access Token with `write:packages` scope: `GITHUB_TOKEN` env var
  - Generate at: https://github.com/settings/tokens → `repo` + `write:packages`
- Docker logged into GHCR: `echo "${GITHUB_TOKEN}" | docker login ghcr.io -u USERNAME --password-stdin`

## Quick Start

```bash
cd /home/sam/GitHub/herbst-mud
# Set token (one-time per session)
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxx

# Login
echo "${GITHUB_TOKEN}" | docker login ghcr.io -u samrocksc --password-stdin

# Build + push all four images
make deploy-build-push
# OR run the script directly:
bash scripts/build-push-ghcr.sh
```

## What gets built

| Image | Dockerfile | Context | Service |
|-------|-----------|---------|---------|
| `ghcr.io/samrocksc/herbst-mud-ssh` | `./Dockerfile` | `.` | SSH MUD server |
| `ghcr.io/samrocksc/herbst-mud-api` | `./server/Dockerfile` | `./server` | REST API |
| `ghcr.io/samrocksc/herbst-mud-web` | `./web-client/Dockerfile` | `./web-client` | Browser client |
| `ghcr.io/samrocksc/herbst-mud-admin` | `./admin/Dockerfile` | `./admin` | Admin panel |

## Tags pushed

| Tag | When |
|-----|------|
| `latest` | Every build |
| `sha-$(git rev-parse --short HEAD)` | Every build (immutable traceability) |
| `$(git describe --tags --always)` | Only if HEAD is tagged |

## One-off image push

Push a single service without rebuilding all four:

```bash
./scripts/build-push-ghcr.sh ssh      # Only SSH image
./scripts/build-push-ghcr.sh api      # Only API image
./scripts/build-push-ghcr.sh web      # Only web client
./scripts/build-push-ghcr.sh admin    # Only admin panel
```

## Verify images on GHCR

```bash
# List packages in your account
curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
  https://api.github.com/users/samrocksc/packages?package_type=container

# Check image exists
docker pull ghcr.io/samrocksc/herbst-mud-ssh:latest
```

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `denied: unauthenticated: unauthenticated` | Re-run `docker login ghcr.io` with a fresh token |
| `ERROR: denied: requested access to the resource is denied` | Token needs `write:packages` scope |
| Image not showing in GHCR Packages tab | Package visibility defaults to Private; toggle to Public in GitHub UI |
| Build fails on `COPY go.mod` | Run `go mod tidy` in the build context first |
| Build fails on `npm ci` | Delete `node_modules` + `package-lock.json` and reinstall |

## Related

- `.github/workflows/docker.yml` — CI workflow (same build steps, runs on GitHub Actions)
- `docker-compose.prod.yml` — Uses the images this script pushes
- `deploy/setup.sh` — Droplet script that pulls these images

---
