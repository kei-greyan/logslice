package tail

// MultiTailer fans out entries from several Tailers into a single channel,
// tagging each entry with a "_source" field identifying the originating file.
type MultiTailer struct {
	tailers []*Tailer
	Out     chan map[string]any
	Err     chan error
}

// NewMulti creates a MultiTailer that tails all provided file paths.
// Call Start to begin receiving entries.
func NewMulti(paths []string) *MultiTailer {
	tailers := make([]*Tailer, len(paths))
	for i, p := range paths {
		tailers[i] = New(p)
	}
	return &MultiTailer{
		tailers: tailers,
		Out:     make(chan map[string]any, 128),
		Err:     make(chan error, 16),
	}
}

// Start begins tailing all files. Returns the first error encountered
// opening any file; already-started tailers are stopped on error.
func (m *MultiTailer) Start() error {
	for i, tlr := range m.tailers {
		if err := tlr.Start(); err != nil {
			for _, prev := range m.tailers[:i] {
				prev.Stop()
			}
			return err
		}
	}
	for _, tlr := range m.tailers {
		go m.forward(tlr)
	}
	return nil
}

// Stop signals all underlying tailers to shut down.
func (m *MultiTailer) Stop() {
	for _, tlr := range m.tailers {
		tlr.Stop()
	}
}

func (m *MultiTailer) forward(tlr *Tailer) {
	for {
		select {
		case entry, ok := <-tlr.Out:
			if !ok {
				return
			}
			if _, exists := entry["_source"]; !exists {
				entry["_source"] = tlr.path
			}
			m.Out <- entry
		case err, ok := <-tlr.Err:
			if !ok {
				return
			}
			select {
			case m.Err <- err:
			default:
			}
		}
	}
}
