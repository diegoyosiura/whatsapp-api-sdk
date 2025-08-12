package httpx

import (
	"testing"
	"time"
)

// TestBackoffWithinBase ensures the backoff stays within the expected range
// when the calculated value is below the max cap.
func TestBackoffWithinBase(t *testing.T) {
	base := 100 * time.Millisecond
	maxTime := time.Second
	got := backoff(1, base, maxTime)
	if got < 100*time.Millisecond || got > 200*time.Millisecond {
		t.Fatalf("backoff out of range: %v", got)
	}
}

// TestBackoffCapped ensures the backoff respects the max cap when the
// exponential growth would exceed it.
func TestBackoffCapped(t *testing.T) {
	base := 100 * time.Millisecond
	maxTime := 150 * time.Millisecond // much lower than base*2^attempt
	got := backoff(5, base, maxTime)
	if got < maxTime/2 || got > maxTime {
		t.Fatalf("backoff not capped correctly: %v", got)
	}
}
