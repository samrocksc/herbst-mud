# Herbst MUD Refactoring Guide

## Overview

This guide outlines the plan to refactor the Herbst MUD codebase to be more modular, reduce dependencies, and follow the "files should not exceed 100 lines" rule from `docs/CODE.md`.

## Architecture Summary

```
┌─────────────────────────────────────────────────────────────────┐
│                     Herbst MUD                                   │
├──────────────────┬───────────────────┬──────────────────────────┤
│   herbst/        │    server/        │       admin/            │
│   SSH Client     │    REST API       │    Web Admin Panel      │
│   Port: 4444     │    Port: 8080     │       Port: 3000        │
│   Go/BubbleTea   │    Go/Gin         │   TypeScript/React      │
└──────────────────┴───────────────────┴──────────────────────────┘
                        PostgreSQL Database (ent ORM)
```

---

## Phase 1: Shared Database Package

### Problem
Both `herbst/db/` and `server/db/` contain nearly identical ent-generated code. This causes:
- Code duplication (~70k+ lines duplicated)
- Maintenance burden (changes must be made twice)
- Potential schema drift

### Solution
Create a shared `db/` package at the root level.

### Steps

#### 1.1 Create shared db package
```
/db/
  /schema/
    character.go
    room.go
    equipment.go
    user.go
    skill.go
    talent.go
    npctemplate.go
    availabletalent.go
    characterskill.go
    charactertalent.go
  /ent/          (generated code)
  client.go
```

#### 1.2 Update imports
- `herbst/main.go`: Change `import "herbst/db"` → `import "github.com/samrocksc/herbst-mud/db"`
- `server/main.go`: Same change

#### 1.3 Generate ent code once
```bash
cd db && go generate ./...
```

#### 1.4 Remove duplicated directories
- Delete `herbst/db/`
- Delete `server/db/`

### Expected Outcome
- ~70k lines of duplicated code removed
- Single source of truth for database schema
- Easier schema migrations

---

## Phase 2: Server Routes Refactoring

### Problem
`server/routes/character_routes.go` is ~1200 lines - violates the 100-line rule.

### Current Structure
```
server/routes/
  character_routes.go   (~1200 lines) ❌
  equipment_routes.go   (~520 lines)  ❌
  room_routes.go        (~170 lines)  ⚠️
  user_routes.go        (~170 lines)  ⚠️
```

### Target Structure
```
server/routes/
  /character/
    auth.go           (~80 lines) - authentication endpoints
    crud.go           (~80 lines) - character CRUD
    stats.go          (~80 lines) - stats/attributes
    skills.go         (~80 lines) - skill management
    talents.go        (~80 lines) - talent management
    npcs.go           (~80 lines) - NPC listing/creation
    class.go          (~60 lines) - class endpoints
    race.go           (~60 lines) - race endpoints
  /equipment/
    crud.go           (~80 lines) - equipment CRUD
    room.go           (~80 lines) - room equipment endpoints
    reveal.go         (~80 lines) - reveal hidden items
    examine.go        (~80 lines) - examine endpoint
  /room/
    crud.go           (~80 lines) - room CRUD
    characters.go     (~60 lines) - room characters
    look.go           (~60 lines) - room look endpoint
  /user/
    auth.go           (~80 lines) - user authentication
    crud.go           (~80 lines) - user CRUD
```

### Refactoring Pattern

**Before (character_routes.go):**
```go
// 1200 lines of mixed concerns
router.POST("/characters", func(c *gin.Context) { ... })
router.GET("/characters/:id", func(c *gin.Context) { ... })
router.POST("/characters/authenticate", func(c *gin.Context) { ... })
// ... many more endpoints
```

**After (character/crud.go):**
```go
package character

import (
    "github.com/gin-gonic/gin"
    "github.com/samrocksc/herbst-mud/db"
)

// RegisterCRUDRoutes registers character CRUD endpoints
func RegisterCRUDRoutes(router *gin.RouterGroup, client *db.Client) {
    router.POST("/characters", createCharacter(client))
    router.GET("/characters", listCharacters(client))
    router.GET("/characters/:id", getCharacter(client))
    router.PUT("/characters/:id", updateCharacter(client))
    router.DELETE("/characters/:id", deleteCharacter(client))
}

func createCharacter(client *db.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ~30 lines
    }
}

// ... rest split into other files
```

### Migration Steps

1. Create subdirectories under `server/routes/`
2. Move each endpoint group to its own file
3. Export `Register*Routes` functions
4. Update `server/main.go` to call each registration function

---

## Phase 3: Server Main.go Cleanup

### Problem
`server/main.go` contains:
- Server initialization (~50 lines)
- OpenAPI spec inline (~400 lines)
- Route registration scattered throughout

### Target Structure
```
server/
  main.go            (~50 lines) - just initialization
  openapi.go         (~400 lines) - OpenAPI spec
  routes.go          (~30 lines) - route registration
```

### Before
```go
func main() {
    // ... setup
    router.GET("/openapi.json", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "openapi": "3.0.0",
            // ... 400 lines of spec
        })
    })
    // ... routes scattered
}
```

### After
```
// server/main.go
func main() {
    client := db.Init()
    defer client.Close()

    router := gin.Default()
    router.Use(cors.Default())

    routes.RegisterAll(router, client)

    log.Println("Server starting on :8080")
    log.Fatal(router.Run(":8080"))
}

// server/routes.go
func RegisterAll(router *gin.Engine, client *db.Client) {
    router.GET("/healthz", handlers.Healthz)
    router.GET("/openapi.json", handlers.OpenAPI)

    v1 := router.Group("/api/v1")
    character.RegisterRoutes(v1, client)
    equipment.RegisterRoutes(v1, client)
    room.RegisterRoutes(v1, client)
    user.RegisterRoutes(v1, client)
}
```

---

## Phase 4: Herbst Command Refactoring

### Problem
Several herbst command files exceed 100 lines.

### Current Structure
```
herbst/
  cmd_look.go        (~143 lines) ❌
  game_model.go      (~240 lines) ❌
  commands.go        (~100 lines) ✅
```

### Target Structure
```
herbst/
  /commands/
    look.go          (~60 lines) - look command
    examine.go       (~60 lines) - examine command
    search.go        (~60 lines) - search command
    movement.go      (~60 lines) - n/s/e/w/u/d commands
    inventory.go     (~60 lines) - inventory management
    combat.go        (~60 lines) - attack/flee commands
    social.go        (~60 lines) - say/whisper/emote
  /game/
    model.go         (~80 lines) - core model struct
    update.go        (~80 lines) - BubbleTea Update function
    view.go          (~80 lines) - BubbleTea View function
    room.go          (~80 lines) - room rendering
    message.go       (~60 lines) - message formatting
  /ui/
    screens.go       (~80 lines) - welcome/login/register screens
    styles.go        (~60 lines) - lipgloss styles
```

### Pattern

**Before (cmd_look.go):**
```go
// 143 lines mixing multiple commands
func (m *model) handleLookCommand(cmd string) { ... }
func (m *model) handleExamineCommand(cmd string) { ... }
func (m *model) handleSearchCommand(cmd string) { ... }
```

**After (commands/look.go):**
```go
package commands

import "herbst/game"

func HandleLook(cmd string, g *game.Model) string {
    // ~50 lines
}
```

---

## Phase 5: Admin Panel Refactoring

### Problem
Admin route files exceed 100 lines and mix concerns.

### Current Structure
```
admin/src/routes/
  map.tsx            (~500 lines) ❌
  npcs.tsx           (~400 lines) ❌
  items.tsx          (~560 lines) ❌
  dashboard.tsx      (~120 lines) ⚠️
```

### Target Structure
```
admin/src/
  /components/
    Modal.tsx         (✅ already exists)
    /map/
      MapCanvas.tsx   (~80 lines)
      RoomNode.tsx    (~60 lines)
      RoomPanel.tsx   (~80 lines)
      ZLevelNav.tsx   (~60 lines)
    /npc/
      NPCList.tsx    (~60 lines)
      NPCForm.tsx    (~80 lines)
    /items/
      ItemList.tsx   (~60 lines)
      ItemForm.tsx    (~80 lines)
  /hooks/
    useRooms.ts      (~60 lines) - room data fetching
    useNPCs.ts       (~60 lines) - NPC data fetching
    useItems.ts      (~60 lines) - item data fetching
  /routes/
    map.tsx          (~80 lines) - just composition
    npcs.tsx         (~80 lines)
    items.tsx        (~80 lines)
    dashboard.tsx    (~80 lines)
```

### Pattern

**Before (routes/map.tsx):**
```tsx
// 500 lines mixing data fetching, rendering, and state
function MapBuilder() {
  const [rooms, setRooms] = useState([])
  // ... 400 lines of inline logic
  return ( /* 100 lines of JSX */ )
}
```

**After (routes/map.tsx):**
```tsx
// ~80 lines - composition only
import { MapCanvas } from '../components/map/MapCanvas'
import { RoomPanel } from '../components/map/RoomPanel'
import { ZLevelNav } from '../components/map/ZLevelNav'
import { useRooms } from '../hooks/useRooms'

export default function MapBuilder() {
  const { rooms, selectedRoom, selectRoom } = useRooms()
  const [zLevel, setZLevel] = useState(0)

  return (
    <div className="flex h-screen">
      <ZLevelNav level={zLevel} onChange={setZLevel} />
      <MapCanvas rooms={rooms} zLevel={zLevel} onSelect={selectRoom} />
      <RoomPanel room={selectedRoom} />
    </div>
  )
}
```

---

## Phase 6: API Client Abstraction

### Problem
HTTP calls in herbst are scattered throughout command files with hardcoded URLs.

### Current Pattern
```go
// Scattered in multiple files
resp, err := http.Get("http://localhost:8080/rooms/" + roomId)
resp, err := http.Post("http://localhost:8080/characters/auth", ...)
```

### Solution: Centralized API Client

```
herbst/
  /api/
    client.go        (~80 lines) - HTTP client with base URL
    rooms.go         (~60 lines) - room endpoints
    characters.go    (~80 lines) - character endpoints
    equipment.go     (~60 lines) - equipment endpoints
    auth.go          (~60 lines) - authentication
```

### Pattern

**herbst/api/client.go:**
```go
package api

import (
    "net/http"
    "time"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
}
```

**herbst/api/rooms.go:**
```go
package api

func (c *Client) GetRoom(id int) (*Room, error) {
    resp, err := c.httpClient.Get(c.baseURL + "/rooms/" + strconv.Itoa(id))
    // ... ~30 lines
}
```

---

## Phase 7: Combat System Integration

### Problem
`herbst/combat/` package is isolated but not fully integrated with the main game loop.

### Current State
```
herbst/combat/
  manager.go         (~100 lines) - CombatManager struct
  tick.go            (~60 lines)  - Tick-based timing
  config.go          (~40 lines)  - Constants
```

### Integration Needed
- Connect combat manager to game model
- Add combat state to model
- Process combat ticks in Update()

### Target
```
herbst/
  /combat/
    manager.go       (~80 lines)
    tick.go          (~60 lines)
    types.go         (~40 lines) - combat-specific types
  game_model.go      - Add CombatManager field
```

---

## Phase 8: DB Init Consolidation

### Problem
Both `herbst/dbinit/` and `server/dbinit/` exist with similar seeding logic.

### Solution
```
/dbinit/
  main.go            (~60 lines) - entry point
  admin.go           (~40 lines) - admin user
  rooms.go           (~80 lines) - initial rooms
  npcs.go            (~60 lines) - NPC templates
  items.go           (~80 lines) - starting items
  skills.go          (~60 lines) - skill definitions
  talents.go         (~60 lines) - talent definitions
```

### Usage
```go
// server/main.go
import "github.com/samrocksc/herbst-mud/dbinit"

func main() {
    client := db.Init()
    dbinit.SeedAll(client) // Seeds everything
}

// herbst/main.go - same import
```

---

## Implementation Order

1. **Phase 1: Shared DB Package** (Highest Impact)
   - Removes ~70k lines of duplication
   - Single source of truth

2. **Phase 2-3: Server Routes Refactoring**
   - Improves maintainability
   - Follows 100-line rule

3. **Phase 6: API Client Abstraction**
   - Reduces scattered HTTP calls
   - Makes testing easier

4. **Phase 4-5: Command/UI Refactoring**
   - Smaller, testable files
   - Better separation of concerns

5. **Phase 7: Combat Integration**
   - Complete feature

6. **Phase 8: DB Init Consolidation**
   - Single seeding process

---

## Testing Strategy

After each phase:

1. **Unit Tests**: Run `go test ./...` in affected packages
2. **Integration Tests**: Verify API endpoints work
3. **E2E Tests**: Play through SSH client
4. **Admin Panel**: `npm run build && npm run dev`

---

## Rollback Plan

Each phase should be:
1. Done in a separate branch
2. Tested thoroughly
3. Merged only after verification

If issues arise:
- Revert the specific phase commit
- Fix issues in isolation
- Retry merge