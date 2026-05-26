package ratelimit

import (
	"testing"
	"time"
)

func TestNew_AllowsBurst(t *testing.T) {
	l := New(1, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() == true for token %d", i)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() == false after burst exhausted")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	now := time.Now()
	l := New(10, 1) // 10 tokens/sec, burst of 1
	l.clock = func() time.Time { return now }

	// exhaust the single token
	if !l.Allow() {
		t.Fatal("expected first Allow() to succeed")
	}
	if l.Allow() {
		t.Fatal("expected second Allow() to fail")
	}

	// advance time by 200ms → +2 tokens refilled (rate=10/s)
	now = now.Add(200 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected Allow() to succeed after refill")
	}
}

func TestAllow_ZeroRate(t *testing.T) {
	l := New(0, 2)
	if !l.Allow() {
		t.Fatal("expected burst token")
	}
	if !l.Allow() {
		t.Fatal("expected second burst token")
	}
	if l.Allow() {
		t.Fatal("expected no more tokens with zero refill rate")
	}
}

func TestReset_RestoresBurst(t *testing.T) {
	l := New(0, 2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("should be empty before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected token after reset")
	}
}

func TestAllow_TokensCapAtMax(t *testing.T) {
	now := time.Now()
	l := New(100, 2) // burst=2, rate=100/s
	l.clock = func() time.Time { return now }

	// advance by 10 seconds — would give 1000 tokens but max is 2
	now = now.Add(10 * time.Second)
	count := 0
	for l.Allow() {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 tokens (capped at burst), got %d", count)
	}
}
