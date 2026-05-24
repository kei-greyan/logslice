package output

import (
	"io"
	"sync"
)

// Writer wraps a Formatter and writes formatted log entries to an io.Writer.
// It is safe for concurrent use.
type Writer struct {
	mu        sync.Mutex
	w         io.Writer
	formatter *Formatter
}

// NewWriter creates a Writer that formats entries using f and writes to w.
func NewWriter(w io.Writer, f *Formatter) *Writer {
	return &Writer{w: w, formatter: f}
}

// Write formats entry and writes it to the underlying writer.
// Returns an error if formatting or writing fails.
func (wr *Writer) Write(entry map[string]any) error {
	line, err := wr.formatter.Format(entry)
	if err != nil {
		return err
	}
	wr.mu.Lock()
	defer wr.mu.Unlock()
	_, err = io.WriteString(wr.w, line+"\n")
	return err
}

// WriteAll writes all entries in order, stopping on first error.
func (wr *Writer) WriteAll(entries []map[string]any) error {
	for _, e := range entries {
		if err := wr.Write(e); err != nil {
			return err
		}
	}
	return nil
}
