package output

import (
	"github.com/user/logslice/internal/query"
)

// RunPipeline reads entries from ch, applies q, and writes matching entries
// to w. It returns the number of entries written and any write error.
// RunPipeline blocks until ch is closed.
func RunPipeline(ch <-chan map[string]any, q *query.CompiledQuery, w *Writer) (int, error) {
	var count int
	for entry := range ch {
		if !q.Match(entry) {
			continue
		}
		if err := w.Write(entry); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// FanIn merges multiple entry channels into a single channel.
// The returned channel is closed once all input channels are drained.
func FanIn(sources ...<-chan map[string]any) <-chan map[string]any {
	out := make(chan map[string]any, len(sources)*8)
	if len(sources) == 0 {
		close(out)
		return out
	}

	var wg sync.WaitGroup
	for _, src := range sources {
		wg.Add(1)
		go func(c <-chan map[string]any) {
			defer wg.Done()
			for e := range c {
				out <- e
			}
		}(src)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
