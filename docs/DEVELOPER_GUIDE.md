# Developer Guide

> 🔵 Last Updated: 2026-04-04

A comprehensive guide for developers working on the Herbst MUD codebase.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Database Architecture](#database-architecture)
- [Security Configuration](#security-configuration)
- [API Development](#api-development)
- [SSH Server Development](#ssh-server-development)
- [Admin Panel Development](#admin-panel-development)
- [Testing](#testing)
- [Deployment Pipeline](#deployment-pipeline)

---

## Architecture Overview

Herbst MUD uses a **microservices architecture** with three main services:

```
┌─────────────────────────────────────────────────────────────────────┐
│                         HERBST MUD                                  │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐            │
│  │ SSH Server   │   │ REST API     │   │ Admin Panel  │            │
│  │ (herbst/)    │   │ (server/)    │   │ (admin/)     │            │
│  │              │   │              │   │              │            │
│  │ • TUI Game   │   │ • REST       │   │ • React SPA  │            │
│  │ • Commands   │   │ • JWT Auth   │   │ • Management │            │
│  │ • Combat     │   │ • OpenAPI    │   │ • CRUD       │            │
│  │              │   │              │   │              │            │
│  └──────┬───────┘   └──────┬───────┘   └──────▲───────┘            │
│         │                  │                  │                   │
│         └──────────────────┼──────────────────┘                   │
│                            │                                      │
│                     ┌──────────────┐                             │
│                     │ PostgreSQL   │                             │
│                     │ (Neon DB)    │                             │
│                     └──────────────┘                             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
herbst-mud/
│
├── herbst/                    # SSH Server (Go)
│   ├── main.go               # Entry point, SSH server setup
│   ├── model.go              # Bubble Tea model + UI state
│   ├── auth.go               # JWT authentication
│   ├── cmd_*.go              # Command handlers
│   ├── combat*.go            # Combat system
│   ├── db/                   # Ent ORM client
│   └── .env.example          # Environment template
│
├── server/                    # REST API (Go)
│   ├── main.go               # Entry point, Gin router setup
│   ├── routes/               # HTTP route handlers
│   ├── middleware/           # Auth middleware
│   ├── db/                   # Ent schema + client
│   ├── dbinit/               # Data initialization
│   └── .env.example          # Environment template
│
├── admin/                     # Admin Panel (React + Vite)
│   ├── src/
│   │   ├── routes/           # TanStack Router routes
│   │   ├── client/           # Generated API client
│   │   └── components/       # React components
│   ├── Dockerfile
│   └── .env.example
│
├── docs/                       # Documentation
├── features/                   # Gherkin feature files
├── .do/app.yaml              # Digital Ocean spec
└── docker-compose.yml        # Local orchestration
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose (optional)
- PostgreSQL 15+ (or Neon DB)

### Local Setup

```bash
# Clone repository
git clone https://github.com/your-username/herbst-mud.git
cd herbst-mud

# Option 1: Docker (Recommended for quick start)
docker-compose up -d

# Option 2: Native Development

# Terminal 1: Start PostgreSQL
# (Use your local Postgres or:)
docker run -d -p 5432:5432 -e POSTGRES_USER=herbst \
  -e POSTGRES_PASSWORD=herbst_password \
  -e POSTGRES_DB=herbst_mud \
  postgres:15

# Terminal 2: Start REST API
cd server
cp .env.example .env
go mod download
go run .

# Terminal 3: Start SSH Server
cd herbst
cp .env.example .env
go mod download
go run .

# Terminal 4: Start Admin Panel (optional)
cd admin
cp .env.example .env
npm install
npm run dev
```

---

## Database Architecture

### Connection Strategy

The application supports two database connection methods:

#### Option 1: DATABASE_URL (Recommended for Neon)

```go
// server/main.go and herbst/main.go
func getDBConfig() string {
    // Neon DB and managed Postgres set DATABASE_URL
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        return dbURL
    }
    // ... fallback to individual variables
}
```

**Priority Order:**
```
1. DATABASE_URL env var
2. DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
3. Development defaults (localhost only)
```

#### Option 2: Individual Variables

Used for local development or when `DATABASE_URL` is not available:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=herbst
DB_PASSWORD=herbst_password
DB_NAME=herbst_mud
DB_SSL_MODE=require  # or "disable" for local dev
```

### SSL Mode Detection

The code automatically detects SSL requirements:

```go
sslMode := os.Getenv("DB_SSL_MODE")
if sslMode == "" {
    if isDev {  // No env vars set
        sslMode = "disable"     // Local dev
    } else {
        sslMode = "require"       // Production (Neon compatible)
    }
}
```

**Behavior:**

| Environment | DB_SSL_MODE Default | Notes |
|-------------|---------------------|-------|
| Production | `require` | Mandatory for Neon DB |
| Development | `disable` | Local Postgres without SSL |

### Database Schema (Ent ORM)

```go
// server/db/schema/
├── user.go          # User accounts
├── character.go     # Player characters
├── room.go          # Game rooms
├── item.go          # Equipment/items
├── skill.go         # Learned skills
├── talent.go        # Available talents
└── quest.go         # Quest system
```

**Key Relationships:**
- User has many Characters
- Character belongs to Room
- Character has many Items
- Character has many Skills
- Character has many Talents

---

## Security Configuration

### CORS (Cross-Origin Resource Sharing)

CORS origins are configurable via environment variable:

```go
// server/main.go
allowedOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")

router.Use(func(c *gin.Context) {
    origin := c.Request.Header.Get("Origin")
    allowed := false
    for _, o := range strings.Split(allowedOrigins, ",") {
        if strings.TrimSpace(o) == origin || origin == "" {
            allowed = true
            break
        }
    }
    // ... set headers
})
```

**Configuration:**

```bash
# Development
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# Production
CORS_ORIGINS=https://your-domain.com,https://admin.your-domain.com
```

**Allowed Methods:**
- GET, POST, PUT, PATCH, DELETE, OPTIONS

**Allowed Headers:**
- Content-Type, Authorization

### Rate Limiting

Protection against DoS and brute force attacks:

```go
// server/main.go
rate := getEnv("RATE_LIMIT", "100")      // requests
window := getEnv("RATE_WINDOW", "60")    // seconds

limiterStore := memory.NewStore()
limiterRate := limiter.Rate{
    Period: time.Duration(windowInt) * time.Second,
    Limit:  int64(rateInt),
}
rateLimiter := limiter.New(limiterStore, limiterRate)

router.Use(func(c *gin.Context) {
    context, err := rateLimiter.Get(c.Request.Context(), c.ClientIP())
    if context.Reached {
        c.AbortWithStatus(http.StatusTooManyRequests)
        return
    }
    c.Next()
})
```

**Default Limits:**
- 100 requests per 60 seconds per IP
- Configurable via environment variables

### JWT Authentication

```go
// server/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        // Remove "Bearer " prefix
        // Parse and validate token
        // Set user context
    }
}
```

**Token Generation:**

```go
// server/routes/user_routes.go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "exp":     time.Now().Add(time.Hour * 24).Unix(),
})
tokenString, _ := token.SignedString(jwtSecret)
```

**Environment Variable:**

```bash
# Generate secret
JWT_SECRET=$(openssl rand -base64 32)

# Minimum length: 32 characters
```

---

## API Development

### Route Registration

```go
// server/main.go
// Public routes
routes.RegisterRoomRoutes(router, client)
routes.RegisterUserRoutes(router, client)
routes.RegisterCharacterRoutes(router, client)
routes.RegisterEquipmentRoutes(router, client)

// Protected routes
protected := router.Group("/api")
protected.Use(middleware.AuthMiddleware())
{
    // Add protected endpoints here
}
```

### Creating a New Route

```go
// server/routes/my_routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "herbst-server/db"
)

func RegisterMyRoutes(r *gin.Engine, client *db.Client) {
    r.GET("/my-resource", func(c *gin.Context) {
        // Handler logic
    })
    
    r.POST("/my-resource", func(c *gin.Context) {
        // Handler logic
    })
}
```

Don't forget to register in `main.go`:

```go
routes.RegisterMyRoutes(router, client)
```

### OpenAPI Specification

The API serves its own OpenAPI spec at `/openapi.json`:

```go
router.GET("/openapi.json", func(c *gin.Context) {
    c.JSON(http.StatusOK, getOpenAPISpec())
})
```

**Generate TypeScript Client:**

```bash
cd admin
npx @hey-api/openapi-ts \
  -i http://localhost:8080/openapi.json \
  -o src/client
```

---

## SSH Server Development

### Bubble Tea Model

The SSH server uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework:

```go
// herbst/model.go
type model struct {
    session      ssh.Session
    client       *db.Client
    character    *ent.Character
    currentRoom  int
    screen       ScreenState
    textInput    textinput.Model
    messages     []Message
    // ... more fields
}
```

### Command Structure

Commands are implemented as methods on the model:

```go
// herbst/cmd_*.go
func (m *model) handleCommand(cmd string) {
    parts := strings.Fields(cmd)
    if len(parts) == 0 {
        return
    }
    
    switch parts[0] {
    case "look", "l":
        m.handleLook()
    case "move", "go", "n", "s", "e", "w":
        m.handleMovement(parts[0])
    // ... more commands
    }
}
```

### Adding a New Command

1. Create handler in appropriate `cmd_*.go` file
2. Register in main command handler
3. Add help text

Example:

```go
// herbst/cmd_custom.go
func (m *model) handleCustomCommand(args []string) {
    if len(args) < 2 {
        m.AppendMessage("Usage: custom <arg>", "error")
        return
    }
    // Command logic
    m.AppendMessage("Success!", "success")
}

// In main handler:
case "custom":
    m.handleCustomCommand(parts)
```

---

## Admin Panel Development

### Technology Stack

- **Framework:** React 18
- **Router:** TanStack Router
- **Build:** Vite
- **Styling:** Tailwind CSS
- **API Client:** Generated from OpenAPI spec
- **HTTP Client:** @tanstack/react-query

### Project Structure

```
admin/src/
├── routes/              # TanStack file-based routing
│   ├── index.tsx       # Home/Dashboard
│   ├── rooms/          # Room management
│   ├── characters/     # Character management
│   └── users/          # User management
├── client/             # Generated API client
│   ├── services.gen.ts
│   └── types.gen.ts
└── components/         # Reusable components
```

### Generating API Client

After changing API routes:

```bash
cd admin
npx @hey-api/openapi-ts \
  -i http://localhost:8080/openapi.json \
  -o src/client
```

### API Base URL

Set in `.env`:

```bash
VITE_API_BASE_URL=http://localhost:8080
# Production:
VITE_API_BASE_URL=https://api.your-domain.com
```

---

## Testing

### Running Tests

```bash
# Server tests
cd server
go test ./... -v

# SSH client tests
cd herbst
go test ./... -v

# Specific test
go test -run TestAuthentication -v
```

### Test Structure

```go
// server/*_test.go
func TestUserCreation(t *testing.T) {
    // Setup
    client := setupTestDB(t)
    defer client.Close()
    
    // Test
    user, err := client.User.Create().
        SetEmail("test@example.com").
        SetPasswordHash("hashed").
        Save(context.Background())
    
    // Assert
    if err != nil {
        t.Fatalf("failed creating user: %v", err)
    }
    if user.Email != "test@example.com" {
        t.Errorf("unexpected email: got %s, want %s", 
            user.Email, "test@example.com")
    }
}
```

### BDD Testing with Gherkin

```gherkin
# features/player_commands.feature
Feature: Player Commands
  Scenario: Player looks around
    Given the player is in room 1
    When the player types "look"
    Then the player sees the room description
```

---

## Deployment Pipeline

### Local Development Flow

```
1. Make code changes
2. Run tests: go test ./...
3. Verify locally: docker-compose up
4. Commit with 🔵 badge
5. Push to GitHub
6. CI/CD runs tests
7. Digital Ocean auto-deploys (if configured)
```

### Production Deployment

```
1. Update .do/app.yaml (if needed)
2. Commit changes
3. Git push to main
4. Digital Ocean App Platform:
   - Detects push
   - Builds Docker images
   - Runs health checks
   - Deploys new version
5. SSH Server:
   - Manual pull on Droplet
   - Or set up CI/CD
```

### Environment Promotion

| Environment | Database | CORS | Rate Limit |
|-------------|----------|------|------------|
| Local | Local/Docker | * | High |
| Staging | Neon (test) | Staging domain | Low |
| Production | Neon (prod) | Production domain | Standard |

---

## Common Tasks

### Database Migration

Ent automatically migrates on startup:

```go
// Runs automatically in main()
if err := client.Schema.Create(context.Background()); err != nil {
    log.Fatalf("failed creating schema: %v", err)
}
```

### Adding a New Entity

1. Create schema file in `server/db/schema/`
2. Run `go generate ./...`
3. Ent generates client code
4. Use in routes

### Debugging

```bash
# Server logs
make logs-server

# SSH client logs
make logs-herbst

# All services
make logs

# Database queries
docker-compose exec postgres psql -U herbst -d herbst_mud
```

---

## Resources

- [Ent Documentation](https://entgo.io/)
- [Gin Framework](https://gin-gonic.com/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [TanStack Router](https://tanstack.com/router)
- [CODE_ARCHITECTURE.md](../CODE_ARCHITECTURE.md) - Detailed architecture

---

🔵 Document version: 2026-04-04
