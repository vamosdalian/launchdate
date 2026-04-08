package ll2datasyncer

import (
	"testing"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/util"
)

// Mock RateLimiter
type mockRateLimiter struct {
	ch chan struct{}
}

func (m *mockRateLimiter) Allow() bool           { return true }
func (m *mockRateLimiter) Wait() <-chan struct{} { return m.ch }
func (m *mockRateLimiter) Close()                {}

func TestBaseSyncer_Lifecycle(t *testing.T) {
	// Setup
	rl := &mockRateLimiter{ch: make(chan struct{})}
	workCh := make(chan struct{})

	workFunc := func() {
		workCh <- struct{}{}
	}

	syncer := NewBaseSyncer(rl, workFunc)

	// 1. Start
	syncer.Start()

	// 2. Initial work execution
	go func() { rl.ch <- struct{}{} }()
	select {
	case <-workCh:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Work function not called after start")
	}

	// 3. Pause
	syncer.Pause()

	// 4. Verify no work execution while paused
	go func() { rl.ch <- struct{}{} }()
	select {
	case <-workCh:
		t.Fatal("Work function called while paused")
	case <-time.After(100 * time.Millisecond):
		// Success: should timeout
	}

	// 5. Resume
	syncer.Resume()

	// 6. Verify work execution after resume (should pick up the pending signal from step 4)
	select {
	case <-workCh:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Work function not called after resume")
	}

	// 7. Cancel
	syncer.Cancel()

	select {
	case <-syncer.Done():
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Syncer not done after cancel")
	}
}

func TestBaseSyncer_CancelDuringPause(t *testing.T) {
	rl := &mockRateLimiter{ch: make(chan struct{})}
	syncer := NewBaseSyncer(rl, func() {})

	syncer.Start()
	syncer.Pause()

	time.Sleep(10 * time.Millisecond) // Give it a moment to potentially handle pause

	syncer.Cancel()

	select {
	case <-syncer.Done():
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Syncer not done after cancel while paused")
	}
}

// Ensure interface compatibility
func TestBaseSyncer_ImplementsInterface(t *testing.T) {
	var _ Syncer = &BaseSyncer{}
}

// Mock RateLimiter type for reference
var _ util.RateLimiter = &mockRateLimiter{}
