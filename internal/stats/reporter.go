package stats

import (
	"fmt"
	"io"
)

// Reporter writes a human-readable summary to an io.Writer.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{out: w}
}

// Print formats and writes the summary to the reporter's writer.
func (r *Reporter) Print(s Summary) {
	fmt.Fprintf(r.out, "--- logslice summary ---\n")
	fmt.Fprintf(r.out, "  total   : %d\n", s.Total)
	fmt.Fprintf(r.out, "  matched : %d\n", s.Matched)
	fmt.Fprintf(r.out, "  skipped : %d\n", s.Skipped)
	fmt.Fprintf(r.out, "  errored : %d\n", s.Errored)
	fmt.Fprintf(r.out, "  elapsed : %s\n", s.Elapsed.Round(1*millisecond))
}

const millisecond = 1_000_000 // time.Millisecond without importing time
