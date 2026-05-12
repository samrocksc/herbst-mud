# Ticket: RFC-002 - send_message Effect for Room Messaging

**Status:** Draft  
**Created:** 2026-05-12  
**RFC Reference:** RFC-002 (Effects & Hooks System)

---

## Executive Summary

Add a `send_message` effect type that allows NPCs, items, and abilities to broadcast messages to various channels (room, yell, shout, tell, whisper, chat, emote, etc.). This is a core building block for NPC scripting, item effects, and ability triggers.

---

## What It Is

A new effect type `send_message` with a channel dropdown that lets admins configure:

| Field | Options |
|-------|---------|
| **Message** | Text input (the content to send) |
| **Channel** | Dropdown: room, yell, shout, tell, whisper, chat, newbie, trade, ooc, admin, emote |
| **Target** | Optional, depends on channel type |

### Channel Behavior

| Channel | Target Required | Range | Example |
|---------|-----------------|-------|---------|
| `room` | No | Room only | "A merchant hawks his wares." |
| `yell` | No | Zone/area | "Guards! Intruders!" |
| `shout` | No | World | "Auction! Ancient Sword!" |
| `tell` | Yes | Direct to player/NPC | "Meet me at the fountain." |
| `whisper` | Yes | Room-private | "Don't trust the goblin." |
| `chat` | No | Channel | "Anyone seen the goblin shaman?" |
| `newbie` | No | Channel | "How do I equip a sword?" |
| `trade` | No | Channel | "Selling Ancient Key for 500 gold!" |
| `emote` | No | Room (freeform action) | "bows deeply before the elder." |

### Effect Parameters (JSON)

```json
{
  "message": "The town crier bellows: 'Hear ye! The king seeks adventurers!'",
  "channel": "room",
  "target": null
}
```

For `tell`:
```json
{
  "message": "Meet me at the fountain.",
  "channel": "tell",
  "target_type": "player",
  "target_id": 7
}
```

---

## What Needs to Be Done

### 1. Schema Changes (`server/db/schema/effect.go`)

Add `send_message` to valid effect types:

```go
field.String("effect_type").
    Comment("xp_drain|xp_gain|...|send_message|...")
```

### 2. Effect Service (`server/service/effect_service.go`)

Add effect handler for `send_message`:

```go
func (s *effectService) ApplySendMessageEffect(
    ctx context.Context, 
    charID int, 
    effect *db.Effect, 
    targetResolver EffectTargetResolver,
) error {
    params := effect.Parameters
    channel := params["channel"].(string)
    message := params["message"].(string)
    
    // Dispatch based on channel type
    switch channel {
    case "room", "yell", "shout", "emote":
        return s.broadcastToChannel(ctx, channel, message, charID)
    case "tell", "whisper":
        return s.sendDirectMessage(ctx, channel, message, charID, targetResolver)
    case "chat", "newbie", "trade", "ooc", "admin":
        return s.sendToChannel(ctx, channel, message, charID)
    }
    return nil
}
```

### 3. Chat Service (`herbst/chat/service.go`)

Create new chat service for SSH-based messaging:

```go
type ChatService struct {
    transport Transport
}

type Transport interface {
    SendRoom(charID int, message string) error
    SendYell(charID int, message string) error
    SendShout(charID int, message string) error
    SendTell(fromID int, toID int, message string) error
    SendWhisper(fromID int, toID int, message string) error
    SendChannel(channel string, message string) error
    SendEmote(charID int, action string) error
}
```

### 4. Command Handlers (`herbert/cmd_sendmessage.go`)

Add effect-based commands that wrap the chat service:

```go
func cmdSendMessage(c *cmdContext, channel, message string) {
    s := chat.NewChatService(c.charDB, c.transport)
    s.Send(c.charID, channel, message)
}
```

### 5. Admin UI (`admin/src/routes/_auth/effects.*`)

Add `send_message` to effect type dropdown with channel selector:

```tsx
<select value={params.channel} onChange={...}>
  <option value="room">Room</option>
  <option value="yell">Yell</option>
  <option value="shout">Shout</option>
  <option value="tell">Tell</option>
  <option value="whisper">Whisper</option>
  <option value="chat">Chat</option>
  <option value="newbie">Newbie</option>
  <option value="trade">Trade</option>
  <option value="emote">Emote</option>
</select>
```

### 6. Logging (applogs)

All messages logged with:

```json
{
  "level": "INFO",
  "service": "chat",
  "character_id": 7,
  "message": "[room] Gorthak: Hello there, traveler.",
  "metadata": {
    "command": "room",
    "room_id": 42,
    "channel": "room"
  }
}
```

---

## Acceptance Criteria

- [ ] `send_message` effect type added to `Effect` schema
- [ ] EffectService handles all channel types (room, yell, shout, tell, whisper, chat, emote)
- [ ] ChatService in `herbst/chat/` with Transport interface
- [ ] All channel implementations (SSH transport)
- [ ] Admin effect creation form includes channel dropdown
- [ ] Admin effect edit form includes channel dropdown
- [ ] Effects page in admin UI lists send_message effects
- [ ] All messages logged to applogs with `service=chat`

---

## Dependencies

- RFC-002 (Effects & Hooks System) - foundational
- Existing `Effect` and `ActiveEffect` schemas
- `applogs` pipeline for message logging

---

## Notes

- First pass: SSH-only transport (no WebSocket/Redis)
- Pronoun substitution (`{he}`, `{him}`, `{his}`) deferred to future
- Rate limiting deferred to future

---

## Related RFCs

- RFC-002: Effects & Hooks System
- RFC-009: Speak/Chat/Message System (full player commands)
- EFF-010: room_message effect type (referenced in RFC-009)
