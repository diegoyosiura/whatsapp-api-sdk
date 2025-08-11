package httpx

import (
	"math/rand"
	"testing"
	"time"
)

// TestBackoffWithinBase ensures the backoff stays within the expected range
// when the calculated value is below the max cap.
func TestBackoffWithinBase(t *testing.T) {
	rand.Seed(1)
	base := 100 * time.Millisecond
	max := time.Second
	got := backoff(1, base, max)
	if got < 100*time.Millisecond || got > 200*time.Millisecond {
		t.Fatalf("backoff out of range: %v", got)
	}
}

// TestBackoffCapped ensures the backoff respects the max cap when the
// exponential growth would exceed it.
func TestBackoffCapped(t *testing.T) {
	rand.Seed(1)
	base := 100 * time.Millisecond
	max := 150 * time.Millisecond // much lower than base*2^attempt
	got := backoff(5, base, max)
	if got < max/2 || got > max {
		t.Fatalf("backoff not capped correctly: %v", got)
	}
}
