// Package ratelimit provides a token-bucket rate limiter for log output.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which log entries are emitted.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to burst entries immediately and
// refills at ratePerSec entries per second.
func New(ratePerSec float64, burst int) *Limiter {
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     ratePerSec,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Allow returns true if an entry may be emitted right now.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the limiter to a full burst bucket.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.lastTick = l.clock()
}
