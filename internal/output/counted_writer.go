package output

import (
	"encoding/json"

	"github.com/your-org/logslice/internal/stats"
)

// CountedWriter wraps a Writer and records per-entry statistics via a Sink.
type CountedWriter struct {
	w    *Writer
	sink *stats.Sink
}

// NewCountedWriter creates a CountedWriter that delegates writes to w and
// records observations in sink.
func NewCountedWriter(w *Writer, sink *stats.Sink) *CountedWriter {
	return &CountedWriter{w: w, sink: sink}
}

// Write formats and writes a single log entry, incrementing matched on
// success or errored on failure.
func (cw *CountedWriter) Write(entry map[string]json.RawMessage) error {
	cw.sink.IncTotal()
	if err := cw.w.Write(entry); err != nil {
		cw.sink.IncErrored()
		return err
	}
	cw.sink.IncMatched()
	return nil
}

// WriteAll writes a stream of entries from ch, recording stats for each.
// It returns the first write error encountered, if any.
func (cw *CountedWriter) WriteAll(ch <-chan map[string]json.RawMessage) error {
	for entry := range ch {
		if err := cw.Write(entry); err != nil {
			return err
		}
	}
	return nil
}
