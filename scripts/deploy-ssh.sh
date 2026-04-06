#!/bin/bash
#
# HerbSt MUD SSH Server Deployment Script
# For Digital Ocean Droplet
#
# Usage: ./deploy-ssh.sh <DROPLET_NAME> <REGION> <NEON_DATABASE_URL>
# Example: ./deploy-ssh.sh herbst-ssh nyc1 postgres://...neon.tech/herbst_mud

set -e

# Configuration
DROPLET_NAME=${1:-"herbst-ssh"}
REGION=${2:-"nyc1"}
SIZE=${3:-"s-1vcpu-1gb"}  
IMAGE=${4:-"docker-20-04"}  # Ubuntu 20.04 with Docker pre-installed
NEON_DB_URL=${5:-"$DATABASE_URL"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== HerbSt MUD SSH Server Deployment ===${NC}"
echo ""

# Validate inputs
if [ -z "$NEON_DB_URL" ]; then
    echo -e "${RED}Error: DATABASE_URL not set${NC}"
    echo "Please set DATABASE_URL environment variable or pass as argument"
    echo "Example: export DATABASE_URL=postgres://...neon.tech/herbst_mud"
    exit 1
fi

# Check for doctl
if ! command -v doctl &> /dev/null; then
    echo -e "${RED}Error: doctl not found${NC}"
    echo "Install: brew install doctl && doctl auth init"
    exit 1
fi

echo -e "${YELLOW}Creating Droplet: $DROPLET_NAME in $REGION${NC}"
echo "Size: $SIZE | Image: $IMAGE"
echo ""

# Create cloud-init script for initial setup
CLOUD_INIT=$(cat <<EOF
#cloud-config
package_update: true
packages:
  - docker-compose

runcmd:
  - mkdir -p /opt/herbst-mud
  - cd /opt/herbst-mud
  
  # Clone repository
  - git clone https://github.com/your-username/herbst-mud.git . || true
  
  # Create .env file
  - echo "DATABASE_URL=$NEON_DB_URL" > .env
  - echo "API_BASE_URL=http://herbst-mud-api.ondigitalocean.app" >> .env
  
  # Generate SSH host key (if not exists)
  - mkdir -p .ssh
  - ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N "" || true
  
  # Start SSH server
  - docker-compose up -d mud-ssh
  
  # Enable firewall (allow SSH on 4444)
  - ufw allow 4444/tcp
  - ufw allow 22/tcp
  - ufw --force enable

final_message: "HerbSt MUD SSH server setup complete!"
EOF
)

# Create Droplet
echo "Creating Droplet..."
DROPLET_ID=$(doctl compute droplet create "$DROPLET_NAME" \
    --region "$REGION" \
    --size "$SIZE" \
    --image "$IMAGE" \
    --ssh-keys ~/.ssh/id_rsa.pub \
    --user-data "$CLOUD_INIT" \
    --format ID \
    --no-header)

echo -e "${GREEN}Droplet created: $DROPLET_ID${NC}"
echo "Waiting for IP address (this may take 30-60 seconds)..."

# Wait for IP
sleep 10
IP=""
while [ -z "$IP" ]; do
    IP=$(doctl compute droplet get "$DROPLET_ID" --format PublicIPv4 --no-header 2>/dev/null || echo "")
    if [ -z "$IP" ]; then
        echo -n "."
        sleep 5
    fi
done

echo ""
echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo ""
echo -e "SSH Server IP: ${GREEN}$IP${NC}"
echo -e "SSH Port: ${GREEN}4444${NC}"
echo ""
echo "Connect with:"
echo -e "  ${YELLOW}ssh -p 4444 $IP${NC}"
echo ""
echo "Monitor logs:"
echo -e "  ${YELLOW}ssh root@$IP 'docker-compose logs -f mud-ssh'${NC}"
echo ""
echo "View droplet:"
echo -e "  ${YELLOW}https://cloud.digitalocean.com/droplets/$DROPLET_ID${NC}"

# Save connection info
cat > .ssh-droplet-info.txt <<EOF
HerbSt MUD SSH Server
=====================
Droplet Name: $DROPLET_NAME
Droplet ID: $DROPLET_ID
IP Address: $IP
Region: $REGION
Port: 4444
Connect: ssh -p 4444 $IP
EOF

echo -e "${GREEN}Connection info saved to: .ssh-droplet-info.txt${NC}"
