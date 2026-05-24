package main

import (
	"fmt"
	"os"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/query"
	"github.com/user/logslice/internal/reader"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	q, err := query.Parse(cfg.filter, cfg.since, cfg.until)
	if err != nil {
		return fmt.Errorf("invalid query: %w", err)
	}

	fmt, err := output.NewFormatter(cfg.format, os.Stdout)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	var entries <-chan reader.Entry
	var errs <-chan error

	if len(cfg.files) == 0 {
		entries, errs = reader.ReadReader(os.Stdin)
	} else {
		entries, errs = reader.MergeFiles(cfg.files)
		go reader.DrainErrors(errs, os.Stderr)
	}

	for entry := range entries {
		if q.Match(entry) {
			if err := fmt.Write(entry); err != nil {
				return err
			}
		}
	}

	return nil
}
