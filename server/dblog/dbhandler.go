// Package dblog provides an async slog.Handler that batches log entries to the applogs table.
package dblog

import (
	"context"
	"log/slog"
	"os"
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
	client    *db.Client
	level     slog.Leveler
	minLevel  slog.Level
	svcFilter map[string]bool // empty = allow all; populated = allow only listed services
	mu        sync.RWMutex
	ch        chan logEntry
	stop      chan struct{}
	wg        sync.WaitGroup
}

type logEntry struct {
	Level   slog.Level
	Message string
	Service string
}

// NewDBHandler creates a new DBHandler.
//
// The client argument is the Ent client (must be connected). Options are read
// from environment variables:
//
//	LOG_MIN_LEVEL        — minimum level to persist (default: INFO)
//	LOG_SERVICE_FILTER   — comma-separated list of services to allow (default: all)
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

// Enabled reports whether this handler is enabled for the given level.
func (h *DBHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.level != nil {
		return level >= h.level.Level()
	}
	return level >= h.minLevel
}

// Handle enqueues a log record. If the channel is full, the entry is dropped
// (oldest-out semantics are maintained by the buffered channel's natural overflow).
func (h *DBHandler) Handle(_ context.Context, r slog.Record) error {
	// Extract the "service" attr if present
	svc := ""
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "service" {
			svc = a.Value.String()
			return false
		}
		return true
	})

	// Apply service filter
	h.mu.RLock()
	filter := h.svcFilter
	h.mu.RUnlock()
	if len(filter) > 0 && !filter[svc] {
		return nil
	}

	entry := logEntry{
		Level:   r.Level,
		Message: r.Message,
		Service: svc,
	}
	select {
	case h.ch <- entry:
	default:
		// channel full — drop oldest by reading one then pushing
		<-h.ch
		h.ch <- entry
	}
	return nil
}

// WithAttrs returns a new handler with additional attributes.
func (h *DBHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with a group.
func (h *DBHandler) WithGroup(name string) slog.Handler {
	return h
}

// flusher is the background goroutine that batch-inserts log entries.
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
			// Drain remaining entries then stop
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

// flush writes the batch to the database. Retries once on failure.
func (h *DBHandler) flush(batch []logEntry) {
	builders := make([]*db.AppLogCreate, 0, len(batch))
	for _, e := range batch {
		b := h.client.AppLog.Create().
			SetLevel(e.Level.String()).
			SetMessage(e.Message)
		if e.Service != "" {
			b.SetService(e.Service)
		}
		builders = append(builders, b)
	}

	err := h.client.AppLog.CreateBulk(builders...).Exec(context.Background())
	if err != nil {
		// Retry once after a brief pause
		time.Sleep(100 * time.Millisecond)
		err = h.client.AppLog.CreateBulk(builders...).Exec(context.Background())
		if err != nil {
			// Log to stderr to avoid infinite loop (don't call slog here)
			os.Stderr.WriteString("dblog: flush failed: " + err.Error() + "\n")
		}
	}
}

// GracefulShutdown drains pending entries and stops the flusher. Call before
// program exit to ensure no log entries are lost.
func (h *DBHandler) GracefulShutdown() {
	close(h.stop)
	h.wg.Wait()
	close(h.ch)
}
