# Deployment Guide

> 🔵 Last Updated: 2026-04-04

Complete deployment instructions for Herbst MUD on Digital Ocean + Neon DB.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Step-by-Step Deployment](#step-by-step-deployment)
- [Configuration Reference](#configuration-reference)
- [Docker Deployment](#docker-deployment)
- [Troubleshooting](#troubleshooting)

---

## Overview

Herbst MUD is a multi-service application consisting of:

| Service | Purpose | Port | Technology |
|---------|---------|------|------------|
| **SSH Server** | Player game client | 4444 | Go + Charmbracelet |
| **REST API** | Game state & logic | 8080 | Go + Gin |
| **Admin Panel** | Web administration | 80/3000 | React + Vite |
| **PostgreSQL** | Database | 5432 | Neon DB (managed) |

### Deployment Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   SSH Client    │────▶│   REST API       │────▶│   PostgreSQL    │
│   (herbst/)     │     │   (server/)      │     │   Database      │
│   Port 4444     │     │   Port 8080      │     │   (Neon DB)     │
└─────────────────┘     └──────────────────┘     └─────────────────┘
         │                        │
         │                        │
         ▼                        ▼
┌─────────────────┐     ┌──────────────────┐
│   Admin Panel   │────▶│   REST API       │
│   (admin/)      │     │   (same as above)│
│   Port 80       │     │                  │
└─────────────────┘     └──────────────────┘
```

---

## Prerequisites

### Required Accounts

1. **Digital Ocean Account** - [Sign up](https://www.digitalocean.com/)
2. **Neon DB Account** - [Sign up](https://neon.tech/)
3. **GitHub Account** - For repository connection

### Required Tools

```bash
# Install Digital Ocean CLI
curl -fsSL https://docs.digitalocean.com/reference/doctl/how-to/install/ | sudo sh

# Authenticate
doctl auth init
# Enter your Digital Ocean API token when prompted

# Verify
doctl account get
```

### Required Files

Ensure these files exist in your repository:

```
herbst-mud/
├── .do/app.yaml           # Digital Ocean App Platform spec
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # SSH Server Dockerfile
├── server/Dockerfile      # REST API Dockerfile
└── admin/Dockerfile       # Admin Panel Dockerfile
```

---

## Quick Start

For the impatient - one-command deployment once prerequisites are set:

```bash
# 1. Set environment variables
export DATABASE_URL="postgres://your_neon_connection_string"
export JWT_SECRET="$(openssl rand -base64 32)"
export CORS_ORIGINS="https://your-domain.com"

# 2. Deploy to Digital Ocean App Platform
doctl apps create --spec .do/app.yaml

# 3. Configure secrets via Digital Ocean UI
#    - Go to Apps → {your-app} → Settings → Components
#    - Set DATABASE_URL and JWT_SECRET as "Secret" type

# 4. Deploy SSH server to Droplet (see SSH Deployment section)
```

---

## Step-by-Step Deployment

### Step 1: Create Neon Database

1. Log in to [Neon Console](https://console.neon.tech/)
2. Create a new project (e.g., "herbst-mud")
3. Create a database (e.g., "herbst_mud")
4. Copy the **Connection String** with pooled connection:

```
postgres://[user]:[password]@[hostname-pooler]/[dbname]?sslmode=require
```

Example:
```
postgres://herbst_owner:abc123@ep-muddy-butterfly-p12abc-pooler.us-west-2.aws.neon.tech/herbst_mud?sslmode=require
```

⚠️ **Important:** Use the **Pooled Connection** endpoint (-pooler) for better performance.

---

### Step 2: Configure Digital Ocean App Platform

#### Option A: Using doctl (Command Line)

```bash
# Step 2a: Update .do/app.yaml with your repository
cat .do/app.yaml | sed 's/your-username/your-actual-github-username/g' > .do/app.yaml

# Step 2b: Create the app
doctl apps create --spec .do/app.yaml

# Output will include:
# - App ID
# - Default URL (e.g., https://herbst-mud-abc123.ondigitalocean.app)
```

#### Option B: Using Digital Ocean Console (UI)

1. Go to [Digital Ocean Apps](https://cloud.digitalocean.com/apps)
2. Click **Create App**
3. Connect your GitHub repository
4. Select the branch to deploy
5. Configure:
   - **Source Directory:** `/server` (for API), `/admin` (for admin)
   - **Dockerfile Path:** `Dockerfile`
   - **HTTP Port:** 8080 (API), 80 (admin)
6. Add environment variables as **Secrets**:
   - `DATABASE_URL` - Your Neon connection string
   - `JWT_SECRET` - Generate with `openssl rand -base64 32`
   - `CORS_ORIGINS` - Your admin panel domain
7. Click **Next** and deploy

---

### Step 3: Set Environment Secrets

After initial deployment, configure secrets:

```bash
# Get your app ID
APP_ID=$(doctl apps list --format ID --no-header | head -1)

# Update the spec with your secrets
cat > /tmp/env-vars.json << 'EOF'
{
  "envs": [
    {
      "key": "DATABASE_URL",
      "scope": "RUN_TIME",
      "type": "SECRET",
      "value": "YOUR_NEON_URL_HERE"
    },
    {
      "key": "JWT_SECRET",
      "scope": "RUN_TIME",
      "type": "SECRET",
      "value": "YOUR_JWT_SECRET_HERE"
    },
    {
      "key": "CORS_ORIGINS",
      "scope": "RUN_TIME",
      "value": "https://your-app.ondigitalocean.app"
    }
  ]
}
EOF

# Update via Digital Ocean UI - safer for secrets
```

**Via UI:**
1. Go to Apps → {your-app} → Settings
2. Find API service → Edit
3. Add environment variables as **Secret** type
4. Save and redeploy

---

### Step 4: Deploy SSH Server to Droplet

The SSH server cannot run on App Platform (HTTP-only). Deploy to a Droplet:

```bash
# Create a droplet
doctl compute droplet create herbst-mud-ssh \
  --region nyc1 \
  --size s-1vcpu-1gb \
  --image docker-20-04 \
  --ssh-keys <your-ssh-key-id> \
  --enable-monitoring \
  --tag "herbst-mud"

# Get the IP
DROPLET_IP=$(doctl compute droplet get herbst-mud-ssh --format PublicIPv4 --no-header)

# SSH into the droplet
ssh root@$DROPLET_IP
```

On the droplet:

```bash
# Clone the repository
git clone https://github.com/your-username/herbst-mud.git /opt/herbst-mud
cd /opt/herbst-mud

# Create environment file
cat > .env << 'EOF'
DATABASE_URL=postgres://your_neon_connection_string
API_BASE_URL=https://your-api-domain.ondigitalocean.app
SSH_PORT=4444
EOF

# Build and run with Docker
docker-compose up -d mud-ssh

# Or run directly (SSH server only)
cd /opt/herbst-mud
docker build -t herbst-mud-ssh .
docker run -d \
  --name herbst-ssh \
  --restart unless-stopped \
  -p 4444:4444 \
  --env-file .env \
  herbst-mud-ssh
```

---

### Step 5: Configure Firewall

Open the SSH port on the droplet:

```bash
# On the droplet
ufw allow 4444/tcp
ufw allow 22/tcp
ufw enable
```

---

### Step 6: Verify Deployment

Test all services:

```bash
# Test REST API
curl https://your-api-domain.ondigitalocean.app/healthz

# Test Admin Panel
open https://your-admin-domain.ondigitalocean.app

# Test SSH (from client)
ssh -p 4444 your-api-domain.ondigitalocean.app
# Or if using droplet IP
ssh -p 4444 your-droplet-ip
```

---

## Configuration Reference

### Environment Variables

#### Database Connection

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DATABASE_URL` | ✅ Yes* | Full Neon connection string | `postgres://user:pass@host/db?sslmode=require` |
| `DB_HOST` | ❌ No | Database host | `ep-xyz-pooler.neon.tech` |
| `DB_PORT` | ❌ No | Database port | `5432` |
| `DB_USER` | ❌ No | Database user | `herbst_owner` |
| `DB_PASSWORD` | ❌ No | Database password | `secret` |
| `DB_NAME` | ❌ No | Database name | `herbst_mud` |
| `DB_SSL_MODE` | ❌ No | SSL mode | `require` (always for Neon) |

\* `DATABASE_URL` is preferred. Individual variables are fallback.

#### Security

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `JWT_SECRET` | ✅ Yes | Secret for JWT signing | `openssl rand -base64 32` |
| `CORS_ORIGINS` | ⚠️ Prod | Allowed origins | `https://domain.com,https://admin.domain.com` |
| `RATE_LIMIT` | ❌ No | Requests per window | `100` |
| `RATE_WINDOW` | ❌ No | Time window (seconds) | `60` |

#### Server Configuration

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `SERVER_PORT` | ❌ No | API server port | `8080` |
| `SERVER_HOST` | ❌ No | API bind address | `0.0.0.0` |
| `SSH_PORT` | ❌ No | SSH server port | `4444` |
| `API_BASE_URL` | ✅ SSH | URL for REST API | `https://api.domain.com` |
| `VITE_API_BASE_URL` | ✅ Admin | API URL for admin | `https://api.domain.com` |

---

## Docker Deployment

### Local Development

```bash
# Clone repository
git clone https://github.com/your-username/herbst-mud.git
cd herbst-mud

# Set up environment
cp herbst/.env.example herbst/.env
cp server/.env.example server/.env
cp admin/.env.example admin/.env

# Edit .env files with your values
# For local dev, you can use the defaults

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production with Docker

```bash
# Build all images
docker-compose build

# Push to Digital Ocean Container Registry
doctl registry login
docker tag herbst-mud-mud-ssh registry.digitalocean.com/your-registry/herbst-ssh:latest
docker push registry.digitalocean.com/your-registry/herbst-ssh:latest

# On production server, pull and run
docker pull registry.digitalocean.com/your-registry/herbst-ssh:latest
docker run -d \
  --name herbst-ssh \
  -p 4444:4444 \
  -e DATABASE_URL="$DATABASE_URL" \
  -e API_BASE_URL="$API_BASE_URL" \
  registry.digitalocean.com/your-registry/herbst-ssh:latest
```

### Docker Compose with Neon DB

```yaml
# docker-compose.override.yml for Neon DB
version: '3.8'

services:
  mud-ssh:
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - API_BASE_URL=${API_BASE_URL}
      - SSH_PORT=4444
    # Remove postgres dependency for Neon
    depends_on: []

  web:
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
      - CORS_ORIGINS=${CORS_ORIGINS}
      - RATE_LIMIT=100
      - RATE_WINDOW=60
    # Remove postgres dependency for Neon
    depends_on: []

  # Disable local postgres when using Neon
  postgres:
    deploy:
      replicas: 0
```

Run with:
```bash
docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
```

---

## Troubleshooting

### Database Connection Issues

**Error:** `failed connecting to postgres: connection refused`

**Solutions:**
1. Verify `DATABASE_URL` is set correctly
2. Ensure `-pooler` endpoint is used for Neon
3. Check firewall rules allow outbound connections to port 5432
4. Verify `sslmode=require` is in the connection string

```bash
# Test connection from container
docker exec -it herbst-mud-mud-ssh-1 sh
# Install psql
apk add --no-cache postgresql-client
# Test connection
psql "${DATABASE_URL}" -c "SELECT 1;"
```

### JWT Authentication Failures

**Error:** `401 Unauthorized` or `signature is invalid`

**Solutions:**
1. Ensure `JWT_SECRET` is set and >= 32 characters
2. Verify SSH server and API server use the same secret
3. Check token hasn't expired (default: 24 hours)

### SSH Connection Refused

**Error:** `Connection refused` on port 4444

**Solutions:**
1. Verify droplet firewall allows port 4444
2. Check SSH service is running: `docker ps | grep herbst`
3. Verify no other service using port 4444: `netstat -tlnp | grep 4444`

### CORS Errors in Browser

**Error:** `CORS policy: No 'Access-Control-Allow-Origin' header`

**Solutions:**
1. Update `CORS_ORIGINS` to include your admin panel domain
2. Restart API server after changing CORS settings
3. For local dev: `CORS_ORIGINS=http://localhost:3000,http://localhost:5173`

### Rate Limiting (429 Errors)

**Error:** `429 Too Many Requests`

**Solutions:**
1. Adjust `RATE_LIMIT` and `RATE_WINDOW` env vars
2. For load testing: `RATE_LIMIT=10000 RATE_WINDOW=1`
3. In production, keep reasonable limits to prevent abuse

---

## Deployment Checklist

Before going live:

- [ ] Neon DB created and connection string tested
- [ ] `DATABASE_URL` configured as secret in Digital Ocean
- [ ] `JWT_SECRET` generated and configured (>= 32 chars)
- [ ] `CORS_ORIGINS` set to production domains
- [ ] SSH server deployed to Droplet with `DATABASE_URL`
- [ ] Firewall open on port 4444
- [ ] Health check endpoint responding: `/healthz`
- [ ] Admin panel loads and can log in
- [ ] SSH client connects and authenticates
- [ ] Rate limiting enabled for production
- [ ] Backups configured for Neon DB (automatic)

---

## Support

Having issues?

1. Check logs: `docker-compose logs -f` or `doctl apps logs <app-id>`
2. Verify environment variables: `printenv | grep -E "DB_|JWT_"`
3. Test database connection with `psql`
4. Review [CODE_ARCHITECTURE.md](../CODE_ARCHITECTURE.md) for architecture details

---

🔵 Document version: 2026-04-04
