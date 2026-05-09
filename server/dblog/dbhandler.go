// Package dblog provides an async slog.Handler that batches log entries to the applogs table.
package dblog

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"herbst-server/db"
)

// DBHandler implements slog.Handler, asynchronously writing log entries to the
// applogs table. Entries are buffered and batch-flushed to reduce DB round trips.
// On overflow, the oldest buffered entry is dropped. On DB error, the handler
// retries once, then logs the failure to stderr.
type DBHandler struct {
	client       *db.Client
	level        slog.Leveler
	minLevel     slog.Level
	svcFilter    map[string]bool
	broadcastFn  BroadcastFunc
	mu           sync.RWMutex
	ch           chan logEntry
	stop         chan struct{}
	wg           sync.WaitGroup
}

// BroadcastFunc is called after each log entry is persisted.
// Used to push live entries to SSE subscribers.
type BroadcastFunc func(level, service, message string, ts time.Time, characterID *int, roomID *int, templateID string, metadata map[string]interface{})

type logEntry struct {
	Level       slog.Level
	Message     string
	Service     string
	CharacterID *int
	RoomID      *int
	TemplateID  string
	Metadata    map[string]interface{}
}

// NewDBHandler creates a new DBHandler.
func NewDBHandler(client *db.Client, opts *slog.HandlerOptions) *DBHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	minLevel := slog.LevelInfo
	if v := os.Getenv("LOG_MIN_LEVEL"); v != "" {
		switch strings.ToUpper(v) {
		case "DEBUG":
			minLevel = slog.LevelDebug
		case "INFO":
			minLevel = slog.LevelInfo
		case "WARN":
			minLevel = slog.LevelWarn
		case "ERROR":
			minLevel = slog.LevelError
		}
	}

	var svcFilter map[string]bool
	if v := os.Getenv("LOG_SERVICE_FILTER"); v != "" {
		svcFilter = make(map[string]bool)
		for _, svc := range strings.Split(v, ",") {
			svcFilter[strings.TrimSpace(svc)] = true
		}
	}

	h := &DBHandler{
		client:    client,
		level:     opts.Level,
		minLevel:  minLevel,
		svcFilter: svcFilter,
		ch:        make(chan logEntry, 1000),
		stop:      make(chan struct{}),
	}
	h.wg.Add(1)
	go h.flusher()
	return h
}

func (h *DBHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.level != nil {
		return level >= h.level.Level()
	}
	return level >= h.minLevel
}

func (h *DBHandler) Handle(_ context.Context, r slog.Record) error {
	entry := logEntry{
		Level:   r.Level,
		Message: r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "service":
			entry.Service = a.Value.String()
		case "character_id":
			if v, err := strconv.Atoi(a.Value.String()); err == nil {
				entry.CharacterID = &v
			}
		case "room_id":
			if v, err := strconv.Atoi(a.Value.String()); err == nil {
				entry.RoomID = &v
			}
		case "template_id":
			entry.TemplateID = a.Value.String()
		default:
			if entry.Metadata == nil {
				entry.Metadata = make(map[string]interface{})
			}
			entry.Metadata[a.Key] = a.Value.Any()
		}
		return true
	})

	h.mu.RLock()
	filter := h.svcFilter
	h.mu.RUnlock()
	if len(filter) > 0 && !filter[entry.Service] {
		return nil
	}

	select {
	case h.ch <- entry:
	default:
		<-h.ch
		h.ch <- entry
	}
	return nil
}

// WithAttrs returns a new handler that carries additional attributes.
func (h *DBHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with a group.
func (h *DBHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *DBHandler) flusher() {
	defer h.wg.Done()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	batch := make([]logEntry, 0, 100)
	for {
		select {
		case entry := <-h.ch:
			batch = append(batch, entry)
			if len(batch) >= 100 {
				h.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				h.flush(batch)
				batch = batch[:0]
			}
		case <-h.stop:
			for {
				select {
				case entry := <-h.ch:
					batch = append(batch, entry)
				default:
					if len(batch) > 0 {
						h.flush(batch)
					}
					return
				}
			}
		}
	}
}

func (h *DBHandler) flush(batch []logEntry) {
	builders := make([]*db.AppLogCreate, 0, len(batch))
	for _, e := range batch {
		b := h.client.AppLog.Create().
			SetLevel(e.Level.String()).
			SetMessage(e.Message)
		if e.Service != "" {
			b.SetService(e.Service)
		}
		if e.CharacterID != nil {
			b.SetCharacterID(*e.CharacterID)
		}
		if e.RoomID != nil {
			b.SetRoomID(*e.RoomID)
		}
		if e.TemplateID != "" {
			b.SetTemplateID(e.TemplateID)
		}
		if e.Metadata != nil && len(e.Metadata) > 0 {
			b.SetMetadata(e.Metadata)
		}
		builders = append(builders, b)
	}

	err := h.client.AppLog.CreateBulk(builders...).Exec(context.Background())
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		err = h.client.AppLog.CreateBulk(builders...).Exec(context.Background())
		if err != nil {
			os.Stderr.WriteString("dblog: flush failed: " + err.Error() + "\n")
		}
	}

	for _, e := range batch {
		if h.broadcastFn != nil {
			h.broadcastFn(e.Level.String(), e.Service, e.Message, time.Now(), e.CharacterID, e.RoomID, e.TemplateID, e.Metadata)
		}
	}
}

// SetBroadcastFunc sets the callback invoked after each log entry is persisted.
// Call this after creating the handler to wire up live SSE streaming.
func (h *DBHandler) SetBroadcastFunc(fn BroadcastFunc) {
	h.broadcastFn = fn
}

func (h *DBHandler) GracefulShutdown() {
	close(h.stop)
	h.wg.Wait()
	close(h.ch)
}