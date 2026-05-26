// Package ratelimit implements a token-bucket rate limiter for logslice.
//
// It is designed to cap the throughput of log entries flowing through the
// output pipeline, preventing downstream writers or terminals from being
// overwhelmed when processing high-volume log streams.
//
// Usage:
//
//	limiter := ratelimit.New(1000, 50) // 1000 entries/sec, burst of 50
//	filtered := ratelimit.Filter(entryCh, limiter, func() {
//		// optional: count or log dropped entries
//	})
//
// The Limiter is safe for concurrent use.
package ratelimit
