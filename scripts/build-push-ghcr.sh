#!/usr/bin/env bash
#
# HerbSt MUD — Build & Push to GHCR
# Builds one or all service images and pushes to GitHub Container Registry.
#
# Usage:
#   ./scripts/build-push-ghcr.sh              Build + push all four images
#   ./scripts/build-push-ghcr.sh ssh          Build + push only SSH image
#   ./scripts/build-push-ghcr.sh api          Build + push only API image
#   ./scripts/build-push-ghcr.sh web          Build + push only web client
#   ./scripts/build-push-ghcr.sh admin        Build + push only admin panel
#
# Prerequisites:
#   - docker login ghcr.io (with a token that has write:packages scope)
#   - Running from repo root
#

set -euo pipefail

REPO_OWNER="samrocksc"
IMAGE_PREFIX="ghcr.io/${REPO_OWNER}/herbst-mud"

# Compute tags
GIT_SHA=$(git rev-parse --short HEAD)
LATEST_TAG="latest"
SHA_TAG="sha-${GIT_SHA}"

# Colours
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log()  { echo -e "${BLUE}[ghcr]${NC} $*"; }
ok()   { echo -e "${GREEN}[ok]${NC} $*"; }
warn() { echo -e "${YELLOW}[warn]${NC} $*"; }
fail() { echo -e "${RED}[fail]${NC} $*"; exit 1; }

# ── Preflight ─────────────────────────────────────────────────────────
log "Preflight..."

# Running from repo root?
[[ -f "docker-compose.prod.yml" ]] || fail "Run this script from the herbst-mud repo root"

# Docker available?
command -v docker >/dev/null || fail "docker not found"

# Logged into GHCR?
if ! docker system info 2>/dev/null | grep -qi "ghcr.io\|github" >/dev/null 2>&1; then
  # Try a lightweight pull check
  if ! docker pull "${IMAGE_PREFIX}-ssh:latest" >/dev/null 2>&1; then
    warn "May not be authenticated with ghcr.io"
    warn "Run: echo \$GITHUB_TOKEN | docker login ghcr.io -u ${REPO_OWNER} --password-stdin"
  fi
fi

ok "Preflight passed"

# ── Build & push helper ────────────────────────────────────────────────
build_push() {
  local service="$1"
  local context="$2"
  local dockerfile="$3"
  local image="${IMAGE_PREFIX}-${service}"

  log "Building ${service}..."
  log "  Context: ${context}"
  log "  Dockerfile: ${dockerfile}"

  docker buildx build \
    --platform linux/amd64 \
    --file "${dockerfile}" \
    --tag "${image}:${LATEST_TAG}" \
    --tag "${image}:${SHA_TAG}" \
    --push \
    "${context}"

  ok "Pushed ${image}:${LATEST_TAG} + ${SHA_TAG}"
}

# ── Parse args ──────────────────────────────────────────────────────
TARGET="${1:-all}"

log "Target: ${TARGET}"
log "Registry: ${IMAGE_PREFIX}"
log "Tags: ${LATEST_TAG}, ${SHA_TAG}"

case "${TARGET}" in
  all)
    build_push "ssh" "." "Dockerfile"
    build_push "api" "./server" "./server/Dockerfile"
    build_push "web" "./web-client" "./web-client/Dockerfile"
    build_push "admin" "./admin" "./admin/Dockerfile"
    ;;
  ssh)
    build_push "ssh" "." "Dockerfile"
    ;;
  api)
    build_push "api" "./server" "./server/Dockerfile"
    ;;
  web|web-client|client)
    build_push "web" "./web-client" "./web-client/Dockerfile"
    ;;
  admin)
    build_push "admin" "./admin" "./admin/Dockerfile"
    ;;
  *)
    fail "Unknown target: ${TARGET}. Use: all, ssh, api, web, admin"
    ;;
esac

# ── Summary ─────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  All images pushed to GHCR                            ${NC}"
echo -e "${GREEN}══════════════════════════════════════════════════════${NC}"
echo ""
echo "  Registry:  ${IMAGE_PREFIX}"
echo "  Tags:      ${LATEST_TAG}, ${SHA_TAG}"
echo ""
echo "  Verify:"
echo "    docker pull ${IMAGE_PREFIX}-ssh:${SHA_TAG}"
echo "    docker pull ${IMAGE_PREFIX}-api:${SHA_TAG}"
echo "    docker pull ${IMAGE_PREFIX}-web:${SHA_TAG}"
echo "    docker pull ${IMAGE_PREFIX}-admin:${SHA_TAG}"
echo ""
