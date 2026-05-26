package ratelimit

import "github.com/user/logslice/internal/reader"

// Filter wraps an entry channel and forwards only entries permitted by the
// Limiter. Dropped entries are counted via the optional drop callback.
func Filter(
	in <-chan reader.Entry,
	l *Limiter,
	onDrop func(),
) <-chan reader.Entry {
	out := make(chan reader.Entry, cap(in))
	go func() {
		defer close(out)
		for entry := range in {
			if l.Allow() {
				out <- entry
			} else if onDrop != nil {
				onDrop()
			}
		}
	}()
	return out
}
