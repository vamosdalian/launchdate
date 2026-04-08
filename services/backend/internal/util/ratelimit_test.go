package util

import (
	"testing"
	"time"
)

func TestRateLimit_Allow(t *testing.T) {
	// Test with a 100ms rate limit
	rate := 100 * time.Millisecond
	limiter := NewRateLimit(rate)
	defer limiter.Close()

	// 1. First call should be allowed immediately
	if !limiter.Allow() {
		t.Error("Expected first call to be allowed, but it was not")
	}

	// 2. Second call immediately after should be denied
	if limiter.Allow() {
		t.Error("Expected second call to be denied, but it was allowed")
	}

	// 3. Wait for the rate limit duration to pass
	time.Sleep(rate + 10*time.Millisecond) // Sleep a bit longer to avoid race conditions

	// 4. The next call should now be allowed
	if !limiter.Allow() {
		t.Error("Expected call after duration to be allowed, but it was not")
	}
}

func TestRateLimit_Wait(t *testing.T) {
	rate := 100 * time.Millisecond
	limiter := NewRateLimit(rate)
	defer limiter.Close()

	// The first Wait() should consume the token and return immediately.
	<-limiter.Wait()

	// The second Wait() should block for the duration of the rate limit.
	startTime := time.Now()
	<-limiter.Wait()
	elapsed := time.Since(startTime)

	// Check if the blocking time is approximately the rate limit duration.
	// We allow a small margin for the scheduler's delay.
	if elapsed < rate {
		t.Errorf("Expected Wait() to block for at least %v, but it blocked for %v", rate, elapsed)
	}

	if elapsed > rate*2 {
		t.Errorf("Expected Wait() to block for around %v, but it blocked for %v", rate, elapsed)
	}
}

func TestRateLimit_Close(t *testing.T) {
	rate := 50 * time.Millisecond
	limiter := NewRateLimit(rate)

	// Consume the initial token
	if !limiter.Allow() {
		t.Fatal("Could not get initial token")
	}

	// Close the limiter
	limiter.Close()

	// Wait for a period longer than the rate limit
	time.Sleep(rate + 20*time.Millisecond)

	// Since the limiter is closed, the token should not be refilled.
	// Allow() should return false.
	if limiter.Allow() {
		t.Error("Expected Allow() to return false after Close(), but it returned true")
	}
}
