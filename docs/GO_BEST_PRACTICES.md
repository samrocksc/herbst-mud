# Go Best Practices Review - Herbst MUD

> Comprehensive analysis for Extensibility, Modularity, and Efficiency
**Review Date:** 2026-04-05
**Scope:** `server/` and `herbst/` packages

---

## Executive Summary

Our codebase shows **strong foundation** in modularity but has areas needing improvement for extensibility and efficiency. Key wins: good file size discipline (mostly under 200 lines), clear separation of concerns between `herbst/` (client) and `server/` (API). Critical need: extract shared types/interfaces and reduce tight coupling.

---

## 1. Extensibility Analysis

### Current State: ⚠️ NEUTRAL

**Strengths:**
- Route handlers are grouped by domain (`character_routes.go`, `room_routes.go`)
- Bubbletea framework provides good Update/View/Model separation
- Database uses ent ORM which is extensible

**Issues Identified:**

#### 1.1 Missing Interface Abstractions
**Problem:** Routes directly depend on concrete types, making testing and mocking difficult.

```go
// CURRENT: Tight coupling
router.POST("/characters/:id/damage", func(c *gin.Context) {
    client.Character.UpdateOneID(id)... // Hard dependency
})
```

**Recommendation:** Define service interfaces:

```go
// RECOMMENDED: Define in server/services/character_service.go
type CharacterService interface {
    GetCharacter(ctx context.Context, id int) (*Character, error)
    ApplyDamage(ctx context.Context, id int, damage int) error
    Heal(ctx context.Context, id int, amount int) error
}

type characterService struct {
    client *db.Client
}

// Routes depend on interface, not implementation
func RegisterCharacterRoutes(router *gin.Engine, svc CharacterService) {
    // Easy to mock for testing
}
```

#### 1.2 Command Pattern Not Fully Exploited
**Problem:** Commands in `herbst/commands.go` use string matching, limiting extensibility.

```go
// CURRENT: String-based routing
func (m *model) handleCommand(input string) {
    switch command {
    case "look", "l":
        m.handleLookCommand(args)
    // Adding new commands requires editing this switch
    }
}
```

**Recommendation:** Implement command registry pattern:

```go
// RECOMMENDED: herbst/commands/registry.go
type CommandHandler func(m *model, args []string) tea.Cmd

type CommandRegistry struct {
    handlers map[string]CommandHandler
    aliases  map[string]string
}

func (r *CommandRegistry) Register(name string, handler CommandHandler, aliases ...string) {
    r.handlers[name] = handler
    for _, alias := range aliases {
        r.aliases[alias] = name
    }
}

func (r *CommandRegistry) Execute(m *model, command string, args []string) tea.Cmd {
    if canonical, ok := r.aliases[command]; ok {
        command = canonical
    }
    if handler, ok := r.handlers[command]; ok {
        return handler(m, args)
    }
    return m.errorMessage("Unknown command")
}
```

#### 1.3 Skill System Hard Coded
**Problem:** NPC skills and classless skills are hardcoded in multiple files.

**Recommendation:** Centralize in YAML/JSON config with runtime loading:

```go
// server/skills/registry.go
type SkillRegistry struct {
    skills map[string]SkillDefinition
    loader SkillLoader // Interface for different sources
}

type SkillLoader interface {
    Load(ctx context.Context) ([]SkillDefinition, error)
}
```

---

## 2. Modularity Analysis

### Current State: ✅ GOOD

**Strengths:**
- File sizes mostly under 200 lines (excellent!)
- Clear package separation: `routes/`, `db/`, `dbinit/`
- Command handlers split by domain (`cmd_look.go`, `cmd_movement.go`)

**Refinement Opportunities:**

#### 2.1 Shared Types Duplication
**Problem:** Client and server may duplicate type definitions.

**Recommendation:** Create shared types package:

```
/shared/
  types/
    character.go    # Character data structures
    room.go         # Room structures  
    combat.go       # Combat messages
    api.go          # API request/response types
```

Both `herbst/` and `server/` import from shared package. Use build tags if needed:

```go
//go:build !client
// Server-only fields
```

#### 2.2 Model God Struct
**Problem:** `model` struct in `herbst/model.go` has many responsibilities:
- UI State management
- Game state
- Network client
- Combat state

**Current count:** ~30 fields (borderline - acceptable for bubbletea)

**Recommendation:** Group related state into sub-structs:

```go
// RECOMMENDED

type model struct {
    // Core
    session     ssh.Session
    client      *db.Client
    connectedAt time.Time
    
    // UI State (managed by bubbletea)
    ui UIModel
    
    // Game State
    game GameModel
    
    // Combat State
    combat CombatModel
    
    // Screen Management
    screen      ScreenState
    width       int
    height      int
}

type UIModel struct {
    textInput    textinput.Model
    spinner      spinner.Model
    messages     []Message
    isScrolling  bool
    viewport     viewport.Model // for scrollback
}

type GameModel struct {
    currentRoom   int
    roomName      string
    roomDesc      string
    exits         []string
    roomItems     []RoomItem
    roomCharacters []RoomCharacter
}

type CombatModel struct {
    inCombat           bool
    combatTarget     *RoomCharacter
    combatManager    *combat.CombatManager
    combatID         int
    combatLog        []string
    combatQueuedAction string
    combatSkills     *CombatSkillState
}
```

#### 2.3 Route Handler Consolidation
**Problem:** `character_routes.go` at 1976 lines violates 100-line guideline.

**Recommendation:** Split by sub-domain:

```
routes/
  character/
    routes.go           # Registration only (100 lines)
    handlers.go         # HTTP handlers
    combat.go           # Combat endpoints
    inventory.go        # Item/equipment endpoints
    skills.go           # Skill/talent endpoints
    stats.go            # Character stats
```

---

## 3. Efficiency Analysis

### Current State: ⚠️ NEUTRAL

**Strengths:**
- Uses ent ORM (efficient query building)
- Context properly passed through
- No obvious memory leaks

**Improvements Needed:**

#### 3.1 Database Connection Pooling
**Problem:** Not explicit about connection pool settings.

**Recommendation:** Configure in db.Open():

```go
// server/db/client.go
func Open(driver, dsn string) (*Client, error) {
    client, err := ent.Open(driver, dsn, 
        ent.Driver(drv), // custom driver with pooling
    )
    
    // Configure pool (PostgreSQL)
    sqlDB, _ := client.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    return &Client{client}, nil
}
```

#### 3.2 API Response Caching
**Problem:** Room data fetched repeatedly in same session.

**Recommendation:** Add caching layer for read-heavy data:

```go
// server/cache/room_cache.go
type RoomCache struct {
    data  map[int]*CachedRoom
    mu    sync.RWMutex
    ttl   time.Duration
}

type CachedRoom struct {
    room      *ent.Room
    timestamp time.Time
}

func (c *RoomCache) Get(id int) (*ent.Room, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if cached, ok := c.data[id]; ok {
        if time.Since(cached.timestamp) < c.ttl {
            return cached.room, true
        }
    }
    return nil, false
}
```

#### 3.3 Client-Side HTTP Connection Reuse
**Problem:** SSH client may create new connections per request.

**Recommendation:** Use shared http.Client:

```go
// herbst/api/client.go
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// Use httpClient instead of http.Get/Post
```

#### 3.4 JSON Pooling (Advanced)
**Problem:** Frequent JSON marshaling/unmarshaling creates GC pressure.

**Recommendation:** For high-traffic endpoints, use sync.Pool:

```go
var jsonBufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func encodeJSON(w io.Writer, v interface{}) error {
    buf := jsonBufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer jsonBufferPool.Put(buf)
    
    encoder := json.NewEncoder(buf)
    if err := encoder.Encode(v); err != nil {
        return err
    }
    
    _, err := w.Write(buf.Bytes())
    return err
}
```

---

## 4. Specific Recommendations by File

### High Priority

| File | Lines | Issue | Action |
|------|-------|-------|--------|
| `character_routes.go` | 1976 | Too long | Split into routes/character/ package |
| `commands.go` | 170 | Switch dispatch | Implement command registry |
| `model.go` | 210 | Can be improved | Group into sub-structs |

### Medium Priority

| File | Lines | Issue | Action |
|------|-------|-------|--------|
| `game_combat.go` | 678 | Combat logic | Extract combat engine to pkg/ |
| `game_model.go` | 535 | Update/View | Separate pure UI from game logic |
| `equipment_routes.go` | 553 | Route length | Split into handlers/ |

### Low Priority

| File | Lines | Issue | Action |
|------|-------|-------|--------|
| `cmd_regen.go` | 226 | Logic mixed | Extract regen service |
| `npc_skills.go` | ~200 | Hardcoded | Config-driven skills |

---

## 5. Recommended Package Structure

### Current vs Recommended

```
CURRENT:
herbst/
  main.go
  model.go
  commands.go
  cmd_*.go (various)
  game_*.go
  npc_skills.go
  classless_skills.go
  combat_talents.go
  ui_*.go

RECOMMENDED:
herbst/
  main.go
  cmd/
    registry.go
    handlers/
      look.go
      movement.go
      combat.go
      inventory.go
  game/
    state.go          # Game state only
    combat/
      engine.go       # Combat logic
      npc_ai.go
      skills.go
    room/
      loader.go
      renderer.go
  ui/
    model.go          # Bubbletea model
    components/
      input.go
      messages.go
      combat.go
  api/
    client.go         # HTTP client
  shared/
    types/            # Shared with server
```

---

## 6. Implementation Priority Matrix

### Phase 1: Quick Wins (1-2 days)
- [ ] Create `herbst/cmd/registry.go` for command pattern
- [ ] Move combat logic to `herbst/game/combat/engine.go`
- [ ] Add connection pooling to db.Open()

### Phase 2: Structural (1 week)
- [ ] Split `character_routes.go` by domain
- [ ] Group model fields into sub-structs
- [ ] Create `shared/` package for types

### Phase 3: Optimization (2 weeks)
- [ ] Implement service interfaces
- [ ] Add caching layer
- [ ] Config-driven skills system

---

## 7. Testing Strategy

With interfaces, testing becomes trivial:

```go
// Mock implementation
type MockCharacterService struct {
    characters map[int]*Character
}

func (m *MockCharacterService) GetCharacter(ctx context.Context, id int) (*Character, error) {
    if char, ok := m.characters[id]; ok {
        return char, nil
    }
    return nil, ErrNotFound
}

// Test routes with mock
func TestCharacterDamage(t *testing.T) {
    mock := &MockCharacterService{
        characters: {1: {ID: 1, HP: 100}},
    }
    router := setupTestRouter(mock)
    
    // Test without real DB
}
```

---

## Conclusion

Our codebase has **solid foundations** with good file size discipline and domain separation. The main improvements needed are:

1. **Interfaces** for service layer (extensibility)
2. **Shared types** package (modularity)
3. **Command registry** pattern (extensibility)
4. **Connection pooling** (efficiency)

These changes will transform a good codebase into a **maintainable, testable, scalable** system.

---

*Cowabunga, dudes! Clean code = happy developers! 🐢🟣*
