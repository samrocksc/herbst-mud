package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// ─── Message protocol ─────────────────────────────────────────────────────────

// ClientMessage is parsed from the WebSocket client.
type ClientMessage struct {
	Type    string          `json:"type"`    // "command" | "heartbeat" | "subscribe"
	Payload string          `json:"payload"` // raw command text when type="command"
	Data    json.RawMessage `json:"data,omitempty"`
}

// ServerMessage is sent to the WebSocket client.
type ServerMessage struct {
	Type      string      `json:"type"`      // "output" | "system" | "error" | "ping" | "screen"
	Text      string      `json:"text"`      // human-readable content
	Data      interface{} `json:"data,omitempty"` // structured data (e.g. screen payload)
	Timestamp int64       `json:"timestamp"` // Unix ms
}

const (
	MsgOutput = "output"
	MsgSystem = "system"
	MsgError  = "error"
	MsgPing   = "ping"
	MsgScreen = "screen"
)

// ─── Screen payload types ──────────────────────────────────────────────────────

// CharInfo represents a visible character in a room.
type CharInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`    // "npc" | "player"
	ID      int    `json:"id"`
	Hostile bool   `json:"hostile"`
}

// ItemInfo represents a visible item in a room.
type ItemInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Takeable   bool   `json:"takeable"`
	Examinable bool   `json:"examinable"`
}

// RoomExit represents an exit from the room.
type RoomExit struct {
	Direction string `json:"direction"`
	Target    int    `json:"target"`
	Label     string `json:"label"`
}

// RoomScreenPayload is the structured room data sent to the client.
type RoomScreenPayload struct {
	ViewType    string     `json:"view_type"`
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Exits       []RoomExit `json:"exits"`
	Characters  []CharInfo `json:"characters"`
	Items       []ItemInfo `json:"items"`
}

// ─── Connection manager ───────────────────────────────────────────────────────

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // CORS handled by Gin
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	connMu sync.RWMutex
	// connections maps user_id → active websocket connection
	connections = make(map[uint]*WSConn)
)

// WSConn wraps an active WebSocket with metadata.
type WSConn struct {
	Conn        *websocket.Conn
	UserID      uint
	CharacterID int
	World       string
	Send        chan ServerMessage
	done        chan struct{}
}

// ─── Helper functions ─────────────────────────────────────────────────────────

func sendScreen(wsc *WSConn, payload RoomScreenPayload) {
	wsc.send(ServerMessage{
		Type:      MsgScreen,
		Text:      payload.Title,
		Data:      payload,
		Timestamp: time.Now().UnixMilli(),
	})
}

func buildRoomScreen(ctx context.Context, roomID int, worldID string, repos *repository.Container) (RoomScreenPayload, error) {
	rm, err := repos.Room.Get(ctx, roomID)
	if err != nil {
		return RoomScreenPayload{}, err
	}

	// Build exits from room.Exits map (direction -> targetID)
	var exits []RoomExit
	for dir, targetID := range rm.Exits {
		label := fmt.Sprintf("Exit %s", dir)
		// Try to get the target room name for a nicer label
		if tgt, err := repos.Room.Get(ctx, targetID); err == nil {
			label = tgt.Name
		}
		exits = append(exits, RoomExit{
			Direction: dir,
			Target:    targetID,
			Label:     label,
		})
	}

	// Characters in room (NPCs + other players)
	var chars []CharInfo
	rmChars, err := repos.Character.ListByRoom(ctx, roomID)
	if err != nil {
		slog.Error("buildRoomScreen: failed to list characters", "error", err, "room_id", roomID)
	} else {
		for _, ch := range rmChars {
			chType := "player"
			if ch.IsNPC {
				chType = "npc"
			}
			hostile := chType == "npc" // default
			if chType == "npc" && ch.NpcTemplateID != "" {
				tmpl, tmplErr := repos.NPCTemplate.Get(ctx, ch.NpcTemplateID)
				if tmplErr == nil && tmpl.Disposition == "friendly" {
					hostile = false
				}
			}
			chars = append(chars, CharInfo{
				Name:    ch.Name,
				Type:    chType,
				ID:      ch.ID,
				Hostile: hostile,
			})
		}
	}

	// Equipment in room
	var items []ItemInfo = []ItemInfo{}
	rmItems, err := repos.Equipment.ListByRoom(ctx, roomID)
	if err != nil {
		slog.Error("buildRoomScreen: failed to list equipment", "error", err, "room_id", roomID)
	} else {
		for _, it := range rmItems {
			items = append(items, ItemInfo{
				ID:         it.ID,
				Name:       it.Name,
				Takeable:   !it.IsImmovable,
				Examinable: it.IsVisible,
			})
		}
	}

	return RoomScreenPayload{
		ViewType:    "room",
		ID:          roomID,
		Title:       rm.Name,
		Description: rm.Description,
		Exits:       exits,
		Characters:  chars,
		Items:       items,
	}, nil
}

// ─── Handler ──────────────────────────────────────────────────────────────────

// RegisterWSRoutes registers the WebSocket upgrade endpoint.
func RegisterWSRoutes(router *gin.Engine, repos *repository.Container, client *db.Client) {
	router.GET("/ws", wsHandler(repos, client))
}

func wsHandler(repos *repository.Container, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT from query param (WebSocket constructor can't set headers)
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token query param required"})
			return
		}

		userID, _, err := middleware.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Character selection is required to enter the game world
		charIDStr := c.Query("character_id")
		if charIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "character_id query param required"})
			return
		}
		charID, err := strconv.Atoi(charIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character_id"})
			return
		}

		// Verify character ownership via ListByUser (safer than ent edge methods)
		userChars, err := repos.Character.ListByUser(c.Request.Context(), int(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load characters"})
			return
		}
		found := false
		for _, uc := range userChars {
			if uc.ID == charID {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusForbidden, gin.H{"error": "character does not belong to user"})
			return
		}

		char := userChars[0] // We need the actual char object for Name/CurrentWorld
		for _, uc := range userChars {
			if uc.ID == charID {
				char = uc
				break
			}
		}

		// Upgrade to WebSocket
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			slog.Error("ws upgrade failed", "error", err)
			return
		}
		defer ws.Close()

		wsc := &WSConn{
			Conn:        ws,
			UserID:      userID,
			CharacterID: charID,
			World:       char.CurrentWorld,
			Send:        make(chan ServerMessage, 64),
			done:        make(chan struct{}),
		}

		// Register connection
		connMu.Lock()
		if old, ok := connections[userID]; ok {
			// Signal old connection to stop without double-closing its done channel
			old.Conn.Close()
		}
		connections[userID] = wsc
		connMu.Unlock()

		// Cleanup on exit
		defer func() {
			connMu.Lock()
			if connections[userID] == wsc {
				delete(connections, userID)
			}
			connMu.Unlock()
			close(wsc.done)
		}()

		// Welcome + room info
		wsc.send(ServerMessage{
			Type:      MsgSystem,
			Text:      fmt.Sprintf("Welcome, %s! You are in world \"%s\".", char.Name, char.CurrentWorld),
			Timestamp: time.Now().UnixMilli(),
		})

		// Send structured room screen
		roomScreen, err := buildRoomScreen(c.Request.Context(), char.CurrentRoomId, char.CurrentWorld, repos)
		if err == nil {
			sendScreen(wsc, roomScreen)
		} else {
			slog.Warn("failed to load room screen", "error", err)
			wsc.send(ServerMessage{
				Type:      MsgSystem,
				Text:      fmt.Sprintf("\nRoom %d\n(Unable to load room details)", char.CurrentRoomId),
				Timestamp: time.Now().UnixMilli(),
			})
		}

		// Start goroutines
		go wsc.writePump()
		wsc.readPump(repos, client)
	}
}

// ─── Read pump ────────────────────────────────────────────────────────────────

func (wsc *WSConn) readPump(repos *repository.Container, client *db.Client) {
	wsc.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	wsc.Conn.SetPongHandler(func(string) error {
		wsc.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msgBytes, err := wsc.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("ws read error", "user", wsc.UserID, "error", err)
			}
			return
		}

		var msg ClientMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			wsc.send(ServerMessage{Type: MsgError, Text: "Invalid JSON", Timestamp: time.Now().UnixMilli()})
			continue
		}

		switch msg.Type {
		case "heartbeat":
			wsc.send(ServerMessage{Type: MsgPing, Text: "pong", Timestamp: time.Now().UnixMilli()})

		case "command":
			response := handleCommand(msg.Payload, wsc, repos, client)
			wsc.send(ServerMessage{Type: MsgOutput, Text: response, Timestamp: time.Now().UnixMilli()})

		default:
			wsc.send(ServerMessage{Type: MsgError, Text: fmt.Sprintf("Unknown message type: %s", msg.Type), Timestamp: time.Now().UnixMilli()})
		}
	}
}

// ─── Write pump ───────────────────────────────────────────────────────────────

func (wsc *WSConn) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-wsc.Send:
			if !ok {
				wsc.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			wsc.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := wsc.Conn.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			wsc.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := wsc.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-wsc.done:
			return
		}
	}
}

func (wsc *WSConn) send(msg ServerMessage) {
	select {
	case wsc.Send <- msg:
	case <-wsc.done:
	}
}

// ─── Command handler ──────────────────────────────────────────────────────────

func handleCommand(cmd string, wsc *WSConn, repos *repository.Container, client *db.Client) string {
	if cmd == "" {
		return "Type a command and press Enter."
	}

	// examine <target>
	if len(cmd) > 8 && strings.HasPrefix(cmd, "examine ") {
		target := cmd[8:]
		return fmt.Sprintf("You examine %s closely. It's nothing special. (Detailed object descriptions coming in Phase 6.)", target)
	}

	// talk to <target>
	if len(cmd) > 5 && strings.HasPrefix(cmd, "talk ") {
		target := cmd[5:]
		return tryTalk(target, wsc, repos)
	}

	// attack <target>
	if strings.HasPrefix(cmd, "attack ") {
		target := strings.TrimPrefix(cmd, "attack ")
		return tryAttack(target, wsc, repos, client)
	}

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "Type a command and press Enter."
	}

	// Aliases
	base := strings.ToLower(parts[0])

	// Direction aliases → canonical names (must match room.Exits keys)
	directionMap := map[string]string{
		"n": "north", "north": "north",
		"s": "south", "south": "south",
		"e": "east", "east": "east",
		"w": "west", "west": "west",
		"ne": "northeast", "northeast": "northeast",
		"se": "southeast", "southeast": "southeast",
		"sw": "southwest", "southwest": "southwest",
		"nw": "northwest", "northwest": "northwest",
		"u": "up", "up": "up",
		"d": "down", "down": "down",
	}
	if dir, ok := directionMap[base]; ok {
		return tryMove(dir, wsc, repos, client)
	}

	switch base {
	case "l":
		base = "look"
	}

	switch base {
	case "look":
		// Refresh room screen
		ctx := context.Background()
		char, err := repos.Character.Get(ctx, wsc.CharacterID)
		if err != nil {
			slog.Error("look command: failed to get character", "error", err, "character_id", wsc.CharacterID)
			return "You look around, but your surroundings refuse to come into focus."
		}
		roomScreen, err := buildRoomScreen(ctx, char.CurrentRoomId, char.CurrentWorld, repos)
		if err != nil {
			slog.Error("look command: failed to build room screen", "error", err, "room_id", char.CurrentRoomId)
		} else {
			sendScreen(wsc, roomScreen)
		}
		return "You look around."

	case "quit", "exit":
		return "Disconnecting is not yet implemented. Use the browser UI."

	case "who":
		return "You are alone in the void. (Player list coming soon.)"

	case "help":
		return "Available commands: look, who, help, quit, examine <target>, directions (n/s/e/w/u/d). (More coming in Phase 6.)"

	default:
		return fmt.Sprintf("You typed: \"%s\". (Command not yet implemented.)", cmd)
	}
}

func tryMove(dir string, wsc *WSConn, repos *repository.Container, client *db.Client) string {
	ctx := context.Background()

	char, err := repos.Character.Get(ctx, wsc.CharacterID)
	if err != nil {
		slog.Error("tryMove: failed to get character", "error", err, "character_id", wsc.CharacterID)
		return "Cannot find your character."
	}

	rm, err := repos.Room.Get(ctx, char.CurrentRoomId)
	if err != nil {
		slog.Error("tryMove: failed to get room", "error", err, "room_id", char.CurrentRoomId)
		return "You can't figure out where you are."
	}

	// Check if exit exists in room.Exits map
	targetID, ok := rm.Exits[dir]
	if !ok {
		return "You can't go that way."
	}

	// Move the character
	_, err = repos.Character.Update(ctx, char.ID, repository.CharacterUpdates{
		CurrentRoomID: &targetID,
	})
	if err != nil {
		slog.Error("tryMove: failed to update character room", "error", err, "character_id", char.ID, "target_room", targetID)
		return "Something prevents you from moving."
	}

	// Send new room screen
	roomScreen, err := buildRoomScreen(ctx, targetID, char.CurrentWorld, repos)
	if err != nil {
		slog.Error("tryMove: failed to build new room screen", "error", err, "room_id", targetID)
		return fmt.Sprintf("You move %s, but the new room refuses to resolve.", dir)
	}

	sendScreen(wsc, roomScreen)
	return fmt.Sprintf("You move %s to %s.", dir, roomScreen.Title)
}

// ─── talk ────────────────────────────────────────────────────────────────────

func tryTalk(targetName string, wsc *WSConn, repos *repository.Container) string {
	ctx := context.Background()
	char, err := repos.Character.Get(ctx, wsc.CharacterID)
	if err != nil {
		slog.Error("tryTalk: failed to get character", "error", err, "character_id", wsc.CharacterID)
		return "You try to talk, but something is wrong with your character."
	}

	// Find NPC by name in current room (case-insensitive contains match)
	roomChars, err := repos.Character.ListByRoom(ctx, char.CurrentRoomId)
	if err != nil {
		slog.Error("tryTalk: failed to list room characters", "error", err, "room_id", char.CurrentRoomId)
		return fmt.Sprintf("You try to talk to %s, but can't see anyone here.", targetName)
	}

	var targetNPC *db.Character
	for _, ch := range roomChars {
		if ch.IsNPC && strings.Contains(strings.ToLower(ch.Name), strings.ToLower(targetName)) {
			targetNPC = ch
			break
		}
	}
	if targetNPC == nil {
		return fmt.Sprintf("There is no %s here to talk to.", targetName)
	}

	// Fetch NPC template for greeting
	if targetNPC.NpcTemplateID == "" {
		return fmt.Sprintf("%s stares blankly and says nothing.", targetNPC.Name)
	}

	tmpl, err := repos.NPCTemplate.Get(ctx, targetNPC.NpcTemplateID)
	if err != nil {
		slog.Error("tryTalk: failed to get NPC template", "error", err, "template_id", targetNPC.NpcTemplateID)
		return fmt.Sprintf("%s seems unable to speak right now.", targetNPC.Name)
	}

	if tmpl.Greeting != "" {
		return fmt.Sprintf("%s says: \"%s\"", targetNPC.Name, tmpl.Greeting)
	}
	return fmt.Sprintf("%s greets you with a nod.", targetNPC.Name)
}

// ─── attack ────────────────────────────────────────────────────────────────────

func tryAttack(targetName string, wsc *WSConn, repos *repository.Container, client *db.Client) string {
	ctx := context.Background()
	char, err := repos.Character.Get(ctx, wsc.CharacterID)
	if err != nil {
		slog.Error("tryAttack: failed to get character", "error", err, "character_id", wsc.CharacterID)
		return "You try to attack, but something is wrong with your character."
	}

	// Find NPC by name in current room
	roomChars, err := repos.Character.ListByRoom(ctx, char.CurrentRoomId)
	if err != nil {
		slog.Error("tryAttack: failed to list room characters", "error", err, "room_id", char.CurrentRoomId)
		return fmt.Sprintf("You swing at empty air — %s is not here.", targetName)
	}

	var targetNPC *db.Character
	for _, ch := range roomChars {
		if ch.IsNPC && strings.Contains(strings.ToLower(ch.Name), strings.ToLower(targetName)) {
			targetNPC = ch
			break
		}
	}
	if targetNPC == nil {
		return fmt.Sprintf("There is no %s here to attack.", targetName)
	}

	// Check disposition via template
	hostile := true
	var tmpl *db.NPCTemplate
	if targetNPC.NpcTemplateID != "" {
		var tmplErr error
		tmpl, tmplErr = repos.NPCTemplate.Get(ctx, targetNPC.NpcTemplateID)
		if tmplErr == nil {
			hostile = tmpl.Disposition == "hostile"
		}
	}

	if !hostile {
		return fmt.Sprintf("%s is friendly. They recoil in surprise. (Type '/combat confirm' to attack neutral targets — Phase 6.)", targetNPC.Name)
	}

	// Simple damage roll (1d6 + str bonus)
	damage := 1 + (char.Strength / 5)
	if damage < 1 {
		damage = 1
	}

	// Apply damage respecting immortality
	newHP := targetNPC.Hitpoints - damage
	if newHP < 0 {
		newHP = 0
	}
	if targetNPC.IsImmortal && newHP < 1 {
		newHP = 1
	}
	_, err = repos.Character.Update(ctx, targetNPC.ID, repository.CharacterUpdates{
		Hitpoints: &newHP,
	})
	if err != nil {
		slog.Error("tryAttack: failed to apply damage", "error", err, "npc_id", targetNPC.ID, "damage", damage)
		return fmt.Sprintf("You swing your weapon, but your blow seems to pass through %s!", targetNPC.Name)
	}

	// Auto-heal immortal training dummy after 3s
	if targetNPC.IsImmortal && strings.Contains(strings.ToLower(targetNPC.Name), "dummy") {
		go func(npcID int, maxHP int) {
			time.Sleep(3 * time.Second)
			full := maxHP
			_, healErr := repos.Character.Update(context.Background(), npcID, repository.CharacterUpdates{Hitpoints: &full})
			if healErr != nil {
				slog.Error("training dummy auto-heal failed", "error", healErr, "npc_id", npcID)
			}
		}(targetNPC.ID, targetNPC.MaxHitpoints)
	}

	// Build result message
	var msg string
	if newHP == 0 && !targetNPC.IsImmortal {
		msg = fmt.Sprintf("You strike %s with a mighty blow for %d damage! They crumple to the ground.", targetNPC.Name, damage)
	} else {
		msg = fmt.Sprintf("You attack %s for %d damage! (%d/%d HP)", targetNPC.Name, damage, newHP, targetNPC.MaxHitpoints)
	}

	// Training dummy counter-attack: always 0 damage
	if strings.Contains(strings.ToLower(targetNPC.Name), "dummy") {
		msg += "\nThe dummy swings harmlessly at you."
	}

	return msg
}
