# RFC-001: Web-Based Client Architecture for Herbst-MUD

**Status:** Draft
**Author:** Leonardo (with Sam)
**Created:** 2026-05-08
**Target:** herbst-mud SSH game server + herb-server HTTP API

---

## 1. Executive Summary

This RFC proposes adding a web-based client for herbst-mud as a third frontend
option alongside the existing SSH TUI (herbst/) and the React admin UI
(admin/). The web client connects via WebSocket through a thin gateway that
proxies to the existing SSH server, requiring zero changes to the core game
logic.

**Key constraint:** The MUD's game state, combat engine, and command
processing all live inside `herbst/main.go` wired to a `tea.Model`. This
architecture is tightly coupled to the SSH session. Any web client must either
pipe through that same SSH session or share the underlying game engine through
a redesigned service layer.

---

## 2. Current Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
│                                                                  │
│   SSH Client (any terminal)           Browser (future)           │
│         │                                    │                    │
│         └─────────── port 4444 ─────────────┘                   │
│                         │                                        │
└─────────────────────────┼────────────────────────────────────────┘
                          │
              ┌───────────▼───────────┐
              │    herbst/main.go     │  ← charmbracelet/wish SSH
              │   (tea.Model-based)   │    server on :4444
              │                       │
              │  • login/auth         │
              │  • command parsing     │
              │  • combat engine      │
              │  • room/npc/item mgmt  │
              │  • full game state    │
              └───────────────────────┘
                          │
                          │ db.Open() via lib/pq
                          │
              ┌───────────▼───────────┐
              │      PostgreSQL        │
              │    herbst_mud DB       │
              └───────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         ADMIN UI                                 │
│                                                                  │
│   Browser → port 5173 (Vite dev)  or  port 80 (prod)             │
│         │                                                       │
│         └─── GET/POST /api/* ──→ herb-server (Go/Gin) on :8080   │
│                                       │                           │
│                                       └─── db.* ──→ PostgreSQL    │
└──────────────────────────────────────────────────────────────────┘
```

### Relevant existing components

| Binary | Port | Role |
|--------|------|------|
| `herbst` | `:4444` | SSH game server (TUI via bubbletea) |
| `herbst-web` (server/) | `:8080` | REST API for admin UI |
| `admin/` (npm) | `:5173` | React admin panel (dev) |
| `herbst-admin-tui/` | — | Go TUI for admin operations |

---

## 3. Design Goals

1. **Zero changes to the game engine** — `herbst/main.go` and its `tea.Model`
   remain untouched. No risk of regressions in the existing SSH game client.
2. **Auth reuse** — Character creation, login, and session management use the
   same DB-backed logic already in `herbst/`.
3. **Full fidelity** — Web client must support every command and game feature
   the SSH client supports, including combat, inventory, examine, etc.
4. **Graceful degradation** — If the web gateway is unavailable, the SSH
   server keeps working.
5. **Minimal operational overhead** — Single optional gateway binary, or one
   additional endpoint on an existing port.

---

## 4. Option A: WebSocket-SSH Gateway (Recommended for Phase 1)

```
Browser  ──── WebSocket ────  Gateway  ──── SSH ────  herbst :4444
           :8081 or :8080        │                   (unchanged)
           (gateway binary)      │
                                 └─── reads DB directly for
                                     initial character list
```

### How it works

The gateway is a **separate process** that acts as an SSH client toward the
existing herbst server. The browser connects to the gateway via WebSocket.
The gateway authenticates on behalf of the browser user by speaking the same
SSH protocol the terminal client uses.

```
[Browser]  ws://localhost:8081/ssh  ←→  [Gateway]  tcp:localhost:4444  ←→  [herbst]
```

The gateway:
1. Upgrades HTTP → WebSocket
2. Prompts for character selection (or accepts credentials via WS messages)
3. Opens an SSH session to `localhost:4444`
4. Pipes WebSocket bytes ↔ SSH channel bytes verbatim

The herbst server sees a normal SSH client connection and processes it exactly
as it does today. The `tea.Model` inside herbst is unaffected.

### What the gateway does NOT do

- It does **not** read or write the database directly for game state
- It does **not** implement any game logic
- It is a **dumb pipe** at the byte level

### Gateway API (WebSocket messages)

```websocket
Client → Gateway:
  auth:        {"type":"auth","username":"...","password":"..."}
  char_select: {"type":"char_select","character_id":123}
  input:       {"type":"input","data":"look\r"}

Gateway → Client:
  output:      {"type":"output","data":"You are in a dark room.\r\n"}
  auth_ok:     {"type":"auth_ok","characters":[...]}
  char_ok:     {"type":"char_ok"}
  error:       {"type":"error","message":"..."}
```

### Advantages

- herbst server unchanged — fully isolated
- All game logic, combat, auth, command parsing stays in the proven SSH code
- Simple to implement and test
- Gateway crash does not affect SSH players

### Disadvantages

- **Latency:** Double hop (browser→gateway→herbst) adds ~1–5ms
- **No server-push:** The gateway must poll or use SSH channel warnings to
  detect output — but SSH is already full-duplex so this works fine
- **Terminal emulation:** The gateway sends raw ANSI bytes; the browser must
  render them. Can use xterm.js or a lightweight parser.

### Reference implementations

| Project | Approach |
|---------|----------|
| [yudai/gotty](https://github.com/yudai/gotty) | WebSocket → PTY → SSH/Shell. Single binary. |
| [charmbracelet/chtmp](https://github.com/charmbracelet/chtmp) | Charm's own WebSocket→SSH bridge (if released) |
| [dualed/termd](https://github.com/dualed/termd) | WebSocket terminal emulator with SSH support |

---

## 5. Option B: Native WebSocket Transport in herbst (Phase 2)

Merge WebSocket support directly into `herbst/main.go` so both SSH and
WebSocket clients share the same `tea.Model` instance.

```
Browser  ──── WebSocket ────  herbst :4444  (or new port)
                              ├──────────── SSH server (existing)
                              └──────────── WebSocket server (new)
```

### How it would work

`charmbracelet/ssh` does not natively support WebSocket upgrade — but the
underlying `gossh` library allows custom `net.Listener` implementations. A
custom listener could accept both SSH and WebSocket connections on the same
port (via TLS ALPN or connection upgrade detection).

Alternatively, run the WebSocket server on a **separate port** (e.g., `:4445`)
that shares the same `tea.Model` instances via an in-process bus:

```
herbst/
  ssh_server.go    — existing wish.Server on :4444
  ws_server.go     — gorilla/websocket on :4445
  game_bus.go      — in-memory event bus; both servers dispatch to shared model
```

### Advantages

- Single binary, no separate gateway process
- Lower latency (single hop)
- Shared game state without inter-process communication

### Disadvantages

- **Significant refactor** of herbst/main.go — extract game state into a
  shared service layer both SSH and WS handlers can use
- The `tea.Model` in herbst is currently **coupled to a single SSH session** —
  supporting multiple concurrent sessions requires making the game state
  global (or using a session registry), which is a large architectural change
- Risk of regressions in the existing SSH game client
- Much larger implementation effort

---

## 6. Option C: Shared Game Engine (Long-term)

Decouple the game engine from the presentation layer entirely.

```
┌─────────────────────────────────────────────────────────┐
│                   game-engine/                          │
│   (pure Go library: combat, rooms, items, characters)  │
│                                                          │
│   Shares no network code. Exposes service interfaces.   │
└──────────────┬──────────────────────────┬──────────────┘
               │                          │
       ┌───────▼───────┐          ┌────────▼────────┐
       │  herbst/main  │          │  herb-server     │
       │  SSH :4444   │          │  WebSocket :8081 │
       │  (wish+tea)  │          │  (browser client) │
       └───────────────┘          └───────────────────┘
```

This is the cleanest architecture long-term but requires the most work
upfront. It would also enable the admin UI to show live game state without
polling.

---

## 7. Option A Implementation Sketch

### File layout

```
herbst-gateway/
  main.go          — gateway binary
  websocket.go     — WebSocket server (gorilla/websocket)
  ssh_client.go    — SSH client toward herbst :4444
  terminal.go      — ANSI/xterm byte buffer handling
  Makefile         — build target
```

### Dependencies

```go
require (
    github.com/gorilla/websocket v1.5.1
    github.com/charmbracelet/ssh v0.0.0-20250128164007-98fd5ae11894
    golang.org/x/crypto v0.36.0   // for ssh client
)
```

### Minimal main.go

```go
func main() {
    // WebSocket server on :8081
    http.HandleFunc("/ssh", handleSSH)
    go http.ListenAndServe(":8081", nil)

    // SSH client pool (one per active WebSocket connection)
    // ... see full sketch in RFC
}
```

### Frontend (browser)

The browser-side component is a **separate project** or part of `admin/`.
It needs:

1. `WebSocket` connection to `ws://host:8081/ssh`
2. An ANSI/xterm parser (e.g., `xterm.js`) to render output
3. Keyboard input capture → JSON → WebSocket send
4. A simple HTML shell that loads `xterm.js` + a small wrapper script

The frontend is intentionally minimal — the gateway handles all protocol
translation.

---

## 8. Security Considerations

1. **SSH gateway auth:** The gateway authenticates to herbst with credentials
   from the browser. These credentials transit over TLS (if the gateway port
   is behind HTTPS) and then over localhost SSH. No additional exposure.
2. **WebSocket origin:** The gateway should validate the `Origin` header and/or
   require a per-session token from the admin API before upgrading.
3. **Rate limiting:** The gateway should implement per-connection write
   throttling to prevent abusive clients from flooding the SSH session.
4. **No direct DB access from gateway:** It only speaks SSH to herbst, so a
   compromised gateway cannot arbitrarily modify game state.

---

## 9. Recommended Roadmap

### Phase 1 (Option A — minimal viable)
- [ ] Build `herbst-gateway/` as separate Go module
- [ ] WebSocket endpoint that dials herbst :4444 as SSH client
- [ ] JSON message protocol (auth, char_select, input, output)
- [ ] Simple HTML test page with `xterm.js`
- [ ] Deploy alongside existing herbst binary; verify SSH players unaffected

### Phase 2 (browser UX)
- [ ] Character selection UI in browser
- [ ] Proper ANSI color rendering via `xterm.js`
- [ ] Session persistence (reconnect without losing game state)
- [ ] HTTPS termination (TLS for the WebSocket connection)

### Phase 3 (native integration, if desired)
- [ ] Evaluate Option B: WebSocket listener merged into herbst
- [ ] Consider Option C: shared game engine for long-term maintainability

---

## 10. Open Questions

1. Should the gateway live in this repo (`herbst-gateway/`) or separately?
2. Do you want character selection to happen in the browser, or require the
   player to first `ssh herbstmud.local` to create/select a character, then
   connect the web client to an **existing session**?
3. Do you want guest/anonymous play (no auth required), or always account-based?
4. What is the target deployment environment? (Same host? Separate VM? Cloud?)
   This affects whether Option A port exposure is acceptable.

---

## 11. TL;DR

**Option A (WebSocket-SSH gateway)** is the fastest path to a web client with
the lowest risk to the existing SSH game server. The gateway is a ~200-line
separate Go binary that speaks WebSocket to browsers and SSH to the existing
herbst server. The herbst server requires zero changes. The browser needs only
`xterm.js` and a small wrapper to render the ANSI output.
