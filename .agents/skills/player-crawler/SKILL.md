---
name: player-crawler
description: "Use when asked to crawl, smoke-test, or audit the herbst-mud player web client (http://100.67.206.65:5174). Logs in as the sma tester, drives a player session through the React+WebSocket UI, files tickets, applies or files fixes, and writes Honcho conclusions at every phase. Trigger on phrases like '/player-crawler', 'crawl the player UI', 'test the web client', 'crawl N improvements on the player client'."
version: 1.0.0
author: Hermes Agent (Leonardo profile)
license: MIT
platforms: [linux]
metadata:
  hermes:
    tags: [qa, browser, debugging, tickets, web-client, player, herbsmud, honcho]
    related_skills: [ui-crawler, admin-crawler, admin-qa-testing, honcho-cycle, dogfood]
---

# /player-crawler — Player Web Client QA

## Overview

Systematic browser-based QA workflow for the **herbst-mud player web client**
at `http://100.67.206.65:5174`. Logs in as the `sma` test user, drives a
full player session through the React+WebSocket UI, identifies bugs with
concrete root-cause analysis, and either fixes them directly or files
tickets in `tickets/<NAME>.md`. Every tool call and finding is recorded
in Honcho so future invocations of this skill inherit context.

The web client is a vanilla React + Vite + Tailwind v4 SPA
(`web-client/src/`) with WebSocket-driven gameplay
(`useWebSocket` → `useMUDSocket`). No TanStack Router — screens are
state-driven from `App.tsx`. Backend API at `http://localhost:8080`.

## Trigger

- `/player-crawler [task]` — single-task crawl, default 5 improvements
- `/player-crawler [task] [N] improvements` — task with explicit count
- `/player-crawler focus on [area]` — focus a single subsystem (combat,
  equipment, character creation, hotkeys, etc.)
- `crawl the player UI` / `smoke-test the web client` — natural language
  equivalents; default 5 improvements, no explicit task means "crawl
  every reachable screen"

If the user gives a task but no count, **default to 5 improvements**.
The number is a floor, not a ceiling — keep going if you find more.
If the user gives both task and count, follow the count.

## When to Use

Use this skill when:
- Sam asks to crawl, smoke-test, or audit the player web client
- A new web-client feature has shipped and needs verification
- The user wants to find N improvements to the player experience
- The dev environment is at `http://100.67.206.65:5174`

**Do NOT use for:**
- Admin panel testing (use `admin-crawler` or `admin-qa-testing`)
- SSH TUI testing (use the SSH client and the `dogfood` skill)
- Backend-only API testing (use `_manual_qa_/test-*.sh`)

## Honcho Integration (MANDATORY)

Sam's standing imperative: **every multi-step tool call in this skill
must be remembered in Honcho for future invocations to inherit.** This
is not optional. Honcho is the cross-session memory layer — without it,
the next `/player-crawler` invocation starts cold.

### Three-beat Honcho cycle

1. **CHECK (start of crawl):**
   ```
   honcho_profile()                    # peer card snapshot
   honcho_search("player client bugs")
   honcho_search("web-client regressions")
   honcho_reasoning("What did we last find crawling the player UI?", level=low)
   ```
   Capture prior findings so we don't re-discover them.

2. **OBSERVE (during crawl):**
   After EVERY user-facing action or significant tool call, write a
   `honcho_conclude(...)` with a one-line factual statement:
   - "Player UI: /combat log shows 2 stale lines after equipment change"
   - "Player UI: POST /api/characters/2/abilities returns 400 when slot already taken — silent no-op"
   - "Player UI: HotkeyBar key '3' overlaps with browser Ctrl+Shift+J"
   Save **patterns and root causes**, not raw logs or commit SHAs.

3. **SAVE (end of crawl):**
   - Final session summary: pages crawled, findings filed, fixes shipped,
     tickets created, what was deferred
   - One conclusion per major checkpoint, not one per tool call

**Failure mode: if `honcho_conclude` returns "Failed to save
conclusion." with no error body, the ollama embedding key on Sam's
side is expired. Do not loop. Surface it in the final report and
pin the standing imperative locally via the `memory` tool.**

## Credentials & Endpoints

| Field | Value |
|-------|-------|
| Player URL | `http://100.67.206.65:5174` |
| API URL | `http://localhost:8080` |
| Tester username | `sma` |
| Tester password | `sma` |
| Tester email (real form value) | `sma` |
| Default character (Ooze Surfers) | `smack` (faerie, Level 24, trash_mage) |
| Server log | `/tmp/herbst-web.log` |
| Web client dev log | `/tmp/herbst-client.log` |
| Dev world | **Ooze Surfers** (world_id 2) — see honcho cycle |
| sma user_id in DB | `3` |
| World Ooze Surfers rooms | id 1 "A dank ass closet", 2 "Sewer Junction", 3 "A dusty street" |

**Default character selection (Sam, 2026-07-12):** unless the crawl's
task explicitly requires creating or testing a different character,
**always select `smack`** when entering Ooze Surfers. Do not create
new test characters unless the task demands it — the `TestHero*`
characters littering room 1 are from unnecessary scaffolding runs.

**Backend path convention (verified 2026-06-10):** the herbst-mud
backend does **NOT** use an `/api/` prefix on its routes. The paths
are bare under the server root:
- `POST /users/auth` — login (returns JWT)
- `GET  /user-characters/:id` — list a user's characters
- `POST /user-characters/:id` — create a character (field: `world`,
  not `current_world`)
- `GET  /characters/:id` — get one character
- `DELETE /characters/:id` — delete one character
- `GET  /races?world_id=<n>`, `GET /genders?world_id=<n>` — entity lists

The web client (`web-client/src/lib/api.ts`) uses these bare paths
correctly. The admin panel adds a Vite proxy that rewrites the
`/api/*` prefix to the bare paths. Don't get bitten by the
asymmetry — when calling the API from a script, use the bare paths.

**Email-vs-username trap:** the player LoginScreen input is labeled
"Username / Email" and calls `POST /users/auth` with whatever you
type. The backend accepts `sma` directly (no `@` domain needed).

## Pre-Flight: Character Scaffolding

**Standing imperative (Sam, 2026-06-10):** unless the crawl's primary
goal IS character creation, scaffold the test character via direct DB
insert or `user-characters` API — never spend 5 minutes stepping
through the React form for setup. We have observed:
- Forms with two `type=submit` buttons (inner subform + outer save)
  require filtering by visible text — picking the first silently
  submits the wrong form
- WorldStore sentinel `"default"` leaks to the create call (Pattern
  5b in `admin-qa-testing`) — created chars land in world 1, not 2
- Required fields (`race`, `gender`, `name`) must match **name
  strings**, not numeric ids: `race='Ooze'`, `gender='It'` for Ooze
  Surfers (world 2)

### Recommended: scaffold via API + DB

```bash
# 1. Get a token for sma (no /api prefix on the herbst-mud backend)
TOKEN=$(curl -s -X POST http://localhost:8080/users/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"sma@herbstmud.local","password":"sma"}' | jq -r .token)

# 2. Check existing characters (do NOT create duplicates)
curl -s -H "Authorization: Bearer *** \
  http://localhost:8080/user-characters/2 | jq

# 3a. Create via API (use this path normally)
#    The field is `world` (not `current_world`) and the race/gender
#    must be the **name string** ("Ooze", "It"), not a numeric id.
curl -s -X POST -H "Authorization: Bearer *** \
  -H "Content-Type: application/json" \
  -d '{"name":"CrawlerA","race":"Ooze","gender":"It","world":"2"}' \
  http://localhost:8080/user-characters/2 | jq

# 3b. OR direct DB insert (faster, no form pipeline to debug,
#     bypasses the server's letter-only name validation)
DB_PW=$(grep DB_PASSWORD .env | cut -d= -f2)
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  INSERT INTO characters
    (name, race, gender, level, hitpoints, max_hitpoints, current_world,
     starting_room_id, respawn_room_id, current_room_id, user_characters)
  VALUES
    ('Crawler-1', 'Ooze', 'It', 1, 100, 100, '2', 1, 1, 1, 2)
  RETURNING id, name, current_room_id, current_world;
"
```

The Python scaffolder (`references/scaffold-character.py` in this skill)
does this end-to-end with a single command. Use it.

## Crawl Workflow

### Phase 1: Login

**Defaults:** username `sma`, password `sma`. After world selection,
pick character **`smack`** in Ooze Surfers unless the task says otherwise.

```
browser_navigate(url="http://100.67.206.65:5174")
browser_snapshot()
browser_console(clear=true)
browser_type(ref="<email-input>", text="sma")
browser_type(ref="<password-input>", text="sma")
browser_click(ref="<login-button>")
browser_snapshot()                       # verify dashboard/world screen
browser_console()                        # check for errors
# Select Ooze Surfers (world 2)
# Select character "smack"
```

**Detection before wasting turns:**
```js
document.body.innerText.includes("'Email' Error")
// → if true, the form is working, your input is wrong (used bare "sma")
```

**World selection:** after login, the WorldScreen appears (if multiple
worlds exist for the user). **Always select Ooze Surfers (id="2")** —
world 1 (`herbst-mud`) is dev-only.

### Phase 2: Audit Plan

The web client's reachable screens (state-driven from `App.tsx`):

| Screen | File | Routes to test |
|--------|------|---------------|
| Login | `LoginScreen.tsx` | form validation, theme toggle, error states |
| World select | `WorldScreen.tsx` | list, switch, empty state |
| Character list / create | `CharacterScreen.tsx`, `CreateCharacterScreen.tsx` | list, create, delete, validation |
| Game (main) | `GameScreen.tsx`, `RoomScreen.tsx` | movement, look, inventory, channels |
| Combat | `CombatScreen.tsx`, `CombatHUD.tsx`, `CombatActionBar.tsx`, `CombatLog.tsx`, `CombatVitals.tsx`, `CombatTargetList.tsx` | attack, abilities, hotkeys, target switching, death/revive |
| Equipment | `EquipmentScreen.tsx` | equip/unequip, slot validation, broken items |
| Conversation | `ConversationOverlay.tsx` | NPC dialog, choices, exit |
| Input bar | `InputBar.tsx` | text commands, history, autocomplete |
| Hotkey bar | `HotkeyBar.tsx` | key bindings, conflicts, tooltips |
| Scrollback | `Scrollback.tsx` | output rendering, ANSI colors, scroll-to-bottom |

### Phase 3: Per-Screen Crawl

For each screen, follow the standard 5-step:

#### Step 1: Navigate & load
```
browser_navigate(url="http://100.67.206.65:5174")  # or trigger via in-app nav
browser_snapshot()
browser_console()                        # silent JS errors are high-value bugs
tail -3 /tmp/herbst-web.log | tr -d '\000'  # check backend received any fetches
```

#### Step 2: Read (list / state)
- Page renders without errors?
- Data visible? Empty state if no data?
- Hotkey bindings shown? Tooltips populated?

#### Step 3: Mutation (create / equip / attack)
- Use `browser_click` for first clicks on stable page
- After ANY state change (list refresh, panel swap), use programmatic
  click via `browser_console`:
  ```js
  Array.from(document.querySelectorAll('button')).find(b => b.textContent.trim() === 'Equip')?.click()
  ```
  Refs go stale after re-render — see `admin-qa-testing` Pattern "Stale
  ref after re-render."

- For React-controlled selects, use prototype setter + change event:
  ```js
  const sel = document.querySelector('select');
  const setter = Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype,'value').set;
  setter.call(sel, '2');
  sel.dispatchEvent(new Event('change', { bubbles: true }));
  ```

#### Step 4: Verify persistence
- Reload screen — is the change still there?
- `tail -3 /tmp/herbst-web.log` — did a real POST/PUT/DELETE fire?
- Cross-check via API:
  ```bash
  curl -s -H "Authorization: Bearer *** \
    http://localhost:8080/characters/<id> | jq
  ```

#### Step 5: Edge cases
- Empty input → should be rejected (browser native + form validation)
- Negative values on number fields
- XSS: type `<script>alert(1)</script>` into name fields
- Network failure simulation (turn off the backend, see what client shows)
- WebSocket disconnect (kill server mid-session, see if client reconnects)

### Phase 4: Combat Specifics (the most-broken surface)

Combat is where the player client earns its keep. Crawl these flows:

1. **Initiate combat** — walk into a hostile NPC's room
2. **Target selection** — switch targets via `CombatTargetList`; verify
   vitals update
3. **Basic attack** — attack with no ability equipped
4. **Ability use** — equip an ability, verify cooldown, verify effect
   message rendering
5. **Death and respawn** — die in combat, verify `is_alive` flips,
   verify respawn_room_id is honored
6. **Flee** — flee from combat, verify the screen transitions back
7. **Hotkey conflicts** — `HotkeyBar.tsx` likely uses digits 1-5;
   check whether browser shortcuts override (`Ctrl+1` switches tabs)

The combat log (`CombatLog.tsx`) is the audit trail. Verify:
- Timestamps are present
- Color codes from `combat.ts` are applied
- Old lines fall off after some N (test with 200+ events)
- HTML in messages is **escaped** (XSS via NPC name)

### Phase 5: WebSocket Health

The whole game runs over a single WebSocket. The `useWebSocket` hook
in `src/lib/websocket.ts` is the only client → server command pipe.
Test the failure modes:

```bash
# Simulate disconnect: kill the server, watch the client
ps aux | grep herbst-web | grep -v grep
kill <pid>
# client should show "Reconnecting..." or equivalent
# restart: make dev
```

Verify reconnection logic exists. If it doesn't, **P0 ticket**:
"Player client does not auto-reconnect on WebSocket drop."

### Phase 6: Classify and File

#### Severity scale
- **P0 (crash)**: page doesn't render, form crashes, data loss, WebSocket
  can't connect at all
- **P1 (broken)**: feature doesn't work, silent failure, data
  inconsistency, missing critical UI element
- **P2 (cosmetic)**: visual issues, missing tooltips, UX papercuts

#### Ticket template (markdown in `tickets/<NAME>.md`)

```markdown
# [player] <brief description>

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
```

**Sam's required pattern (see memory):** Problem Statement →
Findings Table with severity (P0/P1/P2) → Expected Outcome as
numbered list → Files That Need Changes → Non-Goals → Verification
Steps. No abstract advice.

#### Direct-fix decision

Sam's imperative: "If an improvement is not easy to fix in the
web-client codebase, then create a ticket in `tickets/` to resolve it
later." This means:

- **Fix directly** if it's:
  - A typo, missing `key`, broken conditional
  - A missing tooltip or help text
  - A wrong CSS class or layout glitch
  - A one-file prop wiring fix
  - **Caveat:** if the fix changes behavior visible to players
    (e.g. auto-equip, default race), always ask before applying

- **File a ticket** if it's:
  - Multi-file refactor (component split, hook extraction)
  - Touches the WebSocket protocol
  - Changes a game-balance value
  - Requires backend coordination (API contract change)
  - Risk of regression in adjacent feature

When in doubt, file a ticket. Sam will accept the ticket and decide.

## Honcho Conclusions — Required

After every phase transition, write a `honcho_conclude` with one of
these shapes:

```
honcho_conclude(conclusion="Player UI: <screen> shows <symptom>; root cause: <file:line — short reason>")
honcho_conclude(conclusion="Player UI: filed ticket TICKET-<NAME>-<NNN>.md for <P0/P1/P2 issue>")
honcho_conclude(conclusion="Player UI: shipped direct fix for <issue> in <file> (<N> lines)")
honcho_conclude(conclusion="Player UI session complete: <M> findings, <N> tickets, <K> direct fixes, <world/pages touched>")
```

Avoid imperative phrasing ("always check X") in conclusions. Past
tense, declarative facts.

## Gherkin Scenarios (MANDATORY for every finding)

**Sam's standing rule (2026-06-10):** when a player-client finding
results in a direct fix or a filed ticket, **also create a cucumber-
formatted scenario in `features/`.** This is the regression-guard test
that prevents the bug from coming back.

The convention: `features/player-<screen>-<area>.feature` (or
`features/web-client-<area>.feature`). Check first — extend an
existing file or create a new one:

```bash
ls features/player-*.feature 2>/dev/null
ls features/web-client-*.feature 2>/dev/null
# Extend if exists, else create a new file
```

### Scenario template (player)

```gherkin
Feature: Player Web Client — <Crawled area>

  Background:
    Given I am authenticated as sma
    And I have a character in Ooze Surfers

  Scenario: <Reproduce the bug we found or verify the fix we applied>
    Given I navigate to "http://100.67.206.65:5174/<screen>"
    When I <user action>
    Then <expected outcome>

  Scenario: <Regression guard — verify the bug stays fixed>
    Given <state precondition>
    When <action>
    Then <expected outcome>
```

Scenarios should be:
- **Reproducible** — anyone reading can step through them
- **Regression-focused** — test the specific failure mode, not generic CRUD
- **World-aware** — include `Given I have a character in Ooze Surfers`
  in the Background
- **State-machine aware** — combat, equipment, dialog all have state
  transitions; the scenario should walk the full transition

### Example: combat bug

```gherkin
# features/player-combat-hotkeys.feature

  Scenario: Hotkey 1 fires equipped ability
    Given I am in combat with target "Ooze Grunt"
    And slot 1 has the ability "Haymaker"
    When I press "1"
    Then the combat log should show "You use Haymaker"
    And the cooldown for "Haymaker" should be 3 seconds
```

### Run the scenarios

```bash
cd /home/sam/GitHub/herbst-mud/server && go test -v -run TestFeatures
```

The player-client tests may not be wired into the godog harness yet
(player flow is WebSocket-driven, not pure REST). The Gherkin file is
still valuable as a manual-test checklist and as a regression note.
Surface any harness gap in the session summary.

### Ticket template (with Gherkin reference)

```markdown
# [player] <brief description>

## Problem Statement
<one paragraph: what breaks, how often, who hits it>

## Findings Table
| # | Severity | Title | File:Line | Detail |
|---|----------|-------|-----------|--------|

## Expected Outcome
1. <numbered: what good looks like>

## Files That Need Changes
| File | Change |
|------|--------|

## Gherkin Scenario
See: `features/player-<screen>-<area>.feature` — Scenario "<name>"
```

## Cleanup (MANDATORY at end of crawl)

**Standing rule: delete the test character when done, but NEVER the
user `sma`.** The user explicitly added this guard.

```bash
# Find the test character(s) — they were named "Crawler-N" or
# have a marker
DB_PW=$(grep DB_PASSWORD .env | cut -d= -f2)
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  SELECT id, name, current_world, world_id, level
  FROM characters
  WHERE user_characters=2 AND is_npc=false
    AND (name LIKE 'Crawler-%' OR name LIKE 'QA-%' OR name LIKE 'Test-%')
  ORDER BY id;
"

# Delete by ID (NOT by name — never by user)
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  DELETE FROM characters
  WHERE id IN (<ids>) AND user_characters=2;
"

# Verify the user still exists
PGPASSWORD=*** psql -h localhost -U herbst -d herbst_mud -c "
  SELECT count(*) AS sma_still_present FROM users WHERE email='sma@herbstmud.local';
"
# → MUST be 1
```

The character may also be deletable via the API (no exposed
`deleteCharacter` in `web-client/src/lib/api.ts` currently — check
before the crawl and use direct DB if missing). Add a
`honcho_conclude` recording the cleanup IDs.

## Common Failure Patterns (load from `admin-qa-testing`)

These patterns hit the player client too — load that skill for the
full reproduction recipes:

- **Pattern 5 / 5b: world_id="default" sentinel** — create entities in
  world 1 by accident. Always confirm Ooze Surfers (id 2) is the active
  world before any create.
- **Pattern 10: Vite overlay from import gaps** — `npx tsc --noEmit`
  alone is not enough. Open the page in a browser to confirm no
  "Failed to resolve import" overlay.
- **Pattern 11: off-screen buttons** — `browser_click` does nothing
  silently if the button is below the viewport. Use
  `btn.scrollIntoView()` + `btn.click()`.
- **Pattern "Stale ref after re-render"** — refs go stale after state
  transitions. Use programmatic click via `browser_console` for any
  click after a re-render.

## Troubleshooting Regressions

When Sam reports "this used to work, it's a regression":

1. **Git history first** — `git log --oneline --all -p -- <file>` for
   the suspect file. Look for removed logic that previously guarded the
   behavior. The diff history tells you what was lost and when.
2. **Session search** — `session_search(query="the behavior name")` to
   find the session where the fix was originally applied. Read the
   bookend messages for context on why the fix was needed.
3. **Find the current code path** — `search_files` with content regex
   for the function/handler that should be filtering. Check if the
   filter is still there or was silently dropped during a refactor.
4. **Verify via the API, not just the UI** — for WebSocket state like
   room contents, connect programmatically:
   ```python
   import json, asyncio, websockets
   async def check():
       async with websockets.connect(f"ws://localhost:8080/ws?token={token}&character_id={cid}") as ws:
           for _ in range(5):
               msg = await asyncio.wait_for(ws.recv(), timeout=5)
               data = json.loads(msg)
               if data.get("type") == "screen":
                   print(json.dumps(data["data"]["characters"], indent=2))
                   break
   asyncio.run(check())
   ```
   This bypasses browser form issues and gives you the raw server
   response. Use it to isolate "is this a server bug or a client bug?"
5. **Root cause before fix** — never patch the symptom. If offline PCs
   show in a room, find WHERE the list is built (e.g.
   `buildRoomScreen` → `ListByRoom`) and check if the filter was there
   before. File the Gherkin regression guard even for direct fixes.

### Regression found and fixed (2026-07-12)

- **Bug:** Offline player characters appeared in room character lists.
  All PCs with `current_room_id` set were returned by `ListByRoom` with
  no connection check.
- **Root cause:** `buildRoomScreen` in `ws_routes.go` called
  `repos.Character.ListByRoom()` which returns ALL characters in the
  room. The `connections` map (user_id → WSConn) tracks active WebSocket
  sessions but `buildRoomScreen` never consulted it. A previous fix
  that filtered offline PCs was lost during the ws_routes refactor
  (connection management was rewritten, old close logic removed).
- **Fix:** Added `connectedCharacterIDs()` that snapshots the
  connections map to a `map[int]bool` of connected character IDs.
  `buildRoomScreen` and the `examine` command filter via
  `isCharacterConnected(ch, connected)`. NPCs always pass. `tryTalk`
  and `tryAttack` already filtered by `ch.IsNPC` so they were
  unaffected.
- **Cleanup:** 12 `TestHero*` characters and 25 test users purged from
  the DB. Room 1 went from 15 visible chars to 3 (ChefHuman + smack +
  Theodore Von Rad NPC), with only smack + the NPC visible when smack
  is connected.
- **Gherkin:**
  `features/player-room-offline-pc-visibility.feature`

## Commit and Release Workflow

When a fix is ready to ship:

1. **Stage the files** (never auto-commit per Sam's 2026-06-30 rule):
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
   🟣 fix(player): <brief description>

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
   the running binary needs a rebuild + restart:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   cd /home/sam/GitHub/herbst-mud/server && go build -o herbst-web .
   # kill old process, restart with /tmp/run-herbst.sh
   ```

## References

- `references/scaffold-character.py` — Python scaffolder for the test
  character (avoids the 5-minute form-fill loop)
- `references/websocket-protocol.md` — `ClientMessage` and
  `ServerMessage` types from `web-client/src/types.ts`
- `references/character-deletion.md` — exact SQL/API for safe character
  cleanup at crawl end
- `references/known-bugs.md` — running list of player-client findings
  (append-only, scraped from Honcho)

## Verification Checklist (end of crawl)

- [ ] All required pages loaded without console errors
- [ ] Honcho conclusions written for: start, every phase transition,
      cleanup
- [ ] Tickets filed for all P0/P1 findings (markdown in `tickets/`)
- [ ] **Gherkin scenarios written to `features/player-*.feature` for
      every finding** (Sam's standing rule)
- [ ] Direct fixes applied and verified with browser reload
- [ ] `npx tsc --noEmit` clean in `web-client/` (mandatory after any
      direct fix)
- [ ] **Export/import round-trip validation** if any backend file
      changed (server/, server/dbinit/, server/worldexport/). See
      `references/export-import-roundtrip.md`. Per Sam's standing
      rule (2026-07-10): always validate export → reimport whenever
      modifying the backend.
- [ ] Test character(s) deleted
- [ ] sma user still present in DB (1 row in users table)
- [ ] Final session summary `honcho_conclude` written
- [ ] **If fix shipped:** files staged, commit message drafted in
      `docs/plans/COMMIT_MSG_v<version>.md`, surfaced for Sam
- [ ] **If server code changed:** binary rebuilt and server restarted
      via `/tmp/run-herbst.sh`
