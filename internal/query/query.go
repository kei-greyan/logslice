package query

import (
	"fmt"
	"strings"
	"time"

	"logslice/internal/filter"
)

// Query holds parsed query parameters for filtering log entries.
type Query struct {
	Filter    *filter.Filter
	Since     *time.Time
	Until     *time.Time
	TimeField string
	Limit     int
}

// Options configures query parsing behavior.
type Options struct {
	FilterExpr string
	Since      string
	Until      string
	TimeField  string
	Limit      int
}

// Parse builds a Query from the provided Options.
func Parse(opts Options) (*Query, error) {
	q := &Query{
		TimeField: opts.TimeField,
		Limit:     opts.Limit,
	}

	if q.TimeField == "" {
		q.TimeField = "time"
	}

	if opts.FilterExpr != "" {
		f, err := filter.Parse(opts.FilterExpr)
		if err != nil {
			return nil, fmt.Errorf("invalid filter expression: %w", err)
		}
		q.Filter = f
	}

	if opts.Since != "" {
		t, err := parseTime(opts.Since)
		if err != nil {
			return nil, fmt.Errorf("invalid --since value: %w", err)
		}
		q.Since = &t
	}

	if opts.Until != "" {
		t, err := parseTime(opts.Until)
		if err != nil {
			return nil, fmt.Errorf("invalid --until value: %w", err)
		}
		q.Until = &t
	}

	return q, nil
}

var timeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as a timestamp", s)
}
