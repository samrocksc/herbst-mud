#!/bin/bash
#
# HerbSt MUD — Production Deploy Script
# One-command drop onto a fresh DigitalOcean droplet (Ubuntu 22.04+).
#
# Usage:
#   scp deploy/setup.sh root@DROPLET_IP:/root/
#   ssh root@DROPLET_IP 'bash /root/setup.sh'
#
# Or pipe directly:
#   curl -sL https://raw.githubusercontent.com/samrocksc/herbst-mud/main/deploy/setup.sh | bash
#

set -euo pipefail

# ── Configuration ─────────────────────────────────────────────────
REPO_URL="https://github.com/samrocksc/herbst-mud.git"
INSTALL_DIR="/opt/herbst-mud"
DATA_DIR="${INSTALL_DIR}/data/postgres"

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

# ── Clone / update repo ─────────────────────────────────────────────
if [[ -d "${INSTALL_DIR}/.git" ]]; then
  log "Updating existing repo at ${INSTALL_DIR}..."
  cd "${INSTALL_DIR}"
  git fetch origin main
  git reset --hard origin/main
else
  log "Cloning repo to ${INSTALL_DIR}..."
  git clone --depth 1 "${REPO_URL}" "${INSTALL_DIR}"
  cd "${INSTALL_DIR}"
fi
ok "Repo ready"

# ── Environment ─────────────────────────────────────────────────────
ENV_FILE="${INSTALL_DIR}/.env"
if [[ ! -f "${ENV_FILE}" ]]; then
  warn "No .env file found — creating template"
  cat > "${ENV_FILE}" <<'EOF'
# ── Database ──
DB_USER=herbst
DB_PASSWORD=CHANGE_ME_NOW
DB_NAME=herbst_mud

# ── Auth ──
JWT_SECRET=CHANGE_ME_NOW

# ── CORS (comma-separated public domains) ──
# Example: https://game.yourdomain.com,https://admin.yourdomain.com
CORS_ORIGINS=http://localhost

# ── Cloudflare Tunnel ──
# Get token from: https://one.dash.cloudflare.com/ → Access → Tunnels
CF_TUNNEL_TOKEN=your-token-here
EOF
  chmod 600 "${ENV_FILE}"
  fail "Please edit ${ENV_FILE} with real secrets, then re-run this script"
fi

# Validate required vars
source "${ENV_FILE}"
[[ "${DB_PASSWORD}" != "CHANGE_ME_NOW" ]] || fail "DB_PASSWORD is still default in .env"
[[ "${JWT_SECRET}" != "CHANGE_ME_NOW" ]]   || fail "JWT_SECRET is still default in .env"
[[ -n "${CF_TUNNEL_TOKEN}" ]]               || fail "CF_TUNNEL_TOKEN missing in .env"

ok "Environment validated"

# ── Data volume ─────────────────────────────────────────────────────
mkdir -p "${DATA_DIR}"
ok "Data directory: ${DATA_DIR}"

# ── Build & start ───────────────────────────────────────────────────
log "Building containers (this takes a few minutes)..."
cd "${INSTALL_DIR}"
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml build --no-cache
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
echo -e "${GREEN}  HerbSt MUD deployed and running                      ${NC}"
echo -e "${GREEN}══════════════════════════════════════════════════════${NC}"
echo ""
echo "  Install dir: ${INSTALL_DIR}"
echo "  Data dir:    ${DATA_DIR}"
echo "  SSH port:    4444 (host)"
echo ""
echo "  Useful commands:"
echo "    cd ${INSTALL_DIR}"
echo "    docker compose -f docker-compose.prod.yml logs -f"
echo "    docker compose -f docker-compose.prod.yml down"
echo "    docker compose -f docker-compose.prod.yml up -d"
echo ""
echo "  Next steps:"
echo "    1. Configure Cloudflare Tunnel routes in Zero Trust dashboard:"
echo "       game.yourdomain.com  → http://web-client:80"
echo "       admin.yourdomain.com → http://admin:80"
echo "    2. Update CORS_ORIGINS in .env to your public domains"
echo "    3. Test: ssh -p 4444 root@YOUR_DROPLET_IP"
echo ""
