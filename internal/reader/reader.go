package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Record represents a single parsed JSON log line with its source.
type Record struct {
	Source string
	Data   map[string]interface{}
	Raw    string
}

// ReadFile opens a file and streams JSON records line by line into the returned channel.
// The channel is closed when the file is fully read or an error occurs.
func ReadFile(path string) (<-chan Record, <-chan error) {
	records := make(chan Record, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(records)
		defer close(errs)

		f, err := os.Open(path)
		if err != nil {
			errs <- fmt.Errorf("reader: open %q: %w", path, err)
			return
		}
		defer f.Close()

		if err := scan(path, f, records); err != nil {
			errs <- err
		}
	}()

	return records, errs
}

// ReadReader streams JSON records from an io.Reader, labelling them with source.
func ReadReader(source string, r io.Reader) (<-chan Record, <-chan error) {
	records := make(chan Record, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(records)
		defer close(errs)

		if err := scan(source, r, records); err != nil {
			errs <- err
		}
	}()

	return records, errs
}

func scan(source string, r io.Reader, out chan<- Record) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			// Skip non-JSON lines silently.
			continue
		}
		out <- Record{Source: source, Data: data, Raw: line}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reader: scan %q: %w", source, err)
	}
	return nil
}
