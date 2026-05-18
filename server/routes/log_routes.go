package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/applog"
	"herbst-server/middleware"
)

// logBroadcaster fans out log entries to SSE subscribers.
type logBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[chan string]struct{}
	latestLine  string
}

var broadcaster = &logBroadcaster{
	subscribers: make(map[chan string]struct{}),
}

// subscribe registers a channel to receive log entries.
func (b *logBroadcaster) subscribe(ch chan string) {
	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()
}

// unsubscribe removes a channel.
func (b *logBroadcaster) unsubscribe(ch chan string) {
	b.mu.Lock()
	delete(b.subscribers, ch)
	b.mu.Unlock()
}

// broadcast sends a log line to all subscribers.
func (b *logBroadcaster) broadcast(line string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.latestLine = line
	for ch := range b.subscribers {
		select {
		case ch <- line:
		default:
			// subscriber too slow — skip
		}
	}
}

// BroadcastLogLine publishes a log line to SSE subscribers.
func BroadcastLogLine(level, service, message string, ts time.Time, characterID *int, roomID *int, templateID string, worldID string, metadata map[string]interface{}) {
	line := map[string]interface{}{
		"level":      level,
		"service":    service,
		"message":    message,
		"created_at": ts.Format(time.RFC3339),
	}
	if characterID != nil {
		line["character_id"] = *characterID
	}
	if roomID != nil {
		line["room_id"] = *roomID
	}
	if worldID != "" {
		line["world_id"] = worldID
	}
	if metadata != nil && len(metadata) > 0 {
		line["metadata"] = metadata
	}
	data, _ := json.Marshal(line)
	broadcaster.broadcast(string(data))
}

// RegisterLogRoutes registers log query + SSE routes under the protected group.
// The SSE stream endpoint is registered on the public router with query-param
// token auth since EventSource cannot send custom headers.
func RegisterLogRoutes(router *gin.Engine, protected *gin.RouterGroup, client *db.Client) {
	// SSE stream with query-param auth (EventSource can't send Bearer headers)
	router.GET("/api/logs/stream", streamLogsWithTokenAuth(client))
	// GET /api/logs — query with pagination and filters
	protected.GET("/logs", func(c *gin.Context) {
		level := c.Query("level")       // DEBUG, INFO, WARN, ERROR
		service := c.Query("service")   // filter by service name
		charID := c.Query("character_id")
		roomID := c.Query("room_id")
		templateID := c.Query("template_id")
	worldID := c.Query("world_id")
		limitStr := c.DefaultQuery("limit", "100")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		if limit < 1 || limit > 1000 {
			limit = 100
		}
		offset, _ := strconv.Atoi(offsetStr)
		if offset < 0 {
			offset = 0
		}

		// Build query conditions using generated where predicates.
		var conditions []func(*db.AppLogQuery) *db.AppLogQuery

		conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
			return q.Order(db.Desc(applog.FieldCreatedAt))
		})

		if level != "" {
			conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
				return q.Where(applog.LevelEQ(level))
			})
		}
		if service != "" {
			conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
				return q.Where(applog.ServiceEQ(service))
			})
		}
		if charID != "" {
			if id, err := strconv.Atoi(charID); err == nil {
				conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
					return q.Where(applog.CharacterIDEQ(id))
				})
			}
		}
		if roomID != "" {
			if id, err := strconv.Atoi(roomID); err == nil {
				conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
					return q.Where(applog.RoomIDEQ(id))
				})
			}
		}
	if worldID != "" {
		conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
			return q.Where(applog.WorldIDEQ(worldID))
		})
	}
		if templateID != "" {
			conditions = append(conditions, func(q *db.AppLogQuery) *db.AppLogQuery {
				return q.Where(applog.TemplateIDEQ(templateID))
			})
		}

		query := client.AppLog.Query()
		for _, cond := range conditions {
			query = cond(query)
		}

		// Count total for pagination
		countQ := client.AppLog.Query()
		for _, cond := range conditions {
			countQ = cond(countQ)
		}
		total, err := countQ.Count(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count logs"})
			return
		}

		entries, err := query.Limit(limit).Offset(offset).All(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query logs"})
			return
		}

		// Convert to JSON-friendly output
		type logOut struct {
			ID          int                    `json:"id"`
			Level       string                 `json:"level"`
			Message     string                 `json:"message"`
			Service     string                 `json:"service,omitempty"`
			CharacterID *int                   `json:"character_id,omitempty"`
			RoomID      *int                   `json:"room_id,omitempty"`
			TemplateID  string                 `json:"template_id,omitempty"`
			WorldID     string                 `json:"world_id,omitempty"`
			Metadata    map[string]interface{} `json:"metadata,omitempty"`
			CreatedAt   time.Time              `json:"created_at"`
		}

		result := make([]logOut, 0, len(entries))
		for _, e := range entries {
			lo := logOut{
				ID:        e.ID,
				Level:     e.Level,
				Message:   e.Message,
				CreatedAt: e.CreatedAt,
			}
			if e.Service != "" {
				lo.Service = e.Service
			}
			if e.CharacterID != nil {
				lo.CharacterID = e.CharacterID
			}
			if e.RoomID != nil {
				lo.RoomID = e.RoomID
			}
			if e.WorldID != "" {
				lo.WorldID = e.WorldID
			}
			if e.TemplateID != "" {
				lo.TemplateID = e.TemplateID
			}
			if e.Metadata != nil && len(e.Metadata) > 0 {
				lo.Metadata = e.Metadata
			}
			result = append(result, lo)
		}

		c.JSON(http.StatusOK, gin.H{
			"logs":   result,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
	})

	// GET /api/logs/services — distinct service names
	protected.GET("/logs/services", func(c *gin.Context) {
		entries, err := client.AppLog.Query().
			Select(applog.FieldService).
			Where(applog.ServiceNotNil()).
			All(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list services"})
			return
		}

		seen := make(map[string]bool)
		services := make([]string, 0)
		for _, e := range entries {
			if e.Service != "" && !seen[e.Service] {
				seen[e.Service] = true
				services = append(services, e.Service)
			}
		}
		c.JSON(http.StatusOK, gin.H{"services": services})
	})

		// GET /api/logs/worlds -- distinct world IDs
		protected.GET("/logs/worlds", func(c *gin.Context) {
			entries, err := client.AppLog.Query().
				Select(applog.FieldWorldID).
				Where(applog.WorldIDNotNil()).
				All(context.Background())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list worlds"})
				return
			}

			seen := make(map[string]bool)
			worlds := make([]string, 0)
			for _, e := range entries {
				if e.WorldID != "" && !seen[e.WorldID] {
					seen[e.WorldID] = true
					worlds = append(worlds, e.WorldID)
				}
			}
			c.JSON(http.StatusOK, gin.H{"worlds": worlds})
		})
}

// streamLogsWithTokenAuth validates a token from the query string and serves SSE.
// EventSource cannot send custom headers, so we accept ?token=<jwt> instead.
func streamLogsWithTokenAuth(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token parameter"})
			return
		}

		// Validate the JWT token using the same middleware logic
		userID, isAdmin, err := middleware.ValidateToken(token)
		if err != nil || !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		_ = userID

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.WriteHeader(http.StatusOK)

		ch := make(chan string, 50)
		broadcaster.subscribe(ch)
		defer broadcaster.unsubscribe(ch)

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			fmt.Fprint(c.Writer, "event: error\ndata: SSE not supported\n\n")
			return
		}

		ctx := c.Request.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(c.Writer, "data: %s\n\n", line)
				flusher.Flush()
			}
		}
	}
}
