package debuglog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TokenFetcher returns a JWT for authenticating to the server API.
type TokenFetcher func() (string, error)

// Logger batches debug log entries and POSTs them to the server's
// /api/debug-log endpoint. Each entry is tagged with character_id and
// room_id so the admin /logs page can filter by character.
type Logger struct {
	baseURL      string
	fetchToken   TokenFetcher
	client       *http.Client
	mu           sync.Mutex
	buffer       []entry
	done         chan struct{}
}

type entry struct {
	CharacterID int    `json:"character_id"`
	RoomID     int    `json:"room_id,omitempty"`
	Message     string `json:"message"`
}

// New creates a Logger that POSTs to baseURL + "/api/debug-log".
// fetchToken is called before each flush to get a valid JWT.
// Entries are batched and flushed every 500ms to avoid blocking game ticks.
func New(baseURL string, fetchToken TokenFetcher) *Logger {
	l := &Logger{
		baseURL:    baseURL,
		fetchToken: fetchToken,
		client:     &http.Client{Timeout: 3 * time.Second},
		done:       make(chan struct{}),
	}
	go l.flushLoop()
	return l
}

// Log enqueues a debug log entry. Non-blocking — drops entries if the
// buffer is full rather than stalling the game loop.
func (l *Logger) Log(characterID, roomID int, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.mu.Lock()
	if len(l.buffer) < 500 {
		l.buffer = append(l.buffer, entry{
			CharacterID: characterID,
			RoomID:       roomID,
			Message:      msg,
		})
	}
	l.mu.Unlock()
}

// Stop signals the flush goroutine to stop and flushes remaining entries.
func (l *Logger) Stop() {
	close(l.done)
	l.flush()
}

func (l *Logger) flushLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-l.done:
			return
		case <-ticker.C:
			l.flush()
		}
	}
}

func (l *Logger) flush() {
	l.mu.Lock()
	if len(l.buffer) == 0 {
		l.mu.Unlock()
		return
	}
	batch := l.buffer
	l.buffer = nil
	l.mu.Unlock()

	token, err := l.fetchToken()
	if err != nil {
		return
	}

	for _, e := range batch {
		body, _ := json.Marshal(e)
		req, err := http.NewRequest("POST", l.baseURL+"/api/debug-log", bytes.NewReader(body))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := l.client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}
}