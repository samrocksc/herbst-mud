// Package combat implements the combat system for Herbst MUD.
// It provides tick-based combat with action queues and effect management.
package combat

import (
	"context"
	"log"
	"sync"
	"time"
)

// TickDuration is the default time between combat ticks (1.5 seconds)
const TickDuration = 1500 * time.Millisecond

// Tick represents a single combat tick event
type Tick struct {
	ID        int
	Timestamp time.Time
	Combats   []int // IDs of active combats this tick
}

// TickManager handles the global tick loop for all combat encounters
type TickManager struct {
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	ticker      *time.Ticker
	currentTick int
	running     bool

	// Subscribers to tick events
	subscribers map[int]chan Tick // combatID -> tick channel

	// Callback when a tick occurs
	onTick func(tick Tick)
}

// NewTickManager creates a new tick manager with the default tick duration
func NewTickManager() *TickManager {
	return &TickManager{
		ticker:      time.NewTicker(TickDuration),
		subscribers: make(map[int]chan Tick),
	}
}

// NewTickManagerWithDuration creates a tick manager with a custom duration
func NewTickManagerWithDuration(duration time.Duration) *TickManager {
	return &TickManager{
		ticker:      time.NewTicker(duration),
		subscribers: make(map[int]chan Tick),
	}
}

// Start begins the tick loop
func (tm *TickManager) Start(ctx context.Context) {
	tm.mu.Lock()
	tm.ctx, tm.cancel = context.WithCancel(ctx)
	tm.running = true
	tm.mu.Unlock()

	go func() {
		for {
			select {
			case <-tm.ctx.Done():
				tm.mu.Lock()
				tm.running = false
				tm.mu.Unlock()
				return
			case <-tm.ticker.C:
				tm.processTick()
			}
		}
	}()

	log.Println("[TickManager] Started with duration:", TickDuration)
}

// Stop halts the tick loop
func (tm *TickManager) Stop() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.cancel != nil {
		tm.cancel()
	}
	if tm.ticker != nil {
		tm.ticker.Stop()
	}
	tm.running = false

	// Close all subscriber channels
	for combatID, ch := range tm.subscribers {
		close(ch)
		delete(tm.subscribers, combatID)
	}

	log.Println("[TickManager] Stopped")
}

// Subscribe registers a combat to receive tick events
// Returns a channel that will receive Tick events
func (tm *TickManager) Subscribe(combatID int) <-chan Tick {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	ch := make(chan Tick, 10) // Buffer for non-blocking sends
	tm.subscribers[combatID] = ch

	log.Printf("[TickManager] Combat %d subscribed to ticks", combatID)
	return ch
}

// Unsubscribe removes a combat from receiving tick events
func (tm *TickManager) Unsubscribe(combatID int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if ch, exists := tm.subscribers[combatID]; exists {
		close(ch)
		delete(tm.subscribers, combatID)
		log.Printf("[TickManager] Combat %d unsubscribed from ticks", combatID)
	}
}

// SetOnTick sets a callback function to be called on each tick
func (tm *TickManager) SetOnTick(callback func(tick Tick)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.onTick = callback
}

// processTick handles a single tick event
func (tm *TickManager) processTick() {
	tm.mu.Lock()
	tm.currentTick++
	tick := Tick{
		ID:        tm.currentTick,
		Timestamp: time.Now(),
		Combats:   make([]int, 0, len(tm.subscribers)),
	}

	// Collect active combat IDs
	for combatID := range tm.subscribers {
		tick.Combats = append(tick.Combats, combatID)
	}
	tm.mu.Unlock()

	// Call the onTick callback if set
	if tm.onTick != nil {
		tm.onTick(tick)
	}

	// Notify all subscribers (non-blocking)
	tm.mu.RLock()
	for combatID, ch := range tm.subscribers {
		select {
		case ch <- tick:
			// Sent successfully
		default:
			// Channel full, skip this tick for this combat
			log.Printf("[TickManager] Warning: combat %d tick channel full, skipping tick %d", combatID, tick.ID)
		}
	}
	tm.mu.RUnlock()
}

// GetCurrentTick returns the current tick number
func (tm *TickManager) GetCurrentTick() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.currentTick
}

// IsRunning returns whether the tick manager is actively running
func (tm *TickManager) IsRunning() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.running
}

// GetDuration returns the current tick duration
func (tm *TickManager) GetDuration() time.Duration {
	return TickDuration
}