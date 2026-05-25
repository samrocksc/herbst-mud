#!/bin/bash
#
# HerbSt MUD — Production Deploy Script (GHCR Edition)
# One-command drop onto a fresh DigitalOcean droplet (Ubuntu 22.04+).
#
# Prerequisites: .env file with secrets copied to the droplet first.
#
# Usage:
#   scp .env deploy/docker-compose.prod.yml root@DROPLET_IP:/opt/herbst-mud/
#   ssh root@DROPLET_IP 'bash /opt/herbst-mud/setup.sh'
#

set -euo pipefail

# ── Configuration ─────────────────────────────────────────────────
INSTALL_DIR="/opt/herbst-mud"
DATA_DIR="${INSTALL_DIR}/data/postgres"
COMPOSE_URL="https://raw.githubusercontent.com/samrocksc/herbst-mud/main/docker-compose.prod.yml"

# ── Helpers ─────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log()  { echo -e "${BLUE}[setup]${NC} $*"; }
ok()   { echo -e "${GREEN}[ok]${NC} $*"; }
warn() { echo -e "${YELLOW}[warn]${NC} $*"; }
fail() { echo -e "${RED}[fail]${NC} $*"; exit 1; }

# ── Preflight ───────────────────────────────────────────────────────
log "Preflight checks..."

[[ $EUID -eq 0 ]] || fail "Must run as root"

if ! command -v docker &>/dev/null; then
  log "Installing Docker..."
  apt-get update -qq
  apt-get install -y -qq ca-certificates curl gnupg lsb-release
  mkdir -p /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list
  apt-get update -qq
  apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-compose-plugin
  ok "Docker installed"
else
  ok "Docker already installed"
fi

if ! docker compose version &>/dev/null; then
  fail "docker compose plugin missing — installation failed?"
fi

# ── Directories ─────────────────────────────────────────────────────
mkdir -p "${INSTALL_DIR}" "${DATA_DIR}"
cd "${INSTALL_DIR}"

# ── Environment validation ──────────────────────────────────────────
ENV_FILE="${INSTALL_DIR}/.env"
if [[ ! -f "${ENV_FILE}" ]]; then
  fail "No .env file found at ${ENV_FILE}. Create it first and scp it over."
fi

# Source .env for validation (set +a to restore after)
set -a
source "${ENV_FILE}"
set +a

[[ "${DB_PASSWORD:-}" != "" ]]     || fail "DB_PASSWORD missing in .env"
[[ "${JWT_SECRET:-}" != "" ]]     || fail "JWT_SECRET missing in .env"
[[ "${CF_TUNNEL_TOKEN:-}" != "" ]] || fail "CF_TUNNEL_TOKEN missing in .env"

ok "Environment validated"

# ── Fetch compose file ──────────────────────────────────────────────
log "Fetching docker-compose.prod.yml..."
curl -fsSL "${COMPOSE_URL}" -o "${INSTALL_DIR}/docker-compose.prod.yml"
ok "Compose file downloaded"

# ── Pull and start ──────────────────────────────────────────────────
log "Pulling images from GHCR (this takes a few minutes)..."
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d

ok "All services started"

# ── Post-start health check ─────────────────────────────────────────
sleep 5
log "Health checks..."
if docker compose -f docker-compose.prod.yml ps | grep -q "unhealthy\|Restarting"; then
  warn "Some services may be unhealthy — check logs:"
  echo "  docker compose -f docker-compose.prod.yml logs -f"
else
  ok "All services healthy"
fi

# ── Summary ─────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  HerbSt MUD deployed and running (GHCR edition)       ${NC}"
echo -e "${GREEN}══════════════════════════════════════════════════════${NC}"
echo ""
echo "  Install dir: ${INSTALL_DIR}"
echo "  Data dir:    ${DATA_DIR}"
echo "  SSH port:    4444 (host)"
echo "  Images:      ghcr.io/samrocksc/herbst-mud-*:latest"
echo ""
echo "  Useful commands:"
echo "    cd ${INSTALL_DIR}"
echo "    docker compose -f docker-compose.prod.yml logs -f"
echo "    docker compose -f docker-compose.prod.yml down"
echo "    docker compose -f docker-compose.prod.yml up -d"
echo "    docker compose -f docker-compose.prod.yml pull && docker compose -f docker-compose.prod.yml up -d"
echo ""
echo "  Update (pull latest images):"
echo "    docker compose -f docker-compose.prod.yml pull"
echo "    docker compose -f docker-compose.prod.yml up -d"
echo ""
