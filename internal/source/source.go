package source

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Type represents the kind of log source.
type Type int

const (
	TypeFile Type = iota
	TypeStdin
)

// Source describes a named input stream.
type Source struct {
	Name   string
	Type   Type
	Reader io.ReadCloser
}

// FromArgs resolves a list of path arguments into Sources.
// A single "-" entry (or an empty list) means read from stdin.
func FromArgs(args []string) ([]*Source, error) {
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		return []*Source{
			{Name: "stdin", Type: TypeStdin, Reader: io.NopCloser(os.Stdin)},
		}, nil
	}

	sources := make([]*Source, 0, len(args))
	for _, path := range args {
		if path == "-" {
			sources = append(sources, &Source{
				Name:   "stdin",
				Type:   TypeStdin,
				Reader: io.NopCloser(os.Stdin),
			})
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("source: open %q: %w", path, err)
		}
		name := path
		if idx := strings.LastIndexAny(path, "/\\"); idx >= 0 {
			name = path[idx+1:]
		}
		sources = append(sources, &Source{
			Name:   name,
			Type:   TypeFile,
			Reader: f,
		})
	}
	return sources, nil
}

// CloseAll closes all sources, returning the first error encountered.
func CloseAll(sources []*Source) error {
	var first error
	for _, s := range sources {
		if err := s.Reader.Close(); err != nil && first == nil {
			first = fmt.Errorf("source: close %q: %w", s.Name, err)
		}
	}
	return first
}
