package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls how log entries are rendered.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
	FormatText   Format = "text"
)

// Formatter writes log entries to an io.Writer in the chosen format.
type Formatter struct {
	Out    io.Writer
	Format Format
}

// NewFormatter creates a Formatter with the given output and format.
func NewFormatter(out io.Writer, format Format) *Formatter {
	return &Formatter{Out: out, Format: format}
}

// Write renders a single log entry.
func (f *Formatter) Write(entry map[string]any) error {
	switch f.Format {
	case FormatPretty:
		return f.writePretty(entry)
	case FormatText:
		return f.writeText(entry)
	default:
		return f.writeJSON(entry)
	}
}

func (f *Formatter) writeJSON(entry map[string]any) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.Out, string(b))
	return err
}

func (f *Formatter) writePretty(entry map[string]any) error {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.Out, string(b))
	return err
}

func (f *Formatter) writeText(entry map[string]any) error {
	// Print well-known fields first, then the rest alphabetically.
	priority := []string{"time", "level", "msg", "message"}
	seen := map[string]bool{}
	var parts []string

	for _, k := range priority {
		if v, ok := entry[k]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			seen[k] = true
		}
	}

	keys := make([]string, 0, len(entry))
	for k := range entry {
		if !seen[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, entry[k]))
	}

	_, err := fmt.Fprintln(f.Out, strings.Join(parts, " "))
	return err
}
