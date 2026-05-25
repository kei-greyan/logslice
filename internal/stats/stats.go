package stats

import (
	"sync/atomic"
	"time"
)

// Counter tracks pipeline processing statistics.
type Counter struct {
	total    atomic.Int64
	matched  atomic.Int64
	skipped  atomic.Int64
	errored  atomic.Int64
	startedAt time.Time
}

// New creates a new Counter with the start time set to now.
func New() *Counter {
	return &Counter{startedAt: time.Now()}
}

// IncTotal records one more entry seen.
func (c *Counter) IncTotal() { c.total.Add(1) }

// IncMatched records one more entry that passed the query filter.
func (c *Counter) IncMatched() { c.matched.Add(1) }

// IncSkipped records one more entry that was filtered out.
func (c *Counter) IncSkipped() { c.skipped.Add(1) }

// IncErrored records one more entry that could not be parsed.
func (c *Counter) IncErrored() { c.errored.Add(1) }

// Summary returns a snapshot of current counts.
func (c *Counter) Summary() Summary {
	return Summary{
		Total:    c.total.Load(),
		Matched:  c.matched.Load(),
		Skipped:  c.skipped.Load(),
		Errored:  c.errored.Load(),
		Elapsed:  time.Since(c.startedAt),
	}
}

// Summary is a point-in-time snapshot of processing statistics.
type Summary struct {
	Total   int64
	Matched int64
	Skipped int64
	Errored int64
	Elapsed time.Duration
}
