package stats

import (
	"sync/atomic"
)

// Sink is a write-only view of a Counter, suitable for passing into
// pipeline stages that should only record observations.
type Sink struct {
	c *Counter
}

// NewSink wraps the given Counter as a Sink.
func NewSink(c *Counter) *Sink {
	return &Sink{c: c}
}

// IncTotal records that one log entry was read from a source.
func (s *Sink) IncTotal() {
	atomic.AddInt64(&s.c.total, 1)
}

// IncMatched records that one log entry passed all query filters.
func (s *Sink) IncMatched() {
	atomic.AddInt64(&s.c.matched, 1)
}

// IncSkipped records that one log entry was read but did not match.
func (s *Sink) IncSkipped() {
	atomic.AddInt64(&s.c.skipped, 1)
}

// IncErrored records that one line could not be parsed or processed.
func (s *Sink) IncErrored() {
	atomic.AddInt64(&s.c.errored, 1)
}
