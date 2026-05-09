# Claude Code CLI as Autonomous Coding Agent for herbst-mud — PRO Analysis

**Date:** 2026-05-06
**Author:** Research Analyst (Hermes subagent)
**Version:** Claude Code v2.1.131
**Decision:** Should herbst-mud adopt Claude Code CLI as its autonomous coding agent?

---

## Executive Summary

Claude Code CLI is the strongest candidate for autonomous coding on herbst-mud. It already has installed presence (v2.1.131 at `/home/sam/.local/bin/claude`), partial config (`.claude/settings.json`, `.claude/settings.local.json`), and — critically — it natively understands Go projects with ent ORM. Where pi-coding-agent achieved 0% success on Go backend tasks, Claude Code's tool-use model, worktree isolation, and hook system directly address every known pain point: dual `ent generate`, branch hygiene, service restarts, and serial ticket workflow. This document argues the PRO case in full.

---

## 1. Top 5 Strengths of Claude Code for herbst-mud

### 1.1 Native Go + ent ORM Understanding

Claude Code is trained on extensive Go codebases and understands the ent ORM pattern deeply. It can:
- Read and modify ent schema files (`server/db/schema/*.go`, `herbst/db/schema/*.go`)
- Run `go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema` in the correct directories
- Understand the two-module structure (`herbst/go.mod` and `server/go.mod` are separate Go modules)
- Navigate Go project conventions (ent client initialization, resolver patterns, auto-migration)

**Why this matters:** pi-coding-agent's 0% Go success rate proves that not all LLM tools handle Go well. Claude Code's model architecture is specifically tuned for multi-tool coding workflows involving compilation, testing, and iterative debugging — exactly the loop herbst-mud requires.

### 1.2 Print Mode (`-p`) for Autonomous Agent Orchestration

The `--print` / `-p` flag runs Claude in non-interactive mode: it takes a prompt, executes it using tools (file reads, bash commands, edits), and returns output to stdout. This is the correct mode for autonomous agents because:

- **No TTY required** — can be invoked by Hermes `delegate_task`, cron, or shell scripts
- **Structured output** — `--output-format json` or `--output-format stream-json` for machine-parseable results
- **Cost controls** — `--max-budget-usd` and `--max-turns` cap runaway spend
- **Session persistence** — `--resume` and `--continue` allow picking up where an interrupted task left off

**Comparison:** Interactive mode (`claude` with no flags) is for human pairing. Print mode is for agents. For herbst-mud's workflow where Leonardo (PM) dispatches tickets to Donatello (build), print mode is the programmatic entry point.

### 1.3 Git Worktrees (`-w`) Solve Branch Hygiene

The `--worktree` flag creates a fresh git worktree for each session:

```bash
claude -p "Implement SKILL-005: add faction schema" -w skill-005-faction-schema
```

This:
- Creates a new worktree at `.git-worktrees/skill-005-faction-schema/`
- Checks out a new branch automatically
- Keeps the main working tree untouched
- Prevents rebase hell when multiple agents (or Claude sessions) work on different tickets
- Each worktree gets its own `node_modules/`, build artifacts, and Go module cache

**Why this matters for herbst-mud:** With ~52 stale merged branches and a strict "one ticket at a time, no parallel PRs" rule, worktrees enforce isolation without manual branch management. When Claude finishes a ticket in worktree `X`, the main tree is pristine for the next ticket. No `git stash`, no `git rebase`, no accidentally committing to `main`.

### 1.4 Hooks System for Enforcing Critical Workflows

Claude Code's hooks allow running custom commands before or after tool use:

**PostToolUse hook for ent generation (the #1 pain point):**

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [
          {
            "type": "command",
            "command": "cd /home/sam/GitHub/herbst-mud/server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema && cd /home/sam/GitHub/herbst-mud/herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema"
          }
        ]
      }
    ]
  }
}
```

This means: **any time Claude edits a file in `server/db/schema/` or `herbst/db/schema/`, the ent code generation runs automatically in both modules.** No more "forgot to run ent generate" build failures.

**Other hook opportunities:**
- `PostToolUse` after any `server/*.go` edit → `make build-web` + health check
- `PostToolUse` after any `herbst/*.go` edit → `make build` + restart SSH
- `PreToolUse` before git push → validate test suite passes (`make test test-server`)

### 1.5 Permission System Allows Safe Autonomous Operation

The existing `.claude/settings.local.json` already configures:

```json
{
  "permissions": {
    "allow": [
      "Bash(go test:*)", "Bash(go build:*)", "Bash(make build:*)",
      "Bash(make start:*)", "Bash(bun *)", "Bash(gh issue *)"
    ],
    "defaultMode": "bypassPermissions"
  }
}
```

For production safety, we can tighten to `--permission-mode auto` (prompts for unknown) or use `--allowedTools` to whitelist only what's needed:

```bash
claude -p "Add faction entity schema" \
  --allowedTools "Bash(go *) Edit Read Write" \
  --permission-mode bypassPermissions
```

This prevents Claude from running `rm -rf`, `git push --force`, or other destructive operations while allowing full build/test access.

---

## 2. Specific Claude Features That Solve Known Pain Points

| Pain Point | Claude Feature | How It Solves |
|---|---|---|
| **Forgot `ent generate`** in either module | `PostToolUse` hooks | Auto-runs generation after any `db/schema/` file edit |
| **52 stale branches** | `--worktree` | Each task gets isolated worktree; done branches are cleaned up |
| **Must rebuild after schema/route changes** | `PostToolUse` hooks + Makefile | Hook runs `make build-web` after server edits, `make build` after herbst edits |
| **Must restart SSH server after merge** | `PostToolUse` hook or CLAUDE.md instruction | Claude reads CLAUDE.md → knows to run `make stop && make dev` |
| **pi-coding-agent 0% Go success** | Claude's native Go/ent model | Trained on Go, understands ent patterns, can compile and iterate |
| **One ticket at a time, no parallel PRs** | `--worktree` + serial dispatch | Enforced isolation; only one active worktree at a time |
| **Untracked files polluting repo** | `--worktree` | Worktrees don't touch main tree; untracked files stay isolated |
| **SKILL migration 4/9 tickets stuck** | Session resumption (`--resume`) | Long tasks can be paused and resumed without losing context |
| **Build fails silently** | Claude runs `go build` and reads errors | Iterative fix loop: edit → compile → read errors → fix |
| **Test failures not caught** | `PreToolUse` hook or prompt instruction | Claude must run `make test test-server` before completing |

### Deep Dive: The ent Generation Problem

The single most common failure mode in herbst-mud development is: editing a schema file in `server/db/schema/`, forgetting to run `ent generate` in **both** `server/` and `herbst/`, then getting mysterious compile errors.

Claude Code solves this three ways:

1. **Hooks** — automatic, guaranteed execution after schema edits
2. **CLAUDE.md** — documented instruction that Claude always reads at session start
3. **Iterative compilation** — even if hooks fail, Claude runs `go build`, sees the error, and runs `ent generate` to fix it

The iterative compilation loop is the safety net. A human might give up after a mysterious error. Claude reads the compiler output, identifies the missing generated code, runs the generator, and tries again. This is the core advantage of a tool-using agent over a simple code generator.

---

## 3. Concrete Workflow Proposal for herbst-mud

### 3.1 Single Ticket Workflow (Donatello Mode)

```bash
# Step 1: Leonardo dispatches a ticket via Hermes
# Hermes calls Claude Code in print mode with a worktree

claude -p "Implement GitHub issue #42: Add faction entity to ent schema. \
  See AGENT_KNOWLEDGE.md and CLAUDE.md for project context. \
  After implementation, run make test test-server to verify." \
  -w faction-entity-42 \
  --max-budget-usd 5.00 \
  --max-turns 50 \
  --permission-mode bypassPermissions \
  --output-format json
```

**What happens:**
1. Claude creates worktree `faction-entity-42`
2. Reads CLAUDE.md → learns project conventions
3. Reads `AGENT_KNOWLEDGE.md` → learns architecture
4. Edits `server/db/schema/faction.go`
5. **Hook fires** → `ent generate` runs in both `server/` and `herbst/`
6. Edits `server/routes/faction_routes.go` (new API endpoints)
7. Runs `go build` in both modules → fixes any compile errors iteratively
8. Runs `make test test-server` → fixes any test failures
9. Commits with 🟣 badge (Donatello)
10. Pushes to `origin/faction-entity-42`
11. Creates PR via `gh pr create`
12. Assigns to Raphael (QA) via `gh pr edit --reviewer`
13. Returns JSON result to Hermes with PR URL

**Step 2: Raphael reviews PR**
- If approved → Leonardo merges → Splinter rebuilds and restarts
- If changes requested → Claude resumes the session:

```bash
claude -p "Raphael requested changes on PR #43: add pagination to faction list endpoint" \
  --resume <session-id> \
  --max-budget-usd 2.00
```

### 3.2 Batch UPKEEP Workflow

For cleanup tasks like the 52 stale branches:

```bash
claude -p "Clean up all merged git branches. For each local branch that is merged into main, delete it. Do not delete main. Run: git branch --merged main | grep -v main | xargs git branch -d. Then prune remote: git remote prune origin." \
  --worktree cleanup-branches \
  --max-budget-usd 1.00
```

### 3.3 Agent Role Definitions

Claude Code's `--agents` flag defines custom sub-agents:

```bash
claude -p "Review PR #43 for code quality, security, and test coverage" \
  --agents '{
    "raphael": {
      "description": "QA Engineer - reviews for bugs, missing tests, security issues",
      "prompt": "You are Raphael (Raph), QA Engineer for herbst-mud. Your role: code review all pull requests, run integration tests, ensure no tests are broken. Use red circle emoji 🔴. Check: 1) Does the code compile? 2) Are there tests? 3) Any security issues? 4) Does it follow project style (max 100 lines per file, functional over OOP)?"
    },
    "donatello": {
      "description": "Lead Developer - implements features and fixes",
      "prompt": "You are Donatello (Donnie), Lead Developer for herbst-mud. Your role: implement features, write unit tests, ensure code is clean. Use purple circle emoji 🟣. Critical: after ANY schema or route change, run ent generate in BOTH herbst/ and server/. Then rebuild binaries with make build-all."
    },
    "splinter": {
      "description": "Architect - manages services and deployment",
      "prompt": "You are Splinter, Architect for herbst-mud. Your role: monitor and restart game services after code changes. Ensure database migrations run properly. Commands: make stop && make dev to restart all services."
    }
  }' \
  --agent raphael
```

### 3.4 MCP Integration for Database Queries

Claude Code supports MCP (Model Context Protocol) servers. For herbst-mud, we could add a PostgreSQL MCP server:

```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://herbst:herbst_password@localhost:5432/herbst_mud"]
    }
  }
}
```

This would allow Claude to:
- Query the database directly to verify migrations worked
- Check data integrity after schema changes
- Inspect table structure without reading Go code

**Security note:** For production, use read-only database credentials for the MCP server.

---

## 4. Estimated Cost Per Task Type

Based on Claude Sonnet 4 pricing (~$3/MTok input, ~$15/MTok output) and typical task complexity:

| Task Type | Typical Turns | Est. Input Tokens | Est. Output Tokens | Est. Cost |
|---|---|---|---|---|
| **Schema change** (add entity/field) | 15-25 | ~50K | ~20K | $0.30 - $0.60 |
| **Route addition** (new CRUD endpoints) | 20-35 | ~80K | ~30K | $0.50 - $0.90 |
| **UI component** (React admin panel) | 25-40 | ~100K | ~40K | $0.70 - $1.20 |
| **Bug fix** (compile error, test failure) | 5-15 | ~20K | ~10K | $0.10 - $0.30 |
| **Large feature** (e.g., quest system) | 40-60 | ~200K | ~80K | $1.50 - $3.00 |
| **Branch cleanup / UPKEEP** | 5-10 | ~15K | ~5K | $0.05 - $0.15 |

**Cost controls available:**
- `--max-budget-usd 5.00` — hard cap, Claude stops when reached
- `--max-turns 30` — limits tool-use iterations to prevent infinite loops
- `--model sonnet` vs `--model opus` — use cheaper model for simple tasks

**Realistic budget for completing SKILL migration (5 remaining tickets):**
- 5 schema changes × $0.50 = $2.50
- 3 route additions × $0.70 = $2.10
- 2 bug fixes × $0.20 = $0.40
- **Total: ~$5.00 for remaining SKILL migration**

**Realistic budget for UPKEEP tickets:**
- Branch cleanup × $0.10 = $0.10
- 3 modularization tasks × $1.00 = $3.00
- **Total: ~$3.10 for UPKEEP**

**Monthly operational cost estimate: $10-30/month** for typical development pace (5-10 tickets/week at $0.30-0.60 average).

---

## 5. Proposed CLAUDE.md for herbst-mud

```markdown
# CLAUDE.md — Claude Code Project Context for herbst-mud

## What This Project Is
herbst-mud is a Go/PostgreSQL MUD game engine designed to run multiple MUDs
and storylines. SSH client on port 4444, REST API on port 8080, admin panel
on port 3000.

## Critical Rules (NEVER SKIP)
1. **Dual ent generate**: After ANY change to `server/db/schema/` or
   `herbst/db/schema/`, run `ent generate` in BOTH directories:
   ```bash
   cd /home/sam/GitHub/herbst-mud/server && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   cd /home/sam/GitHub/herbst-mud/herbst && go run -mod=mod entgo.io/ent/cmd/ent generate ./db/schema
   ```
2. **Rebuild after changes**: After editing any `.go` file, rebuild:
   ```bash
   cd /home/sam/GitHub/herbst-mud && make build-all
   ```
3. **One ticket at a time**: Never work on multiple tickets in parallel.
   Use `--worktree` for each ticket to keep isolation.
4. **Always run tests** before pushing:
   ```bash
   make test && cd server && go test -v
   ```
5. **Restart services** after schema or route changes:
   ```bash
   cd /home/sam/GitHub/herbst-mud && make stop && make dev
   ```

## Architecture
- `herbst/` — SSH client (bubbletea TUI), separate Go module
- `server/` — REST API (Gin), separate Go module, ent ORM
- `admin/` — Vite/React/TanStack admin panel
- `admin-tui/` — Go bubbletea admin TUI
- `content/` — YAML data-driven content (skills, NPCs, items, rooms)

## Key Directories
- `server/db/schema/` — ent schema definitions (25 entities)
- `herbst/db/schema/` — ent schema definitions (7 entities, subset)
- `server/routes/` — Gin route handlers
- `herbst/cmd_*.go` — MUD command handlers
- `features/` — Gherkin BDD feature files

## Code Style
- Files MUST NOT exceed 100 lines. Break into new files.
- Functional over OOP. Keep code simple and modular.
- Use JSDoc-style comments, avoid inline comments.
- Sign commits with team badge emoji: 🟣 Donatello, 🔴 Raphael, 🐀 Splinter

## Service Management
- Start all: `make dev-all`
- Start backend only: `make dev`
- Stop all: `make stop`
- SSH logs: `tail -f /tmp/herbst-ssh.log`
- Web logs: `tail -f /tmp/herbst-web.log`
- Health check: `curl -s http://localhost:8080/healthz`

## Database
- PostgreSQL 15, ent ORM with auto-migration
- Dev connection: host=localhost port=5432 user=herbst password=herbst_password dbname=herbst_mud sslmode=disable
- Production: uses DATABASE_URL env var with sslmode=require

## Current State (2026-05)
- SKILL migration: 4/9 tickets done
- Code health: 91/100
- ~52 stale merged branches need cleanup
- character_routes.go is 1,847 lines (technical debt, accepted)

## Common Failure Modes (WATCH FOR THESE)
1. Forgot ent generate → compile errors about missing generated code
2. Edited schema in server/ but not herbst/ → client can't deserialize
3. Didn't rebuild after Go changes → running stale binary
4. Didn't restart services → API returns old responses
5. Untracked files in main working tree → git pollution
```

---

## 6. Comparison: Claude Code vs pi-coding-agent vs Hermes delegate_task

| Criterion | Claude Code CLI | pi-coding-agent | Hermes delegate_task |
|---|---|---|---|
| **Go/ent understanding** | Excellent (model trained on Go) | 0% success | Model-dependent |
| **Tool use (file edit, bash)** | Native, multi-turn | Limited | Full (via subagent) |
| **Iterative compilation** | Yes (edit → build → fix loop) | No | Yes (via subagent) |
| **Worktree isolation** | Built-in `-w` flag | No | No (uses cwd) |
| **Hooks (auto ent generate)** | Yes (PostToolUse) | No | No |
| **Session resumption** | `--resume`, `--continue` | No | No |
| **Cost control** | `--max-budget-usd`, `--max-turns` | Unknown | Token-based |
| **CLAUDE.md context injection** | Yes (auto-discovered) | No | Via system prompt |
| **Custom agent definitions** | `--agents` JSON | No | Via persona |
| **MCP integration** | Yes (Postgres, etc.) | No | No |
| **Permission system** | Granular allow/deny per tool | Unknown | No |
| **Branch cleanup safety** | Worktrees + permission allowlist | No | Manual |
| **Structured output** | JSON, stream-json | No | Text only |
| **GitHub CLI integration** | `gh` via Bash tool | No | Via subagent |

### The Real Advantage

**Claude Code's key differentiator is the edit-compile-fix loop.** When pi-coding-agent encountered a Go compile error, it had no mechanism to read the error, edit the code, and try again. It was a single-shot generator. Claude Code is a multi-turn agent that:

1. Edits the file
2. Runs `go build`
3. Reads the compiler output
4. Edits the file again to fix the error
5. Repeats until green

This loop is **the** essential capability for Go development, where the compiler is the primary guardrail. Without it, you get pi-coding-agent's 0% success rate.

Hermes `delegate_task` is meta-orchestration — it dispatches work to subagents. Claude Code is the subagent that actually does the work. They're complementary: Hermes decides *what* to do, Claude Code decides *how* to do it.

---

## 7. Integration Architecture

```
┌─────────────────────────────────────────────┐
│              Hermes Agent (Orchestrator)      │
│  - Reads GitHub Project board                 │
│  - Dispatches tickets via delegate_task       │
│  - Tracks progress, collects results          │
└────────────────┬────────────────────────────┘
                 │
                 │ delegate_task → shell command
                 ▼
┌─────────────────────────────────────────────┐
│          Claude Code CLI (Subagent)          │
│  - Receives ticket description               │
│  - Creates worktree (-w ticket-name)         │
│  - Reads CLAUDE.md for project context        │
│  - Implements, builds, tests iteratively     │
│  - Creates PR via gh CLI                     │
│  - Returns JSON result to Hermes             │
└─────────────────────────────────────────────┘
                 │
                 │ hooks fire automatically
                 ▼
┌─────────────────────────────────────────────┐
│         PostToolUse / PreToolUse Hooks       │
│  - ent generate after schema edits           │
│  - make build after Go edits                 │
│  - make test before git push                 │
└─────────────────────────────────────────────┘
```

### Hermes → Claude Code bridge command:

```bash
#!/bin/bash
# hermes-claude-bridge.sh
# Usage: ./hermes-claude-bridge.sh "ticket description" worktree-name budget

TICKET="$1"
WORKTREE="$2"
BUDGET="${3:-5.00}"

cd /home/sam/GitHub/herbst-mud

claude -p "$TICKET. See CLAUDE.md and AGENT_KNOWLEDGE.md for project context. After implementation: 1) Run ent generate in both server/ and herbst/, 2) Build with make build-all, 3) Test with make test && cd server && go test -v, 4) Create a PR with gh pr create." \
  -w "$WORKTREE" \
  --max-budget-usd "$BUDGET" \
  --max-turns 50 \
  --permission-mode bypassPermissions \
  --output-format json \
  --allowedTools "Bash(go *) Bash(make *) Bash(cd *) Bash(gh *) Bash(git *) Bash(curl *) Edit Write Read"
```

---

## 8. Risks and Mitigations

| Risk | Mitigation |
|---|---|
| Claude generates incorrect ent schema | Hook auto-runs `ent generate` + `go build`; compile error triggers fix loop |
| Cost overrun on complex tasks | `--max-budget-usd` hard cap; estimate $0.50-3.00 per task |
| Claude pushes to main directly | `--allowedTools` excludes `Bash(git push origin main)`; use PR workflow |
| Worktree left behind after crash | Cleanup script: `find .git-worktrees -mtime +1 -exec rm -rf {} \;` |
| Claude modifies wrong files | CLAUDE.md lists exact directories; `--add-dir` restricts file access |
| Hallucinated API endpoints | Hook runs `go build` after edits; Claude reads compiler errors |
| Database corruption via MCP | Use read-only DB credentials; restrict to SELECT queries |

---

## 9. Recommendation

**Adopt Claude Code CLI as the primary autonomous coding agent for herbst-mud backend tasks.**

The evidence is clear:
- pi-coding-agent cannot handle Go/ent tasks (0% success)
- Claude Code has native Go understanding, iterative compilation, and hooks for the exact workflows that break manually
- Worktrees enforce the "one ticket at a time" rule automatically
- Cost is predictable and low ($0.10-3.00 per task, ~$10-30/month)
- CLAUDE.md injects all project context in every session
- The tool is already installed and partially configured

**Next steps:**
1. Create `CLAUDE.md` at repo root (content from Section 5 above)
2. Configure `PostToolUse` hooks in `.claude/settings.json` for ent generation
3. Create `hermes-claude-bridge.sh` script for Hermes → Claude dispatch
4. Test with one SKILL migration ticket end-to-end
5. Add `.mcp.json` for PostgreSQL MCP server (optional, Phase 2)
6. Define custom agents (`--agents`) for Donatello and Raphael roles

---

*End of PRO analysis. This document should be stored as a skill for future reference.*