package combat

import (
	"fmt"
	"sync"
	"time"
)

// TickLoop manages the combat tick clock that drives all combat timing.
type TickLoop struct {
	mu       sync.RWMutex
	interval time.Duration
	combats  map[int]*Combat
	tickChan chan int
	stopChan chan struct{}
	stopped  bool
	tick     int
}

// NewTickLoop creates a new TickLoop with the specified interval in milliseconds.
func NewTickLoop(intervalMs int64) *TickLoop {
	return &TickLoop{
		interval: time.Duration(intervalMs) * time.Millisecond,
		combats:  make(map[int]*Combat),
		tickChan: make(chan int, 10),
		stopChan: make(chan struct{}),
		stopped:  false,
		tick:     0,
	}
}

// Start begins the tick loop. It runs in a goroutine and processes ticks at the configured interval.
func (t *TickLoop) Start() {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.ProcessTick()
		case <-t.stopChan:
			t.mu.Lock()
			t.stopped = true
			t.mu.Unlock()
			return
		}
	}
}

// ProcessTick processes a single tick cycle for all active combats.
func (t *TickLoop) ProcessTick() {
	t.mu.Lock()
	t.tick++
	currentTick := t.tick
	t.mu.Unlock()

	// Notify listeners of the new tick
	select {
	case t.tickChan <- currentTick:
	default:
		// Channel buffer full, skip notification
	}

	// Process all active combats
	t.mu.RLock()
	combats := make([]*Combat, 0, len(t.combats))
	for _, c := range t.combats {
		combats = append(combats, c)
	}
	t.mu.RUnlock()

	for _, c := range combats {
		c.ProcessTick(currentTick)
	}
}

// Stop halts the tick loop.
func (t *TickLoop) Stop() error {
	select {
	case t.stopChan <- struct{}{}:
		// Signal sent
	default:
		// Already stopped or stopping
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Wait for loop to acknowledge stop
	for !t.stopped {
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// RegisterCombat adds a combat to be tracked by the tick loop.
func (t *TickLoop) RegisterCombat(combat *Combat) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.combats[combat.ID] = combat
}

// UnregisterCombat removes a combat from the tick loop.
func (t *TickLoop) UnregisterCombat(combatID int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.combats, combatID)
}

// GetTick returns the current tick number.
func (t *TickLoop) GetTick() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tick
}

// GetInterval returns the current tick interval.
func (t *TickLoop) GetInterval() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.interval
}

// TickChan returns a channel that receives tick notifications.
func (t *TickLoop) TickChan() <-chan int {
	return t.tickChan
}

// String returns a string representation of the TickLoop.
func (t *TickLoop) String() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return fmt.Sprintf("TickLoop{tick: %d, interval: %v, combats: %d, stopped: %v}",
		t.tick, t.interval, len(t.combats), t.stopped)
}