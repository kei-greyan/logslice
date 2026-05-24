package main

import (
	"flag"
	"fmt"
	"strings"
)

type config struct {
	files  []string
	filter string
	since  string
	until  string
	format string
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var (
		filter = fs.String("filter", "", "filter expression, e.g. 'level=error'")
		since  = fs.String("since", "", "include entries at or after this time (RFC3339)")
		until  = fs.String("until", "", "include entries at or before this time (RFC3339)")
		format = fs.String("format", "pretty", "output format: json, pretty, text")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	validFormats := map[string]bool{"json": true, "pretty": true, "text": true}
	if !validFormats[strings.ToLower(*format)] {
		return nil, fmt.Errorf("unknown format %q: must be json, pretty, or text", *format)
	}

	return &config{
		files:  fs.Args(),
		filter: *filter,
		since:  *since,
		until:  *until,
		format: strings.ToLower(*format),
	}, nil
}
