package query

// Pipeline applies a CompiledQuery to a channel of log entries, forwarding
// only those that match. It closes the output channel when the input is
// exhausted.
func Pipeline(cq *CompiledQuery, in <-chan map[string]any) <-chan map[string]any {
	out := make(chan map[string]any, 64)
	go func() {
		defer close(out)
		for entry := range in {
			if cq.Match(entry) {
				out <- entry
			}
		}
	}()
	return out
}

// PipelineAll chains multiple compiled queries (AND semantics): an entry must
// satisfy every query to pass through.
func PipelineAll(queries []*CompiledQuery, in <-chan map[string]any) <-chan map[string]any {
	out := make(chan map[string]any, 64)
	go func() {
		defer close(out)
		for entry := range in {
			pass := true
			for _, cq := range queries {
				if !cq.Match(entry) {
					pass = false
					break
				}
			}
			if pass {
				out <- entry
			}
		}
	}()
	return out
}
