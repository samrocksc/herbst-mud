package events

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// EventType is a string identifier for an event kind.
type EventType string

// Built-in event types.
const (
	EventNPCDefeated   EventType = "npc_defeated"
	EventCharacterDied EventType = "character_died"
	EventLevelUp       EventType = "level_up"
	EventQuestComplete EventType = "quest_complete"
	EventSkillLearned  EventType = "skill_learned"
)

// XPAwarder abstracts the XP award logic so subscribers don't need a concrete service import.
type XPAwarder interface {
	AwardXP(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error)
	ApplyDeathPenalty(ctx context.Context, characterID, penaltyPercent int) (xpLost, newXP int, err error)
}

// Event is the payload published to the bus.
type Event struct {
	Type      EventType             `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp int64                  `json:"timestamp"`
}

// Subscriber is a function that handles an event.
// Return an error to log it; the bus never stops a subscriber for errors.
type Subscriber func(Event) error

// Bus is the central event bus.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]Subscriber
	logger      *slog.Logger
}

// New creates a new event bus.
func New(logger *slog.Logger) *Bus {
	if logger == nil {
		logger = slog.Default()
	}
	return &Bus{
		subscribers: make(map[EventType][]Subscriber),
		logger:      logger,
	}
}

// Subscribe registers a handler for a specific event type.
// The handler runs asynchronously in its own goroutine per publish.
func (b *Bus) Subscribe(eventType EventType, fn Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], fn)
	b.logger.Info("subscriber registered", "event", eventType)
}

// Publish dispatches an event to all subscribers of that event type.
// Each subscriber runs in its own goroutine and is fire-and-forget.
func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	subs, ok := b.subscribers[event.Type]
	b.mu.RUnlock()

	if !ok || len(subs) == 0 {
		b.logger.Debug("no subscribers for event", "type", event.Type)
		return
	}

	for _, sub := range subs {
		go func(fn Subscriber) {
			if err := fn(event); err != nil {
				b.logger.Error("subscriber error",
					"event", event.Type,
					"error", err,
				)
			}
		}(sub)
	}
}

// Global bus instance — initialised once at startup.
var globalBus *Bus

// Init initialises the global bus. Call once from main.go.
func Init(logger *slog.Logger) {
	globalBus = New(logger)
}

// Default returns the global bus. Panics if not yet initialised.
func Default() *Bus {
	if globalBus == nil {
		panic("events.Default() called before events.Init()")
	}
	return globalBus
}

// Publish is a package-level convenience helper.
func Publish(event Event) {
	Default().Publish(event)
}

// Subscribe is a package-level convenience helper.
func Subscribe(eventType EventType, fn Subscriber) {
	Default().Subscribe(eventType, fn)
}

// Validate checks that the event has a non-empty type and non-nil payload.
func (e Event) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	if e.Payload == nil {
		return fmt.Errorf("event payload cannot be nil")
	}
	return nil
}
