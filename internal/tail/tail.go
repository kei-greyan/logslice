// Package tail provides live-tailing of log files, emitting new lines
// as they are appended.
package tail

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// PollInterval is how often the tailer checks for new content.
const PollInterval = 200 * time.Millisecond

// Tailer reads new JSON lines appended to a file and sends them on Out.
type Tailer struct {
	path   string
	Out    chan map[string]any
	Err    chan error
	stop   chan struct{}
}

// New creates a Tailer for the given file path.
func New(path string) *Tailer {
	return &Tailer{
		path: path,
		Out:  make(chan map[string]any, 64),
		Err:  make(chan error, 8),
		stop: make(chan struct{}),
	}
}

// Start begins tailing the file in a background goroutine.
// It seeks to the end of the file before watching for new lines.
func (t *Tailer) Start() error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return err
	}
	go t.run(f)
	return nil
}

// Stop signals the tailer to shut down and closes the output channel.
func (t *Tailer) Stop() {
	close(t.stop)
}

func (t *Tailer) run(f *os.File) {
	defer f.Close()
	defer close(t.Out)

	dec := json.NewDecoder(f)
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			return
		case <-ticker.C:
			for dec.More() {
				var entry map[string]any
				if err := dec.Decode(&entry); err != nil {
					if err == io.EOF {
						break
					}
					select {
					case t.Err <- err:
					default:
					}
					continue
				}
				select {
				case t.Out <- entry:
				case <-t.stop:
					return
				}
			}
		}
	}
}
