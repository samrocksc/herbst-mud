# RFC-009: Speak/Chat/Message System

**Status:** MVP Implemented
**Author:** Leonardo (with Sam)
**Created:** 2026-05-09
**Implemented:** 2026-05-12
**Builds on:** RFC-002 (Effects & Hooks), RFC-004/005/006 (NPC Dialog), LOGS (applogs pipeline)
**Related:** EFF-010 (room_message effect type)

---

## 1. Executive Summary

This RFC proposes a tiered speech and messaging system for herbst-mud — the
communication backbone that makes a multiplayer world feel alive. It covers
room-based speech (`say`, `yell`, `emote`), private messaging (`tell`,
`reply`), and channel-based chat with player-configurable toggles and colors.

The system is designed as a **first-class game service** built on the existing
`applogs` pipeline and effects infrastructure, with future extensibility for
WebSocket clients and group/guild channels.

---

## 2. Research Summary

Three sub-agents researched classic MUDs (DikuMUD, CircleMUD, ROM, GodWars,
Aardwolf, BatMUD, Discworld), modern games (WoW, FFXIV, ESO, Discord, Twitch),
and distributed chat architectures (Redis pub/sub, event sourcing,
connection management, offline delivery, scaling).

**Key takeaways for herbst-mud design:**

| Insight | Source | Implication |
|---------|--------|-------------|
| Room→Zone→World tiered range is universal | Classic MUDs | `say`/`yell`/`shout` is the right model |
| Per-channel on/off toggles are #1 player demand | Aardwolf community | Every channel must be toggleable |
| Color coding per channel reduces cognitive load | All MUDs | Distinct colors for say/tell/yell/channel |
| `reply` to last tell is essential | All MUDs | Must implement reply command |
| `/me` and emotes need pronoun substitution | Discworld, GodWars | Avoid third-person errors in emotes |
| Rate limiting needs layers | Aardwolf, BatMUD | Per-channel cooldowns + global rate limit |
| WebSocket + Redis pub/sub is the scaling sweet spot | Architecture research | Design for WebSocket transport from day one |
| Separate transport from game logic | Architecture research | Chat service is a standalone layer |
| Offline tells are a beloved feature | Aardwolf | Queue tells for offline players |
| Channel history/scrollback reduces FOMO | Aardwolf, Discord | Per-channel ring buffer |

---

## 3. Speech Tiers (Range-Based)

### 3.1 `say` — Room Only

```
> say Hello there, traveler.

Gorthak says, "Hello there, traveler."
```

Everyone in the same room sees it. The speaker's name is included.

### 3.2 `yell` — Zone/Area

```
> yell Guards! Intruders!

Gorthak yells, "Guards! Intruders!"
```

Visible to everyone in the current zone (connected rooms within N hops, or all
rooms sharing the same `area_id`). Cooldown: 5 seconds.

### 3.3 `shout` — World

```
> shout Selling Ancient Key for 500 gold!

Gorthak shouts, "Selling Ancient Key for 500 gold!"
```

Visible to all connected players in the entire MUD. Cooldown: 30 seconds.
Cost: 5 stamina (configurable via game config).

### 3.4 Range Table

| Command | Range | Cooldown | Cost | Color |
|---------|-------|----------|------|-------|
| `say` | Room | None | None | White/neutral |
| `yell` | Zone (N rooms / area) | 5s | None | Yellow |
| `shout` | World | 30s | 5 stamina | Bold red |

---

## 4. Private Messaging

### 4.1 `tell` — Direct 1:1 Message

```
> tell Gorthak Meet me at the fountain.

You tell Gorthak, "Meet me at the fountain."
```

```
Elder Myrddin tells you, "Meet me at the fountain."
```

- Works across the entire MUD regardless of room
- If target is offline: message is queued and delivered on next login
- Offline tells persist for 7 days, then expire
- Target can be a player or NPC

### 4.2 `reply` — Respond to Last Tell

```
> reply I'm on my way.

You reply to Elder Myrddin, "I'm on my way."
```

Replies to the last person who sent you a tell. Stores last teller per
character in the runtime conversation state.

### 4.3 `whisper` — Room-Private Message

```
> whisper Gorthak Don't trust the goblin.

You whisper to Gorthak, "Don't trust the goblin."
```

Only the sender and the target see it. Others in the room see at most:
"Gorthak whispers something to Sam."

### 4.4 Block/Ignore

```
> ignore Gorthak
You are now ignoring Gorthak.

> unignore Gorthak
You are no longer ignoring Gorthak.
```

Ignored players' tells, whispers, says, and channel messages are not visible
to the blocker. Stored as a `character_ignore` join table.

---

## 5. Emotes

### 5.1 `emote` / `/me` — Custom Action

```
> emote bows deeply before the elder.

Gorthak bows deeply before the elder.
```

Freeform text displayed as an action. The character's name is prepended
automatically.

### 5.2 Social Commands (Pre-Built)

Standard socials from classic MUD heritage:

```
> smile
You smile happily.

> smile Gorthak
You smile at Gorthak.
Gorthak sees: Sam smiles at you.

> laugh
You burst out laughing.

> bow
You bow gracefully.

> wave
You wave.

> hug Gorthak
You hug Gorthak warmly.
Gorthak sees: Sam hugs you warmly.

> grin
You grin evilly.
```

Each social has three variants:
- `self_text`: what the actor sees when no target
- `room_text`: what others in the room see when no target
- `target_self_text`: what the actor sees when targeting someone
- `target_text`: what the target sees
- `target_room_text`: what others in the room see when targeting

Socials are stored in the DB as `social_commands` and are editable via
admin UI. MVP ships with 15-20 curated socials.

### 5.3 Pronoun Substitution

Emotes and socials support pronoun substitution based on the character's
`gender` field (he/she/they):

```
> emote checks {his} equipment one last time.
Gorthak checks his equipment one last time.      (if gender = he_him)
Gorthak checks her equipment one last time.      (if gender = she_her)
Gorthak checks their equipment one last time.    (if gender = they_them)
```

Pronouns available: `{he}`, `{him}`, `{his}`, `{himself}`, and feminine/neutral
equivalents.

---

## 6. Channels (Topic-Based)

### 6.1 Built-In Channels

| Channel | Description | Default | Color |
|---------|-------------|---------|-------|
| `chat` | General off-topic chat | ON | Magenta |
| `newbie` | New player help/questions | ON | Green |
| `trade` | Buying/selling items | OFF | Yellow |
| `ooc` | Out-of-character discussion | ON | Cyan |
| `admin` | Admin announcements (admin-only send, all read) | ON | Gold |

### 6.2 Channel Commands

```
> chat Anyone seen the goblin shaman?
[chat] Gorthak: "Anyone seen the goblin shaman?"

> newbie How do I equip a sword?
[newbie] Gorthak: "How do I equip a sword?"

> channels              — list all channels and their status
> channel chat on       — enable channel
> channel trade off     — disable channel
> channel color chat cyan  — set custom color for a channel
```

### 6.3 Channel Toggles

Stored per-character in the DB:

```go
// server/db/schema/character_channel.go
CharacterChannel {
    character_id    int     (FK → Character)
    channel         string   // "chat", "newbie", "trade", etc.
    enabled         bool     (default: true for chat/newbie, false for trade)
    color           string   // ANSI color name (default from channel definition)
}
```

`channels` command lists all with status indicators.

---

## 7. Transport Architecture

### 7.1 Current State (SSH-Only)

All players connect via SSH to the `herbst/` game server. Messages are delivered
through the existing SSH session's TTY output. No cross-process communication
needed — the game server handles all routing.

**For MVP (SSH-only):**

```
Player types "say hello"
  → cmd_say.go parses input
  → ChatService.SendSay(actorID, roomID, message)
    → Query all characters in roomID
    → Write formatted message to each character's SSH session
```

### 7.2 Future (WebSocket + Web Client)

When the web client arrives (RFC-001), add a chat gateway:

```
Browser (WebSocket) → Chat Gateway :8081 → Redis Pub/Sub → herbst SSH server
```

The chat gateway is a thin bridge: WebSocket ↔ Redis pub/sub. The herbst SSH
server subscribes to relevant channels and fans out to SSH sessions. This
means SSH players and web players share the same chat space.

**Recommended stack for Phase 2:**

```
herbst/ (SSH game server)
  ├── chat/service.go       — ChatService (routes messages)
  ├── chat/transport.go     — Transport interface (SSH / WS / Redis)
  ├── chat/channel.go       — Channel subscription management
  └── cmd_say.go, cmd_tell.go, etc. — command handlers

server/ (REST API)
  └── routes/chat_routes.go — REST endpoints for message history/queries

chat-gateway/ (optional, future)
  └── WebSocket ↔ Redis bridge
```

### 7.3 Transport Interface

```go
// herbst/chat/transport.go
type Transport interface {
    Send(characterID int, message string) error
    Broadcast(roomID int, message string) error
    BroadcastZone(zoneID int, message string) error
    BroadcastWorld(message string) error
    SendChannel(channel string, message string) error
}
```

SSH implementation writes to the `tea.Program` output. Future WebSocket
implementation writes to a connection registry.

---

## 8. Commands Summary

| Command | Range | Cooldown | Cost | Storage |
|---------|-------|----------|------|---------|
| `say <msg>` | Room | None | None | Logged |
| `yell <msg>` | Zone | 5s | None | Logged |
| `shout <msg>` | World | 30s | 5 stamina | Logged |
| `tell <player> <msg>` | Direct | 2s | None | Logged + offline queue |
| `reply <msg>` | Direct (last teller) | 2s | None | Logged |
| `whisper <player> <msg>` | Room-private | None | None | Logged |
| `emote <action>` | Room | None | None | Logged |
| `/<social> [target]` | Room | None | None | Logged |
| `chat <msg>` | Channel | 2s | None | Logged |
| `newbie <msg>` | Channel | 2s | None | Logged |
| `trade <msg>` | Channel | 10s | None | Logged |
| `channels` | — | — | — | — |
| `channel <name> on/off` | — | — | — | Persisted |
| `channel color <name> <color>` | — | — | — | Persisted |
| `ignore <player>` | — | — | — | Persisted |
| `unignore <player>` | — | — | — | Persisted |

---

## 9. Rate Limiting & Spam Prevention

Three layers:

1. **Per-command cooldowns** — `shout` 30s, `yell` 5s, `tell` 2s, `chat` 2s, `trade` 10s
2. **Global rate limit** — 30 messages per minute per character (configurable in game config)
3. **Repeat detection** — identical message within 10s is suppressed

Violations: warn → 60s mute → 5min mute → admin notification.

---

## 10. Message Logging

All messages are logged to `applogs` with:

```json
{
  "level": "INFO",
  "service": "chat",
  "character_id": 7,
  "message": "[say] Gorthak: Hello there, traveler.",
  "metadata": {
    "command": "say",
    "room_id": 42,
    "target_id": null,
    "channel": null
  }
}
```

This enables:
- Admin log viewer filtering by `service=chat`
- Searching for specific player messages
- Debugging communication issues
- Moderation audit trail

---

## 11. DB Schemas

### 11.1 Offline Tell Queue

```go
// server/db/schema/tell_queue.go
TellQueue {
    id            int (auto)
    sender_id     int
    sender_name   string
    recipient_id  int
    message       string
    sent_at       timestamp
    delivered_at  timestamp (nullable)
    expires_at    timestamp (sent_at + 7 days)
}
```

### 11.2 Character Ignores

```go
// server/db/schema/character_ignore.go
CharacterIgnore {
    ignorer_id    int  (FK → Character)
    ignored_id    int  (FK → Character)
    created_at    timestamp
}
```

### 11.3 Character Channels

```go
// server/db/schema/character_channel.go
CharacterChannel {
    character_id  int     (FK → Character)
    channel       string   // "chat", "newbie", etc.
    enabled       bool
    color         string   // ANSI color name
}
```

### 11.4 Social Commands

```go
// server/db/schema/social_command.go
SocialCommand {
    name              string    // "smile", "bow", "laugh"
    self_text         string    // "You smile happily."
    room_text         string    // "Gorthak smiles happily."
    target_self_text  string    // "You smile at {target}."
    target_text       string    // "{actor} smiles at you."
    target_room_text  string    // "{actor} smiles at {target}."
}
```

---

## 12. NPC Speech Integration

NPCs use the same `say` and `emote` infrastructure as players. Hooks (RFC-002)
can trigger NPC speech via the `room_message` effect type (EFF-010).

```
Hook (on NPC template "Town Crier"):
  event: "on_timer" (future, or on_enter_room)
  effect: room_message
  parameters: {"text": "The town crier bellows: 'Hear ye! The king seeks adventurers!'"}
```

NPC say color is distinct from player say (green vs white) so players can
visually distinguish NPC dialogue.

---

## 13. Admin UI

- **Chat log viewer**: filter `/logs` by `service=chat`, see all messages
- **Channel management**: enable/disable channels, set defaults, manage colors
- **Social command editor**: CRUD on socials with preview/test
- **Player moderation**: view ignore lists, force-unmute, review reports

---

## 14. Acceptance Criteria

### Core (MVP for SSH client)

- [ ] `say`, `yell`, `shout` commands working with correct ranges
- [ ] `tell`, `reply`, `whisper` private messaging
- [ ] `emote` and basic socials (smile, laugh, bow, wave, nod, grin, sigh)
- [ ] Pronoun substitution (`{he}`, `{him}`, `{his}`)
- [ ] `chat`, `newbie` channels with on/off toggle
- [ ] `trade` channel (default off, 10s cooldown)
- [ ] `ignore`/`unignore` player blocking
- [ ] Offline tell queueing with 7-day expiry
- [ ] Per-command cooldowns enforced
- [ ] Global rate limit (30 msg/min)
- [ ] All messages logged to applogs with `service=chat`
- [ ] NPC speech distinct color from player speech

### Backend
- [ ] `character_channel.go`, `character_ignore.go`, `social_command.go`, `tell_queue.go` ent schemas
- [ ] `ChatService` in `herbst/chat/` with `Transport` interface
- [ ] All command handlers in `herbst/` (cmd_say, cmd_tell, cmd_emote, etc.)
- [ ] REST endpoints for channel toggles, ignore management, tell queue

### Admin UI
- [ ] Chat filter in `/logs` page (service=chat)
- [ ] Social command editor

### Future (Phase 2, deferred)
- [ ] WebSocket transport for web client chat
- [ ] Redis pub/sub for multi-server scaling
- [ ] `trade`, `ooc`, `admin` channels
- [ ] Channel history/scrollback (ring buffer)
- [ ] Message reporting / moderation tools
- [ ] Chat tabs in web client
