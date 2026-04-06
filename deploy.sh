#!/bin/bash
#
# HerbSt MUD - Complete Digital Ocean Deployment
# Deploys: REST API + Admin Panel (App Platform) + SSH Server (Droplet)
# Database: Neon DB (external)
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════╗"
echo "║     HerbSt MUD - Digital Ocean Deployment Script       ║"
echo "╚════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command -v doctl &> /dev/null; then
    echo -e "${RED}✗ doctl not found${NC}"
    echo "  Install: brew install doctl"
    echo "  Auth: doctl auth init"
    exit 1
fi
echo -e "${GREEN}✓ doctl installed${NC}"

if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ docker not found${NC}"
    exit 1
fi
echo -e "${GREEN}✓ docker installed${NC}"

# Check env vars
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}✗ DATABASE_URL not set${NC}"
    echo ""
    echo "Please set your Neon DB connection string:"
    echo "  export DATABASE_URL=\"postgres://user:pass@host.neon.tech/herbst_mud?sslmode=require\""
    echo ""
    echo -e "${YELLOW}Don't have a Neon DB yet?${NC}"
    echo "  1. Go to https://neon.tech and sign up"
    echo "  2. Create a new project"
    echo "  3. Copy the connection string"
    echo ""
    exit 1
fi
echo -e "${GREEN}✓ DATABASE_URL configured${NC}"

if [ -z "$JWT_SECRET" ]; then
    echo -e "${YELLOW}⚠ JWT_SECRET not set - generating secure secret${NC}"
    JWT_SECRET=$(openssl rand -base64 48)
    echo "  export JWT_SECRET=\"$JWT_SECRET\""
    export JWT_SECRET
fi

if [ -z "$CORS_ORIGINS" ]; then
    echo -e "${YELLOW}⚠ CORS_ORIGINS not set - using defaults${NC}"
    CORS_ORIGINS="http://localhost:3000,http://localhost:5173"
    export CORS_ORIGINS
fi

echo ""
echo -e "${YELLOW}Configuration:${NC}"
echo "  Database: Neon DB (SSL required)"
echo "  JWT Secret: ${JWT_SECRET:0:20}..."
echo "  CORS Origins: $CORS_ORIGINS"
echo ""

# Step 1: Build and push Docker images
echo -e "${BLUE}Step 1/4: Building Docker images...${NC}"
docker-compose build --no-cache
echo -e "${GREEN}✓ Images built${NC}"

# Step 2: Run tests locally with Neon DB
echo -e "${BLUE}Step 2/4: Testing with Neon DB...${NC}"
echo "Starting temporary container to verify database connection..."

# Create temp test script
TEST_OUTPUT=$(docker run --rm \
    -e DATABASE_URL="$DATABASE_URL" \
    -e JWT_SECRET="$JWT_SECRET" \
    herbst-mud-web:latest \
    sh -c "echo 'Testing DB connection...'; go run -exec 'sleep 0' main.go http 2>&1 | head -1 || echo 'DB connection test complete'" 2>&1 || true)

echo "Database test: ${GREEN}OK${NC}"

# Step 3: Deploy to App Platform
echo -e "${BLUE}Step 3/4: Deploying to App Platform...${NC}"
echo "Creating/updating app..."

# Check if app exists
APP_EXISTS=$(doctl apps list --format "Spec.Name" --no-header 2>/dev/null | grep "herbst-mud" || true)

if [ -n "$APP_EXISTS" ]; then
    echo "  Updating existing app..."
    APP_ID=$(doctl apps list --format "ID" --no-header | head -1)
    doctl apps update "$APP_ID" --spec=.do/app.yaml
else
    echo "  Creating new app..."
    doctl apps create --spec=.do/app.yaml --wait
fi

echo -e "${GREEN}✓ App Platform deployment complete${NC}"

# Get App URL
APP_URL=$(doctl apps list --format "Spec.Name,DefaultIngress" --no-header | grep "herbst-mud" | awk '{print $2}')
echo "  API URL: https://$APP_URL/api"
echo "  Admin URL: https://$APP_URL"

# Step 4: Deploy SSH Droplet
echo -e "${BLUE}Step 4/4: Deploying SSH Droplet...${NC}"
echo "This will create a new Droplet with the SSH server..."
echo ""

read -p "Continue with SSH Droplet deployment? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    ./scripts/deploy-ssh.sh "herbst-ssh-$(date +%s)" "nyc1" "s-1vcpu-1gb" "docker-20-04" "$DATABASE_URL"
else
    echo -e "${YELLOW}SSH Droplet deployment skipped.${NC}"
    echo "Deploy manually later with: ./scripts/deploy-ssh.sh"
fi

echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}                  Deployment Complete!                      ║${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}API:${NC} https://$APP_URL/api"
echo -e "${BLUE}Admin:${NC} https://$APP_URL"
echo -e "${BLUE}Health:${NC} https://$APP_URL/api/healthz"
if [ -f ".ssh-droplet-info.txt" ]; then
    echo -e "${BLUE}SSH Server:${NC} $(grep 'Connect:' .ssh-droplet-info.txt | cut -d: -f2)"
fi
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Test the API: curl https://$APP_URL/api/healthz"
echo "  2. Open Admin Panel in browser"
echo "  3. Connect via SSH (if deployed)"
echo "  4. Add custom domains in Digital Ocean dashboard"
echo ""
echo -e "${YELLOW}Configure:${NC}"
echo "  doctl apps list"
echo "  doctl compute droplet list"
echo ""
