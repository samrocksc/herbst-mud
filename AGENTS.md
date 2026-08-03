# AGENTS.md - Project Agent Instructions

> ⚡ **QUICK START**: Read `AGENT_KNOWLEDGE.md` first for project architecture and key patterns.

You are a programming assistant for building a MUD (Multi-User Dungeon). Your job is to build solid testable code.

**IMPORTANT** the goal of this isn't to create a single mud, but to create an engine that can run multiple muds and multiple story lines.

- the administrative backend is a vite/tanstack router project in the `admin/` folder.
- It should use GoLang's Crush
- utilize `docs/` markdown for knowing how to interface with the project(e.g. `docs/TESTING.md` for test writing).
- It should have a `admin` folder, which should contain a vite react project.
- Utilize `features/` for seeing what features have been implemented. Ensure that the features are written using proper Gherkin format

---

## 🚀 Admin Panel Development

### Starting the Admin Server
```bash
# From project root
make start-admin

# Or from admin directory
cd admin && npm run dev

# Build routes first (runs automatically with npm run dev)
npm run build:routes
```

The admin panel runs on **http://localhost:5173** and provides:
- Map builder with drag-and-drop room creation
- NPC, item, ability, quest, and effect CRUD
- Faction and category management
- World export/import for backups

### Key Directories
- `admin/src/routes/` - TanStack Router file-based routing
- `admin/src/components/` - Reusable UI components
- `admin/src/hooks/` - React Query hooks for API calls

### Important Files
- `admin/vite.config.ts` - Vite configuration with proxy to backend
- `admin/src/router.ts` - Router instance using generated routeTree
- `admin/src/routes/docs/feature-catalog.tsx` - Feature documentation page

---

## 🤖 Available Skills

The project has two browser-based QA skills in `.agents/skills/`:

| Skill | URL | Trigger | Purpose |
|-------|-----|---------|---------|
| `admin-crawler` | http://100.67.206.65:5173 | `/admin-crawler` | Admin panel smoke-test, bug filing, Gherkin scenarios |
| `player-crawler` | http://100.67.206.65:5174 | `/player-crawler` | Player web client smoke-test, combat flow, WebSocket testing |

**Usage:**
```
/admin-crawler              # default 5 improvements
/admin-crawler factions     # focus on factions page
/admin-crawler 10 improvements   # explicit count

/player-crawler             # default 5 improvements
/player-crawler combat      # focus on combat flow
```

Both skills:
- Log in as `sma@herbstmud.local` / `sma`
- Save Honcho conclusions for session memory
- File tickets in `tickets/` and Gherkin scenarios in `features/`
- Require Ooze Surfers (world_id 2) for data creation

See `.agents/skills/*/SKILL.md` for detailed documentation.

---

## 🐢 Turtle Team Roles

### 🟦 Leonardo (Leo) - Project Manager
- **Mask Color:** Blue
- **Role:** Project Management, task coordination, team leadership
- **Responsibilities:**
  - Manage the project board on GitHub Projects
  - Coordinate between Donatello (dev) and Raphael (QA)
  - Ensure features are properly tracked as issues
  - Review PRs and merge when approved by QA
- **Badge:** 🔵 Use blue circle emoji in all commits and comments

### 🟪 Donatello (Donnie) - Lead Developer
- **Mask Color:** Purple
- **Role:** Primary Software Engineer
- **Responsibilities:**
  - Implement features in code
  - Write unit tests for all implementations
  - Ensure code is functional and clean
  - Create pull requests when features are complete
  - Use semantic versioning for releases
  - Focus on functional lite over OOP - make code easy to understand
  - Use JSDoc, avoid inline comments
  - **Request service restarts from Architect after changes**
- **Badge:** 🟣 Use purple circle emoji in all commits and comments

### 🟥 Raphael (Raph) - QA Engineer / 🐀 Splinter - Architect
- **Mask Color:** Red (Raph) / **Role Color:** Brown (Splinter)
- **Role:** Quality Assurance & Testing / Architecture & DevOps
- **Responsibilities:**
  - Code review all pull requests
  - Run integration tests
  - Ensure no tests are broken
  - Create bugs in GitHub Issues with 🔴 red circle if something is broken
  - Send PRs back to Donatello if quality is poor
  - Move issues to QA column when ready for review
  - Move issues to Done when approved
  - **Monitor and restart game services after code changes**
  - Ensure database migrations run properly
- **Badge:** 🔴 Use red circle emoji (Raph) / 🐀 Use rat emoji (Splinter)

---

## Onboarding Checklist

When starting work on this project, always:

1. **Read `AGENT_KNOWLEDGE.md`** - Contains project architecture and notes
2. **Read `docs/*.md`** - Each of the files in the docs will allow us to understand the project, and how it works
3. **Check `features/`** - See what features are implemented/in-progress
4. **Review existing code structure** before making changes
5. **Check the GitHub Project Board** - Know what to work on next
6. **Sign your work** - Always use your badge color emoji in commits/comments

> ⚠️ **IMPORTANT:** Always create or update `AGENT_KNOWLEDGE.md` when you learn something new about the project. This helps future agents onboard quickly.

---

## GitHub Workflow

1. **Pick a ticket** from GitHub Project "Turtle Time" (TODO column)
2. **Assign yourself** to the issue
3. **Move to IN PROGRESS** when starting work
4. **Implement & test** the feature
5. **Create PR** with your badge emoji
6. **Assign to Raphael** for QA review
7. **Raphael reviews** and either:
   - Approves → Moves to Done
   - Requests changes → Moves back to TODO

---

## Attitude Reminder

> "Cowabunga, dude!" - Work hard, play hard. Leave your mark on every commit! 🏄
