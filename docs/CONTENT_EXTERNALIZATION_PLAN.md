# Donatello's 8-Week Content Externalization Project

> Transform hardcoded game content into data-driven configuration for true multi-MUD engine capability

**Project Owner:** Donatello (🟣)
**Start Date:** 2026-04-07
**Goal:** Make Herbst MUD a true multi-MUD engine by externalizing all content to JSON/YAML configs

---

## Executive Summary

**Current State:** Skills, NPCs, items, and talents are hardcoded in Go files
**Target State:** All content lives in `content/` directory as JSON/YAML, loadable at runtime

**Benefits:**
- New MUDs can be created by editing files, not code
- Content creators don't need Go knowledge
- Faster iteration on game design
- Database remains user state only (clean separation)

---

## Week 1: Content Architecture & Foundation

### Goals
- Design content schema for all game entities
- Create content loading framework
- Implement YAML/JSON validation

### Deliverables

#### 1.1 Content Directory Structure
```
content/
├── _schema/              # JSON Schema for validation
│   ├── skill.schema.json
│   ├── npc.schema.json
│   ├── item.schema.json
│   └── room.schema.json
├── _shared/              # Content shared across MUDs
│   └── classless/        # Base classless skills
├── default/              # Default Herbst MUD content
│   ├── skills/
│   │   ├── combat/
│   │   ├── magic/
│   │   └── classless/
│   ├── npcs/
│   │   ├── templates/    # NPC stat templates
│   │   └── spawn/        # Spawn locations/rules
│   ├── items/
│   │   ├── weapons/
│   │   ├── armor/
│   │   └── consumables/
│   ├── rooms/
│   │   ├── areas/        # Area definitions
│   │   └── connections/  # Exit mappings
│   └── quests/
└── custom/               # Custom MUDs (future)
    └── example-mud/
```

#### 1.2 Content Loader Framework
**File:** `server/content/loader.go`

```go
package content

type Manager struct {
    basePath string
    skills   map[string]*SkillDef
    npcs     map[string]*NPCTemplate
    items    map[string]*ItemDef
    rooms    map[string]*RoomTemplate
}

func Load(basePath string) (*Manager, error)
func (m *Manager) Reload() error
func (m *Manager) Validate() []ValidationError
```

#### 1.3 Schema Definitions (YAML-first)
**Example:** `content/_schema/skill.schema.yaml`

```yaml
$schema: "http://json-schema.org/draft-07/schema#"
title: Skill Definition
type: object
required:
  - id
  - name
  - type
properties:
  id:
    type: string
    pattern: "^[a-z_]+$"
  name:
    type: string
    maxLength: 50
  type:
    enum: [combat, passive, active, triggered]
  description:
    type: string
    maxLength: 500
  effects:
    type: array
    items:
      type: object
      properties:
        type:
          enum: [damage, heal, buff, debuff, dot, hot]
        value:
          type: number
        duration:
          type: integer
          description: "Duration in ticks"
  requirements:
    type: object
    properties:
      level:
        type: integer
        minimum: 1
      class:
        type: string
      skill_prereq:
        type: string
  cooldown:
    type: integer
    minimum: 0
  mana_cost:
    type: integer
    minimum: 0
  stamina_cost:
    type: integer
    minimum: 0
```

### Week 1 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Design content schema for all entities | 4 | Architect |
| Implement content loader framework | 6 | Donatello |
| Create JSON Schema validators | 4 | Donatello |
| Write documentation | 2 | Architect |
| **Total** | **16 hours** | |

---

## Week 2: Skill System Externalization

### Goals
- Move hardcoded skills to YAML/JSON
- Implement skill registry with runtime loading
- Update combat system to use external skills

### Deliverables

#### 2.1 Skill Content Files
**File:** `content/default/skills/classless/concentrate.yaml`

```yaml
id: concentrate
name: "Concentrate"
type: combat
description: "Focus your mind, increasing critical hit chance"
effects:
  - type: buff
    stat: crit_chance
    value: 5
    duration: 3
cooldown: 3
stamina_cost: 5
visual:
  icon: "⚡"
  animation: "buff_self"
  sound: "concentrate.wav"
```

**File:** `content/default/skills/classless/haymaker.yaml`

```yaml
id: haymaker
name: "Haymaker"
type: combat
description: "A powerful haymaker punch with bonus damage"
effects:
  - type: damage
    value: 15
    scaling:
      stat: strength
      ratio: 0.5
cooldown: 4
stamina_cost: 8
mana_cost: 0
visual:
  icon: "👊"
  animation: "punch_heavy"
  sound: "punch_hard.wav"
```

#### 2.2 Skill Registry Refactor
**File:** `server/content/skill_registry.go`

```go
package content

import (
    "fmt"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

// SkillRegistry loads and manages skill definitions
type SkillRegistry struct {
    skills map[string]*SkillDef
}

// LoadFromDirectory loads all YAML files in directory recursively
func (r *SkillRegistry) LoadFromDirectory(path string) error {
    return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && (filepath.Ext(p) == ".yaml" || filepath.Ext(p) == ".yml") {
            return r.loadFile(p)
        }
        return nil
    })
}

// Get retrieves a skill by ID (case insensitive)
func (r *SkillRegistry) Get(id string) (*SkillDef, error)

// GetByClass returns all skills for a class
func (r *SkillRegistry) GetByClass(class string) []*SkillDef

// Validate checks all skills for consistency
func (r *SkillRegistry) Validate() []ValidationError
```

#### 2.3 Combat Integration
Update `game_combat.go` to use `content.SkillRegistry` instead of hardcoded skills.

### Week 2 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Create skill YAML files for all classless skills | 2 | Donatello |
| Implement skill registry | 4 | Donatello |
| Refactor combat system to use registry | 6 | Donatello |
| Update client to load skills from server | 4 | Donatello |
| Testing & validation | 2 | Team |
| **Total** | **18 hours** | |

### Migration Plan
```bash
# Old code to remove:
herbst/classless_skills.go    # Hardcoded definitions
herbst/cmd_skills.go          # Hardcoded skill logic

# New code:
server/content/skill_registry.go
content/default/skills/**/*.yaml
```

---

## Week 3: NPC Template System

### Goals
- Externalize NPC stats, AI behavior, skills
- Create spawn location configs
- Implement NPC factory pattern

### Deliverables

#### 3.1 NPC Template Files
**File:** `content/default/npcs/templates/aragorn.yaml`

```yaml
id: aragorn
template_name: "Aragorn, the Ranger"
description: "A weathered ranger with keen eyes"
level: 10
stats:
  strength: 45
  dexterity: 50
  constitution: 42
  intelligence: 35
  wisdom: 48
  charisma: 40
hp:
  base: 100
  con_multiplier: 8
  level_bonus: 10
skills:
  classless:
    slot_1: concentrate
    slot_2: back_off
    slot_3: scream
    slot_4: slap
  special:
    - druid_heal  # References skill ID
ai:
  type: druid
  aggression: passive
  heal_threshold: 0.3
  flee_threshold: 0.1
  speech:
    greeting: "The forest welcomes you, traveler."
    attack: "You shall not pass!"
    defeat: "..."  # Dramatic silence
visual:
  icon: "🧝"
  color: "#228B22"
```

#### 3.2 Spawn Configuration
**File:** `content/default/npcs/spawn/cross_way.yaml`

```yaml
area: cross_way
spawns:
  - npc_id: aragorn
    room_id: 2  # Fountain Plaza
    respawn_time: 300  # 5 minutes
    max_concurrent: 1
    position: center
    
  - npc_id: gimli
    room_id: 2
    respawn_time: 600
    max_concurrent: 1
    
  - npc_id: combat_dummy
    room_id: 5
    respawn_time: 0  # Instant respawn
    flags: [invincible, no_attack]
```

#### 3.3 NPC Factory
**File:** `server/content/npc_factory.go`

```go
package content

// NPCFactory creates Character entities from templates
type NPCFactory struct {
    templates map[string]*NPCTemplate
}

// Spawn creates a new NPC instance from template
func (f *NPCFactory) Spawn(templateID string, roomID int) (*ent.Character, error)

// CheckSpawnConditions checks if spawn should occur
func (f *NPCFactory) CheckSpawnConditions(spawn SpawnConfig) bool
```

### Week 3 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Create NPC template schema | 3 | Architect |
| Design spawn configuration | 2 | Architect |
| Implement NPC factory | 4 | Donatello |
| Migrate existing NPCs to YAML | 3 | Donatello |
| Implement spawn system | 4 | Donatello |
| **Total** | **16 hours** | |

---

## Week 4: Item System Externalization

### Goals
- Move all item definitions to config
- Implement item effects system
- Create equipment slot system

### Deliverables

#### 4.1 Item Content Files
**File:** `content/default/items/weapons/iron_sword.yaml`

```yaml
id: iron_sword
name: "Iron Sword"
description: "A basic iron sword with a well-worn grip"
level: 1
slot: weapon
type: one_handed_slash
damage:
  min: 3
  max: 7
  speed: 2.0  # Attacks per tick
requirements:
  level: 1
  class: [warrior, survivor]
  strength: 10
stats:
  strength: +2
  crit_chance: +1%
effects:
  on_hit:
    - type: bleed
      chance: 0.05
      duration: 3
      damage: 2
visual:
  icon: "⚔️"
  color: "#C0C0C0"
  equip_animation: "draw_sword"
value:
  buy: 50
  sell: 10
  weight: 3.5
rarity: common
durability:
  max: 100
  repairable: true
```

**File:** `content/default/items/consumables/health_potion.yaml`

```yaml
id: health_potion
name: "Health Potion"
description: "Restores health when consumed"
type: consumable
effects:
  - type: heal
    value: 25
    scaling:
      stat: constitution
      ratio: 0.2
uses: 1
level: 1
weight: 0.1
value:
  buy: 25
  sell: 5
stackable: true
max_stack: 10
visual:
  icon: "🍷"
  color: "#FF0000"
  use_animation: "drink_potion"
  effect: "heal_sparkle"
```

#### 4.2 Item Registry
**File:** `server/content/item_registry.go`

```go
package content

type ItemRegistry struct {
    items map[string]*ItemDef
}

// GetBySlot returns all items for a specific equipment slot
func (r *ItemRegistry) GetBySlot(slot SlotType) []*ItemDef

// GetDropTable returns items for a monster/area loot
func (r *ItemRegistry) GetDropTable(tableID string) DropTable

// CanEquip checks if character can equip item
func (r *ItemRegistry) CanEquip(item *ItemDef, char *ent.Character) (bool, string)
```

### Week 4 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Design item schema (weapons, armor, consumables) | 3 | Architect |
| Create item YAML files | 4 | Donatello |
| Implement item registry | 4 | Donatello |
| Update equipment system | 6 | Donatello |
| **Total** | **17 hours** | |

---

## Week 5: Room & Area System

### Goals
- Externalize room definitions
- Create area templates
- Implement room loader with connections

### Deliverables

#### 5.1 Room Content Files
**File:** `content/default/rooms/areas/cross_way.yaml`

```yaml
area_id: cross_way
name: "Cross-Way"
description: "A bustling crossroads where adventurers gather"
rooms:
  - id: cross_way_center
    name: "Cross-Way"
    description: |
      You stand at the center of a crossroads. Ancient cobblestones 
      stretch out in all directions. A weathered signpost points the way.
    flags: [starting_room, safe_zone]
    
  - id: fountain_plaza
    name: "Fountain Plaza"
    description: |
      A grand fountain bubbles in the center of the plaza. 
      The water shimmers with an otherworldly glow.
    exits:
      west: cross_way_center
      east: fountain  # Leads to special room
    
  - id: north_gate
    name: "North Gate"
    description: |
      A sturdy iron gate guards the northern road.
      Beyond, the wilderness stretches into mist.
    exits:
      south: cross_way_center
      north: null  # Leads to area: wilderness/north_road

exits:
  - from: cross_way_center
    to: fountain_plaza
    direction: east
    description: "A cobblestone path leads east toward the fountain"
    
  - from: cross_way_center
    to: north_gate
    direction: north
    description: "The northern road leads to adventure"
```

#### 5.2 Room Loader
**File:** `server/content/room_loader.go`

```go
package content

// RoomLoader creates room entities from YAML
type RoomLoader struct {
    client *ent.Client
}

// LoadArea loads an entire area from YAML
func (l *RoomLoader) LoadArea(path string) error

// ApplyExits creates exit edges between rooms
func (l *RoomLoader) ApplyExits(area AreaConfig) error

// Validate checks for orphaned rooms, missing exits
func (l *RoomLoader) Validate() error
```

### Week 5 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Design room/area schema | 3 | Architect |
| Create room YAML files | 3 | Donatello |
| Implement room loader | 4 | Donatello |
| Add exit validation | 2 | Donatello |
| Migrate existing rooms | 2 | Donatello |
| **Total** | **14 hours** | |

---

## Week 6: Quest System

### Goals
- Design quest framework
- Externalize quest definitions
- Implement quest manager

### Deliverables

#### 6.1 Quest Content Files
**File:** `content/default/quests/prologue.yaml`

```yaml
quest_id: prologue_first_steps
title: "First Steps"
description: "Explore the Cross-Way and meet its inhabitants"
level: 1

giver:
  npc_id: aragorn
  dialog:
    offer: "Welcome to Cross-Way, traveler. Won't you explore a bit?"
    accept: "Excellent! Speak with me again when you've made the rounds."
    complete: "You've met our resident! Good job."

objectives:
  - type: explore
    description: "Visit all rooms in Cross-Way"
    target_rooms:
      - cross_way_center
      - fountain_plaza
      - north_gate
    reward_progress: 25
    
  - type: npc_talk
    description: "Speak with Aragorn"
    target_npc: aragorn
    dialog_trigger: greeting
    
  - type: kill
    description: "Defeat a training dummy"
    target: combat_dummy
    count: 1
    
rewards:
  xp: 100
  gold: 50
  items:
    - id: health_potion
      count: 3
  reputation:
    cross_way: +10
```

### Week 6 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Design quest system architecture | 4 | Architect |
| Implement quest manager | 6 | Donatello |
| Create quest YAML files | 3 | Donatello |
| Quest UI integration | 3 | Donatello |
| **Total** | **16 hours** | |

---

## Week 7: Content Hot-Reload & Admin Tools

### Goals
- Implement runtime content reloading
- Create content validation tools
- Build admin content editor (Web UI)

### Deliverables

#### 7.1 Hot-Reload System
**File:** `server/content/watcher.go`

```go
package content

// Watcher monitors content files for changes
type Watcher struct {
    manager *Manager
    watcher *fsnotify.Watcher
}

// Start begins watching content directory
func (w *Watcher) Start() error

// OnChange callback when file changes
func (w *Watcher) OnChange(path string) error
```

#### 7.2 Admin Content API
**File:** `server/routes/content_routes.go`

```go
// Content Management Endpoints
GET    /admin/content/skills           // List all skills
GET    /admin/content/skills/:id       // Get skill definition
POST   /admin/content/validate         // Validate content
POST   /admin/content/reload           // Hot-reload content
GET    /admin/content/errors           // Get validation errors
```

#### 7.3 Admin Panel Updates
**File:** `admin/src/routes/content.tsx`

- Content browser tree view
- YAML editor with validation
- Preview pane
- "Apply Changes" button with hot-reload

### Week 7 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Implement file watcher | 4 | Donatello |
| Create content admin API | 4 | Donatello |
| Build content UI components | 6 | Donatello |
| Testing hot-reload | 2 | Team |
| **Total** | **16 hours** | |

---

## Week 8: Testing, Optimization & Documentation

### Goals
- Comprehensive testing of content system
- Performance optimization
- Complete documentation

### Deliverables

#### 8.1 Content Testing Framework
**File:** `content/test/validation_test.go`

```go
// All content files validate against schema
func TestAllContentValid(t *testing.T)

// All skill references exist
func TestSkillReferencesValid(t *testing.T)

// No circular room references
func TestNoCircularRoomReferences(t *testing.T)

// NPC skills exist in registry
func TestNPCSkillsValid(t *testing.T)
```

#### 8.2 Performance Benchmarks
- Content load time < 5 seconds for full MUD
- Hot-reload latency < 2 seconds
- Memory usage < 100MB for content cache

#### 8.3 Documentation
- Content Author Guide (for non-programmers)
- Schema Reference
- Migration Guide from Hardcoded
- Troubleshooting

### Week 8 Tasks
| Task | Est Hours | Owner |
|------|-----------|-------|
| Write content validation tests | 4 | Donatello |
| Performance profiling & optimization | 4 | Donatello |
| Content author guide | 4 | Architect |
| Final integration testing | 4 | Team |
| **Total** | **16 hours** | |

---

## Technical Specifications

### Content Loading Priority
1. `_shared/` loaded first (base content)
2. `default/` or active MUD loaded
3. Custom overrides applied

### Validation Pipeline
```
YAML Parse → Schema Validate → Cross-Reference Check → Load Cache
     ↓              ↓                    ↓                  ↓
  Syntax       Types/Required      ID References       In-Memory
  Errors       Constraints          Valid Targets       Registry
```

### Performance Targets
| Metric | Target |
|--------|--------|
| Cold load | < 5 seconds |
| Hot reload | < 2 seconds |
| Validation | < 1 second |
| Memory per MUD | < 50 MB |

---

## Risk & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Content complexity too high | Medium | High | Start simple, iterate |
| Performance issues | Low | Medium | Profile early, cache aggressively |
| Backward compatibility | High | High | Dual loading mode during transition |
| Schema changes break content | Medium | Medium | Version schemas, migration scripts |

---

## Success Criteria

✅ All hardcoded skills moved to YAML
✅ All NPCs defined in templates
✅ All items externalized
✅ New MUD creatable by copying content folder
✅ Hot-reload working
✅ Admin UI for content editing
✅ Non-programmer can add content

---

## Week-by-Week Summary

| Week | Focus | Key Deliverable | Hours |
|------|-------|-----------------|-------|
| 1 | Architecture | Loader framework, schemas | 16 |
| 2 | Skills | All skills in YAML, registry | 18 |
| 3 | NPCs | Templates, spawn system | 16 |
| 4 | Items | Item definitions, equipment | 17 |
| 5 | Rooms | Area/room YAML, loader | 14 |
| 6 | Quests | Quest framework, files | 16 |
| 7 | Admin | Hot-reload, admin UI | 16 |
| 8 | Polish | Tests, docs, optimization | 16 |
| **Total** | | | **129 hours** (≈ 16 weeks @ 8 hrs/week) |

---

## Immediate Next Steps

### Week 1 Day 1
1. Create content directory structure
2. Design skill schema
3. Implement base loader

### Dependencies
- `gopkg.in/yaml.v3` for YAML parsing
- `github.com/xeipuuv/gojsonschema` for validation
- `github.com/fsnotify/fsnotify` for file watching

---

*Week 1 starts Monday. Ready to externalize! 🐢🟣*

**Cowabunga! The possibilities are endless with data-driven content!** 🟣
