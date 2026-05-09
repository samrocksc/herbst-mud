package questservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Service manages quest definitions and progress via the REST API.
type Service struct {
	cache    map[int]QuestDef
	mu       sync.RWMutex
	restBase string
	client   *http.Client
	logger   *slog.Logger
}

// NewService creates a new quest service.
func NewService(restBase string, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		cache:    make(map[int]QuestDef),
		restBase: restBase,
		client:   &http.Client{Timeout: 10 * time.Second},
		logger:   logger,
	}
}

// RefreshCache loads all quest definitions from the REST API.
func (s *Service) RefreshCache(ctx context.Context) error {
	var quests []QuestDef
	if err := s.getJSON(ctx, "/api/quests", &quests); err != nil {
		return fmt.Errorf("refresh quests: %w", err)
	}
	s.mu.Lock()
	newCache := make(map[int]QuestDef, len(quests))
	for _, q := range quests {
		newCache[q.ID] = q
	}
	s.cache = newCache
	s.mu.Unlock()
	s.logger.Info("quest cache refreshed", "quests", len(quests))
	return nil
}

// GetQuest returns a cached quest definition by ID.
func (s *Service) GetQuest(id int) (QuestDef, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.cache[id]
	return q, ok
}

// GetCharacterQuests fetches a character's active quest progress via REST.
func (s *Service) GetCharacterQuests(charID int) ([]QuestProgress, error) {
	url := fmt.Sprintf("%s/api/characters/%d/quests", s.restBase, charID)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %d", url, resp.StatusCode)
	}
	var result struct {
		Quests []QuestProgress `json:"quests"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Quests, nil
}

// StartRefreshLoop starts a background goroutine that refreshes the cache.
func (s *Service) StartRefreshLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.RefreshCache(context.Background()); err != nil {
				s.logger.Error("quest cache refresh failed", "error", err)
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
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}