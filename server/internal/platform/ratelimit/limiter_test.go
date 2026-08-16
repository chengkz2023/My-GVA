package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterWindowBudget(t *testing.T) {
	l := NewLimiter(2, time.Minute)
	if !l.Allow("k") || !l.Allow("k") {
		t.Fatal("fresh key should be allowed")
	}
	l.Take("k")
	l.Take("k")
	if l.Allow("k") {
		t.Fatal("key at budget should be rejected")
	}
}

func TestLimiterReset(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	l.Take("k")
	if l.Allow("k") {
		t.Fatal("key at budget should be rejected")
	}
	l.Reset("k")
	if !l.Allow("k") {
		t.Fatal("reset key should be allowed")
	}
}

func TestLimiterWindowExpiry(t *testing.T) {
	l := NewLimiter(1, time.Millisecond)
	l.Take("k")
	time.Sleep(2 * time.Millisecond)
	if !l.Allow("k") {
		t.Fatal("expired window should be allowed")
	}
}

func TestLimiterIndependentKeys(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	l.Take("a")
	if !l.Allow("b") {
		t.Fatal("other key should be unaffected")
	}
}
