---
name: admin-crawler
description: "Use when asked to crawl, smoke-test, or audit the herbst-mud admin panel (http://100.67.206.65:5173). Logs in as sma, drives admin workflows through the Vite/React UI, files tickets with root-cause analysis, applies or files fixes, writes Gherkin scenarios to features/, and writes Honcho conclusions at every phase. Trigger on phrases like '/admin-crawler', 'crawl the admin panel', 'test the admin UI', 'crawl N improvements on the admin'."
version: 1.0.0
author: Hermes Agent (Leonardo profile)
license: MIT
platforms: [linux]
metadata:
  hermes:
    tags: [qa, browser, debugging, tickets, admin, herbsmud, honcho, cucumber]
    related_skills: [ui-crawler, player-crawler, admin-qa-testing, honcho-cycle, dogfood]
---

# /admin-crawler — Admin Panel QA

## Overview

Systematic browser-based QA workflow for the **herbst-mud admin panel**
at `http://100.67.206.65:5173`. Logs in as the `sma` test user,
drives admin workflows through the Vite + React + TanStack Router UI,
identifies bugs with concrete root-cause analysis, and either fixes
them directly or files tickets in `tickets/<NAME>.md`. Every tool call
and finding is recorded in Honcho so future invocations of this skill
inherit context. **For every direct fix or filed ticket, also write a
Gherkin scenario to `features/admin-<entity>-<area>.feature`** so the
work is captured as a regression-guard test.

This skill is a **focused command wrapper** around the broader
`admin-qa-testing` skill. The umbrella skill covers manual scripts,
test configuration, and 12 documented failure patterns. This skill
adds:

1. The **explicit improvement count** loop (default 5, configurable)
2. The **cucumber scenario generation** step for every finding
3. The **Honcho conclusion** cycle at every phase
4. The **direct-DB scaffold** for non-character-creation crawls
5. A **short, repeatable** crawl loop for `/admin-crawler` triggers

The admin panel is a Vite/React/TanStack-Query SPA (`admin/src/`)
with file-based routing. Server log at `/tmp/herbst-web.log`.

## Trigger

- `/admin-crawler [task]` — single-task crawl, default 5 improvements
- `/admin-crawler [task] [N] improvements` — task with explicit count
- `/admin-crawler focus on [area]` — focus a single entity or page
- `/admin-crawler [entity]` — e.g. `/admin-crawler factions` crawls
  every factions page (list, new, detail, categories) and adjacent
  surfaces (cross-form data flow)
- `crawl the admin` / `smoke-test the admin panel` — natural language
  equivalents; default 5 improvements, no explicit task means "crawl
  every reachable entity"

If the user gives a task but no count, **default to 5 improvements**.
The number is a floor, not a ceiling — keep going if you find more.
If the user gives both task and count, follow the count.

## When to Use

Use this skill when:
- Sam asks to crawl, smoke-test, or audit the admin panel
- A new admin feature has shipped and needs verification
- The user wants to find N improvements to the admin experience
- The dev environment is at `http://100.67.206.65:5173`

**Do NOT use for:**
- Player web client testing (use `player-crawler`)
- SSH TUI testing (use the SSH client and the `dogfood` skill)
- Backend-only API testing (use `_manual_qa_/test-*.sh`)

## Honcho Integration (MANDATORY)

Sam's standing imperative: **every multi-step tool call in this skill
must be remembered in Honcho for future invocations to inherit.** This
is not optional. Honcho is the cross-session memory layer — without it,
the next `/admin-crawler` invocation starts cold.

### Three-beat Honcho cycle

1. **CHECK (start of crawl):**
   ```
   honcho_profile()                          # peer card snapshot
   honcho_search("admin panel bugs")
   honcho_search("admin regressions")
   honcho_reasoning("What did we last find crawling admin?", level=low)
   ```
   Capture prior findings so we don't re-discover them.

2. **OBSERVE (during crawl):**
   After EVERY user-facing action or significant tool call, write a
   `honcho_conclude(...)` with a one-line factual statement:
   - "Admin: /recipes new form submits to /api/recipes with the
     wrong content-type, server returns 415"
   - "Admin: faction member count on dashboard does not update
     after assigning a new player"
   - "Admin: world_id=default sentinel causes /api/races to 404
     even though races exist in world 2"
   Save **patterns and root causes**, not raw logs or commit SHAs.

3. **SAVE (end of crawl):**
   - Final session summary: pages crawled, findings filed, fixes shipped,
     tickets created, Gherkin scenarios written, what was deferred
   - One conclusion per major checkpoint, not one per tool call

**Failure mode: if `honcho_conclude` returns "Failed to save
conclusion." with no error body, the ollama embedding key on Sam's
side is expired. Do not loop. Surface it in the final report and
pin the standing imperative locally via the `memory` tool.**

## Credentials & Endpoints

| Field | Value |
|-------|-------|
| Admin URL | `http://100.67.206.65:5173` |
| API URL | `http://localhost:8080` |
| Tester email (real form value) | `sma` |
| Tester password | `sma` |
| Server log | `/tmp/herbst-web.log` |
| Admin dev log | `/tmp/herbst-web.log` (same) |
| Production world | **Ooze Surfers** (world_id 2) — see honcho cycle |
| Dev-only world | herbst-mud (world_id 1) — avoid creating data here |

**Email-vs-username trap:** the admin login form label says "Username"
but the backend `binding:"required,email"` accepts bare `sma` (no `@`
domain needed).

## Pre-Flight: Confirm the Active World

**Standing imperative (Sam, 2026-06-10):** before creating any
world-scoped entity via the admin panel (abilities, items, NPCs,
races, factions, etc.), confirm the WorldStore is on Ooze Surfers
(id 2). The default sentinel `"default"` maps to world 1 (herbst-mud,
dev-only) and silently creates data in the wrong world.

```bash
# On the Dashboard page, check the world dropdown value
browser_console(expression="document.querySelector('select')?.value")
# → MUST be "2" (Ooze Surfers), NEVER "1" or "default"
```

The `WorldStoreContext` writes the world **id number** to
`localStorage.herbst_current_world`. When the dropdown shows:
- "herbst-mud" → id 1 → dev-only
- "Ooze Surfers" → id 2 → production

**Always select Ooze Surfers (id 2) before creating production data.**

## Direct-DB Scaffold (when crawl isn't about character creation)

**Standing imperative (Sam, 2026-06-10):** unless the crawl's primary
goal IS creating characters, scaffold test data via direct DB inserts
and skip the React form fill. Use the scaffolder in
`player-crawler/references/scaffold-character.py` for any user-scoped
character. For admin-side world data (abilities, items, etc.), use
direct SQL.

```bash
DB_PW=$(grep DB_PASSWORD /home/sam/GitHub/herbst-mud/.env | cut -d= -f2)
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  INSERT INTO abilities (name, description, ability_type, ability_class,
    mana_cost, stamina_cost, hp_cost, cooldown, world_id, ...)
  VALUES (...)
  RETURNING id;
"
```

The player-crawler scaffolder script (`../../player-crawler/references/scaffold-character.py`)
also has a `list_chars` helper you can call as a module:

```bash
python3 -c "
import sys; sys.path.insert(0, '/home/sam/GitHub/herbst-mud/.agents/skills/player-crawler/references')
from scaffold_character import list_chars
import json
print(json.dumps(list_chars(), indent=2))
"
```

## Crawl Workflow

### Phase 0: Cross-Form Data-Flow Audit (when user names a shared concept)

If the crawl is about a concept that spans multiple admin pages
(e.g. "tags on abilities and races and the tags page"), first trace
the data flow:

1. **Grep for the concept name** across `server/routes/`, `admin/src/hooks/`,
   `server/db/schema/`. Look for entities, fields, tables, hooks, and
   components that all use the same word.
2. **Map the four possible storage forms:**
   - Canonical entity (e.g. `tag` table with m2m edges) — admin page exists
   - String field on another entity (e.g. `ability.required_tag`) — typed
     in a form field
   - String on a join table (e.g. `faction_required_tag.RequiredTag`) —
     hidden behind an admin flow
   - String array on an entity (e.g. `quest.rewards.tag_adds`) — multi-select
3. **Count: how many distinct things share the same name?** If >1, the user
   almost certainly wants them unified. File a meta-ticket that names all
   N concepts.
4. **Then** test the named page in Phases 1-3 below.

**Why**: Sam's "tags" request uncovered 4 separate "tag" concepts in
herbst-mud (entity + 3 string fields). Testing only `/tags` would
have missed the ability/race disconnect entirely.

### Phase 1: Login

```
browser_navigate(url="http://100.67.206.65:5173")
browser_snapshot()                         # get page structure
browser_console(clear=true)                # start clean
browser_type(ref="<username-input>", text="sma")
browser_type(ref="<password-input>", text="sma")
browser_click(ref="<login-button>")
browser_snapshot()                         # verify dashboard
browser_console()                          # check for errors
```

After login, **navigate to the Dashboard** and confirm the world
dropdown is on Ooze Surfers (id 2) before any creation work.

### Phase 2: Audit Plan

The admin panel's reachable entities (per `admin/src/routes/_auth/`):

| Top-level | Sub-pages |
|-----------|-----------|
| Dashboard | Overview stats, world selector |
| Content | NPCs, Items, Abilities, Effects, Skills, Quests, Recipes, Races, Genders, Triggers, Factions, Map |
| Players | Players, Characters, XP |
| Social | Socials, Channels |
| System | Config, Worlds, Tags, Logs, Docs |

### Phase 3: Per-Page Crawl

For each page, follow the standard 5-step (full reference:
`admin-qa-testing` skill, "Workflow: Testing Admin Forms"):

#### Step 1: Navigate & load
```
browser_navigate(url="http://100.67.206.65:5173/<route>")
browser_snapshot()
browser_console()                        # silent JS errors = high-value bugs
tail -3 /tmp/herbst-web.log | tr -d '\000'  # check backend received fetches
```

#### Step 2: Read (list / state)
- Page renders without errors?
- Data visible? Empty state if no data?
- Pagination / filtering work?

#### Step 3: Mutation (create / edit / delete)
- Use `browser_click` for first clicks on a stable page
- After ANY state change, use programmatic click via `browser_console`:
  ```js
  Array.from(document.querySelectorAll('button')).find(b => b.textContent.trim() === 'Save')?.click()
  ```
  Refs go stale after re-render.

- For React-controlled selects, prototype-setter pattern:
  ```js
  const sel = document.querySelector('select');
  const setter = Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype,'value').set;
  setter.call(sel, '2');
  sel.dispatchEvent(new Event('change', { bubbles: true }));
  ```

- For long forms, the form-fill helper from
  `admin-qa-testing/references/browser-automation-form-fill.md` is the
  workhorse. Install `__fill` and `__submit` per page, then call.

#### Step 4: Verify persistence
- Reload the list — is the change still there?
- `tail -3 /tmp/herbst-web.log` — did a real POST/PUT/DELETE fire?
- Cross-check via API:
  ```bash
  curl -s -H "Authorization: Bearer *** \
    "http://localhost:8080/api/abilities?world_id=2" | jq
  ```

#### Step 5: Edge cases
- Empty required field → should be rejected
- Negative numbers on cost/level fields
- XSS: `<script>alert(1)</script>` in name fields
- World switching: create in world 2, switch to world 1, verify isolation
- Concurrent edits: open same record in two tabs, edit both, save both

### Phase 4: Classify and File

#### Severity scale
- **P0 (crash)**: page doesn't render, form crashes, data loss
- **P1 (broken)**: feature doesn't work, silent failure, data
  inconsistency
- **P2 (cosmetic)**: visual issues, missing tooltips, UX papercuts

#### Direct-fix decision (Sam's standing rule)

> "If an improvement is not easy to fix in the web-client codebase,
> then create a ticket in `tickets/` to resolve it later."

- **Fix directly** if:
  - Typo, missing `key`, broken conditional
  - Missing tooltip or help text
  - Wrong CSS class or layout glitch
  - One-file prop wiring fix
  - **Caveat:** if the fix changes behavior visible to admins (e.g.
    auto-fill, default values), ask first

- **File a ticket** if:
  - Multi-file refactor (component split, hook extraction)
  - Touches the API contract
  - Changes game-balance values
  - Risk of regression in adjacent feature

When in doubt, file a ticket.

### Phase 5: MANDATORY — Write a Gherkin Scenario

**Sam's standing rule (2026-06-10):** when tasked with anything that
results in a direct fix or a filed ticket, **also create a cucumber-
formatted scenario in `features/`**. This is the regression-guard test
that prevents the bug from coming back.

The convention: `features/admin-<entity>-<area>.feature`. Append to
the existing entity file if one exists; create a new one otherwise.

Check first:
```bash
ls features/admin-*.feature | grep -i <entity>
# If file exists, extend it. If not, create a new one.
```

#### Scenario template (admin)

```gherkin
Feature: Admin <Entity> — <Crawled area>

  Background:
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"

  Scenario: <Reproduce the bug we found or verify the fix we applied>
    Given I navigate to "/<entity>"
    When I <user action>
    Then <expected outcome>

  Scenario: <Regression guard — verify the bug stays fixed>
    Given <state precondition>
    When <action>
    Then <expected outcome>
```

The scenarios should be:
- **Reproducible** — anyone reading can step through them
- **Regression-focused** — test the specific failure mode, not generic CRUD
- **World-aware** — include `Given the active world is "Ooze Surfers"`
  in the Background when the entity is world-scoped
- **Cross-form** — when a concept spans multiple pages, add scenarios
  for the cross-form propagation (e.g. "create a tag at /tags, verify
  it appears in the ability create form dropdown")

#### Example: bug we found in this session

If we find that `/admin/abilities/new` silently fails to POST because
of a Vite proxy miss:

```gherkin
# Append to features/admin-ability-crud.feature

  Scenario: Create an ability — admin panel submits to API correctly
    Given I am authenticated as an admin
    And the active world is "Ooze Surfers"
    When I navigate to "/abilities/new"
    And I fill the ability form with:
      | name            | TestCrawlerAbility |
      | ability_type    | combat             |
      | ability_class   | active             |
      | mana_cost       | 5                  |
    And I click "Save"
    Then the response status should be 201
    And "TestCrawlerAbility" should appear in "/abilities" list
```

#### Run the scenarios

```bash
cd /home/sam/GitHub/herbst-mud/server && go test -v -run TestFeatures
# or target a specific feature
go test -v -run TestFeatures godog -t "Admin Ability"
```

If the scenarios don't run (test harness missing for the new area),
the file is still valuable as documentation and as a manual-test
checklist. Note the gap in the session summary.

### Phase 6: Ticket Template (markdown in `tickets/<NAME>.md`)

```markdown
# [admin] <brief description>

## Problem Statement
<one paragraph: what breaks, how often, who hits it>

## Findings Table
| # | Severity | Title | File:Line | Detail |
|---|----------|-------|-----------|--------|

## Expected Outcome
1. <numbered: what good looks like>
2. <...>

## Files That Need Changes
| File | Change |
|------|--------|

## Non-Goals
- <out of scope>

## Verification Steps
- [ ] <how to confirm the fix>
- [ ] <...>

## Gherkin Scenario
See: `features/admin-<entity>-<area>.feature` — Scenario "<name>"
```

The Gherkin file reference at the end of the ticket makes the
regression-guard test discoverable from the ticket.

## Honcho Conclusions — Required

After every phase transition, write a `honcho_conclude` with one of
these shapes:

```
honcho_conclude(conclusion="Admin: <page> shows <symptom>; root cause: <file:line — short reason>")
honcho_conclude(conclusion="Admin: filed ticket TICKET-<NAME>-<NNN>.md and Gherkin scenario features/admin-<entity>-<area>.feature for <P0/P1/P2 issue>")
honcho_conclude(conclusion="Admin: shipped direct fix for <issue> in <file> (<N> lines); added Gherkin guard")
honcho_conclude(conclusion="Admin session complete: <M> findings, <N> tickets, <K> direct fixes, <G> Gherkin scenarios, <pages touched>")
```

Avoid imperative phrasing ("always check X") in conclusions. Past
tense, declarative facts.

## Common Failure Patterns (load from `admin-qa-testing`)

These 12 patterns are the most common admin-panel failures — load
that skill for the full reproduction recipes:

- **Pattern 1: Missing `entity.new.tsx` route file** — TanStack Router
  treats `/entity/new` as a catch-all param
- **Pattern 2: `Button` `disabled={undefined}` renders as disabled** —
  one-line destructuring fix
- **Pattern 3: JSX `{identifier}` in plain text** — crashes the page
- **Pattern 4: Form lacks self-documenting tooltips** — add `tooltip` props
- **Pattern 5/5b: `world_id="default"` sentinel** — 404 OR silent wrong
  world. Apply the 3-layer fix.
- **Pattern 6: 2-click delete confirm races with state update** — use a modal
- **Pattern 7: Schema truth vs derived graph** — Z-level compute bug in map
- **Pattern 8: Missing `<Outlet />` on parent list route** — child routes
  don't render
- **Pattern 9: Card-based sub-item editors** — UX recipe
- **Pattern 10: Subagent-shipped UI has a Vite/tsc gap** — tsc alone is
  not enough; load the page in a browser
- **Pattern 11: Off-screen buttons** — `browser_click` silently fails
  below the viewport
- **Pattern 12: Tailwind `flex` + `lg:block` overrides** — silent
  layout break

**Always load `admin-qa-testing` first when starting an
`/admin-crawler` run.** That skill is the umbrella; this one is the
command wrapper.

### Pattern 13: CORS_ORIGINS missing Tailnet IPs (P1 — browser login silently fails)

**Symptom:** The admin panel login form submits, the server log shows
`OPTIONS /users/auth → 204`, but no `POST` follows. The page stays on
the login form with no error message. `browser_console` shows no JS
errors.

**Root cause:** `.env` has `CORS_ORIGINS=http://localhost:3000,http://localhost:5173`
but the browser is accessing the admin panel via the Tailnet IP
`http://100.67.206.65:5173`. The CORS middleware checks the `Origin`
header against the allowed list — `100.67.206.65:5173` is NOT in the
list, so `Access-Control-Allow-Origin` is never set. The browser
receives the 204 preflight but blocks the actual POST because the
origin is not allowed.

**Detection:**
```bash
# Check what CORS_ORIGINS is set to
grep CORS_ORIGINS /home/sam/GitHub/herbst-mud/.env
# If it doesn't include the Tailnet IP + port you're accessing from, that's the bug

# Verify the preflight response is missing Allow-Origin
curl -v -X OPTIONS http://100.67.206.65:8080/users/auth \
  -H "Origin: http://100.67.206.65:5173" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" 2>&1 | grep Allow-Origin
# → empty = origin not allowed
```

**Fix:** Add all access URLs to `CORS_ORIGINS` in `.env`:
```
CORS_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:5174,http://100.67.206.65:5173,http://100.67.206.65:5174
```
Then restart the server (the env is read at startup).

**When to check:** Any time the browser can't log in but `curl` to the
API works fine. The OPTIONS preflight succeeds (204) but the POST
never fires. This also affects the player web client on port 5174.

### Pattern 14: WebSocket command routing — dead code functions never called (P1)

**Symptom:** A feature works in the codebase — the function exists,
compiles, and is well-documented — but the player never experiences
it. The server log shows no errors. The function is simply never
called from the command dispatcher.

**Root cause:** During refactors, command handlers get moved to
separate files (e.g. `ws_dialog_choice.go`, `ws_conversation.go`) but
the `handleCommand` function in `ws_routes.go` is never updated to
route the new command prefix to the new handler. The handler becomes
dead code.

**Real case (2026-07-12):** `handleDialogChoice` in
`ws_dialog_choice.go` was complete and correct — it fetched dialog
nodes, applied effects, and called `sendConversationScreen`. But
`handleCommand` in `ws_routes.go` had no route for the `"dialog "`
prefix. The client sent `dialog theodore_von_rad theodore_entry 1`
and the server returned "Command not yet implemented." The entire
NPC conversation system appeared broken.

**Detection:**
```bash
# Find all command handler functions in the routes package
grep -rn 'func handle\|func try\|func send' server/routes/ws_*.go | grep -v _test

# Check which are actually called from handleCommand
grep -n 'handle\|try\|send' server/routes/ws_routes.go | grep -E 'return |='
```

If a handler function exists in a `ws_*.go` file but is never
referenced from `handleCommand`, it's dead code. Wire it in.

**Fix pattern:** Add a prefix match in `handleCommand`:
```go
// dialog <template_id> <node_id> [<choice_index>]
if strings.HasPrefix(cmd, "dialog ") {
    parts := strings.Fields(cmd)
    if len(parts) < 3 { return "Invalid dialog command." }
    choiceStr := ""
    if len(parts) >= 4 { choiceStr = parts[3] }
    return handleDialogChoice(parts[1], parts[2], choiceStr, wsc, repos, client)
}
```

### Pattern 15: tryTalk only returns text, never opens conversation overlay (P1)

**Symptom:** Player types `talk <NPC>` and gets a text line like
`Theodore Von Rad says: "Hello there"` in the scrollback, but the
conversation overlay never opens. The NPC has dialog nodes in the DB
but they're never used.

**Root cause:** `tryTalk` in `ws_routes.go` fetches the NPC template
and returns `tmpl.Greeting` as a text string. It never checks for
dialog nodes or calls `sendConversationScreen`. The conversation
overlay code (`ws_conversation.go`) exists but is only called from
`handleDialogChoice` — which itself may be dead code (Pattern 14).

**Detection:**
```bash
# Check if tryTalk queries dialog nodes
grep -A5 'func tryTalk' server/routes/ws_routes.go | grep -i dialog
# Empty = tryTalk doesn't use the dialog system

# Verify the NPC has dialog nodes
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  SELECT id, npc_text, is_entry FROM dialog_nodes
  WHERE npc_template_dialog_nodes = (
    SELECT id FROM npc_templates WHERE id = '<template_slug>'
  );
"
```

**Fix:** After fetching the NPC template in `tryTalk`, query
`repos.DialogNode.ListByTemplate(ctx, tmpl.ID)`. If nodes exist, find
the entry node (`is_entry=true`) and call
`sendConversationScreen(wsc, tmpl.Name, tmpl.ID, nodes, entryID)`.
Return empty string (the screen payload is sent via WebSocket). Fall
back to greeting text only if no dialog nodes exist.

**Dialog node schema reference:**
- Table: `dialog_nodes`
- Fields: `id` (string, e.g. `theodore_entry`), `npc_text`,
  `responses` (JSON array of `{label, next_node_id, condition,
  quest_offer_id, decline_node_id, effects}`), `is_entry` (bool),
  `entry_condition` (SPICE expression), `on_enter_effects` (int array)
- Edge: `npc_template` → `NPCTemplate` (required, unique)
- Repo: `repos.DialogNode.ListByTemplate(ctx, templateID)`
- Screen sender: `sendConversationScreen(wsc, npcName, templateID, nodes, currentNodeID)`

## Server Restart After Go Code Changes

When a fix touches `server/*.go` files, the running binary must be
rebuilt and restarted:

```bash
# 1. Build
export PATH=$PATH:/usr/local/go/bin
cd /home/sam/GitHub/herbst-mud/server && go build -o herbst-web .

# 2. Kill old process
pkill -9 -f 'herbst-web'
sleep 2

# 3. Restart via the startup script
# (background=true for long-lived server process)
bash /tmp/run-herbst.sh > /tmp/herbst-web.log 2>&1 &

# 4. Verify
sleep 2
curl -sS -o /dev/null -w "healthz=%{http_code}" http://localhost:8080/healthz
# → healthz=200
```

The startup script at `/tmp/run-herbst.sh` sources `.env`, unsets
`DATABASE_URL` (to use localhost Postgres, not Docker), and runs
`exec server/herbst-web`. If the script doesn't exist, create it:
```bash
#!/bin/bash
set -a
source /home/sam/GitHub/herbst-mud/.env
set +a
unset DATABASE_URL
export DB_HOST=localhost
export DB_SSL_MODE=disable
exec /home/sam/GitHub/herbst-mud/server/herbst-web
```

## Commit and Release Workflow

When a fix is ready to ship:

1. **Stage the files** (never auto-commit per Sam's rule):
   ```bash
   cd /home/sam/GitHub/herbst-mud
   git add <files>
   git status  # verify only intended files are staged
   ```

2. **Draft the commit message** to a file:
   ```bash
   # Write to docs/plans/COMMIT_MSG_v<version>.md
   ```
   Format:
   ```
   🟣 fix(admin): <brief description>

   <body: what broke, root cause, what changed, files>
   ```

3. **Surface for Sam** — present the staged files and commit message.
   Sam runs:
   ```bash
   git commit -F docs/plans/COMMIT_MSG_v<version>.md
   git tag -a v<version> -m 'v<version>: <description>'
   git push origin main --tags
   ```

4. **Version bump** — check the latest tag:
   ```bash
   git describe --tags --abbrev=0  # latest version
   ```
   Increment the patch for bug fixes, minor for features.

5. **Server restart after commit** — if the fix touches server code,
   rebuild + restart as described above.

## References

- `references/audit-checklist.md` — condensed per-page checklist for
  fast crawls
- `references/world-store-debug.md` — WorldStoreContext gotchas, the
  three-layer world_id fix, sentinel detection
- `references/gherkin-recipes.md` — copy-paste Gherkin scenarios for
  common admin flows

## Verification Checklist (end of crawl)

- [ ] All required pages loaded without console errors
- [ ] Honcho conclusions written for: start, every phase transition,
      end
- [ ] Tickets filed for all P0/P1 findings (markdown in `tickets/`)
- [ ] **Gherkin scenarios written to `features/admin-*.feature` for
      every finding** (Sam's standing rule)
- [ ] Direct fixes applied, Gherkin regression guards added
- [ ] `cd admin && npx tsc --noEmit` clean (mandatory after any
      direct fix)
- [ ] **Open the dev page in a browser** and confirm no Vite overlay
      (tsc alone is insufficient — Pattern 10)
- [ ] `make build-all` clean if backend Go files changed
- [ ] **Server rebuilt and restarted** if any `server/*.go` file
      changed (see "Server Restart After Go Code Changes" section)
- [ ] **CORS_ORIGINS includes all access URLs** (localhost + Tailnet
      IPs for ports 5173 and 5174) — check if browser login fails
      but curl works (Pattern 13)
- [ ] **WebSocket command handlers wired** — verify all `func handle*`
      and `func try*` in `ws_*.go` are called from `handleCommand`
      (Pattern 14)
- [ ] **Export/import round-trip validation** if any backend file
      changed (server/, server/dbinit/, server/worldexport/). See
      `references/export-import-roundtrip.md`. Per Sam's standing
      rule (2026-07-10): always validate export → reimport whenever
      modifying the backend.
- [ ] **If fix shipped:** files staged, commit message drafted in
      `docs/plans/COMMIT_MSG_v<version>.md`, surfaced for Sam
- [ ] Final session summary `honcho_conclude` written
- [ ] No test data left in production world (Ooze Surfers / id 2)
      unless explicitly intended
