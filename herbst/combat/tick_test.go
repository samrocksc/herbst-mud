package combat

import (
	"testing"
)

func TestTickManager_StartStop(t *testing.T) {
	tm := NewTickManager()
	
	if tm.IsRunning() {
		t.Error("TickManager should not be running initially")
	}
	
	ctx := t.Context()
	tm.Start(ctx)
	
	if !tm.IsRunning() {
		t.Error("TickManager should be running after Start()")
	}
	
	tm.Stop()
	
	if tm.IsRunning() {
		t.Error("TickManager should not be running after Stop()")
	}
}

func TestTickManager_Subscribe(t *testing.T) {
	tm := NewTickManager()
	ctx := t.Context()
	tm.Start(ctx)
	defer tm.Stop()
	
	ch := tm.Subscribe(1)
	if ch == nil {
		t.Error("Subscribe should return a channel")
	}
	
	// Check subscriber is registered
	tm.mu.RLock()
	_, exists := tm.subscribers[1]
	tm.mu.RUnlock()
	
	if !exists {
		t.Error("Combat should be subscribed")
	}
}

func TestTickManager_Unsubscribe(t *testing.T) {
	tm := NewTickManager()
	ctx := t.Context()
	tm.Start(ctx)
	defer tm.Stop()
	
	_ = tm.Subscribe(1)
	tm.Unsubscribe(1)
	
	tm.mu.RLock()
	_, exists := tm.subscribers[1]
	tm.mu.RUnlock()
	
	if exists {
		t.Error("Combat should be unsubscribed")
	}
}

func TestTickManager_GetCurrentTick(t *testing.T) {
	tm := NewTickManager()
	
	// Initial tick should be 0
	if tick := tm.GetCurrentTick(); tick != 0 {
		t.Errorf("Expected initial tick to be 0, got %d", tick)
	}
}