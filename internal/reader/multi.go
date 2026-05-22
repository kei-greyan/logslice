package reader

import (
	"sync"
)

// MergeFiles reads multiple log files concurrently and merges their records
// into a single channel, preserving per-file source labels.
func MergeFiles(paths []string) (<-chan Record, <-chan error) {
	out := make(chan Record, 128)
	errs := make(chan error, len(paths))

	var wg sync.WaitGroup
	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			records, fileErrs := ReadFile(path)
			for rec := range records {
				out <- rec
			}
			if err := <-fileErrs; err != nil {
				errs <- err
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

// DrainErrors collects all errors from the error channel into a slice.
func DrainErrors(errs <-chan error) []error {
	var result []error
	for err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}
