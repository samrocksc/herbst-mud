# Ticket: RFC-009 - SSH Messaging Commands (say, yell, shout, tell, whisper, emote)

**Status:** Draft  
**Created:** 2026-05-12  
**RFC Reference:** RFC-009 (Speak/Chat/Message System)

---

## Executive Summary

Implement player-facing messaging commands for SSH client: `say`, `yell`, `shout`, `tell`, `reply`, `whisper`, `emote`, and channel chat (`chat`, `newbie`, `trade`, `ooc`, `admin`). This is the user-facing interface for communication.

---

## Commands to Implement

### Room-based Speech

| Command | Range | Cooldown | Cost | Example |
|---------|-------|----------|------|---------|
| `say <msg>` | Room | None | None | `say Hello there` |
| `yell <msg>` | Zone | 5s | None | `yell Guards!` |
| `shout <msg>` | World | 30s | 5 stamina | `shout Auction!` |

### Private Messaging

| Command | Range | Cooldown | Cost | Example |
|---------|-------|----------|------|---------|
| `tell <player> <msg>` | Direct | 2s | None | `tell Gorthak Meet me` |
| `reply <msg>` | Last teller | 2s | None | `reply I'm on my way` |
| `whisper <player> <msg>` | Room-private | None | None | `whisper Gorthak Don't trust` |

### Emotes

| Command | Range | Cooldown | Cost | Example |
|---------|-------|----------|------|---------|
| `emote <action>` | Room | None | None | `emote bows deeply` |
| `/<social> [target]` | Room | None | None | `/smile`, `/bow Gorthak` |

### Channel Chat

| Command | Channel | Cooldown | Example |
|---------|---------|----------|---------|
| `chat <msg>` | General | 2s | `chat Anyone here?` |
| `newbie <msg>` | Newbie help | 2s | `newbie How do I equip?` |
| `trade <msg>` | Trading | 10s | `trade Selling sword` |
| `ooc <msg>` | OOC | 2s | `ooc Meeting at 9pm` |
| `admin <msg>` | Admin only | 2s | `admin Debug log please` |

### Channel Management

| Command | Description |
|---------|-------------|
| `channels` | List all channels and status |
| `channel <name> on/off` | Enable/disable channel |
| `channel color <name> <color>` | Set custom color |
| `ignore <player>` | Block messages from player |
| `unignore <player>` | Unblock player |

---

## Implementation Plan

### Phase 1: Backend (herbst-server/)

#### 1. Database Schema

Create new tables:

```go
// server/db/schema/character_channel.go
CharacterChannel {
    character_id  int     (FK → Character)
    channel       string   // "chat", "newbie", "trade", etc.
    enabled       bool     // default true for chat/newbie
    color         string   // ANSI color name
}

// server/db/schema/character_ignore.go
CharacterIgnore {
    ignorer_id    int  (FK → Character)
    ignored_id    int  (FK → Character)
    created_at    timestamp
}

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

#### 2. Services

**ChatService** (`server/service/chat_service.go`):

```go
type ChatService interface {
    SendSay(charID, roomID int, message string) (*MessageResult, error)
    SendYell(charID, roomID int, message string) (*MessageResult, error)
    SendShout(charID int, message string) (*MessageResult, error)
    SendTell(fromID, toID int, message string) (*MessageResult, error)
    SendWhisper(fromID, toID int, message string) (*MessageResult, error)
    SendEmote(charID int, action string) (*MessageResult, error)
    SendChannel(channel, message string, charID int) (*MessageResult, error)
    
    // Channel management
    GetChannels(charID int) ([]ChannelState, error)
    SetChannelEnabled(charID int, channel string, enabled bool) error
    SetChannelColor(charID int, channel string, color string) error
    
    // Ignore system
    IgnorePlayer(charID, ignoredID int) error
    UnignorePlayer(charID, ignoredID int) error
    GetIgnoredPlayers(charID int) ([]int, error)
    
    // Offline tells
    QueueOfflineTell(fromID int, recipientName string, message string) error
    DeliverQueuedTells(charID int) ([]QueuedTell, error)
}
```

**SocialCommandService** (`server/service/social_command_service.go`):

```go
type SocialCommandService interface {
    GetSocials() ([]*SocialCommand, error)
    GetSocial(name string) (*SocialCommand, error)
    CreateSocial(input CreateSocialInput) (*SocialCommand, error)
    UpdateSocial(name string, input UpdateSocialInput) (*SocialCommand, error)
    DeleteSocial(name string) error
}
```

#### 3. Routes

**REST Routes** (`server/routes/chat_routes.go`):

```go
POST   /api/chat/say       // Send room message
POST   /api/chat/yell      // Send zone message
POST   /api/chat/shout     // Send world message
POST   /api/chat/tell      // Send direct message
POST   /api/chat/whisper   // Send whisper
POST   /api/chat/emote     // Send emote
POST   /api/chat/channel   // Send channel message

GET    /api/channels       // List all channels
POST   /api/channels/:name/enabled  // Toggle channel
POST   /api/channels/:name/color    // Set channel color

POST   /api/ignore         // Ignore player
DELETE /api/ignore/:id     // Unignore player
GET    /api/ignored        // List ignored players

GET    /api/tell-queue     // List queued tells
POST   /api/tell-queue     // Send tell (if recipient offline)
```

---

### Phase 2: SSH Commands (herbert/)

#### Command Handlers

Create new command files:

- `herbert/cmd_say.go`
- `herbert/cmd_yell.go`
- `herbert/cmd_shout.go`
- `herbert/cmd_tell.go`
- `herbert/cmd_reply.go`
- `herbert/cmd_whisper.go`
- `herbert/cmd_emote.go`
- `herbert/cmd_social.go` (for /smile, /bow, etc.)
- `herbert/cmd_channel.go`
- `herbert/cmd_ignore.go`

#### Chat Service Integration

```go
// herbert/chat/chat.go
type ChatService struct {
    clientID  int
    roomID    int
    transport Transport
}

type Transport interface {
    Send(msg string)
    BroadcastRoom(msg string)
    BroadcastZone(msg string)
    BroadcastWorld(msg string)
}
```

---

### Phase 3: Admin UI

#### Social Command Editor

Add social command CRUD to admin panel:

```tsx
// admin/src/routes/_auth/socials.tsx
// Table: name, self_text, room_text, target_text
// Create/Edit form with all variants
```

#### Channel Configuration

Add channel management to admin:

```tsx
// admin/src/routes/_auth/channels.tsx
// Table: channel name, default enabled, color, cooldown
// Toggle default, set colors
```

#### Chat Log Viewer

Filter existing `/logs` page by `service=chat`:

```tsx
// Already exist in log viewer
// Add filter: service=chat
```

---

## Acceptance Criteria

### Core Commands (SSH Client)
- [ ] `say <msg>` — visible to room
- [ ] `yell <msg>` — visible to zone (5s cooldown)
- [ ] `shout <msg>` — visible to world (30s cooldown, 5 stamina cost)
- [ ] `tell <player> <msg>` — direct message with offline queue
- [ ] `reply <msg>` — respond to last teller
- [ ] `whisper <player> <msg>` — room-private
- [ ] `emote <action>` — freeform action
- [ ] `/<social>` — pre-built socials (smile, laugh, bow, wave, etc.)

### Channel System
- [ ] `chat <msg>` — general channel (2s cooldown)
- [ ] `newbie <msg>` — newbie help (2s cooldown)
- [ ] `trade <msg>` — trading (10s cooldown)
- [ ] `channels` command — list all channels
- [ ] `channel <name> on/off` — toggle channels

### Ignore System
- [ ] `ignore <player>` — block messages
- [ ] `unignore <player>` — unblock player
- [ ] Ignored players' messages not visible

### Backend
- [ ] `character_channel` table created
- [ ] `character_ignore` table created
- [ ] `tell_queue` table created
- [ ] `social_command` table created
- [ ] ChatService with all methods
- [ ] REST endpoints for all operations
- [ ] All messages logged to applogs with `service=chat`

### Admin UI
- [ ] Social command editor (CRUD)
- [ ] Channel management page
- [ ] Chat filter in `/logs` (service=chat)

### Player Experience
- [ ] All commands work in SSH client
- [ ] Color-coded channel output
- [ ] Offline tells queued and delivered
- [ ] Cooldowns enforced
- [ ] Rate limiting (30 msg/min global)

---

## Dependencies

- RFC-002 (send_message effect) — NPC scripting uses same channels
- Existing `applogs` pipeline
- Character entity with `name`, `race`, `gender` fields

---

## Open Questions

1. **Pronoun substitution**: `{he}`, `{him}`, `{his}` based on `gender` field — implement MVP or defer?
2. **Social command database storage**: Store in DB or hardcode initial 15-20?
3. **Channel history/scrollback**: Ring buffer for each channel — MVP or Phase 2?
4. **WebSocket transport**: Defer for web client (RFC-001)?

---

## Related RFCs

- RFC-002: Effects & Hooks System (NPC messaging uses same channels)
- RFC-009: Full messaging system specification
