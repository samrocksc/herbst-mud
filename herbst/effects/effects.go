package effects

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// EffectDef is a cached effect definition from the REST API.
type EffectDef struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	EffectType   string                 `json:"effect_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	StackMode    string                 `json:"stack_mode"`
	StackLimit   int                   `json:"stack_limit"`
	IsPermanent  bool                  `json:"is_permanent"`
	DurationSecs int                   `json:"duration_secs"`
	Messages     map[string]string     `json:"messages"`
}

// HookDef is a cached hook definition from the REST API.
type HookDef struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Event        string `json:"event"`
	Target       string `json:"target"`
	Condition    string `json:"condition"`
	Enabled      bool   `json:"enabled"`
	EffectID     int    `json:"effect_id"`
	NPCTemplateID string `json:"npc_template_id"`
}

// Service manages effect definitions and hooks in memory.
type Service struct {
	effectCache map[int]EffectDef
	hookCache   map[string][]HookDef // keyed by event name
	mu          sync.RWMutex
	restBase    string
	httpClient  *http.Client
	lastRefresh time.Time
	logger      *slog.Logger
	messageBus  *MessageBus
}

// NewService creates a new effects service.
func NewService(restBase string, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		effectCache: make(map[int]EffectDef),
		hookCache:   make(map[string][]HookDef),
		restBase:    restBase,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		logger:      logger,
		messageBus:  NewMessageBus(),
	}
}

// RefreshCache loads all effects and hooks from the REST API.
func (s *Service) RefreshCache(ctx context.Context) error {
	var effects []EffectDef
	if err := s.getJSON(ctx, "/api/effects", &effects); err != nil {
		return fmt.Errorf("refresh effects: %w", err)
	}
	var hooks []HookDef
	if err := s.getJSON(ctx, "/api/hooks", &hooks); err != nil {
		return fmt.Errorf("refresh hooks: %w", err)
	}

	s.mu.Lock()
	newEffectCache := make(map[int]EffectDef, len(effects))
	for _, e := range effects {
		newEffectCache[e.ID] = e
	}
	newHookCache := make(map[string][]HookDef)
	for _, h := range hooks {
		newHookCache[h.Event] = append(newHookCache[h.Event], h)
	}
	s.effectCache = newEffectCache
	s.hookCache = newHookCache
	s.lastRefresh = time.Now()
	s.mu.Unlock()

	s.logger.Info("effects cache refreshed", "effects", len(effects), "hooks", len(hooks))
	return nil
}

// GetEffect returns a cached effect by ID.
func (s *Service) GetEffect(id int) (EffectDef, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.effectCache[id]
	return e, ok
}

// GetHooksForEvent returns all enabled hooks for a given event.
func (s *Service) GetHooksForEvent(event string) []HookDef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hookCache[event]
}

// StartRefreshLoop starts a background goroutine that refreshes the cache periodically.
func (s *Service) StartRefreshLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.RefreshCache(context.Background()); err != nil {
				s.logger.Error("effects cache refresh failed", "error", err)
			}
		}
	}()
}

func (s *Service) getJSON(ctx context.Context, path string, target interface{}) error {
	url := s.restBase + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}