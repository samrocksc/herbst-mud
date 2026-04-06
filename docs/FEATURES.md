# Features

> 🔵 Last Updated: 2026-04-05

Feature tracking for the Herbst MUD project.

**GitHub Issues:** https://github.com/samrocksc/herbst-mud/issues  
**Project Board:** https://github.com/users/samrocksc/projects/2

---

## Quick Status Summary

| Status | Count | Description |
|--------|-------|-------------|
| 🟢 Complete | 12 | Features in production |
| 🟡 In Progress | 2 | Currently being developed |
| 🔴 Planned | 7 | On the roadmap |

---

## Complete Features 🟢

### Infrastructure

| # | Feature | Description | Status |
|---|---------|-------------|--------|
| 01 | Initial Server Scaffolding | Go REST API with Gin framework | ✅ Complete |
| 02 | Initial SSH Server Scaffolding | Go SSH server with Charmbracelet/Wish | ✅ Complete |
| 03 | Initial Admin Scaffolding | Vite + React + TanStack Router | ✅ Complete |
| 04 | OpenAPI Client Generation | Auto-generated TypeScript clients | ✅ Complete |
| 05 | Database Setup | PostgreSQL with Ent ORM | ✅ Complete |

### Core Game Systems

| # | Feature | Description | Status |
|---|---------|-------------|--------|
| 06 | CRUD Rooms | Create, read, update, delete game rooms | ✅ Complete |
| 14 | SSH Server Behavior | TUI rendering, commands, state management | ✅ Complete |

### Combat & Gameplay

| Feature | Description | Status |
|---------|-------------|--------|
| Combat System | Tick-based combat, damage calculation | ✅ Complete |
| Skills System | Learnable, improvable skills | ✅ Complete |
| Talents System | Swappable special abilities | ✅ Complete |
| Corpse System | Looting, searching bodies | ✅ Complete |
| NPC System | Interactive NPCs (Gizmo healing) | ✅ Complete |
| NPC Skills | AI-driven special abilities (Aragorn's heal) | ✅ Complete |
| Invincibility | 0 max HP = no damage taken | ✅ Complete |
| Immortal Mode | Takes damage but HP never below 1 | ✅ Complete |
| **Game Export/Import** | Export/import world data (rooms, NPCs, skills) | ✅ Complete |

### Deployment & Operations

| Feature | Description | Status |
|---------|-------------|--------|
| Docker Support | 3 Dockerfiles + docker-compose | ✅ Complete |
| Digital Ocean App Platform | .do/app.yaml spec | ✅ Complete |
| Neon DB Support | DATABASE_URL, SSL mode | ✅ Complete |
| CORS Configuration | Configurable origins | ✅ Complete |
| Rate Limiting | Per-IP request limiting | ✅ Complete |
| JWT Authentication | Environment-based secrets | ✅ Complete |

---

## In Progress 🟡

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 07 | CRUD Characters | 🟡 In Progress | Basic CRUD done, equipment integration ongoing |
| 08 | User CRUD Operations | 🟡 In Progress | Auth working, admin features pending |

---

## In Progress 🟡

| Feature | Description | Status |
|---------|-------------|--------|
| **Content Externalization** | YAML-based content for multi-MUD support | 🚧 Week 2 Complete |

### Content System Roadmap

**8-Week Plan:**

| Week | Status | Description |
|------|--------|-------------|
| 1 | ✅ Complete | Content architecture, schemas, loader framework |
| 2 | ✅ Complete | Skill externalization (5 classless skills) |
| 3 | ✅ Complete | NPC template system (5 templates + 2 NPC skills) |
| 4 | ✅ Complete | Item externalization (8 items + loot tables) |
| 5 | ✅ Complete | Room/area system (9 rooms + exit network) |
| 6 | ✅ Complete | **Quest system** - 4 quests with steps, rewards, NPC integration |
| 7 | ✅ Complete | Hot-reload + admin API (fsnotify watcher, validation, preview) | + admin API (fsnotify watcher, validation, preview) |
| 8 | ✅ Complete | Multi-world support (3 worlds: default, cyberpunk, steampunk) | |

**Week 6 Delivered:**
- ✅ 4 quest templates (tutorial, fetch, chain, combat)
- ✅ Quest steps with ordering and prerequisites
- ✅ Rewards: XP, items, reputation
- ✅ NPC quest offering (`quests_offered`)
- ✅ Quest validation against existing content

**Quest Types:**
- exploration: Visit rooms
- fetch: Collect items
- talk: Interact with NPCs
- chain: Multi-part with prerequisites
- kill: Defeat targets

**Current Content Totals:**
- 7 Skills
- 5 NPC Templates
- 8 Items
- 9 Rooms
- **4 Quests**

**Week 7 Proposed:**
- Hot-reload content changes
- Admin content editor API
- Content validation endpoints

---

## Planned Features 🔴

### Authentication & Characters

| # | Feature | Priority | Dependencies |
|---|---------|----------|--------------|
| 09 | Character Authentication | 🔴 High | #07, #08 |
| 10 | Character Creation | 🔴 High | #09, #19 |
| 11 | Character Login | 🔴 High | #09 |

### Data Structure

| # | Feature | Notes |
|---|---------|-------|
| 12 | Data Structure V1 | Entity relationships finalized |
| 13 | Room Navigation | Direction-based movement |

### Character Systems

| # | Feature | Status |
|---|---------|--------|
| 15 | User Entity | 🟡 Partial |
| 16 | Character Entity | 🟡 Partial |
| 17 | Room Entity | 🟢 Complete |
| 18 | Class System | 🔴 Planned |
| 19 | Race System | 🔴 Planned |
| 20 | Gender System | 🔴 Planned |

### Admin Panel Features

| Feature | Status | Description |
|---------|--------|-------------|
| Room API Integration | 🟡 In Progress | Visual room editor |
| Drag-Drop Room Creation | 🔴 Planned | Graphical room builder |
| Exit Edge Components | 🔴 Planned | Visual exit management |
| Room Edit Panel | 🔴 Planned | Room property editor |

---

## Feature Dependency Graph

```
Infrastructure
├── Database Setup ───┬── CRUD Rooms ───────┬── Room Navigation
├── Server Scaffolding  │                     └── Combat System
├── SSH Scaffolding     └── CRUD Characters ───┬── Skills System
└── Admin Scaffolding      ├── Authentication  ├── Talents System
                           ├── Creation        └── Equipment System
                           └── Login
```

---

## Definitions

### Status Meanings

| Symbol | Status | Definition |
|--------|--------|------------|
| 🟢 | Complete | Feature is implemented, tested, and merged |
| 🟡 | In Progress | Feature is being actively developed |
| 🔴 | Planned | Feature is designed but not started |
| ⚪ | Backlog | Feature idea, not designed yet |

### Priority Levels

| Level | Description |
|-------|-------------|
| Critical | Blocks other work or production deployment |
| High | Important for user experience |
| Medium | Enhancement, nice to have |
| Low | Future consideration |

---

## Contributing Features

### Adding a New Feature

1. Create a GitHub Issue with feature description
2. Add to project board under "To Do"
3. Write Gherkin BDD test in `features/XX_feature_name.feature`
4. Assign to developer
5. Move to "In Progress" when starting
6. Create PR when ready
7. QA reviews and moves to "Done" when approved

### Feature Format

```gherkin
Feature: Feature Name

  Background:
    Given the system is initialized

  Scenario: Main success scenario
    Given some precondition
    When an action happens
    Then expected outcome occurs

  Scenario: Error handling
    Given some precondition
    When invalid action happens
    Then error is returned
```

---

## Testing

All features require:
- Gherkin BDD tests in `features/`
- Unit tests in `*_test.go` files
- Integration tests where applicable

See `docs/GHERKIN_TESTING.md` for testing guidelines.

---

🔵 Document version: 2026-04-05
