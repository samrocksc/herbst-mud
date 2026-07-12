package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/dblog"
	"herbst-server/service"
)
// roomClient holds both service container and raw DB client
// for handlers still needing direct DB access during migration.
type roomClient struct {
	svc *service.Container
	db  *db.Client
}

// getRoomCharacters returns all characters (NPCs and players) in a room.
func getRoomCharacters(svc *service.Container) gin.HandlerFunc {
	// Uses service for room lookup, but direct DB for character query
	// until CharacterService.GetCharactersInRoom is implemented.
	rc := &roomClient{svc: svc}
	return rc.getCharacters
}

func (rc *roomClient) getCharacters(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Invalid room id", "error", err, slog.String("service", "rooms"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	characters, err := rc.db.Character.Query().
		Where(character.CurrentRoomId(id)).
		WithUser().All(c.Request.Context())
	if err != nil {
		dblog.Error("Failed to get room characters", err, slog.String("service", "rooms"), slog.Int("room_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type charView struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		IsNPC         bool   `json:"isNPC"`
		Level         int    `json:"level"`
		Class         string `json:"class"`
		Race          string `json:"race"`
		Hp            int    `json:"hp"`
		MaxHp         int    `json:"maxHp"`
		NpcTemplateID string `json:"npcTemplateId,omitempty"`
		XpValue       int    `json:"xpValue,omitempty"`
		LastSeenAt    string `json:"lastSeenAt,omitempty"`
	}
	result := make([]charView, 0, len(characters))

	userID, _ := c.Get("user_id")
	var currentUserID int
	if uid, ok := userID.(uint); ok {
		currentUserID = int(uid)
	}

	for _, ch := range characters {
		if !ch.IsNPC {
			if (currentUserID != 0 && ch.Edges.User != nil && ch.Edges.User.ID == currentUserID) ||
				ch.LastSeenAt == nil || time.Since(*ch.LastSeenAt) > 2*time.Minute {
				continue
			}
		}
		cv := charView{
			ID: ch.ID, Name: ch.Name, IsNPC: ch.IsNPC,
			Level: ch.Level, Class: ch.Class, Race: ch.Race,
			Hp: ch.Hitpoints, MaxHp: ch.MaxHitpoints,
			NpcTemplateID: ch.NpcTemplateID,
		}
		if ch.LastSeenAt != nil {
			cv.LastSeenAt = ch.LastSeenAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if ch.IsNPC && ch.NpcTemplateID != "" {
			tmpl, err := rc.db.NPCTemplate.Get(c.Request.Context(), ch.NpcTemplateID)
			if err == nil && tmpl.XpValue > 0 {
				cv.XpValue = tmpl.XpValue
			} else {
				cv.XpValue = ch.Level * 10
			}
		}
		result = append(result, cv)
	}
	c.JSON(http.StatusOK, result)
}

// getRoomLook returns a composite room view with characters and items.
func getRoomLook(svc *service.Container) gin.HandlerFunc {
	rc := &roomClient{svc: svc}
	return rc.getLook
}

func (rc *roomClient) getLook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Invalid room id for look", "error", err, slog.String("service", "rooms"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	room, err := rc.svc.Room.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	characters, _ := rc.db.Character.Query().
		Where(character.CurrentRoomId(id)).
		WithUser().All(c.Request.Context())
	npcs, players := partitionCharacters(c, characters)
	equipments, _ := rc.db.Equipment.Query().All(c.Request.Context())
	items := filterVisibleItems(equipments, id)
	c.JSON(http.StatusOK, gin.H{
		"id":             room.ID,
		"name":           room.Name,
		"description":    room.Description,
		"isStartingRoom": room.IsStartingRoom,
		"exits":          room.Exits,
		"items":          items,
		"npcs":           npcs,
		"players":        players,
		"z_level":        0,
	})
}
